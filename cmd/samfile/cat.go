package main

import (
	"log"
	"os"

	"github.com/winfreddy88/samfile"
)

func cat(arguments map[string]interface{}) {
	imageName := arguments["-i"].(string)
	file := arguments["-f"].(string)
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
		if file != filename {
			continue
		}
		fileFound = true
		f, err := diskImage.File(filename)
		if err != nil {
			// if extracting all files, just warn
			log.Fatalf("Failed to extract %q from disk image %q: %v", filename, imageName, err)
		}
		_, _ = os.Stdout.Write(f.Body)
	}
	if !fileFound {
		log.Fatalf("File %q not found in disk image %q", file, imageName)
	}
}
