# DOS family `78bc2964b7516db9`

Self-contained materialisation of one DOS family from the SAM
Coupé corpus. The family is the equivalence class of slot-0 DOS
bodies clustered at 1.5% byte-diff (see
`docs/dos-families.md` for the full table).

## Identity

- **Family-head SHA16:** `78bc2964b7516db9`
- **Variants in family:** 1
- **Disks in family:** 31
- **Body length(s):** 8077
- **Load address(es):** 0x078009 (p31)
- **Execution address(es):** (unset / 0xFF)

## Files in this directory

- `body.bin` — exact slot-0 body of the family-head SHA
  (`78bc2964b7516db9`). Header-decoded, so byte 0 is the first byte the
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
z80dasm body.bin -a -t -o 0x078009 > body.z80.s
```

If the family has multiple load addresses, pick the one that
covers the binary you care about. See `references/README.md`
for the SAM ROM v3.0 annotated disassembly, which is the
single biggest aid when reading these bodies.

## Variants in this family

| Variant SHA16 | Disks | Length | Load | Exec | Within-fam diff vs head |
|---|---:|---:|---|---|---:|
| `78bc2964b7516db9` | 31 | 8077 | 0x078009 (p31) | — | head |

## Sample disks

- Sam Adventure Club Issue 01 (Nov 1991)
- Sam Adventure Club Issue 02 (Jan 1992)
- Sam Adventure Club Issue 03 (Mar 1992)
- Sam Adventure Club Issue 03 (Mar 1992) _a1_
- Sam Adventure Club Issue 04 (May 1992)
- Sam Adventure Club Issue 04 (May 1992) _a1_
- Sam Adventure Club Issue 04 (May 1992) _a2_
- Sam Adventure Club Issue 05 (Aug 1992)
- Sam Adventure Club Issue 05 (Aug 1992) _a1_
- Sam Adventure Club Issue 06 (Sep 1992)

... and 21 more.
