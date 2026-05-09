// rules_ft_array.go
package samfile

// §8 Array rules (catalog docs/disk-validity-rules.md §8).
// Rules in this file check FT_NUM_ARRAY (17) and FT_STR_ARRAY (18)
// invariants. They apply to all dialects.

// ----- ARRAY-FILETYPEINFO-TLBYTE-NAME -----
// For FT_NUM_ARRAY (17) and FT_STR_ARRAY (18), dir bytes 0xDD-0xE7
// hold the array's TLBYTE (type/length byte) followed by its 10-byte
// name. The rule warns when all 11 bytes are zero — that indicates a
// writer didn't populate the array metadata at SAVE time.
func init() {
	Register(Rule{
		ID:          "ARRAY-FILETYPEINFO-TLBYTE-NAME",
		Severity:    SeverityStructural,
		Description: "FT_NUM_ARRAY/FT_STR_ARRAY FileTypeInfo (dir 0xDD-0xE7) is not all zero",
		Citation:    "rom-disasm:22354-22357",
		Check:       checkArrayFileTypeInfoTLBYTEName,
	})
}

func checkArrayFileTypeInfoTLBYTEName(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_NUM_ARRAY && fe.Type != FT_STR_ARRAY {
			return
		}
		allZero := true
		for _, b := range fe.FileTypeInfo {
			if b != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			findings = append(findings, Finding{
				RuleID: "ARRAY-FILETYPEINFO-TLBYTE-NAME", Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  "array file FileTypeInfo (dir 0xDD-0xE7) is all zero; TLBYTE + name not populated",
				Citation: "rom-disasm:22354-22357",
			})
		}
	})
	return findings
}
