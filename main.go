package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type APIUrls struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID         int `json:"id"`
	ArtistID   int `json:"artistId"`
	LocationID int `json:"locationId"`
	DateID     int `json:"dateId"`
}

var (
	artists   []Artist
	locations []Location
	dates     []Date
	relations []Relation
)

func main() {
	http.HandleFunc("/", serveFiles)
	http.HandleFunc("/api/artists", artistHandler)
	http.HandleFunc("/api/locations", locationHandler)
	http.HandleFunc("/api/dates", dateHandler)
	http.HandleFunc("/api/relations", relationHandler)

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchAPI(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

func init() {
	apis := APIUrls{
		Artists:   "https://groupietrackers.herokuapp.com/api/artists",
		Locations: "https://groupietrackers.herokuapp.com/api/locations",
		Dates:     "https://groupietrackers.herokuapp.com/api/dates",
		Relation:  "https://groupietrackers.herokuapp.com/api/relation",
	}

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		fetchAPI(apis.Artists, &artists)
	}()
	go func() {
		defer wg.Done()
		fetchAPI(apis.Locations, &locations)
	}()
	go func() {
		defer wg.Done()
		fetchAPI(apis.Dates, &dates)
	}()
	go func() {
		defer wg.Done()
		fetchAPI(apis.Relation, &relations)
	}()

	wg.Wait()
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(artists)
}

func locationHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(locations)
}

func dateHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(dates)
}

func relationHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(relations)
}
