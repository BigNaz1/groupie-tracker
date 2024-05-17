package Groupi

import (
	"html/template"
	"net/http"
	"strconv"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFoundHandler(w)
		return
	}
	HandleHome(w)
}

func HandleArtist(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Path[len("/artist/"):]
	id, err := strconv.Atoi(artistID)
	if err != nil {
		BadRequestHandler(w)
		return
	}

	var artist Artist
	for _, a := range Artists {
		if a.ID == id {
			artist = a
			break
		}
	}

	if artist.ID == 0 {
		NotFoundHandler(w)
		return
	}

	// Simulate an internal server error for a specific artist ID
	if artist.ID == 123 {
		InternalServerErrorHandler(w)
		return
	}

	var relation struct {
		DatesLocations map[string][]string `json:"datesLocations"`
	}
	for _, r := range Relations.Index {
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
		InternalServerErrorHandler(w)
		return
	}
}

func NotFoundHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	tmpl, err := template.ParseFiles("templates/404.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func BadRequestHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl, err := template.ParseFiles("templates/400.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func InternalServerErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl, err := template.ParseFiles("templates/500.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func HandleHome(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	data := struct {
		Artists []Artist
	}{
		Artists: Artists,
	}

	tmpl.Execute(w, data)
}
