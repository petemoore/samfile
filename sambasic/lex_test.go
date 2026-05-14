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
				{typ: itemLineNumber, val: "0"},
				{typ: itemEOL},
				{typ: itemEOF},
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
		{"GO TO single space", "10 GO TO 20\n", []byte{byte(GO_TO), '2', '0', 0x0E, 0x00, 0x00, 0x14, 0x00, 0x00}},
		{"GOTO no space", "10 GOTO 20\n", []byte{byte(GO_TO), '2', '0', 0x0E, 0x00, 0x00, 0x14, 0x00, 0x00}},
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

func TestLexNumber_Hex(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantVal string
		wantFP  [5]byte
		wantErr string
	}{
		{"hex-FF", "10 PRINT &FF\n", "&FF", [5]byte{0x00, 0x00, 0xFF, 0x00, 0x00}, ""},
		{"hex-lowercase", "10 PRINT &ff\n", "&ff", [5]byte{0x00, 0x00, 0xFF, 0x00, 0x00}, ""},
		{"hex-leading-zeros", "10 PRINT &0000FF\n", "&0000FF", [5]byte{0x00, 0x00, 0xFF, 0x00, 0x00}, ""},
		{"bare-amp", "10 PRINT &\n", "", [5]byte{}, "expected hex digits after &"},
		{"hex-then-letter", "10 PRINT &FFG\n", "", [5]byte{}, `bad number syntax: "&FFG"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			if tt.wantErr != "" {
				if got[len(got)-1].typ != itemError {
					t.Fatalf("expected itemError, got %v", got)
				}
				if got[len(got)-1].val != tt.wantErr {
					t.Errorf("error = %q, want %q", got[len(got)-1].val, tt.wantErr)
				}
				return
			}
			var num *item
			for i := range got {
				if got[i].typ == itemNumber {
					num = &got[i]
					break
				}
			}
			if num == nil {
				t.Fatalf("no itemNumber; got %v", got)
			}
			if num.val != tt.wantVal {
				t.Errorf("val = %q, want %q", num.val, tt.wantVal)
			}
			if len(num.bytes) < 5 {
				t.Fatalf("bytes too short: %v", num.bytes)
			}
			var got5 [5]byte
			copy(got5[:], num.bytes[len(num.bytes)-5:])
			if got5 != tt.wantFP {
				t.Errorf("FP bytes = % X, want % X", got5, tt.wantFP)
			}
		})
	}
}

func TestLexNumber_Decimal(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantVal string
		wantFP  [5]byte
		wantErr string
	}{
		{"int", "10 PRINT 42\n", "42", [5]byte{0x00, 0x00, 0x2A, 0x00, 0x00}, ""},
		{"zero", "10 PRINT 0\n", "0", [5]byte{0x00, 0x00, 0x00, 0x00, 0x00}, ""},
		{"max", "10 PRINT 65535\n", "65535", [5]byte{0x00, 0x00, 0xFF, 0xFF, 0x00}, ""},
		{"trailing-letter", "10 PRINT 1G\n", "", [5]byte{}, `bad number syntax: "1G"`},
		{"trailing-colon-ok", "10 PRINT 1:PRINT 2\n", "", [5]byte{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			if tt.wantErr != "" {
				if got[len(got)-1].typ != itemError {
					t.Fatalf("expected itemError, got %v", got)
				}
				if got[len(got)-1].val != tt.wantErr {
					t.Errorf("error = %q, want %q", got[len(got)-1].val, tt.wantErr)
				}
				return
			}
			var num *item
			for i := range got {
				if got[i].typ == itemNumber {
					num = &got[i]
					break
				}
			}
			if tt.wantVal == "" {
				return
			}
			if num == nil {
				t.Fatalf("no itemNumber emitted; got %v", got)
			}
			if num.val != tt.wantVal {
				t.Errorf("val = %q, want %q", num.val, tt.wantVal)
			}
			if len(num.bytes) < 5 {
				t.Fatalf("bytes too short: %v", num.bytes)
			}
			var got5 [5]byte
			copy(got5[:], num.bytes[len(num.bytes)-5:])
			if got5 != tt.wantFP {
				t.Errorf("FP bytes = % X, want % X", got5, tt.wantFP)
			}
		})
	}
}

func TestLexNumber_LeadingDot(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantVal string
	}{
		{"leading-dot", "10 PRINT .5\n", ".5"},
		{"leading-dot-int-part", "10 PRINT .32\n", ".32"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			var num *item
			for i := range got {
				if got[i].typ == itemNumber {
					num = &got[i]
					break
				}
			}
			if num == nil {
				t.Fatalf("no itemNumber emitted; got %v", got)
			}
			if num.val != tt.wantVal {
				t.Errorf("val = %q, want %q", num.val, tt.wantVal)
			}
		})
	}
}

func TestLexNumber_Scientific(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantVal string
		wantErr string
	}{
		{"basic", "10 PRINT 1E5\n", "1E5", ""},
		{"with-sign", "10 PRINT 1E+5\n", "1E+5", ""},
		{"negative-exp", "10 PRINT 1E-3\n", "1E-3", ""},
		{"with-fraction", "10 PRINT 1.5E3\n", "1.5E3", ""},
		{"lowercase-e", "10 PRINT 1e5\n", "1e5", ""},
		{"trailing-letter", "10 PRINT 1E5G\n", "", `bad number syntax: "1E5G"`},
		{"incomplete", "10 PRINT 1E:PRINT 2\n", "", `bad number syntax: "1E"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			if tt.wantErr != "" {
				if got[len(got)-1].typ != itemError {
					t.Fatalf("expected itemError, got %v", got)
				}
				if got[len(got)-1].val != tt.wantErr {
					t.Errorf("error = %q, want %q", got[len(got)-1].val, tt.wantErr)
				}
				return
			}
			var num *item
			for i := range got {
				if got[i].typ == itemNumber {
					num = &got[i]
					break
				}
			}
			if num == nil {
				t.Fatalf("no itemNumber; got %v", got)
			}
			if num.val != tt.wantVal {
				t.Errorf("val = %q, want %q", num.val, tt.wantVal)
			}
		})
	}
}

func TestLexNumber_Binary(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantVal string
		wantFP  [5]byte
		wantErr string
	}{
		{"basic", "10 PRINT BIN 1010\n", "1010", [5]byte{0x00, 0x00, 0x0A, 0x00, 0x00}, ""},
		{"zero", "10 PRINT BIN 0\n", "0", [5]byte{0x00, 0x00, 0x00, 0x00, 0x00}, ""},
		{"max-16-bits", "10 PRINT BIN 1111111111111111\n", "1111111111111111", [5]byte{0x00, 0x00, 0xFF, 0xFF, 0x00}, ""},
		{"too-many-bits", "10 PRINT BIN 11111111111111111\n", "", [5]byte{}, "binary literal too large"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			if tt.wantErr != "" {
				if got[len(got)-1].typ != itemError {
					t.Fatalf("expected itemError, got %v", got)
				}
				if got[len(got)-1].val != tt.wantErr {
					t.Errorf("error = %q, want %q", got[len(got)-1].val, tt.wantErr)
				}
				return
			}
			var num *item
			for i := range got {
				if got[i].typ == itemNumber {
					num = &got[i]
					break
				}
			}
			if num == nil {
				t.Fatalf("no itemNumber; got %v", got)
			}
			if num.val != tt.wantVal {
				t.Errorf("val = %q, want %q", num.val, tt.wantVal)
			}
			var got5 [5]byte
			copy(got5[:], num.bytes[len(num.bytes)-5:])
			if got5 != tt.wantFP {
				t.Errorf("FP bytes = % X, want % X", got5, tt.wantFP)
			}
		})
	}
}

func TestLexComment_REM(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantTail []byte
	}{
		{"basic", "10 REM hello\n", []byte("hello")},
		{"with-keyword-inside", "10 REM PRINT \"x\"\n", []byte(`PRINT "x"`)},
		{"with-quotes-and-numbers", "10 REM \"a\":42\n", []byte(`"a":42`)},
		{"bare-REM", "10 REM\n", []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			var tail []byte
			afterREM := false
			for _, it := range got {
				if it.typ == itemKeyword && len(it.bytes) == 1 && it.bytes[0] == byte(REM) {
					afterREM = true
					continue
				}
				if !afterREM {
					continue
				}
				if it.typ == itemEOL {
					break
				}
				if it.bytes != nil {
					tail = append(tail, it.bytes...)
				} else {
					tail = append(tail, []byte(it.val)...)
				}
			}
			if !bytesEqual(tail, tt.wantTail) {
				t.Errorf("tail = %q, want %q", tail, tt.wantTail)
			}
		})
	}
}

// TestLexComment_BytePreservingNonASCII guards against UTF-8 mangling in
// REM bodies: a byte > 0x7F that is not a valid UTF-8 prefix would,
// under the old `byte(l.next())` path, decode to utf8.RuneError (0xFFFD)
// and truncate to 0xFD. The detokeniser emits raw bytes from disk verbatim
// — e.g. the SAM glyph at 0x82 in `REM Door<0x82> 1992/3` — so the REM
// body lexer must walk bytes, not runes. Corpus example: "Banzai - The
// Demos & Utils by Dan Doore (1994).mgt"/DigiShow line 20 offset 83.
func TestLexComment_BytePreservingNonASCII(t *testing.T) {
	in := "10 REM Door\x82 1992/3\n"
	wantTail := []byte{'D', 'o', 'o', 'r', 0x82, ' ', '1', '9', '9', '2', '/', '3'}

	got := collectItems(in)
	var tail []byte
	afterREM := false
	for _, it := range got {
		if it.typ == itemKeyword && len(it.bytes) == 1 && it.bytes[0] == byte(REM) {
			afterREM = true
			continue
		}
		if !afterREM {
			continue
		}
		if it.typ == itemEOL {
			break
		}
		if it.bytes != nil {
			tail = append(tail, it.bytes...)
		} else {
			tail = append(tail, []byte(it.val)...)
		}
	}
	if !bytesEqual(tail, wantTail) {
		t.Errorf("REM tail = % X, want % X", tail, wantTail)
	}
}

// TestLexString_BytePreservingNonASCII guards the same byte-preservation
// invariant inside string literals: e.g. `"X\x80Y"` must round-trip as
// `"`, 'X', 0x80, 'Y', `"` — not `"`, 'X', 0xFD, 'Y', `"`. Corpus
// examples in the `String | other | got=0xFD want=0x80..0x84` buckets.
func TestLexString_BytePreservingNonASCII(t *testing.T) {
	in := "10 \"X\x80Y\x82Z\"\n"
	want := []byte{'"', 'X', 0x80, 'Y', 0x82, 'Z', '"'}

	got := collectItems(in)
	var s []byte
	for _, it := range got {
		if it.typ == itemString {
			if it.bytes != nil {
				s = append(s, it.bytes...)
			} else {
				s = append(s, []byte(it.val)...)
			}
		}
	}
	if !bytesEqual(s, want) {
		t.Errorf("string bytes = % X, want % X", s, want)
	}
}

func TestLexKeyword_NoMidIdentifierKeyword(t *testing.T) {
	// In expression context, 'crem' must not be tokenised as 'c'+'REM'.
	// All five letters of "crem" should emit as plain literal bytes.
	got := collectItems("10 LET crem=1\n")
	// Find the bytes between LET keyword and '=' in the body
	var sawREM bool
	for _, it := range got {
		if it.typ == itemKeyword && len(it.bytes) == 1 && it.bytes[0] == byte(REM) {
			// REM keyword should NOT appear in this input
			sawREM = true
			break
		}
	}
	if sawREM {
		t.Error("'crem' was wrongly tokenised with REM keyword inside")
	}
}

func TestLexProcCall_Placeholder(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		wantBytes []byte
	}{
		{"bare-X", "10 X\n", []byte{'X', 0x0E, 0xFD, 0xFD, 0xFD, 0x00, 0x00}},
		{"two-procs", "10 X:Y\n", []byte{'X', 0x0E, 0xFD, 0xFD, 0xFD, 0x00, 0x00, ':', 'Y', 0x0E, 0xFD, 0xFD, 0xFD, 0x00, 0x00}},
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
				t.Errorf("body = % X, want % X", body, tt.wantBytes)
			}
		})
	}
}

func TestLexString_ControlEscape(t *testing.T) {
	got := collectItems("10 PRINT \"a{7}b\"\n")
	// Find the itemString and verify its bytes contain the literal 0x07.
	// (val carries the verbatim source for error messages; bytes carries
	// the resolved on-disk byte sequence.)
	var found bool
	for _, it := range got {
		if it.typ == itemString {
			want := []byte{0x22, 'a', 0x07, 'b', 0x22}
			if bytesEqual(it.bytes, want) {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("string with {7} did not produce expected bytes; got %v", got)
	}
}

func TestLexComment_ControlEscape(t *testing.T) {
	got := collectItems("10 REM x{7}y\n")
	// Concatenate all body bytes after REM keyword.
	var tail []byte
	afterREM := false
	for _, it := range got {
		if it.typ == itemKeyword && len(it.bytes) == 1 && it.bytes[0] == byte(REM) {
			afterREM = true
			continue
		}
		if !afterREM {
			continue
		}
		if it.typ == itemEOL {
			break
		}
		if it.bytes != nil {
			tail = append(tail, it.bytes...)
		} else {
			tail = append(tail, []byte(it.val)...)
		}
	}
	want := []byte{'x', 0x07, 'y'}
	if !bytesEqual(tail, want) {
		t.Errorf("REM tail = % X, want % X", tail, want)
	}
}

func TestLexControlEscape(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		wantBytes []byte
	}{
		{"in-body", "10 {7}\n", []byte{0x07}},
		{"max-255", "10 {255}\n", []byte{0xFF}},
		{"zero", "10 {0}\n", []byte{0x00}},
		{"invalid-not-decimal", "10 {abc}\n", nil},
		{"invalid-out-of-range", "10 {256}\n", nil},
		{"unterminated", "10 {7\n", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectItems(tt.in)
			if tt.wantBytes == nil {
				if got[len(got)-1].typ != itemError {
					t.Fatalf("expected itemError, got %v", got)
				}
				return
			}
			var bytes []byte
			i := 0
			for i < len(got) && got[i].typ != itemLineNumber {
				i++
			}
			i++
			for i < len(got) && got[i].typ != itemEOL {
				if got[i].bytes != nil {
					bytes = append(bytes, got[i].bytes...)
				}
				i++
			}
			if !bytesEqual(bytes, tt.wantBytes) {
				t.Errorf("bytes = %v, want %v", bytes, tt.wantBytes)
			}
		})
	}
}
