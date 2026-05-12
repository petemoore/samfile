# Per-rule coverage on DOS family `9bc0fb4b949109e8`

This report restricts the audit pipeline's `checks` table to
the 291 disks in this family and recomputes
per-rule pass/fail rates. Compare the **family fail-rate**
column to the **all-corpus fail-rate**: if a rule was calibrated
against this DOS, the family fail-rate should be at or near 0%.

Rules where the family fail-rate is materially > 0% are either:

- genuinely buggy / over-strict on this DOS itself, or
- catching real corruption / writer bugs on the specific disks
  in this family.

## Summary

- **Family head SHA:** `9bc0fb4b949109e8`
- **Disks in family:** 291
- **Membership threshold:** 1.5% byte-diff
- **Rules with ≥ 1 family fail:** 34
- **Rules with 0 family fails (clean):** 15

## Table (sorted by family fail-rate, desc)

| Rule | Severity | Family applies | Family fails | Family % | All-corpus % | Δ |
|---|---|---:|---:|---:|---:|---:|
| `BODY-MIRROR-AT-DIR-D3-DB` | cosmetic | 11228 | 3220 | 28.7% | 29.2% | -0.6pp |
| `BOOT-OWNER-AT-T4S1` | structural | 291 | 34 | 11.7% | 10.6% | +1.1pp |
| `BASIC-LINE-NUMBER-BE` | structural | 3209 | 153 | 4.8% | 9.8% | -5.0pp |
| `BASIC-VARS-GAP-INVARIANT` | cosmetic | 3209 | 137 | 4.3% | 13.6% | -9.3pp |
| `SCREEN-LENGTH-MATCHES-MODE` | structural | 1281 | 48 | 3.7% | 5.7% | -1.9pp |
| `SCREEN-MODE-AT-0xDD` | structural | 1281 | 26 | 2.0% | 3.4% | -1.4pp |
| `BODY-STARTPAGE-MATCHES-DIR` | cosmetic | 11228 | 207 | 1.8% | 8.4% | -6.6pp |
| `CODE-LOAD-FITS-IN-MEMORY` | fatal | 5344 | 91 | 1.7% | 1.9% | -0.2pp |
| `DISK-NOT-EMPTY` | inconsistency | 291 | 4 | 1.4% | 2.1% | -0.8pp |
| `BODY-PAGEOFFSET-MATCHES-DIR` | cosmetic | 11228 | 151 | 1.3% | 8.9% | -7.6pp |
| `CHAIN-SECTOR-COUNT-MINIMAL` | cosmetic | 11228 | 134 | 1.2% | 1.7% | -0.6pp |
| `BODY-EXEC-DIV16K-MATCHES-DIR` | cosmetic | 5344 | 63 | 1.2% | 11.0% | -9.9pp |
| `BODY-EXEC-MOD16K-LO-MATCHES-DIR` | cosmetic | 5344 | 61 | 1.1% | 10.9% | -9.8pp |
| `COSMETIC-RESERVEDA-FF` | cosmetic | 11228 | 113 | 1.0% | 1.7% | -0.7pp |
| `BODY-LENGTHMOD16K-MATCHES-DIR` | cosmetic | 11228 | 104 | 0.9% | 8.6% | -7.6pp |
| `BODY-TYPE-MATCHES-DIR` | cosmetic | 11228 | 104 | 0.9% | 8.5% | -7.6pp |
| `BODY-PAGES-MATCHES-DIR` | cosmetic | 11228 | 99 | 0.9% | 6.7% | -5.9pp |
| `DIR-SAM-WITHIN-CAPACITY` | inconsistency | 11228 | 90 | 0.8% | 1.3% | -0.5pp |
| `BODY-PAGE-LE-31` | structural | 11228 | 85 | 0.8% | 0.7% | +0.1pp |
| `DIR-NAME-PADDING` | cosmetic | 11228 | 80 | 0.7% | 1.6% | -0.9pp |
| `BASIC-PROG-END-SENTINEL` | structural | 3209 | 22 | 0.7% | 5.4% | -4.7pp |
| `CODE-EXEC-WITHIN-LOADED-RANGE` | structural | 5344 | 35 | 0.7% | 1.2% | -0.6pp |
| `CHAIN-MATCHES-SAM` | inconsistency | 11228 | 54 | 0.5% | 9.1% | -8.6pp |
| `CODE-FILETYPEINFO-EMPTY` | cosmetic | 5344 | 23 | 0.4% | 1.3% | -0.8pp |
| `DIR-SECTORS-MATCHES-MAP` | structural | 11228 | 48 | 0.4% | 1.2% | -0.8pp |
| `CROSS-NO-DUPLICATE-NAMES` | inconsistency | 11228 | 26 | 0.2% | 0.1% | +0.1pp |
| `DIR-SECTORS-MATCHES-CHAIN` | structural | 11228 | 26 | 0.2% | 8.8% | -8.5pp |
| `CHAIN-TERMINATOR-ZERO-ZERO` | structural | 11228 | 21 | 0.2% | 1.0% | -0.8pp |
| `DIR-FIRST-SECTOR-VALID` | fatal | 11228 | 20 | 0.2% | 1.1% | -0.9pp |
| `BODY-PAGEOFFSET-8000H-FORM` | cosmetic | 11228 | 12 | 0.1% | 3.7% | -3.6pp |
| `BASIC-STARTLINE-WITHIN-PROG` | cosmetic | 3209 | 1 | 0.0% | 0.1% | -0.1pp |
| `DIR-NAME-NOT-EMPTY` | inconsistency | 11228 | 2 | 0.0% | 0.0% | -0.0pp |
| `BODY-BYTES-5-6-CANONICAL-FF` | cosmetic | 11228 | 1 | 0.0% | 0.1% | -0.1pp |
| `CHAIN-NO-CYCLE` | structural | 11228 | 1 | 0.0% | 0.0% | +0.0pp |
| `ARRAY-FILETYPEINFO-TLBYTE-NAME` | structural | 641 | 0 | 0.0% | 0.1% | -0.1pp |
| `BASIC-FILETYPEINFO-TRIPLETS` | structural | 3209 | 0 | 0.0% | 1.1% | -1.1pp |
| `BASIC-MGTFLAGS-20` | inconsistency | 3209 | 0 | 0.0% | 1.3% | -1.3pp |
| `BASIC-STARTLINE-FF-DISABLES` | structural | 3209 | 0 | 0.0% | 0.9% | -0.9pp |
| `BOOT-ENTRY-POINT-AT-9` |  | 291 | 0 | 0.0% | 0.0% | +0.0pp |
| `BOOT-SIGNATURE-AT-256` |  | 291 | 0 | 0.0% | 0.0% | +0.0pp |
| `CODE-LOAD-ABOVE-ROM` |  | 5344 | 0 | 0.0% | 0.0% | +0.0pp |
| `CROSS-DIRECTORY-AREA-UNUSED` |  | 291 | 0 | 0.0% | 0.0% | +0.0pp |
| `CROSS-NO-SECTOR-OVERLAP` |  | 291 | 0 | 0.0% | 0.0% | +0.0pp |
| `DIR-ERASED-IS-ZERO` |  | 11228 | 0 | 0.0% | 0.0% | +0.0pp |
| `DIR-SECTORS-NONZERO` | structural | 11228 | 0 | 0.0% | 0.0% | -0.0pp |
| `DIR-TYPE-BYTE-IS-KNOWN` |  | 11228 | 0 | 0.0% | 0.0% | +0.0pp |
| `DISK-DIRECTORY-TRACKS` |  | 291 | 0 | 0.0% | 0.0% | +0.0pp |
| `DISK-SECTOR-RANGE` |  | 291 | 0 | 0.0% | 0.0% | +0.0pp |
| `DISK-TRACK-SIDE-ENCODING` |  | 291 | 0 | 0.0% | 0.0% | +0.0pp |
| `ZXSNAP-LENGTH-49152` |  | 0 | 0 | N/A | N/A | — |
| `ZXSNAP-LOAD-ADDR-16384` |  | 0 | 0 | N/A | N/A | — |
