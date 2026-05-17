package state

import (
	"fmt"
	"gorim/internal/state"
	"gorim/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddMod(t *testing.T) {
	assert := assert.New(t)
	s := state.Initialize()

	assert.Equal(0, len(s.ModList))
	s.AddMod(types.InternalMod{Name: "test", PackageId: "test.mod"})
	assert.Equal(1, len(s.ModList))
}

func TestAddMods(t *testing.T) {
	assert := assert.New(t)
	s := state.Initialize()

	assert.Equal(0, len(s.ModList))
	nMods := 10
	var modsToAdd []types.InternalMod = make([]types.InternalMod, nMods)
	for i := range nMods - 1 {
		modsToAdd = append(modsToAdd, types.InternalMod{Name: "testmod", PackageId: fmt.Sprintf("test_%d", i)})
	}

	s.AddMods(modsToAdd)
	assert.Equal(nMods, len(s.ModList))
}

func TestEnableMods(t *testing.T) {
	testMod := types.InternalMod{Name: "Test Mod", PackageId: "test.mod"}
	assert := assert.New(t)
	s := state.Initialize()
	s.AddProfile("test_default")
	s.ChangeProfile(s.Profiles[0].Name)

	//add a mod
	assert.True(len(s.ModList) <= 0, "expected empty modlist")
	s.AddMod(testMod)
	assert.True(len(s.ModList) > 0, "expected non empty modlist")

	//plugin
	assert.True(len(s.ActiveProfile.PluginList) == 0)
	s.EnableMod(testMod.PackageId, true)
	assert.True(len(s.ActiveProfile.PluginList) > 0)

}

func TestEnableDisableMods(t *testing.T) {
	testMod := []types.InternalMod{
		{Name: "Test Mod", PackageId: "test.mod"},
		{Name: "Test Mod2", PackageId: "test.mod.new"},
		{Name: "Test Mod3", PackageId: "test.new.mod"},
	}
	assert := assert.New(t)
	s := state.Initialize()
	s.Profiles = append(s.Profiles, types.Profile{Name: "test_default"})
	s.ChangeProfile(s.Profiles[0].Name)

	//add a mod
	assert.True(len(s.ModList) <= 0, "expected empty modlist")
	s.AddMods(testMod)
	assert.True(len(s.ModList) > 0, "expected non empty modlist")
	assert.True(len(s.ModList) == len(testMod), "expected non empty modlist")

	//plugin
	assert.True(len(s.ActiveProfile.PluginList) == 0)
	s.EnableMod(testMod[0].PackageId, true)
	assert.True(len(s.ActiveProfile.PluginList) > 0)

	//disable
	s.EnableMod(testMod[0].PackageId, false)
	assert.True(len(s.ActiveProfile.PluginList) == 0)

	//enable all
	s.EnableAll(true)
	assert.True(len(s.ActiveProfile.PluginList) == len(testMod))
}
