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
