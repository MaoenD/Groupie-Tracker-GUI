package Functions

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var creationDateRange *widget.Slider

func CreateBlockContent() fyne.CanvasObject {
	// Chemin de l'image à charger
	imagePath := "public/img/world_map1.jpg"

	// Charger l'image depuis le chemin spécifié
	image := canvas.NewImageFromFile(imagePath)

	// Vérifier si l'image a été chargée avec succès
	if image == nil {
		// Afficher un message d'erreur si l'image n'a pas pu être chargée
		fmt.Println("Impossible de charger l'image:", imagePath)
		return nil
	}

	// Définir le mode de remplissage de l'image pour qu'elle s'étende pour remplir l'espace
	image.FillMode = canvas.ImageFillStretch

	// Créer un conteneur pour organiser l'image et le texte avec une disposition de bordure
	blockContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil),
		image,
	)

	// Créer un bouton vide pour ajouter une interaction (à remplir selon la logique spécifique)
	button := widget.NewButton("", func() {
		// Action à effectuer lorsque le bouton est cliqué (à remplir selon les besoins)
		fmt.Print("Logique de map à intégrer ici")
	})
	button.Importance = widget.LowImportance // Définir l'importance du bouton comme faible
	button.Resize(image.MinSize())           // Redimensionner le bouton pour qu'il ait la même taille que l'image

	// Ajouter le bouton au contenu du bloc
	blockContent.Add(button)

	// Créer des étiquettes de titre et de description pour le contenu
	title := widget.NewLabel("Geolocation feature")
	description := widget.NewLabel("Find out where and when your favorite artists performed around the globe.")
	description.Wrapping = fyne.TextWrapWord // Activer le wrapping du texte pour la description

	// Créer un conteneur pour organiser les étiquettes de texte avec une disposition verticale
	textContainer := container.New(layout.NewVBoxLayout(),
		title,
		description,
	)

	// Ajouter le conteneur de texte au contenu du bloc
	blockContent.Add(textContainer)

	// Retourner le contenu du bloc
	return blockContent
}

func RefreshContent(searchBar *widget.Entry, searchResultCountLabel *widget.Label, artistsContainer *fyne.Container, artists []Artist, myApp fyne.App) {
	// Réinitialiser le texte de la barre de recherche
	searchBar.SetText("")

	// Effacer tous les objets existants dans le conteneur des artistes
	artistsContainer.Objects = nil

	// Parcourir les artistes et les organiser en cartes dans des conteneurs
	for i := 0; i < len(artists); i += 3 {
		rowContainer := container.NewHBox()    // Créer un conteneur de ligne horizontale pour les cartes d'artiste
		columnContainer := container.NewVBox() // Créer un conteneur de colonne verticale pour les lignes de cartes

		space := widget.NewLabel("") // Créer un espace vide pour l'espacement entre les cartes

		// Ajouter des espaces entre les cartes pour l'espacement visuel
		rowContainer.Add(space)
		rowContainer.Add(space)
		rowContainer.Add(space)

		// Parcourir les artistes pour créer les cartes d'artiste dans la ligne actuelle
		for j := i; j < i+3 && j < len(artists); j++ {
			card := CreateCardGeneralInfo(artists[j], myApp) // Créer une carte d'artiste pour l'artiste actuel
			rowContainer.Add(card)                           // Ajouter la carte à la ligne

			// Ajouter un espace entre les cartes si ce n'est pas la dernière carte dans la ligne
			if j < i+2 && j < len(artists) {
				rowContainer.Add(space)
			}
		}

		columnContainer.Add(rowContainer)     // Ajouter la ligne de cartes au conteneur de colonne
		artistsContainer.Add(columnContainer) // Ajouter le conteneur de colonne au conteneur des artistes
	}

	artistsContainer.Refresh() // Rafraîchir le conteneur des artistes pour afficher les modifications

	// Réinitialiser le texte du label de comptage des résultats de recherche
	searchResultCountLabel.SetText("")
}

func Filter(myApp fyne.App) {
	// Vérifie si la fenêtre des filtres est déjà ouverte
	if myWindow != nil {
		// Ferme la fenêtre des filtres si elle est ouverte
		myWindow.Close()
	}
	// Initialise les filtres de l'application
	initializeFilters(myApp)
}

func initializeFilters(myApp fyne.App) {
	// Initialisation des valeurs minimales et maximales pour les années de création et de sortie des premiers albums
	minCreationYear = Artists[0].YearStarted
	maxCreationYear = Artists[0].YearStarted
	minFirstAlbumYear := Artists[0].FirstAlbum.Year()
	maxFirstAlbumYear = Artists[0].FirstAlbum.Year()

	// Initialisation des emplacements de concerts disponibles
	concertLocations = make([]string, 0)
	locationsMap := make(map[string]bool)
	for _, artist := range Artists {
		// Mise à jour des valeurs minimales et maximales pour les années de création
		if artist.YearStarted < minCreationYear {
			minCreationYear = artist.YearStarted
		}
		if artist.YearStarted > maxCreationYear {
			maxCreationYear = artist.YearStarted
		}
		// Mise à jour des valeurs minimales et maximales pour les années du premier album
		if year := artist.FirstAlbum.Year(); year < minFirstAlbumYear {
			minFirstAlbumYear = year
		}
		if year := artist.FirstAlbum.Year(); year > maxFirstAlbumYear {
			maxFirstAlbumYear = year
		}
		// Recherche des emplacements de concerts uniques
		for _, concert := range artist.NextConcerts {
			if _, found := locationsMap[concert.Location]; !found {
				concertLocations = append(concertLocations, concert.Location)
				locationsMap[concert.Location] = true
			}
		}
	}

	// Création des widgets pour les filtres
	creationDateRange = widget.NewSlider(float64(minCreationYear), float64(maxCreationYear))
	firstAlbumDateRange = widget.NewSlider(float64(minFirstAlbumYear), float64(maxFirstAlbumYear))
	radioSoloGroup = widget.NewRadioGroup([]string{"Solo", "Group"}, func(selected string) {
		if selected == "Group" {
			numMembersBox.Show()
		} else {
			numMembersBox.Hide()
		}
	})
	numMembersCheck = widget.NewCheckGroup([]string{"2", "3", "4", "5", "6+"}, func(selected []string) {
		selectedNumMembers = selected
	})

	// Configuration de la sélection du nombre de membres
	numMembersBox = container.NewHBox()
	for _, option := range []string{"2", "3", "4", "5", "6+"} {
		option := option
		check := widget.NewCheck(option, func(checked bool) {
			selected := numMembersCheck.Selected
			if checked {
				selected = append(selected, option)
			} else {
				for i, val := range selected {
					if val == option {
						selected = append(selected[:i], selected[i+1:]...)
						break
					}
				}
			}
			numMembersCheck.SetSelected(selected)
		})
		for _, selectedOption := range selectedNumMembers {
			if selectedOption == option {
				check.SetChecked(true)
				break
			}
		}
		numMembersBox.Add(check)
	}

	// Masquage de la sélection du nombre de membres par défaut
	numMembersBox.Hide()

	// Sélection des emplacements de concerts disponibles
	locationsSelect = widget.NewSelect(concertLocations, func(selected string) {})

	// Sélection des valeurs initiales pour les filtres
	radioSoloGroup.SetSelected(selectedRadioValue)
	numMembersCheck.SetSelected(selectedNumMembers)
	locationsSelect.SetSelected(selectedLocationValue)
	numMembersCheck.SetSelected(selectedNumMembers)
	creationDateRange.SetValue(savedCreationRange)
	firstAlbumDateRange.SetValue(savedFirstAlbumRange)
	numMembersCheck.SetSelected(selectedNumMembers)

	// Sauvegarde des valeurs initiales des filtres
	savedCreationRange = creationDateRange.Value
	savedFirstAlbumRange = firstAlbumDateRange.Value
	savedNumMembers = selectedNumMembers

	// Création de la fenêtre des filtres
	myWindow = myApp.NewWindow("Groupie Tracker GUI Filters")
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.SetFixedSize(true)

	// Configuration des boutons de réinitialisation et d'application des filtres
	reset := widget.NewButton("Reset Filters", func() {
		selectedRadioValue = ""
		selectedNumMembers = nil
		selectedLocationValue = ""

		creationDateRange.SetValue(float64(minCreationYear))
		firstAlbumDateRange.SetValue(float64(minFirstAlbumYear))

		radioSoloGroup.SetSelected(selectedRadioValue)
		numMembersCheck.SetSelected(selectedNumMembers)
		locationsSelect.SetSelected(selectedLocationValue)
		applyFilter()
	})

	applyButton := widget.NewButton("Apply Filters", func() {
		applyFilter()
		myWindow.Close()
	})

	// Création des étiquettes pour les plages de dates de création et de sortie du premier album
	creationDateRangeLabel := widget.NewLabel(fmt.Sprintf("Creation Date Range: %d - %d", minCreationYear, maxCreationYear))
	firstAlbumDateRangeLabel := widget.NewLabel(fmt.Sprintf("First Album Date Range: %d - %d", minFirstAlbumYear, maxFirstAlbumYear))

	// Fonction de mise à jour des étiquettes des plages de dates
	updateLabels := func() {
		creationRange := int(creationDateRange.Value)
		firstAlbumRange := int(firstAlbumDateRange.Value)

		creationDateRangeLabel.SetText(fmt.Sprintf("Creation Date Range: %d - %d", creationRange, maxCreationYear))
		firstAlbumDateRangeLabel.SetText(fmt.Sprintf("First Album Date Range: %d - %d", firstAlbumRange, maxFirstAlbumYear))
	}

	// Mise à jour des étiquettes lors du changement de valeur des sliders
	creationDateRange.OnChanged = func(value float64) {
		updateLabels()
	}

	firstAlbumDateRange.OnChanged = func(value float64) {
		updateLabels()
	}

	// Création du conteneur des widgets de filtres
	filtersContainer := container.NewVBox(
		reset,
		creationDateRangeLabel, creationDateRange,
		firstAlbumDateRangeLabel, firstAlbumDateRange,
		radioSoloGroup,
		numMembersBox,
		locationsSelect,
		applyButton,
	)

	// Configuration de la fenêtre principale avec les widgets de filtres
	myWindow.SetContent(filtersContainer)
	myWindow.CenterOnScreen()
	myWindow.Show()
	windowOpened = true
}

func applyFilter() saveFilter {
	// Stocker les valeurs sélectionnées dans les variables correspondantes
	selectedRadioValue = radioSoloGroup.Selected
	selectedLocationValue = locationsSelect.Selected

	// Enregistrer les membres sélectionnés dans savedNumMembers
	savedNumMembers = selectedNumMembers

	// Enregistrer les plages de dates sélectionnées dans savedCreationRange et savedFirstAlbumRange
	savedCreationRange = creationDateRange.Value
	savedFirstAlbumRange = firstAlbumDateRange.Value

	// Enregistrer les valeurs sélectionnées dans savedFilter
	savedFilter = saveFilter{
		RadioSelected:      selectedRadioValue,
		NumMembersSelected: selectedNumMembers,
		LocationSelected:   selectedLocationValue,
		CreationRange:      savedCreationRange,
		FirstAlbumRange:    savedFirstAlbumRange,
	}

	// Afficher les valeurs sélectionnées dans la console
	fmt.Printf("Radio sélectionné: %s, Membres sélectionnés: %v, Localisation sélectionnée: %s, savedCreationRange: %f, savedFirstAlbumRange: %f\n", selectedRadioValue, selectedNumMembers, selectedLocationValue, savedCreationRange, savedFirstAlbumRange)

	// Réinitialiser la sélection de l'emplacement lors de l'application du filtre
	selectedLocationValue = "" // Réinitialisation de la sélection de l'emplacement

	return savedFilter // Retourner les filtres sauvegardés
}

func SecondPage(artist Artist, myApp fyne.App) {
	// Créer une nouvelle fenêtre pour afficher les informations de l'artiste
	myWindow := myApp.NewWindow("Information - " + artist.Name)

	// Définir l'icône de la fenêtre avec le logo de l'application
	logo, _ := fyne.LoadResourceFromPath("public/img/logo.png")
	myWindow.SetIcon(logo)

	// Calculer la couleur moyenne de l'image de l'artiste
	averageColor := getAverageColor(artist.Image)

	// Créer un rectangle de fond avec la couleur moyenne calculée
	background := canvas.NewRectangle(averageColor)
	background.SetMinSize(fyne.NewSize(300, 300))
	background.Resize(fyne.NewSize(296, 296))

	// Charger l'image de l'artiste et la redimensionner pour l'affichage
	image := canvas.NewImageFromFile(artist.Image)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(320, 320))
	image.Resize(fyne.NewSize(220, 220))

	// Créer des étiquettes pour afficher les informations de l'artiste
	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	yearLabel := widget.NewLabel(fmt.Sprintf("Year Started: %d", artist.YearStarted))
	debutAlbumLabel := widget.NewLabel(fmt.Sprintf("Debut Album: %s", artist.FirstAlbum.Format("02-Jan-2006")))
	membersLabel := widget.NewLabel(fmt.Sprintf("Members: %s", strings.Join(artist.Members, ", ")))
	lastConcertLabel := widget.NewLabel(fmt.Sprintf("Last Concert: %s - %s", artist.LastConcert.Date.Format("02-Jan-2006"), artist.LastConcert.Location))
	nextConcertLabel := widget.NewLabel("Next Concert:")

	// Vérifier s'il y a un prochain concert et l'afficher s'il y en a un
	if len(artist.NextConcerts) > 0 {
		nextConcertLabel.Text += fmt.Sprintf(" %s - %s", artist.NextConcerts[0].Date.Format("02-Jan-2006"), artist.NextConcerts[0].Location)
	} else {
		nextConcertLabel.Text += " No upcoming concerts" // Affichage si aucun événement à venir
	}

	// Aligner les étiquettes au centre
	nameLabel.Alignment = fyne.TextAlignCenter
	yearLabel.Alignment = fyne.TextAlignCenter
	debutAlbumLabel.Alignment = fyne.TextAlignCenter
	membersLabel.Alignment = fyne.TextAlignCenter
	lastConcertLabel.Alignment = fyne.TextAlignCenter
	nextConcertLabel.Alignment = fyne.TextAlignCenter

	// Créer un conteneur pour organiser les informations de l'artiste
	infoContainer := container.NewVBox(
		image,            // Ajout de l'image
		nameLabel,        // Ajout du nom
		yearLabel,        // Ajout de l'année de commencement
		debutAlbumLabel,  // Ajout de la date de l'album
		membersLabel,     // Ajout des noms des membres
		lastConcertLabel, // Ajout de la date du dernier concert
		nextConcertLabel, // Ajout du label du prochain concert
	)

	// Définir une taille fixe pour le conteneur d'informations
	infoContainer.Resize(fyne.NewSize(300, 200))

	// Créer un conteneur pour la carte de l'artiste avec le rectangle de fond et le conteneur d'informations
	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background, infoContainer)
	cardContent.Resize(fyne.NewSize(300, 300))

	// Définir le contenu de la fenêtre et l'afficher au centre de l'écran
	myWindow.SetContent(cardContent)
	myWindow.CenterOnScreen()
	myWindow.Show()
}
