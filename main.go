package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")

	myWindow.Resize(fyne.NewSize(800, 600))

	myWindow.ShowAndRun()
}
