package main

import (
	"fmt"
	"os"
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

	fmt.Printf("Hello %s\n", "World!")
}

func err(strfmt string, args ...interface{}) {
	fmt.Printf(strfmt+"\n", args...)
	os.Exit(1)
}
