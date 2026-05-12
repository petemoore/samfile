// rules_ft_zxsnap.go
package samfile

import "fmt"

// §10 ZX snapshot rules (catalog docs/disk-validity-rules.md §10).
// Rules in this file check FT_ZX_SNAPSHOT (5) invariants per
// SAMDOS-2's snapshot save convention: 49152-byte body and 0x4000
// load address (samdos/src/d.s:629-630, :612 → :662).
//
// Scoped to DialectSAMDOS2: corpus data (audit run 2026-05-12)
// showed both rules failing 192/193 times across 25 disks. Of those
// 25, 22 are masterdos-dialect and 3 unknown-dialect; zero are
// samdos2. The masterdos ZX-snapshot entries are typically +D /
// Disciple-format imports that use a different Pages/PageOffset
// encoding — decoded Length/Start don't match 49152/0x4000 but the
// files are still loadable via their original-OS loader. The rules
// are correct for SAMDOS-2 but over-strict elsewhere.

// ----- ZXSNAP-LENGTH-49152 -----
// FT_ZX_SNAPSHOT has a 49,152-byte body (48 KiB ZX RAM).
func init() {
	Register(Rule{
		ID:            "ZXSNAP-LENGTH-49152",
		Severity:      SeverityStructural,
		Dialects:      []Dialect{DialectSAMDOS2},
		Description:   "FT_ZX_SNAPSHOT body is exactly 49152 bytes (48 KiB ZX RAM); SAMDOS-2 only",
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
		Dialects:      []Dialect{DialectSAMDOS2},
		Description:   "FT_ZX_SNAPSHOT decoded start address is 0x4000 (16384, ZX RAM base); SAMDOS-2 only",
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
