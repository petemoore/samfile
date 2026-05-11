# sambasic Package Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `sambasic` sub-package with Go types for constructing SAM BASIC programs, plus `NewDiskImage()` and `AddBasicFile()` on the existing `samfile.DiskImage`.

**Architecture:** New `sambasic/` package exports `File`, `Line`, `Token` interface, keyword constants, and serialization methods. The existing `samfile` package gains `NewDiskImage()` and `AddBasicFile()` that delegate to `sambasic.File.Bytes()` for body construction and the existing `addFile()` for sector allocation. The keyword table in `keywords.go` is replaced by an import from `sambasic`.

**Tech Stack:** Go 1.19+, no new dependencies. Tests use `testing` stdlib. Real-disk roundtrip tests scan `~/Downloads/GoodSamC2/*.dsk` (plus `*.mgt`).

---

## File Map

| File | Action | Responsibility |
|------|--------|----------------|
| `sambasic/tokens.go` | Create | Token interface, SingleByteKeyword, TwoByteKeyword, Num, Str, Literal types and their Bytes() methods; Number() and String() factories |
| `sambasic/keywords.go` | Create | All SAM BASIC v3 keyword constants; KeywordName() lookup function |
| `sambasic/file.go` | Create | File, Line types; ProgBytes(), Bytes(), NVARSOffset(), NUMENDOffset(), SAVARSOffset(), line serialization |
| `sambasic/tokens_test.go` | Create | Token serialization tests |
| `sambasic/keywords_test.go` | Create | Keyword constant spot-checks, KeywordName() coverage |
| `sambasic/file_test.go` | Create | Line/File serialization tests, build-disk.sh byte-match test, offset tests |
| `samfile.go` | Modify | Add NewDiskImage(), AddBasicFile() |
| `sambasic.go` | Modify | Import sambasic for keyword lookup, replace inline keywords slice |
| `keywords.go` | Delete | Replaced by sambasic/keywords.go |
| `samfile_test.go` | Modify | Add NewDiskImage test, AddBasicFile integration test, multi-file test |
| `sambasic/roundtrip_test.go` | Create | Exhaustive real-disk roundtrip tests over GoodSamC2 collection |
| `sambasic/parse.go` | Create | Parse() function: raw BASIC body bytes → sambasic.File (needed for roundtrip tests) |

---

### Task 1: Token Types and Serialization

**Files:**
- Create: `sambasic/tokens.go`
- Create: `sambasic/tokens_test.go`

- [ ] **Step 1: Write the failing tests for token serialization**

```go
// sambasic/tokens_test.go
package sambasic

import (
	"bytes"
	"testing"
)

func TestSingleByteKeywordBytes(t *testing.T) {
	k := SingleByteKeyword(0xB3) // CLEAR
	got := k.Bytes()
	want := []byte{0xB3}
	if !bytes.Equal(got, want) {
		t.Errorf("SingleByteKeyword(0xB3).Bytes() = %x; want %x", got, want)
	}
}

func TestTwoByteKeywordBytes(t *testing.T) {
	k := TwoByteKeyword(0x6C) // CODE
	got := k.Bytes()
	want := []byte{0xFF, 0x6C}
	if !bytes.Equal(got, want) {
		t.Errorf("TwoByteKeyword(0x6C).Bytes() = %x; want %x", got, want)
	}
}

func TestNumBytes(t *testing.T) {
	n := Number(32767)
	got := n.Bytes()
	// "32767" ASCII + 0x0E + [0x00, 0x00, 0xFF, 0x7F, 0x00]
	want := []byte{
		'3', '2', '7', '6', '7',
		0x0E, 0x00, 0x00, 0xFF, 0x7F, 0x00,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Number(32767).Bytes() = %x; want %x", got, want)
	}
}

func TestNumBytesZero(t *testing.T) {
	n := Number(0)
	got := n.Bytes()
	want := []byte{
		'0',
		0x0E, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Number(0).Bytes() = %x; want %x", got, want)
	}
}

func TestNumBytes32768(t *testing.T) {
	n := Number(32768)
	got := n.Bytes()
	want := []byte{
		'3', '2', '7', '6', '8',
		0x0E, 0x00, 0x00, 0x00, 0x80, 0x00,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Number(32768).Bytes() = %x; want %x", got, want)
	}
}

func TestStrBytes(t *testing.T) {
	s := String("stub")
	got := s.Bytes()
	want := []byte("stub")
	if !bytes.Equal(got, want) {
		t.Errorf("String(\"stub\").Bytes() = %x; want %x", got, want)
	}
}

func TestLiteralBytes(t *testing.T) {
	l := Literal(':')
	got := l.Bytes()
	want := []byte{0x3A}
	if !bytes.Equal(got, want) {
		t.Errorf("Literal(':').Bytes() = %x; want %x", got, want)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./sambasic/ -v -run 'TestSingleByte|TestTwoByte|TestNum|TestStr|TestLiteral'`
Expected: compilation errors — types not defined yet.

- [ ] **Step 3: Write minimal implementation**

```go
// sambasic/tokens.go
package sambasic

import "strconv"

type Token interface {
	Bytes() []byte
}

type SingleByteKeyword byte

func (k SingleByteKeyword) Bytes() []byte {
	return []byte{byte(k)}
}

type TwoByteKeyword byte

func (k TwoByteKeyword) Bytes() []byte {
	return []byte{0xFF, byte(k)}
}

type Num struct {
	Display string
	Value   [5]byte
}

func (n *Num) Bytes() []byte {
	result := []byte(n.Display)
	result = append(result, 0x0E, n.Value[0], n.Value[1], n.Value[2], n.Value[3], n.Value[4])
	return result
}

func Number(n uint16) *Num {
	return &Num{
		Display: strconv.Itoa(int(n)),
		Value:   [5]byte{0x00, 0x00, byte(n & 0xFF), byte(n >> 8), 0x00},
	}
}

type Str []byte

func (s *Str) Bytes() []byte {
	return []byte(*s)
}

func String(s string) *Str {
	v := Str(s)
	return &v
}

type Literal byte

func (l Literal) Bytes() []byte {
	return []byte{byte(l)}
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./sambasic/ -v -run 'TestSingleByte|TestTwoByte|TestNum|TestStr|TestLiteral'`
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add sambasic/tokens.go sambasic/tokens_test.go
git commit -m "feat(sambasic): add Token types and serialization"
```

---

### Task 2: Keyword Constants

**Files:**
- Create: `sambasic/keywords.go`
- Create: `sambasic/keywords_test.go`

- [ ] **Step 1: Write the failing tests**

```go
// sambasic/keywords_test.go
package sambasic

import (
	"bytes"
	"testing"
)

func TestKeywordConstants(t *testing.T) {
	cases := []struct {
		name  string
		token Token
		want  []byte
	}{
		{"PI", PI, []byte{0x85}},
		{"CLEAR", CLEAR, []byte{0xB3}},
		{"LOAD", LOAD, []byte{0x95}},
		{"CALL", CALL, []byte{0xE4}},
		{"IF_SHORT", IF_SHORT, []byte{0xF4}},
		{"ZOOM", ZOOM, []byte{0xF6}},
		{"CODE", CODE, []byte{0xFF, 0x6C}},
		{"UDG", UDG, []byte{0xFF, 0x63}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.token.Bytes()
			if !bytes.Equal(got, c.want) {
				t.Errorf("%s.Bytes() = %x; want %x", c.name, got, c.want)
			}
		})
	}
}

func TestKeywordName(t *testing.T) {
	cases := []struct {
		name     string
		token    byte
		extended bool
		want     string
		wantOK   bool
	}{
		{"CLEAR is 0xB3", 0xB3, false, "CLEAR", true},
		{"PI is 0x85", 0x85, false, "PI", true},
		{"CODE is 0xFF+0x6C", 0x6C, true, "CODE", true},
		{"below range", 0x20, false, "", false},
		{"above range", 0xF7, false, "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := KeywordName(c.token, c.extended)
			if ok != c.wantOK {
				t.Errorf("KeywordName(0x%02x, %v) ok = %v; want %v", c.token, c.extended, ok, c.wantOK)
			}
			if got != c.want {
				t.Errorf("KeywordName(0x%02x, %v) = %q; want %q", c.token, c.extended, got, c.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./sambasic/ -v -run 'TestKeyword'`
Expected: compilation errors — constants and KeywordName not defined.

- [ ] **Step 3: Write the keyword constants and lookup function**

Create `sambasic/keywords.go` with all SAM BASIC v3 keyword constants. The constants are derived from the existing `keywords` string slice in the root `keywords.go` file. Each entry at index `i` in that slice corresponds to token byte `0x3B + i`. Tokens 0x85–0xF6 are `SingleByteKeyword`; tokens reached via the 0xFF escape (where the escape-byte index is also `0x3B + i`) are `TwoByteKeyword`.

The existing `keywords` slice has 188 entries (indices 0–187). Token bytes range from 0x3B+0=0x3B to 0x3B+187=0xF6. Tokens below 0x85 (indices 0–73, i.e. 0x3B–0x84) are only reachable via the two-byte 0xFF escape. Tokens 0x85–0xF6 (indices 74–187) are single-byte.

```go
// sambasic/keywords.go
package sambasic

// Single-byte keywords: tokens 0x85–0xF6.
const (
	PI          SingleByteKeyword = 0x85
	RND         SingleByteKeyword = 0x86
	POINT       SingleByteKeyword = 0x87
	FREE        SingleByteKeyword = 0x88
	LENGTH      SingleByteKeyword = 0x89
	ITEM        SingleByteKeyword = 0x8A
	ATTR        SingleByteKeyword = 0x8B
	FN          SingleByteKeyword = 0x8C
	BIN         SingleByteKeyword = 0x8D
	XMOUSE      SingleByteKeyword = 0x8E
	YHOUSE      SingleByteKeyword = 0x8F
	XPEN        SingleByteKeyword = 0x90
	YPEN        SingleByteKeyword = 0x91
	RAMTOP      SingleByteKeyword = 0x92
	// 0x93 reserved
	INSTR       SingleByteKeyword = 0x94
	LOAD        SingleByteKeyword = 0x95
	SCREEN_STR  SingleByteKeyword = 0x96
	MEM_STR     SingleByteKeyword = 0x97
	// 0x98 reserved
	PATH_STR    SingleByteKeyword = 0x99
	STRING_STR  SingleByteKeyword = 0x9A
	// 0x9B, 0x9C reserved
	SIN         SingleByteKeyword = 0x9D
	COS         SingleByteKeyword = 0x9E
	TAN         SingleByteKeyword = 0x9F
	ASN         SingleByteKeyword = 0xA0
	ACS         SingleByteKeyword = 0xA1
	ATN         SingleByteKeyword = 0xA2
	LN          SingleByteKeyword = 0xA3
	EXP         SingleByteKeyword = 0xA4
	ABS         SingleByteKeyword = 0xA5
	SGN         SingleByteKeyword = 0xA6
	SQR         SingleByteKeyword = 0xA7
	INT_KW      SingleByteKeyword = 0xA8
	USR         SingleByteKeyword = 0xA9
	IN_KW       SingleByteKeyword = 0xAA
	PEEK        SingleByteKeyword = 0xAB
	LPEEK       SingleByteKeyword = 0xAC
	DVAR        SingleByteKeyword = 0xAD
	SVAR        SingleByteKeyword = 0xAE
	BUTTON      SingleByteKeyword = 0xAF
	EOF_KW      SingleByteKeyword = 0xB0
	PTR         SingleByteKeyword = 0xB1
	// 0xB2 reserved
	CLEAR       SingleByteKeyword = 0xB3
	// 0xB4 reserved
	LEN_KW      SingleByteKeyword = 0xB5
	// ... continue for all through 0xF6
	// (Full listing elided here — implementer must derive from
	// the keywords slice, mapping index 74..187 to 0x85..0xF6)
)
```

**IMPORTANT:** The above is illustrative. The implementer MUST carefully map every entry in the existing `keywords` slice (in root `keywords.go`) to the correct byte value. The mapping is: keywords slice index `i` → token byte `0x3B + i`. If `0x3B + i >= 0x85`, it's a `SingleByteKeyword(0x3B + i)`. If `0x3B + i < 0x85`, the keyword is only reachable via the 0xFF escape as `TwoByteKeyword(0x3B + i)`.

Here is the complete mapping from the existing `keywords` slice. Each keyword's index `i` and token byte `0x3B + i`:

| Index | Byte | Name | Type |
|-------|------|------|------|
| 0 | 0x3B | PI | TwoByteKeyword |
| 1 | 0x3C | RND | TwoByteKeyword |
| 2 | 0x3D | POINT | TwoByteKeyword |
| 3 | 0x3E | FREE | TwoByteKeyword |
| 4 | 0x3F | LENGTH | TwoByteKeyword |
| 5 | 0x40 | ITEM | TwoByteKeyword |
| 6 | 0x41 | ATTR | TwoByteKeyword |
| 7 | 0x42 | FN | TwoByteKeyword |
| 8 | 0x43 | BIN | TwoByteKeyword |
| 9 | 0x44 | XMOUSE | TwoByteKeyword |
| 10 | 0x45 | YHOUSE | TwoByteKeyword |
| 11 | 0x46 | XPEN | TwoByteKeyword |
| 12 | 0x47 | YPEN | TwoByteKeyword |
| 13 | 0x48 | RAMTOP | TwoByteKeyword |
| 14 | 0x49 | <<Reserved>> | — |
| 15 | 0x4A | INSTR | TwoByteKeyword |
| 16 | 0x4B | INKEY$ | TwoByteKeyword |
| 17 | 0x4C | SCREEN$ | TwoByteKeyword |
| 18 | 0x4D | MEM$ | TwoByteKeyword |
| 19 | 0x4E | <<Reserved>> | — |
| 20 | 0x4F | PATH$ | TwoByteKeyword |
| 21 | 0x50 | STRING$ | TwoByteKeyword |
| 22–23 | 0x51–0x52 | <<Reserved>> | — |
| 24 | 0x53 | SIN | TwoByteKeyword |
| 25 | 0x54 | COS | TwoByteKeyword |
| 26 | 0x55 | TAN | TwoByteKeyword |
| 27 | 0x56 | ASN | TwoByteKeyword |
| 28 | 0x57 | ACS | TwoByteKeyword |
| 29 | 0x58 | ATN | TwoByteKeyword |
| 30 | 0x59 | LN | TwoByteKeyword |
| 31 | 0x5A | EXP | TwoByteKeyword |
| 32 | 0x5B | ABS | TwoByteKeyword |
| 33 | 0x5C | SGN | TwoByteKeyword |
| 34 | 0x5D | SQR | TwoByteKeyword |
| 35 | 0x5E | INT | TwoByteKeyword |
| 36 | 0x5F | USR | TwoByteKeyword |
| 37 | 0x60 | IN | TwoByteKeyword |
| 38 | 0x61 | PEEK | TwoByteKeyword |
| 39 | 0x62 | LPEEK | TwoByteKeyword |
| 40 | 0x63 | DVAR | TwoByteKeyword |
| 41 | 0x64 | SVAR | TwoByteKeyword |
| 42 | 0x65 | BUTTON | TwoByteKeyword |
| 43 | 0x66 | EOF | TwoByteKeyword |
| 44 | 0x67 | PTR | TwoByteKeyword |
| 45 | 0x68 | <<Reserved>> | — |
| 46 | 0x69 | UDG | TwoByteKeyword |
| 47 | 0x6A | <<Reserved>> | — |
| 48 | 0x6B | LEN | TwoByteKeyword |
| 49 | 0x6C | CODE | TwoByteKeyword |
| 50 | 0x6D | VAL$ | TwoByteKeyword |
| 51 | 0x6E | VAL | TwoByteKeyword |
| 52 | 0x6F | TRUNC$ | TwoByteKeyword |
| 53 | 0x70 | CHR$ | TwoByteKeyword |
| 54 | 0x71 | STR$ | TwoByteKeyword |
| 55 | 0x72 | BIN$ | TwoByteKeyword |
| 56 | 0x73 | HEX$ | TwoByteKeyword |
| 57 | 0x74 | USR$ | TwoByteKeyword |
| 58 | 0x75 | <<Reserved>> | — |
| 59 | 0x76 | NOT | TwoByteKeyword |
| 60–62 | 0x77–0x79 | <<Reserved>> | — |
| 63 | 0x7A | MOD | TwoByteKeyword |
| 64 | 0x7B | DIV | TwoByteKeyword |
| 65 | 0x7C | BOR | TwoByteKeyword |
| 66 | 0x7D | <<Reserved>> | — |
| 67 | 0x7E | BAND | TwoByteKeyword |
| 68 | 0x7F | OR | TwoByteKeyword |
| 69 | 0x80 | AND | TwoByteKeyword |
| 70 | 0x81 | <> | TwoByteKeyword |
| 71 | 0x82 | <= | TwoByteKeyword |
| 72 | 0x83 | >= | TwoByteKeyword |
| 73 | 0x84 | <<Reserved>> | — |

From index 74 onwards (byte 0x85+), these are SingleByteKeyword:

| Index | Byte | Name |
|-------|------|------|
| 74 | 0x85 | USING |
| 75 | 0x86 | WRITE |
| 76 | 0x87 | AT |
| 77 | 0x88 | TAB |
| 78 | 0x89 | OFF |
| 79 | 0x8A | WHILE |
| 80 | 0x8B | UNTIL |
| 81 | 0x8C | LINE |
| 82 | 0x8D | THEN |
| 83 | 0x8E | TO |
| 84 | 0x8F | STEP |
| 85 | 0x90 | DIR |
| 86 | 0x91 | FORMAT |
| 87 | 0x92 | ERASE |
| 88 | 0x93 | MOVE |
| 89 | 0x94 | SAVE |
| 90 | 0x95 | LOAD |
| 91 | 0x96 | MERGE |
| 92 | 0x97 | VERIFY |
| 93 | 0x98 | OPEN |
| 94 | 0x99 | CLOSE |
| 95 | 0x9A | CIRCLE |
| 96 | 0x9B | PLOT |
| 97 | 0x9C | LET |
| 98 | 0x9D | BLITZ |
| 99 | 0x9E | BORDER |
| 100 | 0x9F | CLS |
| 101 | 0xA0 | PALETTE |
| 102 | 0xA1 | PEN |
| 103 | 0xA2 | PAPER |
| 104 | 0xA3 | FLASH |
| 105 | 0xA4 | BRIGHT |
| 106 | 0xA5 | INVERSE |
| 107 | 0xA6 | OVER |
| 108 | 0xA7 | FATPIX |
| 109 | 0xA8 | CSIZE |
| 110 | 0xA9 | BLOCKS |
| 111 | 0xAA | MODE |
| 112 | 0xAB | GRAB |
| 113 | 0xAC | PUT |
| 114 | 0xAD | BEEP |
| 115 | 0xAE | SOUND |
| 116 | 0xAF | NEW |
| 117 | 0xB0 | RUN |
| 118 | 0xB1 | STOP |
| 119 | 0xB2 | CONTINUE |
| 120 | 0xB3 | CLEAR |
| 121 | 0xB4 | GO TO |
| 122 | 0xB5 | GO SUB |
| 123 | 0xB6 | RETURN |
| 124 | 0xB7 | REM |
| 125 | 0xB8 | READ |
| 126 | 0xB9 | DATA |
| 127 | 0xBA | RESTORE |
| 128 | 0xBB | PRINT |
| 129 | 0xBC | LPRINT |
| 130 | 0xBD | LIST |
| 131 | 0xBE | LLIST |
| 132 | 0xBF | DUMP |
| 133 | 0xC0 | FOR |
| 134 | 0xC1 | NEXT |
| 135 | 0xC2 | PAUSE |
| 136 | 0xC3 | DRAW |
| 137 | 0xC4 | DEFAULT |
| 138 | 0xC5 | DIM |
| 139 | 0xC6 | INPUT |
| 140 | 0xC7 | RANDOMIZE |
| 141 | 0xC8 | DEF FN |
| 142 | 0xC9 | DEF KEYCODE |
| 143 | 0xCA | DEF PROC |
| 144 | 0xCB | END PROC |
| 145 | 0xCC | RENUM |
| 146 | 0xCD | DELETE |
| 147 | 0xCE | REF |
| 148 | 0xCF | COPY |
| 149 | 0xD0 | <<Reserved>> |
| 150 | 0xD1 | KEYIN |
| 151 | 0xD2 | LOCAL |
| 152 | 0xD3 | LOOP IF |
| 153 | 0xD4 | DO |
| 154 | 0xD5 | LOOP |
| 155 | 0xD6 | EXIT IF |
| 156 | 0xD7 | IF (long) |
| 157 | 0xD8 | IF (short) |
| 158 | 0xD9 | ELSE (long) |
| 159 | 0xDA | ELSE (short) |
| 160 | 0xDB | END IF |
| 161 | 0xDC | KEY |
| 162 | 0xDD | ON ERROR |
| 163 | 0xDE | ON |
| 164 | 0xDF | GET |
| 165 | 0xE0 | OUT |
| 166 | 0xE1 | POKE |
| 167 | 0xE2 | DPOKE |
| 168 | 0xE3 | RENAME |
| 169 | 0xE4 | CALL |
| 170 | 0xE5 | ROLL |
| 171 | 0xE6 | SCROLL |
| 172 | 0xE7 | SCREEN |
| 173 | 0xE8 | DISPLAY |
| 174 | 0xE9 | BOOT |
| 175 | 0xEA | LABEL |
| 176 | 0xEB | FILL |
| 177 | 0xEC | WINDOW |
| 178 | 0xED | AUTO |
| 179 | 0xEE | POP |
| 180 | 0xEF | RECORD |
| 181 | 0xF0 | DEVICE |
| 182 | 0xF1 | PROTECT |
| 183 | 0xF2 | HIDE |
| 184 | 0xF3 | ZAP |
| 185 | 0xF4 | POW |
| 186 | 0xF5 | BOOM |
| 187 | 0xF6 | ZOOM |

**CRITICAL CROSS-CHECK:** The existing `keywords` slice has two entries each for IF and ELSE (long/short forms, indices 156/157 and 158/159). Both pairs share the same display string but differ in token byte. Name the Go constants `IF_LONG`/`IF_SHORT` and `ELSE_LONG`/`ELSE_SHORT`.

Also note: the first 74 entries (index 0–73, bytes 0x3B–0x84) exist in the original keywords table and ARE used during detokenization of both single-byte tokens 0x85+ AND two-byte 0xFF-escaped tokens. When the detokenizer sees byte 0x85, it does `keywords[0x85 - 0x3B]` = `keywords[74]` = `"USING"`. When it sees 0xFF followed by 0x3B, it does `keywords[0x3B - 0x3B]` = `keywords[0]` = `"PI"`. The same table serves both paths.

The `KeywordName(tokenByte byte, extended bool) (string, bool)` function must replicate this: for `extended=false` it looks up index `tokenByte - 0x3B` (valid for 0x85 ≤ tokenByte ≤ 0xF6); for `extended=true` it looks up index `tokenByte - 0x3B` (valid for tokenByte ≥ 0x3B). Returns `("", false)` for reserved or out-of-range entries.

```go
// sambasic/keywords.go
package sambasic

// keywordTable is the SAM BASIC v3 keyword table, indexed by
// (token_byte - 0x3B). Used by both single-byte (0x85-0xF6) and
// two-byte (0xFF + idx) token forms. Empty string means reserved.
var keywordTable = [...]string{
	"PI",            // 0x3B  (index 0)
	"RND",           // 0x3C  (index 1)
	"POINT",         // 0x3D
	"FREE",          // 0x3E
	"LENGTH",        // 0x3F
	"ITEM",          // 0x40
	"ATTR",          // 0x41
	"FN",            // 0x42
	"BIN",           // 0x43
	"XMOUSE",        // 0x44
	"YHOUSE",        // 0x45
	"XPEN",          // 0x46
	"YPEN",          // 0x47
	"RAMTOP",        // 0x48
	"",              // 0x49  reserved
	"INSTR",         // 0x4A
	"INKEY$",        // 0x4B
	"SCREEN$",       // 0x4C
	"MEM$",          // 0x4D
	"",              // 0x4E  reserved
	"PATH$",         // 0x4F
	"STRING$",       // 0x50
	"",              // 0x51  reserved
	"",              // 0x52  reserved
	"SIN",           // 0x53
	"COS",           // 0x54
	"TAN",           // 0x55
	"ASN",           // 0x56
	"ACS",           // 0x57
	"ATN",           // 0x58
	"LN",            // 0x59
	"EXP",           // 0x5A
	"ABS",           // 0x5B
	"SGN",           // 0x5C
	"SQR",           // 0x5D
	"INT",           // 0x5E
	"USR",           // 0x5F
	"IN",            // 0x60
	"PEEK",          // 0x61
	"LPEEK",         // 0x62
	"DVAR",          // 0x63
	"SVAR",          // 0x64
	"BUTTON",        // 0x65
	"EOF",           // 0x66
	"PTR",           // 0x67
	"",              // 0x68  reserved
	"UDG",           // 0x69
	"",              // 0x6A  reserved
	"LEN",           // 0x6B
	"CODE",          // 0x6C
	"VAL$",          // 0x6D
	"VAL",           // 0x6E
	"TRUNC$",        // 0x6F
	"CHR$",          // 0x70
	"STR$",          // 0x71
	"BIN$",          // 0x72
	"HEX$",          // 0x73
	"USR$",          // 0x74
	"",              // 0x75  reserved
	"NOT",           // 0x76
	"",              // 0x77  reserved
	"",              // 0x78  reserved
	"",              // 0x79  reserved
	"MOD",           // 0x7A
	"DIV",           // 0x7B
	"BOR",           // 0x7C
	"",              // 0x7D  reserved
	"BAND",          // 0x7E
	"OR",            // 0x7F
	"AND",           // 0x80
	"<>",            // 0x81
	"<=",            // 0x82
	">=",            // 0x83
	"",              // 0x84  reserved
	"USING",         // 0x85  (index 74, first SingleByteKeyword)
	"WRITE",         // 0x86
	"AT",            // 0x87
	"TAB",           // 0x88
	"OFF",           // 0x89
	"WHILE",         // 0x8A
	"UNTIL",         // 0x8B
	"LINE",          // 0x8C
	"THEN",          // 0x8D
	"TO",            // 0x8E
	"STEP",          // 0x8F
	"DIR",           // 0x90
	"FORMAT",        // 0x91
	"ERASE",         // 0x92
	"MOVE",          // 0x93
	"SAVE",          // 0x94
	"LOAD",          // 0x95
	"MERGE",         // 0x96
	"VERIFY",        // 0x97
	"OPEN",          // 0x98
	"CLOSE",         // 0x99
	"CIRCLE",        // 0x9A
	"PLOT",          // 0x9B
	"LET",           // 0x9C
	"BLITZ",         // 0x9D
	"BORDER",        // 0x9E
	"CLS",           // 0x9F
	"PALETTE",       // 0xA0
	"PEN",           // 0xA1
	"PAPER",         // 0xA2
	"FLASH",         // 0xA3
	"BRIGHT",        // 0xA4
	"INVERSE",       // 0xA5
	"OVER",          // 0xA6
	"FATPIX",        // 0xA7
	"CSIZE",         // 0xA8
	"BLOCKS",        // 0xA9
	"MODE",          // 0xAA
	"GRAB",          // 0xAB
	"PUT",           // 0xAC
	"BEEP",          // 0xAD
	"SOUND",         // 0xAE
	"NEW",           // 0xAF
	"RUN",           // 0xB0
	"STOP",          // 0xB1
	"CONTINUE",      // 0xB2
	"CLEAR",         // 0xB3
	"GO TO",         // 0xB4
	"GO SUB",        // 0xB5
	"RETURN",        // 0xB6
	"REM",           // 0xB7
	"READ",          // 0xB8
	"DATA",          // 0xB9
	"RESTORE",       // 0xBA
	"PRINT",         // 0xBB
	"LPRINT",        // 0xBC
	"LIST",          // 0xBD
	"LLIST",         // 0xBE
	"DUMP",          // 0xBF
	"FOR",           // 0xC0
	"NEXT",          // 0xC1
	"PAUSE",         // 0xC2
	"DRAW",          // 0xC3
	"DEFAULT",       // 0xC4
	"DIM",           // 0xC5
	"INPUT",         // 0xC6
	"RANDOMIZE",     // 0xC7
	"DEF FN",        // 0xC8
	"DEF KEYCODE",   // 0xC9
	"DEF PROC",      // 0xCA
	"END PROC",      // 0xCB
	"RENUM",         // 0xCC
	"DELETE",        // 0xCD
	"REF",           // 0xCE
	"COPY",          // 0xCF
	"",              // 0xD0  reserved
	"KEYIN",         // 0xD1
	"LOCAL",         // 0xD2
	"LOOP IF",       // 0xD3
	"DO",            // 0xD4
	"LOOP",          // 0xD5
	"EXIT IF",       // 0xD6
	"IF",            // 0xD7  long IF
	"IF",            // 0xD8  short IF
	"ELSE",          // 0xD9  long ELSE
	"ELSE",          // 0xDA  short ELSE
	"END IF",        // 0xDB
	"KEY",           // 0xDC
	"ON ERROR",      // 0xDD
	"ON",            // 0xDE
	"GET",           // 0xDF
	"OUT",           // 0xE0
	"POKE",          // 0xE1
	"DPOKE",         // 0xE2
	"RENAME",        // 0xE3
	"CALL",          // 0xE4
	"ROLL",          // 0xE5
	"SCROLL",        // 0xE6
	"SCREEN",        // 0xE7
	"DISPLAY",       // 0xE8
	"BOOT",          // 0xE9
	"LABEL",         // 0xEA
	"FILL",          // 0xEB
	"WINDOW",        // 0xEC
	"AUTO",          // 0xED
	"POP",           // 0xEE
	"RECORD",        // 0xEF
	"DEVICE",        // 0xF0
	"PROTECT",       // 0xF1
	"HIDE",          // 0xF2
	"ZAP",           // 0xF3
	"POW",           // 0xF4
	"BOOM",          // 0xF5
	"ZOOM",          // 0xF6
}

// SingleByteKeyword constants (0x85–0xF6).
const (
	USING      SingleByteKeyword = 0x85
	WRITE      SingleByteKeyword = 0x86
	AT         SingleByteKeyword = 0x87
	TAB        SingleByteKeyword = 0x88
	OFF        SingleByteKeyword = 0x89
	WHILE      SingleByteKeyword = 0x8A
	UNTIL      SingleByteKeyword = 0x8B
	LINE       SingleByteKeyword = 0x8C
	THEN       SingleByteKeyword = 0x8D
	TO         SingleByteKeyword = 0x8E
	STEP       SingleByteKeyword = 0x8F
	DIR        SingleByteKeyword = 0x90
	FORMAT     SingleByteKeyword = 0x91
	ERASE      SingleByteKeyword = 0x92
	MOVE       SingleByteKeyword = 0x93
	SAVE       SingleByteKeyword = 0x94
	LOAD       SingleByteKeyword = 0x95
	MERGE      SingleByteKeyword = 0x96
	VERIFY     SingleByteKeyword = 0x97
	OPEN       SingleByteKeyword = 0x98
	CLOSE      SingleByteKeyword = 0x99
	CIRCLE     SingleByteKeyword = 0x9A
	PLOT       SingleByteKeyword = 0x9B
	LET        SingleByteKeyword = 0x9C
	BLITZ      SingleByteKeyword = 0x9D
	BORDER     SingleByteKeyword = 0x9E
	CLS        SingleByteKeyword = 0x9F
	PALETTE    SingleByteKeyword = 0xA0
	PEN        SingleByteKeyword = 0xA1
	PAPER      SingleByteKeyword = 0xA2
	FLASH      SingleByteKeyword = 0xA3
	BRIGHT     SingleByteKeyword = 0xA4
	INVERSE    SingleByteKeyword = 0xA5
	OVER       SingleByteKeyword = 0xA6
	FATPIX     SingleByteKeyword = 0xA7
	CSIZE      SingleByteKeyword = 0xA8
	BLOCKS     SingleByteKeyword = 0xA9
	MODE       SingleByteKeyword = 0xAA
	GRAB       SingleByteKeyword = 0xAB
	PUT        SingleByteKeyword = 0xAC
	BEEP       SingleByteKeyword = 0xAD
	SOUND      SingleByteKeyword = 0xAE
	NEW        SingleByteKeyword = 0xAF
	RUN        SingleByteKeyword = 0xB0
	STOP       SingleByteKeyword = 0xB1
	CONTINUE   SingleByteKeyword = 0xB2
	CLEAR      SingleByteKeyword = 0xB3
	GO_TO      SingleByteKeyword = 0xB4
	GO_SUB     SingleByteKeyword = 0xB5
	RETURN     SingleByteKeyword = 0xB6
	REM        SingleByteKeyword = 0xB7
	READ       SingleByteKeyword = 0xB8
	DATA       SingleByteKeyword = 0xB9
	RESTORE    SingleByteKeyword = 0xBA
	PRINT      SingleByteKeyword = 0xBB
	LPRINT     SingleByteKeyword = 0xBC
	LIST       SingleByteKeyword = 0xBD
	LLIST      SingleByteKeyword = 0xBE
	DUMP       SingleByteKeyword = 0xBF
	FOR        SingleByteKeyword = 0xC0
	NEXT       SingleByteKeyword = 0xC1
	PAUSE      SingleByteKeyword = 0xC2
	DRAW       SingleByteKeyword = 0xC3
	DEFAULT    SingleByteKeyword = 0xC4
	DIM        SingleByteKeyword = 0xC5
	INPUT      SingleByteKeyword = 0xC6
	RANDOMIZE  SingleByteKeyword = 0xC7
	DEF_FN     SingleByteKeyword = 0xC8
	DEF_KEYCODE SingleByteKeyword = 0xC9
	DEF_PROC   SingleByteKeyword = 0xCA
	END_PROC   SingleByteKeyword = 0xCB
	RENUM      SingleByteKeyword = 0xCC
	DELETE     SingleByteKeyword = 0xCD
	REF        SingleByteKeyword = 0xCE
	COPY       SingleByteKeyword = 0xCF
	// 0xD0 reserved
	KEYIN      SingleByteKeyword = 0xD1
	LOCAL      SingleByteKeyword = 0xD2
	LOOP_IF    SingleByteKeyword = 0xD3
	DO         SingleByteKeyword = 0xD4
	LOOP       SingleByteKeyword = 0xD5
	EXIT_IF    SingleByteKeyword = 0xD6
	IF_LONG    SingleByteKeyword = 0xD7
	IF_SHORT   SingleByteKeyword = 0xD8
	ELSE_LONG  SingleByteKeyword = 0xD9
	ELSE_SHORT SingleByteKeyword = 0xDA
	END_IF     SingleByteKeyword = 0xDB
	KEY        SingleByteKeyword = 0xDC
	ON_ERROR   SingleByteKeyword = 0xDD
	ON         SingleByteKeyword = 0xDE
	GET        SingleByteKeyword = 0xDF
	OUT        SingleByteKeyword = 0xE0
	POKE       SingleByteKeyword = 0xE1
	DPOKE      SingleByteKeyword = 0xE2
	RENAME     SingleByteKeyword = 0xE3
	CALL       SingleByteKeyword = 0xE4
	ROLL       SingleByteKeyword = 0xE5
	SCROLL     SingleByteKeyword = 0xE6
	SCREEN     SingleByteKeyword = 0xE7
	DISPLAY    SingleByteKeyword = 0xE8
	BOOT       SingleByteKeyword = 0xE9
	LABEL      SingleByteKeyword = 0xEA
	FILL       SingleByteKeyword = 0xEB
	WINDOW     SingleByteKeyword = 0xEC
	AUTO_KW    SingleByteKeyword = 0xED
	POP        SingleByteKeyword = 0xEE
	RECORD     SingleByteKeyword = 0xEF
	DEVICE     SingleByteKeyword = 0xF0
	PROTECT    SingleByteKeyword = 0xF1
	HIDE       SingleByteKeyword = 0xF2
	ZAP        SingleByteKeyword = 0xF3
	POW        SingleByteKeyword = 0xF4
	BOOM       SingleByteKeyword = 0xF5
	ZOOM       SingleByteKeyword = 0xF6
)

// TwoByteKeyword constants (0xFF + byte).
const (
	PI_2B      TwoByteKeyword = 0x3B
	RND_2B     TwoByteKeyword = 0x3C
	POINT_2B   TwoByteKeyword = 0x3D
	FREE_2B    TwoByteKeyword = 0x3E
	LENGTH_2B  TwoByteKeyword = 0x3F
	ITEM_2B    TwoByteKeyword = 0x40
	ATTR_2B    TwoByteKeyword = 0x41
	FN_2B      TwoByteKeyword = 0x42
	BIN_2B     TwoByteKeyword = 0x43
	XMOUSE_2B  TwoByteKeyword = 0x44
	YHOUSE_2B  TwoByteKeyword = 0x45
	XPEN_2B    TwoByteKeyword = 0x46
	YPEN_2B    TwoByteKeyword = 0x47
	RAMTOP_2B  TwoByteKeyword = 0x48
	INSTR_2B   TwoByteKeyword = 0x4A
	INKEY_2B   TwoByteKeyword = 0x4B
	SCREEN_2B  TwoByteKeyword = 0x4C
	MEM_2B     TwoByteKeyword = 0x4D
	PATH_2B    TwoByteKeyword = 0x4F
	STRING_2B  TwoByteKeyword = 0x50
	SIN_2B     TwoByteKeyword = 0x53
	COS_2B     TwoByteKeyword = 0x54
	TAN_2B     TwoByteKeyword = 0x55
	ASN_2B     TwoByteKeyword = 0x56
	ACS_2B     TwoByteKeyword = 0x57
	ATN_2B     TwoByteKeyword = 0x58
	LN_2B      TwoByteKeyword = 0x59
	EXP_2B     TwoByteKeyword = 0x5A
	ABS_2B     TwoByteKeyword = 0x5B
	SGN_2B     TwoByteKeyword = 0x5C
	SQR_2B     TwoByteKeyword = 0x5D
	INT_2B     TwoByteKeyword = 0x5E
	USR_2B     TwoByteKeyword = 0x5F
	IN_2B      TwoByteKeyword = 0x60
	PEEK_2B    TwoByteKeyword = 0x61
	LPEEK_2B   TwoByteKeyword = 0x62
	DVAR_2B    TwoByteKeyword = 0x63
	SVAR_2B    TwoByteKeyword = 0x64
	BUTTON_2B  TwoByteKeyword = 0x65
	EOF_2B     TwoByteKeyword = 0x66
	PTR_2B     TwoByteKeyword = 0x67
	UDG        TwoByteKeyword = 0x69
	LEN_2B     TwoByteKeyword = 0x6B
	CODE       TwoByteKeyword = 0x6C
	VAL_STR    TwoByteKeyword = 0x6D
	VAL_2B     TwoByteKeyword = 0x6E
	TRUNC_STR  TwoByteKeyword = 0x6F
	CHR_STR    TwoByteKeyword = 0x70
	STR_STR    TwoByteKeyword = 0x71
	BIN_STR    TwoByteKeyword = 0x72
	HEX_STR    TwoByteKeyword = 0x73
	USR_STR    TwoByteKeyword = 0x74
	NOT_2B     TwoByteKeyword = 0x76
	MOD_2B     TwoByteKeyword = 0x7A
	DIV_2B     TwoByteKeyword = 0x7B
	BOR_2B     TwoByteKeyword = 0x7C
	BAND_2B    TwoByteKeyword = 0x7E
	OR_2B      TwoByteKeyword = 0x7F
	AND_2B     TwoByteKeyword = 0x80
	NOT_EQUAL  TwoByteKeyword = 0x81
	LESS_EQUAL TwoByteKeyword = 0x82
	GR_EQUAL   TwoByteKeyword = 0x83
)

// KeywordName returns the display string for a keyword token.
// If extended is false, tokenByte is a single-byte token (0x85–0xF6).
// If extended is true, tokenByte is the byte after the 0xFF escape.
// Returns ("", false) for reserved or out-of-range tokens.
func KeywordName(tokenByte byte, extended bool) (string, bool) {
	if tokenByte < 0x3B {
		return "", false
	}
	idx := int(tokenByte - 0x3B)
	if idx >= len(keywordTable) {
		return "", false
	}
	if !extended && tokenByte < 0x85 {
		return "", false
	}
	name := keywordTable[idx]
	if name == "" {
		return "", false
	}
	return name, true
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./sambasic/ -v -run 'TestKeyword'`
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add sambasic/keywords.go sambasic/keywords_test.go
git commit -m "feat(sambasic): add keyword constants and lookup"
```

---

### Task 3: Line and File Serialization

**Files:**
- Create: `sambasic/file.go`
- Create: `sambasic/file_test.go`

- [ ] **Step 1: Write the failing tests**

```go
// sambasic/file_test.go
package sambasic

import (
	"bytes"
	"testing"
)

func TestLineBytes(t *testing.T) {
	// Line 10 with just "CLEAR 32767"
	line := Line{
		Number: 10,
		Tokens: []Token{
			CLEAR,
			Number(32767),
		},
	}
	got := line.Bytes()
	want := []byte{
		0x00, 0x0a, // line number 10 big-endian
		0x0d, 0x00, // body length 13 little-endian (1 + 5+6 + 1 = 13)
		0xb3,                               // CLEAR
		'3', '2', '7', '6', '7',           // display "32767"
		0x0e, 0x00, 0x00, 0xff, 0x7f, 0x00, // numeric form
		0x0d, // line terminator
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Line.Bytes():\n  got  %x\n  want %x", got, want)
	}
}

func TestBuildDiskAutoProgram(t *testing.T) {
	// Reproduce exactly: 10 CLEAR 32767: LOAD "stub" CODE 32768: CALL 32768
	f := &File{
		Lines: []Line{
			{
				Number: 10,
				Tokens: []Token{
					CLEAR,
					String("32767"),
					Number(32767),
					Literal(':'),
					LOAD,
					Literal('"'),
					String("stub"),
					Literal('"'),
					CODE,
					String("32768"),
					Number(32768),
					Literal(':'),
					CALL,
					String("32768"),
					Number(32768),
				},
			},
		},
		StartLine: 10,
	}

	// Expected PROG section from build-disk.sh:
	wantProg := []byte{
		0x00, 0x0a, 0x2f, 0x00, // line 10, body len 47
		0xb3,                                           // CLEAR
		'3', '2', '7', '6', '7',                       // "32767"
		0x0e, 0x00, 0x00, 0xff, 0x7f, 0x00,            // num(32767)
		0x3a,                                           // :
		0x95,                                           // LOAD
		0x22,                                           // "
		's', 't', 'u', 'b',                            // stub
		0x22,                                           // "
		0xff, 0x6c,                                     // CODE (two-byte)
		'3', '2', '7', '6', '8',                       // "32768"
		0x0e, 0x00, 0x00, 0x00, 0x80, 0x00,            // num(32768)
		0x3a,                                           // :
		0xe4,                                           // CALL
		'3', '2', '7', '6', '8',                       // "32768"
		0x0e, 0x00, 0x00, 0x00, 0x80, 0x00,            // num(32768)
		0x0d,                                           // line terminator
		0xff,                                           // end-of-program sentinel
	}

	gotProg := f.ProgBytes()
	if !bytes.Equal(gotProg, wantProg) {
		t.Errorf("ProgBytes():\n  got  %x\n  want %x", gotProg, wantProg)
	}

	// Full body = PROG + 92 zeros (vars) + 512 zeros (gap)
	gotBody := f.Bytes()
	wantBody := make([]byte, len(wantProg)+92+512)
	copy(wantBody, wantProg)
	if !bytes.Equal(gotBody, wantBody) {
		t.Errorf("Bytes() length = %d; want %d", len(gotBody), len(wantBody))
		if len(gotBody) >= len(wantProg) && len(wantBody) >= len(wantProg) {
			if !bytes.Equal(gotBody[:len(wantProg)], wantProg) {
				t.Errorf("PROG section mismatch")
			}
		}
	}
	if len(gotBody) != 656 {
		t.Errorf("body length = %d; want 656", len(gotBody))
	}
}

func TestFileOffsets(t *testing.T) {
	f := &File{
		Lines: []Line{
			{Number: 10, Tokens: []Token{CLEAR, Number(32767)}},
		},
	}
	progLen := uint32(len(f.ProgBytes()))
	if f.NVARSOffset() != progLen {
		t.Errorf("NVARSOffset() = %d; want %d", f.NVARSOffset(), progLen)
	}
	if f.NUMENDOffset() != progLen+92 {
		t.Errorf("NUMENDOffset() = %d; want %d", f.NUMENDOffset(), progLen+92)
	}
	if f.SAVARSOffset() != progLen+92+512 {
		t.Errorf("SAVARSOffset() = %d; want %d", f.SAVARSOffset(), progLen+92+512)
	}
}

func TestFileOffsetsCustomVars(t *testing.T) {
	f := &File{
		Lines:       []Line{{Number: 1, Tokens: []Token{RUN}}},
		NumericVars: make([]byte, 200),
		Gap:         make([]byte, 1024),
	}
	progLen := uint32(len(f.ProgBytes()))
	if f.NVARSOffset() != progLen {
		t.Errorf("NVARSOffset() = %d; want %d", f.NVARSOffset(), progLen)
	}
	if f.NUMENDOffset() != progLen+200 {
		t.Errorf("NUMENDOffset() = %d; want %d", f.NUMENDOffset(), progLen+200)
	}
	if f.SAVARSOffset() != progLen+200+1024 {
		t.Errorf("SAVARSOffset() = %d; want %d", f.SAVARSOffset(), progLen+200+1024)
	}
	bodyLen := uint32(len(f.Bytes()))
	if bodyLen != progLen+200+1024 {
		t.Errorf("Bytes() length = %d; want %d", bodyLen, progLen+200+1024)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./sambasic/ -v -run 'TestLine|TestBuildDisk|TestFileOffsets'`
Expected: compilation errors — File, Line types not defined.

- [ ] **Step 3: Write the implementation**

```go
// sambasic/file.go
package sambasic

const (
	defaultNumericVarsSize = 92
	defaultGapSize         = 512
)

type File struct {
	Lines           []Line
	NumericVars     []byte
	Gap             []byte
	StringArrayVars []byte
	StartLine       uint16
}

type Line struct {
	Number uint16
	Tokens []Token
}

func (l *Line) Bytes() []byte {
	data := []byte{}
	for _, t := range l.Tokens {
		data = append(data, t.Bytes()...)
	}
	data = append(data, 0x0D)
	result := []byte{
		byte(l.Number >> 8),
		byte(l.Number & 0xFF),
		byte(len(data) & 0xFF),
		byte(len(data) >> 8),
	}
	return append(result, data...)
}

func (f *File) ProgBytes() []byte {
	result := []byte{}
	for _, line := range f.Lines {
		result = append(result, line.Bytes()...)
	}
	result = append(result, 0xFF)
	return result
}

func (f *File) numericVars() []byte {
	if f.NumericVars != nil {
		return f.NumericVars
	}
	return make([]byte, defaultNumericVarsSize)
}

func (f *File) gap() []byte {
	if f.Gap != nil {
		return f.Gap
	}
	return make([]byte, defaultGapSize)
}

func (f *File) stringArrayVars() []byte {
	if f.StringArrayVars != nil {
		return f.StringArrayVars
	}
	return nil
}

func (f *File) Bytes() []byte {
	result := f.ProgBytes()
	result = append(result, f.numericVars()...)
	result = append(result, f.gap()...)
	result = append(result, f.stringArrayVars()...)
	return result
}

func (f *File) NVARSOffset() uint32 {
	return uint32(len(f.ProgBytes()))
}

func (f *File) NUMENDOffset() uint32 {
	return f.NVARSOffset() + uint32(len(f.numericVars()))
}

func (f *File) SAVARSOffset() uint32 {
	return f.NUMENDOffset() + uint32(len(f.gap()))
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./sambasic/ -v -run 'TestLine|TestBuildDisk|TestFileOffsets'`
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add sambasic/file.go sambasic/file_test.go
git commit -m "feat(sambasic): add File/Line types with serialization and offset accessors"
```

---

### Task 4: NewDiskImage and AddBasicFile

**Files:**
- Modify: `samfile.go` — add `NewDiskImage()` and `AddBasicFile()`
- Modify: `samfile_test.go` — add tests

- [ ] **Step 1: Write the failing tests**

Add to `samfile_test.go`:

```go
func TestNewDiskImage(t *testing.T) {
	di := NewDiskImage()
	for i, b := range di {
		if b != 0 {
			t.Fatalf("NewDiskImage()[%d] = 0x%02x; want 0x00", i, b)
		}
	}
	if len(di) != 819200 {
		t.Fatalf("NewDiskImage() length = %d; want 819200", len(di))
	}
}
```

Add import of `sambasic` package at top of test file:

```go
import (
	// ... existing imports ...
	"github.com/petemoore/samfile/v3/sambasic"
)
```

```go
func TestAddBasicFileRoundTrip(t *testing.T) {
	f := &sambasic.File{
		Lines: []sambasic.Line{
			{
				Number: 10,
				Tokens: []sambasic.Token{
					sambasic.CLEAR,
					sambasic.String("32767"),
					sambasic.Number(32767),
					sambasic.Literal(':'),
					sambasic.LOAD,
					sambasic.Literal('"'),
					sambasic.String("stub"),
					sambasic.Literal('"'),
					sambasic.CODE,
					sambasic.String("32768"),
					sambasic.Number(32768),
					sambasic.Literal(':'),
					sambasic.CALL,
					sambasic.String("32768"),
					sambasic.Number(32768),
				},
			},
		},
		StartLine: 10,
	}

	di := NewDiskImage()
	if err := di.AddBasicFile("auto", f); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}

	// Read it back
	readBack, err := di.File("auto")
	if err != nil {
		t.Fatalf("File(\"auto\"): %v", err)
	}

	// Body should match f.Bytes()
	wantBody := f.Bytes()
	if !bytes.Equal(readBack.Body, wantBody) {
		t.Errorf("readback body length = %d; want %d", len(readBack.Body), len(wantBody))
	}

	// File header: type 0x10
	if readBack.Header.Type != FT_SAM_BASIC {
		t.Errorf("header type = %d; want %d (FT_SAM_BASIC)", readBack.Header.Type, FT_SAM_BASIC)
	}

	// Check directory entry
	var fe *FileEntry
	for _, e := range di.DiskJournal() {
		if e.Used() && e.Name.String() == "auto" {
			fe = e
			break
		}
	}
	if fe == nil {
		t.Fatal("auto entry not found in disk journal")
	}
	if fe.Type != FT_SAM_BASIC {
		t.Errorf("fe.Type = %d; want %d", fe.Type, FT_SAM_BASIC)
	}
	if fe.MGTFlags != 0x20 {
		t.Errorf("MGTFlags = 0x%02x; want 0x20", fe.MGTFlags)
	}
	if fe.SAMBASICStartLine != 10 {
		t.Errorf("SAMBASICStartLine = %d; want 10", fe.SAMBASICStartLine)
	}

	// Check page-form triplets via section-size accessors
	if fe.ProgramLength() != f.NVARSOffset() {
		t.Errorf("ProgramLength() = %d; want %d", fe.ProgramLength(), f.NVARSOffset())
	}
	if fe.ProgramLength()+fe.NumericVariablesSize() != f.NUMENDOffset() {
		t.Errorf("ProgramLength+NumericVarsSize = %d; want %d",
			fe.ProgramLength()+fe.NumericVariablesSize(), f.NUMENDOffset())
	}
	if fe.ProgramLength()+fe.NumericVariablesSize()+fe.GapSize() != f.SAVARSOffset() {
		t.Errorf("ProgramLength+NumericVarsSize+GapSize = %d; want %d",
			fe.ProgramLength()+fe.NumericVariablesSize()+fe.GapSize(), f.SAVARSOffset())
	}

	// Detokenize without error
	sb := NewSAMBasic(readBack.Body)
	if err := sb.Output(); err != nil {
		t.Errorf("SAMBasic.Output() on readback: %v", err)
	}
}

func TestAddBasicFileNoAutoRun(t *testing.T) {
	f := &sambasic.File{
		Lines: []sambasic.Line{
			{Number: 10, Tokens: []sambasic.Token{sambasic.PRINT, sambasic.String("hello")}},
		},
		StartLine: 0xFFFF,
	}
	di := NewDiskImage()
	if err := di.AddBasicFile("test", f); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}
	var fe *FileEntry
	for _, e := range di.DiskJournal() {
		if e.Used() && e.Name.String() == "test" {
			fe = e
			break
		}
	}
	if fe == nil {
		t.Fatal("test entry not found")
	}
	if fe.ExecutionAddressDiv16K != 0xFF {
		t.Errorf("ExecutionAddressDiv16K = 0x%02x; want 0xFF (no auto-run)", fe.ExecutionAddressDiv16K)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test -v -run 'TestNewDiskImage|TestAddBasicFile'`
Expected: compilation errors — NewDiskImage and AddBasicFile not defined.

- [ ] **Step 3: Write the implementation**

Add to `samfile.go`, in the imports section add `"github.com/petemoore/samfile/v3/sambasic"`, then add:

```go
// NewDiskImage returns a pointer to a zeroed 819200-byte DiskImage,
// equivalent to a blank formatted MGT floppy.
func NewDiskImage() *DiskImage {
	return &DiskImage{}
}

// pageForm3Byte encodes a value as a 3-byte page-form triplet:
// [page, offset_lo, offset_hi] with offset in 8000H REL PAGE FORM.
// Used for the BASIC dir-entry triplets at 0xDD/0xE0/0xE3.
func pageForm3Byte(value uint32) [3]byte {
	page := byte(value / 16384)
	offset := uint16(value%16384) | 0x8000
	return [3]byte{page, byte(offset & 0xFF), byte(offset >> 8)}
}

// AddBasicFile writes a SAM BASIC file to the disk image, allocating
// a free directory slot and required sectors. name is the 10-char
// space-padded filename. file provides the program lines, variable
// areas, and start-line metadata.
func (di *DiskImage) AddBasicFile(name string, file *sambasic.File) error {
	body := file.Bytes()

	fe := &FileEntry{
		Type:                   FT_SAM_BASIC,
		StartAddressPage:       0,
		StartAddressPageOffset: 0x9CD5,
		MGTFlags:               0x20,
	}

	// Start line / auto-run encoding
	if file.StartLine == 0xFFFF {
		fe.ExecutionAddressDiv16K = 0xFF
		fe.ExecutionAddressMod16K = 0xFFFF
		fe.SAMBASICStartLine = 0xFFFF
	} else {
		fe.ExecutionAddressDiv16K = 0x00
		fe.ExecutionAddressMod16K = file.StartLine
		fe.SAMBASICStartLine = file.StartLine
	}

	// Page-form triplets in FileTypeInfo[0..8]
	nvars := pageForm3Byte(file.NVARSOffset())
	numend := pageForm3Byte(file.NUMENDOffset())
	savars := pageForm3Byte(file.SAVARSOffset())
	copy(fe.FileTypeInfo[0:3], nvars[:])
	copy(fe.FileTypeInfo[3:6], numend[:])
	copy(fe.FileTypeInfo[6:9], savars[:])

	// Mirror the body header into MGTFutureAndPast[1..9] (dir bytes 0xD3-0xDB)
	pages := uint8(len(body) >> 14)
	lengthMod16K := uint16(len(body) & 0x3FFF)
	fe.MGTFutureAndPast[1] = byte(FT_SAM_BASIC)
	fe.MGTFutureAndPast[2] = byte(lengthMod16K)
	fe.MGTFutureAndPast[3] = byte(lengthMod16K >> 8)
	fe.MGTFutureAndPast[4] = 0xD5 // PageOffset lo
	fe.MGTFutureAndPast[5] = 0x9C // PageOffset hi
	fe.MGTFutureAndPast[6] = 0xFF // unused (exec marker)
	fe.MGTFutureAndPast[7] = 0xFF // unused (exec marker)
	fe.MGTFutureAndPast[8] = pages
	fe.MGTFutureAndPast[9] = 0x00 // StartPage

	return di.addFile(name, fe, body)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test -v -run 'TestNewDiskImage|TestAddBasicFile'`
Expected: all PASS.

- [ ] **Step 5: Run the full test suite for regressions**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./... -v`
Expected: all existing tests PASS, no regressions.

- [ ] **Step 6: Commit**

```bash
git add samfile.go samfile_test.go
git commit -m "feat: add NewDiskImage() and AddBasicFile()"
```

---

### Task 5: Harmonize Keyword Table

**Files:**
- Modify: `sambasic.go` — import sambasic for keyword lookup
- Delete: `keywords.go` — replaced by sambasic/keywords.go

- [ ] **Step 1: Modify SAMBasic.Output() to use sambasic.KeywordName()**

Replace the keyword lookup in `sambasic.go`. Change the import block to add `"github.com/petemoore/samfile/v3/sambasic"`. Then replace the two keyword lookups:

For the single-byte path (currently `case b >= 0x85 && b <= 0xf6`):
```go
case b >= 0x85 && b <= 0xf6:
	name, ok := sambasic.KeywordName(b, false)
	if !ok {
		return fmt.Errorf("basic-to-text: keyword index %d out of range", b-0x3b)
	}
	if !spaceBefore {
		fmt.Print(" ")
	}
	fmt.Print(name + " ")
	spaceBefore = true
```

For the 0xFF escape path (currently `case b == 0xff`):
```go
case b == 0xff:
	c++
	if index+uint32(c) >= n {
		return fmt.Errorf("basic-to-text: truncated input: 0xff keyword escape at end of input (offset %d)", index+uint32(c))
	}
	b := basic.Data[index+uint32(c)]
	name, ok := sambasic.KeywordName(b, true)
	if !ok {
		return fmt.Errorf("basic-to-text: invalid keyword byte 0x%02x after 0xff escape at offset %d", b, index+uint32(c))
	}
	if !spaceBefore {
		fmt.Print(" ")
	}
	fmt.Print(name + " ")
	spaceBefore = true
```

- [ ] **Step 2: Delete keywords.go**

```bash
git rm keywords.go
```

- [ ] **Step 3: Run the full test suite**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./... -v`
Expected: all tests PASS — the detokenizer produces identical output via the shared keyword table.

- [ ] **Step 4: Commit**

```bash
git add sambasic.go
git commit -m "refactor: replace keywords.go with sambasic.KeywordName()"
```

---

### Task 6: Parse() for Roundtrip Tests

**Files:**
- Create: `sambasic/parse.go`
- Create: `sambasic/parse_test.go`

- [ ] **Step 1: Write failing parse test**

```go
// sambasic/parse_test.go
package sambasic

import (
	"bytes"
	"testing"
)

func TestParseRoundTrip(t *testing.T) {
	// Build a program, serialize, parse back, re-serialize, compare.
	original := &File{
		Lines: []Line{
			{
				Number: 10,
				Tokens: []Token{
					CLEAR,
					String("32767"),
					Number(32767),
					Literal(':'),
					LOAD,
					Literal('"'),
					String("stub"),
					Literal('"'),
					CODE,
					String("32768"),
					Number(32768),
				},
			},
		},
		StartLine: 10,
	}

	body := original.Bytes()

	parsed, err := Parse(body)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	// Re-serialize: the PROG section should match.
	// We also need to set the same variable areas to get identical Bytes().
	parsed.NumericVars = original.NumericVars
	parsed.Gap = original.Gap
	parsed.StringArrayVars = original.StringArrayVars
	parsed.StartLine = original.StartLine

	got := parsed.Bytes()
	if !bytes.Equal(got, body) {
		t.Errorf("roundtrip mismatch: len(got)=%d, len(want)=%d", len(got), len(body))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./sambasic/ -v -run 'TestParseRoundTrip'`
Expected: compilation error — Parse not defined.

- [ ] **Step 3: Write Parse() implementation**

Parse takes the raw BASIC file body (everything after the 9-byte header — the same bytes that `File.Bytes()` produces) and reconstructs a `File`. It walks the tokenized lines reproducing each byte as the appropriate Token type: keyword bytes become `SingleByteKeyword`/`TwoByteKeyword`, `0x0E` sequences become `Num`, and everything else becomes `Literal`. The post-PROG variable areas are captured as raw `[]byte`.

```go
// sambasic/parse.go
package sambasic

import "fmt"

// Parse reconstructs a File from raw BASIC body bytes (the body
// after the 9-byte FileHeader, i.e. the output of File.Bytes()).
// Lines are parsed token-by-token; the post-PROG variable areas
// (NumericVars, Gap, StringArrayVars) are NOT populated — they
// remain nil, since the caller must supply separate offset info
// from the directory entry to split them. The raw bytes after the
// 0xFF sentinel are returned in StringArrayVars as a single blob
// for roundtrip convenience (callers that know the offsets can
// re-split).
//
// For roundtrip fidelity, every byte in the tokenized line body
// is preserved: keywords as SingleByteKeyword/TwoByteKeyword,
// numeric literals as Num (with Display captured from the bytes
// before the 0x0E marker), and all other bytes as Literal.
func Parse(body []byte) (*File, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("parse: empty body")
	}

	f := &File{}
	pos := 0
	n := len(body)

	for pos < n {
		if body[pos] == 0xFF {
			pos++
			break
		}
		if pos+3 >= n {
			return nil, fmt.Errorf("parse: truncated line header at offset %d", pos)
		}
		lineNum := uint16(body[pos])<<8 | uint16(body[pos+1])
		lineLen := int(body[pos+2]) | int(body[pos+3])<<8
		pos += 4

		if pos+lineLen > n {
			return nil, fmt.Errorf("parse: line %d body extends past input", lineNum)
		}

		line := Line{Number: lineNum}
		end := pos + lineLen
		i := pos
		for i < end {
			b := body[i]
			switch {
			case b == 0x0D:
				// line terminator — don't emit as token, Line.Bytes() adds it
				i++
			case b == 0xFF && i+1 < end:
				line.Tokens = append(line.Tokens, TwoByteKeyword(body[i+1]))
				i += 2
			case b >= 0x85 && b <= 0xF6:
				line.Tokens = append(line.Tokens, SingleByteKeyword(b))
				i++
			case b == 0x0E:
				// Numeric literal: preceding display text was already
				// emitted as Literal tokens. Capture the 5-byte value.
				if i+6 > end {
					return nil, fmt.Errorf("parse: truncated numeric form at offset %d", i)
				}
				// Walk backwards through already-emitted Literal tokens
				// to find the display string (digits before 0x0E).
				display := []byte{}
				for len(line.Tokens) > 0 {
					last, ok := line.Tokens[len(line.Tokens)-1].(Literal)
					if !ok {
						break
					}
					if last >= '0' && last <= '9' || last == '.' || last == '-' || last == 'E' || last == 'e' {
						display = append([]byte{byte(last)}, display...)
						line.Tokens = line.Tokens[:len(line.Tokens)-1]
					} else {
						break
					}
				}
				num := &Num{
					Display: string(display),
				}
				copy(num.Value[:], body[i+1:i+6])
				line.Tokens = append(line.Tokens, num)
				i += 6
			default:
				line.Tokens = append(line.Tokens, Literal(b))
				i++
			}
		}
		f.Lines = append(f.Lines, line)
		pos = end
	}

	// Everything after the 0xFF sentinel is the variable areas.
	// Store as a single blob — the caller splits using dir-entry offsets.
	if pos < n {
		trailer := make([]byte, n-pos)
		copy(trailer, body[pos:])
		// For roundtrip: split into NumericVars + Gap + StringArrayVars
		// requires external offset info. For now, detect the canonical
		// 92+512 case heuristically, or store as NumericVars if total is
		// exactly 604.
		total := len(trailer)
		if total == 604 {
			f.NumericVars = trailer[:92]
			f.Gap = trailer[92:]
		} else if total > 0 {
			// Can't split without offset info — store entire trailer as
			// NumericVars so Bytes() reproduces it. Gap and StringArrayVars
			// are set to empty to avoid double-counting.
			f.NumericVars = trailer
			f.Gap = []byte{}
		}
	} else {
		// No trailer at all
		f.NumericVars = []byte{}
		f.Gap = []byte{}
	}

	return f, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./sambasic/ -v -run 'TestParseRoundTrip'`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add sambasic/parse.go sambasic/parse_test.go
git commit -m "feat(sambasic): add Parse() for roundtrip testing"
```

---

### Task 7: Exhaustive Real-Disk Roundtrip Tests

**Files:**
- Create: `sambasic/roundtrip_test.go`

- [ ] **Step 1: Write the exhaustive roundtrip test**

This test scans `~/Downloads/GoodSamC2/` for all `.dsk` and `.mgt` files, loads each as an MGT image (skipping EDSK), and for every `FT_SAM_BASIC` file: extracts body bytes, parses with `Parse()`, re-serializes, and compares.

```go
// sambasic/roundtrip_test.go
package sambasic_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	samfile "github.com/petemoore/samfile/v3"
	"github.com/petemoore/samfile/v3/sambasic"
)

func TestExhaustiveRealDiskRoundtrip(t *testing.T) {
	// Collect disk images: scan ~/Downloads/ and ~/git/ recursively
	// for .dsk and .mgt files that are exactly 819200 bytes (the
	// canonical MGT image size). This catches raw MGT files regardless
	// of extension and skips archives, EDSK headers, etc.
	//
	// Additionally, a testdata/mgt/ directory under this package holds
	// a curated set of at least 40 .mgt files downloaded from public
	// SAM Coupé archives (World of Sam, etc.) for CI-reproducible
	// testing. The filesystem scan is additive — it finds any extras
	// on the developer's machine.
	var diskPaths []string

	home := os.Getenv("HOME")
	searchDirs := []string{
		filepath.Join(home, "Downloads"),
		filepath.Join(home, "git"),
	}

	// Also include the vendored test corpus
	if wd, err := os.Getwd(); err == nil {
		searchDirs = append(searchDirs, filepath.Join(wd, "testdata", "mgt"))
	}

	seen := map[string]bool{}
	for _, dir := range searchDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			ext := filepath.Ext(path)
			if ext != ".dsk" && ext != ".mgt" {
				return nil
			}
			// All valid MGT images are exactly 819200 bytes
			if info.Size() != 819200 {
				return nil
			}
			if !seen[path] {
				seen[path] = true
				diskPaths = append(diskPaths, path)
			}
			return nil
		})
	}

	if len(diskPaths) == 0 {
		t.Skip("no .dsk or .mgt files found; skipping real-disk tests")
	}

	var disksScanned, basicFilesTotal, passCount, failCount int

	for _, path := range diskPaths {
		entry := filepath.Base(path)
		di, err := samfile.Load(path)
		if err != nil {
			// Skip EDSK or unreadable files
			continue
		}
		disksScanned++

		dj := di.DiskJournal()
		for slot, fe := range dj {
			if !fe.Used() || fe.Type != samfile.FT_SAM_BASIC {
				continue
			}
			basicFilesTotal++
			fileName := fe.Name.String()
			t.Run(entry+"/"+fileName, func(t *testing.T) {
				file, err := di.File(fileName)
				if err != nil {
					failCount++
					t.Errorf("slot %d: File(%q): %v", slot, fileName, err)
					return
				}

				parsed, err := sambasic.Parse(file.Body)
				if err != nil {
					failCount++
					t.Errorf("slot %d: Parse(%q): %v", slot, fileName, err)
					return
				}

				// Re-serialize: for roundtrip, we need exact byte match of
				// the PROG section. The variable areas may not split
				// identically without dir-entry offsets, so compare full body.
				got := parsed.Bytes()
				if !bytes.Equal(got, file.Body) {
					failCount++
					t.Errorf("slot %d %q: roundtrip mismatch: got %d bytes, want %d bytes",
						slot, fileName, len(got), len(file.Body))
					// Find first difference
					minLen := len(got)
					if len(file.Body) < minLen {
						minLen = len(file.Body)
					}
					for i := 0; i < minLen; i++ {
						if got[i] != file.Body[i] {
							t.Errorf("  first diff at offset %d: got 0x%02x, want 0x%02x", i, got[i], file.Body[i])
							break
						}
					}
				} else {
					passCount++
				}
			})
		}
	}

	t.Logf("Summary: %d disks scanned, %d BASIC files tested, %d passed, %d failed",
		disksScanned, basicFilesTotal, passCount, failCount)
}
```

- [ ] **Step 2: Run the test**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./sambasic/ -v -run 'TestExhaustiveRealDiskRoundtrip' -timeout 300s`
Expected: The test runs, reports statistics. Some files may fail initially due to non-canonical trailer splits — those failures guide `Parse()` refinements.

- [ ] **Step 3: Fix any Parse() issues discovered**

Iterate on `sambasic/parse.go` to handle edge cases discovered by the real-disk corpus. Common issues:
- Trailer splits other than 92+512 (use directory entry's page-form triplets to guide the split)
- Programs with actual variable data
- Empty programs (just 0xFF sentinel)

This may require updating `Parse()` to accept optional offset parameters, or adding a `ParseWithOffsets(body []byte, nvarsOff, numendOff, savarsOff uint32)` variant. Update the roundtrip test to pass offsets from the directory entry.

- [ ] **Step 4: Re-run until the pass rate is satisfactory**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./sambasic/ -v -run 'TestExhaustiveRealDiskRoundtrip' -timeout 300s`
Expected: high pass rate. Document any remaining failures (corrupt disks, non-standard ROMs).

- [ ] **Step 5: Commit**

```bash
git add sambasic/roundtrip_test.go sambasic/parse.go
git commit -m "test(sambasic): add exhaustive real-disk roundtrip tests"
```

---

### Task 8: Multi-File Integration Test

**Files:**
- Modify: `samfile_test.go`

- [ ] **Step 1: Write the multi-file integration test**

This replicates the `build-disk.sh` layout: samdos2 (CODE) + auto (BASIC) + stub (CODE) on one disk.

```go
func TestMultiFileBasicAndCode(t *testing.T) {
	di := NewDiskImage()

	// Add a CODE file first (simulates samdos2-like file)
	codeBody := bytes.Repeat([]byte{0xAA}, 1000)
	if err := di.AddCodeFile("samdos2", codeBody, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile(samdos2): %v", err)
	}

	// Add BASIC AUTO file
	basicFile := &sambasic.File{
		Lines: []sambasic.Line{
			{
				Number: 10,
				Tokens: []sambasic.Token{
					sambasic.CLEAR,
					sambasic.String("32767"),
					sambasic.Number(32767),
					sambasic.Literal(':'),
					sambasic.LOAD,
					sambasic.Literal('"'),
					sambasic.String("stub"),
					sambasic.Literal('"'),
					sambasic.CODE,
					sambasic.String("32768"),
					sambasic.Number(32768),
					sambasic.Literal(':'),
					sambasic.CALL,
					sambasic.String("32768"),
					sambasic.Number(32768),
				},
			},
		},
		StartLine: 10,
	}
	if err := di.AddBasicFile("auto", basicFile); err != nil {
		t.Fatalf("AddBasicFile(auto): %v", err)
	}

	// Add a second CODE file (simulates stub)
	stubBody := bytes.Repeat([]byte{0xBB}, 100)
	if err := di.AddCodeFile("stub", stubBody, 0x8000, 0x8000); err != nil {
		t.Fatalf("AddCodeFile(stub): %v", err)
	}

	// Read all three back — no corruption
	for _, name := range []string{"samdos2", "auto", "stub"} {
		if _, err := di.File(name); err != nil {
			t.Errorf("File(%q): %v", name, err)
		}
	}

	// Verify sector maps don't overlap
	dj := di.DiskJournal()
	used := dj.UsedFileEntries()
	if len(used) != 3 {
		t.Errorf("used entries = %d; want 3", len(used))
	}
	for i := 0; i < len(used); i++ {
		for j := i + 1; j < len(used); j++ {
			a := dj[used[i]].SectorAddressMap
			b := dj[used[j]].SectorAddressMap
			for k := 0; k < len(a); k++ {
				if a[k]&b[k] != 0 {
					t.Errorf("sector maps for slots %d and %d overlap at byte %d: 0x%02x & 0x%02x",
						used[i], used[j], k, a[k], b[k])
				}
			}
		}
	}
}
```

- [ ] **Step 2: Run the test**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test -v -run 'TestMultiFileBasicAndCode'`
Expected: PASS.

- [ ] **Step 3: Run full test suite**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./... -v`
Expected: all PASS.

- [ ] **Step 4: Commit**

```bash
git add samfile_test.go
git commit -m "test: add multi-file BASIC+CODE integration test"
```

---

### Task 9: Final Lint and Vet

- [ ] **Step 1: Run go vet**

Run: `cd /Users/pmoore/git/samfile-create-basic && go vet ./...`
Expected: no issues.

- [ ] **Step 2: Run go build to verify compilation**

Run: `cd /Users/pmoore/git/samfile-create-basic && go build ./...`
Expected: clean build.

- [ ] **Step 3: Run full test suite one final time**

Run: `cd /Users/pmoore/git/samfile-create-basic && go test ./... -v -count=1`
Expected: all PASS.

- [ ] **Step 4: Commit any remaining fixes**

Only if lint/vet surfaced issues that needed fixing.
