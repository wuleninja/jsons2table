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

	// if config.NewColumns != nil {
	// 	for _, newDurationColumn := range config.NewColumns.NewDurations {
	// 		if errInsert := commonDef.insertDuration(newDurationColumn); errInsert != nil {
	// 			err("error while adding an extra duration column: %s", errInsert)
	// 		}
	// 	}
	// 	for _, newAdditionColumn := range config.NewColumns.NewSums {
	// 		if errInsert := commonDef.insertSum(newAdditionColumn); errInsert != nil {
	// 			err("error while adding an extra sum column: %s", errInsert)
	// 		}
	// 	}
	// }
	//
	// if debugMode {
	// 	commonDef.reorder()
	// 	commonDef.displayOrdered(0, showKind)
	// }

	return nil
}

// initialising a new inserted column
func (commonDef *fileMap) initInsertedColumn(newCol *newColumn) (*chainedProperty, error) {

	// the column we're inserting right after
	previousProp := commonDef.allChainedProperties[newCol.PutAfter]
	if previousProp == nil {
		return nil, fmt.Errorf("this column does not seem to exist: %s! You might wanna watch for typos", newCol.PutAfter)
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
	localDef.chainedProperties[newCol.Name] = newProperty

	// chaining
	newProperty.insertAfter(previousProp)

	// global registration if the property
	commonDef.register(newProperty)

	return newProperty, nil
}

// initialising a new inserted duration
func (commonDef *fileMap) insertDuration(newCol *newDuration) error {
	newProp, err := commonDef.initInsertedColumn(&newCol.newColumn)
	if err != nil {
		return err
	}
	newProp.kind = reflect.Float64
	newProp.computationDef = newCol
	return nil
}

// initialising a new inserted sum
func (commonDef *fileMap) insertSum(newCol *newSum) error {
	newProp, err := commonDef.initInsertedColumn(&newCol.newColumn)
	if err != nil {
		return err
	}
	newProp.kind = reflect.Float64
	newProp.computationDef = newCol
	return nil
}
