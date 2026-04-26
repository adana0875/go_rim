package pages

import (
	"gorim/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewActionsPanel(appState *state.AppState) *fyne.Container {
	applyButton := widget.NewButton("Apply Changes", func() { applyChanges(appState) })

	mainContent := container.NewHBox(applyButton)
	spacedContent := container.NewCenter(mainContent)

	return spacedContent
}

func applyChanges(state *state.AppState) {

}
