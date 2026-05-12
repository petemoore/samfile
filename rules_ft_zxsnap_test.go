package samfile

import "testing"

// buildZXSnapDisk returns a samfile-built disk where slot 0 is a
// 49152-byte file morphed into FT_ZX_SNAPSHOT with start address
// 0x4000. AddCodeFile load 0x4000 sets fe.StartAddressPage so that
// StartAddress() decodes to 0x4000.
func buildZXSnapDisk(t *testing.T) (*DiskImage, *DiskJournal) {
	t.Helper()
	di := NewDiskImage()
	if err := di.AddCodeFile("ZXSNAP", make([]byte, 49152), 0x4000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].Type = FT_ZX_SNAPSHOT
	di.WriteFileEntry(dj, 0)
	return di, di.DiskJournal()
}

func TestZXSnapLength49152Positive(t *testing.T) {
	di, _ := buildZXSnapDisk(t)
	findings := checkZXSnapLength49152(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("49152-byte ZX snapshot: %d findings; want 0", len(findings))
	}
}

func TestZXSnapLength49152Negative(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("ZXSNAP", make([]byte, 100), 0x4000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].Type = FT_ZX_SNAPSHOT
	di.WriteFileEntry(dj, 0)
	findings := checkZXSnapLength49152(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "ZXSNAP-LENGTH-49152" {
		t.Fatalf("got %d findings, first=%+v; want 1 ZXSNAP-LENGTH-49152", len(findings), findings)
	}
}

func TestZXSnapLoadAddr16384Positive(t *testing.T) {
	di, _ := buildZXSnapDisk(t)
	findings := checkZXSnapLoadAddr16384(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("ZX snapshot at 0x4000: %d findings; want 0", len(findings))
	}
}

func TestZXSnapLoadAddr16384Negative(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("ZXSNAP", make([]byte, 49152), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].Type = FT_ZX_SNAPSHOT
	di.WriteFileEntry(dj, 0)
	findings := checkZXSnapLoadAddr16384(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "ZXSNAP-LOAD-ADDR-16384" {
		t.Fatalf("got %d findings, first=%+v; want 1 ZXSNAP-LOAD-ADDR-16384", len(findings), findings)
	}
}
