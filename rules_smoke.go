package samfile

import "fmt"

// DISK-NOT-EMPTY is the Phase 1 smoke-test rule: it fires on a disk
// with zero occupied directory entries. The "real" rule catalog has
// dozens of these; this is the one we wire up end-to-end in Phase 1
// to prove the registry + Verify plumbing works. Severity is
// inconsistency rather than fatal because an empty disk is unusual
// but technically valid SAM-format output.
func init() {
	Register(Rule{
		ID:            "DISK-NOT-EMPTY",
		Severity:      SeverityInconsistency,
		Dialects:      nil, // all dialects
		Description:   "disk has at least one occupied directory entry",
		Citation:      "docs/disk-validity-rules.md",
		Check:         checkDiskNotEmpty,
		Applicability: &RuleApplicability{Scope: DiskScope},
	})
}

func checkDiskNotEmpty(ctx *CheckContext) []Finding {
	used := ctx.Journal.UsedFileEntries()
	if len(used) > 0 {
		return nil
	}
	return []Finding{{
		RuleID:   "DISK-NOT-EMPTY",
		Severity: SeverityInconsistency,
		Location: DiskWideLocation(),
		Message:  fmt.Sprintf("disk has 0 occupied directory entries (all %d slots are free)", len(ctx.Journal.FreeFileEntries())),
		Citation: "docs/disk-validity-rules.md",
	}}
}
