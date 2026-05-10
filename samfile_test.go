package samfile

import (
	"bytes"
	"testing"
)

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
