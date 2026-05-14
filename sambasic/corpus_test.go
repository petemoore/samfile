//go:build corpus

package sambasic_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/petemoore/samfile/v3"
	"github.com/petemoore/samfile/v3/sambasic"
)

// corpusDir is the user's local SAM disk corpus root. Adjust if your
// path differs.
var corpusDir = filepath.Join(os.Getenv("HOME"), "sam-corpus", "disks")

// maskProcFnPlaceholders returns a copy of body with bytes 3-5 of every
// `0E FD FD ??` and `0E FE FE ??` 6-byte buffer zeroed. These bytes are
// re-patched by LDPROG at LOAD time (grammar spec §6.5) and are not
// reproducible from text input alone.
func maskProcFnPlaceholders(body []byte) []byte {
	out := make([]byte, len(body))
	copy(out, body)
	for i := 0; i+5 < len(out); i++ {
		if out[i] != 0x0E {
			continue
		}
		if out[i+1] == 0xFD && out[i+2] == 0xFD {
			out[i+3] = 0x00
			out[i+4] = 0x00
			out[i+5] = 0x00
			i += 5
			continue
		}
		if out[i+1] == 0xFE && out[i+2] == 0xFE {
			out[i+3] = 0x00
			out[i+4] = 0x00
			out[i+5] = 0x00
			i += 5
		}
	}
	return out
}

type corpusResult struct {
	disk    string
	file    string
	outcome string
	detail  string
}

func TestCorpusRoundTrip(t *testing.T) {
	if _, err := os.Stat(corpusDir); err != nil {
		t.Skipf("corpus dir %s not present: %v", corpusDir, err)
	}
	disks, err := filepath.Glob(filepath.Join(corpusDir, "*.mgt"))
	if err != nil {
		t.Fatal(err)
	}
	if len(disks) == 0 {
		t.Skipf("no .mgt files in %s", corpusDir)
	}

	var results []corpusResult
	for _, diskPath := range disks {
		diskName := filepath.Base(diskPath)
		di, err := samfile.Load(diskPath)
		if err != nil {
			t.Logf("load %s: %v", diskPath, err)
			continue
		}
		for _, fe := range di.DiskJournal() {
			if !fe.Used() || fe.Type != samfile.FT_SAM_BASIC {
				continue
			}
			results = append(results, roundTripOne(di, fe, diskName))
		}
	}

	counts := map[string]int{}
	for _, r := range results {
		counts[r.outcome]++
	}
	t.Logf("corpus summary: %v over %d files", counts, len(results))

	if counts["diverged"] > 0 || counts["parse-error"] > 0 {
		shown := map[string]int{}
		for _, r := range results {
			if r.outcome == "match" {
				continue
			}
			if shown[r.outcome] >= 10 {
				continue
			}
			shown[r.outcome]++
			t.Errorf("[%s] %s/%s: %s", r.outcome, r.disk, r.file, r.detail)
		}
	}
}

// roundTripOne processes a single FT_SAM_BASIC entry, recovering from
// any panic in the underlying samfile/sambasic code so one corrupt
// disk does not abort the whole corpus walk.
func roundTripOne(di *samfile.DiskImage, fe *samfile.FileEntry, diskName string) (r corpusResult) {
	name := fe.Name.String()
	r = corpusResult{disk: diskName, file: name}
	defer func() {
		if p := recover(); p != nil {
			r.outcome = "panic"
			r.detail = fmt.Sprintf("%v", p)
		}
	}()
	f, err := di.File(name)
	if err != nil {
		r.outcome = "file-error"
		r.detail = err.Error()
		return
	}
	text, err := bodyToText(f.Body)
	if err != nil {
		r.outcome = "detok-error"
		r.detail = err.Error()
		return
	}
	got, err := sambasic.ParseTextString(text)
	if err != nil {
		r.outcome = "parse-error"
		r.detail = err.Error()
		return
	}
	gotBytes := got.ProgBytes()
	wantBytes := f.Body
	if idx := bytes.IndexByte(wantBytes, 0xFF); idx >= 0 {
		wantBytes = wantBytes[:idx+1]
	}
	if bytes.Equal(maskProcFnPlaceholders(gotBytes), maskProcFnPlaceholders(wantBytes)) {
		r.outcome = "match"
		return
	}
	r.outcome = "diverged"
	r.detail = fmt.Sprintf("got %d bytes, want %d bytes", len(gotBytes), len(wantBytes))
	return
}

// bodyToText converts an FT_SAM_BASIC body to its plain-text listing
// via the existing SAMBasic.Output() path. Output() writes to stdout;
// capture via an in-process pipe.
func bodyToText(body []byte) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	stdout := os.Stdout
	os.Stdout = w

	// Run Output() in a goroutine so we can read the pipe concurrently
	// without filling its OS buffer (~64KB typically).
	type result struct {
		err error
	}
	done := make(chan result, 1)
	go func() {
		sb := samfile.NewSAMBasic(body)
		err := sb.Output()
		w.Close()
		done <- result{err: err}
	}()

	var buf bytes.Buffer
	_, copyErr := buf.ReadFrom(r)
	res := <-done

	os.Stdout = stdout
	if res.err != nil {
		return "", res.err
	}
	if copyErr != nil {
		return "", copyErr
	}
	return buf.String(), nil
}
