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

/********************************************************************************/
/*********************************** DISPLAY ************************************/
/********************************************************************************/

func CreateCardGeneralInfo(artist Artist, relation Relation, myApp fyne.App) fyne.CanvasObject {
	response, err := http.Get(artist.Image)
	if err != nil {
		fmt.Println("Failed to load image:", err)
		return nil
	}
	defer response.Body.Close()

	// Lire les données de l'image
	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		// Gérer l'erreur lors de la lecture des données de l'image
		fmt.Println("Failed to read image data:", err)
		return nil
	}

	// Obtenir le type de fichier de l'image à partir de l'URL
	parts := strings.Split(artist.Image, ".")
	fileType := parts[len(parts)-1]

	// Créer une image à partir des données lues et du type de fichier
	image := canvas.NewImageFromReader(strings.NewReader(string(imageData)), fileType)

	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(120, 120))
	image.Resize(fyne.NewSize(120, 120))

	r, g, b, a := getAverageColor(image)

	// Création du fond avec la couleur moyenne
	background := canvas.NewRectangle(color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
	background.SetMinSize(fyne.NewSize(300, 300))
	background.Resize(fyne.NewSize(296, 296))
	background.CornerRadius = 20

	// Création du bouton "Plus d'informations"
	button := widget.NewButton("          Plus d'informations          ", func() {
		SecondPage(artist, relation, myApp)
	})

	// Création du bouton de like
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

	// Charger l'icône initiale du bouton en fonction de l'état initial du favori
	likeButton.SetIcon(LoadImageResource(likeIcon))

	// Création du conteneur des boutons
	buttonsContainer := container.NewHBox(
		widget.NewLabel("  "),
		container.NewBorder(nil, layout.NewSpacer(), nil, likeButton, button),
	)

	// Création du label du nom de l'artiste
	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Création du label de l'année de début
	yearLabel := widget.NewLabel(fmt.Sprintf("%d", artist.CreationDate))

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

func GenerateSearchSuggestions(text string, scrollContainer *fyne.Container, artists []Artist, relation Relation, myApp fyne.App, limit int) int {
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
						SecondPage(a, relation, myApp)
					}
				}(artist))
				artistButton.Importance = widget.LowImportance
				scrollContainer.Add(artistButton)
			}

			// Vérifier si le texte de recherche correspond à l'année de commencement de l'artiste
			if strconv.Itoa(artist.CreationDate) == text {
				// Incrémenter le compteur et ajouter un bouton d'artiste avec l'année de commencement au conteneur de défilement
				count++
				artistButton := widget.NewButton(artist.Name+" (Year Started: "+text+")", func(a Artist) func() {
					return func() {
						SecondPage(a, relation, myApp)
					}
				}(artist))
				artistButton.Importance = widget.LowImportance
				scrollContainer.Add(artistButton)
			}

			// Extraction des années de la string FirstAlbum au format "DD-MM-YYYY" pour comparer
			albumYearParts := strings.Split(artist.FirstAlbum, "-")
			if len(albumYearParts) == 3 && albumYearParts[2] == text {
				// Incrémenter le compteur et ajouter un bouton d'artiste avec l'année de l'album de début au conteneur de défilement
				count++
				artistButton := widget.NewButton(artist.Name+" (Debut Album: "+albumYearParts[2]+")", func(a Artist) func() {
					return func() {
						SecondPage(a, relation, myApp)
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
								SecondPage(a, relation, myApp)
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
				for _, location := range concert.Locations {
					if strings.Contains(strings.ToLower(string(location)), strings.ToLower(text)) {
						count++
						artistButton := widget.NewButton(artist.Name+" (Concert Location: "+string(location)+")", func(a Artist) func() {
							return func() {
								SecondPage(a, relation, myApp)
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
	return count
}

func Recherche(searchBar *widget.Entry, scrollContainer *fyne.Container, artists []Artist, relation Relation, myApp fyne.App) {
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
				strconv.Itoa(artist.CreationDate) == searchText ||
				strings.Contains(artist.FirstAlbum, searchText) || // Adjusted for string comparison
				checkMemberName(artist.Members, searchText) ||
				checkConcertLocation(artist.NextConcerts, searchText) { // Adjusted to use artist.NextConcerts
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
	if filter.CreationRange > 0 && float64(artist.CreationDate) < filter.CreationRange {
		return false
	}

	// Vérifier si l'année de sortie du premier album de l'artiste est dans la plage sélectionnée par l'utilisateur
	if filter.FirstAlbumRange > 0 {
		albumYear := parseYearFromAlbum(artist.FirstAlbum) // You would define this function based on your date format
		if float64(albumYear) < filter.FirstAlbumRange {
			return false
		}
	}
	return true
}

func artistHasConcertLocation(artist Artist, location string) bool {
	// Vérifier si l'artiste a un concert à l'emplacement spécifié
	for _, concert := range artist.NextConcerts {
		for _, loc := range concert.Locations {
			if strings.EqualFold(string(loc), location) { // Vérifier sans tenir compte de la casse
				return true
			}
		}
	}
	return false
}

func parseYearFromAlbum(albumDate string) int {
	parts := strings.Split(albumDate, "-")
	year, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0
	}
	return year
}
