package samfile

import "testing"

// Helper for §1 tests: a clean single-file disk with no chain-link
// anomalies. Returns the journal so tests can patch it.
func cleanSingleFileDisk(t *testing.T, name string, dataLen int) (*DiskImage, *DiskJournal) {
	t.Helper()
	di := NewDiskImage()
	data := make([]byte, dataLen)
	if err := di.AddCodeFile(name, data, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile(%q, len=%d): %v", name, dataLen, err)
	}
	return di, di.DiskJournal()
}

func TestDiskDirectoryTracksPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDiskDirectoryTracks(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDiskDirectoryTracksNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	// Patch FirstSector.Track to 2 (in the directory area).
	dj[0].FirstSector.Track = 2
	di.WriteFileEntry(dj, 0)
	findings := checkDiskDirectoryTracks(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) < 1 {
		t.Fatalf("got %d findings; want >= 1", len(findings))
	}
	if findings[0].RuleID != "DISK-DIRECTORY-TRACKS" || findings[0].Severity != SeverityStructural {
		t.Errorf("findings[0] = %+v", findings[0])
	}
}

func TestDiskTrackSideEncodingPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDiskTrackSideEncoding(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDiskTrackSideEncodingNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].FirstSector.Track = 0x60 // in the invalid 0x50-0x7F range
	di.WriteFileEntry(dj, 0)
	findings := checkDiskTrackSideEncoding(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) < 1 || findings[0].RuleID != "DISK-TRACK-SIDE-ENCODING" {
		t.Fatalf("got %d findings, first=%+v; want at least one DISK-TRACK-SIDE-ENCODING",
			len(findings), findings)
	}
	if findings[0].Severity != SeverityFatal {
		t.Errorf("Severity = %v; want fatal", findings[0].Severity)
	}
}

func TestDiskSectorRangePositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkDiskSectorRange(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: got %d findings; want 0", len(findings))
	}
}

func TestDiskSectorRangeNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].FirstSector.Sector = 11 // out of range
	di.WriteFileEntry(dj, 0)
	findings := checkDiskSectorRange(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) < 1 || findings[0].RuleID != "DISK-SECTOR-RANGE" {
		t.Fatalf("got %d findings, first=%+v; want at least one DISK-SECTOR-RANGE",
			len(findings), findings)
	}
	if findings[0].Severity != SeverityFatal {
		t.Errorf("Severity = %v; want fatal", findings[0].Severity)
	}
}
