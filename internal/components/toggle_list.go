package components

import (
	"gorim/internal/types"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewModList(items map[string]types.InternalMod) *widget.List {
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return items[keys[i]].Order < items[keys[j]].Order
	})

	list := widget.NewList(
		func() int {
			return len(items)
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
				mod.Enabled = checked
			}
		},
	)

	return list
}

func NewPluginList(items map[string]types.InternalPlugin) *widget.List {

	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return items[keys[i]].Order < items[keys[j]].Order
	})

	list := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(&widget.Check{}, widget.NewLabel("Plugins"))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)

			key := keys[i]
			mod := items[key]

			label.SetText(mod.Name)

			check.OnChanged = func(checked bool) {
				mod.Enabled = checked
			}
			check.SetChecked(mod.Enabled)
		},
	)

	return list
}
