package audio

import "encoding/binary"

// id3v2Tag builds an ID3v2.3 tag (prepended to MP3 frames) from tags.
// Returns nil if nothing to write. Text frames use ISO-8859-1 (encoding 0x00).
func id3v2Tag(t Tags) []byte {
	var frames []byte
	text := func(id, v string) {
		if v == "" {
			return
		}
		payload := append([]byte{0x00}, []byte(v)...)
		frames = append(frames, id3Frame(id, payload)...)
	}
	txxx := func(desc, v string) {
		if v == "" {
			return
		}
		payload := []byte{0x00}
		payload = append(payload, []byte(desc)...)
		payload = append(payload, 0x00)
		payload = append(payload, []byte(v)...)
		frames = append(frames, id3Frame("TXXX", payload)...)
	}
	text("TIT2", t.Title)
	text("TALB", t.Album)
	text("TCON", t.Genre)
	if t.Comment != "" {
		payload := []byte{0x00, 'e', 'n', 'g', 0x00} // enc + lang + empty short-desc
		payload = append(payload, []byte(t.Comment)...)
		frames = append(frames, id3Frame("COMM", payload)...)
	}
	txxx("SOURCE_DISK", t.SourceDisk)
	txxx("SOURCE_FILE", t.SourceFile)
	txxx("SOURCE_SHA256", t.SourceSHA256)
	txxx("SOFTWARE", t.Software)
	if len(frames) == 0 {
		return nil
	}
	hdr := make([]byte, 10)
	copy(hdr, "ID3")
	hdr[3], hdr[4] = 3, 0 // v2.3.0
	hdr[5] = 0            // flags
	putSynchsafe(hdr[6:], len(frames))
	return append(hdr, frames...)
}

// id3Frame builds one ID3v2.3 frame: 4-byte id, 4-byte big-endian size,
// 2-byte flags, payload. (Frame sizes are plain big-endian in v2.3.)
func id3Frame(id string, payload []byte) []byte {
	f := make([]byte, 10, 10+len(payload))
	copy(f, id)
	binary.BigEndian.PutUint32(f[4:], uint32(len(payload)))
	return append(f, payload...)
}

// putSynchsafe writes n as a 28-bit synchsafe integer into b[0:4].
func putSynchsafe(b []byte, n int) {
	b[0] = byte(n >> 21 & 0x7F)
	b[1] = byte(n >> 14 & 0x7F)
	b[2] = byte(n >> 7 & 0x7F)
	b[3] = byte(n & 0x7F)
}
