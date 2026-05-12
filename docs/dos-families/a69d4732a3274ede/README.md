# DOS family `a69d4732a3274ede`

Self-contained materialisation of one DOS family from the SAM
Coupé corpus. The family is the equivalence class of slot-0 DOS
bodies clustered at 1.5% byte-diff (see
`docs/dos-families.md` for the full table).

## Identity

- **Family-head SHA16:** `a69d4732a3274ede`
- **Variants in family:** 4
- **Disks in family:** 113
- **Body length(s):** 8078
- **Load address(es):** 0x038009 (p15), 0x078009 (p31)
- **Execution address(es):** (unset / 0xFF)

## Files in this directory

- `body.bin` — exact slot-0 body of the family-head SHA
  (`a69d4732a3274ede`). Header-decoded, so byte 0 is the first byte the
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
z80dasm body.bin -a -t -o 0x038009 > body.z80.s
```

If the family has multiple load addresses, pick the one that
covers the binary you care about. See `references/README.md`
for the SAM ROM v3.0 annotated disassembly, which is the
single biggest aid when reading these bodies.

## Variants in this family

| Variant SHA16 | Disks | Length | Load | Exec | Within-fam diff vs head |
|---|---:|---:|---|---|---:|
| `a69d4732a3274ede` | 107 | 8078 | 0x038009 (p15) | — | head |
| `fd6fa2869b471a1e` | 4 | 8078 | 0x078009 (p31) | — | 0.433% |
| `9a3a718414aa71d0` | 1 | 8078 | 0x038009 (p15) | — | 1.114% |
| `8a73776828631b6b` | 1 | 8078 | 0x038009 (p15) | — | 0.198% |

## Sample disks

- Blast Turbo_ by James R Curry (1995) (PD)
- COMMIX V2.00 by S. Grodkowski (1995) (PD)
- COMMIX V2.01 by S. Grodkowski (1995) (PD)
- COMMIX V2.02 by S. Grodkowski (1995) (PD)
- Easydisc V4.9 (1995) (Saturn Software)
- FRED Magazine Issue 05 (1990)
- FRED Magazine Issue 05 (1990) _a1_
- FRED Magazine Issue 06 (1990)
- FRED Magazine Issue 06 (1990) _a1_
- FRED Magazine Issue 07 (1990)

... and 103 more.
