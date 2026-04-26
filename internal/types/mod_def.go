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
	Order     int
}

type InternalPlugin struct {
	Name    string
	Enabled bool
	Order   int
}
