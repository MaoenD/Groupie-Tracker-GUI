package Functions

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"math"
	"strings"
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
