package sambasic

import (
	"testing"
)

func TestEncodeFP_IntegerFastPath(t *testing.T) {
	tests := []struct {
		in   string
		want [5]byte
	}{
		{"0", [5]byte{0x00, 0x00, 0x00, 0x00, 0x00}},
		{"1", [5]byte{0x00, 0x00, 0x01, 0x00, 0x00}},
		{"2", [5]byte{0x00, 0x00, 0x02, 0x00, 0x00}},
		{"3", [5]byte{0x00, 0x00, 0x03, 0x00, 0x00}},
		{"255", [5]byte{0x00, 0x00, 0xFF, 0x00, 0x00}},
		{"256", [5]byte{0x00, 0x00, 0x00, 0x01, 0x00}},
		{"32767", [5]byte{0x00, 0x00, 0xFF, 0x7F, 0x00}},
		{"32768", [5]byte{0x00, 0x00, 0x00, 0x80, 0x00}},
		{"65535", [5]byte{0x00, 0x00, 0xFF, 0xFF, 0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := encodeFP(tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("encodeFP(%q) = % X, want % X", tt.in, got, tt.want)
			}
		})
	}
}

func TestEncodeFP_Hex(t *testing.T) {
	tests := []struct {
		in   string
		want [5]byte
	}{
		{"&80", [5]byte{0x00, 0x00, 0x80, 0x00, 0x00}},
		{"&FFFF", [5]byte{0x00, 0x00, 0xFF, 0xFF, 0x00}},
		{"&FF", [5]byte{0x00, 0x00, 0xFF, 0x00, 0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := encodeFP(tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("encodeFP(%q) = % X, want % X", tt.in, got, tt.want)
			}
		})
	}
}
