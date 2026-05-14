package sambasic

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

// ParseError describes a failure to tokenise SAM BASIC text input.
type ParseError struct {
	Line int    // 1-based source line
	Col  int    // 1-based source column
	Msg  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d, col %d: %s", e.Line, e.Col, e.Msg)
}

// ParseText reads SAM BASIC source from r and returns a *File whose
// program bytes match what the SAM BASIC editor would store after the
// equivalent typing session, including line-edit semantics (sort by
// line number ascending, last-write-wins, bare-line-number delete).
func ParseText(r io.Reader) (*File, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseTextString(string(src))
}

// ParseTextString is the string-input variant of ParseText.
func ParseTextString(src string) (*File, error) {
	l := lex(src)
	l.state = lexStart

	edits := map[uint16]Line{}

	var curNum uint16
	var bodyBytes []byte
	bodyEmpty := true
	haveLineNumber := false

	for {
		it := l.nextItem()
		switch it.typ {
		case itemError:
			return nil, &ParseError{Line: it.line, Col: it.col, Msg: it.val}
		case itemEOF:
			if haveLineNumber {
				if bodyEmpty {
					delete(edits, curNum)
				} else {
					edits[curNum] = finalise(Line{Number: curNum, Tokens: tokenizeBytes(bodyBytes)})
				}
			}
			return assembleFile(edits), nil
		case itemLineNumber:
			if haveLineNumber {
				if bodyEmpty {
					delete(edits, curNum)
				} else {
					edits[curNum] = finalise(Line{Number: curNum, Tokens: tokenizeBytes(bodyBytes)})
				}
			}
			n, err := strconv.ParseUint(it.val, 10, 16)
			if err != nil {
				return nil, &ParseError{Line: it.line, Col: it.col, Msg: fmt.Sprintf("invalid line number: %s", it.val)}
			}
			curNum = uint16(n)
			bodyBytes = nil
			bodyEmpty = true
			haveLineNumber = true
		case itemEOL:
			if haveLineNumber {
				if bodyEmpty {
					delete(edits, curNum)
				} else {
					edits[curNum] = finalise(Line{Number: curNum, Tokens: tokenizeBytes(bodyBytes)})
				}
				haveLineNumber = false
			}
			bodyBytes = nil
			bodyEmpty = true
		default:
			if it.bytes != nil {
				bodyBytes = append(bodyBytes, it.bytes...)
			} else {
				bodyBytes = append(bodyBytes, []byte(it.val)...)
			}
			bodyEmpty = false
		}
	}
}

// tokenizeBytes converts a raw byte slice into a sequence of Token
// values suitable for sambasic.File.Bytes(). For v1 we use a single
// `literal` token per byte; output bytes are unchanged either way.
func tokenizeBytes(bytes []byte) []Token {
	tokens := make([]Token, 0, len(bytes))
	for _, b := range bytes {
		tokens = append(tokens, literal(b))
	}
	return tokens
}

func assembleFile(edits map[uint16]Line) *File {
	nums := make([]uint16, 0, len(edits))
	for n := range edits {
		nums = append(nums, n)
	}
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	lines := make([]Line, len(nums))
	for i, n := range nums {
		lines[i] = edits[n]
	}
	return &File{Lines: lines}
}

// finalise applies SAM's two-pass byte-level patches to a tokenised line:
//
//	LIF (0xD7) → SIF (0xD8) iff THEN (0x8D) appears later in the line.
//	LELSE (0xD9) → ELSE (0xDA) iff THEN appeared earlier.
//	INK (0xFF, as a standalone 1-byte token) → PEN (0xA1) unconditionally.
//
// See grammar spec §1, §3.3, §3.10.
//
// The INK rewrite needs to distinguish a standalone 0xFF (the INK keyword,
// at table slot 0xFF — see grammar §3.3) from a 0xFF that introduces a
// 2-byte keyword pair. Two-byte prefixes are always followed by a byte in
// 0x3B..0x84 (the 2-byte keyword index range), so we skip those.
func finalise(line Line) Line {
	body := flattenTokens(line.Tokens)
	hasTHEN := false
	for i := 0; i < len(body); i++ {
		b := body[i]
		if b == 0x0E {
			// Skip 0x0E marker + 5 FP bytes; these may contain arbitrary
			// bytes (including 0xFF) that must not be reinterpreted.
			i += 5
			continue
		}
		if b == 0x8D {
			hasTHEN = true
		}
	}
	patched := make([]byte, len(body))
	copy(patched, body)
	inkByte := byte(INK)
	penByte := byte(PEN)
	for i := 0; i < len(patched); i++ {
		b := patched[i]
		if b == 0x0E {
			// Skip 0x0E marker + 5 FP bytes.
			i += 5
			continue
		}
		// Skip 2-byte keyword pairs so we don't mistake the 0xFF prefix
		// for a standalone INK token.
		if b == 0xFF && i+1 < len(patched) {
			next := patched[i+1]
			if next >= 0x3B && next <= 0x84 {
				i++ // consume the index byte
				continue
			}
		}
		switch {
		case b == 0xD7 && hasTHEN:
			patched[i] = 0xD8
		case b == 0xD9 && hasTHEN:
			patched[i] = 0xDA
		case b == inkByte:
			patched[i] = penByte
		}
	}
	out := make([]Token, len(patched))
	for i, b := range patched {
		out[i] = literal(b)
	}
	return Line{Number: line.Number, Tokens: out}
}

// flattenTokens returns the concatenation of every token's Bytes().
func flattenTokens(tokens []Token) []byte {
	var out []byte
	for _, t := range tokens {
		out = append(out, t.Bytes()...)
	}
	return out
}
