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
		content:    commonDef.initConfigMap(true),
	}

	return config
}

func (commonDef *fileMap) initConfigMap(isBlue bool) *configMap {

	result := &configMap{
		items:   map[string]*configItem{},
		subMaps: map[string]*configMap{},
	}

	for i, property := range commonDef.orderedProperties {

		// are we even ?
		thisEven := math.Mod(float64(i), 2) == 0

		// about to "compute" foregound and background colors
		var fg, bg string

		// dealing with the top level
		if commonDef.parent == nil {

			if commonDef.getHeight() > 1 {
				fg, bg = getColorForLevel(1, thisEven, false)
			} else {
				fg, bg = getColorForLevel(1, math.Mod(float64(i), 4) == 0, math.Mod(float64(i), 2) == 1)
			}

		} else {

			// dealing with the general level
			fg, bg = getColorForLevel(commonDef.getDepth(), isBlue, !thisEven)
		}

		// new item
		newItem := &configItem{
			foreground: fg,
			background: bg,
		}

		// adding the config for the current prop
		result.items[property] = newItem

		// linkng the prop to its config
		commonDef.chainedProperties[property].conf = newItem

		// going deeper
		if subMap := commonDef.subMaps[property]; subMap != nil {
			if commonDef.parent == nil {
				result.subMaps[property] = subMap.initConfigMap(thisEven)
			} else {
				result.subMaps[property] = subMap.initConfigMap(isBlue)
			}
		}
	}

	return result
}

// building a color corresponding to a level, and to an odd / even arg
func getColorForLevel(level int, isBlue bool, isGrey bool) (fg string, bg string) {

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

	if isGrey {
		return fg, greys[lvl]
	}

	if isBlue {
		return fg, blues[lvl]
	}
	return fg, greens[lvl]
}
