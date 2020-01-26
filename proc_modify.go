//------------------------------------------------------------------------------
// handling the modified columns for each JSON maps
//------------------------------------------------------------------------------

package main

import "fmt"

// handling value changes within each original json file
func (commonDef *fileMap) handleModifiedColumns(config *j2tConfig, jsonMaps []*fileMap) error {

	for _, jsonMap := range jsonMaps {
		if errHandle := commonDef.handleModifiedColumn(config, jsonMap); errHandle != nil {
			return errHandle
		}
	}
	return nil
}

// initialising a new inserted column
func (commonDef *fileMap) handleModifiedColumn(config *j2tConfig, jsonMap *fileMap) error {

	for _, modifConfig := range config.ModifiedColumns {

		// checking the existence of the columns mentioned here
		propDef := commonDef.allChainedProperties[modifConfig.SetColumn]
		if propDef == nil {
			return fmt.Errorf("column '%s' does not exist", modifConfig.SetColumn)
		}
		if commonDef.allChainedProperties[modifConfig.When] == nil {
			return fmt.Errorf("column '%s' does not exist", modifConfig.When)
		}

		// ok let's get the property to change, and the property serving as a condition
		prop := jsonMap.findProp(modifConfig.SetColumn)
		when := jsonMap.findProp(modifConfig.When)

		// the property to change exists in this map, and so does the conditional prorp
		if prop != nil && when != nil {

			// if this equals the configured condition then we can perform the modification
			if when.stringValue() == modifConfig.Equals {
				prop.setValue(propDef.kind, modifConfig.ToValue)
			}
		}
	}

	return nil
}
