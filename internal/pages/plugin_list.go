package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"gorim/internal/components"
	"gorim/internal/state"
	"image/color"
)

func NewPluginList(appState *state.AppState) *fyne.Container {

	modContainer := components.NewPluginList(appState)

	appState.AddModStateWatcher(func([]state.ModDelegate) {
		fyne.Do(func() {
			modContainer.Refresh()
		})
	})

	mainContainer := container.NewPadded(modContainer)

	sizer := canvas.NewRectangle(color.Transparent)
	sizer.SetMinSize(fyne.NewSize(500, 500))

	return container.New(layout.NewStackLayout(), sizer, mainContainer)
}
