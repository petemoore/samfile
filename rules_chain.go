package samfile

import "fmt"

// §3 Sector-chain rules + §15 CHAIN-SECTOR-COUNT-MINIMAL (catalog
// docs/disk-validity-rules.md §3 + §15). Rules in this file walk
// each used file's sector chain and check link integrity, cycle
// freedom, and consistency with the SectorAddressMap. They apply
// to all dialects.
//
// walkChain (private) is shared with rules_cross.go via the same
// package; it is the single canonical chain-walker for Phase 3
// rules so per-rule walking stays simple.

// chainStep is one entry in a sector chain walk.
type chainStep struct {
	Sector Sector // the sector that was read at this step (copy, not pointer)
	Next   Sector // the (track, sector) link at bytes 510-511 of Sector
}

// chainWalkResult is the outcome of a walkChain call.
type chainWalkResult struct {
	Steps      []chainStep // in walk order
	Terminated bool        // true iff a (0, 0) link was encountered
	Cycle      *Sector     // first sector revisited, if any (nil = no cycle)
	Bailed     bool        // true iff the walk hit the 1560-step cap without terminating or cycling
}

// walkChain follows the link chain starting at first for at most 1560
// steps (the disk's data-sector capacity, an absolute upper bound on a
// terminated chain). It records each sector visited and the (track,
// sector) link found at its bytes 510-511. The walk halts on:
//
//   - a (0, 0) terminator (Terminated = true);
//   - a revisited sector (Cycle = &<first repeat>);
//   - hitting the 1560-step cap (Bailed = true).
//
// On read error from SectorData, the walk halts at the current Steps
// length without setting any flag. Callers can still use Steps to see
// what was reachable. (No findings are surfaced for read errors here;
// DISK-TRACK-SIDE-ENCODING / DIR-FIRST-SECTOR-VALID catch the underlying
// out-of-range track byte that triggers the SectorData error.)
func walkChain(di *DiskImage, first *Sector) chainWalkResult {
	var result chainWalkResult
	visited := make(map[Sector]bool)
	cur := *first
	for steps := 0; steps < 1560; steps++ {
		if visited[cur] {
			c := cur
			result.Cycle = &c
			return result
		}
		visited[cur] = true
		sd, err := di.SectorData(&cur)
		if err != nil {
			return result
		}
		fp := sd.FilePart()
		next := *fp.NextSector
		result.Steps = append(result.Steps, chainStep{Sector: cur, Next: next})
		if next.Track == 0 && next.Sector == 0 {
			result.Terminated = true
			return result
		}
		cur = next
	}
	result.Bailed = true
	return result
}

// ----- CHAIN-TERMINATOR-ZERO-ZERO -----
func init() {
	Register(Rule{
		ID:          "CHAIN-TERMINATOR-ZERO-ZERO",
		Severity:    SeverityStructural,
		Description: "each used file's sector chain ends with a (0, 0) link",
		Citation:    "samdos/src/b.s:104-110",
		Check:       checkChainTerminatorZeroZero,
	})
}

func checkChainTerminatorZeroZero(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		result := walkChain(ctx.Disk, fe.FirstSector)
		if !result.Terminated {
			lastSec := fe.FirstSector
			if n := len(result.Steps); n > 0 {
				s := result.Steps[n-1].Sector
				lastSec = &s
			}
			msg := "chain does not terminate"
			if result.Cycle != nil {
				msg = fmt.Sprintf("chain has a cycle (revisited %v)", result.Cycle)
			} else if result.Bailed {
				msg = "chain exceeds 1560 steps without (0, 0) link"
			}
			findings = append(findings, Finding{
				RuleID:   "CHAIN-TERMINATOR-ZERO-ZERO",
				Severity: SeverityStructural,
				Location: SectorLocation(slot, fe.Name.String(), lastSec, 510),
				Message:  msg,
				Citation: "samdos/src/b.s:104-110",
			})
		}
	})
	return findings
}

// ----- CHAIN-NO-CYCLE -----
func init() {
	Register(Rule{
		ID:          "CHAIN-NO-CYCLE",
		Severity:    SeverityStructural,
		Description: "each used file's sector chain has no revisited sectors",
		Citation:    "samfile.go:743-754",
		Check:       checkChainNoCycle,
	})
}

func checkChainNoCycle(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		result := walkChain(ctx.Disk, fe.FirstSector)
		if result.Cycle != nil {
			findings = append(findings, Finding{
				RuleID:   "CHAIN-NO-CYCLE",
				Severity: SeverityStructural,
				Location: SectorLocation(slot, fe.Name.String(), result.Cycle, 510),
				Message:  fmt.Sprintf("chain cycles: sector %v is revisited", result.Cycle),
				Citation: "samfile.go:743-754",
			})
		}
	})
	return findings
}

// ----- CHAIN-MATCHES-SAM -----
func init() {
	Register(Rule{
		ID:          "CHAIN-MATCHES-SAM",
		Severity:    SeverityStructural,
		Description: "the set of sectors walked by the chain equals the bits set in the SectorAddressMap",
		Citation:    "samdos/src/c.s:1306-1343",
		Check:       checkChainMatchesSAM,
	})
}

func checkChainMatchesSAM(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		result := walkChain(ctx.Disk, fe.FirstSector)
		walked := make(map[Sector]bool, len(result.Steps))
		for _, st := range result.Steps {
			walked[st.Sector] = true
		}
		mapSet := make(map[Sector]bool)
		for _, sec := range fe.SectorAddressMap.UsedSectors() {
			mapSet[*sec] = true
		}
		// Symmetric difference: any sector in one set but not the other.
		for s := range walked {
			if !mapSet[s] {
				findings = append(findings, Finding{
					RuleID:   "CHAIN-MATCHES-SAM",
					Severity: SeverityStructural,
					Location: SlotLocation(slot, fe.Name.String()),
					Message:  fmt.Sprintf("sector %v is visited by the chain but not set in the SectorAddressMap", s),
					Citation: "samdos/src/c.s:1306-1343",
				})
				return // one finding per slot is enough; the disagreement is the signal
			}
		}
		for s := range mapSet {
			if !walked[s] {
				findings = append(findings, Finding{
					RuleID:   "CHAIN-MATCHES-SAM",
					Severity: SeverityStructural,
					Location: SlotLocation(slot, fe.Name.String()),
					Message:  fmt.Sprintf("sector %v is set in the SectorAddressMap but not visited by the chain", s),
					Citation: "samdos/src/c.s:1306-1343",
				})
				return
			}
		}
	})
	return findings
}

// ----- CHAIN-SECTOR-COUNT-MINIMAL -----
func init() {
	Register(Rule{
		ID:          "CHAIN-SECTOR-COUNT-MINIMAL",
		Severity:    SeverityCosmetic,
		Description: "used file occupies exactly ceil((9 + body length) / 510) sectors (no padding sectors)",
		Citation:    "samfile.go:919",
		Check:       checkChainSectorCountMinimal,
	})
}

func checkChainSectorCountMinimal(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		bodyLen := int(fe.Length())
		required := uint16((bodyLen + 9 + 509) / 510)
		if fe.Sectors != required {
			findings = append(findings, Finding{
				RuleID:   "CHAIN-SECTOR-COUNT-MINIMAL",
				Severity: SeverityCosmetic,
				Location: SlotLocation(slot, fe.Name.String()),
				Message: fmt.Sprintf("file uses %d sectors but %d would suffice (bodyLen=%d)",
					fe.Sectors, required, bodyLen),
				Citation: "samfile.go:919",
			})
		}
	})
	return findings
}
