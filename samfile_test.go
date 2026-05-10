package samfile

import (
	"bytes"
	"testing"
)

// TestFileTypeInfoPageFormDecoding pins down the SAM Coupé "PAGEFORM"
// encoding of the three SAM-BASIC FileTypeInfo length fields.
//
// The ROM's RDTHREE helper at sam-coupe_rom-v3.0_annotated-disassembly.txt:
// 7654-7659 reads three bytes from a directory entry into Z80 registers
// (C, E, D) — i.e. byte 0 → page register, bytes 1-2 → little-endian
// 16-bit address. PAGEFORM (sam-coupe_rom-v3.0_annotated-disassembly.txt:
// 7578-7589) then treats that as a 19-bit linear value: page * 16384 +
// (address - 0x8000), with bit 15 of the address always 1 (set by SCF;
// RR H) and bit 14 always 0. So the linear length is
//
//     page*16384 + (raw_addr & 0x3fff)
//
// Tech Manual L4370-4382 names the three fields ("program length
// excluding variables", "...plus numeric variables", "...plus numeric
// variables and the gap before string and array variables") but is
// silent on the byte encoding; the ROM disasm resolves the ambiguity.
//
// Falsification fixtures: bytes captured directly from real SAMDOS-
// written BASIC files on disks downloaded from
// ftp.nvg.ntnu.no/pub/sam-coupe/disks/utils/. Each file's third length
// (PROG+NVARS+GAP) must be ≤ the file's total Length, and for programs
// with no string/array variables it equals Length exactly.
func TestFileTypeInfoPageFormDecoding(t *testing.T) {
	cases := []struct {
		name        string
		info        [11]byte
		fileLength  uint32
		wantProgLen uint32
		wantNumVar  uint32
		wantStrArr  uint32
	}{
		{
			name:        "Auto Font (FontLoader.dsk, no string/array vars)",
			info:        [11]byte{0x01, 0xc3, 0x8c, 0x01, 0x0d, 0x8e, 0x01, 0x1f, 0x8f, 0x20, 0xff},
			fileLength:  20255,
			wantProgLen: 19651,
			wantNumVar:  19981,
			wantStrArr:  20255,
		},
		{
			name:        "Shredder (FileShredderv1.2.dsk, no string/array vars)",
			info:        [11]byte{0x00, 0xbf, 0x9b, 0x00, 0x3a, 0x9c, 0x00, 0x1b, 0x9e},
			fileLength:  7707,
			wantProgLen: 7103,
			wantNumVar:  7226,
			wantStrArr:  7707,
		},
		{
			name:        "AUTOCOMMS (CommsLoader.dsk)",
			info:        [11]byte{0x00, 0x22, 0x9f, 0x00, 0x30, 0xa0, 0x00, 0x7e, 0xa1},
			fileLength:  8861,
			wantProgLen: 7970,
			wantNumVar:  8240,
			wantStrArr:  8574,
		},
		{
			name:        "synthetic page=0, addr=0x8000",
			info:        [11]byte{0x00, 0x00, 0x80, 0x00, 0x00, 0x80, 0x00, 0x00, 0x80},
			wantProgLen: 0,
			wantNumVar:  0,
			wantStrArr:  0,
		},
		{
			name:        "synthetic page=31, addr=0xBFFF (max)",
			info:        [11]byte{0x1f, 0xff, 0xbf, 0x1f, 0xff, 0xbf, 0x1f, 0xff, 0xbf},
			wantProgLen: 0x7FFFF,
			wantNumVar:  0x7FFFF,
			wantStrArr:  0x7FFFF,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fe := &FileEntry{FileTypeInfo: c.info}
			if got := fe.ProgramLength(); got != c.wantProgLen {
				t.Errorf("ProgramLength = %d; want %d", got, c.wantProgLen)
			}
			if got := fe.NumericVariableOffset(); got != c.wantNumVar {
				t.Errorf("NumericVariableOffset = %d; want %d", got, c.wantNumVar)
			}
			if got := fe.StringArrayVariableOffset(); got != c.wantStrArr {
				t.Errorf("StringArrayVariableOffset = %d; want %d", got, c.wantStrArr)
			}
			if c.fileLength != 0 && fe.ProgramLength() > c.fileLength {
				t.Errorf("sanity: ProgramLength %d > file Length %d (impossible)",
					fe.ProgramLength(), c.fileLength)
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
