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

func SendMessageEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	cors.EnableCORS(&response)
	fmt.Println("send Message Endpoint")
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

	// find student by unique id
	studentCollection := client.Database(dbName).Collection("students")

	studentFilter := bson.D{{"unique_id", student.Unique_ID}}
	err = studentCollection.FindOne(context.TODO(), studentFilter).Decode(&student)
	if err != nil {
		http.Error(response, "Issue searching DB to get Student by unique Id "+ err.Error(), http.StatusBadRequest)
		return
	}

	// generate unique Id for the sent message
	messageID, err := uuid.NewV4()
	if err != nil {
		http.Error(response, "Unique ID failed to generate for new Message object" , 400)
		return
	}
	messageIDString := messageID.String()

	// generate a timestamp for the message
	t := time.Now()
	timeStamp:= t.Format("Mon _2 Jan 15:04 2006")
	//fmt.Println(timeStamp)


	// decoding JSON put data into Message Struct
	var message Message
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&message)
	//if there was an error panic
	if err != nil {
		http.Error(response, "JSON failed to decode request body to Message object " + err.Error() , 400)
		return
	}
	// Encode message struct fully
	message.Message_ID = messageIDString
	message.Sender = student
	message.Sent_At = timeStamp

	//fmt.Println(message)

	// search parameter to find the chat in the db
	searchParams := "chat_id"
	chatFilter := bson.D{{searchParams, chatId}}

	// update the chat object to have an array of message type
	update := bson.M{"$push": bson.M{"messages": message}}

	_, err = chatCollection.UpdateOne(context.TODO(), chatFilter, update)
	if err != nil {
		http.Error(response, "Problem updating Chat to include latest message " + err.Error(), 400)
		return
	}

	// send success message
	response.WriteHeader(200)
	response.Write([]byte(`{ "message": "Successfully added Message messages array in Chat object" }`))

}
