//------------------------------------------------------------------------------
// the code here is responsible for finding a common definition for
// all the maps resulting from unmarshalling all the JSON files
//------------------------------------------------------------------------------

package main

import "fmt"

// digesting all the file maps, to build a common definition for them
func merge(jsonMaps []*fileMap) *fileMap {

	commonDef := &fileMap{name: "Common definition", allChainedProperties: map[path]*chainedProperty{}}

	// digesting each JSON
	for _, jsonMap := range jsonMaps {

		// some debugging log
		// if debugMode {
		// 	jsonMap.displayOrdered(0, showValue)
		// }

		// digesting this JSON map into the common definition
		commonDef.digest(jsonMap)

		// some debugging log
		// if debugMode {
		// 	commonDef.reorder()
		// 	commonDef.displayOrdered(0, showKind)
		// }
	}

	// reordering the common definition
	commonDef.reorder()

	// controlling that this JSON map is shown some respect, regarding the common definition
	for _, jsonMap := range jsonMaps {
		commonDef.control(jsonMap)
	}

	// retunring, for what's next, i.e. using this common definition to create tables
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
			currentProperty := &chainedProperty{
				owner:     commonDef,
				name:      propertyName,
				kind:      jsonMap.getPropertyKind(propertyName),
				maxLength: len(propertyName),
			}
			commonDef.chainedProperties[propertyName] = currentProperty

			// chaining with the preceding property
			if i > 0 {
				currentProperty.linkAfter(commonDef.chainedProperties[jsonMap.orderedProperties[i-1]], false)
			}

			// global registration if the property
			commonDef.register(currentProperty)

			// initialising the stats for this property
			currentProperty.initStat(jsonMap)
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
				currentProperty := &chainedProperty{
					owner: commonDef,
					name:  propertyName,
					kind:  jsonMap.getPropertyKind(propertyName),
					addOn: true,
				}
				commonDef.chainedProperties[propertyName] = currentProperty

				// insertion of this new chained property at the proper place - a big part of the magic is here
				currentProperty.link(jsonMap.chainedProperties[propertyName])

				// global registration if the property
				commonDef.register(currentProperty)

				// initialising the stats for this property
				currentProperty.initStat(jsonMap)

			} else {

				// if the kind of the property from the definition, and the one from the file map differ,
				// then we have a problem about a "common" definition
				if existingProperty.kind != jsonMap.getPropertyKind(propertyName) {
					msg := fmt.Sprintf("\n\nWe have a problem here with property: %s"+
						"\ntype of '%s' is %s,"+
						"\nbut type of '%s' is %s",
						propertyName,
						existingProperty.getPath(), existingProperty.kind,
						jsonMap.chainedProperties[propertyName].getPath(), jsonMap.getPropertyKind(propertyName))
					if continueMode {
						println(msg)
					} else {
						err(msg)
					}
				}

				// we do want to update the stats though
				existingProperty.updateStat(jsonMap)
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
	commonDef.propertyIndexes = map[string]int{}

	// going from chained property to chained property
	index := 1
	for currentProperty := commonDef.oneChainedProperty().root(); currentProperty != nil; currentProperty = currentProperty.next {
		commonDef.orderedProperties = append(commonDef.orderedProperties, currentProperty.name)
		commonDef.propertyIndexes[currentProperty.name] = index
		index++
	}

	// dealing with the submaps
	for _, submap := range commonDef.subMaps {
		submap.reorder()
	}
}

// controlling that that the order of the properties is preserved in the common definition
func (commonDef *fileMap) control(jsonMap *fileMap) {

	// going through the JSON map to check that the order between 2 consecutive items
	// is maintained in the common definition
	for i := 1; i < len(jsonMap.orderedProperties); i++ {
		propertyName := jsonMap.orderedProperties[i]
		previousName := jsonMap.orderedProperties[i-1]
		if commonDef.propertyIndexes[propertyName] < commonDef.propertyIndexes[previousName] {
			property := jsonMap.chainedProperties[propertyName]
			previous := jsonMap.chainedProperties[previousName]
			err("\nThere's a problem here: "+
				"\nwe have '%s' before '%s'"+
				"\nbut the order is inversed in the common definition!",
				property.getPath(), previous.getPath(),
			)
		}
	}

	// dealing with the submaps
	for _, submap := range jsonMap.subMaps {
		commonDef.subMaps[submap.name].control(submap)
	}
}

// global-indexing each final properties, which will also correspond to a column number in the forecoming Excel file
func (commonDef *fileMap) index(currentIndex *int) {
	for _, property := range commonDef.orderedProperties {
		if subMap := commonDef.subMaps[property]; subMap != nil {
			subMap.index(currentIndex)
		} else {
			commonDef.chainedProperties[property].index = *currentIndex
			*currentIndex = *currentIndex + 1
		}
	}
}
