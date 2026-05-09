package samfile

import (
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
