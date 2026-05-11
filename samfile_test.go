package samfile

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/petemoore/samfile/v3/sambasic"
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

// TestLoadRejectsEDSK pins the user-visible error when samfile is handed
// an Extended CPC DSK image instead of a raw MGT image.
//
// EDSK files start with the 34-byte ASCII magic
// "EXTENDED CPC DSK File\r\nDisk-Info\r\n" at offset 0 (per
// https://www.cpcwiki.eu/index.php/Format:DSK_disk_image_file_format).
// Without rejection, samfile silently treats the first 819200 bytes as
// raw MGT — the directory decode happens to land on plausible bytes
// (file names look right) but file-body sector reads at MGT offsets
// return garbage, because EDSK interleaves track-info blocks with
// sector data.
//
// The fix detects the magic at the entry to Load() and returns an
// error pointing the user at samdisk, which round-trips EDSK<->MGT.
func TestLoadRejectsEDSK(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fake.dsk")

	// 256-byte fake EDSK: magic at offset 0, zeros elsewhere. The
	// magic prefix alone is sufficient for detection; the rest of the
	// disk-info block is not inspected.
	fakeEDSK := make([]byte, 256)
	copy(fakeEDSK, []byte("EXTENDED CPC DSK File\r\nDisk-Info\r\n"))
	if err := os.WriteFile(path, fakeEDSK, 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() returned nil error on EDSK input; want rejection")
	}
	msg := err.Error()
	if !strings.Contains(msg, "EDSK") {
		t.Errorf("error message %q does not mention EDSK", msg)
	}
	if !strings.Contains(msg, "samdisk") {
		t.Errorf("error message %q does not mention the samdisk conversion command", msg)
	}
	if !strings.Contains(msg, "simonowen.com") {
		t.Errorf("error message %q does not include the samdisk URL", msg)
	}
}

// TestSAMBasicOutputRejectsEmptyInput pins the no-panic contract for
// (*SAMBasic).Output(): empty input must produce a clear error, not an
// index-out-of-range panic at sambasic.go:23.
//
// Reproduces the user-visible bug where
//
//	samfile cat -i edsk.dsk -f some-file | samfile basic-to-text
//
// piped 0 bytes (because the EDSK reader returned garbage that didn't
// match any file) and basic-to-text panicked with
// "index out of range [0] with length 0".
//
// Empty input is a degenerate but unavoidable case at a CLI pipe
// boundary: any upstream that writes nothing must not crash this stage.
func TestSAMBasicOutputRejectsEmptyInput(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Output() panicked on empty input: %v", r)
		}
	}()

	sb := NewSAMBasic([]byte{})
	err := sb.Output()
	if err == nil {
		t.Fatal("Output() returned nil error on empty input; want a clear rejection")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("error %q does not mention empty input", err.Error())
	}
}

// TestSAMBasicOutputBoundsChecks pins the no-panic contract for
// (*SAMBasic).Output() on malformed-but-non-empty input. The empty-input
// guard added in commit 00877fa handles len==0; this test covers the
// decode loop's per-byte accesses, which were also unchecked.
//
// Each case is a deliberately-truncated SAM BASIC blob that, before the
// bounds-check fix, would panic with "index out of range" inside
// sambasic.go. After the fix, all return a clear error mentioning
// truncation or invalid input. None should panic; defer/recover catches
// any escapee.
func TestSAMBasicOutputBoundsChecks(t *testing.T) {
	cases := []struct {
		name string
		data []byte
	}{
		{
			name: "single non-sentinel byte: panics on Data[1] for line-no MSB",
			data: []byte{0x00},
		},
		{
			name: "two-byte input: panics on Data[2] for line-len LSB",
			data: []byte{0x00, 0x01},
		},
		{
			name: "three-byte input: panics on Data[3] for line-len MSB",
			data: []byte{0x00, 0x01, 0x05},
		},
		{
			name: "header complete but no body and no sentinel after",
			data: []byte{0x00, 0x01, 0x00, 0x00},
		},
		{
			name: "lineLen says 10 but body is only 3 bytes",
			data: []byte{0x00, 0x01, 0x0a, 0x00, 0x41, 0x42, 0x43},
		},
		{
			name: "0xff keyword escape with no following byte",
			data: []byte{0x00, 0x01, 0x01, 0x00, 0xff},
		},
		{
			name: "no 0xff sentinel anywhere, body fits but program never terminates",
			data: []byte{0x00, 0x01, 0x02, 0x00, 0x41, 0x0d},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Output() panicked on truncated input %v: %v", c.data, r)
				}
			}()
			// Capture stdout so partial output during decode doesn't pollute
			// the test runner. We discard it; we only care about the no-panic
			// + error-returned contract.
			oldStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("os.Pipe: %v", err)
			}
			os.Stdout = w
			defer func() { os.Stdout = oldStdout }()
			drained := make(chan struct{})
			go func() {
				_, _ = io.Copy(io.Discard, r)
				close(drained)
			}()

			outErr := NewSAMBasic(c.data).Output()
			_ = w.Close()
			<-drained

			if outErr == nil {
				t.Fatalf("Output() returned nil error on truncated input %v; want non-nil", c.data)
			}
		})
	}
}

func TestNewDiskImage(t *testing.T) {
	di := NewDiskImage()
	for i, b := range di {
		if b != 0 {
			t.Fatalf("NewDiskImage()[%d] = 0x%02x; want 0x00", i, b)
		}
	}
	if len(di) != 819200 {
		t.Fatalf("NewDiskImage() length = %d; want 819200", len(di))
	}
}

func TestAddBasicFileRoundTrip(t *testing.T) {
	f := &sambasic.File{
		Lines: []sambasic.Line{
			{
				Number: 10,
				Tokens: []sambasic.Token{
					sambasic.CLEAR,
					sambasic.Number(32767),
					sambasic.String(":"),
					sambasic.LOAD,
					sambasic.String(`"`),
					sambasic.String("stub"),
					sambasic.String(`"`),
					sambasic.CODE,
					sambasic.Number(32768),
					sambasic.String(":"),
					sambasic.CALL,
					sambasic.Number(32768),
				},
			},
		},
		StartLine: 10,
	}

	di := NewDiskImage()
	if err := di.AddBasicFile("auto", f); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}

	readBack, err := di.File("auto")
	if err != nil {
		t.Fatalf("File(\"auto\"): %v", err)
	}

	wantBody := f.Bytes()
	if !bytes.Equal(readBack.Body, wantBody) {
		t.Errorf("readback body length = %d; want %d", len(readBack.Body), len(wantBody))
	}

	if readBack.Header.Type != FT_SAM_BASIC {
		t.Errorf("header type = %d; want %d (FT_SAM_BASIC)", readBack.Header.Type, FT_SAM_BASIC)
	}

	var fe *FileEntry
	for _, e := range di.DiskJournal() {
		if e.Used() && e.Name.String() == "auto" {
			fe = e
			break
		}
	}
	if fe == nil {
		t.Fatal("auto entry not found in disk journal")
	}
	if fe.Type != FT_SAM_BASIC {
		t.Errorf("fe.Type = %d; want %d", fe.Type, FT_SAM_BASIC)
	}
	if fe.MGTFlags != 0x20 {
		t.Errorf("MGTFlags = 0x%02x; want 0x20", fe.MGTFlags)
	}
	if fe.SAMBASICStartLine != 10 {
		t.Errorf("SAMBASICStartLine = %d; want 10", fe.SAMBASICStartLine)
	}

	if fe.ProgramLength() != f.NVARSOffset() {
		t.Errorf("ProgramLength() = %d; want %d", fe.ProgramLength(), f.NVARSOffset())
	}
	if fe.ProgramLength()+fe.NumericVariablesSize() != f.NUMENDOffset() {
		t.Errorf("ProgramLength+NumericVarsSize = %d; want %d",
			fe.ProgramLength()+fe.NumericVariablesSize(), f.NUMENDOffset())
	}
	if fe.ProgramLength()+fe.NumericVariablesSize()+fe.GapSize() != f.SAVARSOffset() {
		t.Errorf("ProgramLength+NumericVarsSize+GapSize = %d; want %d",
			fe.ProgramLength()+fe.NumericVariablesSize()+fe.GapSize(), f.SAVARSOffset())
	}

	sb := NewSAMBasic(readBack.Body)
	if err := sb.Output(); err != nil {
		t.Errorf("SAMBasic.Output() on readback: %v", err)
	}
}

func TestAddBasicFileNoAutoRun(t *testing.T) {
	f := &sambasic.File{
		Lines: []sambasic.Line{
			{Number: 10, Tokens: []sambasic.Token{sambasic.PRINT, sambasic.String("hello")}},
		},
		StartLine: 0xFFFF,
	}
	di := NewDiskImage()
	if err := di.AddBasicFile("test", f); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}
	var fe *FileEntry
	for _, e := range di.DiskJournal() {
		if e.Used() && e.Name.String() == "test" {
			fe = e
			break
		}
	}
	if fe == nil {
		t.Fatal("test entry not found")
	}
	if fe.ExecutionAddressDiv16K != 0xFF {
		t.Errorf("ExecutionAddressDiv16K = 0x%02x; want 0xFF (no auto-run)", fe.ExecutionAddressDiv16K)
	}
}

func TestMultiFileBasicAndCode(t *testing.T) {
	di := NewDiskImage()

	codeBody := bytes.Repeat([]byte{0xAA}, 1000)
	if err := di.AddCodeFile("samdos2", codeBody, 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile(samdos2): %v", err)
	}

	basicFile := &sambasic.File{
		Lines: []sambasic.Line{
			{
				Number: 10,
				Tokens: []sambasic.Token{
					sambasic.CLEAR,
					sambasic.Number(32767),
					sambasic.String(":"),
					sambasic.LOAD,
					sambasic.String(`"`),
					sambasic.String("stub"),
					sambasic.String(`"`),
					sambasic.CODE,
					sambasic.Number(32768),
					sambasic.String(":"),
					sambasic.CALL,
					sambasic.Number(32768),
				},
			},
		},
		StartLine: 10,
	}
	if err := di.AddBasicFile("auto", basicFile); err != nil {
		t.Fatalf("AddBasicFile(auto): %v", err)
	}

	stubBody := bytes.Repeat([]byte{0xBB}, 100)
	if err := di.AddCodeFile("stub", stubBody, 0x8000, 0x8000); err != nil {
		t.Fatalf("AddCodeFile(stub): %v", err)
	}

	for _, name := range []string{"samdos2", "auto", "stub"} {
		if _, err := di.File(name); err != nil {
			t.Errorf("File(%q): %v", name, err)
		}
	}

	dj := di.DiskJournal()
	used := dj.UsedFileEntries()
	if len(used) != 3 {
		t.Errorf("used entries = %d; want 3", len(used))
	}
	for i := 0; i < len(used); i++ {
		for j := i + 1; j < len(used); j++ {
			a := dj[used[i]].SectorAddressMap
			b := dj[used[j]].SectorAddressMap
			for k := 0; k < len(a); k++ {
				if a[k]&b[k] != 0 {
					t.Errorf("sector maps for slots %d and %d overlap at byte %d: 0x%02x & 0x%02x",
						used[i], used[j], k, a[k], b[k])
				}
			}
		}
	}
}

// firstSectorBytes returns the first sector's payload for the named
// file. Convenience wrapper for body-header assertions: AddCodeFile
// writes the 9-byte FileHeader at body bytes 0..8 = first-sector
// bytes 0..8, so callers can index directly.
func firstSectorBytes(t *testing.T, di *DiskImage, name string) [512]byte {
	t.Helper()
	for _, fe := range di.DiskJournal() {
		if fe.Used() && fe.Name.String() == name {
			sd, err := di.SectorData(fe.FirstSector)
			if err != nil {
				t.Fatalf("SectorData(%v): %v", fe.FirstSector, err)
			}
			return [512]byte(*sd)
		}
	}
	t.Fatalf("file %q not present on disk", name)
	return [512]byte{}
}

func usedFileEntry(t *testing.T, di *DiskImage, name string) *FileEntry {
	t.Helper()
	for _, fe := range di.DiskJournal() {
		if fe.Used() && fe.Name.String() == name {
			return fe
		}
	}
	t.Fatalf("file %q not present on disk", name)
	return nil
}

// TestFileHeaderRawEmitsExecutionAddress pins down the body-header
// bytes 5-6 encoding. The previous implementation hard-coded these
// to 0x00 0x00, which broke ROM's LOAD-CODE auto-exec gate at
// rom-disasm:22471-22484: a CODE file loaded via `LOAD ... CODE addr`
// is auto-executed unless BOTH dir byte 0xF2 AND body-header byte 6
// are 0xFF. Without this fix, a no-auto-exec file emitted by
// AddCodeFile would mis-fire its post-load auto-exec on real SAM.
func TestFileHeaderRawEmitsExecutionAddress(t *testing.T) {
	cases := []struct {
		name             string
		execDiv16K       byte
		execMod16KLo     byte
		wantByte5        byte
		wantByte6        byte
	}{
		{"no auto-exec (sentinel 0xff 0xff)", 0xFF, 0xFF, 0xFF, 0xFF},
		{"auto-exec at 0x4000 (page 0, offset 0)", 0x00, 0x00, 0x00, 0x00},
		{"auto-exec at 0x9039 (page 1, low byte 0x39)", 0x01, 0x39, 0x01, 0x39},
		{"auto-exec at 0xC000 (page 2, low byte 0x00)", 0x02, 0x00, 0x02, 0x00},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fh := &FileHeader{
				Type:                     FT_CODE,
				LengthMod16K:             100,
				PageOffset:               0x8000,
				ExecutionAddressDiv16K:   c.execDiv16K,
				ExecutionAddressMod16KLo: c.execMod16KLo,
				Pages:                    0,
				StartPage:                1,
			}
			raw := fh.Raw()
			if raw[5] != c.wantByte5 {
				t.Errorf("raw[5] = 0x%02x; want 0x%02x", raw[5], c.wantByte5)
			}
			if raw[6] != c.wantByte6 {
				t.Errorf("raw[6] = 0x%02x; want 0x%02x", raw[6], c.wantByte6)
			}
		})
	}
}

// TestAddCodeFileBodyHeaderAutoExecGate is the user-facing check for
// the auto-exec gate: AddCodeFile with executionAddress=0 must produce
// a body header whose bytes 5-6 are both 0xFF, so ROM's LOAD-CODE
// path returns cleanly to BASIC after a load.
func TestAddCodeFileBodyHeaderAutoExecGate(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("F", []byte("test"), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	first := firstSectorBytes(t, di, "F")
	if first[5] != 0xFF || first[6] != 0xFF {
		t.Errorf("body header bytes 5-6 = 0x%02x 0x%02x; want 0xFF 0xFF (no auto-exec)", first[5], first[6])
	}
}

// TestAddCodeFileBodyHeaderExecAddrMirror pairs with the auto-exec
// gate test above: when executionAddress IS set, the body header
// must mirror ExecutionAddressDiv16K at byte 5 and ExecutionAddress-
// Mod16K's low byte at byte 6 (the high byte doesn't fit in the
// 9-byte body header — the dir entry's 0xF3-0xF4 are authoritative).
func TestAddCodeFileBodyHeaderExecAddrMirror(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("F", make([]byte, 1024), 0x8000, 0x8200); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	fe := usedFileEntry(t, di, "F")
	first := firstSectorBytes(t, di, "F")
	if first[5] != fe.ExecutionAddressDiv16K {
		t.Errorf("body header byte 5 = 0x%02x; want 0x%02x (fe.ExecutionAddressDiv16K)", first[5], fe.ExecutionAddressDiv16K)
	}
	if first[6] != byte(fe.ExecutionAddressMod16K&0xFF) {
		t.Errorf("body header byte 6 = 0x%02x; want 0x%02x (fe.ExecutionAddressMod16K low)", first[6], fe.ExecutionAddressMod16K&0xFF)
	}
}

// TestAddCodeFileMirrorsMGTFutureAndPast verifies that addFile mirrors
// the 9-byte body header into the directory entry's MGTFutureAndPast
// field at offsets 1..9. Real disks saved by ROM SAVE always carry
// this mirror; previously AddCodeFile left the region zeroed and only
// AddBasicFile populated it.
func TestAddCodeFileMirrorsMGTFutureAndPast(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("F", []byte("hello world"), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	fe := usedFileEntry(t, di, "F")
	first := firstSectorBytes(t, di, "F")
	for i := 0; i < 9; i++ {
		if fe.MGTFutureAndPast[i+1] != first[i] {
			t.Errorf("MGTFutureAndPast[%d] = 0x%02x; want 0x%02x (body header byte %d)",
				i+1, fe.MGTFutureAndPast[i+1], first[i], i)
		}
	}
	if fe.MGTFutureAndPast[0] != 0x00 {
		t.Errorf("MGTFutureAndPast[0] = 0x%02x; want 0x00 (reserved)", fe.MGTFutureAndPast[0])
	}
}

// TestAddBasicFileStillMirrorsMGTFutureAndPast guards against a
// regression where the AddBasicFile path stopped populating
// MGTFutureAndPast after the addFile-level mirror was added: the
// mirror covers both code paths, so AddBasicFile output should still
// carry the body-header mirror in the dir entry.
func TestAddBasicFileStillMirrorsMGTFutureAndPast(t *testing.T) {
	di := NewDiskImage()
	bf := &sambasic.File{
		Lines:     []sambasic.Line{{Number: 10, Tokens: []sambasic.Token{sambasic.PRINT, sambasic.String("hi")}}},
		StartLine: 10,
	}
	if err := di.AddBasicFile("auto", bf); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}
	fe := usedFileEntry(t, di, "auto")
	first := firstSectorBytes(t, di, "auto")
	for i := 0; i < 9; i++ {
		if fe.MGTFutureAndPast[i+1] != first[i] {
			t.Errorf("MGTFutureAndPast[%d] = 0x%02x; want 0x%02x (body header byte %d)",
				i+1, fe.MGTFutureAndPast[i+1], first[i], i)
		}
	}
}

// TestSetStartAddressPageUnusedBits verifies the override mechanism
// for canonical-SAVE byte parity: SetStartAddressPageUnusedBits must
// set the upper 3 bits of StartAddressPage in both the directory
// entry (raw[0xEC]) and the matching body-header byte 8 while
// preserving the low 5 bits (the actual page index derived by
// AddCodeFile from the load address).
//
// The canonical FRED 02 / Defender samdos2 install records 0x7D there
// (= 3<<5 | 0x1D = page 29 with the top two bits set). The bits are
// unread by ROM and SAMDOS; this method exists for byte-perfect
// parity with historical disk images.
func TestSetStartAddressPageUnusedBits(t *testing.T) {
	cases := []struct {
		name             string
		bits             uint8
		wantStartAddrPg  byte
	}{
		{"zero (default)", 0, 0x1D},
		{"bit 5 only", 1, 0x3D},
		{"FRED 02 / Defender samdos2 (bits 5+6)", 3, 0x7D},
		{"bit 7 only", 4, 0x9D},
		{"all three (0xE0)", 7, 0xFD},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			di := NewDiskImage()
			const loadAddress = uint32(491529) // page 29 + offset 9
			if err := di.AddCodeFile("samdos2", make([]byte, 10000), loadAddress, 0); err != nil {
				t.Fatalf("AddCodeFile: %v", err)
			}
			if got := usedFileEntry(t, di, "samdos2").StartAddressPage; got != 0x1D {
				t.Fatalf("StartAddressPage before override = 0x%02x; want 0x1D (page 29, no unused bits)", got)
			}
			if err := di.SetStartAddressPageUnusedBits("samdos2", c.bits); err != nil {
				t.Fatalf("SetStartAddressPageUnusedBits: %v", err)
			}
			fe := usedFileEntry(t, di, "samdos2")
			if fe.StartAddressPage != c.wantStartAddrPg {
				t.Errorf("dir entry StartAddressPage = 0x%02x; want 0x%02x", fe.StartAddressPage, c.wantStartAddrPg)
			}
			first := firstSectorBytes(t, di, "samdos2")
			if first[8] != c.wantStartAddrPg {
				t.Errorf("body header byte 8 = 0x%02x; want 0x%02x", first[8], c.wantStartAddrPg)
			}
			// ROM masks to 0x1F when reading, so decoded Start is unchanged.
			if got, want := fe.StartAddress(), loadAddress; got != want {
				t.Errorf("decoded Start after override = %d; want %d (unused bits must not affect the address)", got, want)
			}
		})
	}
}

// TestSetStartAddressPageUnusedBitsThreeWayConsistency is a property
// test: after SetStartAddressPageUnusedBits, the StartPage byte must
// agree across all three on-disk locations — dir 0xEC, body byte 8,
// and the MGTFutureAndPast mirror at dir 0xDB (SAMDOS source
// f.s:462-471, c.s:1376-1379; see BODY-MIRROR-AT-DIR-D3-DB in
// disk-validity-rules.md).
func TestSetStartAddressPageUnusedBitsThreeWayConsistency(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("F", make([]byte, 500), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	if err := di.SetStartAddressPageUnusedBits("F", 5); err != nil {
		t.Fatalf("SetStartAddressPageUnusedBits: %v", err)
	}
	fe := usedFileEntry(t, di, "F")
	first := firstSectorBytes(t, di, "F")

	dirEC := fe.StartAddressPage     // dir byte 0xEC
	dirDB := fe.MGTFutureAndPast[9]  // dir byte 0xDB (mirror of body byte 8)
	bodyB8 := first[8]               // body header byte 8

	if dirEC != bodyB8 || dirEC != dirDB {
		t.Errorf("three-way mismatch: dir[0xEC]=0x%02x, body[8]=0x%02x, MGTFutureAndPast[9]=0x%02x — all must be equal",
			dirEC, bodyB8, dirDB)
	}
}

// TestSetStartAddressPageUnusedBitsOutOfRange confirms that values
// above 7 are rejected — the API is for 3 bits only.
func TestSetStartAddressPageUnusedBitsOutOfRange(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("F", []byte("test"), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	for _, bad := range []uint8{8, 16, 32, 0x60, 0x7D, 0xFF} {
		if err := di.SetStartAddressPageUnusedBits("F", bad); err == nil {
			t.Errorf("SetStartAddressPageUnusedBits(%d) returned nil; want out-of-range error", bad)
		}
	}
}

// TestSetStartAddressPageUnusedBitsMissingFile confirms the helper
// errors cleanly when the named file isn't on disk.
func TestSetStartAddressPageUnusedBitsMissingFile(t *testing.T) {
	di := NewDiskImage()
	err := di.SetStartAddressPageUnusedBits("nope", 3)
	if err == nil {
		t.Fatal("SetStartAddressPageUnusedBits returned nil error for missing file")
	}
	if !strings.Contains(err.Error(), "nope") {
		t.Errorf("error message = %q; want it to mention the missing filename", err.Error())
	}
}

