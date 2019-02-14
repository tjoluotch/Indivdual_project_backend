package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"math/rand"
	"mygosource/ind_proj_backend/cors"
	"mygosource/ind_proj_backend/twillio"
	"net/http"
	"strconv"
	"time"
)

func AuthenticateEndpoint(response http.ResponseWriter, request *http.Request){
	//CORS
	cors.EnableCORS(&response)

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

	// set response header
	response.Header().Set("Content-Type", "application/json")


	var student Student
	// decoding JSON post data to student
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&student)
	//if there was an error panic
	if err != nil {
		panic(err)
	}

	// Create search parameter of student id in bson format
	searchParams := "student_id"
	filter := bson.D{{searchParams, student.Student_ID}}

	studentCollection := client.Database(dbName).Collection("students")
	//find a particular student using student id
	result := studentCollection.FindOne(context.Background(), filter)
	var studentFound Student
	// decode that student into Student struct
	err = result.Decode(&studentFound)
	// user does not exist is return to client
	if err != nil {
		http.Error(response, "User does not exist ", http.StatusBadRequest)
		return
	}

	// encoding json object for returning to the client
	//jsonStudent, err := json.Marshal(&studentFound)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	// step 2 create a jwt from student found
	claims := CustomsClaimsStudent{
		studentFound.FirstName,
		studentFound.Student_ID,
		studentFound.Unique_ID,
		jwt.StandardClaims{
			Issuer: "golang api",
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
	}

	//generate randomm 6 digit int to be used as secret
	rand.Seed(time.Now().Unix())
	randInt := rand.Intn(999999)
	randCode:= strconv.Itoa(randInt)
	fmt.Println("Random Code for signing key ",randInt)

	// twillio message sent for auth
	err = twillio.SendTwillioMessage(randCode, studentFound.Phone_No)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	mySigningKey := []byte(randCode)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(response, "Issue with JWT creation: Error message -> " + err.Error(), http.StatusBadRequest)
		return
	}

	//response.Header().Add("Authorization", "Bearer " + tokenString)
	json.NewEncoder(response).Encode(JwToken{Token: tokenString})
	//response.Write([]byte("Token sent, Login successful"))
}
