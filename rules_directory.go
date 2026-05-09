package samfile

// §2 Directory-entry rules (catalog docs/disk-validity-rules.md §2).
// Rules in this file check internal consistency of each of the 80
// directory entries: type byte, filename padding, sector count vs
// chain length vs SectorAddressMap popcount. They apply to all
// dialects.
