package Functions

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
)

/********************************************************************************/
/************************************* OUTIL ************************************/
/********************************************************************************/
func getAverageColor(imagePath string) color.Color {
	// Ouvrir le fichier image
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return color.Black
	}
	defer file.Close()

	// Décoder l'image
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Erreur lors du décodage de l'image:", err)
		return color.Black
	}

	// Initialiser les variables pour les composantes de couleur totales et le nombre total de pixels
	var totalRed, totalGreen, totalBlue float64
	totalPixels := 0

	// Parcourir tous les pixels de l'image
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Obtenir la couleur du pixel
			pixelColor := img.At(x, y)
			r, g, b, _ := pixelColor.RGBA()

			// Convertir les composantes de couleur en valeurs flottantes normalisées
			red := float64(r) / 65535.0
			green := float64(g) / 65535.0
			blue := float64(b) / 65535.0

			// Ajouter les composantes de couleur aux totaux
			totalRed += red
			totalGreen += green
			totalBlue += blue

			// Incrémenter le nombre total de pixels
			totalPixels++
		}
	}

	// Calculer les moyennes des composantes de couleur
	averageRed := totalRed / float64(totalPixels)
	averageGreen := totalGreen / float64(totalPixels)
	averageBlue := totalBlue / float64(totalPixels)

	// Mettre à l'échelle les valeurs de couleur moyennes à l'échelle de 0 à 255
	averageRed = averageRed * 255
	averageGreen = averageGreen * 255
	averageBlue = averageBlue * 255

	// Créer une couleur moyenne avec les composantes de couleur calculées
	averageColor := color.RGBA{
		R: uint8(averageRed),
		G: uint8(averageGreen),
		B: uint8(averageBlue),
		A: 255,
	}

	return averageColor
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
		if strings.Contains(strings.ToLower(concert.Location), searchText) {
			return true // Retourner true si une correspondance est trouvée
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
