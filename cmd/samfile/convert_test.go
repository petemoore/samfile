package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	docopt "github.com/docopt/docopt-go"
	"github.com/petemoore/samfile/v3"
	"github.com/petemoore/samfile/v3/sambasic"
)

// blankDiskWithCode writes a fresh MGT image containing one CODE file
// ("CODE") with a trivial body, and returns its path.
func blankDiskWithCode(t *testing.T) string {
	t.Helper()
	di := samfile.NewDiskImage()
	// AddCodeFile requires load address >= 16384; body must be non-empty.
	body := bytes.Repeat([]byte{0x00}, 4)
	if err := di.AddCodeFile("CODE", body, 16384, 16384); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	path := filepath.Join(t.TempDir(), "code.mgt")
	if err := di.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path
}

// blankDiskWithBasic writes a fresh MGT image containing one SAM BASIC
// file ("PROG") with a single-line listing, and returns its path.
func blankDiskWithBasic(t *testing.T) string {
	t.Helper()
	f, err := sambasic.ParseTextString("10 REM hello\n")
	if err != nil {
		t.Fatalf("ParseTextString: %v", err)
	}
	di := samfile.NewDiskImage()
	if err := di.AddBasicFile("PROG", f); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}
	path := filepath.Join(t.TempDir(), "basic.mgt")
	if err := di.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path
}

// parseArgs is a convenience wrapper for docopt.Parse used in tests.
func parseArgs(t *testing.T, argv []string) map[string]any {
	t.Helper()
	args, err := docopt.Parse(usage("samfile"), argv, true, "samfile", false, true)
	if err != nil {
		t.Fatalf("docopt.Parse: %v", err)
	}
	return args
}

// ----- parseConvertOptions unit tests -----

func TestParseConvertOptions_Defaults(t *testing.T) {
	args := parseArgs(t, []string{"cat", "-i", "x.mgt", "-f", "F", "-c"})
	opts, err := parseConvertOptions(args, "x.mgt", []byte("dummy"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.audioFormat != "mp3" {
		t.Errorf("default audioFormat = %q, want %q", opts.audioFormat, "mp3")
	}
	if opts.loops != 4 {
		t.Errorf("default loops = %d, want 4", opts.loops)
	}
	if opts.diskName != "x.mgt" {
		t.Errorf("diskName = %q, want %q", opts.diskName, "x.mgt")
	}
	if opts.diskSHA256 == "" {
		t.Error("diskSHA256 should not be empty")
	}
}

func TestParseConvertOptions_ExplicitFormat(t *testing.T) {
	for _, fmt := range []string{"wav", "flac", "mp3"} {
		args := parseArgs(t, []string{"cat", "-i", "x.mgt", "-f", "F", "-c", "-a", fmt})
		opts, err := parseConvertOptions(args, "x.mgt", []byte("dummy"))
		if err != nil {
			t.Fatalf("format %q: unexpected error: %v", fmt, err)
		}
		if opts.audioFormat != fmt {
			t.Errorf("format %q: audioFormat = %q", fmt, opts.audioFormat)
		}
	}
}

func TestParseConvertOptions_InvalidFormat(t *testing.T) {
	for _, bad := range []string{"m4a", "ogg", "aac", "bogus"} {
		args := parseArgs(t, []string{"cat", "-i", "x.mgt", "-f", "F", "-c", "-a", bad})
		_, err := parseConvertOptions(args, "x.mgt", []byte("dummy"))
		if err == nil {
			t.Errorf("format %q: expected error, got nil", bad)
			continue
		}
		msg := err.Error()
		if !strings.Contains(msg, "wav") || !strings.Contains(msg, "flac") || !strings.Contains(msg, "mp3") {
			t.Errorf("format %q: error %q does not mention supported formats wav/flac/mp3", bad, msg)
		}
	}
}

func TestParseConvertOptions_LoopsExplicit(t *testing.T) {
	args := parseArgs(t, []string{"cat", "-i", "x.mgt", "-f", "F", "-c", "--loops", "8"})
	opts, err := parseConvertOptions(args, "x.mgt", []byte("dummy"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.loops != 8 {
		t.Errorf("loops = %d, want 8", opts.loops)
	}
}

func TestParseConvertOptions_LoopsZero(t *testing.T) {
	args := parseArgs(t, []string{"cat", "-i", "x.mgt", "-f", "F", "-c", "--loops", "0"})
	_, err := parseConvertOptions(args, "x.mgt", []byte("dummy"))
	if err == nil {
		t.Fatal("expected error for --loops 0, got nil")
	}
	if !strings.Contains(err.Error(), "positive") {
		t.Errorf("error %q does not mention 'positive'", err.Error())
	}
}

func TestParseConvertOptions_LoopsNegative(t *testing.T) {
	args := parseArgs(t, []string{"cat", "-i", "x.mgt", "-f", "F", "-c", "--loops", "-1"})
	_, err := parseConvertOptions(args, "x.mgt", []byte("dummy"))
	if err == nil {
		t.Fatal("expected error for --loops -1, got nil")
	}
}

func TestParseConvertOptions_LoopsNonNumeric(t *testing.T) {
	args := parseArgs(t, []string{"cat", "-i", "x.mgt", "-f", "F", "-c", "--loops", "abc"})
	_, err := parseConvertOptions(args, "x.mgt", []byte("dummy"))
	if err == nil {
		t.Fatal("expected error for --loops abc, got nil")
	}
}

func TestParseConvertOptions_SHA256IsPerDisk(t *testing.T) {
	// The same imageBytes must yield the same sha256 regardless of invocation.
	imageBytes := []byte("some disk content")
	args := parseArgs(t, []string{"cat", "-i", "disk.mgt", "-f", "F", "-c"})
	opts1, err := parseConvertOptions(args, "disk.mgt", imageBytes)
	if err != nil {
		t.Fatal(err)
	}
	opts2, err := parseConvertOptions(args, "disk.mgt", imageBytes)
	if err != nil {
		t.Fatal(err)
	}
	if opts1.diskSHA256 != opts2.diskSHA256 {
		t.Error("diskSHA256 is not deterministic across calls with same input")
	}
	// Different content must yield different sha256.
	opts3, err := parseConvertOptions(args, "disk.mgt", []byte("different"))
	if err != nil {
		t.Fatal(err)
	}
	if opts3.diskSHA256 == opts1.diskSHA256 {
		t.Error("diskSHA256 must differ for different imageBytes")
	}
}

// ----- Behavioural regression tests -----

// TestCatConvertCodeFileRaw verifies that cat -c on a plain CODE file
// produces raw bytes (no conversion), identical to cat without -c.
func TestCatConvertCodeFileRaw(t *testing.T) {
	img := blankDiskWithCode(t)

	catArgs := func(extra ...string) map[string]any {
		argv := append([]string{"cat", "-i", img, "-f", "CODE"}, extra...)
		return parseArgs(t, argv)
	}

	rawOut := captureStdout(t, func() { cat(catArgs()) })
	convertedOut := captureStdout(t, func() { cat(catArgs("-c")) })

	if !bytes.Equal(rawOut, convertedOut) {
		t.Errorf("cat -c on a CODE file should produce raw bytes, but output differs:\n  raw      %d bytes\n  converted %d bytes", len(rawOut), len(convertedOut))
	}
}

// TestCatConvertBasicToText verifies that cat -c on a BASIC file produces
// plain-text output (starts with a line number, not binary).
func TestCatConvertBasicToText(t *testing.T) {
	img := blankDiskWithBasic(t)

	args := parseArgs(t, []string{"cat", "-i", img, "-f", "PROG", "-c"})
	out := captureStdout(t, func() { cat(args) })

	text := string(out)
	// Detokenised output should contain the line number and keyword.
	if !strings.Contains(text, "10") || !strings.Contains(text, "REM") {
		t.Errorf("cat -c on BASIC file: expected detokenised text with '10 REM', got: %q", text)
	}
}

// TestExtractConvertCodeFileRaw verifies that extract -c on a CODE file
// extracts raw bytes (filename unchanged, no extra extension).
func TestExtractConvertCodeFileRaw(t *testing.T) {
	img := blankDiskWithCode(t)
	target := t.TempDir()

	args := parseArgs(t, []string{"extract", "-i", img, "-t", target, "-c"})
	extract(args)

	// The raw CODE file should be present with no extra extension.
	rawPath := filepath.Join(target, "CODE")
	if _, err := os.Stat(rawPath); err != nil {
		t.Fatalf("expected %q to exist after extract -c on CODE file: %v", rawPath, err)
	}
	// No .txt/.png/.mp3 variant should appear.
	entries, _ := os.ReadDir(target)
	for _, e := range entries {
		name := e.Name()
		if name != "CODE" {
			t.Errorf("unexpected file in extract target: %q (only 'CODE' expected)", name)
		}
	}
}
