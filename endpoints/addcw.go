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

// Find Module by module id and then add cw object to it
func AddCourseworkEndpoint(response http.ResponseWriter, request *http.Request) {

	//CORS
	cors.EnableCORS(&response)
	fmt.Println("Add Coursework")
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
		return
	}
	moduleCollection := client.Database(dbName).Collection("modules")

	// generate unique id for the new Coursework object to be added
	cw_ID, err := uuid.NewV4()
	if err != nil {
		http.Error(response, "Unique ID failed to generate for new coursework object" , 400)
		return
	}
	cwIDString := cw_ID.String()

	// decoding JSON put data into Courswork with module struct
	var coursework Coursework
	var courseworkWithModID CourseWorkWithModule
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&courseworkWithModID)
	//if there was an error panic
	if err != nil {
		http.Error(response, "JSON failed to decode request body to Coursework with module object " + err.Error() , 400)
		return
	}

	// Encode coursework struct fully
	coursework.Coursework_ID = cwIDString
	coursework.Coursework_Description = courseworkWithModID.Coursework_Description
	coursework.Due_Date = courseworkWithModID.Due_Date

	// Search parameters to find the module in db
	searchParams := "module_id"
	moduleUniqueID := &courseworkWithModID.Module_ID

	// encode search parameters as bson
	filter := bson.D{{searchParams,moduleUniqueID}}

	// update to module object to have an array of Courseworks
	update := bson.M{"$push": bson.M{"cwks": coursework}}

	_, err = moduleCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(response, "Problem updating Module to include latest coursework " + err.Error(), 400)
		return
	}

	// send success message
	response.WriteHeader(200)
	response.Write([]byte(`{ "message": "Successfully added Coursework to Module object" }`))
}
