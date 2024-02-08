package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Artist struct to hold information about artists or groups
type Artist struct {
	Name        string
	Image       string
	YearStarted int
	DebutAlbum  time.Time
	Members     []string
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Menu - Groupie Tracker")

	// Exemple de données pour les artistes
	artists := []Artist{
		{Name: "Michael Jackson", Image: "public/michaeljackson.jpg", YearStarted: 1964, DebutAlbum: time.Date(1972, time.November, 13, 0, 0, 0, 0, time.UTC), Members: []string{"Michael Jackson"}},
		{Name: "Queen", Image: "public/queen.png", YearStarted: 1970, DebutAlbum: time.Date(1973, time.July, 13, 0, 0, 0, 0, time.UTC), Members: []string{"Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"}},
		{Name: "Pink Floyd", Image: "public/pinkfloyd.jpeg", YearStarted: 1965, DebutAlbum: time.Date(1967, time.August, 5, 0, 0, 0, 0, time.UTC), Members: []string{"Syd Barrett", "Roger Waters", "Richard Wright", "Nick Mason"}},
	}

	// Champ de recherche
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Artists...")

	// Zone d'affichage des résultats de recherches
	searchResults := container.NewVBox()

	// Bouton de recherche
	searchButton := widget.NewButton("Search", func() {
		// Récupérer le texte de la recherche
		searchText := searchBar.Text

		// Recherche d'artistes correspondants
		var foundArtists []Artist
		for _, artist := range artists {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchText)) {
				foundArtists = append(foundArtists, artist)
			}
		}

		// Afficher les résultats
		searchResultsObjects := make([]fyne.CanvasObject, 0)
		if len(foundArtists) > 0 {
			resultText := fmt.Sprintf("Artist found: %s\n\n", searchText)
			searchResultsObjects = append(searchResultsObjects, widget.NewLabel(resultText))
			for _, artist := range foundArtists {
				card := createCardGeneralInfo(artist)
				searchResultsObjects = append(searchResultsObjects, card)
			}
		} else {
			searchResultsObjects = append(searchResultsObjects, widget.NewLabel("No artist found: "+searchText))
		}

		searchResults.Objects = searchResultsObjects
		searchResults.Refresh()
	})

	// Ajouter l'écouteur d'événements clavier au champ de recherche
	searchBar.OnSubmitted = func(_ string) {
		searchButton.OnTapped()
	}

	// Création du conteneur pour afficher les artistes
	artistsContainer := container.NewHBox()

	// Ajouter les cartes des artistes dans le conteneur
	for _, artist := range artists {
		card := createCardGeneralInfo(artist)
		artistsContainer.Add(card)
	}

	// Création du conteneur principal avec la couleur de fond spécifiée
	content := container.NewVBox(
		searchBar,
		searchButton,
		searchResults,
		artistsContainer, // Utilisation du conteneur HBox pour afficher les cartes côte à côte
	)

	// Centrer les cartes dans la fenêtre
	centeredContent := container.New(layout.NewCenterLayout(), content)

	// Créer un rectangle pour représenter le fond coloré
	background := canvas.NewRectangle(color.NRGBA{R: 0x5C, G: 0x64, B: 0x73, A: 0xFF})
	background.Resize(fyne.NewSize(1080, 720)) // Taille pour remplir toute la fenêtre

	// Créer un conteneur pour contenir le fond coloré
	backgroundContainer := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), background)

	// Ajouter le contenu principal au conteneur de fond
	backgroundContainer.Add(centeredContent)

	// Afficher la fenêtre
	myWindow.SetContent(backgroundContainer)
	myWindow.Resize(fyne.NewSize(1080, 720)) // ajustement size
	myWindow.ShowAndRun()
}

func createCard(artist Artist) fyne.CanvasObject {
	// Redimensionner l'image
	image := canvas.NewImageFromFile(artist.Image)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(120, 120))
	image.Resize(fyne.NewSize(120, 120))

	// Nom de l'artiste en gras et plus gros
	nameLabel := canvas.NewText(artist.Name, theme.TextColor())

	// Date de début de l'artiste en plus petit
	yearLabel := canvas.NewText(fmt.Sprintf("Year Started: %d", artist.YearStarted), theme.TextColor())

	// Album de début de l'artiste
	debutLabel := canvas.NewText(fmt.Sprintf("Debut Album: %s", artist.DebutAlbum.Format("2006-01-02")), theme.TextColor())

	// Créer le conteneur pour les informations sur l'artiste
	infoContainer := container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		nameLabel,
		layout.NewSpacer(),
		container.NewHBox(
			yearLabel,
			layout.NewSpacer(),
			debutLabel,
		),
		layout.NewSpacer(),
		layout.NewSpacer(),
	)

	// Définir la taille fixe pour le conteneur d'informations
	infoContainer.Resize(fyne.NewSize(300, 180))

	// Créer le conteneur pour la carte de l'artiste
	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), image, infoContainer)
	cardContent.Resize(fyne.NewSize(300, 300))

	// Créer un rectangle pour le contour avec des coins arrondis
	border := canvas.NewRectangle(color.Transparent) // Définir une couleur transparente pour le remplissage
	border.SetMinSize(fyne.NewSize(300, 300))
	border.Resize(fyne.NewSize(296, 296)) // Redimensionner légèrement la bordure pour inclure les coins arrondis
	border.StrokeColor = color.Black      // Définir la couleur de la bordure
	border.StrokeWidth = 3                // Définir l'épaisseur de la bordure
	border.CornerRadius = 20              // Définir les coins arrondis

	// Ajouter le rectangle de contour à la carte
	cardContent.AddObject(border)

	return cardContent
}

func createCardGeneralInfo(artist Artist) fyne.CanvasObject {
	// Redimensionner l'image
	image := canvas.NewImageFromFile(artist.Image)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(120, 120))
	image.Resize(fyne.NewSize(120, 120))

	// Obtenir la couleur moyenne de l'image
	averageColor := getAverageColor(artist.Image)

	// Créer un rectangle coloré pour l'arrière-plan de la carte
	background := canvas.NewRectangle(averageColor)
	background.Resize(fyne.NewSize(300, 300))

	// Nom de l'artiste en gras et plus gros
	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Date de début de l'artiste en plus petit
	yearLabel := widget.NewLabel(fmt.Sprintf("%d", artist.YearStarted))

	// Créer un conteneur HBox pour afficher les labels avec un espace entre eux
	labelsContainer := container.NewHBox(
		nameLabel,
		yearLabel,
	)

	// Membres du groupe, s'il y en a, en plus petit
	var membersText string
	if len(artist.Members) == 1 {
		membersText = "Solo Artist"
	} else if len(artist.Members) > 0 {
		membersText = "Members:\n " + strings.Join(artist.Members, ", ")
	}
	membersLabel := widget.NewLabel(membersText)
	membersLabel.Wrapping = fyne.TextWrapWord // Activer le wrapping du texte

	// Créer le conteneur pour les informations sur l'artiste
	infoContainer := container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(), // Ajout d'un espace vertical
		image,
		labelsContainer,    // Placer les labels sur la même ligne avec un espace entre eux
		membersLabel,       // Afficher les membres du groupe
		layout.NewSpacer(), // Ajout d'un petit espace vertical
	)

	// Définir la taille fixe pour le conteneur d'informations
	infoContainer.Resize(fyne.NewSize(300, 180))

	// Créer le conteneur pour la carte de l'artiste
	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background, infoContainer)
	cardContent.Resize(fyne.NewSize(300, 300))

	// Créer un rectangle pour le contour avec des coins arrondis
	border := canvas.NewRectangle(color.Transparent) // Définir une couleur transparente pour le remplissage
	border.SetMinSize(fyne.NewSize(300, 300))
	border.Resize(fyne.NewSize(296, 296)) // Redimensionner légèrement la bordure pour inclure les coins arrondis
	border.StrokeColor = color.Black      // Définir la couleur de la bordure
	border.StrokeWidth = 3                // Définir l'épaisseur de la bordure
	border.CornerRadius = 20              // Définir les coins arrondis

	// Ajouter le rectangle de contour à la carte
	cardContent.AddObject(border)

	return cardContent
}

// Fonction pour obtenir la couleur moyenne d'une image
func getAverageColor(imagePath string) color.Color {
	// Ouvrir le fichier image
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return color.Black // Retourner noir en cas d'erreur
	}
	defer file.Close()

	// Décoder l'image
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Erreur lors du décodage de l'image:", err)
		return color.Black // Retourner noir en cas d'erreur
	}

	// Initialiser les variables pour stocker la somme des composantes de couleur
	var totalRed, totalGreen, totalBlue uint32
	totalPixels := 0

	// Parcourir tous les pixels de l'image
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Obtenir la couleur du pixel
			pixelColor := img.At(x, y)
			r, g, b, _ := pixelColor.RGBA()

			// Ajouter les composantes de couleur à la somme totale
			totalRed += r
			totalGreen += g
			totalBlue += b

			// Incrémenter le nombre total de pixels
			totalPixels++
		}
	}

	// Calculer la moyenne des composantes de couleur en divisant par le nombre total de pixels
	averageRed := totalRed / uint32(totalPixels)
	averageGreen := totalGreen / uint32(totalPixels)
	averageBlue := totalBlue / uint32(totalPixels)

	// Créer et retourner la couleur moyenne
	averageColor := color.RGBA{
		R: uint8(averageRed >> 8),
		G: uint8(averageGreen >> 8),
		B: uint8(averageBlue >> 8),
		A: 130, // Opacité
	}

	return averageColor
}
