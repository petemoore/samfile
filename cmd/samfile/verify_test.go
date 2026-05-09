package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/petemoore/samfile/v3"
)

// TestVerifyCmdOnPopulatedDisk runs the CLI subcommand against a
// disk built in-memory with one CODE file. Expected: DISK-NOT-EMPTY
// does not fire; output reports no findings; exit code (returned
// as nil error) is clean.
func TestVerifyCmdOnPopulatedDisk(t *testing.T) {
	di := samfile.NewDiskImage()
	if err := di.AddCodeFile("F", []byte("hello"), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "test.mgt")
	if err := di.Save(imgPath); err != nil {
		t.Fatalf("Save: %v", err)
	}

	stdout, err := captureVerify(t, imgPath)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if !strings.Contains(stdout, "no findings") {
		t.Errorf("expected 'no findings' in output; got:\n%s", stdout)
	}
}

// TestVerifyCmdOnEmptyDisk runs against an empty disk. Expected:
// DISK-NOT-EMPTY fires; output mentions the rule ID; error is nil
// (inconsistency does not gate exit code).
func TestVerifyCmdOnEmptyDisk(t *testing.T) {
	di := samfile.NewDiskImage()
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "test.mgt")
	if err := di.Save(imgPath); err != nil {
		t.Fatalf("Save: %v", err)
	}

	stdout, err := captureVerify(t, imgPath)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if !strings.Contains(stdout, "DISK-NOT-EMPTY") {
		t.Errorf("expected 'DISK-NOT-EMPTY' in output; got:\n%s", stdout)
	}
}

// captureVerify invokes the verify subcommand and returns its stdout.
func captureVerify(t *testing.T, imgPath string) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	err := runVerify(imgPath)

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String(), err
}
