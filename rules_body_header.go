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
func init() {
	Register(Rule{
		ID:          "BODY-TYPE-MATCHES-DIR",
		Severity:    SeverityInconsistency,
		Description: "body header Type byte equals directory-entry Type (attribute bits masked)",
		Citation:    "samdos/src/c.s:1395-1408",
		Check:       checkBodyTypeMatchesDir,
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
			"BODY-TYPE-MATCHES-DIR", SeverityInconsistency, "samdos/src/c.s:1395-1408", "type",
			slot, fe.Name.String(),
			uint8(fe.Type)&0x1F, hdr[0],
		)...)
	})
	return findings
}

// ----- BODY-EXEC-DIV16K-MATCHES-DIR -----
func init() {
	Register(Rule{
		ID:          "BODY-EXEC-DIV16K-MATCHES-DIR",
		Severity:    SeverityStructural,
		Description: "body header ExecutionAddressDiv16K (byte 5) equals dir-entry ExecutionAddressDiv16K",
		Citation:    "rom-disasm:22471-22484",
		Check:       checkBodyExecDiv16KMatchesDir,
	})
}

func checkBodyExecDiv16KMatchesDir(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		findings = append(findings, bodyDirMirrorFinding(
			"BODY-EXEC-DIV16K-MATCHES-DIR", SeverityStructural, "rom-disasm:22471-22484", "ExecutionAddressDiv16K",
			slot, fe.Name.String(),
			fe.ExecutionAddressDiv16K, hdr[5],
		)...)
	})
	return findings
}

// ----- BODY-EXEC-MOD16K-LO-MATCHES-DIR -----
func init() {
	Register(Rule{
		ID:          "BODY-EXEC-MOD16K-LO-MATCHES-DIR",
		Severity:    SeverityInconsistency,
		Description: "body header ExecutionAddressMod16KLo (byte 6) equals low byte of dir-entry ExecutionAddressMod16K",
		Citation:    "rom-disasm:22472",
		Check:       checkBodyExecMod16KLoMatchesDir,
	})
}

func checkBodyExecMod16KLoMatchesDir(ctx *CheckContext) []Finding {
	var findings []Finding
	forEachUsedSlot(ctx, func(slot int, fe *FileEntry) {
		hdr, err := bodyHeaderRaw(ctx.Disk, fe)
		if err != nil {
			return
		}
		findings = append(findings, bodyDirMirrorFinding(
			"BODY-EXEC-MOD16K-LO-MATCHES-DIR", SeverityInconsistency, "rom-disasm:22472", "ExecutionAddressMod16KLo",
			slot, fe.Name.String(),
			uint8(fe.ExecutionAddressMod16K&0xFF), hdr[6],
		)...)
	})
	return findings
}

// ----- BODY-PAGES-MATCHES-DIR -----
func init() {
	Register(Rule{
		ID:          "BODY-PAGES-MATCHES-DIR",
		Severity:    SeverityInconsistency,
		Description: "body header Pages (byte 7) equals dir-entry Pages",
		Citation:    "samdos/src/c.s:1376-1379",
		Check:       checkBodyPagesMatchesDir,
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
			"BODY-PAGES-MATCHES-DIR", SeverityInconsistency, "samdos/src/c.s:1376-1379", "Pages",
			slot, fe.Name.String(),
			fe.Pages, hdr[7],
		)...)
	})
	return findings
}

// ----- BODY-STARTPAGE-MATCHES-DIR -----
func init() {
	Register(Rule{
		ID:          "BODY-STARTPAGE-MATCHES-DIR",
		Severity:    SeverityInconsistency,
		Description: "body header StartPage (byte 8) equals dir-entry StartAddressPage",
		Citation:    "samdos/src/c.s:1376-1379",
		Check:       checkBodyStartPageMatchesDir,
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
			"BODY-STARTPAGE-MATCHES-DIR", SeverityInconsistency, "samdos/src/c.s:1376-1379", "StartAddressPage",
			slot, fe.Name.String(),
			fe.StartAddressPage, hdr[8],
		)...)
	})
	return findings
}

// ----- BODY-LENGTHMOD16K-MATCHES-DIR -----
func init() {
	Register(Rule{
		ID:          "BODY-LENGTHMOD16K-MATCHES-DIR",
		Severity:    SeverityInconsistency,
		Description: "body header LengthMod16K (bytes 1-2 LE) equals dir-entry LengthMod16K",
		Citation:    "samdos/src/c.s:1376-1379",
		Check:       checkBodyLengthMod16KMatchesDir,
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
				Severity: SeverityInconsistency,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("body LengthMod16K = 0x%04x but dir says 0x%04x", actual, fe.LengthMod16K),
				Citation: "samdos/src/c.s:1376-1379",
			})
		}
	})
	return findings
}

// ----- BODY-PAGEOFFSET-MATCHES-DIR -----
func init() {
	Register(Rule{
		ID:          "BODY-PAGEOFFSET-MATCHES-DIR",
		Severity:    SeverityInconsistency,
		Description: "body header PageOffset (bytes 3-4 LE) equals dir-entry StartAddressPageOffset",
		Citation:    "samdos/src/c.s:1376-1379",
		Check:       checkBodyPageOffsetMatchesDir,
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
				Severity: SeverityInconsistency,
				Location: SlotLocation(slot, fe.Name.String()),
				Message:  fmt.Sprintf("body PageOffset = 0x%04x but dir says 0x%04x", actual, fe.StartAddressPageOffset),
				Citation: "samdos/src/c.s:1376-1379",
			})
		}
	})
	return findings
}

// ----- BODY-MIRROR-AT-DIR-D3-DB -----
func init() {
	Register(Rule{
		ID:          "BODY-MIRROR-AT-DIR-D3-DB",
		Severity:    SeverityInconsistency,
		Description: "dir bytes 0xD3..0xDB mirror body header bytes 0..8 (and dir byte 0xD2 is 0)",
		Citation:    "samdos/src/f.s:462-471",
		Check:       checkBodyMirrorAtDirD3DB,
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
				Severity: SeverityInconsistency,
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
					Severity: SeverityInconsistency,
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
		ID:          "BODY-PAGEOFFSET-8000H-FORM",
		Severity:    SeverityCosmetic,
		Description: "body-header PageOffset has bit 15 set (8000H-form convention)",
		Citation:    "sam-coupe_tech-man_v3-0.txt:3037-3052",
		Check:       checkBodyPageOffset8000HForm,
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
		ID:          "BODY-PAGE-LE-31",
		Severity:    SeverityStructural,
		Description: "body-header StartPage's low 5 bits encode an on-disk page index (0..30)",
		Citation:    "samfile.go:248-249",
		Check:       checkBodyPageLE31,
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
		ID:          "BODY-BYTES-5-6-CANONICAL-FF",
		Severity:    SeverityCosmetic,
		Description: "when body[5]==0xFF (no auto-exec), real SAVE writes body[6]==0xFF too",
		Citation:    "rom-disasm:22076-22080",
		Check:       checkBodyBytes56CanonicalFF,
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
