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
// FileTypeInfo (dir 0xDD-0xE7) is unused for FT_CODE. Three conventions
// are observed in the wild:
//
//   - 0x00 × 11 — samfile's AddCodeFile leaves the struct zero-init.
//   - 0xFF × 11 — real ROM SAMDOS-2 SAVE 0xFF-fills 14 bytes from
//     dir offset 0xDC via HDCLP2 (rom-disasm:22076-22080), which
//     covers MGTFlags + the entire FileTypeInfo region + the first
//     two bytes of ReservedA.
//   - 0x20 — HDR space-fill leakage: ROM `HDCLP` at
//     rom-disasm:22070-22074 fills 25 bytes of the HDR header buffer
//     with 0x20 (space) during the names-area initialisation. When
//     the dir entry is written, FileTypeInfo bytes 0xDD-0xE7 land in
//     a region the HDR buffer overlaps, and the space-fill bleeds
//     through for many writers. Empirical: 16,727 of 16,932 corpus
//     fires (99%) are byte 0x20 — the dominance is striking enough
//     that 0x20 is plainly a third canonical "unused" marker.
//
// All three are legitimate "unused" markers. The rule warns only
// when a byte is in NONE of these conventions — anything other than
// 0x00, 0xFF, or 0x20 — which would suggest the slot was once a
// different file type whose FileTypeInfo bytes weren't cleared on
// overwrite.
func init() {
	Register(Rule{
		ID:          "CODE-FILETYPEINFO-EMPTY",
		Severity:    SeverityCosmetic,
		Description: "FT_CODE file's FileTypeInfo (dir 0xDD-0xE7) is uniformly 0x00 (samfile), 0xFF (ROM SAMDOS-2), or 0x20 (HDR space-fill leakage) — unused-marker convention",
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
			// Iteration 1 SCOPE: 0x20 is accepted as a third "unused"
			// marker (ROM HDCLP names-area space-fill leakage, 99% of
			// corpus fires).
			if b != 0x00 && b != 0xFF && b != 0x20 {
				findings = append(findings, Finding{
					RuleID:   "CODE-FILETYPEINFO-EMPTY",
					Severity: SeverityCosmetic,
					Location: SlotLocation(slot, fe.Name.String()),
					Message:  fmt.Sprintf("CODE file has 0x%02x in FileTypeInfo (dir 0xDD-0xE7) — none of samfile's 0x00, ROM SAVE's 0xFF, or HDR space-fill 0x20", b),
					Citation: "samfile.go:798-827",
				})
				return // one finding per slot
			}
		}
	})
	return findings
}
