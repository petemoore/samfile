# DOS family `152b811ed65b651d` — masterdos-v2.3

Self-contained materialisation of one DOS family from the SAM
Coupé corpus. The family is the equivalence class of slot-0 DOS
bodies clustered at 1.5% byte-diff (see
`docs/dos-families.md` for the full table).

## Identity

- **Family-head SHA16:** `152b811ed65b651d`
- **Variants in family:** 2
- **Disks in family:** 6
- **Body length(s):** 15750
- **Load address(es):** 0x008000 (p3), 0x01c009 (p8)
- **Execution address(es):** (unset / 0xFF)

## Files in this directory

- `body.bin` — exact slot-0 body of the family-head SHA
  (`152b811ed65b651d`). Header-decoded, so byte 0 is the first byte the
  ROM would copy to the body's load address.
- `body.hex` — xxd-style hex dump of `body.bin` (big families only).
- `variants/*.md` — byte-diff summary of each variant against the
  head, including the differing byte ranges in hex so the variant
  body can be reconstructed from head + diff. Only written for
  families with at least 5 disks; extract any variant body with
  `tools/audit/extract_dos.py <sha>`.
- `src/` — commented original assembly source (only present when
  upstream source is known to assemble to a binary in this family).

## Source binding

- **Binary-of-record SHA16:** `152b811ed65b651d` (full SHA `152b811ed65b651df25e29f49e15340bec84ef3deebfba4eaa6cd76bfbb31fae`)
- **Upstream source:** `/Users/pmoore/git/masterdos` (src/*.asm)
- **Reference binary in upstream:** `/Users/pmoore/git/masterdos/res/MDOS23.bin`

### Notes from upstream README

> Source from https://github.com/dandoore/masterdos (Dan Doore).
> `src/masterdos23.asm` assembles to body.bin. The README
> notes v2.2 and v2.3 differ only at points labelled `Fix_*`
> in the source.

## Variants in this family

| Variant SHA16 | Disks | Length | Load | Exec | Within-fam diff vs head |
|---|---:|---:|---|---|---:|
| `152b811ed65b651d` ←source-of-record | 5 | 15750 | 0x008000 (p3) | — | head |
| `de2280bba1c6ba10` | 1 | 15750 | 0x01c009 (p8) | — | 0.673% |

## Sample disks

- CometAssembler1.8EdwinBlink
- FRED Magazine - Morkography (1992)
- Flight of Fantasy and Occult Connection Adventures (19xx)
- H-DOS V2.12 HD Loader V2.0 (1996)
- MasterDOS V2.1 (19xx)
- Recover-E (1995)
