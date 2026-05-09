package samfile

import (
	"fmt"
	"strings"
)

// §4 Cross-entry consistency rules (catalog docs/disk-validity-rules.md
// §4). Rules in this file compare data across multiple directory
// slots: shared sectors, duplicate names, references into the
// directory area. They apply to all dialects.

// ----- CROSS-NO-SECTOR-OVERLAP -----
func init() {
	Register(Rule{
		ID:          "CROSS-NO-SECTOR-OVERLAP",
		Severity:    SeverityFatal,
		Description: "no two used files claim the same data sector",
		Citation:    "samdos/src/c.s:895-951",
		Check:       checkCrossNoSectorOverlap,
	})
}

func checkCrossNoSectorOverlap(ctx *CheckContext) []Finding {
	var findings []Finding
	// owner[sector] = list of (slot, filename) entries claiming this sector.
	type claim struct {
		Slot int
		Name string
	}
	owner := make(map[Sector][]claim)
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		name := fe.Name.String()
		for _, sec := range fe.SectorAddressMap.UsedSectors() {
			owner[*sec] = append(owner[*sec], claim{slot, name})
		}
	})
	for sec, claims := range owner {
		if len(claims) < 2 {
			continue
		}
		s := sec
		findings = append(findings, Finding{
			RuleID:   "CROSS-NO-SECTOR-OVERLAP",
			Severity: SeverityFatal,
			Location: SectorLocation(claims[0].Slot, claims[0].Name, &s, -1),
			Message: fmt.Sprintf("sector %v is claimed by %d slots (first: %d %q, second: %d %q)",
				s, len(claims), claims[0].Slot, claims[0].Name, claims[1].Slot, claims[1].Name),
			Citation: "samdos/src/c.s:895-951",
		})
	}
	return findings
}

// ----- CROSS-NO-DUPLICATE-NAMES -----
func init() {
	Register(Rule{
		ID:          "CROSS-NO-DUPLICATE-NAMES",
		Severity:    SeverityInconsistency,
		Description: "no two used directory entries share the same filename (case-insensitive)",
		Citation:    "samdos/src/c.s:1196-1219",
		Check:       checkCrossNoDuplicateNames,
	})
}

func checkCrossNoDuplicateNames(ctx *CheckContext) []Finding {
	var findings []Finding
	seen := make(map[string]int) // lowercased trimmed name -> first slot to use it
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		key := strings.ToLower(strings.TrimSpace(fe.Name.String()))
		if key == "" {
			return // empty names handled by DIR-NAME-NOT-EMPTY
		}
		if prev, ok := seen[key]; ok {
			findings = append(findings, Finding{
				RuleID:   "CROSS-NO-DUPLICATE-NAMES",
				Severity: SeverityInconsistency,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("filename %q duplicates slot %d", key, prev),
				Citation: "samdos/src/c.s:1196-1219",
			})
			return
		}
		seen[key] = slot
	})
	return findings
}

// ----- CROSS-DIRECTORY-AREA-UNUSED -----
func init() {
	Register(Rule{
		ID:          "CROSS-DIRECTORY-AREA-UNUSED",
		Severity:    SeverityStructural,
		Description: "no chain link in any used file references a directory-area sector (tracks 0-3 of side 0)",
		Citation:    "samfile.go:984-987",
		Check:       checkCrossDirectoryAreaUnused,
	})
}

func checkCrossDirectoryAreaUnused(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		result := walkChain(ctx.Disk, fe.FirstSector)
		for _, st := range result.Steps {
			// Directory area is side 0 cylinders 0..3 only (Tech Manual L4340-4343);
			// side 1 cylinders 0..3 (tracks 0x80..0x83) are valid data sectors.
			if st.Sector.Track < 4 {
				s := st.Sector
				findings = append(findings, Finding{
					RuleID:   "CROSS-DIRECTORY-AREA-UNUSED",
					Severity: SeverityStructural,
					Location: SectorLocation(slot, fe.Name.String(), &s, -1),
					Message:  fmt.Sprintf("chain visits %v which is in the directory area", s),
					Citation: "samfile.go:984-987",
				})
				return // one finding per slot
			}
		}
	})
	return findings
}
