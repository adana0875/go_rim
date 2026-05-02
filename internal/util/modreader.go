package util

import (
	"encoding/xml"
	"errors"
	"fmt"
	"gorim/internal/pages"
	"gorim/internal/state"
	"gorim/internal/types"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// synchronous load - also enables mods subscribed to and loaded
func InitializePaths(inputs pages.InputParams, appState *state.AppState, callback func(result types.LoadedResult)) {
	collectMods(inputs.WorkshopPath, appState)
	collectPluginList(inputs.Modpath, appState)

	//check plugins for if we have it in modlist.
	//if we have an item in plugin list, we should be subscribed to it
	for _, plugin := range appState.PluginList {
		if strings.Contains(plugin, "ludeon.rimworld") {
			continue
		}
		_, ok := appState.ModList[plugin]
		if !ok {
			log.Println("Error: modlist doest have plugin : ", plugin)
			callback(types.LoadFailure)
		}

		//enable the mod
		appState.EnableMod(plugin, true)

	}
	callback(types.LoadSuccess)
}

func collectMods(pathToWorkshop string, appState *state.AppState) {
	mods, err := lookUpDirectoryMods(pathToWorkshop)
	if err != nil {
		log.Println("Failed to fetch mods")
		return
	}

	//add mods to state
	appState.AddMods(mods)
}

func collectPluginList(pathToMods string, appState *state.AppState) {
	err, active := lookupCurrentMods(pathToMods)

	if err != nil {
		log.Println("Failed to get active mods")
		return
	}

	appState.AddPlugins(active)
}

// Function to look up installed workshop mods - the order wont matter here
func lookUpDirectoryMods(path string) ([]types.InternalMod, error) {
	var mods []types.InternalMod

	files, err := os.ReadDir(path)
	if err != nil {
		return mods, err
	}

	var validDirs []string
	//go through the files, only get the ones that are presumably mods
	for _, item := range files {
		//expect mods to be a directory
		if item.IsDir() {
			validDirs = append(validDirs, item.Name())
		}
	}

	var conformingMods []string
	//get the folder for the mod
	for _, item := range validDirs {
		files, err := os.ReadDir(fmt.Sprintf("%s/%s", path, item))
		if err != nil {
			log.Println("failed to read mod ", item)
			continue
		}
		var foundFilesDir bool = false
		for _, file := range files {
			if file.Name() == "About" {
				conformingMods = append(conformingMods, fmt.Sprintf("%s/%s/About", path, item))
				foundFilesDir = true
			}
		}
		if !foundFilesDir {
			log.Printf("Couldnt find files for mod %s: %s", item, files)
		}
	}

	for _, item := range conformingMods {
		if strings.Contains(item, "3237397753") {
			log.Println("found mod")
		}

		files, err := os.ReadDir(item)
		if err != nil {
			log.Println("Failed to read folder for mod ", item)
			return mods, err
		}

		for _, file := range files {
			name := strings.ToLower(file.Name())
			if name == "about.xml" || name == "about" {
				data, err := os.ReadFile(filepath.Join(item, file.Name()))
				if err != nil {
					log.Println("Failed to read file ", file.Name())
					continue
				}

				var mod types.ModDef
				err = xml.Unmarshal(data, &mod)
				if err != nil {
					log.Printf("Failed to read mod data for file %s: %s", item, err)
					continue
				}

				mods = append(mods, types.InternalMod{Name: mod.Name, Enabled: false, PackageId: strings.ToLower(mod.PackageId), LoadAfter: mod.LoadOrder})
			}
		}

	}
	return mods, nil
}

func lookupCurrentMods(path string) (error, []types.InternalPlugin) {
	var pluginList []types.InternalPlugin
	//read the directory
	files, err := os.ReadDir(path)
	if err != nil {
		log.Println("cant read directory for modlist")
		return err, pluginList
	}

	var rawPlugins []string
	var parsed bool = false
	for _, file := range files {
		if file.Name() == "ModsConfig" || file.Name() == "ModsConfig.xml" {
			data, err := os.ReadFile(filepath.Join(path, "ModsConfig.xml"))
			if err != nil {
				log.Println("Failed to read mods xml: ", err)
				return err, pluginList
			}

			var config types.ModsConfig
			err = xml.Unmarshal(data, &config)

			if err != nil {
				log.Println("Failed to unmarshal plugin file: ", err)
				return err, pluginList
			}

			rawPlugins = config.ActiveMods
			parsed = true
			break
		}
	}

	if parsed == false {
		return errors.New("failed to get mods"), pluginList
	}

	for idx, mod := range rawPlugins {
		pluginList = append(pluginList, types.InternalPlugin{Name: strings.ToLower(mod), Enabled: true, Order: idx})
	}

	return nil, pluginList

}
