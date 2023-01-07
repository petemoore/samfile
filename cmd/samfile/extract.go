package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/petemoore/samfile"
)

func extract(arguments map[string]interface{}) {
	imageName := arguments["-i"].(string)
	target := "."
	if arguments["-t"] != nil {
		target = arguments["-t"].(string)
	}
	fileInfo, statError := os.Stat(target)
	if statError != nil {
		log.Fatalf("Target directory must be an existing directory: %v not found", target)
	}
	if !fileInfo.IsDir() {
		log.Fatalf("Target directory must be an existing directory: %v exists, but is not a directory", target)
	}
	diskImage, err := samfile.Load(imageName)
	if err != nil {
		log.Fatal(err)
	}
	dir := diskImage.DiskJournal()
	fileFound := false
	for _, diskfile := range dir {
		if !diskfile.Used() {
			continue
		}
		fileFound = true
		filename := diskfile.Name.String()
		f, err := diskImage.File(filename)
		if err != nil {
			log.Printf("WARNING: Could not extract %q: %v", filename, err)
			continue
		}
		localFile := filepath.Join(target, strings.Replace(filename, string([]rune{os.PathSeparator}), "#", -1))
		log.Printf("Saving file %q from disk image %q to file %q", filename, imageName, localFile)
		err = os.WriteFile(localFile, f.Body, 0666)
		if err != nil {
			log.Fatalf("Failed to write file %q: %v", localFile, err)
		}
	}
	if !fileFound {
		log.Printf("WARNING: no files found in disk image %q so nothing extracted.", imageName)
	}
}
