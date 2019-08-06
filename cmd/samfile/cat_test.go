package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	docopt "github.com/docopt/docopt-go"
)

func TestCatEnolaGayFileFromETrackerDisk(t *testing.T) {

	oldStdout := os.Stdout // keep backup of the real stdout
	defer func(f *os.File) {
		os.Stdout = f // restoring original stdout
	}(oldStdout)
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	h := sha256.New()
	d := make(chan (struct{}))
	go func() {
		io.Copy(h, r)
		close(d)
	}()

	imageFile, err := filepath.Abs(filepath.Join("..", "..", "testdata", "ETrackerv1.2.mgt"))
	if err != nil {
		log.Fatal(err)
	}
	samFile := "ENOLA_G .M"
	expectedSHA256 := "7ff534304084bf3638fb30e00b26105124e3e229981d43e5ec3b4c74cb527d06"

	command := []string{
		"cat",
		"-i",
		imageFile,
		"-f",
		samFile,
	}
	arguments, err := docopt.Parse(usage("samfile"), command, true, "samfile", false, true)
	if err != nil {
		t.Fatal(err)
	}
	cat(arguments)
	w.Close()
	<-d

	actualSHA256 := fmt.Sprintf("%x", h.Sum(nil))

	if actualSHA256 != expectedSHA256 {
		t.Fatalf("Extracted file %q from disk image %q and expected it to have SHA256 %q but it had SHA256 %q", samFile, imageFile, expectedSHA256, actualSHA256)
	}
}
