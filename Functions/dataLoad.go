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
	var apiResponse APIResponse
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse.Index, nil
}

func LoadRelations(url string) ([]Relation, error) {
	var apiResponse RelationsResponse
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse.Index, nil
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
	// Charger les données des artistes
	artists, err := LoadArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}

	// Charger les données des relations
	relations, err := LoadRelations(relationsURL)
	if err != nil {
		return nil, err
	}

	// Charger les données des emplacements
	locations, err := LoadLocations(locationsURL)
	if err != nil {
		return nil, err
	}

	// Créer une carte pour mapper les ID des artistes avec leurs emplacements
	artistLocations := make(map[int][]string)
	for _, location := range locations {
		artistLocations[location.ID] = location.Locations
	}

	// Initialisation d'une slice pour contenir les concerts
	var concerts []Concert

	// Parcours des relations et création des concerts
	for _, rel := range relations {
		for _, artist := range artists {
			if artist.ID == rel.ID { // Vérification de l'ID de l'artiste dans les relations
				// Vérifier si l'artiste a des concerts
				if dates, found := rel.DatesLocations[artist.Name]; found {
					// Parcours des dates pour créer les concerts
					for _, date := range dates {
						// Récupérer les emplacements de chaque artiste à partir de la liste des emplacements
						locations, ok := artistLocations[artist.ID]
						if !ok {
							continue // Si aucun emplacement trouvé pour l'artiste, passer au suivant
						}

						concert := Concert{
							ID:        artist.ID,
							Locations: locations, // Utilisation des emplacements de l'artiste
							Dates:     []string{date},
						}
						concerts = append(concerts, concert)
					}
				}
			}
		}
	}

	return concerts, nil
}
