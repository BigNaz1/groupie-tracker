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
}

type Location struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type Date struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

type Relation struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

var artists []Artist
var locations Location
var dates Date
var relations Relation

func main() {
	// Load data from APIs
	err := loadData()
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

func loadData() error {
	// Load artists data
	err := loadFromAPI("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		return err
	}

	// Load locations data
	err = loadFromAPI("https://groupietrackers.herokuapp.com/api/locations", &locations)
	if err != nil {
		return err
	}

	// Load dates data
	err = loadFromAPI("https://groupietrackers.herokuapp.com/api/dates", &dates)
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
func handleHome(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/home.html"))
    data := struct {
        Artists   []Artist
        Locations map[int][]string
        Dates     map[int][]string
    }{
        Artists: artists,
    }

    data.Locations = make(map[int][]string)
    for _, location := range locations.Index {
        data.Locations[location.ID] = location.Locations
    }

    data.Dates = make(map[int][]string)
    for _, date := range dates.Index {
        data.Dates[date.ID] = date.Dates
    }

    tmpl.Execute(w, data)
}