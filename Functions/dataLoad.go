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

func LoadLocations(url string) ([]Location, error) {
	var location []Location
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &location)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func LoadRelations(url string) ([]Concert, error) {
	var concerts []Concert
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&concerts)
	if err != nil {
		return nil, err
	}

	return concerts, nil
}

// coucou
func LoadDate(url string) ([]Dates, error) {
	var dates []Dates
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &dates)
	if err != nil {
		return nil, err
	}

	return dates, nil
}
