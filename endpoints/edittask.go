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

func EditTaskEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	cors.EnableCORS(&response)
	fmt.Println("Edit Task")
	// Decode jwt claims into student Model
	decoded := context2.Get(request, "decoded")
	var student Student
	claims := decoded.(jwt.MapClaims)
	mapstructure.Decode(claims, &student)

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

	// Decode request body into data structure
	var moduleIDWithTask ModuleWithTaskForEditTaskEP
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&moduleIDWithTask)

	//if there was an error panic
	if err != nil {
		http.Error(response, "JSON failed to decode request body to Task with module object " + err.Error() , 400)
		return
	}

	fmt.Printf("Data structure %v\n", &moduleIDWithTask)

	// Encode Task struct
	var task Task
	task.Description = moduleIDWithTask.Description
	task.Status = moduleIDWithTask.Status
	task.Task_ID = moduleIDWithTask.Task_ID

	// Search parameters to find the task in the db to be updated


	// encode search parameters as bson
	filter := bson.D{{"tasks.task_id", task.Task_ID}}

	// update the task object within the module
	update := bson.M{"$set": bson.M{"tasks.$": task}}

	// avoid upsert or duplication


	_, err = moduleCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(response, "Problem Editing Task see edittask.go " + err.Error(), 400)
		return
	}

	// send success message
	response.WriteHeader(200)
	response.Write([]byte(`{ "message": "Successfully Edited task object" }`))
}
