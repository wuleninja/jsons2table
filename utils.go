//------------------------------------------------------------------------------
// Utilities
//------------------------------------------------------------------------------

package main

import (
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
