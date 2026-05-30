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
		body := f.Body
		if convert, _ := arguments["-c"].(bool); convert {
			if out, ok, err := convertBody(body, diskfile); err != nil {
				log.Fatalf("failed to convert %q from disk image %q: %v", filename, imageName, err)
			} else if ok {
				body = out
			}
		}
		_, _ = os.Stdout.Write(body)
	}
	if !fileFound {
		log.Fatalf("file %q not found in disk image %q", file, imageName)
	}
}
