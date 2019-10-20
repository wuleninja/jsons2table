//------------------------------------------------------------------------------
// the structure with which
//------------------------------------------------------------------------------

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// cf. https://stackoverflow.com/a/48301733
type orderedMap struct {
	name            string                 // for the root maps, keeping the file name here
	parent          *orderedMap            // for a submap, keeping its parent
	subMaps         map[string]*orderedMap // the maps belonging to a map
	values          map[string]interface{} // the pure values (non-map) within a map
	orderedKeys     []string               // keeping track of the original order of the keys
	originalContent map[string]interface{} // the raw info within this mal
	chainedKeys     map[string]*chainedKey // the keys kept as chained values
}

type chainedKey struct {
	key    string
	before *chainedKey
	after  *chainedKey
}

// UnmarshalJSON : keeping the keys order
func (thisMap *orderedMap) UnmarshalJSON(bytes []byte) error {

	// classic unmarshalling to get the content
	if err := json.Unmarshal(bytes, &thisMap.originalContent); err != nil {
		return err
	}

	// recursively building the ordered map
	thisMap.build(bytes)

	return nil
}

// recursively building an ordered map
func (thisMap *orderedMap) build(originalBytes []byte) *orderedMap {

	// some initialisation
	thisMap.subMaps = map[string]*orderedMap{}
	thisMap.values = map[string]interface{}{}
	thisMap.chainedKeys = map[string]*chainedKey{}

	// the index for each key of this map, regarding the given original bytes
	keyIndexes := map[string]int{}

	// going through the content of THIS map
	for key, value := range thisMap.originalContent {

		// adding the key to the others
		thisMap.orderedKeys = append(thisMap.orderedKeys, key)

		// what the key looks like in JSON
		keyAsBytes, _ := json.Marshal(key)

		// keeping in memory the position of the key within the original JSON content
		keyIndexes[key] = bytes.Index(originalBytes, keyAsBytes)

		// dealing with the value
		if reflect.TypeOf(value).Kind() == reflect.Map {

			// adding a new, built, sub-map
			thisMap.subMaps[key] = (&orderedMap{
				parent:          thisMap,
				originalContent: reflect.ValueOf(value).Interface().(map[string]interface{}),
			}).build(originalBytes)

		} else {

			// adding a simple value
			thisMap.values[key] = value
		}
	}

	// sorting the keys, to reflect the original order
	sort.Slice(thisMap.orderedKeys, func(i, j int) bool {
		return keyIndexes[thisMap.orderedKeys[i]] < keyIndexes[thisMap.orderedKeys[j]]
	})

	// building the chained keys
	for _, key := range thisMap.orderedKeys {
		thisMap.chainedKeys[key] = &chainedKey{key: key}
	}
	for i, key := range thisMap.orderedKeys {
		currentChainedKey := thisMap.chainedKeys[key]
		if i > 0 {
			currentChainedKey.before = thisMap.chainedKeys[thisMap.orderedKeys[i-1]]
		}
		if i < len(thisMap.orderedKeys)-1 {
			currentChainedKey.after = thisMap.chainedKeys[thisMap.orderedKeys[i+1]]
		}
	}

	// returning this map
	return thisMap
}

// MarshalJSON : override
func (thisMap orderedMap) MarshalJSON() ([]byte, error) {

	// we don't really care about marshalling here (for now)
	return json.Marshal(thisMap)
}

// showing an ordered map content
func (thisMap orderedMap) display(indent int) {
	for _, key := range thisMap.orderedKeys {
		chainedKey := thisMap.chainedKeys[key]
		beforeKey := "-"
		if chainedKey.before != nil {
			beforeKey = chainedKey.before.key
		}
		afterKey := "-"
		if chainedKey.after != nil {
			afterKey = chainedKey.after.key
		}
		if submap, found := thisMap.subMaps[key]; found {
			println(fmt.Sprintf(strings.Repeat("  ", indent)+"'%s': [ before = %s / after = %s ]", key, beforeKey, afterKey))
			submap.display(indent + 1)
		} else {
			println(fmt.Sprintf(strings.Repeat("  ", indent)+"'%s': '%v' [ before = %s / after = %s ]", key, thisMap.values[key], beforeKey, afterKey))
		}
	}
}
