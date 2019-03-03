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



func CreateChatEndpoint(response http.ResponseWriter, request *http.Request) {

	//CORS
	cors.EnableCORS(&response)
	fmt.Println("Add Chat to chat space")
	// Decode jwt claims into student Model
	decoded := context2.Get(request, "decoded")
	var student Student
	claims := decoded.(jwt.MapClaims)
	mapstructure.Decode(claims, &student)

	// Get chat name -> add student as a member -> Generate chat id -> save to Db -> Send success message

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
	chatCollection := client.Database(dbName).Collection("chatspace")

	// generate unique id for chat and change to string
	chatID, err := uuid.NewV4()
	if err != nil {
		http.Error(response, "Unique ID failed to generate for new Chat Object" , 400)
		return
	}
	chatIDString := chatID.String()

	var chat Chat
	// decoding JSON put data to module
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&chat)
	//if there was an error panic
	if err != nil {
		http.Error(response, "JSON failed to decode request body to Chat object " + err.Error() , 400)
		return
	}
	// save chat id to chat object
	chat.Chat_ID = chatIDString
	// add Founding student(their student id) as chat member
	chat.Chat_Members = append(chat.Chat_Members, student.Student_ID)

	// encoding chat to bson
	data, err := bson.Marshal(chat)
	if err != nil {
		http.Error(response,"Problem encoding Chat struct into BSON " + err.Error(), 400 )
		return
	}

	// adding chat object to db
	_, err = chatCollection.InsertOne(context.Background(), data)
	if err != nil {
		response.WriteHeader(400)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}

	//fmt.Printf("Chat object %v\n", &chat)

	// send success message
	response.WriteHeader(200)
	response.Write([]byte(`{ "message": "Successfully added Chat group to chatspace" }`))
}
