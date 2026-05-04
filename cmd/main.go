package main

import (
	"encoding/json"
	"gorim/internal/components"
	"gorim/internal/pages"
	"gorim/internal/state"
	"gorim/internal/types"
	"gorim/internal/util"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func readFile() (pages.InputParams, error) {
	data, err := os.ReadFile("init.json")
	if err != nil {
		log.Println("Failed to read init")
		return pages.InputParams{}, err
	}

	var inputFormat pages.InputParams
	err = json.Unmarshal(data, &inputFormat)
	if err != nil {
		return pages.InputParams{}, err
	}
	return inputFormat, nil
}

func readRules(path string) (types.CommunityRules, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Println("failed to read rules ", err)
		return types.CommunityRules{}, err
	}

	var rules types.CommunityRules
	err = json.Unmarshal(file, &rules)
	if err != nil {
		log.Println("failed to marshal file ", err)
		return types.CommunityRules{}, err
	}
	return rules, nil
}

func main() {
	state := state.AppState{ModList: map[string]types.InternalMod{}, PluginList: []string{}}

	params, err := readFile()
	if err != nil {
		panic("Unable to read input file")
	}

	rules, err := readRules(params.RulesPath)
	if err != nil {
		log.Println("Failed to find rules")
	}
	state.Rules = rules

	//init
	go util.InitializePaths(params, &state, func(types.LoadedResult) {

	})

	a := app.New()
	w := a.NewWindow("")

	toolbar := components.CreateToolbar()

	modInfo := pages.InputPanel(params, &state)

	pluginInfo := pages.NewPluginList(&state)

	mainPanel := container.NewHSplit(modInfo, pluginInfo)

	contentbox := container.NewBorder(container.NewVBox(toolbar, pages.NewBarPanel(&state)), pages.NewActionsPanel(&state), nil, nil, mainPanel)
	w.SetContent(contentbox)
	w.Resize(fyne.NewSize(1200, 800))

	w.ShowAndRun()

}
