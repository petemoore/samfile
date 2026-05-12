# DOS family `20e1c593dfd98cca`

Self-contained materialisation of one DOS family from the SAM
Coupé corpus. The family is the equivalence class of slot-0 DOS
bodies clustered at 1.5% byte-diff (see
`docs/dos-families.md` for the full table).

## Identity

- **Family-head SHA16:** `20e1c593dfd98cca`
- **Variants in family:** 3
- **Disks in family:** 24
- **Body length(s):** 15700
- **Load address(es):** 0x010000 (p5)
- **Execution address(es):** (unset / 0xFF)

## Files in this directory

- `body.bin` — exact slot-0 body of the family-head SHA
  (`20e1c593dfd98cca`). Header-decoded, so byte 0 is the first byte the
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

No upstream source identified for this family yet. The body
must be disassembled directly from `body.bin` if you need to
reason about its semantics. Disassemble with any z80
disassembler, e.g.

```
z80dasm body.bin -a -t -o 0x010000 > body.z80.s
```

If the family has multiple load addresses, pick the one that
covers the binary you care about. See `references/README.md`
for the SAM ROM v3.0 annotated disassembly, which is the
single biggest aid when reading these bodies.

## Variants in this family

| Variant SHA16 | Disks | Length | Load | Exec | Within-fam diff vs head |
|---|---:|---:|---|---|---:|
| `20e1c593dfd98cca` | 17 | 15700 | 0x010000 (p5) | — | head |
| `7e456418b0e18330` | 6 | 15700 | 0x010000 (p5) | — | 0.006% |
| `e80ef68803505adf` | 1 | 15700 | 0x010000 (p5) | — | 0.019% |

## Sample disks

- Metempsychosis Unreleased Demo - Wizard (19xx)
- Mouse Flash 1.1 for MDOS (19xx)
- Outwrite V2.0 (1992) (Chezron Software)
- PAX Disk 1 (1996) (Glenco)
- SC Filer V2.0 (1991) (Steve_s Software)
- Sam Adventure System Test Disk (1992) (Axxent Software)
- Sam D I C E V1.1 for MasterDOS (1991) (Kobrahsoft)
- Spectrum 128 Music Disk 2 (19xx) (PD)
- Spectrum Games Compilation 02 (1992)
- Spectrum Games Compilation 03 (1992)

... and 14 more.
