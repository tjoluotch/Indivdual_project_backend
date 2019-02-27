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


// Find Module by module id and then add task to it
func AddTaskEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	cors.EnableCORS(&response)
	fmt.Println("Add Task")
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

	// generate unique id for the new Task to be added
	taskID, err := uuid.NewV4()
	if err != nil {
		http.Error(response, "Unique ID failed to generate for new Task object" , 400)
		return
	}
	taskIDString := taskID.String()


	// decoding JSON put data into Task with module struct
	var task Task
	var taskWithModuleID ModuleWithTask
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&taskWithModuleID)
	//if there was an error panic
	if err != nil {
		http.Error(response, "JSON failed to decode request body to Task with module object " + err.Error() , 400)
		return
	}
	// Encoded task struct fully
	task.Task_ID = taskIDString
	task.Status = taskWithModuleID.Status
	task.Description = taskWithModuleID.Description

	// Search parameters to find the module in db
	searchParams := "module_id"
	moduleUniqueID := &taskWithModuleID.Module_ID

	// encode search parameters as bson
	filter := bson.D{{searchParams,moduleUniqueID}}

	// update the module object to have array of tasks
	update := bson.M{"$push": bson.M{"tasks": task}}
	_, err = moduleCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(response, "Problem updating Module to include latest task " + err.Error(), 400)
		return
	}

	// send success message
	response.WriteHeader(200)
	response.Write([]byte(`{ "message": "Successfully added Task to Module object" }`))

	//fmt.Printf("the Task: Description %v, Task Status %v, The module ID: %v\n", taskWithModuleID.Description, taskWithModuleID.Status, taskWithModuleID.Module_ID)
}
