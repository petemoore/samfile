package audio

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/mewkiz/flac"
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

func TestEncodeFLACLossless(t *testing.T) {
	pcm := tone(8820)
	var buf bytes.Buffer
	if err := Encode(&buf, pcm, 44100, "flac", Options{}); err != nil {
		t.Fatal(err)
	}

	stream, err := flac.New(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("flac.New: %v", err)
	}
	defer stream.Close()

	var got []int16
	for {
		f, err := stream.ParseNext()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("ParseNext: %v", err)
		}
		if len(f.Subframes) != 2 {
			t.Fatalf("frame has %d subframes, want 2", len(f.Subframes))
		}
		left := f.Subframes[0].Samples
		right := f.Subframes[1].Samples
		for i := range left {
			got = append(got, int16(left[i]), int16(right[i]))
		}
	}

	if len(got) != len(pcm) {
		t.Fatalf("decoded %d samples, want %d", len(got), len(pcm))
	}
	for i := range pcm {
		if got[i] != pcm[i] {
			t.Fatalf("sample[%d]: decoded %d, want %d (not lossless)", i, got[i], pcm[i])
		}
	}
}

func TestEncodeMP3(t *testing.T) {
	if raceEnabled {
		// shine-mp3's l3subband uses unsafe pointer arithmetic that trips Go's
		// -race/checkptr instrumentation (fatal "bad pointer"). MP3 output is
		// correct in normal builds; only the checkptr-instrumented build crashes.
		t.Skip("shine-mp3 uses unsafe pointer arithmetic incompatible with -race/checkptr")
	}
	pcm := tone(44100) // 1s
	var buf bytes.Buffer
	tags := Tags{Title: "e1", Album: "FRED 51.mgt", SourceSHA256: "abc123", Software: "samfile"}
	if err := Encode(&buf, pcm, 44100, "mp3", Options{Tags: tags}); err != nil {
		t.Fatal(err)
	}
	b := buf.Bytes()
	if len(b) < 200 {
		t.Fatalf("MP3 suspiciously small: %d bytes", len(b))
	}
	if string(b[0:3]) != "ID3" {
		t.Fatalf("missing ID3v2 tag")
	}
	if !bytes.Contains(b, []byte("e1")) || !bytes.Contains(b, []byte("abc123")) {
		t.Fatalf("ID3 missing expected tag values")
	}
	sync := false
	for i := 0; i+1 < len(b); i++ {
		if b[i] == 0xFF && b[i+1]&0xE0 == 0xE0 {
			sync = true
			break
		}
	}
	if !sync {
		t.Fatalf("no MP3 frame sync found")
	}
	t.Logf("MP3 output size: %d bytes", len(b))
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
