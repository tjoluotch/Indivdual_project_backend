package main

import (
	"github.com/gorilla/mux"
	"log"
	"mygosource/ind_proj_backend/endpoints"
	"net/http"
	"os"
	"time"
)

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


func main() {

	LogFileSetup()

	router := mux.NewRouter()


	router.HandleFunc("/api/signup", endpoints.CreateStudentAccountEndpoint).Methods("POST")
	router.HandleFunc("/api/authenticate", endpoints.AuthenticateEndpoint).Methods("POST")
	router.HandleFunc("/api/phonecode", endpoints.CheckMobileCodeEndpoint).Methods("POST")

	router.HandleFunc("/api/getstudent", endpoints.GetStudentEndpoint).Methods("GET")
	router.HandleFunc("/api/middleware-test", endpoints.ValidationMiddleware(endpoints.TestEndpoint)).Methods("GET")

	router.HandleFunc("/api/addmodule", endpoints.ValidationMiddleware(endpoints.AddModuleEndpoint)).Methods("PUT")

	//add a new route that gets the password as input along with the jwt from local storage and uses this to unlock this.
	// if JWT is decoded send back 200 along with the student object if JWT is not decoded send back 400

	//log server running
	log.Printf("server running on port %v", 12345)
	log.Fatal(http.ListenAndServe(":12345",router))


}