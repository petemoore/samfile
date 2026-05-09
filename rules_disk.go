package samfile

// §1 Disk-level rules (catalog docs/disk-validity-rules.md §1).
// Rules in this file check that every track and sector reference on
// disk lies within the documented MGT geometry. They apply to all
// dialects.
