package sambasic

import (
	"testing"
)

// TestConform_ExactMatch verifies that schema-produced bytes conform
// to themselves.
func TestConform_ExactMatch(t *testing.T) {
	_, sch, err := ParseTextSchema(`10 PRINT "hello"` + "\n")
	if err != nil {
		t.Fatalf("ParseTextSchema: %v", err)
	}
	want := assembleBytesFromSchema(sch)
	if mm := Conform(want, sch); len(mm) != 0 {
		t.Errorf("expected exact match, got %d mismatches:\n%s", len(mm), formatMismatches(mm))
	}
}

// TestConform_ProcPlaceholder_DifferentTrailingBytes verifies that a
// PROC call placeholder with different LOAD-rebuilt trailing bytes
// still conforms.
func TestConform_ProcPlaceholder_DifferentTrailingBytes(t *testing.T) {
	_, sch, err := ParseTextSchema(`10 myproc` + "\n")
	if err != nil {
		t.Fatalf("ParseTextSchema: %v", err)
	}
	want := assembleBytesFromSchema(sch)
	// Find the PROC placeholder segment and rewrite its trailing 3
	// bytes (page + 2 LOAD-rebuilt address bytes) to a valid resolved
	// form (page = 0x80|page).
	for _, seg := range sch.Segments {
		if seg.Kind != SegProcCallPlaceholder {
			continue
		}
		// Layout: idBytes + 0E FD FD <page> <lo> <hi>
		want[seg.Offset+seg.Length-3] = 0x81 // page byte with high bit set
		want[seg.Offset+seg.Length-2] = 0x42 // arbitrary lo
		want[seg.Offset+seg.Length-1] = 0x99 // arbitrary hi
		break
	}
	if mm := Conform(want, sch); len(mm) != 0 {
		t.Errorf("expected PROC placeholder to conform with rewritten trailing bytes, got %d mismatches:\n%s", len(mm), formatMismatches(mm))
	}
}

// TestConform_ProcPlaceholder_BadTypeByte verifies a non-PROC type
// byte does NOT conform.
func TestConform_ProcPlaceholder_BadTypeByte(t *testing.T) {
	_, sch, err := ParseTextSchema(`10 myproc` + "\n")
	if err != nil {
		t.Fatalf("ParseTextSchema: %v", err)
	}
	want := assembleBytesFromSchema(sch)
	for _, seg := range sch.Segments {
		if seg.Kind != SegProcCallPlaceholder {
			continue
		}
		// Corrupt the first type byte (FD) to something else.
		want[seg.Offset+seg.Length-5] = 0xFE
		break
	}
	if mm := Conform(want, sch); len(mm) == 0 {
		t.Errorf("expected PROC placeholder with FE type byte to NOT conform")
	}
}

// TestConform_NumberFP_EquivalentEncodings verifies that the integer
// fast-path encoding {00 00 LSB MSB 00} and the general normalised
// form for the same uint16 value both conform to a schema produced
// for that integer.
func TestConform_NumberFP_EquivalentEncodings(t *testing.T) {
	_, sch, err := ParseTextSchema(`10 PRINT 5` + "\n")
	if err != nil {
		t.Fatalf("ParseTextSchema: %v", err)
	}
	want := assembleBytesFromSchema(sch)
	if mm := Conform(want, sch); len(mm) != 0 {
		t.Fatalf("baseline: expected exact match, got %d mismatches:\n%s", len(mm), formatMismatches(mm))
	}
	// Find the NumberFP segment and swap to a general-form encoding
	// of the value 5: 5 = 1.25 * 2^2, so biased exp = 0x83, mantissa
	// (with implicit leading 1 cleared) = 0x40000000 then shifted
	// down by 1 bit gives mant bytes = {0x40, 0x00, 0x00, 0x00}.
	// (See encodeFloatToSAM in fp.go.)
	encoded, err := encodeFloatToSAM(5.0)
	if err != nil {
		t.Fatalf("encodeFloatToSAM: %v", err)
	}
	for _, seg := range sch.Segments {
		if seg.Kind != SegNumberFP {
			continue
		}
		// Layout: 0E + 5 FP bytes
		copy(want[seg.Offset+1:seg.Offset+6], encoded[:])
		break
	}
	if mm := Conform(want, sch); len(mm) != 0 {
		t.Errorf("expected general-form encoding of 5 to conform, got %d mismatches:\n%s", len(mm), formatMismatches(mm))
	}
}

// TestConform_StringBytes_MustMatchExactly verifies that string
// content must match exactly.
func TestConform_StringBytes_MustMatchExactly(t *testing.T) {
	_, sch, err := ParseTextSchema(`10 PRINT "hello"` + "\n")
	if err != nil {
		t.Fatalf("ParseTextSchema: %v", err)
	}
	want := assembleBytesFromSchema(sch)
	// Find the SegString segment and tweak a byte.
	for _, seg := range sch.Segments {
		if seg.Kind != SegString {
			continue
		}
		// "hello" -> swap the 'h' for 'H'.
		want[seg.Offset+1] = 'H'
		break
	}
	if mm := Conform(want, sch); len(mm) == 0 {
		t.Errorf("expected a string byte mismatch")
	}
}

// TestConform_LengthMismatch verifies that trailing bytes past the
// schema's reach surface as a mismatch.
func TestConform_LengthMismatch(t *testing.T) {
	_, sch, err := ParseTextSchema(`10 PRINT 1` + "\n")
	if err != nil {
		t.Fatalf("ParseTextSchema: %v", err)
	}
	want := append(assembleBytesFromSchema(sch), 0xAA, 0xBB)
	mm := Conform(want, sch)
	if len(mm) == 0 {
		t.Errorf("expected trailing-bytes mismatch")
	}
}

// assembleBytesFromSchema concatenates every segment's Bytes in
// offset order to produce a candidate byte slice. Helper for the
// schema tests above.
func assembleBytesFromSchema(sch *Schema) []byte {
	out := make([]byte, 0)
	end := 0
	for _, seg := range sch.Segments {
		segEnd := seg.Offset + seg.Length
		if segEnd > end {
			end = segEnd
		}
	}
	out = make([]byte, end)
	for _, seg := range sch.Segments {
		copy(out[seg.Offset:seg.Offset+seg.Length], seg.Bytes)
	}
	return out
}

func formatMismatches(mm []Mismatch) string {
	var s string
	for _, m := range mm {
		s += "  off=" + itoa(m.Offset) + " kind=" + m.SegmentKind.String() + " desc=" + m.Description + "\n"
	}
	return s
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
