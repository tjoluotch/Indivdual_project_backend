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

func GetChatByIDEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	cors.EnableCORS(&response)
	fmt.Println("Get Chat by ID")
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
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error Connecting to MongoDB at context.WtihTimeout: Err-> %v\n ", err)
		return
	}
	chatCollection := client.Database(dbName).Collection("chatspace")

	// search filter: search chat by id
	filter := bson.D{{"chat_id", chatId}}

	var result Chat

	err = chatCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		http.Error(response, "Issue searching DB to get Chat by Id", http.StatusBadRequest)
		return
	}

	// Working - decode back to json
	json.NewEncoder(response).Encode(result)
}