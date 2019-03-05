package main

import (
	"github.com/gorilla/mux"
	"log"
	"mygosource/ind_proj_backend/endpoints"
	"mygosource/ind_proj_backend/envar"
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

	envar.Variables()

	httpPort := os.Getenv("HTTP_PORT")

	router := mux.NewRouter()


	router.HandleFunc("/api/signup", endpoints.CreateStudentAccountEndpoint).Methods("POST")
	router.HandleFunc("/api/authenticate", endpoints.AuthenticateEndpoint).Methods("POST")
	router.HandleFunc("/api/phonecode", endpoints.CheckMobileCodeEndpoint).Methods("POST")

	router.HandleFunc("/api/getstudent", endpoints.GetStudentEndpoint).Methods("GET")
	router.HandleFunc("/api/middleware-test", endpoints.ValidationMiddleware(endpoints.TestEndpoint)).Methods("GET")

	router.HandleFunc("/api/addmodule", endpoints.ValidationMiddleware(endpoints.AddModuleEndpoint)).Methods("PUT")
	router.HandleFunc("/api/getmodules", endpoints.ValidationMiddleware(endpoints.GetModulesEndpoint)).Methods("GET")

	router.HandleFunc("/api/addtask", endpoints.ValidationMiddleware(endpoints.AddTaskEndpoint)).Methods("PUT")
	router.HandleFunc("/api/editask", endpoints.ValidationMiddleware(endpoints.EditTaskEndpoint)).Methods("PUT")

	router.HandleFunc("/api/addcwk", endpoints.ValidationMiddleware(endpoints.AddCourseworkEndpoint)).Methods("PUT")

	router.HandleFunc("/api/createchat", endpoints.ValidationMiddleware(endpoints.CreateChatEndpoint)).Methods("PUT")
	router.HandleFunc("/api/getmychats", endpoints.ValidationMiddleware(endpoints.GetMyChatsEndpoint)).Methods("GET")
	router.HandleFunc("/api/getchatbyid", endpoints.ValidationMiddleware(endpoints.GetChatByIDEndpoint)).Methods("GET")


	//log server running
	log.Printf("server running on port %v", 12345)
	log.Fatal(http.ListenAndServe(httpPort,router))


}