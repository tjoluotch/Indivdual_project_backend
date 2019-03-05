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

func EditGroupNameEndpoint(response http.ResponseWriter, request *http.Request) {
	// CORS
	cors.EnableCORS(&response)
	fmt.Println("Edit group name Endpoint")
	// Decode jwt claims into student Model
	decoded := context2.Get(request, "decoded")
	var student Student
	claims := decoded.(jwt.MapClaims)
	mapstructure.Decode(claims, &student)

	// test to see if request header is coming through
	chatId := request.Header.Get("getChat")
	fmt.Println("Chat ID:", chatId)

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

	// Decode request body into chat struct
	var chat Chat
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&chat)
	//if there was an error
	if err != nil {
		http.Error(response, "JSON failed to decode request body to Chat object " + err.Error() , 400)
		return
	}
	//fmt.Println(&chat)

	// search parameters to find the chat object to be edited
	searchParams := "chat_id"
	filter := bson.D{{searchParams, chatId}}
	// update the chat name
	update := bson.M{"$set": bson.M{"chat_name": chat.Chat_Name}}

	_,err = chatCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(response, "Problem Editing Chat name see editchatnem.go " + err.Error(), 400)
		return
	}

	// send success message
	response.WriteHeader(200)
	response.Write([]byte(`{ "message": "Successfully Renamed revision group chat" }`))
}
