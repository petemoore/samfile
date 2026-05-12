// rules_ft_screen.go
package samfile

import "fmt"

// §9 SCREEN rules (catalog docs/disk-validity-rules.md §9).
// Rules in this file check FT_SCREEN (20) invariants: mode byte
// and body-length-vs-mode geometry. They apply to all dialects.

// ----- SCREEN-MODE-AT-0xDD -----
// For FT_SCREEN, dir byte 0xDD (= FileTypeInfo[0]) is the screen mode
// (1-4 on SAM).
func init() {
	Register(Rule{
		ID:          "SCREEN-MODE-AT-0xDD",
		Severity:    SeverityStructural,
		Description: "FT_SCREEN dir[0xDD] (mode byte) is in 1..4",
		Citation:    "rom-disasm:22259",
		Check:       checkScreenModeAt0xDD,
	})
}

func checkScreenModeAt0xDD(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SCREEN {
			return
		}
		mode := fe.FileTypeInfo[0]
		if mode < 1 || mode > 4 {
			findings = append(findings, Finding{
				RuleID: "SCREEN-MODE-AT-0xDD", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("SCREEN mode byte = %d (expected 1..4)", mode),
				Citation: "rom-disasm:22259",
			})
		}
	})
	return findings
}

// ----- SCREEN-LENGTH-MATCHES-MODE -----
// For FT_SCREEN, body Length() must be at least the documented
// screen-data size for the given mode: modes 1 and 2 → 6912 bytes,
// modes 3 and 4 → 24576 bytes. ROM SCREEN$ SAVE typically appends a
// palette + sysvars trailer (16 bytes of CLUT + LINE/ATTR/state) so
// real-world MODE 3/4 screens are commonly 24576+41 = 24617 bytes;
// LOAD SCREEN$ ignores the trailer. (Skipped when mode is
// out-of-range; SCREEN-MODE-AT-0xDD catches that.)
//
// Iteration 1 REWORD: previously required strict equality with the
// canonical size, which fired on the 75% of corpus MODE 3/4
// screens that include the standard ROM-SAVE palette trailer.
// New rule: fire when Length < min (genuinely-truncated) or
// Length > min + 512 (suspiciously-long). The 512-byte slack
// accommodates the documented trailer plus reasonable wiggle room
// while still catching the rare mode-mismatch case (e.g. mode 2
// dir byte but mode-3-sized body).
func init() {
	Register(Rule{
		ID:          "SCREEN-LENGTH-MATCHES-MODE",
		Severity:    SeverityStructural,
		Description: "FT_SCREEN body length is within [min, min+512] for its mode (1-2: min=6912; 3-4: min=24576) — slack accommodates the ROM SCREEN$ SAVE palette+sysvars trailer",
		Citation:    "sam-coupe_tech-man_v3-0.txt",
		Check:       checkScreenLengthMatchesMode,
	})
}

func checkScreenLengthMatchesMode(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_SCREEN {
			return
		}
		mode := fe.FileTypeInfo[0]
		var minSize uint32
		switch mode {
		case 1, 2:
			minSize = 6912
		case 3, 4:
			minSize = 24576
		default:
			return // SCREEN-MODE-AT-0xDD reports the bad mode
		}
		const trailerSlack = 512 // palette + sysvars + headroom
		length := fe.Length()
		switch {
		case length < minSize:
			findings = append(findings, Finding{
				RuleID: "SCREEN-LENGTH-MATCHES-MODE", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("SCREEN mode %d body length %d is shorter than mode minimum %d bytes", mode, length, minSize),
				Citation: "sam-coupe_tech-man_v3-0.txt",
			})
		case length > minSize+trailerSlack:
			findings = append(findings, Finding{
				RuleID: "SCREEN-LENGTH-MATCHES-MODE", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("SCREEN mode %d body length %d exceeds mode maximum %d bytes (min %d + %d trailer slack)", mode, length, minSize+trailerSlack, minSize, trailerSlack),
				Citation: "sam-coupe_tech-man_v3-0.txt",
			})
		}
	})
	return findings
}
