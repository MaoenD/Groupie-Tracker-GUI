package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type artistInformation struct {
	Artist string
	Title  string
	Year   int
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")

	// Champ de recherche
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Artists...")

	// Zone d'affichage des résultats de recherche
	searchResults := widget.NewLabel("Search results appear here")

	// Supposons que allArtiste soit une slice d'albums
	allArtiste := []artistInformation{
		{Artist: "Michael Jackson", Title: "Thriller", Year: 1982},
		{Artist: "Michael Jackson", Title: "Billie Jean", Year: 1982},
		{Artist: "Queen", Title: "Bohemian Rhapsody", Year: 1975},
		{Artist: "Queen", Title: "We Will Rock You", Year: 1977},
		{Artist: "AC/DC", Title: "Back in Black", Year: 1980},
		{Artist: "AC/DC", Title: "Highway to Hell", Year: 1979},
	}

	// Créer un dictionnaire pour stocker les albums par artiste
	albumsByArtist := make(map[string][]artistInformation)

	// Remplir le dictionnaire avec les albums
	for _, album := range allArtiste {
		// Si l'artiste est déjà présent dans le dictionnaire
		if albums, found := albumsByArtist[album.Artist]; found {
			// Ajouter le nouvel album à la liste existante
			albumsByArtist[album.Artist] = append(albums, album)
		} else {
			// Créer une nouvelle liste avec le nouvel album
			albumsByArtist[album.Artist] = []artistInformation{album}
		}
	}

	// Bouton de recherche
	searchButton := widget.NewButton("Search", func() {

		// Récupérer le texte de la recherche
		searchText := searchBar.Text

		// Vérifier si l'artiste est trouvé dans le dictionnaire
		albums, found := albumsByArtist[searchText]

		// Afficher les résultats
		if found {
			resultText := fmt.Sprintf("Artist found: %s\n\nAlbums:\n", searchText)
			for _, album := range albums {
				resultText += fmt.Sprintf("• %s (%d)\n", album.Title, album.Year)
			}
			searchResults.SetText(resultText)
		} else {
			searchResults.SetText("Artist not found: " + searchText)
		}
	})

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
