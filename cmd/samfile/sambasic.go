package main

import (
	"bytes"
	"io"
	"os"

	"github.com/winfreddy88/samfile"
)

func basicToText(arguments map[string]interface{}) {
	var buf bytes.Buffer
	io.Copy(&buf, os.Stdin)
	sb := samfile.NewSAMBasic(buf.Bytes())
	sb.Output()
}
