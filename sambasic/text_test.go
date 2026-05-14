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
