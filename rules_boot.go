package samfile

import "fmt"

// §11 Boot-file rules (catalog docs/disk-validity-rules.md §11).
// Rules in this file check that a disk's boot file (the slot whose
// FirstSector is at track 4, sector 1) carries the bytes ROM BOOTEX
// expects: a "BOOT" signature at offset 256-259 of T4S1, and
// plausible Z80 code at body offset 0 (sector offset 9). They apply
// to all dialects.

// bootSlot returns the (slot index, FileEntry) of the disk's boot file
// — the used slot whose FirstSector is (track 4, sector 1). Returns
// found=false when no used slot owns T4S1 (a non-bootable disk).
// BOOT-OWNER-AT-T4S1 produces a finding in that case; the other two
// §11 rules silently skip (their checks are conditional on a boot
// file existing).
func bootSlot(dj *DiskJournal) (slot int, fe *FileEntry, found bool) {
	for idx, e := range dj {
		if e == nil || !e.Used() {
			continue
		}
		if e.FirstSector != nil && e.FirstSector.Track == 4 && e.FirstSector.Sector == 1 {
			return idx, e, true
		}
	}
	return -1, nil, false
}

// ----- BOOT-OWNER-AT-T4S1 -----
// For an image to be bootable on real SAM hardware, some directory
// entry's FirstSector must be (track 4, sector 1) so that the ROM
// BOOTEX (rom-disasm:20473-20598) reads the right sector at &8000.
// Fires on a single disk-wide finding when no used slot owns T4S1.
//
// Note: data-only / archive disks legitimately have no boot file; this
// rule's "fatal" severity flags non-bootability, not corruption.
// Phase 7's corpus-validation pass may demote to cosmetic if archive
// disks dominate the corpus.
func init() {
	Register(Rule{
		ID:          "BOOT-OWNER-AT-T4S1",
		Severity:    SeverityFatal,
		Description: "some used directory entry has FirstSector (4, 1) so the disk is bootable on SAM hardware",
		Citation:    "rom-disasm:20473-20598",
		Check:       checkBootOwnerAtT4S1,
	})
}

func checkBootOwnerAtT4S1(ctx *CheckContext) []Finding {
	if _, _, found := bootSlot(ctx.Journal); found {
		return nil
	}
	return []Finding{{
		RuleID:   "BOOT-OWNER-AT-T4S1",
		Severity: SeverityFatal,
		Location: DiskWideLocation(),
		Message:  "no used slot has FirstSector (track 4, sector 1); disk is not bootable on real SAM hardware",
		Citation: "rom-disasm:20473-20598",
	}}
}

// ----- BOOT-SIGNATURE-AT-256 -----
// For ROM BOOTEX to dispatch to the loaded sector, bytes 256-259 of
// T4S1 must spell "BOOT" — case-insensitively, with bit 7 ignored
// (the ROM compares (disk_byte XOR expected_byte) AND 0x5F per
// rom-disasm:20582-20598). Only applies when a boot owner exists;
// BOOT-OWNER-AT-T4S1 reports the no-owner case separately.
func init() {
	Register(Rule{
		ID:          "BOOT-SIGNATURE-AT-256",
		Severity:    SeverityFatal,
		Description: "T4S1 bytes 256-259 spell \"BOOT\" (case-insensitive, bit 7 ignored)",
		Citation:    "rom-disasm:20582-20598",
		Check:       checkBootSignatureAt256,
	})
}

func checkBootSignatureAt256(ctx *CheckContext) []Finding {
	slot, fe, found := bootSlot(ctx.Journal)
	if !found {
		return nil // BOOT-OWNER-AT-T4S1 reports the underlying issue
	}
	sd, err := ctx.Disk.SectorData(fe.FirstSector)
	if err != nil {
		return nil // §1 rules report the underlying sector problem
	}
	// ROM compares with `XOR expected; AND 0x5F` — 0x5F = 0b01011111
	// masks bits 5 (case) and 7 (BASIC-keyword high bit) before the
	// zero check. So we apply the same mask here.
	expected := [4]byte{'B', 'O', 'O', 'T'}
	for i := 0; i < 4; i++ {
		if (sd[256+i]^expected[i])&0x5F != 0 {
			return []Finding{{
				RuleID:   "BOOT-SIGNATURE-AT-256",
				Severity: SeverityFatal,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("T4S1 boot signature mismatch at byte %d: got 0x%02x, expected 0x%02x (masked with 0x5F)", 256+i, sd[256+i], expected[i]),
				Citation: "rom-disasm:20582-20598",
			}}
		}
	}
	return nil
}

// ----- BOOT-ENTRY-POINT-AT-9 -----
// After signature match, ROM does JP 8009H. The sector buffer is at
// 0x8000-0x81FF, so 0x8009 is sector-buffer offset 9 = body offset 0
// (after the 9-byte body header). The byte at body offset 0 must
// therefore be valid Z80 code. We can't enforce "valid Z80 opcode"
// precisely from one byte, but 0x00 (NOP — unlikely as the first
// boot-code byte by design) and 0xFF (unwritten / no-code marker)
// are useful negative signals. Cosmetic per the catalog's test sketch.
func init() {
	Register(Rule{
		ID:          "BOOT-ENTRY-POINT-AT-9",
		Severity:    SeverityCosmetic,
		Description: "T4S1 body byte 0 (sector offset 9) is not 0x00 or 0xFF — a heuristic plausibility check for Z80 boot code",
		Citation:    "rom-disasm:20598",
		Check:       checkBootEntryPointAt9,
	})
}

func checkBootEntryPointAt9(ctx *CheckContext) []Finding {
	slot, fe, found := bootSlot(ctx.Journal)
	if !found {
		return nil
	}
	sd, err := ctx.Disk.SectorData(fe.FirstSector)
	if err != nil {
		return nil
	}
	b := sd[9]
	if b == 0x00 || b == 0xFF {
		return []Finding{{
			RuleID:   "BOOT-ENTRY-POINT-AT-9",
			Severity: SeverityCosmetic,
			Location: SlotLocation(slot, fe.Name.String()),
			Message:  fmt.Sprintf("T4S1 body byte 0 = 0x%02x (heuristic warn): 0x00 is NOP and 0xFF is the unwritten-byte default — both unusual for real boot-code entry", b),
			Citation: "rom-disasm:20598",
		}}
	}
	return nil
}
