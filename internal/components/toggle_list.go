package components

import (
	"gorim/internal/state"
	"gorim/internal/types"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewModList(items map[string]types.InternalMod, appState *state.AppState) *widget.List {
	var mods []types.InternalMod
	for _, mod := range items {
		mods = append(mods, mod)
	}

	getSortedKeys := func() []string {
		sort.Sort(types.ModByOrder(mods))
		keys := make([]string, 0, len(mods))
		for _, k := range mods {
			keys = append(keys, k.PackageId)
		}
		return keys
	}

	keys := getSortedKeys()
	list := widget.NewList(
		func() int {
			return len(mods)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(&widget.Check{}, widget.NewLabel("Mods"))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)

			key := keys[i]
			mod := items[key]

			label.SetText(mod.Name)

			check.SetChecked(mod.Enabled)
			check.OnChanged = func(checked bool) {
				appState.EnableMod(mod.PackageId, checked)
			}
		},
	)

	return list
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
