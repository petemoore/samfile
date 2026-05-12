# Verify Audit Framework — Design Spec

**Date:** 2026-05-12
**Status:** Draft (pre-implementation)

## Problem

Today's verify pipeline emits findings (failures) but not checks (attempts). Three consequences:

1. **No denominators.** A rule that fires 100 times looks the same whether it was applicable 100 times (100% fail) or 100,000 times (0.1%). Today we can't tell.
2. **No per-check context.** We don't capture which subjects were checked, what their attributes were, or what passed. So we can't ask "is this rule only failing on samdos2 disks?" or "does this fire on every ZX snapshot?".
3. **No association analysis.** Without (subject × rule × outcome × attributes) rows, we can't mine for patterns like "fails when file_type=ZX_SNAP AND first_track=4". Today's investigation is anecdotal — find a finding, eyeball it, guess.

The corpus-reclassification work (`docs/notes/corpus-reclassification-2026-05-12.md`) compensated by manual SQL queries, but every reclassification decision relied on guessing whether the residual fires were "real signal" or "noise from broken disks". With ~50 rules and ~250K potential checks across the 800-disk corpus, this doesn't scale.

## Goal

Build infrastructure that:

1. Records every check (pass and fail) with a full per-subject attribute snapshot.
2. Lets us compute pass-rate and conditional fail-rate slices trivially in pandas.
3. Surfaces patterns ("rule X fails 100% of the time when feature Y = Z") with high enough confidence to drive rule-fix decisions.
4. Reuses across rule changes — every catalog edit gets a fresh audit run.

The endpoint is not the framework. The endpoint is: **a small set of catalog and code fixes grounded in source code (samdos / ROM disasm / samfile) and validated against this corpus**.

## Non-goals

- Real-time / interactive verify. Batch only.
- Deep learning. Classical pandas + sklearn + mlxtend is sufficient at 800-disk scale.
- A new web UI. Markdown reports + SQLite for ad-hoc queries.
- Generalizing beyond MGT disks.

## Decisions (from brainstorming)

| Choice | Decision |
|---|---|
| Lifespan | Ongoing infrastructure (re-runs after every rule change). |
| Code split | Go-side instrumentation in samfile repo; Python analysis in `~/sam-corpus/`. |
| Instrumentation style | Rule struct grows `Applies(subject) → bool` + `Scope` field; framework drives iteration and emits per-subject Check events. |
| Output format from samfile | `samfile verify --format jsonl` emits one JSON object per Check; default text output unchanged. |
| Pattern-mining stack | pandas + sklearn `DecisionTreeClassifier` + mlxtend `apriori` / `association_rules`. |
| PR strategy | One bundled draft PR: framework + all rule fixes. |
| Stop conditions | Loop until a pass produces no high-confidence patterns. Ambiguous patterns documented under `needs-human` and skipped. |
| ROM disasm source | `~/git/sam-aarch64/docs/sam/sam-coupe_rom-v3.0_annotated-disassembly.txt`. |

## Architecture

```
[800 disks] ──► samfile verify ──JSONL──► sam-corpus/ingest.py ──► findings.db
              (rules emit Check events                              ├ checks  (new)
              for pass + fail + N/A)                                ├ findings (kept)
                                                                    └ disks   (kept)
                                                                       │
                                                                       ▼
                                                          sam-corpus/mine.py
                                                          ├ coverage.md       (fallback floor)
                                                          ├ disk-health.md    (fallback floor)
                                                          ├ conditional.md
                                                          ├ disk-clusters.md
                                                          └ patterns.md
                                                                       │
                                                                       ▼
                                                          autonomous fix loop
                                                          (hypothesis → source-ground → fix → re-run)
```

Four stages with clean boundaries:

1. **Go-side instrumentation** (samfile): rules declare scope + applicability; framework drives iteration; emits Check events.
2. **JSONL emission** (samfile): new `--format jsonl` CLI flag.
3. **Ingest** (sam-corpus/Python): JSONL → `checks` SQLite table with denormalized attribute columns.
4. **Mine + report** (sam-corpus/Python): five Markdown reports + a SQLite database for ad-hoc queries.

## Attribute schema

Three subject scopes; framework records `(scope, ref, outcome, attrs)` per check.

### Disk-scope

Per disk:
- `sha256`
- `dialect` (`samdos` / `samdos2` / `masterdos` / `other`)
- `boot_signature_present` (bool)
- `used_slot_count`
- `free_sector_count`
- `sector_map_populated_pct`
- `total_findings_fatal`, `total_findings_structural`, `total_findings_inconsistency`, `total_findings_cosmetic`
- `source_label` (from `manifest.csv`)

### Slot-scope

Per directory entry (0..79):
- `slot_index`
- `filename`
- `file_type` (decoded constant name, e.g. `FT_CODE`, `FT_SAM_BASIC`, `FT_SCREEN`, `FT_ARRAY`, `FT_ZX_SNAP`)
- `file_type_byte` (raw 0..255)
- `file_length` (decoded `LengthMod16K`)
- `page_offset_form` (`0x8000` / `0x4000` / `other`)
- `pages`
- `mgt_flags`
- `chain_length` (sector count)
- `sectors_count` (dir field, for cross-check vs chain)
- `first_track`, `first_sector`, `first_side` (0/1)
- `has_autorun_or_autoexec` (bool — auto-RUN line set for BASIC, auto-exec address set for CODE)
- `dir_mirror_populated` (bool — dir 0xD3..0xDB non-zero)
- `slot_is_erased` (Type==0)
- `file_type_info_hex` (FileTypeInfo bytes, hex string)

### Chain-step-scope

Per sector in a file's chain:
- `slot_index` (which file)
- `chain_position` (`first` / `intermediate` / `last` / `orphan`)
- `chain_index` (0..N-1)
- `track`, `sector`, `side`
- `next_track`, `next_sector`
- `on_sam_map` (bool — appears in side-0 / side-1 SAM sector-allocation maps)
- `on_dir_sam_map` (bool — appears in this dir entry's SectorAddressMap)
- `distance_from_dir_tracks` (min cylinder distance to tracks 0..3)

## Go-side changes (samfile repo)

### `Rule` struct refactor

Add `Scope` and `Applies`; change `Check` signature so it operates on one subject at a time and returns nil for pass.

```go
type SubjectScope int

const (
    DiskScope SubjectScope = iota
    SlotScope
    ChainStepScope
)

type Subject interface {
    Ref() string                  // e.g. "slot=5" or "track=4,sector=1"
    Attributes() map[string]any   // denormalized attribute snapshot
}

type Rule struct {
    ID          string
    Description string
    Severity    Severity
    Citation    string
    Scope       SubjectScope
    Applies     func(*CheckContext, Subject) bool   // NEW
    Check       func(*CheckContext, Subject) *Finding // NEW signature
}
```

### Framework

The verify driver enumerates subjects for each rule's scope, calls `Applies`, and on true calls `Check`. Per (rule, subject) it emits a `CheckEvent`:

```go
type CheckEvent struct {
    Version  int                    `json:"v"`
    Disk     string                 `json:"disk"`
    RuleID   string                 `json:"rule_id"`
    Scope    string                 `json:"scope"`
    Ref      string                 `json:"ref"`
    Outcome  string                 `json:"outcome"`  // pass | fail | not_applicable
    Attrs    map[string]any         `json:"attrs"`
    Finding  *Finding               `json:"finding,omitempty"`
}
```

The existing 51 rules each receive a mechanical refactor: extract the iteration logic (currently embedded in each rule's check function) into `Applies`, and reshape the check body to operate on one subject. Typical change is ~5 lines per rule.

### CLI

```
samfile verify --format jsonl -i disk.mgt > disk.jsonl
samfile verify -i disk.mgt                    # unchanged text output
```

Default behaviour preserved. JSONL is opt-in.

### Tests

- Each refactored rule keeps its existing unit tests, exercised through the new `Applies` / `Check` pair.
- New framework-level test asserting that `not_applicable` + `pass` + `fail` counts sum to the universe of subjects for that scope.

## Python-side changes (`~/sam-corpus/`)

### `ingest.py`

Reads `outputs-jsonl/*.jsonl` (one file per corpus disk), writes a new `checks` table to `findings.db`. Schema:

```sql
CREATE TABLE checks (
    disk TEXT,
    rule_id TEXT,
    scope TEXT,
    ref TEXT,
    outcome TEXT,              -- 'pass' | 'fail' | 'not_applicable'
    -- denormalized disk attributes
    dialect TEXT,
    boot_signature_present INTEGER,
    used_slot_count INTEGER,
    -- denormalized slot attributes (NULL if disk-scope)
    slot_index INTEGER,
    file_type TEXT,
    file_type_byte INTEGER,
    file_length INTEGER,
    page_offset_form TEXT,
    first_track INTEGER,
    first_sector INTEGER,
    first_side INTEGER,
    has_autorun_or_autoexec INTEGER,
    dir_mirror_populated INTEGER,
    chain_length INTEGER,
    -- denormalized chain-step attributes (NULL if disk/slot-scope)
    chain_position TEXT,
    chain_index INTEGER,
    track INTEGER,
    sector INTEGER,
    side INTEGER,
    on_sam_map INTEGER,
    on_dir_sam_map INTEGER,
    -- finding payload (NULL if pass / not_applicable)
    severity TEXT,
    message TEXT,
    citation TEXT
);
CREATE INDEX idx_checks_rule ON checks(rule_id);
CREATE INDEX idx_checks_outcome ON checks(outcome);
CREATE INDEX idx_checks_disk ON checks(disk);
```

`disks` and `findings` tables stay; existing `analyze.py` keeps working.

### `mine.py`

Emits five Markdown reports under `~/sam-corpus/analyses/`. Reports are produced in fail-safe order — the two fallback-floor reports first, then richer analyses.

#### Fallback floor (must always succeed)

**`coverage.md` — per-rule priority report:**

| Rule | Severity | Applies | Fails | Fail-rate | Disks affected |
|---|---|---:|---:|---:|---:|

Sorted by fail-rate desc. Plain pandas groupby; no clustering, no ML.

**`disk-health.md` — per-disk priority report:**

| Disk | Total findings | Fatal | Structural | Structural pass-rate | Distinct rules fired |
|---|---:|---:|---:|---:|---:|

`structural_pass_rate` = `passes / (passes + fails)` over checks where the rule severity is `structural`. Sorted asc — worst disks first. Disks with `structural_pass_rate < 0.5` and high `distinct_rules_fired` are the candidate "not really a disk" cluster.

#### Richer analyses (best-effort, graceful failure)

**`conditional.md`:** per-rule, table of `(attribute, value, conditional_fail_rate, baseline_fail_rate, support, Z-score)`. Rows surface only when `|Z| > 2` OR `conditional_fail_rate ∈ {0%, 100%}` OR `support ≥ 10`.

**`disk-clusters.md`:** hierarchical clustering of disks by rule-fire pattern (sklearn `AgglomerativeClustering`). Co-occurrence matrix between rules with Jaccard ≥ 0.7. Surfaces the "totally broken disks" cluster explicitly.

**`patterns.md`:** per-rule `DecisionTreeClassifier(max_depth=4)` trained on attributes → fail/pass; emits human-readable splits. Apriori-mined association rules at confidence ≥ 0.9, support ≥ 10 distinct disks.

Each rich report ends with a **`needs-human`** section: patterns where signal is real but I cannot find a source-code citation that confirms the explanation.

### `run_audit.sh` (in samfile `tools/audit/`)

```
1. cd ~/git/samfile && go build -o ~/sam-corpus/samfile-audit ./cmd/samfile
2. cd ~/sam-corpus && mkdir -p outputs-jsonl analyses
3. for disk in disks/*.mgt; do
     ~/sam-corpus/samfile-audit verify --format jsonl -i "$disk" \
       > "outputs-jsonl/$(basename "$disk" .mgt).jsonl"
   done
4. python3 ~/git/samfile/tools/audit/ingest.py
5. python3 ~/git/samfile/tools/audit/mine.py
```

Idempotent — re-runnable on every framework or rule change.

## Autonomous fix loop

After the framework PR's code is in place and reports run, the loop:

```
1. Run analyses.
2. For each high-confidence pattern:
     a. Form hypothesis.
     b. Read source (~/git/samdos/src/*.s + ROM disasm + samfile.go + rules_*.go).
     c. Confirm or refute. If confirmed AND fix is unambiguous → patch the rule on the feature branch.
     d. If refuted → discard. If ambiguous → append to needs-human.md.
3. Re-run analyses.
4. If new high-confidence patterns surface, go to 2. Else stop.
```

**High-confidence threshold (both required):**

(a) Statistical signal: one of
    - `conditional_fail_rate ≥ 80%` on a non-trivial slice (support ≥ 10 distinct disks)
    - `conditional_fail_rate = 0%` on a non-trivial slice (support ≥ 10)
    - Apriori association rule with `confidence = 1.0`, `support ≥ 10 distinct disks`

(b) Citation: a line range in `~/git/samdos/src/*.s` or the ROM disasm or `samfile.go` whose semantics, when read in context, explain the pattern.

Patterns missing either criterion → `needs-human.md`. Never act without both.

## PR / branch

Single feature branch `feat/verify-audit-framework`. Single draft PR containing:

- The Go-side framework refactor.
- All 51 rules refactored.
- The JSONL CLI flag.
- Python tooling under `tools/audit/` in the samfile repo (`ingest.py`, `mine.py`, `run_audit.sh`) — version controlled, ships with the PR. `~/sam-corpus/` is the runtime workspace where disks, JSONL outputs, and the SQLite DB live; the audit scripts are invoked from there. Existing `~/sam-corpus/{gather,analyze}.py` are untouched.
- All rule-fix commits, each with citations to source.
- Updated catalog (`docs/disk-validity-rules.md`) where rule semantics change.
- A summary report in the PR body linking to the five analysis reports.

Draft remains draft until user review. Standard CLAUDE.md global rules apply (monitor CI, fix failures autonomously, never mark ready without explicit approval).

## Risks

- **Framework refactor scope.** Touching all 51 rules is large. Mitigation: mechanical per-rule shape change; each rule keeps its tests; CI gates the framework PR.
- **Hallucinated grounding.** I could "find" a citation that doesn't actually support the rule fix. Mitigation: every rule-fix commit must quote the source line range verbatim in the commit message; user reviews diff before merge.
- **Pattern-mining noise.** Decision trees and Apriori can produce spurious rules at small support. Mitigation: high support + confidence thresholds; fallback-floor reports stand alone if richer reports turn noisy.
- **Loop divergence.** Re-running could surface infinite "patterns" that are all artifacts of earlier fixes. Mitigation: stop condition is strict (no NEW high-confidence patterns).

## Open questions for user review

None at design time. Implementation may surface choices that need a call; those land in needs-human.md rather than blocking.
