// rules_ft_zxsnap.go
package samfile

import "fmt"

// §10 ZX snapshot rules (catalog docs/disk-validity-rules.md §10).
// Rules in this file check FT_ZX_SNAPSHOT (5) invariants: 48 KiB
// body length and 0x4000 load address. The catalog tags these as
// SAMDOS-2 specific (the constants live in SAMDOS source); we run
// them on all dialects because the ZX snapshot format is itself
// dialect-agnostic.

// ----- ZXSNAP-LENGTH-49152 -----
// FT_ZX_SNAPSHOT has a 49,152-byte body (48 KiB ZX RAM).
func init() {
	Register(Rule{
		ID:            "ZXSNAP-LENGTH-49152",
		Severity:      SeverityStructural,
		Description:   "FT_ZX_SNAPSHOT body is exactly 49152 bytes (48 KiB ZX RAM)",
		Citation:      "samdos/src/d.s:660-661",
		Check:         checkZXSnapLength49152,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: typedSlot(FT_ZX_SNAPSHOT)},
	})
}

func checkZXSnapLength49152(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_ZX_SNAPSHOT {
			return
		}
		if fe.Length() != 49152 {
			findings = append(findings, Finding{
				RuleID: "ZXSNAP-LENGTH-49152", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("ZX snapshot body length = %d; expected 49152", fe.Length()),
				Citation: "samdos/src/d.s:660-661",
			})
		}
	})
	return findings
}

// ----- ZXSNAP-LOAD-ADDR-16384 -----
// FT_ZX_SNAPSHOT load address is 0x4000 (ZX RAM base).
func init() {
	Register(Rule{
		ID:            "ZXSNAP-LOAD-ADDR-16384",
		Severity:      SeverityStructural,
		Description:   "FT_ZX_SNAPSHOT decoded start address is 0x4000 (16384, ZX RAM base)",
		Citation:      "samdos/src/d.s:660-663",
		Check:         checkZXSnapLoadAddr16384,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: typedSlot(FT_ZX_SNAPSHOT)},
	})
}

func checkZXSnapLoadAddr16384(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_ZX_SNAPSHOT {
			return
		}
		if fe.StartAddress() != 16384 {
			findings = append(findings, Finding{
				RuleID: "ZXSNAP-LOAD-ADDR-16384", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("ZX snapshot start address = 0x%05x; expected 0x4000 (16384)", fe.StartAddress()),
				Citation: "samdos/src/d.s:660-663",
			})
		}
	})
	return findings
}
