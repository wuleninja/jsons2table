//------------------------------------------------------------------------------
// this code is about writing the XLSX file from the common definition,
// with 1 line per original JSON file
//------------------------------------------------------------------------------

package main

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	excel "github.com/360EntSecGroup-Skylar/excelize"
)

const (
	mainSheetName = "results"
)

// writing the excel file
func (commonDef *fileMap) writeExcel(conf *config, jsonMaps []*fileMap) error {

	// creating the file and the main sheet
	excelFile := excel.NewFile()
	excelFile.SetSheetName("Sheet1", mainSheetName)

	// writing out the headers
	headerLine := commonDef.getHeight()
	if errHeader := commonDef.writeHeaders(excelFile, headerLine); errHeader != nil {
		err("could not write the headers. Cause: %s", errHeader)
	}

	// styling
	if errStyle := commonDef.styleHeaders(excelFile, conf); errStyle != nil {
		err("could not style the sheet. Cause: %s", errStyle)
	}

	// writing the content
	if errContent := commonDef.writeLines(excelFile, jsonMaps, headerLine); errContent != nil {
		err("could not write the content. Cause: %s", errContent)
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

		debug("\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		debug("writing line for with '%s'", commonDef.getFullName())

		// always following the order
		for _, property := range commonDef.orderedProperties {

			if subMap := commonDef.subMaps[property]; subMap != nil {
				if errWrite := subMap.writeLine(excelFile, jsonMap.subMaps[property], headerLine, currentLine, even); errWrite != nil {
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

					// yes, so let's copy it into the excel file
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

				// oh, maybe we could do a bit of styling here
				if even {
					style, errNewStyle := excelFile.NewStyle(
						fmt.Sprintf(`{"fill":{"type":"pattern","color":["%s"],"pattern":1}}`, getLightenedColor(commonProp.conf.background, 80)))
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
func (commonDef *fileMap) styleHeaders(excelFile *excel.File, conf *config) error {

	// dealing with the frozen panes
	excelFile.SetPanes(mainSheetName, fmt.Sprintf(`{"freeze":true,"split":true,"x_split":1,"y_split":%d}`, commonDef.getHeight()))

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

		debug("dealing with %s", prop.getFullName())

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
				confItem.background, confItem.foreground))
		if errNewStyle != nil {
			return errNewStyle
		}
		if errSet := excelFile.SetCellStyle(mainSheetName, coord, coord, style); errSet != nil {
			return errSet
		}
	}

	return nil
}

var colorLightened = map[string]string{}

// lightening a color
func getLightenedColor(color string, steps int64) string {

	colorKey := fmt.Sprintf("%s/%d", color, steps)

	if lightenedColor := colorLightened[colorKey]; lightenedColor != "" {
		return lightenedColor
	}

	usePound := false

	if color[0] == '#' {
		color = color[1:len(color)]
		usePound = true
	}

	R, _ := strconv.ParseInt(color[0:2], 16, 8)
	G, _ := strconv.ParseInt(color[2:4], 16, 8)
	B, _ := strconv.ParseInt(color[4:6], 16, 8)

	R = R + steps
	G = G + steps
	B = B + steps

	if R > 255 {
		R = 255
	} else if R < 0 {
		R = 0
	}

	if G > 255 {
		G = 255
	} else if G < 0 {
		G = 0
	}

	if B > 255 {
		B = 255
	} else if B < 0 {
		B = 0
	}

	RR := fmt.Sprintf("%02s", fmt.Sprintf("%x", R))
	GG := fmt.Sprintf("%02s", fmt.Sprintf("%x", G))
	BB := fmt.Sprintf("%02s", fmt.Sprintf("%x", B))

	result := RR + GG + BB
	if usePound {
		result = "#" + result
	}

	// "caching" for later faster retrieval
	colorLightened[colorKey] = result

	return result
}
