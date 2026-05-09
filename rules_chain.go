package samfile

// §3 Sector-chain rules + §15 CHAIN-SECTOR-COUNT-MINIMAL (catalog
// docs/disk-validity-rules.md §3 + §15). Rules in this file walk
// each used file's sector chain and check link integrity, cycle
// freedom, and consistency with the SectorAddressMap. They apply
// to all dialects.
//
// walkChain (private) is shared with rules_cross.go via the same
// package; it is the single canonical chain-walker for Phase 3
// rules so per-rule walking stays simple.
