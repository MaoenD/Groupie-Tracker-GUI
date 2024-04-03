package Functions

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var CreationDateRange *widget.Slider

func CreateBlockContent() fyne.CanvasObject {
	// Chemin de l'image à charger
	imagePath := "public/img/world_map.jpg"

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

func RefreshContent(searchBar *widget.Entry, searchResultCountLabel *widget.Label, artistsContainer *fyne.Container, relation Relation, artists []Artist, myApp fyne.App) {
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
			card := CreateCardGeneralInfo(artists[j], relation, myApp) // Créer une carte d'artiste pour l'artiste actuel
			rowContainer.Add(card)                                     // Ajouter la carte à la ligne

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
	artists, err := LoadArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Printf("Failed to load artists: %v", err)
		return // Exit if there was an error fetching the artist data
	}

	// Fetch location data (assuming this returns data relevant for the filter, e.g., concert locations)
	concerts, err := CombineData("https://groupietrackers.herokuapp.com/api/locations", "https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Printf("Failed to load locations: %v", err)
		return // Exit if there was an error fetching the location data
	}
	// Initialise les filtres de l'application
	initializeFilters(myApp, artists, concerts)
}

func initializeFilters(myApp fyne.App, artists []Artist, concerts []Concert) {
	// Initialisation des valeurs minimales et maximales pour les années de création et de sortie des premiers albums
	minCreationYear := artists[0].CreationDate
	maxCreationYear := artists[0].CreationDate

	// On asumme que le format est "DD-MM-YYYY" et converti en int
	parseYear := func(dateStr string) int {
		parts := strings.Split(dateStr, "-")
		year, _ := strconv.Atoi(parts[2])
		return year
	}

	minFirstAlbumYear := parseYear(artists[0].FirstAlbum)
	maxFirstAlbumYear := parseYear(artists[0].FirstAlbum)

	for _, artist := range artists {
		// Mise à jour des valeurs minimales et maximales pour les années de création
		if artist.CreationDate < minCreationYear {
			minCreationYear = artist.CreationDate
		}
		if artist.CreationDate > maxCreationYear {
			maxCreationYear = artist.CreationDate
		}

		// Parsing de l'année FirstAlbum de la string date
		albumYear := parseYear(artist.FirstAlbum)
		if albumYear < minFirstAlbumYear {
			minFirstAlbumYear = albumYear
		}
		if albumYear > maxFirstAlbumYear {
			maxFirstAlbumYear = albumYear
		}
	}

	// Initialisation des emplacements de concerts disponibles
	concertLocations := make([]string, 0)
	locationsMap := make(map[string]bool)
	for _, artist := range artists {
		// Mise à jour des valeurs minimales et maximales pour les années de création
		if artist.CreationDate < minCreationYear {
			minCreationYear = artist.CreationDate
		}
		if artist.CreationDate > maxCreationYear {
			maxCreationYear = artist.CreationDate
		}
		// Mise à jour des valeurs minimales et maximales pour les années du premier album
		albumYear := parseYear(artist.FirstAlbum)
		if albumYear < minFirstAlbumYear {
			minFirstAlbumYear = albumYear
		}
		if albumYear > maxFirstAlbumYear {
			maxFirstAlbumYear = albumYear
		}
		// Recherche des emplacements de concerts uniques
		for _, concert := range concerts {
			// Concaténer les emplacements de concert en une seule chaîne
			locationStr := strings.Join(concert.Locations, ", ")
			if _, found := locationsMap[locationStr]; !found {
				concertLocations = append(concertLocations, locationStr)
				locationsMap[locationStr] = true
			}
		}
	}

	// Création des widgets pour les filtres
	CreationDateRange = widget.NewSlider(float64(minCreationYear), float64(maxCreationYear))
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
	CreationDateRange.SetValue(savedCreationRange)
	firstAlbumDateRange.SetValue(savedFirstAlbumRange)
	numMembersCheck.SetSelected(selectedNumMembers)

	// Sauvegarde des valeurs initiales des filtres
	savedCreationRange = CreationDateRange.Value
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

		CreationDateRange.SetValue(float64(minCreationYear))
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
	CreationDateRangeLabel := widget.NewLabel(fmt.Sprintf("Creation Date Range: %d - %d", minCreationYear, maxCreationYear))
	firstAlbumDateRangeLabel := widget.NewLabel(fmt.Sprintf("First Album Date Range: %d - %d", minFirstAlbumYear, maxFirstAlbumYear))

	// Fonction de mise à jour des étiquettes des plages de dates
	updateLabels := func() {
		creationRange := int(CreationDateRange.Value)
		firstAlbumRange := int(firstAlbumDateRange.Value)

		CreationDateRangeLabel.SetText(fmt.Sprintf("Creation Date Range: %d - %d", creationRange, maxCreationYear))
		firstAlbumDateRangeLabel.SetText(fmt.Sprintf("First Album Date Range: %d - %d", firstAlbumRange, maxFirstAlbumYear))
	}

	// Mise à jour des étiquettes lors du changement de valeur des sliders
	CreationDateRange.OnChanged = func(value float64) {
		updateLabels()
	}

	firstAlbumDateRange.OnChanged = func(value float64) {
		updateLabels()
	}

	// Création du conteneur des widgets de filtres
	filtersContainer := container.NewVBox(
		reset,
		CreationDateRangeLabel, CreationDateRange,
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
	savedCreationRange = CreationDateRange.Value
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

func SecondPage(artist Artist, relation Relation, myApp fyne.App) {
	myWindow := myApp.NewWindow("Information - " + artist.Name)

	logo, _ := fyne.LoadResourceFromPath("public/img/logo.png")
	myWindow.SetIcon(logo)

	response, err := http.Get(artist.Image)
	if err != nil {
		log.Println("Failed to load artist image:", err)
		return
	}
	defer response.Body.Close()

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed to read image data:", err)
		return
	}

	image := canvas.NewImageFromReader(strings.NewReader(string(imageData)), "image/jpeg")
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(320, 320))
	image.Resize(fyne.NewSize(220, 220))

	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	yearLabel := widget.NewLabel(fmt.Sprintf("Year Started: %d", artist.CreationDate))
	debutAlbumLabel := widget.NewLabel(fmt.Sprintf("Debut Album: %s", artist.FirstAlbum))
	membersLabel := widget.NewLabel(fmt.Sprintf("Members: %s", strings.Join(artist.Members, ", ")))

	// Concert information
	concertInfo := container.NewVBox()
	for location, dates := range relation.DatesLocations {
		for _, date := range dates {
			formattedLocation := formatLocation(location)
			formattedDate := formatDate(date)
			concertLabel := widget.NewLabel(fmt.Sprintf("🗺️Location: %s    📅Date: %s", formattedLocation, formattedDate))
			concertInfo.Add(concertLabel)
		}
	}

	scrollContainer := container.NewScroll(concertInfo)
	scrollContainer.SetMinSize(fyne.NewSize(400, 200))

	content := container.NewVBox(
		image,
		nameLabel,
		yearLabel,
		debutAlbumLabel,
		membersLabel,
		scrollContainer,
	)

	myWindow.SetContent(content)
	myWindow.Show()
}

func formatLocation(location string) string {
	// Supprimer les tirets et les underscores
	location = strings.ReplaceAll(location, "_", " ")
	// Mettre la première lettre de chaque mot en majuscule
	titleCase := cases.Title(language.English)
	location = titleCase.String(location)

	// Trouver le nom du pays après le premier tiret
	parts := strings.Split(location, "-")
	if len(parts) > 1 {
		// Ajouter le nom du pays entre parenthèses
		location = fmt.Sprintf("%s (%s)", parts[0], parts[1])
	}

	return location
}

// Fonction pour formater la date
func formatDate(date string) string {
	// Convertir le format "DD-MM-YYYY" en "JJ/MM/AAAA"
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return date // Retourner la date telle quelle si le format est incorrect
	}
	return fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
}
