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
	//try and find profile
	for i := range this.Profiles {
		if this.Profiles[i].Name == profile {
			//swap in the profile - this should swap from storage everything about profile
			this.ActiveProfile = &this.Profiles[i]
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
