package samfile

import (
	"bytes"
	"testing"
)

// TestFileTypeInfoPageFormDecoding pins down the SAM Coupé "PAGEFORM"
// encoding of the three SAM-BASIC FileTypeInfo length fields, and the
// four section sizes derived from them.
//
// The on-disk encoding stores three cumulative offsets (end of program,
// end of numeric variables, end of gap). Each is a 19-bit length in
// PAGEFORM: byte 0 = page count (16384 each), bytes 1-2 = LE address in
// section C with low 14 bits significant. See ROM disasm RDTHREE
// (sam-coupe_rom-v3.0_annotated-disassembly.txt:7654-7659) and PAGEFORM
// (sam-coupe_rom-v3.0_annotated-disassembly.txt:7578-7589).
//
// The four section sizes (program, numeric vars, gap, string/array
// vars) must sum to the file's total Length — that's the structural
// invariant we check on every fixture below.
//
// Falsification fixtures: bytes captured directly from real SAMDOS-
// written BASIC files on disks downloaded from
// ftp.nvg.ntnu.no/pub/sam-coupe/disks/utils/. Tech Manual L4370-4382
// describes the fields as cumulative lengths.
func TestFileTypeInfoPageFormDecoding(t *testing.T) {
	cases := []struct {
		name         string
		info         [11]byte
		pages        uint8
		lengthMod16K uint16
		wantProgLen  uint32
		wantNVarsSz  uint32
		wantGapSz    uint32
		wantSAVSz    uint32
	}{
		{
			name:         "Auto Font (FontLoader.dsk, no string/array vars)",
			info:         [11]byte{0x01, 0xc3, 0x8c, 0x01, 0x0d, 0x8e, 0x01, 0x1f, 0x8f, 0x20, 0xff},
			pages:        1,
			lengthMod16K: 0x0f1f, // total Length = 20255
			wantProgLen:  19651,
			wantNVarsSz:  330,
			wantGapSz:    274,
			wantSAVSz:    0,
		},
		{
			name:         "Shredder (FileShredderv1.2.dsk, no string/array vars)",
			info:         [11]byte{0x00, 0xbf, 0x9b, 0x00, 0x3a, 0x9c, 0x00, 0x1b, 0x9e},
			pages:        0,
			lengthMod16K: 0x1e1b, // total Length = 7707
			wantProgLen:  7103,
			wantNVarsSz:  123,
			wantGapSz:    481,
			wantSAVSz:    0,
		},
		{
			name:         "AUTOCOMMS (CommsLoader.dsk, has string/array vars)",
			info:         [11]byte{0x00, 0x22, 0x9f, 0x00, 0x30, 0xa0, 0x00, 0x7e, 0xa1},
			pages:        0,
			lengthMod16K: 0x229d, // total Length = 8861
			wantProgLen:  7970,
			wantNVarsSz:  270,
			wantGapSz:    334,
			wantSAVSz:    287,
		},
		{
			name:         "empty: page=0 addr=0x8000 throughout, total=0",
			info:         [11]byte{0x00, 0x00, 0x80, 0x00, 0x00, 0x80, 0x00, 0x00, 0x80},
			pages:        0,
			lengthMod16K: 0,
			wantProgLen:  0,
			wantNVarsSz:  0,
			wantGapSz:    0,
			wantSAVSz:    0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fe := &FileEntry{
				FileTypeInfo: c.info,
				Pages:        c.pages,
				LengthMod16K: c.lengthMod16K,
			}
			if got := fe.ProgramLength(); got != c.wantProgLen {
				t.Errorf("ProgramLength = %d; want %d", got, c.wantProgLen)
			}
			if got := fe.NumericVariablesSize(); got != c.wantNVarsSz {
				t.Errorf("NumericVariablesSize = %d; want %d", got, c.wantNVarsSz)
			}
			if got := fe.GapSize(); got != c.wantGapSz {
				t.Errorf("GapSize = %d; want %d", got, c.wantGapSz)
			}
			if got := fe.StringArrayVariablesSize(); got != c.wantSAVSz {
				t.Errorf("StringArrayVariablesSize = %d; want %d", got, c.wantSAVSz)
			}
			sum := fe.ProgramLength() + fe.NumericVariablesSize() +
				fe.GapSize() + fe.StringArrayVariablesSize()
			if sum != fe.Length() {
				t.Errorf("section sizes sum to %d; want fe.Length() = %d", sum, fe.Length())
			}
		})
	}
}


// TestAddCodeFile8000HFormPageOffset asserts that AddCodeFile stores the
// page offset in 8000H-BFFFH form, so disks built with samfile load at
// the correct address when read by SAMDOS.
//
// Tech Manual v3.0 L4326-4329 (file-header section) and L4388-4392
// (directory-entry section) both specify the encoding:
//
//   - Byte 8 (file header) / 236 (directory entry): "starting page
//     number ... AND this with 1FH to get the page number in the
//     range 0 to 31".
//   - Bytes 3-4 (file header) / 237-238 (directory entry): "PAGE
//     OFFSET (8000-BFFFH)".
//   - Decode: start = page * 16384 + raw_offset - 0x4000.
//
// Empirically confirmed against three real SAMDOS-written disks
// (CommsLoader.dsk, FileShredderv1.2.dsk, FontLoader.dsk): 38 of 38
// directory entries have bit 15 set in bytes 0xed/0xee. Without the
// fix, samfile-written disks load ~16K below the intended address
// when read by SAMDOS — even though samfile's own reader (which masks
// `& 0x3fff`) reads them back correctly.
//
// This test runs samfile-written stored bytes through the Tech Manual
// decode formula directly, bypassing samfile's reader-side masking.
func TestAddCodeFile8000HFormPageOffset(t *testing.T) {
	cases := []struct {
		name        string
		loadAddress uint32
	}{
		{"first RAM byte (0x4000)", 0x4000},
		{"section C boundary (0xC000)", 0xC000},
		{"FontLoader 'Font Code' equivalent (82000)", 82000},
		{"end of last RAM page (0x7FFFC)", 0x7FFFC},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			di := &DiskImage{}
			if err := di.AddCodeFile("F", []byte("test"), c.loadAddress, 0); err != nil {
				t.Fatalf("AddCodeFile: %v", err)
			}
			var fe *FileEntry
			for _, e := range di.DiskJournal() {
				if e.Used() && e.Name.String() == "F" {
					fe = e
					break
				}
			}
			if fe == nil {
				t.Fatal("F entry not found in disk journal")
			}
			rawOffset := fe.StartAddressPageOffset
			if rawOffset&0x8000 == 0 {
				t.Errorf("raw page offset 0x%04x has bit 15 clear; "+
					"Tech Manual L4390 specifies range 8000-BFFFH",
					rawOffset)
			}
			if rawOffset >= 0xC000 {
				t.Errorf("raw page offset 0x%04x outside 8000-BFFFH range",
					rawOffset)
			}
			techManualStart := uint32(fe.StartAddressPage&0x1f)*0x4000 +
				uint32(rawOffset) - 0x4000
			if techManualStart != c.loadAddress {
				t.Errorf("Tech Manual decode of stored bytes "+
					"(page=0x%02x, offset=0x%04x) = 0x%05x; "+
					"load address was 0x%05x",
					fe.StartAddressPage, rawOffset,
					techManualStart, c.loadAddress)
			}
		})
	}
}

// TestAddCodeFileExecutionAddressDiv16KConvention pins the byte-level
// encoding of the execution-address triplet at directory-entry bytes
// 0xf2-0xf4.
//
// Tech Manual v3.0 L4396 lists "EXECUTION ADDRESS" at directory bytes
// 242-244 (UIFA 37-39) for CODE files, but is silent on the exact
// byte-level encoding. The encoding therefore has to be derived
// empirically from real SAMDOS-written disks.
//
// Three independent disks all carry COMET.COD with start address
// 36921 (= 0x9039) and stored execution-address bytes
// 0xf2 = 0x02, 0xf3-0xf4 = 0x39 0x90 (= 0x9039 LE):
//
//   - CometAssembler1.8EdwinBlink.dsk
//   - comet18(1)/Comet18.dsk
//   - GoodSamC2/comet.dsk
//
// Decoding byte 0xf2 directly as `addr / 16384` and bytes 0xf3-0xf4
// (LE) as `(addr mod 16384) | 0x8000` reproduces 0x9039 exactly.
// Decoding byte 0xf2 with a -1 offset (the convention StartPage byte
// 0xec uses, per L4388) gives 53305 — past the end of the 12231-byte
// COMET body. So the convention for execution-address bytes is:
//
//   - byte 0xf2 = (executionAddress / 16384), NO -1 offset
//   - bytes 0xf3-0xf4 = (executionAddress mod 16384) | 0x8000
//     (8000H-form, same as StartAddressPageOffset)
//
// Without the fix, AddCodeFile stored byte 0xf2 with `(addr>>14) - 1`,
// inherited by copy-paste from the StartPage writer. SAMDOS auto-RUN
// of a samfile-written file would then jump 16K below the intended
// entry point.
func TestAddCodeFileExecutionAddressDiv16KConvention(t *testing.T) {
	cases := []struct {
		name             string
		loadAddress      uint32
		executionAddress uint32
		wantDiv16K       uint8
		wantMod16K       uint16
	}{
		{"first RAM byte (0x4000)", 0x4000, 0x4000, 1, 0x8000},
		{"typical user code (0x6000)", 0x6000, 0x6000, 1, 0xA000},
		{"COMET.COD (0x9039)", 0x9039, 0x9039, 2, 0x9039},
		{"section C boundary (0xC000)", 0xC000, 0xC000, 3, 0x8000},
		{"end of last RAM page (0x7FFFC)", 0x7FFFC, 0x7FFFC, 31, 0xBFFC},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			di := &DiskImage{}
			if err := di.AddCodeFile("F", []byte("test"), c.loadAddress, c.executionAddress); err != nil {
				t.Fatalf("AddCodeFile: %v", err)
			}
			var fe *FileEntry
			for _, e := range di.DiskJournal() {
				if e.Used() && e.Name.String() == "F" {
					fe = e
					break
				}
			}
			if fe == nil {
				t.Fatal("F entry not found in disk journal")
			}
			if fe.ExecutionAddressDiv16K != c.wantDiv16K {
				t.Errorf("byte 0xf2 (Div16K) = 0x%02x; want 0x%02x. "+
					"Real SAMDOS-written disks (3 COMET.COD samples) store "+
					"`addr/16384` directly; the StartPage-style `-1` offset "+
					"does NOT apply to the execution-address byte.",
					fe.ExecutionAddressDiv16K, c.wantDiv16K)
			}
			if fe.ExecutionAddressMod16K != c.wantMod16K {
				t.Errorf("bytes 0xf3-0xf4 (Mod16K LE) = 0x%04x; want 0x%04x "+
					"(8000H-form, like StartAddressPageOffset)",
					fe.ExecutionAddressMod16K, c.wantMod16K)
			}
		})
	}
}

// TestAddCodeFileExecutionAddressRoundTrip is the user-facing
// round-trip assertion for the execution-address writer/reader pair:
// what AddCodeFile stores, ExecutionAddress() must read back unchanged.
//
// Without the fix to AddCodeFile, the writer applies `-1` to the page
// byte (inherited from the StartPage writer), but the reader (corrected
// in commit e64f5d5 "Execution Address Page off by one") does not add
// `+1` back. Round-trip is then off by 16K.
func TestAddCodeFileExecutionAddressRoundTrip(t *testing.T) {
	cases := []struct {
		name             string
		loadAddress      uint32
		executionAddress uint32
	}{
		{"exec at load (0x4000)", 0x4000, 0x4000},
		{"exec at load (0x6000)", 0x6000, 0x6000},
		{"exec at load (COMET 0x9039)", 0x9039, 0x9039},
		{"exec at load (0xC000)", 0xC000, 0xC000},
		{"exec at load (end of RAM 0x7FFFC)", 0x7FFFC, 0x7FFFC},
		{"exec mid-body (load 0x6000, exec 0x6002)", 0x6000, 0x6002},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			di := &DiskImage{}
			if err := di.AddCodeFile("F", []byte("test"), c.loadAddress, c.executionAddress); err != nil {
				t.Fatalf("AddCodeFile: %v", err)
			}
			var fe *FileEntry
			for _, e := range di.DiskJournal() {
				if e.Used() && e.Name.String() == "F" {
					fe = e
					break
				}
			}
			if fe == nil {
				t.Fatal("F entry not found in disk journal")
			}
			got := fe.ExecutionAddress()
			if got != c.executionAddress {
				t.Errorf("ExecutionAddress() = 0x%05x; want 0x%05x "+
					"(writer/reader round-trip mismatch)",
					got, c.executionAddress)
			}
		})
	}
}


// TestSAMMaskExhaustive iterates over every valid data sector — the 1560
// in the SAM domain (T4..T79 S1..S10 on side 0, T128..T207 S1..S10 on
// side 1) — and asserts SAMMask returns (bitOffset/8, 1<<(bitOffset%8)).
//
// Bit-numbering convention from Tech Manual v3.0 L4405-4413: "SAMDOS
// allocates 195 bytes to the sector address map, giving 1560 bits ...
// Bit 0 of the first byte is allocated to track 4 sector 1."
func TestSAMMaskExhaustive(t *testing.T) {
	bitOffset := 0
	check := func(track, sector uint8) {
		s := &Sector{Track: track, Sector: sector}
		offset, mask := s.SAMMask()
		wantOffset := uint8(bitOffset >> 3)
		wantMask := byte(1) << (bitOffset & 0x07)
		if offset != wantOffset || mask != wantMask {
			t.Errorf("SAMMask(T%dS%d) = (%d, 0x%02x); want (%d, 0x%02x) [bitOffset=%d]",
				track, sector, offset, mask, wantOffset, wantMask, bitOffset)
		}
		bitOffset++
	}
	for track := uint8(4); track <= 79; track++ {
		for sector := uint8(1); sector <= 10; sector++ {
			check(track, sector)
		}
	}
	if bitOffset != 760 {
		t.Fatalf("after side 0 enumeration, bitOffset = %d; want 760", bitOffset)
	}
	for track := uint8(128); track <= 207; track++ {
		for sector := uint8(1); sector <= 10; sector++ {
			check(track, sector)
		}
	}
	if bitOffset != 1560 {
		t.Fatalf("after side 1 enumeration, bitOffset = %d; want 1560", bitOffset)
	}
}

// TestSAMMaskBoundaries spot-checks the bit-offset boundaries derived
// from the Tech Manual's SAM-encoding text (L4405-4413).
func TestSAMMaskBoundaries(t *testing.T) {
	cases := []struct {
		name           string
		sector         Sector
		wantByteOffset uint8
		wantMask       uint8
	}{
		{"T4S1 first side-0 data sector", Sector{Track: 4, Sector: 1}, 0, 0x01},
		{"T4S10", Sector{Track: 4, Sector: 10}, 1, 0x02},
		{"T79S10 last side-0 data sector (bit 759)", Sector{Track: 79, Sector: 10}, 94, 0x80},
		{"T128S1 first side-1 data sector (bit 760)", Sector{Track: 128, Sector: 1}, 95, 0x01},
		{"T207S10 last side-1 data sector (bit 1559)", Sector{Track: 207, Sector: 10}, 194, 0x80},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			offset, mask := c.sector.SAMMask()
			if offset != c.wantByteOffset || mask != c.wantMask {
				t.Errorf("SAMMask = (%d, 0x%02x); want (%d, 0x%02x)",
					offset, mask, c.wantByteOffset, c.wantMask)
			}
		})
	}
}

// TestSAMMaskBugRegression pins the bit-mask values for sectors that the
// pre-fix expression `1 << bitOffset & 0x07` silently dropped.
//
// Go's operator precedence (Go spec § Operators: `<<` and `&` are both
// multiplicative, left-to-right) parses `1 << bitOffset & 0x07` as
// `(1 << bitOffset) & 0x07`, which is zero whenever bitOffset%8 >= 3.
// The corrected expression `1 << (bitOffset & 0x07)` evaluates to
// `1 << (bitOffset%8)` — always a single non-zero bit.
//
// Without the fix every sub-test below fails with mask == 0.
func TestSAMMaskBugRegression(t *testing.T) {
	cases := []struct {
		name      string
		sector    Sector
		bitOffset int
	}{
		{"T4S4 bitOffset 3", Sector{Track: 4, Sector: 4}, 3},
		{"T4S5 bitOffset 4", Sector{Track: 4, Sector: 5}, 4},
		{"T4S6 bitOffset 5", Sector{Track: 4, Sector: 6}, 5},
		{"T4S7 bitOffset 6", Sector{Track: 4, Sector: 7}, 6},
		{"T4S8 bitOffset 7", Sector{Track: 4, Sector: 8}, 7},
		{"T5S2 bitOffset 11", Sector{Track: 5, Sector: 2}, 11},
		{"T5S10 bitOffset 19", Sector{Track: 5, Sector: 10}, 19},
		{"T128S4 bitOffset 763", Sector{Track: 128, Sector: 4}, 763},
		{"T207S10 bitOffset 1559", Sector{Track: 207, Sector: 10}, 1559},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			offset, mask := c.sector.SAMMask()
			wantOffset := uint8(c.bitOffset >> 3)
			wantMask := byte(1) << (c.bitOffset & 0x07)
			if offset != wantOffset {
				t.Errorf("byte offset = %d; want %d", offset, wantOffset)
			}
			if mask == 0 {
				t.Fatalf("mask = 0; pre-fix `1 << bitOffset & 0x07` returns 0 for bitOffset %d "+
					"(bit-within-byte = %d). Correct expression: `1 << (bitOffset & 0x07)`.",
					c.bitOffset, c.bitOffset%8)
			}
			if mask != wantMask {
				t.Errorf("mask = 0x%02x; want 0x%02x", mask, wantMask)
			}
		})
	}
}

// TestAddFileMultiFileNoCorruption is the integration regression for the
// user-visible bite of the SAMMask bug.
//
// With the bug, AddCodeFile of file A records only sectors whose
// bitOffset%8 ∈ {0,1,2} in A's per-file SAM. The disk-wide free map
// (Tech Manual L4419-4420: bitwise OR of all per-file SAMs) therefore
// reports A's other sectors as free, and the next AddCodeFile silently
// allocates them — overwriting A's content. Reading A back then either
// fails (chain pointer dereferenced into the corrupted sector lands on
// a NextSector = (0,0)) or returns truncated/wrong bytes.
//
// 5120-byte body A spans 11 sectors (T4S1..T5S1, bitOffsets 0..10) — the
// last 8 cover bit-within-byte positions 3..7 of byte 0 and 0..2 of byte
// 1, exercising both the buggy and unaffected slots. 200-byte body B
// fits in a single sector and, under the bug, is allocated to T4S4
// (the first false-free slot at bitOffset 3) — directly inside A's chain.
func TestAddFileMultiFileNoCorruption(t *testing.T) {
	di := &DiskImage{}

	bodyA := make([]byte, 5120)
	for i := range bodyA {
		bodyA[i] = byte(i & 0xff)
	}
	bodyB := bytes.Repeat([]byte{0xBB}, 200)

	if err := di.AddCodeFile("A", bodyA, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile(A): %v", err)
	}
	if err := di.AddCodeFile("B", bodyB, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile(B): %v", err)
	}

	t.Run("file A read-back matches input", func(t *testing.T) {
		f, err := di.File("A")
		if err != nil {
			t.Fatalf("di.File(\"A\"): %v\n"+
				"This usually means B's AddCodeFile overwrote a sector inside A's chain — the SAMMask bug.",
				err)
		}
		if !bytes.Equal(f.Body, bodyA) {
			n := len(bodyA)
			if len(f.Body) < n {
				n = len(f.Body)
			}
			for i := 0; i < n; i++ {
				if f.Body[i] != bodyA[i] {
					t.Fatalf("file A body differs at offset %d: got 0x%02x; want 0x%02x "+
						"(B's AddCodeFile overwrote one of A's sectors — SAMMask bug)",
						i, f.Body[i], bodyA[i])
				}
			}
			t.Fatalf("file A body length differs: got %d; want %d", len(f.Body), len(bodyA))
		}
	})

	t.Run("file B read-back matches input", func(t *testing.T) {
		f, err := di.File("B")
		if err != nil {
			t.Fatalf("di.File(\"B\"): %v", err)
		}
		if !bytes.Equal(f.Body, bodyB) {
			t.Fatalf("file B body differs from input")
		}
	})
}
