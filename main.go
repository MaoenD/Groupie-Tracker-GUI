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
		searchResults.SetText("Results for: " + searchBar.Text)
	})

	// Disposition widget
	content := container.NewVBox(
		searchBar,
		searchButton,
		searchResults,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 300)) // ajustement size
	myWindow.ShowAndRun()
}
