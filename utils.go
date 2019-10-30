//------------------------------------------------------------------------------
// Utilities
//------------------------------------------------------------------------------

package main

import (
	"fmt"
	"strconv"

	excel "github.com/360EntSecGroup-Skylar/excelize"
)

//------------------------------------------------------------------------------
// Excel file access
//------------------------------------------------------------------------------

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

//------------------------------------------------------------------------------
// Dealing with colors
//------------------------------------------------------------------------------

var black = "#000000"
var white = "#FFFFFF"

var colors = []string{
	"#1A5276",
	"#0E6655",
	"#9A7D0A",
	"#873600",
	"#7B241C",
	"#5B2C6F",
	"#283747",
	"#196F3D",
}

var colorChart = map[string]string{}

// lightening or darkening a color
func getAdjustedColor(color string, addedLight int) string {

	colorKey := fmt.Sprintf("%s/%d", color, addedLight)

	if lightenedColor := colorChart[colorKey]; lightenedColor != "" {
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

	R = R + int64(addedLight)
	G = G + int64(addedLight)
	B = B + int64(addedLight)

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
	colorChart[colorKey] = result

	return result
}
