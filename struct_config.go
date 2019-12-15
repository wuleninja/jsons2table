//------------------------------------------------------------------------------
// the code here is about the JSON config file (but with a .conf extension)
// that contains useful, customizable info about the Excel file generation
//------------------------------------------------------------------------------

package main

import "os"

// the struct for the config file
type j2tConfig struct {
	folderPath string
	folderInfo os.FileInfo
	content    *configMap
	General    *generalConfig    `json:"General"`
	NewColumns *newColumnsConfig `json:"NewColumns"`
}

type configMap struct {
	items   map[string]*configItem
	subMaps map[string]*configMap
}

type configItem struct {
	background string
}

type generalConfig struct {
	TrueValue  string `json:"TrueValue"`
	FalseValue string `json:"FalseValue"`
}

type newColumnsConfig struct {
	NewDurations []*newDuration `json:"NewDurations"`
	NewSums      []*newSum      `json:"NewSums"`
}

type newColumn struct {
	Name     string `json:"Name"`
	PutAfter path   `json:"PutAfter"`
}

type newDuration struct {
	newColumn
	From path `json:"From"`
	To   path `json:"To"`
}

type newSum struct {
	newColumn
	AddTogether       []path `json:"AddTogether"`
	SubstractTogether []path `json:"SubstractTogether"`
}
