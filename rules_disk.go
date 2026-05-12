package samfile

import "fmt"

// §1 Disk-level rules (catalog docs/disk-validity-rules.md §1).
// Rules in this file check that every track and sector reference on
// disk lies within the documented MGT geometry. They apply to all
// dialects.

// trackSectorRefs returns every (track, sector) link reachable from
// ctx — first-sector references from used dir entries plus the
// next-link bytes (510-511) of every sector in every used file's
// chain. Bounded by the disk's 1560-sector capacity per chain so a
// cyclic or truncated chain cannot hang the iteration. Used by
// DISK-DIRECTORY-TRACKS, DISK-TRACK-SIDE-ENCODING, DISK-SECTOR-RANGE.
//
// Each returned ref carries enough context for a Finding's Location
// (slot index, slot name, the sector itself, and the byte offset
// within the sector where the link byte lives — 0 for a first-sector
// reference, 510 for a chain link's track byte, 511 for sector byte).
//
// Errors from SectorData (only fire on out-of-range raw track values
// that bypass the dir entry's parse path) are silently ignored —
// DISK-TRACK-SIDE-ENCODING will catch them via the dir entry's own
// first-sector reference.
type sectorRef struct {
	Slot         int
	Filename     string
	Sector       Sector // copy (not pointer) so the value is independent of any pool
	ByteOffset   int    // 0 (first-sector) or 510 (chain link track) or 511 (chain link sector)
	IsTerminator bool   // true when this ref is the (0, 0) chain terminator — skip range checks
}

func trackSectorRefs(ctx *CheckContext) []sectorRef {
	var refs []sectorRef
	for _, slot := range ctx.Journal.UsedFileEntries() {
		fe := ctx.Journal[slot]
		name := fe.Name.String()
		// First-sector reference from the dir entry.
		refs = append(refs, sectorRef{Slot: slot, Filename: name, Sector: *fe.FirstSector, ByteOffset: 0})
		// Walk the chain. Bound by 1560 (disk capacity) to defend
		// against cycles / missing terminators; CHAIN-NO-CYCLE will
		// also catch those.
		cur := fe.FirstSector
		for steps := 0; steps < 1560; steps++ {
			sd, err := ctx.Disk.SectorData(cur)
			if err != nil {
				break
			}
			fp := sd.FilePart()
			nextSec := *fp.NextSector
			isTerm := nextSec.Track == 0 && nextSec.Sector == 0
			refs = append(refs,
				sectorRef{Slot: slot, Filename: name, Sector: nextSec, ByteOffset: 510, IsTerminator: isTerm},
				sectorRef{Slot: slot, Filename: name, Sector: nextSec, ByteOffset: 511, IsTerminator: isTerm},
			)
			if isTerm {
				break
			}
			cur = fp.NextSector
		}
	}
	return refs
}

func init() {
	Register(Rule{
		ID:          "DISK-DIRECTORY-TRACKS",
		Severity:    SeverityStructural,
		Description: "no file references a sector in the directory area (tracks 0-3 of side 0)",
		Citation:    "sam-coupe_tech-man_v3-0.txt:4340-4343",
		Check:       checkDiskDirectoryTracks,
	})
}

func checkDiskDirectoryTracks(ctx *CheckContext) []Finding {
	var findings []Finding
	for _, ref := range trackSectorRefs(ctx) {
		if ref.IsTerminator {
			continue // (0, 0) terminator is allowed even though Track=0 is in [0..3]
		}
		// Directory area is side 0 cylinders 0..3 only (Tech Manual L4340-4343);
		// side 1 cylinders 0..3 (tracks 0x80..0x83) are valid data sectors. The
		// (0,0) chain terminator is already filtered out above.
		if ref.Sector.Track < 4 {
			findings = append(findings, Finding{
				RuleID:   "DISK-DIRECTORY-TRACKS",
				Severity: SeverityStructural,
				Location: SectorLocation(ref.Slot, ref.Filename, &ref.Sector, ref.ByteOffset),
				Message:  fmt.Sprintf("track 0x%02x references the directory area (tracks 0-3 of side 0)", ref.Sector.Track),
				Citation: "sam-coupe_tech-man_v3-0.txt:4340-4343",
			})
		}
	}
	return findings
}

func init() {
	Register(Rule{
		ID:          "DISK-TRACK-SIDE-ENCODING",
		Severity:    SeverityFatal,
		Description: "every track byte references a physical cylinder 0-79 on side 0 or side 1",
		Citation:    "samfile.go:393-394",
		Check:       checkDiskTrackSideEncoding,
	})
}

func checkDiskTrackSideEncoding(ctx *CheckContext) []Finding {
	var findings []Finding
	for _, ref := range trackSectorRefs(ctx) {
		if ref.ByteOffset == 511 {
			continue // sector-number byte, not the track byte
		}
		if ref.IsTerminator {
			continue
		}
		t := ref.Sector.Track
		if (t >= 80 && t < 128) || t >= 208 {
			findings = append(findings, Finding{
				RuleID:   "DISK-TRACK-SIDE-ENCODING",
				Severity: SeverityFatal,
				Location: SectorLocation(ref.Slot, ref.Filename, &ref.Sector, ref.ByteOffset),
				Message:  fmt.Sprintf("track 0x%02x is in the invalid range (valid: 0x00-0x4F or 0x80-0xCF)", t),
				Citation: "samfile.go:393-394",
			})
		}
	}
	return findings
}

func init() {
	Register(Rule{
		ID:          "DISK-SECTOR-RANGE",
		Severity:    SeverityFatal,
		Description: "every sector number is in range 1-10 (or 0 for the chain terminator)",
		Citation:    "samfile.go:389-392",
		Check:       checkDiskSectorRange,
	})
}

func checkDiskSectorRange(ctx *CheckContext) []Finding {
	var findings []Finding
	for _, ref := range trackSectorRefs(ctx) {
		if ref.ByteOffset == 510 {
			continue // track byte, not the sector byte
		}
		if ref.IsTerminator {
			continue
		}
		s := ref.Sector.Sector
		if s < 1 || s > 10 {
			findings = append(findings, Finding{
				RuleID:   "DISK-SECTOR-RANGE",
				Severity: SeverityFatal,
				Location: SectorLocation(ref.Slot, ref.Filename, &ref.Sector, ref.ByteOffset),
				Message:  fmt.Sprintf("sector 0x%02x is out of range (valid: 1-10)", s),
				Citation: "samfile.go:389-392",
			})
		}
	}
	return findings
}
