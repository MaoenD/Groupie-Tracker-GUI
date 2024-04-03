package Functions

import (
	"fmt"
	"math"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

/********************************************************************************/
/************************************* OUTIL ************************************/
/********************************************************************************/

func getAverageColor(img *canvas.Image) (r, g, b, a uint32) {
	width := int(img.MinSize().Width)
	height := int(img.MinSize().Height)
	var totalR, totalG, totalB, totalA uint64

	// Parcours de tous les pixels de l'image pour calculer la somme des composantes de couleur
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colorRGBA := img.Image.At(x, y)
			r, g, b, a := colorRGBA.RGBA()
			totalR += uint64(r >> 8)
			totalG += uint64(g >> 8)
			totalB += uint64(b >> 8)
			totalA += uint64(a >> 8)
		}
	}

	// Calcul des moyennes arrondies
	pixelCount := uint64(width * height)
	r = uint32(totalR / pixelCount)
	g = uint32(totalG / pixelCount)
	b = uint32(totalB / pixelCount)
	a = uint32(totalA / pixelCount)

	// Vérifier si la moyenne des composantes de couleur est proche de 255
	maxComponent := uint32(255)
	tolerance := uint32(50) // Valeur de tolérance pour ajuster la couleur
	if r >= maxComponent-tolerance && g >= maxComponent-tolerance && b >= maxComponent-tolerance {
		// Réduire la composante rouge, verte et bleue pour éviter le blanc pur
		r = uint32(math.Max(0, float64(r)-float64(tolerance)))
		g = uint32(math.Max(0, float64(g)-float64(tolerance)))
		b = uint32(math.Max(0, float64(b)-float64(tolerance)))
	}

	return r, g, b, a
}

func LoadImageResource(path string) fyne.Resource {
	// Charger une ressource (image) à partir du chemin spécifié
	image, err := fyne.LoadResourceFromPath(path)
	// Vérifier s'il y a eu une erreur lors du chargement de l'image
	if err != nil {
		// Afficher un message d'erreur en cas d'échec du chargement de l'image
		fmt.Println("Erreur lors du chargement de l'icône:", err)
		// Retourner nil si une erreur s'est produite lors du chargement de l'image
		return nil
	}
	// Retourner la ressource (image) chargée avec succès
	return image
}

func checkMemberName(members []string, searchText string) bool {
	// Parcourir chaque membre de la liste
	for _, member := range members {
		// Vérifier si le nom du membre contient le texte de recherche (ignorer la casse)
		if strings.Contains(strings.ToLower(member), searchText) {
			return true // Retourner true si une correspondance est trouvée
		}
	}
	return false // Retourner false si aucune correspondance n'est trouvée
}

func checkConcertLocation(concerts []Concert, searchText string) bool {
	// Parcourir chaque concert dans la liste
	for _, concert := range concerts {
		// Vérifier si le lieu du concert contient le texte de recherche (ignorer la casse)
		for _, location := range concert.Locations {
			if strings.Contains(strings.ToLower(string(location)), searchText) {
				return true // Retourner true si une correspondance est trouvée
			}
		}
	}
	return false // Retourner false si aucune correspondance n'est trouvée
}

func contains(slice []string, str string) bool {
	// Vérifier si une chaîne est présente dans une slice de chaînes
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

const layoutDate = "02-01-2006"

func parseFirstAlbumDate(albumDate string) (time.Time, error) {
	return time.Parse(layoutDate, albumDate)
}

// UpdateArtistConcertInfo met à jour les informations sur les concerts de l'artiste en fonction des dates de concert fournies
func UpdateArtistConcertInfo(artist *Artist, concertDates []string, locations []Location) {
	// Initialise la variable pour stocker les détails du dernier concert
	var lastConcertDetails Concert

	// Parcourir les dates de concert pour trouver les concerts de l'artiste
	for _, date := range concertDates {
		// Convertir la date de chaîne à l'objet time.Time
		concertDate, err := time.Parse("02-01-2006", date) // Assurez-vous que le format correspond à celui des dates de concert dans les données
		if err != nil {
			// Gérer l'erreur si la date ne peut pas être analysée
			fmt.Println("Erreur lors de l'analyse de la date:", err)
			continue
		}

		// Parcourir les emplacements pour vérifier si l'artiste a joué à un concert à cette date
		for _, location := range locations {
			for _, eventDate := range location.DatesURL {
				// Convertir la date de chaîne à l'objet time.Time
				parsedEventDate, err := time.Parse("02-01-2006", string(eventDate)) // Assurez-vous que le format correspond à celui des dates de concert dans les données
				if err != nil {
					// Gérer l'erreur si la date ne peut pas être analysée
					fmt.Println("Erreur lors de l'analyse de la date de l'événement:", err)
					continue
				}

				// Vérifier si la date de concert correspond à la date de l'événement et si l'artiste a joué à cet endroit
				if parsedEventDate.Equal(concertDate) {
					// L'artiste a joué à cet endroit à la date du concert, mettre à jour les détails du dernier concert
					lastConcertDetails.Locations = append(lastConcertDetails.Locations, location.Locations...)
					lastConcertDetails.Dates = append(lastConcertDetails.Dates, string(eventDate))
				}
			}
		}
	}

	// Mettre à jour les informations du dernier concert de l'artiste
	artist.LastConcert = lastConcertDetails
}
