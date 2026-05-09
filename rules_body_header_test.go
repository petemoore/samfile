package samfile

import "testing"

// mutateFirstSectorByte patches one byte of slot 0's first sector
// payload (e.g. body header bytes). It's a small utility for the
// body-header tests' negative cases; raw byte-level mutation is
// the only way to disturb the body header without re-running the
// whole AddCodeFile path (which would re-mirror to the dir entry).
func mutateFirstSectorByte(t *testing.T, di *DiskImage, byteOffset int, newValue byte) {
	t.Helper()
	fe := di.DiskJournal()[0]
	sd, err := di.SectorData(fe.FirstSector)
	if err != nil {
		t.Fatalf("SectorData: %v", err)
	}
	sd[byteOffset] = newValue
	di.WriteSector(fe.FirstSector, sd)
}

// ----- BODY-TYPE-MATCHES-DIR -----

func TestBodyTypeMatchesDirPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkBodyTypeMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestBodyTypeMatchesDirNegative(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	// Patch body byte 0 (Type) to a value the dir doesn't reflect.
	mutateFirstSectorByte(t, di, 0, 0x05) // body says ZX_SNAPSHOT, dir says CODE
	findings := checkBodyTypeMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "BODY-TYPE-MATCHES-DIR" {
		t.Fatalf("got %d findings, first=%+v; want 1 BODY-TYPE-MATCHES-DIR", len(findings), findings)
	}
}

// ----- BODY-EXEC-DIV16K-MATCHES-DIR -----

func TestBodyExecDiv16KMatchesDirPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkBodyExecDiv16KMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestBodyExecDiv16KMatchesDirNegative(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	// Dir's ExecutionAddressDiv16K is 0xFF for AddCodeFile(..., 0) (no auto-exec); 0x7E differs.
	mutateFirstSectorByte(t, di, 5, 0x7E)
	findings := checkBodyExecDiv16KMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "BODY-EXEC-DIV16K-MATCHES-DIR" {
		t.Fatalf("got %d findings, first=%+v; want 1 BODY-EXEC-DIV16K-MATCHES-DIR", len(findings), findings)
	}
}

// ----- BODY-EXEC-MOD16K-LO-MATCHES-DIR -----

func TestBodyExecMod16KLoMatchesDirPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkBodyExecMod16KLoMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestBodyExecMod16KLoMatchesDirNegative(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	// Dir's ExecutionAddressMod16K low byte is 0xFF for no-exec; 0xAA differs.
	mutateFirstSectorByte(t, di, 6, 0xAA)
	findings := checkBodyExecMod16KLoMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "BODY-EXEC-MOD16K-LO-MATCHES-DIR" {
		t.Fatalf("got %d findings, first=%+v; want 1 BODY-EXEC-MOD16K-LO-MATCHES-DIR", len(findings), findings)
	}
}

// ----- BODY-PAGES-MATCHES-DIR -----

func TestBodyPagesMatchesDirPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkBodyPagesMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestBodyPagesMatchesDirNegative(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	// Dir's Pages for a 100-byte CODE file is 0; 0x99 differs.
	mutateFirstSectorByte(t, di, 7, 0x99)
	findings := checkBodyPagesMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "BODY-PAGES-MATCHES-DIR" {
		t.Fatalf("got %d findings, first=%+v; want 1 BODY-PAGES-MATCHES-DIR", len(findings), findings)
	}
}

// ----- BODY-STARTPAGE-MATCHES-DIR -----

func TestBodyStartPageMatchesDirPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkBodyStartPageMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestBodyStartPageMatchesDirNegative(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	// Dir's StartAddressPage for load 0x8000 is 1; 0x99 differs.
	mutateFirstSectorByte(t, di, 8, 0x99)
	findings := checkBodyStartPageMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "BODY-STARTPAGE-MATCHES-DIR" {
		t.Fatalf("got %d findings, first=%+v; want 1 BODY-STARTPAGE-MATCHES-DIR", len(findings), findings)
	}
}

// ----- BODY-LENGTHMOD16K-MATCHES-DIR -----

func TestBodyLengthMod16KMatchesDirPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkBodyLengthMod16KMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestBodyLengthMod16KMatchesDirNegative(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	// Patching byte 1 alone disagrees with the dir's parsed 16-bit LengthMod16K.
	mutateFirstSectorByte(t, di, 1, 0xAA)
	findings := checkBodyLengthMod16KMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "BODY-LENGTHMOD16K-MATCHES-DIR" {
		t.Fatalf("got %d findings, first=%+v; want 1 BODY-LENGTHMOD16K-MATCHES-DIR", len(findings), findings)
	}
}

// ----- BODY-PAGEOFFSET-MATCHES-DIR -----

func TestBodyPageOffsetMatchesDirPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkBodyPageOffsetMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestBodyPageOffsetMatchesDirNegative(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	// Same shape as LengthMod16K, byte 3.
	mutateFirstSectorByte(t, di, 3, 0xAA)
	findings := checkBodyPageOffsetMatchesDir(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "BODY-PAGEOFFSET-MATCHES-DIR" {
		t.Fatalf("got %d findings, first=%+v; want 1 BODY-PAGEOFFSET-MATCHES-DIR", len(findings), findings)
	}
}

// ----- BODY-MIRROR-AT-DIR-D3-DB -----

func TestBodyMirrorAtDirD3DBPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkBodyMirrorAtDirD3DB(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestBodyMirrorAtDirD3DBNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].MGTFutureAndPast[0] = 0xFF
	di.WriteFileEntry(dj, 0)
	findings := checkBodyMirrorAtDirD3DB(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "BODY-MIRROR-AT-DIR-D3-DB" {
		t.Fatalf("got %d findings, first=%+v; want 1 BODY-MIRROR-AT-DIR-D3-DB", len(findings), findings)
	}
}
