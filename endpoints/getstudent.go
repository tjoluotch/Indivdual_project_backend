package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"mygosource/ind_proj_backend/cors"
	"net/http"
	"strings"
	"time"
)

func GetStudentEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	cors.EnableCORS(&response)
	// set response header
	response.Header().Set("Content-Type", "application/json")

	// Get the request headers
	authHeader := strings.Split(request.Header.Get("Authorization"), "Bearer ")
	authTok := authHeader[1]
	signK := request.Header.Get("signK")

	fmt.Println("Token " + authTok)
	fmt.Println("Key " + signK)

	//DB access layer
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

	// decode JWT if it is successful go to the handler, if it is not successful send an error message and return
	token, _ := jwt.Parse(authTok, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in decoding JWT")
		}
		return []byte(signK),nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var student Student
		mapstructure.Decode(claims, &student)

		// Use student unique_id to find student in database
		uniqueID := &student.Unique_ID

		// Create search parameter of unique_id in bson format
		searchParams := "unique_id"
		filter := bson.D{{searchParams, uniqueID}}

		studentCollection := client.Database(dbName).Collection("students")
		//find a particular student using student id
		result := studentCollection.FindOne(context.Background(), filter)
		var studentFound Student
		// decode that student into Student struct
		err = result.Decode(&studentFound)
		// user does not exist error is returned to client
		if err != nil {
			http.Error(response, "Problem finding User by Unique ID on GetStudentEndpoint", http.StatusBadRequest)
			return
		}
		json.NewEncoder(response).Encode(studentFound)
	} else {
		http.Error(response, "Invalid Authorization token & code", 400)
	}
}
