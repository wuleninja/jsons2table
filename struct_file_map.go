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

// cf. https://stackoverflow.com/a/48301733
type fileMap struct {
	name              string                      // for the root maps, keeping the file name here; for submaps, keeping the property name
	parent            *fileMap                    // for a submap, keeping its parent
	subMaps           map[string]*fileMap         // the maps belonging to a map
	values            map[string]interface{}      // the pure values (non-map) within a map
	orderedProperties []string                    // keeping track of the original order of the properties
	originalContent   map[string]interface{}      // the raw info within this mal
	chainedProperties map[string]*chainedProperty // the properties kept as chained values
	propertyIndexes   map[string]int              // the index of each property in this map
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
