package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type MainScreen struct {
	navigationBar *NavigationBar
	statusBar     *StatusBar
	contentArea   *fyne.Container
}

func NewMainScreen(window fyne.Window) *MainScreen {
	ms := &MainScreen{
		contentArea: container.NewStack(),
	}

	ms.navigationBar = NewNavigationBar()
	ms.statusBar = NewStatusBar()

	ms.buildUI()
	return ms
}

func (ms *MainScreen) buildUI() {
	mainContainer := container.NewBorder(
		ms.navigationBar.Content(),
		ms.statusBar.Content(),
		nil, nil,
		ms.contentArea,
	)
	getScreenManager().window.SetContent(mainContainer)
}

func (ms *MainScreen) ShowPage(page Page) {
	ms.contentArea.Objects = []fyne.CanvasObject{page.Content()}
	ms.contentArea.Refresh()
}
