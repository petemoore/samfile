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
