package main

import (
	"fmt"
	"groupie-tracker-gui/Functions"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

/********************************************************************************/
/************************************* MAIN *************************************/
/********************************************************************************/

func main() {
	// Créer une nouvelle application
	myApp := app.New()

	// Créer une nouvelle fenêtre avec le titre "Menu - Groupie Tracker"
	myWindow := myApp.NewWindow("Menu - Groupie Tracker")

	// Charger l'icône de l'application
	logoApp, _ := fyne.LoadResourceFromPath("public/img/logo.png")

	// Définir l'icône de la fenêtre
	myWindow.SetIcon(logoApp)

	// Charger les données des artistes, des lieux et des relations
	artists, err := Functions.LoadArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatalf("Failed to load artists: %v", err)
	}
	fmt.Printf("Loaded %d artists\n", len(artists))

	locations, err := Functions.LoadLocations("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		log.Fatalf("Failed to load locations: %v", err)
	}
	fmt.Printf("Loaded locations for %d artists\n", len(locations))

	relations, err := Functions.LoadRelations("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Fatalf("Failed to load relations: %v", err)
	}
	fmt.Printf("Loaded relations for %d artists\n", len(relations))

	// Créer une zone de recherche avec un champ de texte
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Artists...")
	searchBar.Resize(fyne.NewSize(1080, searchBar.MinSize().Height))

	// Créer une boîte pour afficher les résultats de recherche
	searchResults := container.NewVBox()

	// Créer une boîte pour contenir les artistes
	artistsContainer := container.NewVBox()

	// Créer un conteneur de défilement pour les artistes
	scrollContainer := container.NewVScroll(artistsContainer)
	scrollContainer.SetMinSize(fyne.NewSize(1080, 720))

	// Créer un bouton de recherche
	searchButton := widget.NewButton("Search", func() {
		// Exécuter la fonction de recherche avec les paramètres appropriés
		Functions.Recherche(searchBar, artistsContainer, Functions.Artists, myApp)

		// Effacer le texte de la zone de recherche après la recherche
		searchBar.SetText("")
	})

	// Créer une étiquette pour afficher le nombre de résultats de recherche
	searchResultCountLabel := widget.NewLabel("")

	// Créer un bouton pour afficher le logo
	logoButton := widget.NewButtonWithIcon("", (Functions.LoadImageResource("public/img/logo.png")), func() {
		// Rafraîchir le contenu de la recherche
		Functions.RefreshContent(searchBar, searchResultCountLabel, artistsContainer, Functions.Artists, myApp)
	})

	// Créer un bouton pour filtrer les résultats de recherche
	filterButton := widget.NewButton("Filtrer", func() {
		// Exécuter la fonction de filtrage
		Functions.Filter(myApp)
	})

	// Créer un conteneur pour organiser la zone de recherche et les boutons associés
	searchBarContainer := container.NewVBox(
		container.NewBorder(nil, nil, logoButton, filterButton, searchBar),
		searchButton,
		searchResultCountLabel,
	)
	searchBarContainer.Resize(searchBarContainer.MinSize())

	// Définir l'action à effectuer lorsque la touche "Entrée" est pressée dans la zone de recherche
	searchBar.OnSubmitted = func(_ string) {
		searchButton.OnTapped()
	}

	// Définir l'action à effectuer lorsque le contenu de la zone de recherche change
	searchBar.OnChanged = func(text string) {
		// Générer des suggestions de recherche basées sur le texte saisi
		count := Functions.GenerateSearchSuggestions(text, searchResults, Functions.Artists, myApp, 5)

		// Mettre à jour l'étiquette de comptage des résultats de recherche
		if count != 0 {
			searchResultCountLabel.SetText(fmt.Sprintf("Results for '%s':", text))
		} else {
			searchResultCountLabel.SetText("No result")
		}
	}

	// Organiser les artistes en cartes dans des conteneurs de lignes et de colonnes
	for i := 0; i < len(Functions.Artists); i += 3 {
		rowContainer := container.NewHBox()
		columnContainer := container.NewVBox()

		space := widget.NewLabel("")

		rowContainer.Add(space)
		rowContainer.Add(space)
		rowContainer.Add(space)

		for j := i; j < i+3 && j < len(Functions.Artists); j++ {
			card := Functions.CreateCardGeneralInfo(Functions.Artists[j], myApp)
			rowContainer.Add(card)

			if j < i+2 && j < len(Functions.Artists) {
				rowContainer.Add(space)
			}
		}

		columnContainer.Add(rowContainer)
		artistsContainer.Add(columnContainer)
	}

	// Créer le contenu de bloc
	blockContent := Functions.CreateBlockContent()

	// Créer le contenu de la fenêtre
	content := container.NewVBox(
		searchBarContainer,
		searchResults,
		blockContent,
		scrollContainer,
	)

	// Définir la disposition dynamique pour le contenu
	dynamicLayout := layout.NewVBoxLayout()
	content.Layout = dynamicLayout

	// Centrer le contenu dans un conteneur
	centeredContent := container.New(layout.NewCenterLayout(), content)

	// Créer un arrière-plan rectangulaire
	background := canvas.NewRectangle(color.NRGBA{R: 0x5C, G: 0x64, B: 0x73, A: 0xFF})
	background.Resize(fyne.NewSize(1080, 720))

	// Créer un conteneur pour l'arrière-plan avec une disposition de bordure
	backgroundContainer := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background)

	// Ajouter le contenu centré à l'arrière-plan
	backgroundContainer.Add(centeredContent)

	// Définir le contenu de la fenêtre
	myWindow.SetContent(backgroundContainer)
	myWindow.Resize(fyne.NewSize(1080, 720))

	// Afficher la fenêtre et exécuter l'application
	myWindow.ShowAndRun()
}
