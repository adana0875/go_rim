package state

import (
	"gorim/internal/state"
	"gorim/internal/types"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

type modtest struct {
	profileName string
	mods        []string
}

func TestProfile(t *testing.T) {
	assert := assert.New(t)
	s := state.Initialize()
	s.Profiles = append(s.Profiles, types.Profile{Name: "test_default"})

	assert.True(len(s.Profiles) > 0)
	assert.NotEmpty(s.Profiles[0])
	if assert.NotEmpty(s.Profiles[0].Name) {
		s.ChangeProfile(s.Profiles[0].Name)
		assert.NotNil(s.ActiveProfile)
		assert.Equal(s.ActiveProfile.Name, s.Profiles[0].Name)
	}
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
