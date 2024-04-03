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

// CreationDateRange represents the slider widget for selecting the creation date range.
var CreationDateRange *widget.Slider

// CreateBlockContent creates a canvas object containing image, button, title, and description.
func CreateBlockContent() fyne.CanvasObject {
	// Path to the image to load
	imagePath := "public/img/world_map.jpg"

	// Load the image from the specified path
	image := canvas.NewImageFromFile(imagePath)

	// Check if the image was loaded successfully
	if image == nil {
		// Print an error message if the image failed to load
		fmt.Println("Failed to load image:", imagePath)
		return nil
	}

	// Set the image fill mode to stretch to fill the space
	image.FillMode = canvas.ImageFillStretch

	// Create a container to organize the image and text with a border layout
	blockContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil),
		image,
	)

	// Create an empty button for adding interaction (to be filled as per specific logic)
	button := widget.NewButton("", func() {
		// Action to perform when the button is clicked (to be filled as per needs)
		fmt.Print("Map logic to be integrated here")
	})
	button.Importance = widget.LowImportance // Set the button importance as low
	button.Resize(image.MinSize())           // Resize the button to match the image size

	// Add the button to the block content
	blockContent.Add(button)

	// Create title and description labels for the content
	title := widget.NewLabel("Geolocation feature")
	description := widget.NewLabel("Find out where and when your favorite artists performed around the globe.")
	description.Wrapping = fyne.TextWrapWord // Enable text wrapping for the description

	// Create a container to organize the text labels with a vertical layout
	textContainer := container.New(layout.NewVBoxLayout(),
		title,
		description,
	)

	// Add the text container to the block content
	blockContent.Add(textContainer)

	// Return the block content
	return blockContent
}

// RefreshContent resets the search bar text and updates the content based on search results.
func RefreshContent(searchBar *widget.Entry, searchResultCountLabel *widget.Label, artistsContainer *fyne.Container, relation Relation, artists []Artist, myApp fyne.App) {
	// Reset the search bar text
	searchBar.SetText("")

	// Clear all existing objects in the artists container
	artistsContainer.Objects = nil

	// Iterate over artists and organize them into cards in containers
	for i := 0; i < len(artists); i += 3 {
		rowContainer := container.NewHBox()    // Create a horizontal row container for artist cards
		columnContainer := container.NewVBox() // Create a vertical column container for card rows

		space := widget.NewLabel("") // Create an empty space for spacing between cards

		// Add spaces between cards for visual spacing
		rowContainer.Add(space)
		rowContainer.Add(space)
		rowContainer.Add(space)

		// Iterate over artists to create artist cards in the current row
		for j := i; j < i+3 && j < len(artists); j++ {
			card := CreateCardGeneralInfo(artists[j], relation, myApp) // Create an artist card for the current artist
			rowContainer.Add(card)                                     // Add the card to the row

			// Add space between cards if it's not the last card in the row
			if j < i+2 && j < len(artists) {
				rowContainer.Add(space)
			}
		}

		columnContainer.Add(rowContainer)     // Add the row of cards to the column container
		artistsContainer.Add(columnContainer) // Add the column container to the artists container
	}

	artistsContainer.Refresh() // Refresh the artists container to display the changes

	// Reset the search result count label text
	searchResultCountLabel.SetText("")
}

// Filter initializes and displays the filter window.
func Filter(myApp fyne.App) {
	// Check if the filter window is already open
	if myWindow != nil {
		// Close the filter window if it's open
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
	// Initialize the application filters
	initializeFilters(myApp, artists, concerts)
}

// initializeFilters initializes the filter options based on artist and concert data.
func initializeFilters(myApp fyne.App, artists []Artist, concerts []Concert) {
	// Initialize minimum and maximum values for creation and first album release years
	minCreationYear := artists[0].CreationDate
	maxCreationYear := artists[0].CreationDate

	// Function to parse year from date string (format: "DD-MM-YYYY")
	parseYear := func(dateStr string) int {
		parts := strings.Split(dateStr, "-")
		year, _ := strconv.Atoi(parts[2])
		return year
	}

	minFirstAlbumYear := parseYear(artists[0].FirstAlbum)
	maxFirstAlbumYear := parseYear(artists[0].FirstAlbum)

	// Iterate through artists to update min/max creation and first album years
	for _, artist := range artists {
		// Update min and max values for creation years
		if artist.CreationDate < minCreationYear {
			minCreationYear = artist.CreationDate
		}
		if artist.CreationDate > maxCreationYear {
			maxCreationYear = artist.CreationDate
		}

		// Parse year from FirstAlbum date string
		albumYear := parseYear(artist.FirstAlbum)
		if albumYear < minFirstAlbumYear {
			minFirstAlbumYear = albumYear
		}
		if albumYear > maxFirstAlbumYear {
			maxFirstAlbumYear = albumYear
		}
	}

	// Initialize available concert locations
	concertLocations := make([]string, 0)
	locationsMap := make(map[string]bool)
	for _, artist := range artists {
		// Update min and max values for creation years
		if artist.CreationDate < minCreationYear {
			minCreationYear = artist.CreationDate
		}
		if artist.CreationDate > maxCreationYear {
			maxCreationYear = artist.CreationDate
		}
		// Update min and max values for first album years
		albumYear := parseYear(artist.FirstAlbum)
		if albumYear < minFirstAlbumYear {
			minFirstAlbumYear = albumYear
		}
		if albumYear > maxFirstAlbumYear {
			maxFirstAlbumYear = albumYear
		}
		// Find unique concert locations
		for _, concert := range concerts {
			// Concatenate concert locations into a single string
			locationStr := strings.Join(concert.Locations, ", ")
			if _, found := locationsMap[locationStr]; !found {
				formattedLocation := formatLocation(locationStr)
				concertLocations = append(concertLocations, formattedLocation)
				locationsMap[locationStr] = true
			}
		}
	}

	// Create widgets for filters
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

	// Configure the selection of the number of members
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
		// Set the initial state of the check box based on selectedNumMembers
		for _, selectedOption := range selectedNumMembers {
			if selectedOption == option {
				check.SetChecked(true)
				break
			}
		}
		numMembersBox.Add(check)
	}

	// Hide the number of members selection by default
	numMembersBox.Hide()

	// Select available concert locations
	locationsSelect = widget.NewSelect(concertLocations, func(selected string) {})

	// Set initial values for the filters
	radioSoloGroup.SetSelected(selectedRadioValue)
	numMembersCheck.SetSelected(selectedNumMembers)
	locationsSelect.SetSelected(selectedLocationValue)
	numMembersCheck.SetSelected(selectedNumMembers)
	CreationDateRange.SetValue(savedCreationRange)
	firstAlbumDateRange.SetValue(savedFirstAlbumRange)
	numMembersCheck.SetSelected(selectedNumMembers)

	// Save initial filter values
	savedCreationRange = CreationDateRange.Value
	savedFirstAlbumRange = firstAlbumDateRange.Value
	savedNumMembers = selectedNumMembers

	// Create the filters window
	myWindow = myApp.NewWindow("Groupie Tracker GUI Filters")
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.SetFixedSize(true)

	// Configure reset and apply filter buttons
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

	// Create labels for creation date and first album date ranges
	CreationDateRangeLabel := widget.NewLabel(fmt.Sprintf("Creation Date Range: %d - %d", minCreationYear, maxCreationYear))
	firstAlbumDateRangeLabel := widget.NewLabel(fmt.Sprintf("First Album Date Range: %d - %d", minFirstAlbumYear, maxFirstAlbumYear))

	// Function to update labels for date ranges
	updateLabels := func() {
		creationRange := int(CreationDateRange.Value)
		firstAlbumRange := int(firstAlbumDateRange.Value)

		CreationDateRangeLabel.SetText(fmt.Sprintf("Creation Date Range: %d - %d", creationRange, maxCreationYear))
		firstAlbumDateRangeLabel.SetText(fmt.Sprintf("First Album Date Range: %d - %d", firstAlbumRange, maxFirstAlbumYear))
	}

	// Update labels when slider values change
	CreationDateRange.OnChanged = func(value float64) {
		updateLabels()
	}

	firstAlbumDateRange.OnChanged = func(value float64) {
		updateLabels()
	}

	// Create container for filter widgets
	filtersContainer := container.NewVBox(
		reset,
		CreationDateRangeLabel, CreationDateRange,
		firstAlbumDateRangeLabel, firstAlbumDateRange,
		radioSoloGroup,
		numMembersBox,
		locationsSelect,
		applyButton,
	)

	// Configure main window with filter widgets
	myWindow.SetContent(filtersContainer)
	myWindow.CenterOnScreen()
	myWindow.Show()
	windowOpened = true
}

// applyFilter function applies the selected filters and saves them for future use.
func applyFilter() saveFilter {
	// Store the selected values in their corresponding variables
	selectedRadioValue = radioSoloGroup.Selected
	selectedLocationValue = locationsSelect.Selected

	// Save the selected members in savedNumMembers
	savedNumMembers = selectedNumMembers

	// Save the selected date ranges in savedCreationRange and savedFirstAlbumRange
	savedCreationRange = CreationDateRange.Value
	savedFirstAlbumRange = firstAlbumDateRange.Value

	// Save the selected values in savedFilter
	savedFilter = saveFilter{
		RadioSelected:      selectedRadioValue,
		NumMembersSelected: selectedNumMembers,
		LocationSelected:   selectedLocationValue,
		CreationRange:      savedCreationRange,
		FirstAlbumRange:    savedFirstAlbumRange,
	}

	// Print the selected values to the console
	fmt.Printf("Radio selected: %s, Selected members: %v, Selected location: %s, savedCreationRange: %f, savedFirstAlbumRange: %f\n", selectedRadioValue, selectedNumMembers, selectedLocationValue, savedCreationRange, savedFirstAlbumRange)

	// Reset the location selection when applying the filter
	selectedLocationValue = "" // Reset location selection

	return savedFilter // Return the saved filters
}

// SecondPage function displays detailed information about an artist on the second page.
func SecondPage(artist Artist, relation Relation, myApp fyne.App) {
	myWindow := myApp.NewWindow("Information - " + artist.Name)

	logo, _ := fyne.LoadResourceFromPath("public/img/logo.png")
	myWindow.SetIcon(logo)

	// Load artist's image from URL
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

	// Create an image from the loaded data
	image := canvas.NewImageFromReader(strings.NewReader(string(imageData)), "image/jpeg")
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(320, 320))
	image.Resize(fyne.NewSize(220, 220))

	// Create labels for artist's information
	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	yearLabel := widget.NewLabel(fmt.Sprintf("Year Started: %d", artist.CreationDate))
	debutAlbumLabel := widget.NewLabel(fmt.Sprintf("Debut Album: %s", artist.FirstAlbum))
	membersLabel := widget.NewLabel(fmt.Sprintf("Members: %s", strings.Join(artist.Members, ", ")))

	// Create container for concert information
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

	// Arrange content vertically
	content := container.NewVBox(
		image,
		nameLabel,
		yearLabel,
		debutAlbumLabel,
		membersLabel,
		scrollContainer,
	)

	myWindow.SetContent(content) // Set the content of the window
	myWindow.Show()              // Show the window
}

// formatLocation function formats the location string.
func formatLocation(location string) string {
	// Remove dashes and underscores
	location = strings.ReplaceAll(location, "_", " ")
	// Capitalize the first letter of each word
	titleCase := cases.Title(language.English)
	location = titleCase.String(location)

	// Find the country name after the first dash
	parts := strings.Split(location, "-")
	if len(parts) > 1 {
		// Add the country name in parentheses
		location = fmt.Sprintf("%s (%s)", parts[0], parts[1])
	}

	return location
}

// formatDate function formats the date string.
func formatDate(date string) string {
	// Convert from "DD-MM-YYYY" to "MM/DD/YYYY" format
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return date // Return the date as is if the format is incorrect
	}
	return fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
}
