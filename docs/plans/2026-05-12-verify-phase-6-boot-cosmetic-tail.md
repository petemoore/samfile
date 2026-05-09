# Verify Phase 6 — Boot-File & Cosmetic-Tail Rules

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the 4 remaining substantively-checkable rules from the catalog: 3 boot-file rules (§11) and 1 cosmetic ReservedA rule (§14). After this lands the registry holds 51 rules total and the catalog is closed for v1 implementation; Phase 7's corpus-validation pass classifies empirical violation rates and reclassifies severities as needed.

**Architecture:** Two new files per the established Phase 3-5 convention: `rules_boot.go` (§11, 3 rules) and `rules_cosmetic.go` (§14, 1 rule). The three boot rules share a common shape — find the slot whose `FirstSector == (4, 1)`, then inspect the sector contents — so they share a private helper `bootSlot(*DiskJournal) (slot int, fe *FileEntry, found bool)`. No new global helpers; everything else reuses existing infrastructure (`forEachUsedSlot`, `cleanSingleFileDisk`).

**Tech Stack:** Go 1.22+. No new dependencies.

**Context for the engineer:**

Read these first, in order:

1. `docs/specs/2026-05-11-verify-feature-design.md` §"Implementation order" Phase 6: "remaining ~16 rules. After this lands, the catalog is fully realised."
2. `docs/disk-validity-rules.md` §11 (boot-file rules) and §14 (cosmetic / canonical-output rules).
3. `rules_ft_code.go` from Phase 4 — the pattern of "filter on `fe.Type == X`" applies here too where relevant.
4. `samfile.go:597-605` — `(*FileEntry).Used` so you understand which slots `forEachUsedSlot` visits.
5. `samfile.go:388-400` — `(*DiskImage).SectorData` is what the boot rules call to read T4S1.

**Phase 6 scope: 4 rules.** The catalog has 14 entries across §11-§15; ten of them are deferred as documentation or already covered by earlier phases:

| Catalog rule | Phase 6 status |
|---|---|
| `BOOT-OWNER-AT-T4S1` | IMPLEMENT |
| `BOOT-SIGNATURE-AT-256` | IMPLEMENT |
| `BOOT-ENTRY-POINT-AT-9` | IMPLEMENT (cosmetic — heuristic) |
| `BOOT-FILE-TYPE-IGNORED` | DEFER — catalog explicitly says "not a validity check; a note for verify to not flag" |
| `ATTR-HIDDEN-NOT-LISTED` | DEFER — `samfile ls --hidden` hint, not a validity check |
| `ATTR-PROTECTED-NO-OVERWRITE` | DEFER — semantic, not a validity check |
| `ATTR-ERASED-SUPPRESSES-ALL` | DEFER — meta-rule about how `forEachUsedSlot` orchestrates; already implemented in Phase 3 |
| `DIALECT-MASTERDOS-MGTFLAGS` | DEFER — drives Phase 2's `mgtFlagsDialect`; no separate rule needed |
| `DIALECT-MASTERDOS-GAP-2156` | DEFER — drives Phase 5's `BASIC-VARS-GAP-INVARIANT` |
| `DIALECT-SAMDOS-1-TYPE-3` | DEFER — documentation (don't flag type 3 as invalid) |
| `DIALECT-HOOK-128-DEAD-CODE` | DEFER — documentation (do not treat missing AUTO as bootability issue) |
| `COSMETIC-RESERVEDA-FF` | IMPLEMENT |
| `COSMETIC-RESERVEDB-FILL` | DEFER — catalog explicitly says "not a rule; document only" |
| `COSMETIC-STARTPAGE-DECORATIVE-BITS` | DEFER — comparison convention for byte-perfect diffing |

After Task 4 the registry holds 51 rules total (1 smoke + 19 Phase-3 + 15 Phase-4 + 12 Phase-5 + 4 Phase-6).

**Phase 6 standing rules** (same as Phase 3-5):

- Use `g` not plain `git` for commits.
- Every rule's `Citation` is a real `file:line`; copy verbatim from the plan.
- Each rule ships with positive + negative tests.
- Draft PR only.
- All rules use `Dialects: nil` (apply to all dialects).

---

## File Structure

| Path | Action | Responsibility |
|---|---|---|
| `rules_boot.go` | Create | §11 boot-file rules (3 rules) + `bootSlot` helper. |
| `rules_boot_test.go` | Create | Positive + negative tests for §11 (6 tests). |
| `rules_cosmetic.go` | Create | §14 cosmetic rule (1 rule). |
| `rules_cosmetic_test.go` | Create | Positive + negative tests for §14 (2-3 tests, including a regression for the dual-convention acceptance). |
| `rules_smoke_test.go` | Modify | `TestRegistryGrowth` count update 47 → 51. |

---

## The boot-slot helper

Add at the top of `rules_boot.go`:

```go
package samfile

import "fmt"

// bootSlot returns the (slot index, FileEntry) of the disk's boot file
// — the used slot whose FirstSector is (track 4, sector 1). Returns
// found=false when no used slot owns T4S1 (a non-bootable disk).
// BOOT-OWNER-AT-T4S1 produces a finding in that case; the other two
// §11 rules silently skip (their checks are conditional on a boot
// file existing).
func bootSlot(dj *DiskJournal) (slot int, fe *FileEntry, found bool) {
	for idx, e := range dj {
		if e == nil || !e.Used() {
			continue
		}
		if e.FirstSector != nil && e.FirstSector.Track == 4 && e.FirstSector.Sector == 1 {
			return idx, e, true
		}
	}
	return -1, nil, false
}
```

---

## Task 1: Skeleton + registry-growth gate update

**Files:**
- Create: `rules_boot.go`, `rules_cosmetic.go` (skeletons).
- Modify: `rules_smoke_test.go` — update `TestRegistryGrowth` count to 51.

- [ ] **Step 1: Create the rule-file skeletons**

`rules_boot.go` gets the section comment, `import "fmt"`, AND the `bootSlot` helper.

```go
// rules_boot.go
package samfile

import "fmt"

// §11 Boot-file rules (catalog docs/disk-validity-rules.md §11).
// Rules in this file check that a disk's boot file (the slot whose
// FirstSector is at track 4, sector 1) carries the bytes ROM BOOTEX
// expects: a "BOOT" signature at offset 256-259 of T4S1, and
// plausible Z80 code at body offset 0 (sector offset 9). They apply
// to all dialects.

// ...bootSlot helper as above...
```

`rules_cosmetic.go`:

```go
// rules_cosmetic.go
package samfile

import "fmt"

// §14 Cosmetic / canonical-output rules (catalog docs/disk-validity-rules.md §14).
// Rules in this file warn when dir-entry bytes diverge from
// the conventions real ROM SAVE produces, without affecting
// runtime behaviour. They apply to all dialects.
```

- [ ] **Step 2: Update registry-growth gate**

In `rules_smoke_test.go`:

```go
func TestRegistryGrowth(t *testing.T) {
	if got := len(Rules()); got != 51 {
		t.Errorf("len(Rules()) = %d; want 51 (1 smoke + 19 phase-3 + 15 phase-4 + 12 phase-5 + 4 phase-6 rules)", got)
	}
}
```

- [ ] **Step 3: Build + test**

```
cd /Users/pmoore/git/samfile-verify-phase-6 && go build ./... && go test -run TestRegistryGrowth -v ./...
```
Expected: build silent; test FAILs with `len(Rules()) = 47; want 51`.

- [ ] **Step 4: Full suite**

```
go test ./...
```
Expected: only `TestRegistryGrowth` fails.

- [ ] **Step 5: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-6 && \
g add rules_boot.go rules_cosmetic.go rules_smoke_test.go && \
g commit -m "verify: phase 6 skeleton (boot + cosmetic files + bootSlot helper)

Two rule-file skeletons for the §11 boot-file rules and the §14
cosmetic tail. rules_boot.go also includes the bootSlot helper —
returns the used slot whose FirstSector is (4, 1) along with a
found flag, used by all three boot rules.

TestRegistryGrowth's count bumps from 47 to 51 (adds 3 boot
rules + 1 cosmetic). Deliberately failing until Tasks 2-4 land
the remaining rules."
```

---

## Task 2: §11 boot-file rules (3 rules)

**Files:**
- Modify: `rules_boot.go` — register and implement 3 rules.
- Modify: `rules_boot_test.go` — create with 6 tests.

### Rules

```go
// ----- BOOT-OWNER-AT-T4S1 -----
// For an image to be bootable on real SAM hardware, some directory
// entry's FirstSector must be (track 4, sector 1) so that the ROM
// BOOTEX (rom-disasm:20473-20598) reads the right sector at &8000.
// Fires on a single disk-wide finding when no used slot owns T4S1.
//
// Note: data-only / archive disks legitimately have no boot file; this
// rule's "fatal" severity flags non-bootability, not corruption.
// Phase 7's corpus-validation pass may demote to cosmetic if archive
// disks dominate the corpus.
func init() {
	Register(Rule{
		ID:          "BOOT-OWNER-AT-T4S1",
		Severity:    SeverityFatal,
		Description: "some used directory entry has FirstSector (4, 1) so the disk is bootable on SAM hardware",
		Citation:    "rom-disasm:20473-20598",
		Check:       checkBootOwnerAtT4S1,
	})
}

func checkBootOwnerAtT4S1(ctx *CheckContext) []Finding {
	if _, _, found := bootSlot(ctx.Journal); found {
		return nil
	}
	return []Finding{{
		RuleID:   "BOOT-OWNER-AT-T4S1",
		Severity: SeverityFatal,
		Location: DiskWideLocation(),
		Message:  "no used slot has FirstSector (track 4, sector 1); disk is not bootable on real SAM hardware",
		Citation: "rom-disasm:20473-20598",
	}}
}

// ----- BOOT-SIGNATURE-AT-256 -----
// For ROM BOOTEX to dispatch to the loaded sector, bytes 256-259 of
// T4S1 must spell "BOOT" — case-insensitively, with bit 7 ignored
// (the ROM compares (disk_byte XOR expected_byte) AND 0x5F per
// rom-disasm:20582-20598). Only applies when a boot owner exists;
// BOOT-OWNER-AT-T4S1 reports the no-owner case separately.
func init() {
	Register(Rule{
		ID:          "BOOT-SIGNATURE-AT-256",
		Severity:    SeverityFatal,
		Description: "T4S1 bytes 256-259 spell \"BOOT\" (case-insensitive, bit 7 ignored)",
		Citation:    "rom-disasm:20582-20598",
		Check:       checkBootSignatureAt256,
	})
}

func checkBootSignatureAt256(ctx *CheckContext) []Finding {
	slot, fe, found := bootSlot(ctx.Journal)
	if !found {
		return nil // BOOT-OWNER-AT-T4S1 reports the underlying issue
	}
	sd, err := ctx.Disk.SectorData(fe.FirstSector)
	if err != nil {
		return nil // §1 rules report the underlying sector problem
	}
	// ROM compares with `XOR expected; AND 0x5F` — 0x5F = 0b01011111
	// masks bits 5 (case) and 7 (BASIC-keyword high bit) before the
	// zero check. So we apply the same mask here.
	expected := [4]byte{'B', 'O', 'O', 'T'}
	for i := 0; i < 4; i++ {
		if (sd[256+i]^expected[i])&0x5F != 0 {
			return []Finding{{
				RuleID:   "BOOT-SIGNATURE-AT-256",
				Severity: SeverityFatal,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("T4S1 boot signature mismatch at byte %d: got 0x%02x, expected 0x%02x (masked with 0x5F)", 256+i, sd[256+i], expected[i]),
				Citation: "rom-disasm:20582-20598",
			}}
		}
	}
	return nil
}

// ----- BOOT-ENTRY-POINT-AT-9 -----
// After signature match, ROM does JP 8009H. The sector buffer is at
// 0x8000-0x81FF, so 0x8009 is sector-buffer offset 9 = body offset 0
// (after the 9-byte body header). The byte at body offset 0 must
// therefore be valid Z80 code. We can't enforce "valid Z80 opcode"
// precisely from one byte, but 0x00 (NOP — unlikely as the first
// boot-code byte by design) and 0xFF (unwritten / no-code marker)
// are useful negative signals. Cosmetic per the catalog's test sketch.
func init() {
	Register(Rule{
		ID:          "BOOT-ENTRY-POINT-AT-9",
		Severity:    SeverityCosmetic,
		Description: "T4S1 body byte 0 (sector offset 9) is not 0x00 or 0xFF — a heuristic plausibility check for Z80 boot code",
		Citation:    "rom-disasm:20598",
		Check:       checkBootEntryPointAt9,
	})
}

func checkBootEntryPointAt9(ctx *CheckContext) []Finding {
	slot, fe, found := bootSlot(ctx.Journal)
	if !found {
		return nil
	}
	sd, err := ctx.Disk.SectorData(fe.FirstSector)
	if err != nil {
		return nil
	}
	b := sd[9]
	if b == 0x00 || b == 0xFF {
		return []Finding{{
			RuleID:   "BOOT-ENTRY-POINT-AT-9",
			Severity: SeverityCosmetic,
			Location: SlotLocation(slot, fe.Name.String()),
			Message:  fmt.Sprintf("T4S1 body byte 0 = 0x%02x; expected a real Z80 opcode (0x00 = NOP and 0xFF = unwritten are implausible boot entries)", b),
			Citation: "rom-disasm:20598",
		}}
	}
	return nil
}
```

### Tests

Create `rules_boot_test.go`:

```go
package samfile

import "testing"

// buildBootableDisk builds a samfile-built disk where slot 0's first
// sector is (4, 1) — i.e. AddCodeFile's first allocation. Body bytes
// 256-259 are patched to "BOOT" and body byte 0 (sector offset 9) is
// patched to a real opcode (0xC3 = JP nn) so all three §11 rules pass
// on the positive case.
func buildBootableDisk(t *testing.T) (*DiskImage, *DiskJournal) {
	t.Helper()
	di := NewDiskImage()
	if err := di.AddCodeFile("samdos2", make([]byte, 400), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	first := di.DiskJournal()[0].FirstSector
	sd, err := di.SectorData(first)
	if err != nil {
		t.Fatalf("SectorData: %v", err)
	}
	// "BOOT" at sector offset 256-259.
	copy(sd[256:260], []byte{'B', 'O', 'O', 'T'})
	// Real opcode at sector offset 9 (body offset 0). 0xC3 = JP nn.
	sd[9] = 0xC3
	di.WriteSector(first, sd)
	return di, di.DiskJournal()
}

func TestBootOwnerAtT4S1Positive(t *testing.T) {
	di, _ := buildBootableDisk(t)
	findings := checkBootOwnerAtT4S1(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("disk with T4S1 owner: %d findings; want 0", len(findings))
	}
}

func TestBootOwnerAtT4S1Negative(t *testing.T) {
	// Empty disk → no used slot owns T4S1 → rule fires.
	di := NewDiskImage()
	findings := checkBootOwnerAtT4S1(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BOOT-OWNER-AT-T4S1" {
		t.Fatalf("got %d findings, first=%+v; want 1 BOOT-OWNER-AT-T4S1", len(findings), findings)
	}
	if findings[0].Severity != SeverityFatal {
		t.Errorf("Severity = %v; want fatal", findings[0].Severity)
	}
}

func TestBootSignatureAt256Positive(t *testing.T) {
	di, _ := buildBootableDisk(t)
	findings := checkBootSignatureAt256(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("disk with BOOT signature: %d findings; want 0", len(findings))
	}
}

func TestBootSignatureAt256Negative(t *testing.T) {
	di, _ := buildBootableDisk(t)
	first := di.DiskJournal()[0].FirstSector
	sd, _ := di.SectorData(first)
	// Corrupt one byte of the signature.
	sd[257] = 'X'
	di.WriteSector(first, sd)
	findings := checkBootSignatureAt256(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BOOT-SIGNATURE-AT-256" {
		t.Fatalf("got %d findings, first=%+v; want 1 BOOT-SIGNATURE-AT-256", len(findings), findings)
	}
}

func TestBootSignatureAt256CaseInsensitive(t *testing.T) {
	// Lowercase "boot" must also match (ROM AND 0x5F mask).
	di, _ := buildBootableDisk(t)
	first := di.DiskJournal()[0].FirstSector
	sd, _ := di.SectorData(first)
	copy(sd[256:260], []byte{'b', 'o', 'o', 't'})
	di.WriteSector(first, sd)
	findings := checkBootSignatureAt256(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("lowercase 'boot' (ROM case-insensitive): %d findings; want 0", len(findings))
	}
}

func TestBootEntryPointAt9Positive(t *testing.T) {
	di, _ := buildBootableDisk(t)
	findings := checkBootEntryPointAt9(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("body[0] = 0xC3 (JP nn): %d findings; want 0", len(findings))
	}
}

func TestBootEntryPointAt9Negative(t *testing.T) {
	di, _ := buildBootableDisk(t)
	first := di.DiskJournal()[0].FirstSector
	sd, _ := di.SectorData(first)
	sd[9] = 0xFF // unwritten marker — implausible entry
	di.WriteSector(first, sd)
	findings := checkBootEntryPointAt9(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BOOT-ENTRY-POINT-AT-9" {
		t.Fatalf("got %d findings, first=%+v; want 1 BOOT-ENTRY-POINT-AT-9", len(findings), findings)
	}
}
```

That's 7 tests (3 positive + 3 standard negative + 1 case-insensitive regression).

- [ ] **Step 1: Implement the 3 rules + 7 tests**

- [ ] **Step 2: Build + run**

```
cd /Users/pmoore/git/samfile-verify-phase-6 && go test -run 'TestBoot' -v ./...
```
Expected: 7 PASS.

Full suite:

```
go test ./...
```
Expected: all pass EXCEPT `TestRegistryGrowth` (now 50; want 51).

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-6 && \
g add rules_boot.go rules_boot_test.go && \
g commit -m "verify: §11 boot-file rules (3 rules)

  BOOT-OWNER-AT-T4S1     fatal      a used slot owns T4S1
  BOOT-SIGNATURE-AT-256  fatal      T4S1[256..259] spells 'BOOT'
                                    (case-insensitive, bit 7 ignored)
  BOOT-ENTRY-POINT-AT-9  cosmetic   T4S1 body byte 0 is not 0x00/0xFF

All three share a bootSlot(dj) helper that returns the used slot
whose FirstSector is (4, 1). BOOT-SIGNATURE and BOOT-ENTRY-POINT
silently skip when no boot owner exists — BOOT-OWNER reports the
missing-owner case separately.

BOOT-SIGNATURE applies ROM BOOTEX's mask (XOR expected; AND 0x5F)
so 'BOOT', 'boot', 'BOOt', and the BTWD literal 0x42 0x4F 0x4F 0xD4
all match. Regression test TestBootSignatureAt256CaseInsensitive
pins this."
```

---

## Task 3: §14 cosmetic rule (1 rule)

**Files:**
- Modify: `rules_cosmetic.go` — register and implement 1 rule.
- Modify: `rules_cosmetic_test.go` — create with 3 tests.

### Rule

```go
// ----- COSMETIC-RESERVEDA-FF -----
// Real ROM SAVE 0xFF-fills 14 bytes from dir offset 0xDC (HDCLP2 at
// rom-disasm:22076-22080), which covers MGTFlags + FileTypeInfo + the
// first two bytes of ReservedA. The catalog describes ReservedA (dir
// 0xE8-0xEB, 4 bytes) as fully 0xFF-filled by real SAVE. samfile's
// AddCodeFile leaves ReservedA at struct-zero (0x00). Both
// conventions are observed in the wild; the rule warns only when a
// byte is in NEITHER set — i.e. anything outside {0x00, 0xFF}.
//
// Same dual-acceptance pattern as Phase 4's CODE-FILETYPEINFO-EMPTY:
// real-ROM-SAVE byte == 0xFF, samfile byte == 0x00, both legitimate.
func init() {
	Register(Rule{
		ID:          "COSMETIC-RESERVEDA-FF",
		Severity:    SeverityCosmetic,
		Description: "ReservedA (dir 0xE8-0xEB) is uniformly 0x00 (samfile) or 0xFF (ROM SAMDOS-2)",
		Citation:    "rom-disasm:22076-22080",
		Check:       checkCosmeticReservedAFF,
	})
}

func checkCosmeticReservedAFF(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		for i, b := range fe.ReservedA {
			if b != 0x00 && b != 0xFF {
				findings = append(findings, Finding{
					RuleID:   "COSMETIC-RESERVEDA-FF",
					Severity: SeverityCosmetic,
					Location: SlotLocation(slot, fe.Name.String()),
					Message:  fmt.Sprintf("ReservedA[%d] (dir 0x%02x) = 0x%02x — neither samfile's 0x00 nor ROM SAVE's 0xFF", i, 0xE8+i, b),
					Citation: "rom-disasm:22076-22080",
				})
				return // one finding per slot
			}
		}
	})
	return findings
}
```

### Tests

Create `rules_cosmetic_test.go`:

```go
package samfile

import "testing"

func TestCosmeticReservedAFFPositive(t *testing.T) {
	// samfile's AddCodeFile leaves ReservedA at 0x00 × 4 — accepted.
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkCosmeticReservedAFF(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("samfile-built (ReservedA=0x00×4): %d findings; want 0", len(findings))
	}
}

func TestCosmeticReservedAFFAcceptsAllFF(t *testing.T) {
	// Real ROM SAMDOS-2 SAVE 0xFF-fills these bytes via HDCLP2.
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	for i := range dj[0].ReservedA {
		dj[0].ReservedA[i] = 0xFF
	}
	di.WriteFileEntry(dj, 0)
	findings := checkCosmeticReservedAFF(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("ReservedA=0xFF×4 (ROM SAVE convention): %d findings; want 0", len(findings))
	}
}

func TestCosmeticReservedAFFNegative(t *testing.T) {
	// A byte that's neither 0x00 nor 0xFF fires the rule.
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].ReservedA[2] = 0xAA
	di.WriteFileEntry(dj, 0)
	findings := checkCosmeticReservedAFF(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "COSMETIC-RESERVEDA-FF" {
		t.Fatalf("got %d findings, first=%+v; want 1 COSMETIC-RESERVEDA-FF", len(findings), findings)
	}
}
```

- [ ] **Step 1: Implement the rule + 3 tests**

- [ ] **Step 2: Full suite**

```
cd /Users/pmoore/git/samfile-verify-phase-6 && go test ./...
```
Expected: all green; `TestRegistryGrowth` now reports 51 and PASSES. Catalog implementation is closed.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-6 && \
g add rules_cosmetic.go rules_cosmetic_test.go && \
g commit -m "verify: §14 COSMETIC-RESERVEDA-FF (1 rule)

Closes Phase 6 at 4 rules. ReservedA (dir 0xE8-0xEB) is the last
4-byte block real ROM SAMDOS-2 SAVE 0xFF-fills via HDCLP2
(rom-disasm:22076-22080). samfile's AddCodeFile leaves it at
struct-zero. The rule accepts both conventions (same dual-
acceptance pattern as Phase 4's CODE-FILETYPEINFO-EMPTY); fires
only on bytes outside {0x00, 0xFF}.

TestRegistryGrowth now passes at 51 (1 smoke + 19 phase-3 + 15
phase-4 + 12 phase-5 + 4 phase-6 rules). The catalog of validity
rules is now fully implemented; Phase 7 is the empirical
corpus-validation pass."
```

---

## Task 4: Final verification + push + draft PR + monitor CI

- [ ] **Step 1: Full suite + vet**

```
cd /Users/pmoore/git/samfile-verify-phase-6 && go test ./... && go vet ./...
```
Expected: all green, vet silent.

- [ ] **Step 2: Build the CLI**

```
cd /Users/pmoore/git/samfile-verify-phase-6 && go build -o /tmp/samfile-phase6 ./cmd/samfile
```

- [ ] **Step 3: Run verify on the M0 boot disk**

```
[ -f /Users/pmoore/git/sam-aarch64/build/test.mgt ] && /tmp/samfile-phase6 verify -i /Users/pmoore/git/sam-aarch64/build/test.mgt 2>/dev/null | head -30 || echo "no M0 disk"
```

Expected: `detected dialect: samdos2`. The boot file is `samdos2` at T4S1 — BOOT-OWNER passes. The body bytes 256-259 should spell "BOOT" (the samdos2 binary has this signature near offset 256). BOOT-ENTRY-POINT-AT-9 should pass (samdos2's first body byte is real Z80 code, not 0x00 / 0xFF). COSMETIC-RESERVEDA-FF should pass on every slot (samfile-built CODE files leave ReservedA at 0x00, ROM SAVE'd files leave it at 0xFF; both accepted).

If anything fires on M0, investigate — Phase 6 rules should be clean on a samfile-built bootable disk.

- [ ] **Step 4: Run verify on the testdata corpus**

```
/tmp/samfile-phase6 verify -i /Users/pmoore/git/samfile-verify-phase-6/testdata/ETrackerv1.2.mgt 2>/dev/null | grep -E '^[A-Z]+ \(' && /tmp/samfile-phase6 verify -i /Users/pmoore/git/samfile-verify-phase-6/testdata/ETrackerv1.2.mgt 2>/dev/null | tail -3
```

Expected: finding count similar to Phase 5's 465, plus a small bump from the 4 new rules.

- [ ] **Step 5: Push**

```bash
cd /Users/pmoore/git/samfile-verify-phase-6 && g push -u origin feat/verify-phase-6-boot-cosmetic-tail
```

- [ ] **Step 6: Open the draft PR**

```bash
cd /Users/pmoore/git/samfile-verify-phase-6 && gh pr create --draft --base master \
  --title "verify: Phase 6 — boot-file & cosmetic-tail rules (4 rules; catalog complete)" \
  --body "$(cat <<'EOF'
Phase 6 of `samfile verify` (spec: `docs/specs/2026-05-11-verify-feature-design.md`, plan: `docs/plans/2026-05-12-verify-phase-6-boot-cosmetic-tail.md`). Implements the 4 remaining substantively-checkable rules: 3 from §11 (boot file) and 1 from §14 (cosmetic). After this lands the registry holds 51 rules total and the catalog is fully realised for v1; Phase 7 is the empirical corpus-validation pass.

## Rules added

**§11 boot file** (3): `BOOT-OWNER-AT-T4S1`, `BOOT-SIGNATURE-AT-256`, `BOOT-ENTRY-POINT-AT-9`

**§14 cosmetic tail** (1): `COSMETIC-RESERVEDA-FF`

Severity distribution: 2 fatal, 0 structural, 0 inconsistency, 2 cosmetic.

## Deliberately deferred

Ten catalog entries across §11-§15 are documentation or already covered by earlier phases:

- `BOOT-FILE-TYPE-IGNORED` — catalog says "not a validity check; a note for verify to not flag".
- `ATTR-HIDDEN-NOT-LISTED`, `ATTR-PROTECTED-NO-OVERWRITE`, `ATTR-ERASED-SUPPRESSES-ALL` — meta-rules / semantic notes; the erased-suppression behaviour is already baked into Phase 3's `forEachUsedSlot`.
- `DIALECT-MASTERDOS-MGTFLAGS` — drives Phase 2's `mgtFlagsDialect`.
- `DIALECT-MASTERDOS-GAP-2156` — drives Phase 5's `BASIC-VARS-GAP-INVARIANT`.
- `DIALECT-SAMDOS-1-TYPE-3` and `DIALECT-HOOK-128-DEAD-CODE` — documentation entries.
- `COSMETIC-RESERVEDB-FILL` — catalog explicitly says "not a rule; document only".
- `COSMETIC-STARTPAGE-DECORATIVE-BITS` — a comparison convention for byte-perfect diffing, not a runtime check.

§15 entries all map to rules already implemented in §1-§10.

## Architecture

Two new files: `rules_boot.go` (§11) and `rules_cosmetic.go` (§14). One private helper:

- `bootSlot(*DiskJournal) (slot int, fe *FileEntry, found bool)` finds the used slot whose `FirstSector` is `(4, 1)`. `BOOT-OWNER-AT-T4S1` produces a finding when `found == false`; `BOOT-SIGNATURE-AT-256` and `BOOT-ENTRY-POINT-AT-9` silently skip in that case.

`COSMETIC-RESERVEDA-FF` follows the dual-acceptance pattern Phase 4 established for `CODE-FILETYPEINFO-EMPTY`: accept both `0x00` (samfile's struct-zero) and `0xFF` (real ROM SAMDOS-2 SAVE's HDCLP2 0xFF-fill); warn on anything else.

## CLI smoke

- **M0 boot disk** (`../sam-aarch64/build/test.mgt`): [fill in].
- **`testdata/ETrackerv1.2.mgt`**: [fill in].

## Test plan

- [x] `go test ./...` — all green (10 new tests: 7 §11 + 3 §14)
- [x] `go vet ./...` — clean
- [x] CLI smoke against `testdata/ETrackerv1.2.mgt` produces a well-formed report
- [x] CLI smoke against the M0 boot disk reports `samdos2`
- [ ] GitHub Actions CI green

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 7: Monitor CI**

```bash
cd /Users/pmoore/git/samfile-verify-phase-6 && gh pr checks --watch
```

- [ ] **Step 8: Hand off**

Reply with the PR URL, CI status, and CLI smoke results.

---

## Self-review notes

**Spec coverage walk-through:**

| Spec requirement (Phase 6) | Where in plan |
|---|---|
| Boot-file rules (§11) | Task 2 — 3 rules |
| Cosmetic tail (§14) | Task 3 — 1 rule |
| "After this lands, the catalog is fully realised" | Confirmed; 51 rules registered, 10 deferred-as-documentation entries listed in PR body |

4 in-scope rules, 10 explicitly deferred. Spec covered.

**Type / signature consistency:**

- `bootSlot(dj *DiskJournal) (slot int, fe *FileEntry, found bool)` — used by 3 §11 rules.
- `forEachUsedSlot` (Phase 3) — used by `COSMETIC-RESERVEDA-FF`.
- `cleanSingleFileDisk` (Phase 3) — used by §14 tests.
- `buildBootableDisk` — new helper in `rules_boot_test.go`, used by all 3 boot-rule positive tests + 1 case-insensitive regression test.

All consistent.

**Rule severity sanity check (4 rules total):**

| Severity | Count | Rules |
|---|---|---|
| Fatal | 2 | BOOT-OWNER-AT-T4S1, BOOT-SIGNATURE-AT-256 |
| Cosmetic | 2 | BOOT-ENTRY-POINT-AT-9, COSMETIC-RESERVEDA-FF |

Total: 2 + 2 = 4 ✓. Registry final after Task 3 = 51.
