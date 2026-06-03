package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/petemoore/samfile/v3"
)

func TestConvertETrackerTuneToWAV(t *testing.T) {
	body, err := os.ReadFile("../../etunes/testdata/m01")
	if err != nil {
		t.Fatal(err)
	}
	fe := &samfile.FileEntry{Type: samfile.FT_CODE}
	opts := convertOptions{
		audioFormat: "wav", loops: 1,
		diskName: "FRED 51.mgt", diskSHA256: "deadbeef", fileName: "e1", version: "test",
	}

	name, ok := convertedName("m01", fe, body, opts)
	if !ok || name != "m01.wav" {
		t.Fatalf("name=%q ok=%v, want m01.wav/true", name, ok)
	}
	out, ok, err := convertBody(body, fe, opts)
	if err != nil || !ok {
		t.Fatalf("convertBody: ok=%v err=%v", ok, err)
	}
	if len(out) < 12 || string(out[0:4]) != "RIFF" || string(out[8:12]) != "WAVE" {
		t.Fatalf("not a WAV file: len=%d", len(out))
	}
	if !bytes.Contains(out, []byte("FRED 51.mgt")) || !bytes.Contains(out, []byte("e1")) {
		t.Fatalf("WAV missing provenance tags (Album/Title)")
	}
	// one loop of m01: intro=0 frames, loop=1760 frames.
	// PCM = 1760 frames * 882 samples/frame * 2 channels * 2 bytes/sample.
	const wantPCM = 1760 * 882 * 2 * 2
	if !bytes.Contains(out, []byte("data")) {
		t.Fatalf("WAV has no 'data' chunk")
	}
	if len(out) < wantPCM {
		t.Fatalf("WAV PCM too small: %d < %d", len(out), wantPCM)
	}
	t.Logf("WAV output size: %d bytes (expected PCM data: %d bytes)", len(out), wantPCM)
}

func TestConvertETrackerTuneDefaultMP3Name(t *testing.T) {
	body, err := os.ReadFile("../../etunes/testdata/m01")
	if err != nil {
		t.Fatal(err)
	}
	fe := &samfile.FileEntry{Type: samfile.FT_CODE}
	opts := convertOptions{audioFormat: "mp3", loops: 4, version: "test"}
	name, ok := convertedName("m01", fe, body, opts)
	if !ok || name != "m01.mp3" {
		t.Fatalf("name=%q ok=%v, want m01.mp3/true", name, ok)
	}
}
