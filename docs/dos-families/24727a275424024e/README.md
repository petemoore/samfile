# DOS family `24727a275424024e`

Self-contained materialisation of one DOS family from the SAM
Coupé corpus. The family is the equivalence class of slot-0 DOS
bodies clustered at 1.5% byte-diff (see
`docs/dos-families.md` for the full table).

## Identity

- **Family-head SHA16:** `24727a275424024e`
- **Variants in family:** 1
- **Disks in family:** 1
- **Body length(s):** 15700
- **Load address(es):** 0x010000 (p5)
- **Execution address(es):** (unset / 0xFF)

## Files in this directory

- `body.bin` — exact slot-0 body of the family-head SHA
  (`24727a275424024e`). Header-decoded, so byte 0 is the first byte the
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
| `24727a275424024e` | 1 | 15700 | 0x010000 (p5) | — | head |

## Sample disks

- Mike AJ Disc 6-Edwin (19xx)
