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

## Ground-truth oracle: `llist-capture`

The SAM ROM's `LLIST` is the authoritative rendering of a SAM BASIC
program. We can capture it for any corpus file via:

```bash
~/git/sam-aarch64/tools/llist-capture.sh <source.mgt> <basic-name> [<output.txt>]
```

The tool builds a one-shot test disk that boots SAMDOS, loads the
corpus file with auto-RUN forced to a synthesised line at 65279
which runs `LLIST 1 TO 65278: CALL 16384`. The CALL transfers to
a 2-byte DI; HALT stub at the top-left of the screen file, which
makes SimCoupé exit cleanly via `-exitonhalt 1`. Output is captured
via SimCoupé's parallel-port-to-file mechanism (auto-named
`simc####.txt` in `~/Documents/SimCoupe/`).

**Decoupling the round-trip:** the captured LLIST output is the
ground truth for what `samfile basic-to-text` should emit. So
the round-trip can be split into two independent tests:

  - **`basic-to-text` correctness** — diff `samfile basic-to-text`
    against the LLIST capture. Tool: `llist-vs-b2t.sh`.
  - **`text-to-basic` correctness** — feed LLIST output to
    `samfile text-to-basic` and compare against the original
    body bytes. (Or use `samfile basic-to-text` output if the
    detok side is known-clean for this file.)

This lets us locate which side a divergence belongs to, instead
of always being uncertain whether a corpus-test failure is a
detok bug or a lexer bug.

### Known differences (not basic-to-text bugs)

When comparing LLIST output to `basic-to-text` output, expect these
formatting differences that do NOT indicate a basic-to-text bug:

  - **Current-line marker `>`** after the first line number in
    the LLIST output — driven by the EPPC sysvar (0x5C49). Now
    emitted by samfile in `--lossy` mode.
  - **Line wrapping at column ~80** for long lines in LLIST
    output — printer-width formatting that doesn't affect program
    semantics. Reproduced in `--lossy` mode.
  - **Stray `?` syntax-error marker** — LLIST inserts a flashing
    `?` at the byte whose memory address equals the XPTR sysvar
    (0x5AA3-0x5AA4). XPTR is persistent runtime state pointing to
    the last syntax-error location, set by whichever parser path
    the ROM was last running. Its value is *not* encoded in the
    saved BASIC body, so samfile cannot reproduce it from file
    bytes alone.

    **Fix in `tools/llist-capture/main.go`**: the auto-RUN line
    now does `POKE 23203,0: POKE 23204,0:` before `LLIST`, which
    zeros XPTR. Subsequent parser activity may store XPTR into
    line 65279 (our control line), but LLIST's `1 TO 65278` range
    excludes that, so no `?` artefact appears in captures.

  - **Line 0 / lines > 65278 excluded** — the harness invokes
    `LLIST 1 TO 65278`, not `LLIST 0 TO 65278`. The 0-lower-bound
    form triggers a ROM slow-path that makes LLIST runtime jump
    from <15s to 60+s on non-trivial programs, breaking the
    sweep's timing budget. samfile's lossy mode mirror-skips
    lines outside `[1, 65278]` so byte parity is preserved at the
    cost of not testing the 329 corpus files' line-0 content.

Anything else is worth investigating.

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
