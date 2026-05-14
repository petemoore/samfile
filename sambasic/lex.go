package sambasic

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const eof = -1

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemEOL
	itemLineNumber
	itemKeyword
	itemNumber
	itemString
	itemLiteral
	itemControlEscape
	itemProcCallPlaceholder
	itemFnCallPlaceholder
)

type item struct {
	typ   itemType
	val   string
	bytes []byte
	line  int
	col   int
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input       string
	pos         int
	start       int
	width       int
	line        int
	col         int
	startLine   int
	startCol    int
	state       stateFn
	items       chan item
	stmtInitial bool
}

func lex(input string) *lexer {
	l := &lexer{
		input:       input,
		line:        1,
		col:         1,
		startLine:   1,
		startCol:    1,
		items:       make(chan item, 2),
		stmtInitial: true,
	}
	return l
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += w
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return r
}

func (l *lexer) backup() {
	if l.width == 0 {
		return
	}
	l.pos -= l.width
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	if r == '\n' {
		l.line--
		// col tracking after backup over a newline is approximate; refine if a test demands it
	} else {
		l.col--
	}
	l.width = 0
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) ignore() {
	l.start = l.pos
	l.startLine = l.line
	l.startCol = l.col
}

func (l *lexer) emit(t itemType) {
	l.items <- item{
		typ:  t,
		val:  l.input[l.start:l.pos],
		line: l.startLine,
		col:  l.startCol,
	}
	l.start = l.pos
	l.startLine = l.line
	l.startCol = l.col
}

// emitBytes pushes an item with pre-resolved disk bytes (e.g. a keyword's
// tokenised form). The caller is responsible for ensuring l.pos has been
// advanced past the consumed source span before calling, so that start is
// correctly reset for the next token.
func (l *lexer) emitBytes(t itemType, b []byte, val string) {
	l.items <- item{
		typ:   t,
		val:   val,
		bytes: b,
		line:  l.startLine,
		col:   l.startCol,
	}
	l.start = l.pos
	l.startLine = l.line
	l.startCol = l.col
}

func (l *lexer) errorf(format string, args ...any) stateFn {
	l.items <- item{
		typ:  itemError,
		val:  fmt.Sprintf(format, args...),
		line: l.line,
		col:  l.col,
	}
	return nil
}

// nextItem drives the state machine inline and returns the next item.
//
// IMPORTANT for state-function authors: state functions must emit at
// most 2 items per invocation (the items channel buffer size). A state
// function that wants to emit more must emit one item and return either
// itself or another state function — nextItem will drain the channel
// before invoking the returned state. The canonical pattern is in
// lexBodyLoop: read one rune, emit at most one item, return self. The
// alternative would be to run the state machine in a goroutine with an
// unbuffered (or larger-buffer) channel; we deliberately don't, to
// keep the lexer single-goroutine and lifecycle-simple.
func (l *lexer) nextItem() item {
	for {
		select {
		case it := <-l.items:
			return it
		default:
			if l.state == nil {
				return item{typ: itemEOF, line: l.line, col: l.col}
			}
			l.state = l.state(l)
		}
	}
}

// lexStart strips leading whitespace and dispatches to line-number parsing
// or emits EOF on empty/whitespace-only input.
func lexStart(l *lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			l.emit(itemEOF)
			return nil
		}
		// Treat any byte < 0x21 as whitespace per the editor's GTCH3 skip.
		// Strict reading would also skip 0x00..0x1F; we handle 0x20 ' ',
		// 0x09 '\t', 0x0D '\r', 0x0A '\n'. Newlines that come before any
		// line number are simply skipped (empty leading lines).
		if r == ' ' || r == '\t' || r == '\r' {
			l.ignore()
			continue
		}
		if r == '\n' {
			l.ignore()
			continue
		}
		// First significant character: dispatch.
		l.backup()
		l.start = l.pos
		l.startLine = l.line
		l.startCol = l.col
		return lexLineNumber
	}
}

// lexLineNumber consumes a decimal digit run and emits itemLineNumber.
// Range: 1..0xFEFF (65279). Line 0 is reserved.
func lexLineNumber(l *lexer) stateFn {
	const digits = "0123456789"
	if !l.accept(digits) {
		// First non-WS char isn't a digit.
		return l.errorf("expected line number")
	}
	l.acceptRun(digits)
	text := l.input[l.start:l.pos]
	// Parse to validate range.
	var n uint64
	for _, c := range text {
		n = n*10 + uint64(c-'0')
		if n > 0xFFFF {
			return l.errorf("line number out of range: %s", text)
		}
	}
	if n == 0 {
		return l.errorf("line number 0 is reserved")
	}
	if n > 0xFEFF {
		return l.errorf("line number out of range: %s", text)
	}
	l.emit(itemLineNumber)
	return lexBody
}

// lexBody is entered immediately after itemLineNumber. It applies the
// post-line-number one-space-drop (grammar spec §2.3) and then dispatches
// to lexBodyLoop. For now (until later tasks add keyword/number/string/etc.
// handling), every non-newline byte is emitted as itemLiteral.
func lexBody(l *lexer) stateFn {
	l.stmtInitial = true
	// One-space-drop: examine the first byte after the line number digits.
	r := l.next()
	if r == eof {
		l.emit(itemEOL)
		l.emit(itemEOF)
		return nil
	}
	if r == '\n' || r == '\r' {
		l.backup()
		l.start = l.pos
		l.startLine = l.line
		l.startCol = l.col
		l.next()
		l.ignore()
		l.emit(itemEOL)
		return lexStart
	}
	if r == ' ' {
		// b1 is space; peek at b2.
		next := l.peek()
		if next == '\n' || next == '\r' || next == eof {
			// b2 is the terminator: PRESERVE this space — emit as literal.
			l.emit(itemLiteral)
		} else {
			// b2 is something else: DROP this space.
			l.ignore()
		}
	} else {
		// b1 is not a space: emit it as the first body byte.
		l.backup()
	}
	return lexBodyLoop
}

// lexBodyLoop handles the rest of the body byte-by-byte. Each non-newline
// byte is emitted as an itemLiteral for now; later tasks replace this with
// keyword/number/string/etc. dispatch.
//
// We return to nextItem after each emit (rather than looping internally) so
// the inline state-machine driver in nextItem can drain the buffered item
// channel between emits. An internal for-loop would deadlock on bodies of
// more than two bytes once the channel buffer fills.
func lexBodyLoop(l *lexer) stateFn {
	r := l.next()
	if r == eof {
		l.emit(itemEOL)
		l.emit(itemEOF)
		return nil
	}
	if r == '\n' || r == '\r' {
		l.backup()
		l.start = l.pos
		l.startLine = l.line
		l.startCol = l.col
		l.next()
		l.ignore()
		l.emit(itemEOL)
		return lexStart
	}
	if r == '"' {
		l.backup()
		l.stmtInitial = false
		return lexString
	}
	if r == ' ' {
		// Leading-space-drop: if a keyword starts right after this space,
		// drop the space (don't emit it). Don't touch stmtInitial here:
		// we're still in statement-initial position.
		if l.pos < len(l.input) && isAlpha(rune(l.input[l.pos])) {
			if _, _, _, isKW := lookupKeyword(l.input, l.pos); isKW {
				l.ignore()
				return lexBodyLoop
			}
		}
		l.emit(itemLiteral)
		l.stmtInitial = false
		return lexBodyLoop
	}
	if isAlpha(r) {
		l.backup()
		return lexKeyword
	}
	if r == '&' {
		l.backup()
		l.stmtInitial = false
		return lexNumber
	}
	if r >= '0' && r <= '9' {
		l.backup()
		l.stmtInitial = false
		return lexNumber
	}
	if r == '{' {
		l.backup()
		l.stmtInitial = false
		return lexControlEscape
	}
	if r == ':' {
		l.emit(itemLiteral)
		l.stmtInitial = true
		return lexBodyLoop
	}
	l.emit(itemLiteral)
	l.stmtInitial = false
	return lexBodyLoop
}

// lexNumber scans a decimal numeric literal and emits itemNumber.
// Handles integer-only here; Tasks 8-10 extend to hex, scientific,
// and BIN forms. Enforces digit-then-letter rejection per grammar
// spec §4.2: after the digits, the next character must not be a
// letter or underscore.
func lexNumber(l *lexer) stateFn {
	// Hex literal: &[0-9A-Fa-f]+
	if l.peek() == '&' {
		l.next() // consume &
		const hexDigits = "0123456789abcdefABCDEF"
		if !l.accept(hexDigits) {
			return l.errorf("expected hex digits after &")
		}
		l.acceptRun(hexDigits)
		if r := l.peek(); r != eof && (isAlpha(r) || r == '_') {
			l.next()
			return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}
		return emitNumberFP(l)
	}
	// Decimal integer / decimal-with-fraction / scientific.
	const digits = "0123456789"
	l.acceptRun(digits)
	if l.peek() == '.' {
		l.next()
		l.acceptRun(digits)
	}
	if r := l.peek(); r == 'e' || r == 'E' {
		l.next()
		if r := l.peek(); r == '+' || r == '-' {
			l.next()
		}
		if !l.accept(digits) {
			return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}
		l.acceptRun(digits)
	}
	if r := l.peek(); r != eof && (isAlpha(r) || r == '_') {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	return emitNumberFP(l)
}

// emitNumberFP emits an itemNumber whose bytes are the visible literal
// followed by the 0x0E marker and the 5-byte FP form returned by
// encodeFP.
func emitNumberFP(l *lexer) stateFn {
	literal := l.input[l.start:l.pos]
	fp, err := encodeFP(literal)
	if err != nil {
		return l.errorf("%s", err.Error())
	}
	out := make([]byte, 0, len(literal)+6)
	out = append(out, []byte(literal)...)
	out = append(out, 0x0E)
	out = append(out, fp[:]...)
	l.emitBytes(itemNumber, out, literal)
	return lexBodyLoop
}

func isAlpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isAlphaNum(r rune) bool {
	return isAlpha(r) || (r >= '0' && r <= '9') || r == '_'
}

// keywordFold returns the input string with ASCII letters uppercased
// (the SAM editor's AND 0DFH fold). Non-letter bytes pass through
// unchanged.
func keywordFold(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 0x20
		}
		b[i] = c
	}
	return string(b)
}

// lookupKeyword tries to match the input starting at pos against the
// keyword table (zero-or-one space rule for multi-word entries; ASCII
// case-fold). Returns the matched canonical form, the byte index after
// the match (in the input string), the keyword bytes to emit, and true
// if matched. Enforces the word-boundary rule from GTTOK6.
//
// On ties (e.g. IF appears at 0xD7 and 0xD8), the first table entry
// wins because we use a strict `<=` in the bestEnd comparison.
func lookupKeyword(input string, pos int) (canonical string, endPos int, bytes []byte, ok bool) {
	folded := keywordFold(input[pos:])
	bestEnd := 0
	bestIdx := -1
	for i, name := range keywordTable {
		if name == "" {
			continue
		}
		end, matched := matchKeyword(folded, name)
		if !matched {
			continue
		}
		if end <= bestEnd {
			continue
		}
		if !checkWordBoundary(folded, end, name) {
			continue
		}
		bestEnd = end
		bestIdx = i
	}
	if bestIdx < 0 {
		return "", pos, nil, false
	}
	canonical = keywordTable[bestIdx]
	endPos = pos + bestEnd
	keywordByte := byte(0x3B + bestIdx)
	if keywordByte >= 0x85 {
		bytes = []byte{keywordByte}
	} else {
		bytes = []byte{0xFF, keywordByte}
	}
	return canonical, endPos, bytes, true
}

// matchKeyword returns the byte length matched in input against the
// keyword name. Spaces in name match zero or one input space.
func matchKeyword(input, name string) (int, bool) {
	i := 0
	for j := 0; j < len(name); j++ {
		c := name[j]
		if c == ' ' {
			if i < len(input) && input[i] == ' ' {
				i++
			}
			continue
		}
		if i >= len(input) || input[i] != c {
			return 0, false
		}
		i++
	}
	return i, true
}

// checkWordBoundary enforces that the byte after the matched keyword is
// not a continuation of a longer identifier. Skips the check when the
// keyword's last char is `$`, `=`, or `>` (per GTTOK6).
func checkWordBoundary(input string, end int, name string) bool {
	if end >= len(input) {
		return true
	}
	lastCh := name[len(name)-1]
	if lastCh == '$' || lastCh == '=' || lastCh == '>' {
		return true
	}
	next := input[end]
	if next >= 'A' && next <= 'Z' {
		return false
	}
	if next >= '0' && next <= '9' {
		return false
	}
	if next == '_' {
		return false
	}
	return true
}

// advanceColOver walks the input bytes from oldPos to newPos updating
// l.line and l.col so that emit's start-of-token stamp stays sensible
// after a keyword match that jumped pos directly.
func (l *lexer) advanceColOver(oldPos, newPos int) {
	for i := oldPos; i < newPos; i++ {
		if l.input[i] == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
	}
}

// lexKeyword scans an alphabetic run and tries to tokenise it as a
// keyword. On match: emits itemKeyword with the resolved bytes and
// applies the one-trailing-space drop. The one-leading-space drop is
// handled in lexBodyLoop by NOT emitting a single space when followed
// by a keyword.
//
// On no-match: emits the next single byte of the alphabetic run as
// itemLiteral and returns self (to preserve the yield-after-emit rule)
// until the run is exhausted. Task 14 will replace this with
// bare-identifier handling.
func lexKeyword(l *lexer) stateFn {
	canonical, endPos, kwBytes, ok := lookupKeyword(l.input, l.pos)
	if !ok {
		// Not a keyword. If we're at the start of a statement, this is a
		// bare-identifier PROC call: consume the whole alphanumeric run
		// and emit the identifier + 6-byte placeholder. Otherwise, emit
		// one byte and recurse (expression-context identifier path).
		if l.stmtInitial {
			// Consume the whole alphanumeric run.
			for {
				r := l.next()
				if r == eof || !isAlphaNum(r) {
					if r != eof {
						l.backup()
					}
					break
				}
			}
			identifier := l.input[l.start:l.pos]
			out := make([]byte, 0, len(identifier)+6)
			out = append(out, []byte(identifier)...)
			out = append(out, 0x0E, 0xFD, 0xFD, 0xFD, 0x00, 0x00)
			l.emitBytes(itemProcCallPlaceholder, out, identifier)
			l.stmtInitial = false
			return lexBodyLoop
		}
		// Expression-context: this is an identifier (not a keyword). Consume
		// the entire alphanumeric run and emit as a single itemLiteral, with
		// no further keyword-match attempts inside the run. Per TOKMAIN spec
		// §3.1-§3.2, keywords only match at word boundaries; we must not
		// re-attempt keyword lookup mid-identifier.
		for {
			r := l.next()
			if r == eof {
				break
			}
			if !isAlphaNum(r) {
				l.backup()
				break
			}
		}
		l.emit(itemLiteral)
		l.stmtInitial = false
		return lexBodyLoop
	}
	// Keyword match: jump pos to endPos, walking the consumed range to
	// keep col/line tracking accurate.
	l.advanceColOver(l.pos, endPos)
	l.pos = endPos
	// One-trailing-space drop.
	if l.pos < len(l.input) && l.input[l.pos] == ' ' {
		l.pos++
		l.col++
	}
	l.stmtInitial = false
	l.emitBytes(itemKeyword, kwBytes, l.input[l.start:l.pos])
	if canonical == "BIN" {
		return lexBinaryDigits
	}
	if canonical == "REM" {
		return lexComment
	}
	return lexBodyLoop
}

// lexBinaryDigits scans a run of binary digits (0 or 1) following a BIN
// keyword and emits an itemNumber with a 5-byte FP form. The run is
// limited to 16 bits per grammar spec §4.1.
func lexBinaryDigits(l *lexer) stateFn {
	// Skip any leading spaces between BIN and the digits. (The keyword
	// trailing-space-drop already consumed one space; additional ones
	// stay in the buffer.)
	for l.pos < len(l.input) && l.input[l.pos] == ' ' {
		l.pos++
		l.start = l.pos
		l.col++
		l.startLine = l.line
		l.startCol = l.col
	}
	const binDigits = "01"
	if !l.accept(binDigits) {
		return l.errorf("expected binary digits after BIN")
	}
	l.acceptRun(binDigits)
	literal := l.input[l.start:l.pos]
	if len(literal) > 16 {
		return l.errorf("binary literal too large")
	}
	if r := l.peek(); r != eof && (isAlpha(r) || r == '_') {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	var v uint64
	for _, c := range literal {
		v = (v << 1) | uint64(c-'0')
	}
	var fp [5]byte
	fp[2] = byte(v & 0xFF)
	fp[3] = byte((v >> 8) & 0xFF)
	out := make([]byte, 0, len(literal)+6)
	out = append(out, []byte(literal)...)
	out = append(out, 0x0E)
	out = append(out, fp[:]...)
	l.emitBytes(itemNumber, out, literal)
	return lexBodyLoop
}

// lexString scans a "..."-delimited string literal. Two consecutive "
// characters inside the string are stored verbatim (the doubled-quote
// escape that the run-time SQUOTE collapses to a single "). An
// unterminated string at end-of-line or end-of-input is accepted per
// E.4 empirical: store all bytes up to but not including the line
// terminator. Emits one itemString carrying the verbatim bytes
// including the opening " and (if present) the closing ".
func lexString(l *lexer) stateFn {
	// Consume opening quote.
	if r := l.next(); r != '"' {
		return l.errorf(`lexString entered without "`)
	}
	for {
		r := l.next()
		if r == eof || r == '\n' || r == '\r' {
			if r != eof {
				l.backup()
			}
			l.emit(itemString)
			return lexBodyLoop
		}
		if r == '"' {
			// Look ahead: another " means doubled-quote escape; consume
			// both and stay in string mode.
			if l.peek() == '"' {
				l.next()
				continue
			}
			// True closing quote.
			l.emit(itemString)
			return lexBodyLoop
		}
	}
}

// lexComment consumes the rest of the current line as a single raw
// itemLiteral (no keyword/number/string tokenisation inside), then
// hands back to a small finaliser state.
func lexComment(l *lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			// Emit literal here (if any), then return a finaliser that
			// emits the closing EOL and EOF on subsequent invocations.
			if l.pos > l.start {
				l.emit(itemLiteral)
				return lexCommentEndEOF
			}
			l.emit(itemEOL)
			l.emit(itemEOF)
			return nil
		}
		if r == '\n' || r == '\r' {
			l.backup()
			if l.pos > l.start {
				l.emit(itemLiteral)
				return lexCommentEndNewline
			}
			l.start = l.pos
			l.startLine = l.line
			l.startCol = l.col
			l.next()
			l.ignore()
			l.emit(itemEOL)
			return lexStart
		}
	}
}

// lexCommentEndEOF emits the closing EOL and EOF after a comment that
// ran to end-of-input.
func lexCommentEndEOF(l *lexer) stateFn {
	l.emit(itemEOL)
	l.emit(itemEOF)
	return nil
}

// lexCommentEndNewline consumes the line terminator after a comment
// body and emits the closing EOL, then returns to lexStart.
func lexCommentEndNewline(l *lexer) stateFn {
	l.start = l.pos
	l.startLine = l.line
	l.startCol = l.col
	l.next()
	l.ignore()
	l.emit(itemEOL)
	return lexStart
}

// lexControlEscape parses the {N} sugar where N is a decimal byte value
// 0..255. Emits itemControlEscape with the single raw byte.
func lexControlEscape(l *lexer) stateFn {
	if l.next() != '{' {
		return l.errorf("lexControlEscape entered without {")
	}
	digitStart := l.pos
	for {
		r := l.next()
		if r == eof || r == '\n' || r == '\r' {
			return l.errorf("unterminated control escape")
		}
		if r == '}' {
			break
		}
		if r < '0' || r > '9' {
			return l.errorf("invalid control escape: %q", l.input[l.start:l.pos])
		}
	}
	digits := l.input[digitStart : l.pos-1] // exclude closing }
	if digits == "" {
		return l.errorf("invalid control escape: %q", l.input[l.start:l.pos])
	}
	v, err := strconv.ParseUint(digits, 10, 16)
	if err != nil || v > 255 {
		return l.errorf("control escape out of range: %s", digits)
	}
	l.emitBytes(itemControlEscape, []byte{byte(v)}, l.input[l.start:l.pos])
	return lexBodyLoop
}
