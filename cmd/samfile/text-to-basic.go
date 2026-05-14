package main

import (
	"log"
	"os"

	"github.com/petemoore/samfile/v3/sambasic"
)

func textToBasic(_ map[string]any) {
	f, err := sambasic.ParseText(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(f.ProgBytes()); err != nil {
		log.Fatal(err)
	}
}
