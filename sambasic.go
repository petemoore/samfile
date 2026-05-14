package samfile

import (
	"fmt"
	"io"
	"os"

	"github.com/petemoore/samfile/v3/sambasic"
)

type (
	// SAMBasic wraps the body bytes of a SAM BASIC (FT_SAM_BASIC,
	// type 16) file so they can be detokenised back into a plain
	// text listing. Data must be the file body without its 9-byte
	// FileHeader prefix — i.e. the Body field of a File returned
	// by DiskImage.File, or equivalent.
	//
	// The body is a sequence of tokenised lines terminated by a
	// 0xFF end-of-program sentinel. Each line has a 2-byte
	// big-endian line number, a 2-byte little-endian body length,
	// the tokenised body and a 0x0D line terminator. Keyword
	// tokens are single bytes in the range 0x85..0xF6 or the
	// two-byte sequence 0xFF, <idx>; numeric literals carry a
	// 5-byte "invisible" floating-point representation introduced
	// by 0x0E (which Output silently skips).
	//
	// Lossy controls Output's behaviour for control bytes and
	// formatting that the SAM ROM's LLIST routine filters/expands.
	// See Output() for the full spec.
	SAMBasic struct {
		Data []byte
		// Lossy = true makes Output match the SAM ROM's LLIST byte
		// stream as closely as possible (per the v3 ROM disassembly).
		// This includes filtering control-attribute sequences
		// (INK/PAPER/BRIGHT/INVERSE/OVER/FLASH/AT/TAB), emitting CR LF
		// line terminators, wrapping at column 80 with a 6-space
		// continuation indent, and applying the FLAGS-bit-0
		// "leading space before keyword" rule.
		//
		// Lossy = false (default) keeps a round-trip-faithful
		// rendering: control bytes are emitted as {N} escapes so they
		// survive parsing back through text-to-basic; no line wrapping;
		// LF-only terminators.
		Lossy bool
		// EPPC sets the line number that gets the `>` cursor marker in
		// lossy mode (mirrors the SAM ROM's EPPC sys-var at 0x5C49 —
		// "edit-position-program-counter"). The ROM's LIST/LLIST
		// command updates EPPC at 0x0681 to the FIRST line in the
		// requested range before printing; OUTLINE then prints `>`
		// only on the line whose number equals EPPC. If no listed
		// line has that number, no `>` appears in the output.
		//
		// For samfile callers: EPPC defaults to 1 because our
		// LLIST-capture harness always invokes `LLIST 1 TO 65278`.
		// Set EPPC explicitly if you want a different cursor line.
		EPPC uint16
	}
)

// NewSAMBasic wraps a SAM BASIC body for detokenisation. data is
// taken by reference, not copied.
func NewSAMBasic(data []byte) *SAMBasic {
	return &SAMBasic{
		Data: data,
	}
}

// Output writes basic.Data as a plain-text BASIC listing to stdout.
//
// In FAITHFUL mode (Lossy=false): each line is prefixed with a
// 5-space-padded decimal line number, keyword tokens are expanded via
// the v3 SAM BASIC keyword table, the invisible 5-byte numeric form
// after each 0x0E byte is skipped, and 0x0D becomes a newline.
// Control characters below 0x20 (other than 0x0D and 0x0E) are
// rendered as "{N}" so text-to-basic can re-create them. Keyword
// emission inserts a leading space when the previous output byte was
// non-space, so output like `125LOAD` (digit immediately followed by
// a keyword byte on disk) becomes `125 LOAD` and is parseable.
//
// In LOSSY mode (Lossy=true): byte-for-byte equivalent of the SAM
// ROM's LLIST command output to stream 3 (printer). Specifically:
//
//   - Line number column: 5-digit zero-suppressed, space-padded
//     (matches PRNUMB2 at ROM 0xF5B4).
//   - First listed line gets `>` after the line-number column (matches
//     OUTLINE's EPPC check at 0xF328); other lines get a space.
//   - Keyword bytes emit their name; trailing space if name ends in a
//     letter or `$` (per POMSG4 at 0xDD48); leading space ONLY when
//     FLAGS bit 0 = 0 (last printed char wasn't space; per POGEN at
//     0xDD29).
//   - 0x0E + next 5 bytes (invisible FP form) is silently skipped
//     (matches RDCN at 0x00A1).
//   - 0x0D terminates the line as CR LF (matches LPRENT at 0xDEC3).
//   - Control bytes 0x10..0x15 (INK/PAPER/FLASH/BRIGHT/INVERSE/OVER)
//     swallow themselves + 1 operand byte and emit nothing (CC1OP at
//     0xDEDB via dispatch table CCPTB at 0xDDDA).
//   - Control bytes 0x16..0x17 (AT/TAB) swallow themselves + 2
//     operand bytes and emit nothing (CC2OPS at 0xDEE0).
//   - Other control bytes (0x00..0x05, 0x07..0x0C, 0x0F, 0x18..0x1F)
//     emit literal `?` (PRQUERY via CCPTB).
//   - Column count is tracked; when it exceeds 79 (PRRHS default per
//     0xFC30) the output wraps with CR LF + 6-space indent (INDOPEN
//     at 0xF46A).
//   - Inside string literals (between `"`s) keyword expansion and
//     0x0E skipping are suppressed (INQUFG-equivalent); control-byte
//     filtering still applies (path goes through PROM1 regardless of
//     INQUFG).
//
// Returns an error if the input is empty, truncated, or contains an
// out-of-range keyword index.
func (basic *SAMBasic) Output() error {
	if len(basic.Data) == 0 {
		return fmt.Errorf("basic-to-text: empty input; expected SAM BASIC bytes on stdin")
	}
	eppc := basic.EPPC
	if eppc == 0 && !basic.Lossy {
		eppc = 1
	}
	// In lossy mode, eppc stays at 0 unless explicitly set. This
	// matches the ROM at 0x0681 (`LD (EPPC),HL ;NEW EPPC=FIRST
	// LISTED LINE`) which sets EPPC = the lower bound of the LIST
	// range, not the first actual line. Our LLIST-capture harness
	// uses `LLIST 0 TO 65278`, so EPPC=0: the `>` cursor appears on
	// line 0 if it exists, otherwise nowhere.
	s := &outputState{
		out:   os.Stdout,
		rhs:   79,
		eppc:  eppc,
		lossy: basic.Lossy,
	}
	n := uint32(len(basic.Data))
	index := uint32(0)
	for {
		if index >= n {
			return fmt.Errorf("basic-to-text: truncated input: missing 0xff end-of-program sentinel after offset %d", index)
		}
		if basic.Data[index] == 0xff {
			break
		}
		if index+3 >= n {
			return fmt.Errorf("basic-to-text: truncated input: incomplete line header at offset %d (need 4 bytes, have %d)", index, n-index)
		}
		lineNo := uint16(basic.Data[index])<<8 | uint16(basic.Data[index+1])
		lineLen := uint16(basic.Data[index+2]) | uint16(basic.Data[index+3])<<8
		index += 4

		s.startLine(lineNo)

		for c := uint16(0); c < lineLen; c++ {
			if index+uint32(c) >= n {
				return fmt.Errorf("basic-to-text: truncated input: line body for line %d extends past input (offset %d, length %d)", lineNo, index+uint32(c), n)
			}
			b := basic.Data[index+uint32(c)]
			consumed, err := s.handleByte(b, basic.Data, index+uint32(c), n, lineLen-c-1)
			if err != nil {
				return err
			}
			c += consumed
		}
		index += uint32(lineLen)
	}
	return nil
}

// outputState holds the LLIST-emulation state machine's variables.
// Field names mirror the SAM ROM's sys-var conventions where useful.
type outputState struct {
	out     io.Writer
	col     int    // current column (0-based)
	rhs     int    // PRRHS — wrap-after column (default 79)
	inQuote bool   // INQUFG — inside `"..."` string literal
	flagsSp bool   // FLAGS bit 0 — "last printed char was space"
	eppc    uint16 // emit `>` on the line whose number matches eppc
	lossy   bool
}

// putRaw writes a single byte directly, bypassing column tracking and
// wrap. Used for line endings and wrap-emit themselves.
func (s *outputState) putRaw(b byte) {
	_, _ = s.out.Write([]byte{b})
}

// emit writes one byte to the output. It advances the column, wraps
// at the RHS (lossy mode only), and tracks FLAGS bit 0. Do not call
// for line terminators / wrap markers — use putRaw + state resets.
//
// FLAGS bit 0 semantics differ between modes:
//
//   - Lossy: matches ROM PRASCII (0xDC11-DC17) — bit 0 = (last byte
//     printed was 0x20). Used for the POGEN "leading space if !flagsSp"
//     check, so disk-space bytes correctly suppress the extra leading
//     space LLIST would otherwise insert.
//
//   - Faithful: bit 0 stays false after any literal byte (including
//     spaces). This forces the leading-space-before-keyword rule to
//     ALWAYS fire when the previous emit was a literal — matching the
//     pre-refactor `spaceBefore` tracking, so that when the disk has
//     a literal space byte between content and a keyword, b2t emits
//     2 spaces. text-to-basic's leading-space-drop then consumes one,
//     leaving one — which is the byte the disk had.
func (s *outputState) emit(b byte) {
	if s.lossy && s.col > s.rhs {
		s.wrap()
	}
	s.putRaw(b)
	s.col++
	if s.lossy {
		s.flagsSp = (b == 0x20)
	} else {
		s.flagsSp = false
	}
}

// emitKeywordSpace writes a single space that originates from a
// keyword's leading/trailing-space convention (POGEN/POMSG4). Always
// sets flagsSp=true so the very next keyword doesn't get a redundant
// leading space. Distinguished from emit(0x20) which represents a
// literal space byte from the disk.
func (s *outputState) emitKeywordSpace() {
	if s.lossy && s.col > s.rhs {
		s.wrap()
	}
	s.putRaw(0x20)
	s.col++
	s.flagsSp = true
}

// emitString emits a series of bytes via emit() (column-tracked).
func (s *outputState) emitString(str string) {
	for i := 0; i < len(str); i++ {
		s.emit(str[i])
	}
}

// wrap emits CR LF + 6-space continuation indent. Per LPRENT at
// 0xDEC3 + INDOPEN at 0xF46A. Lossy mode only.
func (s *outputState) wrap() {
	s.putRaw(0x0D)
	s.putRaw(0x0A)
	for i := 0; i < 6; i++ {
		s.putRaw(0x20)
	}
	s.col = 6
	s.flagsSp = true
}

// endLine emits the line terminator. Lossy mode = CR LF (per LPRENT
// + AFTERCR default). Faithful mode = LF only.
func (s *outputState) endLine() {
	if s.lossy {
		s.putRaw(0x0D)
	}
	s.putRaw(0x0A)
	s.col = 0
	s.flagsSp = false
	s.inQuote = false
}

// startLine emits the line-number prefix and `>` cursor (or space).
// Mirrors PRNUMB2 + OUTLN3/OUTLN25 (ROM 0xF5B4 / 0xF35F / 0xF356).
func (s *outputState) startLine(lineNo uint16) {
	// 5-digit zero-suppressed, space-padded line number.
	fmt.Fprintf(s.out, "%5d", lineNo)
	s.col = 5
	s.flagsSp = false
	if s.lossy && lineNo == s.eppc {
		// EPPC marker — per OUTLN25 at 0xF356 the `>` itself sets
		// "no leading space" via SET 0,(HL) at 0xF35A.
		s.putRaw('>')
		s.col++
		s.flagsSp = true
		return
	}
	// Trailing space after the line number — sets "no leading space"
	// for the very first body byte / keyword.
	s.emitKeywordSpace()
}

// handleByte processes one body byte. data + offset + remaining are
// provided so the handler can peek/consume successor bytes (used for
// the 0xFF 2-byte keyword escape, 0x0E + 5-byte FP form, and the
// CC1OP / CC2OPS attribute-operand swallows). Returns the number of
// EXTRA bytes consumed past `b` (caller increments c by 1 + result).
func (s *outputState) handleByte(b byte, data []byte, offset, n uint32, remaining uint16) (uint16, error) {
	// In-string handling: most special bytes are emitted verbatim.
	// Control bytes are still filtered per CCPTB (LLIST path goes
	// through PROM1 regardless of INQUFG state).
	if s.inQuote {
		switch {
		case b == 0x0d:
			// Implicit string close at line end.
			s.endLine()
			return 0, nil
		case b == 0x22:
			// Doubled-quote `""` = embedded `"` literal; lone `"` = close.
			if remaining > 0 && data[offset+1] == 0x22 {
				s.emit('"')
				s.emit('"')
				return 1, nil
			}
			s.emit('"')
			s.inQuote = false
			return 0, nil
		case b < 0x20:
			return s.handleControl(b, data, offset, remaining)
		case s.lossy && b >= 0xA9 && b <= 0xFE:
			// POFUDG at ROM 0xDDB4 — inside a string literal (INQUFG=1),
			// bytes 0xA9..0xFE get SUB 0xA9 before reaching the printer
			// via OPCHAR. The SAM ROM path is PRGR80 → POUDGH → POUDG
			// (DEVICE=2, LPRTV=0 no-op) → PUDGS → POFUDG → SUB 0xA9 →
			// POUDG1 → PRINTMN1 → IOPENT stores OPCHAR = post-SUB value
			// → ENDOUTP → ENDOP2 → CHBOP sends OPCHAR to channel B →
			// SENDA → port out. Empirically: byte 0xE9 (BOOT) inside a
			// string emits as 0x40 (`@`) on the printer; byte 0xFE
			// would emit as 0x55 (`U`). Verified against SAM81 line
			// 130 on B-DOS V1.7D.
			//
			// Bytes 0x85..0xA8 inside strings pass through as-is
			// (block-graphics path at DD99 and UDG path at DDB8 leave
			// OPCHAR holding the original byte) — handled by the
			// default case below.
			s.emit(b - 0xA9)
			return 0, nil
		default:
			s.emit(b)
			return 0, nil
		}
	}

	// Outside a string.
	switch {
	case b == 0xff:
		if remaining == 0 {
			return 0, fmt.Errorf("basic-to-text: truncated input: 0xff keyword escape at end of input (offset %d)", offset)
		}
		kwByte := data[offset+1]
		name, ok := sambasic.KeywordName(kwByte, true)
		if !ok {
			return 0, fmt.Errorf("basic-to-text: invalid keyword byte 0x%02x after 0xff escape at offset %d", kwByte, offset+1)
		}
		s.emitKeywordTwoByte(kwByte, name)
		return 1, nil
	case b == 0x0e:
		// 5-byte FP form — silently skipped (RDCN at 0x00A1).
		// In faithful mode this is also dropped (the 5 bytes are
		// reconstructible from the visible digit run that preceded).
		if remaining < 5 {
			return 0, fmt.Errorf("basic-to-text: truncated input: 0x0e FP form at end of line (offset %d, remaining %d)", offset, remaining)
		}
		return 5, nil
	case b == 0x0d:
		s.endLine()
		return 0, nil
	case b == 0x22:
		s.emit('"')
		s.inQuote = true
		return 0, nil
	case b < 0x20:
		return s.handleControl(b, data, offset, remaining)
	case b >= 0x85 && b <= 0xf6:
		name, ok := sambasic.KeywordName(b, false)
		if !ok {
			return 0, fmt.Errorf("basic-to-text: keyword index %d out of range", b-0x3b)
		}
		s.emitKeyword(name)
		return 0, nil
	default:
		s.emit(b)
		return 0, nil
	}
}

// emitKeyword writes a 1-byte keyword (range 0x85..0xF6) using the
// POBTL dispatch (PSTFF2 path falls through to POBTL at 0xDD28 for
// 1-byte commands): leading space if FLAGS bit 0 = 0, trailing
// space if the name ends in a letter or `$` (POMSG4 at 0xDD48).
func (s *outputState) emitKeyword(name string) {
	s.emitKeywordWithRule(name, kwLeading|kwTrailing)
}

// emitKeywordTwoByte writes a 2-byte keyword (FF nn) using the
// per-token-class dispatch from PSTFF2 at 0xDCEE:
//
//   - nn in 0x3B..0x52 (POIMFN range, immediate fns including
//     SCREEN$, CHR$, PI, INKEY$, etc.):
//     * nn == 0x42 (FN) or nn == 0x43 (BIN) → POTRNL (trailing only)
//     * otherwise → POMSP2 (no spaces at all)
//   - nn in 0x53..0x79 (POFPCFN range, fp-calc functions like SIN,
//     COS, TAN, INT, LN, EXP) → POTRNL
//   - nn in 0x7A..0x80 (POBTL range, MOD/DIV/AND-area binary ops) → POBTL
//   - nn >= 0x81 (operators <>, <=, >=, etc.) → POMSP2
//
// PITOK = 0x3B, FNTOK = 0x42, BINTOK = 0x43, SINTOK = 0x53,
// MODTOK = 0x7A, ANDTOK = 0x80 — all from the v3 ROM source.
func (s *outputState) emitKeywordTwoByte(nn byte, name string) {
	var rule uint8
	switch {
	case nn < 0x3B:
		// Below PITOK — shouldn't happen for valid disks, but be defensive.
		rule = kwLeading | kwTrailing
	case nn < 0x53:
		// POIMFN range. FN / BIN get POTRNL; others get POMSP2.
		if nn == 0x42 || nn == 0x43 {
			rule = kwTrailing
		} else {
			rule = 0
		}
	case nn < 0x7A:
		// POFPCFN range → POTRNL.
		rule = kwTrailing
	case nn < 0x81:
		// POBTL range (MOD..AND).
		rule = kwLeading | kwTrailing
	default:
		// POMSP2 (operators).
		rule = 0
	}
	s.emitKeywordWithRule(name, rule)
}

const (
	kwLeading  uint8 = 1 << 0
	kwTrailing uint8 = 1 << 1
)

func (s *outputState) emitKeywordWithRule(name string, rule uint8) {
	if name == "" {
		return
	}
	// In faithful mode, leading space fires for ALL keywords whenever
	// previous emit wasn't a keyword-trailing space, so that disk-
	// literal spaces between content and a keyword round-trip through
	// text-to-basic correctly. POMSP2 (no-leading) and POTRNL (no-
	// leading) are only enforced in lossy mode.
	leading := rule&kwLeading != 0
	trailing := rule&kwTrailing != 0
	if !s.lossy {
		leading = true
	}
	if leading && !s.flagsSp {
		s.emitKeywordSpace()
	}
	s.emitString(name)
	if trailing {
		last := name[len(name)-1]
		if isAlpha(last) || last == '$' {
			s.emitKeywordSpace()
		}
	} else if !s.lossy {
		// Faithful mode: emit trailing space for ALL keywords whose
		// last char is letter or `$`, even POMSP2 ones, so the next
		// keyword has a clean separator.
		last := name[len(name)-1]
		if isAlpha(last) || last == '$' {
			s.emitKeywordSpace()
		}
	}
}

func isAlpha(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

// handleControl processes a control byte (b < 0x20) per the SAM ROM's
// PRCRLCDS dispatch table (0xDDC4 / CCPTB at 0xDDDA).
//
// In LOSSY mode:
//
//	0x00..0x05  -> `?`  (PRQUERY)
//	0x06        -> pad with spaces to next 16-column tab stop (PRCOMMA)
//	0x07        -> `?`  (PRQUERY in CCPTB table)
//	0x08..0x0C  -> `?`  (cursor / delete control codes emit `?` on printer)
//	0x0D        -> CR LF  (caller normally hits this via the `b == 0x0d` switch
//	               branch; this case is here defensively)
//	0x0E        -> skip 5 FP bytes (also handled in the main switch)
//	0x0F        -> `?`  (CCPTB)
//	0x10..0x15  -> consume self + 1 operand, emit nothing (CC1OP)
//	0x16..0x17  -> consume self + 2 operands, emit nothing (CC2OPS)
//	0x18..0x1F  -> `?`  (PRQUERY; out-of-table range)
//
// In FAITHFUL mode, control bytes (other than 0x0D and 0x0E which are
// handled in the main switch) are emitted as `{N}` escapes for
// round-trippability.
func (s *outputState) handleControl(b byte, data []byte, offset uint32, remaining uint16) (uint16, error) {
	if !s.lossy {
		s.flagsSp = false
		fmt.Fprintf(s.out, "{%d}", int(b))
		// {N} emits multiple chars; conservatively assume non-space
		// for column tracking and bit 0. Column is tracked via emit()
		// elsewhere; for {N} we approximate.
		s.col += digitsOf(int(b)) + 2
		return 0, nil
	}
	switch b {
	case 0x06:
		// PRCOMMA (ROM 0xDDEC) uses screen WINDLHS=0/WINDRHS=31, NOT
		// the printer PRRHS=79. Two sub-rules:
		//
		//   col <= 31: pad to next 16-col tab stop, capped at 32
		//              (= WINDRHS+1). Math is `(col & 0xF0) + 16`.
		//   col >  31: PC25 path at 0xDE17 unconditionally emits
		//              exactly 16 spaces via OPSPLP (DJNZ loop). The
		//              individual space-emits go through PROM1's
		//              NLENTRY, so if col crosses 79 mid-loop the
		//              wrap fires (LPRENT + INDOPEN 6-space indent)
		//              and the remaining DJNZ iterations continue
		//              emitting spaces AFTER the indent. That's why
		//              the SNOOKER line 20 "continuation indent"
		//              looks like 14 spaces — it's 6 from INDOPEN
		//              plus the rest of the 16-space PC25 loop.
		//
		// Trace via ROM L21577 (PRCOMMA) + L21300+ (PRGR80 path) +
		// the agent's analysis in this commit's working notes.
		const windRHS = 31
		if s.col > windRHS {
			// PC25: 16 unconditional spaces, individually emit so
			// wrap can fire mid-loop.
			for i := 0; i < 16; i++ {
				s.emit(0x20)
			}
		} else {
			target := (s.col & 0xF0) + 16
			if target > windRHS+1 {
				target = windRHS + 1
			}
			for s.col < target {
				s.emit(0x20)
			}
		}
		return 0, nil
	case 0x10, 0x11, 0x12, 0x13, 0x14, 0x15:
		// CC1OP — consume control + 1 operand byte, emit nothing.
		if remaining > 0 {
			return 1, nil
		}
		return 0, nil
	case 0x16, 0x17:
		// CC2OPS — consume control + 2 operand bytes, emit nothing.
		if remaining >= 2 {
			return 2, nil
		}
		return remaining, nil
	default:
		// PRQUERY paths (0x00..0x05, 0x07..0x0C, 0x0F, 0x18..0x1F).
		s.emit('?')
		return 0, nil
	}
}

func digitsOf(n int) int {
	if n == 0 {
		return 1
	}
	d := 0
	for n > 0 {
		d++
		n /= 10
	}
	return d
}
