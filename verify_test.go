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

func TestFindingShape(t *testing.T) {
	f := Finding{
		RuleID:   "TEST-RULE",
		Severity: SeverityStructural,
		Location: SlotLocation(2, "stub"),
		Message:  "expected X, got Y",
		Citation: "samdos/src/c.s:1306-1343",
	}
	if f.RuleID != "TEST-RULE" {
		t.Errorf("RuleID = %q; want TEST-RULE", f.RuleID)
	}
	if f.Severity != SeverityStructural {
		t.Errorf("Severity = %v; want structural", f.Severity)
	}
	if f.Location.Slot != 2 || f.Location.Filename != "stub" {
		t.Errorf("Location fields wrong; got %+v", f.Location)
	}
	if f.Message == "" || f.Citation == "" {
		t.Errorf("Message and Citation should be populated")
	}
}

func TestRegisterAndIterate(t *testing.T) {
	// Snapshot the registry, clear it for this test, restore after.
	saved := append([]Rule(nil), allRules...)
	allRules = nil
	defer func() { allRules = saved }()

	r1 := Rule{ID: "TEST-A", Severity: SeverityCosmetic, Description: "a", Citation: "x:1"}
	r2 := Rule{ID: "TEST-B", Severity: SeverityFatal, Description: "b", Citation: "x:2"}
	Register(r1)
	Register(r2)

	got := Rules()
	if len(got) != 2 {
		t.Fatalf("Rules() returned %d entries; want 2", len(got))
	}
	if got[0].ID != "TEST-A" || got[1].ID != "TEST-B" {
		t.Errorf("Rules() out of registration order: %+v", got)
	}
}

func TestRegisterRejectsDuplicateID(t *testing.T) {
	saved := append([]Rule(nil), allRules...)
	allRules = nil
	defer func() { allRules = saved }()

	Register(Rule{ID: "DUP", Severity: SeverityFatal, Description: "x", Citation: "x:1"})
	defer func() {
		if r := recover(); r == nil {
			t.Error("Register with duplicate ID did not panic")
		}
	}()
	Register(Rule{ID: "DUP", Severity: SeverityFatal, Description: "y", Citation: "x:2"})
}

func TestCheckContextShape(t *testing.T) {
	di := NewDiskImage()
	dj := di.DiskJournal()
	ctx := &CheckContext{
		Disk:    di,
		Journal: dj,
		Dialect: DialectSAMDOS2,
	}
	if ctx.Disk != di {
		t.Errorf("ctx.Disk wrong")
	}
	if ctx.Journal != dj {
		t.Errorf("ctx.Journal wrong")
	}
	if ctx.Dialect != DialectSAMDOS2 {
		t.Errorf("ctx.Dialect = %v; want samdos2", ctx.Dialect)
	}
}
