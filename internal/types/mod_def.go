package types

import "encoding/xml"

type ModDef struct {
	XMlName      xml.Name       `xml:"ModMetaData"`
	Name         string         `xml:"name"`
	Author       string         `xml:"author"`
	Url          string         `xml:"url"`
	Versions     []string       `xml:"supportedVersions>li"`
	Dependencies []Dependencies `xml:"modDependencies>li"`
	LoadOrder    []string       `xml:"loadAfter>li"`
	PackageId    string         `xml:"packageId"`
	Description  string         `xml:"description"`
}

type Dependencies struct {
	PackageId   string `xml:"packageId"`
	DisplayName string `xml:"displayName"`
	WorkshopUrl string `xml:"steamWorkshopUrl"`
	DownloadUrl string `xml:"downloadUrl"`
}

type InternalMod struct {
	Name      string
	PackageId string
	Enabled   bool
	LoadAfter []string
	Order     int
}

type ModByOrder []InternalMod

func (plugin ModByOrder) Len() int           { return len(plugin) }
func (plugin ModByOrder) Swap(i, j int)      { plugin[i], plugin[j] = plugin[j], plugin[i] }
func (plugin ModByOrder) Less(i, j int) bool { return plugin[i].Order < plugin[j].Order }

type InternalPlugin struct {
	Name    string
	Enabled bool
	Order   int
}

type PluginByOrder []InternalPlugin

func (plugin PluginByOrder) Len() int           { return len(plugin) }
func (plugin PluginByOrder) Swap(i, j int)      { plugin[i], plugin[j] = plugin[j], plugin[i] }
func (plugin PluginByOrder) Less(i, j int) bool { return plugin[i].Order < plugin[j].Order }
