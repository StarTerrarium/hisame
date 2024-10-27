package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// AnimeListPage represents the page displaying the user's anime list.
type AnimeListPage struct {
	content fyne.CanvasObject
}

// NewAnimeListPage creates a new instance of AnimeListPage.
func NewAnimeListPage() *AnimeListPage {
	alp := &AnimeListPage{}
	alp.content = alp.buildContent()
	return alp
}

// Content returns the root content object of the AnimeListPage.
func (alp *AnimeListPage) Content() fyne.CanvasObject {
	return alp.content
}

// buildContent constructs the UI elements for the AnimeListPage.
func (alp *AnimeListPage) buildContent() fyne.CanvasObject {
	// Placeholder content
	label := widget.NewLabel("Anime List Page - Content Coming Soon")
	content := container.NewCenter(label)
	return content
}
