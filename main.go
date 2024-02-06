package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

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

// Venue struct to hold information about concert venues
type Venue struct {
	Name     string
	Location string
	Capacity int
}

// ConcertDate struct to hold information about concert dates
type ConcertDate struct {
	Date    time.Time
	IsPast  bool // Indicates whether the concert date is in the past or future
	Artists []Artist
	Venue   Venue
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")

	// Champ de recherche
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Artists...")

	// Zone d'affichage des résultats de recherche
	searchResults := container.NewVBox()

	// Bouton de recherche
	searchButton := widget.NewButton("Search", func() {
		// Récupérer le texte de la recherche
		searchText := searchBar.Text

		// Exemple de données pour les artistes, les lieux et les dates de concert
		// Vous pouvez les remplacer par vos propres données
		artists := []Artist{
			{Name: "Michael Jackson", Image: "michael_jackson.jpg", YearStarted: 1964, DebutAlbum: time.Date(1972, time.November, 13, 0, 0, 0, 0, time.UTC), Members: []string{"Michael Jackson"}},
			{Name: "Queen", Image: "public/queen.png", YearStarted: 1970, DebutAlbum: time.Date(1973, time.July, 13, 0, 0, 0, 0, time.UTC), Members: []string{"Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"}},
		}

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
				card := createCard(artist)
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

	// Disposition widget
	content := container.NewVBox(
		searchBar,
		searchButton,
		searchResults,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(1080, 720)) // ajustement size
	myWindow.ShowAndRun()
}

func createCard(artist Artist) fyne.CanvasObject {
	image := canvas.NewImageFromFile(artist.Image)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(100, 100))
	image.Resize(fyne.NewSize(100, 100))

	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	yearLabel := widget.NewLabel(fmt.Sprintf("Year Started: %d", artist.YearStarted))
	debutLabel := widget.NewLabel(fmt.Sprintf("Debut Album: %s", artist.DebutAlbum.Format("2006-01-02")))
	membersLabel := widget.NewLabel(fmt.Sprintf("Members: %s", strings.Join(artist.Members, ", ")))

	// Ajouter des espaces entre les éléments
	spacer := layout.NewSpacer()

	content := container.NewVBox(
		layout.NewSpacer(),
		nameLabel,
		spacer,
		yearLabel,
		spacer,
		debutLabel,
		spacer,
		membersLabel,
		layout.NewSpacer(),
	)
	content.Resize(fyne.NewSize(200, 200))

	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), image, content)
	cardContent.Resize(fyne.NewSize(200, 200))

	// Créer un rectangle pour le contour avec des coins arrondis
	border := canvas.NewRectangle(color.Transparent) // Définir une couleur transparente pour le remplissage
	border.SetMinSize(fyne.NewSize(200, 200))
	border.Resize(fyne.NewSize(196, 196)) // Redimensionner légèrement la bordure pour inclure les coins arrondis
	border.StrokeColor = color.Black      // Définir la couleur de la bordure
	border.StrokeWidth = 3                // Définir l'épaisseur de la bordure
	border.CornerRadius = 20              // Définir les coins arrondis

	// Ajouter le rectangle de contour à la carte
	cardContent.Add(border)

	return cardContent
}
