//------------------------------------------------------------------------------
// the code here is responsible for finding a common definition for
// all the maps resulting from unmarshalling all the JSON files
//------------------------------------------------------------------------------

package main

import "fmt"

// digesting all the file maps, to build a common definition for them
func merge(fileMaps []*fileMap) *fileMap {

	commonDef := &fileMap{name: "Common definition"}

	for _, fileMap := range fileMaps {
		if debugMode {
			fileMap.displayOrdered(0, showValue)
		}
		commonDef.digest(fileMap)
		if debugMode {
			commonDef.reorder()
			commonDef.displayOrdered(0, showKind)
		}
	}

	commonDef.reorder()

	return commonDef
}

// puts all the definition contained in the given file map, into the common definition
func (commonDef *fileMap) digest(jsonMap *fileMap) *fileMap {

	// initialisation is needed for the common definition, if it hasn't been done yet
	if commonDef.chainedProperties == nil {

		// initialisation
		commonDef.chainedProperties = map[string]*chainedProperty{}
		commonDef.subMaps = map[string]*fileMap{}

		// ADDING all the chained properties
		for i, propertyName := range jsonMap.orderedProperties {

			// init of the chained property
			currentProperty := &chainedProperty{owner: commonDef, name: propertyName, kind: jsonMap.getPropertyKind(propertyName)}
			commonDef.chainedProperties[propertyName] = currentProperty

			// chaining with the preceding property
			if i > 0 {
				currentProperty.linkAfter(commonDef.chainedProperties[jsonMap.orderedProperties[i-1]], false)
			}
		}

		// initialising the submaps recursively this way:
		for propertyName, submap := range jsonMap.subMaps {
			commonDef.subMaps[propertyName] = (&fileMap{parent: commonDef, name: propertyName}).digest(submap)
		}

	} else { // the definition has already been built once with one file map, so we have to add the missing info here

		// UPGRADING with the missing chained properties
		for _, propertyName := range jsonMap.orderedProperties {

			existingProperty := commonDef.chainedProperties[propertyName]

			// there's a property missing in the whole chain
			if existingProperty == nil {

				// init of a new chained property
				currentProperty := &chainedProperty{owner: commonDef, name: propertyName, kind: jsonMap.getPropertyKind(propertyName), addOn: true}
				commonDef.chainedProperties[propertyName] = currentProperty

				// insertion of this new chained property at the proper place - a big part of the magic is here
				currentProperty.link(jsonMap.chainedProperties[propertyName])

			} else {

				// if the kinds of the properties from the definition, and the current file map differ, then we have a problem
				if existingProperty.kind != jsonMap.getPropertyKind(propertyName) {
					panic(fmt.Sprintf("\n\nWe have a problem here with property: %s"+
						"\ntype of '%s' is %s,"+
						"\nbut type of '%s' is %s",
						propertyName,
						existingProperty.FullString(), existingProperty.kind,
						jsonMap.chainedProperties[propertyName].FullString(), jsonMap.getPropertyKind(propertyName),
					))
				}
			}
		}

		// upgrading the submaps recursively this way:
		for propertyName, submap := range jsonMap.subMaps {
			if commonDef.subMaps[propertyName] == nil {
				commonDef.subMaps[propertyName] = (&fileMap{parent: commonDef, name: propertyName}).digest(submap)
			} else {
				commonDef.subMaps[propertyName].digest(submap)
			}
		}
	}

	return commonDef
}

// building the array of ordered properties from the chained properties
func (commonDef *fileMap) reorder() {

	// (re)init
	commonDef.orderedProperties = []string{}

	// going from chained property to chained property
	for currentProperty := commonDef.oneChainedProperty().root(); currentProperty != nil; currentProperty = currentProperty.next {
		commonDef.orderedProperties = append(commonDef.orderedProperties, currentProperty.name)
	}

	// dealing with the submaps
	for _, submap := range commonDef.subMaps {
		submap.reorder()
	}
}
