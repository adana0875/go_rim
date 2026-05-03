package pages

import (
	"encoding/xml"
	"gorim/internal/state"
	"gorim/internal/types"
	"log"
	"os"

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
	log.Println("Appling changes to a new ModsConfig.xml")

	newMods := types.ModsConfig{Version: state.RimWorldVersion, ActiveMods: state.PluginList, KnownExpansions: state.KnownExpansion}

	file, err := os.Create("newmodsconfig.xml")
	if err != nil {
		log.Println("Failed to write file: ", err)
		return
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", " ")

	if err := encoder.Encode(newMods); err != nil {
		log.Println("Failed to write modsconfig to file: ", err)
		return
	}
}
