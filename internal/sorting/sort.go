package sorting

import (
	"gorim/internal/types"
	"log"
	"maps"
	"slices"
	"strings"

	"github.com/oko/toposort"
)

var beforeSeparator string = "== before.base =="
var baseSeparator string = "== base =="
var afterSeparator string = "== after.base =="

type TopoNode struct {
	ID string
}

func (t *TopoNode) Id() string {
	return t.ID
}

var _ toposort.Node = &TopoNode{}
var baseList []string = []string{types.BaseGame, types.RoyaltyDLC, types.IdeologyDLC, types.BiotechDLC, types.AnomalyDLC, types.OdysseyDLC}

func TopoSortList(lookup map[string]types.InternalMod, plugins []string, providedRules types.CommunityRules) ([]string, error) {

	var nodesMap map[string]toposort.Node = map[string]toposort.Node{}

	//for reading community rules

	t := toposort.NewTopology()

	//group nodes
	beforeBase := &TopoNode{ID: beforeSeparator}
	t.AddNode(beforeBase)
	baseGroup := &TopoNode{ID: baseSeparator}
	t.AddNode(baseGroup)
	afterBase := &TopoNode{ID: afterSeparator}
	t.AddNode(afterBase)

	//organize groups
	t.AddEdge(beforeBase, baseGroup)
	t.AddEdge(baseGroup, afterBase)

	//sorting functionality
	//go through all plugins

	// sorting
	// default plugin (base game, dlc) are put into a group

	//there are additional markers for plugins which go before base game, and after base game
	// these all go before the marker
	for _, plugin := range plugins {
		modNode := &TopoNode{ID: plugin}
		nodesMap[plugin] = modNode
		t.AddNode(modNode)
	}

	for _, plugin := range plugins {
		//handle default plugins
		if slices.Contains(baseList, plugin) {
			t.AddEdge(nodesMap[plugin], baseGroup)
		}
		//handle dependencies
		if mod, ok := lookup[plugin]; ok {
			var loadAfter []string = []string{}
			var loadBefore []string = []string{}

			//has a community rule
			if item, ok := providedRules.Rules[plugin]; ok {
				loadAfter = append(loadAfter, slices.Collect(maps.Keys(item.LoadAfter))...)
				loadBefore = append(loadBefore, slices.Collect(maps.Keys(item.LoadBefore))...)

			}
			//has steam-provided dependency
			if len(mod.LoadAfter) > 0 {
				loadAfter = append(loadAfter, mod.LoadAfter...)
			}
			if len(mod.LoadBefore) > 0 {
				loadBefore = append(loadBefore, mod.LoadBefore...)
			}

			//assign the group
			loadsBefore := false
			for _, dep := range loadBefore {
				d := strings.ToLower(dep)
				log.Println("item ", d)

				if slices.Contains(baseList, d) {
					if !loadsBefore {
						loadsBefore = true
					}
				} else {
					//check if we are subscribed
					if _, ok := lookup[d]; !ok {
						log.Println("not subscribed to  ", d)
						continue
					}

					//check in case we dont have the dependency as a node
					if _, ok := nodesMap[d]; !ok {
						node := &TopoNode{ID: d}
						nodesMap[d] = node
					}
					t.AddEdge(nodesMap[plugin], nodesMap[d])
				}
			}

			loadsAfter := false
			for _, dep := range loadAfter {
				d := strings.ToLower(dep)
				log.Println("after item ", d)

				if slices.Contains(baseList, d) {
					if !loadsAfter {
						loadsAfter = true
					}
				} else {

					if _, ok := lookup[d]; !ok {
						log.Println("not sub after ", d)
						continue
					}

					if _, ok := nodesMap[d]; !ok {
						node := &TopoNode{ID: d}
						nodesMap[d] = node
					}
					t.AddEdge(nodesMap[d], nodesMap[plugin])
				}
			}
			log.Printf("plugin %s before[%v], after[%v]", plugin, loadsBefore, loadsAfter)

			//add additional ruling - if we have discovered the plugin loads before base game put that, else it will load after base game
			if loadsBefore {
				t.AddEdge(nodesMap[plugin], beforeBase)
			} else {
				t.AddEdge(baseGroup, nodesMap[plugin])
				t.AddEdge(nodesMap[plugin], afterBase)
			}
		}
	}

	sorted, err := t.Sort()
	if err != nil {
		log.Println("Failed to sort topology")
		return plugins, err
	}

	var sortedPlugins []string = make([]string, len(sorted))
	for i, s := range sorted {
		sortedPlugins[i] = s.Id()
	}

	//pull out the temporary separators
	beforeIdx := slices.Index(sortedPlugins, beforeSeparator)
	sortedPlugins = slices.Delete(sortedPlugins, beforeIdx, beforeIdx+1)

	baseIdx := slices.Index(sortedPlugins, baseSeparator)
	sortedPlugins = slices.Delete(sortedPlugins, baseIdx, baseIdx+1)

	afterIdx := slices.Index(sortedPlugins, afterSeparator)
	sortedPlugins = slices.Delete(sortedPlugins, afterIdx, afterIdx+1)
	return sortedPlugins, nil

}

func includesBaseGame(deps []string) bool {
	return slices.Contains(deps, types.BaseGame)
}
