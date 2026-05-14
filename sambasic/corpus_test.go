//go:build corpus

package sambasic_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/petemoore/samfile/v3"
	"github.com/petemoore/samfile/v3/sambasic"
)

// corpusDir is the user's local SAM disk corpus root. Adjust if your
// path differs.
var corpusDir = filepath.Join(os.Getenv("HOME"), "sam-corpus", "disks")

// truncateAtProgEnd walks the SAM BASIC line-header chain to find the
// true 0xFF program-end sentinel. The naive `bytes.IndexByte(body,
// 0xFF)` is wrong because 0xFF is the prefix of every 2-byte keyword.
// Format: each line is [MSB LSB LenLo LenHi] + body bytes + 0x0D, then
// the next line, terminated by 0xFF.
func truncateAtProgEnd(body []byte) []byte {
	pos := 0
	for pos < len(body) {
		if body[pos] == 0xFF {
			return body[:pos+1]
		}
		if pos+3 >= len(body) {
			return body
		}
		lineLen := int(body[pos+2]) | int(body[pos+3])<<8
		pos += 4 + lineLen
	}
	return body
}

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

	// Populated only for "diverged" outcomes.
	divOffset    int    // first differing byte offset
	divGot       byte   // byte in masked got at divOffset (0 if past end)
	divWant      byte   // byte in masked want at divOffset (0 if past end)
	divGotLen    int    // total length of masked got
	divWantLen   int    // total length of masked want
	divCtxBefore []byte // up to 4 bytes of want immediately before divOffset
	divCtxAfter  []byte // up to 4 bytes of want immediately after divOffset
	divShape     string // shape classification (extra-space-around-keyword etc.)
	divSignature string // full signature string used as the category key
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

	// Categorised report runs regardless of pass/fail and must not
	// mutate test outcome — it writes /tmp/corpus-failures-report.json
	// and logs a ranked breakdown.
	writeCategorisedReport(t, results)

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
	wantBytes := truncateAtProgEnd(f.Body)
	if bytes.Equal(maskProcFnPlaceholders(gotBytes), maskProcFnPlaceholders(wantBytes)) {
		r.outcome = "match"
		return
	}
	r.outcome = "diverged"
	r.detail = fmt.Sprintf("got %d bytes, want %d bytes", len(gotBytes), len(wantBytes))
	mg := maskProcFnPlaceholders(gotBytes)
	mw := maskProcFnPlaceholders(wantBytes)
	classifyDivergence(&r, mg, mw)
	return
}

// isKeywordByte reports whether b is a single-byte SAM BASIC keyword.
// Single-byte keywords occupy 0x85..0xF6 inclusive. The byte 0xFF is the
// prefix of every 2-byte keyword.
func isKeywordByte(b byte) bool {
	return b >= 0x85 && b <= 0xF6
}

// classifyDivergence fills the divergence fields of r given the masked
// got and want byte slices. It is defensive: if either slice is empty
// the result is "length-mismatch" / "other".
func classifyDivergence(r *corpusResult, got, want []byte) {
	r.divGotLen = len(got)
	r.divWantLen = len(want)
	off := firstDiff(got, want)
	r.divOffset = off

	// Context window from want.
	startCtx := off - 4
	if startCtx < 0 {
		startCtx = 0
	}
	r.divCtxBefore = append([]byte(nil), want[startCtx:off]...)
	endCtx := off + 1 + 4
	if off >= len(want) {
		// off may equal len(want) when got is longer; clip.
		r.divCtxAfter = nil
	} else {
		if endCtx > len(want) {
			endCtx = len(want)
		}
		r.divCtxAfter = append([]byte(nil), want[off+1:endCtx]...)
	}

	var g, w byte
	if off < len(got) {
		g = got[off]
	}
	if off < len(want) {
		w = want[off]
	}
	r.divGot = g
	r.divWant = w

	r.divShape = classifyShape(got, want, off)
	r.divSignature = buildSignature(r)
}

// firstDiff returns the index of the first differing byte between a and
// b. If one is a strict prefix of the other, it returns the length of
// the shorter.
func firstDiff(a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}

// classifyShape inspects the byte at off in got vs want (plus a small
// neighbourhood) and returns a coarse shape label. It does not look at
// the exact bytes — only their kinds — so it groups similar failures.
func classifyShape(got, want []byte, off int) string {
	// length-mismatch wins if one side ran out at the diff point.
	gPast := off >= len(got)
	wPast := off >= len(want)
	if gPast || wPast {
		return "length-mismatch"
	}

	g := got[off]
	w := want[off]

	// Helper: is byte at i in s part of a 2-byte 0xFF nn keyword
	// (either position)? i must be in range.
	isFFKeyword := func(s []byte, i int) bool {
		if i < 0 || i >= len(s) {
			return false
		}
		if s[i] == 0xFF {
			return true
		}
		if i > 0 && s[i-1] == 0xFF {
			return true
		}
		return false
	}
	isKW := func(b byte) bool { return isKeywordByte(b) || b == 0xFF }

	// different-FP-bytes: diff lies within the 5 byte FP payload
	// following an 0x0E in either side. Walk back up to 5 bytes
	// looking for 0x0E in want.
	for k := 1; k <= 5; k++ {
		if off-k >= 0 && want[off-k] == 0x0E {
			return "different-FP-bytes"
		}
	}

	// line-header-byte: diff sits inside the 4-byte line header
	// [MSB LSB LenLo LenHi] that follows a CR (0x0D). This catches
	// the bulk of downstream-length-drift symptoms — the real
	// upstream divergence is on the preceding line. Detecting it
	// here keeps those cascade results from drowning out genuine
	// per-byte shape categories.
	if isAtLineHeader(want, off) || isAtLineHeader(got, off) {
		return "line-header-byte"
	}

	// extra-space-around-keyword: one side has 0x20 here, and a
	// keyword byte sits immediately on the other side of the
	// insertion in want or got.
	if g == 0x20 && w != 0x20 {
		// got inserted a space; the matching position in want is w,
		// and the space is "extra" before/after the keyword run.
		if isKW(w) || (off > 0 && isKW(want[off-1])) || isFFKeyword(want, off) {
			return "extra-space-around-keyword"
		}
		if w == 0x3A || (off > 0 && want[off-1] == 0x3A) {
			return "extra-space-around-colon"
		}
	}
	if w == 0x20 && g != 0x20 {
		// got dropped a space; the want side has 0x20 next to a
		// keyword / colon.
		if isKW(g) || (off > 0 && isKW(got[off-1])) || isFFKeyword(got, off) {
			return "extra-space-around-keyword"
		}
		if g == 0x3A || (off > 0 && got[off-1] == 0x3A) {
			return "extra-space-around-colon"
		}
	}

	// different-keyword-byte: both sides hold a keyword byte but a
	// different one.
	if isKW(g) && isKW(w) {
		return "different-keyword-byte"
	}

	// literal-vs-keyword-byte: one side has an ASCII literal byte
	// where the other has a keyword byte.
	gKW := isKW(g)
	wKW := isKW(w)
	gASCII := g >= 0x20 && g < 0x7F
	wASCII := w >= 0x20 && w < 0x7F
	if (gKW && wASCII) || (wKW && gASCII) {
		return "literal-vs-keyword-byte"
	}

	// "other" covers the long tail. We deliberately do NOT fall back
	// to "length-mismatch" here just because total lengths differ:
	// the byte pair at the diff point is the more actionable signal,
	// and total length mismatches usually follow from a single early
	// divergence rather than indicating their own bug class.
	return "other"
}

// isAtLineHeader reports whether index off in s sits inside the
// 4-byte line header [MSB LSB LenLo LenHi] that immediately follows a
// 0x0D line terminator. It walks back up to 4 bytes looking for a 0x0D
// and checks that the bytes between it and off look like the start of
// a header.
func isAtLineHeader(s []byte, off int) bool {
	if off < 1 || off >= len(s) {
		return false
	}
	for k := 1; k <= 4; k++ {
		j := off - k
		if j < 0 {
			break
		}
		if s[j] == 0x0D {
			// k=1: off is MSB; k=2: off is LSB; k=3: off is LenLo;
			// k=4: off is LenHi. In every case off is part of the
			// header.
			return true
		}
	}
	return false
}

// buildSignature combines the shape and a small fingerprint of the
// context bytes into a stable category key. The context bytes are
// rendered as their "kind" (KW, FF, FP, SP, COLON, NL, CR, NUM, ALPHA,
// or hex) so that signatures group similar but not byte-identical
// diffs.
func buildSignature(r *corpusResult) string {
	var sb strings.Builder
	sb.WriteString(r.divShape)
	sb.WriteString(" | got=")
	sb.WriteString(kindOf(r.divGot))
	sb.WriteString(" want=")
	sb.WriteString(kindOf(r.divWant))
	sb.WriteString(" | ctx=")
	for _, b := range r.divCtxBefore {
		sb.WriteString(kindOf(b))
	}
	sb.WriteString(">")
	for _, b := range r.divCtxAfter {
		sb.WriteString(kindOf(b))
	}
	return sb.String()
}

// kindOf maps a byte to a short symbolic kind, used to build stable
// category signatures.
func kindOf(b byte) string {
	switch b {
	case 0x00:
		return "NUL"
	case 0x0D:
		return "CR"
	case 0x0E:
		return "FP"
	case 0x20:
		return "SP"
	case 0x3A:
		return "COLON"
	case 0xFF:
		return "FF"
	}
	if isKeywordByte(b) {
		return "KW"
	}
	if b >= '0' && b <= '9' {
		return "NUM"
	}
	if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
		return "ALPHA"
	}
	if b >= 0x20 && b < 0x7F {
		return fmt.Sprintf("%q", rune(b))
	}
	return fmt.Sprintf("0x%02X", b)
}

// errorCategory returns a short fingerprint of a parse/detok error
// message. It strips known wrapper prefixes ("basic-to-text: ", etc.)
// and location prefixes of the form "line N, col M: ", then keeps the
// first phrase (up to the first `: ` or `,`) and replaces numeric
// literals with placeholders so messages differing only in numbers
// group together.
func errorCategory(msg string) string {
	s := strings.TrimSpace(msg)
	// Peel off known wrapper prefixes that just identify the
	// subsystem and carry no categorical signal.
	for _, p := range []string{"basic-to-text: ", "sambasic: ", "parse: "} {
		if strings.HasPrefix(s, p) {
			s = s[len(p):]
		}
	}
	// Strip a leading "line N, col M:" location, which varies per
	// file but does not change the kind of error.
	if strings.HasPrefix(s, "line ") {
		if i := strings.Index(s, ":"); i >= 0 {
			// Be permissive: only strip if everything before ':' is
			// "line N" or "line N, col M".
			loc := s[:i]
			rest := strings.TrimSpace(s[i+1:])
			if locLooksLikePosition(loc) {
				s = rest
			}
		}
	}
	cut := len(s)
	if i := strings.Index(s, ": "); i >= 0 && i < cut {
		cut = i
	}
	if i := strings.Index(s, ","); i >= 0 && i < cut {
		cut = i
	}
	phrase := strings.TrimSpace(s[:cut])
	// Replace numeric literals with placeholders.
	var b strings.Builder
	i := 0
	for i < len(phrase) {
		c := phrase[i]
		if i+1 < len(phrase) && c == '0' && (phrase[i+1] == 'x' || phrase[i+1] == 'X') {
			b.WriteString("0xNN")
			i += 2
			for i < len(phrase) {
				cc := phrase[i]
				if (cc >= '0' && cc <= '9') || (cc >= 'a' && cc <= 'f') || (cc >= 'A' && cc <= 'F') {
					i++
					continue
				}
				break
			}
			continue
		}
		if c >= '0' && c <= '9' {
			b.WriteString("N")
			for i < len(phrase) && phrase[i] >= '0' && phrase[i] <= '9' {
				i++
			}
			continue
		}
		b.WriteByte(c)
		i++
	}
	return b.String()
}

// locLooksLikePosition reports whether s is of the form "line N" or
// "line N, col M" (digits only). Used by errorCategory to recognise
// safely strippable location prefixes.
func locLooksLikePosition(s string) bool {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "line ") {
		return false
	}
	rest := strings.TrimSpace(s[5:])
	// Optional ", col M" suffix.
	if comma := strings.Index(rest, ","); comma >= 0 {
		num := strings.TrimSpace(rest[:comma])
		tail := strings.TrimSpace(rest[comma+1:])
		if !allDigits(num) {
			return false
		}
		if !strings.HasPrefix(tail, "col ") {
			return false
		}
		return allDigits(strings.TrimSpace(tail[4:]))
	}
	return allDigits(rest)
}

func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// hexCtx returns a hex/ASCII rendering of a byte slice for inclusion in
// human-readable reports.
func hexCtx(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	parts := make([]string, len(b))
	for i, x := range b {
		parts[i] = fmt.Sprintf("%02X", x)
	}
	return strings.Join(parts, " ")
}

// categoryEntry captures one row of the JSON report.
type categoryEntry struct {
	Outcome   string                   `json:"outcome"`
	Signature string                   `json:"signature"`
	Count     int                      `json:"count"`
	Examples  []map[string]interface{} `json:"examples"`
}

// writeCategorisedReport writes /tmp/corpus-failures-report.json and
// logs a human-readable ranked summary via t.Logf. It must not mutate
// the test pass/fail state — it is purely informational.
func writeCategorisedReport(t *testing.T, results []corpusResult) {
	t.Helper()

	type bucket struct {
		outcome   string
		signature string
		entries   []*corpusResult
	}

	buckets := map[string]*bucket{}
	for i := range results {
		r := &results[i]
		var sig string
		switch r.outcome {
		case "diverged":
			sig = r.divSignature
		case "parse-error", "detok-error":
			sig = errorCategory(r.detail)
		default:
			continue
		}
		key := r.outcome + "::" + sig
		b, ok := buckets[key]
		if !ok {
			b = &bucket{outcome: r.outcome, signature: sig}
			buckets[key] = b
		}
		b.entries = append(b.entries, r)
	}

	ordered := make([]*bucket, 0, len(buckets))
	for _, b := range buckets {
		ordered = append(ordered, b)
	}
	sort.Slice(ordered, func(i, j int) bool {
		if len(ordered[i].entries) != len(ordered[j].entries) {
			return len(ordered[i].entries) > len(ordered[j].entries)
		}
		if ordered[i].outcome != ordered[j].outcome {
			return ordered[i].outcome < ordered[j].outcome
		}
		return ordered[i].signature < ordered[j].signature
	})

	report := make([]categoryEntry, 0, len(ordered))
	for _, b := range ordered {
		entry := categoryEntry{
			Outcome:   b.outcome,
			Signature: b.signature,
			Count:     len(b.entries),
		}
		n := 3
		if n > len(b.entries) {
			n = len(b.entries)
		}
		for _, r := range b.entries[:n] {
			ex := map[string]interface{}{
				"disk": r.disk,
				"file": r.file,
			}
			if r.outcome == "diverged" {
				ex["offset"] = r.divOffset
				ex["got_byte"] = fmt.Sprintf("0x%02X", r.divGot)
				ex["want_byte"] = fmt.Sprintf("0x%02X", r.divWant)
				ex["got_len"] = r.divGotLen
				ex["want_len"] = r.divWantLen
				ex["ctx_before"] = hexCtx(r.divCtxBefore)
				ex["ctx_after"] = hexCtx(r.divCtxAfter)
				ex["shape"] = r.divShape
			} else {
				ex["detail"] = r.detail
			}
			entry.Examples = append(entry.Examples, ex)
		}
		report = append(report, entry)
	}

	// Write JSON report. Failure to write is logged but does not fail
	// the test.
	const reportPath = "/tmp/corpus-failures-report.json"
	if data, err := json.MarshalIndent(report, "", "  "); err != nil {
		t.Logf("marshal report: %v", err)
	} else if err := os.WriteFile(reportPath, data, 0o644); err != nil {
		t.Logf("write report: %v", err)
	} else {
		t.Logf("wrote categorised failure report to %s", reportPath)
	}

	// Human-readable summary, ordered by count desc.
	totalFail := 0
	for _, b := range ordered {
		totalFail += len(b.entries)
	}
	t.Logf("=== categorised failure summary (%d failure entries across %d categories) ===", totalFail, len(ordered))
	for i, b := range ordered {
		pct := 0.0
		if totalFail > 0 {
			pct = 100.0 * float64(len(b.entries)) / float64(totalFail)
		}
		t.Logf("#%-3d [%-12s] %5d (%5.1f%%) %s", i+1, b.outcome, len(b.entries), pct, b.signature)
		max := 3
		if max > len(b.entries) {
			max = len(b.entries)
		}
		for _, r := range b.entries[:max] {
			if r.outcome == "diverged" {
				t.Logf("        ex: %s/%s @%d got=0x%02X want=0x%02X ctx=[%s]>[%s] (got_len=%d want_len=%d)",
					r.disk, r.file, r.divOffset, r.divGot, r.divWant,
					hexCtx(r.divCtxBefore), hexCtx(r.divCtxAfter),
					r.divGotLen, r.divWantLen)
			} else {
				t.Logf("        ex: %s/%s -- %s", r.disk, r.file, r.detail)
			}
		}
	}
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
