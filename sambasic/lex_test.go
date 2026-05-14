package sambasic

import (
	"testing"
)

func TestLexerPrimitives_NextBackupPeek(t *testing.T) {
	l := lex("AB")
	if r := l.next(); r != 'A' {
		t.Fatalf("next() = %q, want 'A'", r)
	}
	if r := l.peek(); r != 'B' {
		t.Fatalf("peek() = %q, want 'B'", r)
	}
	if r := l.next(); r != 'B' {
		t.Fatalf("next() = %q, want 'B'", r)
	}
	l.backup()
	if r := l.peek(); r != 'B' {
		t.Fatalf("peek after backup = %q, want 'B'", r)
	}
	if r := l.next(); r != 'B' {
		t.Fatalf("re-consumed = %q, want 'B'", r)
	}
	if r := l.next(); r != eof {
		t.Fatalf("at end, next() = %q, want eof", r)
	}
}

func TestLexerPrimitives_LineColTracking(t *testing.T) {
	l := lex("a\nbc")
	l.next() // 'a': line 1 col 1
	if l.line != 1 || l.col != 2 {
		t.Fatalf("after 'a': line=%d col=%d, want 1,2", l.line, l.col)
	}
	l.next() // '\n': advances line
	if l.line != 2 || l.col != 1 {
		t.Fatalf("after '\\n': line=%d col=%d, want 2,1", l.line, l.col)
	}
	l.next() // 'b'
	if l.line != 2 || l.col != 2 {
		t.Fatalf("after 'b': line=%d col=%d, want 2,2", l.line, l.col)
	}
}

func TestLexerPrimitives_AcceptAcceptRun(t *testing.T) {
	l := lex("0123abc")
	if !l.accept("0123456789") {
		t.Fatal("accept digit failed")
	}
	l.acceptRun("0123456789")
	if l.pos != 4 {
		t.Fatalf("pos after digits = %d, want 4", l.pos)
	}
	if !l.accept("abc") {
		t.Fatal("accept letter failed")
	}
	if l.accept("0123") {
		t.Fatal("accept after non-match should fail")
	}
}

func TestLexerPrimitives_EmitAndDrain(t *testing.T) {
	l := lex("hello")
	go func() {
		l.start = 0
		l.pos = 5
		l.emit(itemLiteral)
		close(l.items)
	}()
	it := <-l.items
	if it.typ != itemLiteral || it.val != "hello" {
		t.Fatalf("got %+v, want literal 'hello'", it)
	}
}

func TestLexerPrimitives_EmitStampsStartPosition(t *testing.T) {
	l := lex("abc\nXY")
	// Consume "abc"
	l.next()
	l.next()
	l.next()
	// At this point line=1, col=4. The token starts at line=1, col=1.
	go func() {
		l.emit(itemLiteral)
		close(l.items)
	}()
	it := <-l.items
	if it.line != 1 || it.col != 1 {
		t.Errorf("emit stamped line=%d col=%d, want 1,1 (start of token)", it.line, it.col)
	}
	if it.val != "abc" {
		t.Errorf("val = %q, want 'abc'", it.val)
	}
}

func TestLexerPrimitives_ErrorfStampsCurrentPosition(t *testing.T) {
	// errorf should stamp the CURRENT position (where the bad char is),
	// not the start-of-token position.
	l := lex("abXY")
	l.next()
	l.next()
	l.next()
	// line=1, col=4 (just past 'X')
	go func() {
		l.errorf("test")
		close(l.items)
	}()
	it := <-l.items
	if it.line != 1 || it.col != 4 {
		t.Errorf("errorf stamped line=%d col=%d, want 1,4 (current position)", it.line, it.col)
	}
}

func TestLexerPrimitives_Errorf(t *testing.T) {
	l := lex("xyz")
	l.line = 5
	l.col = 7
	go func() {
		l.errorf("bad thing %q", "z")
		close(l.items)
	}()
	it := <-l.items
	if it.typ != itemError {
		t.Fatalf("got typ %v, want itemError", it.typ)
	}
	if it.val != `bad thing "z"` {
		t.Fatalf("got val %q, want 'bad thing \"z\"'", it.val)
	}
	if it.line != 5 || it.col != 7 {
		t.Fatalf("got line/col %d/%d, want 5/7", it.line, it.col)
	}
}

// collectItems runs the lexer's state machine to completion and returns all
// items emitted (including the final itemEOF or itemError).
func collectItems(input string) []item {
	l := lex(input)
	l.state = lexStart
	var items []item
	for {
		it := l.nextItem()
		items = append(items, it)
		if it.typ == itemEOF || it.typ == itemError {
			return items
		}
	}
}

func TestLexLineNumber(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []item
	}{
		{
			name: "simple",
			in:   "10\n",
			want: []item{
				{typ: itemLineNumber, val: "10"},
				{typ: itemEOL},
				{typ: itemEOF},
			},
		},
		{
			name: "leading whitespace",
			in:   "   42\n",
			want: []item{
				{typ: itemLineNumber, val: "42"},
				{typ: itemEOL},
				{typ: itemEOF},
			},
		},
		{
			name: "leading zeros",
			in:   "0010\n",
			want: []item{
				{typ: itemLineNumber, val: "0010"},
				{typ: itemEOL},
				{typ: itemEOF},
			},
		},
		{
			name: "max",
			in:   "65279\n",
			want: []item{
				{typ: itemLineNumber, val: "65279"},
				{typ: itemEOL},
				{typ: itemEOF},
			},
		},
		{
			name: "above max",
			in:   "65280\n",
			want: []item{
				{typ: itemError, val: "line number out of range: 65280"},
			},
		},
		{
			name: "zero",
			in:   "0\n",
			want: []item{
				{typ: itemError, val: "line number 0 is reserved"},
			},
		},
		{
			name: "no digits",
			in:   "PRINT\n",
			want: []item{
				{typ: itemError, val: "expected line number"},
			},
		},
		{
			name: "empty input",
			in:   "",
			want: []item{
				{typ: itemEOF},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d items, want %d: got=%v", len(got), len(tt.want), got)
			}
			for i, w := range tt.want {
				if got[i].typ != w.typ || got[i].val != w.val {
					t.Errorf("item[%d] = %+v, want typ=%v val=%q", i, got[i], w.typ, w.val)
				}
			}
		})
	}
}

func TestLexBody_OneSpaceDrop(t *testing.T) {
	tests := []struct {
		name string
		in   string
		// literalBytes is what we expect the body to emit before EOL,
		// represented as the concatenation of every emitted item's val.
		literalBytes string
	}{
		{"no-space-then-X", "10X\n", "X"},
		{"one-space-X-dropped", "10 X\n", "X"},
		{"two-spaces-one-dropped", "10  X\n", " X"},
		{"one-space-bare-preserved", "10 \n", " "},
		{"bare-no-body", "10\n", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			// Skip itemLineNumber, accumulate val of intervening items
			// until itemEOL.
			var body string
			i := 0
			for i < len(got) && got[i].typ != itemLineNumber {
				i++
			}
			i++ // past line number
			for i < len(got) && got[i].typ != itemEOL {
				body += got[i].val
				i++
			}
			if body != tt.literalBytes {
				t.Errorf("body = %q, want %q", body, tt.literalBytes)
			}
		})
	}
}

func TestLexKeyword_Basic(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		wantBytes []byte
	}{
		{"PRINT", "10 PRINT \"hi\"\n", []byte{byte(PRINT), '"', 'h', 'i', '"'}},
		{"print-lowercase", "10 print \"hi\"\n", []byte{byte(PRINT), '"', 'h', 'i', '"'}},
		{"Print-mixed", "10 Print \"hi\"\n", []byte{byte(PRINT), '"', 'h', 'i', '"'}},
		{"LOAD-no-space-before-quote", "10 LOAD\"foo\"\n", []byte{byte(LOAD), '"', 'f', 'o', 'o', '"'}},
		{"GO TO single space", "10 GO TO 20\n", []byte{byte(GO_TO), '2', '0'}},
		{"GOTO no space", "10 GOTO 20\n", []byte{byte(GO_TO), '2', '0'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			var body []byte
			i := 0
			for i < len(got) && got[i].typ != itemLineNumber {
				i++
			}
			i++
			for i < len(got) && got[i].typ != itemEOL {
				if got[i].bytes != nil {
					body = append(body, got[i].bytes...)
				} else {
					body = append(body, []byte(got[i].val)...)
				}
				i++
			}
			if !bytesEqual(body, tt.wantBytes) {
				t.Errorf("body = %v, want %v", body, tt.wantBytes)
			}
		})
	}
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestLexKeyword_SpaceDrop(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		wantN int // total body byte count (excluding EOL marker)
	}{
		// PRINT(1) + "(1) + a(1) + "(1) + :(1) + PRINT(1) + "(1) + b(1) + "(1) = 9
		{"colon-no-space", "10 PRINT \"a\":PRINT \"b\"\n", 9},
		// 9 + 1 space before colon = 10
		{"space-before-colon", "10 PRINT \"a\" :PRINT \"b\"\n", 10},
		{"two-colons", "10 :: PRINT \"x\"\n", 6}, // : : PRINT " x " = 6
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			var n int
			i := 0
			for i < len(got) && got[i].typ != itemLineNumber {
				i++
			}
			i++
			for i < len(got) && got[i].typ != itemEOL {
				if got[i].bytes != nil {
					n += len(got[i].bytes)
				} else {
					n += len(got[i].val)
				}
				i++
			}
			if n != tt.wantN {
				t.Errorf("body byte count = %d, want %d", n, tt.wantN)
			}
		})
	}
}

func TestLexString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string // expected concatenation of itemString val(s)
	}{
		{"empty-string", `10 ""` + "\n", `""`},
		{"simple", `10 "hello"` + "\n", `"hello"`},
		{"doubled-quote", `10 "a""b"` + "\n", `"a""b"`},
		{"three-quotes-unterminated", `10 """` + "\n", `"""`},
		{"unterminated-eol", `10 "abc` + "\n", `"abc`},
		{"unterminated-eof", `10 "abc`, `"abc`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			var s string
			for _, it := range got {
				if it.typ == itemString {
					s += it.val
				}
			}
			if s != tt.want {
				t.Errorf("string item val = %q, want %q", s, tt.want)
			}
		})
	}
}
