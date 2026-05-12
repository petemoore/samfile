# Per-rule coverage on strict-SHA cohort `3cca541beb3f9fe9`

This report restricts the audit pipeline's `checks` table to
the 30 disks in this cohort and recomputes
per-rule pass/fail rates. Compare the **cohort fail-rate**
column to the **all-corpus fail-rate**: if a rule was calibrated
against this DOS, the cohort fail-rate should be at or near 0%.

Rules where the cohort fail-rate is materially > 0% are either:

- genuinely buggy / over-strict on this DOS itself, or
- catching real corruption / writer bugs on the specific disks
  in this cohort.

With `--strict-sha`, the cohort contains only disks whose slot-0
body matches a single exact SHA — no clustering. This is the
falsifiable test for whether a 'convention' rule is enforced by
a single SAMDOS-2 build: if it fires in this cohort, the build
itself doesn't enforce it.

## Summary

- **Cohort:** strict-SHA cohort `3cca541beb3f9fe9`
- **Disks in cohort:** 30
- **Rules with ≥ 1 cohort fail:** 13
- **Rules with 0 cohort fails (clean):** 36

## Table (sorted by cohort fail-rate, desc)

| Rule | Severity | Cohort applies | Cohort fails | Cohort % | All-corpus % | Δ |
|---|---|---:|---:|---:|---:|---:|
| `BODY-MIRROR-AT-DIR-D3-DB` | cosmetic | 877 | 227 | 25.9% | 29.2% | -3.4pp |
| `BASIC-VARS-GAP-INVARIANT` | cosmetic | 216 | 43 | 19.9% | 13.6% | +6.3pp |
| `BOOT-OWNER-AT-T4S1` | structural | 30 | 4 | 13.3% | 10.6% | +2.7pp |
| `BASIC-LINE-NUMBER-BE` | structural | 216 | 25 | 11.6% | 9.8% | +1.8pp |
| `DISK-NOT-EMPTY` | inconsistency | 30 | 2 | 6.7% | 2.1% | +4.5pp |
| `SCREEN-LENGTH-MATCHES-MODE` | structural | 100 | 4 | 4.0% | 5.7% | -1.7pp |
| `CHAIN-SECTOR-COUNT-MINIMAL` | cosmetic | 877 | 12 | 1.4% | 1.7% | -0.4pp |
| `BODY-PAGEOFFSET-MATCHES-DIR` | cosmetic | 877 | 9 | 1.0% | 8.9% | -7.9pp |
| `BODY-STARTPAGE-MATCHES-DIR` | cosmetic | 877 | 9 | 1.0% | 8.4% | -7.4pp |
| `BASIC-PROG-END-SENTINEL` | structural | 216 | 1 | 0.5% | 5.4% | -5.0pp |
| `DIR-NAME-PADDING` | cosmetic | 877 | 4 | 0.5% | 1.6% | -1.1pp |
| `DIR-SAM-WITHIN-CAPACITY` | inconsistency | 877 | 3 | 0.3% | 1.3% | -1.0pp |
| `CODE-LOAD-FITS-IN-MEMORY` | fatal | 393 | 1 | 0.3% | 1.9% | -1.7pp |
| `ARRAY-FILETYPEINFO-TLBYTE-NAME` | structural | 67 | 0 | 0.0% | 0.1% | -0.1pp |
| `BASIC-FILETYPEINFO-TRIPLETS` | structural | 216 | 0 | 0.0% | 1.1% | -1.1pp |
| `BASIC-MGTFLAGS-20` | inconsistency | 216 | 0 | 0.0% | 1.3% | -1.3pp |
| `BASIC-STARTLINE-FF-DISABLES` | structural | 216 | 0 | 0.0% | 0.9% | -0.9pp |
| `BASIC-STARTLINE-WITHIN-PROG` | cosmetic | 216 | 0 | 0.0% | 0.1% | -0.1pp |
| `BODY-BYTES-5-6-CANONICAL-FF` | cosmetic | 877 | 0 | 0.0% | 0.1% | -0.1pp |
| `BODY-EXEC-DIV16K-MATCHES-DIR` | cosmetic | 393 | 0 | 0.0% | 11.0% | -11.0pp |
| `BODY-EXEC-MOD16K-LO-MATCHES-DIR` | cosmetic | 393 | 0 | 0.0% | 10.9% | -10.9pp |
| `BODY-LENGTHMOD16K-MATCHES-DIR` | cosmetic | 877 | 0 | 0.0% | 8.6% | -8.6pp |
| `BODY-PAGE-LE-31` | structural | 877 | 0 | 0.0% | 0.7% | -0.7pp |
| `BODY-PAGEOFFSET-8000H-FORM` | cosmetic | 877 | 0 | 0.0% | 3.7% | -3.7pp |
| `BODY-PAGES-MATCHES-DIR` | cosmetic | 877 | 0 | 0.0% | 6.7% | -6.7pp |
| `BODY-TYPE-MATCHES-DIR` | cosmetic | 877 | 0 | 0.0% | 8.5% | -8.5pp |
| `BOOT-ENTRY-POINT-AT-9` |  | 30 | 0 | 0.0% | 0.0% | +0.0pp |
| `BOOT-SIGNATURE-AT-256` |  | 30 | 0 | 0.0% | 0.0% | +0.0pp |
| `CHAIN-MATCHES-SAM` | inconsistency | 877 | 0 | 0.0% | 9.1% | -9.1pp |
| `CHAIN-NO-CYCLE` | structural | 877 | 0 | 0.0% | 0.0% | -0.0pp |
| `CHAIN-TERMINATOR-ZERO-ZERO` | structural | 877 | 0 | 0.0% | 1.0% | -1.0pp |
| `CODE-EXEC-WITHIN-LOADED-RANGE` | structural | 393 | 0 | 0.0% | 1.2% | -1.2pp |
| `CODE-FILETYPEINFO-EMPTY` | cosmetic | 393 | 0 | 0.0% | 1.3% | -1.3pp |
| `CODE-LOAD-ABOVE-ROM` |  | 393 | 0 | 0.0% | 0.0% | +0.0pp |
| `COSMETIC-RESERVEDA-FF` | cosmetic | 877 | 0 | 0.0% | 1.7% | -1.7pp |
| `CROSS-DIRECTORY-AREA-UNUSED` |  | 30 | 0 | 0.0% | 0.0% | +0.0pp |
| `CROSS-NO-DUPLICATE-NAMES` | inconsistency | 877 | 0 | 0.0% | 0.1% | -0.1pp |
| `CROSS-NO-SECTOR-OVERLAP` |  | 30 | 0 | 0.0% | 0.0% | +0.0pp |
| `DIR-ERASED-IS-ZERO` |  | 877 | 0 | 0.0% | 0.0% | +0.0pp |
| `DIR-FIRST-SECTOR-VALID` | fatal | 877 | 0 | 0.0% | 1.1% | -1.1pp |
| `DIR-NAME-NOT-EMPTY` | inconsistency | 877 | 0 | 0.0% | 0.0% | -0.0pp |
| `DIR-SECTORS-MATCHES-CHAIN` | structural | 877 | 0 | 0.0% | 8.8% | -8.8pp |
| `DIR-SECTORS-MATCHES-MAP` | structural | 877 | 0 | 0.0% | 1.2% | -1.2pp |
| `DIR-SECTORS-NONZERO` | structural | 877 | 0 | 0.0% | 0.0% | -0.0pp |
| `DIR-TYPE-BYTE-IS-KNOWN` |  | 877 | 0 | 0.0% | 0.0% | +0.0pp |
| `DISK-DIRECTORY-TRACKS` |  | 30 | 0 | 0.0% | 0.0% | +0.0pp |
| `DISK-SECTOR-RANGE` |  | 30 | 0 | 0.0% | 0.0% | +0.0pp |
| `DISK-TRACK-SIDE-ENCODING` |  | 30 | 0 | 0.0% | 0.0% | +0.0pp |
| `SCREEN-MODE-AT-0xDD` | structural | 100 | 0 | 0.0% | 3.4% | -3.4pp |
| `ZXSNAP-LENGTH-49152` |  | 0 | 0 | N/A | N/A | — |
| `ZXSNAP-LOAD-ADDR-16384` |  | 0 | 0 | N/A | N/A | — |
