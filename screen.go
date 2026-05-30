package samfile

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"io"
)

// samIntensity maps a SAM 3-bit colour-channel value (0..7) onto an
// 8-bit RGB channel. The table is round(value/7*255) for value 0..7,
// matching SimCoupe's IO::Palette() at maxintensity 255 and the
// reference scrimage decoder (scrimage/palette.py).
var samIntensity = [8]uint8{0x00, 0x24, 0x49, 0x6d, 0x92, 0xb6, 0xdb, 0xff}

// SAMPaletteColor returns the RGB colour for one of the 128 SAM Coupé
// hardware palette indices (i in 0..127). The bit-twiddling matches
// SimCoupe's IO::Palette() and scrimage/palette.py's get_sam_palette:
// each of the red/green/blue channels is assembled from three palette
// bits into a 0..7 intensity, then expanded through samIntensity.
func SAMPaletteColor(i uint8) color.RGBA {
	r := ((i & 0x02) << 0) | ((i & 0x20) >> 3) | ((i & 0x08) >> 3)
	g := ((i & 0x04) >> 1) | ((i & 0x40) >> 4) | ((i & 0x08) >> 3)
	b := ((i & 0x01) << 1) | ((i & 0x10) >> 2) | ((i & 0x08) >> 3)
	return color.RGBA{
		R: samIntensity[r&0x07],
		G: samIntensity[g&0x07],
		B: samIntensity[b&0x07],
		A: 0xff,
	}
}

// Screen mode 4 geometry: 256×192 pixels, 2 pixels packed per byte
// (high nibble = left pixel, low nibble = right pixel), so the display
// dump is 256/2 * 192 = 24576 bytes. After the display bytes the
// SCREEN$ body carries the saved hardware state used to restore the
// screen on LOAD.
const (
	screenMode4Width    = 256
	screenMode4Height   = 192
	screenMode4DataLen  = screenMode4Width / 2 * screenMode4Height // 24576
	screenMode4CLUTLen  = 16                                       // 16-entry colour look-up table
	screenMode4TotalLen = 24617                                    // canonical SAVE SCREEN$ length
)

// Screen wraps the body bytes of a SAM SCREEN$ (FT_SCREEN, type 20)
// file so they can be decoded into a Go image. Data must be the file
// body without its 9-byte FileHeader prefix — i.e. the Body field of
// a File returned by DiskImage.File, or a raw SCREEN$ dump as written
// by BASIC's SAVE "name" SCREEN$.
//
// Only MODE 4 (the common photographic/16-colour mode) is decoded.
// The body layout for MODE 4 is:
//
//	24576 bytes  display memory (2 px/byte, high nibble left)
//	   16 bytes  CLUT — palette index per 4-bit pixel value (palette A)
//	    4 bytes  mode/border bytes
//	   16 bytes  CLUT copy used for the FLASH alternate (palette B)
//	    4 bytes  mode/border copy
//	    n bytes  optional line-interrupt records (4 bytes each)
//	    1 byte   0xFF terminator
//
// total 24617 bytes when there are no line interrupts. When palette A
// and palette B differ the screen was flashing; PNG output renders
// palette A only, while GIF output (Mode = ScreenGIF) animates both.
type Screen struct {
	Data []byte
	// Mode is the SAM screen MODE byte as stored in the directory
	// entry's FileTypeInfo[0] (0..3 for MODE 1..4). A raw SCREEN$ body
	// on stdin has no directory entry, so callers default this to 3
	// (MODE 4). Only MODE 4 is currently decodable.
	Mode uint8
}

// NewScreen wraps a SCREEN$ body for decoding. data is taken by
// reference, not copied. mode is the SAM MODE byte (0..3 for MODE
// 1..4); pass 3 for the common MODE 4 case.
func NewScreen(data []byte, mode uint8) *Screen {
	return &Screen{Data: data, Mode: mode}
}

// trailer parses the post-display portion of a MODE 4 body, returning
// palette A, palette B and the line-interrupt records (each a 4-byte
// {position, palette-index, colourA, colourB} tuple). Missing optional
// sections are returned empty rather than erroring, so screens saved
// with a short trailer still decode.
func (s *Screen) trailer() (paletteA, paletteB []byte, lineInts []byte) {
	t := s.Data[screenMode4DataLen:]
	clip := func(lo, hi int) []byte {
		if lo > len(t) {
			return nil
		}
		if hi > len(t) {
			hi = len(t)
		}
		return t[lo:hi]
	}
	paletteA = clip(0, 16)
	paletteB = clip(20, 36)
	// Line interrupts sit between offset 40 and the trailing 0xFF.
	li := clip(40, len(t))
	if n := len(li); n > 0 && li[n-1] == 0xFF {
		li = li[:n-1]
	}
	lineInts = li
	return
}

// renderMode4 decodes the display bytes into a 256×192 *image.RGBA
// using the supplied CLUT (16 SAM palette indices). Line interrupts,
// if present, recolour the image from each interrupt's scan line
// downward (matching scrimage's apply_line_interrupts); palette
// selects which of the interrupt's two stored colours to use (0 = A,
// 1 = B).
func (s *Screen) renderMode4(clut, lineInts []byte, palette int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, screenMode4Width, screenMode4Height))
	// Build the per-CLUT-entry RGB colours.
	var clutColor [screenMode4CLUTLen]color.RGBA
	for i := 0; i < screenMode4CLUTLen; i++ {
		idx := uint8(0)
		if i < len(clut) {
			idx = clut[i]
		}
		clutColor[i] = SAMPaletteColor(idx)
	}
	display := s.Data[:screenMode4DataLen]
	for y := 0; y < screenMode4Height; y++ {
		for bx := 0; bx < screenMode4Width/2; bx++ {
			b := display[y*(screenMode4Width/2)+bx]
			high := (b >> 4) & 0x0f
			low := b & 0x0f
			img.SetRGBA(bx*2, y, clutColor[high])
			img.SetRGBA(bx*2+1, y, clutColor[low])
		}
	}
	applyLineInterrupts(img, clutColor, lineInts, palette)
	return img
}

// applyLineInterrupts re-colours img from each interrupt's scan line
// downward, mirroring scrimage's apply_line_interrupts. Each record is
// {y, clutIndex, colourA, colourB}: every pixel below row y that
// currently matches the CLUT entry's colour is repainted in the
// interrupt's chosen SAM colour, and the CLUT entry is then updated so
// later interrupts see the new colour.
func applyLineInterrupts(img *image.RGBA, clutColor [screenMode4CLUTLen]color.RGBA, lineInts []byte, palette int) {
	for i := 0; i+3 < len(lineInts); i += 4 {
		y := int(lineInts[i])
		c := int(lineInts[i+1])
		var newIdx uint8
		if palette == 0 {
			newIdx = lineInts[i+2]
		} else {
			newIdx = lineInts[i+3]
		}
		if c < 0 || c >= screenMode4CLUTLen {
			continue
		}
		old := clutColor[c]
		newColor := SAMPaletteColor(newIdx)
		for yp := y + 1; yp < screenMode4Height; yp++ {
			for x := 0; x < screenMode4Width; x++ {
				if img.RGBAAt(x, yp) == old {
					img.SetRGBA(x, yp, newColor)
				}
			}
		}
		clutColor[c] = newColor
	}
}

// validateMode4 returns an error if the body is too short to be a
// MODE 4 screen or the mode is unsupported.
func (s *Screen) validateMode4() error {
	if s.Mode != 3 {
		return fmt.Errorf("only MODE 4 (mode byte 3) SCREEN$ decoding is supported; got mode byte %d (MODE %d)", s.Mode, s.Mode+1)
	}
	if len(s.Data) < screenMode4DataLen+screenMode4CLUTLen {
		return fmt.Errorf("SCREEN$ body too short for MODE 4: have %d bytes, need at least %d", len(s.Data), screenMode4DataLen+screenMode4CLUTLen)
	}
	return nil
}

// flashing reports whether palette A and palette B differ, i.e. the
// screen was saved with a flashing (alternating-palette) effect.
func flashing(a, b []byte) bool {
	if len(a) != len(b) || len(a) == 0 {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return true
		}
	}
	return false
}

// WritePNG decodes the MODE 4 screen and writes it as a PNG to w.
// Flashing screens are rendered with palette A only.
func (s *Screen) WritePNG(w io.Writer) error {
	if err := s.validateMode4(); err != nil {
		return err
	}
	paletteA, _, lineInts := s.trailer()
	img := s.renderMode4(paletteA, lineInts, 0)
	return png.Encode(w, img)
}

// WriteGIF decodes the MODE 4 screen and writes it as a GIF to w. If
// the screen is flashing (palette A != palette B) the GIF animates the
// two palettes; otherwise it is a single static frame.
func (s *Screen) WriteGIF(w io.Writer) error {
	if err := s.validateMode4(); err != nil {
		return err
	}
	paletteA, paletteB, lineInts := s.trailer()
	frameA := s.renderMode4(paletteA, lineInts, 0)
	if !flashing(paletteA, paletteB) {
		return gif.Encode(w, frameA, nil)
	}
	frameB := s.renderMode4(paletteB, lineInts, 1)
	g := &gif.GIF{
		Image:     []*image.Paletted{rgbaToPaletted(frameA), rgbaToPaletted(frameB)},
		Delay:     []int{33, 33}, // ~3 Hz flash, in hundredths of a second
		LoopCount: 0,
	}
	return gif.EncodeAll(w, g)
}

// rgbaToPaletted converts an RGBA image to a paletted image using the
// distinct colours it contains (a MODE 4 frame has at most 16), so GIF
// frames render without dithering.
func rgbaToPaletted(img *image.RGBA) *image.Paletted {
	seen := map[color.RGBA]bool{}
	var pal color.Palette
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if !seen[c] {
				seen[c] = true
				pal = append(pal, c)
			}
		}
	}
	if len(pal) == 0 {
		pal = color.Palette{color.RGBA{A: 0xff}}
	}
	out := image.NewPaletted(b, pal)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			out.SetColorIndex(x, y, uint8(pal.Index(img.RGBAAt(x, y))))
		}
	}
	return out
}
