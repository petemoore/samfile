package samfile

import "fmt"

// Severity ranks findings by impact, lowest to highest.
//
// Severity names are stable public API; the numeric values
// assigned by iota are NOT — new severities may be inserted
// between existing ones, shifting integer values. Don't serialise
// the raw int.
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
	return fmt.Sprintf("Severity(%d)", s)
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
	return fmt.Sprintf("Dialect(%d)", d)
}

// Location pinpoints a Finding on the disk. Construct one via the
// DiskWideLocation, SlotLocation, or SectorLocation factories — they
// set the "not applicable" sentinels correctly. The zero value of
// Location is NOT a valid disk-wide location (Slot=0 is a real slot).
type Location struct {
	Slot       int     // -1 if not applicable, else 0..79
	Sector     *Sector // nil if not applicable
	ByteOffset int     // -1 if not applicable, else byte offset within Sector
	Filename   string  // file name from the slot's directory entry, embedded for message formatting; "" when Slot is -1
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

	// Legacy single-shot check. Mutually exclusive with CheckSubject.
	// Kept for rules that haven't been migrated to the per-subject
	// model. If CheckSubject is non-nil, Check is ignored.
	Check func(ctx *CheckContext) []Finding

	// Scope identifies the per-subject iteration scope when
	// CheckSubject is set. Ignored for legacy (Check-only) rules.
	Scope SubjectScope

	// Applies reports whether subject is eligible for this rule. If
	// nil, the rule applies to all subjects of its Scope. The number
	// of applicable subjects is the denominator for fail-rate
	// analysis in the audit pipeline.
	Applies func(ctx *CheckContext, subject Subject) bool

	// CheckSubject evaluates one applicable subject. Returns nil for
	// pass, a Finding pointer for fail. The framework calls this once
	// per applicable subject and emits a CheckEvent for each call.
	CheckSubject func(ctx *CheckContext, subject Subject) *Finding
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
	Disk     *DiskImage
	Journal  *DiskJournal
	Dialect  Dialect
	recorder *EventRecorder
}

// SetRecorder installs an EventRecorder that receives a CheckEvent
// per (rule, applicable subject) pair during Verify. nil is fine —
// the text-output path leaves the recorder unset.
func (ctx *CheckContext) SetRecorder(r *EventRecorder) { ctx.recorder = r }

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
//
// FilterOpts intentionally has no Dialect field: a VerifyReport
// has a single Dialect for the whole disk, so per-finding
// dialect filtering would be meaningless. If you need to scope
// by dialect, do so on the report producer (or pass --dialect
// on the CLI).
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
// Verify calls DetectDialect to infer the dialect that wrote di,
// then runs every registered rule whose Dialects slice is empty
// (all-dialects) or contains the detected dialect. Rules scoped to a
// dialect other than the one detected are skipped. DetectDialect is
// conservative: when it returns DialectUnknown (empty or ambiguous
// disks), only all-dialects rules run.
func (di *DiskImage) Verify() VerifyReport {
	return di.verifyInternal(nil)
}

// VerifyWithRecorder runs verify and captures every Check event into
// the provided recorder (in addition to producing the VerifyReport).
// Use this from the JSONL CLI path.
func (di *DiskImage) VerifyWithRecorder(rec *EventRecorder) VerifyReport {
	return di.verifyInternal(rec)
}

func (di *DiskImage) verifyInternal(rec *EventRecorder) VerifyReport {
	dialect := DetectDialect(di)
	ctx := &CheckContext{
		Disk:     di,
		Journal:  di.DiskJournal(),
		Dialect:  dialect,
		recorder: rec,
	}
	report := VerifyReport{Dialect: dialect}
	for _, rule := range allRules {
		if !ruleAppliesToDialect(rule, dialect) {
			continue
		}
		if rule.CheckSubject != nil {
			subjects := ctx.subjectsForScope(rule.Scope)
			for _, subj := range subjects {
				applicable := rule.Applies == nil || rule.Applies(ctx, subj)
				if !applicable {
					rec.Record(CheckEvent{
						RuleID: rule.ID, Scope: rule.Scope.String(),
						Ref: subj.Ref(), Outcome: "not_applicable",
						Attrs: subj.Attributes(),
					})
					continue
				}
				finding := rule.CheckSubject(ctx, subj)
				outcome := "pass"
				if finding != nil {
					outcome = "fail"
					report.Findings = append(report.Findings, *finding)
				}
				rec.Record(CheckEvent{
					RuleID: rule.ID, Scope: rule.Scope.String(),
					Ref: subj.Ref(), Outcome: outcome,
					Attrs: subj.Attributes(), Finding: finding,
				})
			}
			continue
		}
		// Legacy path: rule provides Check only.
		legacyFindings := rule.Check(ctx)
		report.Findings = append(report.Findings, legacyFindings...)
		for i := range legacyFindings {
			f := legacyFindings[i]
			ref := "disk"
			if f.Location.Slot >= 0 {
				ref = fmt.Sprintf("slot=%d", f.Location.Slot)
			}
			rec.Record(CheckEvent{
				RuleID: rule.ID, Scope: "legacy",
				Ref: ref, Outcome: "fail", Finding: &f,
			})
		}
	}
	return report
}

// subjectsForScope enumerates every Subject of the given scope on
// the current disk. DiskScope yields one DiskSubject; SlotScope
// yields 80 SlotSubjects (every dir entry, used or erased).
// ChainStepScope is not yet implemented — rules that need per-step
// iteration should declare SlotScope and walk the chain internally.
func (ctx *CheckContext) subjectsForScope(scope SubjectScope) []Subject {
	switch scope {
	case DiskScope:
		return []Subject{&DiskSubject{Journal: ctx.Journal, Disk: ctx.Disk, Dialect: ctx.Dialect}}
	case SlotScope:
		out := make([]Subject, 0, 80)
		for i := 0; i < 80; i++ {
			fe := (*ctx.Journal)[i]
			out = append(out, &SlotSubject{SlotIndex: i, FileEntry: fe, Disk: ctx.Disk, Journal: ctx.Journal})
		}
		return out
	}
	return nil
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
