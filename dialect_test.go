package samfile

import (
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
	// A disk where every used slot has MGTFlags in {0x00, 0x20}
	// (vanilla SAMDOS-2) yields no opinion from mgtFlagsDialect.
	di := NewDiskImage()
	if err := di.AddCodeFile("CODE", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (CODE, MGTFlags=0): %v", err)
	}
	// AddBasicFile sets MGTFlags=0x20 — exercise both bytes of the
	// SAMDOS-2 set.
	// (We do not need a real BASIC body; patch the second slot's
	// MGTFlags directly to keep the test minimal.)
	if err := di.AddCodeFile("BASIC", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (BASIC stub): %v", err)
	}
	dj := di.DiskJournal()
	dj[1].MGTFlags = 0x20
	di.WriteFileEntry(dj, 1)

	if got := mgtFlagsDialect(di.DiskJournal()); got != DialectUnknown {
		t.Errorf("mgtFlagsDialect(vanilla MGTFlags) = %v; want unknown", got)
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
