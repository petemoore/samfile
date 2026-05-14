# SAM BASIC round-trip caveats

When `samfile basic-to-text` followed by `samfile text-to-basic` does
not reproduce the original body bytes exactly, the difference falls
into one of four buckets:

1. **Lexer bug** — text-to-basic produced the wrong bytes. Fix the
   lexer.
2. **Detokeniser bug** — basic-to-text didn't faithfully reproduce
   what the SAM ROM's `LIST` command would have shown. Fix
   `sambasic.SAMBasic.Output()`.
3. **Acceptable variation** — the difference has no functional
   effect (e.g. the SAM ROM rebuilds the value at LOAD time, or
   multiple encodings represent the same value). Add a schema rule
   that permits the variation, plus an entry in this file
   documenting why it's safe.
4. **Corrupt corpus body** — the source disk's bytes are
   genuinely invalid (e.g. truncated, sector errors, mis-typed
   file). Document and exclude.

This document records every acceptable variation (bucket 3),
with the schema rule that permits it and the rationale that the
variation has no functional consequence.

---

## Triage process

This is the **single source of truth** for how the round-trip fix
loop is run. Process notes that live only in memory or chat
scrollback are unreliable and may be lost.

### 1. Run the corpus test and print the scoreboard

```bash
go test -tags corpus -count=1 -run TestCorpusRoundTrip ./sambasic
```

The summary line reports
`match:N diverged:N detok-error:N parse-error:N file-error:N panic:N`.
After **every** change (including no-op refactors), emit a 4-line
scoreboard:

```
Corpus status:
  Total:    7028 BASIC programs
  Passing:  <match count>
  Failing:  <Total - match count>
  Δ:        +/- since previous iteration
```

`Total` and `Passing` come from the summary; `Failing` is the sum
of all non-`match` outcomes. Δ is the signed delta from the
previously-reported figure in this session (`+0` is still worth
printing — silence is ambiguous).

### 2. Pick the first failure

Take the **first** failure listed in the test output. Do not sort
or cluster-rank by frequency. Every bug needs fixing eventually,
and most fixes incidentally clear other failures from the same
root cause — pre-categorisation is overhead with no net saving.

### 3. Reproduce it locally

Identify the disk and file from the failure line, e.g.:

```
[diverged] 18 Rated Poker for 512k (19xx) (Supplement Software).mgt/FILE3:
    got 21851 bytes, want 21853 bytes, 6899 mismatches
```

Reproduce in-process: load the disk, get `f.Body`, run it through
`bodyToText` then `sambasic.ParseTextSchema`, compare via
`sambasic.Conform`. The boilerplate is small — keep a scratch
`/tmp/probe-*.go` around. The goal is to extract:

- The first divergence offset and the bytes around it (got and
  want, in hex).
- The decoded interpretation (line number, line body, what each
  byte means — keyword tokens vs literals vs FP-form, etc.).
- The source-text line that produced the divergence.

### 4. Classify into one of the four buckets

Form a hypothesis: which bucket does this divergence belong to?

| Bucket | Indicator |
|--------|-----------|
| 1. Lexer bug | Source text is correct SAM BASIC; the lexer's output disagrees with what TOKMAIN would produce. Fix `sambasic/lex.go`. |
| 2. Detokeniser bug | The text `basic-to-text` produced doesn't match what SAM ROM's `LIST` would have shown. Fix `sambasic.SAMBasic.Output()`. |
| 3. Acceptable variation | Two encodings have identical runtime effect (e.g. ROM rebuilds the value at LOAD). Add a schema rule and an entry below. |
| 4. Corrupt / non-stock corpus body | The disk's bytes aren't reachable from the stock editor (e.g. produced by an external tool, hand-poked, sector errors). Exclude the file. |

### 5. Before designing around bucket 3 or 4: verify in SimCoupé

If the hypothesis is "acceptable variation" or "corrupt body",
**do not** add escape mechanisms to basic-to-text or schema
caveats here without first confirming the disk's bytes are
reachable from the stock editor.

Procedure:

1. Note the disk path, the affected line number, and the
   divergent byte sequence (hex + decoded meaning).
2. In SimCoupé: insert the disk, `LOAD "FILE"`, `LIST <line>`,
   place the cursor on the line, and press **Enter** to re-submit
   it through TOKMAIN.
3. Save the disk back out, re-extract the file, inspect the byte
   at the affected position.
4. Apply the rule:
   - **If SimCoupé rewrites the bytes** to match what the
     lexer/grammar say (e.g. tokenises the identifier as a
     keyword), the disk's original bytes are **not reachable**
     from the editor → **bucket 4** → exclude the file from the
     corpus test (skip-list in `corpus_test.go` with a comment
     citing this SimCoupé result).
   - **If SimCoupé preserves the bytes**, the editor accepts the
     form → **bucket 3** → design a schema rule that matches the
     editor's actual mechanism, and add a section below.

This step is non-optional. Adding caveats or escape encodings for
content the stock editor wouldn't produce pollutes basic-to-text
and the schema with cases that have no real-world equivalent.

### 6. Apply the fix

- **Bucket 1 / 2**: write a focused failing unit test first if
  practical; fix in code; verify the unit test passes; then
  re-run the corpus.
- **Bucket 3**: add the schema rule in `sambasic/schema.go`, add
  the caveat entry under "Currently allowed variations" below
  with the SimCoupé citation, and verify the corpus result.
- **Bucket 4**: extend the skip-list in `corpus_test.go` with a
  one-line comment referencing the SimCoupé test.

### 7. Re-run the corpus and print the scoreboard

Go back to step 1. Repeat until everything passes or is documented.

---

## Currently allowed variations

### PROC-call placeholder trailing bytes (`SegProcCallPlaceholder`)

**What varies:** After a bare-identifier statement (a procedure
call), the editor stores a 6-byte invisible form
`0E FD FD <type> <addrLo> <addrHi>`. The leading 3 bytes
(`0E FD FD`) are stable; the trailing 3 bytes are rebuilt by
`LDPROG` (ROM L22699) on every LOAD via the `DOCOMP` pass — they
record the current program-store address of the target PROC, or
`0xFF` if unresolved.

**Why safe:** The corpus may store any of three states for these
bytes: pre-COMPILE (`FD ?? ??`), resolved (`0x80|page LSB MSB`),
or unresolved (`FF <stale> <stale>`). All three are functionally
equivalent post-LOAD because LDPROG rewrites them unconditionally
from current program layout. Our lexer emits the pre-COMPILE
placeholder (`FD 00 00`).

**Schema rule:** `Conform` checks the leading 4 bytes
(`<id> 0E FD FD`) exactly, then requires byte 4 to be `0x80|page`,
`0xFF`, or `0xFD`; bytes 5–6 are unconstrained.

**Reference:** `docs/sambasic-grammar.md` §6.5;
`sambasic/schema.go` (`SegProcCallPlaceholder` case).

### FN-call placeholder trailing bytes (`SegFnCallPlaceholder`)

Same as PROC except the marker is `0E FE FE` and the placeholder
pattern is `FE 00 00`. Same rationale and rule.

### Numeric FP encoding equivalence (`SegNumberFP`)

**What varies:** For an integer-valued literal in 0..65535, the
SAM 5-byte invisible-FP form may be the integer fast-path
(`00 00 LSB MSB 00`) or the general FP form (biased exponent +
normalised mantissa). Multiple encodings represent the same
value; the runtime treats them identically.

**Why safe:** `LISTING` / `PRINT` / arithmetic all decode the
5-byte form to the same scalar; the visible ASCII rendering
alongside is what the user sees.

**Schema rule:** `Conform` decodes both `want` and the schema's
recorded FP bytes; if the numeric values match (integers exactly,
floats within a tiny relative epsilon), the variation is accepted.

**Reference:** `docs/sambasic-grammar.md` §4.3;
`sambasic/schema.go` (`SegNumberFP` case via `decodeFP5`).
