package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorim/internal/components"
	"gorim/internal/state"
	"gorim/internal/types"
	"image/color"
	"strings"
	"unicode/utf8"
)

type InputParams struct {
	Modpath      string `json:"modPath"`
	WorkshopPath string `json:"workshopPath"`
}

func getStringChar(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

// Main Modlist panel
func InputPanel(opts InputParams, appState *state.AppState) *fyne.Container {

	modContainer, refreshList := components.NewModList(appState.DisplayedMods, appState)

	appState.AddModStateWatcher(func(mod []state.ModDelegate) {
		fyne.Do(func() {
			modContainer.Refresh()
		})

	})

	sizer := canvas.NewRectangle(color.Transparent)
	sizer.SetMinSize(fyne.NewSize(500, 500))

	checkAll := widget.NewCheck("All", func(value bool) {
		appState.EnableAll(value)
	})

	search := widget.NewEntry()
	search.SetPlaceHolder("mod name...")
	search.OnChanged = func(text string) {
		searchTerm := strings.ToLower(text)
		var filterList map[string]types.InternalMod = map[string]types.InternalMod{}
		if len([]byte(searchTerm)) > 0 && ([]byte(searchTerm)[0] == []byte("*")[0]) {
			for _, item := range appState.ModList {
				if strings.Contains(strings.ToLower(item.Name), getStringChar(searchTerm)) {
					filterList[item.PackageId] = item
				}
				appState.DisplayedMods = filterList
			}

		} else { //starts with
			for _, item := range appState.ModList {
				if strings.HasPrefix(strings.ToLower(item.Name), searchTerm) {
					filterList[item.PackageId] = item
				}
			}
			appState.DisplayedMods = filterList
		}
		refreshList()
	}
	infoPanel := container.NewVBox(widget.NewLabel("Mods"), container.NewHBox(checkAll, search))
	return container.NewBorder(infoPanel, nil, nil, nil, container.NewPadded(modContainer))
}
