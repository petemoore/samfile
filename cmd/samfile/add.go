package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/petemoore/samfile/v3"
)

func add(arguments map[string]any) {
	file := arguments["-f"].(string)
	fileInfo, statError := os.Stat(file)
	if statError != nil {
		log.Fatalf("file %v not found", file)
	}
	if fileInfo.IsDir() {
		log.Fatalf("target directory must be an existing file: %v exists, but is a directory", file)
	}
	imageName := arguments["-i"].(string)
	diskImage, err := samfile.Load(imageName)
	if err != nil {
		log.Fatal(err)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	// -s/--screen adds the input as a SCREEN$ file with the given SAM
	// MODE (1..4), stored as mode byte 0..3.
	if arguments["--screen"] != nil {
		modeStr := arguments["--screen"].(string)
		mode, err := strconv.Atoi(modeStr)
		if err != nil || mode < 1 || mode > 4 {
			log.Fatalf("-s/--screen MODE must be 1..4, got %q", modeStr)
		}
		err = diskImage.AddScreenFile(filepath.Base(file), data, uint8(mode-1))
		if err != nil {
			log.Fatal(err)
		}
		if err = diskImage.Save(imageName); err != nil {
			log.Fatal(err)
		}
		return
	}
	loadAddressStr := arguments["-l"].(string)
	loadAddress, err := strconv.Atoi(loadAddressStr)
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
	err = diskImage.AddCodeFile(filepath.Base(file), data, uint32(loadAddress), uint32(executionAddress))
	if err != nil {
		log.Fatal(err)
	}
	err = diskImage.Save(imageName)
	if err != nil {
		log.Fatal(err)
	}
}
