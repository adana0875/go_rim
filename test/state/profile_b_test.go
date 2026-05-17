package state

import (
	"gorim/internal/state"
	"gorim/internal/types"
	"gorim/test"
	"log"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

//more advanced profile tests

type verifyStruct struct {
	name            string
	numMods         int
	expectedPlugins []string
}

func TestChangeProfile(t *testing.T) {
	log.Printf("Running Test [%s]", t.Name())
	assert := assert.New(t)
	profiles := []string{"4", "5", "6"}

	s := state.Initialize()
	s.AddProfile(profiles...)

	assert.Equal(len(profiles), len(s.Profiles))

	//create and add random mods
	numMods := 6
	var mods []types.InternalMod
	for _ = range numMods {
		mods = append(mods, test.CreateRandomMod())
	}

	s.AddMods(mods)
	assert.Equal(numMods, len(s.ModList))

	s.ChangeProfile("4")
	s.EnableMod(mods[2].PackageId, true)

	assert.Equal(1, len(s.ActiveProfile.PluginList))

	for i := range s.Profiles {
		if s.Profiles[i].Name == "4" {
			assert.Equal(1, len(s.Profiles[i].PluginList))
		} else {
			assert.Equal(0, len(s.Profiles[i].PluginList))
		}
	}

	s.ChangeProfile("5")
	assert.Equal(0, len(s.ActiveProfile.PluginList))
	s.EnableMod(mods[0].PackageId, true)
	s.EnableMod(mods[1].PackageId, true)

	assert.Equal(2, len(s.ActiveProfile.PluginList))

	for i := range s.Profiles {
		if s.Profiles[i].Name == "5" {
			assert.Equal(2, len(s.Profiles[i].PluginList))
		}
		if s.Profiles[i].Name == "4" {
			assert.Equal(1, len(s.Profiles[i].PluginList))
		}
	}

}

// tests multiple profiles with a different set of mods
func TestProfileWithMods(t *testing.T) {
	log.Printf("Running Test [%s]", t.Name())
	type verifyStruct struct {
		name            string
		numMods         int
		expectedPlugins []string
	}
	assert := assert.New(t)
	profiles := []string{"4", "5", "6"}

	s := state.Initialize()
	s.AddProfile(profiles...)

	assert.Equal(len(profiles), len(s.Profiles))

	//create and add random mods
	numMods := 6
	var mods []types.InternalMod
	for _ = range numMods {
		mods = append(mods, test.CreateRandomMod())
	}

	s.AddMods(mods)
	assert.Equal(numMods, len(s.ModList))

	//establish the profile + their expected mods list
	var expect map[string]verifyStruct = map[string]verifyStruct{}

	for _, p := range profiles {
		modsToUse := rand.IntN(numMods)

		v := verifyStruct{name: p, numMods: modsToUse, expectedPlugins: []string{}}
		for i := 0; i < modsToUse; i++ {
			v.expectedPlugins = append(v.expectedPlugins, mods[i].PackageId)
		}
		expect[p] = v
	}

	s.AddProfile("none")
	log.Printf("%s: Enabling Mods", t.Name())

	//enable the required mods for each profile
	for _, p := range profiles {
		s.ChangeProfile(p)
		for _, mod := range expect[p].expectedPlugins {
			s.EnableMod(mod, true)
		}
		//assert the changes are applied immediately
		assert.Equal(expect[p].numMods, len(s.ActiveProfile.PluginList))
	}

	s.ChangeProfile("none")
	log.Printf("%s: Testing...\n", t.Name())
	//assert that the underlying plugins arent swept away by changing profiles
	for _, profile := range s.Profiles {
		assert.Equal(len(expect[profile.Name].expectedPlugins), len(profile.PluginList))
	}

	//assert that changing to these profiles doesnt result in any strange behavior
	for _, profile := range s.Profiles {
		s.ChangeProfile(profile.Name)

		assert.Equal(profile.Name, s.ActiveProfile.Name)
		assert.Equal(len(profile.PluginList), len(s.ActiveProfile.PluginList))
		// assert.Equal(len(expect[profile.Name].expectedPlugins), len(s.ActiveProfile.PluginList))

	}

}
