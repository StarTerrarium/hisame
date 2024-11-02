package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

type NavigationBar struct {
	content fyne.CanvasObject

	// Buttons
	animeButton    *widget.Button
	mangaButton    *widget.Button
	searchButton   *widget.Button
	settingsButton *widget.Button
	logoutButton   *widget.Button
}

func NewNavigationBar() *NavigationBar {
	nb := &NavigationBar{}
	nb.content = nb.buildContent()
	return nb
}

func (nb *NavigationBar) Content() fyne.CanvasObject {
	return nb.content
}

func (nb *NavigationBar) buildContent() fyne.CanvasObject {
	// Left side buttons
	nb.animeButton = widget.NewButton("Anime", func() {
		logrus.Debug("Anime navigation button clicked")
		getScreenManager().ShowPage(NewAnimeListPage())
	})
	nb.mangaButton = widget.NewButton("Manga", func() {
		logrus.Debug("Manga navigation button clicked")
		// Implement Manga page navigation when ready
	})
	nb.searchButton = widget.NewButton("Search/Add", func() {
		logrus.Debug("Search navigation button clicked")
		// Implement Search/Add page navigation when ready
	})

	// Right side buttons
	nb.settingsButton = widget.NewButton("Settings", func() {
		logrus.Debug("Settings navigation button clicked")
		// Implement Settings page navigation when ready
	})
	nb.logoutButton = widget.NewButton("Logout", func() {
		logrus.Debug("Logout button clicked")
		getScreenManager().HandleLogout()
	})

	// Initially disable all buttons
	if !getScreenManager().isAuth {
		buttonsToDisable := []*widget.Button{nb.animeButton, nb.mangaButton, nb.searchButton, nb.settingsButton, nb.logoutButton}
		for _, btn := range buttonsToDisable {
			btn.Disable()
		}
	}

	// Left and right containers
	leftContainer := container.NewHBox(nb.animeButton, nb.mangaButton, nb.searchButton)
	rightContainer := container.NewHBox(nb.settingsButton, nb.logoutButton)

	// Spacer between left and right
	spacer := layout.NewSpacer()

	// Combine into a single container
	navBar := container.NewHBox(leftContainer, spacer, rightContainer)
	return navBar
}

func (nb *NavigationBar) UpdateAuthenticationState(isAuthenticated bool) {
	buttonsToToggle := []*widget.Button{nb.animeButton, nb.mangaButton, nb.searchButton, nb.settingsButton, nb.logoutButton}
	for _, btn := range buttonsToToggle {
		if isAuthenticated {
			btn.Enable()
		} else {
			btn.Disable()
		}
	}
}
