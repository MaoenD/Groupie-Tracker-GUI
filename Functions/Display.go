package Functions

import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// CreateCardGeneralInfo creates a card with general information about an artist
func CreateCardGeneralInfo(artist Artist, relation Relation, myApp fyne.App) fyne.CanvasObject {
	// Load image of the artist from URL
	response, err := http.Get(artist.Image)
	if err != nil {
		fmt.Println("Failed to load image:", err)
		return nil
	}
	defer response.Body.Close()

	// Read image data
	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read image data:", err)
		return nil
	}

	// Determine file type of the image from URL
	parts := strings.Split(artist.Image, ".")
	fileType := parts[len(parts)-1]

	// Create image from read data and file type
	image := canvas.NewImageFromReader(strings.NewReader(string(imageData)), fileType)

	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(120, 120))
	image.Resize(fyne.NewSize(120, 120))

	// Calculate average color of the image
	r, g, b, a := getAverageColor(image)

	// Create background with the average color
	background := canvas.NewRectangle(color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
	background.SetMinSize(fyne.NewSize(300, 300))
	background.Resize(fyne.NewSize(296, 296))
	background.CornerRadius = 20

	// Create "More Information" button
	button := widget.NewButton("          More informations          ", func() {
		SecondPage(artist, relation, myApp)
	})

	// Create like button
	var likeButton *widget.Button
	var likeIcon string
	if artist.Favorite {
		likeIcon = "public/img/likeOn.ico"
	} else {
		likeIcon = "public/img/likeOff.ico"
	}
	likeButton = widget.NewButton("", func() {
		artist.Favorite = !artist.Favorite
		if artist.Favorite {
			likeButton.SetIcon(LoadImageResource("public/img/likeOn.ico"))
		} else {
			likeButton.SetIcon(LoadImageResource("public/img/likeOff.ico"))
		}
	})
	likeButton.SetIcon(LoadImageResource(likeIcon))

	// Create container for buttons
	buttonsContainer := container.NewHBox(
		widget.NewLabel("  "),
		container.NewBorder(nil, layout.NewSpacer(), nil, likeButton, button),
	)

	// Create label for artist name
	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Create label for creation year
	yearLabel := widget.NewLabel(fmt.Sprintf("%d", artist.CreationDate))

	// Create container for name and year labels
	labelsContainer := container.NewHBox(nameLabel, yearLabel)

	// Create label for group members
	var membersText string
	if len(artist.Members) == 1 {
		membersText = "Solo Artist \n"
	} else if len(artist.Members) > 0 {
		membersText = "Members : " + strings.Join(artist.Members, ", ")
	}
	membersLabel := widget.NewLabel(membersText)
	membersLabel.Wrapping = fyne.TextWrapWord

	// Create container for artist information
	infoContainer := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), image, labelsContainer, membersLabel, layout.NewSpacer(), buttonsContainer, layout.NewSpacer())
	infoContainer.Resize(fyne.NewSize(300, 180))

	// Create card content with background and artist information
	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background, infoContainer)
	cardContent.Resize(fyne.NewSize(300, 300))

	// Create border around card content
	border := canvas.NewRectangle(color.Transparent)
	border.SetMinSize(fyne.NewSize(300, 300))
	border.Resize(fyne.NewSize(296, 296))
	border.StrokeColor = color.Black
	border.StrokeWidth = 3
	border.CornerRadius = 20

	cardContent.Add(border)

	return cardContent
}

// GenerateSearchSuggestions generates search suggestions based on user input
func GenerateSearchSuggestions(text string, scrollContainer *fyne.Container, artists []Artist, relation Relation, myApp fyne.App, limit int) int {
	// Clear previous objects from scroll container
	scrollContainer.Objects = nil

	// Check if search text is empty or no artists in the list
	if text == "" || len(artists) == 0 {
		return 0
	}

	// Extract category tag and search text
	category, searchText := extractCategoryAndText(text)

	// Counter for displayed suggestions
	count := 0

	// Iterate over all artists in the list
	for _, artist := range artists {
		// Check if displayed suggestion count reaches the specified limit
		if count >= limit {
			break
		}

		// Check if artist matches saved filters
		if artistMatchesFilters(artist, savedFilter) {
			// Check if artist name, creation year, debut album year, member name, or concert location matches search text
			if category == "" && (strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchText)) ||
				strconv.Itoa(artist.CreationDate) == searchText ||
				checkDebutAlbumYear(artist.FirstAlbum, searchText) ||
				checkMemberName(artist.Members, searchText) ||
				checkConcertLocation(artist.NextConcerts, searchText)) ||
				category == "a" && strings.Contains(strings.ToLower(artist.Name), searchText) ||
				category == "l" && checkConcertLocation(artist.NextConcerts, searchText) ||
				category == "m" && checkMemberName(artist.Members, searchText) {
				// Increment counter and add artist button to scroll container
				count++
				artistButton := widget.NewButton(artist.Name, func(a Artist) func() {
					return func() {
						SecondPage(a, relation, myApp)
					}
				}(artist))
				artistButton.Importance = widget.LowImportance
				scrollContainer.Add(artistButton)
			}
		}
	}
	return count
}

// extractCategoryAndText extracts category tag and search text from input text
func extractCategoryAndText(text string) (string, string) {
	if len(text) < 3 || text[1] != '/' {
		return "", text
	}
	category := text[0:1]
	searchText := text[3:]
	return category, searchText
}

// checkDebutAlbumYear checks if debut album year matches search text
func checkDebutAlbumYear(firstAlbum string, searchText string) bool {
	albumYearParts := strings.Split(firstAlbum, "-")
	if len(albumYearParts) == 3 && albumYearParts[2] == searchText {
		return true
	}
	return false
}

// Recherche performs search based on user input
func Recherche(searchBar *widget.Entry, scrollContainer *fyne.Container, artists []Artist, relation Relation, myApp fyne.App) {
	// Convert search text to lowercase for case-insensitive search
	searchText := strings.ToLower(searchBar.Text)

	// Extract category tag and search text
	category, text := extractCategoryAndText(searchText)

	// Create container to store found artists
	artistsContainer := container.NewVBox()

	// List of found artists
	var foundArtists []Artist

	// Iterate over all artists in the list
	for _, artist := range artists {
		// Check if artist matches saved filters
		if artistMatchesFilters(artist, savedFilter) {
			// Check if artist name, creation year, debut album year, member name, or concert location matches search text
			if category == "" && (strings.Contains(strings.ToLower(artist.Name), text) ||
				strconv.Itoa(artist.CreationDate) == text ||
				strings.Contains(artist.FirstAlbum, text) || // Adjusted for string comparison
				checkMemberName(artist.Members, text) ||
				checkConcertLocation(artist.NextConcerts, text)) || // Adjusted to use artist.NextConcerts
				category == "a" && strings.Contains(strings.ToLower(artist.Name), text) ||
				category == "l" && checkConcertLocation(artist.NextConcerts, text) ||
				category == "m" && checkMemberName(artist.Members, text) {
				foundArtists = append(foundArtists, artist)
			}
		}
	}

	// Check if there are found artists
	if len(foundArtists) > 0 {
		// Display found artists in groups of 3 in row and column containers
		for i := 0; i < len(foundArtists); i += 3 {
			rowContainer := container.NewHBox()
			columnContainer := container.NewVBox()

			space := widget.NewLabel("")

			rowContainer.Add(space)
			rowContainer.Add(space)
			rowContainer.Add(space)
			for j := i; j < i+3 && j < len(foundArtists); j++ {
				card := CreateCardGeneralInfo(foundArtists[j], relation, myApp)
				rowContainer.Add(card)

				if j < i+2 && j < len(foundArtists) {
					rowContainer.Add(space)
				}
			}

			columnContainer.Add(rowContainer)
			artistsContainer.Add(columnContainer)
		}
	} else {
		// Display message indicating no result found
		noResultLabel := widget.NewLabel("No result found")
		artistsContainer.Add(noResultLabel)
	}

	// Update objects in scroll container with found artists
	scrollContainer.Objects = []fyne.CanvasObject{artistsContainer}
	scrollContainer.Refresh()
}

// artistMatchesFilters checks if artist matches user filters
func artistMatchesFilters(artist Artist, filter saveFilter) bool {
	// Check if artist matches user filters
	if filter.RadioSelected != "" {
		// Check if user selected "Solo" and artist is a group, or vice versa
		if filter.RadioSelected == "Solo" && len(artist.Members) > 1 {
			return false
		} else if filter.RadioSelected == "Group" && len(artist.Members) <= 1 {
			return false
		}
	}

	// Check if number of artist members matches any of user selections
	if len(filter.NumMembersSelected) > 0 && !contains(filter.NumMembersSelected, strconv.Itoa(len(artist.Members))) {
		return false
	}

	// Check if artist concert location matches user selection
	if filter.LocationSelected != "" && !artistHasConcertLocation(artist, filter.LocationSelected) {
		return false
	}

	// Check if artist creation year is in user-selected range
	if filter.CreationRange > 0 && float64(artist.CreationDate) < filter.CreationRange {
		return false
	}

	// Check if artist first album release year is in user-selected range
	if filter.FirstAlbumRange > 0 {
		albumYear := parseYearFromAlbum(artist.FirstAlbum)
		if float64(albumYear) < filter.FirstAlbumRange {
			return false
		}
	}
	return true
}

// artistHasConcertLocation checks if artist has concert at specified location
func artistHasConcertLocation(artist Artist, location string) bool {
	// Check if artist has concert at specified location
	for _, concert := range artist.NextConcerts {
		for _, loc := range concert.Locations {
			if strings.EqualFold(string(loc), location) { // Check without case sensitivity
				return true
			}
		}
	}
	return false
}

// parseYearFromAlbum parses year from album date
func parseYearFromAlbum(albumDate string) int {
	parts := strings.Split(albumDate, "-")
	year, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0
	}
	return year
}
