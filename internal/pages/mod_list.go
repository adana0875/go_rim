package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"gorim/internal/components"
	"gorim/internal/state"
	"gorim/internal/types"
	"image/color"
)

type InputParams struct {
	Modpath      string `json:"modPath"`
	WorkshopPath string `json:"workshopPath"`
}

// Main Modlist panel
func InputPanel(opts InputParams, state *state.AppState) *fyne.Container {

	modContainer := components.NewModList(state.ModList)

	//sub to state updates refreshing
	state.AddModWatcher(func(mods []types.InternalMod) {
		fyne.Do(func() {
			modContainer.Refresh()
		})
	})

	sizer := canvas.NewRectangle(color.Transparent)
	sizer.SetMinSize(fyne.NewSize(500, 500))

	return container.New(layout.NewStackLayout(), sizer, container.NewPadded(modContainer))
}
