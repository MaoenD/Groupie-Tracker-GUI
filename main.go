package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")

	// Champ de recherche
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Artists...")

	// Zone d'affichage des résultats de recherche
	searchResults := widget.NewLabel("Search results appear here")

	// Bouton de recherche
	searchButton := widget.NewButton("Search", func() {

		// Récupérer le texte de la recherche
		searchText := searchBar.Text

		// Comparer avec une liste dans un tableau
		artists := []string{"artist1", "artist2", "artist3"} // Remplacez ceci par votre propre liste d'artistes

		var found bool
		for _, artist := range artists {
			if searchText == artist {
				found = true
				break
			}
		}

		// Afficher les résultats
		if found {
			searchResults.SetText("Artist found: " + searchText)
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
