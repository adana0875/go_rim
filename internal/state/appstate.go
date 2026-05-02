package state

import (
	"gorim/internal/types"
	"log"
	"slices"
)

type AppState struct {
	ModList           map[string]types.InternalMod
	PluginList        []string
	ModlistChanges    []func([]types.InternalMod)
	PluginChanges     []func([]string)
	ModEnabledChanges []func(string, bool)
}

func (state *AppState) AddModWatcher(fn func([]types.InternalMod)) {
	state.ModlistChanges = append(state.ModlistChanges, fn)
}

func (state *AppState) AddPluginWatcher(fn func([]string)) {
	state.PluginChanges = append(state.PluginChanges, fn)
}

func (state *AppState) AddModStateWatcher(fn func(string, bool)) {
	state.ModEnabledChanges = append(state.ModEnabledChanges, fn)
}

func (state *AppState) AddPlugin(plugin types.InternalPlugin) {
	if slices.Contains(state.PluginList, plugin.Name) {
		log.Println("plugin already in plugins list")
		return
	}
	plugins := append(state.PluginList, plugin.Name)
	state.PluginList = plugins
	go state.runPluginDelegates([]string{plugin.Name})
}

func (state *AppState) AddMod(mod types.InternalMod) {
	_, has := state.ModList[mod.PackageId]
	if has {
		log.Println("mod already in map")
		return
	}
	state.ModList[mod.PackageId] = mod
	go state.runSubs([]types.InternalMod{mod})
}

func (state *AppState) AddPlugins(plugins []types.InternalPlugin) {
	var addedPlugins []string
	for _, plugin := range plugins {
		if !slices.Contains(state.PluginList, plugin.Name) {
			addedPlugins = append(addedPlugins, plugin.Name)
		}
	}
	state.PluginList = append(state.PluginList, addedPlugins...)
	go state.runPluginDelegates(addedPlugins)
}

func (state *AppState) AddMods(mods []types.InternalMod) {
	var addedMods []types.InternalMod
	for _, mod := range mods {
		if _, has := state.ModList[mod.PackageId]; !has {
			state.ModList[mod.PackageId] = mod
			addedMods = append(addedMods, mod)
		}
	}
	go state.runSubs(addedMods)
}

// enable mod functionality
// if we are enabling, update its status and add it to the ModList
// if we are disabling, delete it from the modlist
func (state *AppState) EnableMod(name string, enabled bool) {
	//check if we are subscribed to it
	mod, ok := state.ModList[name]
	if !ok {
		log.Println("Unable to find mod : ", name)
		return
	}

	// update the state in modlist var
	mod.Enabled = enabled
	state.ModList[name] = mod

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
		log.Println("add mod")
		state.PluginList = append(state.PluginList, name)
	}
	state.runModChangeDelegate(name, enabled)
}

func (state *AppState) runSubs(adding []types.InternalMod) {
	for _, delegate := range state.ModlistChanges {
		delegate(adding)
	}
}

func (state *AppState) runPluginDelegates(adding []string) {
	for _, delegate := range state.PluginChanges {
		delegate(adding)
	}
}

func (state *AppState) runModChangeDelegate(mod string, newState bool) {
	for _, delegate := range state.ModEnabledChanges {
		delegate(mod, newState)
	}
}
