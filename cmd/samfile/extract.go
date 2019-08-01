package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/winfreddy88/samfile"
)

func extract(arguments map[string]interface{}) {
	imageName := arguments["IMAGE"].(string)
	file := ""
	extractAllFiles := true
	if arguments["FILE"] != nil {
		file = arguments["FILE"].(string)
		extractAllFiles = false
	}
	target := "."
	if arguments["--dest"] != nil {
		target = arguments["--dest"].(string)
	}
	targetExists := false
	targetIsDir := false
	fileInfo, statError := os.Stat(target)
	if statError == nil {
		targetExists = true
		targetIsDir = fileInfo.IsDir()
	}
	if extractAllFiles {
		if !targetExists {
			log.Fatalf("When extracting all files, TARGET must be an existing directory: %v not found", target)
		}
		if !targetIsDir {
			log.Fatalf("When extracting all files, TARGET must be an existing directory: %v exists, but is not a directory", target)
		}
	}
	diskImage, err := samfile.Load(imageName)
	if err != nil {
		log.Fatal(err)
	}
	dir := diskImage.DirectoryListing()
	fileFound := false
	for _, diskfile := range dir {
		if diskfile.Free() {
			continue
		}
		filename := diskfile.Name.String()
		if file != filename && !extractAllFiles {
			continue
		}
		fileFound = true
		f, err := diskImage.File(filename)
		if err != nil {
			// if extracting all files, just warn
			if extractAllFiles {
				log.Printf("WARNING: Could not extract %q: %v", filename, err)
				continue
			}
			log.Fatalf("Failed to extract %q from disk image %q: %v", filename, imageName, err)
		}
		localFile := ""
		if targetExists && targetIsDir {
			localFile = filepath.Join(target, strings.Replace(filename, string([]rune{os.PathSeparator}), "#", -1))
		} else {
			localFile = target
		}
		log.Printf("Saving file %q from disk image %q to file %q", filename, imageName, localFile)
		err = ioutil.WriteFile(localFile, f.Body, 0666)
		if err != nil {
			log.Fatalf("Failed to write file %q: %v", localFile, err)
		}
	}
	if !fileFound {
		if extractAllFiles {
			log.Printf("WARNING: no files found in disk image %q so nothing extracted.", imageName)
		} else {
			log.Fatalf("File %q not found in disk image %q", file, imageName)
		}
	}
}
