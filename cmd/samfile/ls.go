package main

import (
	"log"

	"github.com/winfreddy88/samfile"
)

func ls(arguments map[string]interface{}) {
	imageName := arguments["IMAGE"].(string)
	diskImage, err := samfile.Load(imageName)
	if err != nil {
		log.Fatal(err)
	}
	dir := diskImage.DirectoryListing()
	dir.Output()
}
