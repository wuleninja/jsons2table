//------------------------------------------------------------------------------
// the structure with which
//------------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// the "path" is the complete name for a property or a submap, that includes the parents' names
// e.g. "SubMap_2.SubMap_7.Property_11"
type path string

// cf. https://stackoverflow.com/a/48301733
type fileMap struct {
	name                 string                      // for the root maps, keeping the file name here; for submaps, keeping the property name
	parent               *fileMap                    // for a submap, keeping its parent
	subMaps              map[string]*fileMap         // the maps belonging to a map
	values               map[string]interface{}      // the pure values (non-map) within a map
	orderedProperties    []string                    // keeping track of the original order of the properties
	originalContent      map[string]interface{}      // the raw info within this map
	chainedProperties    map[string]*chainedProperty // the properties kept as chained values
	propertyIndexes      map[string]int              // the index of each property in this map
	height               int                         // this data tree's height
	depth                int                         // this data tree's depth
	path                 path                        // this map's full path within the common definition
	allChainedProperties map[path]*chainedProperty   // indexing all the chained properties from the root common definition
}

// UnmarshalJSON : keeping the properties' order
func (thisMap *fileMap) UnmarshalJSON(bytes []byte) error {

	// classic unmarshalling to get the content
	if err := json.Unmarshal(bytes, &thisMap.originalContent); err != nil {
		return err
	}

	// recursively building the ordered map
	thisMap.build(bytes)

	return nil
}

// MarshalJSON : override
func (thisMap *fileMap) MarshalJSON() ([]byte, error) {

	// we don't really care about marshalling here (for now)
	return json.Marshal(thisMap)
}

// showing an ordered map content, with respect to the order of the properties
func (thisMap *fileMap) displayOrdered(indent int, valueFn func(*fileMap, string) interface{}) {
	if indent == 0 {
		println("\n======================================================================================================= " + thisMap.name + "\n")
	}
	for _, propertyName := range thisMap.orderedProperties {
		chainedProp := thisMap.chainedProperties[propertyName]
		indexString := ""
		if index := thisMap.propertyIndexes[propertyName]; index != 0 {
			indexString = fmt.Sprintf("%d - ", index)
		}
		if submap, found := thisMap.subMaps[propertyName]; found {
			println(fmt.Sprintf(strings.Repeat("  ", indent)+"%s'%s': [ previous = %s / next = %s ]",
				indexString, propertyName, chainedProp.previous, chainedProp.next))
			submap.displayOrdered(indent+1, valueFn)
		} else {
			println(fmt.Sprintf(strings.Repeat("  ", indent)+"%s'%s': '%v' [ previous = %s / next = %s ]",
				indexString, propertyName, valueFn(thisMap, propertyName), chainedProp.previous, chainedProp.next))
		}
	}
}

// returns one chained property at random
func (thisMap *fileMap) oneChainedProperty() *chainedProperty {
	for _, prop := range thisMap.chainedProperties {
		if !prop.addOn {
			return prop
		}
	}
	return nil
}

// getting the kind for a property of the given built file map
func (thisMap *fileMap) getPropertyKind(propertyName string) reflect.Kind {
	if thisMap.subMaps[propertyName] != nil {
		return reflect.Map
	}
	return reflect.ValueOf(thisMap.values[propertyName]).Kind()
}

// showing the value for a property of the given built file map
func showValue(thisMap *fileMap, propertyName string) interface{} {
	return thisMap.values[propertyName]
}

// showing the kind for a property of the given built file map
func showKind(thisMap *fileMap, propertyName string) interface{} {
	return thisMap.chainedProperties[propertyName].kind
}

// getting this map's height
func (thisMap *fileMap) getHeight() int {

	// we already know, let's return
	if thisMap.height > 0 {
		return thisMap.height
	}

	// nothing under this, so the height is 1
	if len(thisMap.subMaps) == 0 {
		thisMap.height = 1

	} else { // this is 1 + the max height under this
		for _, subMap := range thisMap.subMaps {
			if newHeight := subMap.getHeight(); newHeight > thisMap.height {
				thisMap.height = newHeight
			}
		}
		thisMap.height = 1 + thisMap.height
	}

	return thisMap.height
}

// getting this map's depth
func (thisMap *fileMap) getDepth() int {

	// we already know, let's return
	if thisMap.depth > 0 {
		return thisMap.depth
	}

	// nothing over this, so the depth is 1
	if thisMap.parent == nil {
		thisMap.depth = 1

	} else { // this is 1 + the max height under this
		thisMap.depth = 1 + thisMap.parent.getDepth()
	}

	return thisMap.depth
}

// getting the full name
func (thisMap *fileMap) getFullName() string {
	result := thisMap.name
	for parent := thisMap.parent; parent != nil; parent = parent.parent {
		result = parent.name + " / " + result
	}
	return result
}

// getting the first index for this file map
func (thisMap *fileMap) getFirstIndex() int {
	firstProperty := thisMap.orderedProperties[0]
	if subMap := thisMap.subMaps[firstProperty]; subMap != nil {
		return subMap.getFirstIndex()
	}
	return thisMap.chainedProperties[firstProperty].index
}

// getting the last index for this file map
func (thisMap *fileMap) getLastIndex() int {
	lastProperty := thisMap.orderedProperties[len(thisMap.orderedProperties)-1]
	if subMap := thisMap.subMaps[lastProperty]; subMap != nil {
		return subMap.getLastIndex()
	}
	return thisMap.chainedProperties[lastProperty].index
}

// getting the common definition from any of its submaps
func (thisMap *fileMap) root() *fileMap {
	if thisMap.parent == nil {
		return thisMap
	}
	return thisMap.parent.root()
}

// getting the path for this (sub-)map
func (thisMap *fileMap) getPath() path {
	if thisMap.path == "" {
		if thisMap.parent == nil {
			thisMap.path = path("")
		} else {
			thisMap.path = path(fmt.Sprintf("%s%s/", thisMap.parent.getPath(), thisMap.name))
		}
	}
	return thisMap.path
}

// indexing the given property at the root level
func (thisMap *fileMap) register(prop *chainedProperty) {
	thisMap.root().allChainedProperties[prop.getPath()] = prop
}

// returns a property indexed on the COMMON definition thanks to its path
func (thisMap *fileMap) getProp(propPath path) *chainedProperty {
	return thisMap.root().allChainedProperties[propPath]
}

// returns a property in a given JSON MAP (not the common definition) thanks to its path
func (thisMap *fileMap) findProp(propPath path) *chainedProperty {

	pathAsString := string(propPath)

	// we're not at the right level yet
	if sepIndex := strings.Index(pathAsString, "/"); sepIndex > 0 {
		subMapName := pathAsString[0:sepIndex]
		if subMap := thisMap.subMaps[subMapName]; subMap != nil {
			propSubPath := path(pathAsString[(sepIndex + 1):len(pathAsString)])
			return subMap.findProp(propSubPath)
		}
		return nil
	}
	return thisMap.chainedProperties[pathAsString]
}
