package Functions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Define global variables to store data retrieved from API
var artists []Artist
var response json.RawMessage
var locationResponse LocationResponse
var locations []Location
var relationsResponse RelationsResponse
var relations []Relation

// LoadArtists loads artist data from the specified URL
func LoadArtists(url string) ([]Artist, error) {
	// Perform GET request to fetch artist data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON response into artists slice
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// LoadLocations loads location data from the specified URL
func LoadLocations(url string) ([]Location, error) {
	// Perform GET request to fetch location data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON response into a generic RawMessage
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Failed to unmarshal to RawMessage: %v", err)
		return nil, err
	}

	// Attempt to unmarshal RawMessage into LocationResponse struct
	err = json.Unmarshal(response, &locationResponse)
	if err != nil {
		log.Printf("Unexpected data format, not an object as expected: %v", err)
	} else {
		locations = locationResponse.Index
	}

	return locations, nil
}

// LoadRelations loads relation data from the specified URL
func LoadRelations(url string) ([]Relation, error) {
	// Perform GET request to fetch relation data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON response into a generic RawMessage
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Failed to unmarshal to RawMessage: %v", err)
		return nil, err
	}

	// Attempt to unmarshal RawMessage into RelationsResponse struct
	err = json.Unmarshal(response, &relationsResponse)
	if err != nil {
		log.Printf("Unexpected data format, not an object as expected: %v", err)
	} else {
		relations = relationsResponse.Index
	}

	return relations, nil
}

// CombineData combines location and relation data to form concert data
func CombineData(locationsURL, relationsURL string) ([]Concert, error) {
	// Fetch locations
	locations, err := LoadLocations(locationsURL)
	if err != nil {
		return nil, err
	}

	// Fetch relations
	relations, err := LoadRelations(relationsURL)
	if err != nil {
		return nil, err
	}

	// Initialize slice to contain concert information
	var concerts []Concert

	// Map locations with dates
	for _, location := range locations {
		for _, rel := range relations {
			for loc, dates := range rel.DatesLocations {
				if contains(location.Locations, loc) {
					// Create a slice containing the current location
					locSlice := []string{loc}
					concert := Concert{
						ID:        location.ID,
						Locations: locSlice, // Assign the slice containing the location
						Dates:     dates,
					}
					concerts = append(concerts, concert)
				}
			}
		}
	}

	return concerts, nil
}
