package pages

import (
	"gorim/internal/components"
	"gorim/internal/state"
	"gorim/internal/types"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// func addModToPluginList(mod types.InternalMod, state *state.AppState, modList *widget.List) {
// 	if modList == nil {
// 		log.Println("Cannot access list")
// 	}
// 	fyne.Do(func() {
// 		state.PluginList = append(state.PluginList, mod)
// 		modList.Refresh()
// 	})
// }
//
// func addPlugins(mods []types.InternalMod, state *state.AppState, modList *widget.List) {
// 	if modList == nil {
// 		log.Println("Cannot access list")
// 	}
// 	fyne.Do(func() {
// 		newList := append(state.PluginList, mods...)
// 		state.PluginList = newList
// 		modList.Refresh()
// 	})
// }

func NewPluginList(state *state.AppState) *fyne.Container {

	modContainer := components.NewPluginList(state.PluginList)

	state.AddPluginWatcher(func(mods []types.InternalPlugin) {
		fyne.Do(func() {
			modContainer.Refresh()
		})
	})

	mainContainer := container.NewPadded(modContainer)

	sizer := canvas.NewRectangle(color.Transparent)
	sizer.SetMinSize(fyne.NewSize(500, 500))

	return container.New(layout.NewStackLayout(), sizer, mainContainer)
}
