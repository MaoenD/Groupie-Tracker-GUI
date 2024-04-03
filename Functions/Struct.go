package Functions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// APIResponse represents the response structure for API requests.
type APIResponse struct {
	Index []Location `json:"index"`
}

// RelationsResponse represents the response structure for relations API requests.
type RelationsResponse struct {
	Index []Relation `json:"index"`
}

// LocationResponse represents the response structure for location API requests.
type LocationResponse struct {
	Index []Location `json:"index"`
}

// Concert represents information about a concert, including its ID, locations, and dates.
type Concert struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     []string `json:"dates"`
}

// Artist represents information about an artist, including their ID, image, name, members, creation date,
// first album, concert locations URL, concert dates URL, relations URL, last concert, next concerts, and favorite status.
type Artist struct {
	ID              int       `json:"id"`
	Image           string    `json:"image"`
	Name            string    `json:"name"`
	Members         []string  `json:"members"`
	CreationDate    int       `json:"creationDate"`
	FirstAlbum      string    `json:"firstAlbum"`
	LocationsURL    string    `json:"locations"`
	ConcertDatesURL string    `json:"concertDates"`
	RelationsURL    string    `json:"relations"`
	LastConcert     Concert   `json:"lastConcert"`
	NextConcerts    []Concert `json:"nextConcerts"`
	Favorite        bool
}

// saveFilter represents the saved filters for artist search.
type saveFilter struct {
	RadioSelected      string   // RadioSelected stores the selected artist type (Solo or Group).
	NumMembersSelected []string // NumMembersSelected contains the selected number of members.
	LocationSelected   string   // LocationSelected stores the selected concert location.
	CreationRange      float64  // CreationRange stores the selected creation date range.
	FirstAlbumRange    float64  // FirstAlbumRange stores the selected first album date range.
}

// Relation represents information about a relation, including its ID and dates and locations.
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Dates represents a list of dates.
type Dates struct {
	Dates []string `json:"dates"`
	Date  string
}

// Location represents information about a location, including its ID, locations, and dates URL.
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	DatesURL  string   `json:"dates"`
}

var (
	minCreationYear     int                // Minimum creation year among all artists
	maxCreationYear     int                // Maximum creation year among all artists
	minFirstAlbumYear   int                // Oldest first album year among all artists
	maxFirstAlbumYear   int                // Newest first album year among all artists
	concertLocations    []string           // List of unique concert locations
	YearStartedRange    *widget.Slider     // Slider to select the range of creation years
	firstAlbumDateRange *widget.Slider     // Slider to select the range of first album years
	radioSoloGroup      *widget.RadioGroup // Radio group to select between Solo and Group
	numMembersCheck     *widget.CheckGroup // Check group to select number of members
	numMembersBox       *fyne.Container    // Container to display number of members checkboxes
	locationsSelect     *widget.Select     // Selector to choose concert location
	myWindow            fyne.Window        // Application window
	windowOpened        bool               // Indicates if the window is opened

	selectedRadioValue    string   // Selected value in the radio group
	selectedNumMembers    []string // Selected values in the check group
	selectedLocationValue string   // Selected value in the location selector
	savedCreationRange    float64  // Saved selected creation year range
	savedFirstAlbumRange  float64  // Saved selected first album year range
	savedNumMembers       []string // Saved selected number of members
)

var savedFilter saveFilter
