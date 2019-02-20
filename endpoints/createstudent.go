package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/satori/go.uuid"
	"log"
	"mygosource/ind_proj_backend/cors"
	"net/http"
	"time"
)

// Get the form data entered by client; FirstName, LastName, phone Number,
// assign the person a unique i.d
// check to see if that user isn't in the database already
// if they are send an error message with the a  'bad' response code
// if they aren't in db add to db and send a message with success
func CreateStudentAccountEndpoint(response http.ResponseWriter, request *http.Request){


	//CORS
	cors.EnableCORS(&response)

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
		response.WriteHeader(400)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}
	// encoding json object for returning to the client
	jsonStudent, err := json.Marshal(student)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	response.Write(jsonStudent)
}
