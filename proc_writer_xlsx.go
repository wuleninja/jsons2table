//------------------------------------------------------------------------------
// this code is about writing the XLSX file from the common definition,
// with 1 line per original JSON file
//------------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"
	"reflect"

	excel "github.com/360EntSecGroup-Skylar/excelize"
)

const (
	mainSheetName = "results"
)

// writing the excel file
func (commonDef *fileMap) writeExcel(folderPath string, folderInfo os.FileInfo, jsonMaps []*fileMap) error {

	// creating the file and the main sheet
	excelFile := excel.NewFile()
	excelFile.SetSheetName("Sheet1", mainSheetName)

	// writing out the headers
	headerLine := commonDef.getHeight()
	if errHeader := commonDef.writeHeaders(excelFile, headerLine); errHeader != nil {
		err("could not write the headers. Cause: %s", errHeader)
	}

	// writing the content
	if errContent := commonDef.writeLines(excelFile, jsonMaps, headerLine); errContent != nil {
		err("could not write the content. Cause: %s", errContent)
	}

	// savinf the file
	excelFileName := fmt.Sprintf("%s/%s.xlsx", folderPath, folderInfo.Name())
	errSave := excelFile.SaveAs(excelFileName)
	if errSave != nil {
		err("could not save file '%s'. Cause: %s", excelFileName, errSave)
	}

	// we're done
	println(fmt.Sprintf("(re-)created file '%s'", excelFileName))
	return nil
}

// writing the Excel file's headers
func (commonDef *fileMap) writeHeaders(excelFile *excel.File, headerLine int) error {

	debug("\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	if debugMode {
		debug("dealing with section %s", commonDef.getFullName())
	}

	// following the order
	for _, property := range commonDef.orderedProperties {

		if subMap := commonDef.subMaps[property]; subMap != nil {
			if errWrite := subMap.writeHeaders(excelFile, headerLine); errWrite != nil {
				return errWrite
			}
		} else {
			prop := commonDef.chainedProperties[property]
			debug("dealing with property n°%d = %s", prop.index, prop.getFullName())
			setString(excelFile, headerLine, prop.index, property)
		}
	}

	return nil
}

// writing the Excel file's lines, 1 line per JSON file
func (commonDef *fileMap) writeLines(excelFile *excel.File, jsonMaps []*fileMap, headerLine int) error {
	for i, jsonMap := range jsonMaps {
		if errWrite := commonDef.writeLine(excelFile, jsonMap, headerLine, headerLine+i+1); errWrite != nil {
			return errWrite
		}
		println(fmt.Sprintf("successfully treated JSON file: %s", jsonMap.name))
	}
	return nil
}

// writing out 1 JSON file
func (commonDef *fileMap) writeLine(excelFile *excel.File, jsonMap *fileMap, headerLine, currentLine int) error {

	excelFile.SetCellValue("", "", "")

	if jsonMap != nil {

		debug("\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		debug("writing line for with '%s'", commonDef.getFullName())

		// always following the order
		for _, property := range commonDef.orderedProperties {

			if subMap := commonDef.subMaps[property]; subMap != nil {
				if errWrite := subMap.writeLine(excelFile, jsonMap.subMaps[property], headerLine, currentLine); errWrite != nil {
					return errWrite
				}
			} else {

				// checking we're in the right column !
				commonProp := commonDef.chainedProperties[property]
				debug("dealing with property %d - '%s'", commonProp.index, commonProp.getFullName())
				if header := getString(excelFile, headerLine, commonProp.index); header != property {
					err("We have a problem here : at column %d, header says '%s', but we're dealing with '%s'",
						commonProp.index, header, property)
				}

				// does the current JSON have this property ?
				if jsonProp := jsonMap.chainedProperties[property]; jsonProp != nil {
					if commonProp.kind == reflect.Bool {
						setBool(excelFile, currentLine, commonProp.index, jsonMap.values[property].(bool))
					} else if commonProp.kind == reflect.String {
						setString(excelFile, currentLine, commonProp.index, jsonMap.values[property].(string))
					} else if commonProp.kind == reflect.Float64 {
						setFloat(excelFile, currentLine, commonProp.index, jsonMap.values[property].(float64))
					} else {
						err("case unhandled: '%s' (type = %v)", jsonProp.getFullName(), commonProp.kind)
					}
				}
			}
		}
	}

	return nil
}

// getting the cell for the given row and col
func getCell(row int, col int) string {
	coord, errCoord := excel.CoordinatesToCellName(col, row)
	if errCoord != nil {
		err("cound not get coordinates at row '%s' and column '%s'. Cause: %s", row, col, errCoord)
	}
	return coord
}

// getting a string value from the main sheet
func getString(excelFile *excel.File, row int, col int) string {
	value, errGet := excelFile.GetCellValue(mainSheetName, getCell(row, col))
	if errGet != nil {
		err("error while getting at row %d and column %d", row, col)
	}
	return value
}

// setting a bool value into the main sheet
func setBool(excelFile *excel.File, row int, col int, value bool) {
	if errSet := excelFile.SetCellBool(mainSheetName, getCell(row, col), value); errSet != nil {
		err("error while setting value '%v' at row %d and column %d", value, row, col)
	}
}

// setting a string value into the main sheet
func setString(excelFile *excel.File, row int, col int, value string) {
	if errSet := excelFile.SetCellStr(mainSheetName, getCell(row, col), value); errSet != nil {
		err("error while setting value '%v' at row %d and column %d", value, row, col)
	}
}

// setting a float value into the main sheet
func setFloat(excelFile *excel.File, row int, col int, value float64) {
	if errSet := excelFile.SetCellFloat(mainSheetName, getCell(row, col), value, -1, 64); errSet != nil {
		err("error while setting value '%v' at row %d and column %d", value, row, col)
	}
}
