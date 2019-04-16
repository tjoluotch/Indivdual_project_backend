package endpoints

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dgrijalva/jwt-go"
	context2 "github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/satori/go.uuid"
	"log"
	"mygosource/ind_proj_backend/cors"
	"mygosource/ind_proj_backend/envar"
	"net/http"
	"os"
	"time"
)

// Upload student's file to AWS S3 Bucket
func UploadFileEndpoint(response http.ResponseWriter, request *http.Request) {

	//CORS
	cors.EnableCORS(&response)
	fmt.Println("Upload File endpoint")
	// Decode jwt claims into student Model
	decoded := context2.Get(request, "decoded")
	var student Student
	claims := decoded.(jwt.MapClaims)
	mapstructure.Decode(claims, &student)

	envar.Variables()

	awsAccessKey := os.Getenv("AWS_ACESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	token := ""

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, token)

	_, err := creds.Get()
	if err != nil {
		fmt.Println("Cannot Access s3 credentials")
		http.Error(response, "Cannot Access s3 Credentials " + err.Error() , 400)
		return
	}
	bucket := "studently-indivdualproj"
	cfg := aws.NewConfig().WithRegion("eu-west-2").WithCredentials(creds)
	session,err := session.NewSession(cfg)
	svc := s3.New(session, cfg)

	file, handler, err := request.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the file")
		http.Error(response, "Error Retrieving the file " + err.Error() , 400)
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Upload File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Save to S3 bucket
	path := "/docspace/"+ "/"+ student.Student_ID + "/" + handler.Filename
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(path),
		Body:file,
		ContentLength: aws.Int64(handler.Size),
		ContentType: aws.String(handler.Header.Get("Content-Type")),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Println("Error Uploading the file to aws s3 bucket")
		http.Error(response, "Error Uploading the file to aws s3 bucket" + err.Error() , 400)
		return
	}
	fmt.Printf("succesfully uploaded %q into %v\n", handler.Filename, bucket)
	fmt.Printf("response %s", awsutil.StringValue(resp))

	// Get Object URL
	paramsGet := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(path),
	}
	req, _:= svc.GetObjectRequest(paramsGet)

	docUrl, err := req.Presign(604799 * time.Second)
	if err != nil {
		fmt.Println("Error Getting URL of Object")
		http.Error(response, "Error getting URL of Object" + err.Error() , 400)
		return
	}

	fmt.Print("URL is: " + docUrl)

	// Open db and save document

	// Db opening section
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatalf("Error connecting to mongoDB client Host: Err-> %v\n ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error Connecting to MongoDB at context.WtihTimeout: Err-> %v\n ", err)
	}
	docSpaceCollection := client.Database(dbName).Collection("docspace")

	var doc Document
	// generate unique id and change to string
	docID, err := uuid.NewV4()
	if err != nil {
		http.Error(response, "Unique ID failed to generate for new Module Object" , 400)
		return
	}
	// Set doc struct values
	docIDString := docID.String()
	doc.Doc_ID = docIDString
	doc.Doc_Name = handler.Filename
	doc.Can_View = append(doc.Can_View, student.Student_ID)
	doc.Url = docUrl
	doc.AWS_Bucket = bucket
	doc.AWS_Key = path

	// encoding doc to bson
	data, err := bson.Marshal(doc)
	if err != nil {
		http.Error(response,"Problem encoding Module struct into BSON " + err.Error(), 400 )
		return
	}

	// adding doc to DB
	_, err = docSpaceCollection.InsertOne(context.Background(),data)
	if err != nil {
		response.WriteHeader(400)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}
	// send success message
	response.WriteHeader(200)
	response.Write([]byte(`{ "message": "Successfully added Document to DocSpace" }`))
}