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

// Upload student's file to AWS S3 Bucket
func GetUploadsEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	cors.EnableCORS(&response)
	fmt.Println("Get Uploads endpoint")
	// Decode jwt claims into student Model
	decoded := context2.Get(request, "decoded")
	var student Student
	claims := decoded.(jwt.MapClaims)
	mapstructure.Decode(claims, &student)

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
	docCollection := client.Database(dbName).Collection("docspace")

	// Filter to find the docs the member can view
	documentFilter := bson.D{{"can_view", student.Student_ID}}

	// store results in Doc slice: initialising, so if empty, then an empty array is returned rather than nil.
		results := []Document{}

	// search db for all docs the student can view
	cursor, err := docCollection.Find(context.TODO(), documentFilter, nil)
	if err != nil {
		http.Error(response, "Issue searching DB to get Chat groups student is a part of: " + err.Error(), http.StatusBadRequest)
		return
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO()) {

		var element Document
		err = cursor.Decode(&element)
		if err != nil {
			http.Error(response, "Issue decoding one of the Chats to Struct " + err.Error(), http.StatusBadRequest)
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

	// Working - decode back to json
	json.NewEncoder(response).Encode(results)
}