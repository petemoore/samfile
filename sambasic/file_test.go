package sambasic

import (
	"bytes"
	"testing"
)

func TestLineBytes(t *testing.T) {
	line := Line{
		Number: 10,
		Tokens: []Token{
			CLEAR,
			Number(32767),
		},
	}
	got := line.Bytes()
	want := []byte{
		0x00, 0x0a, // line number 10 big-endian
		0x0d, 0x00, // body length 13 little-endian
		0xb3,                               // CLEAR
		'3', '2', '7', '6', '7',           // display "32767"
		0x0e, 0x00, 0x00, 0xff, 0x7f, 0x00, // numeric form
		0x0d, // line terminator
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Line.Bytes():\n  got  %x\n  want %x", got, want)
	}
}

func TestBuildDiskAutoProgram(t *testing.T) {
	f := &File{
		Lines: []Line{
			{
				Number: 10,
				Tokens: []Token{
					CLEAR,
					Number(32767),
					Literal(':'),
					LOAD,
					Literal('"'),
					String("stub"),
					Literal('"'),
					CODE,
					Number(32768),
					Literal(':'),
					CALL,
					Number(32768),
				},
			},
		},
		StartLine: 10,
	}

	wantProg := []byte{
		0x00, 0x0a, 0x2f, 0x00, // line 10, body len 47
		0xb3,                                // CLEAR
		'3', '2', '7', '6', '7',            // "32767"
		0x0e, 0x00, 0x00, 0xff, 0x7f, 0x00, // num(32767)
		0x3a,                                // :
		0x95,                                // LOAD
		0x22,                                // "
		's', 't', 'u', 'b',                 // stub
		0x22,                                // "
		0xff, 0x6c,                          // CODE (two-byte)
		'3', '2', '7', '6', '8',            // "32768"
		0x0e, 0x00, 0x00, 0x00, 0x80, 0x00, // num(32768)
		0x3a,                                // :
		0xe4,                                // CALL
		'3', '2', '7', '6', '8',            // "32768"
		0x0e, 0x00, 0x00, 0x00, 0x80, 0x00, // num(32768)
		0x0d, // line terminator
		0xff, // end-of-program sentinel
	}

	gotProg := f.ProgBytes()
	if !bytes.Equal(gotProg, wantProg) {
		t.Errorf("ProgBytes():\n  got  %x\n  want %x", gotProg, wantProg)
	}

	gotBody := f.Bytes()
	wantBody := make([]byte, len(wantProg)+92+512)
	copy(wantBody, wantProg)
	if !bytes.Equal(gotBody, wantBody) {
		t.Errorf("Bytes() length = %d; want %d", len(gotBody), len(wantBody))
		if len(gotBody) >= len(wantProg) && len(wantBody) >= len(wantProg) {
			if !bytes.Equal(gotBody[:len(wantProg)], wantProg) {
				t.Errorf("PROG section mismatch")
			}
		}
	}
	if len(gotBody) != 656 {
		t.Errorf("body length = %d; want 656", len(gotBody))
	}
}

func TestFileOffsets(t *testing.T) {
	f := &File{
		Lines: []Line{
			{Number: 10, Tokens: []Token{CLEAR, Number(32767)}},
		},
	}
	progLen := uint32(len(f.ProgBytes()))
	if f.NVARSOffset() != progLen {
		t.Errorf("NVARSOffset() = %d; want %d", f.NVARSOffset(), progLen)
	}
	if f.NUMENDOffset() != progLen+92 {
		t.Errorf("NUMENDOffset() = %d; want %d", f.NUMENDOffset(), progLen+92)
	}
	if f.SAVARSOffset() != progLen+92+512 {
		t.Errorf("SAVARSOffset() = %d; want %d", f.SAVARSOffset(), progLen+92+512)
	}
}

func TestFileOffsetsCustomVars(t *testing.T) {
	f := &File{
		Lines:       []Line{{Number: 1, Tokens: []Token{RUN}}},
		NumericVars: make([]byte, 200),
		Gap:         make([]byte, 1024),
	}
	progLen := uint32(len(f.ProgBytes()))
	if f.NVARSOffset() != progLen {
		t.Errorf("NVARSOffset() = %d; want %d", f.NVARSOffset(), progLen)
	}
	if f.NUMENDOffset() != progLen+200 {
		t.Errorf("NUMENDOffset() = %d; want %d", f.NUMENDOffset(), progLen+200)
	}
	if f.SAVARSOffset() != progLen+200+1024 {
		t.Errorf("SAVARSOffset() = %d; want %d", f.SAVARSOffset(), progLen+200+1024)
	}
	bodyLen := uint32(len(f.Bytes()))
	if bodyLen != progLen+200+1024 {
		t.Errorf("Bytes() length = %d; want %d", bodyLen, progLen+200+1024)
	}
}
