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
