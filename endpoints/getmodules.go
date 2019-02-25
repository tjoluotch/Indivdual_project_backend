package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	context2 "github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"mygosource/ind_proj_backend/cors"
	"net/http"
	"time"
)


// Get Modules for each student by using the student id
func GetModulesEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	cors.EnableCORS(&response)
	fmt.Println("Add module")
	// Decode jwt claims into student Model
	decoded := context2.Get(request, "decoded")
	var student Student
	claims := decoded.(jwt.MapClaims)
	mapstructure.Decode(claims, &student)

	// step 1 authenticate user

	// Db opening section
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatalf("Error connecting to mongoDB client Host: Err-> %v\n ", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error Connecting to MongoDB at context.WtihTimeout: Err-> %v\n ", err)
		return
	}
	moduleCollection := client.Database(dbName).Collection("modules")


	// Then add a filter to find the documents by.
	moduleFilter := bson.D{{"student_id", student.Student_ID}}

	// Slice to store decoded module documents
	var results []Module

	// search db for all modules belonging to a student
	cursor, err := moduleCollection.Find(context.TODO(), moduleFilter, nil)
	if err != nil {
		http.Error(response, "Issue searching DB to get Student's modules", http.StatusBadRequest)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO())  {

		// create a value into which the single document can be decoded
		var element Module
		err = cursor.Decode(&element)
		if err != nil {
			http.Error(response, "Issue decoding one of the Modules to Struct " + err.Error(), http.StatusBadRequest)
			return
		}
		results = append(results, element)
	}

	if err = cursor.Err(); err != nil {
		http.Error(response, "Issue with the Cursor " + err.Error(), http.StatusBadRequest)
		return
	}

	// Close the cursor once finished
	cursor.Close(context.TODO())

	fmt.Printf("Found multiple documents (Student modules): %+v\n", results)

	// Working - decode back to json
	json.NewEncoder(response).Encode(results)
}
