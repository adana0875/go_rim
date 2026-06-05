package state

import "gorim/internal/types"

func Initialize() AppState {
	return AppState{
		ModList:       map[string]types.InternalMod{},
		DisplayedMods: map[string]types.InternalMod{},
	}
}
