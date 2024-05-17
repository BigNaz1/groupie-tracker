package main

import (
	"log"
	"net/http"

	"GroupiTracker/Groupi"
)

func main() {
	// Load data from APIs
	err := Groupi.LoadData()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Set up routes
	mux.HandleFunc("/", Groupi.HandleRoot)
	mux.HandleFunc("/artist/", Groupi.HandleArtist)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	log.Println("Server started on http://localhost:5500")
	log.Fatal(http.ListenAndServe(":5500", mux))
}
