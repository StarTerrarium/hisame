package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type StatusBar struct {
	content    fyne.CanvasObject
	leftLabel  *widget.Label
	rightLabel *widget.Label
}

func NewStatusBar() *StatusBar {
	sb := &StatusBar{
		leftLabel:  widget.NewLabel(""),
		rightLabel: widget.NewLabel(""),
	}
	sb.content = sb.buildContent()
	return sb
}

func (sb *StatusBar) Content() fyne.CanvasObject {
	return sb.content
}

func (sb *StatusBar) buildContent() fyne.CanvasObject {
	leftContainer := container.NewHBox(sb.leftLabel)
	rightContainer := container.NewHBox(sb.rightLabel)

	// Spacer between left and right
	spacer := layout.NewSpacer()

	// Combine into a single container
	statusBar := container.NewHBox(leftContainer, spacer, rightContainer)
	return statusBar
}

func (sb *StatusBar) UpdateLeft(text string) {
	sb.leftLabel.SetText(text)
}

func (sb *StatusBar) UpdateRight(text string) {
	sb.rightLabel.SetText(text)
}
