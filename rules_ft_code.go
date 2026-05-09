package samfile

// §6 FT_CODE rules (catalog docs/disk-validity-rules.md §6).
// Rules in this file check FT_CODE-specific invariants: the file's
// load address is above ROM, the loaded region fits in SAM's 512 KiB
// address space, the execution address (if not opted out) lies within
// the loaded region, and dir-entry FileTypeInfo is unused (cosmetic).
// Each Check function filters on fe.Type == FT_CODE at the top.
