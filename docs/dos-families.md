# DOS families

Slot-0 DOS body variants clustered into families. Two variants
are in the same family iff they have the same body length and
their byte-wise diff is below **1.5%** of the body
length.

## Rationale

Three real causes of small-percentage variation between same-
DOS slot-0 bodies are all captured under one threshold:

1. **Per-magazine / per-disk launcher data.** Magazine-specific
   embedded auto-launch programs in the slot-0 data section.
   Pure data, not code. Typical diff: well under 1%.
2. **Memory-config rebase.** Same DOS code reassembled for a
   different SAM RAM page (e.g. page 14 vs page 30). Identical
   instructions, different page-selector constants and sysvar
   pointers. Typical diff: ~0.6%.
3. **Build / patch differences.** Bugfixes, branding, build
   stamps. Sub-percent typically.

A pair of bodies whose diff exceeds the threshold is *not* in
the same family at this threshold. Cross-length variants (e.g.
MasterDOS 15700 vs 15750) are kept as separate families — a
length change implies an inserted region whose semantics need
verifying before assuming compatibility.

## Summary

- **Threshold:** 1.5% byte-diff
- **Total bootable disks:** 554
- **Total unique slot-0 SHAs:** 179
- **Total families:** 45
- **Families with ≥ 5 disks:** 8
- **Families with 1 disk only (long-tail customs):** 31

## Table

| Rank | Family head SHA | Variants | Disks | Lengths | Load addresses (page) | Max within-variance | Sample disk |
|---:|---|---:|---:|---|---|---:|---|
| 1 | `9bc0fb4b949109e8` | 126 | 291 | 10000 | 0x008000 (p3), 0x008030 (p3), 0x00e0be (p4), 0x038009 (p15), 0x078009 (p31), 0x080009 (p33) | 3.00% | Curse of the Serpent_s Eye_ The (1994) (Dream Worl |
| 2 | `a69d4732a3274ede` | 4 | 113 | 8078 | 0x038009 (p15), 0x078009 (p31) | 1.20% | Blast Turbo_ by James R Curry (1995) (PD) |
| 3 | `13f6279c4d62e8be` | 3 | 31 | 15750 | 0x010000 (p5) | 0.63% | E-Tracker Program Disk V1.2 (19xx) (FRED Publishin |
| 4 | `78bc2964b7516db9` | 1 | 31 | 8077 | 0x078009 (p31) | 0.00% | Sam Adventure Club Issue 01 (Nov 1991) |
| 5 | `20e1c593dfd98cca` | 3 | 24 | 15700 | 0x010000 (p5) | 0.02% | Spectrum 128 Music Disk 2 (19xx) (PD) |
| 6 | `254ae17a87efb171` | 1 | 6 | 8077 | 0x078009 (p31) | 0.00% | Best of ENCELADUS_ Birthday Pack Edition (19xx) (R |
| 7 | `152b811ed65b651d` | 2 | 6 | 15750 | 0x008000 (p3), 0x01c009 (p8) | 0.67% | CometAssembler1.8EdwinBlink |
| 8 | `1b0b65f8a9545787` | 1 | 5 | 10157 | 0x008009 (p3) | 0.00% | Blitz Magazine Issue 2 (1997) (Persona) |
| 9 | `e7ead976f53c6003` | 1 | 4 | 8192 | 0x078009 (p31) | 0.00% | Mono Clipart Samples V1.0 (Nov 1995) (Steve_s Soft |
| 10 | `91e0f98622d2a6b0` | 1 | 3 | 32631 | 0x010000 (p5) | 0.00% | DS12 Duff Capers Music Demo (2003) (PD) |
| 11 | `21106301f8545821` | 1 | 3 | 8077 | 0x078009 (p31) | 0.00% | ENCELADUS Magazine Issue 09 (Feb 1992) (Relion Sof |
| 12 | `c76f8e68b0d0301b` | 1 | 2 | 15800 | 0x008009 (p3) | 0.00% | B-DOS V1.7N (1999) (Martijn Groen _ Edwin Blink) ( |
| 13 | `470699700014483a` | 2 | 2 | 9999 | 0x078009 (p31) | 1.11% | Metempsychosis Unreleased Demo - Demos4metemp (19x |
| 14 | `0f4d767f9db34845` | 1 | 2 | 8077 | 0x078009 (p31) | 0.00% | Sam Adventure Club Issue 09b (Mar 1993) _a1_ |
| 15 | `78843b6a4b894771` | 1 | 1 | 10191 | 0x008009 (p3) | 0.00% | B-DOS V1.5A (1997) (Martijn Groen _ Edwin Blink) ( |
| 16 | `7166b6af2054107e` | 1 | 1 | 14000 | 0x008009 (p3) | 0.00% | B-DOS V1.7D (1999) (Martijn Groen _ Edwin Blink) ( |
| 17 | `521478fd84761030` | 1 | 1 | 15800 | 0x008009 (p3) | 0.00% | B-DOS V1.7J (1999) (Martijn Groen _ Edwin Blink) ( |
| 18 | `f0047a502d0d54d9` | 1 | 1 | 501 | 0x00e0be (p4) | 0.00% | Banzai Pictures I by Dan Doore (1994) (PD) |
| 19 | `39f8558204cb3981` | 1 | 1 | 10000 | 0x008009 (p3) | 0.00% | Blinky Samples Disk 4 (1997) (Edwin Blink) |
| 20 | `587fa1d449e85ef3` | 1 | 1 | 67976 | 0x008000 (p3) | 0.00% | E-Tunes Player (19xx) (Andrew Collier) |
| 21 | `16b08ca76ac9bf6c` | 1 | 1 | 9000 | 0x078009 (p31) | 0.00% | ENCELADUS - Complete Guide to SAMBASIC Parts 1-7 ( |
| 22 | `571793c2f6a53f92` | 1 | 1 | 10000 | 0x078009 (p31) | 0.00% | ENCELADUS Magazine Issue 01 (Oct 1990) (Relion Sof |
| 23 | `50edb1b9a5308f85` | 1 | 1 | 10000 | 0x078009 (p31) | 0.00% | ENCELADUS Magazine Issue 02 (Dec 1990) (Relion Sof |
| 24 | `6a2f65a44273122f` | 1 | 1 | 8077 | 0x078009 (p31) | 0.00% | ENCELADUS Magazine Issue 03 (Feb 1991) (Relion Sof |
| 25 | `3d31391ff91d110b` | 1 | 1 | 8192 | 0x078009 (p31) | 0.00% | ENCELADUS Magazine Issue 04 (Apr 1991) (Relion Sof |
| 26 | `f9e25435a04c5542` | 1 | 1 | 8077 | 0x078009 (p31) | 0.00% | ENCELADUS Magazine Issue 10 (Apr 1992) (Relion Sof |
| 27 | `25b3b8c3de323fc8` | 1 | 1 | 32631 | 0x010000 (p5) | 0.00% | EXPLOSION - ZX SPECTRUM 48 Emulator _ COMMANDER (1 |
| 28 | `c98ea212d3f15722` | 1 | 1 | 10157 | 0x008009 (p3) | 0.00% | Entropy Demo (1992) (PD) _a2_ |
| 29 | `5a9d78bd06d11350` | 1 | 1 | 36957 | 0x008000 (p3) | 0.00% | FRED Magazine Issue 65 Menu (1995) |
| 30 | `dc5bc13f03508224` | 1 | 1 | 32631 | 0x010000 (p5) | 0.00% | MDOS _ MBASIC for Formatting Discs in 2 Drives (19 |
| 31 | `c3202ec6d71daf64` | 1 | 1 | 107317 | 0x008000 (p3) | 0.00% | MNEMOtech Demo 1 (19xx) (PD) |
| 32 | `f450085de2d9c53a` | 1 | 1 | 154354 | 0x008000 (p3) | 0.00% | MNEMOtech Demo 2 (19xx) (PD) |
| 33 | `487854350502cf42` | 1 | 1 | 14000 | 0x008009 (p3) | 0.00% | Megaboot V2.3 (Atom HD Interface) (1999) (M.Groen) |
| 34 | `16d35cdb1c766e7f` | 1 | 1 | 8976 | 0x078009 (p31) | 0.00% | Metempsychosis pdm12 (19xx) |
| 35 | `24727a275424024e` | 1 | 1 | 15700 | 0x010000 (p5) | 0.00% | Mike AJ Disc 6-Edwin (19xx) |
| 36 | `6e4c75fbba87c8ee` | 1 | 1 | 15800 | 0x008009 (p3) | 0.00% | Open 3D V082 by Tobermory (2001) (PD) |
| 37 | `08160038384ce831` | 1 | 1 | 32631 | 0x010000 (p5) | 0.00% | Ore Warz II (1990) (William McGugan) |
| 38 | `1ae0eda46245dfa8` | 1 | 1 | 15700 | 0x010000 (p5) | 0.00% | Sam D I C E V1.0 for MasterDOS (1991) (Kobrahsoft) |
| 39 | `68b90ca31c5f14e8` | 1 | 1 | 73567 | 0x008000 (p3) | 0.00% | Sam Mines (19xx) (PD) |
| 40 | `4dc74e1fc51f82bf` | 1 | 1 | 8100 | 0x007530 (p2) | 0.00% | Sam Prime (19xx) |
| 41 | `b3e1f498510fc710` | 1 | 1 | 9792 | 0x017d00 (p6) | 0.00% | Samsational Complete Guide to SAM PD Software (199 |
| 42 | `c1dc81fc3674eed2` | 1 | 1 | 10000 | 0x038009 (p15) | 0.00% | Spectrum Emulator (Sept 04) (1990) |
| 43 | `a3a7f8bf24d650ef` | 1 | 1 | 15700 | 0x010000 (p5) | 0.00% | TurboMON V1.0 (19xx) |
| 44 | `bec2a8d41401e03d` | 1 | 1 | 10000 | 0x038009 (p15) | 0.00% | Zenith Edition 1 (19xx) (Zenith Graphics) |
| 45 | `88a75d769da6a53f` | 1 | 1 | 11044 | 0x008009 (p3) | 0.00% | trinity |
