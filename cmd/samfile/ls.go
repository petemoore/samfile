package main

import (
	"log"

	"github.com/petemoore/samfile/v2"
)

func ls(arguments map[string]interface{}) {
	imageName := arguments["-i"].(string)
	diskImage, err := samfile.Load(imageName)
	if err != nil {
		log.Fatal(err)
	}
	dir := diskImage.DiskJournal()
	dir.Output()
}
