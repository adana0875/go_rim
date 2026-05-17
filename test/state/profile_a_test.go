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

// test basic profile within the state
func TestProfile(t *testing.T) {
	log.Printf("Running Test [%s]", t.Name())
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

// test adding multiple profiles works
func TestMultiProfile(t *testing.T) {
	log.Printf("Running Test [%s]", t.Name())
	p := []string{"default", "test1"}
	assert := assert.New(t)
	s := state.Initialize()
	s.AddProfile(p...)

	assert.True(len(s.Profiles) > 0)
	assert.True(len(s.Profiles) == len(p))

	name := p[0]
	s.ChangeProfile(name)

	assert.True(s.ActiveProfile.Name == name)

	name = p[len(p)-1]
	s.ChangeProfile(name)
	assert.True(s.ActiveProfile.Name == name)
}

// test removing a profile works
func TestRemoveProfile(t *testing.T) {
	log.Printf("Running Test [%s]", t.Name())
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
