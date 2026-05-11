package sambasic

import (
	"bytes"
	"testing"
)

func TestKeywordConstants(t *testing.T) {
	cases := []struct {
		name  string
		token Token
		want  []byte
	}{
		{"PI", PI_2B, []byte{0xFF, 0x3B}},
		{"CLEAR", CLEAR, []byte{0xB3}},
		{"LOAD", LOAD, []byte{0x95}},
		{"CALL", CALL, []byte{0xE4}},
		{"IF_SHORT", IF_SHORT, []byte{0xD8}},
		{"ZOOM", ZOOM, []byte{0xF6}},
		{"CODE", CODE, []byte{0xFF, 0x6C}},
		{"UDG", UDG, []byte{0xFF, 0x69}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.token.Bytes()
			if !bytes.Equal(got, c.want) {
				t.Errorf("%s.Bytes() = %x; want %x", c.name, got, c.want)
			}
		})
	}
}

func TestKeywordName(t *testing.T) {
	cases := []struct {
		name     string
		token    byte
		extended bool
		want     string
		wantOK   bool
	}{
		{"CLEAR is 0xB3", 0xB3, false, "CLEAR", true},
		{"USING is 0x85", 0x85, false, "USING", true},
		{"CODE is 0xFF+0x6C", 0x6C, true, "CODE", true},
		{"PI is 0xFF+0x3B", 0x3B, true, "PI", true},
		{"below range", 0x20, false, "", false},
		{"above range", 0xF7, false, "", false},
		{"reserved 0x84", 0x84, false, "", false},
		{"reserved via extended", 0x49, true, "", false},
		{"single-byte below 0x85", 0x3B, false, "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := KeywordName(c.token, c.extended)
			if ok != c.wantOK {
				t.Errorf("KeywordName(0x%02x, %v) ok = %v; want %v", c.token, c.extended, ok, c.wantOK)
			}
			if got != c.want {
				t.Errorf("KeywordName(0x%02x, %v) = %q; want %q", c.token, c.extended, got, c.want)
			}
		})
	}
}

func TestKeywordTableLength(t *testing.T) {
	if len(keywordTable) != 188 {
		t.Errorf("keywordTable has %d entries; want 188", len(keywordTable))
	}
}
