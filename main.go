package main

import (
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

// Artist struct to hold information about artists or groups
type Artist struct {
	Name        string
	Image       string
	YearStarted int
	DebutAlbum  time.Time
	Members     []string
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Menu - Groupie Tracker")

	// Exemple de données pour les artistes
	artists := []Artist{
		{Name: "Michael Jackson", Image: "public/michaeljackson.jpg", YearStarted: 1964, DebutAlbum: time.Date(1972, time.November, 13, 0, 0, 0, 0, time.UTC), Members: []string{"Michael Jackson"}},
		{Name: "Queen", Image: "public/queen.jpg", YearStarted: 1970, DebutAlbum: time.Date(1973, time.July, 13, 0, 0, 0, 0, time.UTC), Members: []string{"Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"}},
		{Name: "Pink Floyd", Image: "public/pinkfloyd.jpeg", YearStarted: 1965, DebutAlbum: time.Date(1967, time.August, 5, 0, 0, 0, 0, time.UTC), Members: []string{"Syd Barrett", "Roger Waters", "Richard Wright", "Nick Mason"}},
		{Name: "The Beatles", Image: "public/thebeatles.jpg", YearStarted: 1960, DebutAlbum: time.Date(1963, time.March, 22, 0, 0, 0, 0, time.UTC), Members: []string{"John Lennon", "Paul McCartney", "George Harrison", "Ringo Starr"}},
		{Name: "Elvis Presley", Image: "public/elvispresley.jpg", YearStarted: 1954, DebutAlbum: time.Date(1956, time.March, 23, 0, 0, 0, 0, time.UTC), Members: []string{"Elvis Presley"}},
		{Name: "The Rolling Stones", Image: "public/therollingstones.jpg", YearStarted: 1962, DebutAlbum: time.Date(1964, time.April, 17, 0, 0, 0, 0, time.UTC), Members: []string{"Mick Jagger", "Keith Richards", "Charlie Watts", "Ronnie Wood"}},
		{Name: "Led Zeppelin", Image: "public/ledzeppelin.jpg", YearStarted: 1968, DebutAlbum: time.Date(1969, time.January, 12, 0, 0, 0, 0, time.UTC), Members: []string{"Robert Plant", "Jimmy Page", "John Paul Jones", "John Bonham"}},
		{Name: "AC/DC", Image: "public/acdc.jpg", YearStarted: 1973, DebutAlbum: time.Date(1975, time.February, 17, 0, 0, 0, 0, time.UTC), Members: []string{"Angus Young", "Brian Johnson", "Phil Rudd", "Cliff Williams", "Stevie Young"}},
		{Name: "Nirvana", Image: "public/nirvana.jpg", YearStarted: 1987, DebutAlbum: time.Date(1989, time.June, 15, 0, 0, 0, 0, time.UTC), Members: []string{"Kurt Cobain", "Krist Novoselic", "Dave Grohl"}},
		{Name: "The Beach Boys", Image: "public/thebeachboys.jpg", YearStarted: 1961, DebutAlbum: time.Date(1962, time.October, 1, 0, 0, 0, 0, time.UTC), Members: []string{"Brian Wilson", "Mike Love", "Al Jardine", "Bruce Johnston", "David Marks"}},
		{Name: "The Who", Image: "public/thewho.jpg", YearStarted: 1964, DebutAlbum: time.Date(1965, time.December, 3, 0, 0, 0, 0, time.UTC), Members: []string{"Roger Daltrey", "Pete Townshend", "John Entwistle", "Keith Moon"}},
		{Name: "David Bowie", Image: "public/davidbowie.jpg", YearStarted: 1962, DebutAlbum: time.Date(1967, time.June, 1, 0, 0, 0, 0, time.UTC), Members: []string{"David Bowie"}},
		{Name: "Metallica", Image: "public/metallica.jpg", YearStarted: 1981, DebutAlbum: time.Date(1983, time.July, 25, 0, 0, 0, 0, time.UTC), Members: []string{"James Hetfield", "Lars Ulrich", "Kirk Hammett", "Robert Trujillo"}},
	}

	// Champ de recherche
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Artists...")

	// Zone d'affichage des résultats de recherches
	searchResults := container.NewVBox()

	// Bouton de recherche
	searchButton := widget.NewButton("Search", func() {
		// Récupérer le texte de la recherche
		searchText := searchBar.Text

		// Recherche d'artistes correspondants
		var foundArtists []Artist
		for _, artist := range artists {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchText)) {
				foundArtists = append(foundArtists, artist)
			}
		}

		// Afficher les résultats
		searchResultsObjects := make([]fyne.CanvasObject, 0)
		if len(foundArtists) > 0 {
			resultText := fmt.Sprintf("Artist found: %s\n\n", searchText)
			searchResultsObjects = append(searchResultsObjects, widget.NewLabel(resultText))
			for _, artist := range foundArtists {
				card := createCardGeneralInfo(artist)
				searchResultsObjects = append(searchResultsObjects, card)
			}
		} else {
			searchResultsObjects = append(searchResultsObjects, widget.NewLabel("No artist found: "+searchText))
		}

		searchResults.Objects = searchResultsObjects
		searchResults.Refresh()
	})

	// Ajouter l'écouteur d'événements clavier au champ de recherche
	searchBar.OnSubmitted = func(_ string) {
		searchButton.OnTapped()
	}

	// Création du conteneur pour afficher les artistes
	artistsContainer := container.NewVBox()

	// Ajouter les cartes des artistes par groupes de 3 dans des conteneurs de grille
	for i := 0; i < len(artists); i += 3 {
		// Créer un nouveau conteneur de grille pour chaque ligne d'artistes
		rowContainer := container.NewGridWithColumns(3)
		// Ajouter les trois artistes de cette ligne dans la grille
		for j := i; j < i+3 && j < len(artists); j++ {
			card := createCardGeneralInfo(artists[j])
			rowContainer.Add(card)
		}
		// Ajouter un espacement après chaque ligne
		artistsContainer.Add(rowContainer)
		artistsContainer.Add(layout.NewSpacer()) // Ajouter un espacement
	}

	// Créer un conteneur de défilement pour le conteneur principal
	scrollContainer := container.NewVScroll(artistsContainer)
	scrollContainer.SetMinSize(fyne.NewSize(1080, 720)) // Taille minimale pour activer le défilement

	// Création du conteneur principal avec la couleur de fond spécifiée
	content := container.NewVBox(
		searchBar,
		searchButton,
		searchResults,
		scrollContainer, // Utilisation du conteneur de défilement pour les artistes
	)

	// Centrer les cartes dans la fenêtre
	centeredContent := container.New(layout.NewCenterLayout(), content)

	// Créer un rectangle pour représenter le fond coloré
	background := canvas.NewRectangle(color.NRGBA{R: 0x5C, G: 0x64, B: 0x73, A: 0xFF})
	background.Resize(fyne.NewSize(1080, 720)) // Taille pour remplir toute la fenêtre

	// Créer un conteneur pour contenir le fond coloré
	backgroundContainer := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), background)

	// Ajouter le contenu principal au conteneur de fond
	backgroundContainer.Add(centeredContent)

	// Afficher la fenêtre
	myWindow.SetContent(backgroundContainer)
	myWindow.Resize(fyne.NewSize(1080, 720)) // ajustement size
	myWindow.ShowAndRun()
}

func createBlockContent() fyne.CanvasObject {
	// Chemin de l'image
	imagePath := "public/world_map.jpg"

	// Charger l'image
	image := canvas.NewImageFromFile(imagePath)

	// Vérifier les erreurs lors du chargement de l'image
	if image == nil {
		fmt.Println("Impossible de charger l'image:", imagePath)
		return nil
	}

	// Redimensionner l'image à la taille spécifiée
	image.Resize(fyne.NewSize(1000, 120))

	// Créer le texte
	title := widget.NewLabel("Geolocation feature")
	description := widget.NewLabel("Find out where and when your favorite artists will be performing around the globe.")
	description.Wrapping = fyne.TextWrapWord // Activer le wrapping du texte

	// Créer un conteneur pour organiser le texte
	textContainer := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		title,
		description,
	)

	// Créer un conteneur pour organiser l'image et le texte côte à côte
	content := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		image,
		textContainer,
	)

	// Créer un conteneur pour la carte de contenu avec un rectangle de contour arrondi
	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), content)
	cardContent.Resize(fyne.NewSize(1000, 120))

	// Créer un rectangle de contour avec des coins arrondis
	border := canvas.NewRectangle(color.Transparent) // Définir une couleur transparente pour le remplissage
	border.SetMinSize(fyne.NewSize(1000, 120))
	border.Resize(fyne.NewSize(996, 116)) // Redimensionner légèrement la bordure pour inclure les coins arrondis
	border.StrokeColor = color.Black      // Définir la couleur de la bordure
	border.StrokeWidth = 3                // Définir l'épaisseur de la bordure
	border.CornerRadius = 20              // Définir les coins arrondis

	// Ajouter le rectangle de contour à la carte de contenu
	cardContent.AddObject(border)

	return cardContent
}

func createCardGeneralInfo(artist Artist) fyne.CanvasObject {
	// Redimensionner l'image
	image := canvas.NewImageFromFile(artist.Image)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(120, 120))
	image.Resize(fyne.NewSize(120, 120))

	// Obtenir la couleur moyenne de l'image
	averageColor := getAverageColor(artist.Image)

	// Créer un rectangle coloré pour l'arrière-plan de la carte
	background := canvas.NewRectangle(averageColor)
	background.SetMinSize(fyne.NewSize(300, 300))
	background.Resize(fyne.NewSize(296, 296)) // Redimensionner légèrement la bordure pour inclure les coins arrondis
	background.CornerRadius = 20              // Définir les coins arrondis

	// Nom de l'artiste en gras et plus gros
	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Date de début de l'artiste en plus petit
	yearLabel := widget.NewLabel(fmt.Sprintf("%d", artist.YearStarted))

	// Créer un conteneur HBox pour afficher les labels avec un espace entre eux
	labelsContainer := container.NewHBox(
		nameLabel,
		yearLabel,
	)

	// Membres du groupe, s'il y en a, en plus petit
	var membersText string
	if len(artist.Members) == 1 {
		membersText = "Solo Artist"
	} else if len(artist.Members) > 0 {
		membersText = "Members:\n " + strings.Join(artist.Members, ", ")
	}
	membersLabel := widget.NewLabel(membersText)
	membersLabel.Wrapping = fyne.TextWrapWord // Activer le wrapping du texte

	// Créer le conteneur pour les informations sur l'artiste
	infoContainer := container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(), // Ajout d'un espace vertical
		image,
		labelsContainer,    // Placer les labels sur la même ligne avec un espace entre eux
		membersLabel,       // Afficher les membres du groupe
		layout.NewSpacer(), // Ajout d'un petit espace vertical
	)

	// Définir la taille fixe pour le conteneur d'informations
	infoContainer.Resize(fyne.NewSize(300, 180))

	// Créer le conteneur pour la carte de l'artiste
	cardContent := container.New(layout.NewBorderLayout(nil, nil, nil, nil), background, infoContainer)
	cardContent.Resize(fyne.NewSize(300, 300))

	// Créer un rectangle pour le contour avec des coins arrondis
	border := canvas.NewRectangle(color.Transparent) // Définir une couleur transparente pour le remplissage
	border.SetMinSize(fyne.NewSize(300, 300))
	border.Resize(fyne.NewSize(296, 296)) // Redimensionner légèrement la bordure pour inclure les coins arrondis
	border.StrokeColor = color.Black      // Définir la couleur de la bordure
	border.StrokeWidth = 3                // Définir l'épaisseur de la bordure
	border.CornerRadius = 20              // Définir les coins arrondis

	// Ajouter le rectangle de contour à la carte
	cardContent.AddObject(border)

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
