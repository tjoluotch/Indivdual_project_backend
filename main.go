package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

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

func SignUpEndpoint(response http.ResponseWriter, request *http.Request){
	response.Write([]byte("Individual project!"))
}

func main() {

	LogFileSetup()


	router := mux.NewRouter()
	// CORS setup so that frontend can access this API
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/signup", SignUpEndpoint).Methods("POST")
	//log server running
	log.Printf("server running on port %v", 6000)
	log.Fatal(http.ListenAndServe(":6000", handlers.CORS(headers, methods, origins)(router)))

}