//------------------------------------------------------------------------------
// the code here is responsible for parsing all the JSON files, and
// building the common structure for them
//------------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// main directory scanning function
func scanDir(folderPath string) (results []*fileMap, err error) {

	// checking the file infos
	fileInfos, errDir := ioutil.ReadDir(folderPath)
	if errDir != nil {
		return nil, fmt.Errorf("error while checking the directory: %s", errDir)
	}

	// iterating over each file
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), ".json") {
			fileMap, errScan := scanFile(folderPath, fileInfo.Name())
			if errScan != nil {
				return nil, fmt.Errorf("error while treating file: %s. Cause: %s", fileInfo.Name(), errScan)
			}
			results = append(results, fileMap)
		}
	}

	return
}

// main file scanning function
func scanFile(folderPath string, fileName string) (*fileMap, error) {

	// the file path
	filePath := folderPath + string(os.PathSeparator) + fileName

	// opening the file
	file, errOpen := os.Open(filePath)
	if errOpen != nil {
		return nil, fmt.Errorf("error while opening the file at path: %s. Cause: %s", filePath, errOpen)
	}
	defer file.Close()

	// reading the file
	fileBytes, errRead := ioutil.ReadAll(file)
	if errRead != nil {
		return nil, fmt.Errorf("error while reading the file at path: %s. Cause: %s", filePath, errRead)
	}

	// unmarshalling the JSON content
	rootMap := &fileMap{name: strings.Replace(fileName, ".json", "", 1)}
	if errUnmarshall := json.Unmarshal(fileBytes, rootMap); errUnmarshall != nil {
		return nil, fmt.Errorf("error while parsing the file at path: %s. Cause: %s", filePath, errUnmarshall)
	}

	// we're fine
	return rootMap, nil
}
