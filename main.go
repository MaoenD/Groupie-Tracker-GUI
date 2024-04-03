package main

import (
	"fmt"
	"groupie-tracker-gui/Functions" // Import custom functions package

	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var artistRelations Functions.Relation // Declare variable outside loop

func main() {
	// Create a new application
	myApp := app.New()

	// Create a new window with the title "Menu - Groupie Tracker"
	myWindow := myApp.NewWindow("Menu - Groupie Tracker")

	// Load the application icon
	logoApp, _ := fyne.LoadResourceFromPath("public/img/logo.png")

	// Set the window icon
	myWindow.SetIcon(logoApp)

	// Load artist, location, and relation data
	artists, err := Functions.LoadArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatalf("Failed to load artists: %v", err)
	}

	relations, err := Functions.LoadRelations("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Fatalf("Failed to load relations: %v", err)
	}

	// Create a search area with a text field
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Artists...")
	searchBar.Resize(fyne.NewSize(1080, searchBar.MinSize().Height))

	// Create a box to display search results
	searchResults := container.NewVBox()

	// Create a box to contain artists
	artistsContainer := container.NewVBox()

	// Create a scroll container for artists
	scrollContainer := container.NewVScroll(artistsContainer)
	scrollContainer.SetMinSize(fyne.NewSize(1080, 720))

	// Organize artists into cards within rows and columns
	for i := 0; i < len(artists); i += 3 {
		rowContainer := container.NewHBox()
		columnContainer := container.NewVBox()

		space := widget.NewLabel("")

		rowContainer.Add(space)
		rowContainer.Add(space)
		rowContainer.Add(space)

		// Iterate over each group of 3 artists
		for j := i; j < i+3 && j < len(artists); j++ {
			// Find corresponding relation for the current artist
			for _, rel := range relations {
				if rel.ID == artists[j].ID {
					artistRelations = rel
					break
				}
			}

			// Create artist card with corresponding relation
			card := Functions.CreateCardGeneralInfo(artists[j], artistRelations, myApp)
			rowContainer.Add(card)

			if j < i+2 && j < len(artists) {
				rowContainer.Add(space)
			}
		}

		columnContainer.Add(rowContainer)
		artistsContainer.Add(columnContainer)
	}

	// Create a search button
	searchButton := widget.NewButton("Search", func() {
		// Execute search function with appropriate parameters
		Functions.Recherche(searchBar, artistsContainer, artists, artistRelations, myApp)

		// Clear search bar text after search
		searchBar.SetText("")
	})

	// Create a label to display search result count
	searchResultCountLabel := widget.NewLabel("")

	// Create a button to display the logo
	logoButton := widget.NewButtonWithIcon("", (Functions.LoadImageResource("public/img/logo.png")), func() {
		// Refresh search content
		Functions.RefreshContent(searchBar, searchResultCountLabel, artistsContainer, artistRelations, artists, myApp)
	})

	// Create a button to filter search results
	filterButton := widget.NewButton("Filter", func() {
		// Execute filtering function
		Functions.Filter(myApp)
	})

	// Create a container to organize search area and associated buttons
	searchBarContainer := container.NewVBox(
		container.NewBorder(nil, nil, logoButton, filterButton, searchBar),
		searchButton,
		searchResultCountLabel,
	)
	searchBarContainer.Resize(searchBarContainer.MinSize())

	// Define action when "Enter" key is pressed in search area
	searchBar.OnSubmitted = func(_ string) {
		searchButton.OnTapped()
	}

	// Define action when content in search area changes
	searchBar.OnChanged = func(text string) {
		// Generate search suggestions based on entered text
		count := Functions.GenerateSearchSuggestions(text, searchResults, artists, artistRelations, myApp, 5)

		// Update search result count label
		if count != 0 {
			searchResultCountLabel.SetText(fmt.Sprintf("Results for '%s':", text))
		} else {
			searchResultCountLabel.SetText("No result")
		}
	}

	// Create block content
	blockContent := Functions.CreateBlockContent()

	// Create window content
	content := container.NewVBox(
		searchBarContainer,
		searchResults,
		blockContent,
		scrollContainer,
	)

	// Set dynamic layout for content
	dynamicLayout := layout.NewVBoxLayout()
	content.Layout = dynamicLayout

	// Center content within a container
	centeredContent := container.New(layout.NewCenterLayout(), content)

	// Create a rectangular background
	background := canvas.NewRectangle(color.NRGBA{R: 0x5C, G: 0x64, B: 0x73, A: 0xFF})
	background.Resize(fyne.NewSize(1080, 720))

	// Create a container for the background with border layout
	backgroundContainer := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background)

	// Add centered content to the background
	backgroundContainer.Add(centeredContent)

	// Set window content
	myWindow.SetContent(backgroundContainer)
	myWindow.Resize(fyne.NewSize(1080, 720))

	// Show and run the window
	myWindow.ShowAndRun()
}
