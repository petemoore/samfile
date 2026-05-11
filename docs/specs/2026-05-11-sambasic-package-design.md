# SAM BASIC Program Construction — `sambasic` Package

Add a `sambasic` sub-package to samfile that provides native Go types for
constructing SAM BASIC programs and writing them to MGT disk images. Also
add `NewDiskImage()` for creating blank disks. The goal is to replace
hand-rolled Python/shell scripts (like `sam-aarch64/tools/build-disk.sh`)
with idiomatic Go using typed constants and compile-time checking.

## Background

samfile currently supports adding CODE files (`DiskImage.AddCodeFile`) and
detokenizing BASIC files to text (`SAMBasic.Output`), but has no way to
construct BASIC programs or create blank disk images. The
`sam-aarch64` project works around this with a 265-line bash/Python script
that hand-rolls directory entries, sector chains, and tokenized BASIC.

The `spectrum4/utils/zxbasic` package provides a proven API pattern for
constructing ZX Spectrum BASIC programs from Go types. SAM BASIC is
structurally similar (line-number + length + tokenized body + 0x0D
terminator, 0xFF end-of-program sentinel) but uses a different keyword
token encoding (single-byte 0x85-0xF6, two-byte 0xFF+index).

## Package: `sambasic`

New sub-package at `sambasic/` within the samfile module.

### Core types

```go
package sambasic

// File represents a complete SAM BASIC file as stored on disk.
// It encompasses both the in-memory program layout (PROG section
// + variable areas) and filesystem metadata (start line).
//
// The on-disk body (after the 9-byte FileHeader) is structured as:
//
//   PROG section:  tokenized lines + 0xFF sentinel
//   NumericVars:   (NUMEND - NVARS) bytes
//   Gap:           (SAVARS - NUMEND) bytes
//   StringArrayVars: remaining bytes
//
// When NumericVars is nil, the canonical 92 zero bytes are used.
// When Gap is nil, the canonical 512 zero bytes are used.
// When StringArrayVars is nil, it is empty.
// These defaults match the 94% empirical canonical case documented
// in sam-aarch64/docs/notes/sam-basic-save-format.md.
type File struct {
    Lines           []Line
    NumericVars     []byte   // nil -> 92 zero bytes
    Gap             []byte   // nil -> 512 zero bytes
    StringArrayVars []byte   // nil -> empty
    StartLine       uint16   // auto-RUN line; 0xFFFF = no auto-run
}

type Line struct {
    Number uint16
    Tokens []Token
}

// Token is implemented by SingleByteKeyword, TwoByteKeyword,
// Num, Str, and Literal.
type Token interface {
    Bytes() []byte
}

// SingleByteKeyword represents tokens in the range 0x85-0xF6.
type SingleByteKeyword byte

// TwoByteKeyword represents tokens encoded as 0xFF followed by
// an index byte. The constant value is the index byte.
type TwoByteKeyword byte

// Num is a numeric literal: display text + 0x0E + 5-byte value.
type Num struct {
    Display string
    Value   [5]byte
}

// Str is a literal string or character sequence.
type Str []byte

// Literal is a single raw byte (colon 0x3A, quote marks, etc).
type Literal byte
```

### Keyword constants

Every SAM BASIC v3 keyword becomes a named Go constant. Keywords are
split across `SingleByteKeyword` (direct tokens 0x85-0xF6) and
`TwoByteKeyword` (0xFF-prefixed tokens).

Examples:
```go
const (
    CLEAR SingleByteKeyword = 0xB3
    LOAD  SingleByteKeyword = 0x95
    CALL  SingleByteKeyword = 0xE4
    // ... all 0x85-0xF6 keywords
)

const (
    CODE TwoByteKeyword = 0x6C
    // ... all 0xFF-prefixed keywords
)
```

The full set is derived from the existing `keywords` string slice in
`keywords.go`, which maps token indices to display strings.

### Keyword-to-string mapping

A single bidirectional mapping exported from `sambasic` that provides:
- Token byte(s) -> display string (for detokenization)
- Display string -> token byte(s) (for future text-to-BASIC parsing)

The existing `keywords` string slice in `keywords.go` is replaced by
(or derived from) this mapping, so that `SAMBasic.Output()` and the
new tokenizer share the same source of truth.

### Serialization methods

**`(k SingleByteKeyword) Bytes() []byte`** — returns `[]byte{byte(k)}`.

**`(k TwoByteKeyword) Bytes() []byte`** — returns `[]byte{0xFF, byte(k)}`.

**`(n *Num) Bytes() []byte`** — returns display bytes + `0x0E` + 5 value
bytes (same pattern as zxbasic).

**`(s *Str) Bytes() []byte`** — returns the raw byte slice.

**`(l Literal) Bytes() []byte`** — returns `[]byte{byte(l)}`.

**`(l *Line) Bytes() []byte`** — concatenates all token bytes, appends
`0x0D` terminator, prepends 4-byte header (2-byte BE line number +
2-byte LE body length).

**`(f *File) ProgBytes() []byte`** — returns the PROG section: all line
bytes concatenated + `0xFF` end-of-program sentinel.

**`(f *File) Bytes() []byte`** — returns the complete file body (everything
after the 9-byte FileHeader): PROG section + numeric vars (defaulted if
nil) + gap (defaulted if nil) + string/array vars (defaulted if nil).

### Computed offset accessors

These return the byte offsets relative to the start of the file body,
needed by `DiskImage.AddBasicFile` to populate the directory entry's
three page-form triplets at offsets 0xDD/0xE0/0xE3.

**`(f *File) NVARSOffset() uint32`** — length of PROG section.

**`(f *File) NUMENDOffset() uint32`** — NVARSOffset + len(NumericVars)
(or 92 if nil).

**`(f *File) SAVARSOffset() uint32`** — NUMENDOffset + len(Gap) (or 512
if nil).

### Factory functions

**`Number(n uint16) *Num`** — creates a small-integer numeric literal
using the `[0x0E, 0x00, sign, lo, hi, 0x00]` encoding.

**`String(s string) *Str`** — wraps a Go string as a `Str` token.

### Example usage

```go
import "github.com/petemoore/samfile/v3/sambasic"

f := &sambasic.File{
    Lines: []sambasic.Line{
        {
            Number: 10,
            Tokens: []sambasic.Token{
                sambasic.CLEAR,
                sambasic.Number(32767),
                sambasic.Literal(':'),
                sambasic.LOAD,
                sambasic.Literal('"'),
                sambasic.String("stub"),
                sambasic.Literal('"'),
                sambasic.CODE,
                sambasic.Number(32768),
                sambasic.Literal(':'),
                sambasic.CALL,
                sambasic.Number(32768),
            },
        },
    },
    StartLine: 10,
}
```

## Package: `samfile` additions

### `NewDiskImage() *DiskImage`

Returns a pointer to a zeroed 819200-byte `DiskImage`. A blank MGT disk
is all zeros — no special formatting required.

### `DiskImage.AddBasicFile(name string, file *sambasic.File) error`

Writes a SAM BASIC file to the disk image, allocating a free directory
slot and required sectors. Parallel to the existing `AddCodeFile`.

Steps:
1. Call `file.Bytes()` to get the complete file body.
2. Construct the 9-byte body header:
   - Type: 0x10 (FT_SAM_BASIC)
   - LengthMod16K: `len(body) & 0x3FFF`
   - PageOffset: 0x9CD5 (PROG = 0x5CD5 in 8000H REL PAGE FORM)
   - Bytes 5-6: 0xFF, 0xFF (unused/no-exec marker)
   - Pages: `len(body) >> 14`
   - StartPage: 0
3. Create a `FileEntry` with:
   - Type = FT_SAM_BASIC
   - StartAddressPage = 0
   - StartAddressPageOffset = 0x9CD5
   - FileTypeInfo: three 3-byte page-form triplets encoding
     NVARSOffset, NUMENDOffset, SAVARSOffset
   - MGTFlags = 0x20
   - SAMBASICStartLine = file.StartLine (or 0xFFFF for no auto-run)
   - ExecutionAddressDiv16K = 0xFF when no auto-run, or 0x00 when
     StartLine is set
4. Call internal `addFile()` to allocate sectors and write.

## Harmonizing the keyword table

The `sambasic` package exports the keyword mapping. The existing
`SAMBasic.Output()` in the main `samfile` package imports `sambasic`
and uses it for keyword lookup, replacing the current `keywords` string
slice in `keywords.go`. This ensures both directions (construction and
detokenization) use a single source of truth.

## Testing strategy

### 1. Synthetic byte-match test

Build the `auto` program from `build-disk.sh` using Go types:
`10 CLEAR 32767: LOAD "stub" CODE 32768: CALL 32768`

Serialize via `File.Bytes()` and compare byte-for-byte against the
known-good output from `build-disk.sh`. This validates tokenization,
line framing, numeric form encoding, and trailer generation.

### 2. Exhaustive real disk roundtrip tests

Systematically loop through every `.mgt` file in
`~/Downloads/GoodSamC2/`, load each disk image, enumerate all files
via `DiskJournal`, and for every file with type `FT_SAM_BASIC`:

1. Extract raw body bytes via `DiskImage.File()`
2. Parse the body into a `sambasic.File` (requires a
   `Parse([]byte) (*File, error)` function — needed for roundtrip
   verification even if not a primary deliverable)
3. Re-serialize with `File.Bytes()`
4. Compare: serialized bytes must equal original body bytes

The test should report per-disk, per-file pass/fail and aggregate
statistics (total disks scanned, total BASIC files tested, pass/fail
counts). This gives comprehensive real-world coverage — the GoodSamC2
collection contains hundreds of disks with diverse BASIC programs.

### 3. Detokenizer consistency test

Build a program from Go types, serialize to bytes, pass through
`SAMBasic.Output()`, verify the text listing matches expected output.

### 4. Blank disk test

`NewDiskImage()` produces exactly 819200 zero bytes.

### 5. AddBasicFile integration test

Create blank disk, add a BASIC file via `AddBasicFile`, read it back
with `DiskImage.File()`, verify:
- Directory entry fields (type, name, sectors, start line, page-form
  triplets, MGTFlags)
- Body header bytes
- Body content matches `File.Bytes()` output
- File can be detokenized by `SAMBasic.Output()` without error

### 6. Multi-file integration test

Create blank disk, add samdos2 (CODE), auto (BASIC), stub (CODE) —
replicating the `build-disk.sh` layout. Verify all files coexist
without sector allocation conflicts.

## Out of scope

- Text-to-BASIC parsing (parsing `"10 CLEAR 32767"` into tokens)
- CLI commands (library API only)
- Floating-point display formatting beyond the raw 5-byte form
- Sophisticated variable-area types (just `[]byte` for now)
- MasterDOS 2156-byte trailer variant

## References

- `sam-aarch64/docs/notes/sam-basic-save-format.md` — canonical trailer
  recipe and empirical validation
- `sam-aarch64/docs/notes/sam-file-header.md` — 9-byte header and
  directory entry layout for type-16 files
- `sam-aarch64/tools/build-disk.sh` — reference Python implementation
- `spectrum4/utils/zxbasic/zxbasic.go` — API pattern for typed BASIC
  program construction
- SAM Coupe Technical Manual v3.0
- SAM ROM v3.0 annotated disassembly
