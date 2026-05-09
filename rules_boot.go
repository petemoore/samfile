package samfile

import "fmt"

// §11 Boot-file rules (catalog docs/disk-validity-rules.md §11).
// Rules in this file check that a disk's boot file (the slot whose
// FirstSector is at track 4, sector 1) carries the bytes ROM BOOTEX
// expects: a "BOOT" signature at offset 256-259 of T4S1, and
// plausible Z80 code at body offset 0 (sector offset 9). They apply
// to all dialects.

// bootSlot returns the (slot index, FileEntry) of the disk's boot file
// — the used slot whose FirstSector is (track 4, sector 1). Returns
// found=false when no used slot owns T4S1 (a non-bootable disk).
// BOOT-OWNER-AT-T4S1 produces a finding in that case; the other two
// §11 rules silently skip (their checks are conditional on a boot
// file existing).
func bootSlot(dj *DiskJournal) (slot int, fe *FileEntry, found bool) {
	for idx, e := range dj {
		if e == nil || !e.Used() {
			continue
		}
		if e.FirstSector != nil && e.FirstSector.Track == 4 && e.FirstSector.Sector == 1 {
			return idx, e, true
		}
	}
	return -1, nil, false
}

// Ensure fmt is used (will be used by rules added in Task 2).
var _ = fmt.Sprintf
