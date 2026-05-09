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

// Rule is a registered validity check. Check is invoked once per
// Verify run and returns zero or more Findings. Rule values are
// immutable after registration.
type Rule struct {
	ID          string // catalog-stable, e.g. "DISK-NOT-EMPTY"
	Severity    Severity
	Dialects    []Dialect // dialects the rule applies to; nil/empty = all
	Description string    // one-line summary, used in human output
	Citation    string    // file:line of the strongest evidence
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
