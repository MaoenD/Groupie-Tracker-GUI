package main

import ( //importation des bibliotheques nécessaires

	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"
)

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
			searchResultCountLabel.SetText("No result")
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

	content := container.NewVBox(
		searchBarContainer,
		searchResults,
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
		if strings.Contains(strings.ToLower(artist.Name), searchText) ||
			strconv.Itoa(artist.YearStarted) == searchText ||
			strconv.Itoa(artist.DebutAlbum.Year()) == searchText ||
			checkMemberName(artist.Members, searchText) ||
			checkConcertLocation(artist.NextConcerts, searchText) {
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

func refreshContent(searchBar *widget.Entry, searchResultCountLabel *widget.Label, artistsContainer *fyne.Container, artists []Artist, myApp fyne.App) {
	searchBar.SetText("") // Effacer le texte de la barre de recherche

	// Réinitialiser le contenu de artistsContainer avec les résultats de recherche initiaux
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

	// Actualiser l'affichage du contenu
	artistsContainer.Refresh()

	// Actualiser le texte du label du nombre de résultats
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
	numMembersCheck = widget.NewCheckGroup([]string{"All", "2", "3", "4", "5", "6+"}, func(selected []string) {})
	numMembersCheck.SetSelected([]string{"All"})
	numMembersBox = container.NewHBox()
	for _, option := range []string{"All", "2", "3", "4", "5", "6+"} {
		numMembersBox.Add(widget.NewCheck(option, func(checked bool) {}))
	}
	numMembersBox.Hide()
	locationsSelect = widget.NewSelect(concertLocations, func(selected string) {})
	locationsSelect.Resize(fyne.NewSize(200, 150))

	radioSoloGroup.SetSelected(selectedRadioValue)
	numMembersCheck.SetSelected(selectedNumMembers)
	locationsSelect.SetSelected(selectedLocationValue)

	creationDateRange.SetValue(savedCreationRange)
	firstAlbumDateRange.SetValue(savedFirstAlbumRange)
	savedCreationRange = creationDateRange.Value
	savedFirstAlbumRange = firstAlbumDateRange.Value

	myWindow = myApp.NewWindow("Groupie Tracker GUI Filters")
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.SetFixedSize(true)

	reset := widget.NewButton("Reset Filters", func() {})
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

func applyFilter() {
	selectedRadioValue = radioSoloGroup.Selected
	selectedLocationValue = locationsSelect.Selected

	selectedNumMembers = numMembersCheck.Selected

	fmt.Printf("Radio sélectionné: %s, Membres sélectionnés: %v, Localisation sélectionnée: %s\n", selectedRadioValue, selectedNumMembers, selectedLocationValue)

	savedCreationRange = creationDateRange.Value
	savedFirstAlbumRange = firstAlbumDateRange.Value
}
