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

var artists = []Artist{ // Définir les données des artistes
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
	{Name: "Groshi", Image: "public/yoshi.png", YearStarted: 2000, DebutAlbum: time.Date(2001, time.July, 27, 0, 0, 0, 0, time.UTC), Members: []string{"Mario Mario", "Luigi Mario", "Toad le champi"}, LastConcert: Concert{Date: time.Date(2023, time.December, 9, 0, 0, 0, 0, time.UTC), Location: "Tokyo - Japon"}, NextConcerts: []Concert{{Date: time.Date(2026, time.April, 30, 0, 0, 0, 0, time.UTC), Location: "Paris - France"}, {Date: time.Date(2024, time.July, 3, 0, 0, 0, 0, time.UTC), Location: "London"}}},
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Menu - Groupie Tracker")

	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Artists...")

	searchResults := container.New(layout.NewVBoxLayout())

	searchButton := widget.NewButton("Search", func() {
		recherche(searchBar, searchResults, artists)
	})

	searchBar.OnSubmitted = func(_ string) {
		searchButton.OnTapped()
	}

	searchBar.OnChanged = func(text string) {
		generateSearchSuggestions(text, searchResults, artists)
	}

	artistsContainer := container.NewVBox()

	for i := 0; i < len(artists); i += 3 {
		rowContainer := container.NewHBox()
		columnContainer := container.NewVBox()

		space := widget.NewLabel("")

		rowContainer.Add(space)
		rowContainer.Add(space)
		rowContainer.Add(space)
		for j := i; j < i+3 && j < len(artists); j++ {
			card := createCardGeneralInfo(artists[j])
			rowContainer.Add(card)

			if j < i+2 && j < len(artists) {
				rowContainer.Add(space)
			}
		}

		columnContainer.Add(rowContainer)
		artistsContainer.Add(columnContainer)
	}

	scrollContainer := container.NewVScroll(artistsContainer)
	scrollContainer.SetMinSize(fyne.NewSize(1080, 720))

	content := container.NewVBox(
		searchBar,
		searchButton,
		searchResults,
		scrollContainer,
	)

	centeredContent := container.New(layout.NewCenterLayout(), content)

	background := canvas.NewRectangle(color.NRGBA{R: 0x5C, G: 0x64, B: 0x73, A: 0xFF})
	background.Resize(fyne.NewSize(1080, 720))

	backgroundContainer := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background)

	backgroundContainer.Add(centeredContent)

	myWindow.SetContent(backgroundContainer)
	myWindow.Resize(fyne.NewSize(1080, 720))
	myWindow.ShowAndRun()
}

func generateSearchSuggestions(text string, searchResults *fyne.Container, artists []Artist) {
	searchResults.Objects = nil

	if text == "" {
		return
	}

	var found bool
	var correspondingResultAdded bool

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(text)) {
			found = true

			if !correspondingResultAdded {
				correspondingResultLabel := widget.NewLabel("Corresponding result: ")
				searchResults.Add(correspondingResultLabel)

				correspondingResultAdded = true
			}

			artistButton := widget.NewButton(artist.Name, func() {
				fmt.Println(artist.Name)
				fmt.Print("Affiche toutes les informations de l'artiste (nouvelle page)") // Nouvelle page de Giovanni
			})

			searchResults.Add(layout.NewSpacer())

			searchResults.Add(artistButton)
		}
	}

	if !found {
		noResultLabel := widget.NewLabel("No result")
		searchResults.Add(noResultLabel)
	}
}

func recherche(searchBar *widget.Entry, scrollContainer *fyne.Container, artists []Artist) {
	searchText := searchBar.Text

	scrollContainer.Objects = nil

	var foundArtists []Artist
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchText)) {
			foundArtists = append(foundArtists, artist)
		}
	}

	if len(foundArtists) > 0 {
		for i := 0; i < len(foundArtists); i += 3 {
			rowContainer := container.NewHBox()
			columnContainer := container.NewVBox()

			space := widget.NewLabel("")

			rowContainer.Add(space)
			rowContainer.Add(space)
			rowContainer.Add(space)
			for j := i; j < i+3 && j < len(foundArtists); j++ {
				card := createCardGeneralInfo(foundArtists[j])
				rowContainer.Add(card)

				if j < i+2 && j < len(foundArtists) {
					rowContainer.Add(space)
				}
			}

			columnContainer.Add(rowContainer)
			scrollContainer.Add(columnContainer)
		}
	} else {
		noResultLabel := widget.NewLabel("No result found")
		scrollContainer.Add(noResultLabel)
	}

	scrollContainer.Refresh()
}

func createCardGeneralInfo(artist Artist) fyne.CanvasObject {
	image := canvas.NewImageFromFile(artist.Image)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(120, 120))
	image.Resize(fyne.NewSize(120, 120))

	averageColor := getAverageColor(artist.Image)

	background := canvas.NewRectangle(averageColor)
	background.SetMinSize(fyne.NewSize(300, 300))
	background.Resize(fyne.NewSize(296, 296))
	background.CornerRadius = 20

	button := widget.NewButton("More information", func() {
		fmt.Println(artist.Name)
		fmt.Print("Affiche toutes les informations de l'artiste (nouvelle page)") //nouvelle page de Giovanni
	})

	buttonFavorie := widget.NewButton("Add to favorite", func() {
		fmt.Println("Ajouter aux favoris")
		//fonction à intégré ici
	})

	buttonsContainer := container.NewHBox(
		layout.NewSpacer(),
		button,
		layout.NewSpacer(),
		buttonFavorie,
		layout.NewSpacer(),
	)

	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	yearLabel := widget.NewLabel(fmt.Sprintf("%d", artist.YearStarted))

	labelsContainer := container.NewHBox(nameLabel, yearLabel)

	var membersText string
	if len(artist.Members) == 1 {
		membersText = "Solo Artist\n"
	} else if len(artist.Members) > 0 {
		membersText = "Members: " + strings.Join(artist.Members, ", ")
	}
	membersLabel := widget.NewLabel(membersText)
	membersLabel.Wrapping = fyne.TextWrapWord

	infoContainer := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), image, labelsContainer, membersLabel, layout.NewSpacer(), buttonsContainer, layout.NewSpacer())
	infoContainer.Resize(fyne.NewSize(300, 180))

	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background, infoContainer)
	cardContent.Resize(fyne.NewSize(300, 300))

	border := canvas.NewRectangle(color.Transparent)
	border.SetMinSize(fyne.NewSize(300, 300))
	border.Resize(fyne.NewSize(296, 296))
	border.StrokeColor = color.Black
	border.StrokeWidth = 3
	border.CornerRadius = 20

	cardContent.Add(border)

	return cardContent
}

/* func createCardAllInfo(artist Artist) fyne.CanvasObject {
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
} */

func getAverageColor(imagePath string) color.Color {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return color.Black
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Erreur lors du décodage de l'image:", err)
		return color.Black
	}

	var totalRed, totalGreen, totalBlue float64
	totalPixels := 0

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelColor := img.At(x, y)
			r, g, b, _ := pixelColor.RGBA()

			red := float64(r) / 65535.0
			green := float64(g) / 65535.0
			blue := float64(b) / 65535.0

			totalRed += red
			totalGreen += green
			totalBlue += blue

			totalPixels++
		}
	}

	averageRed := totalRed / float64(totalPixels)
	averageGreen := totalGreen / float64(totalPixels)
	averageBlue := totalBlue / float64(totalPixels)

	averageRed = averageRed * 255
	averageGreen = averageGreen * 255
	averageBlue = averageBlue * 255

	averageColor := color.RGBA{
		R: uint8(averageRed),
		G: uint8(averageGreen),
		B: uint8(averageBlue),
		A: 255,
	}

	return averageColor
}
