package state

import (
	"gorim/internal/types"
	"log"
	"slices"
)

type ModDelegate struct {
	PackageName string
	Enabled     bool
}

type AppState struct {
	RimWorldVersion   string
	KnownExpansion    []string
	Profiles          []string
	ModList           map[string]types.InternalMod
	DisplayedMods     map[string]types.InternalMod
	PluginList        []string
	ModEnabledChanges []func([]ModDelegate)
	Rules             types.CommunityRules
}

func (state *AppState) AddModStateWatcher(fn func([]ModDelegate)) {
	state.ModEnabledChanges = append(state.ModEnabledChanges, fn)
}

func (state *AppState) AddPlugin(plugin types.InternalPlugin) {
	if slices.Contains(state.PluginList, plugin.Name) {
		log.Println("plugin already in plugins list")
		return
	}
	plugins := append(state.PluginList, plugin.Name)
	state.PluginList = plugins
}

func (state *AppState) AddMod(mod types.InternalMod) {
	_, has := state.ModList[mod.PackageId]
	if has {
		log.Println("mod already in map")
		return
	}
	state.ModList[mod.PackageId] = mod
}

func (state *AppState) AddPlugins(plugins []types.InternalPlugin) {
	var addedPlugins []string
	for _, plugin := range plugins {
		if !slices.Contains(state.PluginList, plugin.Name) {
			addedPlugins = append(addedPlugins, plugin.Name)
		}
	}
	state.PluginList = append(state.PluginList, addedPlugins...)
}

func (state *AppState) AddMods(mods []types.InternalMod) {
	var addedMods []types.InternalMod
	for _, mod := range mods {
		if _, has := state.ModList[mod.PackageId]; !has {
			state.ModList[mod.PackageId] = mod
			addedMods = append(addedMods, mod)
		}
	}
}

// enable mod functionality
// if we are enabling, update its status and add it to the ModList
// if we are disabling, delete it from the modlist
func (state *AppState) EnableMod(name string, enabled bool) {
	state.enableMod(name, enabled, true)
}

func (state *AppState) enableMod(name string, enabled bool, delegate bool) {
	//check if we are subscribed to it
	mod, ok := state.ModList[name]
	if !ok {
		log.Println("Unable to find mod : ", name)
		return
	}

	// update the state in modlist var
	mod.Enabled = enabled
	state.ModList[name] = mod

	//get index of mod
	index := slices.Index(state.PluginList, name)

	//if removing
	if !enabled {
		index := slices.Index(state.PluginList, name)
		if index == -1 {
			log.Println("cant find mod")
			return
		}
		state.PluginList = slices.Delete(state.PluginList, index, index+1)
	} else {
		//only add a new one if we dont have this in there - it goes at the end
		if index == -1 {
			state.PluginList = append(state.PluginList, name)
		}
	}
	if delegate {
		d := ModDelegate{PackageName: name, Enabled: enabled}
		state.runModChangeDelegate([]ModDelegate{d})
	}
}

func (state *AppState) EnableAll(enable bool) {
	var d []ModDelegate
	for _, mod := range state.ModList {
		state.enableMod(mod.PackageId, enable, false)
		d = append(d, ModDelegate{PackageName: mod.Name, Enabled: enable})
	}

	state.runModChangeDelegate(d)
}

func (state *AppState) runModChangeDelegate(mods []ModDelegate) {
	for _, delegate := range state.ModEnabledChanges {
		delegate(mods)
	}
}

func (state *AppState) SwapPlugin(curPos int, newPos int) {
	//create temp arr
	t := state.PluginList
	//keep within bounds
	newPos = max(0, min(newPos, len(t)-1))
	if newPos == curPos {
		return
	}

	//delete the item at its current position, and insert it into new pos
	cur := t[curPos]
	t = slices.Delete(t, curPos, curPos+1)
	t = slices.Insert(t, newPos, cur)
	state.PluginList = t
}
