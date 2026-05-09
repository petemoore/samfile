package samfile

import (
	"testing"
)

func TestSeverityOrdering(t *testing.T) {
	if !(SeverityCosmetic < SeverityInconsistency &&
		SeverityInconsistency < SeverityStructural &&
		SeverityStructural < SeverityFatal) {
		t.Fatalf("Severity constants out of order: cosmetic=%d inconsistency=%d structural=%d fatal=%d",
			SeverityCosmetic, SeverityInconsistency, SeverityStructural, SeverityFatal)
	}
}

func TestSeverityString(t *testing.T) {
	cases := []struct {
		sev  Severity
		want string
	}{
		{SeverityCosmetic, "cosmetic"},
		{SeverityInconsistency, "inconsistency"},
		{SeverityStructural, "structural"},
		{SeverityFatal, "fatal"},
	}
	for _, c := range cases {
		if got := c.sev.String(); got != c.want {
			t.Errorf("Severity(%d).String() = %q; want %q", c.sev, got, c.want)
		}
	}
}

func TestDialectString(t *testing.T) {
	cases := []struct {
		d    Dialect
		want string
	}{
		{DialectUnknown, "unknown"},
		{DialectSAMDOS1, "samdos1"},
		{DialectSAMDOS2, "samdos2"},
		{DialectMasterDOS, "masterdos"},
	}
	for _, c := range cases {
		if got := c.d.String(); got != c.want {
			t.Errorf("Dialect(%d).String() = %q; want %q", c.d, got, c.want)
		}
	}
}

func TestLocationDiskWide(t *testing.T) {
	loc := DiskWideLocation()
	if !loc.IsDiskWide() {
		t.Errorf("DiskWideLocation().IsDiskWide() = false; want true")
	}
	if loc.Slot != -1 || loc.Sector != nil || loc.ByteOffset != -1 || loc.Filename != "" {
		t.Errorf("DiskWideLocation() should leave all fields unset; got %+v", loc)
	}
}

func TestLocationSlot(t *testing.T) {
	loc := SlotLocation(3, "IN")
	if loc.IsDiskWide() {
		t.Errorf("SlotLocation(3).IsDiskWide() = true; want false")
	}
	if loc.Slot != 3 {
		t.Errorf("Slot = %d; want 3", loc.Slot)
	}
	if loc.Filename != "IN" {
		t.Errorf("Filename = %q; want %q", loc.Filename, "IN")
	}
	if loc.Sector != nil || loc.ByteOffset != -1 {
		t.Errorf("SlotLocation should leave sector + byte unset; got %+v", loc)
	}
}

func TestLocationSector(t *testing.T) {
	sec := &Sector{Track: 6, Sector: 3}
	loc := SectorLocation(2, "stub", sec, 8)
	if loc.IsDiskWide() {
		t.Errorf("SectorLocation.IsDiskWide() = true; want false")
	}
	if loc.Slot != 2 || loc.Filename != "stub" || loc.Sector != sec || loc.ByteOffset != 8 {
		t.Errorf("SectorLocation fields wrong; got %+v", loc)
	}
}
