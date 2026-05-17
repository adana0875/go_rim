package state

import (
	"gorim/internal/state"
	"gorim/internal/types"
	"log"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

type modtest struct {
	profileName string
	mods        []string
}

func profileToString(profile []types.Profile) []string {
	var items []string
	for _, p := range profile {
		items = append(items, p.Name)
	}
	return items
}

func TestProfile(t *testing.T) {
	assert := assert.New(t)
	s := state.Initialize()
	s.AddProfile("test_default")

	assert.True(len(s.Profiles) > 0)
	assert.NotEmpty(s.Profiles[0])
	assert.NotNil(s.Profiles[0].PluginList)
	if assert.NotEmpty(s.Profiles[0].Name) {
		s.ChangeProfile(s.Profiles[0].Name)
		assert.NotNil(s.ActiveProfile)
		assert.Equal(s.ActiveProfile.Name, s.Profiles[0].Name)
	}
}

func TestRemoveProfile(t *testing.T) {
	assert := assert.New(t)
	s := state.Initialize()
	profiles := []string{"0", "1", "2", "3"}
	remove := profiles[2]

	for _, p := range profiles {
		s.AddProfile(p)
	}

	assert.Equal(len(s.Profiles), len(profiles))

	s.RemoveProfile(remove)
	assert.Equal(len(s.Profiles), len(profiles)-1)

	names := profileToString(s.Profiles)
	assert.False(slices.Contains(names, remove))
	assert.True(slices.Contains(names, profiles[0]))
	assert.True(slices.Contains(names, profiles[1]))
	assert.True(slices.Contains(names, profiles[3]))
}

func TestMultiProfile(t *testing.T) {
	p := []types.Profile{{Name: "default"}, {Name: "test1"}}
	assert := assert.New(t)
	s := state.Initialize()
	s.Profiles = append(s.Profiles, p...)

	assert.True(len(s.Profiles) > 0)
	assert.True(len(s.Profiles) == len(p))

	name := p[0].Name
	s.ChangeProfile(name)

	assert.True(s.ActiveProfile.Name == name)

	name = p[len(p)-1].Name
	s.ChangeProfile(name)
	assert.True(s.ActiveProfile.Name == name)
}

// tests multiple profiles with a different set of mods
func TestProfileWithMods(t *testing.T) {
	assert := assert.New(t)
	tests := []modtest{
		{profileName: "test1", mods: []string{"test1mod1", "test1mod2"}},
		{profileName: "test2", mods: []string{"test2mod1"}},
		{profileName: "test3", mods: []string{"test3mod1", "test3mod2", "test3mod3"}},
	}

	s := state.Initialize()
	for _, test := range tests {
		p := types.Profile{Name: test.profileName}
		s.Profiles = append(s.Profiles, p)
	}
	assert.True(len(s.Profiles) > 0)
	assert.True(len(s.Profiles) == len(tests))

	var mods []string
	for _, test := range tests {
		mods = append(mods, test.mods...)
	}

	for _, mod := range mods {
		s.AddMod(types.InternalMod{PackageId: mod})
	}

	assert.True(len(s.ModList) == len(mods))

	for _, profile := range tests {
		s.ChangeProfile(profile.profileName)
		assert.True(len(s.ActiveProfile.PluginList) == 0)

		for _, mod := range profile.mods {
			//enable profile specific mods
			s.EnableMod(mod, true)
		}
		assert.Equal(len(s.ActiveProfile.PluginList), len(profile.mods))

		//make sure the active plugins are expected
		for _, plugin := range s.ActiveProfile.PluginList {
			assert.True(slices.Contains(profile.mods, plugin))
		}
	}
}

// test adding mods only affects the plugins of that profile
func TestAddAndSwitch(t *testing.T) {
	assert := assert.New(t)

	s := state.Initialize()
	mods := []types.InternalMod{
		{Name: "mod1", PackageId: "mod.1"},
		{Name: "mod2", PackageId: "mod.2"},
		{Name: "mod3", PackageId: "mod.3"},
	}
	s.AddProfile([]string{"default", "test", "test2"}...)

	s.AddMods(mods)
	s.ChangeProfile("default")

	assert.Equal(s.ActiveProfile.Name, "default")
	assert.Equal(len(s.ModList), len(mods))
	assert.Equal(len(s.ActiveProfile.PluginList), 0)

	s.ChangeProfile("test")
	assert.Equal(s.ActiveProfile.Name, "test")
	assert.Equal(len(s.ModList), len(mods))
	assert.Equal(len(s.ActiveProfile.PluginList), 0)
	s.EnableAll(true)

	log.Println("Profiles: ", s.Profiles)
	log.Println("Mods: ", s.ModList)
	log.Println("active: ", s.ActiveProfile)
	assert.Equal(len(s.ActiveProfile.PluginList), len(mods))

	s.ChangeProfile("default")
	log.Println("Profiles: ", s.Profiles)
	assert.Equal(s.ActiveProfile.Name, "default")
	assert.Equal(len(s.ModList), len(mods))
	assert.Equal(len(s.ActiveProfile.PluginList), 0)

	s.ChangeProfile("test2")
	assert.Equal(s.ActiveProfile.Name, "test2")
	assert.Equal(len(s.ModList), len(mods))
	assert.Equal(len(s.ActiveProfile.PluginList), 0)
	s.EnableMod(mods[0].PackageId, true)

	assert.Equal(len(s.ActiveProfile.PluginList), 1)

	s.ChangeProfile("test")
	assert.Equal(s.ActiveProfile.Name, "test")
	assert.Equal(len(s.ModList), len(mods))
	assert.Equal(len(mods), len(s.ActiveProfile.PluginList))

	s.ChangeProfile("test2")
	assert.Equal(s.ActiveProfile.Name, "test2")
	assert.Equal(len(s.ModList), len(mods))
	assert.Equal(1, len(s.ActiveProfile.PluginList))
}
