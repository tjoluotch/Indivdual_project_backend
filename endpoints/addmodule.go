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
	"github.com/satori/go.uuid"
	"log"
	"mygosource/ind_proj_backend/cors"
	"net/http"
	"time"
)

// Use this when adding tasks to module

/*
func AddModuleEndpoint(response http.ResponseWriter, request *http.Request) {
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
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error Connecting to MongoDB at context.WtihTimeout: Err-> %v\n ", err)
	}
	studentCollection := client.Database(dbName).Collection("students")


	// generate unique id and change to string
	moduleID, err := uuid.NewV4()
	if err != nil {
		http.Error(response, "Unique ID failed to generate" , 400)
		return
	}
	moduleIDString := moduleID.String()

	var module Module
	// decoding JSON post data to module
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&module)
	//if there was an error panic
	if err != nil {
		http.Error(response, "JSON failed to decode to Module " + err.Error() , 400)
		return
	}
	// save unique id to module
	module.Module_ID = moduleIDString
	fmt.Println(&module)


	if err != nil {
		log.Fatalf("Problem encoding Student struct into BSON: Err-> %v\n ",err)
		return
	}

	// search parameters to find the student
	searchParams := "unique_id"
	studentUniqueID := &student.Unique_ID

	// find the student and add the module
	filter := bson.D{{searchParams, studentUniqueID}}


	if err != nil {
		http.Error(response, "Problem updating converting to BSON", 400)
		return
	}



	//value,_ := bson.Marshal(moduleSlice)


	// fixed
	// This is the way to add to the array
	update := bson.M{"$push": bson.M{"modules": module}}

	_, err = studentCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(response, "Problem updating student to include latest module " + err.Error(), 400)
		return
	}

	response.WriteHeader(200)
	response.Write([]byte("successfully added module"))
}
*/

func AddModuleEndpoint(response http.ResponseWriter, request *http.Request) {
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
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error Connecting to MongoDB at context.WtihTimeout: Err-> %v\n ", err)
	}
	moduleCollection := client.Database(dbName).Collection("modules")

	// generate unique id and change to string
	moduleID, err := uuid.NewV4()
	if err != nil {
		http.Error(response, "Unique ID failed to generate" , 400)
		return
	}
	moduleIDString := moduleID.String()

	var module Module
	// decoding JSON post data to module
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&module)
	//if there was an error panic
	if err != nil {
		http.Error(response, "JSON failed to decode to Module " + err.Error() , 400)
		return
	}
	// save unique id to module
	module.Module_ID = moduleIDString

	//save student id to module
	module.Student_ID = student.Student_ID

	//print module to console
	fmt.Println(&module)

	// encoding module to bson
	data, err := bson.Marshal(module)
	if err != nil {
		log.Fatalf("Problem encoding Student struct into BSON: Err-> %v\n ",err)
		return
	}

	// adding module to DB
	_, err = moduleCollection.InsertOne(context.Background(),data)
	if err != nil {
		response.WriteHeader(400)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}

	// send success message
	response.WriteHeader(200)
	response.Write([]byte(`{ "message": "Successfully added module" }`))
}