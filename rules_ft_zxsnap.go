// rules_ft_zxsnap.go
package samfile

// §10 ZX snapshot rules (catalog docs/disk-validity-rules.md §10).
// Rules in this file check FT_ZX_SNAPSHOT (5) invariants: 48 KiB
// body length and 0x4000 load address. The catalog tags these as
// SAMDOS-2 specific (the constants live in SAMDOS source); we run
// them on all dialects because the ZX snapshot format is itself
// dialect-agnostic.
