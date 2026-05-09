package samfile

import (
	"fmt"
	"math/bits"
)

// §2 Directory-entry rules (catalog docs/disk-validity-rules.md §2).
// Rules in this file check internal consistency of each of the 80
// directory entries: type byte, filename padding, sector count vs
// chain length vs SectorAddressMap popcount. They apply to all
// dialects.

// forEachUsedSlot loops over every used directory slot in registration order
// and invokes fn for each. A small helper that keeps the per-rule
// Check function's loop body focused on the actual invariant.
func forEachUsedSlot(ctx *CheckContext, fn func(slot int, fe *FileEntry)) {
	for _, slot := range ctx.Journal.UsedFileEntries() {
		fn(slot, ctx.Journal[slot])
	}
}

// ----- DIR-TYPE-BYTE-IS-KNOWN -----
func init() {
	Register(Rule{
		ID:          "DIR-TYPE-BYTE-IS-KNOWN",
		Severity:    SeverityInconsistency,
		Description: "directory type byte (low 5 bits, attribute bits masked) is one of the documented file types",
		Citation:    "samdos/src/e.s:322-355",
		Check:       checkDirTypeByteIsKnown,
	})
}

// dirKnownTypes is the SAM-public set after masking off HIDDEN + PROTECTED.
// 0 is omitted: erased slots are caught by Used(), not here.
var dirKnownTypes = map[uint8]bool{
	5: true, 16: true, 17: true, 18: true, 19: true, 20: true,
}

func checkDirTypeByteIsKnown(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		t := uint8(fe.Type) & 0x1F
		// Iteration 1 FIX: type byte 0 is the erased-slot sentinel
		// and is handled (with the structural severity it deserves)
		// by DIR-ERASED-IS-ZERO. The catalog's test sketch explicitly
		// lists 0 in the accepted set, so this rule should not also
		// fire on it. Skipping here removes a 100% double-fire across
		// 2,492 corpus findings.
		if t == 0 {
			return
		}
		if !dirKnownTypes[t] {
			findings = append(findings, Finding{
				RuleID:   "DIR-TYPE-BYTE-IS-KNOWN",
				Severity: SeverityInconsistency,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("masked type byte 0x%02x is not a documented file type (expected one of 5, 16-20)", t),
				Citation: "samdos/src/e.s:322-355",
			})
		}
	})
	return findings
}

// ----- DIR-ERASED-IS-ZERO -----
// Iteration 2 DEMOTE + REWORD (structural → inconsistency). The rule
// fires on slots where the type byte is 0 (so fdhf at c.s:1133-1143
// treats it as free) but FirstSector.Track is non-zero (so name /
// chain / SAM are still populated). That's the canonical
// "DEL/ERASE leaves the file recoverable" archaeology: SAMDOS DEL
// (and ROM ERASE) zero only the type byte, leaving the rest of the
// dir entry intact. SAMDOS treats the slot as free per fdhf, so the
// dir walk is well-defined — but the orphaned filename + chain
// disagree with the type byte's "this slot is unused" claim. That's
// "two views of the same fact disagree" → inconsistency, not
// "disk-walk invariant violated" → structural.
//
// 43% of the corpus has at least one such slot. Keeping it at
// structural drowns out genuine structural corruption (e.g.
// CHAIN-NO-CYCLE, 3 disks).
//
// The message text is also reworded: the previous "used slot"
// framing was self-contradictory (Used() returns true here, but
// SAMDOS itself treats the slot as free), and obscured the
// archaeological pattern. The new wording names the actual
// signature ("type byte 0x00 (erased) but filename/chain are still
// populated") and labels the pattern ("probably a DEL'd file with
// recoverable header").
func init() {
	Register(Rule{
		ID:          "DIR-ERASED-IS-ZERO",
		Severity:    SeverityInconsistency,
		Description: "a used directory slot has a non-zero type byte",
		Citation:    "samdos/src/c.s:1133-1143",
		Check:       checkDirErasedIsZero,
	})
}

func checkDirErasedIsZero(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if uint8(fe.Type) == 0 {
			findings = append(findings, Finding{
				RuleID:   "DIR-ERASED-IS-ZERO",
				Severity: SeverityInconsistency,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  "slot has type byte 0x00 (erased) but filename/chain are still populated (probably a DEL'd file with recoverable header)",
				Citation: "samdos/src/c.s:1133-1143",
			})
		}
	})
	return findings
}

// ----- DIR-NAME-PADDING -----
func init() {
	Register(Rule{
		ID:          "DIR-NAME-PADDING",
		Severity:    SeverityCosmetic,
		Description: "filename bytes are printable ASCII or space-padded",
		Citation:    "sam-coupe_tech-man_v3-0.txt:4358-4359",
		Check:       checkDirNamePadding,
	})
}

func checkDirNamePadding(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		for i, b := range fe.Name {
			if b == 0x20 || (b >= 0x21 && b < 0x7F) {
				continue
			}
			findings = append(findings, Finding{
				RuleID:   "DIR-NAME-PADDING",
				Severity: SeverityCosmetic,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("filename byte %d is 0x%02x (expected printable ASCII or 0x20 space)", i, b),
				Citation: "sam-coupe_tech-man_v3-0.txt:4358-4359",
			})
			return // one finding per slot; further byte-by-byte detail belongs in a diagnostic
		}
	})
	return findings
}

// ----- DIR-NAME-NOT-EMPTY -----
func init() {
	Register(Rule{
		ID:          "DIR-NAME-NOT-EMPTY",
		Severity:    SeverityInconsistency,
		Description: "a used slot has at least one non-space, non-FF character in its 10-byte name",
		Citation:    "rom-disasm:22093-22105",
		Check:       checkDirNameNotEmpty,
	})
}

func checkDirNameNotEmpty(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		empty := true
		for _, b := range fe.Name {
			if b != 0x20 && b != 0xFF && b != 0 {
				empty = false
				break
			}
		}
		if empty {
			findings = append(findings, Finding{
				RuleID:   "DIR-NAME-NOT-EMPTY",
				Severity: SeverityInconsistency,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  "filename is all spaces / 0xFF / 0x00 (no visible characters)",
				Citation: "rom-disasm:22093-22105",
			})
		}
	})
	return findings
}

// ----- DIR-FIRST-SECTOR-VALID -----
func init() {
	Register(Rule{
		ID:          "DIR-FIRST-SECTOR-VALID",
		Severity:    SeverityFatal,
		Description: "directory entry's FirstSector points at a valid data sector",
		Citation:    "samfile.go:611-616",
		Check:       checkDirFirstSectorValid,
	})
}

func checkDirFirstSectorValid(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		fs := fe.FirstSector
		t := fs.Track
		s := fs.Sector
		// Directory area is side 0 cylinders 0..3 only (Tech Manual L4340-4343);
		// side 1 cylinders 0..3 (tracks 0x80..0x83) are valid data sectors.
		// Side 0 (0x00..0x4F): cylinders 0..3 are the directory area, 4..79 are data.
		// Side 1 (0x80..0xCF): all 80 cylinders are data.
		validTrack := (t >= 4 && t < 80) || (t >= 128 && t < 208)
		validSector := s >= 1 && s <= 10
		if !validTrack || !validSector {
			findings = append(findings, Finding{
				RuleID:   "DIR-FIRST-SECTOR-VALID",
				Severity: SeverityFatal,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("FirstSector (track=0x%02x, sector=%d) is not a valid data sector", t, s),
				Citation: "samfile.go:611-616",
			})
		}
	})
	return findings
}

// ----- DIR-SECTORS-MATCHES-CHAIN -----
func init() {
	Register(Rule{
		ID:          "DIR-SECTORS-MATCHES-CHAIN",
		Severity:    SeverityStructural,
		Description: "dir-entry Sectors count equals the number of sectors visited walking the chain to the (0,0) terminator",
		Citation:    "samfile.go:743-754",
		Check:       checkDirSectorsMatchesChain,
	})
}

func checkDirSectorsMatchesChain(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		result := walkChain(ctx.Disk, fe.FirstSector)
		count := uint16(len(result.Steps))
		if count != fe.Sectors {
			findings = append(findings, Finding{
				RuleID:   "DIR-SECTORS-MATCHES-CHAIN",
				Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("dir Sectors=%d, but chain walk visited %d sectors", fe.Sectors, count),
				Citation: "samfile.go:743-754",
			})
		}
	})
	return findings
}

// ----- DIR-SECTORS-MATCHES-MAP -----
func init() {
	Register(Rule{
		ID:          "DIR-SECTORS-MATCHES-MAP",
		Severity:    SeverityStructural,
		Description: "dir-entry Sectors count equals the popcount of the per-slot SectorAddressMap",
		Citation:    "sam-coupe_tech-man_v3-0.txt:4405-4414",
		Check:       checkDirSectorsMatchesMap,
	})
}

func checkDirSectorsMatchesMap(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		pop := 0
		for _, b := range fe.SectorAddressMap {
			pop += bits.OnesCount8(b)
		}
		if uint16(pop) != fe.Sectors {
			findings = append(findings, Finding{
				RuleID:   "DIR-SECTORS-MATCHES-MAP",
				Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("dir Sectors=%d, but SectorAddressMap has popcount=%d", fe.Sectors, pop),
				Citation: "sam-coupe_tech-man_v3-0.txt:4405-4414",
			})
		}
	})
	return findings
}

// ----- DIR-SECTORS-NONZERO -----
func init() {
	Register(Rule{
		ID:          "DIR-SECTORS-NONZERO",
		Severity:    SeverityStructural,
		Description: "a used dir entry's Sectors count is at least 1",
		Citation:    "samdos/src/c.s:919-951",
		Check:       checkDirSectorsNonzero,
	})
}

func checkDirSectorsNonzero(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Sectors == 0 {
			findings = append(findings, Finding{
				RuleID:   "DIR-SECTORS-NONZERO",
				Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  "used slot has Sectors=0 (must be at least 1 for the body header)",
				Citation: "samdos/src/c.s:919-951",
			})
		}
	})
	return findings
}

// ----- DIR-SAM-WITHIN-CAPACITY -----
// 195 bytes × 8 = 1560 bits matches the disk's data-sector count
// exactly; the catalog's test sketch checks byte 194 & 0xE0 == 0
// (top 3 bits of byte 194), per Tech Manual L4405-4406.
func init() {
	Register(Rule{
		ID:          "DIR-SAM-WITHIN-CAPACITY",
		Severity:    SeverityInconsistency,
		Description: "SectorAddressMap byte 194's top 3 bits (1557-1559) are clear (no sector beyond disk capacity)",
		Citation:    "sam-coupe_tech-man_v3-0.txt:4405-4406",
		Check:       checkDirSAMWithinCapacity,
	})
}

func checkDirSAMWithinCapacity(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.SectorAddressMap[194]&0xE0 != 0 {
			findings = append(findings, Finding{
				RuleID:   "DIR-SAM-WITHIN-CAPACITY",
				Severity: SeverityInconsistency,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("SectorAddressMap[194]=0x%02x has bits beyond bit 1559 set", fe.SectorAddressMap[194]),
				Citation: "sam-coupe_tech-man_v3-0.txt:4405-4406",
			})
		}
	})
	return findings
}
