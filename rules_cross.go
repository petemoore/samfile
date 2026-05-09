package samfile

// §4 Cross-entry consistency rules (catalog docs/disk-validity-rules.md
// §4). Rules in this file compare data across multiple directory
// slots: shared sectors, duplicate names, references into the
// directory area. They apply to all dialects.
