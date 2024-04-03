package search

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// affichage liste éléments
func DisplayL[T any](items []T, convertFunc func(T) string) *widget.List {
	list := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			item := items[id]
			label.SetText(convertFunc(item))
		},
	)
	return list
}
