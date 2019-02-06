package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Student struct {
	FirstName string `json:"first_name,omitempty"`
	LastName string	 `json:"last_name,omitempty"`
	Phone_No string	 `json:"phone_no,omitempty"`
	Student_ID string `json:"student_id,omitempty"`
	Unique_ID string `json:"unique_id,omitempty"`
}

type JwToken struct {
	Token string `json:"token,omitempty"`
}

type CustomsClaimsStudent struct {
	FirstName string `json:"first_name"`
	Student_ID string `json:"student_id"`
	Unique_ID string `json:"unique_id"`
	jwt.StandardClaims
}

type PhoneCode struct {
	Code string `json:"phoneCode"`
}

// Global variables to be able to use in each endpoint
var client *mongo.Client

const(
	dbName string = "Studently"
	privKeyPath string = "keys/app.rsa"
	pubKeyPath string = "keys/app.rsa.pub"
)

var VerifyKey, SignKey []byte



func LogFileSetup() {
	// Setting format and parsing current time in this format
	currentTime := time.Now().Format(time.RFC1123)

	// create a new file if one doesn't exists and append data to this file when writing.
	file, err := os.OpenFile("info.log"+currentTime, os.O_RDWR |os.O_CREATE|os.O_APPEND, 0666)

	// if there's an error with the opening of the log file log Fatal
	if err != nil {
		log.Fatal(err)
	}


	// set output to the file
	log.SetOutput(file)
	// log with date time and file location
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

func initKeys() {
	var err error

	SignKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading Private Key")
		return
	}

	VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading Public Key")
		return
	}
}


func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}



// Get the form data entered by client; FirstName, LastName, phone Number,
// assign the person a unique i.d
// check to see if that user isn't in the database already
// if they are send an error message with the a  'bad' response code
// if they aren't in db add to db and send a message with success
func CreateStudentAccountEndpoint(response http.ResponseWriter, request *http.Request){


	//CORS
	enableCORS(&response)

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


	response.Header().Set("Content-Type", "application/json")
	var student Student
	// decoding JSON post data to
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&student)
	//if there was an error panic
	if err != nil {
		panic(err)
	}
	// save unique id to student
	unId, _ := uuid.NewV4()
	idString := unId.String()
	student.Unique_ID = idString
	//print student to console
	fmt.Println(&student)
	//encode student into BSON
	data, err := bson.Marshal(student)
	if err != nil {
		log.Fatalf("Problem encoding Student struct into BSON: Err-> %v\n ",err)
	}
	studentCollection := client.Database(dbName).Collection("students")
	_, err = studentCollection.InsertOne(context.Background(),data)
	if err != nil {
		response.WriteHeader(501)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}
	// encoding json object for returning to the client
	jsonStudent, err := json.Marshal(student)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	response.Write(jsonStudent)
}

func CreateJwtokenEndpoint(response http.ResponseWriter, request *http.Request){
	//CORS
	enableCORS(&response)

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
	err = sendTwillioMessage(randCode, studentFound.Phone_No)
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

func sendTwillioMessage(code, phone_no string) error {

	// Set account keys & information
	accountSid := "AC32cd443ee4fc285c6a8d1b30805ae462"
	authToken := "8342021b04ecfd7990cfe31807ab56f4"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	twillioNo := "+447480534149"

	loginMessage := "Thanks for using Studently, please enter this Code: " + code

	// Pack up the data for the message
	msgData := url.Values{}
	msgData.Set("To", phone_no)
	msgData.Set("From", twillioNo)
	msgData.Set("Body", loginMessage)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Create HTTP request client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
			return err
		}
	} else {
		fmt.Println(resp.Status)
		err := errors.New("twillio didn't execute the SMS")
		return err
	}
	return nil
}

func CheckMobileCodeEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	enableCORS(&response)

	var phoneCode PhoneCode
	auth := request.Header.Get("Authorization")
	fmt.Println("Auth ", auth)
	r := request.Body

	json.NewDecoder(r).Decode(&phoneCode)
	fmt.Println("Code ", phoneCode.Code)

	newauth := strings.Split(auth, "Bearer ")
	fmt.Println("New Auth:", newauth[1])

	// decode JWT if it is successful go to the handler, if it is not successful send an error message and return
	token, _ := jwt.Parse(newauth[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in decoding JWT")
		}
		return []byte(phoneCode.Code),nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var student Student
		mapstructure.Decode(claims, &student)
		json.NewEncoder(response).Encode(student)
	} else {
		http.Error(response, "Invalid Authorization token", 400)
	}
}

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// middleware code
		// decode jwt through Authorization header
		handler.ServeHTTP(response,request)
	}
}

func GetStudentEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	enableCORS(&response)
	// set response header
	response.Header().Set("Content-Type", "application/json")


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

func main() {

	LogFileSetup()
	initKeys()
	router := mux.NewRouter()


	router.HandleFunc("/api/signup", CreateStudentAccountEndpoint).Methods("POST")
	router.HandleFunc("/api/authenticate", CreateJwtokenEndpoint).Methods("POST")
	router.HandleFunc("/api/phonecode", CheckMobileCodeEndpoint).Methods("POST")
	router.HandleFunc("/api/getstudent", GetStudentEndpoint).Methods("GET")

	//add a new route that gets the password as input along with the jwt from local storage and uses this to unlock this.
	// if JWT is decoded send back 200 along with the student object if JWT is not decoded send back 400

	//log server running
	log.Printf("server running on port %v", 12345)
	log.Fatal(http.ListenAndServe(":12345",router))


}