package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {

	// getting the args
	args := os.Args[1:]
	if len(args) == 0 {
		err("Missing the folder name here! Please provide a path to a valid directory")
	}

	// getting the folder path, which should be a valid directory
	folderPath := args[0]
	if _, errPath := os.Stat(folderPath); os.IsNotExist(errPath) {
		err("'%s' is not a valid directory!", folderPath)
	}

	// scanning all the files within the JSON folder
	fileMaps, errScan := scanDir(folderPath)
	if errScan != nil {
		err("error while scanning: %s", errScan)
	}

	// a bit of sorting
	sort.Slice(fileMaps, func(i int, j int) bool {
		return fileMaps[i].name < fileMaps[j].name
	})

	// merging all the maps to determine the common definition
	// TODO
}

func err(strfmt string, args ...interface{}) {
	fmt.Printf(strfmt+"\n", args...)
	os.Exit(1)
}
