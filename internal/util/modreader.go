package util

import (
	"encoding/xml"
	"errors"
	"fmt"
	"gorim/internal/state"
	"gorim/internal/types"
	"log"
	"os"
	"path/filepath"
)

func CollectMods(pathToWorkshop string, appState *state.AppState) {
	mods, err := lookUpDirectoryMods(pathToWorkshop)
	if err != nil {
		log.Println("Failed to fetch mods")
		return
	}

	//add mods to state
	appState.AddMods(mods)
}

func CollectPluginList(pathToMods string, appState *state.AppState) {
	err, active := lookupCurrentMods(pathToMods)

	if err != nil {
		log.Println("Failed to get active mods")
		return
	}

	appState.AddPlugins(active)
}

// Function to look up installed workshop mods
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
		for _, file := range files {
			if file.Name() == "About" {
				conformingMods = append(conformingMods, fmt.Sprintf("%s/%s/About", path, item))
			}
		}
	}

	for _, item := range conformingMods {
		files, err := os.ReadDir(item)
		if err != nil {
			log.Println("Failed to read folder for mod ", item)
			return mods, err
		}

		for _, file := range files {
			if file.Name() == "About" || file.Name() == "About.xml" {

				data, err := os.ReadFile(filepath.Join(item, files[0].Name()))
				if err != nil {
					log.Println("Failed to read file ", files[0].Name())
					continue
				}

				var mod types.ModDef
				err = xml.Unmarshal(data, &mod)
				if err != nil {
					log.Println("Failed to read mod data ", err)
					continue
				}

				mods = append(mods, types.InternalMod{Name: mod.Name, Enabled: false, PackageId: mod.PackageId})
			}
		}

	}
	return mods, nil
}

func lookupCurrentMods(path string) (error, []types.InternalPlugin) {
	var pluginList []types.InternalPlugin
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

	for _, mod := range rawPlugins {
		pluginList = append(pluginList, types.InternalPlugin{Name: mod, Enabled: true})
	}

	return nil, pluginList

}
