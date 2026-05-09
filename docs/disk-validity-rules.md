# SAM Coupé MGT disk-image validity rules

Catalog of validity rules for SAMDOS .mgt floppy images, drawn from the
SAMDOS source (`~/git/samdos/src/*.s`), the SAM Coupé Technical Manual
v3.0 and ROM v3.0 annotated disassembly, samfile's existing parser, and
sam-aarch64's accumulated forensic notes.

Intended consumer: a future `samfile verify` sub-command. Source phase
only — corpus validation comes later.

## Per-rule schema

```
### <RULE_ID> — <one-line summary>

- What:               constraint statement
- Severity:           fatal | structural | inconsistency | cosmetic
- Source authority:   SAMDOS-code | ROM | Tech-Manual | samfile-implicit |
                      empirical-convention
- Citation:           file:line + a verbatim snippet for the strongest
                      evidence
- Dialect:            SAMDOS-2 | MasterDOS | SAMDOS-1 | all
- Suppressed by:      (optional) HIDDEN / PROTECTED / ERASED
- Test sketch:        how to check
- Open questions:     (optional)
```

Severity legend:
- `fatal`         — image will not boot / SAMDOS will reject or corrupt.
- `structural`    — disk-walk invariant; violation produces undefined
                    behaviour or sector reuse.
- `inconsistency` — two views of the same fact disagree; usually
                    cosmetic at runtime but indicates a buggy writer.
- `cosmetic`      — Tech-Manual convention; spec-compliant either way.

Source-authority legend:
- `SAMDOS-code`         — a code path in `~/git/samdos/src/*.s` reads or
                          writes the byte.
- `ROM`                 — a code path in the SAM ROM v3.0 reads or
                          writes the byte.
- `Tech-Manual`         — only the SAM Coupé Tech Manual mentions it; no
                          code citation found. Possibly aspirational
                          (cf. hook 128 / BTHK auto-RUN, samdos2-auto-
                          run-analysis.md).
- `samfile-implicit`    — samfile's parser or writer relies on it; not
                          necessarily a SAMDOS or ROM constraint.
- `empirical-convention`— observed in real-world disks; no formal
                          authority documents it.

---

## 0. Two PR-12 hypotheses, verified up front

### MGTFutureAndPast mirror of body header — CONFIRMED

PR-12 hypothesised that `MGTFutureAndPast[0]` (dir byte `0xD2`) is
reserved and stays zero, and `MGTFutureAndPast[1..9]` (dir bytes
`0xD3..0xDB`) mirror body-header bytes 0..8. The SAMDOS source confirms
this:

- `samdos/src/f.s:462-471` (`svhd`):

  ```asm
  svhd:  ld hl,hd001        ; SAMDOS in-RAM 9-byte header copy
         ld de,fsa+211      ; offset 211 (= 0xD3) of dir-entry buffer
         ld b,9             ; nine bytes
  svhd1: ld a,(hl)
         ld (de),a
         call sbyt          ; ALSO write to file body via sector chain
         inc hl
         inc de
         djnz svhd1
         ret
  ```

  So SAVE writes the same 9 bytes to dir offset 0xD3 AND to the file
  body's leading 9 bytes. `fsa+211 = 0xD3` because `fsa` is the start
  of the 256-byte directory-entry RAM buffer in SAMDOS workspace.

- `samdos/src/c.s:1376-1379` (`gtfle`):

  ```asm
  ld (ix+rptl),211        ; offset 211 (= 0xD3) inside the dir buffer
  call grpnt
  ld bc,9
  ldir                    ; → into hd001..page1 (the in-RAM 9-byte cache)
  ```

  So LOAD reads the same 9 bytes back from dir offset 0xD3 into
  SAMDOS's 9-byte header cache (`hd001..page1` defined at
  `samdos/src/b.s:255-260`).

- Byte 0xD2 is never read or written by SAMDOS. It is genuinely
  reserved.

Conclusion: `MGTFutureAndPast[0]==0` and `MGTFutureAndPast[1..9] ==
body[0..8]` are real invariants. Tech Manual L4366-4368 ("MGT FUTURE
AND PAST ... not used by SAMDOS") is a documentation error. See
`sam-disk-format.md` §2.4 and `sam-file-header.md` §2.

### ExecutionAddressDiv16K + ExecutionAddressMod16K mirror at body-header bytes 5-6 — CONFIRMED, with a wrinkle

PR-12 hypothesised that body-header byte 5 mirrors
`ExecutionAddressDiv16K` and body-header byte 6 mirrors the low byte
of `ExecutionAddressMod16K`. ROM's `LOAD CODE` auto-exec gate at
`E281-E294` (`rom-disasm:22467-22484`) confirms the byte 5 part:

```
E281 3A254B   LD A,(HDR+HDN+6)   ;byte 37 = REQUESTED exec page (from
                                 ; SAVE's HDR — i.e. dir byte 0xF2)
E287 FEFF     CP 0FFH
E289 2009     JR NZ,HDLDEX       ;JR IF LOAD EXEC OVER-RIDES SAVED EXEC
E28B 3A754B   LD A,(HDL+HDN+6)   ;byte 37 of HDL = LOADED exec page
                                 ; (populated from body header byte 5
                                 ; via SAMDOS's dschd → hd001..)
E28E 2A764B   LD HL,(HDL+HDN+7)
E291 FEFF     CP 0FFH
E293 C8       RET Z              ;RET IF NO EXEC ADDR
E294 CD7912   HDLDEX: CALL PDPSR2 ;page form decode and JP
```

The gate has two paths:

1. Requested-exec path (`HDR+HDN+6`): SAVE-time / `LOAD CODE n,...,exec`
   override populated from dir byte 0xF2. If non-FF, take HDLDEX
   immediately.
2. Loaded-exec path (`HDL+HDN+6`): from the body-header byte 5
   (SAMDOS's `hd001..` cache → HDL via dschd / hconr). If non-FF, take
   HDLDEX; if FF, RET (no auto-exec).

So the auto-exec is gated on BOTH the directory entry's byte 0xF2 AND
the body-header's byte 5. Each must be `0xFF` to disable auto-exec.

The "low byte of ExecutionAddressMod16K mirrors at body-header byte 6"
part of PR-12 prose is half-right: body-header bytes 5-6 are the
LOADED exec-page + low-half of the LOADED exec-address, and HDL+HDN+7
is a 16-bit LE pair starting at body header byte 6. Body header byte
6 is therefore a copy of the low byte of `ExecutionAddressMod16K` (or
the low byte of HDR+HDN+7 in HDL/HDR layout terms). But the high byte
of `ExecutionAddressMod16K` lives only in the directory entry at
0xF4 — there is no body-header byte for it. The body header is
9 bytes long and stops at byte 8. So the mirror is *partial*:

| Field                          | Body header byte | Dir entry byte |
|--------------------------------|------------------|----------------|
| `ExecutionAddressDiv16K`       | 5                | 0xF2           |
| `ExecutionAddressMod16K_lo`    | 6                | 0xF3           |
| `ExecutionAddressMod16K_hi`    | — (not stored)   | 0xF4           |

This is congruent with samfile's `FileHeader.ExecutionAddressMod16KLo`
field (`samfile.go:194`) and the dir-entry's full 16-bit
`ExecutionAddressMod16K` (`samfile.go:117-118`).

Conclusion: PR-12 is right about byte 5; byte 6 mirrors only the low
byte. A `verify` rule should check the byte-5 equality and the
byte-6-equals-low-byte-of-0xF3..0xF4 equality. See
`test-mgt-byte-layout.md` for the boot-disk byte-by-byte breakdown.

---

## 1. Disk-level rules

### DISK-IMAGE-SIZE — Image is exactly 819,200 bytes

- What: An .mgt image is 80 cylinders × 2 sides × 10 sectors × 512
  bytes = 819,200 bytes. Files shorter than this are zero-padded by
  `samfile Load`; files longer than this are truncated.
- Severity: structural
- Source authority: Tech-Manual + samfile-implicit
- Citation: `sam-coupe_tech-man_v3-0.txt:4262-4275`:

  ```
  formatted as double sided, 80 track per side, 10 sectors per track,
  ... 1560 data sectors of 512 bytes (798720 bytes).
  ```

  samfile: `samfile.go:71` `DiskImage [819200]byte` and
  `samfile.go:362-371` (`Load`) zero-pads / truncates.
- Dialect: all
- Test sketch: assert `len(image) == 819200` before parsing.

### DISK-NOT-EDSK — Reject EDSK-format images

- What: Files prefixed with `"EXTENDED CPC DSK File"` are EDSK, not
  MGT; they must be converted (samdisk) before samfile can read them.
- Severity: fatal
- Source authority: samfile-implicit
- Citation: `samfile.go:355-368`:

  ```go
  var edskMagic = []byte("EXTENDED CPC DSK File")
  ...
  return nil, fmt.Errorf("error: EDSK format not supported; ...")
  ```
- Dialect: all
- Test sketch: check first 21 bytes don't match `EDSK` magic.

### DISK-LAYOUT-CYL-INTERLEAVED — Sectors are cylinder-interleaved

- What: Image byte offset for (track, sector) is
  `((track & 0x80) >> 7) * 5120 + (sector - 1) * 512 + (track & 0x7f) * 10240`.
  Side 0 of each cylinder precedes side 1.
- Severity: structural (defines how every other rule reads bytes)
- Source authority: samfile-implicit + empirical-convention
- Citation: `samfile.go:993-995`:

  ```go
  func (sector *Sector) Offset() int {
      return int(sector.Track>>7)*5120 + (int(sector.Sector)-1)*512 + int(sector.Track&0x7f)*10240
  }
  ```

  SimCoupé Base/Disk.cpp:164 uses the same formula. Tech Manual is
  silent on file-image layout; this is the de-facto MGT convention.
- Dialect: all
- Test sketch: not a rule to enforce; rather a precondition for the
  reader. Cite when documenting other rules.
- Open questions: Tech Manual does not formally specify the image-
  file layout (only the physical disk); cylinder-interleave is
  SimCoupé/samfile convention. No "official" .mgt spec exists.

### DISK-DIRECTORY-TRACKS — Tracks 0-3 (side 0) are reserved for the directory

- What: The 4 directory tracks (T0..T3 of side 0) hold 80 directory
  slots packed two per sector. No file's first sector may land in this
  region. Bytes 510-511 of these sectors are NOT chain links — they
  fall in slot 2k+1's `ReservedB` (dir bytes 0xFE-0xFF).
- Severity: structural
- Source authority: Tech-Manual
- Citation: `sam-coupe_tech-man_v3-0.txt:4340-4343`:

  ```
  The first 4 tracks of the disk are allocated to the disk directory,
  starting at track 0, sector 1. These 4 tracks give us 40 sectors
  each split into two 256 bytes entries.
  ```
- Dialect: all
- Test sketch: assert every `FileEntry.FirstSector.Track` is in
  {4..79, 128..207}.

### DISK-TRACK-SIDE-ENCODING — Valid track byte ranges

- What: Track-byte bit 7 encodes side; valid values are `0x00..0x4F`
  (side 0 cylinders 0..79) and `0x80..0xCF` (side 1 cylinders 0..79).
  Values `0x50..0x7F` and `0xD0..0xFF` are invalid.
- Severity: fatal
- Source authority: samfile-implicit + empirical-convention
- Citation: `samfile.go:393-394`:

  ```go
  if (sector.Track >= 80 && sector.Track < 128) || sector.Track >= 208 {
      return nil, fmt.Errorf("track out of range: %v ...", sector.Track)
  }
  ```
- Dialect: all
- Test sketch: every track byte encountered (directory entry first-
  sector, sector chain links, any other position) is in {0..79,
  128..207}.

### DISK-SECTOR-RANGE — Sector numbers are 1..10

- What: Sector numbers are 1-based and run 1..10. Zero is the chain
  terminator (paired with track 0) but is never a valid sector
  reference for a live data sector.
- Severity: fatal
- Source authority: samfile-implicit
- Citation: `samfile.go:389-392`:

  ```go
  if sector.Sector < 1 || sector.Sector > 10 {
      ...
      return nil, fmt.Errorf("sector out of range: %v", sector.Sector)
  }
  ```
- Dialect: all
- Test sketch: every sector byte encountered in a live link is 1..10.
  A chain-terminator (0, 0) pair is allowed at the tail of a file.

---

## 2. Directory-entry rules

### DIR-SLOT-COUNT — Exactly 80 slots at fixed offsets

- What: The directory holds 80 256-byte slots at image offsets
  `slot * 256` for `slot ∈ {0..79}`. There is no growth; max 80
  files.
- Severity: structural
- Source authority: Tech-Manual + samfile-implicit
- Citation: Tech Manual L4340-4343; `samfile.go:438-446`
  (`DiskJournal` walks tracks 0..3, sectors 1..10, two slots per
  sector).
- Dialect: all
- Test sketch: parse 80 slots from `image[slot*256 : slot*256+256]`.

### DIR-TYPE-BYTE-IS-KNOWN — Type byte's low 5 bits are a documented type

- What: After masking off HIDDEN (bit 7) and PROTECTED (bit 6), the
  remaining `type & 0x3F` must be either 0 (erased) or one of the
  documented type values. SAMDOS's `drtab` (`samdos/src/e.s:322-355`)
  recognises 1..12 and 16..20 for DIR display; the SAM ROM
  `E019`-block (rom-disasm:22032) lists 16, 17, 18, 19, 20 as the
  public file types and the Tech Manual L4304-4314 lists 5 (SNP),
  16, 17, 18, 19, 20.
- Severity: inconsistency
- Source authority: SAMDOS-code + ROM + Tech-Manual
- Citation: `samdos/src/e.s:322-355` (SAMDOS DIR-output type table),
  `rom-disasm:22032`, Tech Manual L4304-4314.
- Dialect: all
- Test sketch: after masking flags, type is in {0, 5, 16, 17, 18,
  19, 20}; warn (don't error) on types 1-4, 6-12 (ZX-compat /
  MasterDOS / SAMDOS-1) and unknown values.
- Open questions: types 1-12 in SAMDOS's DIR table correspond to ZX
  / Plus-D legacy file types. Real-world MGT disks rarely use them.
  Type 6 ("MD.FILE") is MasterDOS; type 8 ("SPECIAL") is opaque. A
  verify rule should probably warn rather than fail on these.

### DIR-TYPE-MASKING — HIDDEN (bit 7) and PROTECTED (bit 6) are attribute bits

- What: Bits 7 and 6 of the type byte are HIDDEN and PROTECTED flags
  respectively. They must be masked before comparing against the
  documented type values.
- Severity: structural (any rule that switches on file type)
- Source authority: Tech-Manual + SAMDOS-code
- Citation: Tech Manual L4351-4356:

  ```
  If the byte is 0 then the file has been erased. If the file is
  HIDDEN then bit 7 is set. If the file is PROTECTED then bit 6 is set.
  ```

  SAMDOS `c.s:1034-1040` reads `bit 6,a` to display protected files
  with `*` prefix; `e.s:248,263,282` mask `pntyp:and &1f` before
  switching on type.
- Dialect: all
- Test sketch: when interpreting `Type`, use `Type & 0x1F` (or
  `& 0x3F` if you want to keep bits visible) before comparing.

### DIR-ERASED-IS-ZERO — Erased slot is exactly Type==0

- What: A slot with type byte 0 is erased / free. SAMDOS checks this
  literally; bit 7 + low-5 zero is still erased.
- Severity: structural
- Source authority: SAMDOS-code + Tech-Manual
- Citation: `samdos/src/c.s:1133-1143` (`fdhf` — "TEST FOR FREE
  DIRECTORY SPACE"): `ld a,(hl); and a; jr nz,fdhd`.
  Tech Manual L4351-4354.
- Dialect: all
- Test sketch: slot is free iff `data[0x00] == 0`.

### DIR-NAME-PADDING — Filename is 10 bytes, space-padded ASCII

- What: Dir bytes 0x01..0x0A hold a 10-byte filename. SAMDOS pads
  short names with `0x20` (space). Matching is case-insensitive.
- Severity: cosmetic (padding) / structural (length)
- Source authority: SAMDOS-code + Tech-Manual
- Citation: Tech Manual L4358-4359. SAMDOS `f.s:835-846` pads with
  `&20`:

  ```asm
  ld hl,nstr1
  ld a,15
  evnm1: ld (hl),&20
         inc hl
         dec a
         jr nz,evnm1
  ```

  Case-insensitive match: `c.s:1155-1180` (`cknam`) uses `xor (hl);
  and &df` to ignore the case bit.
- Dialect: all
- Test sketch: slot is used → name bytes are either ASCII or `0x20`;
  warn on control chars `< 0x20` or `>= 0x80`.

### DIR-NAME-NOT-EMPTY — Used slot has non-blank name

- What: A used slot (`Type != 0`) should have at least one non-space
  character in bytes 0x01-0x0A. SAMDOS uses an FF-prefix sentinel for
  "null name" at save time (`rom-disasm:22094` `LD (HL),0FFH ;ASSUME
  NAME IS NULL`), so an all-FF or all-space name is not necessarily
  invalid but is suspicious.
- Severity: inconsistency
- Source authority: ROM
- Citation: `rom-disasm:22093-22105`: SAVE refuses an empty name with
  error 18 "Invalid file name".
- Dialect: all
- Test sketch: warn if used slot has all-space or all-FF name.

### DIR-SECTORS-BIG-ENDIAN — Sector count at 0x0B-0x0C is big-endian

- What: The 16-bit sector count at dir bytes 0x0B-0x0C is stored
  big-endian (high byte at 0x0B, low byte at 0x0C). Almost every
  other field is little-endian; this is a known SAM quirk.
- Severity: structural
- Source authority: samfile-implicit + Tech-Manual
- Citation: Tech Manual L4360-4361; `samfile.go:473`:

  ```go
  Sectors: uint16(data[0x0b])<<8 | uint16(data[0x0c]),  // big endian!
  ```
- Dialect: all
- Test sketch: read Sectors as big-endian when parsing.

### DIR-FIRST-SECTOR-VALID — First-sector points at a data sector

- What: Dir bytes 0x0D-0x0E (FirstSector.Track, FirstSector.Sector)
  must be a valid data sector: track in {4..79, 128..207}, sector in
  {1..10}. Track 0 means "erased" per samfile's `Used()` heuristic.
- Severity: fatal
- Source authority: samfile-implicit (Tech-Manual implies but does
  not enforce)
- Citation: `samfile.go:597-605`: `if fe.FirstSector.Track == 0 {
  return false }`; `samfile.go:611-616`: `Output` errors on
  `Track < 4`.
- Dialect: all
- Test sketch: for every used slot, `FirstSector.Track in
  {4..79, 128..207}` and `FirstSector.Sector in {1..10}`.

### DIR-SECTORS-MATCHES-CHAIN — Sector count equals length of chain walk

- What: The count at 0x0B-0x0C must equal the number of sectors
  visited when walking the chain from FirstSector to the (0,0)
  terminator.
- Severity: structural
- Source authority: samfile-implicit
- Citation: `samfile.go:743-754` reads `fe.Sectors` chunks and stops:

  ```go
  for {
      copy(raw[510*i:], filepart.Data[:])
      i++
      if i == fe.Sectors {
          break
      }
      sectorData, err = di.SectorData(filepart.NextSector)
      ...
  }
  ```
- Dialect: all
- Test sketch: walk chain, count visited sectors, compare to
  `Sectors`. If walk runs past a (0,0) terminator before reaching
  `Sectors`, file is short. If walk hits `Sectors` but next link is
  not (0,0), file is long.

### DIR-SECTORS-MATCHES-MAP — Sector count equals popcount of SectorAddressMap

- What: The number of bits set in dir bytes 0x0F-0xD1 (the 195-byte
  / 1560-bit SectorAddressMap) must equal the sector count at
  0x0B-0x0C.
- Severity: structural
- Source authority: Tech-Manual + samfile-implicit
- Citation: Tech Manual L4405-4414. samfile derives the map from
  bitset operations in `addFile` (`samfile.go:958-960`):

  ```go
  offset, mask := freeSectors[i].SAMMask()
  fe.SectorAddressMap[offset] |= byte(mask)
  ```
- Dialect: all
- Test sketch: `popcount(map) == Sectors`.

### DIR-SECTORS-NONZERO — Used slot has Sectors >= 1

- What: A used slot must own at least one sector (the body must hold
  at least the 9-byte header).
- Severity: structural
- Source authority: SAMDOS-code
- Citation: `samdos/src/c.s:919-951` (`fns6`) — allocator increments
  `cntl/cnth` (the sector counter) before returning; a SAVE always
  goes through at least one `fnfs` call, so Sectors >= 1.
- Dialect: all
- Test sketch: every used slot has `Sectors >= 1`.

### DIR-SAM-WITHIN-CAPACITY — All bits in SectorAddressMap are in 1..1559

- What: Bits 1560..1599 (high 5 bits of byte 194 plus a notional
  byte 195) lie outside the disk's data-sector domain. Tech Manual
  L4405-4406 explicitly says 1560 bits = 1560 sectors. Bits beyond
  must be zero.
- Severity: inconsistency
- Source authority: Tech-Manual
- Citation: Tech Manual L4405-4406:

  ```
  SAMDOS allocates 195 bytes to the sector address map, giving 1560
  bits, which is the exact number of sectors available for storage
  on the drive.
  ```
- Dialect: all
- Test sketch: `byte 194 & 0xE0 == 0` (top 3 bits beyond bit 1559
  are clear). Plus 195 bytes total — no overflow beyond that
  position.

---

## 3. Sector-chain rules

### CHAIN-LINK-AT-510-511 — Bytes 510-511 of every data sector are the next-sector link

- What: Each data sector reserves the last 2 bytes for a (track,
  sector) link. The payload occupies bytes 0..509 (510 bytes).
- Severity: structural
- Source authority: Tech-Manual + SAMDOS-code
- Citation: Tech Manual L4277-4280:

  ```
  Although each data sector can hold 512 bytes, only 510 bytes of
  them are available for storage. The last two bytes of the data
  sector are used by the DOS to locate the next part of the file
  stored. Byte 511 holds the next track used by the file, while
  byte 512 holds the next sector.
  ```

  SAMDOS `b.s:33`: `ld hl,&8000+510` — `dos:` loader walks file
  body and reads (track, sector) at offset 510 of each sector.
- Dialect: all
- Test sketch: extract `track = sector_data[510]`, `sector =
  sector_data[511]` to find next link.

### CHAIN-TERMINATOR-ZERO-ZERO — Last sector has (0, 0) link

- What: The last sector of a file's chain has bytes 510-511 set to
  `(0, 0)`. Track 0 is reserved for the directory area, so this is
  unambiguous.
- Severity: structural
- Source authority: SAMDOS-code
- Citation: `samdos/src/b.s:104-110`:

  ```asm
  dos8:  dec hl
         ld e,(hl)
         dec hl
         ld d,(hl)
         ld a,d
         or e
         jr nz,dos     ; chain continues only if d|e is non-zero
  ```
- Dialect: all
- Test sketch: walk chain; last sector's `(byte 510, byte 511) ==
  (0, 0)`.

### CHAIN-NO-CYCLE — Chain has no cycles within `Sectors` walks

- What: Following links must reach the (0, 0) terminator in
  exactly `Sectors` steps without revisiting any sector.
- Severity: structural
- Source authority: samfile-implicit
- Citation: `samfile.go:743-754` (chain walk in `File`) does not
  detect cycles, but its termination depends on `i == fe.Sectors`;
  a cycle within fewer than `Sectors` steps would silently overwrite
  buffer contents in samfile but break SAMDOS.
- Dialect: all
- Test sketch: maintain a visited-sector set during walk; flag a
  repeat as a cycle.

### CHAIN-MATCHES-SAM — Walked sectors equal the bits set in SectorAddressMap

- What: The set of sectors visited during a chain walk must equal
  the set of bits set in dir bytes 0x0F-0xD1.
- Severity: structural
- Source authority: SAMDOS-code + samfile-implicit
- Citation: `samdos/src/c.s:1306-1343` (`cfsm` — Close File Sector
  Map): the allocator-side that writes both the SAM map and the
  chain links simultaneously; if these diverged, SAMDOS's "find next
  free sector" (`fnfs` at `c.s:895-951`) would re-allocate the same
  sector.
- Dialect: all
- Test sketch: bitmap returned by `SectorAddressMap.UsedSectors()`
  matches the multiset of sectors visited during the chain walk.

### CHAIN-FIRST-MATCHES-DIR — Chain starts at dir-entry FirstSector

- What: The first sector of the chain is exactly the (track,
  sector) at dir bytes 0x0D-0x0E.
- Severity: structural
- Source authority: SAMDOS-code
- Citation: `samdos/src/c.s:1183-1267` (`ofsm` opens a new file:
  calls `fnfs` then writes `(ix+ftrk),d` and `(ix+fsct),e` at
  c.s:1256-1257); these later become the dir-entry's 0x0D-0x0E.
- Dialect: all
- Test sketch: trivial.

---

## 4. Cross-entry consistency

### CROSS-NO-SECTOR-OVERLAP — No two used slots share a data sector

- What: The disk-wide allocation map (bitwise OR of every used
  slot's `SectorAddressMap`) has no overlap — each set bit must
  come from exactly one slot.
- Severity: fatal
- Source authority: SAMDOS-code
- Citation: SAMDOS allocator at `samdos/src/c.s:895-951` (`fnfs`):
  it scans the merged `sam` array (built from every dir entry) for
  free bits; a duplicate set bit means SAMDOS will refuse to reuse
  a free sector. The release notes for samfile v2.1.0 explicitly
  list this as a planned check.
- Dialect: all
- Test sketch: for every bit position, count which slots claim it;
  flag bits with count > 1.

### CROSS-NO-DUPLICATE-NAMES — No two used slots have the same name

- What: Filenames must be unique (case-insensitive) across used
  slots. SAMDOS's SAVE asks "OVERWRITE?" when it finds a duplicate
  (`samdos/src/c.s:1196-1219`), so on-disk duplicates indicate a
  partially-completed operation or a hand-rolled image.
- Severity: inconsistency
- Source authority: SAMDOS-code
- Citation: `samdos/src/c.s:1196-1219` (`ofsm` overwrite path) and
  the samfile v2.1.0 release notes.
- Dialect: all
- Test sketch: compare `Name.String()` (trimmed) case-insensitively
  across all used slots.

### CROSS-DIRECTORY-AREA-UNUSED — No bit in any SectorAddressMap covers tracks 0..3 (side 0)

- What: The directory tracks are excluded from the sector-address
  map by construction (`bitOffset = ... - 40`). It is impossible
  to encode a directory sector in the map. But check that the
  per-slot chain does not include T0..T3 either.
- Severity: structural
- Source authority: samfile-implicit
- Citation: `samfile.go:984-987` (`SAMMask` formula). For chains:
  any link `(track, sector)` with `track < 4` would be invalid.
- Dialect: all
- Test sketch: every link in every chain has `(track & 0x7F) >= 4`.

---

## 5. Body-header rules (the 9-byte file header)

### BODY-HEADER-AT-FIRST-SECTOR — Body starts with the 9-byte header

- What: The first 9 bytes of every used file's first sector are the
  body header in the layout documented at `sam-file-header.md` §1.
- Severity: structural
- Source authority: SAMDOS-code + Tech-Manual
- Citation: Tech Manual L4284-4295; SAMDOS `f.s:462-471` (`svhd`)
  writes 9 bytes via `sbyt` (= the sector-chain payload stream).
  `samfile.go:756-766` reconstructs `FileHeader` from `raw[0..8]`.
- Dialect: all
- Test sketch: read first 9 bytes of `FirstSector` payload.

### BODY-TYPE-MATCHES-DIR — Body header byte 0 == dir Type byte (masked)

- What: Body header byte 0 must equal dir byte 0 with HIDDEN /
  PROTECTED bits masked off. (SAMDOS keeps these in sync via
  `svhd` writing the same `hd001` value to dir 0xD3 and body byte
  0; the dir's own type byte is set separately by `ofsm`.)
- Severity: inconsistency
- Source authority: SAMDOS-code
- Citation: `samdos/src/c.s:1395-1408` (`gtfle`-derived `gtfl1`
  block) and `samdos/src/b.s:255` (`hd001:  defb &13`).
- Dialect: all
- Test sketch: `(dir[0] & 0x1F) == body[0]` for every used slot.

### BODY-LENGTHMOD16K-MATCHES-DIR — Body bytes 1-2 == dir bytes 0xF0-0xF1

- What: Body header `LengthMod16K` (LE 16-bit at bytes 1-2) must
  equal the dir entry's mirrored field at 0xF0-0xF1.
- Severity: inconsistency
- Source authority: SAMDOS-code
- Citation: `samdos/src/c.s:1376-1379` (`gtfle` reads 9 bytes from
  dir+211 into `hd001..page1` — which IS `hd001/hd0b1/hd0d1/.../
  page1`). And `samdos/src/c.s:1457-1472` reads `hd0b2/hd0d2`
  back into the DIFA representation — which IS dir 0xEC-0xF1.
- Dialect: all
- Test sketch: `body[1..3] == dir[0xF0..0xF2]` (and the mirror at
  0xD3-0xDB).

### BODY-PAGEOFFSET-MATCHES-DIR — Body bytes 3-4 == dir bytes 0xED-0xEE

- What: Body `PageOffset` (LE 16-bit at bytes 3-4) mirrors dir
  `StartAddressPageOffset` at 0xED-0xEE.
- Severity: inconsistency
- Source authority: SAMDOS-code
- Citation: same as BODY-LENGTHMOD16K-MATCHES-DIR; the 9-byte
  mirror at dir+211 includes bytes 3-4.
- Dialect: all
- Test sketch: `body[3..5] == dir[0xED..0xEF]` (and the mirror at
  0xD3-0xDB).

### BODY-EXEC-DIV16K-MATCHES-DIR — Body byte 5 == dir byte 0xF2

- What: Body byte 5 (`ExecutionAddressDiv16K`) mirrors dir byte
  0xF2. The ROM auto-exec gate (rom-disasm:22467-22484) checks BOTH
  for `0xFF` before deciding to auto-execute, so the two must agree
  to avoid surprises.
- Severity: structural (mismatch can cause unwanted auto-exec)
- Source authority: ROM
- Citation: ROM disasm L22471-22484 (see §0 above).
- Dialect: all
- Suppressed by: PROTECTED files take a different path
  (rom-disasm:22467-22469: `BIT 1,A; JR NZ,HDNSTP`); the requested
  exec is then skipped and only the LOADED path is honoured.
- Test sketch: `body[5] == dir[0xF2]`.

### BODY-EXEC-MOD16K-LO-MATCHES-DIR — Body byte 6 == dir byte 0xF3

- What: Body byte 6 holds only the low byte of
  `ExecutionAddressMod16K`. It mirrors dir byte 0xF3. The high byte
  (dir 0xF4) has no body-header counterpart.
- Severity: inconsistency
- Source authority: ROM + samfile-implicit
- Citation: ROM disasm L22472: `LD HL,(HDR+HDN+7)` — 16-bit LE pair
  starting at HDR offset 38 (= byte after HDR+HDN+6 = byte 37). The
  body header's byte 6 maps to HDL's `HDR+HDN+7` low byte via
  SAMDOS `dschd` (`h.s:74-90`). samfile models this via
  `FileHeader.ExecutionAddressMod16KLo` (samfile.go:194).
- Dialect: all
- Test sketch: `body[6] == (dir[0xF3] & 0xFF)`.

### BODY-PAGES-MATCHES-DIR — Body byte 7 == dir byte 0xEF

- What: Body `Pages` mirrors dir `Pages` at 0xEF.
- Severity: inconsistency
- Source authority: SAMDOS-code
- Citation: 9-byte mirror at dir+211 (see PR-12 hypothesis #1).
- Dialect: all
- Test sketch: `body[7] == dir[0xEF]`.

### BODY-STARTPAGE-MATCHES-DIR — Body byte 8 == dir byte 0xEC

- What: Body `StartPage` mirrors dir `StartAddressPage` at 0xEC.
  Note that only the low 5 bits are functional (page index 0..31);
  bits 5-7 are decorative and may differ between byte-perfect
  ROM-SAVE output and synthetic writers (cf. `sam-stub-audit.md`).
- Severity: inconsistency
- Source authority: SAMDOS-code + samfile-implicit
- Citation: 9-byte mirror at dir+211. samfile's `Start()`
  masks `(StartPage & 0x1F)+1`:

  ```go
  return uint32(fileHeader.PageOffset&0x3fff) |
         uint32((fileHeader.StartPage&0x1f)+1)<<14
  ```
- Dialect: all
- Test sketch: `body[8] == dir[0xEC]` (warn if decorative bits 5-7
  differ even though the low-5-bit page indices match).

### BODY-MIRROR-AT-DIR-D3-DB — Dir bytes 0xD3..0xDB exactly mirror body bytes 0..8

- What: SAMDOS keeps a verbatim 9-byte body-header cache at dir
  offset 0xD3-0xDB. See §0 for the verification. Byte 0xD2 is
  unused; expect zero.
- Severity: inconsistency (writer that omits this still works because
  SAMDOS re-reads from body on next access, but it deviates from
  canonical SAVE output)
- Source authority: SAMDOS-code
- Citation: §0, `samdos/src/f.s:462-471` + `c.s:1376-1379`.
- Dialect: all
- Test sketch: `dir[0xD2] == 0`; `dir[0xD3..0xDC] == body[0..9]`.
- Open questions: must `dir[0xD2]` be exactly zero, or merely
  ignored? Real-SAVE leaves it at 0 (verified empirically against
  FRED/Defender disks per `sam-stub-audit.md`).

### BODY-PAGEOFFSET-8000H-FORM — Bit 15 of PageOffset is set (convention)

- What: Real ROM SAVE writes the on-disk PageOffset with bit 15 set
  ("8000H form"); the offset value is taken from a section-C
  address `0x8000-0xBFFF`. Both samfile `Start()` and the ROM
  PDPSR2 decoder mask `& 0x3FFF` before use, so a bit-15-clear
  value still parses, but it deviates from convention.
- Severity: cosmetic
- Source authority: Tech-Manual + empirical-convention
- Citation: Tech Manual L3037-3052 (PAGE FORM convention). Real
  example: CHOMPER from `~/Downloads/GoodSamC2/x.mgt` body bytes
  `... d5 9c ...` (= 0x9CD5, bit 15 set). See
  `sam-file-header.md` §1 worked example.
- Dialect: all
- Test sketch: warn if `(body[3] | body[4] << 8) & 0x8000 == 0` for
  a non-zero offset.

### BODY-PAGE-LE-31 — StartPage's low 5 bits encode a page index 0..30

- What: Decoded page is `(StartPage & 0x1F) + 1`, so the linear
  start-page index after samfile's `+1` is 1..32 (i.e. 0x4000..
  0x84000). SAM has 512 KiB = 32 × 16 KiB pages; index 32 is the
  off-disk pseudo-page used as a marker.
- Severity: structural
- Source authority: samfile-implicit + ROM
- Citation: `samfile.go:248-249`:

  ```go
  return uint32(fileHeader.PageOffset&0x3fff) | uint32((fileHeader.StartPage&0x1f)+1)<<14
  ```

  ROM `PDPSR2` at rom-disasm:4499-4527 decodes the same format
  for `LOAD CODE` (exec).
- Dialect: all
- Test sketch: derived load address `(StartPage & 0x1F) * 16384 +
  (PageOffset & 0x3FFF)` is in 0x4000..0x7FFFF.

### BODY-BYTES-5-6-CANONICAL-FF — Real SAVE writes 0xFF 0xFF when no auto-exec

- What: When `ExecutionAddressDiv16K == 0xFF`, real ROM SAVE writes
  `0xFF 0xFF` for bytes 5-6 of the body header. samfile's
  `FileHeader.Raw()` (samfile.go:1011-1023) currently emits
  `0x00 0x00`. Both parse as "no auto-exec" because the gate at
  rom-disasm:22473 only checks byte 5 (and byte 5 alone for
  HDR+HDN+6 path), but the convention is `FF FF`.
- Severity: cosmetic
- Source authority: empirical-convention
- Citation: CHOMPER body bytes 5-6 = `ff ff` from
  `~/Downloads/GoodSamC2/x.mgt`. samfile.go:1011-1023 emits zeros;
  Pete's `sam-file-header.md` §1 flags this.
- Dialect: all
- Test sketch: cosmetic warning when byte 5 == 0xFF but byte 6 ==
  0x00.

---

## 6. File-type-specific: CODE (FT_CODE = 19)

### CODE-LOAD-ABOVE-ROM — Load address >= 0x4000

- What: A CODE file's decoded load address must be at least
  0x4000 (ROM occupies 0x0000-0x3FFF).
- Severity: fatal (loading into ROM corrupts BASIC's view of
  paging)
- Source authority: samfile-implicit + Tech-Manual
- Citation: Tech Manual L4316-4329 (ROM-skip arithmetic);
  `samfile.go:799-801`:

  ```go
  if loadAddress < 1<<14 {
      return fmt.Errorf("load address %v of %q is in ROM ...", loadAddress, name)
  }
  ```
- Dialect: all
- Test sketch: derived `Start()` >= 0x4000.

### CODE-LOAD-FITS-IN-MEMORY — Load address + length <= 0x80000

- What: Body length + load address must not overshoot SAM's 512 KiB
  address space.
- Severity: fatal
- Source authority: samfile-implicit
- Citation: `samfile.go:802-804`:

  ```go
  if int(loadAddress) > 1<<19-len(data) {
      return fmt.Errorf("load address %v of %v byte file %q higher than maximum allowed %v", ...)
  }
  ```

  Release notes v2.1.0: "verifying that code blocks load into
  memory without overshooting 512KB RAM limit".
- Dialect: all
- Test sketch: `Start() + Length() <= 0x80000`.

### CODE-EXEC-WITHIN-LOADED-RANGE — Execution address falls within loaded region

- What: If `ExecutionAddressDiv16K != 0xFF`, the decoded
  execution address must be in `[loadAddress, loadAddress + length)`.
- Severity: structural
- Source authority: samfile-implicit
- Citation: `samfile.go:805-810`:

  ```go
  if executionAddress > 0 && executionAddress < loadAddress {
      return fmt.Errorf("execution address %v of %q lower than load address %v", ...)
  }
  if int(executionAddress) >= int(loadAddress)+len(data) {
      return fmt.Errorf("execution address %v of %q is higher than the memory region it is loaded to (...)", ...)
  }
  ```

  Release notes v2.1.0: "verifying that execution addresses of
  code files are within code block load address".
- Dialect: all
- Test sketch: `loadAddr <= execAddr < loadAddr + length`.

### CODE-EXEC-FF-DISABLES — ExecutionAddressDiv16K == 0xFF means no auto-exec

- What: Setting dir byte 0xF2 to `0xFF` opts out of auto-exec on
  `LOAD CODE`. The ROM gate at rom-disasm:22473 (`CP 0FFH; JR
  NZ,HDLDEX`) takes this branch.
- Severity: structural (controls auto-exec behaviour)
- Source authority: ROM
- Citation: rom-disasm:22471-22479 — see §0.
- Dialect: all
- Test sketch: when `dir[0xF2] == 0xFF`, the file is
  auto-exec-disabled (no further check needed on 0xF3-0xF4).

### CODE-FILETYPEINFO-EMPTY — Dir 0xDD-0xE7 unused for CODE

- What: For FT_CODE, dir bytes 0xDD-0xE7 (`FileTypeInfo`) are
  unused. samfile's `AddCodeFile` leaves them at zero.
- Severity: cosmetic
- Source authority: samfile-implicit
- Citation: `samfile.go:798-827` (`AddCodeFile` does not set
  `FileTypeInfo`).
- Dialect: all
- Test sketch: warn if any byte in `dir[0xDD..0xE8]` is non-zero
  for a CODE file (unlikely; just a sanity check).

---

## 7. File-type-specific: SAM BASIC (FT_SAM_BASIC = 16)

### BASIC-FILETYPEINFO-TRIPLETS — Dir 0xDD-0xE5 holds three 3-byte PAGEFORM lengths

- What: For FT_SAM_BASIC, dir bytes 0xDD-0xE5 hold three
  page-form 3-byte lengths in order: `(NVARS-PROG)`, `(NUMEND-PROG)`,
  `(SAVARS-PROG)`. Each triplet is `[page, offset_lo, offset_hi]`.
- Severity: structural (used by ROM LOAD to restore NVARS/NUMEND/
  SAVARS sysvars)
- Source authority: ROM
- Citation: rom-disasm:22163-22180 — see §0. ROM SAVE writes
  these at `HDR+16/+19/+22` from `NVARS-PROG`, `NUMEND-PROG`,
  `SAVARS-PROG`. samfile decodes via `ProgramLength()` /
  `NumericVariablesSize()` / `GapSize()` (samfile.go:674-699).
- Dialect: all
- Test sketch: triplets are non-zero for FT_SAM_BASIC files;
  decoded `(NVARS-PROG) <= (NUMEND-PROG) <= (SAVARS-PROG) <=
  Length()`.

### BASIC-VARS-GAP-INVARIANT — `SAVARS-NVARS` is typically 604 (or 2156)

- What: Empirically, 93.9% of well-formed BASIC files have
  `(NUMEND-PROG) + (SAVARS-NUMEND) - (NVARS-PROG) == 604`
  (canonical SAMDOS) or `2156` (MasterDOS variant). See
  `sam-basic-save-format.md`.
- Severity: cosmetic (warn-only)
- Source authority: empirical-convention
- Citation: `sam-basic-save-format.md`, scan of 161 disks under
  `~/Downloads/`: 632/673 BASIC files satisfy 604; 32/673 satisfy
  2156. Mechanism: ROM `CLRSR` allocates 46 byte-pointers + 26
  PSVTAB + 20 PSVT2 + 20 PSVT2 (= 92 bytes) for the vars area, and
  the first numeric-var creation triggers a `MAKEROOM 0x0200` (=
  512 bytes) gap allocation per rom-disasm:10240-10255.
- Dialect: SAMDOS-2 = 604; MasterDOS = 2156
- Test sketch: warn if `SAVARS-NVARS` differs from 604 (SAMDOS) or
  2156 (MasterDOS) by more than a few bytes.

### BASIC-PROG-END-SENTINEL — Last byte of program area is 0xFF

- What: The tokenised program ends with a 0xFF sentinel byte. ROM
  detokeniser (sambasic.go:56-58) and SAMDOS-compatible writers all
  rely on this.
- Severity: structural
- Source authority: samfile-implicit + Tech-Manual
- Citation: `samfile/sambasic.go:56-58`:

  ```go
  if basic.Data[index] == 0xff {
      break
  }
  ```

  `sambasic/file.go:36-42` (ProgBytes): always appends `0xFF`.
- Dialect: all
- Test sketch: body byte at offset `NVARS-PROG - 1` (= last byte
  of program area) is 0xFF; equivalently the byte at `(decoded
  NVARS-PROG offset) - 1`.

### BASIC-LINE-NUMBER-BE — Line numbers stored big-endian

- What: Within the program area, each line starts with a 2-byte
  big-endian line number followed by a 2-byte little-endian line
  length and then `lineLen` bytes of tokenised body (terminated by
  0x0D within those bytes).
- Severity: structural
- Source authority: samfile-implicit + ROM
- Citation: `samfile/sambasic.go:62-63`:

  ```go
  lineNo := uint16(basic.Data[index])<<8 | uint16(basic.Data[index+1])
  lineLen := uint16(basic.Data[index+2]) | uint16(basic.Data[index+3])<<8
  ```

  `sambasic/file.go:21-34` Line.Bytes writes the same way.
- Dialect: all
- Test sketch: walk lines from PROG to NVARS; assert each line's
  length matches its actual size, all line numbers are within
  documented range (1..16383 for SAM, 1..9999 typical).

### BASIC-STARTLINE-FF-DISABLES — ExecutionAddressDiv16K == 0xFF means no auto-RUN

- What: For FT_SAM_BASIC, dir byte 0xF2 acts as an autorun marker:
  `0xFF` = no auto-RUN; `0x00` = auto-RUN at line `dir[0xF3..0xF4]`.
  ROM SAVE writes this at rom-disasm:22136-22141:

  ```
  LD HL,HDR+HDN+6        ;PTR TO AUTO-RUN LINE AREA
  LD (HL),0              ;FLAG 'AUTORUN'
  INC HL; LD (HL),C; INC HL; LD (HL),B  ;PLACE LINE NO.
  ```
- Severity: structural
- Source authority: ROM
- Citation: rom-disasm:22136-22141.
- Dialect: all
- Test sketch: for FT_SAM_BASIC, `dir[0xF2] in {0x00, 0xFF}`; if
  `0x00`, `dir[0xF3..0xF5]` is a valid line number.

### BASIC-STARTLINE-WITHIN-PROG — Auto-RUN line is in the saved program

- What: If auto-RUN is enabled, the line number at 0xF3-0xF4 should
  correspond to an actual line number in the program area.
- Severity: cosmetic (warn-only — auto-RUN of a missing line just
  errors with "Statement lost", not a corruption)
- Source authority: empirical-convention
- Citation: ROM auto-RUN dispatch via `NEW PPC` sysvar after LOAD;
  if line doesn't exist, BASIC errors. No code citation enforces
  pre-LOAD validation.
- Dialect: all
- Test sketch: walk program, collect line numbers, check
  `dir[0xF3..0xF5]` is among them.

### BASIC-MGTFLAGS-20 — MGTFlags is typically 0x20 for BASIC files

- What: Real-world BASIC files have `MGTFlags == 0x20`. The exact
  semantics are undocumented (Tech Manual L4369: "MGT use only"),
  but Defender and 50%+ of canonical disks set this. Our M0 disk
  requires it for boot (`test-mgt-byte-layout.md` §slot 1).
- Severity: inconsistency (empirically load-bearing for M0 — escalate
  to structural if a code citation surfaces)
- Source authority: empirical-convention
- Citation: `test-mgt-byte-layout.md` 0xDC slot-1 comment ("MGTFlags
  = 0x20 — required for M0 boot. Defender, pete-made.mgt, and 50%+
  of canonical disks set this. Semantics not fully documented but
  clearly load-bearing.").
- Dialect: SAMDOS-2 + MasterDOS (probably; unverified)
- Test sketch: warn if `dir[0xDC] != 0x20` for a BASIC file.
- Open questions: what code path reads `MGTFlags`? No grep hit in
  SAMDOS source. Likely consumed by MasterDOS or the BASIC ROM
  somewhere we haven't traced. Until then this is convention.

---

## 8. File-type-specific: arrays (FT_NUM_ARRAY = 17, FT_STR_ARRAY = 18)

### ARRAY-FILETYPEINFO-TLBYTE-NAME — Dir 0xDD-0xE7 holds TLBYTE + array name

- What: For FT_NUM_ARRAY / FT_STR_ARRAY, dir bytes 0xDD-0xE7 hold
  the array's TLBYTE (type/length byte) followed by its 10-byte
  name.
- Severity: structural (used by `MERGE` to recover variable names)
- Source authority: ROM
- Citation: rom-disasm:22354-22357 (`E1D7`): `LD HL,TLBYTE; LD
  DE,HDR+16; LD BC,11; LDIR`. Tech Manual L4371-4372.
- Dialect: all
- Test sketch: warn if all 11 bytes are zero for an array file.
- Open questions: TLBYTE bit-layout (type code + length high bits)
  is not fully documented in this audit; the ROM E019 block calls it
  "TLBYTE" but doesn't reverse it. Skip until needed.

---

## 9. File-type-specific: SCREEN (FT_SCREEN = 20)

### SCREEN-MODE-AT-0xDD — Dir byte 0xDD holds the screen mode

- What: For FT_SCREEN, dir byte 0xDD is the screen MODE (1-4 on
  SAM). The remaining bytes 0xDE-0xE7 are unused.
- Severity: structural
- Source authority: ROM
- Citation: rom-disasm:22259 (`E146`): `LD (HDR+16),A ;MODE`.
  Tech Manual L4373-4374.
- Dialect: all
- Test sketch: `dir[0xDD] in {1, 2, 3, 4}` for FT_SCREEN.

### SCREEN-LENGTH-MATCHES-MODE — Body length matches mode's screen size

- What: For FT_SCREEN, the body length should match the documented
  screen size for the given mode: mode 1 = 6912 bytes, mode 2 = 6912,
  modes 3-4 = 24576 bytes.
- Severity: structural
- Source authority: Tech-Manual
- Citation: Tech Manual modes table.
- Dialect: all
- Test sketch: cross-reference mode byte and `Length()`.
- Open questions: exact mode/length mapping; needs Tech Manual
  cross-check at finalisation time.

---

## 10. File-type-specific: ZX snapshot (FT_ZX_SNAPSHOT = 5)

### ZXSNAP-LENGTH-49152 — Body is 49,152 bytes

- What: A ZX 48K snapshot has a 49,152-byte body (48 KiB main RAM).
- Severity: structural
- Source authority: SAMDOS-code
- Citation: `samdos/src/d.s:660-661` (`snlen: defw 49152`).
  SAMDOS NMI snapshot save uses this constant unconditionally.
- Dialect: SAMDOS-2 (snapshot save is in SAMDOS source)
- Test sketch: `Length() == 49152` for FT_ZX_SNAPSHOT.

### ZXSNAP-LOAD-ADDR-16384 — Body header start = 16384

- What: ZX snapshot start address is 0x4000 (ZX RAM base).
- Severity: structural
- Source authority: SAMDOS-code
- Citation: `samdos/src/d.s:660-663` (`snadd: defw 16384`).
- Dialect: SAMDOS-2
- Test sketch: derived `Start() == 16384`.

---

## 11. Boot-file rules (slot 0 + T4S1)

### BOOT-OWNER-AT-T4S1 — Some slot's FirstSector == (4, 1) on a bootable disk

- What: For an image to be bootable on real SAM hardware, some
  directory entry's FirstSector must be (4, 1) so that the ROM
  BOOTEX reads the right sector at `&8000`.
- Severity: fatal (bootability)
- Source authority: ROM
- Citation: rom-disasm:20473-20598 (`BOOTEX`):

  ```
  LD DE, 0401H    ; set track=4, sector=1
  RSAD: ...       ; read 512 bytes to HL=8000H
  ```
- Dialect: all
- Test sketch: search all used slots for one whose `FirstSector ==
  Sector{Track: 4, Sector: 1}`. Bootability requires at least one
  such slot.

### BOOT-SIGNATURE-AT-256 — T4S1 bytes 256-259 are "BOOT" (case-insensitive)

- What: For ROM BOOTEX to dispatch to the loaded sector, bytes
  256-259 of T4S1 must spell `B O O T` (any case; bit 7 ignored).
- Severity: fatal (bootability)
- Source authority: ROM
- Citation: rom-disasm:20582-20598 (`BTNOE`/`BTCK`/`BTLY`):

  ```
  LD DE, 80FFH; LD HL, BTWD; LD B, 4
  BTCK: INC DE; LD A,(DE); XOR (HL); AND 5FH; JR Z,BTLY
        RST 8 / DB 53
  ```

  `AND 5FH` masks bits 5 and 7, so case + bit 7 are ignored.
  `BTWD` is at FB94H: `42 4F 4F D4` = "BOOT" (last byte has bit 7
  set per BASIC keyword convention).
- Dialect: all
- Test sketch: read T4S1 bytes 256..259; compare with `0x42 0x4F
  0x4F 0x54` after AND 0x5F.
- Open questions: what offset within the file body is byte 256 of
  T4S1? With the standard 9-byte body header, byte 256 of the
  sector = body offset 247. samdos2 is engineered so that its
  body byte 247 happens to be a token-table entry containing
  "BOOT" (`samdos2.reference.bin` offset 0xF7).

### BOOT-ENTRY-POINT-AT-9 — JP 8009H expects code at body offset 0

- What: After signature match, ROM does `JP 8009H`. The sector
  buffer is at 0x8000-0x81FF, so 0x8009 is sector-buffer offset 9
  = body offset 0 (after the 9-byte header). The file body's first
  byte must therefore be valid Z80 code.
- Severity: fatal (bootability)
- Source authority: ROM
- Citation: rom-disasm:20598 (final `JP 8009H` of BOOTEX); SAMDOS
  source `b.s:27` sets `org.adjust = 9` so real code is at body
  offset 0.
- Dialect: all
- Test sketch: derived requirement; verify body byte 0 of the T4S1
  file is a plausible Z80 opcode (cosmetic warn).

### BOOT-FILE-TYPE-IGNORED — Type byte of boot file is irrelevant

- What: ROM BOOTEX does not consult the directory entry at all; it
  reads T4S1 raw and looks for the signature. Any type byte works.
- Severity: cosmetic (note for `verify` to not flag)
- Source authority: ROM
- Citation: rom-disasm:20473-20598 — no `RST 8`/dir-read path.
  `sam-disk-format.md` §5.5 documents this explicitly. SAMDOS's
  own SAVE writes type 3 for itself (samdos source `b.s:14-22`)
  but build-disk.sh writes type 19 (CODE) to keep `samfile ls`
  happy.
- Dialect: all
- Test sketch: don't reject T4S1 files by type.

---

## 12. HIDDEN / PROTECTED / ERASED handling

### ATTR-HIDDEN-NOT-LISTED — Bit 7 of type suppresses DIR display

- What: HIDDEN files (bit 7 of type) are skipped by SAMDOS's DIR
  listing.
- Severity: structural (semantic, not validity)
- Source authority: SAMDOS-code
- Citation: `samdos/src/c.s:1023-1026` (`fdh4` block):

  ```asm
  bit 7,a       ; HIDDEN
  jp nz,fdhd    ; skip to next entry
  ```
- Dialect: all
- Test sketch: not a validity check; a hint for `samfile ls
  --hidden`.

### ATTR-PROTECTED-NO-OVERWRITE — Bit 6 of type prevents ERASE/OVERWRITE

- What: PROTECTED files (bit 6) cannot be erased unless the user
  confirms (`samdos/src/e.s:233-235` in eraz3 path checks `bit 6`
  before allowing erase). Auto-exec is also affected: rom-disasm:22467-
  22469 (`HDL+HFG bit 1`) gates the requested-exec branch.
- Severity: structural (semantic)
- Source authority: SAMDOS-code + ROM
- Citation: `samdos/src/e.s:232-237` (`eraz3` beep-and-skip path);
  rom-disasm:22467-22469.
- Dialect: all
- Test sketch: not a validity check.

### ATTR-ERASED-SUPPRESSES-ALL — Type == 0 means free; skip all body / chain rules

- What: When `Type == 0` the slot is free; sector chain, body
  header, all `FileTypeInfo`/`StartAddress*` mirrors etc. are
  irrelevant. Only the dir bytes 0x00 (=0) and possibly 0x0D (which
  samfile uses as an additional "Track==0 = free" heuristic, not
  enforced by SAMDOS) matter.
- Severity: structural (suppresses other rules)
- Source authority: SAMDOS-code
- Citation: `samdos/src/c.s:1133-1143` (`fdhf`); samfile
  `Used()` at `samfile.go:597-605`.
- Dialect: all
- Test sketch: skip all per-file rules for erased slots.

---

## 13. Dialect notes

### DIALECT-MASTERDOS-MGTFLAGS — MasterDOS uses additional MGTFlags bits

- What: MasterDOS sets bits beyond `0x20` in `MGTFlags` to track
  per-file attributes that SAMDOS-2 ignores. Exact bit semantics
  undocumented in our corpus.
- Severity: cosmetic
- Source authority: empirical-convention
- Citation: none in our corpus; out-of-scope until MasterDOS-side
  source is available.
- Dialect: MasterDOS
- Open questions: MasterDOS source location, MGTFlags bit map.

### DIALECT-MASTERDOS-GAP-2156 — MasterDOS BASIC files have `SAVARS-NVARS == 2156`

- What: BASIC files saved by MasterDOS have a 2156-byte vars+gap
  trailer, vs SAMDOS-2's 604.
- Severity: cosmetic (warn-only)
- Source authority: empirical-convention
- Citation: `sam-basic-save-format.md`: 32/673 files in 161-disk
  scan satisfy the 2156 invariant.
- Dialect: MasterDOS
- Test sketch: when reporting BASIC trailer size, accept both 604
  and 2156 without warning.

### DIALECT-SAMDOS-1-TYPE-3 — SAMDOS-1 "auto-include header" type 3

- What: SAMDOS source's `if defined (include-header)` block at
  `b.s:14-22` emits a 9-byte header beginning with type byte 3 for
  SAMDOS itself. The shipped samdos2 binary does NOT include this
  header; it is a build-time option for older SAMDOS variants.
  Type 3 in `samdos/src/e.s:330` is also the SAMDOS DIR-display
  alias for "ZX $.ARRAY", so a real type-3 file gets displayed as
  "ZX $.ARRAY" — collision is intentional in SAMDOS (the DIR
  display table doesn't care).
- Severity: inconsistency
- Source authority: SAMDOS-code
- Citation: `samdos/src/b.s:12-22`, `samdos/src/e.s:330-331`.
- Dialect: SAMDOS-1
- Test sketch: don't flag type 3 as invalid; report as "ZX $.ARRAY
  or SAMDOS-1 header".

### DIALECT-HOOK-128-DEAD-CODE — Hook 128 (BTHK) does not auto-RUN in samdos2

- What: Tech Manual L4524 ("INIT 128 dec ... Initialise and look
  for AUTO file") is aspirational. samdos2's `init` /`initx`
  (h.s:215-218) just sets CURCMD=LOAD and returns. The `hauto`
  routine (h.s:224) that would actually load AUTO is dead code —
  no caller exists in the shipped samdos2 binary. FRED-style boot
  disks supply their own T4S1 bootstrap.
- Severity: cosmetic (documents an aspirational rule a verifier
  should NOT enforce)
- Source authority: SAMDOS-code
- Citation: `samdos/src/h.s:215-237`; `samdos2-auto-run-analysis.md`.
- Dialect: SAMDOS-2
- Test sketch: do not treat absence of AUTO file as a bootability
  issue.
- Open questions: did SAMDOS-1's hook 128 implement the Tech Manual
  spec? No source available to check.

---

## 14. Cosmetic / canonical-output rules (warn-only)

### COSMETIC-RESERVEDA-FF — Real SAVE writes 0xFF in dir 0xE8-0xEB

- What: `ReservedA` (4 bytes at dir 0xE8-0xEB) is filled with `0xFF`
  by real ROM SAVE (rom-disasm:22078-22080: the HDR-initialisation
  loop writes `0xFF` to 14 bytes after the name-area space-fill).
- Severity: cosmetic
- Source authority: ROM
- Citation: rom-disasm:22076-22080:

  ```
  LD B,14
  HDCLP2: LD (HL),0FFH       ;CLEAR REST WITH FFH
          INC HL
          DJNZ HDCLP2
  ```
- Dialect: all
- Test sketch: cosmetic warning when `dir[0xE8..0xEC] != [0xFF,
  0xFF, 0xFF, 0xFF]`.

### COSMETIC-RESERVEDB-FILL — Real SAVE writes specific bytes in dir 0xF5-0xFF

- What: `ReservedB` (11 bytes at dir 0xF5-0xFF) includes 8 bytes
  "spare" + 2 bytes "MGT future" per Tech Manual L4399-4400. Real
  SAVE leaves these at the HDR's post-initialisation pattern (mix
  of `0xFF`-fill from HDCLP2 and zero-init from the SAVE-time
  computation).
- Severity: cosmetic
- Source authority: Tech-Manual
- Citation: Tech Manual L4399-4400.
- Dialect: all
- Test sketch: not a rule; document only.

### COSMETIC-STARTPAGE-DECORATIVE-BITS — Bits 5-7 of dir 0xEC differ between writers

- What: The functional `StartPage` field is bits 0-4 only; bits 5-7
  are decorative. Real-SAVE output on FRED 02 / Defender disks
  records samdos2's StartAddressPage as `0x7D` (= 0x60 decorative
  bits + 0x1D page index); samfile's `AddCodeFile` writes only
  `0x1D`. The decoded address is the same; only byte-equal
  diffs differ.
- Severity: cosmetic
- Source authority: empirical-convention
- Citation: `samfile.go:845-857` (`SetStartAddressPageRaw` exists
  precisely to support byte-perfect parity); `sam-stub-audit.md`
  documents the convention.
- Dialect: all
- Test sketch: when comparing dir entries byte-for-byte, mask off
  bits 5-7 of `dir[0xEC]` first.

---

## 15. Rules carried over from the samfile v2.1.0 release notes

The v2.1.0 release notes (2023-01-09, https://github.com/petemoore/samfile/releases/tag/v2.1.0)
listed planned integrity checks. Mapping to this catalog:

| Release-note item | Catalog rule(s) |
|---|---|
| verifying that no sectors are shared across file entries | CROSS-NO-SECTOR-OVERLAP |
| verifying that Sector Address Maps are consistent with sector references | DIR-SECTORS-MATCHES-MAP + CHAIN-MATCHES-SAM |
| verifying that no two files have the same name | CROSS-NO-DUPLICATE-NAMES |
| verifying that the data in the file headers are consistent with the same data in the file entries | all BODY-*-MATCHES-DIR rules |
| verifying that execution addresses of code files are within code block load address | CODE-EXEC-WITHIN-LOADED-RANGE |
| verifying that code blocks load into memory without overshooting 512KB RAM limit | CODE-LOAD-FITS-IN-MEMORY |
| verifying that file names do not contain unprintable characters | DIR-NAME-PADDING (partial — extend with a stricter charset check) |
| verifying that disk images contain a suitable dos file | BOOT-OWNER-AT-T4S1 + BOOT-SIGNATURE-AT-256 |
| verifying that SAM BASIC files contain SAM BASIC (etc.) | BASIC-FILETYPEINFO-TRIPLETS + BASIC-PROG-END-SENTINEL + BASIC-LINE-NUMBER-BE |
| verifying that files do not contain empty sectors (= sector count is lowest possible) | New: derived from `ceil((9 + Length()) / 510) == Sectors`. Add as `CHAIN-SECTOR-COUNT-MINIMAL`. |

### CHAIN-SECTOR-COUNT-MINIMAL — Sector count is `ceil((9 + Length()) / 510)`

- What: A well-allocated file uses exactly `ceil((9 + body length) /
  510)` sectors — no padding sectors. Per release notes:
  "verifying that files do not contain empty sectors".
- Severity: cosmetic (warn-only — extra trailing sectors waste
  space but don't break anything)
- Source authority: samfile-implicit
- Citation: `samfile.go:919`:

  ```go
  requiredSectorCount := (len(data) + 9 + 509) / 510
  ```
- Dialect: all
- Test sketch: `Sectors == ceil((Length() + 9) / 510)`.

---

## 16. Sources index

- **SAMDOS source** (`~/git/samdos/src/`):
  - `a.s` — top-of-build constants (port addresses, dir-entry RAM
    offsets `chbtlo..chflag..recnum`, file-buffer layout `chbtlo,
    chbthi, chrec, chname, chflag, chdriv, recflg, recnum, rclnlo,
    rclnhi`).
  - `b.s:7-22` — optional 9-byte body header (type 3) for SAMDOS
    itself.
  - `b.s:27` — `org.adjust = 9` (boot entry point).
  - `b.s:33-126` — `dos:` bootstrap loader (the code that runs
    after ROM BOOTEX's `JP 8009H`).
  - `b.s:255-260` — `hd001..page1` 9-byte body-header RAM cache.
  - `b.s:497-540` — `samhk` hook-code dispatch table.
  - `c.s:1133-1143` — `fdhf` ("test for free directory space").
  - `c.s:1183-1267` — `ofsm` (open file sector map, allocator
    entry).
  - `c.s:1306-1343` — `cfsm` (close file sector map, finaliser).
  - `c.s:1376-1379` — `gtfle` reads 9-byte body-header cache from
    dir offset 211.
  - `c.s:1457-1472` — `gtfle` reads `hd0b2/hd0d2` into DIFA.
  - `c.s:895-951` — `fnfs` (find next free sector).
  - `c.s:1155-1180` — `cknam` (case-insensitive filename match).
  - `d.s:660-663` — ZX snapshot length & start address constants.
  - `e.s:322-356` — `drtab` (DIR-display type table: 1-12, 16-20).
  - `f.s:462-471` — `svhd` (save 9-byte body header to dir +211
    AND to file body via `sbyt`).
  - `f.s:494-497` — `ldhd` (load 9-byte body header).
  - `h.s:74-90` — `dschd` (load body header, populate hd001..).
  - `h.s:201-237` — `autnam` template, `init`/`initx`/`hauto`
    (hook 128 — dead code).
  - `h.s:336-361` — `hconr` (UIFA → DIFA / hd001).

- **SAM ROM v3.0 annotated disassembly**
  (`docs/sam/sam-coupe_rom-v3.0_annotated-disassembly.txt`):
  - L4499-4527 — `PDPSR2` (REL PAGE FORM decoder for LOAD CODE
    exec).
  - L14852-14861 — `TSURPG`/`SELURPG` (upper-page selector).
  - L20453-20471 — BOOT token handler (`D8CD`).
  - L20473-20598 — `BOOTEX` (raw-sector load + signature check).
  - L22025-22054 — `E019` HDR/HDL header buffer documentation.
  - L22136-22141 — BASIC autorun-line setup at HDR+HDN+6.
  - L22163-22180 — SAVE writes three BASIC prog-length triplets.
  - L22247 — `LD A,19` (CODE type on SAVE/LOAD CODE path).
  - L22259 — `LD (HDR+16),A ;MODE` (SCREEN$ mode store).
  - L22354-22357 — `LD HL,TLBYTE; LD BC,11; LDIR` (array
    TLBYTE+name to HDR+16).
  - L22467-22484 — LOAD CODE auto-exec gate (HDLDEX,
    R1OFFCLBC).
  - L22057-22119 — `SLMVC`/`HDR2` (SAVE/LOAD entry, HDR init).
  - L26919 — `BTWD` "BOOT" keyword bytes.

- **Tech Manual v3.0**
  (`docs/sam/sam-coupe_tech-man_v3-0.txt`):
  - L2974-3068 — 80-byte HDR/HDL buffer.
  - L3037-3052 — REL PAGE FORM convention (8000H-form offset).
  - L4262-4275 — disk geometry.
  - L4277-4280 — sector chain bytes 510-511.
  - L4284-4332 — 9-byte body header layout.
  - L4304-4314 — file-type values.
  - L4338-4400 — 256-byte directory entry.
  - L4403-4427 — sector address map / BAM.
  - L4524 — Hook 128 INIT spec (aspirational; cf.
    samdos2-auto-run-analysis.md).
  - L4548 — Hook code explanations.

- **sam-aarch64 notes** (`docs/notes/`):
  - `sam-disk-format.md` — geometry, sector map, BOOT mechanism.
  - `sam-file-header.md` — 9-byte body header, dir entry,
    HDR/HDL buffer.
  - `sam-basic-save-format.md` — BASIC trailer (604 / 2156
    invariant).
  - `samdos2-auto-run-analysis.md` — hook 128 is dead code.
  - `test-mgt-byte-layout.md` — byte-by-byte M0 boot-disk dump.
  - `sam-stub-audit.md` — decorative-bit-vs-functional-bit
    StartPage convention.
  - `sam-paging.md` — REL PAGE FORM encoder/decoder.

- **samfile** (`~/git/samfile/`):
  - `samfile.go:71` — `DiskImage [819200]byte`.
  - `samfile.go:355-368` — EDSK rejection.
  - `samfile.go:389-398` — sector / track range validation.
  - `samfile.go:438-446` — directory walk.
  - `samfile.go:470-496` — `FileEntryFrom`.
  - `samfile.go:597-605` — `Used()` heuristic.
  - `samfile.go:611-647` — `Output()` per-type printer.
  - `samfile.go:674-723` — BASIC / CODE accessors.
  - `samfile.go:736-771` — chain walk.
  - `samfile.go:798-827` — `AddCodeFile` invariants.
  - `samfile.go:865-895` — `AddBasicFile`.
  - `samfile.go:913-965` — `addFile`.
  - `samfile.go:984-995` — `SAMMask` (and its operator-precedence
    bug, see `sam-disk-format.md` §3.4).
  - `sambasic/file.go` — BASIC body assembly + NVARS/NUMEND/SAVARS
    offsets.
  - `sambasic.go` — tokeniser-aware printer.
  - Release notes v2.1.0 — initial wish-list of integrity checks.

- **SimCoupé** (`~/git/simcoupe/Base/Disk.{cpp,h}`):
  - `Disk.h:28-41` — MGT geometry constants.
  - `Disk.cpp:164` — cylinder-interleaved offset formula.
