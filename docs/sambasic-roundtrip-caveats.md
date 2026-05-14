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
