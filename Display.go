package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strconv"
	"strings"
)

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
