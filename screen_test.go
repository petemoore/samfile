package samfile

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// makeMode4Body builds a synthetic MODE 4 SCREEN$ body: every pixel
// set to nibble value clutIdx, a CLUT mapping that nibble to SAM
// palette index palIdx, and the canonical 24617-byte trailer (palette
// A == palette B, no line interrupts).
func makeMode4Body(clutIdx, palIdx uint8) []byte {
	body := make([]byte, screenMode4TotalLen)
	packed := clutIdx<<4 | clutIdx
	for i := 0; i < screenMode4DataLen; i++ {
		body[i] = packed
	}
	// CLUT A at +0, CLUT B at +20.
	body[screenMode4DataLen+int(clutIdx)] = palIdx
	body[screenMode4DataLen+20+int(clutIdx)] = palIdx
	body[len(body)-1] = 0xFF
	return body
}

func TestSAMPaletteColorMatchesReference(t *testing.T) {
	// Spot values cross-checked against scrimage/palette.py
	// get_sam_palette: index 0 = black, 127 = full white. Index 2
	// (0b010) sets the high red bit -> r=2; index 1 (0b001) sets the
	// high blue bit -> b=2; both map through samIntensity.
	cases := []struct {
		i    uint8
		want color.RGBA
	}{
		{0, color.RGBA{0x00, 0x00, 0x00, 0xff}},
		{127, color.RGBA{0xff, 0xff, 0xff, 0xff}},
		{0x02, color.RGBA{samIntensity[2], 0x00, 0x00, 0xff}},
		{0x01, color.RGBA{0x00, 0x00, samIntensity[2], 0xff}},
	}
	for _, c := range cases {
		got := SAMPaletteColor(c.i)
		if got != c.want {
			t.Errorf("SAMPaletteColor(%d) = %v, want %v", c.i, got, c.want)
		}
	}
}

func TestScreenWritePNGDimensionsAndColor(t *testing.T) {
	const clutIdx, palIdx = 5, 0x02 // pure red
	body := makeMode4Body(clutIdx, palIdx)
	scr := NewScreen(body, 3)
	var buf bytes.Buffer
	if err := scr.WritePNG(&buf); err != nil {
		t.Fatalf("WritePNG: %v", err)
	}
	img, err := png.Decode(&buf)
	if err != nil {
		t.Fatalf("decode PNG: %v", err)
	}
	if img.Bounds() != image.Rect(0, 0, 256, 192) {
		t.Fatalf("PNG bounds = %v, want 256x192", img.Bounds())
	}
	want := SAMPaletteColor(palIdx)
	r, g, b, _ := img.At(0, 0).RGBA()
	if uint8(r>>8) != want.R || uint8(g>>8) != want.G || uint8(b>>8) != want.B {
		t.Fatalf("pixel (0,0) = (%d,%d,%d), want (%d,%d,%d)", r>>8, g>>8, b>>8, want.R, want.G, want.B)
	}
}

func TestScreenWritePNGRejectsNonMode4(t *testing.T) {
	scr := NewScreen(makeMode4Body(0, 0), 0) // MODE 1
	if err := scr.WritePNG(&bytes.Buffer{}); err == nil {
		t.Fatal("expected error for non-MODE-4 screen, got nil")
	}
}

func TestScreenWritePNGRejectsShortBody(t *testing.T) {
	scr := NewScreen(make([]byte, 100), 3)
	if err := scr.WritePNG(&bytes.Buffer{}); err == nil {
		t.Fatal("expected error for short body, got nil")
	}
}

// TestAddScreenFileRoundTrip adds a SCREEN$ file to a blank disk,
// re-reads the directory, and confirms the type and body survive.
func TestAddScreenFileRoundTrip(t *testing.T) {
	di := NewDiskImage()
	body := makeMode4Body(7, 0x49)
	if err := di.AddScreenFile("TITLE", body, 3); err != nil {
		t.Fatalf("AddScreenFile: %v", err)
	}
	dj := di.DiskJournal()
	var found *FileEntry
	for _, fe := range dj {
		if fe.Used() && fe.Name.String() == "TITLE" {
			found = fe
			break
		}
	}
	if found == nil {
		t.Fatal("TITLE not found in directory after AddScreenFile")
	}
	if found.Type != FT_SCREEN {
		t.Errorf("Type = %v, want %v", found.Type, FT_SCREEN)
	}
	if found.FileTypeInfo[0] != 3 {
		t.Errorf("stored mode = %d, want 3", found.FileTypeInfo[0])
	}
	f, err := di.File("TITLE")
	if err != nil {
		t.Fatalf("File: %v", err)
	}
	if !bytes.Equal(f.Body, body) {
		t.Fatalf("body round-trip mismatch: got %d bytes, want %d", len(f.Body), len(body))
	}
}

// TestWriteFileEntryBeyond20Files guards the cylinder-interleave fix:
// adding more than 20 files (which spill past directory track 0) must
// keep all of them readable. The pre-fix flat index<<8 offset silently
// dropped entries 20+.
func TestWriteFileEntryBeyond20Files(t *testing.T) {
	di := NewDiskImage()
	const n = 25
	for i := 0; i < n; i++ {
		// Small distinct CODE files; load address arbitrary.
		data := []byte{byte(i), byte(i + 1), byte(i + 2)}
		name := "F" + string(rune('A'+i))
		if err := di.AddCodeFile(name, data, 0x8000, 0); err != nil {
			t.Fatalf("AddCodeFile %q: %v", name, err)
		}
	}
	dj := di.DiskJournal()
	used := 0
	for _, fe := range dj {
		if fe.Used() {
			used++
		}
	}
	if used != n {
		t.Fatalf("after adding %d files, directory shows %d used entries (cylinder-interleave write bug?)", n, used)
	}
	// Confirm a beyond-track-0 file (index >= 20) reads back its body.
	f, err := di.File("F" + string(rune('A'+22)))
	if err != nil {
		t.Fatalf("File for 23rd added file: %v", err)
	}
	if !bytes.Equal(f.Body, []byte{22, 23, 24}) {
		t.Fatalf("23rd file body = %v, want [22 23 24]", f.Body)
	}
}
