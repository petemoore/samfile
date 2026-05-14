package sambasic

import (
	"fmt"
)

// Schema describes the structure of bytes emitted by ParseText, with
// each segment carrying provenance and rules for acceptable variation.
// Corpus comparison uses Conform(want, schema) instead of byte equality
// so that intentional variations (LOAD-rebuilt PROC placeholders,
// equivalent FP encodings, etc.) don't count as divergences while still
// catching real mismatches.
type Schema struct {
	Segments []SchemaSegment
}

// SchemaSegment is a contiguous run of bytes in the produced output
// classified by its provenance. Offsets are absolute into the full
// ProgBytes() concatenation.
type SchemaSegment struct {
	Offset int         // byte offset of this segment in the produced output
	Length int         // segment length in bytes
	Kind   SegmentKind
	Bytes  []byte // the exact bytes produced (for diagnostics; conformance rules vary)
	Detail any    // kind-specific data; nil if not applicable
}

// SegmentKind enumerates the schema's notion of byte provenance.
type SegmentKind int

const (
	SegLineHeader            SegmentKind = iota // 4 bytes: MSB LSB LenLo LenHi
	SegLineTerminator                           // 1 byte: 0x0D
	SegKeyword                                  // 1 or 2 bytes
	SegNumberVisible                            // visible ASCII digit run
	SegNumberFP                                 // 0x0E + 5 invisible FP bytes
	SegString                                   // string body verbatim (delimiters + "" doubling kept)
	SegREMBody                                  // REM body verbatim
	SegLiteral                                  // any other single byte (operators, separators, identifier chars)
	SegControlEscape                            // single raw byte from {N} sugar
	SegProcCallPlaceholder                      // identifier bytes + 0E FD FD FD 00 00 (last 3 LOAD-rebuilt)
	SegFnCallPlaceholder                        // identifier bytes + 0E FE FE FE 00 00 (last 3 LOAD-rebuilt)
	SegProgramEnd                               // 1 byte: 0xFF
)

// String returns a short human-readable label for the segment kind.
func (k SegmentKind) String() string {
	switch k {
	case SegLineHeader:
		return "LineHeader"
	case SegLineTerminator:
		return "LineTerminator"
	case SegKeyword:
		return "Keyword"
	case SegNumberVisible:
		return "NumberVisible"
	case SegNumberFP:
		return "NumberFP"
	case SegString:
		return "String"
	case SegREMBody:
		return "REMBody"
	case SegLiteral:
		return "Literal"
	case SegControlEscape:
		return "ControlEscape"
	case SegProcCallPlaceholder:
		return "ProcCallPlaceholder"
	case SegFnCallPlaceholder:
		return "FnCallPlaceholder"
	case SegProgramEnd:
		return "ProgramEnd"
	}
	return fmt.Sprintf("SegmentKind(%d)", int(k))
}

// NumberDetail carries the parsed numeric value for a SegNumberFP
// segment so that Conform can compare across equivalent encodings.
// IsInt is true when the literal fits the integer fast-path
// (uint16, encoded as {0x00, 0x00, LSB, MSB, 0x00}). Otherwise the
// general FP form was used and Float carries the parsed value.
type NumberDetail struct {
	IsInt bool
	Int   uint16
	Float float64
}

// PlaceholderDetail carries the identifier text for a PROC/FN call
// placeholder. The identifier text is also already present in the
// segment's Bytes prefix; this struct just makes intent explicit.
type PlaceholderDetail struct {
	Identifier string
	IDLen      int // length of the identifier prefix in the segment's Bytes
}

// Mismatch records a single point of divergence between a candidate
// byte slice and a Schema.
type Mismatch struct {
	Offset      int
	SegmentKind SegmentKind
	Description string
	WantBytes   []byte // bytes at the diff site in `want`
	GotBytes    []byte // schema's expected bytes at the same site
}

// Conform checks whether `want` (corpus bytes) is a valid encoding of
// the program described by `schema`. Returns nil if it conforms,
// otherwise a slice of mismatches with offset, segment kind, and
// human-readable description.
//
// Per-segment rules:
//   - SegLineHeader, SegLineTerminator, SegKeyword, SegLiteral,
//     SegControlEscape, SegProgramEnd, SegString, SegREMBody,
//     SegNumberVisible: exact byte equality.
//   - SegNumberFP: byte 0 of `want` must be 0x0E; bytes 1..5 must
//     decode to the same numeric value as the schema's Detail.
//   - SegProcCallPlaceholder / SegFnCallPlaceholder: identifier
//     prefix bytes must match exactly; then 0x0E; then FD FD (PROC)
//     or FE FE (FN); then a "page" byte of 0xFF (unresolved) or any
//     value with the high bit set (resolved); then any 2 bytes
//     (LOAD-rebuilt address).
func Conform(want []byte, schema *Schema) []Mismatch {
	var mismatches []Mismatch
	expectedLen := 0
	for _, seg := range schema.Segments {
		end := seg.Offset + seg.Length
		if end > expectedLen {
			expectedLen = end
		}
	}
	for _, seg := range schema.Segments {
		segEnd := seg.Offset + seg.Length
		if segEnd > len(want) {
			// Segment runs past the end of want.
			haveLen := 0
			if seg.Offset < len(want) {
				haveLen = len(want) - seg.Offset
			}
			var wb []byte
			if seg.Offset < len(want) {
				wb = append([]byte(nil), want[seg.Offset:]...)
			}
			mismatches = append(mismatches, Mismatch{
				Offset:      seg.Offset,
				SegmentKind: seg.Kind,
				Description: fmt.Sprintf("%s segment of %d bytes runs past want (have %d)", seg.Kind, seg.Length, haveLen),
				WantBytes:   wb,
				GotBytes:    append([]byte(nil), seg.Bytes...),
			})
			continue
		}
		ws := want[seg.Offset:segEnd]
		switch seg.Kind {
		case SegNumberFP:
			if m, ok := checkNumberFP(seg, ws); !ok {
				mismatches = append(mismatches, m)
			}
		case SegProcCallPlaceholder, SegFnCallPlaceholder:
			if m, ok := checkPlaceholder(seg, ws); !ok {
				mismatches = append(mismatches, m)
			}
		default:
			if !bytesEqualSlice(ws, seg.Bytes) {
				mismatches = append(mismatches, Mismatch{
					Offset:      seg.Offset,
					SegmentKind: seg.Kind,
					Description: fmt.Sprintf("%s bytes differ", seg.Kind),
					WantBytes:   append([]byte(nil), ws...),
					GotBytes:    append([]byte(nil), seg.Bytes...),
				})
			}
		}
	}
	if len(want) > expectedLen {
		mismatches = append(mismatches, Mismatch{
			Offset:      expectedLen,
			SegmentKind: SegProgramEnd,
			Description: fmt.Sprintf("want has %d trailing bytes past schema end", len(want)-expectedLen),
			WantBytes:   append([]byte(nil), want[expectedLen:]...),
			GotBytes:    nil,
		})
	}
	return mismatches
}

// checkNumberFP validates a 6-byte 0x0E + FP segment in want against
// the schema's recorded numeric value.
func checkNumberFP(seg SchemaSegment, ws []byte) (Mismatch, bool) {
	if len(ws) < 6 || ws[0] != 0x0E {
		return Mismatch{
			Offset:      seg.Offset,
			SegmentKind: seg.Kind,
			Description: "NumberFP missing 0x0E marker",
			WantBytes:   append([]byte(nil), ws...),
			GotBytes:    append([]byte(nil), seg.Bytes...),
		}, false
	}
	detail, ok := seg.Detail.(NumberDetail)
	if !ok {
		// No detail recorded: fall back to exact byte equality.
		if !bytesEqualSlice(ws, seg.Bytes) {
			return Mismatch{
				Offset:      seg.Offset,
				SegmentKind: seg.Kind,
				Description: "NumberFP bytes differ (no detail)",
				WantBytes:   append([]byte(nil), ws...),
				GotBytes:    append([]byte(nil), seg.Bytes...),
			}, false
		}
		return Mismatch{}, true
	}
	if detail.IsInt {
		// Accept any encoding that decodes to the same uint16 value.
		// Integer fast-path: {0x00, 0x00, LSB, MSB, 0x00}.
		fp := ws[1:6]
		if fp[0] == 0x00 && fp[1] == 0x00 && fp[4] == 0x00 {
			v := uint16(fp[2]) | uint16(fp[3])<<8
			if v == detail.Int {
				return Mismatch{}, true
			}
		}
		// Otherwise try a general-form decode and compare values.
		f, ok := decodeFP5(fp)
		if ok && f == float64(detail.Int) {
			return Mismatch{}, true
		}
		return Mismatch{
			Offset:      seg.Offset,
			SegmentKind: seg.Kind,
			Description: fmt.Sprintf("NumberFP encodes a different value (want=%d)", detail.Int),
			WantBytes:   append([]byte(nil), ws...),
			GotBytes:    append([]byte(nil), seg.Bytes...),
		}, false
	}
	// Float case.
	fp := ws[1:6]
	f, ok := decodeFP5(fp)
	if !ok {
		return Mismatch{
			Offset:      seg.Offset,
			SegmentKind: seg.Kind,
			Description: "NumberFP unable to decode FP bytes",
			WantBytes:   append([]byte(nil), ws...),
			GotBytes:    append([]byte(nil), seg.Bytes...),
		}, false
	}
	if floatNearlyEqual(f, detail.Float) {
		return Mismatch{}, true
	}
	return Mismatch{
		Offset:      seg.Offset,
		SegmentKind: seg.Kind,
		Description: fmt.Sprintf("NumberFP encodes a different value (got=%g want=%g)", f, detail.Float),
		WantBytes:   append([]byte(nil), ws...),
		GotBytes:    append([]byte(nil), seg.Bytes...),
	}, false
}

// checkPlaceholder validates a PROC/FN placeholder against the
// grammar's rules. The first idLen bytes are the identifier and must
// match exactly. The next 6 bytes are 0x0E, two type bytes (FD FD or
// FE FE), a "page" byte, and two LOAD-rebuilt address bytes.
func checkPlaceholder(seg SchemaSegment, ws []byte) (Mismatch, bool) {
	detail, ok := seg.Detail.(PlaceholderDetail)
	if !ok {
		// No detail: fall back to exact match.
		if !bytesEqualSlice(ws, seg.Bytes) {
			return Mismatch{
				Offset:      seg.Offset,
				SegmentKind: seg.Kind,
				Description: fmt.Sprintf("%s bytes differ (no detail)", seg.Kind),
				WantBytes:   append([]byte(nil), ws...),
				GotBytes:    append([]byte(nil), seg.Bytes...),
			}, false
		}
		return Mismatch{}, true
	}
	if len(ws) != detail.IDLen+6 {
		return Mismatch{
			Offset:      seg.Offset,
			SegmentKind: seg.Kind,
			Description: fmt.Sprintf("%s segment length %d != expected %d", seg.Kind, len(ws), detail.IDLen+6),
			WantBytes:   append([]byte(nil), ws...),
			GotBytes:    append([]byte(nil), seg.Bytes...),
		}, false
	}
	// Identifier prefix must match exactly.
	if !bytesEqualSlice(ws[:detail.IDLen], seg.Bytes[:detail.IDLen]) {
		return Mismatch{
			Offset:      seg.Offset,
			SegmentKind: seg.Kind,
			Description: fmt.Sprintf("%s identifier prefix differs", seg.Kind),
			WantBytes:   append([]byte(nil), ws[:detail.IDLen]...),
			GotBytes:    append([]byte(nil), seg.Bytes[:detail.IDLen]...),
		}, false
	}
	tail := ws[detail.IDLen:]
	if tail[0] != 0x0E {
		return Mismatch{
			Offset:      seg.Offset + detail.IDLen,
			SegmentKind: seg.Kind,
			Description: fmt.Sprintf("%s missing 0x0E marker", seg.Kind),
			WantBytes:   append([]byte(nil), tail...),
			GotBytes:    append([]byte(nil), seg.Bytes[detail.IDLen:]...),
		}, false
	}
	var wantType byte
	if seg.Kind == SegProcCallPlaceholder {
		wantType = 0xFD
	} else {
		wantType = 0xFE
	}
	if tail[1] != wantType || tail[2] != wantType {
		return Mismatch{
			Offset:      seg.Offset + detail.IDLen + 1,
			SegmentKind: seg.Kind,
			Description: fmt.Sprintf("%s type bytes are %02X %02X, want %02X %02X", seg.Kind, tail[1], tail[2], wantType, wantType),
			WantBytes:   append([]byte(nil), tail...),
			GotBytes:    append([]byte(nil), seg.Bytes[detail.IDLen:]...),
		}, false
	}
	// Page byte: 0xFF (unresolved) or any byte with the high bit set
	// (resolved-to-page, per grammar spec §6.5).
	page := tail[3]
	if page != 0xFF && (page&0x80) == 0 {
		return Mismatch{
			Offset:      seg.Offset + detail.IDLen + 3,
			SegmentKind: seg.Kind,
			Description: fmt.Sprintf("%s page byte 0x%02X has high bit clear and is not 0xFF", seg.Kind, page),
			WantBytes:   append([]byte(nil), tail...),
			GotBytes:    append([]byte(nil), seg.Bytes[detail.IDLen:]...),
		}, false
	}
	// tail[4] and tail[5] are LOAD-rebuilt address bytes: accept any.
	return Mismatch{}, true
}

// decodeFP5 decodes a 5-byte SAM BASIC FP form into a float64 value.
// Returns (value, true) on success. The decoder accepts both the
// integer fast-path encoding {0x00, 0x00, LSB, MSB, 0x00} and the
// general normalised form.
func decodeFP5(fp []byte) (float64, bool) {
	if len(fp) != 5 {
		return 0, false
	}
	// Integer fast-path: byte 0 == 0x00.
	if fp[0] == 0x00 {
		// Treat as signed-magnitude uint16-ish; SAM stores positive
		// integers as {0x00, 0x00, LSB, MSB, 0x00}. We accept that
		// shape literally; anything else with exp=0 is unusual.
		if fp[1] == 0x00 && fp[4] == 0x00 {
			v := uint16(fp[2]) | uint16(fp[3])<<8
			return float64(v), true
		}
		// Fall through; treat as zero or unusual.
		if fp[1] == 0 && fp[2] == 0 && fp[3] == 0 && fp[4] == 0 {
			return 0, true
		}
		return 0, false
	}
	// General form: exponent in fp[0] (biased by 0x80); 32-bit
	// mantissa across fp[1..4], with the top bit being the sign and
	// the implicit leading-1 mantissa bit restored.
	biased := int(fp[0])
	e := biased - 0x80
	mant := uint32(fp[1])<<24 | uint32(fp[2])<<16 | uint32(fp[3])<<8 | uint32(fp[4])
	negative := mant&0x80000000 != 0
	mant |= 0x80000000 // restore implicit leading-1
	m := float64(mant) / float64(1<<32)
	v := m
	if e > 0 {
		for i := 0; i < e; i++ {
			v *= 2.0
		}
	} else {
		for i := 0; i < -e; i++ {
			v /= 2.0
		}
	}
	if negative {
		v = -v
	}
	return v, true
}

// floatNearlyEqual compares two floats with a small relative
// tolerance, treating exact 0 equality specially.
func floatNearlyEqual(a, b float64) bool {
	if a == b {
		return true
	}
	if a == 0 || b == 0 {
		// One side is zero; require absolute closeness.
		diff := a - b
		if diff < 0 {
			diff = -diff
		}
		return diff < 1e-30
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	ref := a
	if ref < 0 {
		ref = -ref
	}
	rb := b
	if rb < 0 {
		rb = -rb
	}
	if rb > ref {
		ref = rb
	}
	return diff/ref < 1e-9
}

func bytesEqualSlice(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
