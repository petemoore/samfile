package main

import (
	"log"
	"os"

	"github.com/petemoore/samfile/v3"
)

func cat(arguments map[string]any) {
	imageName := arguments["-i"].(string)
	file := arguments["-f"].(string)
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
		filename := diskfile.Name.String()
		if file != filename {
			continue
		}
		fileFound = true
		f, err := diskImage.File(filename)
		if err != nil {
			log.Fatalf("failed to extract %q from disk image %q: %v", filename, imageName, err)
		}
		_, _ = os.Stdout.Write(f.Body)
	}
	if !fileFound {
		log.Fatalf("file %q not found in disk image %q", file, imageName)
	}
}
