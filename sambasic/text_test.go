package sambasic

import (
	"strings"
	"testing"
)

func TestParseText_BasicLine(t *testing.T) {
	f, err := ParseTextString(`10 PRINT "hello"` + "\n")
	if err != nil {
		t.Fatalf("ParseText error: %v", err)
	}
	if len(f.Lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(f.Lines))
	}
	if f.Lines[0].Number != 10 {
		t.Errorf("line number = %d, want 10", f.Lines[0].Number)
	}
}

func TestParseText_EditSemantics(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantNum []uint16
	}{
		{"sort", "10 X\n5 Y\n15 Z\n", []uint16{5, 10, 15}},
		{"dedup-last-wins", "10 X\n10 Y\n", []uint16{10}},
		{"bare-deletes", "10 X\n10\n", nil},
		{"bare-space-preserves", "10 X\n10 \n", []uint16{10}},
		{"empty-input", "", nil},
		{"whitespace-only-input", "  \n\t\n", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := ParseTextString(tt.in)
			if err != nil {
				t.Fatalf("ParseText error: %v", err)
			}
			if len(f.Lines) != len(tt.wantNum) {
				t.Fatalf("got %d lines, want %d", len(f.Lines), len(tt.wantNum))
			}
			for i, n := range tt.wantNum {
				if f.Lines[i].Number != n {
					t.Errorf("line[%d] = %d, want %d", i, f.Lines[i].Number, n)
				}
			}
		})
	}
}

func TestParseText_LineNumberOutOfRange(t *testing.T) {
	_, err := ParseTextString("65280 PRINT \"x\"\n")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	pe, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("error type = %T, want *ParseError", err)
	}
	if !strings.Contains(pe.Msg, "out of range") {
		t.Errorf("error msg %q does not mention 'out of range'", pe.Msg)
	}
}

func TestFinalise_IfThenPatch(t *testing.T) {
	f, err := ParseTextString("10 IF A=1 THEN PRINT \"x\"\n")
	if err != nil {
		t.Fatalf("ParseText: %v", err)
	}
	bytes := f.Lines[0].Bytes()
	if len(bytes) < 5 {
		t.Fatalf("body too short: %v", bytes)
	}
	if bytes[4] != 0xD8 {
		t.Errorf("first body byte = %#x, want 0xD8 (SIF)", bytes[4])
	}
}

func TestFinalise_IfNoThen(t *testing.T) {
	f, err := ParseTextString("10 IF A=1: PRINT \"x\"\n")
	if err != nil {
		t.Fatalf("ParseText: %v", err)
	}
	bytes := f.Lines[0].Bytes()
	if bytes[4] != 0xD7 {
		t.Errorf("first body byte = %#x, want 0xD7 (LIF)", bytes[4])
	}
}

func TestFinalise_ElsePatch(t *testing.T) {
	f, err := ParseTextString("10 IF A=1 THEN PRINT \"x\":ELSE PRINT \"y\"\n")
	if err != nil {
		t.Fatalf("ParseText: %v", err)
	}
	bytes := f.Lines[0].Bytes()
	var found byte
	for i := 4; i < len(bytes); i++ {
		if bytes[i] == 0xD9 || bytes[i] == 0xDA {
			found = bytes[i]
			break
		}
	}
	if found != 0xDA {
		t.Errorf("ELSE byte = %#x, want 0xDA", found)
	}
}

func TestFinalise_InkPatch(t *testing.T) {
	f, err := ParseTextString("10 INK 2\n")
	if err != nil {
		t.Fatalf("ParseText: %v", err)
	}
	bytes := f.Lines[0].Bytes()
	if bytes[4] != 0xA1 {
		t.Errorf("first body byte = %#x, want 0xA1 (PEN)", bytes[4])
	}
}

func TestParseText_BuildDiskAutoRunFixture(t *testing.T) {
	src := `10 CLEAR 32767
20 LOAD "stub" CODE 32768
30 CALL 32768
`
	got, err := ParseTextString(src)
	if err != nil {
		t.Fatalf("ParseText: %v", err)
	}

	want := &File{
		Lines: []Line{
			{Number: 10, Tokens: []Token{CLEAR, Number(32767)}},
			{Number: 20, Tokens: []Token{LOAD, String(`"stub"`), CODE, Number(32768)}},
			{Number: 30, Tokens: []Token{CALL, Number(32768)}},
		},
	}

	gotBytes := got.ProgBytes()
	wantBytes := want.ProgBytes()

	if !bytesEqual(gotBytes, wantBytes) {
		t.Errorf("ProgBytes mismatch\ngot:  % X\nwant: % X", gotBytes, wantBytes)
	}
}
