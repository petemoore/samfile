# Verify Audit Framework Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Instrument the verify pipeline so every rule emits a Check event (pass / fail / not_applicable) per subject with attribute snapshots; ingest those into a SQLite `checks` table and emit five Markdown reports under `~/sam-corpus/analyses/`, prioritising fallback-floor reports (`coverage.md`, `disk-health.md`) that need only plain pandas counts; then loop on high-confidence patterns with source-citation grounding to land rule fixes on `feat/verify-audit-framework`.

**Architecture:** Add `Scope` + `Applies` + `CheckSubject` to the Go `Rule` struct. Framework enumerates per-scope subjects (Disk / Slot / ChainStep), evaluates Applies (denominator) and CheckSubject (numerator + finding payload), emits `CheckEvent`s. A new `--format jsonl` flag on `samfile verify` writes one JSON event per line. Python `tools/audit/{ingest,mine}.py` slurps JSONL into `findings.db` and produces reports.

**Tech Stack:** Go 1.21+ (samfile repo). Python 3.11+ with pandas, sklearn (`DecisionTreeClassifier`), mlxtend (`apriori`, `association_rules`). SQLite 3 (`findings.db`).

**Reference paths:**
- samfile repo: `~/git/samfile/`
- corpus workspace: `~/sam-corpus/`
- SAMDOS source: `~/git/samdos/src/*.s`
- ROM disasm: `~/git/sam-aarch64/docs/sam/sam-coupe_rom-v3.0_annotated-disassembly.txt`
- Design spec: `docs/specs/2026-05-12-verify-audit-framework-design.md`

---

## Phase 1 — Foundation (Go)

### Task 1: Subject interface + scope types

**Files:**
- Create: `subject.go`

- [ ] **Step 1:** Create `subject.go` with:

```go
package samfile

// SubjectScope identifies the kind of subject a Rule operates on.
type SubjectScope int

const (
	DiskScope SubjectScope = iota
	SlotScope
	ChainStepScope
)

func (s SubjectScope) String() string {
	switch s {
	case DiskScope:
		return "disk"
	case SlotScope:
		return "slot"
	case ChainStepScope:
		return "chain_step"
	}
	return "unknown"
}

// Subject is the universal interface for anything a rule can check —
// a whole disk, a directory slot, or a single sector in a file's
// chain. Ref() returns a stable string id; Attributes() returns the
// denormalised attribute snapshot recorded with every Check event.
type Subject interface {
	Ref() string
	Attributes() map[string]any
}
```

- [ ] **Step 2:** Commit.

```bash
git add subject.go
git commit -m "feat(verify): add Subject interface + SubjectScope enum"
```

---

### Task 2: DiskSubject

**Files:**
- Create: `subject_disk.go`

- [ ] **Step 1:** Create `subject_disk.go`:

```go
package samfile

// DiskSubject wraps a DiskJournal + DiskImage for disk-scope rule
// evaluation. Single instance per Verify() run.
type DiskSubject struct {
	Journal *DiskJournal
	Disk    *DiskImage
	Dialect Dialect
}

func (s *DiskSubject) Ref() string { return "disk" }

func (s *DiskSubject) Attributes() map[string]any {
	j := s.Journal
	used := 0
	for _, fe := range j.UsedSlots() {
		_ = fe
		used++
	}
	bootSig := false
	if sd, err := s.Disk.SectorData(&Sector{Track: 4, Sector: 1}); err == nil {
		if string(sd[256:259]) == "BOOT" {
			bootSig = true
		}
	}
	return map[string]any{
		"dialect":                s.Dialect.String(),
		"boot_signature_present": bootSig,
		"used_slot_count":        used,
	}
}
```

- [ ] **Step 2:** Verify it compiles.

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 3:** Commit.

```bash
git add subject_disk.go
git commit -m "feat(verify): add DiskSubject with dialect/boot/used-slot attrs"
```

---

### Task 3: SlotSubject

**Files:**
- Create: `subject_slot.go`

- [ ] **Step 1:** Create `subject_slot.go`:

```go
package samfile

import (
	"encoding/hex"
	"fmt"
)

// SlotSubject wraps one directory slot's FileEntry for slot-scope
// rule evaluation. SlotIndex is the 0..79 position; FileEntry is
// the parsed dir entry. For erased slots (Type==0) the FileEntry
// still exists but most fields are zero.
type SlotSubject struct {
	SlotIndex int
	FileEntry *FileEntry
	Disk      *DiskImage
	Journal   *DiskJournal
}

func (s *SlotSubject) Ref() string { return fmt.Sprintf("slot=%d", s.SlotIndex) }

func (s *SlotSubject) Attributes() map[string]any {
	fe := s.FileEntry
	pageOffsetForm := "other"
	switch fe.StartAddressPageOffset & 0xC000 {
	case 0x8000:
		pageOffsetForm = "0x8000"
	case 0x4000:
		pageOffsetForm = "0x4000"
	case 0x0000:
		pageOffsetForm = "0x0000"
	case 0xC000:
		pageOffsetForm = "0xC000"
	}
	dirMirror := s.Disk[directoryByteOffset(s.SlotIndex)+0xD3 : directoryByteOffset(s.SlotIndex)+0xDC]
	dirMirrorPopulated := false
	for _, b := range dirMirror {
		if b != 0 {
			dirMirrorPopulated = true
			break
		}
	}
	hasAutoRun := fe.ExecutionAddressDiv16K != 0xFF
	return map[string]any{
		"slot_index":              s.SlotIndex,
		"filename":                fe.Name.String(),
		"file_type":               fe.Type.String(),
		"file_type_byte":          int(fe.Type),
		"file_length":             int(fe.LengthMod16K),
		"page_offset_form":        pageOffsetForm,
		"pages":                   int(fe.Pages),
		"mgt_flags":               int(fe.MGTFlags),
		"first_track":             int(fe.FirstSector.Track),
		"first_sector":            int(fe.FirstSector.Sector),
		"first_side":              int(fe.FirstSector.Track >> 7),
		"has_autorun_or_autoexec": hasAutoRun,
		"dir_mirror_populated":    dirMirrorPopulated,
		"slot_is_erased":          fe.Type == 0,
		"file_type_info_hex":      hex.EncodeToString(fe.FileTypeInfo[:]),
		"sectors_count":           int(fe.Sectors),
	}
}

// directoryByteOffset returns the byte offset of dir slot N (0..79)
// within the MGT disk image, accounting for the cylinder-interleaved
// layout. Tracks 0..3 of side 0 hold the directory; each track has
// 20 slots; each slot is 256 bytes.
func directoryByteOffset(slot int) int {
	track := slot / 20
	slotInTrack := slot % 20
	return track*10240 + slotInTrack*256
}
```

- [ ] **Step 2:** Verify it compiles.

```bash
go build ./...
```

- [ ] **Step 3:** Commit.

```bash
git add subject_slot.go
git commit -m "feat(verify): add SlotSubject with slot-attribute snapshot"
```

---

### Task 4: ChainStepSubject

**Files:**
- Create: `subject_chain.go`

- [ ] **Step 1:** Create `subject_chain.go`:

```go
package samfile

import "fmt"

// ChainStepSubject wraps one step in a file's sector chain. Track
// and Sector identify the sector itself; ChainIndex is the 0..N-1
// position within this slot's chain.
type ChainStepSubject struct {
	SlotIndex   int
	ChainIndex  int
	Track       uint8
	Sector      uint8
	NextTrack   uint8
	NextSector  uint8
	Position    string // first | intermediate | last | orphan
	OnSAMMap    bool
	OnDirSAMMap bool
	Disk        *DiskImage
	Journal     *DiskJournal
}

func (s *ChainStepSubject) Ref() string {
	return fmt.Sprintf("slot=%d,chain=%d,track=%d,sector=%d", s.SlotIndex, s.ChainIndex, s.Track, s.Sector)
}

func (s *ChainStepSubject) Attributes() map[string]any {
	side := 0
	if s.Track&0x80 != 0 {
		side = 1
	}
	// Distance from dir tracks: dir is tracks 0..3 of side 0.
	dirDistance := 999
	t := int(s.Track & 0x7F)
	if side == 0 {
		if t >= 4 {
			dirDistance = t - 3
		} else {
			dirDistance = 0
		}
	} else {
		dirDistance = t + 1
	}
	return map[string]any{
		"slot_index":               s.SlotIndex,
		"chain_index":              s.ChainIndex,
		"chain_position":           s.Position,
		"track":                    int(s.Track),
		"sector":                   int(s.Sector),
		"side":                     side,
		"next_track":               int(s.NextTrack),
		"next_sector":              int(s.NextSector),
		"on_sam_map":               s.OnSAMMap,
		"on_dir_sam_map":           s.OnDirSAMMap,
		"distance_from_dir_tracks": dirDistance,
	}
}
```

- [ ] **Step 2:** Verify it compiles.

```bash
go build ./...
```

- [ ] **Step 3:** Commit.

```bash
git add subject_chain.go
git commit -m "feat(verify): add ChainStepSubject with chain-step attributes"
```

---

### Task 5: Extend `Rule` struct with Scope + Applies + CheckSubject

**Files:**
- Modify: `verify.go:119-126` (Rule struct)

- [ ] **Step 1:** Replace the `Rule` struct with:

```go
type Rule struct {
	ID          string // catalog-stable, e.g. "DISK-NOT-EMPTY"
	Severity    Severity
	Dialects    []Dialect // dialects the rule applies to; nil/empty = all
	Description string    // one-line summary, used in human output
	Citation    string    // file:line of the strongest evidence

	// Legacy single-shot check. Mutually exclusive with Scope + CheckSubject.
	// Kept for rules that haven't been migrated to the per-subject model.
	Check func(ctx *CheckContext) []Finding

	// Per-subject scope. If unset (zero value = DiskScope) AND CheckSubject
	// is nil, the rule is treated as legacy (Check is called once per Verify).
	Scope SubjectScope

	// Applies reports whether subject is eligible for this rule. If nil,
	// the rule applies to all subjects of its Scope. Counted as the
	// denominator for fail-rate analysis.
	Applies func(ctx *CheckContext, subject Subject) bool

	// CheckSubject evaluates one applicable subject. Returns nil for
	// pass, a Finding pointer for fail. The framework calls this once
	// per applicable subject and emits a CheckEvent for each call.
	CheckSubject func(ctx *CheckContext, subject Subject) *Finding
}
```

- [ ] **Step 2:** Verify it compiles.

```bash
go build ./...
```

Expected: compiles. Existing rules still register fine because all new fields are optional.

- [ ] **Step 3:** Commit.

```bash
git add verify.go
git commit -m "feat(verify): extend Rule struct with Scope/Applies/CheckSubject (back-compat)"
```

---

### Task 6: Framework iteration + CheckEvent emission

**Files:**
- Create: `check_event.go`
- Modify: `verify.go` (rewrite `(*DiskImage).Verify()`, add iteration helpers)

- [ ] **Step 1:** Create `check_event.go`:

```go
package samfile

// CheckEvent is the per-subject record emitted by Verify when a rule
// declares Scope + CheckSubject. One event per (rule, applicable
// subject) pair. Legacy rules (Check-only) emit one event per finding
// they produce, with Outcome=fail and a synthetic Subject ref.
type CheckEvent struct {
	Version int            `json:"v"`
	Disk    string         `json:"disk"`
	RuleID  string         `json:"rule_id"`
	Scope   string         `json:"scope"`
	Ref     string         `json:"ref"`
	Outcome string         `json:"outcome"` // pass | fail | not_applicable
	Attrs   map[string]any `json:"attrs"`
	Finding *Finding       `json:"finding,omitempty"`
}

// EventRecorder collects CheckEvents during a Verify run. The verify
// CLI installs a recorder when --format jsonl is set; for normal text
// output the recorder is nil and no events are kept.
type EventRecorder struct {
	Events []CheckEvent
	Disk   string
}

func (r *EventRecorder) Record(e CheckEvent) {
	if r == nil {
		return
	}
	e.Version = 1
	e.Disk = r.Disk
	r.Events = append(r.Events, e)
}
```

- [ ] **Step 2:** Modify `verify.go` `CheckContext` to add recorder + subject iterators. Append after the existing `CheckContext` struct (around line 164):

```go
// Recorder, if non-nil, receives a CheckEvent per (rule, applicable
// subject) when a Rule declares Scope + CheckSubject. The text-output
// path leaves this nil.
func (ctx *CheckContext) SetRecorder(r *EventRecorder) { ctx.recorder = r }

func (ctx *CheckContext) recorderRef() *EventRecorder { return ctx.recorder }
```

And add a field to the struct itself:

```go
type CheckContext struct {
	Disk     *DiskImage
	Journal  *DiskJournal
	Dialect  Dialect
	recorder *EventRecorder
}
```

- [ ] **Step 3:** Add `subjectsForScope`:

```go
// subjectsForScope enumerates every Subject of the given scope on
// the current disk. DiskScope yields one DiskSubject; SlotScope
// yields 80 SlotSubjects (every dir entry, used or erased);
// ChainStepScope yields one ChainStepSubject per sector in every
// used file's chain.
func (ctx *CheckContext) subjectsForScope(scope SubjectScope) []Subject {
	switch scope {
	case DiskScope:
		return []Subject{&DiskSubject{Journal: ctx.Journal, Disk: ctx.Disk, Dialect: ctx.Dialect}}
	case SlotScope:
		out := make([]Subject, 0, 80)
		for i := 0; i < 80; i++ {
			fe := ctx.Journal.SlotAt(i)
			out = append(out, &SlotSubject{SlotIndex: i, FileEntry: fe, Disk: ctx.Disk, Journal: ctx.Journal})
		}
		return out
	case ChainStepScope:
		var out []Subject
		for i := 0; i < 80; i++ {
			fe := ctx.Journal.SlotAt(i)
			if fe.Type == 0 {
				continue
			}
			chain := ctx.Journal.ChainForSlot(i)
			for idx, step := range chain {
				pos := "intermediate"
				if idx == 0 {
					pos = "first"
				}
				if idx == len(chain)-1 {
					if pos == "first" {
						pos = "first" // single-sector files
					} else {
						pos = "last"
					}
				}
				out = append(out, &ChainStepSubject{
					SlotIndex:  i,
					ChainIndex: idx,
					Track:      step.Track,
					Sector:     step.Sector,
					NextTrack:  step.NextTrack,
					NextSector: step.NextSector,
					Position:   pos,
					Disk:       ctx.Disk,
					Journal:    ctx.Journal,
				})
			}
		}
		return out
	}
	return nil
}
```

This task assumes `DiskJournal` has `SlotAt(int) *FileEntry` and `ChainForSlot(int) []ChainStep` helpers. **If they don't exist, defer this task and implement them first.** Verify:

```bash
grep -nE "func.*DiskJournal.*SlotAt|func.*DiskJournal.*ChainForSlot|func.*DiskJournal.*UsedSlots" *.go
```

If missing, this task expands to also create those helpers (read the existing journal struct to see what's already there).

- [ ] **Step 4:** Rewrite `Verify()`:

```go
func (di *DiskImage) Verify() VerifyReport {
	return di.verifyInternal(nil)
}

// VerifyWithRecorder runs verify and captures every Check event into
// the provided recorder (in addition to the existing VerifyReport).
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
					rec.Record(CheckEvent{RuleID: rule.ID, Scope: rule.Scope.String(), Ref: subj.Ref(), Outcome: "not_applicable", Attrs: subj.Attributes()})
					continue
				}
				finding := rule.CheckSubject(ctx, subj)
				outcome := "pass"
				if finding != nil {
					outcome = "fail"
					report.Findings = append(report.Findings, *finding)
				}
				rec.Record(CheckEvent{RuleID: rule.ID, Scope: rule.Scope.String(), Ref: subj.Ref(), Outcome: outcome, Attrs: subj.Attributes(), Finding: finding})
			}
			continue
		}
		// Legacy path: rule provides Check only.
		legacyFindings := rule.Check(ctx)
		report.Findings = append(report.Findings, legacyFindings...)
		for _, f := range legacyFindings {
			rec.Record(CheckEvent{RuleID: rule.ID, Scope: "legacy", Ref: f.Location.String(), Outcome: "fail", Finding: &f})
		}
	}
	return report
}
```

- [ ] **Step 5:** Run existing tests.

```bash
go test ./...
```

Expected: all pass — the new code path is only exercised when CheckSubject is set, and no rule has it set yet.

- [ ] **Step 6:** Commit.

```bash
git add check_event.go verify.go
git commit -m "feat(verify): add per-subject framework iteration + CheckEvent recording"
```

---

### Task 7: `--format jsonl` CLI flag

**Files:**
- Modify: `cmd/samfile/verify.go`

- [ ] **Step 1:** Read existing `cmd/samfile/verify.go` to find the flag-parsing block.

```bash
grep -n "flag\." cmd/samfile/verify.go | head -20
```

- [ ] **Step 2:** Add a `--format` flag (string, default `"text"`). After flag parsing, if `format == "jsonl"`:

```go
rec := &samfile.EventRecorder{Disk: filepath.Base(input)}
_ = di.VerifyWithRecorder(rec)
enc := json.NewEncoder(os.Stdout)
for _, ev := range rec.Events {
    if err := enc.Encode(ev); err != nil {
        return err
    }
}
return nil
```

Otherwise the existing text path runs unchanged.

- [ ] **Step 3:** Build and smoke-test.

```bash
go build -o /tmp/samfile-audit ./cmd/samfile
/tmp/samfile-audit verify --format jsonl -i testdata/ETrackerv1.2.mgt | head -5
```

Expected: 0+ JSONL lines on stdout (all legacy events for now — they're "fail" events from existing rules).

- [ ] **Step 4:** Commit.

```bash
git add cmd/samfile/verify.go
git commit -m "feat(samfile): add --format jsonl flag to verify command"
```

---

## Phase 2 — Rule migrations (one task per rule file)

Each task migrates one rule file from `Check func(ctx) []Finding` to `Scope + Applies + CheckSubject`. The existing tests must still pass after migration. Per-rule pattern: extract the iteration that's currently inside the rule body into `Applies` (returns true when the subject is eligible — typically a type/dialect/state check), and reshape `Check`'s body to operate on the single passed-in Subject, returning `*Finding`.

**Migration helper template:**

```go
Register(Rule{
    ID:          "EXAMPLE-RULE",
    Severity:    SeverityCosmetic,
    Description: "...",
    Citation:    "...",
    Scope:       SlotScope,
    Applies: func(ctx *CheckContext, subj Subject) bool {
        s := subj.(*SlotSubject)
        return s.FileEntry.Type == FT_CODE  // adjust per-rule
    },
    CheckSubject: func(ctx *CheckContext, subj Subject) *Finding {
        s := subj.(*SlotSubject)
        // ... existing per-slot check body ...
        if mismatch {
            return &Finding{RuleID: "EXAMPLE-RULE", Severity: SeverityCosmetic, Location: SlotLocation(s.SlotIndex, s.FileEntry.Name.String()), Message: ...}
        }
        return nil
    },
})
```

### Task 8: Migrate `rules_smoke.go` (1 rule)

- [ ] **Step 1:** Migrate `DISK-NOT-EMPTY` to `Scope: DiskScope, Applies: always-true, CheckSubject: per-disk`.
- [ ] **Step 2:** Run `go test ./...`. Expected: pass.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_smoke to per-subject framework`.

### Task 9: Migrate `rules_disk.go` (3 rules)

- [ ] **Step 1:** Migrate `DISK-DIRECTORY-TRACKS`, `DISK-TRACK-SIDE-ENCODING`, `DISK-SECTOR-RANGE`. All three are `DiskScope`. The two encoding/range rules iterate sectors internally — move that iteration to `ChainStepScope` if it makes sense, otherwise keep as `DiskScope` and walk internally.
- [ ] **Step 2:** Run `go test ./rules_disk_test.go ./...`. Expected: pass.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_disk to per-subject framework`.

### Task 10: Migrate `rules_directory.go` (9 rules)

- [ ] **Step 1:** Migrate each rule. Most are `SlotScope` with `Applies` filtering on slot-state (used vs erased).
- [ ] **Step 2:** Run tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_directory to per-subject framework`.

### Task 11: Migrate `rules_chain.go` (4 rules)

- [ ] **Step 1:** Migrate. `CHAIN-TERMINATOR-ZERO-ZERO`, `CHAIN-NO-CYCLE`, `CHAIN-MATCHES-SAM`: `ChainStepScope` makes sense for the first; the others are slot-scope (they reason about the whole chain). Pick the scope that produces the right denominator and use the chain-step or slot subject accordingly.
- [ ] **Step 2:** Tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_chain to per-subject framework`.

### Task 12: Migrate `rules_cross.go` (3 rules)

- [ ] **Step 1:** Migrate. These are inherently disk-scope (compare across slots).
- [ ] **Step 2:** Tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_cross to per-subject framework`.

### Task 13: Migrate `rules_body_header.go` (11 rules)

- [ ] **Step 1:** Migrate. All `SlotScope`. `Applies` filters: most need `Type != 0`; some FT_CODE-only (the EXEC-* rules already have this check inside).
- [ ] **Step 2:** Tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_body_header to per-subject framework`.

### Task 14: Migrate `rules_ft_code.go` (4 rules)

- [ ] **Step 1:** Migrate. All `SlotScope`, `Applies: Type == FT_CODE`.
- [ ] **Step 2:** Tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_ft_code to per-subject framework`.

### Task 15: Migrate `rules_ft_basic.go` (7 rules)

- [ ] **Step 1:** Migrate. All `SlotScope`, `Applies: Type == FT_SAM_BASIC`.
- [ ] **Step 2:** Tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_ft_basic to per-subject framework`.

### Task 16: Migrate `rules_ft_array.go` (1 rule)

- [ ] **Step 1:** Migrate. `SlotScope`, `Applies: Type == FT_NUMERIC_ARRAY || Type == FT_STRING_ARRAY`.
- [ ] **Step 2:** Tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_ft_array to per-subject framework`.

### Task 17: Migrate `rules_ft_screen.go` (2 rules)

- [ ] **Step 1:** Migrate. `SlotScope`, `Applies: Type == FT_SCREEN`.
- [ ] **Step 2:** Tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_ft_screen to per-subject framework`.

### Task 18: Migrate `rules_ft_zxsnap.go` (2 rules)

- [ ] **Step 1:** Migrate. `SlotScope`, `Applies: Type == FT_ZX_SNAP_48K`.
- [ ] **Step 2:** Tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_ft_zxsnap to per-subject framework`.

### Task 19: Migrate `rules_boot.go` (3 rules)

- [ ] **Step 1:** Migrate. `DiskScope`. `Applies: always-true`.
- [ ] **Step 2:** Tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_boot to per-subject framework`.

### Task 20: Migrate `rules_cosmetic.go` (1 rule)

- [ ] **Step 1:** Migrate. Probably `DiskScope` or `SlotScope` depending on the rule body.
- [ ] **Step 2:** Tests.
- [ ] **Step 3:** Commit: `refactor(verify): migrate rules_cosmetic to per-subject framework`.

---

## Phase 3 — Python pipeline

### Task 21: `tools/audit/ingest.py`

**Files:**
- Create: `tools/audit/ingest.py`

- [ ] **Step 1:** Create `tools/audit/ingest.py`:

```python
#!/usr/bin/env python3
"""Ingest JSONL CheckEvents from samfile --format jsonl into the
`checks` table of ~/sam-corpus/findings.db.

Reads: ~/sam-corpus/outputs-jsonl/<disk>.jsonl
Writes: ~/sam-corpus/findings.db (creates / replaces `checks` table).
Existing `disks` and `findings` tables are not touched.
"""

from __future__ import annotations

import json
import sqlite3
from pathlib import Path

CORPUS = Path.home() / "sam-corpus"
JSONL_DIR = CORPUS / "outputs-jsonl"
DB = CORPUS / "findings.db"

COLUMNS = [
    "disk", "rule_id", "scope", "ref", "outcome",
    # disk attrs
    "dialect", "boot_signature_present", "used_slot_count",
    # slot attrs
    "slot_index", "filename", "file_type", "file_type_byte",
    "file_length", "page_offset_form", "pages", "mgt_flags",
    "first_track", "first_sector", "first_side",
    "has_autorun_or_autoexec", "dir_mirror_populated",
    "slot_is_erased", "file_type_info_hex", "sectors_count",
    # chain-step attrs
    "chain_position", "chain_index", "track", "sector", "side",
    "next_track", "next_sector", "on_sam_map", "on_dir_sam_map",
    "distance_from_dir_tracks",
    # finding payload
    "severity", "message", "citation",
]


def column_value(event: dict, col: str):
    if col in ("disk", "rule_id", "scope", "ref", "outcome"):
        return event.get(col)
    attrs = event.get("attrs") or {}
    if col in attrs:
        v = attrs[col]
        if isinstance(v, bool):
            return int(v)
        return v
    finding = event.get("finding") or {}
    if col == "severity":
        return finding.get("Severity")
    if col == "message":
        return finding.get("Message")
    if col == "citation":
        return finding.get("Citation")
    return None


def main() -> None:
    conn = sqlite3.connect(DB)
    c = conn.cursor()
    c.execute("DROP TABLE IF EXISTS checks")
    column_decls = []
    for col in COLUMNS:
        # outcome / scope / file_type are TEXT; numerics stay INTEGER.
        if col in ("disk", "rule_id", "scope", "ref", "outcome",
                   "dialect", "filename", "file_type", "page_offset_form",
                   "file_type_info_hex", "chain_position",
                   "severity", "message", "citation"):
            column_decls.append(f"{col} TEXT")
        else:
            column_decls.append(f"{col} INTEGER")
    c.execute(f"CREATE TABLE checks ({', '.join(column_decls)})")
    c.execute("CREATE INDEX idx_checks_rule ON checks(rule_id)")
    c.execute("CREATE INDEX idx_checks_outcome ON checks(outcome)")
    c.execute("CREATE INDEX idx_checks_disk ON checks(disk)")
    c.execute("CREATE INDEX idx_checks_type ON checks(file_type)")
    placeholders = ", ".join("?" for _ in COLUMNS)
    insert = f"INSERT INTO checks ({', '.join(COLUMNS)}) VALUES ({placeholders})"
    n = 0
    for path in sorted(JSONL_DIR.glob("*.jsonl")):
        for line in path.read_text().splitlines():
            if not line.strip():
                continue
            event = json.loads(line)
            event.setdefault("disk", path.stem)
            row = [column_value(event, col) for col in COLUMNS]
            c.execute(insert, row)
            n += 1
    conn.commit()
    conn.close()
    print(f"ingested {n} CheckEvents into {DB}")


if __name__ == "__main__":
    main()
```

- [ ] **Step 2:** Commit.

```bash
git add tools/audit/ingest.py
git commit -m "feat(audit): add ingest.py — JSONL CheckEvents -> checks table"
```

---

### Task 22: `tools/audit/mine.py` — fallback-floor reports (coverage + disk-health)

**Files:**
- Create: `tools/audit/mine.py`

- [ ] **Step 1:** Create `tools/audit/mine.py` (initial version with the two fallback reports only):

```python
#!/usr/bin/env python3
"""Generate audit reports from ~/sam-corpus/findings.db's `checks` table.

Reports (written under ~/sam-corpus/analyses/):
  Fallback floor (plain counts):
    coverage.md     — per-rule applies/fails/fail-rate
    disk-health.md  — per-disk findings + structural-pass-rate
  Richer (best-effort, see Task 24):
    conditional.md
    disk-clusters.md
    patterns.md
"""

from __future__ import annotations

import sqlite3
from pathlib import Path

import pandas as pd

CORPUS = Path.home() / "sam-corpus"
DB = CORPUS / "findings.db"
OUT = CORPUS / "analyses"


def load_checks() -> pd.DataFrame:
    conn = sqlite3.connect(DB)
    df = pd.read_sql_query("SELECT * FROM checks", conn)
    conn.close()
    return df


def report_coverage(checks: pd.DataFrame) -> None:
    rows = []
    for rule_id, grp in checks.groupby("rule_id"):
        applicable = grp[grp["outcome"].isin(["pass", "fail"])]
        fails = grp[grp["outcome"] == "fail"]
        applies_n = len(applicable)
        fails_n = len(fails)
        disks = grp.loc[grp["outcome"] == "fail", "disk"].nunique()
        rate = (100.0 * fails_n / applies_n) if applies_n else 0.0
        sev = fails["severity"].mode().iloc[0] if not fails.empty else ""
        rows.append((rule_id, sev, applies_n, fails_n, rate, disks))
    df = pd.DataFrame(rows, columns=["rule_id", "severity", "applies", "fails", "fail_rate_pct", "disks_affected"])
    df = df.sort_values("fail_rate_pct", ascending=False)
    md = ["# Per-rule coverage and failure rate", ""]
    md.append("| Rule | Severity | Applies | Fails | Fail-rate | Disks |")
    md.append("|---|---|---:|---:|---:|---:|")
    for _, r in df.iterrows():
        md.append(f"| `{r.rule_id}` | {r.severity} | {r.applies} | {r.fails} | {r.fail_rate_pct:.1f}% | {r.disks_affected} |")
    (OUT / "coverage.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'coverage.md'} ({len(df)} rules)")


def report_disk_health(checks: pd.DataFrame) -> None:
    rows = []
    for disk, grp in checks.groupby("disk"):
        total_findings = (grp["outcome"] == "fail").sum()
        fatal = ((grp["outcome"] == "fail") & (grp["severity"] == "fatal")).sum()
        structural = ((grp["outcome"] == "fail") & (grp["severity"] == "structural")).sum()
        struct_checks = grp[(grp["outcome"].isin(["pass", "fail"])) & (grp["severity"].isin(["structural", "fatal"]) | grp["rule_id"].isin([]))]
        # Use the per-rule severity from the catalog, not the per-finding severity, for the denominator.
        # Workaround: count outcomes per-rule-severity by joining via rule severity table built from fails.
        # Simpler: use all applicable checks where the rule produced a fail with severity fatal/structural anywhere in the corpus.
        struct_pass = (struct_checks["outcome"] == "pass").sum()
        struct_fail = (struct_checks["outcome"] == "fail").sum()
        struct_rate = (struct_pass / (struct_pass + struct_fail)) if (struct_pass + struct_fail) else 1.0
        distinct_rules_fired = grp.loc[grp["outcome"] == "fail", "rule_id"].nunique()
        rows.append((disk, total_findings, fatal, structural, struct_rate, distinct_rules_fired))
    df = pd.DataFrame(rows, columns=["disk", "total_findings", "fatal", "structural", "structural_pass_rate", "distinct_rules_fired"])
    df = df.sort_values(["structural_pass_rate", "total_findings"], ascending=[True, False])
    md = ["# Per-disk health", ""]
    md.append("Sorted by structural pass-rate (worst first). Disks at the top are the candidate 'not really a disk' cluster.")
    md.append("")
    md.append("| Disk | Total findings | Fatal | Structural | Struct pass-rate | Distinct rules fired |")
    md.append("|---|---:|---:|---:|---:|---:|")
    for _, r in df.iterrows():
        md.append(f"| {r.disk[:60]} | {r.total_findings} | {r.fatal} | {r.structural} | {r.structural_pass_rate:.2f} | {r.distinct_rules_fired} |")
    (OUT / "disk-health.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'disk-health.md'} ({len(df)} disks)")


def main() -> None:
    OUT.mkdir(parents=True, exist_ok=True)
    checks = load_checks()
    print(f"loaded {len(checks)} CheckEvents")
    report_coverage(checks)
    report_disk_health(checks)


if __name__ == "__main__":
    main()
```

- [ ] **Step 2:** Commit.

```bash
git add tools/audit/mine.py
git commit -m "feat(audit): add mine.py with coverage + disk-health fallback reports"
```

---

### Task 23: `tools/audit/mine.py` — richer reports

**Files:**
- Modify: `tools/audit/mine.py`

- [ ] **Step 1:** Append richer reports to `mine.py` and add their calls in `main()`:

```python
def report_conditional(checks: pd.DataFrame) -> None:
    """Per-rule conditional fail-rate per attribute value vs baseline."""
    attr_cols = [c for c in checks.columns if c not in ("disk", "rule_id", "scope", "ref", "outcome", "severity", "message", "citation")]
    rows = []
    for rule_id, grp in checks.groupby("rule_id"):
        applicable = grp[grp["outcome"].isin(["pass", "fail"])]
        if len(applicable) < 10:
            continue
        baseline = (applicable["outcome"] == "fail").mean()
        for col in attr_cols:
            if applicable[col].isna().all():
                continue
            for val, sub in applicable.groupby(col):
                if len(sub) < 10:
                    continue
                cond = (sub["outcome"] == "fail").mean()
                # surface only when conditional rate deviates strongly
                if (cond in (0.0, 1.0)) or abs(cond - baseline) > 0.3:
                    rows.append((rule_id, col, str(val), len(sub), baseline, cond, cond - baseline))
    df = pd.DataFrame(rows, columns=["rule_id", "attribute", "value", "support", "baseline_fail_rate", "conditional_fail_rate", "delta"])
    df = df.sort_values(["rule_id", "delta"], ascending=[True, False])
    md = ["# Conditional failure rates per attribute", ""]
    md.append("Surfaces (rule, attribute, value) combinations where the conditional fail-rate differs substantially from the rule's baseline. |delta| > 0.3 or conditional rate ∈ {0%, 100%}, support ≥ 10.")
    md.append("")
    md.append("| Rule | Attribute | Value | Support | Baseline | Conditional | Δ |")
    md.append("|---|---|---|---:|---:|---:|---:|")
    for _, r in df.iterrows():
        md.append(f"| `{r.rule_id}` | {r.attribute} | {r.value} | {r.support} | {r.baseline_fail_rate*100:.1f}% | {r.conditional_fail_rate*100:.1f}% | {r.delta*100:+.1f}pp |")
    (OUT / "conditional.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'conditional.md'} ({len(df)} rows)")


def report_disk_clusters(checks: pd.DataFrame) -> None:
    """Hierarchical clustering of disks by rule-fire vector."""
    try:
        from sklearn.cluster import AgglomerativeClustering
    except ImportError:
        (OUT / "disk-clusters.md").write_text("# Disk clusters\n\nsklearn not installed; skipping.\n")
        return
    pivot = (checks[checks["outcome"] == "fail"]
             .pivot_table(index="disk", columns="rule_id", values="ref", aggfunc="count", fill_value=0))
    if pivot.empty:
        (OUT / "disk-clusters.md").write_text("# Disk clusters\n\nNo failure events to cluster.\n")
        return
    n_clusters = min(8, max(2, len(pivot) // 100))
    cl = AgglomerativeClustering(n_clusters=n_clusters)
    pivot["cluster"] = cl.fit_predict((pivot > 0).astype(int).values)
    md = ["# Disk clusters by rule-fire pattern", ""]
    for cid, grp in pivot.groupby("cluster"):
        top_rules = (grp.drop(columns="cluster") > 0).sum().sort_values(ascending=False).head(5)
        md.append(f"## Cluster {cid} — {len(grp)} disks")
        md.append("")
        md.append("Most-fired rules in this cluster:")
        for rule, n in top_rules.items():
            md.append(f"- `{rule}` ({n} of {len(grp)} disks)")
        md.append("")
        md.append("Example disks:")
        for disk in list(grp.index)[:5]:
            md.append(f"- {disk[:60]}")
        md.append("")
    (OUT / "disk-clusters.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'disk-clusters.md'} ({n_clusters} clusters)")


def report_patterns(checks: pd.DataFrame) -> None:
    """Per-rule decision-tree splits."""
    try:
        from sklearn.tree import DecisionTreeClassifier, export_text
    except ImportError:
        (OUT / "patterns.md").write_text("# Patterns\n\nsklearn not installed; skipping.\n")
        return
    attr_cols = [c for c in checks.columns if c not in ("disk", "rule_id", "scope", "ref", "outcome", "severity", "message", "citation", "filename", "file_type_info_hex", "ref")]
    md = ["# Per-rule patterns (decision-tree splits)", ""]
    for rule_id, grp in checks.groupby("rule_id"):
        applicable = grp[grp["outcome"].isin(["pass", "fail"])].copy()
        if len(applicable) < 50:
            continue
        y = (applicable["outcome"] == "fail").astype(int)
        if y.nunique() < 2:
            continue
        X = applicable[attr_cols].copy()
        # Coerce categoricals to category codes; coerce booleans to int.
        for c in X.columns:
            if X[c].dtype == object:
                X[c] = X[c].astype("category").cat.codes
            else:
                X[c] = X[c].fillna(-1)
        try:
            tree = DecisionTreeClassifier(max_depth=4, min_samples_leaf=10).fit(X, y)
        except ValueError:
            continue
        md.append(f"## `{rule_id}` (applies={len(applicable)}, fail-rate={y.mean()*100:.1f}%)")
        md.append("")
        md.append("```")
        md.append(export_text(tree, feature_names=list(X.columns), max_depth=4).strip())
        md.append("```")
        md.append("")
    (OUT / "patterns.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'patterns.md'}")
```

And update `main()`:

```python
def main() -> None:
    OUT.mkdir(parents=True, exist_ok=True)
    checks = load_checks()
    print(f"loaded {len(checks)} CheckEvents")
    report_coverage(checks)
    report_disk_health(checks)
    try:
        report_conditional(checks)
    except Exception as e:
        print(f"conditional.md failed: {e}")
    try:
        report_disk_clusters(checks)
    except Exception as e:
        print(f"disk-clusters.md failed: {e}")
    try:
        report_patterns(checks)
    except Exception as e:
        print(f"patterns.md failed: {e}")
```

- [ ] **Step 2:** Commit.

```bash
git add tools/audit/mine.py
git commit -m "feat(audit): add conditional / disk-clusters / patterns reports (best-effort)"
```

---

### Task 24: `tools/audit/run_audit.sh`

**Files:**
- Create: `tools/audit/run_audit.sh`

- [ ] **Step 1:** Create:

```bash
#!/bin/bash
# Rebuild samfile-verify, regenerate JSONL across the corpus,
# ingest, and mine. Idempotent.
set -euo pipefail

REPO="${HOME}/git/samfile"
CORPUS="${HOME}/sam-corpus"

cd "$REPO"
go build -o "$CORPUS/samfile-audit" ./cmd/samfile

mkdir -p "$CORPUS/outputs-jsonl" "$CORPUS/analyses"
rm -f "$CORPUS/outputs-jsonl"/*.jsonl

count=0
for disk in "$CORPUS/disks"/*.mgt; do
    [ -f "$disk" ] || continue
    name=$(basename "$disk" .mgt)
    "$CORPUS/samfile-audit" verify --format jsonl -i "$disk" \
        > "$CORPUS/outputs-jsonl/$name.jsonl" 2>/dev/null || true
    count=$((count + 1))
done
echo "ran samfile-audit on $count disks"

python3 "$REPO/tools/audit/ingest.py"
python3 "$REPO/tools/audit/mine.py"
echo "reports in $CORPUS/analyses/"
```

- [ ] **Step 2:** Mark executable.

```bash
chmod +x tools/audit/run_audit.sh
```

- [ ] **Step 3:** Commit.

```bash
git add tools/audit/run_audit.sh
git commit -m "feat(audit): add run_audit.sh end-to-end driver"
```

---

## Phase 4 — End-to-end run + autonomous discovery loop

### Task 25: End-to-end smoke test on a few disks

- [ ] **Step 1:** Build and run on 3 sample disks:

```bash
go build -o /tmp/samfile-audit ./cmd/samfile
mkdir -p /tmp/audit-smoke
for disk in ~/sam-corpus/disks/*.mgt; do
    /tmp/samfile-audit verify --format jsonl -i "$disk" > "/tmp/audit-smoke/$(basename "$disk" .mgt).jsonl" 2>/dev/null
    break
done
ls -la /tmp/audit-smoke/
head -3 /tmp/audit-smoke/*.jsonl
```

Expected: at least one .jsonl file with valid JSON objects per line.

- [ ] **Step 2:** Sanity-check that pass + fail + not_applicable counts make sense for one rule:

```bash
jq -s 'group_by(.rule_id) | map({rule:.[0].rule_id, total:length, fail:map(select(.outcome=="fail"))|length, pass:map(select(.outcome=="pass"))|length, na:map(select(.outcome=="not_applicable"))|length})' /tmp/audit-smoke/*.jsonl | head -40
```

Expected: per-rule pass + fail + n/a counts.

### Task 26: Full corpus run

- [ ] **Step 1:**

```bash
~/git/samfile/tools/audit/run_audit.sh
```

Expected: ~800 .jsonl files in `~/sam-corpus/outputs-jsonl/`, `checks` table populated, five Markdown reports in `~/sam-corpus/analyses/`.

- [ ] **Step 2:** Eyeball `coverage.md` and `disk-health.md`. Sanity checks:
  - High-fail-rate rules at the top: are they known cases? (e.g. for rules already known to be over-strict, do they appear?)
  - Worst-health disks: do the same ~10 disks fire across most rules? That's the "broken disk cluster".

### Task 27: Autonomous discovery + fix loop

This is the iterative loop. For each pass:

- [ ] **Step 1:** Read `coverage.md` + `conditional.md` + `patterns.md`. Identify high-confidence patterns per the spec's threshold:
  - (statistical): conditional fail-rate ≥ 80% on support ≥ 10 distinct disks, OR conditional fail-rate = 0% on support ≥ 10, OR Apriori confidence=1.0 support ≥ 10
  - (citation): a line range in `~/git/samdos/src/*.s`, ROM disasm, or samfile that explains the pattern

- [ ] **Step 2:** For each high-confidence pattern:
  - Form a hypothesis ("rule X mis-fires when attribute Y because rule Z is over-strict / under-strict")
  - Open the cited source range and confirm the explanation
  - If confirmed: patch the rule (catalog + verify code + tests). Commit with the citation in the commit message.
  - If refuted: discard.
  - If ambiguous: add to `~/sam-corpus/analyses/needs-human.md`.

- [ ] **Step 3:** After acting on all high-confidence patterns, re-run:

```bash
~/git/samfile/tools/audit/run_audit.sh
```

- [ ] **Step 4:** If a new pass surfaces no new high-confidence patterns, stop. Otherwise return to Step 1.

### Task 28: PR

- [ ] **Step 1:** Push the branch.

```bash
git push -u origin feat/verify-audit-framework
```

- [ ] **Step 2:** Create draft PR.

```bash
gh pr create --draft --title "feat(verify): per-subject audit framework + corpus-driven rule fixes" --body "$(cat <<'EOF'
## Summary

Adds per-subject Check-event instrumentation to the verify pipeline, ingests events into `~/sam-corpus/findings.db`, and uses the resulting `checks` table to discover patterns that drive rule-fix decisions. Five Markdown reports produced under `~/sam-corpus/analyses/`. All rule fixes in this PR are grounded by source citations from SAMDOS / ROM disasm / samfile.

Design spec: `docs/specs/2026-05-12-verify-audit-framework-design.md`.
Implementation plan: `docs/plans/2026-05-12-verify-audit-framework.md`.

## Framework changes

- `subject.go`, `subject_disk.go`, `subject_slot.go`, `subject_chain.go` — Subject interface + three scope-specific implementations.
- `verify.go` — extended `Rule` struct with `Scope` + `Applies` + `CheckSubject`; framework iterates subjects per scope and records `CheckEvent`s.
- `check_event.go` — CheckEvent + EventRecorder.
- `cmd/samfile/verify.go` — new `--format jsonl` flag.
- All N rules migrated from `Check func(ctx) []Finding` to `Scope + Applies + CheckSubject`. Existing tests pass.

## Python pipeline

- `tools/audit/ingest.py` — JSONL → `checks` SQLite table.
- `tools/audit/mine.py` — five reports (coverage, disk-health, conditional, disk-clusters, patterns).
- `tools/audit/run_audit.sh` — end-to-end driver.

## Rule fixes

(Filled in by the autonomous loop; each commit cites the source range that grounded the fix.)

## Test plan

- [ ] `go test ./...` passes
- [ ] `samfile verify --format jsonl -i testdata/ETrackerv1.2.mgt` emits valid JSONL
- [ ] `tools/audit/run_audit.sh` runs end-to-end on the local corpus
- [ ] Reports under `~/sam-corpus/analyses/` look coherent
- [ ] CI green

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 3:** Monitor CI until all checks complete. Fix failures autonomously (amend the responsible commit, force-push). Leave the PR as draft for Pete's review.

---

## Self-review

- All tasks reference exact file paths.
- All code-touching steps include the code to write.
- Each task ends in a commit.
- Migration tasks group rules by file (review unit = one file's worth of rules).
- The autonomous loop (Task 27) is explicit about when to stop and what counts as "act-on" confidence.
- The plan covers every section of the spec: Subject schema (Tasks 1-4), Rule struct + framework (Tasks 5-7), JSONL CLI (Task 7), migration (Tasks 8-20), Python ingest (Task 21), reports including fallback floor (Tasks 22-23), driver (Task 24), end-to-end (Tasks 25-26), discovery/fix loop (Task 27), PR (Task 28).
