package Functions

import (
	"encoding/json"
	"io"
	"net/http"
)

// Charge les données des artistes depuis une URL spécifiée
func LoadArtists(url string) ([]Artist, error) {
	var artists []Artist
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func LoadLocations(url string) ([]Concert, error) {
	var locations []Concert
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func LoadRelations(url string) ([]Relation, error) {
	var relations []Relation
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &relations)
	if err != nil {
		return nil, err
	}

	return relations, nil
}
