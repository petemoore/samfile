# Verify Phase 3 — Disk, Directory, Chain & Cross-Entry Rules

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement 19 of the catalog's "structural-not-file-type-specific" rules so `Verify` produces real findings on real disks. After this phase, `samfile verify` is genuinely useful for the inspector / archivist audience even though file-type rules (FT_CODE, FT_SAM_BASIC, …) still come later.

**Architecture:** Each rule is a `samfile.Rule` registered at package `init()` time via `Register`. Rules are grouped by catalog section across four new files (`rules_disk.go`, `rules_directory.go`, `rules_chain.go`, `rules_cross.go`), each paired with a `*_test.go`. A private `walkChain` helper in `rules_chain.go` is shared by chain and cross-entry rules so each rule's Check function stays focused. No changes to Phase 1's registry plumbing or Phase 2's `DetectDialect`.

**Tech Stack:** Go 1.22+, standard library only. Existing `samfile` API surface (`DiskJournal`, `FileEntry`, `Sector`, `SectorAddressMap`, `SectorData.FilePart`, `Filename.String`) covers everything.

**Context for the engineer:**

Read these first, in order:

1. `docs/specs/2026-05-11-verify-feature-design.md` §"Implementation order" (Phase 3 of 6): this phase exercises the foundation against ~22 structural rules before the file-type rules in Phases 4–6 add another ~45.
2. `docs/disk-validity-rules.md` §1 (Disk-level), §2 (Directory-entry), §3 (Sector-chain), §4 (Cross-entry), §15 (CHAIN-SECTOR-COUNT-MINIMAL). The catalog gives every rule's What / Severity / Source authority / Citation / Test sketch.
3. `verify.go` — Phase 1's `Rule` / `CheckContext` / `Finding` / `Register` plumbing. Each rule below plugs into this exactly the same way the Phase-1 smoke rule does.
4. `rules_smoke.go` — the canonical "one rule, one Check function, one `init() { Register(...) }`" pattern this phase follows.
5. `samfile.go:80-217` — `FileEntry`, `SectorAddressMap`, `Sector` structs. `samfile.go:731-770` — the chain-walk pattern (you'll generalise it into `walkChain` in Task 4).

**Phase 3 scope: 19 rules.** Counts in the catalog look higher because several entries are preconditions or parser invariants that can't fail at Verify time. They are explicitly **deferred** below; the plan body explains why for each.

| Catalog rule | Phase 3 status |
|---|---|
| DISK-IMAGE-SIZE | DEFER — `Load` zero-pads / truncates to 819,200 (`samfile.go:362-371`); never reaches Verify. |
| DISK-NOT-EDSK | DEFER — `Load` rejects EDSK before Verify is reachable (`samfile.go:355-368`). |
| DISK-LAYOUT-CYL-INTERLEAVED | DEFER — precondition, not a check (catalog says so). |
| DIR-SLOT-COUNT | DEFER — `DiskJournal()` always returns 80 entries (`samfile.go:438-446`); not falsifiable post-parse. |
| DIR-TYPE-MASKING | DEFER — precondition, not a check (catalog says so). |
| DIR-SECTORS-BIG-ENDIAN | DEFER — parser invariant; `FileEntryFrom` always reads BE. Not falsifiable post-parse. |
| CHAIN-LINK-AT-510-511 | DEFER — precondition, not a check. |
| CHAIN-FIRST-MATCHES-DIR | DEFER — tautology post-parse: `samfile.File` starts the walk at `fe.FirstSector`. To falsify it we'd need to compare raw dir-bytes 0x0D-0x0E to a separately-extracted first sector, which is circular. Skip. |

The 19 rules that **are** in scope (3 + 9 + 3 + 3 + 1 = 19):

- §1: DISK-DIRECTORY-TRACKS, DISK-TRACK-SIDE-ENCODING, DISK-SECTOR-RANGE (3)
- §2: DIR-TYPE-BYTE-IS-KNOWN, DIR-ERASED-IS-ZERO, DIR-NAME-PADDING, DIR-NAME-NOT-EMPTY, DIR-FIRST-SECTOR-VALID, DIR-SECTORS-MATCHES-CHAIN, DIR-SECTORS-MATCHES-MAP, DIR-SECTORS-NONZERO, DIR-SAM-WITHIN-CAPACITY (9)
- §3: CHAIN-TERMINATOR-ZERO-ZERO, CHAIN-NO-CYCLE, CHAIN-MATCHES-SAM (3)
- §4: CROSS-NO-SECTOR-OVERLAP, CROSS-NO-DUPLICATE-NAMES, CROSS-DIRECTORY-AREA-UNUSED (3)
- §15: CHAIN-SECTOR-COUNT-MINIMAL (1)

After Task 6 the registry holds 20 rules total (Phase-1 smoke + 19 from Phase 3).

**Phase 3 standing rules:**

- Use `g` (the user's alias) not plain `git` for commits — it preserves authorship timestamps.
- Every rule's `Citation` field cites a real `file:line` location (per Pete's "samfile PRs must cite sources" rule). The exact citations are pre-filled in each rule block below; copy them verbatim.
- Test fabrication uses the inline pattern from Phase 2 (`NewDiskImage` + `AddCodeFile` + targeted byte patches + `WriteFileEntry`). No helper packages, no committed test disks beyond `testdata/ETrackerv1.2.mgt`.
- Each rule ships with two unit tests (positive: clean disk → 0 findings; negative: one targeted byte flip → exactly 1 finding with the right RuleID and severity).
- Draft PR only. The PR-creation step is explicit in Task 9; do not run it earlier.

---

## File Structure

| Path | Action | Responsibility |
|---|---|---|
| `rules_disk.go` | Create | §1 disk-level rules: 3 rules covering valid track/sector ranges across all link points. |
| `rules_disk_test.go` | Create | Positive + negative tests for §1 rules. |
| `rules_directory.go` | Create | §2 directory-entry rules: 9 rules covering name padding, type byte, sector count consistency. |
| `rules_directory_test.go` | Create | Positive + negative tests for §2 rules. |
| `rules_chain.go` | Create | §3 chain rules: 3 rules plus the `walkChain` helper used by chain and cross-entry rules. Also houses §15 CHAIN-SECTOR-COUNT-MINIMAL. |
| `rules_chain_test.go` | Create | Positive + negative tests for §3 + §15 rules, plus unit tests for `walkChain` itself. |
| `rules_cross.go` | Create | §4 cross-entry rules: 3 rules that compare across slots. |
| `rules_cross_test.go` | Create | Positive + negative tests for §4 rules. |
| `verify_test.go` | Modify | Add a `TestVerifyOnTestdataCorpus` integration test asserting Verify returns a populated, panic-free report on `testdata/ETrackerv1.2.mgt`. |

Phase 1's `verify.go` and Phase 2's `dialect.go` are **not** modified.

---

## The Rule Template (read this once, apply to every rule)

Every rule in this phase follows the same shape. Skim this template once; the per-rule tasks reference back to it.

```go
// In rules_<section>.go:

func init() {
    Register(Rule{
        ID:          "RULE-ID",             // catalog-stable, UPPER-KEBAB
        Severity:    SeverityXxx,           // from catalog
        Dialects:    nil,                   // nil = all dialects (Phase 3 rules are dialect-agnostic)
        Description: "one-line summary",    // matches catalog's "What" field, paraphrased
        Citation:    "file:line",           // copied verbatim from this plan
        Check:       checkRuleId,
    })
}

func checkRuleId(ctx *CheckContext) []Finding {
    var findings []Finding
    // ... iterate ctx.Journal / ctx.Disk as needed
    // ... append to findings on violation
    return findings // nil is fine if no findings
}
```

```go
// In rules_<section>_test.go:

func TestRuleIdPositive(t *testing.T) {
    di := NewDiskImage()
    // ... build a clean disk
    findings := checkRuleId(&CheckContext{
        Disk: di, Journal: di.DiskJournal(), Dialect: DetectDialect(di),
    })
    if len(findings) != 0 {
        t.Errorf("checkRuleId on clean disk returned %d findings; want 0", len(findings))
    }
}

func TestRuleIdNegative(t *testing.T) {
    di := NewDiskImage()
    // ... build a disk that deliberately violates this one rule
    findings := checkRuleId(&CheckContext{
        Disk: di, Journal: di.DiskJournal(), Dialect: DetectDialect(di),
    })
    if len(findings) != 1 {
        t.Fatalf("checkRuleId on bad disk returned %d findings; want 1", len(findings))
    }
    if findings[0].RuleID != "RULE-ID" {
        t.Errorf("RuleID = %q; want %q", findings[0].RuleID, "RULE-ID")
    }
    if findings[0].Severity != SeverityXxx {
        t.Errorf("Severity = %v; want %v", findings[0].Severity, SeverityXxx)
    }
    // Optionally also assert findings[0].Location is what the rule should produce.
}
```

Two conventions for the negative tests:

1. **One targeted byte flip per test.** Don't combine multiple violations in one fixture — a different rule may also fire and the assertion `len(findings) != 1` would catch the wrong condition.
2. **Mutate via `dj[slot].Field = …` then `di.WriteFileEntry(dj, slot)`.** This was the Phase-2 pattern and it works for every dir-entry field. Raw-byte mutation (`di[offset] = …`) is needed only when patching sector payload bytes — call out where this is used.

---

## Task 1: Create empty file skeletons + registry-growth smoke test

**Why this task exists:** lock in the four new files before any rule lands, and add one assertion that the registry actually reaches 20 entries once Phase 3 is complete. The skeleton commit is small and reviewable; subsequent commits then each add a coherent batch of rules.

**Files:**
- Create: `rules_disk.go`, `rules_directory.go`, `rules_chain.go`, `rules_cross.go` (each: package declaration only, no rules yet).
- Modify: `rules_smoke_test.go` — add a registry-growth assertion that pins the expected rule count for Phase 3.

- [ ] **Step 1: Create the four rule-file skeletons**

Each file gets exactly this content (substitute the section comment):

```go
// rules_disk.go
package samfile

// §1 Disk-level rules (catalog docs/disk-validity-rules.md §1).
// Rules in this file check that every track and sector reference on
// disk lies within the documented MGT geometry. They apply to all
// dialects.
```

```go
// rules_directory.go
package samfile

// §2 Directory-entry rules (catalog docs/disk-validity-rules.md §2).
// Rules in this file check internal consistency of each of the 80
// directory entries: type byte, filename padding, sector count vs
// chain length vs SectorAddressMap popcount. They apply to all
// dialects.
```

```go
// rules_chain.go
package samfile

// §3 Sector-chain rules + §15 CHAIN-SECTOR-COUNT-MINIMAL (catalog
// docs/disk-validity-rules.md §3 + §15). Rules in this file walk
// each used file's sector chain and check link integrity, cycle
// freedom, and consistency with the SectorAddressMap. They apply
// to all dialects.
//
// walkChain (private) is shared with rules_cross.go via the same
// package; it is the single canonical chain-walker for Phase 3
// rules so per-rule walking stays simple.
```

```go
// rules_cross.go
package samfile

// §4 Cross-entry consistency rules (catalog docs/disk-validity-rules.md
// §4). Rules in this file compare data across multiple directory
// slots: shared sectors, duplicate names, references into the
// directory area. They apply to all dialects.
```

- [ ] **Step 2: Pin expected registry count**

After this phase is fully implemented, the registry will hold the Phase-1 smoke rule plus 19 Phase-3 rules = 20 entries. The catalog's order is fixed; the registration order matches the catalog's section order.

Add this test to `rules_smoke_test.go` (append at the end):

```go
// TestPhase3RegistryGrowth pins the expected rule count once Phase 3
// is fully implemented. It will fail in Task 1 (only 1 rule registered)
// and pass once Tasks 2-6 land the remaining 19 rules. This is a
// regression gate: if any rule is accidentally removed or never
// registered, this test fails.
func TestPhase3RegistryGrowth(t *testing.T) {
    if got := len(Rules()); got != 20 {
        t.Errorf("len(Rules()) = %d; want 20 (1 smoke + 19 phase-3 rules)", got)
    }
}
```

- [ ] **Step 3: Verify the skeleton compiles**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && go build ./...`
Expected: no output, exit 0.

- [ ] **Step 4: Verify the new test fails as expected**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && go test -run TestPhase3RegistryGrowth -v ./...`
Expected: FAIL with `len(Rules()) = 1; want 20`. (This is the regression gate working — it'll pass once the rules land.)

Other tests must still pass; run the full suite to confirm nothing else regresses:

```
go test ./...
```

Expected: only `TestPhase3RegistryGrowth` fails; everything else green.

- [ ] **Step 5: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-3 && \
g add rules_disk.go rules_directory.go rules_chain.go rules_cross.go rules_smoke_test.go && \
g commit -m "verify: phase 3 skeleton (four new rule files + registry gate)

Adds empty rules_{disk,directory,chain,cross}.go for the §1/§2/§3/§4
catalog sections plus a TestPhase3RegistryGrowth test that pins the
final rule count at 20 (1 smoke + 19 phase-3 rules).

The test deliberately fails after this commit; it turns green
incrementally as tasks 2-6 register their rules. This is the
regression gate that catches a rule accidentally never being
registered."
```

---

## Task 2: §1 Disk-level rules (3 rules)

**Why this task exists:** §1 is the smallest section and gives the implementer a chance to internalise the Rule template before tackling the larger directory-entry batch. All three rules iterate the same thing (every track byte / sector byte found on disk) but check different invariants.

**Files:**
- Modify: `rules_disk.go` — append three rules.
- Modify: `rules_disk_test.go` — create, with two tests per rule (six total).

**Helper: enumerating "every track byte on disk".** Three sources of track/sector references:

1. `fe.FirstSector.Track` and `.Sector` for every used dir entry.
2. Each sector's payload byte 510 (track) and 511 (sector) — the next-link.
3. Bits set in each used slot's `SectorAddressMap` — these are pre-validated by the `SAMMask` formula and don't need separate range checks. Skip.

For sources 1 and 2 the per-rule check is "iterate all `(track, sector)` references, flag invalid". Rather than duplicate the iteration across three rules, write one shared private helper in `rules_disk.go`:

```go
// trackSectorRefs returns every (track, sector) link reachable from
// ctx — first-sector references from used dir entries plus the
// next-link bytes (510-511) of every sector in every used file's
// chain. Bounded by the disk's 1560-sector capacity per chain so a
// cyclic or truncated chain cannot hang the iteration. Used by
// DISK-DIRECTORY-TRACKS, DISK-TRACK-SIDE-ENCODING, DISK-SECTOR-RANGE.
//
// Each returned ref carries enough context for a Finding's Location
// (slot index, slot name, the sector itself, and the byte offset
// within the sector where the link byte lives — 0 for a first-sector
// reference, 510 for a chain link's track byte, 511 for sector byte).
//
// Errors from SectorData (only fire on out-of-range raw track values
// that bypass the dir entry's parse path) are silently ignored —
// DISK-TRACK-SIDE-ENCODING will catch them via the dir entry's own
// first-sector reference.
type sectorRef struct {
    Slot       int
    Filename   string
    Sector     Sector // copy (not pointer) so the value is independent of any pool
    ByteOffset int    // 0 (first-sector) or 510 (chain link track) or 511 (chain link sector)
    IsTerminator bool // true when this ref is the (0, 0) chain terminator — skip range checks
}

func trackSectorRefs(ctx *CheckContext) []sectorRef {
    var refs []sectorRef
    for _, slot := range ctx.Journal.UsedFileEntries() {
        fe := ctx.Journal[slot]
        name := fe.Name.String()
        // First-sector reference from the dir entry.
        refs = append(refs, sectorRef{Slot: slot, Filename: name, Sector: *fe.FirstSector, ByteOffset: 0})
        // Walk the chain. Bound by 1560 (disk capacity) to defend
        // against cycles / missing terminators; CHAIN-NO-CYCLE will
        // also catch those.
        cur := fe.FirstSector
        for steps := 0; steps < 1560; steps++ {
            sd, err := ctx.Disk.SectorData(cur)
            if err != nil {
                break
            }
            fp := sd.FilePart()
            nextSec := *fp.NextSector
            isTerm := nextSec.Track == 0 && nextSec.Sector == 0
            refs = append(refs,
                sectorRef{Slot: slot, Filename: name, Sector: nextSec, ByteOffset: 510, IsTerminator: isTerm},
                sectorRef{Slot: slot, Filename: name, Sector: nextSec, ByteOffset: 511, IsTerminator: isTerm},
            )
            if isTerm {
                break
            }
            cur = fp.NextSector
        }
    }
    return refs
}
```

The `(0, 0)` terminator is a special case — it IS a valid chain terminator but its track and sector bytes are 0, which is technically "invalid" for a data sector. The `IsTerminator` flag lets each rule decide whether to skip it.

Now the three rules.

### Rule 1: DISK-DIRECTORY-TRACKS

Catalog: every used `FirstSector.Track` is in `{4..79, 128..207}`. Tracks 0-3 (side 0) hold the directory; no file's first sector can land there.

```go
func init() {
    Register(Rule{
        ID:          "DISK-DIRECTORY-TRACKS",
        Severity:    SeverityStructural,
        Description: "no file references a sector in the directory area (tracks 0-3 of side 0)",
        Citation:    "sam-coupe_tech-man_v3-0.txt:4340-4343",
        Check:       checkDiskDirectoryTracks,
    })
}

func checkDiskDirectoryTracks(ctx *CheckContext) []Finding {
    var findings []Finding
    for _, ref := range trackSectorRefs(ctx) {
        if ref.IsTerminator {
            continue // (0, 0) terminator is allowed even though Track=0 is in [0..3]
        }
        if (ref.Sector.Track & 0x7F) < 4 {
            findings = append(findings, Finding{
                RuleID:   "DISK-DIRECTORY-TRACKS",
                Severity: SeverityStructural,
                Location: SectorLocation(ref.Slot, ref.Filename, &ref.Sector, ref.ByteOffset),
                Message:  fmt.Sprintf("track 0x%02x references the directory area (tracks 0-3 of side 0)", ref.Sector.Track),
                Citation: "sam-coupe_tech-man_v3-0.txt:4340-4343",
            })
        }
    }
    return findings
}
```

(`&ref.Sector` captures the iteration variable's local copy because `Sector` in `sectorRef` is by-value, not by-pointer; the address is stable within the loop body.)

### Rule 2: DISK-TRACK-SIDE-ENCODING

Catalog: valid track byte ranges are `0x00..0x4F` (side 0 cylinders 0-79) and `0x80..0xCF` (side 1 cylinders 0-79). `0x50..0x7F` and `0xD0..0xFF` are invalid.

```go
func init() {
    Register(Rule{
        ID:          "DISK-TRACK-SIDE-ENCODING",
        Severity:    SeverityFatal,
        Description: "every track byte references a physical cylinder 0-79 on side 0 or side 1",
        Citation:    "samfile.go:393-394",
        Check:       checkDiskTrackSideEncoding,
    })
}

func checkDiskTrackSideEncoding(ctx *CheckContext) []Finding {
    var findings []Finding
    for _, ref := range trackSectorRefs(ctx) {
        if ref.ByteOffset == 511 {
            continue // sector-number byte, not the track byte
        }
        if ref.IsTerminator {
            continue
        }
        t := ref.Sector.Track
        if (t >= 80 && t < 128) || t >= 208 {
            findings = append(findings, Finding{
                RuleID:   "DISK-TRACK-SIDE-ENCODING",
                Severity: SeverityFatal,
                Location: SectorLocation(ref.Slot, ref.Filename, &ref.Sector, ref.ByteOffset),
                Message:  fmt.Sprintf("track 0x%02x is in the invalid range (valid: 0x00-0x4F or 0x80-0xCF)", t),
                Citation: "samfile.go:393-394",
            })
        }
    }
    return findings
}
```

### Rule 3: DISK-SECTOR-RANGE

Catalog: every sector byte in a live link is 1..10. The `(0, 0)` chain terminator is allowed.

```go
func init() {
    Register(Rule{
        ID:          "DISK-SECTOR-RANGE",
        Severity:    SeverityFatal,
        Description: "every sector number is in range 1-10 (or 0 for the chain terminator)",
        Citation:    "samfile.go:389-392",
        Check:       checkDiskSectorRange,
    })
}

func checkDiskSectorRange(ctx *CheckContext) []Finding {
    var findings []Finding
    for _, ref := range trackSectorRefs(ctx) {
        if ref.ByteOffset == 510 {
            continue // track byte, not the sector byte
        }
        if ref.IsTerminator {
            continue
        }
        s := ref.Sector.Sector
        if s < 1 || s > 10 {
            findings = append(findings, Finding{
                RuleID:   "DISK-SECTOR-RANGE",
                Severity: SeverityFatal,
                Location: SectorLocation(ref.Slot, ref.Filename, &ref.Sector, ref.ByteOffset),
                Message:  fmt.Sprintf("sector 0x%02x is out of range (valid: 1-10)", s),
                Citation: "samfile.go:389-392",
            })
        }
    }
    return findings
}
```

### Tests

Create `rules_disk_test.go`:

```go
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
```

- [ ] **Step 1: Add the three rules and tests**

Apply the four changes above: three `init()` + `Register(...)` blocks and three `check…` functions in `rules_disk.go`; create `rules_disk_test.go`.

Note: `rules_disk.go` also needs `import "fmt"` at the top (for the `fmt.Sprintf` calls in Message construction). Add the import block:

```go
package samfile

import "fmt"

// §1 Disk-level rules (catalog docs/disk-validity-rules.md §1).
// ...
```

- [ ] **Step 2: Build and test**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && go build ./... && go test ./...`
Expected: all tests pass EXCEPT `TestPhase3RegistryGrowth`, which now reports 4 rules instead of 20 (still failing — will turn green at Task 6).

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-3 && \
g add rules_disk.go rules_disk_test.go && \
g commit -m "verify: §1 disk-level rules (track + sector range checks)

Adds three rules covering every track / sector reference on disk
(directory entry first sectors plus chain-link bytes 510-511):

  DISK-DIRECTORY-TRACKS    structural — no file references tracks 0-3
  DISK-TRACK-SIDE-ENCODING fatal      — track byte in {0x00-0x4F, 0x80-0xCF}
  DISK-SECTOR-RANGE        fatal      — sector byte in {1..10} (0 = terminator)

All three share a private trackSectorRefs helper that enumerates the
union of (a) first-sector references from used dir entries and (b)
chain-link bytes from each used file's sector walk. The walker is
bounded by 1560 steps (disk capacity) so a malformed chain cannot
hang the iteration; CHAIN-NO-CYCLE will catch the actual cycle in
a later task."
```

---

## Task 3: §2 Directory-entry rules (9 rules)

**Why this task exists:** the directory-entry rules form the largest single section. They're internally similar — each iterates `ctx.Journal.UsedFileEntries()` and tests one invariant per slot — so the bulk is repetition over a shared loop pattern. Doing them in one commit keeps the diff cohesive.

**Files:**
- Modify: `rules_directory.go` — register and implement 9 rules.
- Modify: `rules_directory_test.go` — create, with positive + negative tests per rule.

### Helper

Add at the top of `rules_directory.go`:

```go
package samfile

import (
    "fmt"
    "math/bits"
    "strings"
)

// usedSlot loops over every used directory slot in registration order
// and invokes fn for each. A small helper that keeps the per-rule
// Check function's loop body focused on the actual invariant.
func forEachUsedSlot(ctx *CheckContext, fn func(slot int, fe *FileEntry)) {
    for _, slot := range ctx.Journal.UsedFileEntries() {
        fn(slot, ctx.Journal[slot])
    }
}
```

### Rules (in catalog order)

For brevity, each rule below is given as a {Register block, Check function} pair. The full code is what to paste into `rules_directory.go`. The tests follow the same pattern as §1's: each rule gets a positive + negative test.

```go
// ----- DIR-TYPE-BYTE-IS-KNOWN -----
func init() {
    Register(Rule{
        ID:          "DIR-TYPE-BYTE-IS-KNOWN",
        Severity:    SeverityInconsistency,
        Description: "directory type byte (low 5 bits, attribute bits masked) is one of the documented file types",
        Citation:    "samdos/src/e.s:322-355",
        Check:       checkDirTypeByteIsKnown,
    })
}

// dirKnownTypes is the SAM-public set after masking off HIDDEN + PROTECTED.
// 0 is omitted: erased slots are caught by Used(), not here.
var dirKnownTypes = map[uint8]bool{
    5: true, 16: true, 17: true, 18: true, 19: true, 20: true,
}

func checkDirTypeByteIsKnown(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        t := uint8(fe.Type) & 0x1F
        if !dirKnownTypes[t] {
            findings = append(findings, Finding{
                RuleID:   "DIR-TYPE-BYTE-IS-KNOWN",
                Severity: SeverityInconsistency,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  fmt.Sprintf("masked type byte 0x%02x is not a documented file type (expected one of 5, 16-20)", t),
                Citation: "samdos/src/e.s:322-355",
            })
        }
    })
    return findings
}

// ----- DIR-ERASED-IS-ZERO -----
// Used() already encodes the rule but the catalog asks us to check the
// inverse statement: any slot whose raw Type byte is exactly 0x00 but
// whose other fields look populated (FirstSector non-zero) is suspicious.
// Phase 3 implements only the forward check: a used slot must NOT have
// Type == 0. (Empty Type 0 + Track 0 = legitimately free, which is the
// common case.)
func init() {
    Register(Rule{
        ID:          "DIR-ERASED-IS-ZERO",
        Severity:    SeverityStructural,
        Description: "a used directory slot has a non-zero type byte",
        Citation:    "samdos/src/c.s:1133-1143",
        Check:       checkDirErasedIsZero,
    })
}

func checkDirErasedIsZero(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        if uint8(fe.Type) == 0 {
            findings = append(findings, Finding{
                RuleID:   "DIR-ERASED-IS-ZERO",
                Severity: SeverityStructural,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  "used slot has type byte 0x00, which is the erased-slot sentinel",
                Citation: "samdos/src/c.s:1133-1143",
            })
        }
    })
    return findings
}

// ----- DIR-NAME-PADDING -----
func init() {
    Register(Rule{
        ID:          "DIR-NAME-PADDING",
        Severity:    SeverityCosmetic,
        Description: "filename bytes are printable ASCII or space-padded",
        Citation:    "sam-coupe_tech-man_v3-0.txt:4358-4359",
        Check:       checkDirNamePadding,
    })
}

func checkDirNamePadding(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        for i, b := range fe.Name {
            if b == 0x20 || (b >= 0x21 && b < 0x7F) {
                continue
            }
            findings = append(findings, Finding{
                RuleID:   "DIR-NAME-PADDING",
                Severity: SeverityCosmetic,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  fmt.Sprintf("filename byte %d is 0x%02x (expected printable ASCII or 0x20 space)", i, b),
                Citation: "sam-coupe_tech-man_v3-0.txt:4358-4359",
            })
            return // one finding per slot; further byte-by-byte detail belongs in a diagnostic
        }
    })
    return findings
}

// ----- DIR-NAME-NOT-EMPTY -----
func init() {
    Register(Rule{
        ID:          "DIR-NAME-NOT-EMPTY",
        Severity:    SeverityInconsistency,
        Description: "a used slot has at least one non-space, non-FF character in its 10-byte name",
        Citation:    "rom-disasm:22093-22105",
        Check:       checkDirNameNotEmpty,
    })
}

func checkDirNameNotEmpty(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        empty := true
        for _, b := range fe.Name {
            if b != 0x20 && b != 0xFF && b != 0 {
                empty = false
                break
            }
        }
        if empty {
            findings = append(findings, Finding{
                RuleID:   "DIR-NAME-NOT-EMPTY",
                Severity: SeverityInconsistency,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  "filename is all spaces / 0xFF / 0x00 (no visible characters)",
                Citation: "rom-disasm:22093-22105",
            })
        }
    })
    return findings
}

// ----- DIR-FIRST-SECTOR-VALID -----
func init() {
    Register(Rule{
        ID:          "DIR-FIRST-SECTOR-VALID",
        Severity:    SeverityFatal,
        Description: "directory entry's FirstSector points at a valid data sector",
        Citation:    "samfile.go:611-616",
        Check:       checkDirFirstSectorValid,
    })
}

func checkDirFirstSectorValid(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        fs := fe.FirstSector
        t := fs.Track
        s := fs.Sector
        // Side bit (0x80) is informational; mask it off for the cylinder check.
        cyl := t & 0x7F
        validTrack := (t < 80 || (t >= 128 && t < 208)) && cyl >= 4
        validSector := s >= 1 && s <= 10
        if !validTrack || !validSector {
            findings = append(findings, Finding{
                RuleID:   "DIR-FIRST-SECTOR-VALID",
                Severity: SeverityFatal,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  fmt.Sprintf("FirstSector (track=0x%02x, sector=%d) is not a valid data sector", t, s),
                Citation: "samfile.go:611-616",
            })
        }
    })
    return findings
}

// ----- DIR-SECTORS-MATCHES-CHAIN -----
// This rule walks each used slot's chain and compares the visited
// count to fe.Sectors. The walk is bounded by 1560 steps and uses the
// same single-step iteration pattern as trackSectorRefs / walkChain.
//
// Because walkChain has not yet landed (Task 4), this rule uses an
// inline walk to stay self-contained. Once Task 4 introduces walkChain,
// this rule's Check function can switch to it; that's a Task 4 follow-up.
func init() {
    Register(Rule{
        ID:          "DIR-SECTORS-MATCHES-CHAIN",
        Severity:    SeverityStructural,
        Description: "dir-entry Sectors count equals the number of sectors visited walking the chain to the (0,0) terminator",
        Citation:    "samfile.go:743-754",
        Check:       checkDirSectorsMatchesChain,
    })
}

func checkDirSectorsMatchesChain(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        // Inline minimal walk: count sectors until (0,0) or 1560 cap.
        var count uint16
        cur := fe.FirstSector
        for steps := 0; steps < 1560 && cur != nil; steps++ {
            count++
            sd, err := ctx.Disk.SectorData(cur)
            if err != nil {
                break
            }
            fp := sd.FilePart()
            if fp.NextSector.Track == 0 && fp.NextSector.Sector == 0 {
                break
            }
            cur = fp.NextSector
        }
        if count != fe.Sectors {
            findings = append(findings, Finding{
                RuleID:   "DIR-SECTORS-MATCHES-CHAIN",
                Severity: SeverityStructural,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  fmt.Sprintf("dir Sectors=%d, but chain walk visited %d sectors", fe.Sectors, count),
                Citation: "samfile.go:743-754",
            })
        }
    })
    return findings
}

// ----- DIR-SECTORS-MATCHES-MAP -----
func init() {
    Register(Rule{
        ID:          "DIR-SECTORS-MATCHES-MAP",
        Severity:    SeverityStructural,
        Description: "dir-entry Sectors count equals the popcount of the per-slot SectorAddressMap",
        Citation:    "sam-coupe_tech-man_v3-0.txt:4405-4414",
        Check:       checkDirSectorsMatchesMap,
    })
}

func checkDirSectorsMatchesMap(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        pop := 0
        for _, b := range fe.SectorAddressMap {
            pop += bits.OnesCount8(b)
        }
        if uint16(pop) != fe.Sectors {
            findings = append(findings, Finding{
                RuleID:   "DIR-SECTORS-MATCHES-MAP",
                Severity: SeverityStructural,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  fmt.Sprintf("dir Sectors=%d, but SectorAddressMap has popcount=%d", fe.Sectors, pop),
                Citation: "sam-coupe_tech-man_v3-0.txt:4405-4414",
            })
        }
    })
    return findings
}

// ----- DIR-SECTORS-NONZERO -----
func init() {
    Register(Rule{
        ID:          "DIR-SECTORS-NONZERO",
        Severity:    SeverityStructural,
        Description: "a used dir entry's Sectors count is at least 1",
        Citation:    "samdos/src/c.s:919-951",
        Check:       checkDirSectorsNonzero,
    })
}

func checkDirSectorsNonzero(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        if fe.Sectors == 0 {
            findings = append(findings, Finding{
                RuleID:   "DIR-SECTORS-NONZERO",
                Severity: SeverityStructural,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  "used slot has Sectors=0 (must be at least 1 for the body header)",
                Citation: "samdos/src/c.s:919-951",
            })
        }
    })
    return findings
}

// ----- DIR-SAM-WITHIN-CAPACITY -----
// SectorAddressMap is 195 bytes = 1560 bits. Disk capacity is 1560
// data sectors. So bits beyond bit-1559 must be zero. Bit-1559 is
// the high bit of byte 194 (bit 7 of byte 194). The rule is: bits
// 1560..1567 (bits 0..7 of a notional byte 195) cannot exist in the
// 195-byte array — already enforced by length. So the only check is
// inside byte 194: bits beyond bit 7 of byte 194... wait, all 8 bits
// of byte 194 ARE in range (bits 1552-1559). So the catalog's "top 3
// bits beyond bit 1559 are clear" wording is about WHICH disks have
// the 1560 bits; if there were a byte 195 it'd be the overflow.
//
// Re-reading the catalog: "byte 194 & 0xE0 == 0 (top 3 bits beyond
// bit 1559 are clear)". So the rule treats the top 3 bits of byte
// 194 as the overflow zone. Implement literally per the catalog.
func init() {
    Register(Rule{
        ID:          "DIR-SAM-WITHIN-CAPACITY",
        Severity:    SeverityInconsistency,
        Description: "SectorAddressMap byte 194's top 3 bits (1557-1559) are clear (no sector beyond disk capacity)",
        Citation:    "sam-coupe_tech-man_v3-0.txt:4405-4406",
        Check:       checkDirSAMWithinCapacity,
    })
}

func checkDirSAMWithinCapacity(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        if fe.SectorAddressMap[194]&0xE0 != 0 {
            findings = append(findings, Finding{
                RuleID:   "DIR-SAM-WITHIN-CAPACITY",
                Severity: SeverityInconsistency,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  fmt.Sprintf("SectorAddressMap[194]=0x%02x has bits beyond bit 1559 set", fe.SectorAddressMap[194]),
                Citation: "sam-coupe_tech-man_v3-0.txt:4405-4406",
            })
        }
    })
    return findings
}

// strings is imported above; if the linter complains it's unused after
// you remove a reference, drop the import.
var _ = strings.TrimSpace
```

(The `var _ = strings.TrimSpace` line keeps the `strings` import live in case a Check function uses it. If no rule in this file uses `strings`, remove the import; if any does, remove the dead `var _`. Verify before committing.)

### Tests

Create `rules_directory_test.go` with **one positive + one negative test per rule**. The positive test uses `cleanSingleFileDisk` (defined in `rules_disk_test.go` — accessible because both files are in the same package). The negative test patches one field via `dj[0].Field = …` + `WriteFileEntry`.

A template for each negative test (apply per-rule):

```go
func TestDirTypeByteIsKnownNegative(t *testing.T) {
    di, dj := cleanSingleFileDisk(t, "TEST", 100)
    dj[0].Type = FileType(7) // not in {5, 16-20}
    di.WriteFileEntry(dj, 0)
    findings := checkDirTypeByteIsKnown(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 1 || findings[0].RuleID != "DIR-TYPE-BYTE-IS-KNOWN" {
        t.Fatalf("got %d findings, first=%+v; want 1 DIR-TYPE-BYTE-IS-KNOWN", len(findings), findings)
    }
}
```

Concrete negative-test field patches for the remaining 8 rules (positive tests for all 9 use `cleanSingleFileDisk` then assert `len(findings) == 0` and need no further notes):

| Rule | Negative-test mutation |
|---|---|
| DIR-ERASED-IS-ZERO | `dj[0].Type = FileType(0)`. (The slot is still considered "used" by samfile because FirstSector.Track != 0.) Wait — `Used()` returns false for `FileType(0).String()` which starts with "UNKNOWN"... actually FileType(0) is FT_ERASED whose String returns "ERASED" — not prefixed with "UNKNOWN". So `Used()` returns true because FirstSector.Track != 0. Good — this slot is "used" with Type=0, exactly the violation. |
| DIR-NAME-PADDING | `dj[0].Name = Filename{'A', 0x01, 'B', ' ', ' ', ' ', ' ', ' ', ' ', ' '}` (0x01 control char). |
| DIR-NAME-NOT-EMPTY | `dj[0].Name = Filename{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}` (all spaces). |
| DIR-FIRST-SECTOR-VALID | `dj[0].FirstSector.Sector = 99`. |
| DIR-SECTORS-MATCHES-CHAIN | `dj[0].Sectors = 99` (real chain is shorter). |
| DIR-SECTORS-MATCHES-MAP | `dj[0].Sectors = 99` (map popcount is real allocation). |
| DIR-SECTORS-NONZERO | `dj[0].Sectors = 0`. |
| DIR-SAM-WITHIN-CAPACITY | `dj[0].SectorAddressMap[194] = 0xE0` (set top 3 bits). |

For each negative test, after the mutation: call `di.WriteFileEntry(dj, 0)`, build a fresh CheckContext from `di.DiskJournal()` (to re-read the patched state), and assert `len(findings) == 1` with the right `RuleID`.

- [ ] **Step 1: Implement all 9 rules and tests**

Add the 9 Register/check pairs to `rules_directory.go` (with the import block above) and create `rules_directory_test.go` with 18 tests (9 positive + 9 negative).

- [ ] **Step 2: Build + test**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && go build ./... && go test ./...`
Expected: all tests pass EXCEPT `TestPhase3RegistryGrowth`, which now reports 13 rules instead of 20.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-3 && \
g add rules_directory.go rules_directory_test.go && \
g commit -m "verify: §2 directory-entry rules (9 rules)

Adds nine rules covering internal consistency of each directory
entry. Each rule iterates ctx.Journal.UsedFileEntries() via the
forEachUsedSlot helper:

  DIR-TYPE-BYTE-IS-KNOWN     inconsistency
  DIR-ERASED-IS-ZERO         structural
  DIR-NAME-PADDING           cosmetic
  DIR-NAME-NOT-EMPTY         inconsistency
  DIR-FIRST-SECTOR-VALID     fatal
  DIR-SECTORS-MATCHES-CHAIN  structural
  DIR-SECTORS-MATCHES-MAP    structural
  DIR-SECTORS-NONZERO        structural
  DIR-SAM-WITHIN-CAPACITY    inconsistency

DIR-SECTORS-MATCHES-CHAIN uses an inline minimal chain walk; this
will be folded into the shared walkChain helper landing in the
next commit."
```

---

## Task 4: walkChain helper + §3 chain rules (3 rules)

**Why this task exists:** §3 needs a chain walker that records every visited sector (for cycle detection and SectorAddressMap comparison) — strictly more than `samfile.File`'s by-count walk. Defining the walker once means CHAIN-NO-CYCLE, CHAIN-MATCHES-SAM, and the §4 cross-entry rules can all read from the same data.

**Files:**
- Modify: `rules_chain.go` — add `walkChain` helper and three rules.
- Modify: `rules_chain_test.go` — create, with `walkChain` unit tests + positive/negative tests per rule.

### The walkChain helper

```go
package samfile

import "fmt"

// chainStep is one entry in a sector chain walk.
type chainStep struct {
    Sector Sector // the sector that was read at this step (copy, not pointer)
    Next   Sector // the (track, sector) link at bytes 510-511 of Sector
}

// chainWalkResult is the outcome of a walkChain call.
type chainWalkResult struct {
    Steps       []chainStep // in walk order
    Terminated  bool        // true iff a (0, 0) link was encountered
    Cycle       *Sector     // first sector revisited, if any (nil = no cycle)
    Bailed      bool        // true iff the walk hit the 1560-step cap without terminating or cycling
}

// walkChain follows the link chain starting at first for at most 1560
// steps (the disk's data-sector capacity, an absolute upper bound on a
// terminated chain). It records each sector visited and the (track,
// sector) link found at its bytes 510-511. The walk halts on:
//
//   - a (0, 0) terminator (Terminated = true);
//   - a revisited sector (Cycle = &<first repeat>);
//   - hitting the 1560-step cap (Bailed = true).
//
// On read error from SectorData, the walk halts at the current Steps
// length without setting any flag. Callers can still use Steps to see
// what was reachable. (No findings are surfaced for read errors here;
// DISK-TRACK-SIDE-ENCODING / DIR-FIRST-SECTOR-VALID catch the underlying
// out-of-range track byte that triggers the SectorData error.)
func walkChain(di *DiskImage, first *Sector) chainWalkResult {
    var result chainWalkResult
    visited := make(map[Sector]bool)
    cur := *first
    for steps := 0; steps < 1560; steps++ {
        if visited[cur] {
            c := cur
            result.Cycle = &c
            return result
        }
        visited[cur] = true
        sd, err := di.SectorData(&cur)
        if err != nil {
            return result
        }
        fp := sd.FilePart()
        next := *fp.NextSector
        result.Steps = append(result.Steps, chainStep{Sector: cur, Next: next})
        if next.Track == 0 && next.Sector == 0 {
            result.Terminated = true
            return result
        }
        cur = next
    }
    result.Bailed = true
    return result
}
```

### Rules

```go
// ----- CHAIN-TERMINATOR-ZERO-ZERO -----
func init() {
    Register(Rule{
        ID:          "CHAIN-TERMINATOR-ZERO-ZERO",
        Severity:    SeverityStructural,
        Description: "each used file's sector chain ends with a (0, 0) link",
        Citation:    "samdos/src/b.s:104-110",
        Check:       checkChainTerminatorZeroZero,
    })
}

func checkChainTerminatorZeroZero(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        result := walkChain(ctx.Disk, fe.FirstSector)
        if !result.Terminated {
            lastSec := fe.FirstSector
            if n := len(result.Steps); n > 0 {
                s := result.Steps[n-1].Sector
                lastSec = &s
            }
            msg := "chain does not terminate"
            if result.Cycle != nil {
                msg = fmt.Sprintf("chain has a cycle (revisited %v)", result.Cycle)
            } else if result.Bailed {
                msg = "chain exceeds 1560 steps without (0, 0) link"
            }
            findings = append(findings, Finding{
                RuleID:   "CHAIN-TERMINATOR-ZERO-ZERO",
                Severity: SeverityStructural,
                Location: SectorLocation(slot, fe.Name.String(), lastSec, 510),
                Message:  msg,
                Citation: "samdos/src/b.s:104-110",
            })
        }
    })
    return findings
}

// ----- CHAIN-NO-CYCLE -----
func init() {
    Register(Rule{
        ID:          "CHAIN-NO-CYCLE",
        Severity:    SeverityStructural,
        Description: "each used file's sector chain has no revisited sectors",
        Citation:    "samfile.go:743-754",
        Check:       checkChainNoCycle,
    })
}

func checkChainNoCycle(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        result := walkChain(ctx.Disk, fe.FirstSector)
        if result.Cycle != nil {
            findings = append(findings, Finding{
                RuleID:   "CHAIN-NO-CYCLE",
                Severity: SeverityStructural,
                Location: SectorLocation(slot, fe.Name.String(), result.Cycle, 510),
                Message:  fmt.Sprintf("chain cycles: sector %v is revisited", result.Cycle),
                Citation: "samfile.go:743-754",
            })
        }
    })
    return findings
}

// ----- CHAIN-MATCHES-SAM -----
func init() {
    Register(Rule{
        ID:          "CHAIN-MATCHES-SAM",
        Severity:    SeverityStructural,
        Description: "the set of sectors walked by the chain equals the bits set in the SectorAddressMap",
        Citation:    "samdos/src/c.s:1306-1343",
        Check:       checkChainMatchesSAM,
    })
}

func checkChainMatchesSAM(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        result := walkChain(ctx.Disk, fe.FirstSector)
        walked := make(map[Sector]bool, len(result.Steps))
        for _, st := range result.Steps {
            walked[st.Sector] = true
        }
        mapSet := make(map[Sector]bool)
        for _, sec := range fe.SectorAddressMap.UsedSectors() {
            mapSet[*sec] = true
        }
        // Symmetric difference: any sector in one set but not the other.
        for s := range walked {
            if !mapSet[s] {
                findings = append(findings, Finding{
                    RuleID:   "CHAIN-MATCHES-SAM",
                    Severity: SeverityStructural,
                    Location: SlotLocation(slot, fe.Name.String()),
                    Message:  fmt.Sprintf("sector %v is visited by the chain but not set in the SectorAddressMap", s),
                    Citation: "samdos/src/c.s:1306-1343",
                })
                return // one finding per slot is enough; the disagreement is the signal
            }
        }
        for s := range mapSet {
            if !walked[s] {
                findings = append(findings, Finding{
                    RuleID:   "CHAIN-MATCHES-SAM",
                    Severity: SeverityStructural,
                    Location: SlotLocation(slot, fe.Name.String()),
                    Message:  fmt.Sprintf("sector %v is set in the SectorAddressMap but not visited by the chain", s),
                    Citation: "samdos/src/c.s:1306-1343",
                })
                return
            }
        }
    })
    return findings
}
```

### Tests

Create `rules_chain_test.go`:

```go
package samfile

import "testing"

// walkChain unit tests come first — exercise the helper directly with
// fabricated chains so the rule tests can trust the walker behaves.

func TestWalkChainClean(t *testing.T) {
    di, _ := cleanSingleFileDisk(t, "TEST", 600) // 2 sectors worth of payload
    result := walkChain(di, di.DiskJournal()[0].FirstSector)
    if !result.Terminated {
        t.Errorf("clean chain: Terminated = false; want true")
    }
    if result.Cycle != nil {
        t.Errorf("clean chain: Cycle = %v; want nil", result.Cycle)
    }
    if result.Bailed {
        t.Errorf("clean chain: Bailed = true; want false")
    }
    if len(result.Steps) < 2 {
        t.Errorf("clean chain: %d steps; want >= 2", len(result.Steps))
    }
}

func TestWalkChainCycleDetection(t *testing.T) {
    di, dj := cleanSingleFileDisk(t, "TEST", 600)
    fe := dj[0]
    first := fe.FirstSector
    // Force sector 1's NextSector to point back at itself.
    sd, err := di.SectorData(first)
    if err != nil {
        t.Fatalf("SectorData: %v", err)
    }
    raw := sd[:]
    raw[510] = first.Track
    raw[511] = first.Sector
    di.WriteSector(first, sd)

    result := walkChain(di, first)
    if result.Cycle == nil {
        t.Errorf("Cycle = nil; want non-nil (chain points at itself)")
    }
    if result.Terminated {
        t.Errorf("Terminated = true; want false")
    }
}

// Now the three rule tests.

func TestChainTerminatorZeroZeroPositive(t *testing.T) {
    di, _ := cleanSingleFileDisk(t, "TEST", 100)
    findings := checkChainTerminatorZeroZero(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 0 {
        t.Errorf("clean chain: %d findings; want 0", len(findings))
    }
}

func TestChainTerminatorZeroZeroNegative(t *testing.T) {
    di, _ := cleanSingleFileDisk(t, "TEST", 100)
    first := di.DiskJournal()[0].FirstSector
    // Overwrite the terminator with a fake link (it would loop forever
    // if walkChain weren't bounded).
    sd, _ := di.SectorData(first)
    raw := sd[:]
    raw[510] = first.Track
    raw[511] = first.Sector
    di.WriteSector(first, sd)

    findings := checkChainTerminatorZeroZero(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 1 || findings[0].RuleID != "CHAIN-TERMINATOR-ZERO-ZERO" {
        t.Fatalf("got %d findings, first=%+v; want 1 CHAIN-TERMINATOR-ZERO-ZERO", len(findings), findings)
    }
}

func TestChainNoCyclePositive(t *testing.T) {
    di, _ := cleanSingleFileDisk(t, "TEST", 100)
    findings := checkChainNoCycle(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 0 {
        t.Errorf("clean chain: %d findings; want 0", len(findings))
    }
}

func TestChainNoCycleNegative(t *testing.T) {
    di, _ := cleanSingleFileDisk(t, "TEST", 100)
    first := di.DiskJournal()[0].FirstSector
    sd, _ := di.SectorData(first)
    raw := sd[:]
    raw[510] = first.Track
    raw[511] = first.Sector
    di.WriteSector(first, sd)

    findings := checkChainNoCycle(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 1 || findings[0].RuleID != "CHAIN-NO-CYCLE" {
        t.Fatalf("got %d findings, first=%+v; want 1 CHAIN-NO-CYCLE", len(findings), findings)
    }
}

func TestChainMatchesSAMPositive(t *testing.T) {
    di, _ := cleanSingleFileDisk(t, "TEST", 1500) // ~3 sectors worth
    findings := checkChainMatchesSAM(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 0 {
        t.Errorf("clean disk: %d findings; want 0", len(findings))
    }
}

func TestChainMatchesSAMNegative(t *testing.T) {
    di, dj := cleanSingleFileDisk(t, "TEST", 1500)
    // Clear a bit that IS set in the map so walked > mapSet.
    for i, b := range dj[0].SectorAddressMap {
        if b != 0 {
            dj[0].SectorAddressMap[i] &^= (b & -b) // clear the lowest set bit
            break
        }
    }
    di.WriteFileEntry(dj, 0)
    findings := checkChainMatchesSAM(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 1 || findings[0].RuleID != "CHAIN-MATCHES-SAM" {
        t.Fatalf("got %d findings, first=%+v; want 1 CHAIN-MATCHES-SAM", len(findings), findings)
    }
}
```

A couple of notes about the test code:

- `SectorData` returns `*SectorData` whose underlying type is `[512]byte`; `sd[:]` gives a slice view. After mutating, `di.WriteSector(sec, sd)` writes the 512 bytes back.
- `dj[0].SectorAddressMap[i] &^= (b & -b)` clears the lowest set bit of `b`. Go's `&^` is bit-clear (`a &^ b == a & ~b`); `b & -b` isolates the lowest set bit. This produces a one-bit disagreement between the map and the actual chain — perfect for the rule.

- [ ] **Step 1: Implement walkChain + three rules + their tests**

- [ ] **Step 2: Build + test**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && go build ./... && go test ./...`
Expected: all tests pass EXCEPT `TestPhase3RegistryGrowth`, which now reports 16 rules instead of 20.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-3 && \
g add rules_chain.go rules_chain_test.go && \
g commit -m "verify: §3 sector-chain rules + walkChain helper (3 rules)

Adds the canonical walkChain helper and three rules built on it:

  CHAIN-TERMINATOR-ZERO-ZERO structural — chain ends with (0, 0)
  CHAIN-NO-CYCLE             structural — no sector revisited
  CHAIN-MATCHES-SAM          structural — walked set == map bits

walkChain is bounded at 1560 steps (disk capacity) so cycles and
chains without terminators cannot hang the iteration; both
conditions are reported via the chainWalkResult flags. The helper
is also used by §4 cross-entry rules in the next commit."
```

---

## Task 5: §4 Cross-entry rules (3 rules)

**Why this task exists:** §4 rules read across multiple slots — they need data that's expensive to compute per-slot. Doing them in their own commit (after `walkChain` lands) keeps the dependency graph linear.

**Files:**
- Modify: `rules_cross.go` — register and implement 3 rules.
- Modify: `rules_cross_test.go` — create.

### Rules

```go
package samfile

import "fmt"

// ----- CROSS-NO-SECTOR-OVERLAP -----
func init() {
    Register(Rule{
        ID:          "CROSS-NO-SECTOR-OVERLAP",
        Severity:    SeverityFatal,
        Description: "no two used files claim the same data sector",
        Citation:    "samdos/src/c.s:895-951",
        Check:       checkCrossNoSectorOverlap,
    })
}

func checkCrossNoSectorOverlap(ctx *CheckContext) []Finding {
    var findings []Finding
    // owner[sector] = list of (slot, filename) entries claiming this sector.
    type claim struct{ Slot int; Name string }
    owner := make(map[Sector][]claim)
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        name := fe.Name.String()
        for _, sec := range fe.SectorAddressMap.UsedSectors() {
            owner[*sec] = append(owner[*sec], claim{slot, name})
        }
    })
    for sec, claims := range owner {
        if len(claims) < 2 {
            continue
        }
        s := sec
        findings = append(findings, Finding{
            RuleID:   "CROSS-NO-SECTOR-OVERLAP",
            Severity: SeverityFatal,
            Location: SectorLocation(claims[0].Slot, claims[0].Name, &s, -1),
            Message:  fmt.Sprintf("sector %v is claimed by %d slots (first: %d %q, second: %d %q)",
                s, len(claims), claims[0].Slot, claims[0].Name, claims[1].Slot, claims[1].Name),
            Citation: "samdos/src/c.s:895-951",
        })
    }
    return findings
}

// ----- CROSS-NO-DUPLICATE-NAMES -----
func init() {
    Register(Rule{
        ID:          "CROSS-NO-DUPLICATE-NAMES",
        Severity:    SeverityInconsistency,
        Description: "no two used directory entries share the same filename (case-insensitive)",
        Citation:    "samdos/src/c.s:1196-1219",
        Check:       checkCrossNoDuplicateNames,
    })
}

func checkCrossNoDuplicateNames(ctx *CheckContext) []Finding {
    var findings []Finding
    seen := make(map[string]int) // lowercased trimmed name -> first slot to use it
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        key := strings.ToLower(strings.TrimSpace(fe.Name.String()))
        if key == "" {
            return // empty names handled by DIR-NAME-NOT-EMPTY
        }
        if prev, ok := seen[key]; ok {
            findings = append(findings, Finding{
                RuleID:   "CROSS-NO-DUPLICATE-NAMES",
                Severity: SeverityInconsistency,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  fmt.Sprintf("filename %q duplicates slot %d", key, prev),
                Citation: "samdos/src/c.s:1196-1219",
            })
            return
        }
        seen[key] = slot
    })
    return findings
}

// ----- CROSS-DIRECTORY-AREA-UNUSED -----
func init() {
    Register(Rule{
        ID:          "CROSS-DIRECTORY-AREA-UNUSED",
        Severity:    SeverityStructural,
        Description: "no chain link in any used file references a directory-area sector (tracks 0-3 of side 0)",
        Citation:    "samfile.go:984-987",
        Check:       checkCrossDirectoryAreaUnused,
    })
}

func checkCrossDirectoryAreaUnused(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        result := walkChain(ctx.Disk, fe.FirstSector)
        for _, st := range result.Steps {
            if (st.Sector.Track & 0x7F) < 4 {
                s := st.Sector
                findings = append(findings, Finding{
                    RuleID:   "CROSS-DIRECTORY-AREA-UNUSED",
                    Severity: SeverityStructural,
                    Location: SectorLocation(slot, fe.Name.String(), &s, -1),
                    Message:  fmt.Sprintf("chain visits %v which is in the directory area", s),
                    Citation: "samfile.go:984-987",
                })
                return // one finding per slot
            }
        }
    })
    return findings
}
```

Add `import "strings"` to `rules_cross.go` (used by CROSS-NO-DUPLICATE-NAMES).

### Tests

Create `rules_cross_test.go`:

```go
package samfile

import "testing"

func TestCrossNoSectorOverlapPositive(t *testing.T) {
    di := NewDiskImage()
    if err := di.AddCodeFile("A", make([]byte, 100), 0x8000, 0); err != nil {
        t.Fatalf("AddCodeFile A: %v", err)
    }
    if err := di.AddCodeFile("B", make([]byte, 100), 0x8000, 0); err != nil {
        t.Fatalf("AddCodeFile B: %v", err)
    }
    findings := checkCrossNoSectorOverlap(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 0 {
        t.Errorf("two distinct files: %d findings; want 0", len(findings))
    }
}

func TestCrossNoSectorOverlapNegative(t *testing.T) {
    di := NewDiskImage()
    if err := di.AddCodeFile("A", make([]byte, 100), 0x8000, 0); err != nil {
        t.Fatalf("AddCodeFile A: %v", err)
    }
    if err := di.AddCodeFile("B", make([]byte, 100), 0x8000, 0); err != nil {
        t.Fatalf("AddCodeFile B: %v", err)
    }
    dj := di.DiskJournal()
    // Copy slot 0's map into slot 1 so they claim overlapping sectors.
    dj[1].SectorAddressMap = dj[0].SectorAddressMap
    di.WriteFileEntry(dj, 1)
    findings := checkCrossNoSectorOverlap(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) < 1 || findings[0].RuleID != "CROSS-NO-SECTOR-OVERLAP" {
        t.Fatalf("got %d findings, first=%+v; want at least one CROSS-NO-SECTOR-OVERLAP",
            len(findings), findings)
    }
}

func TestCrossNoDuplicateNamesPositive(t *testing.T) {
    di := NewDiskImage()
    if err := di.AddCodeFile("A", make([]byte, 100), 0x8000, 0); err != nil {
        t.Fatalf(": %v", err)
    }
    if err := di.AddCodeFile("B", make([]byte, 100), 0x8000, 0); err != nil {
        t.Fatalf(": %v", err)
    }
    findings := checkCrossNoDuplicateNames(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 0 {
        t.Errorf("distinct names: %d findings; want 0", len(findings))
    }
}

func TestCrossNoDuplicateNamesNegative(t *testing.T) {
    di := NewDiskImage()
    if err := di.AddCodeFile("A", make([]byte, 100), 0x8000, 0); err != nil {
        t.Fatalf(": %v", err)
    }
    if err := di.AddCodeFile("B", make([]byte, 100), 0x8000, 0); err != nil {
        t.Fatalf(": %v", err)
    }
    dj := di.DiskJournal()
    // Rename slot 1 to "A" so it duplicates slot 0.
    copy(dj[1].Name[:], "A         ")
    di.WriteFileEntry(dj, 1)
    findings := checkCrossNoDuplicateNames(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 1 || findings[0].RuleID != "CROSS-NO-DUPLICATE-NAMES" {
        t.Fatalf("got %d findings, first=%+v; want 1 CROSS-NO-DUPLICATE-NAMES",
            len(findings), findings)
    }
}

func TestCrossDirectoryAreaUnusedPositive(t *testing.T) {
    di, _ := cleanSingleFileDisk(t, "TEST", 100)
    findings := checkCrossDirectoryAreaUnused(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 0 {
        t.Errorf("clean disk: %d findings; want 0", len(findings))
    }
}

func TestCrossDirectoryAreaUnusedNegative(t *testing.T) {
    di, _ := cleanSingleFileDisk(t, "TEST", 600) // 2 sectors so there's a chain link
    first := di.DiskJournal()[0].FirstSector
    // Point sector 0's next-link at a directory-area sector (T2 S5).
    sd, _ := di.SectorData(first)
    raw := sd[:]
    raw[510] = 2 // track 2, in directory area
    raw[511] = 5
    di.WriteSector(first, sd)
    findings := checkCrossDirectoryAreaUnused(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) < 1 || findings[0].RuleID != "CROSS-DIRECTORY-AREA-UNUSED" {
        t.Fatalf("got %d findings, first=%+v; want at least one CROSS-DIRECTORY-AREA-UNUSED",
            len(findings), findings)
    }
}
```

- [ ] **Step 1: Implement the three rules and tests**

- [ ] **Step 2: Build + test**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && go build ./... && go test ./...`
Expected: all tests pass EXCEPT `TestPhase3RegistryGrowth`, which now reports 19 rules instead of 20.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-3 && \
g add rules_cross.go rules_cross_test.go && \
g commit -m "verify: §4 cross-entry rules (3 rules)

Adds three rules that compare data across used directory slots:

  CROSS-NO-SECTOR-OVERLAP    fatal         — no two files share a sector
  CROSS-NO-DUPLICATE-NAMES   inconsistency — names unique (case-insensitive)
  CROSS-DIRECTORY-AREA-UNUSED structural   — no chain visits tracks 0-3

CROSS-NO-SECTOR-OVERLAP uses each file's SectorAddressMap to detect
double-claims without re-walking chains. CROSS-DIRECTORY-AREA-UNUSED
reuses the walkChain helper from Task 4 to inspect every visited
sector, not just the dir-entry first-sector (which DISK-DIRECTORY-TRACKS
already covers)."
```

---

## Task 6: §15 CHAIN-SECTOR-COUNT-MINIMAL (1 rule)

**Why this task exists:** the catalog carryover rule from samfile v2.1.0's wish list ("verifying that files do not contain empty sectors"). Phase 3's last rule; lives in `rules_chain.go` because it's a chain-shape check.

**Files:**
- Modify: `rules_chain.go` — append one rule.
- Modify: `rules_chain_test.go` — append two tests.

```go
// ----- CHAIN-SECTOR-COUNT-MINIMAL -----
func init() {
    Register(Rule{
        ID:          "CHAIN-SECTOR-COUNT-MINIMAL",
        Severity:    SeverityCosmetic,
        Description: "used file occupies exactly ceil((9 + body length) / 510) sectors (no padding sectors)",
        Citation:    "samfile.go:919",
        Check:       checkChainSectorCountMinimal,
    })
}

func checkChainSectorCountMinimal(ctx *CheckContext) []Finding {
    var findings []Finding
    forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
        bodyLen := int(fe.Length())
        required := uint16((bodyLen + 9 + 509) / 510)
        if fe.Sectors != required {
            findings = append(findings, Finding{
                RuleID:   "CHAIN-SECTOR-COUNT-MINIMAL",
                Severity: SeverityCosmetic,
                Location: SlotLocation(slot, fe.Name.String()),
                Message:  fmt.Sprintf("file uses %d sectors but %d would suffice (bodyLen=%d)",
                    fe.Sectors, required, bodyLen),
                Citation: "samfile.go:919",
            })
        }
    })
    return findings
}
```

Tests:

```go
func TestChainSectorCountMinimalPositive(t *testing.T) {
    di, _ := cleanSingleFileDisk(t, "TEST", 1500)
    findings := checkChainSectorCountMinimal(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 0 {
        t.Errorf("clean disk: %d findings; want 0", len(findings))
    }
}

func TestChainSectorCountMinimalNegative(t *testing.T) {
    di, dj := cleanSingleFileDisk(t, "TEST", 100)
    dj[0].Sectors += 1 // claim one more sector than the body needs
    di.WriteFileEntry(dj, 0)
    findings := checkChainSectorCountMinimal(&CheckContext{
        Disk: di, Journal: di.DiskJournal(),
    })
    if len(findings) != 1 || findings[0].RuleID != "CHAIN-SECTOR-COUNT-MINIMAL" {
        t.Fatalf("got %d findings, first=%+v; want 1 CHAIN-SECTOR-COUNT-MINIMAL",
            len(findings), findings)
    }
}
```

- [ ] **Step 1: Append the rule and tests**

- [ ] **Step 2: Build + test**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && go test ./...`
Expected: all 19 Phase-3 rules registered; `TestPhase3RegistryGrowth` now reports 20 and PASSES. Every rule's positive + negative test passes.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-3 && \
g add rules_chain.go rules_chain_test.go && \
g commit -m "verify: §15 CHAIN-SECTOR-COUNT-MINIMAL (1 rule)

Closes phase 3's rule count at 19. Carryover from samfile v2.1.0
release notes ('verifying that files do not contain empty sectors'):
warn when Sectors > ceil((9 + bodyLen) / 510), i.e. when the file
occupies more sectors than strictly necessary. Cosmetic severity:
extra trailing sectors waste space but don't break anything.

TestPhase3RegistryGrowth now passes (20 rules registered)."
```

---

## Task 7: Integration test — Verify on the committed corpus

**Why this task exists:** every test so far has been unit-level (fabricated disks). One smoke test that drives `(*DiskImage).Verify()` end-to-end on `testdata/ETrackerv1.2.mgt` catches any rule that panics, double-fires, or misses a real-world byte layout. The assertion shape is "no panic and the report is well-formed" rather than "specific findings present" — we don't have ground truth on which catalog rules this disk should trip.

**Files:**
- Modify: `verify_test.go` — append one test.

```go
func TestVerifyOnTestdataCorpus(t *testing.T) {
    const path = "testdata/ETrackerv1.2.mgt"
    if _, err := os.Stat(path); err != nil {
        t.Skipf("corpus image not present (%v); skipping", err)
    }
    di, err := Load(path)
    if err != nil {
        t.Fatalf("Load(%q): %v", path, err)
    }
    report := di.Verify()
    // Smoke shape: Dialect is set to one of the four documented
    // values, and every Finding's RuleID is a registered rule.
    switch report.Dialect {
    case DialectUnknown, DialectSAMDOS1, DialectSAMDOS2, DialectMasterDOS:
        // ok
    default:
        t.Errorf("Dialect = %v; not a documented value", report.Dialect)
    }
    knownIDs := make(map[string]bool)
    for _, r := range Rules() {
        knownIDs[r.ID] = true
    }
    for i, f := range report.Findings {
        if !knownIDs[f.RuleID] {
            t.Errorf("Findings[%d].RuleID = %q is not registered", i, f.RuleID)
        }
        if f.Citation == "" {
            t.Errorf("Findings[%d].Citation is empty (rule %s)", i, f.RuleID)
        }
    }
    t.Logf("verify(%s): dialect=%s, %d findings", path, report.Dialect, len(report.Findings))
}
```

Make sure `verify_test.go` imports `os` (it may not already). If you have to add it, place it alongside the existing imports.

- [ ] **Step 1: Append the test and any missing import**

- [ ] **Step 2: Run**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && go test -run TestVerifyOnTestdataCorpus -v ./...`
Expected: PASS with a log line like `verify(testdata/ETrackerv1.2.mgt): dialect=unknown, N findings` (N may be 0 or any positive integer; the test does not assert a specific N).

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-3 && \
g add verify_test.go && \
g commit -m "verify: integration smoke test against testdata corpus

Runs (*DiskImage).Verify() on testdata/ETrackerv1.2.mgt and asserts
that the report is structurally well-formed: dialect is one of the
four documented values, every Finding cites a registered RuleID,
and every Finding has a non-empty Citation. Does not assert a
specific finding count — that's the corpus-validation pass in
phase 7."
```

---

## Task 8: Final verification, CLI smoke, push, draft PR, monitor CI

**Why this task exists:** the gate before opening the PR. Pete's standing rule (`memory/feedback_correctness_over_workarounds.md`) is to verify the change actually works end-to-end — running CLI smoke against the M0 boot disk surfaced a real bug in Phase 2 and may do so again here.

**Files:** none modified.

- [ ] **Step 1: Run the full test suite + vet**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && go test ./... && go vet ./...`
Expected: all green, vet silent.

- [ ] **Step 2: Build the CLI**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && go build -o /tmp/samfile-phase3 ./cmd/samfile`
Expected: silent, exit 0.

- [ ] **Step 3: Run verify on the committed testdata corpus**

Run: `/tmp/samfile-phase3 verify -i /Users/pmoore/git/samfile-verify-phase-3/testdata/ETrackerv1.2.mgt`
Expected: a non-empty findings list (this disk likely trips at least DIR-* rules). Inspect the output by eye:

- Every line cites a real rule ID.
- No panics, no nonsense byte values in the messages.
- The `detected dialect:` line is one of the four documented values.

- [ ] **Step 4: Run verify on the M0 boot disk (if present)**

Run: `[ -f /Users/pmoore/git/sam-aarch64/build/test.mgt ] && /tmp/samfile-phase3 verify -i /Users/pmoore/git/sam-aarch64/build/test.mgt || echo 'no M0 disk; skipping'`
Expected: the dialect line reports `samdos2` (Phase 2 confirmed). The findings list may or may not be empty — that's data, not a pass/fail criterion. If any finding looks structurally suspicious (e.g. a CHAIN rule firing on a clean boot disk, an out-of-range value in a Message), stop and investigate before pushing.

- [ ] **Step 5: Push**

```bash
cd /Users/pmoore/git/samfile-verify-phase-3 && \
g push -u origin feat/verify-phase-3-disk-dir-chain-rules
```

- [ ] **Step 6: Open the draft PR**

```bash
cd /Users/pmoore/git/samfile-verify-phase-3 && gh pr create --draft \
  --base master \
  --title "verify: Phase 3 — disk, directory, chain & cross-entry rules (20 rules)" \
  --body "$(cat <<'EOF'
Phase 3 of `samfile verify` (spec: `docs/specs/2026-05-11-verify-feature-design.md`, plan: `docs/plans/2026-05-12-verify-phase-3-disk-dir-chain-rules.md`). Implements 19 of the catalog's structural rules — every rule that doesn't depend on file-type specifics. After this lands the registry holds 20 rules total (Phase-1 smoke + 19 from Phase 3); file-type rules (FT_CODE, FT_SAM_BASIC, ...) follow in Phases 4-6.

Rules added, grouped by catalog section:

- §1 disk-level (3): DISK-DIRECTORY-TRACKS, DISK-TRACK-SIDE-ENCODING, DISK-SECTOR-RANGE
- §2 directory-entry (9): DIR-TYPE-BYTE-IS-KNOWN, DIR-ERASED-IS-ZERO, DIR-NAME-PADDING, DIR-NAME-NOT-EMPTY, DIR-FIRST-SECTOR-VALID, DIR-SECTORS-MATCHES-CHAIN, DIR-SECTORS-MATCHES-MAP, DIR-SECTORS-NONZERO, DIR-SAM-WITHIN-CAPACITY
- §3 chain (3): CHAIN-TERMINATOR-ZERO-ZERO, CHAIN-NO-CYCLE, CHAIN-MATCHES-SAM
- §4 cross-entry (3): CROSS-NO-SECTOR-OVERLAP, CROSS-NO-DUPLICATE-NAMES, CROSS-DIRECTORY-AREA-UNUSED
- §15 carryover (1): CHAIN-SECTOR-COUNT-MINIMAL

Four catalog entries are deliberately deferred (DISK-IMAGE-SIZE, DISK-NOT-EDSK, DIR-SLOT-COUNT, DIR-SECTORS-BIG-ENDIAN, CHAIN-FIRST-MATCHES-DIR) — they are preconditions enforced by `Load` / `DiskJournal`'s parser or tautologies post-parse, with no path to fail at Verify time. See the plan for one-line rationales.

Architecture:

- One file per catalog section: `rules_disk.go`, `rules_directory.go`, `rules_chain.go`, `rules_cross.go`. Each rule registers in an `init()` block alongside its Check function.
- A new private `walkChain` helper in `rules_chain.go` is the single canonical sector-chain walker. Bounded at 1560 steps (disk capacity); records each visited sector and reports termination / cycle / bail status in a `chainWalkResult` struct. Used by §3 and §4 rules.
- A small `forEachUsedSlot` helper in `rules_directory.go` keeps per-rule loop bodies focused on the actual invariant.
- A new regression gate `TestPhase3RegistryGrowth` pins the registry count at 20 so a future rule that's accidentally never registered fails immediately.

Smoke-tested in this environment:

- `samfile verify -i testdata/ETrackerv1.2.mgt` — runs to completion, dialect `unknown`, N findings (depending on the disk's actual state — the test asserts shape, not content).
- `samfile verify -i ../sam-aarch64/build/test.mgt` (the M0 boot disk) — dialect `samdos2`, findings list inspected.

## Test plan

- [x] `go test ./...` — all green (one positive + one negative test per rule = 38 new unit tests, plus walkChain unit tests, plus the integration smoke test on the corpus)
- [x] `go vet ./...` — clean
- [x] CLI smoke against testdata/ETrackerv1.2.mgt produces a well-formed report
- [x] CLI smoke against the M0 boot disk (sam-aarch64 build/test.mgt) is well-formed
- [ ] GitHub Actions CI green

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 7: Monitor CI**

Run: `cd /Users/pmoore/git/samfile-verify-phase-3 && gh pr checks --watch`
Expected: every check passes. If any fails, diagnose locally with the same `go test` / `go vet` commands and iterate per Pete's standing rule (fix small issues into the relevant commit; escalate design questions). Do NOT mark the PR ready-for-review until Pete approves.

- [ ] **Step 8: Hand off**

Reply with the PR URL, CI status, and any noteworthy findings from the corpus / M0 smoke runs.

---

## Self-review notes

**Spec coverage walk-through:**

| Spec requirement (§"Implementation order" Phase 3) | Where in plan |
|---|---|
| Disk-level rules (§1) | Task 2 — 3 rules |
| Directory-entry rules (§2) | Task 3 — 9 rules |
| Sector-chain rules (§3) | Task 4 — 3 rules + walkChain helper |
| Cross-entry rules (§4) | Task 5 — 3 rules |
| §15 carryover CHAIN-SECTOR-COUNT-MINIMAL | Task 6 — 1 rule |
| "Exercise the foundation and shake out any API issues" | Task 7 — integration smoke test on corpus |
| All rules dialect-agnostic (no Dialects field) | Every Register block above has `Dialects: nil` implicit (field omitted) — Phase 3 rules apply to all dialects per the spec |

19 in-scope rules, 8 explicitly-deferred catalog entries (with rationale in the plan body), one regression gate (`TestPhase3RegistryGrowth`). Spec covered.

**Placeholder scan:** every Check function is spelled out. Every test has the exact mutation. Every commit message is given verbatim. No TBDs.

**Type / signature consistency:**

- All Check functions match `func(ctx *CheckContext) []Finding`.
- `walkChain(di *DiskImage, first *Sector) chainWalkResult` — used in `rules_chain.go` Tasks 4 and 6, and in `rules_cross.go` Task 5. Same signature each time.
- `forEachUsedSlot(ctx *CheckContext, fn func(slot int, fe *FileEntry))` — used across rules_directory.go and rules_cross.go. Same shape.
- `trackSectorRefs(ctx *CheckContext) []sectorRef` — used only in rules_disk.go (Task 2). Returns by-value Sector copies so each rule's Location can take `&ref.Sector` safely.
- `cleanSingleFileDisk(t *testing.T, name string, dataLen int) (*DiskImage, *DiskJournal)` — declared in `rules_disk_test.go` (Task 2), used in all four `*_test.go` files (Tasks 2-6). Returns the disk + its journal so tests can mutate slot 0 directly.

All consistent. No floating references.

**Rule severity sanity check (19 rules total):**

| Severity | Count | Rules at this severity in Phase 3 |
|---|---|---|
| Fatal | 4 | DISK-TRACK-SIDE-ENCODING, DISK-SECTOR-RANGE, DIR-FIRST-SECTOR-VALID, CROSS-NO-SECTOR-OVERLAP |
| Structural | 9 | DISK-DIRECTORY-TRACKS, DIR-ERASED-IS-ZERO, DIR-SECTORS-MATCHES-CHAIN, DIR-SECTORS-MATCHES-MAP, DIR-SECTORS-NONZERO, CHAIN-TERMINATOR-ZERO-ZERO, CHAIN-NO-CYCLE, CHAIN-MATCHES-SAM, CROSS-DIRECTORY-AREA-UNUSED |
| Inconsistency | 4 | DIR-TYPE-BYTE-IS-KNOWN, DIR-NAME-NOT-EMPTY, DIR-SAM-WITHIN-CAPACITY, CROSS-NO-DUPLICATE-NAMES |
| Cosmetic | 2 | DIR-NAME-PADDING, CHAIN-SECTOR-COUNT-MINIMAL |

Total: 4 + 9 + 4 + 2 = 19 ✓. Registry final count after Task 6 = 20 (Phase-1 smoke + 19 Phase-3 rules).
