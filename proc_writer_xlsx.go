//------------------------------------------------------------------------------
// this code is about writing the XLSX file from the common definition,
// with 1 line per original JSON file
//------------------------------------------------------------------------------

package main

import (
	"fmt"
	"math"
	"reflect"

	excel "github.com/360EntSecGroup-Skylar/excelize"
)

const (
	mainSheetName = "results"
)

// writing the excel file
func (commonDef *fileMap) writeExcel(conf *j2tConfig, jsonMaps []*fileMap) error {

	// reordering - just to be sure - then computing the index for each final property contained within the definition
	commonDef.reorder()
	currentIndex := 1
	commonDef.index(&currentIndex)

	// creating the file and the main sheet
	excelFile := excel.NewFile()
	excelFile.SetSheetName("Sheet1", mainSheetName)

	// writing out the headers
	headerLine := commonDef.getHeight()
	if errHeader := commonDef.writeHeaders(excelFile, headerLine); errHeader != nil {
		err("could not write the headers. Cause: %s", errHeader)
	}

	// styling the headers
	if errStyle := commonDef.styleHeaders(excelFile, conf); errStyle != nil {
		err("could not style the sheet. Cause: %s", errStyle)
	}

	// writing the content
	if errContent := commonDef.writeLines(excelFile, jsonMaps, headerLine); errContent != nil {
		err("could not write the content. Cause: %s", errContent)
	}

	// writing some stats
	footerLine := headerLine + len(jsonMaps) + 2
	if errStat := commonDef.writeStats(excelFile, headerLine, footerLine, len(jsonMaps)); errStat != nil {
		err("could not write the stats. Cause: %s", errStat)
	}

	// saving the file
	excelFileName := fmt.Sprintf("%s/%s.xlsx", conf.folderPath, conf.folderInfo.Name())
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

	log("\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	if debugMode {
		log("dealing with section %s", commonDef.getFullName())
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
			log("dealing with property n°%d = %s", prop.index, prop.getPath())

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
			columnSize := math.Ceil(float64(prop.maxLength) * 1.15)
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
		if errWrite := commonDef.writeLine(excelFile, jsonMap, headerLine, headerLine+i+1, i%2 == 0); errWrite != nil {
			return errWrite
		}
		println(fmt.Sprintf("successfully treated JSON file: %s", jsonMap.name))
	}
	return nil
}

// writing out 1 JSON file
func (commonDef *fileMap) writeLine(excelFile *excel.File, jsonMap *fileMap, headerLine, currentLine int, even bool) error {

	// excelFile.SetCellValue("", "", "")

	if jsonMap != nil {

		log("\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		log("writing line for with '%s'", commonDef.getFullName())

		// always following the order
		for _, property := range commonDef.orderedProperties {

			if subMap := commonDef.subMaps[property]; subMap != nil {
				if errWrite := subMap.writeLine(excelFile, jsonMap.subMaps[property], headerLine, currentLine, even); errWrite != nil {
					return errWrite
				}
			} else {

				// checking we're in the right column !
				commonProp := commonDef.chainedProperties[property]
				log("dealing with property %d - '%s'", commonProp.index, commonProp.getPath())
				if header := getString(excelFile, headerLine, commonProp.index); header != property {
					err("We have a problem here : at column %d, header says '%s', but we're dealing with '%s'",
						commonProp.index, header, property)
				}

				// does the current JSON have this property ?
				if jsonProp := jsonMap.chainedProperties[property]; jsonProp != nil {

					// yes, so let's copy it into the excel file
					if commonProp.kind == reflect.Bool {
						setBool(excelFile, currentLine, commonProp.index, jsonMap.values[property].(bool))
					} else if commonProp.kind == reflect.String {
						setString(excelFile, currentLine, commonProp.index, jsonMap.values[property].(string))
					} else if commonProp.kind == reflect.Float64 {
						if value := jsonMap.values[property].(float64); value != -999999 {
							setFloat(excelFile, currentLine, commonProp.index, value)
						}
					} else {
						err("case unhandled: '%s' (type = %v)", jsonProp.getPath(), commonProp.kind)
					}
				}

				// oh, maybe we could do a bit of styling here
				if even {
					style, errNewStyle := excelFile.NewStyle(
						fmt.Sprintf(`{"fill":{"type":"pattern","color":["%s"],"pattern":1}}`, getAdjustedColor(commonProp.conf.background, 90)))
					if errNewStyle != nil {
						return errNewStyle
					}
					coord := getCell(currentLine, commonProp.index)
					if errSet := excelFile.SetCellStyle(mainSheetName, coord, coord, style); errSet != nil {
						return errSet
					}
				}
			}
		}
	}

	return nil
}

// apply a basic style on the Excel file
func (commonDef *fileMap) styleHeaders(excelFile *excel.File, conf *j2tConfig) error {

	style := fmt.Sprintf(`{"freeze":true,"split":false,"x_split":1,"y_split":%d,"top_left_cell":"B10","active_pane":"topRight","panes":[{"sqref":"B10","active_cell":"B10","pane":"topRight"}]}`, commonDef.getHeight())
	println(style)

	// dealing with the frozen panes
	excelFile.SetPanes(mainSheetName, style)

	// applying some colors
	if errColor := commonDef.applyColor(excelFile, conf.content); errColor != nil {
		return errColor
	}

	return nil
}

// applying colors configured in the given config map for the current common definition map
func (commonDef *fileMap) applyColor(excelFile *excel.File, confMap *configMap) error {

	// coloring each property
	for propertyName, prop := range commonDef.chainedProperties {

		// getting the right config item
		confItem := confMap.items[propertyName]

		// where should the style apply ?
		var coord string

		// is it a section we're dealing with ?
		if subMap := commonDef.subMaps[propertyName]; subMap != nil {

			coord = getCell(subMap.getDepth()-1, subMap.getFirstIndex())

			// going under
			subMap.applyColor(excelFile, confMap.subMaps[propertyName])

		} else {
			coord = getCell(prop.owner.getDepth(), prop.index)
		}

		// applying the style to the current property
		style, errNewStyle := excelFile.NewStyle(
			fmt.Sprintf(`{"fill":{"type":"pattern","color":["%s"],"pattern":1}, "font":{"color":"%s"}, "alignment":{"horizontal":"center"}}`,
				confItem.background, white))
		if errNewStyle != nil {
			return errNewStyle
		}
		if errSet := excelFile.SetCellStyle(mainSheetName, coord, coord, style); errSet != nil {
			return errSet
		}
	}

	return nil
}
