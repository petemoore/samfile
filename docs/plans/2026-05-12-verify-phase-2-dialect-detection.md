# Verify Phase 2 — Dialect Detection Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace Phase 1's hardcoded `dialect := DialectUnknown` with a public `DetectDialect(*DiskImage) Dialect` that conservatively infers SAMDOS-2 / MasterDOS / SAMDOS-1 from disk content, falling back to `DialectUnknown` on any ambiguity.

**Architecture:** Two cheap independent signals — `bootFileDialect` (examines the slot whose `FirstSector == (4, 1)`: matches on name and type) and `mgtFlagsDialect` (scans every used slot's `MGTFlags` for bits outside `{0x00, 0x20}`). Each signal returns `DialectUnknown` when it has no opinion. `DetectDialect` collects the non-Unknown opinions: if they all agree, return that dialect; otherwise return `DialectUnknown`. The whole heuristic lives in one new file (`dialect.go` ≈ 90 lines) so it stays reviewable in one screen. `(*DiskImage).Verify` is changed in exactly one place to call `DetectDialect(di)` instead of the hardcoded literal.

**Tech Stack:** Go 1.22+, standard library only. The `samfile` package already exposes everything we need: `DiskJournal`, `FileEntry`, `Filename`, `Sector`, `WriteFileEntry`, `NewDiskImage`, `AddCodeFile`, `AddBasicFile`.

**Context for the engineer:**

- **Read first**, in order:
  1. `docs/specs/2026-05-11-verify-feature-design.md` — §Dialect detection (lines ~194-208) for design intent and the "conservative — return Unknown when ambiguous" mandate.
  2. `docs/disk-validity-rules.md` — §13 (`Dialect notes`) for the catalog evidence behind each signal (SAMDOS-1 type 3, MasterDOS extended MGTFlags, SAMDOS-2 vs MasterDOS BASIC gap).
  3. `verify.go` (this PR will modify it) — see how `dialect` is currently a hardcoded `DialectUnknown` in `(*DiskImage).Verify` and how `ruleAppliesToDialect` consumes it.
  4. `samfile.go:80-132` — `FileEntry` struct: `Name`, `Type`, `FirstSector`, `MGTFlags` are public fields you'll read in the heuristic.

- **What "conservative" means here:** the spec says "any ambiguity returns DialectUnknown, which causes Verify to run only rules tagged AllDialects." Concretely:
  - A clean SAMDOS-2 disk → SAMDOS2.
  - A clean MasterDOS disk → MasterDOS.
  - A clean SAMDOS-1 disk → SAMDOS1.
  - An empty disk, a data-only disk with no boot file, or a disk with conflicting signals (e.g. boot file named "samdos2" but a slot with MasterDOS MGTFlags bits set) → Unknown.

- **Signal coverage you are NOT implementing:** the BASIC `SAVARS-NVARS == 2156` MasterDOS marker (catalog `DIALECT-MASTERDOS-GAP-2156`) is deferred — it lives behind the FT_SAM_BASIC content rules in Phase 5. Phase 2 is intentionally limited to disk-level signals readable without parsing file bodies.

- **What Phase 2 does NOT change:** rule registry, Verify exit-code policy, CLI flag surface. The only Verify-level change is one line.

---

## File Structure

| Path | Action | Responsibility |
|---|---|---|
| `dialect.go` | **Create** | Houses `DetectDialect` and the two signal helpers `bootFileDialect` / `mgtFlagsDialect`. Pure functions of `*DiskImage`; no I/O, no allocations beyond a journal call. |
| `dialect_test.go` | **Create** | All unit tests for the three functions plus one integration test against `testdata/ETrackerv1.2.mgt`. |
| `verify.go` | **Modify** (one line) | Replace the hardcoded `dialect := DialectUnknown` with `dialect := DetectDialect(di)`. Update the godoc comment that currently says "Phase 1 always passes DialectUnknown". |
| `verify_test.go` | **Modify** (one comment + one assertion comment) | `TestVerifyReportCarriesDialect` keeps the empty-disk fixture (DetectDialect returns Unknown for empty disks); only the comment that calls out "Phase 1" is updated. |

---

## Task 1: Stub `DetectDialect` and wire it into `Verify`

Why this task exists: Phase 1's `Verify` hardcodes `DialectUnknown`. We need the call site in place before writing detection tests, so each subsequent task can assert on the report's `Dialect`, not on private helpers.

**Files:**
- Create: `dialect.go`
- Modify: `verify.go:270-285` (Verify's body) and the godoc comment above it.

- [ ] **Step 1: Create `dialect.go` with a stub that always returns Unknown**

```go
package samfile

// DetectDialect inspects di and returns the most likely dialect that
// wrote the disk. The heuristic combines independent signals (boot
// file at T4S1, MGTFlags bit patterns across used slots) and returns
// DialectUnknown when those signals are silent or contradict each
// other.
//
// Detection is deliberately conservative: when the result is
// DialectUnknown, Verify only runs rules tagged AllDialects, which is
// always safe. Pass --dialect=NAME on the CLI to override the result
// when the heuristic gets it wrong.
//
// Signals consulted (each returns its own DialectUnknown when it has
// no opinion; see bootFileDialect, mgtFlagsDialect):
//
//   - Boot file name and type — the slot whose FirstSector is (4, 1)
//     identifies the DOS that wrote the disk: "samdos2" → SAMDOS-2,
//     "masterdos"/"masterdos2" → MasterDOS, "samdos" or a type-3 file
//     → SAMDOS-1.
//   - MGTFlags across used slots — bits outside {0x00, 0x20} signal
//     MasterDOS (catalog: DIALECT-MASTERDOS-MGTFLAGS).
//
// Other dialect-distinguishing signals (BASIC SAVARS-NVARS gap,
// FileTypeInfo conventions) are deferred to later phases when the
// file-type rules land.
func DetectDialect(di *DiskImage) Dialect {
	return DialectUnknown
}
```

- [ ] **Step 2: Verify `dialect.go` compiles**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go build ./...`
Expected: no output, exit 0.

- [ ] **Step 3: Replace the hardcoded literal in `Verify`**

In `verify.go`, find the function `(di *DiskImage) Verify`:

```go
func (di *DiskImage) Verify() VerifyReport {
	dialect := DialectUnknown
	ctx := &CheckContext{
```

Change to:

```go
func (di *DiskImage) Verify() VerifyReport {
	dialect := DetectDialect(di)
	ctx := &CheckContext{
```

Also update the godoc comment immediately above the function. The current comment ends with:

```go
// In Phase 1, dialect detection is not yet implemented and Verify
// always passes DialectUnknown to rules; rules whose Dialects slice
// is non-empty and excludes DialectUnknown are skipped. Phase 2
// adds DetectDialect.
```

Replace those four lines with:

```go
// Verify calls DetectDialect to infer the dialect that wrote di,
// then runs every registered rule whose Dialects slice is empty
// (all-dialects) or contains the detected dialect. Rules scoped to a
// dialect other than the one detected are skipped. DetectDialect is
// conservative: when it returns DialectUnknown (empty or ambiguous
// disks), only all-dialects rules run.
```

- [ ] **Step 4: Update the Phase-1 comments in `verify_test.go`**

In `verify_test.go`, find `TestVerifyRespectsDialectScoping`:

```go
	di := NewDiskImage()
	di.Verify() // Phase 1 always passes DialectUnknown
```

Change to:

```go
	di := NewDiskImage()
	di.Verify() // empty disk: DetectDialect returns DialectUnknown
```

And in `TestVerifyReportCarriesDialect`:

```go
	di := NewDiskImage()
	report := di.Verify()
	// Phase 1: dialect detection is not implemented; always DialectUnknown.
	if report.Dialect != DialectUnknown {
		t.Errorf("Dialect = %v; want unknown (detection lands in Phase 2)", report.Dialect)
	}
```

Change the comment + message to:

```go
	di := NewDiskImage()
	report := di.Verify()
	// Empty disk has no signals; DetectDialect returns Unknown.
	if report.Dialect != DialectUnknown {
		t.Errorf("Dialect = %v; want unknown for empty disk", report.Dialect)
	}
```

- [ ] **Step 5: Run existing tests to verify no regression**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test ./...`
Expected: all tests pass. Specifically `TestVerifyReportCarriesDialect`, `TestVerifyRespectsDialectScoping`, and `TestVerifyRunsRegisteredRules` continue to pass because the stub returns Unknown, which is the same value the hardcoded literal supplied.

- [ ] **Step 6: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-2 && \
g add dialect.go verify.go verify_test.go && \
g commit -m "verify: add DetectDialect stub and wire into Verify

Phase 2 step 1: introduce the DetectDialect entry point with a
conservative stub that returns DialectUnknown. Verify now calls
it instead of using a hardcoded literal, and the Phase-1 godoc
comment is updated accordingly. Behaviour is unchanged because
the stub returns the same value the hardcoded literal supplied;
real signal heuristics arrive in the following tasks."
```

(Use `g` not `git`; Pete's alias preserves authorship timestamps.)

---

## Task 2: Empty-disk and disk-with-unknown-boot-file return Unknown

Why this task exists: lock in the conservative-fallback behaviour with explicit tests before adding any real signal logic. These tests stay green for the rest of the plan.

**Files:**
- Create: `dialect_test.go`

- [ ] **Step 1: Write failing tests for empty and unknown-boot-file disks**

Create `dialect_test.go`:

```go
package samfile

import (
	"os"
	"testing"
)

func TestDetectDialectEmptyDisk(t *testing.T) {
	di := NewDiskImage()
	if got := DetectDialect(di); got != DialectUnknown {
		t.Errorf("DetectDialect(empty) = %v; want unknown", got)
	}
}

func TestDetectDialectUnknownBootFileName(t *testing.T) {
	// A disk whose first file is named something neither DOS recognises
	// and whose MGTFlags are vanilla 0 (AddCodeFile leaves MGTFlags at
	// zero) emits no signal. DetectDialect must return Unknown rather
	// than guessing.
	di := NewDiskImage()
	if err := di.AddCodeFile("BOOTER", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	if got := DetectDialect(di); got != DialectUnknown {
		t.Errorf("DetectDialect(unknown boot file) = %v; want unknown", got)
	}
}
```

- [ ] **Step 2: Run the new tests to confirm they pass against the stub**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test -run 'TestDetectDialect(EmptyDisk|UnknownBootFileName)' ./...`
Expected: both PASS (because the stub returns Unknown).

The tests must keep passing as the heuristic grows. They are not "failing tests" in the strict TDD sense — they are *regression guards* establishing the conservative-fallback invariant before signal logic that could accidentally violate it.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-2 && \
g add dialect_test.go && \
g commit -m "verify: lock conservative fallback for DetectDialect

Add regression tests asserting DetectDialect returns Unknown for
(a) an empty disk and (b) a disk whose only file has neither a
recognised boot-file name nor extended MGTFlags. These pin the
'no opinion → Unknown' invariant before the signal heuristics
land in the next commits."
```

---

## Task 3: Boot-file-name signal — SAMDOS-2

Why this task exists: the first signal is the simplest: examine the slot whose `FirstSector == (4, 1)` (the bootable sector ROM BOOTEX reads). If its filename trims+lowercases to `"samdos2"`, the disk was written by SAMDOS-2. Catalog evidence: `docs/disk-validity-rules.md` §11 (`BOOT-OWNER-AT-T4S1`) plus the canonical samdos2 binary shipped as the slot-0 file in this project's own M0 boot disk.

**Files:**
- Modify: `dialect.go` — add `bootFileDialect` helper and wire it into `DetectDialect`.
- Modify: `dialect_test.go` — add positive test.

- [ ] **Step 1: Write the failing test**

Append to `dialect_test.go`:

```go
func TestDetectDialectSamdos2BootFile(t *testing.T) {
	// First file added → allocated at FirstSector (4, 1). Name is the
	// canonical samdos2 filename. The body content does not matter for
	// detection — only the slot name does.
	di := NewDiskImage()
	if err := di.AddCodeFile("samdos2", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	if got := DetectDialect(di); got != DialectSAMDOS2 {
		t.Errorf("DetectDialect(samdos2 boot file) = %v; want samdos2", got)
	}
}
```

- [ ] **Step 2: Run the test to confirm it fails**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test -run TestDetectDialectSamdos2BootFile ./...`
Expected: FAIL with `DetectDialect(samdos2 boot file) = unknown; want samdos2` (stub still returns Unknown).

- [ ] **Step 3: Implement `bootFileDialect` and call it from `DetectDialect`**

Replace the body of `dialect.go` with:

```go
package samfile

import "strings"

// DetectDialect inspects di and returns the most likely dialect that
// wrote the disk. The heuristic combines independent signals (boot
// file at T4S1, MGTFlags bit patterns across used slots) and returns
// DialectUnknown when those signals are silent or contradict each
// other.
//
// Detection is deliberately conservative: when the result is
// DialectUnknown, Verify only runs rules tagged AllDialects, which is
// always safe. Pass --dialect=NAME on the CLI to override the result
// when the heuristic gets it wrong.
//
// Signals consulted (each returns its own DialectUnknown when it has
// no opinion; see bootFileDialect, mgtFlagsDialect):
//
//   - Boot file name and type — the slot whose FirstSector is (4, 1)
//     identifies the DOS that wrote the disk: "samdos2" → SAMDOS-2,
//     "masterdos"/"masterdos2" → MasterDOS, "samdos" or a type-3 file
//     → SAMDOS-1.
//   - MGTFlags across used slots — bits outside {0x00, 0x20} signal
//     MasterDOS (catalog: DIALECT-MASTERDOS-MGTFLAGS).
//
// Other dialect-distinguishing signals (BASIC SAVARS-NVARS gap,
// FileTypeInfo conventions) are deferred to later phases when the
// file-type rules land.
func DetectDialect(di *DiskImage) Dialect {
	dj := di.DiskJournal()
	opinions := []Dialect{
		bootFileDialect(dj),
		mgtFlagsDialect(dj),
	}
	var picked Dialect = DialectUnknown
	for _, o := range opinions {
		if o == DialectUnknown {
			continue
		}
		if picked == DialectUnknown {
			picked = o
			continue
		}
		if picked != o {
			return DialectUnknown // conflict → conservative
		}
	}
	return picked
}

// bootFileDialect examines the slot whose FirstSector is (track 4,
// sector 1) — the sector ROM BOOTEX reads to &8000 (see catalog
// BOOT-OWNER-AT-T4S1). The slot's filename (trimmed, lowercased) and
// masked Type are matched against the canonical DOS bootstraps:
//
//   - "samdos2" or "samdos 2"      → DialectSAMDOS2
//   - "masterdos" or "masterdos2"  → DialectMasterDOS
//   - "samdos" (no trailing 2), or masked Type == 3
//                                  → DialectSAMDOS1
//
// Anything else (including no used slot at T4S1) returns
// DialectUnknown — the signal abstains rather than guesses.
func bootFileDialect(dj *DiskJournal) Dialect {
	for _, fe := range dj {
		if fe == nil || !fe.Used() {
			continue
		}
		if fe.FirstSector == nil ||
			fe.FirstSector.Track != 4 ||
			fe.FirstSector.Sector != 1 {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(fe.Name.String()))
		switch name {
		case "samdos2", "samdos 2":
			return DialectSAMDOS2
		case "masterdos", "masterdos2":
			return DialectMasterDOS
		case "samdos":
			return DialectSAMDOS1
		}
		if uint8(fe.Type)&0x1F == 3 {
			// SAMDOS-1's "auto-include header" sets type 3 on the
			// bootstrap itself (samdos/src/b.s:14-22). Type 3 is
			// otherwise a DIR alias for "ZX $.ARRAY"; restricting
			// this check to the boot slot keeps it unambiguous.
			return DialectSAMDOS1
		}
		return DialectUnknown
	}
	return DialectUnknown
}

// mgtFlagsDialect scans every used slot's MGTFlags. A bit outside the
// SAMDOS-2 set {0x00, 0x20} signals MasterDOS (catalog:
// DIALECT-MASTERDOS-MGTFLAGS, sourced from real-disk observation —
// SAMDOS source has no MGTFlags reader). Filled in by a later task.
func mgtFlagsDialect(dj *DiskJournal) Dialect {
	return DialectUnknown
}
```

- [ ] **Step 4: Run the test to confirm it passes**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test -run TestDetectDialectSamdos2BootFile ./...`
Expected: PASS.

- [ ] **Step 5: Re-run the regression suite**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test ./...`
Expected: all tests pass — empty-disk and unknown-boot-file tests from Task 2 still return Unknown.

- [ ] **Step 6: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-2 && \
g add dialect.go dialect_test.go && \
g commit -m "verify: detect SAMDOS-2 from boot-file name at T4S1

Add the first DetectDialect signal: examine the slot whose
FirstSector is (4, 1) and recognise canonical DOS filenames.
Filename 'samdos2' → DialectSAMDOS2. MasterDOS and SAMDOS-1
branches are sketched in the same switch and exercised by
later tasks. mgtFlagsDialect is stubbed."
```

---

## Task 4: Boot-file-name signal — MasterDOS

Why this task exists: extend the same helper to recognise MasterDOS's bootstrap filename. The branch already exists in `bootFileDialect` from Task 3; this task adds the test that exercises it.

**Files:**
- Modify: `dialect_test.go`

- [ ] **Step 1: Write the failing test**

Append to `dialect_test.go`:

```go
func TestDetectDialectMasterDOSBootFile(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("masterdos2", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	if got := DetectDialect(di); got != DialectMasterDOS {
		t.Errorf("DetectDialect(masterdos2 boot file) = %v; want masterdos", got)
	}
}
```

- [ ] **Step 2: Run the test**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test -run TestDetectDialectMasterDOSBootFile ./...`
Expected: PASS (Task 3's `bootFileDialect` already handles `masterdos2`).

If the test fails, the cause is a typo in the switch in `bootFileDialect`; fix it.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-2 && \
g add dialect_test.go && \
g commit -m "verify: cover MasterDOS boot-file branch in DetectDialect

Exercise the 'masterdos2' name match introduced in the previous
commit. No production change."
```

---

## Task 5: Boot-file-name signal — SAMDOS-1 (name + type-3 paths)

Why this task exists: SAMDOS-1's bootstrap is the trickiest of the three — it can be identified either by filename `"samdos"` (no trailing 2) or by the type-3 "auto-include header" the older SAMDOS variant emits for itself (`samdos/src/b.s:14-22`). Two tests, one for each path. The production code already handles both from Task 3.

**Files:**
- Modify: `dialect_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `dialect_test.go`:

```go
func TestDetectDialectSAMDOS1ByName(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("samdos", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	if got := DetectDialect(di); got != DialectSAMDOS1 {
		t.Errorf("DetectDialect(samdos boot file) = %v; want samdos1", got)
	}
}

func TestDetectDialectSAMDOS1ByType3(t *testing.T) {
	// A bootstrap with an unrecognised filename but masked type 3 is
	// SAMDOS-1's auto-include header (samdos/src/b.s:14-22). Use
	// AddCodeFile, then patch Type to FT(3) via a journal write.
	di := NewDiskImage()
	if err := di.AddCodeFile("oddname", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj.DiskJournal[0].Type = FileType(3)
	di.WriteFileEntry(dj, 0)

	if got := DetectDialect(di); got != DialectSAMDOS1 {
		t.Errorf("DetectDialect(type-3 boot file) = %v; want samdos1", got)
	}
}
```

- [ ] **Step 2: Run the tests**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test -run 'TestDetectDialectSAMDOS1' ./...`
Expected: both PASS.

If `TestDetectDialectSAMDOS1ByType3` fails because the journal mutation does not stick, double-check: `di.WriteFileEntry(dj, 0)` re-serialises slot 0 from `dj.DiskJournal[0]` back into `di`'s bytes. After the call, a fresh `di.DiskJournal()` will read the new type back. If the assertion still fails, the bug is in `bootFileDialect`'s type-3 branch — verify `uint8(fe.Type) & 0x1F` matches `3` for a `FileType(3)` value.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-2 && \
g add dialect_test.go && \
g commit -m "verify: cover SAMDOS-1 boot-file branches (name + type-3)

Exercise the two paths bootFileDialect uses for SAMDOS-1:
filename 'samdos' (no trailing 2) and masked type byte 3 (the
auto-include-header variant from samdos/src/b.s:14-22). No
production change — both branches were sketched in Task 3."
```

---

## Task 6: MGTFlags signal — MasterDOS

Why this task exists: SAMDOS-2 BASIC files have `MGTFlags == 0x20`; SAMDOS-2 CODE files leave it at `0x00`; MasterDOS sets additional bits beyond `0x20`. Any used slot with `MGTFlags & ^0x20 != 0` is therefore a MasterDOS signal (catalog: `DIALECT-MASTERDOS-MGTFLAGS`, §13). This task implements `mgtFlagsDialect` and tests it in isolation, then asserts `DetectDialect` returns MasterDOS when only this signal fires.

**Files:**
- Modify: `dialect.go` — flesh out `mgtFlagsDialect`.
- Modify: `dialect_test.go` — add positive and negative tests.

- [ ] **Step 1: Write the failing test**

Append to `dialect_test.go`:

```go
func TestDetectDialectMasterDOSByMGTFlags(t *testing.T) {
	// AddCodeFile leaves MGTFlags at 0x00 (vanilla SAMDOS-2 CODE
	// convention). Patch MGTFlags to 0x80 — an extended bit outside
	// {0x00, 0x20} — and DetectDialect must report MasterDOS.
	di := NewDiskImage()
	if err := di.AddCodeFile("data", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj.DiskJournal[0].MGTFlags = 0x80
	di.WriteFileEntry(dj, 0)

	if got := DetectDialect(di); got != DialectMasterDOS {
		t.Errorf("DetectDialect(MGTFlags=0x80) = %v; want masterdos", got)
	}
}

func TestMGTFlagsDialectVanillaIsSilent(t *testing.T) {
	// A disk where every used slot has MGTFlags in {0x00, 0x20}
	// (vanilla SAMDOS-2) yields no opinion from mgtFlagsDialect.
	di := NewDiskImage()
	if err := di.AddCodeFile("CODE", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (CODE, MGTFlags=0): %v", err)
	}
	// AddBasicFile sets MGTFlags=0x20 — exercise both bytes of the
	// SAMDOS-2 set.
	// (We do not need a real BASIC body; patch the second slot's
	// MGTFlags directly to keep the test minimal.)
	if err := di.AddCodeFile("BASIC", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (BASIC stub): %v", err)
	}
	dj := di.DiskJournal()
	dj.DiskJournal[1].MGTFlags = 0x20
	di.WriteFileEntry(dj, 1)

	if got := mgtFlagsDialect(di.DiskJournal()); got != DialectUnknown {
		t.Errorf("mgtFlagsDialect(vanilla MGTFlags) = %v; want unknown", got)
	}
}
```

- [ ] **Step 2: Run the tests to confirm the positive case fails**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test -run 'TestDetectDialectMasterDOSByMGTFlags|TestMGTFlagsDialectVanillaIsSilent' ./...`
Expected: `TestDetectDialectMasterDOSByMGTFlags` FAILS (the helper still returns Unknown); `TestMGTFlagsDialectVanillaIsSilent` PASSES (stub returns Unknown).

- [ ] **Step 3: Implement `mgtFlagsDialect`**

In `dialect.go`, replace the stub body:

```go
func mgtFlagsDialect(dj *DiskJournal) Dialect {
	return DialectUnknown
}
```

with:

```go
// mgtFlagsDialect scans every used slot's MGTFlags. A bit outside the
// SAMDOS-2 set {0x00, 0x20} signals MasterDOS (catalog:
// DIALECT-MASTERDOS-MGTFLAGS). Real-disk observation: MasterDOS sets
// per-file attribute bits beyond 0x20 to track its own metadata,
// while SAMDOS-2 leaves MGTFlags at either 0x00 (CODE) or 0x20 (BASIC).
// Returns DialectUnknown when every used slot's MGTFlags is in the
// SAMDOS-2 set, including the trivial empty-disk case.
func mgtFlagsDialect(dj *DiskJournal) Dialect {
	const samdos2Mask uint8 = ^uint8(0x20) // bits the SAMDOS-2 set ignores
	for _, fe := range dj {
		if fe == nil || !fe.Used() {
			continue
		}
		if fe.MGTFlags&samdos2Mask != 0 {
			return DialectMasterDOS
		}
	}
	return DialectUnknown
}
```

Note the mask: `samdos2Mask = ^0x20 = 0xDF`. Any bit set in MGTFlags that lies inside that mask (i.e. is anything other than the bit 0x20) trips MasterDOS. So `0x00` and `0x20` are both silent; `0x80`, `0x40`, `0x01`, `0x21`, `0xA0` all signal MasterDOS.

- [ ] **Step 4: Run the tests to confirm they pass**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test -run 'TestDetectDialectMasterDOSByMGTFlags|TestMGTFlagsDialectVanillaIsSilent' ./...`
Expected: both PASS.

- [ ] **Step 5: Re-run the regression suite**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test ./...`
Expected: all tests pass.

- [ ] **Step 6: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-2 && \
g add dialect.go dialect_test.go && \
g commit -m "verify: detect MasterDOS from extended MGTFlags bits

Implement mgtFlagsDialect: scan every used slot's MGTFlags and
report MasterDOS if any bit outside {0x00, 0x20} is set. The
SAMDOS-2 set is justified by AddCodeFile leaving MGTFlags at
0x00 and AddBasicFile setting it to 0x20; bits beyond that
range are MasterDOS per catalog DIALECT-MASTERDOS-MGTFLAGS.

Vanilla SAMDOS-2 disks remain silent so they can still be
identified by the boot-file signal alone."
```

---

## Task 7: Combine signals — agreement and conflict cases

Why this task exists: the previous tasks each exercised one signal in isolation. This task pins down `DetectDialect`'s combination logic: agreeing signals reinforce each other; conflicting signals collapse to Unknown.

**Files:**
- Modify: `dialect_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `dialect_test.go`:

```go
func TestDetectDialectMasterDOSBothSignalsAgree(t *testing.T) {
	// Boot file "masterdos2" + extended MGTFlags on a second slot —
	// two signals both point at MasterDOS.
	di := NewDiskImage()
	if err := di.AddCodeFile("masterdos2", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (boot): %v", err)
	}
	if err := di.AddCodeFile("payload", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (payload): %v", err)
	}
	dj := di.DiskJournal()
	dj.DiskJournal[1].MGTFlags = 0x40
	di.WriteFileEntry(dj, 1)

	if got := DetectDialect(di); got != DialectMasterDOS {
		t.Errorf("DetectDialect(both signals masterdos) = %v; want masterdos", got)
	}
}

func TestDetectDialectConflictReturnsUnknown(t *testing.T) {
	// Boot file says SAMDOS-2 but a later slot's MGTFlags say
	// MasterDOS. DetectDialect must collapse to Unknown rather than
	// pick a winner.
	di := NewDiskImage()
	if err := di.AddCodeFile("samdos2", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (boot): %v", err)
	}
	if err := di.AddCodeFile("payload", []byte{0xC9}, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile (payload): %v", err)
	}
	dj := di.DiskJournal()
	dj.DiskJournal[1].MGTFlags = 0x80
	di.WriteFileEntry(dj, 1)

	if got := DetectDialect(di); got != DialectUnknown {
		t.Errorf("DetectDialect(conflict samdos2 vs masterdos) = %v; want unknown", got)
	}
}
```

- [ ] **Step 2: Run the tests**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test -run 'TestDetectDialect(MasterDOSBothSignalsAgree|ConflictReturnsUnknown)' ./...`
Expected: both PASS. The combination logic in `DetectDialect` (from Task 3) already handles both cases — agreement reinforces, conflict collapses to Unknown. These are regression tests, not new implementation.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-2 && \
g add dialect_test.go && \
g commit -m "verify: cover signal-agreement and signal-conflict paths

Two new tests: (1) boot-file name and MGTFlags both pointing at
MasterDOS return MasterDOS; (2) boot-file 'samdos2' combined
with a MasterDOS-style MGTFlags bit on a later slot collapses
to Unknown. Exercises DetectDialect's combination logic
end-to-end. No production change."
```

---

## Task 8: Integration test — real-world testdata disk does not panic

Why this task exists: every previous test fabricates disks in memory. We need one test that drives `DetectDialect` through `Load` on a real `.mgt` image to catch panics, nil-pointer paths, or unexpected dialect values that only show up against real-world byte layouts. The repository already commits `testdata/ETrackerv1.2.mgt`.

The assertion shape is "result is a valid Dialect value" rather than "result is SAMDOS2", because we have no out-of-band knowledge of which DOS authored ETrackerv1.2 — and the point of the test is robustness, not classifying that specific image.

**Files:**
- Modify: `dialect_test.go`

- [ ] **Step 1: Confirm the corpus file exists**

Run: `ls /Users/pmoore/git/samfile-verify-phase-2/testdata/ETrackerv1.2.mgt && stat -f '%z bytes' /Users/pmoore/git/samfile-verify-phase-2/testdata/ETrackerv1.2.mgt`
Expected: prints the file path, then `819200 bytes`. If the file is missing or not 819200 bytes, stop and report — the test cannot run.

- [ ] **Step 2: Write the integration test**

Append to `dialect_test.go`:

```go
func TestDetectDialectETrackerCorpus(t *testing.T) {
	// Smoke test against a real-world MGT image. We do not assert a
	// specific dialect — we just assert DetectDialect returns one of
	// the four documented values without panicking. This protects
	// against nil-pointer paths in bootFileDialect / mgtFlagsDialect
	// that fabricated disks might not exercise.
	const path = "testdata/ETrackerv1.2.mgt"
	if _, err := os.Stat(path); err != nil {
		t.Skipf("corpus image not present (%v); skipping", err)
	}
	di, err := Load(path)
	if err != nil {
		t.Fatalf("Load(%q): %v", path, err)
	}
	got := DetectDialect(di)
	switch got {
	case DialectUnknown, DialectSAMDOS1, DialectSAMDOS2, DialectMasterDOS:
		// All four are acceptable; log for diagnostic value.
		t.Logf("DetectDialect(%s) = %s", path, got)
	default:
		t.Errorf("DetectDialect(%s) = %v; not a documented Dialect value", path, got)
	}
}
```

- [ ] **Step 3: Run the test**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test -run TestDetectDialectETrackerCorpus -v ./...`
Expected: PASS, with a log line showing the detected dialect (e.g. `DetectDialect(testdata/ETrackerv1.2.mgt) = unknown`).

If the test panics: the panic site is the bug. Most likely cause is a nil `fe.FirstSector` for a used slot whose first-sector bytes are zero (a corrupted real-world disk) — `bootFileDialect` already guards against this with `fe.FirstSector == nil`, but double-check.

- [ ] **Step 4: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-2 && \
g add dialect_test.go && \
g commit -m "verify: smoke-test DetectDialect against committed corpus

Run DetectDialect on testdata/ETrackerv1.2.mgt and assert the
result is one of the four documented Dialect values. The test
does not assert a specific dialect — its job is to catch
panics or out-of-range returns that only show up against real-
world byte layouts, not fabricated in-memory disks."
```

---

## Task 9: Final verification — full test suite, vet, and a manual CLI smoke run

Why this task exists: belt-and-braces check before opening the PR. Pete's standing rule (`memory/feedback_correctness_over_workarounds.md`) is to verify the change actually works end-to-end, not just that unit tests pass. The CLI `samfile verify` is the user-facing surface that exposes the detected dialect; run it once against the corpus image and confirm the printed `detected dialect:` line reflects the new heuristic.

**Files:** none modified.

- [ ] **Step 1: Run the full test suite**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go test ./...`
Expected: all tests PASS, including the seven new ones added across Tasks 2-8.

- [ ] **Step 2: Run `go vet`**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go vet ./...`
Expected: no output, exit 0.

- [ ] **Step 3: Build the CLI**

Run: `cd /Users/pmoore/git/samfile-verify-phase-2 && go build -o /tmp/samfile ./cmd/samfile`
Expected: no output, exit 0; `/tmp/samfile` exists.

- [ ] **Step 4: Run the verify subcommand against the corpus image**

Run: `/tmp/samfile verify -i /Users/pmoore/git/samfile-verify-phase-2/testdata/ETrackerv1.2.mgt`
Expected: output begins with:

```
samfile verify: results for /Users/pmoore/git/samfile-verify-phase-2/testdata/ETrackerv1.2.mgt
detected dialect: <one of: unknown | samdos1 | samdos2 | masterdos>
```

followed by the Phase 1 smoke rule (`DISK-NOT-EMPTY`) report — which should NOT fire because the corpus image has files. If the dialect line is anything other than one of the four documented strings, stop and investigate before opening the PR.

- [ ] **Step 5: Confirm the M0 boot disk classifies as SAMDOS-2 (optional sanity check)**

The user's M0 boot disk under `/Users/pmoore/git/sam-aarch64/build/test.mgt` (if present) contains a slot-0 file named `samdos2`. Running verify against it should report `detected dialect: samdos2`.

Run: `[ -f /Users/pmoore/git/sam-aarch64/build/test.mgt ] && /tmp/samfile verify -i /Users/pmoore/git/sam-aarch64/build/test.mgt | head -3 || echo "no M0 disk available, skipping"`
Expected: either a `samdos2` dialect line or the skip message. If `build/test.mgt` exists and reports anything other than `samdos2`, that's a strong signal the heuristic is wrong — stop and investigate.

- [ ] **Step 6: Push the branch**

```bash
cd /Users/pmoore/git/samfile-verify-phase-2 && g push -u origin feat/verify-phase-2-dialect-detection
```

Expected: branch pushed to GitHub, tracking set up. Note the URL gh emits for opening a PR.

- [ ] **Step 7: Open a DRAFT PR**

```bash
cd /Users/pmoore/git/samfile-verify-phase-2 && gh pr create --draft \
  --base master \
  --title "verify: Phase 2 — DetectDialect (boot-file + MGTFlags signals)" \
  --body "$(cat <<'EOF'
## Summary

Phase 2 of the `samfile verify` rollout (spec: `docs/specs/2026-05-11-verify-feature-design.md`, plan: `docs/plans/2026-05-12-verify-phase-2-dialect-detection.md`). Replaces Phase 1's hardcoded `dialect := DialectUnknown` with a public `DetectDialect(*DiskImage) Dialect`.

Two cheap, independent signals are combined conservatively — when they disagree, `DetectDialect` returns `DialectUnknown` rather than pick a winner:

- **Boot-file signal** (`bootFileDialect`) — the slot whose `FirstSector == (4, 1)` is matched on name (`samdos2`, `masterdos`/`masterdos2`, `samdos`) and on masked Type (`3` ⇒ SAMDOS-1's auto-include-header).
- **MGTFlags signal** (`mgtFlagsDialect`) — any used slot with `MGTFlags & 0xDF != 0` (anything outside `{0x00, 0x20}`) ⇒ MasterDOS, per catalog rule `DIALECT-MASTERDOS-MGTFLAGS`.

Deferred to later phases (out of scope here):

- BASIC `SAVARS-NVARS == 2156` MasterDOS signal — needs Phase 5's FT_SAM_BASIC rules.
- Confidence levels and `--dialect` override on the CLI — design spec §"Open questions deferred to plan-writing".

## Test plan

- [ ] `go test ./...` passes (8 new DetectDialect tests + existing suite green)
- [ ] `go vet ./...` clean
- [ ] CLI smoke: `samfile verify -i testdata/ETrackerv1.2.mgt` prints a valid `detected dialect:` line
- [ ] (Optional) `samfile verify -i ../sam-aarch64/build/test.mgt` reports `samdos2` if the M0 disk is present

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

Expected: `gh` prints a PR URL. Note it for handoff. Pete reviews; CI runs.

- [ ] **Step 8: Monitor CI to completion**

Per Pete's standing PR-workflow rule (`~/.claude/CLAUDE.md`): watch every check until it finishes — GitHub Actions, Taskcluster Decision Tasks, anything else reporting status. Fix any failures autonomously (small corrections amend into the relevant commit; design questions escalate to Pete).

Run: `gh pr checks` periodically (or `gh pr checks --watch` if available) until all required checks are green.

If `gh pr checks` shows failures: diagnose locally with the same `go test` / `go vet` commands. Iterate fixes in the dev container if they're CI-environment-specific (Pete's standing rule: `iterate-CI-fixes-locally`).

- [ ] **Step 9: Hand off**

Reply to Pete with: the PR URL, the CI status, the corpus-image detected dialect (from Step 4), and any noteworthy decisions made during the run. Do NOT mark the PR ready for review without explicit confirmation from Pete — drafts only until he approves.

---

## Self-review notes

**Spec coverage walk-through (`docs/specs/2026-05-11-verify-feature-design.md` §Dialect detection):**

| Spec requirement | Where in plan |
|---|---|
| `DetectDialect(*DiskImage) Dialect` exists as public API | Task 1 |
| Inspects boot-sector presence/contents | Task 3, 4, 5 (bootFileDialect) |
| Inspects MGT future-and-past patterns | NOT YET — catalog rule §13 is empirical, not citation-backed; deferred (noted in DetectDialect godoc) |
| Inspects MasterDOS-only dir-entry field usage | Task 6 (mgtFlagsDialect) |
| Conservative: ambiguity → DialectUnknown | Task 7 (conflict test pinning behaviour) |
| Verify uses DetectDialect's result | Task 1 (one-line change in verify.go) |
| Phase-2 PR against samfile master | Task 9 step 7 |

The "MGT future-and-past patterns" signal is intentionally out of scope. Catalog §13 has no concrete bit patterns to check against — only the prose hint "MasterDOS sets bits beyond 0x20 in MGTFlags". The mgtFlagsDialect signal covers exactly that. Adding speculative `MGTFutureAndPast` bit checks would weaken the heuristic.

**Placeholder scan:** every step has concrete code, concrete commands, and concrete expected output. No TBDs.

**Type / signature consistency:**

- `DetectDialect(di *DiskImage) Dialect` — same signature in stub (Task 1), godoc (Task 3), CLI consumer (no change), all tests (Tasks 2-8).
- `bootFileDialect(dj *DiskJournal) Dialect` — same signature in Task 3 implementation and Task 6 reference.
- `mgtFlagsDialect(dj *DiskJournal) Dialect` — same signature in Task 3 stub and Task 6 implementation.
- `FileEntry.MGTFlags` is `uint8` — Task 6 mask `^uint8(0x20)` is correctly typed.
- `FileEntry.FirstSector` is `*Sector` — guarded with `fe.FirstSector == nil` in `bootFileDialect`.
- `FileType` is the type of `FileEntry.Type` (`samfile.go:93`); `uint8(fe.Type) & 0x1F` is the correct way to mask attribute bits per catalog `DIR-TYPE-MASKING`.

All consistent.
