package sambasic

import (
	"fmt"
	"strconv"
	"strings"
)

// encodeFP converts an ASCII numeric literal to SAM's 5-byte invisible
// floating-point form. The fast-path for integer-valued literals in
// 0..65535 emits {0x00, 0x00, LSB, MSB, 0x00}. Hex literals (`&...`),
// decimal integers without a fractional part, and scientific notation
// that evaluates to an integer in 0..65535 all use the fast-path.
//
// The ROM's parser absorbs whitespace between digits in fraction parts
// (CONVFRAC2 / CONVFRALP loop at L17A8 uses RST 20H which skips
// 0x00-0x20), in hex digits (AMPERSAND at L18681 uses RST 20H), and
// in BIN digits. So `.0 20` parses as value 0.020. The lexer hands
// us the full literal *with* the embedded spaces; we strip them here
// before passing to strconv.
//
// Cites docs/sambasic-grammar.md §4.3.
func encodeFP(literal string) ([5]byte, error) {
	var out [5]byte
	if literal == "" {
		return out, fmt.Errorf("empty numeric literal")
	}
	// Strip embedded whitespace — the ROM-level value is the digit
	// sequence concatenated. See doc-comment above for which parse
	// paths absorb whitespace.
	if strings.IndexByte(literal, ' ') >= 0 {
		literal = strings.ReplaceAll(literal, " ", "")
		if literal == "" {
			return out, fmt.Errorf("empty numeric literal")
		}
	}
	// Hex: &[0-9A-Fa-f]+
	if literal[0] == '&' {
		if len(literal) == 1 {
			return out, fmt.Errorf("expected hex digits after &")
		}
		v, err := strconv.ParseUint(literal[1:], 16, 32)
		if err != nil {
			return out, fmt.Errorf("bad hex literal: %q", literal)
		}
		if v > 0xFFFFFF {
			return out, fmt.Errorf("hex literal too large")
		}
		if v <= 0xFFFF {
			out[2] = byte(v & 0xFF)
			out[3] = byte((v >> 8) & 0xFF)
			return out, nil
		}
		return encodeFloatToSAM(float64(v))
	}
	// Try integer fast-path.
	if !strings.ContainsAny(literal, ".eE") {
		v, err := strconv.ParseUint(literal, 10, 32)
		if err == nil && v <= 0xFFFF {
			out[2] = byte(v & 0xFF)
			out[3] = byte((v >> 8) & 0xFF)
			return out, nil
		}
		// Falls through to float encoding for large integers.
	}
	if err := validateDecimalScientific(literal); err != nil {
		return out, err
	}
	f, err := strconv.ParseFloat(literal, 64)
	if err != nil {
		return out, fmt.Errorf("bad number syntax: %q", literal)
	}
	return encodeFloatToSAM(f)
}

// validateDecimalScientific enforces the syntactic rules from
// docs/sambasic-grammar.md §4.2.
func validateDecimalScientific(s string) error {
	eIdx := -1
	for i := 0; i < len(s); i++ {
		if s[i] == 'E' || s[i] == 'e' {
			eIdx = i
			break
		}
	}
	if eIdx < 0 {
		hasDigit := false
		for i := 0; i < len(s); i++ {
			if s[i] >= '0' && s[i] <= '9' {
				hasDigit = true
				break
			}
		}
		if !hasDigit {
			return fmt.Errorf("bad number syntax: %q", s)
		}
		return nil
	}
	mantissa := s[:eIdx]
	exponent := s[eIdx+1:]
	hasDigit := false
	for i := 0; i < len(mantissa); i++ {
		if mantissa[i] >= '0' && mantissa[i] <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return fmt.Errorf("bad number syntax: %q", s)
	}
	if len(exponent) == 0 {
		return fmt.Errorf("bad number syntax: %q", s)
	}
	expDigits := exponent
	if exponent[0] == '+' || exponent[0] == '-' {
		expDigits = exponent[1:]
	}
	if len(expDigits) == 0 {
		return fmt.Errorf("bad number syntax: %q", s)
	}
	for i := 0; i < len(expDigits); i++ {
		if expDigits[i] < '0' || expDigits[i] > '9' {
			return fmt.Errorf("bad number syntax: %q", s)
		}
	}
	return nil
}

// encodeFloatToSAM converts a non-negative float64 to SAM's 5-byte
// general FP form. byte 0 = exponent + 0x80 (biased). Mantissa is
// normalised so 0.5 ≤ m < 1.0. The implicit leading-1 mantissa bit
// is replaced by the sign bit (TM p49 "the first bit is always 1,
// allowing it to be actually used as a SGN bit").
func encodeFloatToSAM(f float64) ([5]byte, error) {
	var out [5]byte
	if f == 0 {
		return out, nil
	}
	negative := false
	if f < 0 {
		negative = true
		f = -f
	}
	e := 0
	for f >= 1.0 {
		f /= 2.0
		e++
	}
	for f < 0.5 {
		f *= 2.0
		e--
	}
	biased := e + 0x80
	if biased < 1 || biased > 0xFF {
		return out, fmt.Errorf("exponent out of range")
	}
	out[0] = byte(biased)
	scaled := f * (1 << 32)
	mant := uint32(scaled + 0.5)
	mant &^= 0x80000000
	if negative {
		mant |= 0x80000000
	}
	out[1] = byte(mant >> 24)
	out[2] = byte(mant >> 16)
	out[3] = byte(mant >> 8)
	out[4] = byte(mant)
	return out, nil
}
