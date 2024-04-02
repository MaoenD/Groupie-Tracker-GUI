package main

/********************************************************************************/
/*********************************** IMPORTS ************************************/
/********************************************************************************/
import ( //importation des bibliotheques nécessaires

	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

/********************************************************************************/
/************************************ TYPES *************************************/
/********************************************************************************/
type Concert struct { // Définition de la struct "Concert"
	Date     time.Time
	Location string
}

type Artist struct { // Définition de la struct "Artist"
	Name         string
	Image        string
	YearStarted  int
	DebutAlbum   time.Time
	Members      []string
	LastConcert  Concert
	NextConcerts []Concert
	Favorite     bool
	Type         string
}

type saveFilter struct {
	RadioSelected      string   // RadioSelected stocke le type d'artiste sélectionné (Solo ou Groupe).
	NumMembersSelected []string // NumMembersSelected contient les nombres de membres sélectionnés.
	LocationSelected   string   // LocationSelected stocke l'emplacement de concert sélectionné.
	CreationRange      float64  // CreationRange stocke la plage de dates de création sélectionnée.
	FirstAlbumRange    float64  // FirstAlbumRange stocke la plage de dates du premier album sélectionnée.
}

/********************************************************************************/
/********************************* VARIABLES ************************************/
/********************************************************************************/

var artists = []Artist{ // Définir les données des artistes (de façon statique pour les test)
	{Name: "Michael Jackson", Image: "public/michaeljackson.jpg", YearStarted: 1964, DebutAlbum: time.Date(1972, time.November, 13, 0, 0, 0, 0, time.UTC), Members: []string{"Michael Jackson"}, LastConcert: Concert{Date: time.Date(2009, time.June, 24, 0, 0, 0, 0, time.UTC), Location: "O2 Arena, London, UK"}, NextConcerts: []Concert{{Date: time.Date(2024, time.April, 15, 0, 0, 0, 0, time.UTC), Location: "Madison Square Garden, New York, USA"}, {Date: time.Date(2024, time.July, 10, 0, 0, 0, 0, time.UTC), Location: "Stade de France, Paris, France"}}},
	{Name: "Queen", Image: "public/queen.jpg", YearStarted: 1970, DebutAlbum: time.Date(1973, time.July, 13, 0, 0, 0, 0, time.UTC), Members: []string{"Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"}, LastConcert: Concert{Date: time.Date(2022, time.December, 15, 0, 0, 0, 0, time.UTC), Location: "The O2 Arena, London, UK"}, NextConcerts: []Concert{{Date: time.Date(2024, time.May, 20, 0, 0, 0, 0, time.UTC), Location: "Wembley Stadium, London, UK"}, {Date: time.Date(2024, time.September, 5, 0, 0, 0, 0, time.UTC), Location: "Los Angeles Memorial Coliseum, Los Angeles, USA"}}},
	{Name: "Pink Floyd", Image: "public/pinkfloyd.jpeg", YearStarted: 1965, DebutAlbum: time.Date(1967, time.August, 5, 0, 0, 0, 0, time.UTC), Members: []string{"Syd Barrett", "Roger Waters", "Richard Wright", "Nick Mason"}, LastConcert: Concert{Date: time.Date(1994, time.October, 29, 0, 0, 0, 0, time.UTC), Location: "Earls Court Exhibition Centre, London, UK"}, NextConcerts: []Concert{{Date: time.Date(2024, time.October, 15, 0, 0, 0, 0, time.UTC), Location: "Royal Albert Hall, London, UK"}, {Date: time.Date(2024, time.November, 20, 0, 0, 0, 0, time.UTC), Location: "Madison Square Garden, New York, USA"}}},
	{Name: "The Beatles", Image: "public/thebeatles.jpg", YearStarted: 1960, DebutAlbum: time.Date(1963, time.March, 22, 0, 0, 0, 0, time.UTC), Members: []string{"John Lennon", "Paul McCartney", "George Harrison", "Ringo Starr"}, LastConcert: Concert{Date: time.Date(1969, time.August, 29, 0, 0, 0, 0, time.UTC), Location: "Candlestick Park, San Francisco, USA"}, NextConcerts: []Concert{{Date: time.Date(2024, time.February, 25, 0, 0, 0, 0, time.UTC), Location: "Tokyo Dome, Tokyo, Japan"}, {Date: time.Date(2024, time.May, 5, 0, 0, 0, 0, time.UTC), Location: "Sydney Opera House, Sydney, Australia"}}},
	{Name: "Elvis Presley", Image: "public/elvispresley.jpg", YearStarted: 1954, DebutAlbum: time.Date(1956, time.March, 23, 0, 0, 0, 0, time.UTC), Members: []string{"Elvis Presley"}, LastConcert: Concert{Date: time.Date(1977, time.June, 26, 0, 0, 0, 0, time.UTC), Location: "Market Square Arena, Indianapolis, USA"}, NextConcerts: []Concert{{Date: time.Date(2024, time.September, 2, 0, 0, 0, 0, time.UTC), Location: "MGM Grand Garden Arena, Las Vegas, USA"}, {Date: time.Date(2024, time.November, 30, 0, 0, 0, 0, time.UTC), Location: "O2 Arena, London, UK"}}},
	{Name: "The Rolling Stones", Image: "public/therollingstones.jpg", YearStarted: 1962, DebutAlbum: time.Date(1964, time.April, 17, 0, 0, 0, 0, time.UTC), Members: []string{"Mick Jagger", "Keith Richards", "Charlie Watts", "Ronnie Wood"}, LastConcert: Concert{Date: time.Date(2021, time.August, 30, 0, 0, 0, 0, time.UTC), Location: "Ford Field, Detroit, USA"}, NextConcerts: []Concert{{Date: time.Date(2024, time.June, 5, 0, 0, 0, 0, time.UTC), Location: "Lambeau Field, Green Bay, USA"}, {Date: time.Date(2024, time.August, 12, 0, 0, 0, 0, time.UTC), Location: "Soldier Field, Chicago, USA"}}},
	{Name: "Led Zeppelin", Image: "public/ledzeppelin.jpg", YearStarted: 1968, DebutAlbum: time.Date(1969, time.January, 12, 0, 0, 0, 0, time.UTC), Members: []string{"Robert Plant", "Jimmy Page", "John Paul Jones", "John Bonham"}, LastConcert: Concert{Date: time.Date(2007, time.December, 10, 0, 0, 0, 0, time.UTC), Location: "02 Arena, London, UK"}, NextConcerts: []Concert{{Date: time.Date(2024, time.July, 20, 0, 0, 0, 0, time.UTC), Location: "Wembley Stadium, London, UK"}, {Date: time.Date(2024, time.October, 5, 0, 0, 0, 0, time.UTC), Location: "Stade de France, Paris, France"}}},
	{Name: "AC/DC", Image: "public/acdc.jpg", YearStarted: 1973, DebutAlbum: time.Date(1975, time.February, 17, 0, 0, 0, 0, time.UTC), Members: []string{"Angus Young", "Brian Johnson", "Phil Rudd", "Cliff Williams", "Stevie Young"}, LastConcert: Concert{Date: time.Date(2016, time.September, 20, 0, 0, 0, 0, time.UTC), Location: "Verizon Center, Washington D.C., USA"}, NextConcerts: []Concert{{Date: time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC), Location: "Wells Fargo Center, Philadelphia, USA"}, {Date: time.Date(2024, time.June, 22, 0, 0, 0, 0, time.UTC), Location: "Etihad Stadium, Manchester, UK"}}},
	{Name: "Nirvana", Image: "public/nirvana.jpg", YearStarted: 1987, DebutAlbum: time.Date(1989, time.June, 15, 0, 0, 0, 0, time.UTC), Members: []string{"Kurt Cobain", "Krist Novoselic", "Dave Grohl"}, LastConcert: Concert{Date: time.Date(1994, time.March, 1, 0, 0, 0, 0, time.UTC), Location: "Terminal 1, Munich Airport, Munich, Germany"}, NextConcerts: []Concert{{Date: time.Date(2024, time.August, 8, 0, 0, 0, 0, time.UTC), Location: "Wembley Stadium, London, UK"}, {Date: time.Date(2024, time.October, 12, 0, 0, 0, 0, time.UTC), Location: "Tokyo Dome, Tokyo, Japan"}}},
	{Name: "The Beach Boys", Image: "public/thebeachboys.jpg", YearStarted: 1961, DebutAlbum: time.Date(1962, time.October, 1, 0, 0, 0, 0, time.UTC), Members: []string{"Brian Wilson", "Mike Love", "Al Jardine", "Bruce Johnston", "David Marks"}, LastConcert: Concert{Date: time.Date(2012, time.December, 30, 0, 0, 0, 0, time.UTC), Location: "Alamodome, San Antonio, USA"}, NextConcerts: []Concert{{Date: time.Date(2024, time.May, 28, 0, 0, 0, 0, time.UTC), Location: "The SSE Hydro, Glasgow, UK"}, {Date: time.Date(2024, time.September, 15, 0, 0, 0, 0, time.UTC), Location: "Hollywood Bowl, Los Angeles, USA"}}},
	{Name: "The Who", Image: "public/thewho.jpg", YearStarted: 1964, DebutAlbum: time.Date(1965, time.December, 3, 0, 0, 0, 0, time.UTC), Members: []string{"Roger Daltrey", "Pete Townshend", "John Entwistle", "Keith Moon"}, LastConcert: Concert{Date: time.Date(2017, time.April, 1, 0, 0, 0, 0, time.UTC), Location: "The Colosseum at Caesars Palace, Las Vegas, USA"}, NextConcerts: []Concert{{Date: time.Date(2024, time.June, 30, 0, 0, 0, 0, time.UTC), Location: "PNC Music Pavilion, Charlotte, USA"}, {Date: time.Date(2024, time.September, 25, 0, 0, 0, 0, time.UTC), Location: "Bridgestone Arena, Nashville, USA"}}},
	{Name: "David Bowie", Image: "public/davidbowie.jpg", YearStarted: 1962, DebutAlbum: time.Date(1967, time.June, 1, 0, 0, 0, 0, time.UTC), Members: []string{"David Bowie"}, LastConcert: Concert{Date: time.Date(2004, time.June, 25, 0, 0, 0, 0, time.UTC), Location: "Hurricane Festival, Scheeßel, Germany"}, NextConcerts: []Concert{{Date: time.Date(2024, time.May, 10, 0, 0, 0, 0, time.UTC), Location: "Principality Stadium, Cardiff, UK"}, {Date: time.Date(2024, time.August, 20, 0, 0, 0, 0, time.UTC), Location: "Wembley Stadium, London, UK"}}},
	{Name: "Metallica", Image: "public/metallica.jpg", YearStarted: 1981, DebutAlbum: time.Date(1983, time.July, 25, 0, 0, 0, 0, time.UTC), Members: []string{"James Hetfield", "Lars Ulrich", "Kirk Hammett", "Robert Trujillo"}, LastConcert: Concert{Date: time.Date(2022, time.December, 19, 0, 0, 0, 0, time.UTC), Location: "T-Mobile Arena, Las Vegas, USA"}, NextConcerts: []Concert{{Date: time.Date(2024, time.April, 30, 0, 0, 0, 0, time.UTC), Location: "Estadio Monumental, Buenos Aires, Argentina"}, {Date: time.Date(2024, time.July, 7, 0, 0, 0, 0, time.UTC), Location: "Parque dos Atletas, Rio de Janeiro, Brazil"}}},
}

var (
	minCreationYear     int                // Année de commencement minimale parmi tous les artistes
	maxCreationYear     int                // Année de commencement maximale parmi tous les artistes
	minFirstAlbumYear   int                // Année du premier album la plus ancienne parmi tous les artistes
	maxFirstAlbumYear   int                // Année du premier album la plus récente parmi tous les artistes
	concertLocations    []string           // Liste des emplacements de concerts uniques
	creationDateRange   *widget.Slider     // Curseur pour sélectionner la plage d'années de commencement
	firstAlbumDateRange *widget.Slider     // Curseur pour sélectionner la plage d'années du premier album
	radioSoloGroup      *widget.RadioGroup // Groupe de boutons radio pour sélectionner entre Solo et Groupe
	numMembersCheck     *widget.CheckGroup // Groupe de cases à cocher pour sélectionner le nombre de membres
	numMembersBox       *fyne.Container    // Conteneur pour afficher les cases à cocher du nombre de membres
	locationsSelect     *widget.Select     // Sélecteur pour choisir l'emplacement du concert
	myWindow            fyne.Window        // Fenêtre de l'application
	windowOpened        bool               // Indique si la fenêtre est ouverte

	selectedRadioValue    string   // Valeur sélectionnée dans le groupe de boutons radio
	selectedNumMembers    []string // Valeurs sélectionnées dans le groupe de cases à cocher
	selectedLocationValue string   // Valeur sélectionnée dans le sélecteur d'emplacement
	savedCreationRange    float64  // Plage d'années de commencement sélectionnée (sauvegardée)
	savedFirstAlbumRange  float64  // Plage d'années du premier album sélectionnée (sauvegardée)
	savedNumMembers       []string // Nombre de membres sélectionnés (sauvegardé)
)

var savedFilter saveFilter

/********************************************************************************/
/************************************* MAIN *************************************/
/********************************************************************************/

func main() {
	// Créer une nouvelle application
	myApp := app.New()

	// Créer une nouvelle fenêtre avec le titre "Menu - Groupie Tracker"
	myWindow := myApp.NewWindow("Menu - Groupie Tracker")

	// Charger l'icône de l'application
	logoApp, _ := fyne.LoadResourceFromPath("public/logo.png")

	// Définir l'icône de la fenêtre
	myWindow.SetIcon(logoApp)

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
		recherche(searchBar, artistsContainer, artists, myApp)

		// Effacer le texte de la zone de recherche après la recherche
		searchBar.SetText("")
	})

	// Créer une étiquette pour afficher le nombre de résultats de recherche
	searchResultCountLabel := widget.NewLabel("")

	// Créer un bouton pour afficher le logo
	logoButton := widget.NewButtonWithIcon("", (loadImageResource("public/logo.png")), func() {
		// Rafraîchir le contenu de la recherche
		refreshContent(searchBar, searchResultCountLabel, artistsContainer, artists, myApp)
	})

	// Créer un bouton pour filtrer les résultats de recherche
	filterButton := widget.NewButton("Filtrer", func() {
		// Exécuter la fonction de filtrage
		Filter(myApp)
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
		count := generateSearchSuggestions(text, searchResults, artists, myApp, 5)

		// Mettre à jour l'étiquette de comptage des résultats de recherche
		if count != 0 {
			searchResultCountLabel.SetText(fmt.Sprintf("Results for '%s':", text))
		} else {
			searchResultCountLabel.SetText("No result")
		}
	}

	// Organiser les artistes en cartes dans des conteneurs de lignes et de colonnes
	for i := 0; i < len(artists); i += 3 {
		rowContainer := container.NewHBox()
		columnContainer := container.NewVBox()

		space := widget.NewLabel("")

		rowContainer.Add(space)
		rowContainer.Add(space)
		rowContainer.Add(space)

		for j := i; j < i+3 && j < len(artists); j++ {
			card := createCardGeneralInfo(artists[j], myApp)
			rowContainer.Add(card)

			if j < i+2 && j < len(artists) {
				rowContainer.Add(space)
			}
		}

		columnContainer.Add(rowContainer)
		artistsContainer.Add(columnContainer)
	}

	// Créer le contenu de bloc
	blockContent := createBlockContent()

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

func createBlockContent() fyne.CanvasObject {
	// Chemin de l'image à charger
	imagePath := "public/world_map1.jpg"

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

func refreshContent(searchBar *widget.Entry, searchResultCountLabel *widget.Label, artistsContainer *fyne.Container, artists []Artist, myApp fyne.App) {
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
			card := createCardGeneralInfo(artists[j], myApp) // Créer une carte d'artiste pour l'artiste actuel
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
	minCreationYear = artists[0].YearStarted
	maxCreationYear = artists[0].YearStarted
	minFirstAlbumYear = artists[0].DebutAlbum.Year()
	maxFirstAlbumYear = artists[0].DebutAlbum.Year()

	// Initialisation des emplacements de concerts disponibles
	concertLocations = make([]string, 0)
	locationsMap := make(map[string]bool)
	for _, artist := range artists {
		// Mise à jour des valeurs minimales et maximales pour les années de création
		if artist.YearStarted < minCreationYear {
			minCreationYear = artist.YearStarted
		}
		if artist.YearStarted > maxCreationYear {
			maxCreationYear = artist.YearStarted
		}
		// Mise à jour des valeurs minimales et maximales pour les années du premier album
		if year := artist.DebutAlbum.Year(); year < minFirstAlbumYear {
			minFirstAlbumYear = year
		}
		if year := artist.DebutAlbum.Year(); year > maxFirstAlbumYear {
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
	logo, _ := fyne.LoadResourceFromPath("public/logo.png")
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
	debutAlbumLabel := widget.NewLabel(fmt.Sprintf("Debut Album: %s", artist.DebutAlbum.Format("02-Jan-2006")))
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

/********************************************************************************/
/*********************************** DISPLAY ************************************/
/********************************************************************************/

func createCardGeneralInfo(artist Artist, myApp fyne.App) fyne.CanvasObject {
	// Chargement de l'image de l'artiste
	image := canvas.NewImageFromFile(artist.Image)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(120, 120))
	image.Resize(fyne.NewSize(120, 120))

	// Calcul de la couleur moyenne de l'image
	averageColor := getAverageColor(artist.Image)

	// Création du fond avec la couleur moyenne
	background := canvas.NewRectangle(averageColor)
	background.SetMinSize(fyne.NewSize(300, 300))
	background.Resize(fyne.NewSize(296, 296))
	background.CornerRadius = 20

	// Création du bouton "Plus d'informations"
	button := widget.NewButton("          Plus d'informations          ", func() {
		fmt.Println(artist.Name)
		fmt.Print("Affiche toutes les informations de l'artiste (nouvelle page)")
		SecondPage(artist, myApp)
	})

	// Création du bouton de like
	var likeButton *widget.Button
	var likeIcon string
	if artist.Favorite {
		likeIcon = "public/likeOn.ico"
	} else {
		likeIcon = "public/likeOff.ico"
	}

	likeButton = widget.NewButton("", func() {
		artist.Favorite = !artist.Favorite
		if artist.Favorite {
			likeButton.SetIcon(loadImageResource("public/likeOn.ico"))
		} else {
			likeButton.SetIcon(loadImageResource("public/likeOff.ico"))
		}
	})

	// Charger l'icône initiale du bouton en fonction de l'état initial du favori
	likeButton.SetIcon(loadImageResource(likeIcon))

	// Création du conteneur des boutons
	buttonsContainer := container.NewHBox(
		widget.NewLabel("  "),
		container.NewBorder(nil, layout.NewSpacer(), nil, likeButton, button),
	)

	// Création du label du nom de l'artiste
	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Création du label de l'année de début
	yearLabel := widget.NewLabel(fmt.Sprintf("%d", artist.YearStarted))

	// Création du conteneur pour les labels du nom et de l'année
	labelsContainer := container.NewHBox(nameLabel, yearLabel)

	// Création du texte des membres du groupe
	var membersText string
	if len(artist.Members) == 1 {
		membersText = "Artiste Solo\n"
	} else if len(artist.Members) > 0 {
		membersText = "Membres : " + strings.Join(artist.Members, ", ")
	}
	membersLabel := widget.NewLabel(membersText)
	membersLabel.Wrapping = fyne.TextWrapWord

	// Création du conteneur pour les informations de l'artiste
	infoContainer := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), image, labelsContainer, membersLabel, layout.NewSpacer(), buttonsContainer, layout.NewSpacer())
	infoContainer.Resize(fyne.NewSize(300, 180))

	// Création du contenu de la carte avec le fond et les informations de l'artiste
	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background, infoContainer)
	cardContent.Resize(fyne.NewSize(300, 300))

	// Création de la bordure autour du contenu de la carte
	border := canvas.NewRectangle(color.Transparent)
	border.SetMinSize(fyne.NewSize(300, 300))
	border.Resize(fyne.NewSize(296, 296))
	border.StrokeColor = color.Black
	border.StrokeWidth = 3
	border.CornerRadius = 20

	cardContent.Add(border)

	return cardContent
}

func generateSearchSuggestions(text string, scrollContainer *fyne.Container, artists []Artist, myApp fyne.App, limit int) int {
	// Effacer les objets précédents du conteneur de défilement
	scrollContainer.Objects = nil

	// Vérifier si le texte de recherche est vide ou s'il n'y a pas d'artistes dans la liste
	if text == "" || len(artists) == 0 {
		return 0
	}

	// Variable pour compter le nombre de suggestions affichées
	count := 0

	// Parcourir tous les artistes dans la liste
	for _, artist := range artists {
		// Vérifier si le nombre de suggestions affichées a atteint la limite spécifiée
		if count >= limit {
			break
		}

		// Vérifier si l'artiste correspond aux filtres sauvegardés
		if artistMatchesFilters(artist, savedFilter) {
			// Vérifier si le nom de l'artiste contient le texte de recherche
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(text)) {
				// Incrémenter le compteur et ajouter un bouton d'artiste au conteneur de défilement
				count++
				artistButton := widget.NewButton(artist.Name, func(a Artist) func() {
					return func() {
						SecondPage(a, myApp)
					}
				}(artist))
				artistButton.Importance = widget.LowImportance
				scrollContainer.Add(artistButton)
			} else {
				// Vérifier si le texte de recherche correspond à l'année de commencement de l'artiste
				if strconv.Itoa(artist.YearStarted) == text {
					// Incrémenter le compteur et ajouter un bouton d'artiste avec l'année de commencement au conteneur de défilement
					count++
					artistButton := widget.NewButton(artist.Name+" (Year Started: "+text+")", func(a Artist) func() {
						return func() {
							SecondPage(a, myApp)
						}
					}(artist))
					artistButton.Importance = widget.LowImportance
					scrollContainer.Add(artistButton)
				}

				// Vérifier si le texte de recherche correspond à l'année de l'album de début de l'artiste
				if strconv.Itoa(artist.DebutAlbum.Year()) == text {
					// Incrémenter le compteur et ajouter un bouton d'artiste avec l'année de l'album de début au conteneur de défilement
					count++
					artistButton := widget.NewButton(artist.Name+" (Debut Album: "+strconv.Itoa(artist.DebutAlbum.Year())+")", func(a Artist) func() {
						return func() {
							SecondPage(a, myApp)
						}
					}(artist))
					artistButton.Importance = widget.LowImportance
					scrollContainer.Add(artistButton)
				}

				// Vérifier s'il y a plus d'un membre dans le groupe et si le texte de recherche correspond au nom d'un membre
				if len(artist.Members) > 1 {
					for _, member := range artist.Members {
						if strings.Contains(strings.ToLower(member), strings.ToLower(text)) {
							// Incrémenter le compteur et ajouter un bouton d'artiste avec le nom du membre au conteneur de défilement
							count++
							artistButton := widget.NewButton(artist.Name+" (Member Name: "+member+")", func(a Artist) func() {
								return func() {
									SecondPage(a, myApp)
								}
							}(artist))
							artistButton.Importance = widget.LowImportance
							scrollContainer.Add(artistButton)
							break
						}
					}
				}

				// Vérifier si le texte de recherche correspond à un lieu de concert dans les prochains concerts de l'artiste
				for _, concert := range artist.NextConcerts {
					if strings.Contains(strings.ToLower(concert.Location), strings.ToLower(text)) {
						// Incrémenter le compteur et ajouter un bouton d'artiste avec le lieu du concert au conteneur de défilement
						count++
						artistButton := widget.NewButton(artist.Name+" (Concert Location: "+concert.Location+")", func(a Artist) func() {
							return func() {
								SecondPage(a, myApp)
							}
						}(artist))
						artistButton.Importance = widget.LowImportance
						scrollContainer.Add(artistButton)
						break
					}
				}
			}
		}
	}

	// Vérifier si le nombre de suggestions affichées est inférieur au nombre total d'artistes et si la limite a été atteinte
	if count < len(artists) && count >= limit {
		// Ajouter un bouton "Plus de résultats" pour charger davantage de résultats
		showMoreButton := widget.NewButton("More results", func() {
			generateSearchSuggestions(text, scrollContainer, artists, myApp, limit+5)
		})
		scrollContainer.Add(showMoreButton)
	}

	// Retourner le nombre total de suggestions affichées
	return count
}

func recherche(searchBar *widget.Entry, scrollContainer *fyne.Container, artists []Artist, myApp fyne.App) {
	// Convertir le texte de recherche en minuscules pour une recherche insensible à la casse
	searchText := strings.ToLower(searchBar.Text)

	// Créer un conteneur pour stocker les artistes trouvés
	artistsContainer := container.NewVBox()

	// Liste des artistes trouvés
	var foundArtists []Artist

	// Parcourir tous les artistes dans la liste
	for _, artist := range artists {
		// Vérifier si l'artiste correspond aux filtres sauvegardés
		if artistMatchesFilters(artist, savedFilter) {
			// Vérifier si le nom de l'artiste, l'année de commencement, l'année de l'album de début, le nom d'un membre ou le lieu d'un concert correspond au texte de recherche
			if strings.Contains(strings.ToLower(artist.Name), searchText) ||
				strconv.Itoa(artist.YearStarted) == searchText ||
				strconv.Itoa(artist.DebutAlbum.Year()) == searchText ||
				checkMemberName(artist.Members, searchText) ||
				checkConcertLocation(artist.NextConcerts, searchText) {
				// Ajouter l'artiste à la liste des artistes trouvés
				foundArtists = append(foundArtists, artist)
			}
		}
	}

	// Vérifier s'il y a des artistes trouvés
	if len(foundArtists) > 0 {
		// Afficher les artistes trouvés par groupe de 3 dans des conteneurs de rangées et colonnes
		for i := 0; i < len(foundArtists); i += 3 {
			rowContainer := container.NewHBox()
			columnContainer := container.NewVBox()

			space := widget.NewLabel("")

			rowContainer.Add(space)
			rowContainer.Add(space)
			rowContainer.Add(space)
			for j := i; j < i+3 && j < len(foundArtists); j++ {
				card := createCardGeneralInfo(foundArtists[j], myApp)
				rowContainer.Add(card)

				if j < i+2 && j < len(foundArtists) {
					rowContainer.Add(space)
				}
			}

			columnContainer.Add(rowContainer)
			artistsContainer.Add(columnContainer)
		}
	} else {
		// Afficher un message indiquant qu'aucun résultat n'a été trouvé
		noResultLabel := widget.NewLabel("No result found")
		artistsContainer.Add(noResultLabel)
	}

	// Mettre à jour les objets dans le conteneur de défilement avec les artistes trouvés
	scrollContainer.Objects = []fyne.CanvasObject{artistsContainer}
	scrollContainer.Refresh()
}

func artistMatchesFilters(artist Artist, filter saveFilter) bool {
	// Vérifier si l'artiste correspond aux filtres de l'utilisateur
	if filter.RadioSelected != "" {
		// Vérifier si l'utilisateur a sélectionné "Solo" et l'artiste est un groupe, ou vice versa
		if filter.RadioSelected == "Solo" && len(artist.Members) > 1 {
			return false
		} else if filter.RadioSelected == "Group" && len(artist.Members) <= 1 {
			return false
		}
	}

	// Vérifier si le nombre de membres de l'artiste correspond à l'une des sélections de l'utilisateur
	if len(filter.NumMembersSelected) > 0 && !contains(filter.NumMembersSelected, strconv.Itoa(len(artist.Members))) {
		return false
	}

	// Vérifier si l'emplacement du concert de l'artiste correspond à la sélection de l'utilisateur
	if filter.LocationSelected != "" && !artistHasConcertLocation(artist, filter.LocationSelected) {
		return false
	}

	// Vérifier si l'année de commencement de l'artiste est dans la plage sélectionnée par l'utilisateur
	if filter.CreationRange > 0 && float64(artist.YearStarted) < filter.CreationRange {
		return false
	}

	// Vérifier si l'année de sortie du premier album de l'artiste est dans la plage sélectionnée par l'utilisateur
	if filter.FirstAlbumRange > 0 && float64(artist.DebutAlbum.Year()) < filter.FirstAlbumRange {
		return false
	}
	return true
}

func artistHasConcertLocation(artist Artist, location string) bool {
	// Vérifier si l'artiste a un concert à l'emplacement spécifié
	for _, concert := range artist.NextConcerts {
		if strings.EqualFold(concert.Location, location) { // Vérifier sans tenir compte de la casse
			return true
		}
	}
	return false
}

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

func loadImageResource(path string) fyne.Resource {
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
