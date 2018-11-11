package main

import (
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

func HomeHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Individual project!"))
}

func main() {

	LogFileSetup()

	router := mux.NewRouter()
	// Routes consist of a Path and Handler function

	// Router to HomeHandler
	router.HandleFunc("/", HomeHandler)

	log.Printf("server running on port %v", 6000)
	//Server start listen on port 6000
	log.Fatal(http.ListenAndServe(":6000", nil))



}