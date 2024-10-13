package ui

import (
	"fyne.io/fyne/v2"
)

// ScreenManager acts as a central management tool for changing between the main screens available in the app.
type ScreenManager struct {
	window fyne.Window
}

func NewScreenManager(w fyne.Window) *ScreenManager {
	return &ScreenManager{
		window: w,
	}
}
