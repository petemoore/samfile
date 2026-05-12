package samfile

// SubjectScope identifies the kind of subject a Rule operates on.
type SubjectScope int

const (
	DiskScope SubjectScope = iota
	SlotScope
	ChainStepScope
)

func (s SubjectScope) String() string {
	switch s {
	case DiskScope:
		return "disk"
	case SlotScope:
		return "slot"
	case ChainStepScope:
		return "chain_step"
	}
	return "unknown"
}

// Subject is the universal interface for anything a rule can check —
// a whole disk, a directory slot, or a single sector in a file's
// chain. Ref() returns a stable string id; Attributes() returns the
// denormalised attribute snapshot recorded with every Check event.
type Subject interface {
	Ref() string
	Attributes() map[string]any
}
