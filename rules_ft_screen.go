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
// For FT_SCREEN, body Length() matches the documented screen size for
// the given mode: modes 1 and 2 → 6912 bytes, modes 3 and 4 → 24576
// bytes. (Skipped when mode is out-of-range; SCREEN-MODE-AT-0xDD
// catches that.)
func init() {
	Register(Rule{
		ID:          "SCREEN-LENGTH-MATCHES-MODE",
		Severity:    SeverityStructural,
		Description: "FT_SCREEN body length matches the documented size for its mode (1-2: 6912; 3-4: 24576)",
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
		var expected uint32
		switch mode {
		case 1, 2:
			expected = 6912
		case 3, 4:
			expected = 24576
		default:
			return // SCREEN-MODE-AT-0xDD reports the bad mode
		}
		if fe.Length() != expected {
			findings = append(findings, Finding{
				RuleID: "SCREEN-LENGTH-MATCHES-MODE", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("SCREEN mode %d expects body length %d; got %d", mode, expected, fe.Length()),
				Citation: "sam-coupe_tech-man_v3-0.txt",
			})
		}
	})
	return findings
}
