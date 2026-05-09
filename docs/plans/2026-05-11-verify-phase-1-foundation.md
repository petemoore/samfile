# Verify — Phase 1: Foundation — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the foundation of the `samfile verify` feature — type system (Severity, Dialect, Location, Finding, Rule, CheckContext, VerifyReport, FilterOpts), an internal rule registry, the `(*DiskImage).Verify()` library entry point, and a CLI subcommand scaffold — with one trivial rule wired end-to-end as a smoke test. No real rule implementations beyond the smoke test; those land in phases 3-6.

**Architecture:** Per `docs/specs/2026-05-11-verify-feature-design.md`. Rules register at package init time via a package-private `Register` function. `Verify()` builds a `CheckContext` (disk + journal + dialect), iterates the registry, runs each rule whose `Dialects` slice includes the current dialect (or is empty), and aggregates findings into a `VerifyReport`. Dialect detection itself is deferred to Phase 2 — Phase 1 passes `DialectUnknown` to every rule, so only rules tagged with empty `Dialects` (i.e. "all dialects") fire. The CLI subcommand is a thin formatter wrapping `Verify()`.

**Tech Stack:** Go 1.19+ (samfile's existing minimum). Standard library only. docopt-go is already a samfile dependency for CLI parsing; we follow the existing usage pattern in `cmd/samfile/usage.go`.

**Spec:** `docs/specs/2026-05-11-verify-feature-design.md`
**Catalog:** `docs/disk-validity-rules.md`

---

## File structure

### New files

| Path | Responsibility |
|---|---|
| `verify.go` | Core types (Severity, Dialect, Location, Finding, Rule, CheckContext, VerifyReport, FilterOpts), the rule registry, and `(*DiskImage).Verify()`. |
| `verify_test.go` | Unit tests for the types and `Verify()` plumbing. |
| `rules_smoke.go` | The single smoke-test rule (`DISK-NOT-EMPTY`). Establishes the per-category file convention that phases 3-6 follow (`rules_body_header.go`, `rules_chain.go`, etc.). |
| `rules_smoke_test.go` | Positive + negative tests for the smoke-test rule. |
| `cmd/samfile/verify.go` | CLI subcommand that calls `Verify()` and formats the report for humans. Phase 1 ships the scaffold; flag parsing for `--severity`, `--json`, etc. comes in later phases. |
| `cmd/samfile/verify_test.go` | CLI invocation tests via the existing test pattern. |

### Modified files

| Path | Change |
|---|---|
| `cmd/samfile/main.go` | Dispatch the new `verify` subcommand alongside `add`, `cat`, `extract`, `ls`, `basic-to-text`. |
| `cmd/samfile/usage.go` | Add `verify` to the usage banner and command table. |

### Why split rules out of `verify.go` from day one

The spec organises rules by category, each in its own file (`rules_disk.go`, `rules_dir.go`, `rules_chain.go`, …). The smoke-test rule seeds that pattern — putting it in `rules_smoke.go` from the start avoids a later refactor when phases 3-6 add their files. `verify.go` itself stays focused on the type system and registry.

---

## Task 1: Severity + Dialect enums

**Files:**
- Create: `verify.go`
- Test: `verify_test.go`

- [ ] **Step 1: Write the failing test**

```go
// verify_test.go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'TestSeverity|TestDialect' .`
Expected: build failure with `undefined: SeverityCosmetic` etc.

- [ ] **Step 3: Write minimal implementation**

```go
// verify.go
package samfile

// Severity ranks findings by impact, lowest to highest.
type Severity int

const (
	SeverityCosmetic Severity = iota
	SeverityInconsistency
	SeverityStructural
	SeverityFatal
)

// String returns the lowercase canonical name of the severity,
// matching the names used by the disk-validity-rules.md catalog
// and the CLI's --severity flag.
func (s Severity) String() string {
	switch s {
	case SeverityCosmetic:
		return "cosmetic"
	case SeverityInconsistency:
		return "inconsistency"
	case SeverityStructural:
		return "structural"
	case SeverityFatal:
		return "fatal"
	}
	return "unknown"
}

// Dialect identifies which DOS produced the disk. Phase 1 only
// uses DialectUnknown (dialect detection lands in Phase 2); rules
// are scoped by their Dialects slice, with nil meaning all dialects.
type Dialect int

const (
	DialectUnknown Dialect = iota
	DialectSAMDOS1
	DialectSAMDOS2
	DialectMasterDOS
)

// String returns the lowercase canonical name of the dialect,
// matching the CLI's --dialect flag.
func (d Dialect) String() string {
	switch d {
	case DialectUnknown:
		return "unknown"
	case DialectSAMDOS1:
		return "samdos1"
	case DialectSAMDOS2:
		return "samdos2"
	case DialectMasterDOS:
		return "masterdos"
	}
	return "unknown"
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -count=1 -run 'TestSeverity|TestDialect' .`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add verify.go verify_test.go
git commit -m "verify: add Severity and Dialect enums"
```

---

## Task 2: Location type

**Files:**
- Modify: `verify.go`
- Modify: `verify_test.go`

- [ ] **Step 1: Write the failing test**

```go
// append to verify_test.go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -count=1 -run 'TestLocation' .`
Expected: build failure with `undefined: DiskWideLocation` etc.

- [ ] **Step 3: Write minimal implementation**

```go
// append to verify.go
// Location pinpoints a Finding on the disk. Construct one via the
// DiskWideLocation, SlotLocation, or SectorLocation factories — they
// set the "not applicable" sentinels correctly. The zero value of
// Location is NOT a valid disk-wide location (Slot=0 is a real slot).
type Location struct {
	Slot       int     // -1 if not applicable, else 0..79
	Sector     *Sector // nil if not applicable
	ByteOffset int     // -1 if not applicable, else byte offset within Sector
	Filename   string  // copied from Slot's directory entry when known, for messages
}

// DiskWideLocation returns a Location for findings that apply to the
// disk image as a whole (no specific slot or sector).
func DiskWideLocation() Location {
	return Location{Slot: -1, Sector: nil, ByteOffset: -1}
}

// SlotLocation returns a Location for findings tied to a specific
// directory slot but not a specific sector or byte.
func SlotLocation(slot int, filename string) Location {
	return Location{Slot: slot, Sector: nil, ByteOffset: -1, Filename: filename}
}

// SectorLocation returns a Location for findings tied to a specific
// byte within a specific sector of a specific file.
func SectorLocation(slot int, filename string, sector *Sector, byteOffset int) Location {
	return Location{Slot: slot, Sector: sector, ByteOffset: byteOffset, Filename: filename}
}

// IsDiskWide reports whether loc has no slot, sector, or byte set.
func (loc Location) IsDiskWide() bool {
	return loc.Slot == -1 && loc.Sector == nil && loc.ByteOffset == -1
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -count=1 -run 'TestLocation' .`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add verify.go verify_test.go
git commit -m "verify: add Location type with DiskWide/Slot/Sector factories"
```

---

## Task 3: Finding type

**Files:**
- Modify: `verify.go`
- Modify: `verify_test.go`

- [ ] **Step 1: Write the failing test**

```go
// append to verify_test.go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -count=1 -run 'TestFinding' .`
Expected: build failure with `undefined: Finding`.

- [ ] **Step 3: Write minimal implementation**

```go
// append to verify.go
// Finding is one specific violation produced by one Rule.
//
// Message is the prose summary intended for human readers
// (default CLI output prints it directly). It should be a
// single line including the relevant Expected vs Actual
// values; multi-line context goes in a separate diagnostic.
//
// Citation duplicates the parent Rule's citation for easy
// access without a registry lookup.
type Finding struct {
	RuleID   string
	Severity Severity
	Location Location
	Message  string
	Citation string
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -count=1 -run 'TestFinding' .`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add verify.go verify_test.go
git commit -m "verify: add Finding type"
```

---

## Task 4: Rule type + Register + iteration

**Files:**
- Modify: `verify.go`
- Modify: `verify_test.go`

- [ ] **Step 1: Write the failing test**

```go
// append to verify_test.go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -count=1 -run 'TestRegister' .`
Expected: build failure with `undefined: Rule, Register, Rules, allRules`.

- [ ] **Step 3: Write minimal implementation**

```go
// append to verify.go
// Rule is a registered validity check. Check is invoked once per
// Verify run and returns zero or more Findings. Rule values are
// immutable after registration.
type Rule struct {
	ID          string      // catalog-stable, e.g. "DISK-NOT-EMPTY"
	Severity    Severity
	Dialects    []Dialect   // dialects the rule applies to; nil/empty = all
	Description string      // one-line summary, used in human output
	Citation    string      // file:line of the strongest evidence
	Check       func(ctx *CheckContext) []Finding
}

// allRules is the package-private registry. Rules register at package
// init time via Register; the order is preserved so Verify output is
// deterministic.
var allRules []Rule

// Register adds rule to the package-wide rule registry. Panics if a
// rule with the same ID is already registered (rule IDs must be
// catalog-stable and unique). Intended to be called from init().
func Register(rule Rule) {
	for _, r := range allRules {
		if r.ID == rule.ID {
			panic("samfile: duplicate rule ID registered: " + rule.ID)
		}
	}
	allRules = append(allRules, rule)
}

// Rules returns a copy of the registered rules in registration
// order. Use this for inspection (e.g. CLI help, documentation
// generators); Verify iterates allRules directly.
func Rules() []Rule {
	out := make([]Rule, len(allRules))
	copy(out, allRules)
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -count=1 -run 'TestRegister' .`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add verify.go verify_test.go
git commit -m "verify: add Rule type and Register/Rules registry API"
```

---

## Task 5: CheckContext type

**Files:**
- Modify: `verify.go`
- Modify: `verify_test.go`

- [ ] **Step 1: Write the failing test**

```go
// append to verify_test.go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -count=1 -run 'TestCheckContext' .`
Expected: build failure with `undefined: CheckContext`.

- [ ] **Step 3: Write minimal implementation**

```go
// append to verify.go
// CheckContext is the read-only environment passed to each Rule's
// Check function. All disk inspection should go through ctx — Rules
// must NOT call disk.DiskJournal() themselves (the journal is
// computed once per Verify run and shared). If a future rule needs
// another expensive derivation (e.g. a combined sector map), add
// it as a field on CheckContext and memoise it in Verify.
type CheckContext struct {
	Disk    *DiskImage
	Journal *DiskJournal
	Dialect Dialect
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -count=1 -run 'TestCheckContext' .`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add verify.go verify_test.go
git commit -m "verify: add CheckContext type"
```

---

## Task 6: VerifyReport + filter methods + FilterOpts

**Files:**
- Modify: `verify.go`
- Modify: `verify_test.go`

- [ ] **Step 1: Write the failing test**

```go
// append to verify_test.go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -count=1 -run 'TestVerifyReport' .`
Expected: build failure with `undefined: VerifyReport, FilterOpts`.

- [ ] **Step 3: Write minimal implementation**

```go
// append to verify.go
// VerifyReport is the result of running Verify on a DiskImage.
// Findings is the full ordered slice; the helper methods filter
// without mutating.
type VerifyReport struct {
	Dialect  Dialect
	Findings []Finding
}

// HasFatal reports whether any finding has severity SeverityFatal.
func (r VerifyReport) HasFatal() bool {
	for _, f := range r.Findings {
		if f.Severity == SeverityFatal {
			return true
		}
	}
	return false
}

// HasStructural reports whether any finding has severity
// SeverityStructural or higher.
func (r VerifyReport) HasStructural() bool {
	for _, f := range r.Findings {
		if f.Severity >= SeverityStructural {
			return true
		}
	}
	return false
}

// BySeverity returns findings with exactly the given severity, in
// registration order.
func (r VerifyReport) BySeverity(s Severity) []Finding {
	var out []Finding
	for _, f := range r.Findings {
		if f.Severity == s {
			out = append(out, f)
		}
	}
	return out
}

// ByRule returns findings produced by the rule with the given ID,
// in registration order.
func (r VerifyReport) ByRule(ruleID string) []Finding {
	var out []Finding
	for _, f := range r.Findings {
		if f.RuleID == ruleID {
			out = append(out, f)
		}
	}
	return out
}

// FilterOpts controls VerifyReport.Filter. Zero-value fields act
// as "no constraint".
type FilterOpts struct {
	MinSeverity Severity // findings with severity >= MinSeverity pass
	Rules       []string // if non-empty, only these rule IDs pass
	Slot        *int     // if non-nil, only findings at this slot pass
}

// Filter returns findings matching every set constraint in opts, in
// registration order. An empty FilterOpts returns r.Findings.
func (r VerifyReport) Filter(opts FilterOpts) []Finding {
	var out []Finding
	for _, f := range r.Findings {
		if f.Severity < opts.MinSeverity {
			continue
		}
		if len(opts.Rules) > 0 {
			matched := false
			for _, id := range opts.Rules {
				if f.RuleID == id {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		if opts.Slot != nil && f.Location.Slot != *opts.Slot {
			continue
		}
		out = append(out, f)
	}
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -count=1 -run 'TestVerifyReport' .`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add verify.go verify_test.go
git commit -m "verify: add VerifyReport with HasFatal/BySeverity/ByRule/Filter"
```

---

## Task 7: Verify() method on *DiskImage

**Files:**
- Modify: `verify.go`
- Modify: `verify_test.go`

- [ ] **Step 1: Write the failing test**

```go
// append to verify_test.go
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
		Check: func(ctx *CheckContext) []Finding { allDialects++; return nil },
	})
	Register(Rule{
		ID:          "MASTERDOS-ONLY",
		Severity:    SeverityCosmetic,
		Dialects:    []Dialect{DialectMasterDOS},
		Description: "masterdos only",
		Citation:    "test",
		Check: func(ctx *CheckContext) []Finding { scoped++; return nil },
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -count=1 -run 'TestVerify' .`
Expected: build failure with `di.Verify undefined`.

- [ ] **Step 3: Write minimal implementation**

```go
// append to verify.go
// Verify runs all registered rules against di and returns a report
// describing the disk's structural state. The report is always
// populated — individual rule failures are surfaced as Findings,
// not Go errors. Verify itself does not return an error.
//
// In Phase 1, dialect detection is not yet implemented and Verify
// always passes DialectUnknown to rules; rules whose Dialects slice
// is non-empty and excludes DialectUnknown are skipped. Phase 2
// adds DetectDialect.
func (di *DiskImage) Verify() VerifyReport {
	dialect := DialectUnknown
	ctx := &CheckContext{
		Disk:    di,
		Journal: di.DiskJournal(),
		Dialect: dialect,
	}
	report := VerifyReport{Dialect: dialect}
	for _, rule := range allRules {
		if !ruleAppliesToDialect(rule, dialect) {
			continue
		}
		report.Findings = append(report.Findings, rule.Check(ctx)...)
	}
	return report
}

// ruleAppliesToDialect reports whether rule should run when the
// detected dialect is d. A rule with no Dialects field set (nil or
// empty) applies to all dialects.
func ruleAppliesToDialect(rule Rule, d Dialect) bool {
	if len(rule.Dialects) == 0 {
		return true
	}
	for _, allowed := range rule.Dialects {
		if allowed == d {
			return true
		}
	}
	return false
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -count=1 -run 'TestVerify' .`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add verify.go verify_test.go
git commit -m "verify: add (*DiskImage).Verify() entry point"
```

---

## Task 8: DISK-NOT-EMPTY smoke-test rule

**Files:**
- Create: `rules_smoke.go`
- Create: `rules_smoke_test.go`

- [ ] **Step 1: Write the failing test**

```go
// rules_smoke_test.go
package samfile

import (
	"testing"
)

func TestDiskNotEmptyRulePositive(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("A", []byte("hello"), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	findings := checkDiskNotEmpty(&CheckContext{
		Disk:    di,
		Journal: di.DiskJournal(),
		Dialect: DialectUnknown,
	})
	if len(findings) != 0 {
		t.Errorf("checkDiskNotEmpty on populated disk returned %d findings; want 0", len(findings))
	}
}

func TestDiskNotEmptyRuleNegative(t *testing.T) {
	di := NewDiskImage()
	findings := checkDiskNotEmpty(&CheckContext{
		Disk:    di,
		Journal: di.DiskJournal(),
		Dialect: DialectUnknown,
	})
	if len(findings) != 1 {
		t.Fatalf("checkDiskNotEmpty on empty disk returned %d findings; want 1", len(findings))
	}
	f := findings[0]
	if f.RuleID != "DISK-NOT-EMPTY" {
		t.Errorf("RuleID = %q; want DISK-NOT-EMPTY", f.RuleID)
	}
	if f.Severity != SeverityInconsistency {
		t.Errorf("Severity = %v; want inconsistency", f.Severity)
	}
	if !f.Location.IsDiskWide() {
		t.Errorf("Location.IsDiskWide() = false; want true")
	}
	if f.Message == "" {
		t.Error("Message empty")
	}
}

func TestDiskNotEmptyRegistered(t *testing.T) {
	found := false
	for _, r := range Rules() {
		if r.ID == "DISK-NOT-EMPTY" {
			found = true
			break
		}
	}
	if !found {
		t.Error("DISK-NOT-EMPTY rule not in registry")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -count=1 -run 'TestDiskNotEmpty' .`
Expected: build failure with `undefined: checkDiskNotEmpty`.

- [ ] **Step 3: Write minimal implementation**

```go
// rules_smoke.go
package samfile

import "fmt"

// DISK-NOT-EMPTY is the Phase 1 smoke-test rule: it fires on a disk
// with zero occupied directory entries. The "real" rule catalog has
// dozens of these; this is the one we wire up end-to-end in Phase 1
// to prove the registry + Verify plumbing works. Severity is
// inconsistency rather than fatal because an empty disk is unusual
// but technically valid SAM-format output.
func init() {
	Register(Rule{
		ID:          "DISK-NOT-EMPTY",
		Severity:    SeverityInconsistency,
		Dialects:    nil, // all dialects
		Description: "disk has at least one occupied directory entry",
		Citation:    "docs/disk-validity-rules.md",
		Check:       checkDiskNotEmpty,
	})
}

func checkDiskNotEmpty(ctx *CheckContext) []Finding {
	used := ctx.Journal.UsedFileEntries()
	if len(used) > 0 {
		return nil
	}
	return []Finding{{
		RuleID:   "DISK-NOT-EMPTY",
		Severity: SeverityInconsistency,
		Location: DiskWideLocation(),
		Message:  fmt.Sprintf("disk has 0 occupied directory entries (all %d slots are free)", len(ctx.Journal.FreeFileEntries())),
		Citation: "docs/disk-validity-rules.md",
	}}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -count=1 -run 'TestDiskNotEmpty' .`
Expected: PASS.

Also run the whole package to confirm nothing else broke:
Run: `go test -count=1 .`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add rules_smoke.go rules_smoke_test.go
git commit -m "verify: add DISK-NOT-EMPTY smoke-test rule"
```

---

## Task 9: CLI verify subcommand

**Files:**
- Create: `cmd/samfile/verify.go`
- Create: `cmd/samfile/verify_test.go`

Before writing the subcommand, look at one of the existing subcommands (`cmd/samfile/ls.go` is the simplest) to confirm the call-shape: each subcommand exposes a `Run(args map[string]any) error` function (or similar — match whatever pattern is there).

- [ ] **Step 1: Read existing subcommand for pattern**

Run: `cat cmd/samfile/ls.go cmd/samfile/main.go`
Note: you should now know how subcommands are invoked, what args look like, and how errors are surfaced. The verify subcommand follows the same pattern. If your shape below doesn't match what `ls.go` does, adapt to match the existing convention before writing the test.

- [ ] **Step 2: Write the failing test**

```go
// cmd/samfile/verify_test.go
package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/petemoore/samfile/v3"
)

// TestVerifyCmdOnPopulatedDisk runs the CLI subcommand against a
// disk built in-memory with one CODE file. Expected: DISK-NOT-EMPTY
// does not fire; output reports zero findings; exit code (returned
// as nil error) is clean.
func TestVerifyCmdOnPopulatedDisk(t *testing.T) {
	di := samfile.NewDiskImage()
	if err := di.AddCodeFile("F", []byte("hello"), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "test.mgt")
	if err := di.Save(imgPath); err != nil {
		t.Fatalf("Save: %v", err)
	}

	stdout, err := captureVerify(t, imgPath)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if !strings.Contains(stdout, "0 finding") {
		t.Errorf("expected '0 finding' in output; got:\n%s", stdout)
	}
}

// TestVerifyCmdOnEmptyDisk runs against an empty disk. Expected:
// DISK-NOT-EMPTY fires; output mentions the rule ID; error is nil
// (inconsistency does not gate exit code).
func TestVerifyCmdOnEmptyDisk(t *testing.T) {
	di := samfile.NewDiskImage()
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "test.mgt")
	if err := di.Save(imgPath); err != nil {
		t.Fatalf("Save: %v", err)
	}

	stdout, err := captureVerify(t, imgPath)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if !strings.Contains(stdout, "DISK-NOT-EMPTY") {
		t.Errorf("expected 'DISK-NOT-EMPTY' in output; got:\n%s", stdout)
	}
}

// captureVerify invokes the verify subcommand and returns its stdout.
// Mirrors the helper pattern in cat_test.go (if it exists) — adapt
// to whatever convention cmd/samfile uses.
func captureVerify(t *testing.T, imgPath string) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	err := runVerify(imgPath)

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String(), err
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test -count=1 ./cmd/samfile/`
Expected: build failure with `undefined: runVerify`.

- [ ] **Step 4: Write minimal implementation**

```go
// cmd/samfile/verify.go
package main

import (
	"fmt"

	"github.com/petemoore/samfile/v3"
)

// runVerify is the entry point for the `samfile verify` subcommand.
// Phase 1 implements the minimum useful behaviour: load the disk,
// call Verify, print findings grouped by severity, return nil
// unless a fatal finding is present (in which case return a non-nil
// error so main exits non-zero).
//
// CLI flags (--severity, --json, --dialect, --rule, --quiet, --all)
// are deferred to a later phase. Phase 1 always shows every finding
// regardless of severity, so the smoke test can see DISK-NOT-EMPTY.
func runVerify(imagePath string) error {
	di, err := samfile.Load(imagePath)
	if err != nil {
		return fmt.Errorf("verify: %w", err)
	}

	report := di.Verify()

	fmt.Printf("samfile verify: results for %s\n", imagePath)
	fmt.Printf("detected dialect: %s\n", report.Dialect)
	fmt.Println()

	if len(report.Findings) == 0 {
		fmt.Println("no findings.")
		return nil
	}

	// Group by severity, highest first.
	severities := []samfile.Severity{
		samfile.SeverityFatal,
		samfile.SeverityStructural,
		samfile.SeverityInconsistency,
		samfile.SeverityCosmetic,
	}
	for _, s := range severities {
		findings := report.BySeverity(s)
		if len(findings) == 0 {
			continue
		}
		fmt.Printf("%s (%d):\n", upperString(s.String()), len(findings))
		for _, f := range findings {
			fmt.Printf("  [%s]", f.RuleID)
			if !f.Location.IsDiskWide() {
				fmt.Printf(" slot %d", f.Location.Slot)
				if f.Location.Filename != "" {
					fmt.Printf(" (%s)", f.Location.Filename)
				}
			}
			fmt.Println()
			fmt.Printf("    %s\n", f.Message)
			fmt.Printf("    citation: %s\n", f.Citation)
		}
		fmt.Println()
	}

	fmt.Printf("%d finding(s).\n", len(report.Findings))
	if report.HasFatal() {
		return fmt.Errorf("verify: %d fatal finding(s)", len(report.BySeverity(samfile.SeverityFatal)))
	}
	return nil
}

// upperString returns s with its first byte uppercased. We avoid
// strings.Title (deprecated) and the unicode package since severity
// names are guaranteed ASCII.
func upperString(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 32
	}
	return string(b)
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test -count=1 ./cmd/samfile/`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add cmd/samfile/verify.go cmd/samfile/verify_test.go
git commit -m "verify: add CLI subcommand scaffold (no flags yet)"
```

---

## Task 10: Wire verify into main.go + usage.go

**Files:**
- Modify: `cmd/samfile/main.go`
- Modify: `cmd/samfile/usage.go`

End-to-end binary check rather than unit test for this task — `runVerify` is already covered by Task 9's tests; here we're verifying that the docopt parser knows about `verify` and `main.go` dispatches it correctly. Running the built binary catches both.

- [ ] **Step 1: Read existing dispatch code**

Run: `cat cmd/samfile/main.go cmd/samfile/usage.go`

Note the (a) docopt usage string format, (b) how arguments are extracted from the parsed args map (`args["-i"].(string)`, etc.), and (c) the existing subcommand switch / dispatch chain. The verify subcommand follows the same pattern: one new usage line, one new switch arm.

- [ ] **Step 2: Add verify to usage.go**

Find the `Usage:` block in `usage.go` (it's a docopt string passed to `docopt.ParseArgs`). Add a new line for verify, alongside the existing subcommand lines:

```
  samfile verify -i IMAGE
```

Place it in whatever order the existing subcommands use (alphabetical, declaration-order, whatever). Do **not** rewrite the file from scratch — read it and add the one line.

- [ ] **Step 3: Add verify to main.go**

Find the subcommand-dispatch switch in `main.go`. Add a new arm calling `runVerify`:

```go
case args["verify"].(bool):
    return runVerify(args["-i"].(string))
```

Insertion order matches the rest of the switch. If the switch returns errors and `main()` already handles them (likely — that's the existing pattern), no `main()` changes are needed; if not, follow whatever convention the file uses.

- [ ] **Step 4: Run all unit tests**

Run: `go test -count=1 ./...`
Expected: PASS across `.` and `./cmd/samfile/`. (The smoke test in Task 9 will exercise the same `runVerify` function the dispatch now calls.)

- [ ] **Step 5: Build the binary and check the help output**

Run: `go build -o /tmp/samfile-verify-test ./cmd/samfile`
Run: `/tmp/samfile-verify-test -h | grep verify`
Expected: at least one line mentioning `verify`.

- [ ] **Step 6: Build a sample disk and run the binary end-to-end**

```bash
# Construct an empty disk and a populated one via the package API,
# then exercise the binary against both.
cat <<'EOF' > /tmp/buildtwo.go
package main

import (
	"log"

	"github.com/petemoore/samfile/v3"
)

func main() {
	empty := samfile.NewDiskImage()
	if err := empty.Save("/tmp/empty.mgt"); err != nil {
		log.Fatal(err)
	}
	pop := samfile.NewDiskImage()
	if err := pop.AddCodeFile("F", []byte("hi"), 0x8000, 0); err != nil {
		log.Fatal(err)
	}
	if err := pop.Save("/tmp/pop.mgt"); err != nil {
		log.Fatal(err)
	}
}
EOF
go run /tmp/buildtwo.go
/tmp/samfile-verify-test verify -i /tmp/empty.mgt
echo "exit=$?"
/tmp/samfile-verify-test verify -i /tmp/pop.mgt
echo "exit=$?"
rm /tmp/buildtwo.go /tmp/empty.mgt /tmp/pop.mgt /tmp/samfile-verify-test
```

Expected outputs:

- `empty.mgt`: report mentions `DISK-NOT-EMPTY` under `INCONSISTENCY (1)`; `exit=0`.
- `pop.mgt`: report says `no findings.`; `exit=0`.

If either fails, the dispatch wiring is wrong; check the docopt usage string and the switch arm.

- [ ] **Step 7: Commit**

```bash
git add cmd/samfile/main.go cmd/samfile/usage.go
git commit -m "verify: wire CLI subcommand into main dispatch and usage"
```

---

## Definition of done (Phase 1)

- `go test -count=1 ./...` passes locally.
- CI green on the PR branch.
- `samfile verify -i <empty.mgt>` prints a single `DISK-NOT-EMPTY` finding under "INCONSISTENCY (1):", exits 0.
- `samfile verify -i <populated.mgt>` prints "no findings.", exits 0.
- `samfile -h` lists `verify`.
- `(*samfile.DiskImage).Verify()` is callable from external Go packages and returns a `VerifyReport` with the expected fields.
- No real rule implementations exist yet — only the smoke-test rule. Phase 2 adds `DetectDialect`; phases 3-6 add the catalog's 70 rules. Each phase gets its own plan written when it starts.

## Out of scope for Phase 1

- CLI flags `--severity`, `--all`, `--rule`, `--dialect`, `--json`, `--quiet`. Default Phase 1 output is "show everything"; flag handling is added in a later phase before any user-facing release.
- `DetectDialect` heuristic.
- Any rule from the catalog.
- Testdata corpus in `testdata/mgt/`. Phase 1's tests build disks in-memory.

## After Phase 1 lands

Write Phase 2's plan against `docs/specs/2026-05-11-verify-feature-design.md` § Dialect detection.
