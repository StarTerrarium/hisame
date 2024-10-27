package ui

import "fyne.io/fyne/v2"

// Page defines the interface that all pages must implement.
type Page interface {
	Content() fyne.CanvasObject
}
