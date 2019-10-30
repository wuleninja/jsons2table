//------------------------------------------------------------------------------
// the code here is about building the config object
//------------------------------------------------------------------------------

package main

import (
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

	for index, property := range commonDef.orderedProperties {

		// about to "compute" foregound and background colors
		var bg string

		// dealing with the top level, i.e. the first header line
		if commonDef.parent == nil {
			bg = getAdjustedColor(colors[index%len(colors)], 0)

		} else {

			// dealing with a sub-section level
			level := 10
			if index%2 == 1 {
				level = 30
			}
			bg = getAdjustedColor(commonDef.parent.chainedProperties[commonDef.name].conf.background, level*commonDef.getDepth())
		}

		// new item
		newItem := &configItem{
			background: bg,
		}

		// adding the config for the current prop
		result.items[property] = newItem

		// linking the prop to its config
		commonDef.chainedProperties[property].conf = newItem

		// going deeper
		if subMap := commonDef.subMaps[property]; subMap != nil {
			result.subMaps[property] = subMap.initConfigMap()
		}
	}

	return result
}
