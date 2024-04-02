package main

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
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Menu - Groupie Tracker")
	logoApp, _ := fyne.LoadResourceFromPath("public/logo.png")
	myWindow.SetIcon(logoApp)

	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Artists...")
	searchBar.Resize(fyne.NewSize(1080, searchBar.MinSize().Height))
	searchResults := container.NewVBox()
	artistsContainer := container.NewVBox()
	scrollContainer := container.NewVScroll(artistsContainer)
	scrollContainer.SetMinSize(fyne.NewSize(1080, 720))
	searchButton := widget.NewButton("Search", func() {
		recherche(searchBar, artistsContainer, artists, myApp)
		searchBar.SetText("")
	})

	searchResultCountLabel := widget.NewLabel("")

	logoButton := widget.NewButtonWithIcon("", (loadImageResource("public/logo.png")), func() {
		refreshContent(searchBar, searchResultCountLabel, artistsContainer, artists, myApp)
	})

	filterButton := widget.NewButton("Filtrer", func() {
		Filter(myApp)
	})

	searchBarContainer := container.NewVBox(
		container.NewBorder(nil, nil, logoButton, filterButton, searchBar),
		searchButton,
		searchResultCountLabel,
	)
	searchBarContainer.Resize(searchBarContainer.MinSize())

	searchBar.OnSubmitted = func(_ string) {
		searchButton.OnTapped()
	}

	searchBar.OnChanged = func(text string) {
		count := generateSearchSuggestions(text, searchResults, artists, myApp, 5)
		if count != 0 {
			searchResultCountLabel.SetText(fmt.Sprintf("Results for '%s':", text))
		} else {
			searchResultCountLabel.SetText("")
		}
	}

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

	blockContent := createBlockContent()

	content := container.NewVBox(
		searchBarContainer,
		searchResults,
		blockContent,
		scrollContainer,
	)

	dynamicLayout := layout.NewVBoxLayout()
	content.Layout = dynamicLayout

	centeredContent := container.New(layout.NewCenterLayout(), content)

	background := canvas.NewRectangle(color.NRGBA{R: 0x5C, G: 0x64, B: 0x73, A: 0xFF})
	background.Resize(fyne.NewSize(1080, 720))

	backgroundContainer := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background)

	backgroundContainer.Add(centeredContent)

	myWindow.SetContent(backgroundContainer)
	myWindow.Resize(fyne.NewSize(1080, 720))
	myWindow.ShowAndRun()
}

func createBlockContent() fyne.CanvasObject {
	imagePath := "public/world_map1.jpg"

	image := canvas.NewImageFromFile(imagePath)

	if image == nil {
		fmt.Println("Impossible de charger l'image:", imagePath)
		return nil
	}

	image.FillMode = canvas.ImageFillStretch

	blockContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil),
		image,
	)

	button := widget.NewButton("", func() {
		fmt.Print("Logique de map à intégrer ici")
	})
	button.Importance = widget.LowImportance
	button.Resize(image.MinSize())

	blockContent.Add(button)

	title := widget.NewLabel("Geolocation feature")
	description := widget.NewLabel("Find out where and when your favorite artists will be performing around the globe.")
	description.Wrapping = fyne.TextWrapWord

	textContainer := container.New(layout.NewVBoxLayout(),
		title,
		description,
	)

	blockContent.Add(textContainer)

	return blockContent
}

func refreshContent(searchBar *widget.Entry, searchResultCountLabel *widget.Label, artistsContainer *fyne.Container, artists []Artist, myApp fyne.App) {
	searchBar.SetText("")

	artistsContainer.Objects = nil
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

	artistsContainer.Refresh()

	searchResultCountLabel.SetText("")
}

func checkMemberName(members []string, searchText string) bool {
	for _, member := range members {
		if strings.Contains(strings.ToLower(member), searchText) {
			return true
		}
	}
	return false
}

func checkConcertLocation(concerts []Concert, searchText string) bool {
	for _, concert := range concerts {
		if strings.Contains(strings.ToLower(concert.Location), searchText) {
			return true
		}
	}
	return false
}

func createCardGeneralInfo(artist Artist, myApp fyne.App) fyne.CanvasObject {
	image := canvas.NewImageFromFile(artist.Image)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(120, 120))
	image.Resize(fyne.NewSize(120, 120))

	averageColor := getAverageColor(artist.Image)

	background := canvas.NewRectangle(averageColor)
	background.SetMinSize(fyne.NewSize(300, 300))
	background.Resize(fyne.NewSize(296, 296))
	background.CornerRadius = 20

	button := widget.NewButton("          More information          ", func() {
		fmt.Println(artist.Name)
		fmt.Print("Affiche toutes les informations de l'artiste (nouvelle page)")
		SecondPage(artist, myApp)
	})

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

	// Code existant omis pour des raisons de concision

	// Ajouter le bouton de favori au conteneur de boutons
	buttonsContainer := container.NewHBox(
		widget.NewLabel("  "),
		container.NewBorder(nil, layout.NewSpacer(), nil, likeButton, button),
	)

	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	yearLabel := widget.NewLabel(fmt.Sprintf("%d", artist.YearStarted))

	labelsContainer := container.NewHBox(nameLabel, yearLabel)

	var membersText string
	if len(artist.Members) == 1 {
		membersText = "Solo Artist\n"
	} else if len(artist.Members) > 0 {
		membersText = "Members : " + strings.Join(artist.Members, ", ")
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

func SecondPage(artist Artist, myApp fyne.App) {
	myWindow := myApp.NewWindow("Information - " + artist.Name)

	logo, _ := fyne.LoadResourceFromPath("public/logo.png")
	myWindow.SetIcon(logo)

	averageColor := getAverageColor(artist.Image)

	background := canvas.NewRectangle(averageColor)
	background.SetMinSize(fyne.NewSize(300, 300))
	background.Resize(fyne.NewSize(296, 296))

	image := canvas.NewImageFromFile(artist.Image)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(320, 320))
	image.Resize(fyne.NewSize(220, 220))
	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	yearLabel := widget.NewLabel(fmt.Sprintf("Year Started: %d", artist.YearStarted))
	debutAlbumLabel := widget.NewLabel(fmt.Sprintf("Debut Album: %s", artist.DebutAlbum.Format("02-Jan-2006")))
	membersLabel := widget.NewLabel(fmt.Sprintf("Members: %s", strings.Join(artist.Members, ", ")))
	lastConcertLabel := widget.NewLabel(fmt.Sprintf("Last Concert: %s - %s", artist.LastConcert.Date.Format("02-Jan-2006"), artist.LastConcert.Location))
	nextConcertLabel := widget.NewLabel("Next Concert:")
	if len(artist.NextConcerts) > 0 {
		nextConcertLabel.Text += fmt.Sprintf(" %s - %s", artist.NextConcerts[0].Date.Format("02-Jan-2006"), artist.NextConcerts[0].Location)
	} else {
		nextConcertLabel.Text += " No upcoming concerts" // Affichage si aucun événeement
	}

	nameLabel.Alignment = fyne.TextAlignCenter
	yearLabel.Alignment = fyne.TextAlignCenter
	debutAlbumLabel.Alignment = fyne.TextAlignCenter
	membersLabel.Alignment = fyne.TextAlignCenter
	lastConcertLabel.Alignment = fyne.TextAlignCenter
	nextConcertLabel.Alignment = fyne.TextAlignCenter

	infoContainer := container.NewVBox(
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
	cardContent.Resize(fyne.NewSize(300, 300))

	myWindow.SetContent(cardContent)
	myWindow.CenterOnScreen()
	myWindow.Show()
}

func loadImageResource(path string) fyne.Resource {
	image, err := fyne.LoadResourceFromPath(path)
	if err != nil {
		fmt.Println("Erreur lors du chargement de l'icône:", err)
		return nil
	}
	return image
}

func generateSearchSuggestions(text string, scrollContainer *fyne.Container, artists []Artist, myApp fyne.App, limit int) int {
	scrollContainer.Objects = nil

	if text == "" || len(artists) == 0 {
		return 0
	}

	count := 0

	for _, artist := range artists {
		if count >= limit {
			break
		}

		if artistMatchesFilters(artist, savedFilter) {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(text)) {
				count++
				artistButton := widget.NewButton(artist.Name, func(a Artist) func() {
					return func() {
						SecondPage(a, myApp)
					}
				}(artist))
				artistButton.Importance = widget.LowImportance
				scrollContainer.Add(artistButton)
			} else {
				if strconv.Itoa(artist.YearStarted) == text {
					count++
					artistButton := widget.NewButton(artist.Name+" (Year Started: "+text+")", func(a Artist) func() {
						return func() {
							SecondPage(a, myApp)
						}
					}(artist))
					artistButton.Importance = widget.LowImportance
					scrollContainer.Add(artistButton)
				}

				if strconv.Itoa(artist.DebutAlbum.Year()) == text {
					count++
					artistButton := widget.NewButton(artist.Name+" (Debut Album: "+strconv.Itoa(artist.DebutAlbum.Year())+")", func(a Artist) func() {
						return func() {
							SecondPage(a, myApp)
						}
					}(artist))
					artistButton.Importance = widget.LowImportance
					scrollContainer.Add(artistButton)
				}

				if len(artist.Members) > 1 {
					for _, member := range artist.Members {
						if strings.Contains(strings.ToLower(member), strings.ToLower(text)) {
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

				for _, concert := range artist.NextConcerts {
					if strings.Contains(strings.ToLower(concert.Location), strings.ToLower(text)) {
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

	if count < len(artists) && count >= limit {
		showMoreButton := widget.NewButton("More results", func() {
			generateSearchSuggestions(text, scrollContainer, artists, myApp, limit+5)
		})
		scrollContainer.Add(showMoreButton)
	}

	return count
}

func recherche(searchBar *widget.Entry, scrollContainer *fyne.Container, artists []Artist, myApp fyne.App) {
	searchText := strings.ToLower(searchBar.Text)

	artistsContainer := container.NewVBox()

	var foundArtists []Artist

	for _, artist := range artists {
		if artistMatchesFilters(artist, savedFilter) {
			if strings.Contains(strings.ToLower(artist.Name), searchText) ||
				strconv.Itoa(artist.YearStarted) == searchText ||
				strconv.Itoa(artist.DebutAlbum.Year()) == searchText ||
				checkMemberName(artist.Members, searchText) ||
				checkConcertLocation(artist.NextConcerts, searchText) {
				foundArtists = append(foundArtists, artist)
			}
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
		noResultLabel := widget.NewLabel("No result found")
		artistsContainer.Add(noResultLabel)
	}

	scrollContainer.Objects = []fyne.CanvasObject{artistsContainer}
	scrollContainer.Refresh()
}

func artistMatchesFilters(artist Artist, filter saveFilter) bool {
	if filter.RadioSelected != "" {
		if filter.RadioSelected == "Solo" && len(artist.Members) > 1 {
			return false
		} else if filter.RadioSelected == "Group" && len(artist.Members) <= 1 {
			return false
		}
	}

	if len(filter.NumMembersSelected) > 0 && !contains(filter.NumMembersSelected, strconv.Itoa(len(artist.Members))) {
		return false
	}

	if filter.LocationSelected != "" && !artistHasConcertLocation(artist, filter.LocationSelected) {
		return false
	}

	if filter.CreationRange > 0 && float64(artist.YearStarted) < filter.CreationRange {
		return false
	}

	if filter.FirstAlbumRange > 0 && float64(artist.DebutAlbum.Year()) < filter.FirstAlbumRange {
		return false
	}
	return true
}

func artistHasConcertLocation(artist Artist, location string) bool {
	for _, concert := range artist.NextConcerts {
		if strings.EqualFold(concert.Location, location) {
			return true
		}
	}
	return false
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

var (
	minCreationYear     int
	maxCreationYear     int
	minFirstAlbumYear   int
	maxFirstAlbumYear   int
	concertLocations    []string
	creationDateRange   *widget.Slider
	firstAlbumDateRange *widget.Slider
	radioSoloGroup      *widget.RadioGroup
	numMembersCheck     *widget.CheckGroup
	numMembersBox       *fyne.Container
	locationsSelect     *widget.Select
	myWindow            fyne.Window
	windowOpened        bool

	selectedRadioValue    string
	selectedNumMembers    []string
	selectedLocationValue string

	savedCreationRange   float64
	savedFirstAlbumRange float64
	savedNumMembers      []string
)

func Filter(myApp fyne.App) {
	if myWindow != nil {
		myWindow.Close()
	}
	initializeFilters(myApp)
}

func initializeFilters(myApp fyne.App) {
	minCreationYear = artists[0].YearStarted
	maxCreationYear = artists[0].YearStarted
	minFirstAlbumYear = artists[0].DebutAlbum.Year()
	maxFirstAlbumYear = artists[0].DebutAlbum.Year()

	concertLocations = make([]string, 0)
	locationsMap := make(map[string]bool)
	for _, artist := range artists {
		if artist.YearStarted < minCreationYear {
			minCreationYear = artist.YearStarted
		}
		if artist.YearStarted > maxCreationYear {
			maxCreationYear = artist.YearStarted
		}
		if year := artist.DebutAlbum.Year(); year < minFirstAlbumYear {
			minFirstAlbumYear = year
		}
		if year := artist.DebutAlbum.Year(); year > maxFirstAlbumYear {
			maxFirstAlbumYear = year
		}
		for _, concert := range artist.NextConcerts {
			if _, found := locationsMap[concert.Location]; !found {
				concertLocations = append(concertLocations, concert.Location)
				locationsMap[concert.Location] = true
			}
		}
	}

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

	numMembersBox.Hide()
	locationsSelect = widget.NewSelect(concertLocations, func(selected string) {})
	locationsSelect.Resize(fyne.NewSize(200, 150))

	radioSoloGroup.SetSelected(selectedRadioValue)
	numMembersCheck.SetSelected(selectedNumMembers)
	locationsSelect.SetSelected(selectedLocationValue)
	numMembersCheck.SetSelected(selectedNumMembers)

	creationDateRange.SetValue(savedCreationRange)
	firstAlbumDateRange.SetValue(savedFirstAlbumRange)
	numMembersCheck.SetSelected(selectedNumMembers)
	savedCreationRange = creationDateRange.Value
	savedFirstAlbumRange = firstAlbumDateRange.Value
	savedNumMembers = selectedNumMembers

	myWindow = myApp.NewWindow("Groupie Tracker GUI Filters")
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.SetFixedSize(true)

	reset := widget.NewButton("Reset Filters", func() {
		selectedRadioValue = ""
		selectedNumMembers = nil
		selectedLocationValue = ""

		creationDateRange.SetValue(float64(minCreationYear))
		firstAlbumDateRange.SetValue(float64(minFirstAlbumYear))

		radioSoloGroup.SetSelected(selectedRadioValue)
		numMembersCheck.SetSelected(selectedNumMembers)
		locationsSelect.SetSelected(selectedLocationValue)
	})

	applyButton := widget.NewButton("Apply Filters", func() {
		applyFilter()
		myWindow.Close()
	})

	creationDateRangeLabel := widget.NewLabel(fmt.Sprintf("Creation Date Range: %d - %d", minCreationYear, maxCreationYear))
	firstAlbumDateRangeLabel := widget.NewLabel(fmt.Sprintf("First Album Date Range: %d - %d", minFirstAlbumYear, maxFirstAlbumYear))

	updateLabels := func() {
		creationRange := int(creationDateRange.Value)
		firstAlbumRange := int(firstAlbumDateRange.Value)

		creationDateRangeLabel.SetText(fmt.Sprintf("Creation Date Range: %d - %d", creationRange, maxCreationYear))
		firstAlbumDateRangeLabel.SetText(fmt.Sprintf("First Album Date Range: %d - %d", firstAlbumRange, maxFirstAlbumYear))
	}

	creationDateRange.OnChanged = func(value float64) {
		updateLabels()
	}

	firstAlbumDateRange.OnChanged = func(value float64) {
		updateLabels()
	}

	// Création du conteneur de widgets
	filtersContainer := container.NewVBox(
		reset,
		creationDateRangeLabel, creationDateRange,
		firstAlbumDateRangeLabel, firstAlbumDateRange,
		radioSoloGroup,
		numMembersBox,
		locationsSelect,
		applyButton,
	)

	myWindow.SetContent(filtersContainer)
	myWindow.CenterOnScreen()
	myWindow.Show()
	windowOpened = true
}

type saveFilter struct {
	RadioSelected      string
	NumMembersSelected []string
	LocationSelected   string
	CreationRange      float64
	FirstAlbumRange    float64
}

var savedFilter saveFilter

func applyFilter() saveFilter {
	selectedRadioValue = radioSoloGroup.Selected
	selectedLocationValue = locationsSelect.Selected

	savedNumMembers = selectedNumMembers

	savedCreationRange = creationDateRange.Value
	savedFirstAlbumRange = firstAlbumDateRange.Value

	savedFilter = saveFilter{
		RadioSelected:      selectedRadioValue,
		NumMembersSelected: selectedNumMembers,
		LocationSelected:   selectedLocationValue,
		CreationRange:      savedCreationRange,
		FirstAlbumRange:    savedFirstAlbumRange,
	}

	fmt.Printf("Radio sélectionné: %s, Membres sélectionnés: %v, Localisation sélectionnée: %s, savedCreationRange: %f, savedFirstAlbumRange: %f\n", selectedRadioValue, selectedNumMembers, selectedLocationValue, savedCreationRange, savedFirstAlbumRange)

	return savedFilter
}
