package types

import "encoding/xml"

var RoyaltyDLC string = "ludeon.rimworld.royalty"
var IdeologyDLC string = "ludeon.rimworld.ideology"
var BiotechDLC string = "ludeon.rimworld.biotech"
var AnomalyDLC string = "ludeon.rimworld.anomaly"
var OdysseyDLC string = "ludeon.rimworld.odyssey"
var BaseGame string = "ludeon.rimworld"

type ModsConfig struct {
	XMLName         xml.Name `xml:"ModsConfigData"`
	Version         string   `xml:"version"`
	ActiveMods      []string `xml:"activeMods>li"`
	KnownExpansions []string `xml:"knownExpansions>li"`
}
