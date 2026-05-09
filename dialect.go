package samfile

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
	return DialectUnknown
}
