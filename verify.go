package samfile

// Severity ranks findings by impact, lowest to highest.
type Severity int

const (
	SeverityCosmetic Severity = iota
	SeverityInconsistency
	SeverityStructural
	SeverityFatal
)

// String returns the lowercase canonical name of the severity,
// matching the names used by the disk-validity-rules.md catalog
// and the CLI's --severity flag.
func (s Severity) String() string {
	switch s {
	case SeverityCosmetic:
		return "cosmetic"
	case SeverityInconsistency:
		return "inconsistency"
	case SeverityStructural:
		return "structural"
	case SeverityFatal:
		return "fatal"
	}
	return "unknown"
}

// Dialect identifies which DOS produced the disk. Phase 1 only
// uses DialectUnknown (dialect detection lands in Phase 2); rules
// are scoped by their Dialects slice, with nil meaning all dialects.
type Dialect int

const (
	DialectUnknown Dialect = iota
	DialectSAMDOS1
	DialectSAMDOS2
	DialectMasterDOS
)

// String returns the lowercase canonical name of the dialect,
// matching the CLI's --dialect flag.
func (d Dialect) String() string {
	switch d {
	case DialectUnknown:
		return "unknown"
	case DialectSAMDOS1:
		return "samdos1"
	case DialectSAMDOS2:
		return "samdos2"
	case DialectMasterDOS:
		return "masterdos"
	}
	return "unknown"
}

// Location pinpoints a Finding on the disk. Construct one via the
// DiskWideLocation, SlotLocation, or SectorLocation factories — they
// set the "not applicable" sentinels correctly. The zero value of
// Location is NOT a valid disk-wide location (Slot=0 is a real slot).
type Location struct {
	Slot       int     // -1 if not applicable, else 0..79
	Sector     *Sector // nil if not applicable
	ByteOffset int     // -1 if not applicable, else byte offset within Sector
	Filename   string  // copied from Slot's directory entry when known, for messages
}

// DiskWideLocation returns a Location for findings that apply to the
// disk image as a whole (no specific slot or sector).
func DiskWideLocation() Location {
	return Location{Slot: -1, Sector: nil, ByteOffset: -1}
}

// SlotLocation returns a Location for findings tied to a specific
// directory slot but not a specific sector or byte.
func SlotLocation(slot int, filename string) Location {
	return Location{Slot: slot, Sector: nil, ByteOffset: -1, Filename: filename}
}

// SectorLocation returns a Location for findings tied to a specific
// byte within a specific sector of a specific file.
func SectorLocation(slot int, filename string, sector *Sector, byteOffset int) Location {
	return Location{Slot: slot, Sector: sector, ByteOffset: byteOffset, Filename: filename}
}

// IsDiskWide reports whether loc has no slot, sector, or byte set.
func (loc Location) IsDiskWide() bool {
	return loc.Slot == -1 && loc.Sector == nil && loc.ByteOffset == -1
}
