package samfile

import "testing"

func TestDirTypeByteIsKnownPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDirTypeByteIsKnown(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

// (Previously: TestDirTypeByteIsKnownNegative used FileType(0) to
// trigger the rule. Iteration 1 FIX explicitly skips Type==0
// because DIR-ERASED-IS-ZERO handles that condition with
// structural severity. Under the new behaviour the rule is
// effectively unreachable via the in-memory builder: every
// FileType that passes FileEntry.Used() — i.e. doesn't render as
// `UNKNOWN (N)` via String() — has its low-5 bits already in
// dirKnownTypes. The rule fires on raw corpus disks whose type
// byte was hand-written. TestDirTypeByteIsKnownSkipsErased below
// covers the FIX itself.)

func TestDirTypeByteIsKnownSkipsErased(t *testing.T) {
	// Iteration 1 FIX regression test: type byte 0 is the erased-slot
	// sentinel and is handled by DIR-ERASED-IS-ZERO at structural
	// severity. DIR-TYPE-BYTE-IS-KNOWN must not also fire on it
	// (previously double-fired across 2,492 corpus findings).
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Type = FileType(0)
	di.WriteFileEntry(dj, 0)
	findings := checkDirTypeByteIsKnown(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("Type=0 (erased sentinel): %d findings; want 0 (DIR-ERASED-IS-ZERO handles this)", len(findings))
	}
}

func TestDirErasedIsZeroPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDirErasedIsZero(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDirErasedIsZeroNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Type = FileType(0)
	di.WriteFileEntry(dj, 0)
	findings := checkDirErasedIsZero(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "DIR-ERASED-IS-ZERO" {
		t.Fatalf("got %d findings, first=%+v; want 1 DIR-ERASED-IS-ZERO", len(findings), findings)
	}
}

func TestDirNamePaddingPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDirNamePadding(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDirNamePaddingNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Name = Filename{'A', 0x01, 'B', ' ', ' ', ' ', ' ', ' ', ' ', ' '} // 0x01 control char
	di.WriteFileEntry(dj, 0)
	findings := checkDirNamePadding(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "DIR-NAME-PADDING" {
		t.Fatalf("got %d findings, first=%+v; want 1 DIR-NAME-PADDING", len(findings), findings)
	}
}

func TestDirNameNotEmptyPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDirNameNotEmpty(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDirNameNotEmptyNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Name = Filename{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '} // all spaces
	di.WriteFileEntry(dj, 0)
	findings := checkDirNameNotEmpty(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "DIR-NAME-NOT-EMPTY" {
		t.Fatalf("got %d findings, first=%+v; want 1 DIR-NAME-NOT-EMPTY", len(findings), findings)
	}
}

func TestDirFirstSectorValidPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDirFirstSectorValid(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDirFirstSectorValidNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].FirstSector.Sector = 99
	di.WriteFileEntry(dj, 0)
	findings := checkDirFirstSectorValid(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "DIR-FIRST-SECTOR-VALID" {
		t.Fatalf("got %d findings, first=%+v; want 1 DIR-FIRST-SECTOR-VALID", len(findings), findings)
	}
}

func TestDirSectorsMatchesChainPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDirSectorsMatchesChain(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDirSectorsMatchesChainNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Sectors = 99 // real chain is shorter
	di.WriteFileEntry(dj, 0)
	findings := checkDirSectorsMatchesChain(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "DIR-SECTORS-MATCHES-CHAIN" {
		t.Fatalf("got %d findings, first=%+v; want 1 DIR-SECTORS-MATCHES-CHAIN", len(findings), findings)
	}
}

func TestDirSectorsMatchesMapPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDirSectorsMatchesMap(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDirSectorsMatchesMapNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Sectors = 99 // map popcount is real allocation
	di.WriteFileEntry(dj, 0)
	findings := checkDirSectorsMatchesMap(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "DIR-SECTORS-MATCHES-MAP" {
		t.Fatalf("got %d findings, first=%+v; want 1 DIR-SECTORS-MATCHES-MAP", len(findings), findings)
	}
}

func TestDirSectorsNonzeroPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDirSectorsNonzero(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDirSectorsNonzeroNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Sectors = 0
	di.WriteFileEntry(dj, 0)
	findings := checkDirSectorsNonzero(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "DIR-SECTORS-NONZERO" {
		t.Fatalf("got %d findings, first=%+v; want 1 DIR-SECTORS-NONZERO", len(findings), findings)
	}
}

func TestDirSAMWithinCapacityPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDirSAMWithinCapacity(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDirSAMWithinCapacityNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].SectorAddressMap[194] = 0xE0 // set top 3 bits
	di.WriteFileEntry(dj, 0)
	findings := checkDirSAMWithinCapacity(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "DIR-SAM-WITHIN-CAPACITY" {
		t.Fatalf("got %d findings, first=%+v; want 1 DIR-SAM-WITHIN-CAPACITY", len(findings), findings)
	}
}
