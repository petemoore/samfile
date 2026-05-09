# `samfile verify` — design

Status: approved 2026-05-11. Author: Pete (brainstormed via Claude).
Companion docs:

- `docs/disk-validity-rules.md` — the 71-rule validity catalog this feature
  consumes. The catalog enumerates each rule's ID, severity, dialect,
  source authority, and citation; this spec describes the machinery that
  runs them and surfaces findings.
- The `Disk integrity checks` wishlist in the
  [samfile v2.1.0 release notes](https://github.com/petemoore/samfile/releases/tag/v2.1.0)
  — all ten items from that list are covered by rules in the catalog.

## Why

samfile already reads, writes, and inspects MGT disk images at the
byte / sector / file level. It does not currently tell its caller
whether a given image is *valid*: whether the directory entries
agree with the body headers they mirror, whether sector chains
terminate cleanly, whether two files claim the same sector, whether
the boot sector contains a runnable DOS. `verify` adds that
capability.

Three primary users benefit:

1. **Inspectors / archivists** — feed in an MGT, get a readable
   report of what's structurally wrong or unusual. Primary v1
   audience.
2. **Builders** — `samfile add … && samfile verify` confirms the
   tool produced something correct.
3. **Forensic / recovery tooling** (longer-term) — consumes
   `Verify()` programmatically to identify what *should* be on a
   damaged disk so it can be reconstructed. The same rule set
   that powers detection powers eventual recovery.

## Decisions

These five choices were made during brainstorming; everything else
in this spec follows from them:

1. **Rule scope: all 71 rules in v1**, severity-gated by default.
   Every rule from `docs/disk-validity-rules.md` is implemented;
   the default UI only surfaces `fatal` + `structural` findings.
   `--all` (or `--severity=inconsistency`) brings in lower
   severities. This means corpus validation can classify every
   rule from day one rather than carrying a partial implementation
   through multiple releases.

2. **Primary audience: inspectors / archivists**, default verbose
   human-readable. Output is grouped, prose-style, easy to read in
   a terminal. `--json` emits machine-readable findings. `--quiet`
   suppresses output for shell pipelines. Exit code is `0` unless
   `fatal` findings exist, in which case `1` — `verify` is
   informational, not a CI gate by default.

3. **Dialect: auto-detect** from disk content. The catalog tags
   rules with dialect applicability (SAMDOS-2 / MasterDOS /
   SAMDOS-1 / all). A small heuristic infers the producer dialect
   from disk content (samdos2 presence / MasterDOS-specific field
   patterns / boot-sector signatures). Rules whose dialect tag
   doesn't match the detected dialect are skipped. `--dialect=X`
   overrides the heuristic when it guesses wrong.

4. **Library API: public**. `(di *DiskImage) Verify() VerifyReport`
   is a first-class part of the samfile API. The CLI subcommand is
   a thin formatter on top. Recovery features land later as
   separate consumers of the same report shape.

5. **Test corpus: hybrid**. ~5–10 hand-picked disks committed to
   `testdata/mgt/` (one per dialect, plus a couple of
   deliberately-corrupted negative cases) run in CI on every
   commit. A larger external corpus is fetched on demand via
   `scripts/fetch-corpus.sh` and exercised behind `-tags=corpus`;
   it does not gate CI.

## Architecture

### Public API

```go
// Verify runs all registered rules against di and returns a
// report describing the disk's structural state. The report is
// always populated (no errors are returned from Verify itself —
// individual rule failures become Findings, not errors).
func (di *DiskImage) Verify() VerifyReport

// VerifyReport is the result of running Verify on a DiskImage.
type VerifyReport struct {
    Dialect  Dialect    // detected dialect; never nil
    Findings []Finding  // all findings across all rules, ordered by rule registration
}

// Convenience predicates / filters. The full Findings slice
// stays available for callers who want raw access.
func (r VerifyReport) HasFatal() bool
func (r VerifyReport) HasStructural() bool
func (r VerifyReport) BySeverity(s Severity) []Finding
func (r VerifyReport) ByRule(ruleID string) []Finding
func (r VerifyReport) Filter(opts FilterOpts) []Finding

type FilterOpts struct {
    MinSeverity Severity   // zero value = all
    Rules       []string   // empty = all
    Dialect     Dialect    // zero value = whatever the report has
    Slot        *int       // nil = all slots
}
```

### Core types

```go
// Severity ranks findings by impact, lowest to highest.
type Severity int
const (
    SeverityCosmetic Severity = iota
    SeverityInconsistency
    SeverityStructural
    SeverityFatal
)

// Dialect identifies which DOS produced the disk.
type Dialect int
const (
    DialectUnknown Dialect = iota
    DialectSAMDOS1
    DialectSAMDOS2
    DialectMasterDOS
)

// Rule is a registered validity check. The Check function is
// invoked once per Verify run; it returns zero or more Findings.
// Rule values are immutable after registration.
type Rule struct {
    ID          string      // catalog-stable, e.g. "BODY-EXEC-DIV16K-MATCHES-DIR"
    Severity    Severity
    Dialects    []Dialect   // dialects the rule applies to; nil/empty = all
    Description string      // one-line summary for human output
    Citation    string      // file:line of the strongest evidence
    Check       func(ctx *CheckContext) []Finding
}

// CheckContext is the read-only environment passed to each rule.
// All disk inspection goes through it so rules don't reach into
// di directly — keeps rules unit-testable in isolation.
type CheckContext struct {
    Disk    *DiskImage
    Journal *DiskJournal     // pre-computed, shared across rules
    Dialect Dialect          // result of dialect detection
}

// Finding is one specific violation produced by one Rule.
type Finding struct {
    RuleID   string
    Severity Severity
    Location Location  // where on the disk; zero value if disk-wide
    Message  string    // human-readable, includes Expected vs Actual when applicable
    Citation string    // copied from the Rule for convenience
}

// Location pinpoints a Finding on the disk. Fields are optional —
// disk-wide findings leave them all zero; per-file findings set Slot;
// per-byte findings set Slot + Sector + ByteOffset.
type Location struct {
    Slot       int      // -1 if not applicable
    Sector     *Sector  // nil if not applicable
    ByteOffset int      // -1 if not applicable
    Filename   string   // copied from Slot's entry if Used, for messages
}
```

### Rule registry

Rules register at package init time:

```go
// rules_body_header.go
func init() {
    Register(Rule{
        ID:          "BODY-EXEC-DIV16K-MATCHES-DIR",
        Severity:    SeverityInconsistency,
        Dialects:    nil, // all dialects
        Description: "body-header byte 5 must equal dir-entry ExecutionAddressDiv16K (0xF2)",
        Citation:    "rom-disasm:22467-22484",
        Check:       checkBodyExecDiv16KMatchesDir,
    })
}
```

`Register` is package-private; rules ship with samfile and are not
user-extensible in v1. The registry is iterated in registration
order so report output is deterministic.

### Dialect detection

```go
// DetectDialect inspects di and returns the most likely dialect
// that wrote the disk. Returns DialectUnknown if the heuristic
// can't decide. Documented in detail in docs/disk-validity-rules.md
// §Dialect-notes; the heuristic looks at:
//   - boot sector presence and contents (samdos2 vs masterdos2)
//   - MGT future-and-past patterns in occupied slots
//   - MasterDOS-only field usage in dir entries
//
// Heuristic is deliberately conservative: any ambiguity returns
// DialectUnknown, which causes Verify to run only rules tagged
// `AllDialects`. The user can override with --dialect.
func DetectDialect(di *DiskImage) Dialect
```

### CLI

```
samfile verify -i IMAGE [flags]

Flags:
  --severity=LEVEL    Minimum severity to display
                      (cosmetic|inconsistency|structural|fatal).
                      Default: structural.
  --all               Shorthand for --severity=cosmetic.
  --rule=ID           Show only findings from rule ID. Repeatable.
  --dialect=NAME      Force dialect (samdos1|samdos2|masterdos).
                      Default: auto-detect.
  --json              Emit findings as JSON instead of prose.
  --quiet             Suppress all output; exit code only.
```

Exit codes:

- `0` — no findings at or above the displayed severity.
- `1` — at least one `fatal` finding.

`structural` / `inconsistency` / `cosmetic` findings do *not*
change exit code by themselves. (Builders who want strict CI
gating wrap the call: `samfile verify -i $disk --severity=structural --quiet && echo ok`.)

Default human-readable output is grouped by severity, then by
slot:

```
samfile verify: results for samdos-clean.mgt
detected dialect: SAMDOS 2

FATAL (1):
  [BOOT-T4S1-PRESENT] slot 0 (samdos2):
    Track 4 / Sector 1 is unallocated; ROM BOOT requires it.
    citation: rom-disasm:55501-55613

STRUCTURAL (1):
  [CHAIN-LINK-TERMINATOR] slot 3 (IN):
    Last sector Track 6 / Sector 4 has non-zero next-link
    (track=0x06, sector=0x05); chain does not terminate.
    citation: samdos/src/c.s:1306-1343

2 finding(s) above --severity=structural threshold
(1 fatal, 1 structural). 4 lower-severity finding(s) hidden;
run with --all to see them.
exit 1 (fatal present).
```

JSON output (`--json`) is one finding per object in a top-level
array, plus a metadata header. Schema documented in package godoc.

## Data flow

```
Load(file) → *DiskImage
   │
   ▼
DetectDialect(disk) → Dialect
   │
   ▼
ctx := &CheckContext{Disk, Journal: disk.DiskJournal(), Dialect}
   │
   ▼
for each Rule in registry:
    if len(Rule.Dialects) == 0 || ctx.Dialect ∈ Rule.Dialects:
        report.Findings = append(report.Findings, Rule.Check(ctx)…)
   │
   ▼
VerifyReport{Dialect, Findings}
   │
   ├── library caller: receives report directly
   └── CLI: filters by --severity / --rule, formats, writes to stdout
```

The journal is computed once and shared across all rules — rules
must not call `disk.DiskJournal()` themselves. Same for any other
expensive derivations a rule might want; if a second rule needs
the same intermediate (e.g. combined sector map), it goes on
`CheckContext` as a memoised field.

## Implementation order

Six implementation phases (each landing as its own PR against
samfile master), followed by an open-ended corpus-validation
pass:

1. **Foundation** — `Rule`, `Finding`, `VerifyReport`, `Severity`,
   `Dialect`, `Location` types; `Register` + iteration; CLI
   subcommand scaffold that prints "no rules registered yet";
   one trivial rule wired end-to-end as a smoke test. No real
   rule implementations.
2. **Dialect detection** — `DetectDialect` heuristic + tests
   against the committed corpus. Without this, dialect-scoped
   rules in subsequent phases can't be tested properly.
3. **Disk-level + directory-entry + sector-chain rules** — the
   ~25 rules that don't depend on file-type specifics. These
   exercise the foundation and shake out any API issues before
   we commit to 45 more.
4. **Body-header + FT_CODE rules** — ~16 rules. Includes the two
   PR-12-confirmed mirrors (`BODY-EXEC-DIV16K-MATCHES-DIR` and
   `BODY-EXEC-MOD16K-LO-MATCHES-DIR`) as the simplest
   demonstrations.
5. **FT_SAM_BASIC + array + screen + ZX-snapshot rules** — ~13
   rules. File-type-specific content checks.
6. **Boot file + cross-entry + dialect-specific + cosmetic
   tail** — remaining ~16 rules. After this lands, the catalog
   is fully realised.
7. **Corpus validation pass** — run the external corpus
   (`-tags=corpus`), classify each rule's empirical violation
   rate, demote severities / add dialect filters / adjust as
   needed. One PR per reclassification, citing the corpus disks
   that triggered the change. Continues indefinitely as the
   corpus grows.

Phases 3–6 are roughly equal-size; the foundation phase is
deliberately small so the abstractions stay reviewable in
isolation.

## Testing strategy

Each rule ships with two unit tests:

- **Positive** — fabricate a clean disk in-memory (via
  `NewDiskImage` + `AddCodeFile` / `AddBasicFile` / hand-written
  helpers), assert the rule produces zero findings.
- **Negative** — fabricate a disk that deliberately violates the
  rule (one targeted byte flip), assert the rule produces exactly
  one finding with the expected `RuleID`, severity, and location.

In-line fixture construction is preferred over checked-in
`.mgt` files for unit tests — keeps each test self-contained and
the violation visible at the test site. The committed corpus is
for integration testing: "all rules on a known-good disk produce
zero high-severity findings".

Two `go test` targets:

- `go test ./...` — exercises unit tests + the committed corpus
  in `testdata/mgt/`. Runs in CI on every commit.
- `go test -tags=corpus ./...` — additionally fetches and
  exercises the external corpus via `scripts/fetch-corpus.sh`.
  Runs locally and on a nightly CI job.

## What's *not* in v1

These are explicitly out of scope for the initial `verify` work.
They build on top of it and should be designed as separate
features:

- **Deleted-file listing** — surfacing slot entries with
  `Type == FT_ERASED` whose sector data still contains
  recognisable file content.
- **Sector-remnant salvage** — following dangling next-sector
  links from unallocated sectors into reconstructed chains.
- **Disk-clean** — zeroing unallocated regions to strip remnants.
- **Repair / fix-up commands** — `verify` reports; it never
  mutates the disk.
- **User-extensible rules** — third-party packages registering
  custom rules. v1's `Register` is package-private; opening it
  up needs a versioning story for rule IDs that we're not paying
  for upfront.
- **Multi-format input** — only MGT is accepted (same as today).
  EDSK is still rejected at `Load` time.

## Open questions deferred to plan-writing

These are implementation-level details that don't need to be
settled before plan-writing starts; they'll resolve naturally as
each phase's plan is written.

- Exact CLI flag long/short forms and help-text wording.
- Whether `Finding.Message` is constructed eagerly during Check
  or lazily via a printer interface.
- Whether the "empirical / convention" rules (e.g.
  `BASIC-VARS-GAP-INVARIANT`) get their own severity tier or
  stay tagged `cosmetic` with a documentation note.
- How the corpus fetcher handles broken upstream URLs (skip /
  fail / retry mirror).
- Whether `DetectDialect` returns confidence levels (e.g.
  `(Dialect, confidence float)`) or just the dialect.
