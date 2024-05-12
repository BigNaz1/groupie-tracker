package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

var artists []Artist

func main() {
	// Load artists data from API
	err := loadArtistsData()
	if err != nil {
		log.Fatal(err)
	}

	// Set up routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/artists", handleArtists)

	// Start the server
	log.Println("Server started on http://localhost:5500")
	log.Fatal(http.ListenAndServe(":5500", nil))
}

func loadArtistsData() error {
	// Make a GET request to the API endpoint
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the JSON response into the artists slice
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return err
	}

	return nil
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}

func handleArtists(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/artists.html"))
	tmpl.Execute(w, artists)
}
