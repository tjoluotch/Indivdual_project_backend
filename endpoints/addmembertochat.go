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
	"mygosource/ind_proj_backend/twillio"
	"net/http"
	"time"
)

func AddMemberToGroupEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	cors.EnableCORS(&response)
	fmt.Println("Edit group members")
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
	studentCollection := client.Database(dbName).Collection("students")

	var findStudent Student
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&findStudent)
	//if there was an error panic
	if err != nil {
		http.Error(response, "JSON failed to decode request body to Student object " + err.Error() , 400)
		return
	}

	// search filter: search chat by id
	filterID := bson.D{{"student_id", findStudent.Student_ID}}
	filterFirstName := bson.D{{"first_name", findStudent.FirstName}}

	var result Student

	err = studentCollection.FindOne(context.TODO(), filterFirstName).Decode(&result)
	if err != nil {
		http.Error(response, "Issue searching DB for student by First Name", http.StatusBadRequest)
	}

	err = studentCollection.FindOne(context.TODO(), filterID).Decode(&result)
	if err != nil {
		http.Error(response, "Issue searching DB for student by either First Name or ID", http.StatusBadRequest)
		return
	}

	var original Student
	filterOriginalStudent := bson.D{{"student_id", student.Student_ID}}
	err = studentCollection.FindOne(context.TODO(), filterOriginalStudent).Decode(&original)
	if err != nil {
		http.Error(response, "Issue searching DB for student by either First Name or ID", http.StatusBadRequest)
		return
	}

	// add result to chat
	chatCollection := client.Database(dbName).Collection("chatspace")
	chatId := request.Header.Get("getChat")

	searchParam := bson.D{{"chat_id", chatId}}

	// update the task object within the module
	update := bson.M{"$push": bson.M{"members": result.Student_ID}}

	_, err = chatCollection.UpdateOne(context.TODO(), searchParam, update)
	if err != nil {
		http.Error(response, "Problem Editing Chatspace collection " + err.Error(), 400)
		return
	}

	var chat Chat
	err = chatCollection.FindOne(context.TODO(), searchParam).Decode(&chat)
	if err != nil {
		http.Error(response, "Issue searching DB to get Chat by Id", http.StatusBadRequest)
		return
	}

	// Send twilio SMS regarding Coursework with parameters: phone_no, coursework_desc, coursework_due, firstName
	err = twillio.AddedMemberSMSMessage(&result.Phone_No, &result.FirstName, &original.FirstName, &chat.Chat_Name)
	if err != nil {
		http.Error(response, "Problem Sending coursework sms " + err.Error(), 400)
		return
	}

	// send success message
	response.WriteHeader(200)
	response.Write([]byte(`{ "message": "Successfully Added member object" }`))
}
