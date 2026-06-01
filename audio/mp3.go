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
	// findBitrateIndex is unexported, so we scan the known MPEG bitrate table
	// to find the correct BitrateIndex for this bitrate + MPEG version.
	// MPEG-I (44100/48000/32000 Hz): column 3; MPEG-II: column 2; MPEG-2.5: column 0.
	// bitRates[i][version] — we find i such that bitRates[i][version] == bitrateKbps.
	// Since we can't call the unexported function, derive BitrateIndex from
	// the relationship: BitrateIndex is the same as what NewEncoder stored for 128,
	// so we locate 256 by offsetting — but this is fragile. Instead, reuse the
	// exported CheckConfig return (mpegVersion) and scan ourselves via a local table.
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
	bitsPerSlot := int64(8)
	avg := float64(enc.Mpeg.GranulesPerFrame) * 576.0 / float64(sampleRate) *
		(float64(bitrateKbps) * 1000 / float64(bitsPerSlot))
	enc.Mpeg.WholeSlotsPerFrame = int64(avg)
	enc.Mpeg.FracSlotsPerFrame = avg - float64(enc.Mpeg.WholeSlotsPerFrame)
	enc.Mpeg.Slot_lag = -enc.Mpeg.FracSlotsPerFrame
	if enc.Mpeg.FracSlotsPerFrame == 0 {
		enc.Mpeg.Padding = 0
	}

	// Feed interleaved int16 stereo PCM; shine writes MP3 frames to w.
	return enc.Write(w, pcm)
}
