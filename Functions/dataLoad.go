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
	var apiResponse APIResponse // Modifiez ici
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&apiResponse) // Modifiez ici
	if err != nil {
		return nil, err
	}

	return apiResponse.Index, nil // Modifiez ici
}

func LoadRelations(url string) ([]Relation, error) {
	var relations []Relation
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&relations)
	if err != nil {
		return nil, err
	}

	return relations, nil
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
	relations, err := LoadRelations(relationsURL)
	if err != nil {
		return nil, err
	}

	// Initialisation d'une slice pour contenir les infos
	var concerts []Concert

	// Mapping des locations avec les dates
	for _, location := range locations {
		for _, rel := range relations { // Modifiez ici
			for loc, dates := range rel.DatesLocations { // Modifiez ici
				if contains(location.Locations, loc) {
					// Créer un slice contenant le lieu actuel
					locSlice := []string{loc}
					concert := Concert{
						ID:        location.ID,
						Locations: locSlice, // Assigner le slice contenant le lieu
						Dates:     dates,
					}
					concerts = append(concerts, concert)
				}
			}
		}
	}

	return concerts, nil
}
