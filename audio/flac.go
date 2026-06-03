package audio

import (
	"io"

	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/frame"
	"github.com/mewkiz/flac/meta"
)

const flacBlockSize = 4096

// encodeFLAC encodes interleaved 16-bit stereo PCM to FLAC and writes it to w.
// The output is untagged (no Vorbis comment block). Input pcm is L,R,L,R…
// (2 channels, 16 bits/sample).
func encodeFLAC(w io.Writer, pcm []int16, sampleRate int) error {
	nsamples := len(pcm) / 2 // samples per channel

	info := &meta.StreamInfo{
		BlockSizeMin:  flacBlockSize,
		BlockSizeMax:  flacBlockSize,
		SampleRate:    uint32(sampleRate),
		NChannels:     2,
		BitsPerSample: 16,
		NSamples:      uint64(nsamples),
	}

	enc, err := flac.NewEncoder(w, info)
	if err != nil {
		return err
	}

	for off := 0; off < nsamples; off += flacBlockSize {
		end := off + flacBlockSize
		if end > nsamples {
			end = nsamples
		}
		n := end - off

		left := make([]int32, n)
		right := make([]int32, n)
		for i := 0; i < n; i++ {
			left[i] = int32(pcm[(off+i)*2])
			right[i] = int32(pcm[(off+i)*2+1])
		}

		f := &frame.Frame{
			Header: frame.Header{
				HasFixedBlockSize: true,
				BlockSize:         uint16(n),
				SampleRate:        uint32(sampleRate),
				Channels:          frame.ChannelsLR,
				BitsPerSample:     16,
			},
			Subframes: []*frame.Subframe{
				{
					SubHeader: frame.SubHeader{Pred: frame.PredVerbatim},
					Samples:   left,
					NSamples:  n,
				},
				{
					SubHeader: frame.SubHeader{Pred: frame.PredVerbatim},
					Samples:   right,
					NSamples:  n,
				},
			},
		}

		if err := enc.WriteFrame(f); err != nil {
			return err
		}
	}

	return enc.Close()
}
