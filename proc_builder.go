//------------------------------------------------------------------------------
// the code here is about how we build a file map from a JSON file content
//------------------------------------------------------------------------------

package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
)

// recursively building an ordered map
func (thisMap *fileMap) build(originalBytes []byte) *fileMap {

	// some initialisation
	thisMap.subMaps = map[string]*fileMap{}
	thisMap.values = map[string]interface{}{}
	thisMap.chainedProperties = map[string]*chainedProperty{}
	thisMap.propertyIndexes = map[string]int{}

	// the index for each property of this map, regarding the given original bytes
	propertiesIndexes := map[string]int{}

	// going through the content of THIS map
	for propertyName, value := range thisMap.originalContent {

		// adding the proerty to the others
		thisMap.orderedProperties = append(thisMap.orderedProperties, propertyName)

		// what the propertyName looks like in JSON
		propertyNameAsBytes, _ := json.Marshal(propertyName)

		// keeping in memory the position of the property within the original JSON content
		propertiesIndexes[propertyName] = bytes.Index(originalBytes, propertyNameAsBytes)

		// dealing with the value
		if reflect.TypeOf(value).Kind() == reflect.Map {

			// adding a new, built, sub-map
			thisMap.subMaps[propertyName] = (&fileMap{
				parent:          thisMap,
				name:            propertyName,
				originalContent: reflect.ValueOf(value).Interface().(map[string]interface{}),
			}).build(originalBytes)

		} else {

			// adding a simple value
			thisMap.values[propertyName] = value
		}
	}

	// sorting the properties, to reflect their original order
	sort.Slice(thisMap.orderedProperties, func(i, j int) bool {
		return propertiesIndexes[thisMap.orderedProperties[i]] < propertiesIndexes[thisMap.orderedProperties[j]]
	})

	// building the chained properties
	for i, propertyName := range thisMap.orderedProperties {

		// new chained property into the map
		property := &chainedProperty{owner: thisMap, name: propertyName}
		thisMap.chainedProperties[propertyName] = property
		thisMap.propertyIndexes[propertyName] = i + 1

		// chaining this property with the preceding one
		if i > 0 {
			property.linkAfter(thisMap.chainedProperties[thisMap.orderedProperties[i-1]], false)
		}
	}

	// returning this map
	return thisMap
}
