package samfile

import "fmt"

// §5 Body-header rules (catalog docs/disk-validity-rules.md §5).
// Rules in this file compare the 9-byte body header at each used
// file's first sector against the parsed directory-entry fields it
// is supposed to mirror, plus a handful of byte-level format
// invariants that don't have a dir-entry counterpart. They apply to
// all dialects.
//
// bodyHeaderRaw (private) reads the 9-byte header once per rule
// invocation; bodyDirMirrorFinding (private) standardises the
// "body field X mismatches dir field Y" Finding shape.

// bodyHeaderRaw reads the 9 leading bytes of fe's body — the on-disk
// FileHeader bytes (Type, LengthMod16K-lo, LengthMod16K-hi, PageOffset-lo,
// PageOffset-hi, ExecutionAddressDiv16K, ExecutionAddressMod16KLo, Pages,
// StartPage). Returns an error if fe.FirstSector is unreadable; rules
// should treat that as "no finding" because §1 / §2 rules already report
// the underlying first-sector problem.
//
// This is a thin convenience over SectorData(fe.FirstSector); it does
// not allocate beyond the returned array.
func bodyHeaderRaw(di *DiskImage, fe *FileEntry) ([9]byte, error) {
	var hdr [9]byte
	sd, err := di.SectorData(fe.FirstSector)
	if err != nil {
		return hdr, err
	}
	copy(hdr[:], sd[:9])
	return hdr, nil
}

// bodyDirMirrorFinding compares one expected (dir-derived) value to one
// actual (body-derived) value and returns either nil or a single Finding
// pinpointing the mismatch. The same shape is used by every §5 byte-mirror
// rule — RuleID, Severity, Citation, and a human-readable fieldName feed
// into a uniform message format.
func bodyDirMirrorFinding(
	ruleID string, sev Severity, citation, fieldName string,
	slot int, name string,
	expected, actual uint8,
) []Finding {
	if expected == actual {
		return nil
	}
	return []Finding{{
		RuleID:   ruleID,
		Severity: sev,
		Location: SlotLocation(slot, name),
		Message:  fmt.Sprintf("body %s = 0x%02x but dir says 0x%02x", fieldName, actual, expected),
		Citation: citation,
	}}
}
