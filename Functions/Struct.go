package Functions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

/********************************************************************************/
/************************************ TYPES *************************************/
/********************************************************************************/
type APIResponse struct {
	Index []Location `json:"index"`
}

type RelationsResponse struct {
	Index []Relation `json:"index"`
}

type LocationResponse struct {
	Index []Location `json:"index"`
}

type Concert struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     []string `json:"dates"`
}

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

type saveFilter struct {
	RadioSelected      string   // RadioSelected stocke le type d'artiste sélectionné (Solo ou Groupe).
	NumMembersSelected []string // NumMembersSelected contient les nombres de membres sélectionnés.
	LocationSelected   string   // LocationSelected stocke l'emplacement de concert sélectionné.
	CreationRange      float64  // CreationRange stocke la plage de dates de création sélectionnée.
	FirstAlbumRange    float64  // FirstAlbumRange stocke la plage de dates du premier album sélectionnée.
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Dates struct {
	Dates []string `json:"dates"`
	Date  string
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	DatesURL  string   `json:"dates"`
}

/********************************************************************************/
/********************************* VARIABLES ************************************/
/********************************************************************************/

var (
	minCreationYear     int                // Année de commencement minimale parmi tous les artistes
	maxCreationYear     int                // Année de commencement maximale parmi tous les artistes
	minFirstAlbumYear   int                // Année du premier album la plus ancienne parmi tous les artistes
	maxFirstAlbumYear   int                // Année du premier album la plus récente parmi tous les artistes
	concertLocations    []string           // Liste des emplacements de concerts uniques
	YearStartedRange    *widget.Slider     // Curseur pour sélectionner la plage d'années de commencement
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
