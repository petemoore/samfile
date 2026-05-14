package sambasic

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
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
	f, _, err := parseTextInternal(src, false)
	return f, err
}

// ParseTextSchema is like ParseTextString but also returns a Schema
// describing the produced bytes for conformance testing.
func ParseTextSchema(src string) (*File, *Schema, error) {
	return parseTextInternal(src, true)
}

// bodyItem records one lexer item's contribution to a line's body
// bytes, along with the schema metadata needed to classify it later.
type bodyItem struct {
	kind   SegmentKind
	bytes  []byte
	detail any
}

// lineDraft accumulates body items and concatenated bytes for a line
// before finalisation.
type lineDraft struct {
	items []bodyItem
	body  []byte
}

func (d *lineDraft) append(kind SegmentKind, bytes []byte, detail any) {
	d.items = append(d.items, bodyItem{kind: kind, bytes: append([]byte(nil), bytes...), detail: detail})
	d.body = append(d.body, bytes...)
}

func (d *lineDraft) empty() bool {
	return len(d.body) == 0
}

func parseTextInternal(src string, withSchema bool) (*File, *Schema, error) {
	l := lex(src)
	l.state = lexStart

	edits := map[uint16]Line{}
	drafts := map[uint16]*lineDraft{}

	var curNum uint16
	var draft *lineDraft
	haveLineNumber := false
	// lastKeywordByte tracks the most recent keyword's first byte so
	// that the literal payload of a REM keyword can be reclassified as
	// SegREMBody rather than a run of SegLiteral bytes.
	var lastKeywordByte byte

	flush := func() {
		if !haveLineNumber {
			return
		}
		if draft == nil || draft.empty() {
			delete(edits, curNum)
			delete(drafts, curNum)
			return
		}
		ln, patched := finaliseWithBody(Line{Number: curNum, Tokens: tokenizeBytes(draft.body)})
		edits[curNum] = ln
		if withSchema {
			// Replace each item's bytes with the corresponding slice
			// from the patched body so D7→D8 / D9→DA / FF→A1 rewrites
			// are reflected in the schema.
			pos := 0
			for i := range draft.items {
				ln := len(draft.items[i].bytes)
				draft.items[i].bytes = append([]byte(nil), patched[pos:pos+ln]...)
				pos += ln
			}
			drafts[curNum] = draft
		}
	}

	startLine := func(n uint16) {
		curNum = n
		draft = &lineDraft{}
		haveLineNumber = true
		lastKeywordByte = 0
	}

	for {
		it := l.nextItem()
		switch it.typ {
		case itemError:
			return nil, nil, &ParseError{Line: it.line, Col: it.col, Msg: it.val}
		case itemEOF:
			flush()
			f := assembleFile(edits)
			var sch *Schema
			if withSchema {
				sch = buildSchema(f, drafts)
			}
			return f, sch, nil
		case itemLineNumber:
			flush()
			n, err := strconv.ParseUint(it.val, 10, 16)
			if err != nil {
				return nil, nil, &ParseError{Line: it.line, Col: it.col, Msg: fmt.Sprintf("invalid line number: %s", it.val)}
			}
			startLine(uint16(n))
		case itemEOL:
			flush()
			haveLineNumber = false
			draft = nil
		default:
			if draft == nil {
				// Stray body content before a line number — shouldn't
				// happen given lexStart's structure, but ignore safely.
				break
			}
			var bytes []byte
			if it.bytes != nil {
				bytes = it.bytes
			} else {
				bytes = []byte(it.val)
			}
			kind, detail := classifyItem(it, bytes, lastKeywordByte)
			draft.append(kind, bytes, detail)
			if it.typ == itemKeyword && len(bytes) > 0 {
				lastKeywordByte = bytes[0]
				if len(bytes) > 1 {
					// 2-byte keyword: the meaningful byte is the second one.
					lastKeywordByte = bytes[1]
				}
			} else if it.typ != itemLiteral || len(bytes) != 0 {
				// Any non-keyword emission resets REM tracking; we only
				// want a REM keyword immediately followed by its body.
				if kind == SegREMBody {
					// Don't reset within a REM body emission.
				} else if it.typ == itemKeyword {
					// already handled above
				} else {
					// Once the keyword's body (REM) has been consumed we
					// no longer want subsequent literals classified as
					// REM. The lexComment state pushes only a single
					// itemLiteral for the body, so any further items on
					// a later line are unrelated.
					lastKeywordByte = 0
				}
			}
		}
	}
}

// classifyItem maps a lex item to a SegmentKind and optional Detail.
func classifyItem(it item, bytes []byte, lastKeywordByte byte) (SegmentKind, any) {
	switch it.typ {
	case itemKeyword:
		return SegKeyword, nil
	case itemNumber:
		// Find the 0x0E marker to split the visible literal from the FP form.
		fpIdx := -1
		for i, b := range bytes {
			if b == 0x0E {
				fpIdx = i
				break
			}
		}
		_ = fpIdx // visible-vs-FP split is handled in the schema builder
		detail := numberDetailFromBytes(bytes, it.val)
		return SegNumberFP, detail
	case itemString:
		return SegString, nil
	case itemControlEscape:
		return SegControlEscape, nil
	case itemProcCallPlaceholder:
		return SegProcCallPlaceholder, PlaceholderDetail{Identifier: it.val, IDLen: len(it.val)}
	case itemFnCallPlaceholder:
		return SegFnCallPlaceholder, PlaceholderDetail{Identifier: it.val, IDLen: len(it.val)}
	case itemLiteral:
		// REM-body recognition: when the previous keyword was REM, this
		// literal carries the REM body.
		if lastKeywordByte == byte(REM) {
			return SegREMBody, nil
		}
		return SegLiteral, nil
	}
	return SegLiteral, nil
}

// numberDetailFromBytes parses the visible-literal prefix of an
// itemNumber's bytes to produce a NumberDetail.
func numberDetailFromBytes(bytes []byte, display string) NumberDetail {
	// The FP bytes follow a 0x0E marker. Use the recorded display
	// string as the visible literal when available.
	literal := display
	if literal == "" {
		// Fall back to extracting up to the 0x0E.
		for i, b := range bytes {
			if b == 0x0E {
				literal = string(bytes[:i])
				break
			}
		}
	}
	// Hex form.
	if len(literal) > 0 && literal[0] == '&' {
		v, err := strconv.ParseUint(literal[1:], 16, 64)
		if err == nil && v <= 0xFFFF {
			return NumberDetail{IsInt: true, Int: uint16(v)}
		}
		if err == nil {
			return NumberDetail{IsInt: false, Float: float64(v)}
		}
	}
	// Try integer fast-path.
	if !strings.ContainsAny(literal, ".eE") && literal != "" {
		v, err := strconv.ParseUint(literal, 10, 64)
		if err == nil && v <= 0xFFFF {
			return NumberDetail{IsInt: true, Int: uint16(v)}
		}
	}
	// General float parse.
	f, err := strconv.ParseFloat(literal, 64)
	if err == nil {
		return NumberDetail{IsInt: false, Float: f}
	}
	// Last resort: decode the recorded FP bytes.
	if idx := indexOf(bytes, 0x0E); idx >= 0 && idx+6 <= len(bytes) {
		if val, ok := decodeFP5(bytes[idx+1 : idx+6]); ok {
			return NumberDetail{IsInt: false, Float: val}
		}
	}
	return NumberDetail{}
}

func indexOf(b []byte, target byte) int {
	for i, x := range b {
		if x == target {
			return i
		}
	}
	return -1
}

// buildSchema walks the assembled file and the recorded per-line
// drafts, producing absolute-offset SchemaSegments for every byte in
// f.ProgBytes(). For SegNumberFP items the recorded bytes also include
// the visible literal preceding the 0x0E marker; that gets split into
// a SegNumberVisible + SegNumberFP segment pair.
func buildSchema(f *File, drafts map[uint16]*lineDraft) *Schema {
	sch := &Schema{}
	offset := 0
	for _, line := range f.Lines {
		// Line header: 4 bytes.
		lineBytes := line.Bytes()
		headerLen := 4
		header := lineBytes[:headerLen]
		sch.Segments = append(sch.Segments, SchemaSegment{
			Offset: offset,
			Length: headerLen,
			Kind:   SegLineHeader,
			Bytes:  append([]byte(nil), header...),
		})
		bodyStart := offset + headerLen
		bodyOffset := bodyStart
		// Body: walk the draft's items.
		draft := drafts[line.Number]
		if draft != nil {
			for _, it := range draft.items {
				switch it.kind {
				case SegNumberFP:
					// Split into visible + FP form using the 0x0E marker.
					fpIdx := indexOf(it.bytes, 0x0E)
					if fpIdx > 0 {
						sch.Segments = append(sch.Segments, SchemaSegment{
							Offset: bodyOffset,
							Length: fpIdx,
							Kind:   SegNumberVisible,
							Bytes:  append([]byte(nil), it.bytes[:fpIdx]...),
						})
						bodyOffset += fpIdx
					}
					if fpIdx >= 0 {
						fpLen := len(it.bytes) - fpIdx
						sch.Segments = append(sch.Segments, SchemaSegment{
							Offset: bodyOffset,
							Length: fpLen,
							Kind:   SegNumberFP,
							Bytes:  append([]byte(nil), it.bytes[fpIdx:]...),
							Detail: it.detail,
						})
						bodyOffset += fpLen
					} else {
						// No 0x0E marker — odd; record as-is.
						sch.Segments = append(sch.Segments, SchemaSegment{
							Offset: bodyOffset,
							Length: len(it.bytes),
							Kind:   SegNumberFP,
							Bytes:  append([]byte(nil), it.bytes...),
							Detail: it.detail,
						})
						bodyOffset += len(it.bytes)
					}
				default:
					sch.Segments = append(sch.Segments, SchemaSegment{
						Offset: bodyOffset,
						Length: len(it.bytes),
						Kind:   it.kind,
						Bytes:  append([]byte(nil), it.bytes...),
						Detail: it.detail,
					})
					bodyOffset += len(it.bytes)
				}
			}
		}
		// Line terminator: 1 byte (0x0D).
		sch.Segments = append(sch.Segments, SchemaSegment{
			Offset: bodyOffset,
			Length: 1,
			Kind:   SegLineTerminator,
			Bytes:  []byte{0x0D},
		})
		offset = bodyOffset + 1
	}
	// Program end sentinel.
	sch.Segments = append(sch.Segments, SchemaSegment{
		Offset: offset,
		Length: 1,
		Kind:   SegProgramEnd,
		Bytes:  []byte{0xFF},
	})
	return sch
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
	out, _ := finaliseWithBody(line)
	return out
}

// finaliseWithBody is like finalise but also returns the patched body
// byte slice so callers (e.g. ParseTextSchema) can re-sync schema
// segment bytes with the rewrites finalise applied.
func finaliseWithBody(line Line) (Line, []byte) {
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
	return Line{Number: line.Number, Tokens: out}, patched
}

// flattenTokens returns the concatenation of every token's Bytes().
func flattenTokens(tokens []Token) []byte {
	var out []byte
	for _, t := range tokens {
		out = append(out, t.Bytes()...)
	}
	return out
}
