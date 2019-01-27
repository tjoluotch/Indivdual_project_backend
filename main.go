package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Student struct {
	FirstName string `json:"first_name,omitempty"`
	LastName string	 `json:"last_name,omitempty"`
	Phone_No string	 `json:"phone_no,omitempty"`
	Student_ID string `json:"student_id,omitempty"`
	Unique_ID string `json:"unique_id,omitempty"`
}

type JwToken struct {
	Token string `json:"token,omitempty"`
}

type CustomsClaimsStudent struct {
	FirstName string `json:"first_name"`
	Student_ID string `json:"student_id"`
	Unique_ID string `json:"unique_id"`
	jwt.StandardClaims
}

// Global variables to be able to use in each endpoint
var client *mongo.Client

const(
	dbName string = "Studently"
	privKeyPath string = "keys/app.rsa"
	pubKeyPath string = "keys/app.rsa.pub"
)

var VerifyKey, SignKey []byte



func LogFileSetup() {
	// Setting format and parsing current time in this format
	currentTime := time.Now().Format(time.RFC1123)

	// create a new file if one doesn't exists and append data to this file when writing.
	file, err := os.OpenFile("info.log"+currentTime, os.O_RDWR |os.O_CREATE|os.O_APPEND, 0666)

	// if there's an error with the opening of the log file log Fatal
	if err != nil {
		log.Fatal(err)
	}


	// set output to the file
	log.SetOutput(file)
	// log with date time and file location
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

func initKeys() {
	var err error

	SignKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading Private Key")
		return
	}

	VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading Public Key")
		return
	}
}



// Get the form data entered by client; FirstName, LastName, phone Number,
// assign the person a unique i.d
// check to see if that user isn't in the database already
// if they are send an error message with the a  'bad' response code
// if they aren't in db add to db and send a message with success
func CreateStudentAccountEndpoint(response http.ResponseWriter, request *http.Request){

	//CORS
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

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


	response.Header().Set("Content-Type", "application/json")
	var student Student
	// decoding JSON post data to
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&student)
	//if there was an error panic
	if err != nil {
		panic(err)
	}
	// save unique id to student
	unId, _ := uuid.NewV4()
	idString := unId.String()
	student.Unique_ID = idString
	//print student to console
	fmt.Println(&student)
	//encode student into BSON
	data, err := bson.Marshal(student)
	if err != nil {
		log.Fatalf("Problem encoding Student struct into BSON: Err-> %v\n ",err)
	}
	studentCollection := client.Database(dbName).Collection("students")
	_, err = studentCollection.InsertOne(context.Background(),data)
	if err != nil {
		response.WriteHeader(501)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}
	// encoding json object for returning to the client
	jsonStudent, err := json.Marshal(student)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	response.Write(jsonStudent)
}

func CreateJwtokenEndpoint(response http.ResponseWriter, request *http.Request){
	//CORS
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// step 1 authenticate user

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

	// set response header
	response.Header().Set("Content-Type", "application/json")


	var student Student
	// decoding JSON post data to student
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&student)
	//if there was an error panic
	if err != nil {
		panic(err)
	}

	// Create search parameter of student id in bson format
	searchParams := "student_id"
	filter := bson.D{{searchParams, student.Student_ID}}

	studentCollection := client.Database(dbName).Collection("students")
	//find a particular student using student id
	result := studentCollection.FindOne(context.Background(), filter)
	var studentFound Student
	// decode that student into Student struct
	err = result.Decode(&studentFound)
	// user does not exist is return to client
	if err != nil {
		http.Error(response, "User does not exist ", http.StatusInternalServerError)
		return
	}

	// encoding json object for returning to the client
	//jsonStudent, err := json.Marshal(&studentFound)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	// step 2 create a jwt from student found
	claims := CustomsClaimsStudent{
		studentFound.FirstName,
		studentFound.Student_ID,
		studentFound.Unique_ID,
		jwt.StandardClaims{
			Issuer: "golang api",
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
	}
	mySigningKey := []byte("secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(response, "Issue with JWT creation: Error message -> " + err.Error(), http.StatusInternalServerError)
	}

	response.Header().Add("Authorization", "Bearer " + tokenString)
	json.NewEncoder(response).Encode(JwToken{Token: tokenString})
	//response.Write(jsonStudent)
}



func main() {

	LogFileSetup()
	initKeys()
	router := mux.NewRouter()


	router.HandleFunc("/api/signup", CreateStudentAccountEndpoint).Methods("POST")
	router.HandleFunc("/api/authenticate", CreateJwtokenEndpoint).Methods("POST")
	//log server running
	log.Printf("server running on port %v", 12345)
	log.Fatal(http.ListenAndServe(":12345",router))


}