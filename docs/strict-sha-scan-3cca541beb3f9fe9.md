# Strict-SHA scan: every rule fire on `3cca541beb3f9fe9` disks

Reference binary: `/Users/pmoore/git/samdos/res/samdos2.reference.bin` (sha256 `3cca541beb3f9fe93402a770997945b2be852e69f278d2b176ba0bbc4fbb6077`).

Each of the **30 disks** in the corpus has a slot-0
body whose sha256 matches the reference binary exactly. The
current branch's `samfile verify` was run against each. The
report below lists every rule that fired at least once on this
cohort, with the specific findings.

Because the cohort is one exact SAMDOS-2 build, any fire is
either:

- a real corruption / writer bug on the specific disk, or
- a rule documenting a writer convention that SAMDOS-2 itself
  does not enforce when SAVE'ing files (i.e. a false-positive
  rule that should be rewritten or removed).

## Summary

- **Cohort size:** 30 disks
- **Total events emitted:** 97550
- **Rules with at least one fire (fail):** 13

## Disks in cohort

- `32 Colour Demo by Gordon Wallis (1992) (PD)`
- `Allan Stevens - Home Utilities (1994)`
- `Allan Stevens - Home Utilities - Seven Pack (1994)`
- `Comms Loader (19xx)`
- `F-16 Combat Pilot Demo (1991) (PD)`
- `FRED Magazine Issue 13 (1991) _a1_`
- `FRED Magazine Issue 13 (1991)`
- `FRED Magazine Issue 14 (1991)`
- `Fredatives 3 (1992)`
- `GM-Calc V1.0 (19xx) (GM Software)`
- `GamesMaster 1.52 (1992) (Betasoft)`
- `MasterBASIC V1.7 (1994) (Betasoft)`
- `Metempsychosis Demo 2 (19xx)`
- `Metempsychosis Demo 4 (19xx)`
- `Metempsychosis Demo 5 (19xx)`
- `Metempsychosis Demo 6 (19xx)`
- `Metempsychosis Demo 7 (19xx)`
- `PAW Convertor by Martijn Groen (19xx) (PD)`
- `PacEmu by Simon Owen (2004) (PD)`
- `Sam Adventure System Demo Disk (1992) (Axxent Software)`
- `Sam Amateur Programming _ Electronics Issue 3 (Apr 1992)`
- `Sam Demo (19xx) (PD)`
- `Sam Demo Disk (1990) (Chris White)`
- `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)`
- `Sam Small C by Rumspft (1995) (Fred Publishing)`
- `SamCo Birthday Demos for 512K by Chris White (1991) (PD)`
- `SamCo Birthday Pack Games and Utils (1991) (Revelation)`
- `SupernovaUnfinished`
- `Surprise Demo from SAMCO News Disk 1 (1992) (PD)`
- `test`

## Rules that fired

### `BODY-MIRROR-AT-DIR-D3-DB` — 227/877 fails (25.9%), severity `cosmetic`

Source citation: `samdos/src/f.s:462-471`

**Distinct failure messages:**

- `dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0`
- `dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10`
- `dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12`
- `dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13`
- `dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `32 Colour Demo by Gordon Wallis (1992) (PD)` | `slot=0` | `samdos2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `32 Colour Demo by Gordon Wallis (1992) (PD)` | `slot=15` | `R.M.font` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `32 Colour Demo by Gordon Wallis (1992) (PD)` | `slot=16` | `scrollcode` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `32 Colour Demo by Gordon Wallis (1992) (PD)` | `slot=17` | `samfont` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Allan Stevens - Home Utilities (1994)` | `slot=0` | `CAPsoft` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Allan Stevens - Home Utilities (1994)` | `slot=15` | `NEWS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `Allan Stevens - Home Utilities (1994)` | `slot=3` | `BALL` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Allan Stevens - Home Utilities (1994)` | `slot=7` | `PAGE1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Allan Stevens - Home Utilities - Seven Pack (1994)` | `slot=0` | `CAPsoft` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Allan Stevens - Home Utilities - Seven Pack (1994)` | `slot=15` | `NEWS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `Allan Stevens - Home Utilities - Seven Pack (1994)` | `slot=3` | `BALL` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Allan Stevens - Home Utilities - Seven Pack (1994)` | `slot=7` | `PAGE1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Comms Loader (19xx)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Comms Loader (19xx)` | `slot=1` | `AUTOCOMMS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Comms Loader (19xx)` | `slot=2` | `ccode.cde` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `F-16 Combat Pilot Demo (1991) (PD)` | `slot=0` | `samdos2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `F-16 Combat Pilot Demo (1991) (PD)` | `slot=1` | `auto F16` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `F-16 Combat Pilot Demo (1991) (PD)` | `slot=2` | `title` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `F-16 Combat Pilot Demo (1991) (PD)` | `slot=3` | `gfx` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `F-16 Combat Pilot Demo (1991) (PD)` | `slot=4` | `code` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=10` | `CRACK THIS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 13 (1991)` | `slot=15` | `M-C CAPS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=16` | `DIS_MERGE` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `FRED Magazine Issue 13 (1991)` | `slot=17` | `AD` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=18` | `READPARROT` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `FRED Magazine Issue 13 (1991)` | `slot=20` | `BITBOBSCR` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=23` | `SIMON.SCR` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991)` | `slot=24` | `SIMON.TXT` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=27` | `M/C PT8` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `FRED Magazine Issue 13 (1991)` | `slot=28` | `M/C8` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=31` | `SCROLLSCR` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=32` | `SCROLLCODE` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=35` | `AD2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991)` | `slot=40` | `AD3` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991)` | `slot=41` | `AD1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991)` | `slot=42` | `AD4` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991)` | `slot=47` | `GauntPlay` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=48` | `Gauntlet` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=49` | `ZEB GREEN` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `FRED Magazine Issue 13 (1991)` | `slot=50` | `zeb code` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=51` | `splatsnap` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991)` | `slot=54` | `Banzai.col` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 13 (1991)` | `slot=55` | `Banzai.scr` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=56` | `Banzai.spr` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 13 (1991)` | `slot=57` | `Banzai.bar` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 13 (1991)` | `slot=58` | `Banzai.txt` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=60` | `reader` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991)` | `slot=62` | `# SAM` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=0` | ` FRED13 ` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=10` | `CRACK THIS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=15` | `M-C CAPS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=16` | `DIS_MERGE` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=17` | `AD` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=18` | `READPARROT` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=20` | `BITBOBSCR` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=23` | `SIMON.SCR` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=24` | `SIMON.TXT` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=27` | `M/C PT8` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=28` | `M/C8` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=31` | `SCROLLSCR` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=32` | `SCROLLCODE` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=35` | `AD2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=40` | `AD3` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=41` | `AD1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=42` | `AD4` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=47` | `GauntPlay` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=48` | `Gauntlet` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=49` | `ZEB GREEN` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=50` | `zeb code` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=51` | `splatsnap` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=54` | `Banzai.col` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=55` | `Banzai.scr` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=56` | `Banzai.spr` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=57` | `Banzai.bar` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=58` | `Banzai.txt` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=60` | `reader` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=62` | `# SAM` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 14 (1991)` | `slot=0` | ` FRED14 ` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `FRED Magazine Issue 14 (1991)` | `slot=1` | `DEMO2.MC` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `FRED Magazine Issue 14 (1991)` | `slot=12` | `Copper.pal` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `FRED Magazine Issue 14 (1991)` | `slot=15` | `islandcol` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 14 (1991)` | `slot=21` | `small.fnt` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 14 (1991)` | `slot=25` | `small2.fnt` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 14 (1991)` | `slot=26` | `Enigfont` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `FRED Magazine Issue 14 (1991)` | `slot=28` | `report.utl` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `FRED Magazine Issue 14 (1991)` | `slot=31` | `scr4` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 14 (1991)` | `slot=35` | `LETTERCODE` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 14 (1991)` | `slot=36` | `crunchy` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `FRED Magazine Issue 14 (1991)` | `slot=37` | `CrunchCode` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 14 (1991)` | `slot=38` | `Thick Font` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `FRED Magazine Issue 14 (1991)` | `slot=39` | `Text Code` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 14 (1991)` | `slot=40` | `Text FONT$` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 14 (1991)` | `slot=41` | `Logo code` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 14 (1991)` | `slot=44` | `MASTERTEQ` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `FRED Magazine Issue 14 (1991)` | `slot=47` | `P.MC` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 14 (1991)` | `slot=48` | `BITBOBSCR` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 14 (1991)` | `slot=56` | `BLOCKS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 14 (1991)` | `slot=57` | `MBLOCKS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 14 (1991)` | `slot=65` | `Mouse Code` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `FRED Magazine Issue 14 (1991)` | `slot=66` | `Arrow Data` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `FRED Magazine Issue 14 (1991)` | `slot=68` | `SpaceCode` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Fredatives 3 (1992)` | `slot=0` | `OS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Fredatives 3 (1992)` | `slot=32` | `Pointers` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Fredatives 3 (1992)` | `slot=37` | `Ras_Main A` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Fredatives 3 (1992)` | `slot=42` | `Lem_Main A` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Fredatives 3 (1992)` | `slot=56` | `Kaboom` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Fredatives 3 (1992)` | `slot=58` | `Mdriver2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Fredatives 3 (1992)` | `slot=6` | `Driver` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Fredatives 3 (1992)` | `slot=7` | `Initialise` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `GM-Calc V1.0 (19xx) (GM Software)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `GM-Calc V1.0 (19xx) (GM Software)` | `slot=3` | `dumpld` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `GamesMaster 1.52 (1992) (Betasoft)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `GamesMaster 1.52 (1992) (Betasoft)` | `slot=46` | `INITSDR` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `GamesMaster 1.52 (1992) (Betasoft)` | `slot=47` | `SDRIVER` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `GamesMaster 1.52 (1992) (Betasoft)` | `slot=51` | `mdriver` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 2 (19xx)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 4 (19xx)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 4 (19xx)` | `slot=1` | `ghost11` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 5 (19xx)` | `slot=0` | `SAMDOS2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 5 (19xx)` | `slot=2` | `SCR1` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 5 (19xx)` | `slot=27` | `TREK1` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 6 (19xx)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 6 (19xx)` | `slot=1` | `SCR1` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 6 (19xx)` | `slot=21` | `AUTO` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 7 (19xx)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 7 (19xx)` | `slot=1` | `PICCY1` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 7 (19xx)` | `slot=16` | `PICCY16` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Metempsychosis Demo 7 (19xx)` | `slot=30` | `AUTO` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `PAW Convertor by Martijn Groen (19xx) (PD)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `PAW Convertor by Martijn Groen (19xx) (PD)` | `slot=10` | `PAWOVR   H` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `PAW Convertor by Martijn Groen (19xx) (PD)` | `slot=11` | `Convertor` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `PAW Convertor by Martijn Groen (19xx) (PD)` | `slot=2` | `PAW 48k` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `PAW Convertor by Martijn Groen (19xx) (PD)` | `slot=20` | `code4` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `PacEmu by Simon Owen (2004) (PD)` | `slot=0` | `samdos2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `PacEmu by Simon Owen (2004) (PD)` | `slot=2` | `emulate` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `PacEmu by Simon Owen (2004) (PD)` | `slot=3` | `sprites` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `PacEmu by Simon Owen (2004) (PD)` | `slot=4` | `tiles` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `PacEmu by Simon Owen (2004) (PD)` | `slot=5` | `sound` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Adventure System Demo Disk (1992) (Axxent Software)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=1` | `AUTO` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=10` | `flash1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=11` | `flash2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=12` | `fontld` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=13` | `font` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=14` | `dumpld` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=15` | `UDG Design` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=16` | `btrans` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=17` | `btrans1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=18` | `emulator` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=19` | `skelt.bin` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=2` | `demo` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=20` | `modif.bin` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=21` | `snapt.bin` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=22` | `trans.bin` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=23` | `patch.bin` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=24` | `udg demo` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=25` | `test.udg` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=26` | `gdemo` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=27` | `ademo` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=28` | `idemo` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=29` | `dbt.bas` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=3` | `screen1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=30` | `DTRANS` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=31` | `mand1.scr` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=32` | `mand2.scr` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=33` | `con1.scr` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=34` | `con2.scr` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=35` | `wdemos.dem` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=36` | `LOAD.ME` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=4` | `screen2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=5` | `screen3` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=6` | `screen4` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=7` | `astro` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=8` | `king` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `Sam Demo Disk Issue 3 (Apr 1990) (SAM Computers LTD)` | `slot=9` | `flash` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=1` | `CC` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=4` | `cc.bin` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=1` | `BIRTHDAY` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=10` | `6` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=11` | `7` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=12` | `8` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=13` | `4` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=14` | `10` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=15` | `2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=16` | `POP_LOAD` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=17` | `IMP_LOAD` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=18` | `DEMO2.MC` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=21` | `splatsnap` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=22` | `splatsnap2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=24` | `SM 1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=25` | `SM 2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=3` | `LOAD SCR` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=4` | `3a` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=5` | `3` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=6` | `1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=7` | `11` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=8` | `9` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=9` | `5` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=0` | `samdos2` | dir byte 0xD2 (MGTFutureAndPast[0]) = 0x2a but should be 0 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=1` | `COS_LOAD` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=10` | `gothic2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=11` | `mapgraf2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=12` | `grafix3` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x12 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=13` | `grid3` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x14 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=14` | `CASSCRN` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=15` | `H` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=16` | `H.c` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=19` | `Cruncher` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=2` | `COSFONT` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=20` | `CrunchCode` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=3` | `COSMOS1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=5` | `VDU` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=6` | `FINISH` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=7` | `TITLE` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=8` | `CAS_LOAD` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=9` | `castle2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `SupernovaUnfinished` | `slot=0` | `samdos2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Surprise Demo from SAMCO News Disk 1 (1992) (PD)` | `slot=0` | `samdos2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Surprise Demo from SAMCO News Disk 1 (1992) (PD)` | `slot=1` | `AUTOSURP` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x10 |
| `Surprise Demo from SAMCO News Disk 1 (1992) (PD)` | `slot=3` | `SURPRISE1` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Surprise Demo from SAMCO News Disk 1 (1992) (PD)` | `slot=4` | `SURPRISE2` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Surprise Demo from SAMCO News Disk 1 (1992) (PD)` | `slot=5` | `SURPRISE3` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |
| `Surprise Demo from SAMCO News Disk 1 (1992) (PD)` | `slot=6` | `SURPRISE4` | dir byte 0xd3 (MGTFutureAndPast[1]) = 0x00 but body byte 0 = 0x13 |

_(1523 not-applicable events for this rule on this cohort)_

### `BASIC-VARS-GAP-INVARIANT` — 43/216 fails (19.9%), severity `cosmetic`

Source citation: `sam-basic-save-format.md`

**Distinct failure messages:**

- `BASIC SAVARS-NVARS = 1628; expected 604 for dialect samdos2`
- `BASIC SAVARS-NVARS = 2140; expected 604 for dialect unknown`
- `BASIC SAVARS-NVARS = 2156; expected 604 for dialect samdos2`
- `BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos`
- `BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `FRED Magazine Issue 13 (1991)` | `slot=49` | `ZEB GREEN` | BASIC SAVARS-NVARS = 2140; expected 604 for dialect unknown |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=49` | `ZEB GREEN` | BASIC SAVARS-NVARS = 2140; expected 604 for dialect unknown |
| `Fredatives 3 (1992)` | `slot=10` | `IDOS Setup` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=11` | `Iconcat` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=13` | `Wallow 1.0` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=14` | `Wallow 2.4` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=15` | `SoundProcs` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=2` | `AUTO Hippo` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=26` | `Multifont` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=4` | `Startup` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=55` | `Play Some` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=56` | `Kaboom` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=57` | `Kaboom.Bas` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=67` | `New Errors` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=68` | `Fun Errors` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=69` | `BBC Errors` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=70` | `Starstruck` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=71` | `Credits` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=72` | `CreditDemo` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=74` | `TVPP` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=76` | `New Bordow` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=77` | `Play Rasp` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `Fredatives 3 (1992)` | `slot=79` | `Clear All` | BASIC SAVARS-NVARS = 604; expected 2156 for dialect masterdos |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=10` | `sval` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=11` | `soundFX` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=12` | `sound1` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=16` | `put` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=17` | `put grab1` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=18` | `put grab2` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=20` | `spinner` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=23` | `blocks` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=24` | `tics` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=25` | `hide` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=26` | `exit` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=27` | `locn` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=30` | `README` | BASIC SAVARS-NVARS = 2156; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=4` | `alter` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=6` | `sort` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=7` | `delete` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=8` | `join` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `MasterBASIC V1.7 (1994) (Betasoft)` | `slot=9` | `inarray` | BASIC SAVARS-NVARS = 606; expected 604 for dialect samdos2 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=4` | `COSMISSONE` | BASIC SAVARS-NVARS = 1628; expected 604 for dialect samdos2 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=8` | `CAS_LOAD` | BASIC SAVARS-NVARS = 2156; expected 604 for dialect samdos2 |

_(2184 not-applicable events for this rule on this cohort)_

### `CHAIN-SECTOR-COUNT-MINIMAL` — 12/877 fails (1.4%), severity `cosmetic`

Source citation: `samfile.go:919`

**Distinct failure messages:**

- `file uses 11 sectors but 97 would suffice (bodyLen=49352)`
- `file uses 26 sectors but 97 would suffice (bodyLen=49352)`
- `file uses 27 sectors but 97 would suffice (bodyLen=49352)`
- `file uses 3 sectors but 97 would suffice (bodyLen=49352)`
- `file uses 4 sectors but 9 would suffice (bodyLen=4482)`
- `file uses 4 sectors but 97 would suffice (bodyLen=49352)`
- `file uses 5 sectors but 9 would suffice (bodyLen=4482)`
- `file uses 5 sectors but 97 would suffice (bodyLen=49352)`
- `file uses 6 sectors but 97 would suffice (bodyLen=49352)`
- `file uses 7 sectors but 97 would suffice (bodyLen=49352)`
- `file uses 8 sectors but 97 would suffice (bodyLen=49352)`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `Fredatives 3 (1992)` | `slot=16` | `Mist` | file uses 8 sectors but 97 would suffice (bodyLen=49352) |
| `Fredatives 3 (1992)` | `slot=17` | `Popcorn` | file uses 8 sectors but 97 would suffice (bodyLen=49352) |
| `Fredatives 3 (1992)` | `slot=18` | `Pipetune` | file uses 7 sectors but 97 would suffice (bodyLen=49352) |
| `Fredatives 3 (1992)` | `slot=19` | `House` | file uses 5 sectors but 97 would suffice (bodyLen=49352) |
| `Fredatives 3 (1992)` | `slot=20` | `Equinox` | file uses 26 sectors but 97 would suffice (bodyLen=49352) |
| `Fredatives 3 (1992)` | `slot=21` | `Pyracurse` | file uses 6 sectors but 97 would suffice (bodyLen=49352) |
| `Fredatives 3 (1992)` | `slot=22` | `Puzzletune` | file uses 3 sectors but 97 would suffice (bodyLen=49352) |
| `Fredatives 3 (1992)` | `slot=23` | `Strange` | file uses 11 sectors but 97 would suffice (bodyLen=49352) |
| `Fredatives 3 (1992)` | `slot=24` | `Dolphin` | file uses 4 sectors but 97 would suffice (bodyLen=49352) |
| `Fredatives 3 (1992)` | `slot=25` | `Panther` | file uses 27 sectors but 97 would suffice (bodyLen=49352) |
| `Fredatives 3 (1992)` | `slot=8` | `Disc III` | file uses 4 sectors but 9 would suffice (bodyLen=4482) |
| `Fredatives 3 (1992)` | `slot=9` | `Disc IV` | file uses 5 sectors but 9 would suffice (bodyLen=4482) |

_(1523 not-applicable events for this rule on this cohort)_

### `BODY-STARTPAGE-MATCHES-DIR` — 9/877 fails (1.0%), severity `cosmetic`

Source citation: `samdos/src/c.s:1376-1379`

**Distinct failure messages:**

- `body StartAddressPage = 0x01 but dir says 0x05`
- `body StartAddressPage = 0x01 but dir says 0x0a`
- `body StartAddressPage = 0x01 but dir says 0x61`
- `body StartAddressPage = 0x61 but dir says 0x60`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `FRED Magazine Issue 14 (1991)` | `slot=65` | `Mouse Code` | body StartAddressPage = 0x61 but dir says 0x60 |
| `Fredatives 3 (1992)` | `slot=33` | `Mdriver` | body StartAddressPage = 0x61 but dir says 0x60 |
| `GamesMaster 1.52 (1992) (Betasoft)` | `slot=51` | `mdriver` | body StartAddressPage = 0x61 but dir says 0x60 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=25` | `ECHO` | body StartAddressPage = 0x01 but dir says 0x61 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=26` | `ECHO    .C` | body StartAddressPage = 0x01 but dir says 0x0a |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=27` | `HELLO` | body StartAddressPage = 0x01 but dir says 0x61 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=28` | `HELLO   .C` | body StartAddressPage = 0x01 but dir says 0x0a |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=29` | `HILBERT` | body StartAddressPage = 0x01 but dir says 0x61 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=30` | `HILBERT .C` | body StartAddressPage = 0x01 but dir says 0x05 |

_(1523 not-applicable events for this rule on this cohort)_

### `BODY-PAGEOFFSET-MATCHES-DIR` — 9/877 fails (1.0%), severity `cosmetic`

Source citation: `samdos/src/c.s:1376-1379`

**Distinct failure messages:**

- `body PageOffset = 0x8000 but dir says 0x8f00`
- `body PageOffset = 0x8030 but dir says 0x8000`
- `body PageOffset = 0x8977 but dir says 0x8000`
- `body PageOffset = 0x8a1c but dir says 0x8000`
- `body PageOffset = 0x930d but dir says 0x8000`
- `body PageOffset = 0x9389 but dir says 0x8000`
- `body PageOffset = 0x9f3f but dir says 0x8000`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `FRED Magazine Issue 14 (1991)` | `slot=65` | `Mouse Code` | body PageOffset = 0x8000 but dir says 0x8f00 |
| `Fredatives 3 (1992)` | `slot=33` | `Mdriver` | body PageOffset = 0x8000 but dir says 0x8f00 |
| `GamesMaster 1.52 (1992) (Betasoft)` | `slot=51` | `mdriver` | body PageOffset = 0x8000 but dir says 0x8f00 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=25` | `ECHO` | body PageOffset = 0x8030 but dir says 0x8000 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=26` | `ECHO    .C` | body PageOffset = 0x8977 but dir says 0x8000 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=27` | `HELLO` | body PageOffset = 0x8a1c but dir says 0x8000 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=28` | `HELLO   .C` | body PageOffset = 0x930d but dir says 0x8000 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=29` | `HILBERT` | body PageOffset = 0x9389 but dir says 0x8000 |
| `Sam Small C by Rumspft (1995) (Fred Publishing)` | `slot=30` | `HILBERT .C` | body PageOffset = 0x9f3f but dir says 0x8000 |

_(1523 not-applicable events for this rule on this cohort)_

### `BASIC-LINE-NUMBER-BE` — 6/86 fails (7.0%), severity `structural`

Source citation: `sambasic/parse.go`

**Distinct failure messages:**

- `BASIC line number 0 out of range (1..65535)`
- `BASIC program parse failed: parse: line 48059 body extends past input`
- `BASIC program parse failed: parse: truncated numeric form at offset 241`
- `BASIC program parse failed: parse: truncated numeric form at offset 358`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `SamCo Birthday Demos for 512K by Chris White (1991) (PD)` | `slot=17` | `IMP_LOAD` | BASIC line number 0 out of range (1..65535) |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=1` | `COS_LOAD` | BASIC program parse failed: parse: truncated numeric form at offset 241 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=15` | `H` | BASIC line number 0 out of range (1..65535) |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=3` | `COSMOS1` | BASIC program parse failed: parse: truncated numeric form at offset 358 |
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=4` | `COSMISSONE` | BASIC program parse failed: parse: line 48059 body extends past input |
| `Surprise Demo from SAMCO News Disk 1 (1992) (PD)` | `slot=1` | `AUTOSURP` | BASIC line number 0 out of range (1..65535) |

_(1434 not-applicable events for this rule on this cohort)_

### `BOOT-OWNER-AT-T4S1` — 4/30 fails (13.3%), severity `structural`

Source citation: `rom-disasm:20473-20598`

**Distinct failure messages:**

- `no used slot has FirstSector (track 4, sector 1); disk is not bootable on real SAM hardware`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `FRED Magazine Issue 13 (1991)` | `disk` | `` | no used slot has FirstSector (track 4, sector 1); disk is not bootable on real SAM hardware |
| `Sam Amateur Programming _ Electronics Issue 3 (Apr 1992)` | `disk` | `` | no used slot has FirstSector (track 4, sector 1); disk is not bootable on real SAM hardware |
| `Sam Demo (19xx) (PD)` | `disk` | `` | no used slot has FirstSector (track 4, sector 1); disk is not bootable on real SAM hardware |
| `Sam Demo Disk (1990) (Chris White)` | `disk` | `` | no used slot has FirstSector (track 4, sector 1); disk is not bootable on real SAM hardware |

### `DIR-NAME-PADDING` — 4/877 fails (0.5%), severity `cosmetic`

Source citation: `sam-coupe_tech-man_v3-0.txt:4358-4359`

**Distinct failure messages:**

- `filename byte 0 is 0x11 (expected printable ASCII or 0x20 space)`
- `filename byte 0 is 0x7f (expected printable ASCII or 0x20 space)`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `32 Colour Demo by Gordon Wallis (1992) (PD)` | `slot=11` | `1.9.9.2.` | filename byte 0 is 0x7f (expected printable ASCII or 0x20 space) |
| `Allan Stevens - Home Utilities - Seven Pack (1994)` | `slot=0` | `CAPsoft` | filename byte 0 is 0x11 (expected printable ASCII or 0x20 space) |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=0` | ` FRED13 ` | filename byte 0 is 0x7f (expected printable ASCII or 0x20 space) |
| `FRED Magazine Issue 14 (1991)` | `slot=0` | ` FRED14 ` | filename byte 0 is 0x7f (expected printable ASCII or 0x20 space) |

_(1523 not-applicable events for this rule on this cohort)_

### `SCREEN-LENGTH-MATCHES-MODE` — 4/100 fails (4.0%), severity `structural`

Source citation: `sam-coupe_tech-man_v3-0.txt`

**Distinct failure messages:**

- `SCREEN mode 2 body length 24617 exceeds mode maximum 7424 bytes (min 6912 + 512 trailer slack)`
- `SCREEN mode 2 body length 24681 exceeds mode maximum 7424 bytes (min 6912 + 512 trailer slack)`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `Allan Stevens - Home Utilities (1994)` | `slot=9` | `RECORD1` | SCREEN mode 2 body length 24617 exceeds mode maximum 7424 bytes (min 6912 + 512 trailer slack) |
| `Allan Stevens - Home Utilities - Seven Pack (1994)` | `slot=9` | `RECORD1` | SCREEN mode 2 body length 24617 exceeds mode maximum 7424 bytes (min 6912 + 512 trailer slack) |
| `Fredatives 3 (1992)` | `slot=31` | `Mainscreen` | SCREEN mode 2 body length 24681 exceeds mode maximum 7424 bytes (min 6912 + 512 trailer slack) |
| `GM-Calc V1.0 (19xx) (GM Software)` | `slot=2` | `HL` | SCREEN mode 2 body length 24617 exceeds mode maximum 7424 bytes (min 6912 + 512 trailer slack) |

_(2300 not-applicable events for this rule on this cohort)_

### `DIR-SAM-WITHIN-CAPACITY` — 3/877 fails (0.3%), severity `inconsistency`

Source citation: `sam-coupe_tech-man_v3-0.txt:4405-4406`

**Distinct failure messages:**

- `SectorAddressMap[194]=0x7c has bits beyond bit 1559 set`
- `SectorAddressMap[194]=0x7f has bits beyond bit 1559 set`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `FRED Magazine Issue 13 (1991)` | `slot=44` | `FRED13` | SectorAddressMap[194]=0x7f has bits beyond bit 1559 set |
| `FRED Magazine Issue 13 (1991) _a1_` | `slot=44` | `FRED13` | SectorAddressMap[194]=0x7f has bits beyond bit 1559 set |
| `FRED Magazine Issue 14 (1991)` | `slot=60` | `reviews` | SectorAddressMap[194]=0x7c has bits beyond bit 1559 set |

_(1523 not-applicable events for this rule on this cohort)_

### `DISK-NOT-EMPTY` — 2/30 fails (6.7%), severity `inconsistency`

Source citation: `docs/disk-validity-rules.md`

**Distinct failure messages:**

- `disk has 0 occupied directory entries (all 80 slots are free)`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `Sam Demo (19xx) (PD)` | `disk` | `` | disk has 0 occupied directory entries (all 80 slots are free) |
| `Sam Demo Disk (1990) (Chris White)` | `disk` | `` | disk has 0 occupied directory entries (all 80 slots are free) |

### `CODE-LOAD-FITS-IN-MEMORY` — 1/393 fails (0.3%), severity `fatal`

Source citation: `samfile.go:802-804`

**Distinct failure messages:**

- `CODE load 0x7e000 + length 0x04000 = 0x82000 exceeds SAM's 512 KiB address space`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `FRED Magazine Issue 14 (1991)` | `slot=7` | `OPTION_TXT` | CODE load 0x7e000 + length 0x04000 = 0x82000 exceeds SAM's 512 KiB address space |

_(2007 not-applicable events for this rule on this cohort)_

### `BASIC-PROG-END-SENTINEL` — 1/86 fails (1.2%), severity `structural`

Source citation: `sambasic/file.go:36-42`

**Distinct failure messages:**

- `BASIC program does not end with 0xFF sentinel; body[76693] = 0x00`

**Every fire (disk, slot/ref, filename, message):**

| Disk | Ref | Filename | Message |
|---|---|---|---|
| `SamCo Birthday Pack Games and Utils (1991) (Revelation)` | `slot=4` | `COSMISSONE` | BASIC program does not end with 0xFF sentinel; body[76693] = 0x00 |

_(1434 not-applicable events for this rule on this cohort)_

## Rules that only passed (never fired)

(36 rules — listed for completeness.)

| Rule | passes | not-applicable |
|---|---:|---:|
| `ARRAY-FILETYPEINFO-TLBYTE-NAME` | 67 | 2333 |
| `BASIC-FILETYPEINFO-TRIPLETS` | 216 | 2184 |
| `BASIC-MGTFLAGS-20` | 216 | 2184 |
| `BASIC-STARTLINE-FF-DISABLES` | 216 | 2184 |
| `BASIC-STARTLINE-WITHIN-PROG` | 216 | 2184 |
| `BODY-BYTES-5-6-CANONICAL-FF` | 877 | 1523 |
| `BODY-EXEC-DIV16K-MATCHES-DIR` | 393 | 2007 |
| `BODY-EXEC-MOD16K-LO-MATCHES-DIR` | 393 | 2007 |
| `BODY-LENGTHMOD16K-MATCHES-DIR` | 877 | 1523 |
| `BODY-PAGE-LE-31` | 877 | 1523 |
| `BODY-PAGEOFFSET-8000H-FORM` | 877 | 1523 |
| `BODY-PAGES-MATCHES-DIR` | 877 | 1523 |
| `BODY-TYPE-MATCHES-DIR` | 877 | 1523 |
| `BOOT-ENTRY-POINT-AT-9` | 30 | 0 |
| `BOOT-SIGNATURE-AT-256` | 30 | 0 |
| `CHAIN-MATCHES-SAM` | 877 | 1523 |
| `CHAIN-NO-CYCLE` | 877 | 1523 |
| `CHAIN-TERMINATOR-ZERO-ZERO` | 877 | 1523 |
| `CODE-EXEC-WITHIN-LOADED-RANGE` | 393 | 2007 |
| `CODE-FILETYPEINFO-EMPTY` | 393 | 2007 |
| `CODE-LOAD-ABOVE-ROM` | 393 | 2007 |
| `COSMETIC-RESERVEDA-FF` | 877 | 1523 |
| `CROSS-DIRECTORY-AREA-UNUSED` | 30 | 0 |
| `CROSS-NO-DUPLICATE-NAMES` | 877 | 1523 |
| `CROSS-NO-SECTOR-OVERLAP` | 30 | 0 |
| `DIR-ERASED-IS-ZERO` | 877 | 1523 |
| `DIR-FIRST-SECTOR-VALID` | 877 | 1523 |
| `DIR-NAME-NOT-EMPTY` | 877 | 1523 |
| `DIR-SECTORS-MATCHES-CHAIN` | 877 | 1523 |
| `DIR-SECTORS-MATCHES-MAP` | 877 | 1523 |
| `DIR-SECTORS-NONZERO` | 877 | 1523 |
| `DIR-TYPE-BYTE-IS-KNOWN` | 877 | 1523 |
| `DISK-DIRECTORY-TRACKS` | 30 | 0 |
| `DISK-SECTOR-RANGE` | 30 | 0 |
| `DISK-TRACK-SIDE-ENCODING` | 30 | 0 |
| `SCREEN-MODE-AT-0xDD` | 100 | 2300 |

## Rules that were never applicable

(2 rules — never observed a subject they apply to.)

- `ZXSNAP-LENGTH-49152`
- `ZXSNAP-LOAD-ADDR-16384`

