package main

import (
	"log"

	"github.com/petemoore/samfile/v3"
)

func ls(arguments map[string]any) {
	imageName := arguments["-i"].(string)
	diskImage, err := samfile.Load(imageName)
	if err != nil {
		log.Fatal(err)
	}
	dir := diskImage.DiskJournal()
	dir.Output()
}
