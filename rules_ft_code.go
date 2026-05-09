package samfile

import "fmt"

// §6 FT_CODE rules (catalog docs/disk-validity-rules.md §6).
// Rules in this file check FT_CODE-specific invariants: the file's
// load address is above ROM, the loaded region fits in SAM's 512 KiB
// address space, the execution address (if not opted out) lies within
// the loaded region, and dir-entry FileTypeInfo is unused (cosmetic).
// Each Check function filters on fe.Type == FT_CODE at the top.

// ----- CODE-LOAD-ABOVE-ROM -----
func init() {
	Register(Rule{
		ID:          "CODE-LOAD-ABOVE-ROM",
		Severity:    SeverityFatal,
		Description: "FT_CODE file's load address is at least 0x4000 (above ROM)",
		Citation:    "samfile.go:799-801",
		Check:       checkCodeLoadAboveROM,
	})
}

func checkCodeLoadAboveROM(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_CODE {
			return
		}
		loadAddr := fe.StartAddress()
		if loadAddr < 0x4000 {
			findings = append(findings, Finding{
				RuleID:   "CODE-LOAD-ABOVE-ROM",
				Severity: SeverityFatal,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("CODE load address 0x%05x is below 0x4000 (ROM)", loadAddr),
				Citation: "samfile.go:799-801",
			})
		}
	})
	return findings
}

// ----- CODE-LOAD-FITS-IN-MEMORY -----
func init() {
	Register(Rule{
		ID:          "CODE-LOAD-FITS-IN-MEMORY",
		Severity:    SeverityFatal,
		Description: "FT_CODE file's load address + body length does not exceed SAM's 512 KiB address space",
		Citation:    "samfile.go:802-804",
		Check:       checkCodeLoadFitsInMemory,
	})
}

func checkCodeLoadFitsInMemory(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_CODE {
			return
		}
		loadAddr := fe.StartAddress()
		length := fe.Length()
		if uint64(loadAddr)+uint64(length) > 0x80000 {
			findings = append(findings, Finding{
				RuleID:   "CODE-LOAD-FITS-IN-MEMORY",
				Severity: SeverityFatal,
				Location: SlotLocation(slot, fe.Name.String()),
				Message: fmt.Sprintf("CODE load 0x%05x + length 0x%05x = 0x%05x exceeds SAM's 512 KiB address space",
					loadAddr, length, uint64(loadAddr)+uint64(length)),
				Citation: "samfile.go:802-804",
			})
		}
	})
	return findings
}

// ----- CODE-EXEC-WITHIN-LOADED-RANGE -----
func init() {
	Register(Rule{
		ID:          "CODE-EXEC-WITHIN-LOADED-RANGE",
		Severity:    SeverityStructural,
		Description: "FT_CODE file's execution address (when not 0xFF-disabled) lies within its loaded region",
		Citation:    "samfile.go:805-810",
		Check:       checkCodeExecWithinLoadedRange,
	})
}

func checkCodeExecWithinLoadedRange(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_CODE {
			return
		}
		if fe.ExecutionAddressDiv16K == 0xFF {
			return // 0xFF marker = no auto-exec; nothing to validate
		}
		execAddr := fe.ExecutionAddress()
		loadAddr := fe.StartAddress()
		length := fe.Length()
		if execAddr < loadAddr || execAddr >= loadAddr+length {
			findings = append(findings, Finding{
				RuleID:   "CODE-EXEC-WITHIN-LOADED-RANGE",
				Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message: fmt.Sprintf("CODE exec address 0x%05x is outside loaded region [0x%05x, 0x%05x)",
					execAddr, loadAddr, loadAddr+length),
				Citation: "samfile.go:805-810",
			})
		}
	})
	return findings
}

// ----- CODE-FILETYPEINFO-EMPTY -----
func init() {
	Register(Rule{
		ID:          "CODE-FILETYPEINFO-EMPTY",
		Severity:    SeverityCosmetic,
		Description: "FT_CODE file's FileTypeInfo (dir 0xDD-0xE7) is all zero (samfile convention)",
		Citation:    "samfile.go:798-827",
		Check:       checkCodeFileTypeInfoEmpty,
	})
}

func checkCodeFileTypeInfoEmpty(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_CODE {
			return
		}
		for _, b := range fe.FileTypeInfo {
			if b != 0 {
				findings = append(findings, Finding{
					RuleID:   "CODE-FILETYPEINFO-EMPTY",
					Severity: SeverityCosmetic,
					Location: SlotLocation(slot, fe.Name.String()),
					Message:  "CODE file has non-zero FileTypeInfo (dir 0xDD-0xE7) — samfile leaves these zero",
					Citation: "samfile.go:798-827",
				})
				return // one finding per slot
			}
		}
	})
	return findings
}
