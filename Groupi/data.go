package Groupi

import (
	"encoding/json"
	"net/http"
)

func LoadData() error {
	err := loadFromAPI("https://groupietrackers.herokuapp.com/api/artists", &Artists)
	if err != nil {
		return err
	}

	// Load relations data
	err = loadFromAPI("https://groupietrackers.herokuapp.com/api/relation", &Relations)
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
