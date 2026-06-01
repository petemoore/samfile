package audio

import (
	"fmt"
	"io"

	"github.com/braheezy/shine-mp3/pkg/mp3"
)

func encodeMP3(w io.Writer, pcm []int16, sampleRate, bitrateKbps int, tags Tags) error {
	if bitrateKbps == 0 {
		bitrateKbps = 256
	}

	// Write ID3v2.3 provenance tag before any MP3 frames.
	if tag := id3v2Tag(tags); tag != nil {
		if _, err := w.Write(tag); err != nil {
			return err
		}
	}

	// Validate sample rate and desired bitrate combination before constructing.
	if mp3.CheckConfig(sampleRate, bitrateKbps) < 0 {
		return fmt.Errorf("audio/mp3: unsupported sample rate %d Hz / bitrate %d kbps combination", sampleRate, bitrateKbps)
	}

	// Construct shine encoder: stereo, given sample rate.
	// NewEncoder defaults to 128 kbps; we update the exported bitrate fields
	// to match the desired bitrate before calling Write.
	enc := mp3.NewEncoder(sampleRate, 2)

	// Update bitrate-dependent exported fields.
	// shine's findBitrateIndex is unexported, so we reproduce the lookup here
	// using the same table (bitRates[index][mpegVersion]) from tables.go.
	// Columns: 0=MPEG-2.5, 1=reserved, 2=MPEG-II, 3=MPEG-I.
	// We find index i such that mpegBitRates[i][mpegVersion] == bitrateKbps.
	mpegBitRates := [16][4]int64{
		{-1, -1, -1, -1}, {8, -1, 8, 32}, {16, -1, 16, 40}, {24, -1, 24, 48},
		{32, -1, 32, 56}, {40, -1, 40, 64}, {48, -1, 48, 80}, {56, -1, 56, 96},
		{64, -1, 64, 112}, {-1, -1, 80, 128}, {-1, -1, 96, 160}, {-1, -1, 112, 192},
		{-1, -1, 128, 224}, {-1, -1, 144, 256}, {-1, -1, 160, 320}, {-1, -1, -1, -1},
	}
	mpegVer := int(enc.Mpeg.Version) // 0=2.5, 2=II, 3=I
	bitrateIdx := int64(-1)
	for i, row := range mpegBitRates {
		if row[mpegVer] == int64(bitrateKbps) {
			bitrateIdx = int64(i)
			break
		}
	}
	if bitrateIdx < 0 {
		return fmt.Errorf("audio/mp3: bitrate %d kbps not found in MPEG table for version %d", bitrateKbps, mpegVer)
	}

	enc.Mpeg.Bitrate = int64(bitrateKbps)
	enc.Mpeg.BitrateIndex = bitrateIdx

	// Recompute slot-count fields that depend on bitrate (mirrors NewEncoder logic).
	// enc.Mpeg.BitsPerSlot is always 8 for Layer III (set by NewEncoder).
	avg := float64(enc.Mpeg.GranulesPerFrame) * 576.0 / float64(sampleRate) *
		(float64(bitrateKbps) * 1000 / float64(enc.Mpeg.BitsPerSlot))
	enc.Mpeg.WholeSlotsPerFrame = int64(avg)
	enc.Mpeg.FracSlotsPerFrame = avg - float64(enc.Mpeg.WholeSlotsPerFrame)
	enc.Mpeg.Slot_lag = -enc.Mpeg.FracSlotsPerFrame
	if enc.Mpeg.FracSlotsPerFrame == 0 {
		enc.Mpeg.Padding = 0
	}

	// Condition the PCM before handing it to shine. shine's windowFilterSubband
	// walks the interleaved buffer with unsafe pointer arithmetic and, on the
	// final frame, advances one stride PAST the last sample it reads — for the
	// 2nd (odd-offset) channel that overshoots a frame-exact slice. Untreated
	// this both (a) reads garbage into the final frame of a partial input and
	// (b) forms an out-of-bounds pointer that trips Go's -race/checkptr (a fatal
	// "invalid allocation"/"bad pointer"). So we (1) pad to a whole number of
	// frames, zero-filling the final frame (clean trailing silence, not garbage)
	// and (2) give the backing array a window-sized margin beyond that length so
	// the trailing pointer always lands strictly inside the allocation — never
	// at or past the span limit — regardless of allocation size class. The
	// margin is pure capacity beyond len(padded); shine reads only `total`
	// samples (whole frames), so it is never encoded. (Reported upstream:
	// braheezy/shine-mp3 windowFilterSubband / Write.)
	frame := int(enc.Mpeg.GranulesPerFrame) * 576 * 2 // 576 = GRANULE_SIZE (Layer III)
	nFrames := (len(pcm) + frame - 1) / frame
	total := nFrames * frame
	const slack = 64                             // > stride; window-sized margin
	padded := make([]int16, total+slack)[:total] // len=total (whole frames) + slack cap
	copy(padded, pcm)

	// Feed interleaved int16 stereo PCM; shine writes MP3 frames to w.
	return enc.Write(w, padded)
}
