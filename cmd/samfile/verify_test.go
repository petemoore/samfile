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
//
// Phase 6: BOOT-SIGNATURE-AT-256 requires T4S1 bytes 256-259 to spell
// "BOOT" and BOOT-ENTRY-POINT-AT-9 requires a plausible Z80 opcode at
// sector offset 9, so we patch those bytes after AddCodeFile.
func TestVerifyCmdOnPopulatedDisk(t *testing.T) {
	di := samfile.NewDiskImage()
	if err := di.AddCodeFile("F", []byte("hello"), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	// Patch the boot sector so §11 rules pass: "BOOT" at bytes 256-259
	// and a real Z80 opcode (0xC3 = JP nn) at sector offset 9.
	first := di.DiskJournal()[0].FirstSector
	sd, err := di.SectorData(first)
	if err != nil {
		t.Fatalf("SectorData: %v", err)
	}
	copy(sd[256:260], []byte{'B', 'O', 'O', 'T'})
	sd[9] = 0xC3 // JP nn — a plausible boot-code entry opcode
	di.WriteSector(first, sd)

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
// DISK-NOT-EMPTY fires; output mentions the rule ID.
//
// Phase 6: BOOT-OWNER-AT-T4S1 (fatal) also fires on an empty disk
// because no used slot owns T4S1. runVerify therefore returns an
// error. The test checks the output for DISK-NOT-EMPTY but does not
// assert that error is nil.
func TestVerifyCmdOnEmptyDisk(t *testing.T) {
	di := samfile.NewDiskImage()
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "test.mgt")
	if err := di.Save(imgPath); err != nil {
		t.Fatalf("Save: %v", err)
	}

	stdout, _ := captureVerify(t, imgPath)
	if !strings.Contains(stdout, "DISK-NOT-EMPTY") {
		t.Errorf("expected 'DISK-NOT-EMPTY' in output; got:\n%s", stdout)
	}
}

// TestVerifyCmdOnFatalFinding exercises the fatal-findings exit
// path in runVerify. Phase 1's only rule is severity
// Inconsistency, so this test registers a synthetic fatal rule,
// then confirms runVerify returns a non-nil error mentioning
// "fatal".
//
// Caveat: the rule registry is package-private to samfile and
// exposes no Unregister. We can only call samfile.Register
// additively from this test package, and the registered rule
// will run for every later test in the same process. Two
// workarounds keep things sane:
//
//  1. alreadyRegistered guard — under `go test -count=N` the same
//     test runs N times in one process; without this guard the
//     second iteration would panic on duplicate-ID registration.
//
//  2. Marker-file gate — the rule's Check fires only when slot 0
//     of the disk under inspection contains a file whose name
//     starts with "FATAL-MARK". The sibling tests use plain "F",
//     so the rule is invisible to them even though it's in the
//     registry. This is a Phase-1 limitation; Phase 2+ may
//     revisit it when corpus-validation tests need fully-clean
//     registries (e.g. by adding samfile.Unregister or a
//     TestMain-scoped registry snapshot).
func TestVerifyCmdOnFatalFinding(t *testing.T) {
	alreadyRegistered := false
	for _, r := range samfile.Rules() {
		if r.ID == "TEST-FATAL" {
			alreadyRegistered = true
			break
		}
	}
	if !alreadyRegistered {
		samfile.Register(samfile.Rule{
			ID:          "TEST-FATAL",
			Severity:    samfile.SeverityFatal,
			Description: "always fatal (gated on FATAL-MARK marker file)",
			Citation:    "test",
			Check: func(ctx *samfile.CheckContext) []samfile.Finding {
				// Gate on a marker filename so this rule is
				// invisible to the sibling CLI tests that share
				// the same process registry.
				dj := ctx.Disk.DiskJournal()
				if dj[0] == nil {
					return nil
				}
				name := strings.TrimRight(string(dj[0].Name[:]), " ")
				if !strings.HasPrefix(name, "FATAL-MARK") {
					return nil
				}
				return []samfile.Finding{{
					RuleID:   "TEST-FATAL",
					Severity: samfile.SeverityFatal,
					Location: samfile.DiskWideLocation(),
					Message:  "synthetic fatal for testing",
					Citation: "test",
				}}
			},
		})
	}

	di := samfile.NewDiskImage()
	if err := di.AddCodeFile("FATAL-MARK", []byte("hi"), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "test.mgt")
	if err := di.Save(imgPath); err != nil {
		t.Fatalf("Save: %v", err)
	}

	_, err := captureVerify(t, imgPath)
	if err == nil {
		t.Fatal("runVerify returned nil error; expected non-nil for fatal finding")
	}
	if !strings.Contains(err.Error(), "verify:") {
		t.Errorf("error message = %q; want prefix 'verify:'", err.Error())
	}
	if !strings.Contains(err.Error(), "fatal") {
		t.Errorf("error message = %q; want to mention 'fatal'", err.Error())
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
