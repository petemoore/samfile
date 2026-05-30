package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"testing"

	docopt "github.com/docopt/docopt-go"
	"github.com/petemoore/samfile/v3"
)

// mode4Body builds a minimal valid MODE 4 SCREEN$ body: a solid screen
// of CLUT entry 0, palette index 0x49 (a mid-red), with the canonical
// 24617-byte trailer and no line interrupts.
func mode4Body() []byte {
	body := make([]byte, 24617)
	body[24576] = 0x49 // CLUT A entry 0
	body[24576+20] = 0x49
	body[len(body)-1] = 0xFF
	return body
}

// blankDiskWithScreen writes a fresh MGT image containing one SCREEN$
// file ("SCR") added via AddScreenFile, and returns its path.
func blankDiskWithScreen(t *testing.T) string {
	t.Helper()
	di := samfile.NewDiskImage()
	if err := di.AddScreenFile("SCR", mode4Body(), 3); err != nil {
		t.Fatalf("AddScreenFile: %v", err)
	}
	path := filepath.Join(t.TempDir(), "screen.mgt")
	if err := di.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path
}

// captureStdout runs f with os.Stdout redirected to a pipe and returns
// what f wrote. Mirrors the redirection pattern in cat_test.go.
func captureStdout(t *testing.T, f func()) []byte {
	t.Helper()
	old := os.Stdout
	defer func() { os.Stdout = old }()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	done := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		done <- buf.Bytes()
	}()
	f()
	w.Close()
	return <-done
}

// TestCatConvertScreenToPNG cats a SCREEN$ file with -c and confirms
// the output is a valid 256x192 PNG.
func TestCatConvertScreenToPNG(t *testing.T) {
	img := blankDiskWithScreen(t)
	args, err := docopt.Parse(usage("samfile"),
		[]string{"cat", "-i", img, "-f", "SCR", "-c"}, true, "samfile", false, true)
	if err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() { cat(args) })
	assertPNG256x192(t, out)
}

// TestExtractConvertScreenToPNG extracts a disk with -c and confirms
// the SCREEN$ file lands as SCR.png with valid PNG content.
func TestExtractConvertScreenToPNG(t *testing.T) {
	img := blankDiskWithScreen(t)
	target := t.TempDir()
	args, err := docopt.Parse(usage("samfile"),
		[]string{"extract", "-i", img, "-t", target, "-c"}, true, "samfile", false, true)
	if err != nil {
		t.Fatal(err)
	}
	extract(args)
	out, err := os.ReadFile(filepath.Join(target, "SCR.png"))
	if err != nil {
		t.Fatalf("expected SCR.png in extract target: %v", err)
	}
	assertPNG256x192(t, out)
}

// TestScreenToPNGSubcommand pipes a SCREEN$ body through screen-to-png
// and confirms a valid PNG comes out.
func TestScreenToPNGSubcommand(t *testing.T) {
	old := os.Stdin
	defer func() { os.Stdin = old }()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	go func() {
		w.Write(mode4Body())
		w.Close()
	}()
	args, err := docopt.Parse(usage("samfile"),
		[]string{"screen-to-png"}, true, "samfile", false, true)
	if err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() { screenToPNG(args) })
	assertPNG256x192(t, out)
}

func assertPNG256x192(t *testing.T, data []byte) {
	t.Helper()
	im, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("output is not a valid PNG: %v", err)
	}
	if im.Bounds() != image.Rect(0, 0, 256, 192) {
		t.Fatalf("PNG bounds = %v, want 256x192", im.Bounds())
	}
}
