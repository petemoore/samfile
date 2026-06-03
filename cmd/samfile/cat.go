package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/petemoore/samfile/v3"
	"github.com/petemoore/samfile/v3/audio"
)

func cat(arguments map[string]any) {
	imageName := arguments["-i"].(string)
	file := arguments["-f"].(string)
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
		if convert {
			copts.fileName = filename
			if out, ok, err := convertBody(body, diskfile, copts); err != nil {
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

// parseConvertOptions reads --audio-format and --loops from docopt arguments
// and computes provenance fields from the image path and bytes.
func parseConvertOptions(arguments map[string]any, imageName string, imageBytes []byte) (convertOptions, error) {
	audioFormat := "mp3"
	if v, _ := arguments["--audio-format"].(string); v != "" {
		audioFormat = v
	}
	supported := audio.Supported()
	validFormat := false
	for _, s := range supported {
		if s == audioFormat {
			validFormat = true
			break
		}
	}
	if !validFormat {
		return convertOptions{}, fmt.Errorf("unsupported audio format %q; choose one of: wav, flac, mp3", audioFormat)
	}

	loops := 4
	if v, _ := arguments["--loops"].(string); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return convertOptions{}, fmt.Errorf("--loops must be a positive integer, got %q", v)
		}
		loops = n
	}

	sum := sha256.Sum256(imageBytes)
	diskSHA256 := fmt.Sprintf("%x", sum)
	diskName := filepath.Base(imageName)

	ver := version
	if ver == "" {
		ver = "v3"
	}

	return convertOptions{
		audioFormat: audioFormat,
		loops:       loops,
		diskName:    diskName,
		diskSHA256:  diskSHA256,
		version:     ver,
	}, nil
}
