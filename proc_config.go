//------------------------------------------------------------------------------
// the code here is about building the config object
//------------------------------------------------------------------------------

package main

import (
	"math"
	"os"
)

// getting the config from an existing JSON (.conf) file, or building it from the definition
func (commonDef *fileMap) getOrInitConfig(folderPath string, folderInfo os.FileInfo) *config {

	config := &config{
		folderPath: folderPath,
		folderInfo: folderInfo,
		content:    commonDef.initConfigMap(),
	}

	return config
}

func (commonDef *fileMap) initConfigMap() *configMap {

	result := &configMap{
		items:   map[string]*configItem{},
		subMaps: map[string]*configMap{},
	}

	for i, property := range commonDef.orderedProperties {

		// about to "compute" foregound and background colors
		var fg, bg string
		var bgtype bgType

		// dealing with the top level
		if commonDef.getDepth() == 1 {
			fg, bg, bgtype = getColorForLevel(1, math.Mod(float64(i), 0) == 0, false)

		} else {

		}

		// adding the config for the current prop
		result.items[property] = &configItem{
			foreground:  fg,
			background:  bg,
			bgroundType: bgtype,
		}

		// going deeper
		if subMap := commonDef.subMaps[property]; subMap != nil {
			result.subMaps[property] = subMap.initConfigMap()
		}
	}

	return result
}

type bgType string

const (
	bgTypeGREY  bgType = "grey"
	bgTypeBLUE  bgType = "blue"
	bgTypeGREEN bgType = "green"
)

// building a color corresponding to a level, and to an odd / even arg
func getColorForLevel(level int, even bool, alt bool) (fg string, bg string, bgtype bgType) {

	// from dark to clear
	greens := [8]string{"#186A3B", "#1D8348", "#239B56", "#28B463", "#2ECC71", "#58D68D", "#82E0AA", "#ABEBC6"}
	blues := [8]string{"#21618C", "#2874A6", "#2E86C1", "#3498DB", "#5DADE2", "#85C1E9", "#D6EAF8", "#EBF5FB"}
	greys := [8]string{"#626567", "#797D7F", "#909497", "#A6ACAF", "#BDC3C7", "#CACFD2", "#D7DBDD", "#E5E7E9"}

	if level <= 4 {
		fg = "#FFFFFF"
	} else {
		fg = "#000000"
	}

	lvl := level - 1
	if level >= 8 {
		lvl = 7
	}

	if alt {
		return fg, greys[lvl], bgTypeGREY
	}

	if even {
		return fg, blues[lvl], bgTypeBLUE
	}
	return fg, greens[lvl], bgTypeGREEN
}
