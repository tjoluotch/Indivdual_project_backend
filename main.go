package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func HomeHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Individual project!"))
}

func main() {
	file := os.OpenFile("info.log")
	router := mux.NewRouter()
	// Routes consist of a Path and Handler function

	// Router to HomeHandler
	router.HandleFunc("/", HomeHandler)
	// Server start listen on port 6000
	log.Fatal(http.ListenAndServe(":6000", nil))

}