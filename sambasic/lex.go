package sambasic

import (
	"fmt"
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
	input     string
	pos       int
	start     int
	width     int
	line      int
	col       int
	startLine int
	startCol  int
	state     stateFn
	items     chan item
}

func lex(input string) *lexer {
	l := &lexer{
		input:     input,
		line:      1,
		col:       1,
		startLine: 1,
		startCol:  1,
		items:     make(chan item, 2),
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
	l.emit(itemLiteral)
	return lexBodyLoop
}
