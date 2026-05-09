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

func TestVerifyReportHelpers(t *testing.T) {
	r := VerifyReport{
		Dialect: DialectSAMDOS2,
		Findings: []Finding{
			{RuleID: "A", Severity: SeverityFatal, Location: DiskWideLocation()},
			{RuleID: "B", Severity: SeverityStructural, Location: SlotLocation(2, "stub")},
			{RuleID: "C", Severity: SeverityCosmetic, Location: DiskWideLocation()},
			{RuleID: "B", Severity: SeverityStructural, Location: SlotLocation(3, "IN")},
		},
	}

	if !r.HasFatal() {
		t.Error("HasFatal() = false; want true")
	}
	if !r.HasStructural() {
		t.Error("HasStructural() = false; want true")
	}

	if got := r.BySeverity(SeverityStructural); len(got) != 2 {
		t.Errorf("BySeverity(structural) returned %d; want 2", len(got))
	}
	if got := r.ByRule("B"); len(got) != 2 {
		t.Errorf("ByRule(B) returned %d; want 2", len(got))
	}
	if got := r.Filter(FilterOpts{MinSeverity: SeverityStructural}); len(got) != 3 {
		t.Errorf("Filter(min=structural) returned %d; want 3 (B, A, B in registration order)", len(got))
	}
	if got := r.Filter(FilterOpts{Rules: []string{"A"}}); len(got) != 1 {
		t.Errorf("Filter(rules=[A]) returned %d; want 1", len(got))
	}
}

func TestVerifyReportHasFatalEmpty(t *testing.T) {
	r := VerifyReport{}
	if r.HasFatal() || r.HasStructural() {
		t.Error("empty report should not report Has*")
	}
}

func TestVerifyRunsRegisteredRules(t *testing.T) {
	saved := append([]Rule(nil), allRules...)
	allRules = nil
	defer func() { allRules = saved }()

	called := 0
	Register(Rule{
		ID:          "X-1",
		Severity:    SeverityFatal,
		Description: "always fires",
		Citation:    "test",
		Check: func(ctx *CheckContext) []Finding {
			called++
			return []Finding{{
				RuleID:   "X-1",
				Severity: SeverityFatal,
				Location: DiskWideLocation(),
				Message:  "test finding",
				Citation: "test",
			}}
		},
	})

	di := NewDiskImage()
	report := di.Verify()

	if called != 1 {
		t.Errorf("Check called %d times; want 1", called)
	}
	if len(report.Findings) != 1 {
		t.Fatalf("Findings = %d; want 1", len(report.Findings))
	}
	if report.Findings[0].RuleID != "X-1" {
		t.Errorf("Findings[0].RuleID = %q; want X-1", report.Findings[0].RuleID)
	}
}

func TestVerifyRespectsDialectScoping(t *testing.T) {
	saved := append([]Rule(nil), allRules...)
	allRules = nil
	defer func() { allRules = saved }()

	allDialects := 0
	scoped := 0
	Register(Rule{
		ID:          "ALL",
		Severity:    SeverityCosmetic,
		Description: "all dialects",
		Citation:    "test",
		Check:       func(ctx *CheckContext) []Finding { allDialects++; return nil },
	})
	Register(Rule{
		ID:          "MASTERDOS-ONLY",
		Severity:    SeverityCosmetic,
		Dialects:    []Dialect{DialectMasterDOS},
		Description: "masterdos only",
		Citation:    "test",
		Check:       func(ctx *CheckContext) []Finding { scoped++; return nil },
	})

	di := NewDiskImage()
	di.Verify() // Phase 1 always passes DialectUnknown

	if allDialects != 1 {
		t.Errorf("all-dialects rule called %d times; want 1", allDialects)
	}
	if scoped != 0 {
		t.Errorf("masterdos-only rule called %d times; want 0 (dialect is Unknown)", scoped)
	}
}

func TestVerifyReportCarriesDialect(t *testing.T) {
	saved := append([]Rule(nil), allRules...)
	allRules = nil
	defer func() { allRules = saved }()

	di := NewDiskImage()
	report := di.Verify()
	// Phase 1: dialect detection is not implemented; always DialectUnknown.
	if report.Dialect != DialectUnknown {
		t.Errorf("Dialect = %v; want unknown (detection lands in Phase 2)", report.Dialect)
	}
}
