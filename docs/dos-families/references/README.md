# External reference materials

Per-family directories embed the commented assembly source for
the DOSes whose source we have. The materials below live
outside the samfile repo and are too large or too remote to
ship inline — but every agent reasoning about a family should
know they exist.

## ROM v3.0 annotated disassembly

- **Path:** `/Users/pmoore/git/sam-aarch64/docs/sam/sam-coupe_rom-v3.0_annotated-disassembly.txt`
- **Size:** ~1.1 MB, 27353 lines
- **Use for:** the LOAD-path semantics that every DOS plugs
  into. Grep for `BOOTEX`, `BTCK`, `LDHD`, `GTFLE`, `HCONR`
  to walk the SAVE / LOAD / boot-sector chain.
- **Source:** captured in `~/git/sam-aarch64/docs/sam/` and in
  `~/git/migrate-build-disk-to-go/docs/sam/` (same file).

## SAMDOS source — upstream tokenised archive

- **Upstream:** https://ftp.nvg.ntnu.no/pub/sam-coupe/sources/SamDos2InCometFormatMasterv1.2.zip
- **Contains:** a SAM .dsk image with 28 source files
  (`a1..h2.S`, `comp1..comp5.S`, `gm1`, `gm2`, `ldit1..3`,
  `ld1`). These are in COMET tokenised format — not plain
  text. Use `samfile extract -i <dsk>` to pull them out, then
  a Comet detokeniser (or compare against Stefan Drissen's
  git history at https://github.com/stefandrissen/samdos which
  has the lowered / detokenised plain-text form).
- **Local clean copy:** `/Users/pmoore/git/samdos` (Drissen's repo).
  HEAD assembles to `9bc0fb4b949109e8-samdos2/variants/3cca541beb3f9fe9.bin`.
  Earlier comp1..comp5 versions are reachable via `git log`.

## MasterDOS source

- **Upstream:** https://github.com/dandoore/masterdos (Dan Doore)
- **Local clean copy:** `/Users/pmoore/git/masterdos`
- `src/masterdos23.asm` assembles to
  `152b811ed65b651d-masterdos-v2.3/body.bin`. v2.2 and v2.3
  differ only at points labelled `Fix_*` in the source.

## How to disassemble a DOS body when no source is available

Most non-SAMDOS / non-MasterDOS families have no upstream
source. Disassemble the body directly:

```
# z80dasm: the load address is recorded in each family's README
z80dasm -a -t -o0x008009 body.bin > body.z80.s
```

or use any equivalent z80 disassembler. The annotated ROM
disassembly above is invaluable for cross-referencing CALL
targets and hardware port writes.
