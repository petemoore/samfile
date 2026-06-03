package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/petemoore/samfile/v3"
)

func extract(arguments map[string]any) {
	imageName := arguments["-i"].(string)
	target := "."
	if arguments["-t"] != nil {
		target = arguments["-t"].(string)
	}
	fileInfo, statError := os.Stat(target)
	if statError != nil {
		log.Fatalf("target directory must be an existing directory: %v not found", target)
	}
	if !fileInfo.IsDir() {
		log.Fatalf("target directory must be an existing directory: %v exists, but is not a directory", target)
	}
	diskImage, err := samfile.Load(imageName)
	if err != nil {
		log.Fatal(err)
	}

	convert, _ := arguments["-c"].(bool)

	var copts convertOptions
	if convert {
		imageBytes, err := os.ReadFile(imageName)
		if err != nil {
			log.Fatalf("failed to read disk image %q for provenance: %v", imageName, err)
		}
		copts, err = parseConvertOptions(arguments, imageName, imageBytes)
		if err != nil {
			log.Fatal(err)
		}
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
			log.Printf("warning: could not extract %q: %v", filename, err)
			continue
		}
		body := f.Body
		base := strings.ReplaceAll(filename, string([]rune{os.PathSeparator}), "#")
		if convert {
			copts.fileName = filename
			if out, ok, err := convertBody(body, diskfile, copts); err != nil {
				log.Printf("warning: could not convert %q, extracting raw: %v", filename, err)
			} else if ok {
				body = out
				base, _ = convertedName(base, diskfile, f.Body, copts)
			}
		}
		localFile := filepath.Join(target, base)
		log.Printf("saving file %q from disk image %q to file %q", filename, imageName, localFile)
		err = os.WriteFile(localFile, body, 0666)
		if err != nil {
			log.Fatalf("failed to write file %q: %v", localFile, err)
		}
	}
	if !fileFound {
		log.Printf("warning: no files found in disk image %q so nothing extracted.", imageName)
	}
}
