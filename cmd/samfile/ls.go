package main

import (
	"log"

	"github.com/petemoore/samfile"
)

func ls(arguments map[string]interface{}) {
	imageName := arguments["-i"].(string)
	diskImage, err := samfile.Load(imageName)
	if err != nil {
		log.Fatal(err)
	}
	dir := diskImage.DirectoryListing()
	dir.Output()
}
