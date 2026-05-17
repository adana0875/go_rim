package state

import (
	"gorim/internal/types"
	"log"
	"slices"
)

func (this *AppState) profilesToString() []string {
	var items []string
	for _, p := range this.Profiles {
		items = append(items, p.Name)
	}
	return items
}

// add a new profile
func (this *AppState) AddProfile(profile ...string) {
	items := this.profilesToString()
	for _, p := range profile {
		if !slices.Contains(items, p) {
			this.Profiles = append(this.Profiles, types.Profile{Name: p, PluginList: []string{}})
		}
	}
}

func (this *AppState) RemoveProfile(profile string) {
	var i int = -1
	for idx, p := range this.Profiles {
		if p.Name == profile {
			i = idx
		}
	}

	if i == -1 {
		log.Println("cannot find index of profile ", profile)
	}

	this.Profiles = slices.Delete(this.Profiles, i, i+1)
}

// change profile functionality
func (this *AppState) ChangeProfile(profile string) {
	if this.ActiveProfile != nil {
		this.SaveMods()
	}
	//try and find profile
	for _, item := range this.Profiles {
		if item.Name == profile {
			log.Printf("using profile %s, switching to %v", this.ActiveProfile, item)
			//before removing all plugins, get a list of the ones we want activated
			this.EnableAll(false)
			log.Println("Enabling plugins for profile: ", item.PluginList)
			this.ActivatePlugins(item.PluginList)
			this.ActiveProfile = &item
			return
		}
	}

	log.Println("unable to find profile ", profile)
}

func (this *AppState) SaveMods() {
	var idx int = -1
	for i, p := range this.Profiles {
		if p.Name == this.ActiveProfile.Name {
			idx = i
		}
	}

	if idx == -1 {
		return
	}

	log.Println("saving profile ", this.ActiveProfile.Name)
	this.Profiles[idx] = *this.ActiveProfile
}
