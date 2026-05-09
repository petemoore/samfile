# samfile verify — corpus-driven reclassification

Date: 2026-05-12. Branch: `samfile-corpus-reclass` (16 commits ahead of
`origin/master`). Author: Claude, three rounds of propose/review/implement.

## TL;DR

We pointed `samfile verify` at a deduplicated corpus of **800 real-world
SAM disks** (sourced from `~/Downloads/GoodSamC2/` plus scattered local
disks; deduped by SHA-256 of the 819,200-byte content; format-filtered
on disk size) and used the resulting fire-rates to reclassify rules
whose severities were demonstrably miscalibrated. Each candidate change
was justified against SAMDOS source (`~/git/samdos/src/`), the ROM
disassembly, the Tech Manual, or the catalog at
`~/git/samfile-corpus-reclass/docs/disk-validity-rules.md` — never on
prevalence alone.

The corpus showed three pathological clusters: a single off-by-one mask
bug producing ~40k false structural findings; a `CROSS-NO-SECTOR-OVERLAP`
rule producing 162,727 of the 165,620 fatal findings (98% of fatal
volume) on disks that load fine under SAMDOS; and a six-rule
"body-header mirror" cluster firing as `inconsistency` on 91–97% of
every dialect (samdos1, samdos2, masterdos, unknown). Across three
iterations we landed **16 commits** that net out to:

- **Fatal volume: −98.8%** (165,620 → 1,937)
- **Disks with no fatal findings: 25.5% → 77.6%**
- **Disks clean at the default display threshold (no fatal + no
  structural): 0.1% → 18.6%** (1 → 149 of 800)
- **Total rule count unchanged at 51** — every change is FIX (3 sites,
  one mask bug), DEMOTE (10 rules), SCOPE (4 rules), REWORD (4 rules),
  KEEP (5 rules). No rule was added or removed.

We stopped at iter-3 because every high-fire rule now has a defensible
classification grounded in SAMDOS source. The 4 follow-up issues we
identified are catalog corrections and a `CROSS-NO-SECTOR-OVERLAP`
per-disk summary-mode feature — none of them are severity questions.

---

## The corpus

- **Size**: 800 unique 819,200-byte `.mgt` images (`SELECT COUNT(*) FROM
  disks` against `~/sam-corpus/findings.db` returns 800).
- **Sources**: `~/Downloads/GoodSamC2/` (the bulk) plus scattered local
  disks including the M0 boot disk `~/sam-corpus/disks/test.mgt`.
- **Filter**: 819,200 bytes after stripping the 22-byte SAD signature
  header where present. Anything else (.dsk variants, alternate
  geometries) was excluded.
- **Dedupe**: SHA-256 of the 819,200-byte content.

### Dialect distribution

| Dialect | Disks |
|---|---:|
| unknown | 470 |
| masterdos | 151 |
| samdos2 | 114 |
| samdos1 | 65 |

The 59% "unknown" rate is itself a signal: most disks in the wild
weren't produced by a single canonical SAMDOS-or-MasterDOS writer; they
are ROM-SAVE'd or built by third-party tools (Lemmings/Demo packagers,
sample-disk authors, Chris White's Lemmings tooling, etc.). This is
the load-bearing fact behind several of the demotes below.

---

## Headline metric table

All four severity columns and the disk-health rows come from
`~/sam-corpus/findings.db` queries shown in the per-iteration reports
(`~/sam-corpus/report-iter-{0,1,2,3}.md`).

| Metric | iter-0 (before) | iter-3 (after) | Δ |
|---|---:|---:|---:|
| fatal findings | 165,620 | 1,937 | −98.8% |
| structural findings | 65,254 | 11,532 | −82.3% |
| inconsistency findings | 42,798 | 172,949 | +304.1% (severity shift) |
| cosmetic findings | 25,967 | 41,154 | +58.5% (severity shift) |
| **total findings** | 299,639 | 227,572 | −24.1% |
| disks with zero findings | 1 (0.1%) | 6 (0.8%) | +5 |
| disks with no fatal | 204 (25.5%) | 621 (77.6%) | +417 |
| disks with no fatal + no structural | 1 (0.1%) | 149 (18.6%) | +148 |

The +304%/+58.5% growth on inconsistency/cosmetic is the headline
*reclassification* effect (severity-only demotes); the −24.1% on total
volume is the *bug-fix* effect (the side-1 mask FIX and the
DIR-TYPE-BYTE-IS-KNOWN double-fire FIX). Together they mean the
default-display output (severity ≥ structural by current convention)
now shows roughly 1/30th of the noise it did before, mostly comprising
real load-time hazards rather than mask-bug false positives or
SAMDOS-canonical conventions that 90%+ of real disks ignore.

---

## Per-iteration summary

### Iter-1 — find the high-volume bugs and miscalibrations (11 commits)

Proposal: `~/sam-corpus/proposal-iter-1.md` (4 FIX, 3 DEMOTE, 4 SCOPE,
2 REWORD, 2 KEEP, 1 INVESTIGATE). Review: `~/sam-corpus/review-iter-1.md`
(13 of 15 approved as-is; 2 BODY-EXEC entries had a text inversion that
landed only as a wording fix — code was sound).

The big-ticket items were the side-agnostic cyl-mask FIX shared by
`DIR-FIRST-SECTOR-VALID`, `DISK-DIRECTORY-TRACKS`, and
`CROSS-DIRECTORY-AREA-UNUSED` (`(t & 0x7F) < 4` wrongly excluded side-1
cylinders 0..3 from the data area — see commit `f619d9a`), and the
`CROSS-NO-SECTOR-OVERLAP` fatal→structural demote that alone removed
~163k fatal findings (commit `bff412c`; the rule was tagged fatal but
SAMDOS LOAD doesn't consult the SectorAddressMap, per
`~/git/samdos/src/b.s:104-110`). Eight further commits implemented the
SCOPE/REWORD/DEMOTE entries the reviewer signed off on.

The `BODY-MIRROR-AT-DIR-D3-DB` INVESTIGATE was kept open for iter-2 —
93% violation across all four dialects pointed to a catalog claim that
needed source verification, not a code change.

### Iter-2 — body-mirror cluster + DIR-ERASED-IS-ZERO (2 commits)

Proposal: `~/sam-corpus/proposal-iter-2.md` (6 DEMOTE + 1 reword in one
commit; DIR-ERASED-IS-ZERO DEMOTE + reword in another; 3 INVESTIGATE-
defer; 1 KEEP). Review: `~/sam-corpus/review-iter-2.md` — same actions
approved but the proposer's **technical justification for the
body-mirror cluster demote was inverted**. The fix was a wording rewrite
in commit `2b8aa1a`, citing the correct LOAD path through
`gtfle → hconr → txhed` (dir-derived; body bytes 0..8 are skipped by
`ldhd` and discarded, not authoritative on LOAD).

Both commits landed. The chain-vs-map structural-fan-out cluster
(`CROSS-NO-SECTOR-OVERLAP`, `DIR-SECTORS-MATCHES-CHAIN`,
`CHAIN-MATCHES-SAM`) was deferred to iter-3 — demoting it without a
disk-level summary mode would have left the per-sector fan-out noise
intact regardless of severity.

### Iter-3 — chain/SAM/Sectors cluster (3 commits)

Proposal: `~/sam-corpus/proposal-iter-3.md` (3 DEMOTE for the deferred
trio, 2 REWORD bundled, 5 KEEP, 1 INVESTIGATE-defer). Review:
`~/sam-corpus/review-iter-3.md` — accepted 2 of 3 DEMOTEs but rejected
the third (`DIR-SECTORS-MATCHES-CHAIN`) on the grounds that the spec's
`structural` legend at `disk-validity-rules.md:28-34` is broader than
"SAMDOS-authority-only" (it specifically says "disk-walk invariant;
violation produces undefined behaviour or sector reuse"), and
`samfile.go:743-754` reads exactly `fe.Sectors` chunks — too-large
Sectors reads past the (0,0) terminator into garbage; too-small Sectors
silently truncates. Both are real load-time consumer corruption.

The implemented split:
- Commit `86c08d1`: CROSS-NO-SECTOR-OVERLAP DEMOTE + REWORD (structural
  → inconsistency).
- Commit `3287b03`: CHAIN-MATCHES-SAM DEMOTE (structural →
  inconsistency).
- Commit `f7a1d42`: DIR-SECTORS-MATCHES-CHAIN — message reword only,
  severity stays structural per reviewer's amendment.

---

## All 16 commits

In landing order (oldest first), with the load-bearing citation each
change rests on:

1. **`f619d9a`** — verify: FIX DIR-FIRST-SECTOR-VALID /
   DISK-DIRECTORY-TRACKS / CROSS-DIRECTORY-AREA-UNUSED — drop
   side-agnostic cylinder mask. Citation: Tech Manual L4340-4343
   ("first 4 tracks of side 0 are the directory"); `samfile.go:389-394`
   (range check uses correct ranges).
2. **`bff412c`** — verify: DEMOTE CROSS-NO-SECTOR-OVERLAP (fatal →
   structural). Citation: `~/git/samdos/src/b.s:104-110` (LOAD chain
   walk reads bytes 510-511 only); `~/git/samdos/src/c.s:895-951` (SAVE
   allocator merges SAM map).
3. **`6d14929`** — verify: DEMOTE BOOT-OWNER-AT-T4S1 (fatal →
   structural). Citation: `rules_boot.go:36-39` (rule's own
   self-acknowledgement); spec §Decisions item 1.
4. **`ac9d120`** — verify: DEMOTE BOOT-SIGNATURE-AT-256 (fatal →
   structural). Citation: ROM disasm 20582-20598 (BTCK signature
   check).
5. **`365fa7e`** — verify: SCOPE BODY-EXEC-DIV16K-MATCHES-DIR — skip
   when body[5]==0xFF. Citation: ROM disasm 22467-22484 (auto-exec
   gate); catalog §BODY-BYTES-5-6-CANONICAL-FF.
6. **`ffb30da`** — verify: SCOPE BODY-EXEC-MOD16K-LO-MATCHES-DIR — same
   `body[5]==0xFF` skip. Citation: as #5.
7. **`07d3dc1`** — verify: SCOPE CODE-FILETYPEINFO-EMPTY — accept 0x20.
   Citation: ROM disasm 22070-22074 (HDCLP 25-byte space-fill);
   empirical 99% concentration on 0x20.
8. **`9a1ae0b`** — verify: SCOPE BASIC-LINE-NUMBER-BE /
   BASIC-STARTLINE-FF-DISABLES — widen to 1..65535. Citation:
   `sambasic` parses uint16 line numbers without complaint; corpus
   evidence of legitimate 20000–65000 line numbers.
9. **`67491f4`** — verify: REWORD BASIC-STARTLINE-WITHIN-PROG — fire
   only when auto-RUN line exceeds the highest saved line. Citation:
   SAM BASIC RUN semantics use NEXT-LINE-GE not LINE-EQ.
10. **`277aa42`** — verify: REWORD SCREEN-LENGTH-MATCHES-MODE — accept
    canonical ROM SCREEN$ palette+sysvars trailer. Citation: Tech
    Manual L2156 ("spare 8K following a MODE 3/4 screen"); empirical
    75% of fires are 24576+41 = 24617.
11. **`27216be`** — verify: FIX DIR-TYPE-BYTE-IS-KNOWN — skip Type==0.
    Citation: catalog `disk-validity-rules.md:310-312` lists type 0 as
    allowed in the rule's own test sketch; 100% double-fire with
    DIR-ERASED-IS-ZERO.
12. **`2b8aa1a`** — verify: DEMOTE body-mirror cluster (6 rules:
    BODY-MIRROR-AT-DIR-D3-DB, BODY-TYPE-MATCHES-DIR,
    BODY-LENGTHMOD16K-MATCHES-DIR, BODY-PAGEOFFSET-MATCHES-DIR,
    BODY-PAGES-MATCHES-DIR, BODY-STARTPAGE-MATCHES-DIR) inconsistency →
    cosmetic. Citation: `~/git/samdos/src/c.s:1376-1379` (`gtfle`
    populates 9-byte cache from dir+211); `~/git/samdos/src/h.s:336-361`
    (`hconr` repopulates from `uifa` = dir-derived);
    `~/git/samdos/src/f.s:494-497` (`ldhd` reads and discards body
    bytes 0..8 via `lbyt`); `~/git/samdos/src/h.s:38-56` (`txhed`
    transmits dir entry to ROM HDR/HDL).
13. **`6d282a7`** — verify: DEMOTE + REWORD DIR-ERASED-IS-ZERO
    (structural → inconsistency). Citation: `~/git/samdos/src/c.s:1133-1143`
    (`fdhf` consumer treats Type==0 as free; orphaned filename/chain
    are normal DEL/ERASE archaeology).
14. **`86c08d1`** — verify: DEMOTE + REWORD CROSS-NO-SECTOR-OVERLAP
    (structural → inconsistency). Citation: as #2, plus the spec's
    inconsistency legend ("two views of the same fact disagree") fits
    merged-SAM-map vs per-file chain-walk.
15. **`3287b03`** — verify: DEMOTE CHAIN-MATCHES-SAM (structural →
    inconsistency). Citation: `~/git/samdos/src/b.s:104-110` (LOAD
    never reads per-slot SAM); `~/git/samdos/src/c.s:1306-1343` (`cfsm`
    is SAVE-side bookkeeping); 99.7% co-firing with
    DIR-SECTORS-MATCHES-CHAIN.
16. **`f7a1d42`** — verify: REWORD DIR-SECTORS-MATCHES-CHAIN (severity
    unchanged at structural; reviewer-amended). Citation:
    `samfile.go:743-754` (consumer reads exactly `fe.Sectors` chunks
    and never checks the (0,0) terminator — too-large/too-small
    Sectors both produce undefined behaviour); spec
    `docs/disk-validity-rules.md:28-34` structural legend includes
    "violation produces undefined behaviour or sector reuse", which is
    broader than SAMDOS-authority-only.

---

## What was deferred and why

Four follow-up issues were filed as INVESTIGATE-defer or noted in
iteration reviews — none of them are severity questions and none
require another reclassification loop.

- **CROSS-NO-SECTOR-OVERLAP per-disk summary mode** (iter-2-deferred,
  carried through iter-3). The rule fires once per overlapping sector,
  producing 162,727 findings across only 299 disks (mean 544
  fires/disk; max 1,560 on Comms Loader). After the iter-3 demote to
  inconsistency this no longer drowns out higher tiers, but the
  per-sector fan-out is still poor UX. Right shape: one disk-level
  finding plus per-sector detail behind `--verbose`. This is a feature,
  not a severity reclassification.
- **Structural-fan-out follow-up** for CROSS-NO-SECTOR-OVERLAP /
  DIR-SECTORS-MATCHES-CHAIN / CHAIN-MATCHES-SAM. Iter-2's correlation
  analysis showed 291/299 overlap-cluster disks also fire
  DIR-SECTORS-MATCHES-CHAIN and 315/316 fire CHAIN-MATCHES-SAM — i.e.
  the three rules catch the same writer-side bookkeeping disagreement
  from three angles. After the summary-mode feature lands, re-measure
  whether the three rules still provide independent signal on the
  ~8–29 disks that hit one without the others.
- **Catalog §0 narrative inversion** (iter-2 reviewer flagged; iter-3
  INVESTIGATE-defer). The text at `disk-validity-rules.md:113-119`
  describes the loaded-exec path as "from the body-header byte 5
  (SAMDOS's `hd001..` cache → HDL via dschd / hconr)" — but `hconr`
  populates from `uifa+*` (dir-derived), not from the body. This is
  out of scope for the reclassification work but should be corrected
  in a separate catalog-hygiene PR. See §"The catalog inversion we
  found" below.
- **BODY-EXEC body[5..6] BASIC-mirror investigation** (issue #21
  referenced in the proposal). Iter-1's BODY-EXEC pair scoped the rule
  to skip `body[5]==0xFF`; the remaining real-mismatch cases need a
  follow-up against the ROM RUN/GOTO line-lookup path to confirm
  whether body bytes 5-6 ever feed BASIC's auto-exec.
- **`samfile.go:390` debug.PrintStack noise** (issue #19, separate
  from this work). Noted for completeness; not touched by any commit
  here.

---

## What we deliberately did NOT change (KEEPs)

Five rules with high fire rates were explicitly kept after analysis
because their residuals are real:

- **DIR-FIRST-SECTOR-VALID** (723 fires / 157 disks / 19.6%
  post-iter-1). Iter-1 fixed the mask bug; the remaining fires
  enumerated in the iter-3 evidence (`track=0x02, sector=4`;
  `track=0x96, sector=235`; `track=0xff, sector=255`; etc.) are
  out-of-range first-sector pointers — real corruption.
- **DISK-SECTOR-RANGE** (688 fires / 153 disks / 19.1%). KEPT in
  iter-1, iter-2 and iter-3. The (0,0) terminator is special-cased at
  `rules_disk.go:144-146`; the residual `sector 0x00` /
  `sector 0x56` / etc. messages are out-of-range chain links.
- **BASIC-VARS-GAP-INVARIANT** (937 fires / 197 disks / 24.6%). 77%
  of fires are the cross-dialect signal "boot-classified as masterdos
  but file written with samdos2 gap size" — exactly what the rule
  was designed to flag. Already cosmetic; correct tier.
- **BASIC-LINE-NUMBER-BE** (690 fires / 265 disks / 33.1%
  post-iter-1). Iter-1 widened the range to 1..65535; residual fires
  are line=0 (313 of 690 — genuine errors) and parse failures on
  corrupted programs.
- **BASIC-STARTLINE-FF-DISABLES** (752 fires / 232 disks / 29.0%).
  Same family as BASIC-LINE-NUMBER-BE; iter-1's widening took care
  of the false positives.

---

## Where the reviewer overruled the proposer

Three substantive proposer-vs-reviewer disagreements; in each case the
code change Pete will see in the PR is what the reviewer recommended,
not what the proposer originally proposed.

- **Iter-1 BODY-EXEC pair**: the proposer mislabeled the most common
  firing pattern. The proposal claimed `dir=N, body=0xFF` was the
  bulk of corpus fires and "the genuinely-inconsistent case" was the
  inverse. The reviewer (`review-iter-1.md:154-214`) verified against
  the DB that the largest single bucket was actually
  `body=0x00, dir=0xff` (727 fires) — a *third* pattern not covered
  by the proposer's bullets. Net outcome: the code change
  (`if body[5]==0xFF: skip else compare`) was sound and landed
  unchanged; only the proposal narrative was reworded before
  implementation.
- **Iter-2 body-mirror cluster**: the proposer's claim that "body wins
  on LOAD because `ldhd` overwrites the dir cache from the body" was
  **inverted**. The reviewer (`review-iter-2.md:24-95`) verified
  `~/git/samdos/src/f.s:494-497` directly: `ldhd` loops 9 times
  calling `lbyt` (`~/git/samdos/src/c.s:557-570`), which reads one
  byte from the chain-driven disk-buffer and `jp incrpt` — it
  returns the byte in `A` and does not store it anywhere. So `ldhd`
  *skips* the body header rather than copying it into the cache. The
  correct framing is "body bytes 0..8 are unused on LOAD; the dir
  mirror feeds ROM via `gtfle → hconr → txhed`" (citation chain in
  the iter-2 review §Per-entry verdict). Same DEMOTE conclusion
  (inconsistency → cosmetic) but for the opposite reason. This is
  also what surfaced the §0 catalog-text inversion below.
- **Iter-3 DIR-SECTORS-MATCHES-CHAIN**: the proposer argued the rule
  should follow CROSS-NO-SECTOR-OVERLAP and CHAIN-MATCHES-SAM down to
  inconsistency because the catalog's "Source authority" line says
  `samfile-implicit` and SAMDOS LOAD doesn't read `fe.Sectors`. The
  reviewer (`review-iter-3.md:95-211`) pointed out that the spec
  legend for `structural` at `disk-validity-rules.md:28-34` is broader
  than the SAMDOS-only framing the proposer imported from `fatal` —
  the legend specifically says "disk-walk invariant; violation
  produces undefined behaviour or sector reuse", and `samfile.go`'s
  count-driven loop (`samfile.go:743-754`) produces undefined
  behaviour on either side of the disagreement. Net outcome: only
  the message reword landed (commit `f7a1d42`); severity stayed
  structural.

---

## The catalog inversion we found

While verifying the iter-2 body-mirror cluster demote, the iter-2
reviewer's source-trace surfaced a pre-existing inversion in the
catalog narrative at `disk-validity-rules.md:113-119` (the §0 "PR-12
hypotheses" verification text). The catalog claims `ldhd` *overwrites*
the in-RAM `hd001..page1` cache from body bytes — i.e. that the body
header is read into the cache on LOAD. Source verification shows the
opposite:

- `~/git/samdos/src/f.s:494-497` (`ldhd`): loops 9 times calling
  `lbyt`. `lbyt` is at `~/git/samdos/src/c.s:557-570` — it reads one
  byte from the disk-buffer-via-chain pointer, returns it in `A`, and
  `jp incrpt`. **It never stores the byte anywhere.** The 9-byte body
  header is read past, not into the cache.
- `~/git/samdos/src/h.s:74-90` (`dschd`): called by every LOAD/VERIFY
  path. Calls `ldhd` first (to advance past the header), then stores
  caller-provided `hkhl/hkbc/hkde` (BASIC LOAD destination registers
  captured at the `hook:` entry) into `hd0d1 / pges1 / hd0b1`. **It
  does not populate `hd001` (type byte) or `page1` (start-page byte)
  from anywhere — those stay as whatever the dir-side path put there.**
- `~/git/samdos/src/h.s:336-361` (`hconr`): reloads
  `hd001 / page1 / hd0d1 / pges1 / hd0b1` from `uifa+*` (dir-derived
  data populated by `gtfle` at `~/git/samdos/src/c.s:1376-1379` from
  dir bytes 0xD3-0xDB). The dir mirror is what feeds the cache.
- `~/git/samdos/src/h.s:38-56` (`txhed`): transmits 48 bytes from
  `difa` (the dir-entry buffer) into ROM's HDL/HDR area. Body bytes
  0..8 never enter ROM's view directly.

So on LOAD, **the dir mirror is authoritative** — the body header bytes
0..8 are SAVE-time decoration. A body↔dir mismatch has zero LOAD-time
consequence (which is why the body-mirror cluster demoted cleanly to
cosmetic), but the *reason* is the opposite of what the catalog
currently says. This inversion was discovered through corpus review —
the proposer originally walked the wrong direction, and only the
reviewer's independent source verification caught it. The catalog text
itself is out of scope for this PR; it's INVESTIGATE-defer #3 in the
deferred list above.

---

## Final state

- **51 rules registered** (`rules_*.go`, unchanged from iter-0).
- **M0 boot disk stays clean** (`~/sam-corpus/disks/test.mgt`, samdos2,
  zero findings). The regression gate held every iteration — every
  proposal documents an explicit M0-safety analysis (see
  `proposal-iter-{1,2,3}.md` "Cross-cutting observation: M0 regression
  analysis" sections), and every iteration verified zero new findings
  on M0 before remeasuring the corpus.
- **78% of corpus shows no fatal findings** (621/800, up from 204/800 =
  25.5%).
- **18.6% of corpus is clean at the default-display threshold** (no
  fatal + no structural) — 149/800, up from 1/800 = 0.1%. (The default
  display threshold in `samfile verify` is currently
  `severity ≥ structural`; the inconsistency and cosmetic tiers are
  hidden unless `--verbose` or `--severity` is passed.)
- **Total volume −24.1%** (299,639 → 227,572), with most of the drop
  attributable to the two FIX commits (`f619d9a` and `27216be`)
  removing ~43k false positives rather than the DEMOTEs (which
  preserve every finding).

Verified against the live DB on disk:

```
$ sqlite3 ~/sam-corpus/findings.db \
    "SELECT severity, COUNT(*) FROM findings GROUP BY severity;"
cosmetic|41154
fatal|1937
inconsistency|172949
structural|11532
```

---

## Future work

- **Three INVESTIGATE-deferred follow-up issues** (listed in "What was
  deferred and why" above): CROSS-NO-SECTOR-OVERLAP summary mode;
  structural-fan-out re-measure post-summary-mode; catalog §0 hconr
  narrative correction.
- **Cut samfile 3.1.0** — Pete to do.
- **Bump sam-aarch64 to samfile 3.1.0** — Pete to do.
- **Phase 7 onwards is ongoing** — as the corpus grows (more disks,
  more dialects, more unusual writers) we may surface further
  reclassifications. The corpus + findings DB are checked into
  `~/sam-corpus/` so future iterations can re-run the same
  propose/review/implement loop against an expanded dataset.

---

## Appendix: M0 regression gate

Every iteration's proposal includes an explicit M0 regression analysis
verifying that none of the proposed changes introduce a finding on
`~/sam-corpus/disks/test.mgt`. The gate held end-to-end:

- iter-0 M0 finding count: **0**
- iter-1 M0 finding count: **0** (FIXes widen acceptance; DEMOTEs are
  severity-only; SCOPEs narrow firing sets; M0's dir layout — slot 0
  Type=0x13/CODE, FirstSector=(4,1), dir 0xDD-0xE7 all zero, body
  bytes 0..8 match dir mirror byte-for-byte, no overlapping sectors —
  satisfies every new constraint).
- iter-2 M0 finding count: **0** (body-mirror demotes don't fire on M0
  because M0 IS the canonical mirror disk; DIR-ERASED-IS-ZERO doesn't
  fire because slot 0 is Type=0x13/CODE).
- iter-3 M0 finding count: **0** (chain/SAM/Sectors demotes don't fire
  on M0 because M0 has Sectors=49 matching a 49-sector chain to (0,0),
  per-slot SAM matching the chain walk, and no overlap).

The gate is reproducible by running the iter-3 `samfile verify` build
against `~/sam-corpus/disks/test.mgt` directly.
