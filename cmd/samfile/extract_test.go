package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	docopt "github.com/docopt/docopt-go"
)

func TestExtractEnolaGayFileFromETrackerDisk(t *testing.T) {

	imageFile, err := filepath.Abs(filepath.Join("..", "..", "testdata", "ETrackerv1.2.mgt"))
	if err != nil {
		log.Fatal(err)
	}

	samFile := "ENOLA_G .M"
	expectedSHA256 := "7ff534304084bf3638fb30e00b26105124e3e229981d43e5ec3b4c74cb527d06"

	dir, err := ioutil.TempDir("", "samfile-TestExtract")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	localFile, err := filepath.Abs(filepath.Join(dir, "Enola Gay"))
	if err != nil {
		log.Fatal(err)
	}

	command := []string{
		"extract",
		"--dest",
		localFile,
		imageFile,
		samFile,
	}
	arguments, err := docopt.Parse(usage("samfile"), command, true, "samfile", false, true)
	if err != nil {
		t.Fatal(err)
	}
	extract(arguments)
	f, err := os.Open(localFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		t.Fatal(err)
	}

	actualSHA256 := fmt.Sprintf("%x", h.Sum(nil))

	if actualSHA256 != expectedSHA256 {
		t.Fatalf("Extracted file %q from disk image %q to %q and expected it to have SHA256 %q but it had SHA256 %q", samFile, imageFile, localFile, expectedSHA256, actualSHA256)
	}
}
