# Verify Phase 5 — FT_SAM_BASIC, Array, SCREEN, ZX-Snapshot Rules

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement 12 of the catalog's §7 (BASIC), §8 (array), §9 (SCREEN), §10 (ZX snapshot) rules. After this lands the registry holds 47 rules total (Phase-1 smoke + 19 Phase-3 + 15 Phase-4 + 12 Phase-5); the only remaining phase before corpus validation is Phase 6 (boot-file + cross-entry + dialect-specific + cosmetic tail).

**Architecture:** Four new files following the Phase 3/4 convention (one file per catalog section): `rules_ft_basic.go`, `rules_ft_array.go`, `rules_ft_screen.go`, `rules_ft_zxsnap.go`. A new private helper `bodyData(*DiskImage, *FileEntry) ([]byte, error)` reads each file's body (excluding the 9-byte header) once — needed for the BASIC content rules that walk the program. Each rule filters on its target file type at the top of its Check function (same pattern as §6's FT_CODE filter). `BASIC-VARS-GAP-INVARIANT` consults `ctx.Dialect` directly (samdos2 → 604; masterdos → 2156; unknown → accept both).

**Tech Stack:** Go 1.22+. Adds a dependency on the existing `sambasic` package for `Parse(body []byte) (*File, error)` — used by `BASIC-LINE-NUMBER-BE` and `BASIC-STARTLINE-WITHIN-PROG`.

**Context for the engineer:**

Read these first, in order:

1. `docs/specs/2026-05-11-verify-feature-design.md` §"Implementation order" Phase 5: "~13 rules. File-type-specific content checks."
2. `docs/disk-validity-rules.md` §7 (BASIC), §8 (array), §9 (SCREEN), §10 (ZX snapshot).
3. `samfile.go:80-132` — `FileEntry` struct, particularly `FileTypeInfo`, `SAMBASICStartLine`, `ExecutionAddressMod16K`.
4. `samfile.go:670-720` — `ProgramLength`, `NumericVariablesSize`, `GapSize`, `StringArrayVariablesSize`, `Start()`, `Length()`. These accessors decode the BASIC body layout from FileTypeInfo triplets.
5. `samfile.go:731-770` — `(*DiskImage).File(filename)` is the existing chain-walker that reads body bytes. Your `bodyData` helper is the same loop without the filename-lookup wrapper.
6. `sambasic/parse.go` and `sambasic/file.go` — `Parse(body []byte) (*sambasic.File, error)` returns `(Lines []*Line, StartLine uint16)`. `Line` has a `LineNumber` and `Bytes` accessor. Used by Phase 5 to walk the program for line-number checks.
7. `rules_ft_code.go` from Phase 4 — the canonical "filter on `fe.Type == X` then check invariant" pattern.

**Phase 5 scope: 12 rules.** All catalog entries in §7–§10 are implementable. No deferrals.

The 12 rules **in scope** (7 + 1 + 2 + 2 = 12):

**§7 FT_SAM_BASIC** (7):
- `BASIC-FILETYPEINFO-TRIPLETS` (structural)
- `BASIC-VARS-GAP-INVARIANT` (cosmetic, dialect-aware)
- `BASIC-PROG-END-SENTINEL` (structural)
- `BASIC-LINE-NUMBER-BE` (structural)
- `BASIC-STARTLINE-FF-DISABLES` (structural)
- `BASIC-STARTLINE-WITHIN-PROG` (cosmetic)
- `BASIC-MGTFLAGS-20` (inconsistency)

**§8 Array** (1):
- `ARRAY-FILETYPEINFO-TLBYTE-NAME` (structural)

**§9 SCREEN** (2):
- `SCREEN-MODE-AT-0xDD` (structural)
- `SCREEN-LENGTH-MATCHES-MODE` (structural)

**§10 ZX snapshot** (2):
- `ZXSNAP-LENGTH-49152` (structural)
- `ZXSNAP-LOAD-ADDR-16384` (structural)

After Task 5 the registry holds 47 rules total.

**Phase 5 standing rules** (same as Phase 3/4):

- Use `g` not plain `git` for commits.
- Every rule's `Citation` is a real `file:line`; copy verbatim from the plan.
- Test fabrication uses the inline pattern. Reuse `cleanSingleFileDisk` from `rules_disk_test.go`; introduce one new BASIC-fabrication helper.
- Each rule ships with positive + negative tests.
- Draft PR only; Task 6 handles push/PR/CI.
- All rules use `Dialects: nil` (apply to all dialects). The catalog's per-rule dialect tags are informational; rule LOGIC may consult `ctx.Dialect` (e.g. `BASIC-VARS-GAP-INVARIANT`) but the registry doesn't filter by dialect.

---

## File Structure

| Path | Action | Responsibility |
|---|---|---|
| `rules_ft_basic.go` | Create | §7 BASIC rules: 7 rules + the `bodyData` helper. |
| `rules_ft_basic_test.go` | Create | Positive + negative tests for §7 rules. Introduces a `buildBasicDisk` test helper. |
| `rules_ft_array.go` | Create | §8 array rule (1 rule). |
| `rules_ft_array_test.go` | Create | Positive + negative tests for §8. |
| `rules_ft_screen.go` | Create | §9 SCREEN rules (2 rules). |
| `rules_ft_screen_test.go` | Create | Positive + negative tests for §9. |
| `rules_ft_zxsnap.go` | Create | §10 ZX snapshot rules (2 rules). |
| `rules_ft_zxsnap_test.go` | Create | Positive + negative tests for §10. |
| `rules_smoke_test.go` | Modify | `TestRegistryGrowth` count update 35 → 47. |

---

## The body-data helper

Add at the top of `rules_ft_basic.go`:

```go
package samfile

import (
	"encoding/binary"
	"fmt"

	"github.com/petemoore/samfile/v3/sambasic"
)

// bodyData reads the file body (excluding the 9-byte header) by
// walking fe's sector chain. Mirrors the chain-walk loop in
// (*DiskImage).File but without the filename-lookup wrapper, so
// callers that already have a *FileEntry don't re-iterate the
// directory. Returns ("body bytes", nil) on success or
// (nil, err) when a SectorData call fails — rules treat the error
// as "no finding" because Phase 3's §1/§3 rules already report the
// underlying chain problem.
//
// The returned slice is fe.Length() bytes long; it does NOT include
// the body-header bytes 0..8, matching the convention of samfile.File's
// Body field.
func bodyData(di *DiskImage, fe *FileEntry) ([]byte, error) {
	fileLength := fe.Length()
	raw := make([]byte, fileLength+9)
	sd, err := di.SectorData(fe.FirstSector)
	if err != nil {
		return nil, err
	}
	fp := sd.FilePart()
	i := uint16(0)
	for {
		copy(raw[510*i:], fp.Data[:])
		i++
		if i == fe.Sectors {
			break
		}
		sd, err = di.SectorData(fp.NextSector)
		if err != nil {
			return nil, err
		}
		fp = sd.FilePart()
	}
	return raw[9:], nil
}
```

---

## Task 1: Skeleton + registry-growth gate update + bodyData helper

**Files:**
- Create: `rules_ft_basic.go`, `rules_ft_array.go`, `rules_ft_screen.go`, `rules_ft_zxsnap.go` (skeletons).
- Modify: `rules_smoke_test.go` — update `TestRegistryGrowth` count to 47.

- [ ] **Step 1: Create the four rule-file skeletons**

`rules_ft_basic.go` gets the package decl, section comment, imports, AND the `bodyData` helper above.

The other three files get only the package decl and a section comment (no rules yet). Example for `rules_ft_array.go`:

```go
// rules_ft_array.go
package samfile

import "fmt"

// §8 Array rules (catalog docs/disk-validity-rules.md §8).
// Rules in this file check FT_NUM_ARRAY (17) and FT_STR_ARRAY (18)
// invariants. They apply to all dialects.
```

`rules_ft_screen.go`:

```go
// rules_ft_screen.go
package samfile

import "fmt"

// §9 SCREEN rules (catalog docs/disk-validity-rules.md §9).
// Rules in this file check FT_SCREEN (20) invariants: mode byte
// and body-length-vs-mode geometry. They apply to all dialects.
```

`rules_ft_zxsnap.go`:

```go
// rules_ft_zxsnap.go
package samfile

import "fmt"

// §10 ZX snapshot rules (catalog docs/disk-validity-rules.md §10).
// Rules in this file check FT_ZX_SNAPSHOT (5) invariants: 48 KiB
// body length and 0x4000 load address. The catalog tags these as
// SAMDOS-2 specific (the constants live in SAMDOS source); we run
// them on all dialects because the ZX snapshot format is itself
// dialect-agnostic.
```

- [ ] **Step 2: Update registry-growth gate**

In `rules_smoke_test.go`, update `TestRegistryGrowth`:

```go
func TestRegistryGrowth(t *testing.T) {
	if got := len(Rules()); got != 47 {
		t.Errorf("len(Rules()) = %d; want 47 (1 smoke + 19 phase-3 + 15 phase-4 + 12 phase-5 rules)", got)
	}
}
```

- [ ] **Step 3: Build + test**

```
cd /Users/pmoore/git/samfile-verify-phase-5 && go build ./... && go test -run TestRegistryGrowth -v ./...
```
Expected: build silent; test FAILs with `len(Rules()) = 35; want 47`.

- [ ] **Step 4: Full suite**

```
go test ./...
```
Expected: only `TestRegistryGrowth` fails.

- [ ] **Step 5: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-5 && \
g add rules_ft_basic.go rules_ft_array.go rules_ft_screen.go rules_ft_zxsnap.go rules_smoke_test.go && \
g commit -m "verify: phase 5 skeleton (4 file-type rule files + bodyData helper)

Adds empty rules_ft_{basic,array,screen,zxsnap}.go skeletons for
catalog §7/§8/§9/§10 plus a private bodyData(*DiskImage, *FileEntry)
helper in rules_ft_basic.go. bodyData walks a file's sector chain
and returns the body bytes without the 9-byte header; used by
BASIC content rules that parse the tokenised program.

TestRegistryGrowth's count bumps from 35 to 47 (1 smoke + 19
phase-3 + 15 phase-4 + 12 phase-5 rules). Deliberately failing
until Tasks 2-5 register the remaining rules."
```

---

## Task 2: §7 BASIC rules (7 rules)

**Why this task exists:** §7 is the largest section in Phase 5 and exercises both the `bodyData` helper and the `sambasic.Parse` integration. Doing it in one commit keeps the test-fixture infrastructure (a `buildBasicDisk` helper plus tests) together.

**Files:**
- Modify: `rules_ft_basic.go` — register and implement 7 rules.
- Modify: `rules_ft_basic_test.go` — create with the test helper + 14 tests.

### Test fixture helper

The §7 rules need a real BASIC file fixture. samfile's `AddBasicFile(name, file *sambasic.File)` is the right entry point. A `sambasic.File` needs `Lines []*Line` and `StartLine uint16`. Put the helper at the top of `rules_ft_basic_test.go`:

```go
package samfile

import (
	"testing"

	"github.com/petemoore/samfile/v3/sambasic"
)

// buildBasicDisk returns a samfile-built disk containing one BASIC
// program with one line (10 REM "hi") and auto-RUN at line 10. The
// returned dj is the journal at construction time; callers can
// mutate slot 0 and call di.WriteFileEntry(dj, 0) to test
// negative cases.
//
// The defaults produce a SAMDOS-2-canonical disk: NumericVars=92
// bytes + Gap=512 bytes = SAVARS-NVARS=604 (sambasic/file.go:3-6).
func buildBasicDisk(t *testing.T) (*DiskImage, *DiskJournal) {
	t.Helper()
	bf := &sambasic.File{
		StartLine: 10,
		Lines: []sambasic.Line{
			{Number: 10, Tokens: []sambasic.Token{sambasic.REM, sambasic.String("hi")}},
		},
	}
	di := NewDiskImage()
	if err := di.AddBasicFile("DEMO", bf); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}
	return di, di.DiskJournal()
}

// buildBasicDiskNoAutoRun builds a BASIC disk where StartLine is
// 0xFFFF (no auto-RUN). Used by BASIC-STARTLINE-* rule tests.
func buildBasicDiskNoAutoRun(t *testing.T) (*DiskImage, *DiskJournal) {
	t.Helper()
	bf := &sambasic.File{
		StartLine: 0xFFFF,
		Lines: []sambasic.Line{
			{Number: 10, Tokens: []sambasic.Token{sambasic.REM, sambasic.String("hi")}},
		},
	}
	di := NewDiskImage()
	if err := di.AddBasicFile("DEMO", bf); err != nil {
		t.Fatalf("AddBasicFile (no auto-RUN): %v", err)
	}
	return di, di.DiskJournal()
}
```

API verified against `sambasic/file.go` and `sambasic/keywords.go` at plan-write time: `Lines []Line` (by value, not pointer), `Line.Number` (uint16), `Line.Tokens []Token`, and `sambasic.REM` is `SingleByteKeyword(0xB7)`.

### Rules

```go
// ----- BASIC-FILETYPEINFO-TRIPLETS -----
// For FT_SAM_BASIC, dir bytes 0xDD-0xE5 hold three 3-byte PAGEFORM
// lengths (cumulative offsets into the body): NVARS-PROG, NUMEND-PROG,
// SAVARS-PROG. The decoded values must be non-zero AND satisfy
// NVARS <= NUMEND <= SAVARS <= body Length.
func init() {
	Register(Rule{
		ID:          "BASIC-FILETYPEINFO-TRIPLETS",
		Severity:    SeverityStructural,
		Description: "FT_SAM_BASIC FileTypeInfo (dir 0xDD-0xE5) holds three non-zero, non-decreasing PAGEFORM cumulative offsets bounded by body length",
		Citation:    "rom-disasm:22163-22180",
		Check:       checkBasicFileTypeInfoTriplets,
	})
}

func checkBasicFileTypeInfoTriplets(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		nvars := fe.ProgramLength()                                       // decode(FileTypeInfo[0..2])
		numend := fe.ProgramLength() + fe.NumericVariablesSize()          // decode(FileTypeInfo[3..5])
		savars := fe.ProgramLength() + fe.NumericVariablesSize() + fe.GapSize() // decode(FileTypeInfo[6..8])
		length := fe.Length()
		if nvars == 0 || numend == 0 || savars == 0 {
			findings = append(findings, Finding{
				RuleID: "BASIC-FILETYPEINFO-TRIPLETS", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC file has zero offset in FileTypeInfo triplet (NVARS=%d NUMEND=%d SAVARS=%d)", nvars, numend, savars),
				Citation: "rom-disasm:22163-22180",
			})
			return
		}
		if !(nvars <= numend && numend <= savars && savars <= length) {
			findings = append(findings, Finding{
				RuleID: "BASIC-FILETYPEINFO-TRIPLETS", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC FileTypeInfo offsets out of order (NVARS=%d NUMEND=%d SAVARS=%d length=%d)", nvars, numend, savars, length),
				Citation: "rom-disasm:22163-22180",
			})
		}
	})
	return findings
}

// ----- BASIC-VARS-GAP-INVARIANT -----
// Empirically, SAMDOS-2 BASIC files have SAVARS-NVARS == 604, MasterDOS
// BASIC files have SAVARS-NVARS == 2156 (sam-basic-save-format.md, 161-
// disk scan). Cosmetic; depends on detected dialect — on Unknown, accept
// either value.
func init() {
	Register(Rule{
		ID:          "BASIC-VARS-GAP-INVARIANT",
		Severity:    SeverityCosmetic,
		Description: "FT_SAM_BASIC SAVARS-NVARS equals the dialect-canonical value (604 SAMDOS-2 / 2156 MasterDOS)",
		Citation:    "sam-basic-save-format.md",
		Check:       checkBasicVarsGapInvariant,
	})
}

func checkBasicVarsGapInvariant(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		gap := fe.NumericVariablesSize() + fe.GapSize() // SAVARS - NVARS
		var expected uint32
		switch ctx.Dialect {
		case DialectSAMDOS2:
			expected = 604
		case DialectMasterDOS:
			expected = 2156
		default:
			// Unknown — accept either canonical value, silently skip.
			if gap == 604 || gap == 2156 {
				return
			}
			expected = 604 // for the message; flag the rarer of the two
		}
		if gap != expected {
			findings = append(findings, Finding{
				RuleID: "BASIC-VARS-GAP-INVARIANT", Severity: SeverityCosmetic,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC SAVARS-NVARS = %d; expected %d for dialect %s", gap, expected, ctx.Dialect),
				Citation: "sam-basic-save-format.md",
			})
		}
	})
	return findings
}

// ----- BASIC-PROG-END-SENTINEL -----
// The tokenised program ends with a 0xFF sentinel byte. The byte at
// body[ProgramLength-1] is the sentinel (NVARS-PROG is the program-area
// end offset).
func init() {
	Register(Rule{
		ID:          "BASIC-PROG-END-SENTINEL",
		Severity:    SeverityStructural,
		Description: "FT_SAM_BASIC program-area ends with a 0xFF sentinel byte",
		Citation:    "sambasic/file.go:36-42",
		Check:       checkBasicProgEndSentinel,
	})
}

func checkBasicProgEndSentinel(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		body, err := bodyData(ctx.Disk, fe)
		if err != nil {
			return
		}
		progLen := fe.ProgramLength()
		if progLen == 0 || int(progLen) > len(body) {
			return // BASIC-FILETYPEINFO-TRIPLETS will catch this
		}
		if body[progLen-1] != 0xFF {
			findings = append(findings, Finding{
				RuleID: "BASIC-PROG-END-SENTINEL", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC program does not end with 0xFF sentinel; body[%d] = 0x%02x", progLen-1, body[progLen-1]),
				Citation: "sambasic/file.go:36-42",
			})
		}
	})
	return findings
}

// ----- BASIC-LINE-NUMBER-BE -----
// Walk the program with sambasic.Parse; any parse failure means the
// big-endian line-number / little-endian length / 0x0D-terminator
// invariant doesn't hold somewhere. Also check each line number is
// in 1..16383.
func init() {
	Register(Rule{
		ID:          "BASIC-LINE-NUMBER-BE",
		Severity:    SeverityStructural,
		Description: "FT_SAM_BASIC program parses cleanly and every line number is in 1..16383",
		Citation:    "sambasic/parse.go",
		Check:       checkBasicLineNumberBE,
	})
}

func checkBasicLineNumberBE(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		body, err := bodyData(ctx.Disk, fe)
		if err != nil {
			return
		}
		progLen := fe.ProgramLength()
		if progLen == 0 || int(progLen) > len(body) {
			return
		}
		prog := body[:progLen]
		bf, err := sambasic.Parse(prog)
		if err != nil {
			findings = append(findings, Finding{
				RuleID: "BASIC-LINE-NUMBER-BE", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC program parse failed: %v", err),
				Citation: "sambasic/parse.go",
			})
			return
		}
		for _, ln := range bf.Lines {
			if ln.Number < 1 || ln.Number > 16383 {
				findings = append(findings, Finding{
					RuleID: "BASIC-LINE-NUMBER-BE", Severity: SeverityStructural,
					Location: SlotLocation(slot, fe.Name.String()),
					Message:  fmt.Sprintf("BASIC line number %d out of range (1..16383)", ln.Number),
					Citation: "sambasic/parse.go",
				})
				return // one finding per slot
			}
		}
	})
	return findings
}

// ----- BASIC-STARTLINE-FF-DISABLES -----
// dir[0xF2] (= fe.ExecutionAddressDiv16K) is 0x00 (auto-RUN) or 0xFF
// (no auto-RUN); when 0x00, dir[0xF3..0xF4] (= fe.SAMBASICStartLine) is
// a valid line number (1..16383, not 0xFFFF).
func init() {
	Register(Rule{
		ID:          "BASIC-STARTLINE-FF-DISABLES",
		Severity:    SeverityStructural,
		Description: "FT_SAM_BASIC dir[0xF2] is 0x00 (auto-RUN) or 0xFF (no auto-RUN); when 0x00, the start-line is a valid line number",
		Citation:    "rom-disasm:22136-22141",
		Check:       checkBasicStartLineFFDisables,
	})
}

func checkBasicStartLineFFDisables(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		marker := fe.ExecutionAddressDiv16K
		if marker != 0x00 && marker != 0xFF {
			findings = append(findings, Finding{
				RuleID: "BASIC-STARTLINE-FF-DISABLES", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC auto-RUN marker dir[0xF2] = 0x%02x (expected 0x00 or 0xFF)", marker),
				Citation: "rom-disasm:22136-22141",
			})
			return
		}
		if marker == 0x00 {
			line := fe.SAMBASICStartLine
			if line == 0 || line == 0xFFFF || line > 16383 {
				findings = append(findings, Finding{
					RuleID: "BASIC-STARTLINE-FF-DISABLES", Severity: SeverityStructural,
					Location: SlotLocation(slot, fe.Name.String()),
					Message:  fmt.Sprintf("BASIC auto-RUN enabled (dir[0xF2]=0x00) but start-line %d is invalid (1..16383)", line),
					Citation: "rom-disasm:22136-22141",
				})
			}
		}
	})
	return findings
}

// ----- BASIC-STARTLINE-WITHIN-PROG -----
// When auto-RUN is enabled, the start-line should correspond to an
// actual line in the saved program. Cosmetic — auto-RUN of a missing
// line just errors with "Statement lost", it's not a corruption.
func init() {
	Register(Rule{
		ID:          "BASIC-STARTLINE-WITHIN-PROG",
		Severity:    SeverityCosmetic,
		Description: "FT_SAM_BASIC auto-RUN start-line exists in the saved program",
		Citation:    "rom-disasm:22136-22141",
		Check:       checkBasicStartLineWithinProg,
	})
}

func checkBasicStartLineWithinProg(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		if fe.ExecutionAddressDiv16K != 0x00 {
			return // auto-RUN disabled; nothing to check
		}
		body, err := bodyData(ctx.Disk, fe)
		if err != nil {
			return
		}
		progLen := fe.ProgramLength()
		if progLen == 0 || int(progLen) > len(body) {
			return
		}
		bf, err := sambasic.Parse(body[:progLen])
		if err != nil {
			return // BASIC-LINE-NUMBER-BE reports the parse failure
		}
		want := fe.SAMBASICStartLine
		for _, ln := range bf.Lines {
			if ln.Number == want {
				return
			}
		}
		findings = append(findings, Finding{
			RuleID: "BASIC-STARTLINE-WITHIN-PROG", Severity: SeverityCosmetic,
			Location: SlotLocation(slot, fe.Name.String()),
			Message:  fmt.Sprintf("BASIC auto-RUN line %d not present in the saved program", want),
			Citation: "rom-disasm:22136-22141",
		})
	})
	return findings
}

// ----- BASIC-MGTFLAGS-20 -----
// Real-world BASIC files have MGTFlags == 0x20 (empirical convention,
// 50%+ of canonical disks, required for M0 boot per
// test-mgt-byte-layout.md §slot 1). Inconsistency severity.
func init() {
	Register(Rule{
		ID:          "BASIC-MGTFLAGS-20",
		Severity:    SeverityInconsistency,
		Description: "FT_SAM_BASIC MGTFlags is 0x20 (empirical convention)",
		Citation:    "test-mgt-byte-layout.md",
		Check:       checkBasicMGTFlags20,
	})
}

func checkBasicMGTFlags20(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		if fe.MGTFlags != 0x20 {
			findings = append(findings, Finding{
				RuleID: "BASIC-MGTFLAGS-20", Severity: SeverityInconsistency,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC file MGTFlags = 0x%02x; expected 0x20 (empirical convention)", fe.MGTFlags),
				Citation: "test-mgt-byte-layout.md",
			})
		}
	})
	return findings
}

var _ = binary.LittleEndian // imported above; remove the import if no rule uses it
```

**Engineer note on imports**: `encoding/binary` is imported above in case a rule needs it for multi-byte decoding. If no rule in `rules_ft_basic.go` ends up using it (the current code doesn't), remove the import AND the `var _` line before committing. The linter will complain about unused imports otherwise. (The plan includes the import preemptively because the body-loader pattern often grows to need it; if not, drop it.)

### Tests

Create `rules_ft_basic_test.go`. For each rule, positive + negative tests:

```go
func TestBasicFileTypeInfoTripletsPositive(t *testing.T) {
	di, _ := buildBasicDisk(t)
	findings := checkBasicFileTypeInfoTriplets(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("clean BASIC disk: %d findings; want 0", len(findings))
	}
}

func TestBasicFileTypeInfoTripletsNegative(t *testing.T) {
	di, dj := buildBasicDisk(t)
	// Zero out FileTypeInfo so the triplets decode to all zero.
	for i := range dj[0].FileTypeInfo {
		dj[0].FileTypeInfo[i] = 0
	}
	di.WriteFileEntry(dj, 0)
	findings := checkBasicFileTypeInfoTriplets(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BASIC-FILETYPEINFO-TRIPLETS" {
		t.Fatalf("got %d findings, first=%+v; want 1 BASIC-FILETYPEINFO-TRIPLETS", len(findings), findings)
	}
}
```

Apply the same pattern for the remaining 6 §7 rules. Concrete negative mutations:

| Rule | Negative mutation |
|---|---|
| `BASIC-FILETYPEINFO-TRIPLETS` | Zero FileTypeInfo via journal patch (above). |
| `BASIC-VARS-GAP-INVARIANT` | Use `ctx.Dialect = DialectSAMDOS2` explicitly in the CheckContext (so expected is 604); patch FileTypeInfo to produce a GapSize that makes SAVARS-NVARS = 999 (neither 604 nor 2156). Can use raw pageFormLength manipulation. Simpler: skip the negative test for this dialect-aware rule and explicitly test the three branches (samdos2 + 604 → no finding; samdos2 + 999 → 1 finding; masterdos + 2156 → no finding; unknown + 999 → no finding). |
| `BASIC-PROG-END-SENTINEL` | Read body via `bodyData`, find `body[progLen-1]`, patch that sector byte to 0xAA via `mutateFirstSectorByte` (calculating the sector + offset). Tricky because the program-area sentinel may be in a later sector than the first. **Simpler approach**: build the BASIC file then use `WriteSector` to patch the appropriate sector. Or: scope the fixture to a single-sector body (small BASIC), so the sentinel lands in the first sector at a known body offset. The default `buildBasicDisk` produces a small enough body for this. |
| `BASIC-LINE-NUMBER-BE` | Patch the first 4 bytes of the program area (body[0..3] in `bodyData` terms) to have an obviously-bad line-length byte. Easiest path: build the BASIC disk, then `WriteSector` to corrupt body[0..3] via direct raw-byte write. |
| `BASIC-STARTLINE-FF-DISABLES` | Set `dj[0].ExecutionAddressDiv16K = 0x42` (neither 0x00 nor 0xFF). |
| `BASIC-STARTLINE-WITHIN-PROG` | Default `buildBasicDisk` has line 10 + StartLine=10 → no finding. For the negative: `dj[0].SAMBASICStartLine = 99` (line 99 isn't in the program). |
| `BASIC-MGTFLAGS-20` | `dj[0].MGTFlags = 0x80` (not 0x20). |

For `BASIC-VARS-GAP-INVARIANT`, write three explicit dialect-scoped tests:

```go
func TestBasicVarsGapInvariantSAMDOS2Clean(t *testing.T) {
	di, _ := buildBasicDisk(t)
	// AddBasicFile defaults: SAMDOS-2 gap is 604. ctx.Dialect = SAMDOS2 → no finding.
	findings := checkBasicVarsGapInvariant(&CheckContext{
		Disk: di, Journal: di.DiskJournal(), Dialect: DialectSAMDOS2,
	})
	if len(findings) != 0 {
		t.Errorf("SAMDOS-2 clean BASIC: %d findings; want 0", len(findings))
	}
}

func TestBasicVarsGapInvariantSAMDOS2BadGap(t *testing.T) {
	// Pass a non-canonical NumericVars length so the gap (NumericVars
	// + Gap) becomes 92+1 + 512 = 605 (≠ 604, ≠ 2156). Under SAMDOS-2
	// dialect, the rule fires.
	bf := &sambasic.File{
		StartLine: 10,
		Lines: []sambasic.Line{
			{Number: 10, Tokens: []sambasic.Token{sambasic.REM, sambasic.String("hi")}},
		},
		NumericVars: make([]byte, 93), // default+1 → gap = 605
	}
	di := NewDiskImage()
	if err := di.AddBasicFile("DEMO", bf); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}
	findings := checkBasicVarsGapInvariant(&CheckContext{
		Disk: di, Journal: di.DiskJournal(), Dialect: DialectSAMDOS2,
	})
	if len(findings) != 1 || findings[0].RuleID != "BASIC-VARS-GAP-INVARIANT" {
		t.Fatalf("got %d findings, first=%+v; want 1 BASIC-VARS-GAP-INVARIANT", len(findings), findings)
	}
}

func TestBasicVarsGapInvariantMasterDOSClean(t *testing.T) {
	// Pass Gap=2064 so NumericVars+Gap = 92 + 2064 = 2156 (MasterDOS
	// canonical). Under MasterDOS dialect, no finding.
	bf := &sambasic.File{
		StartLine: 10,
		Lines: []sambasic.Line{
			{Number: 10, Tokens: []sambasic.Token{sambasic.REM, sambasic.String("hi")}},
		},
		Gap: make([]byte, 2064), // 92 default NumericVars + 2064 Gap = 2156
	}
	di := NewDiskImage()
	if err := di.AddBasicFile("DEMO", bf); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}
	findings := checkBasicVarsGapInvariant(&CheckContext{
		Disk: di, Journal: di.DiskJournal(), Dialect: DialectMasterDOS,
	})
	if len(findings) != 0 {
		t.Errorf("MasterDOS clean (gap=2156): %d findings; want 0", len(findings))
	}
}

func TestBasicVarsGapInvariantUnknownDialect(t *testing.T) {
	di, _ := buildBasicDisk(t)
	// Default gap is 604 (SAMDOS-2 canonical). Under Unknown dialect
	// the rule accepts both 604 and 2156 silently.
	findings := checkBasicVarsGapInvariant(&CheckContext{
		Disk: di, Journal: di.DiskJournal(), Dialect: DialectUnknown,
	})
	if len(findings) != 0 {
		t.Errorf("Unknown dialect + canonical gap: %d findings; want 0", len(findings))
	}
}
```

That's 4 tests for the dialect-aware rule instead of the usual 2.

- [ ] **Step 1: Implement the 7 rules + ~16 tests**

- [ ] **Step 2: Build + test**

```
cd /Users/pmoore/git/samfile-verify-phase-5 && go build ./... && go test -run 'TestBasic' -v ./...
```
Expected: all §7 tests PASS.

Then the full suite:

```
go test ./...
```
Expected: all tests pass EXCEPT `TestRegistryGrowth` (now 42 rules; want 47).

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-5 && \
g add rules_ft_basic.go rules_ft_basic_test.go && \
g commit -m "verify: §7 FT_SAM_BASIC rules (7 rules)

Adds seven rules covering BASIC file content and metadata:

  BASIC-FILETYPEINFO-TRIPLETS  structural
  BASIC-VARS-GAP-INVARIANT     cosmetic (dialect-aware)
  BASIC-PROG-END-SENTINEL      structural
  BASIC-LINE-NUMBER-BE         structural
  BASIC-STARTLINE-FF-DISABLES  structural
  BASIC-STARTLINE-WITHIN-PROG  cosmetic
  BASIC-MGTFLAGS-20            inconsistency

BASIC-LINE-NUMBER-BE and BASIC-STARTLINE-WITHIN-PROG use the
sambasic.Parse helper to walk the tokenised program. The dialect-
aware BASIC-VARS-GAP-INVARIANT consults ctx.Dialect: SAMDOS-2
expects gap=604, MasterDOS expects gap=2156, Unknown silently
accepts either canonical value. Test fixture buildBasicDisk uses
samfile's AddBasicFile + sambasic.File for a 1-line program."
```

---

## Task 3: §8 Array rule (1 rule)

**Files:**
- Modify: `rules_ft_array.go` — register and implement 1 rule.
- Modify: `rules_ft_array_test.go` — create with 2 tests.

```go
// ----- ARRAY-FILETYPEINFO-TLBYTE-NAME -----
// For FT_NUM_ARRAY (17) and FT_STR_ARRAY (18), dir bytes 0xDD-0xE7
// hold the array's TLBYTE (type/length byte) followed by its 10-byte
// name. The rule warns when all 11 bytes are zero — that indicates a
// writer didn't populate the array metadata at SAVE time.
func init() {
	Register(Rule{
		ID:          "ARRAY-FILETYPEINFO-TLBYTE-NAME",
		Severity:    SeverityStructural,
		Description: "FT_NUM_ARRAY/FT_STR_ARRAY FileTypeInfo (dir 0xDD-0xE7) is not all zero",
		Citation:    "rom-disasm:22354-22357",
		Check:       checkArrayFileTypeInfoTLBYTEName,
	})
}

func checkArrayFileTypeInfoTLBYTEName(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_NUM_ARRAY && fe.Type != FT_STR_ARRAY {
			return
		}
		allZero := true
		for _, b := range fe.FileTypeInfo {
			if b != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			findings = append(findings, Finding{
				RuleID: "ARRAY-FILETYPEINFO-TLBYTE-NAME", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  "array file FileTypeInfo (dir 0xDD-0xE7) is all zero; TLBYTE + name not populated",
				Citation: "rom-disasm:22354-22357",
			})
		}
	})
	return findings
}
```

Tests (build a CODE disk, morph to array type, check):

```go
package samfile

import "testing"

func TestArrayFileTypeInfoTLBYTENamePositive(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Type = FT_NUM_ARRAY
	dj[0].FileTypeInfo[0] = 0x42 // TLBYTE
	copy(dj[0].FileTypeInfo[1:], []byte("ARR       "))
	di.WriteFileEntry(dj, 0)
	findings := checkArrayFileTypeInfoTLBYTEName(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("array file with populated FileTypeInfo: %d findings; want 0", len(findings))
	}
}

func TestArrayFileTypeInfoTLBYTENameNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Type = FT_STR_ARRAY
	// FileTypeInfo is zero by default (AddCodeFile leaves it that way).
	di.WriteFileEntry(dj, 0)
	findings := checkArrayFileTypeInfoTLBYTEName(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "ARRAY-FILETYPEINFO-TLBYTE-NAME" {
		t.Fatalf("got %d findings, first=%+v; want 1 ARRAY-FILETYPEINFO-TLBYTE-NAME", len(findings), findings)
	}
}
```

- [ ] **Step 1: Implement the rule + 2 tests**

- [ ] **Step 2: Build + test**

```
cd /Users/pmoore/git/samfile-verify-phase-5 && go test -run 'TestArray' -v ./...
```
Expected: 2 PASS.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-5 && \
g add rules_ft_array.go rules_ft_array_test.go && \
g commit -m "verify: §8 array rule (1 rule)

ARRAY-FILETYPEINFO-TLBYTE-NAME (structural): for FT_NUM_ARRAY and
FT_STR_ARRAY, dir bytes 0xDD-0xE7 hold the array's TLBYTE and
10-byte name. Warn when all 11 bytes are zero (writer didn't
populate the array metadata at SAVE time)."
```

---

## Task 4: §9 SCREEN rules (2 rules)

**Files:**
- Modify: `rules_ft_screen.go` — register and implement 2 rules.
- Modify: `rules_ft_screen_test.go` — create with 4 tests.

```go
// ----- SCREEN-MODE-AT-0xDD -----
// For FT_SCREEN, dir byte 0xDD (= FileTypeInfo[0]) is the screen mode
// (1-4 on SAM).
func init() {
	Register(Rule{
		ID:          "SCREEN-MODE-AT-0xDD",
		Severity:    SeverityStructural,
		Description: "FT_SCREEN dir[0xDD] (mode byte) is in 1..4",
		Citation:    "rom-disasm:22259",
		Check:       checkScreenModeAt0xDD,
	})
}

func checkScreenModeAt0xDD(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SCREEN {
			return
		}
		mode := fe.FileTypeInfo[0]
		if mode < 1 || mode > 4 {
			findings = append(findings, Finding{
				RuleID: "SCREEN-MODE-AT-0xDD", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("SCREEN mode byte = %d (expected 1..4)", mode),
				Citation: "rom-disasm:22259",
			})
		}
	})
	return findings
}

// ----- SCREEN-LENGTH-MATCHES-MODE -----
// For FT_SCREEN, body Length() matches the documented screen size for
// the given mode: modes 1 and 2 → 6912 bytes, modes 3 and 4 → 24576
// bytes. (Skipped when mode is out-of-range; SCREEN-MODE-AT-0xDD
// catches that.)
func init() {
	Register(Rule{
		ID:          "SCREEN-LENGTH-MATCHES-MODE",
		Severity:    SeverityStructural,
		Description: "FT_SCREEN body length matches the documented size for its mode (1-2: 6912; 3-4: 24576)",
		Citation:    "sam-coupe_tech-man_v3-0.txt",
		Check:       checkScreenLengthMatchesMode,
	})
}

func checkScreenLengthMatchesMode(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SCREEN {
			return
		}
		mode := fe.FileTypeInfo[0]
		var expected uint32
		switch mode {
		case 1, 2:
			expected = 6912
		case 3, 4:
			expected = 24576
		default:
			return // SCREEN-MODE-AT-0xDD reports the bad mode
		}
		if fe.Length() != expected {
			findings = append(findings, Finding{
				RuleID: "SCREEN-LENGTH-MATCHES-MODE", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("SCREEN mode %d expects body length %d; got %d", mode, expected, fe.Length()),
				Citation: "sam-coupe_tech-man_v3-0.txt",
			})
		}
	})
	return findings
}
```

Tests:

```go
package samfile

import "testing"

func TestScreenModeAt0xDDPositive(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Type = FT_SCREEN
	dj[0].FileTypeInfo[0] = 2
	di.WriteFileEntry(dj, 0)
	findings := checkScreenModeAt0xDD(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("mode=2: %d findings; want 0", len(findings))
	}
}

func TestScreenModeAt0xDDNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Type = FT_SCREEN
	dj[0].FileTypeInfo[0] = 9 // out of range
	di.WriteFileEntry(dj, 0)
	findings := checkScreenModeAt0xDD(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "SCREEN-MODE-AT-0xDD" {
		t.Fatalf("got %d findings, first=%+v; want 1 SCREEN-MODE-AT-0xDD", len(findings), findings)
	}
}

func TestScreenLengthMatchesModePositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 6912-9) // body = 6912 after the 9-byte header allowance? No: Length() returns body length post-header. Use AddCodeFile with 6912 bytes for a clean fixture.
	di2 := NewDiskImage()
	if err := di2.AddCodeFile("SCREEN1", make([]byte, 6912), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di2.DiskJournal()
	dj[0].Type = FT_SCREEN
	dj[0].FileTypeInfo[0] = 1
	di2.WriteFileEntry(dj, 0)
	findings := checkScreenLengthMatchesMode(&CheckContext{Disk: di2, Journal: di2.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("mode 1 + 6912 bytes: %d findings; want 0", len(findings))
	}
	_ = di
}

func TestScreenLengthMatchesModeNegative(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("TEST", make([]byte, 100), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].Type = FT_SCREEN
	dj[0].FileTypeInfo[0] = 1 // expects 6912 bytes; body has 100
	di.WriteFileEntry(dj, 0)
	findings := checkScreenLengthMatchesMode(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "SCREEN-LENGTH-MATCHES-MODE" {
		t.Fatalf("got %d findings, first=%+v; want 1 SCREEN-LENGTH-MATCHES-MODE", len(findings), findings)
	}
}
```

- [ ] **Step 1: Implement the 2 rules + 4 tests**

- [ ] **Step 2: Build + test**

```
cd /Users/pmoore/git/samfile-verify-phase-5 && go test -run 'TestScreen' -v ./...
```
Expected: 4 PASS.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-5 && \
g add rules_ft_screen.go rules_ft_screen_test.go && \
g commit -m "verify: §9 SCREEN rules (2 rules)

SCREEN-MODE-AT-0xDD       structural — mode byte in 1..4
SCREEN-LENGTH-MATCHES-MODE structural — body length matches mode

Mode 1 and 2 use 6912 bytes (the SAM low-res framebuffer);
modes 3 and 4 use 24576 (the high-res framebuffer)."
```

---

## Task 5: §10 ZX snapshot rules (2 rules)

**Files:**
- Modify: `rules_ft_zxsnap.go` — register and implement 2 rules.
- Modify: `rules_ft_zxsnap_test.go` — create with 4 tests.

```go
// ----- ZXSNAP-LENGTH-49152 -----
// FT_ZX_SNAPSHOT has a 49,152-byte body (48 KiB ZX RAM).
func init() {
	Register(Rule{
		ID:          "ZXSNAP-LENGTH-49152",
		Severity:    SeverityStructural,
		Description: "FT_ZX_SNAPSHOT body is exactly 49152 bytes (48 KiB ZX RAM)",
		Citation:    "samdos/src/d.s:660-661",
		Check:       checkZXSnapLength49152,
	})
}

func checkZXSnapLength49152(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_ZX_SNAPSHOT {
			return
		}
		if fe.Length() != 49152 {
			findings = append(findings, Finding{
				RuleID: "ZXSNAP-LENGTH-49152", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("ZX snapshot body length = %d; expected 49152", fe.Length()),
				Citation: "samdos/src/d.s:660-661",
			})
		}
	})
	return findings
}

// ----- ZXSNAP-LOAD-ADDR-16384 -----
// FT_ZX_SNAPSHOT load address is 0x4000 (ZX RAM base).
func init() {
	Register(Rule{
		ID:          "ZXSNAP-LOAD-ADDR-16384",
		Severity:    SeverityStructural,
		Description: "FT_ZX_SNAPSHOT decoded start address is 0x4000 (16384, ZX RAM base)",
		Citation:    "samdos/src/d.s:660-663",
		Check:       checkZXSnapLoadAddr16384,
	})
}

func checkZXSnapLoadAddr16384(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_ZX_SNAPSHOT {
			return
		}
		if fe.StartAddress() != 16384 {
			findings = append(findings, Finding{
				RuleID: "ZXSNAP-LOAD-ADDR-16384", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("ZX snapshot start address = 0x%05x; expected 0x04000 (16384)", fe.StartAddress()),
				Citation: "samdos/src/d.s:660-663",
			})
		}
	})
	return findings
}
```

Tests:

```go
package samfile

import "testing"

// buildZXSnapDisk returns a samfile-built disk where slot 0 is a
// 49152-byte file morphed into FT_ZX_SNAPSHOT with start address
// 0x4000. AddCodeFile load 0x4000 sets fe.StartAddressPage so that
// StartAddress() decodes to 0x4000.
func buildZXSnapDisk(t *testing.T) (*DiskImage, *DiskJournal) {
	t.Helper()
	di := NewDiskImage()
	if err := di.AddCodeFile("ZXSNAP", make([]byte, 49152), 0x4000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].Type = FT_ZX_SNAPSHOT
	di.WriteFileEntry(dj, 0)
	return di, di.DiskJournal()
}

func TestZXSnapLength49152Positive(t *testing.T) {
	di, _ := buildZXSnapDisk(t)
	findings := checkZXSnapLength49152(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("49152-byte ZX snapshot: %d findings; want 0", len(findings))
	}
}

func TestZXSnapLength49152Negative(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("ZXSNAP", make([]byte, 100), 0x4000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].Type = FT_ZX_SNAPSHOT
	di.WriteFileEntry(dj, 0)
	findings := checkZXSnapLength49152(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "ZXSNAP-LENGTH-49152" {
		t.Fatalf("got %d findings, first=%+v; want 1 ZXSNAP-LENGTH-49152", len(findings), findings)
	}
}

func TestZXSnapLoadAddr16384Positive(t *testing.T) {
	di, _ := buildZXSnapDisk(t)
	findings := checkZXSnapLoadAddr16384(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("ZX snapshot at 0x4000: %d findings; want 0", len(findings))
	}
}

func TestZXSnapLoadAddr16384Negative(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("ZXSNAP", make([]byte, 49152), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	dj := di.DiskJournal()
	dj[0].Type = FT_ZX_SNAPSHOT
	di.WriteFileEntry(dj, 0)
	findings := checkZXSnapLoadAddr16384(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "ZXSNAP-LOAD-ADDR-16384" {
		t.Fatalf("got %d findings, first=%+v; want 1 ZXSNAP-LOAD-ADDR-16384", len(findings), findings)
	}
}
```

- [ ] **Step 1: Implement the 2 rules + 4 tests**

- [ ] **Step 2: Full suite**

```
cd /Users/pmoore/git/samfile-verify-phase-5 && go test ./...
```
Expected: all green; `TestRegistryGrowth` now reports 47 and PASSES.

- [ ] **Step 3: Commit**

```bash
cd /Users/pmoore/git/samfile-verify-phase-5 && \
g add rules_ft_zxsnap.go rules_ft_zxsnap_test.go && \
g commit -m "verify: §10 ZX snapshot rules (2 rules)

ZXSNAP-LENGTH-49152      structural — body is exactly 48 KiB
ZXSNAP-LOAD-ADDR-16384   structural — start address is 0x4000

Closes Phase 5 at 12 rules. TestRegistryGrowth now passes
(47 rules registered: 1 smoke + 19 phase-3 + 15 phase-4 +
12 phase-5)."
```

---

## Task 6: Final verification + push + draft PR + monitor CI

- [ ] **Step 1: Full suite + vet**

```
cd /Users/pmoore/git/samfile-verify-phase-5 && go test ./... && go vet ./...
```

- [ ] **Step 2: Build the CLI**

```
cd /Users/pmoore/git/samfile-verify-phase-5 && go build -o /tmp/samfile-phase5 ./cmd/samfile
```

- [ ] **Step 3: Run verify on the M0 boot disk**

```
[ -f /Users/pmoore/git/sam-aarch64/build/test.mgt ] && /tmp/samfile-phase5 verify -i /Users/pmoore/git/sam-aarch64/build/test.mgt 2>/dev/null | head -30 || echo "no M0 disk"
```

Expected: `detected dialect: samdos2`. The M0 disk has a BASIC file (slot 1 `auto`) — Phase 5 rules WILL inspect it. Acceptable outcomes:

- `BASIC-MGTFLAGS-20`: must not fire (M0's BASIC file has MGTFlags=0x20 per the build-disk output).
- `BASIC-FILETYPEINFO-TRIPLETS`: must not fire (samfile populates these correctly for AddBasicFile).
- `BASIC-VARS-GAP-INVARIANT`: depends on the gap size in samfile's writer; if not 604 or 2156, this cosmetic rule fires. Investigate.
- `BASIC-PROG-END-SENTINEL`: must not fire (sambasic always appends 0xFF to the program).
- `BASIC-LINE-NUMBER-BE`: must not fire (M0's BASIC parses cleanly).
- `BASIC-STARTLINE-FF-DISABLES`: must not fire (M0's BASIC has dir[0xF2]=0x00 and a valid line 10).
- `BASIC-STARTLINE-WITHIN-PROG`: must not fire (line 10 exists).

If anything fires on M0 that's not in the "expected to fire" list above, investigate before pushing. Phase 5 should produce zero findings on the M0 boot disk (the disk is a samfile-built BASIC + CODE disk with canonical content).

- [ ] **Step 4: Run verify on the testdata corpus**

```
/tmp/samfile-phase5 verify -i /Users/pmoore/git/samfile-verify-phase-5/testdata/ETrackerv1.2.mgt 2>/dev/null | grep -E '^[A-Z]+ \(' && /tmp/samfile-phase5 verify -i /Users/pmoore/git/samfile-verify-phase-5/testdata/ETrackerv1.2.mgt 2>/dev/null | tail -3
```

Expected: findings count somewhat higher than Phase 4's 492 (the new §7-§10 rules add coverage). Check the distribution by severity makes sense.

- [ ] **Step 5: Push**

```bash
cd /Users/pmoore/git/samfile-verify-phase-5 && g push -u origin feat/verify-phase-5-basic-array-screen-zxsnap
```

- [ ] **Step 6: Open the draft PR**

```bash
cd /Users/pmoore/git/samfile-verify-phase-5 && gh pr create --draft --base master \
  --title "verify: Phase 5 — FT_SAM_BASIC, array, SCREEN & ZX-snapshot rules (12 rules)" \
  --body "$(cat <<'EOF'
Phase 5 of `samfile verify` (spec: `docs/specs/2026-05-11-verify-feature-design.md`, plan: `docs/plans/2026-05-12-verify-phase-5-basic-array-screen-zxsnap.md`). Implements 12 of the catalog's §7 (BASIC), §8 (array), §9 (SCREEN), §10 (ZX snapshot) rules. After this lands the registry holds 47 rules total; Phase 6 (boot-file + cross-entry + dialect-specific + cosmetic tail) is the last implementation phase before Phase 7's corpus-validation pass.

## Rules added

**§7 BASIC** (7): `BASIC-FILETYPEINFO-TRIPLETS`, `BASIC-VARS-GAP-INVARIANT`, `BASIC-PROG-END-SENTINEL`, `BASIC-LINE-NUMBER-BE`, `BASIC-STARTLINE-FF-DISABLES`, `BASIC-STARTLINE-WITHIN-PROG`, `BASIC-MGTFLAGS-20`

**§8 array** (1): `ARRAY-FILETYPEINFO-TLBYTE-NAME`

**§9 SCREEN** (2): `SCREEN-MODE-AT-0xDD`, `SCREEN-LENGTH-MATCHES-MODE`

**§10 ZX snapshot** (2): `ZXSNAP-LENGTH-49152`, `ZXSNAP-LOAD-ADDR-16384`

Severity distribution: 0 fatal, 9 structural, 1 inconsistency, 2 cosmetic.

## Architecture

One file per catalog section: `rules_ft_basic.go` (§7), `rules_ft_array.go` (§8), `rules_ft_screen.go` (§9), `rules_ft_zxsnap.go` (§10). One new private helper:

- `bodyData(*DiskImage, *FileEntry) ([]byte, error)` reads each used file's body (excluding the 9-byte header) by walking its sector chain. Used by `BASIC-PROG-END-SENTINEL`, `BASIC-LINE-NUMBER-BE`, and `BASIC-STARTLINE-WITHIN-PROG` for BASIC body parsing via `sambasic.Parse`.

`BASIC-VARS-GAP-INVARIANT` consults `ctx.Dialect` (SAMDOS-2 expects gap=604, MasterDOS expects gap=2156, Unknown silently accepts either canonical value). All other rules are dialect-agnostic.

## CLI smoke

- **M0 boot disk**: `detected dialect: samdos2`, [findings list — fill in].
- **`testdata/ETrackerv1.2.mgt`**: [count] findings.

## Test plan

- [x] `go test ./...` — all green
- [x] `go vet ./...` — clean
- [x] CLI smoke against `testdata/ETrackerv1.2.mgt` produces a structurally well-formed report
- [x] CLI smoke against the M0 boot disk reports `samdos2` and finds no real bugs (any findings are explainable; see body)
- [ ] GitHub Actions CI green

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 7: Monitor CI**

```bash
cd /Users/pmoore/git/samfile-verify-phase-5 && gh pr checks --watch
```

- [ ] **Step 8: Hand off**

Reply with the PR URL, CI status, M0 disk findings (any bugs surfaced?), and ETracker finding distribution.

---

## Self-review notes

**Spec coverage walk-through:**

| Spec requirement (Phase 5) | Where in plan |
|---|---|
| FT_SAM_BASIC rules (§7) | Task 2 — 7 rules |
| Array rules (§8) | Task 3 — 1 rule |
| SCREEN rules (§9) | Task 4 — 2 rules |
| ZX snapshot rules (§10) | Task 5 — 2 rules |

12 rules in scope, no deferred entries. Spec covered.

**Placeholder scan:** the `BASIC-VARS-GAP-INVARIANT` negative test acknowledges that the FileTypeInfo patch math is engineer-fills-in. The rest is concrete. This is the only "fill in the arithmetic" in the plan; the engineer is responsible for getting a non-canonical gap value (e.g. 605 by adding 1 byte to NumericVariablesSize).

**Type / signature consistency:**

- `bodyData(di *DiskImage, fe *FileEntry) ([]byte, error)` — used in §7 rules.
- `forEachUsedSlot` from Phase 3 — used everywhere.
- `cleanSingleFileDisk` from Phase 3 — used by §8/§9 negative tests.
- `buildBasicDisk`/`buildBasicDiskNoAutoRun` — new helpers introduced in Task 2.
- `buildZXSnapDisk` — new helper introduced in Task 5.
- `ctx.Dialect` is consulted by `BASIC-VARS-GAP-INVARIANT`; all other rules are dialect-agnostic.
- `sambasic.Parse(body []byte) (*sambasic.File, error)` — used by BASIC-LINE-NUMBER-BE and BASIC-STARTLINE-WITHIN-PROG.

All consistent.

**Rule severity sanity check (12 rules total):**

| Severity | Count | Rules |
|---|---|---|
| Fatal | 0 | (none — Phase 5 has no fatal rules) |
| Structural | 9 | BASIC-FILETYPEINFO-TRIPLETS, BASIC-PROG-END-SENTINEL, BASIC-LINE-NUMBER-BE, BASIC-STARTLINE-FF-DISABLES, ARRAY-FILETYPEINFO-TLBYTE-NAME, SCREEN-MODE-AT-0xDD, SCREEN-LENGTH-MATCHES-MODE, ZXSNAP-LENGTH-49152, ZXSNAP-LOAD-ADDR-16384 |
| Inconsistency | 1 | BASIC-MGTFLAGS-20 |
| Cosmetic | 2 | BASIC-VARS-GAP-INVARIANT, BASIC-STARTLINE-WITHIN-PROG |

Total: 0 + 9 + 1 + 2 = 12 ✓. Registry final count after Task 5 = 47 (1 smoke + 19 phase-3 + 15 phase-4 + 12 phase-5).
