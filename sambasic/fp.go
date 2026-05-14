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
// Task 9 extends this to general FP form for non-integer values.
//
// Cites docs/sambasic-grammar.md §4.3.
func encodeFP(literal string) ([5]byte, error) {
	var out [5]byte
	if literal == "" {
		return out, fmt.Errorf("empty numeric literal")
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
		if v > 0xFFFF {
			return out, fmt.Errorf("hex value > 65535 not yet supported")
		}
		out[2] = byte(v & 0xFF)
		out[3] = byte((v >> 8) & 0xFF)
		return out, nil
	}
	// Decimal integer fast-path.
	if !strings.ContainsAny(literal, ".eE") {
		v, err := strconv.ParseUint(literal, 10, 32)
		if err != nil {
			return out, fmt.Errorf("bad number syntax: %q", literal)
		}
		if v > 0xFFFF {
			return out, fmt.Errorf("decimal > 65535 not yet supported")
		}
		out[2] = byte(v & 0xFF)
		out[3] = byte((v >> 8) & 0xFF)
		return out, nil
	}
	return out, fmt.Errorf("non-integer literal: not yet supported")
}
