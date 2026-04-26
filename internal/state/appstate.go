package state

import (
	"gorim/internal/types"
	"log"
)

type AppState struct {
	ModList        map[string]types.InternalMod
	PluginList     map[string]types.InternalPlugin
	ModlistChanges []func([]types.InternalMod)
	PluginChanges  []func([]types.InternalPlugin)
}

func (state *AppState) AddModWatcher(fn func([]types.InternalMod)) {
	state.ModlistChanges = append(state.ModlistChanges, fn)
}

func (state *AppState) AddPluginWatcher(fn func([]types.InternalPlugin)) {
	state.PluginChanges = append(state.PluginChanges, fn)
}

func (state *AppState) AddPlugin(plugin types.InternalPlugin) {
	_, has := state.PluginList[plugin.Name]
	if has {
		log.Println("plugin already in map")
		return
	}
	state.PluginList[plugin.Name] = plugin
	go state.runPluginDelegates([]types.InternalPlugin{plugin})
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
	var addedPlugins []types.InternalPlugin
	for _, plugin := range plugins {
		if _, has := state.PluginList[plugin.Name]; !has {
			state.PluginList[plugin.Name] = plugin
			addedPlugins = append(addedPlugins, plugin)
		}
	}
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

func (state *AppState) runSubs(adding []types.InternalMod) {
	for _, delegate := range state.ModlistChanges {
		delegate(adding)
	}
}

func (state *AppState) runPluginDelegates(adding []types.InternalPlugin) {
	for _, delegate := range state.PluginChanges {
		delegate(adding)
	}
}
