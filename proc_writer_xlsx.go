//------------------------------------------------------------------------------
// this code is about writing the XLSX file from the common definition,
// with 1 line per original JSON file
//------------------------------------------------------------------------------

package main

import (
	"fmt"
	"math"
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

	// styling
	if errStyle := commonDef.style(excelFile); errStyle != nil {
		err("could not style the sheet. Cause: %s", errStyle)
	}

	// saving the file
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

			// writing the name for this section
			setString(excelFile, subMap.getDepth()-1, subMap.getFirstIndex(), subMap.name)

			// merging the whole section
			excelFile.MergeCell(mainSheetName,
				getCell(subMap.getDepth()-1, subMap.getFirstIndex()),
				getCell(subMap.getDepth()-1, subMap.getLastIndex()))

			// dealing with what's below
			if errWrite := subMap.writeHeaders(excelFile, headerLine); errWrite != nil {
				return errWrite
			}

		} else {

			// retrieving the right property
			prop := commonDef.chainedProperties[property]
			debug("dealing with property n°%d = %s", prop.index, prop.getFullName())

			// writing it
			setString(excelFile, prop.owner.getDepth(), prop.index, property)

			// merging till the header line
			excelFile.MergeCell(mainSheetName,
				getCell(prop.owner.getDepth(), prop.index),
				getCell(headerLine, prop.index))

			// adjusting the column size
			column, errCol := excel.ColumnNumberToName(prop.index)
			if errCol != nil {
				return errCol
			}
			columnSize := math.Ceil(float64(len(property)) * 1.15)
			if columnSize < 8 {
				columnSize = 8
			}
			if errWidth := excelFile.SetColWidth(mainSheetName, column, column, columnSize); errWidth != nil {
				return errWidth
			}
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

// apply a basic style on the Excel file
func (commonDef *fileMap) style(excelFile *excel.File) error {

	// dealing with the alignment
	style, errNewStyle := excelFile.NewStyle(`{"alignment":{"horizontal":"center"}}`)
	if errNewStyle != nil {
		return errNewStyle
	}
	if errSet := excelFile.SetCellStyle(mainSheetName, "A1", getCell(commonDef.getHeight()-1, commonDef.getLastIndex()), style); errSet != nil {
		return errSet
	}

	// dealing with the frozen panes
	excelFile.SetPanes(mainSheetName, fmt.Sprintf(`{"freeze":true,"split":true,"x_split":1,"y_split":%d}`, commonDef.getHeight()))

	return nil
}
