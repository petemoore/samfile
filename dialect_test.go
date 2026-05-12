package samfile

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
)

func TestDetectDialectEmptyDisk(t *testing.T) {
	di := NewDiskImage()
	if got := DetectDialect(di); got != DialectUnknown {
		t.Errorf("DetectDialect(empty) = %v; want unknown", got)
	}
}

func TestDetectDialectUnknownBootFileName(t *testing.T) {
	// A disk whose first file is named something neither DOS recognises
	// and whose MGTFlags are vanilla 0 (AddCodeFile leaves MGTFlags at
	// zero) emits no signal. DetectDialect must return Unknown rather
	// than guessing.
	di := NewDiskImage()
	if err := di.AddCodeFile("BOOTER", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	if got := DetectDialect(di); got != DialectUnknown {
		t.Errorf("DetectDialect(unknown boot file) = %v; want unknown", got)
	}
}

func TestDetectDialectSamdos2BootFile(t *testing.T) {
	// First file added → allocated at FirstSector (4, 1). Name is the
	// canonical samdos2 filename. The body content does not matter for
	// detection — only the slot name does.
	di := NewDiskImage()
	if err := di.AddCodeFile("samdos2", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	if got := DetectDialect(di); got != DialectSAMDOS2 {
		t.Errorf("DetectDialect(samdos2 boot file) = %v; want samdos2", got)
	}
}

func TestDetectDialectMasterDOSBootFile(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("masterdos2", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	if got := DetectDialect(di); got != DialectMasterDOS {
		t.Errorf("DetectDialect(masterdos2 boot file) = %v; want masterdos", got)
	}
}

func TestDetectDialectSAMDOS1ByName(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("samdos", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	if got := DetectDialect(di); got != DialectSAMDOS1 {
		t.Errorf("DetectDialect(samdos boot file) = %v; want samdos1", got)
	}
}

func TestDetectDialectSAMDOS1ByType3(t *testing.T) {
	// A bootstrap with an unrecognised filename but masked type 3 is
	// SAMDOS-1's auto-include header (samdos/src/b.s:14-22). Use
	// AddCodeFile, then patch Type to FT(3) via a journal write.
	di := NewDiskImage()
	if err := di.AddCodeFile("oddname", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].Type = FileType(3)
	di.WriteFileEntry(dj, 0)

	if got := DetectDialect(di); got != DialectSAMDOS1 {
		t.Errorf("DetectDialect(type-3 boot file) = %v; want samdos1", got)
	}
}

func TestDetectDialectMasterDOSByMGTFlags(t *testing.T) {
	// AddCodeFile leaves MGTFlags at 0x00 (vanilla SAMDOS-2 CODE
	// convention). Patch MGTFlags to 0x80 — an extended bit outside
	// {0x00, 0x20} — and DetectDialect must report MasterDOS.
	di := NewDiskImage()
	if err := di.AddCodeFile("data", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].MGTFlags = 0x80
	di.WriteFileEntry(dj, 0)

	if got := DetectDialect(di); got != DialectMasterDOS {
		t.Errorf("DetectDialect(MGTFlags=0x80) = %v; want masterdos", got)
	}
}

func TestMGTFlagsDialectVanillaIsSilent(t *testing.T) {
	// A disk where every used slot has MGTFlags in {0x00, 0x20, 0xFF}
	// — the SAMDOS-2 set — yields no opinion from mgtFlagsDialect.
	// Slot 0 keeps AddCodeFile's MGTFlags=0x00 default (samfile-built
	// CODE convention). Slot 1 is patched to 0x20 to stand in for a
	// BASIC file. Slot 2 is patched to 0xFF to stand in for a
	// real-SAMDOS-2 CODE file (the HDCLP2 0xFF-fill from rom-disasm
	// L22076-22080, observed on the M0 boot disk's slot-4 OUT entry).
	di := NewDiskImage()
	if err := di.AddCodeFile("CODE", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (CODE, MGTFlags=0): %v", err)
	}
	if err := di.AddCodeFile("BASIC", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (BASIC stub): %v", err)
	}
	if err := di.AddCodeFile("ROMSAVE", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (ROMSAVE stub): %v", err)
	}
	dj := di.DiskJournal()
	dj[1].MGTFlags = 0x20
	di.WriteFileEntry(dj, 1)
	dj[2].MGTFlags = 0xff
	di.WriteFileEntry(dj, 2)

	if got := mgtFlagsDialect(di.DiskJournal()); got != DialectUnknown {
		t.Errorf("mgtFlagsDialect(vanilla MGTFlags) = %v; want unknown", got)
	}
}

func TestDetectDialectSamdos2WithRomSaveMGTFlags(t *testing.T) {
	// Regression for the M0 boot disk scenario: a SAMDOS-2 boot file
	// at T4S1 plus a real-ROM-SAVE CODE file with MGTFlags=0xFF must
	// detect as SAMDOS-2, not Unknown. Before the SAMDOS-2 set was
	// widened to include 0xFF, mgtFlagsDialect mistook the HDCLP2
	// 0xFF-fill (rom-disasm L22076-22080) for a MasterDOS attribute
	// bit and collapsed the report to Unknown.
	di := NewDiskImage()
	if err := di.AddCodeFile("samdos2", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (samdos2 boot): %v", err)
	}
	if err := di.AddCodeFile("OUT", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (OUT): %v", err)
	}
	dj := di.DiskJournal()
	dj[1].MGTFlags = 0xff
	di.WriteFileEntry(dj, 1)

	if got := DetectDialect(di); got != DialectSAMDOS2 {
		t.Errorf("DetectDialect(samdos2 + 0xFF MGTFlags) = %v; want samdos2", got)
	}
}

func TestDetectDialectMasterDOSBothSignalsAgree(t *testing.T) {
	// Boot file "masterdos2" + extended MGTFlags on a second slot —
	// two signals both point at MasterDOS.
	di := NewDiskImage()
	if err := di.AddCodeFile("masterdos2", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (boot): %v", err)
	}
	if err := di.AddCodeFile("payload", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (payload): %v", err)
	}
	dj := di.DiskJournal()
	dj[1].MGTFlags = 0x40
	di.WriteFileEntry(dj, 1)

	if got := DetectDialect(di); got != DialectMasterDOS {
		t.Errorf("DetectDialect(both signals masterdos) = %v; want masterdos", got)
	}
}

func TestDetectDialectConflictReturnsUnknown(t *testing.T) {
	// Boot file says SAMDOS-2 but a later slot's MGTFlags say
	// MasterDOS. DetectDialect must collapse to Unknown rather than
	// pick a winner.
	di := NewDiskImage()
	if err := di.AddCodeFile("samdos2", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (boot): %v", err)
	}
	if err := di.AddCodeFile("payload", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (payload): %v", err)
	}
	dj := di.DiskJournal()
	dj[1].MGTFlags = 0x80
	di.WriteFileEntry(dj, 1)

	if got := DetectDialect(di); got != DialectUnknown {
		t.Errorf("DetectDialect(conflict samdos2 vs masterdos) = %v; want unknown", got)
	}
}

// withKnownBodySHA registers (sha→dialect) in the package-level map
// for the duration of a test. The deferred cleanup runs even when
// the test fails, so parallel tests stay isolated.
func withKnownBodySHA(t *testing.T, sha string, d Dialect) {
	t.Helper()
	if _, exists := knownBodySHAs[sha]; exists {
		t.Fatalf("test fixture conflict: SHA %s already in knownBodySHAs", sha)
	}
	knownBodySHAs[sha] = d
	t.Cleanup(func() { delete(knownBodySHAs, sha) })
}

// bootableBodyFixture returns a synthetic body whose first 510
// bytes-after-the-9-byte-header place the BOOT signature at T4S1
// disk offsets 256..259 (the ROM BTCK check, rom-disasm:20473-20598).
// Body bytes 247..250 == "BOOT" because AddCodeFile writes
// [9-byte header][body] starting at T4S1 offset 0, so body[247]
// lands at disk offset 9 + 247 = 256. The rest of the body is
// filled with the supplied marker byte so each call produces a
// different sha256.
func bootableBodyFixture(marker byte) []byte {
	body := make([]byte, 600)
	for i := range body {
		body[i] = marker
	}
	copy(body[247:251], []byte{'B', 'O', 'O', 'T'})
	return body
}

func TestBodyShaDialectMatchOverridesEverything(t *testing.T) {
	// Regression for the Fredatives 3 misclassification:
	// - Slot 0 body bytes are byte-identical to samdos2.reference.bin
	// - Slot 0 dir filename is *not* "samdos2" (it's "OS" on that
	//   particular disk)
	// - Other slots have MGTFlags values that mgtFlagsDialect reads
	//   as MasterDOS (bits outside {0x00, 0x20, 0xFF})
	//
	// Before bodyShaDialect, DetectDialect collapsed to MasterDOS.
	// With the body-SHA signal in place, the registered samdos2 SHA
	// is recognised on slot 0 and DetectDialect returns SAMDOS2,
	// ignoring the misleading dir-entry name and the noisy
	// neighbouring MGTFlags.
	body := bootableBodyFixture(0xAA)
	sum := sha256.Sum256(body)
	withKnownBodySHA(t, hex.EncodeToString(sum[:]), DialectSAMDOS2)

	di := NewDiskImage()
	if err := di.AddCodeFile("OS", body, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (boot stand-in): %v", err)
	}
	if err := di.AddCodeFile("noise", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (noise): %v", err)
	}
	dj := di.DiskJournal()
	// Push slot 1's MGTFlags outside the SAMDOS-2 set; without
	// bodyShaDialect this would push DetectDialect to MasterDOS.
	dj[1].MGTFlags = 0x24
	di.WriteFileEntry(dj, 1)

	if got := DetectDialect(di); got != DialectSAMDOS2 {
		t.Errorf("DetectDialect(samdos2 body + non-canonical name + masterdos MGTFlags) = %v; want samdos2", got)
	}
}

func TestBodyShaDialectMatchesWithoutDirEntry(t *testing.T) {
	// Regression for FRED Magazine Issue 13 (1991): SAMDOS-2 body is
	// present at T4S1 but the directory has no slot pointing to it
	// (the BOOT-OWNER-AT-T4S1 rule fires on these). bodyShaDialect
	// must walk T4S1 directly and still match the known SHA so the
	// per-disk dialect comes out right.
	body := bootableBodyFixture(0xBB)
	sum := sha256.Sum256(body)
	withKnownBodySHA(t, hex.EncodeToString(sum[:]), DialectSAMDOS2)

	di := NewDiskImage()
	if err := di.AddCodeFile("OS", body, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	// Erase the dir entry but leave the data sectors intact —
	// simulates a disk whose dir was rewritten without scrubbing T4S1.
	dj := di.DiskJournal()
	dj[0] = FileEntryFrom([256]byte{})
	di.WriteFileEntry(dj, 0)

	if got := DetectDialect(di); got != DialectSAMDOS2 {
		t.Errorf("DetectDialect(SAMDOS-2 body at T4S1, no dir entry) = %v; want samdos2", got)
	}
}

func TestBodyShaDialectAbstainsFallsBackToHeuristics(t *testing.T) {
	// When the slot-0 body sha256 is not registered in
	// knownBodySHAs, bodyShaDialect returns Unknown and DetectDialect
	// falls back to bootFileDialect + mgtFlagsDialect, which means
	// the existing heuristic tests stay correct. This test exercises
	// the same shape as TestDetectDialectSamdos2BootFile but
	// explicitly proves no body-SHA entry is needed for it.
	di := NewDiskImage()
	if err := di.AddCodeFile("samdos2", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	// The 1-byte body hashes to a SHA we don't register; bodyShaDialect
	// abstains and bootFileDialect picks SAMDOS-2 from the filename.
	if got := DetectDialect(di); got != DialectSAMDOS2 {
		t.Errorf("DetectDialect(unknown SHA, samdos2 filename) = %v; want samdos2", got)
	}
}

func TestBodyShaDialectIgnoresEmptyDisk(t *testing.T) {
	// An empty disk has no slot owning T4S1, so bodyShaDialect must
	// return Unknown rather than panicking on the missing FirstSector.
	di := NewDiskImage()
	if got := bodyShaDialect(di); got != DialectUnknown {
		t.Errorf("bodyShaDialect(empty disk) = %v; want unknown", got)
	}
}

func TestDetectDialectETrackerCorpus(t *testing.T) {
	// Smoke test against a real-world MGT image. We do not assert a
	// specific dialect — we just assert DetectDialect returns one of
	// the four documented values without panicking. This protects
	// against nil-pointer paths in bootFileDialect / mgtFlagsDialect
	// that fabricated disks might not exercise.
	const path = "testdata/ETrackerv1.2.mgt"
	if _, err := os.Stat(path); err != nil {
		t.Skipf("corpus image not present (%v); skipping", err)
	}
	di, err := Load(path)
	if err != nil {
		t.Fatalf("Load(%q): %v", path, err)
	}
	got := DetectDialect(di)
	switch got {
	case DialectUnknown, DialectSAMDOS1, DialectSAMDOS2, DialectMasterDOS:
		// All four are acceptable; log for diagnostic value.
		t.Logf("DetectDialect(%s) = %s", path, got)
	default:
		t.Errorf("DetectDialect(%s) = %v; not a documented Dialect value", path, got)
	}
}
