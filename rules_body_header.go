package samfile

import "fmt"

// §5 Body-header rules (catalog docs/disk-validity-rules.md §5).
// Rules in this file compare the 9-byte body header at each used
// file's first sector against the parsed directory-entry fields it
// is supposed to mirror, plus a handful of byte-level format
// invariants that don't have a dir-entry counterpart. They apply to
// all dialects.
//
// bodyHeaderRaw (private) reads the 9-byte header once per rule
// invocation; bodyDirMirrorFinding (private) standardises the
// "body field X mismatches dir field Y" Finding shape.

// bodyHeaderRaw reads the 9 leading bytes of fe's body — the on-disk
// FileHeader bytes (Type, LengthMod16K-lo, LengthMod16K-hi, PageOffset-lo,
// PageOffset-hi, ExecutionAddressDiv16K, ExecutionAddressMod16KLo, Pages,
// StartPage). Returns an error if fe.FirstSector is unreadable; rules
// should treat that as "no finding" because §1 / §2 rules already report
// the underlying first-sector problem.
//
// This is a thin convenience over SectorData(fe.FirstSector); it does
// not allocate beyond the returned array.
func bodyHeaderRaw(di *DiskImage, fe *FileEntry) ([9]byte, error) {
	var hdr [9]byte
	sd, err := di.SectorData(fe.FirstSector)
	if err != nil {
		return hdr, err
	}
	copy(hdr[:], sd[:9])
	return hdr, nil
}

// bodyDirMirrorFinding compares one expected (dir-derived) value to one
// actual (body-derived) value and returns either nil or a single Finding
// pinpointing the mismatch. The same shape is used by every §5 byte-mirror
// rule — RuleID, Severity, Citation, and a human-readable fieldName feed
// into a uniform message format.
func bodyDirMirrorFinding(
	ruleID string, sev Severity, citation, fieldName string,
	slot int, name string,
	expected, actual uint8,
) []Finding {
	if expected == actual {
		return nil
	}
	return []Finding{{
		RuleID:   ruleID,
		Severity: sev,
		Location: SlotLocation(slot, name),
		Message:  fmt.Sprintf("body %s = 0x%02x but dir says 0x%02x", fieldName, actual, expected),
		Citation: citation,
	}}
}

// ----- BODY-TYPE-MATCHES-DIR -----
// Iteration 2 DEMOTE (inconsistency → cosmetic). The body-header
// mirror cluster is save-time-only: SAMDOS LOAD reads the dir entry
// (gtfle c.s:1376-1379 → uifa → hconr h.s:336-361 → ROM HDL/HDR
// via txhed h.s:38-56), and the body's first 9 bytes are read by
// ldhd (f.s:494-497) using lbyt (c.s:557-570) — but lbyt returns
// each byte in A without storing it, so ldhd just skips past the
// 9-byte body header so payload reads start at body byte 9. The
// body header never feeds into ROM's view, so a body↔dir mismatch
// has zero load-time consequence.
func init() {
	Register(Rule{
		ID:            "BODY-TYPE-MATCHES-DIR",
		Severity:      SeverityCosmetic,
		Description:   "body header Type byte equals directory-entry Type (attribute bits masked)",
		Citation:      "samdos/src/c.s:1395-1408",
		Check:         checkBodyTypeMatchesDir,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: usedSlot},
	})
}

func checkBodyTypeMatchesDir(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return // §1 rules already report the underlying first-sector problem
		}
		findings = append(findings, bodyDirMirrorFinding(
			"BODY-TYPE-MATCHES-DIR", SeverityCosmetic, "samdos/src/c.s:1395-1408", "type",
			slot, fe.Name.String(),
			uint8(fe.Type)&0x1F, hdr[0],
		)...)
	})
	return findings
}

// ----- BODY-EXEC-DIV16K-MATCHES-DIR -----
// Scoped to FT_CODE only. Dir bytes 0xF2-0xF4 are repurposed for
// non-CODE file types: FT_SAM_BASIC stores the auto-RUN line number
// there (0x00 marker + 16-bit line), while samfile's CreateHeader
// (samfile.go:921-927) emits 0xFF for body[5] on every non-CODE
// file regardless of dir contents. The "mirror" only holds for
// CODE files where both sides encode the same exec-address.
//
// Iteration 1 SCOPE: skip when body[5] == 0xFF. Per the auto-exec
// gate at rom-disasm:22471-22484, auto-exec is disabled at the body
// level when body[5] is 0xFF — bytes 5-6 are the canonical
// "defer to dir" pattern ROM SAVE writes (catalog
// §BODY-BYTES-5-6-CANONICAL-FF).
//
// Iteration 2 DEMOTE: structural → cosmetic. The body-mirror cluster
// (BODY-MIRROR-AT-DIR-D3-DB and friends, commit 2b8aa1a) was
// demoted to cosmetic after the SAMDOS source-chain analysis
// (samdos/src/f.s:494-497 ldhd → c.s:557-570 lbyt → h.s:336-361
// hconr → h.s:38-56 txhed) proved that body bytes 0..8 are read-
// and-discarded on LOAD; the auto-exec gate reads HDL+HDN+6 which
// is dir-derived (gtfle fills hd001.. from dir 0xD3-0xDB; hconr
// reloads from uifa+* which is dir-side). Body[5] never enters
// ROM's view. The "mismatch can cause unwanted auto-exec" framing
// in the earlier catalog text was wrong; this rule should match
// the rest of the body-mirror cluster. Corpus iter-2: 1897 fires
// across 164 disks, 727 of which are the `body=0x00, dir=0xFF`
// pattern — load-time consequence: none.
func init() {
	Register(Rule{
		ID:            "BODY-EXEC-DIV16K-MATCHES-DIR",
		Severity:      SeverityCosmetic,
		Description:   "FT_CODE body-header ExecutionAddressDiv16K (byte 5) equals dir-entry ExecutionAddressDiv16K (cosmetic: body bytes 0..8 are unused on LOAD; mismatch has zero load-time consequence)",
		Citation:      "rom-disasm:22471-22484; samdos/src/f.s:494-497",
		Check:         checkBodyExecDiv16KMatchesDir,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: typedSlot(FT_CODE)},
	})
}

func checkBodyExecDiv16KMatchesDir(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_CODE {
			return
		}
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		// body[5] == 0xFF is the canonical "defer to dir" pattern that
		// real ROM SAVE writes; dir is authoritative. Not a real
		// mismatch regardless of dir's value — suppress.
		if hdr[5] == 0xFF {
			return
		}
		findings = append(findings, bodyDirMirrorFinding(
			"BODY-EXEC-DIV16K-MATCHES-DIR", SeverityCosmetic, "rom-disasm:22471-22484", "ExecutionAddressDiv16K",
			slot, fe.Name.String(),
			fe.ExecutionAddressDiv16K, hdr[5],
		)...)
	})
	return findings
}

// ----- BODY-EXEC-MOD16K-LO-MATCHES-DIR -----
// Scoped to FT_CODE only — same reasoning as BODY-EXEC-DIV16K-MATCHES-DIR.
// For non-CODE files, dir's ExecutionAddressMod16K holds the auto-RUN
// line (BASIC) or other type-specific data while body[6] is always
// 0xFF (CreateHeader's non-FT_CODE default).
//
// Iteration 1 SCOPE: skip when body[5] == 0xFF. When auto-exec is
// disabled at the body level (body[5]==0xFF — the canonical
// "defer to dir" pattern per BODY-EXEC-DIV16K-MATCHES-DIR), the
// body's byte 6 (lo of mod16K) is meaningless and any value is
// allowed.
//
// Iteration 2 DEMOTE: inconsistency → cosmetic. Same SAMDOS source-
// chain analysis as BODY-EXEC-DIV16K-MATCHES-DIR (samdos/src/f.s:
// 494-497 ldhd reads body bytes 0..8 via lbyt and discards them;
// hconr/txhed feed ROM HDL/HDR from the dir entry). Body[6] never
// enters ROM's view. Corpus iter-2: 1873 fires across 164 disks
// (parallels the div16k pair). Load-time consequence: none.
func init() {
	Register(Rule{
		ID:            "BODY-EXEC-MOD16K-LO-MATCHES-DIR",
		Severity:      SeverityCosmetic,
		Description:   "FT_CODE body-header ExecutionAddressMod16KLo (byte 6) equals low byte of dir-entry ExecutionAddressMod16K (cosmetic: body bytes 0..8 are unused on LOAD; mismatch has zero load-time consequence)",
		Citation:      "rom-disasm:22472; samdos/src/f.s:494-497",
		Check:         checkBodyExecMod16KLoMatchesDir,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: typedSlot(FT_CODE)},
	})
}

func checkBodyExecMod16KLoMatchesDir(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		if fe.Type != FT_CODE {
			return
		}
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		// When body[5]==0xFF, auto-exec is disabled at the body level
		// and byte 6 (lo of mod16K) is meaningless — its value is not
		// a mismatch with dir, just an undefined byte. Suppress.
		if hdr[5] == 0xFF {
			return
		}
		findings = append(findings, bodyDirMirrorFinding(
			"BODY-EXEC-MOD16K-LO-MATCHES-DIR", SeverityCosmetic, "rom-disasm:22472", "ExecutionAddressMod16KLo",
			slot, fe.Name.String(),
			uint8(fe.ExecutionAddressMod16K&0xFF), hdr[6],
		)...)
	})
	return findings
}

// ----- BODY-PAGES-MATCHES-DIR -----
// Iteration 2 DEMOTE (inconsistency → cosmetic). Same family as
// BODY-TYPE-MATCHES-DIR: body byte 7 is skipped (not stored) by
// ldhd (f.s:494-497) on LOAD; the dir mirror is what feeds the
// load path via gtfle/hconr/txhed.
func init() {
	Register(Rule{
		ID:            "BODY-PAGES-MATCHES-DIR",
		Severity:      SeverityCosmetic,
		Description:   "body header Pages (byte 7) equals dir-entry Pages",
		Citation:      "samdos/src/c.s:1376-1379",
		Check:         checkBodyPagesMatchesDir,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: usedSlot},
	})
}

func checkBodyPagesMatchesDir(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		findings = append(findings, bodyDirMirrorFinding(
			"BODY-PAGES-MATCHES-DIR", SeverityCosmetic, "samdos/src/c.s:1376-1379", "Pages",
			slot, fe.Name.String(),
			fe.Pages, hdr[7],
		)...)
	})
	return findings
}

// ----- BODY-STARTPAGE-MATCHES-DIR -----
// Iteration 2 DEMOTE (inconsistency → cosmetic). Same family as
// BODY-TYPE-MATCHES-DIR: body byte 8 is skipped (not stored) by
// ldhd (f.s:494-497) on LOAD. ROM's StartPage view is loaded from
// the dir entry via uifa+31 → page1 through hconr (h.s:346-347).
func init() {
	Register(Rule{
		ID:            "BODY-STARTPAGE-MATCHES-DIR",
		Severity:      SeverityCosmetic,
		Description:   "body header StartPage (byte 8) equals dir-entry StartAddressPage",
		Citation:      "samdos/src/c.s:1376-1379",
		Check:         checkBodyStartPageMatchesDir,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: usedSlot},
	})
}

func checkBodyStartPageMatchesDir(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		findings = append(findings, bodyDirMirrorFinding(
			"BODY-STARTPAGE-MATCHES-DIR", SeverityCosmetic, "samdos/src/c.s:1376-1379", "StartAddressPage",
			slot, fe.Name.String(),
			fe.StartAddressPage, hdr[8],
		)...)
	})
	return findings
}

// ----- BODY-LENGTHMOD16K-MATCHES-DIR -----
// Iteration 2 DEMOTE (inconsistency → cosmetic). Same family as
// BODY-TYPE-MATCHES-DIR: body bytes 1-2 are skipped (not stored)
// by ldhd (f.s:494-497) on LOAD; the dir-side LengthMod16K is what
// the load path uses.
func init() {
	Register(Rule{
		ID:            "BODY-LENGTHMOD16K-MATCHES-DIR",
		Severity:      SeverityCosmetic,
		Description:   "body header LengthMod16K (bytes 1-2 LE) equals dir-entry LengthMod16K",
		Citation:      "samdos/src/c.s:1376-1379",
		Check:         checkBodyLengthMod16KMatchesDir,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: usedSlot},
	})
}

func checkBodyLengthMod16KMatchesDir(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		actual := uint16(hdr[1]) | uint16(hdr[2])<<8
		if actual != fe.LengthMod16K {
			findings = append(findings, Finding{
				RuleID:   "BODY-LENGTHMOD16K-MATCHES-DIR",
				Severity: SeverityCosmetic,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("body LengthMod16K = 0x%04x but dir says 0x%04x", actual, fe.LengthMod16K),
				Citation: "samdos/src/c.s:1376-1379",
			})
		}
	})
	return findings
}

// ----- BODY-PAGEOFFSET-MATCHES-DIR -----
// Iteration 2 DEMOTE (inconsistency → cosmetic). Same family as
// BODY-TYPE-MATCHES-DIR: body bytes 3-4 are skipped (not stored)
// by ldhd (f.s:494-497) on LOAD. Dir 0xED-0xEE is what hconr
// (h.s:349-350) populates hd0d1 from on the load path.
func init() {
	Register(Rule{
		ID:            "BODY-PAGEOFFSET-MATCHES-DIR",
		Severity:      SeverityCosmetic,
		Description:   "body header PageOffset (bytes 3-4 LE) equals dir-entry StartAddressPageOffset",
		Citation:      "samdos/src/c.s:1376-1379",
		Check:         checkBodyPageOffsetMatchesDir,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: usedSlot},
	})
}

func checkBodyPageOffsetMatchesDir(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		actual := uint16(hdr[3]) | uint16(hdr[4])<<8
		if actual != fe.StartAddressPageOffset {
			findings = append(findings, Finding{
				RuleID:   "BODY-PAGEOFFSET-MATCHES-DIR",
				Severity: SeverityCosmetic,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("body PageOffset = 0x%04x but dir says 0x%04x", actual, fe.StartAddressPageOffset),
				Citation: "samdos/src/c.s:1376-1379",
			})
		}
	})
	return findings
}

// ----- BODY-MIRROR-AT-DIR-D3-DB -----
// Iteration 2 DEMOTE (inconsistency → cosmetic). svhd (f.s:462-471)
// writes the same 9 bytes to dir+0xD3 AND to the body header on
// SAVE, so the mirror IS canonical SAVE output. But on LOAD the
// body header bytes 0..8 are unused: ldhd (f.s:494-497) calls lbyt
// (c.s:557-570) which returns each body byte in A *without storing
// it anywhere* — ldhd simply advances the read pointer past the
// 9-byte body header so subsequent ldblk reads start at body byte
// 9 (the payload). Everything ROM sees on LOAD is dir-derived:
// gtfle (c.s:1376-1379) fills the in-RAM cache and uifa from the
// dir entry, hconr (h.s:336-361) reloads hd001/page1/hd0d1/pges1/
// hd0b1 from uifa+* (dir-derived) — not from the body — and txhed
// (h.s:38-56) transmits 48 bytes from difa (the dir-entry buffer)
// into ROM's HDL/HDR area. So a body↔dir mismatch in bytes 0..8
// has zero load-time consequence; only the dir side feeds into the
// load path. 93% of samdos2-written disks omit the mirror, which
// is irreconcilable with "buggy writer" framing.
func init() {
	Register(Rule{
		ID:            "BODY-MIRROR-AT-DIR-D3-DB",
		Severity:      SeverityCosmetic,
		Description:   "dir bytes 0xD3..0xDB mirror body header bytes 0..8 (and dir byte 0xD2 is 0)",
		Citation:      "samdos/src/f.s:462-471",
		Check:         checkBodyMirrorAtDirD3DB,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: usedSlot},
	})
}

func checkBodyMirrorAtDirD3DB(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		// dir byte 0xD2 == 0 (MGTFutureAndPast[0])
		if fe.MGTFutureAndPast[0] != 0 {
			findings = append(findings, Finding{
				RuleID:   "BODY-MIRROR-AT-DIR-D3-DB",
				Severity: SeverityCosmetic,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("dir byte 0xD2 (MGTFutureAndPast[0]) = 0x%02x but should be 0", fe.MGTFutureAndPast[0]),
				Citation: "samdos/src/f.s:462-471",
			})
		}
		// dir bytes 0xD3..0xDB (MGTFutureAndPast[1..9]) mirror body bytes 0..8
		for i := 0; i < 9; i++ {
			if fe.MGTFutureAndPast[1+i] != hdr[i] {
				findings = append(findings, Finding{
					RuleID:   "BODY-MIRROR-AT-DIR-D3-DB",
					Severity: SeverityCosmetic,
					Location: SlotLocation(slot, fe.Name.String()),
					Message: fmt.Sprintf("dir byte 0x%02x (MGTFutureAndPast[%d]) = 0x%02x but body byte %d = 0x%02x",
						0xD3+i, 1+i, fe.MGTFutureAndPast[1+i], i, hdr[i]),
					Citation: "samdos/src/f.s:462-471",
				})
				return // one finding per slot is enough; the disagreement is the signal
			}
		}
	})
	return findings
}

// ----- BODY-PAGEOFFSET-8000H-FORM -----
// Real ROM SAVE writes PageOffset with bit 15 set ("8000H form" / REL
// PAGE FORM convention, Tech Manual L3037-3052). Both samfile.Start()
// and the ROM PDPSR2 decoder mask & 0x3FFF before use, so a bit-15-
// clear value still parses — but it deviates from convention and is
// a useful corpus-validation signal.
func init() {
	Register(Rule{
		ID:            "BODY-PAGEOFFSET-8000H-FORM",
		Severity:      SeverityCosmetic,
		Description:   "body-header PageOffset has bit 15 set (8000H-form convention)",
		Citation:      "sam-coupe_tech-man_v3-0.txt:3037-3052",
		Check:         checkBodyPageOffset8000HForm,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: usedSlot},
	})
}

func checkBodyPageOffset8000HForm(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		pageOffset := uint16(hdr[3]) | uint16(hdr[4])<<8
		// A zero offset is a legitimate "page-aligned" load; only warn
		// when there are bits in the low 14 but bit 15 is clear.
		if pageOffset != 0 && pageOffset&0x8000 == 0 {
			findings = append(findings, Finding{
				RuleID:   "BODY-PAGEOFFSET-8000H-FORM",
				Severity: SeverityCosmetic,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("body PageOffset = 0x%04x is missing bit 15 (8000H-form convention)", pageOffset),
				Citation: "sam-coupe_tech-man_v3-0.txt:3037-3052",
			})
		}
	})
	return findings
}

// ----- BODY-PAGE-LE-31 -----
// body[8] & 0x1F is the page index BEFORE samfile's +1 shift in
// FileHeader.Start(). Index 31 (raw) gives a +1 of 32, which lands
// the load address at 0x80000 (off-disk pseudo-page used as a
// marker, e.g. by SAMBASIC). Real on-disk load addresses use 0..30.
func init() {
	Register(Rule{
		ID:            "BODY-PAGE-LE-31",
		Severity:      SeverityStructural,
		Description:   "body-header StartPage's low 5 bits encode an on-disk page index (0..30)",
		Citation:      "samfile.go:248-249",
		Check:         checkBodyPageLE31,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: usedSlot},
	})
}

func checkBodyPageLE31(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		page := hdr[8] & 0x1F
		if page > 30 {
			findings = append(findings, Finding{
				RuleID:   "BODY-PAGE-LE-31",
				Severity: SeverityStructural,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("body StartPage low-5 bits = %d (>30); +1 shift lands above on-disk pages", page),
				Citation: "samfile.go:248-249",
			})
		}
	})
	return findings
}

// ----- BODY-BYTES-5-6-CANONICAL-FF -----
// When ExecutionAddressDiv16K (body[5]) is 0xFF (the "no auto-exec"
// marker), real ROM SAVE writes 0xFF to body[6] as well — both bytes
// 0xFF are the canonical "no auto-exec" pair. samfile's writer also
// emits 0xFF for body[6] in that case via CreateHeader (samfile.go:921-937),
// so this rule fires only on hand-edited or legacy files where body[5]
// is 0xFF but body[6] is something else. Both forms parse identically;
// the rule documents the {FF, FF} convention.
func init() {
	Register(Rule{
		ID:            "BODY-BYTES-5-6-CANONICAL-FF",
		Severity:      SeverityCosmetic,
		Description:   "when body[5]==0xFF (no auto-exec), real SAVE writes body[6]==0xFF too",
		Citation:      "rom-disasm:22076-22080",
		Check:         checkBodyBytes56CanonicalFF,
		Applicability: &RuleApplicability{Scope: SlotScope, Filter: usedSlot},
	})
}

func checkBodyBytes56CanonicalFF(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		if hdr[5] == 0xFF && hdr[6] != 0xFF {
			findings = append(findings, Finding{
				RuleID:   "BODY-BYTES-5-6-CANONICAL-FF",
				Severity: SeverityCosmetic,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("body[5]=0xFF (no auto-exec) but body[6]=0x%02x; canonical SAVE writes 0xFF here too", hdr[6]),
				Citation: "rom-disasm:22076-22080",
			})
		}
	})
	return findings
}
