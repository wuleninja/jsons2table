//------------------------------------------------------------------------------
// putting some stats into the Excel file
//------------------------------------------------------------------------------

package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// incrementing the stat for a given value
func (thisStat *stat) countUp(value interface{}) *stat {

	// seeing the given value as a string
	key := fmt.Sprintf("%v", value)

	// counting up how many times we've seen it
	thisStat.valueCounts[key] = thisStat.valueCounts[key] + 1

	// updating the property's max length if needed
	if newLength := len(key); newLength > thisStat.owner.maxLength {
		thisStat.owner.maxLength = newLength
	}

	return thisStat
}

// initialising a stat for a property
func (thisProp *chainedProperty) initStat(jsonMap *fileMap) {

	// initialising the stat
	thisProp.statistic = (&stat{
		owner:       thisProp,
		valueCounts: map[string]int{},
	}).countUp(jsonMap.values[thisProp.name])

	// maybe we can determine the type right away
	switch jsonMap.values[thisProp.name].(type) {
	case bool:
		thisProp.statistic.kind = statKindBOOLEAN
	case float64:
		thisProp.statistic.kind = statKindNUMBER
	case string:
		thisProp.statistic.kind = statKindTEXT
	}
}

// updating a stat for a property
func (thisProp *chainedProperty) updateStat(jsonMap *fileMap) {
	thisProp.statistic.countUp(jsonMap.values[thisProp.name])
}

// writing all the stats
func (commonDef *fileMap) writeStats(excelFile *excelize.File, headerLine int, footerLine int, nbRows int) error {

	for _, property := range commonDef.orderedProperties {

		if subMap := commonDef.subMaps[property]; subMap != nil {

			// going under
			if errWrite := subMap.writeStats(excelFile, headerLine, footerLine, nbRows); errWrite != nil {
				return fmt.Errorf("Error while treating section '%s': %s", property, errWrite)
			}
		} else {

			// first we need to know what kind of stat we have
			if errDetect := commonDef.chainedProperties[property].detectStat(); errDetect != nil {
				return errDetect
			}

			// let's write the stat now
			if errWrite := commonDef.chainedProperties[property].writeStat(excelFile, headerLine, footerLine, nbRows); errWrite != nil {
				return errWrite
			}
		}
	}

	return nil
}

// detecting the stat kind
func (thisProp *chainedProperty) detectStat() error {

	// if we still haven't found out about this property's stat type, let's try figuring it out
	if thisProp.statistic.kind == statKindTEXT {

		// iterating over all the possible values
		for value, count := range thisProp.statistic.valueCounts {

			// we can't state anything from an empty value
			if value != "" {

				// do we have a date here ?
				if isDate(value) {

					if thisProp.statistic.kind == statKindTEXT {
						thisProp.statistic.kind = statKindDATE

					} else if thisProp.statistic.kind != statKindDATE {
						return fmt.Errorf("value '%s' in column %s is a date, whereas the the column type has been detected as of type '%s'",
							value, thisProp.getFullName(), thisProp.statistic.kind)
					}

					// if not, we might have a category
				} else if count > 1 { // twice the same text in a column, this is probably a category
					thisProp.statistic.kind = statKindCATEGORY
					return nil
				}
			}
		}
	}

	// for numbers, let's find out if we need decimals of not
	if thisProp.statistic.kind == statKindNUMBER {
		for valueString := range thisProp.statistic.valueCounts {
			value, _ := strconv.ParseFloat(valueString, 64)
			thisProp.statistic.decimal = thisProp.statistic.decimal || float64(int64(value)) != value
		}
	}

	return nil
}

// writing a particular stat
func (thisProp *chainedProperty) writeStat(excelFile *excelize.File, headerLine, footerLine, nbRows int) error {

	// handling the "header" for this stat
	setString(excelFile, footerLine, thisProp.index, thisProp.name)

	style, errNewStyle := excelFile.NewStyle(
		fmt.Sprintf(`{"fill":{"type":"pattern","color":["%s"],"pattern":1}, "font":{"color":"%s"}, "alignment":{"horizontal":"center"}}`,
			thisProp.conf.background, white))
	if errNewStyle != nil {
		return errNewStyle
	}
	coord := getCell(footerLine, thisProp.index)
	if errSet := excelFile.SetCellStyle(mainSheetName, coord, coord, style); errSet != nil {
		return errSet
	}

	// writing out the stat type
	setString(excelFile, footerLine+1, thisProp.index, string(thisProp.statistic.kind))

	style, errNewStyle = excelFile.NewStyle(`{"font":{"italic":true}, "alignment":{"horizontal":"center"}}`)
	if errNewStyle != nil {
		return errNewStyle
	}
	coord = getCell(footerLine+1, thisProp.index)
	if errSet := excelFile.SetCellStyle(mainSheetName, coord, coord, style); errSet != nil {
		return errSet
	}

	// computing the first & last cell coordinates for the current column; and total number of rows
	firstCell := getCell(headerLine+1, thisProp.index)
	lastCell := getCell(footerLine-2, thisProp.index)

	// the stats begin here
	statLine := footerLine + 2

	// writing out the stats for this column
	switch thisProp.statistic.kind {
	case statKindBOOLEAN:
		return thisProp.writeBooleanStats(excelFile, firstCell, lastCell, statLine, nbRows)
	case statKindDATE:
	case statKindCATEGORY:
		return thisProp.writeCategoryStats(excelFile, firstCell, lastCell, statLine, nbRows)
	case statKindNUMBER:
		return thisProp.writeNumberStats(excelFile, firstCell, lastCell, statLine, nbRows)
	}

	return nil
}

// writing formulae useful to treat a boolean statistic
func (thisProp *chainedProperty) writeBooleanStats(excelFile *excelize.File, firstCell, lastCell string, statLine, nbRows int) error {
	if err := thisProp.writeCategoryValue(excelFile, "TRUE", firstCell, lastCell, 0, statLine, nbRows); err != nil {
		return err
	}
	if err := thisProp.writeCategoryValue(excelFile, "FALSE", firstCell, lastCell, 1, statLine, nbRows); err != nil {
		return err
	}
	return thisProp.writeCategoryValue(excelFile, "", firstCell, lastCell, 2, statLine, nbRows)
}

// writing a category value
func (thisProp *chainedProperty) writeCategoryValue(excelFile *excelize.File, value, firstCell, lastCell string, index, statLine, nbRows int) error {

	// which row do we start from ?
	i := statLine + 3*index

	// writing the value name
	valueName := "- NC -"
	if value != "" {
		valueName = value
	}
	setString(excelFile, i, thisProp.index, valueName)
	style, errStyle := excelFile.NewStyle(`{"font":{"bold":true}}`)
	if errStyle != nil {
		return errStyle
	}
	if errSet := excelFile.SetCellStyle(mainSheetName, getCell(i, thisProp.index), getCell(i, thisProp.index), style); errSet != nil {
		return errSet
	}

	// counting the occurrences
	if errSet := excelFile.SetCellFormula(mainSheetName, getCell(i+1, thisProp.index), "COUNTIF("+firstCell+":"+lastCell+", \""+value+"\")"); errSet != nil {
		return errSet
	}

	// computing the percentage
	if errSet := excelFile.SetCellFormula(mainSheetName, getCell(i+2, thisProp.index), getCell(i+1, thisProp.index)+"/"+strconv.Itoa(nbRows)); errSet != nil {
		return errSet
	}
	style, errStyle = excelFile.NewStyle(`{"custom_number_format": "0.0 %"}`)
	if errStyle != nil {
		return errStyle
	}
	if errSet := excelFile.SetCellStyle(mainSheetName, getCell(i+2, thisProp.index), getCell(i+2, thisProp.index), style); errSet != nil {
		return errSet
	}

	// nah, it's ok now
	return nil
}

// writing the stats for a category column
func (thisProp *chainedProperty) writeCategoryStats(excelFile *excelize.File, firstCell, lastCell string, statLine, nbRows int) error {

	// getting an ordered list for the values; sorting is done by value count, descending
	values := []string{}
	for value := range thisProp.statistic.valueCounts {
		values = append(values, value)
	}
	sort.Slice(values, func(i int, j int) bool {
		count1 := thisProp.statistic.valueCounts[values[i]]
		count2 := thisProp.statistic.valueCounts[values[j]]
		if count1 == count2 {
			return values[i] < values[j]
		}
		return count1 > count2
	})

	// now let's rool
	for i, value := range values {
		if err := thisProp.writeCategoryValue(excelFile, value, firstCell, lastCell, i, statLine, nbRows); err != nil {
			return err
		}
	}

	// let's deal with the empty values
	return thisProp.writeCategoryValue(excelFile, "", firstCell, lastCell, len(values), statLine, nbRows)
}

// writing the stats for a number column
func (thisProp *chainedProperty) writeNumberStats(excelFile *excelize.File, firstCell, lastCell string, statLine, nbRows int) error {
	numberFormat := "0"
	if thisProp.statistic.decimal {
		numberFormat = "0.00"
	}
	if err := thisProp.writeNumberStatFn(excelFile, firstCell, lastCell, 0, statLine, "MIN", numberFormat); err != nil {
		return err
	}
	if err := thisProp.writeNumberStatFn(excelFile, firstCell, lastCell, 1, statLine, "MAX", numberFormat); err != nil {
		return err
	}
	if err := thisProp.writeNumberStatFn(excelFile, firstCell, lastCell, 2, statLine, "MEDIAN", numberFormat); err != nil {
		return err
	}
	if err := thisProp.writeNumberStatFn(excelFile, firstCell, lastCell, 3, statLine, "AVERAGE", numberFormat); err != nil {
		return err
	}
	if err := thisProp.writeNumberStatFn(excelFile, firstCell, lastCell, 4, statLine, "STDEV.P", numberFormat); err != nil {
		return err
	}
	if err := thisProp.writeNumberStatFn(excelFile, firstCell, lastCell, 5, statLine, "COUNTA", ""); err != nil {
		return err
	}
	return nil
}

// writing a particular stat for a number column
func (thisProp *chainedProperty) writeNumberStatFn(excelFile *excelize.File,
	firstCell, lastCell string, index, statLine int, function string, customNumberFormat string) error {

	// which row do we start from ?
	i := statLine + 2*index

	// the function name
	setString(excelFile, i, thisProp.index, function)

	// writing down the formula for the function
	if errSet := excelFile.SetCellFormula(mainSheetName, getCell(i+1, thisProp.index), function+"("+firstCell+":"+lastCell+")"); errSet != nil {
		return errSet
	}
	if customNumberFormat != "" {
		style, errStyle := excelFile.NewStyle(fmt.Sprintf(`{"custom_number_format": "%s"}`, customNumberFormat))
		if errStyle != nil {
			return errStyle
		}
		if errSet := excelFile.SetCellStyle(mainSheetName, getCell(i+1, thisProp.index), getCell(i+1, thisProp.index), style); errSet != nil {
			return errSet
		}
	}

	return nil
}

//------------------------------------------------------------------------------
// Dealing with special types of values
//------------------------------------------------------------------------------

// is the given value a date ?
func isDate(value string) bool {
	if value == "" {
		return false
	}
	_, err := time.Parse("02/01/2006", value)
	if err != nil {
		_, err = time.Parse("01/2006", value)
	}
	return err == nil
}
