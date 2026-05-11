package sambasic

import (
	"bytes"
	"testing"
)

func TestSingleByteKeywordBytes(t *testing.T) {
	k := SingleByteKeyword(0xB3) // CLEAR
	got := k.Bytes()
	want := []byte{0xB3}
	if !bytes.Equal(got, want) {
		t.Errorf("SingleByteKeyword(0xB3).Bytes() = %x; want %x", got, want)
	}
}

func TestTwoByteKeywordBytes(t *testing.T) {
	k := TwoByteKeyword(0x6C) // CODE
	got := k.Bytes()
	want := []byte{0xFF, 0x6C}
	if !bytes.Equal(got, want) {
		t.Errorf("TwoByteKeyword(0x6C).Bytes() = %x; want %x", got, want)
	}
}

func TestNumBytes(t *testing.T) {
	n := Number(32767)
	got := n.Bytes()
	want := []byte{
		'3', '2', '7', '6', '7',
		0x0E, 0x00, 0x00, 0xFF, 0x7F, 0x00,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Number(32767).Bytes() = %x; want %x", got, want)
	}
}

func TestNumBytesZero(t *testing.T) {
	n := Number(0)
	got := n.Bytes()
	want := []byte{
		'0',
		0x0E, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Number(0).Bytes() = %x; want %x", got, want)
	}
}

func TestNumBytes32768(t *testing.T) {
	n := Number(32768)
	got := n.Bytes()
	want := []byte{
		'3', '2', '7', '6', '8',
		0x0E, 0x00, 0x00, 0x00, 0x80, 0x00,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Number(32768).Bytes() = %x; want %x", got, want)
	}
}

func TestStrBytes(t *testing.T) {
	s := String("stub")
	got := s.Bytes()
	want := []byte("stub")
	if !bytes.Equal(got, want) {
		t.Errorf("String(\"stub\").Bytes() = %x; want %x", got, want)
	}
}

func TestLiteralBytes(t *testing.T) {
	l := Literal(':')
	got := l.Bytes()
	want := []byte{0x3A}
	if !bytes.Equal(got, want) {
		t.Errorf("Literal(':').Bytes() = %x; want %x", got, want)
	}
}
