package samfile

import "strings"

// DetectDialect inspects di and returns the most likely dialect that
// wrote the disk. The heuristic combines independent signals (boot
// file at T4S1, MGTFlags bit patterns across used slots) and returns
// DialectUnknown when those signals are silent or contradict each
// other.
//
// Detection is deliberately conservative: when the result is
// DialectUnknown, Verify only runs rules tagged AllDialects, which is
// always safe. Pass --dialect=NAME on the CLI to override the result
// when the heuristic gets it wrong.
//
// Signals consulted (each returns its own DialectUnknown when it has
// no opinion; see bootFileDialect, mgtFlagsDialect):
//
//   - Boot file name and type — the slot whose FirstSector is (4, 1)
//     identifies the DOS that wrote the disk: "samdos2" → SAMDOS-2,
//     "masterdos"/"masterdos2" → MasterDOS, "samdos" or a type-3 file
//     → SAMDOS-1.
//   - MGTFlags across used slots — bits outside {0x00, 0x20} signal
//     MasterDOS (catalog: DIALECT-MASTERDOS-MGTFLAGS).
//
// Other dialect-distinguishing signals (BASIC SAVARS-NVARS gap,
// FileTypeInfo conventions) are deferred to later phases when the
// file-type rules land.
func DetectDialect(di *DiskImage) Dialect {
	dj := di.DiskJournal()
	opinions := []Dialect{
		bootFileDialect(dj),
		mgtFlagsDialect(dj),
	}
	var picked Dialect = DialectUnknown
	for _, o := range opinions {
		if o == DialectUnknown {
			continue
		}
		if picked == DialectUnknown {
			picked = o
			continue
		}
		if picked != o {
			return DialectUnknown // conflict → conservative
		}
	}
	return picked
}

// bootFileDialect examines the slot whose FirstSector is (track 4,
// sector 1) — the sector ROM BOOTEX reads to &8000 (see catalog
// BOOT-OWNER-AT-T4S1). The slot's filename (trimmed, lowercased) and
// masked Type are matched against the canonical DOS bootstraps:
//
//   - "samdos2" or "samdos 2"      → DialectSAMDOS2
//   - "masterdos" or "masterdos2"  → DialectMasterDOS
//   - "samdos" (no trailing 2), or masked Type == 3
//                                  → DialectSAMDOS1
//
// Anything else (including no used slot at T4S1) returns
// DialectUnknown — the signal abstains rather than guesses.
func bootFileDialect(dj *DiskJournal) Dialect {
	for _, fe := range dj {
		if fe == nil {
			continue
		}
		if fe.FirstSector == nil ||
			fe.FirstSector.Track != 4 ||
			fe.FirstSector.Sector != 1 {
			continue
		}
		// We have the boot slot. Check type-3 first: SAMDOS-1's
		// auto-include header (samdos/src/b.s:14-22) sets this type on
		// the bootstrap itself. Type 3 is otherwise unused by later DOSes,
		// and restricting the check to the boot slot keeps it unambiguous.
		// Note: FileEntry.Used() treats unknown types as not-used, so we
		// must check type-3 before the Used() guard.
		if uint8(fe.Type)&0x1F == 3 {
			return DialectSAMDOS1
		}
		if !fe.Used() {
			return DialectUnknown
		}
		name := strings.ToLower(strings.TrimSpace(fe.Name.String()))
		switch name {
		case "samdos2", "samdos 2":
			return DialectSAMDOS2
		case "masterdos", "masterdos2":
			return DialectMasterDOS
		case "samdos":
			return DialectSAMDOS1
		}
		return DialectUnknown
	}
	return DialectUnknown
}

// mgtFlagsDialect scans every used slot's MGTFlags. A value outside
// the SAMDOS-2 set {0x00, 0x20, 0xFF} signals MasterDOS (catalog:
// DIALECT-MASTERDOS-MGTFLAGS). Returns DialectUnknown when every used
// slot's MGTFlags is in the SAMDOS-2 set, including the trivial
// empty-disk case.
//
// The three SAMDOS-2 values cover all known writer conventions:
//
//   - 0xFF — what ROM SAMDOS-2 SAVE writes by default. The 14-byte
//     0xFF-fill loop at HDCLP2 (rom-disasm L22076-22080) starts at
//     dir offset 0xDC, which is the MGTFlags byte. Real-SAVE CODE
//     files therefore retain 0xFF; observed on the M0 boot disk's
//     slot-4 OUT file.
//   - 0x20 — what ROM SAMDOS-2 BASIC-SAVE overwrites it to after the
//     HDCLP2 fill (catalog BASIC-MGTFLAGS-20). The "MGT use only"
//     marker bit; Tech Manual L4369.
//   - 0x00 — what samfile.AddCodeFile leaves it at (Go struct
//     zero-init). Not what real ROM SAVE produces, but the
//     convention every other samfile-built CODE file follows.
//
// MasterDOS sets per-file attribute bits in MGTFlags to track its
// own metadata. The exact bit semantics are undocumented in our
// corpus (catalog §13 DIALECT-MASTERDOS-MGTFLAGS), so we treat
// anything outside the SAMDOS-2 set as a MasterDOS signal rather
// than checking specific bit patterns.
func mgtFlagsDialect(dj *DiskJournal) Dialect {
	for _, fe := range dj {
		if fe == nil || !fe.Used() {
			continue
		}
		switch fe.MGTFlags {
		case 0x00, 0x20, 0xff:
			// SAMDOS-2 set — silent.
		default:
			return DialectMasterDOS
		}
	}
	return DialectUnknown
}
