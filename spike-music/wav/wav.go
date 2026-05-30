// Package wav writes 16-bit PCM stereo WAV files (RIFF) — small and dependency
// free, per the spec's "write the RIFF/WAV header + PCM yourself, it's tiny".
package wav

import (
	"encoding/binary"
	"io"
	"os"
)

// WriteStereo16 writes interleaved L,R int16 samples as a 16-bit stereo WAV at
// the given sample rate.
func WriteStereo16(path string, samples []int16, sampleRate int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return writeStereo16(f, samples, sampleRate)
}

func writeStereo16(w io.Writer, samples []int16, sampleRate int) error {
	const channels = 2
	const bitsPerSample = 16
	dataLen := len(samples) * 2 // bytes
	byteRate := sampleRate * channels * bitsPerSample / 8
	blockAlign := channels * bitsPerSample / 8

	le := binary.LittleEndian
	hdr := make([]byte, 44)
	copy(hdr[0:], "RIFF")
	le.PutUint32(hdr[4:], uint32(36+dataLen))
	copy(hdr[8:], "WAVE")
	copy(hdr[12:], "fmt ")
	le.PutUint32(hdr[16:], 16)       // fmt chunk size
	le.PutUint16(hdr[20:], 1)        // PCM
	le.PutUint16(hdr[22:], channels) //
	le.PutUint32(hdr[24:], uint32(sampleRate))
	le.PutUint32(hdr[28:], uint32(byteRate))
	le.PutUint16(hdr[32:], uint16(blockAlign))
	le.PutUint16(hdr[34:], bitsPerSample)
	copy(hdr[36:], "data")
	le.PutUint32(hdr[40:], uint32(dataLen))
	if _, err := w.Write(hdr); err != nil {
		return err
	}

	buf := make([]byte, dataLen)
	for i, s := range samples {
		le.PutUint16(buf[i*2:], uint16(s))
	}
	_, err := w.Write(buf)
	return err
}
