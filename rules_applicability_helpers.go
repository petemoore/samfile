package samfile

// Shared Applicability.Filter helpers used by rule registrations.
// usedSlot accepts every SlotSubject whose dir entry has a non-zero
// type byte (the standard "used file" set). typedSlot returns a
// filter that further narrows to one or more specific file types.

// usedSlot returns true when s is a SlotSubject whose FileEntry has
// a non-zero type byte — i.e. the slot is considered occupied by
// SAMDOS's free-slot test (type byte == 0).
func usedSlot(ctx *CheckContext, s Subject) bool {
	ss, ok := s.(*SlotSubject)
	if !ok || ss.FileEntry == nil {
		return false
	}
	return ss.FileEntry.Type != 0
}

// typedSlot returns a Filter that accepts a SlotSubject only when
// its FileEntry's Type matches one of the provided file types. Used
// by rules that apply to a specific file-type family (e.g. all
// FT_CODE rules, the array-pair rule, etc.).
func typedSlot(types ...FileType) func(*CheckContext, Subject) bool {
	return func(ctx *CheckContext, s Subject) bool {
		ss, ok := s.(*SlotSubject)
		if !ok || ss.FileEntry == nil {
			return false
		}
		for _, t := range types {
			if ss.FileEntry.Type == t {
				return true
			}
		}
		return false
	}
}
