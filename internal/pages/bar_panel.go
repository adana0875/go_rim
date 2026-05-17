package pages

import (
	"fmt"
	"gorim/internal/state"
	"gorim/internal/types"
	"log"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//top level panel which holds the profile selection

func NewBarPanel(appState *state.AppState) *fyne.Container {
	profileSelect, addFunc := newProfileChoices(appState)
	addButton := widget.NewButton("+", func() {
		newName := fmt.Sprintf("new_profile_%d", rand.Int())
		newProfile := types.Profile{Name: newName, PluginList: []string{}}
		appState.Profiles = append(appState.Profiles, newProfile)
		addFunc(newName)
		log.Println("New Profile, appstate has profiles: ", appState.Profiles)
	})
	mainContent := container.NewHBox(profileSelect, addButton)
	spacedContent := container.NewCenter(mainContent)
	return spacedContent
}

func newProfileChoices(appState *state.AppState) (*widget.Select, func(newOption string)) {
	var names []string = make([]string, len(appState.Profiles))
	for idx, item := range appState.Profiles {
		names[idx] = item.Name
	}
	log.Println("Profiles: ", names)
	combo := widget.NewSelect(names, func(value string) {
		log.Println("Clicked Profile ", value)
		//call in to change profile
		appState.ChangeProfile(value)
	})

	addOption := func(newOption string) {
		combo.Options = append(combo.Options, newOption)
		combo.Refresh()
	}

	combo.SetSelected(appState.ActiveProfile.Name)

	return combo, addOption
}
