package samfile

import "testing"

func TestCrossNoSectorOverlapPositive(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("A", make([]byte, 100), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile A: %v", err)
	}
	if err := di.AddCodeFile("B", make([]byte, 100), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile B: %v", err)
	}
	findings := checkCrossNoSectorOverlap(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("two distinct files: %d findings; want 0", len(findings))
	}
}

func TestCrossNoSectorOverlapNegative(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("A", make([]byte, 100), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile A: %v", err)
	}
	if err := di.AddCodeFile("B", make([]byte, 100), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile B: %v", err)
	}
	dj := di.DiskJournal()
	// Copy slot 0's map into slot 1 so they claim overlapping sectors.
	dj[1].SectorAddressMap = dj[0].SectorAddressMap
	di.WriteFileEntry(dj, 1)
	findings := checkCrossNoSectorOverlap(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) < 1 || findings[0].RuleID != "CROSS-NO-SECTOR-OVERLAP" {
		t.Fatalf("got %d findings, first=%+v; want at least one CROSS-NO-SECTOR-OVERLAP",
			len(findings), findings)
	}
}

func TestCrossNoDuplicateNamesPositive(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("A", make([]byte, 100), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile A: %v", err)
	}
	if err := di.AddCodeFile("B", make([]byte, 100), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile B: %v", err)
	}
	findings := checkCrossNoDuplicateNames(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("distinct names: %d findings; want 0", len(findings))
	}
}

func TestCrossNoDuplicateNamesNegative(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("A", make([]byte, 100), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile A: %v", err)
	}
	if err := di.AddCodeFile("B", make([]byte, 100), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile B: %v", err)
	}
	dj := di.DiskJournal()
	// Rename slot 1 to "A" so it duplicates slot 0.
	copy(dj[1].Name[:], "A         ")
	di.WriteFileEntry(dj, 1)
	findings := checkCrossNoDuplicateNames(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "CROSS-NO-DUPLICATE-NAMES" {
		t.Fatalf("got %d findings, first=%+v; want 1 CROSS-NO-DUPLICATE-NAMES",
			len(findings), findings)
	}
}

func TestCrossDirectoryAreaUnusedPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkCrossDirectoryAreaUnused(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestCrossDirectoryAreaUnusedNegative(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 600) // 2 sectors so there's a chain link
	first := di.DiskJournal()[0].FirstSector
	// Point sector 0's next-link at a directory-area sector (T2 S5).
	sd, _ := di.SectorData(first)
	raw := sd[:]
	raw[510] = 2 // track 2, in directory area
	raw[511] = 5
	di.WriteSector(first, sd)
	findings := checkCrossDirectoryAreaUnused(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) < 1 || findings[0].RuleID != "CROSS-DIRECTORY-AREA-UNUSED" {
		t.Fatalf("got %d findings, first=%+v; want at least one CROSS-DIRECTORY-AREA-UNUSED",
			len(findings), findings)
	}
}
