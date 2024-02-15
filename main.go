package main

import ( //importation des bibliotheques nécessaires
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

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
}

func main() { // Fonction Principale, lancement de l'application
	myApp := app.New()                                    // Création de l'application
	myWindow := myApp.NewWindow("Menu - Groupie Tracker") // Création d'une fenêtre et nommage

	artists := []Artist{ // Définir les données des artistes
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

	searchBar := widget.NewEntry()                // Création d'un champs de recherche
	searchBar.SetPlaceHolder("Search Artists...") // Définir un placeholder

	searchResults := container.NewVBox() // Création de la zone d'affichage des résultats de recherches

	searchButton := widget.NewButton("Search", func() { // Création d'un bouton de recherche
		searchText := searchBar.Text // Récupérer le texte de la recherche

		var foundArtists []Artist        // Recherche d'artistes correspondants
		for _, artist := range artists { // Parcours la liste avec tous les artists
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchText)) { //Vérification de résultat avec l'entrée
				foundArtists = append(foundArtists, artist) //L'ajouté à un tableaux
			}
		}

		searchResultsObjects := make([]fyne.CanvasObject, 0) // Afficher les résultats
		if len(foundArtists) > 0 {                           // Si le tableaux n'est pas vide
			resultText := fmt.Sprintf("Artist found for: %s\n\n", searchText)                //affiche l'entrée de recherche
			searchResultsObjects = append(searchResultsObjects, widget.NewLabel(resultText)) // Ajoute au Canva afin de l'affiché
			for _, artist := range foundArtists {                                            // Affichage des informations pour les tous les artistes dans le tableau
				card := createCardAllInfo(artist)                         // Création de la card selon l'artiste trouvé
				searchResultsObjects = append(searchResultsObjects, card) // Ajoute au Canva afin de l'affiché
			}
		} else { // Si le tableaux est pas vide (aucun résultat trouvé)
			searchResultsObjects = append(searchResultsObjects, widget.NewLabel("No artist found: "+searchText)) // Ajoute au Canva afin de l'affiché
		}

		searchResults.Objects = searchResultsObjects //COMPLETE ICI
		searchResults.Refresh()                      //permet de refresh selon la recherche
	})

	searchBar.OnSubmitted = func(_ string) { // Ajouter l'écouteur d'événements clavier au champ de recherche
		searchButton.OnTapped() //COMPLETE ICI
	}

	artistsContainer := container.NewVBox() // Création du conteneur pour afficher les artistes

	for i := 0; i < len(artists); i += 3 { // Ajouter les cartes des artistes par groupes de 3 dans des conteneurs de grille
		rowContainer := container.NewGridWithColumns(3) // Créer un nouveau conteneur de grille pour chaque ligne d'artistes
		for j := i; j < i+3 && j < len(artists); j++ {  // Ajouter les trois artistes de cette ligne dans la grille
			card := createCardGeneralInfo(artists[j]) // Création des cards
			rowContainer.Add(card)                    //Ajouter les cards à la colonne
		}

		/* artistsContainer.Add(layout.NewSpacer()) // Ajouter un espacement */
		artistsContainer.Add(rowContainer) // Ajoute des cards créer dans le container
	}

	scrollContainer := container.NewVScroll(artistsContainer) // Créer un conteneur de défilement pour le conteneur principal
	scrollContainer.SetMinSize(fyne.NewSize(1080, 720))       // Taille minimale pour activer le défilement

	blockContent := mapFeature() // Création du bloc de contenu //MAP FEATURE

	content := container.NewVBox( // Création du conteneur principal avec la couleur de fond spécifiée
		searchBar,       // Ajout de la search bar
		searchButton,    // Ajout du boutton
		searchResults,   // Ajout de les résultats (si trouvés)
		blockContent,    // Ajouter le bloc de contenu //MAP FEATURE
		scrollContainer, // Utilisation du conteneur de défilement pour les artistes
	)

	centeredContent := container.New(layout.NewCenterLayout(), content) // Centrer les cartes dans la fenêtre

	background := canvas.NewRectangle(color.NRGBA{R: 0x5C, G: 0x64, B: 0x73, A: 0xFF}) // Création d'un rectangle pour le background
	background.Resize(fyne.NewSize(1080, 720))                                         // Taille pour remplir toute la fenêtre

	backgroundContainer := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background) // Créer un conteneur pour contenir le fond coloré

	backgroundContainer.Add(centeredContent) // Ajouter le contenu principal au conteneur de fond

	myWindow.SetContent(backgroundContainer) // Afficher la fenêtre
	myWindow.Resize(fyne.NewSize(1080, 720)) // ajustement size
	myWindow.ShowAndRun()                    //run la fenêtre
}

func mapFeature() fyne.CanvasObject {
	imagePath := "public/world_map1.jpg" // Chemin de l'image

	image := canvas.NewImageFromFile(imagePath) // Charger l'image

	if image == nil { // Gestion d'erreurs
		fmt.Println("Impossible de charger l'image:", imagePath)
		return nil
	}

	image.FillMode = canvas.ImageFillStretch // Définir le mode de remplissage à "Contain" pour le recadrage

	title := widget.NewLabel("Geolocation feature")                                                                      // Créer du titre
	description := widget.NewLabel("Find out where and when your favorite artists will be performing around the globe.") // Créer de la description
	/* description.Wrapping = fyne.TextWrapWord   */ // Activer le wrapping du texte

	textContainer := container.New(layout.NewVBoxLayout(), // Création d'un conteneur pour le texte
		title,
		description,
	)

	blockContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), // Créer un conteneur pour organiser l'image et le texte
		image,
		textContainer,
	)

	return blockContent
}

func createCardGeneralInfo(artist Artist) fyne.CanvasObject {
	image := canvas.NewImageFromFile(artist.Image) // Redimensionner l'image
	image.FillMode = canvas.ImageFillContain       // Gestion du fill image
	image.SetMinSize(fyne.NewSize(120, 120))       //Définir la taille minimum de l'image
	image.Resize(fyne.NewSize(120, 120))           //Définir la nouvelle image de l'image

	averageColor := getAverageColor(artist.Image) // Obtenir la couleur moyenne de l'image

	background := canvas.NewRectangle(averageColor) // Création d'un rectangle coloré pour l'arrière-plan de la card
	background.SetMinSize(fyne.NewSize(300, 300))   // Définir la taille minimum du bakcground
	background.Resize(fyne.NewSize(296, 296))       // Redimensionner pour inclure les coin
	background.CornerRadius = 20                    // Définir les coins arrondis

	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}) // Nom de l'artiste en gras et plus gros

	yearLabel := widget.NewLabel(fmt.Sprintf("%d", artist.YearStarted)) // Date de début de l'artiste en plus petit

	labelsContainer := container.NewHBox( // Créer un conteneur HBox pour afficher les labels avec un espace entre eux
		nameLabel,
		yearLabel,
	)

	var membersText string        // Gestion des membres du groupe
	if len(artist.Members) == 1 { // Contion du nombre de membres
		membersText = "Solo Artist"
	} else if len(artist.Members) > 0 { // Contion du nombre de membres
		membersText = "Members:\n " + strings.Join(artist.Members, ", ")
	}
	membersLabel := widget.NewLabel(membersText) // EXPLIQUE ICI
	membersLabel.Wrapping = fyne.TextWrapWord    // Activer le wrapping du texte

	infoContainer := container.New(layout.NewVBoxLayout(), // Créer le conteneur pour les informations sur l'artiste
		layout.NewSpacer(), // Ajout d'un espace vertical
		image,              // Ajout de l'image
		labelsContainer,    // Ajout titre et date
		membersLabel,       // Afficher les membres du groupe
		layout.NewSpacer(), // Ajout d'un petit espace vertical
	)

	infoContainer.Resize(fyne.NewSize(300, 180)) // Définir la taille fixe pour le conteneur d'informations

	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background, infoContainer) // Créer le conteneur pour la card de l'artiste
	cardContent.Resize(fyne.NewSize(300, 300))                                                          //Définir la taille minimum de la card                                               //

	border := canvas.NewRectangle(color.Transparent) // Définir une couleur transparente pour le remplissage
	border.SetMinSize(fyne.NewSize(300, 300))        //Définir la taille minimum de la bordure
	border.Resize(fyne.NewSize(296, 296))            // Redimensionner pour inclure les coin
	border.StrokeColor = color.Black                 // Définir la couleur de la bordure
	border.StrokeWidth = 3                           // Définir l'épaisseur de la bordure
	border.CornerRadius = 20                         // Définir les coins

	cardContent.Add(border) // Ajouter le rectangle de contour à la carte

	return cardContent
}

func createCardAllInfo(artist Artist) fyne.CanvasObject {
	image := canvas.NewImageFromFile(artist.Image) // Redimensionner l'image dans un canva
	image.FillMode = canvas.ImageFillContain       // Définir le Fill image
	image.SetMinSize(fyne.NewSize(120, 120))       // Définir la taille minimum de l'image
	image.Resize(fyne.NewSize(120, 120))           // Définir la nouvelle taille

	averageColor := getAverageColor(artist.Image) // Obtenir la couleur moyenne de l'image

	background := canvas.NewRectangle(averageColor) // Création d'un rectangle coloré pour l'arrière-plan de la card
	background.SetMinSize(fyne.NewSize(300, 300))   // Définir la taille minimum du bakcground
	background.Resize(fyne.NewSize(296, 296))       // Redimensionner pour inclure les coin
	background.CornerRadius = 20                    // Définir les coins arrondis

	// Créer des labels pour afficher les informations sur l'artiste
	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})                                                 // Affichage du nom
	yearLabel := widget.NewLabel(fmt.Sprintf("Year Started: %d", artist.YearStarted))                                                                     // Affichage de l'année de commencement
	debutAlbumLabel := widget.NewLabel(fmt.Sprintf("Debut Album: %s", artist.DebutAlbum.Format("02-Jan-2006")))                                           // Affichage de la date de l'album
	membersLabel := widget.NewLabel(fmt.Sprintf("Members: %s", strings.Join(artist.Members, ", ")))                                                       // Affichage des noms des artistes
	lastConcertLabel := widget.NewLabel(fmt.Sprintf("Last Concert: %s - %s", artist.LastConcert.Date.Format("02-Jan-2006"), artist.LastConcert.Location)) // Affichage de la date du dernier concert
	nextConcertLabel := widget.NewLabel("Next Concert:")                                                                                                  // Affichage de la date du prochain concert
	if len(artist.NextConcerts) > 0 {                                                                                                                     // Condtion pour afficher les dates
		nextConcertLabel.Text += fmt.Sprintf(" %s - %s", artist.NextConcerts[0].Date.Format("02-Jan-2006"), artist.NextConcerts[0].Location) // Affichage des dates et lieux
	} else {
		nextConcertLabel.Text += " No upcoming concerts" // Affichage si aucun événeement
	}

	infoContainer := container.NewVBox( // Créer un conteneur VBox pour organiser les labels verticalement
		image,            // Ajout de l'image
		nameLabel,        // Ajout du nom
		yearLabel,        // Ajout de l'année de commencement
		debutAlbumLabel,  // AJout de la date de l'album
		membersLabel,     // Ajout des noms des artites
		lastConcertLabel, // AJout de la date du dernier concert
		nextConcertLabel, // Ajout du label du prochain concert
	)

	infoContainer.Resize(fyne.NewSize(300, 200)) // Définir la taille fixe pour le conteneur d'informations

	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background, infoContainer) // Créer le conteneur pour la carte de l'artiste
	cardContent.Resize(fyne.NewSize(300, 300))                                                          // Redimensionner le conteneur

	// Créer un rectangle pour le contour avec des coins arrondis
	border := canvas.NewRectangle(color.Transparent) // Définir une couleur transparente pour le remplissage
	border.SetMinSize(fyne.NewSize(300, 300))        // Définir la taille minimum
	border.Resize(fyne.NewSize(296, 296))            // Redimensionner légèrement la bordure pour inclure les coins arrondis
	border.StrokeColor = color.Black                 // Définir la couleur de la bordure
	border.StrokeWidth = 3                           // Définir l'épaisseur de la bordure
	border.CornerRadius = 20                         // Définir les coins arrondis

	cardContent.Add(border) // Ajouter le rectangle de contour à la carte

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
