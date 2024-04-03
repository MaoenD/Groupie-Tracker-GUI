package search

import (
	"strings"
	"time"
)

const layout = "02-01-2006"
const layoutWithPrefix = "*02-01-2006"

// Struct simplifiée
type SearchSuggestion struct {
	Label string
	Type  string
}

// Recherche générale avec suggestions
func SearchWithSuggestions(artists []Artist, query string) []SearchSuggestion {
	var suggestions []SearchSuggestion
	lowerQuery := strings.ToLower(query)

	for _, artist := range artists {
		// Suggestion nom de artiste
		if strings.Contains(strings.ToLower(artist.Name), lowerQuery) {
			suggestions = append(suggestions, SearchSuggestion{Label: artist.Name + " - Artist/Band", Type: "artist"})
		}

		// Suggestions membres
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), lowerQuery) {
				suggestions = append(suggestions, SearchSuggestion{Label: member + " - Member", Type: "member"})
			}
		}
	}
	return suggestions
}

// Filtre par nombre de membres avec chekbox pour sélection multiple
func ByMembersCount(artists []Artist, counts []int) []Artist {
	var filtered []Artist
	countMap := make(map[int]bool)
	for _, count := range counts {
		countMap[count] = true
	}

	for _, artist := range artists {
		if _, exists := countMap[len(artist.Members)]; exists {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

// Filtre par plage de date du premier album.
func FilterByConcertDateRange(artists []Artist, relations []Relation, datesData []Dates, startDate, endDate string) []Artist {
	var filtered []Artist
	startTime, _ := time.Parse(layout, startDate)
	endTime, _ := time.Parse(layout, endDate)

	for _, artist := range artists {
		relation := RelationByArtID(relations, artist.ID)
		if relation == nil {
			continue
		}
		dates := datesData[relation.ID-1]
		for _, dateString := range dates.Dates {
			date, _ := time.Parse("*02-01-2006", dateString)
			if (date.After(startTime) || date.Equal(startTime)) && (date.Before(endTime) || date.Equal(endTime)) {
				filtered = append(filtered, artist)
				break
			}
		}
	}
	return filtered
}

// Filtre par lieux de concerts avec checkbox
func ByConcert_Locations_CheckBox(artists []Artist, relations []Relation, selectedLocations []string) []Artist {
	var filtered []Artist
	for _, artist := range artists {
		relation := RelationByArtID(relations, artist.ID)
		if relation == nil {
			continue
		}
		for _, location := range selectedLocations {
			if _, exists := relation.DatesLocations[location]; exists {
				filtered = append(filtered, artist)
				break
			}
		}
	}
	return filtered
}

// Trouver la relation correspondant à un ID d'artiste
func RelationByArtID(relations []Relation, artistID int) *Relation {
	for _, relation := range relations {
		if relation.ID == artistID {
			return &relation
		}
	}
	return nil
}
