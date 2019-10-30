//------------------------------------------------------------------------------
// the code here is about the JSON config file (but with a .conf extension)
// that contains useful, customizable info about the Excel file generation
//------------------------------------------------------------------------------

package main

import "os"

type config struct {
	folderPath string
	folderInfo os.FileInfo
	content    *configMap
}

type configMap struct {
	items   map[string]*configItem
	subMaps map[string]*configMap
}

type configItem struct {
	foreground string
	background string
}
