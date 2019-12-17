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
	General    *generalConfig     `json:"General"`
	NewColumns []*newColumnConfig `json:"NewColumns"`
}

type configItem struct {
	background string
}

type generalConfig struct {
	TrueValue  string `json:"TrueValue"`
	FalseValue string `json:"FalseValue"`
}

type newColumnConfig struct {
	Name               string             `json:"Name"`
	PutAfter           path               `json:"PutAfter"`
	Formula            string             `json:"Formula"`
	NoStat             bool               `json:"NoStat"`
	formattableFormula string             // the formula ready to be filled with real column coordinates
	columns            []*chainedProperty // the columns involved in the definition of the formula
}
