package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/petemoore/samfile/v3"
)

// screenToPNG reads a raw SCREEN$ body from stdin and writes a PNG (or
// GIF, with --format gif) of the decoded MODE 4 screen to stdout. A
// raw body has no directory entry, so the mode defaults to MODE 4
// (mode byte 3) unless overridden with --mode.
func screenToPNG(arguments map[string]any) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, os.Stdin); err != nil {
		log.Fatal(err)
	}
	mode := uint8(3) // MODE 4
	if v, ok := arguments["--mode"]; ok && v != nil {
		m, err := strconv.Atoi(v.(string))
		if err != nil || m < 1 || m > 4 {
			log.Fatalf("--mode must be 1..4, got %q", v)
		}
		mode = uint8(m - 1)
	}
	scr := samfile.NewScreen(buf.Bytes(), mode)
	format := "png"
	if v, ok := arguments["--format"]; ok && v != nil {
		format = v.(string)
	}
	if err := writeScreen(scr, format, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

// writeScreen encodes scr to w in the given format ("png" or "gif").
func writeScreen(scr *samfile.Screen, format string, w io.Writer) error {
	switch format {
	case "gif":
		return scr.WriteGIF(w)
	case "png", "":
		return scr.WritePNG(w)
	default:
		log.Fatalf("unknown --format %q (want png or gif)", format)
		return nil
	}
}
