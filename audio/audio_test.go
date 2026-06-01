package audio

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func tone(n int) []int16 {
	pcm := make([]int16, n*2)
	for i := 0; i < n; i++ {
		v := int16(8000)
		if (i/50)%2 == 1 {
			v = -8000
		}
		pcm[2*i], pcm[2*i+1] = v, v
	}
	return pcm
}

func TestEncodeWAVRoundTrip(t *testing.T) {
	pcm := tone(4410)
	var buf bytes.Buffer
	if err := Encode(&buf, pcm, 44100, "wav", Options{}); err != nil {
		t.Fatal(err)
	}
	b := buf.Bytes()
	if string(b[0:4]) != "RIFF" || string(b[8:12]) != "WAVE" {
		t.Fatalf("bad WAV header")
	}
	// last 4 bytes of the 44-byte canonical header region hold the data size
	// when there is no INFO chunk: locate the "data" chunk and check its size.
	idx := bytes.Index(b, []byte("data"))
	if idx < 0 {
		t.Fatalf("no data chunk")
	}
	got := binary.LittleEndian.Uint32(b[idx+4:])
	if int(got) != len(pcm)*2 {
		t.Fatalf("data len = %d, want %d", got, len(pcm)*2)
	}
}

func TestEncodeFLAC(t *testing.T) {
	pcm := tone(8820)
	var buf bytes.Buffer
	if err := Encode(&buf, pcm, 44100, "flac", Options{}); err != nil {
		t.Fatal(err)
	}
	b := buf.Bytes()
	if len(b) < 8 || string(b[0:4]) != "fLaC" {
		t.Fatalf("missing fLaC magic")
	}
	if len(b) >= len(pcm)*2 {
		t.Fatalf("FLAC not smaller than raw PCM (%d >= %d)", len(b), len(pcm)*2)
	}
}

func TestEncodeUnsupported(t *testing.T) {
	if err := Encode(&bytes.Buffer{}, tone(10), 44100, "m4a", Options{}); err == nil {
		t.Fatal("want error for m4a")
	}
}

func TestWAVInfoChunk(t *testing.T) {
	pcm := tone(100)
	var buf bytes.Buffer
	tags := Tags{Title: "e1", Album: "FRED 51.mgt", Software: "samfile"}
	if err := Encode(&buf, pcm, 44100, "wav", Options{Tags: tags}); err != nil {
		t.Fatal(err)
	}
	b := buf.Bytes()
	if !bytes.Contains(b, []byte("LIST")) || !bytes.Contains(b, []byte("INFO")) ||
		!bytes.Contains(b, []byte("e1")) || !bytes.Contains(b, []byte("FRED 51.mgt")) {
		t.Fatalf("WAV missing LIST/INFO tags")
	}
	if !bytes.Contains(b, []byte("data")) {
		t.Fatalf("WAV missing data chunk")
	}
}
