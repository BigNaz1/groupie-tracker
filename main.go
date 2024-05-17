package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Relation struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

var artists []Artist
var relations Relation

func main() {
	// Load data from APIs
	err := loadData()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Set up routes
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/artist/", handleArtist)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	log.Println("Server started on http://localhost:5500")
	log.Fatal(http.ListenAndServe(":5500", mux))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w)
		return
	}
	handleHome(w)
}

func notFoundHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	tmpl, err := template.ParseFiles("templates/404.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func loadData() error {
	// Load artists data
	err := loadFromAPI("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		return err
	}

	// Load relations data
	err = loadFromAPI("https://groupietrackers.herokuapp.com/api/relation", &relations)
	if err != nil {
		return err
	}

	return nil
}

func loadFromAPI(url string, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func badRequestHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl, err := template.ParseFiles("templates/400.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func internalServerErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl, err := template.ParseFiles("templates/500.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleHome(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	data := struct {
		Artists []Artist
	}{
		Artists: artists,
	}

	tmpl.Execute(w, data)
}

func handleArtist(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Path[len("/artist/"):]
	id, err := strconv.Atoi(artistID)
	if err != nil {
		badRequestHandler(w)
		return
	}

	var artist Artist
	for _, a := range artists {
		if a.ID == id {
			artist = a
			break
		}
	}

	if artist.ID == 0 {
		notFoundHandler(w)
		return
	}

	// Simulate an internal server error for a specific artist ID
	if artist.ID == 123 {
		internalServerErrorHandler(w)
		return
	}

	var relation struct {
		DatesLocations map[string][]string `json:"datesLocations"`
	}
	for _, r := range relations.Index {
		if r.ID == id {
			relation.DatesLocations = r.DatesLocations
			break
		}
	}

	data := struct {
		Artist         Artist
		DatesLocations map[string][]string
	}{
		Artist:         artist,
		DatesLocations: relation.DatesLocations,
	}

	tmpl := template.Must(template.ParseFiles("templates/artist.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		internalServerErrorHandler(w)
		return
	}
}
