package samfile

import "testing"

// buildBootableDisk builds a samfile-built disk where slot 0's first
// sector is (4, 1) — i.e. AddCodeFile's first allocation. Body bytes
// 256-259 are patched to "BOOT" and body byte 0 (sector offset 9) is
// patched to a real opcode (0xC3 = JP nn) so all three §11 rules pass
// on the positive case.
func buildBootableDisk(t *testing.T) (*DiskImage, *DiskJournal) {
	t.Helper()
	di := NewDiskImage()
	if err := di.AddCodeFile("samdos2", make([]byte, 400), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	first := di.DiskJournal()[0].FirstSector
	sd, err := di.SectorData(first)
	if err != nil {
		t.Fatalf("SectorData: %v", err)
	}
	// "BOOT" at sector offset 256-259.
	copy(sd[256:260], []byte{'B', 'O', 'O', 'T'})
	// Real opcode at sector offset 9 (body offset 0). 0xC3 = JP nn.
	sd[9] = 0xC3
	di.WriteSector(first, sd)
	return di, di.DiskJournal()
}

func TestBootOwnerAtT4S1Positive(t *testing.T) {
	di, _ := buildBootableDisk(t)
	findings := checkBootOwnerAtT4S1(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("disk with T4S1 owner: %d findings; want 0", len(findings))
	}
}

func TestBootOwnerAtT4S1Negative(t *testing.T) {
	// Empty disk → no used slot owns T4S1 → rule fires.
	di := NewDiskImage()
	findings := checkBootOwnerAtT4S1(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BOOT-OWNER-AT-T4S1" {
		t.Fatalf("got %d findings, first=%+v; want 1 BOOT-OWNER-AT-T4S1", len(findings), findings)
	}
	if findings[0].Severity != SeverityStructural {
		t.Errorf("Severity = %v; want structural", findings[0].Severity)
	}
}

func TestBootSignatureAt256Positive(t *testing.T) {
	di, _ := buildBootableDisk(t)
	findings := checkBootSignatureAt256(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("disk with BOOT signature: %d findings; want 0", len(findings))
	}
}

func TestBootSignatureAt256Negative(t *testing.T) {
	di, _ := buildBootableDisk(t)
	first := di.DiskJournal()[0].FirstSector
	sd, _ := di.SectorData(first)
	// Corrupt one byte of the signature.
	sd[257] = 'X'
	di.WriteSector(first, sd)
	findings := checkBootSignatureAt256(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BOOT-SIGNATURE-AT-256" {
		t.Fatalf("got %d findings, first=%+v; want 1 BOOT-SIGNATURE-AT-256", len(findings), findings)
	}
}

func TestBootSignatureAt256CaseInsensitive(t *testing.T) {
	// Lowercase "boot" must also match (ROM AND 0x5F mask).
	di, _ := buildBootableDisk(t)
	first := di.DiskJournal()[0].FirstSector
	sd, _ := di.SectorData(first)
	copy(sd[256:260], []byte{'b', 'o', 'o', 't'})
	di.WriteSector(first, sd)
	findings := checkBootSignatureAt256(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("lowercase 'boot' (ROM case-insensitive): %d findings; want 0", len(findings))
	}
}

func TestBootEntryPointAt9Positive(t *testing.T) {
	di, _ := buildBootableDisk(t)
	findings := checkBootEntryPointAt9(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("body[0] = 0xC3 (JP nn): %d findings; want 0", len(findings))
	}
}

func TestBootEntryPointAt9Negative(t *testing.T) {
	di, _ := buildBootableDisk(t)
	first := di.DiskJournal()[0].FirstSector
	sd, _ := di.SectorData(first)
	sd[9] = 0xFF // unwritten marker — implausible entry
	di.WriteSector(first, sd)
	findings := checkBootEntryPointAt9(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BOOT-ENTRY-POINT-AT-9" {
		t.Fatalf("got %d findings, first=%+v; want 1 BOOT-ENTRY-POINT-AT-9", len(findings), findings)
	}
}
