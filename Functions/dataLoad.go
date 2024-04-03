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
	var locations []Location
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func LoadRelations(url string) (Relation, error) {
	var relation Relation
	resp, err := http.Get(url)
	if err != nil {
		return Relation{}, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&relation)
	if err != nil {
		return Relation{}, err
	}

	return relation, nil
}

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

func CombineData(locationsURL, relationsURL string) ([]Concert, error) {
	// Fetch location
	locations, err := LoadLocations(locationsURL)
	if err != nil {
		return nil, err
	}

	// Fetch relation
	relation, err := LoadRelations(relationsURL)
	if err != nil {
		return nil, err
	}

	// Initialisation d'une slice pour contenir les infos
	var concerts []Concert

	// Mapping des locations avec les dates
	for _, location := range locations {
		for loc, dates := range relation.DatesLocations {

			if contains(location.Locations, loc) {
				concert := Concert{
					ID:        location.ID,
					Locations: loc,
					Dates:     dates,
				}
				concerts = append(concerts, concert)
			}
		}
	}

	return concerts, nil
}
