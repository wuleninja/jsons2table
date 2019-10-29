package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	_ "github.com/tealeg/xlsx"
)

var debugMode bool
var continueMode bool

func debug(strfmt string, params ...interface{}) {
	if debugMode {
		println(fmt.Sprintf(strfmt, params...))
	}
}

func main() {

	// a bit of doc
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] <folder_path>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\navailable flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\narguments:\n")
		fmt.Fprintf(os.Stderr, "  folder_path: mandatory - the path to the folder containing the JSON files\n")
		fmt.Fprintf(os.Stderr, "\n")
	}

	// adding the flags
	flag.BoolVar(&debugMode, "debug", false, "runs the program in debug mode, i.e. with debug messages")
	flag.BoolVar(&continueMode, "continue", false, "runs the program without stopping at the merging step")
	flag.Parse()

	// controlling the args
	if flag.NArg() == 0 {
		println("\n/!\\ the folder path is missing!\n")
		flag.Usage()
		os.Exit(1)
	}
	if flag.NArg() > 1 {
		err("too many arguments, we only need the folder path here!")
	}

	// getting the folder path, which should be a valid directory
	folderPath := flag.Arg(0)
	folderInfo, errPath := os.Stat(folderPath)
	if os.IsNotExist(errPath) {
		err("'%s' is not a valid directory!", folderPath)
	}

	// scanning all the files within the JSON folder
	jsonMaps, errScan := scanDir(folderPath)
	if errScan != nil {
		err("error while scanning: %s", errScan)
	}

	// a bit of sorting, to make sure the treatment is always the same
	sort.Slice(jsonMaps, func(i int, j int) bool {
		return jsonMaps[i].name < jsonMaps[j].name
	})

	// merging all the maps to determine the common definition
	commonDef := merge(jsonMaps)

	// computing the common definition data tree height
	debug("Common definition height is %d", commonDef.getHeight())

	// retrieving or initialising the config
	config := commonDef.getOrInitConfig(folderPath, folderInfo)

	// writing the Excel file
	commonDef.writeExcel(config, jsonMaps)
}

// fatal error handling
func err(strfmt string, args ...interface{}) {
	fmt.Printf(strfmt+"\n", args...)
	os.Exit(1)
}
