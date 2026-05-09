package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/petemoore/samfile/v3"
)

func basicToText(_ map[string]any) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, os.Stdin); err != nil {
		log.Fatal(err)
	}
	sb := samfile.NewSAMBasic(buf.Bytes())
	if err := sb.Output(); err != nil {
		log.Fatal(err)
	}
}
