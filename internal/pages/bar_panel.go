package pages

import (
	"fmt"
	"gorim/internal/state"
	"log"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewBarPanel(appState *state.AppState) *fyne.Container {
	profileSelect := newProfileChoices(appState)
	addButton := widget.NewButton("+", func() {
		appState.Profiles = append(appState.Profiles, fmt.Sprintf("Profile New %d", rand.Int()))
		profileSelect.Refresh()
		log.Println("New Profile")
	})
	mainContent := container.NewHBox(profileSelect, addButton)
	spacedContent := container.NewCenter(mainContent)
	return spacedContent
}

func newProfileChoices(appState *state.AppState) *widget.Select {
	combo := widget.NewSelect(appState.Profiles, func(value string) {
		log.Println("Clicked Profile ", value)
	})

	return combo
}
