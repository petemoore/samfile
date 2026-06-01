// Package audio encodes 16-bit stereo PCM (interleaved L,R) to audio files.
// All encoders are pure Go (no cgo, no external binaries).
package audio

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Tags is provenance metadata written into the output where the container
// supports it (MP3 ID3v2, WAV LIST/INFO; FLAC untagged for now).
type Tags struct {
	Title, Album, Genre, Comment         string
	SourceDisk, SourceFile, SourceSHA256 string
	Software                             string
}

// Options tunes optional encoder parameters.
type Options struct {
	MP3BitrateKbps int  // MP3 only; 0 ⇒ default 256
	Tags           Tags // provenance metadata
}

// Supported lists the format names Encode accepts.
func Supported() []string { return []string{"wav", "flac", "mp3"} }

// Ext returns the file extension (without dot) for a supported format.
func Ext(format string) string { return format }

// Encode writes interleaved 16-bit stereo PCM to w in the named format.
func Encode(w io.Writer, pcm []int16, sampleRate int, format string, opts Options) error {
	switch format {
	case "wav":
		return encodeWAV(w, pcm, sampleRate, opts.Tags)
	case "flac":
		return encodeFLAC(w, pcm, sampleRate) // untagged for now
	case "mp3":
		return encodeMP3(w, pcm, sampleRate, opts.MP3BitrateKbps, opts.Tags)
	default:
		return fmt.Errorf("unsupported audio format %q (supported: wav, flac, mp3; for m4a/opus/etc. emit wav and pipe to ffmpeg)", format)
	}
}

func encodeWAV(w io.Writer, pcm []int16, sampleRate int, tags Tags) error {
	const channels, bits = 2, 16
	dataLen := len(pcm) * 2
	info := wavInfoChunk(tags) // "LIST"+size+"INFO"+subchunks, or nil
	le := binary.LittleEndian
	riffSize := 4 + (8 + 16) + len(info) + (8 + dataLen) // WAVE + fmt + LIST + data
	var h []byte
	put := func(b []byte) { h = append(h, b...) }
	u32 := func(v uint32) { var t [4]byte; le.PutUint32(t[:], v); h = append(h, t[:]...) }
	u16 := func(v uint16) { var t [2]byte; le.PutUint16(t[:], v); h = append(h, t[:]...) }
	put([]byte("RIFF"))
	u32(uint32(riffSize))
	put([]byte("WAVE"))
	put([]byte("fmt "))
	u32(16)
	u16(1)
	u16(channels)
	u32(uint32(sampleRate))
	u32(uint32(sampleRate * channels * bits / 8))
	u16(channels * bits / 8)
	u16(bits)
	put(info) // may be empty
	put([]byte("data"))
	u32(uint32(dataLen))
	if _, err := w.Write(h); err != nil {
		return err
	}
	buf := make([]byte, dataLen)
	for i, s := range pcm {
		le.PutUint16(buf[i*2:], uint16(s))
	}
	_, err := w.Write(buf)
	return err
}

// wavInfoChunk builds a RIFF LIST/INFO chunk from tags, or nil if all empty.
func wavInfoChunk(t Tags) []byte {
	type kv struct{ id, val string }
	items := []kv{
		{"INAM", t.Title}, {"IPRD", t.Album}, {"IGNR", t.Genre},
		{"ICMT", t.Comment}, {"ISFT", t.Software},
	}
	var body []byte
	le := binary.LittleEndian
	for _, it := range items {
		if it.val == "" {
			continue
		}
		v := append([]byte(it.val), 0) // null-terminated
		if len(v)%2 == 1 {
			v = append(v, 0) // pad to even
		}
		var sz [4]byte
		le.PutUint32(sz[:], uint32(len(v)))
		body = append(body, []byte(it.id)...)
		body = append(body, sz[:]...)
		body = append(body, v...)
	}
	if len(body) == 0 {
		return nil
	}
	out := []byte("LIST")
	var sz [4]byte
	le.PutUint32(sz[:], uint32(4+len(body))) // "INFO" + body
	out = append(out, sz[:]...)
	out = append(out, []byte("INFO")...)
	return append(out, body...)
}
