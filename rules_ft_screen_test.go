package samfile

import "testing"

func TestScreenModeAt0xDDPositive(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Type = FT_SCREEN
	dj[0].FileTypeInfo[0] = 2
	di.WriteFileEntry(dj, 0)
	findings := checkScreenModeAt0xDD(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("mode=2: %d findings; want 0", len(findings))
	}
}

func TestScreenModeAt0xDDNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Type = FT_SCREEN
	dj[0].FileTypeInfo[0] = 9 // out of range
	di.WriteFileEntry(dj, 0)
	findings := checkScreenModeAt0xDD(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "SCREEN-MODE-AT-0xDD" {
		t.Fatalf("got %d findings, first=%+v; want 1 SCREEN-MODE-AT-0xDD", len(findings), findings)
	}
}

func TestScreenLengthMatchesModePositive(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("SCREEN1", make([]byte, 6912), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].Type = FT_SCREEN
	dj[0].FileTypeInfo[0] = 1
	di.WriteFileEntry(dj, 0)
	findings := checkScreenLengthMatchesMode(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("mode 1 + 6912 bytes: %d findings; want 0", len(findings))
	}
}

func TestScreenLengthMatchesModeNegative(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("TEST", make([]byte, 100), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].Type = FT_SCREEN
	dj[0].FileTypeInfo[0] = 1 // expects 6912 bytes; body has 100
	di.WriteFileEntry(dj, 0)
	findings := checkScreenLengthMatchesMode(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "SCREEN-LENGTH-MATCHES-MODE" {
		t.Fatalf("got %d findings, first=%+v; want 1 SCREEN-LENGTH-MATCHES-MODE", len(findings), findings)
	}
}
