package samfile

import "testing"

func TestCosmeticReservedAFFPositive(t *testing.T) {
	// samfile's AddCodeFile leaves ReservedA at 0x00 × 4 — accepted.
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkCosmeticReservedAFF(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("samfile-built (ReservedA=0x00×4): %d findings; want 0", len(findings))
	}
}

func TestCosmeticReservedAFFAcceptsAllFF(t *testing.T) {
	// Real ROM SAMDOS-2 SAVE 0xFF-fills these bytes via HDCLP2.
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	for i := range dj[0].ReservedA {
		dj[0].ReservedA[i] = 0xFF
	}
	di.WriteFileEntry(dj, 0)
	findings := checkCosmeticReservedAFF(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("ReservedA=0xFF×4 (ROM SAVE convention): %d findings; want 0", len(findings))
	}
}

func TestCosmeticReservedAFFNegative(t *testing.T) {
	// A byte that's neither 0x00 nor 0xFF fires the rule.
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].ReservedA[2] = 0xAA
	di.WriteFileEntry(dj, 0)
	findings := checkCosmeticReservedAFF(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "COSMETIC-RESERVEDA-FF" {
		t.Fatalf("got %d findings, first=%+v; want 1 COSMETIC-RESERVEDA-FF", len(findings), findings)
	}
}
