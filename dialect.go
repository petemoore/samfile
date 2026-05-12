package samfile

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// knownBodySHAs maps the sha256 of slot-0's file body — the exact bytes
// ROM BOOTEX would load into RAM — to the DOS dialect it represents.
// Entries are added only when the SHA can be matched byte-for-byte
// against an upstream reference binary or a corpus-confirmed family
// variant; we do not guess. Expand the map as more references land.
//
// Per-family directory tree at docs/dos-families/<sha>-<label>/ holds
// the binding chain: src/ has the upstream assembly, body.bin is the
// reference body, README.md records the verification.
var knownBodySHAs = map[string]Dialect{
	// ~/git/samdos/res/samdos2.reference.bin
	// = the binary the upstream samdos2 source assembles to
	// = source-of-record for the samdos2 family at family-head 9bc0fb4b…
	"3cca541beb3f9fe93402a770997945b2be852e69f278d2b176ba0bbc4fbb6077": DialectSAMDOS2,
	// ~/git/masterdos/res/MDOS23.bin
	// = the binary the upstream masterdos23.asm assembles to
	// = head of the masterdos-v2.3 family
	"152b811ed65b651df25e29f49e15340bec84ef3deebfba4eaa6cd76bfbb31fae": DialectMasterDOS,
}

// DetectDialect inspects di and returns the most likely dialect that
// wrote the disk. The heuristic combines independent signals (boot
// file at T4S1, MGTFlags bit patterns across used slots) and returns
// DialectUnknown when those signals are silent or contradict each
// other.
//
// Detection is deliberately conservative: when the result is
// DialectUnknown, Verify only runs rules tagged AllDialects, which is
// always safe. Pass --dialect=NAME on the CLI to override the result
// when the heuristic gets it wrong.
//
// Signals consulted (each returns its own DialectUnknown when it has
// no opinion; see bodyShaDialect, bootFileDialect, mgtFlagsDialect):
//
//   - **Body-SHA** (top priority, authoritative when matched) — the
//     slot-0 file body's sha256 is looked up in knownBodySHAs. A hit
//     is conclusive: the DOS *code* at T4S1 is byte-for-byte one we
//     have a source binding for. Returned immediately — other signals
//     (slot filename, neighbouring MGTFlags) can only mislead about
//     the DOS that will actually run when the disk boots.
//   - **Boot file name and type** — the slot whose FirstSector is
//     (4, 1) identifies the DOS by the dir-entry's filename:
//     "samdos2" → SAMDOS-2, "masterdos"/"masterdos2" → MasterDOS,
//     "samdos" or a type-3 file → SAMDOS-1.
//   - **MGTFlags across used slots** — bits outside {0x00, 0x20, 0xFF}
//     signal MasterDOS (catalog: DIALECT-MASTERDOS-MGTFLAGS). A
//     per-slot signal that gets promoted to a disk-level claim;
//     fragile when one disk contains files written by mixed DOSes.
//
// Other dialect-distinguishing signals (BASIC SAVARS-NVARS gap,
// FileTypeInfo conventions) are deferred to later phases when the
// file-type rules land.
func DetectDialect(di *DiskImage) Dialect {
	// Body-SHA is authoritative when matched: it identifies the DOS
	// *code* at T4S1, not a dir-entry convention. Skip the heuristics
	// entirely on a hit. Necessary because real disks (Fredatives 3,
	// many magazines) have a SAMDOS-2 body but a non-canonical slot-0
	// filename like "OS" — bootFileDialect would abstain, MGTFlags
	// would over-claim MasterDOS based on neighbouring slots, and the
	// disk would end up mis-classified.
	if d := bodyShaDialect(di); d != DialectUnknown {
		return d
	}
	dj := di.DiskJournal()
	opinions := []Dialect{
		bootFileDialect(dj),
		mgtFlagsDialect(dj),
	}
	var picked Dialect = DialectUnknown
	for _, o := range opinions {
		if o == DialectUnknown {
			continue
		}
		if picked == DialectUnknown {
			picked = o
			continue
		}
		if picked != o {
			return DialectUnknown // conflict → conservative
		}
	}
	return picked
}

// bodyShaDialect hashes the DOS body at T4S1 — the bytes ROM BOOTEX
// actually loads into RAM and jumps to — and looks the result up in
// knownBodySHAs. Returns DialectUnknown when:
//
//   - the disk is not ROM-bootable (T4S1 fails the BOOT signature
//     check at offset 256..259, so the question "what DOS wrote it"
//     has no answer at this layer), or
//   - the chain walk hits an unreadable sector or loops, or
//   - the body header's declared length is bogus (out of range, or
//     exceeds what the walked chain produced), or
//   - the resulting sha256 has no entry in knownBodySHAs.
//
// The walk starts from T4S1 directly rather than from a dir entry —
// some real disks (FRED Magazine 13, "scrubbed" disks) have the DOS
// body installed at T4S1 but no surviving dir slot that owns it. ROM
// BOOTEX doesn't consult the directory; nor should we when asking
// what code lives at T4S1.
//
// We do **not** sha256(T4S1 alone) — that would conflate two SAMDOS
// builds whose first 512 bytes happen to match. The body is walked
// per the file-header contract (samdos/src/c.s:1376-1379), exactly
// like tools/audit/survey_dos.py does in Python, so the SHAs
// produced here match what docs/dos-families/ records.
func bodyShaDialect(di *DiskImage) Dialect {
	body, ok := t4s1ChainBody(di)
	if !ok {
		return DialectUnknown
	}
	sum := sha256.Sum256(body)
	if d, found := knownBodySHAs[hex.EncodeToString(sum[:])]; found {
		return d
	}
	return DialectUnknown
}

// t4s1ChainBody walks the sector chain starting at (track 4,
// sector 1) and returns the body bytes — exactly Length() bytes
// starting at chain offset 9, per the 9-byte file header — that
// ROM BOOTEX would copy to the body's load address.
//
// Returns ok=false when T4S1 is unreadable or fails the BOOT
// signature check (XOR with "BOOT" AND 0x5F == 0 at bytes
// 256..259), the chain has an unreadable sector or loops, the body
// header declares an impossible length, or fewer payload bytes are
// available than the header demands.
//
// Mirrors tools/audit/survey_dos.py's decode_slot0_loadexec.
func t4s1ChainBody(di *DiskImage) ([]byte, bool) {
	t4s1 := &Sector{Track: 4, Sector: 1}
	sd, err := di.SectorData(t4s1)
	if err != nil {
		return nil, false
	}
	// BOOT signature check at bytes 256..259, per ROM BTCK
	// (rom-disasm:20473-20598). If T4S1 fails this check the ROM
	// won't boot into the body, so it isn't "DOS code" for dialect
	// purposes regardless of what bytes happen to be there.
	boot := [4]byte{'B', 'O', 'O', 'T'}
	for i := 0; i < 4; i++ {
		if (sd[256+i]^boot[i])&0x5F != 0 {
			return nil, false
		}
	}
	fp := sd.FilePart()
	raw := make([]byte, 0, 32*510)
	raw = append(raw, fp.Data[:510]...)
	visited := map[Sector]bool{*t4s1: true}
	// Decode body length from the 9-byte header we just read. The
	// header is byte 0..8 of the chain payload (raw[0..9]). Same
	// formula as FileHeader.Length(): low 14 bits of LengthMod16K
	// (raw[1]|raw[2]<<8) plus (pages = raw[7]) << 14.
	if len(raw) < 9 {
		return nil, false
	}
	lengthMod16K := uint32(raw[1]) | uint32(raw[2])<<8
	pages := uint32(raw[7])
	length := (lengthMod16K & 0x3FFF) | (pages << 14)
	if length == 0 {
		return nil, false
	}
	// Walk forward until we have at least 9 + length bytes, or the
	// chain terminates / errors / loops. Cap at a generous 1024
	// sectors (~512KB body) to bound runtime.
	for len(raw) < int(length)+9 {
		ns := *fp.NextSector
		if ns.Track == 0 && ns.Sector == 0 {
			return nil, false // chain ended before declared length
		}
		if visited[ns] {
			return nil, false // loop
		}
		if len(visited) > 1024 {
			return nil, false
		}
		visited[ns] = true
		sd, err = di.SectorData(&ns)
		if err != nil {
			return nil, false
		}
		fp = sd.FilePart()
		raw = append(raw, fp.Data[:510]...)
	}
	return raw[9 : 9+int(length)], true
}

// bootFileDialect examines the slot whose FirstSector is (track 4,
// sector 1) — the sector ROM BOOTEX reads to &8000 (see catalog
// BOOT-OWNER-AT-T4S1). The slot's filename (trimmed, lowercased) and
// masked Type are matched against the canonical DOS bootstraps:
//
//   - "samdos2" or "samdos 2"      → DialectSAMDOS2
//   - "masterdos" or "masterdos2"  → DialectMasterDOS
//   - "samdos" (no trailing 2), or masked Type == 3
//                                  → DialectSAMDOS1
//
// Anything else (including no used slot at T4S1) returns
// DialectUnknown — the signal abstains rather than guesses.
func bootFileDialect(dj *DiskJournal) Dialect {
	for _, fe := range dj {
		if fe == nil {
			continue
		}
		if fe.FirstSector == nil ||
			fe.FirstSector.Track != 4 ||
			fe.FirstSector.Sector != 1 {
			continue
		}
		// We have the boot slot. Check type-3 first: SAMDOS-1's
		// auto-include header (samdos/src/b.s:14-22) sets this type on
		// the bootstrap itself. Type 3 is otherwise unused by later DOSes,
		// and restricting the check to the boot slot keeps it unambiguous.
		// Note: FileEntry.Used() treats unknown types as not-used, so we
		// must check type-3 before the Used() guard.
		if uint8(fe.Type)&0x1F == 3 {
			return DialectSAMDOS1
		}
		if !fe.Used() {
			return DialectUnknown
		}
		name := strings.ToLower(strings.TrimSpace(fe.Name.String()))
		switch name {
		case "samdos2", "samdos 2":
			return DialectSAMDOS2
		case "masterdos", "masterdos2":
			return DialectMasterDOS
		case "samdos":
			return DialectSAMDOS1
		}
		return DialectUnknown
	}
	return DialectUnknown
}

// mgtFlagsDialect scans every used slot's MGTFlags. A value outside
// the SAMDOS-2 set {0x00, 0x20, 0xFF} signals MasterDOS (catalog:
// DIALECT-MASTERDOS-MGTFLAGS). Returns DialectUnknown when every used
// slot's MGTFlags is in the SAMDOS-2 set, including the trivial
// empty-disk case.
//
// The three SAMDOS-2 values cover all known writer conventions:
//
//   - 0xFF — what ROM SAMDOS-2 SAVE writes by default. The 14-byte
//     0xFF-fill loop at HDCLP2 (rom-disasm L22076-22080) starts at
//     dir offset 0xDC, which is the MGTFlags byte. Real-SAVE CODE
//     files therefore retain 0xFF; observed on the M0 boot disk's
//     slot-4 OUT file.
//   - 0x20 — what ROM SAMDOS-2 BASIC-SAVE overwrites it to after the
//     HDCLP2 fill (catalog BASIC-MGTFLAGS-20). The "MGT use only"
//     marker bit; Tech Manual L4369.
//   - 0x00 — what samfile.AddCodeFile leaves it at (Go struct
//     zero-init). Not what real ROM SAVE produces, but the
//     convention every other samfile-built CODE file follows.
//
// MasterDOS sets per-file attribute bits in MGTFlags to track its
// own metadata. The exact bit semantics are undocumented in our
// corpus (catalog §13 DIALECT-MASTERDOS-MGTFLAGS), so we treat
// anything outside the SAMDOS-2 set as a MasterDOS signal rather
// than checking specific bit patterns.
func mgtFlagsDialect(dj *DiskJournal) Dialect {
	for _, fe := range dj {
		if fe == nil || !fe.Used() {
			continue
		}
		switch fe.MGTFlags {
		case 0x00, 0x20, 0xff:
			// SAMDOS-2 set — silent.
		default:
			return DialectMasterDOS
		}
	}
	return DialectUnknown
}
