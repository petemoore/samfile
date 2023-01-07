package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/petemoore/samfile/v2"
)

func add(arguments map[string]interface{}) {
	file := arguments["-f"].(string)
	fileInfo, statError := os.Stat(file)
	if statError != nil {
		log.Fatalf("File %v not found", file)
	}
	if fileInfo.IsDir() {
		log.Fatalf("Target directory must be an existing file: %v exists, but is a directory", file)
	}
	loadAddressStr := arguments["-l"].(string)
	loadAddress, err := strconv.Atoi(loadAddressStr)
	if err != nil {
		log.Fatal(err)
	}
	imageName := arguments["-i"].(string)
	diskImage, err := samfile.Load(imageName)
	if err != nil {
		log.Fatal(err)
	}
	executionAddress := 0
	if arguments["-e"] != nil {
		executionAddressStr := arguments["-e"].(string)
		executionAddress, err = strconv.Atoi(executionAddressStr)
		if err != nil {
			log.Fatal(err)
		}
	}
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	err = diskImage.AddCodeFile(filepath.Base(file), data, uint32(loadAddress), uint32(executionAddress))
	if err != nil {
		log.Fatal(err)
	}
	err = diskImage.Save(imageName)
	if err != nil {
		log.Fatal(err)
	}
}
