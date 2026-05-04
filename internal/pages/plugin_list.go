package pages

import (
	"gorim/internal/components"
	"gorim/internal/sorting"
	"gorim/internal/state"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NewPluginList(appState *state.AppState) *fyne.Container {

	var selectedPlugin int = -1

	modContainer := components.NewPluginList(appState)

	modContainer.OnSelected = func(item int) {
		selectedPlugin = item
	}

	appState.AddModStateWatcher(func([]state.ModDelegate) {
		fyne.Do(func() {
			modContainer.Refresh()
		})
	})

	mainContainer := container.NewPadded(modContainer)
	up := widget.NewButton(" Up ", func() {
		shiftAmount := selectedPlugin + -3
		appState.SwapPlugin(selectedPlugin, shiftAmount)
		modContainer.Refresh()
		modContainer.Select(shiftAmount)
	})

	down := widget.NewButton(" Down ", func() {
		shiftAmount := selectedPlugin + 2
		appState.SwapPlugin(selectedPlugin, shiftAmount)
		modContainer.Refresh()
		modContainer.Select(shiftAmount)
	})

	sort := widget.NewButton("Sort", func() {
		newList, err := sorting.TopoSortList(appState.ModList, appState.PluginList, appState.Rules)
		if err != nil {
			log.Println("Failed to sort: ", err)
		}
		appState.PluginList = newList
		fyne.Do(func() { modContainer.Refresh() })
	})
	buttons := container.NewHBox(up, down, sort)
	box := container.NewBorder(nil, buttons, nil, nil, mainContainer)

	sizer := canvas.NewRectangle(color.Transparent)
	sizer.SetMinSize(fyne.NewSize(500, 500))

	return container.New(layout.NewStackLayout(), sizer, box)
}
