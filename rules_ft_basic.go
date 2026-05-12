package samfile

import (
	"fmt"

	"github.com/petemoore/samfile/v3/sambasic"
)

// §7 FT_SAM_BASIC rules (catalog docs/disk-validity-rules.md §7).
// Rules in this file check FT_SAM_BASIC invariants: FileTypeInfo triplets,
// VARS/gap sizes, program sentinel byte, line-number encoding, auto-RUN
// start-line validity, and MGTFlags convention. They apply to all dialects
// (BASIC-VARS-GAP-INVARIANT consults ctx.Dialect internally).

// bodyData reads the file body (excluding the 9-byte header) by
// walking fe's sector chain. Mirrors the chain-walk loop in
// (*DiskImage).File but without the filename-lookup wrapper, so
// callers that already have a *FileEntry don't re-iterate the
// directory. Returns ("body bytes", nil) on success or
// (nil, err) when a SectorData call fails — rules treat the error
// as "no finding" because Phase 3's §1/§3 rules already report the
// underlying chain problem.
//
// The returned slice is fe.Length() bytes long; it does NOT include
// the body-header bytes 0..8, matching the convention of samfile.File's
// Body field.
func bodyData(di *DiskImage, fe *FileEntry) ([]byte, error) {
	fileLength := fe.Length()
	raw := make([]byte, fileLength+9)
	sd, err := di.SectorData(fe.FirstSector)
	if err != nil {
		return nil, err
	}
	fp := sd.FilePart()
	i := uint16(0)
	for {
		copy(raw[510*i:], fp.Data[:])
		i++
		if i == fe.Sectors {
			break
		}
		sd, err = di.SectorData(fp.NextSector)
		if err != nil {
			return nil, err
		}
		fp = sd.FilePart()
	}
	return raw[9:], nil
}

// ----- BASIC-FILETYPEINFO-TRIPLETS -----
// For FT_SAM_BASIC, dir bytes 0xDD-0xE5 hold three 3-byte PAGEFORM
// lengths (cumulative offsets into the body): NVARS-PROG, NUMEND-PROG,
// SAVARS-PROG. The decoded values must be non-zero AND satisfy
// NVARS <= NUMEND <= SAVARS <= body Length.
func init() {
	Register(Rule{
		ID:          "BASIC-FILETYPEINFO-TRIPLETS",
		Severity:    SeverityStructural,
		Description: "FT_SAM_BASIC FileTypeInfo (dir 0xDD-0xE5) holds three non-zero, non-decreasing PAGEFORM cumulative offsets bounded by body length",
		Citation:    "rom-disasm:22163-22180",
		Check:       checkBasicFileTypeInfoTriplets,
	})
}

func checkBasicFileTypeInfoTriplets(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		nvars := fe.ProgramLength()                                              // decode(FileTypeInfo[0..2])
		numend := fe.ProgramLength() + fe.NumericVariablesSize()                // decode(FileTypeInfo[3..5])
		savars := fe.ProgramLength() + fe.NumericVariablesSize() + fe.GapSize() // decode(FileTypeInfo[6..8])
		length := fe.Length()
		if nvars == 0 || numend == 0 || savars == 0 {
			findings = append(findings, Finding{
				RuleID: "BASIC-FILETYPEINFO-TRIPLETS", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC file has zero offset in FileTypeInfo triplet (NVARS=%d NUMEND=%d SAVARS=%d)", nvars, numend, savars),
				Citation: "rom-disasm:22163-22180",
			})
			return
		}
		if !(nvars <= numend && numend <= savars && savars <= length) {
			findings = append(findings, Finding{
				RuleID: "BASIC-FILETYPEINFO-TRIPLETS", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC FileTypeInfo offsets out of order (NVARS=%d NUMEND=%d SAVARS=%d length=%d)", nvars, numend, savars, length),
				Citation: "rom-disasm:22163-22180",
			})
		}
	})
	return findings
}

// ----- BASIC-VARS-GAP-INVARIANT -----
// Empirically, SAMDOS-2 BASIC files have SAVARS-NVARS == 604, MasterDOS
// BASIC files have SAVARS-NVARS == 2156 (sam-basic-save-format.md, scan
// of 161 disks). Cosmetic; depends on detected dialect — on Unknown,
// accept either value.
func init() {
	Register(Rule{
		ID:          "BASIC-VARS-GAP-INVARIANT",
		Severity:    SeverityCosmetic,
		Description: "FT_SAM_BASIC SAVARS-NVARS equals the dialect-canonical value (604 SAMDOS-2 / 2156 MasterDOS)",
		Citation:    "sam-basic-save-format.md",
		Check:       checkBasicVarsGapInvariant,
	})
}

func checkBasicVarsGapInvariant(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		gap := fe.NumericVariablesSize() + fe.GapSize() // SAVARS - NVARS
		var expected uint32
		switch ctx.Dialect {
		case DialectSAMDOS2:
			expected = 604
		case DialectMasterDOS:
			expected = 2156
		default:
			// Unknown — accept either canonical value, silently skip.
			if gap == 604 || gap == 2156 {
				return
			}
			expected = 604 // for the message; prefer the SAMDOS-2 canonical value
		}
		if gap != expected {
			findings = append(findings, Finding{
				RuleID: "BASIC-VARS-GAP-INVARIANT", Severity: SeverityCosmetic,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC SAVARS-NVARS = %d; expected %d for dialect %s", gap, expected, ctx.Dialect),
				Citation: "sam-basic-save-format.md",
			})
		}
	})
	return findings
}

// ----- BASIC-PROG-END-SENTINEL -----
// The tokenised program ends with a 0xFF sentinel byte. The byte at
// body[ProgramLength-1] is the sentinel (NVARS-PROG is the program-area
// end offset).
func init() {
	Register(Rule{
		ID:          "BASIC-PROG-END-SENTINEL",
		Severity:    SeverityStructural,
		Description: "FT_SAM_BASIC program-area ends with a 0xFF sentinel byte",
		Citation:    "sambasic/file.go:36-42",
		Check:       checkBasicProgEndSentinel,
	})
}

func checkBasicProgEndSentinel(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		body, err := bodyData(ctx.Disk, fe)
		if err != nil {
			return
		}
		progLen := fe.ProgramLength()
		if progLen == 0 || int(progLen) > len(body) {
			return // BASIC-FILETYPEINFO-TRIPLETS will catch this
		}
		if body[progLen-1] != 0xFF {
			findings = append(findings, Finding{
				RuleID: "BASIC-PROG-END-SENTINEL", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC program does not end with 0xFF sentinel; body[%d] = 0x%02x", progLen-1, body[progLen-1]),
				Citation: "sambasic/file.go:36-42",
			})
		}
	})
	return findings
}

// ----- BASIC-LINE-NUMBER-BE -----
// Walk the program with sambasic.Parse; any parse failure means the
// big-endian line-number / little-endian length / 0x0D-terminator
// invariant doesn't hold somewhere. Also check each line number is
// non-zero (line 0 doesn't exist in SAM BASIC).
//
// Iteration 1 SCOPE: widened from 1..16383 to 1..65535. The 16383
// cap came from a samfile-specific 0x3FFF mask, but SAM BASIC stores
// line numbers as 16-bit big-endian — line numbers above 16383 are
// legitimate. Corpus evidence: 1,098 fires / 385 disks, with several
// hundred messages reading `line number 20000/50000/60000` — real
// "library" / "internal" line numbers in published BASIC programs.
// Line 0 (309 fires) remains structurally invalid: BASIC has no
// line 0, so it is dropped from the accepted range.
//
// Line.Number is `uint16`, so the upper bound 65535 is implicit in
// the type; only the `== 0` check remains as an explicit guard.
func init() {
	Register(Rule{
		ID:          "BASIC-LINE-NUMBER-BE",
		Severity:    SeverityStructural,
		Description: "FT_SAM_BASIC program parses cleanly and every line number is in 1..65535 (uint16 BE; widened from 1..16383 in iteration 1)",
		Citation:    "sambasic/parse.go",
		Check:       checkBasicLineNumberBE,
	})
}

func checkBasicLineNumberBE(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		body, err := bodyData(ctx.Disk, fe)
		if err != nil {
			return
		}
		progLen := fe.ProgramLength()
		if progLen == 0 || int(progLen) > len(body) {
			return
		}
		prog := body[:progLen]
		bf, err := sambasic.Parse(prog)
		if err != nil {
			findings = append(findings, Finding{
				RuleID: "BASIC-LINE-NUMBER-BE", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC program parse failed: %v", err),
				Citation: "sambasic/parse.go",
			})
			return
		}
		for _, ln := range bf.Lines {
			// Line 0 doesn't exist in SAM BASIC; the upper bound 65535
			// is implicit in Line.Number's uint16 type.
			if ln.Number == 0 {
				findings = append(findings, Finding{
					RuleID: "BASIC-LINE-NUMBER-BE", Severity: SeverityStructural,
					Location: SlotLocation(slot, fe.Name.String()),
					Message:  fmt.Sprintf("BASIC line number %d out of range (1..65535)", ln.Number),
					Citation: "sambasic/parse.go",
				})
				return // one finding per slot
			}
		}
	})
	return findings
}

// ----- BASIC-STARTLINE-FF-DISABLES -----
// dir[0xF2] (= fe.ExecutionAddressDiv16K) is 0x00 (auto-RUN) or 0xFF
// (no auto-RUN); when 0x00, dir[0xF3..0xF4] (= fe.SAMBASICStartLine) is
// a valid line number (1..16383, not 0xFFFF).
func init() {
	Register(Rule{
		ID:          "BASIC-STARTLINE-FF-DISABLES",
		Severity:    SeverityStructural,
		Description: "FT_SAM_BASIC dir[0xF2] is 0x00 (auto-RUN) or 0xFF (no auto-RUN); when 0x00, the start-line is a valid line number",
		Citation:    "rom-disasm:22136-22141",
		Check:       checkBasicStartLineFFDisables,
	})
}

func checkBasicStartLineFFDisables(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		marker := fe.ExecutionAddressDiv16K
		if marker != 0x00 && marker != 0xFF {
			findings = append(findings, Finding{
				RuleID: "BASIC-STARTLINE-FF-DISABLES", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC auto-RUN marker dir[0xF2] = 0x%02x (expected 0x00 or 0xFF)", marker),
				Citation: "rom-disasm:22136-22141",
			})
			return
		}
		if marker == 0x00 {
			line := fe.SAMBASICStartLine
			// Iteration 1 SCOPE: line numbers are 16-bit BE; widen the
			// accepted range from 1..16383 to 1..65534 (excluding 0xFFFF
			// which is the no-auto-RUN sentinel). Companion to
			// BASIC-LINE-NUMBER-BE.
			if line == 0 || line == 0xFFFF {
				findings = append(findings, Finding{
					RuleID: "BASIC-STARTLINE-FF-DISABLES", Severity: SeverityStructural,
					Location: SlotLocation(slot, fe.Name.String()),
					Message:  fmt.Sprintf("BASIC auto-RUN enabled (dir[0xF2]=0x00) but start-line %d is invalid (1..65534; 0xFFFF disables auto-RUN)", line),
					Citation: "rom-disasm:22136-22141",
				})
			}
		}
	})
	return findings
}

// ----- BASIC-STARTLINE-WITHIN-PROG -----
// When auto-RUN is enabled, the start-line should not exceed the
// highest line number in the saved program. ROM BASIC's RUN N uses
// NEXT-LINE-GE semantics (the lookup finds the first line whose
// number is >= N), so `RUN N` where N is less than or equal to the
// highest line in the program starts at the first line at or after
// N — this is exactly the canonical "RUN 1 to start from the
// beginning" idiom. Only `RUN N` where N is greater than every
// saved line is a real bug: there's no line at or after N, and
// BASIC errors out.
//
// Iteration 1 REWORD: previously fired on "start-line not present
// in the saved program" — which mis-described the rule's intent
// because the canonical "RUN 1 with first line 10" pattern (78% of
// 3,074 corpus fires) is not an error, just a marker for "start
// from the beginning". The catalog's "Statement lost" framing was
// also wrong: SAM BASIC's NEXT-LINE-GE lookup means RUN N at or
// below the lowest line is a no-op, not an error.
//
// Severity stays cosmetic.
func init() {
	Register(Rule{
		ID:          "BASIC-STARTLINE-WITHIN-PROG",
		Severity:    SeverityCosmetic,
		Description: "FT_SAM_BASIC auto-RUN start-line is at or below the highest saved line (RUN's NEXT-LINE-GE lookup tolerates start-lines below the lowest saved line)",
		Citation:    "rom-disasm:22136-22141",
		Check:       checkBasicStartLineWithinProg,
	})
}

func checkBasicStartLineWithinProg(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		if fe.ExecutionAddressDiv16K != 0x00 {
			return // auto-RUN disabled; nothing to check
		}
		body, err := bodyData(ctx.Disk, fe)
		if err != nil {
			return
		}
		progLen := fe.ProgramLength()
		if progLen == 0 || int(progLen) > len(body) {
			return
		}
		bf, err := sambasic.Parse(body[:progLen])
		if err != nil {
			return // BASIC-LINE-NUMBER-BE reports the parse failure
		}
		if len(bf.Lines) == 0 {
			return // empty program; BASIC-LINE-NUMBER-BE / parse rules cover this
		}
		want := fe.SAMBASICStartLine
		// Find the highest line number in the saved program.
		var highest uint16
		for _, ln := range bf.Lines {
			if ln.Number > highest {
				highest = ln.Number
			}
		}
		// Iteration 1 REWORD: SAM BASIC RUN N uses NEXT-LINE-GE
		// semantics. `want` at or below `highest` always resolves
		// to a saved line (the first line whose number is >= want);
		// only `want > highest` produces no line to run.
		if want > highest {
			findings = append(findings, Finding{
				RuleID: "BASIC-STARTLINE-WITHIN-PROG", Severity: SeverityCosmetic,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC auto-RUN line %d is greater than the highest saved line %d; BASIC's NEXT-LINE-GE lookup will find no line to run", want, highest),
				Citation: "rom-disasm:22136-22141",
			})
		}
	})
	return findings
}

// ----- BASIC-MGTFLAGS-20 -----
// Real-world BASIC files have MGTFlags == 0x20 (empirical convention,
// 50%+ of canonical disks, required for M0 boot per
// test-mgt-byte-layout.md §slot 1). Inconsistency severity.
func init() {
	Register(Rule{
		ID:          "BASIC-MGTFLAGS-20",
		Severity:    SeverityInconsistency,
		Description: "FT_SAM_BASIC MGTFlags is 0x20 (empirical convention)",
		Citation:    "test-mgt-byte-layout.md",
		Check:       checkBasicMGTFlags20,
	})
}

func checkBasicMGTFlags20(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SAM_BASIC {
			return
		}
		if fe.MGTFlags != 0x20 {
			findings = append(findings, Finding{
				RuleID: "BASIC-MGTFLAGS-20", Severity: SeverityInconsistency,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("BASIC file MGTFlags = 0x%02x; expected 0x20 (empirical convention)", fe.MGTFlags),
				Citation: "test-mgt-byte-layout.md",
			})
		}
	})
	return findings
}
