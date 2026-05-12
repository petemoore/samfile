package samfile

import "fmt"

// §14 Cosmetic / canonical-output rules (catalog docs/disk-validity-rules.md §14).
// Rules in this file warn when dir-entry bytes diverge from
// the conventions real ROM SAVE produces, without affecting
// runtime behaviour. They apply to all dialects.

// ----- COSMETIC-RESERVEDA-FF -----
// Real ROM SAVE 0xFF-fills 14 bytes from dir offset 0xDC (HDCLP2 at
// rom-disasm:22076-22080), which covers MGTFlags + FileTypeInfo + the
// first two bytes of ReservedA. The catalog describes ReservedA (dir
// 0xE8-0xEB, 4 bytes) as fully 0xFF-filled by real SAVE. samfile's
// AddCodeFile leaves ReservedA at struct-zero (0x00). Both
// conventions are observed in the wild; the rule warns only when a
// byte is in NEITHER set — i.e. anything outside {0x00, 0xFF}.
//
// Same dual-acceptance pattern as Phase 4's CODE-FILETYPEINFO-EMPTY:
// real-ROM-SAVE byte == 0xFF, samfile byte == 0x00, both legitimate.
func init() {
	Register(Rule{
		ID:            "COSMETIC-RESERVEDA-FF",
		Severity:      SeverityCosmetic,
		Description:   "ReservedA (dir 0xE8-0xEB) is uniformly 0x00 (samfile) or 0xFF (ROM SAMDOS-2)",
		Citation:      "rom-disasm:22076-22080",
		Check:         checkCosmeticReservedAFF,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: usedSlot},
	})
}

func checkCosmeticReservedAFF(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		for i, b := range fe.ReservedA {
			if b != 0x00 && b != 0xFF {
				findings = append(findings, Finding{
					RuleID:   "COSMETIC-RESERVEDA-FF",
					Severity: SeverityCosmetic,
					Location: SlotLocation(slot, fe.Name.String()),
					Message:  fmt.Sprintf("ReservedA[%d] (dir 0x%02x) = 0x%02x — neither samfile's 0x00 nor ROM SAVE's 0xFF", i, 0xE8+i, b),
					Citation: "rom-disasm:22076-22080",
				})
				return // one finding per slot
			}
		}
	})
	return findings
}
