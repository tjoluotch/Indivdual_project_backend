package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/couchbase/gocb.v1"
	"log"
	"net/http"
	"os"
	"time"
)

type Member struct {
	FirstName string `json:"first_name,omitempty"`
	LastName string	 `json:"last_name,omitempty"`
	U_ID string		 `json:"unique_id,omitempty"`
	Phone_No string	 `json:"phone_no,omitempty"`
	Student_ID string `json: "student_id, omitempty"`
}

// Global variables to be able to use in each endpoint
var bucket *gocb.Bucket
var cluster *gocb.Cluster


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

func enableCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// Get the form data entered by client; FirstName, LastName, phone Number,
// assign the person a unique i.d
// check to see if that user isn't in the database already
// if they are send an error message with the a  'bad' response code
// if they aren't in db add to db and send a message with success
func CreateStudentAccountEndpoint(response http.ResponseWriter, request *http.Request){
	/*
	var member Member
	var n1qlParams []interface{}
	_ = json.NewDecoder(request.Body).Decode(&member)
	query := gocb.NewN1qlQuery("INSERT INTO stu_prod_hub (KEY,VALUE) values ($1)")
	*/
	enableCORS(&response, request)
	fmt.Println(request.Body)
	fmt.Println("im in")
}

func main() {

	LogFileSetup()

	//DB logger for Couchbase - prints to Command Line
	gocb.SetLogger(gocb.DefaultStdioLogger())

	// Db access Layer
	dbUsername := "admin"
	dbPass := "admin1997"
	cluster, errClust := gocb.Connect("couchbase://127.0.0.1")
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: dbUsername,
		Password: dbPass,
	})
	// throw error if there is a problem with opening the DB cluster
	if errClust != nil {
		log.Fatal("Error: Problem with setting up cluster:",errClust)
	}
	// dereference cluster pointer
	log.Printf("DataBase cluster setup at: %v", *cluster)

	bucket, errBuc := cluster.OpenBucket("stu_prod_hub", "")
	if errBuc != nil {
		log.Fatal("Error: Problem with DB Bucket",errBuc)
	}
	log.Printf("Bucket setup correctly at: %v", *bucket)


	router := mux.NewRouter()


	router.HandleFunc("/api/signup", CreateStudentAccountEndpoint).Methods("PUT")
	//log server running
	log.Printf("server running on port %v", 8080)
	log.Fatal(http.ListenAndServe(":8080", router))


}