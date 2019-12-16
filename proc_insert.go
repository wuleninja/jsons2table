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
		for _, newDurationColumn := range config.NewColumns.NewDurations {
			if errInsert := commonDef.insertDuration(newDurationColumn); errInsert != nil {
				err("error while adding an extra duration column: %s", errInsert)
			}
		}
		for _, newAdditionColumn := range config.NewColumns.NewSums {
			if errInsert := commonDef.insertSum(newAdditionColumn); errInsert != nil {
				err("error while adding an extra sum column: %s", errInsert)
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

	// initialising the config
	newProperty.conf = &configItem{
		background: getAdjustedColor(previousProp.conf.background, 15),
	}

	// initialising the stat - this will depend on the computation definition later on (when we have time)
	newProperty.statistic = &stat{
		owner: newProperty,
		kind:  statKindNUMBER, // this will have to be changed, depending on the property type
	}

	return newProperty, nil
}

// initialising a new inserted duration
func (commonDef *fileMap) insertDuration(newCol *newDuration) error {
	newProp, errInit := commonDef.initInsertedColumn(&newCol.newColumn)
	if errInit != nil {
		return errInit
	}
	newProp.kind = reflect.Float64
	newProp.computationDef = newCol
	// checking the duration is well configured
	if commonDef.allChainedProperties[newCol.FromDate] == nil {
		err("error in the configuration of new duration '%s': from date '%s' does not exist. This may be a typo.", newCol.Name, newCol.FromDate)
	}
	if commonDef.allChainedProperties[newCol.ToDate] == nil {
		err("error in the configuration of new duration '%s': to date '%s' does not exist. This may bea typo.", newCol.Name, newCol.ToDate)
	}
	return nil
}

// initialising a new inserted sum
func (commonDef *fileMap) insertSum(newCol *newSum) error {
	newProp, errInit := commonDef.initInsertedColumn(&newCol.newColumn)
	if errInit != nil {
		return errInit
	}
	// checking the columns are well configured
	for _, added := range newCol.AddTogether {
		if commonDef.allChainedProperties[added] == nil {
			err("error in the configuration of new sum '%s': column '%s' does not exist. This may be a typo.", newCol.Name, added)
		}
	}
	for _, substracted := range newCol.AddTogether {
		if commonDef.allChainedProperties[substracted] == nil {
			err("error in the configuration of new sum '%s': column '%s' does not exist. This may be a typo.", newCol.Name, substracted)
		}
	}
	newProp.kind = reflect.Float64
	newProp.computationDef = newCol
	return nil
}
