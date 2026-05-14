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

func TestEncodeFP_GeneralForm(t *testing.T) {
	tests := []struct {
		in      string
		wantErr string
	}{
		// Successful cases — we only check no-error here; exact bytes are
		// validated via the corpus round-trip in Task 19.
		{"0.5", ""},
		{"1.5", ""},
		{"100000", ""},
		{"1E5", ""},
		{"1.5E3", ""},
		{"1E38", ""},
		{".5", ""},
		{"1.", ""},
		// Errors:
		{"1E", `bad number syntax: "1E"`},
		{"1E+", `bad number syntax: "1E+"`},
		{".E5", `bad number syntax: ".E5"`},
		{"1E-300", "exponent out of range"},
		{"1E300", "exponent out of range"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			_, err := encodeFP(tt.in)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}
