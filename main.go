package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type User struct {
	FirstName string `json:"first_name,omitempty"`
	LastName string	 `json:"last_name,omitempty"`
	U_ID string		 `json:"unique_id,omitempty"`
	Phone_No string	 `json:"phone_no,omitempty"`
}

func LogFileSetup() {
	// create a new file if one doesn't exists and append data to this file when writing.
	file, err := os.OpenFile("info.log", os.O_RDWR |os.O_CREATE|os.O_APPEND, 0666)

	// if there's an error with the opening of the log file log Fatal
	if err != nil {
		log.Fatal(err)
	}


	// set output to the file
	log.SetOutput(file)
	// log with date time and file location
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

// Get the form data entered by client; FirstName, LastName, phone Number,
// assign the person a unique i.d
// check to see if that user isn't in the database already
// if they are send an error message with the a  'bad' response code
// if they aren't in db add to db and send a message with success
func SignUpEndpoint(response http.ResponseWriter, request *http.Request){
	log.Print("Entered func SignUpEndpoint")
	log.Printf("Request: %v", request)
	response.Write([]byte("Individual project!"))
}

func main() {

	LogFileSetup()


	router := mux.NewRouter()
	// CORS setup so that frontend can access this API
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/signup", SignUpEndpoint)
	//log server running
	log.Printf("server running on port %v", 8080)
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router)))

}