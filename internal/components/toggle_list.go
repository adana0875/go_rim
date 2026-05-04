package components

import (
	"gorim/internal/state"
	"gorim/internal/types"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewModList(items map[string]types.InternalMod, appState *state.AppState) (*widget.List, func()) {
	var keys []string

	//get the mods from the internal mods
	getKeys := func() []string {
		var mods []types.InternalMod
		for _, mod := range appState.DisplayedMods {
			mods = append(mods, mod)
		}

		sort.Stable(types.ModByOrder(mods))
		result := make([]string, 0, len(mods))
		for _, k := range mods {
			result = append(result, k.PackageId)
		}
		return result
	}

	keys = getKeys()

	list := widget.NewList(
		func() int {
			return len(appState.DisplayedMods)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(&widget.Check{}, widget.NewLabel("Mods"), widget.NewLabel("PluginName"))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)
			pluginLabel := box.Objects[2].(*widget.Label)

			key := keys[i]
			mod := appState.ModList[key]

			label.SetText(mod.Name)
			pluginLabel.SetText(mod.PackageId)

			check.OnChanged = nil
			check.SetChecked(mod.Enabled)
			check.OnChanged = func(checked bool) {
				appState.EnableMod(mod.PackageId, checked)
			}
		},
	)

	refreshList := func() {
		keys = getKeys()
		list.Refresh()
	}

	return list, refreshList
}

// declaration a plugin list
// should follow the app state of plugins
func NewPluginList(appState *state.AppState) *widget.List {
	//get the active plugins - probs not necessary
	var activeItems []types.InternalMod
	for _, plugin := range appState.ModList {
		if plugin.Enabled {
			activeItems = append(activeItems, plugin)
		}
	}

	list := widget.NewList(
		func() int {
			return len(appState.PluginList)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Plugins"))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)

			label.SetText(appState.PluginList[i])
		},
	)

	return list

}
