//------------------------------------------------------------------------------
// inserting the new columns within the common definition
//------------------------------------------------------------------------------

package main

import (
	"fmt"
	"reflect"
)

// inserting new columns, as configured in the config file
func (commonDef *fileMap) insertNewColumns(config *j2tConfig) error {

	if config.NewColumns != nil {
		for _, newColumn := range config.NewColumns {
			if errInsert := commonDef.insertNewColumn(newColumn); errInsert != nil {
				err("error while adding an extra column: %s", errInsert)
			}
		}
	}

	if debugMode {
		commonDef.reorder()
		commonDef.displayOrdered(0, showKind)
	}

	return nil
}

// initialising a new inserted column
func (commonDef *fileMap) insertNewColumn(newCol *newColumnConfig) error {

	// the column we're inserting right after
	previousProp := commonDef.allChainedProperties[newCol.PutAfter]
	if previousProp == nil {
		return fmt.Errorf("error in the configuration of the new column '%s': this property does not seem to exist: %s!"+
			" You might wanna watch for typos", newCol.Name, newCol.PutAfter)
	}

	// getting the map we're going to insert into
	localDef := previousProp.owner

	// init of the chained property
	newProperty := &chainedProperty{
		owner:     localDef,
		name:      newCol.Name,
		computed:  true,
		maxLength: len(newCol.Name),
	}

	// linking the property at the same level as the previous prop
	localDef.chainedProperties[newCol.Name] = newProperty

	// chaining
	newProperty.insertAfter(previousProp)

	// global registration if the property
	commonDef.register(newProperty)

	// initialising the config
	newProperty.conf = &configItem{
		background: getAdjustedColor(previousProp.conf.background, 2), // almost keeping the same color as the one before
	}

	// initialising the stat - this will depend on the computation definition later on (when we have time)
	newProperty.statistic = &stat{
		owner: newProperty,
		kind:  statKindNUMBER, // this will have to be changed, depending on the property type
	}

	// preparing the stats
	newProperty.kind = reflect.Float64

	// keeping track of the computation definition on the property itself
	newProperty.computationDef = newCol

	// parsing the formula, to check it, and build the format string that's going to be filed sometime later with real column coordinates
	currentPropName := []rune{}
	currentPropReading := false
	formattableFormula := []rune{}
	for _, rune := range []rune(newCol.Formula) {
		switch rune {
		case '{':
			currentPropReading = true
			formattableFormula = append(formattableFormula, '%', 's')
		case '}':
			currentPropReading = false
			prop := commonDef.getProp(path(currentPropName))
			if prop == nil {
				return fmt.Errorf("error in the configuration of the new column '%s': this property does not seem to exist: %s!"+
					" You might wanna watch for typos", newCol.Name, string(currentPropName))
			}
			newCol.columns = append(newCol.columns, prop)
			currentPropName = nil
		default:
			if currentPropReading {
				currentPropName = append(currentPropName, rune)
			} else {
				formattableFormula = append(formattableFormula, rune)
			}
		}
	}

	// at this point, we should have finished building the formattable formula
	newCol.formattableFormula = string(formattableFormula)

	return nil
}

// func logouille(r rune, isReadingProp bool, currentProp []rune, currentFormula []rune) {
// 	println(fmt.Sprintf("Read: %r, Inside prop ? %t, Current Prop: "))
// }
