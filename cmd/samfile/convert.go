package main

import (
	"bytes"
	"fmt"

	"github.com/petemoore/samfile/v3"
	"github.com/petemoore/samfile/v3/audio"
	"github.com/petemoore/samfile/v3/etunes"
)

// convertOptions carries audio settings + provenance for -c on E-Tracker tunes.
type convertOptions struct {
	audioFormat string // wav|flac|mp3
	loops       int    // >= 1
	diskName    string // .mgt basename (provenance)
	diskSHA256  string // sha256 of the .mgt image bytes
	fileName    string // the SAM on-disk filename of this file
	version     string // samfile version string
}

// convertedName returns the output filename for a file body that is
// being converted (the -c/--convert option). SAM BASIC files gain a
// ".txt" suffix (they hold detokenised text), SCREEN$ files gain a
// ".png" suffix, E-Tracker modules gain an audio extension, and every
// other type keeps its raw name. The bool reports whether conversion
// applies to fe at all.
func convertedName(base string, fe *samfile.FileEntry, body []byte, opts convertOptions) (string, bool) {
	switch {
	case fe.Type == samfile.FT_SAM_BASIC:
		return base + ".txt", true
	case fe.Type == samfile.FT_SCREEN:
		return base + ".png", true
	case etunes.IsETrackerModule(body):
		return base + "." + audio.Ext(opts.audioFormat), true
	default:
		return base, false
	}
}

// convertBody renders a file body according to its directory type: SAM
// BASIC bodies are detokenised to a plain-text listing, SCREEN$ bodies
// are rendered to a PNG, and E-Tracker modules are rendered to audio.
// Any other type (or a type that cannot be decoded) returns ok=false so
// the caller falls back to the raw body. fe supplies the type and, for
// SCREEN$, the stored MODE byte. opts carries audio settings and
// provenance for E-Tracker tunes.
func convertBody(body []byte, fe *samfile.FileEntry, opts convertOptions) (out []byte, ok bool, err error) {
	switch {
	case fe.Type == samfile.FT_SAM_BASIC:
		var buf bytes.Buffer
		sb := samfile.NewSAMBasic(body)
		if err := sb.WriteText(&buf); err != nil {
			return nil, false, err
		}
		return buf.Bytes(), true, nil
	case fe.Type == samfile.FT_SCREEN:
		var buf bytes.Buffer
		scr := samfile.NewScreen(body, fe.FileTypeInfo[0])
		if err := scr.WritePNG(&buf); err != nil {
			return nil, false, fmt.Errorf("decoding SCREEN$: %w", err)
		}
		return buf.Bytes(), true, nil
	case etunes.IsETrackerModule(body):
		pcm, sr, err := etunes.Render(body, opts.loops)
		if err != nil {
			return nil, false, fmt.Errorf("decoding E-Tracker tune: %w", err)
		}
		tags := audio.Tags{
			Title:        opts.fileName,
			Album:        opts.diskName,
			Genre:        "Chiptune",
			Comment:      "Decoded from SAM Coupe E-Tracker module by samfile " + opts.version,
			SourceDisk:   opts.diskName,
			SourceFile:   opts.fileName,
			SourceSHA256: opts.diskSHA256,
			Software:     "samfile " + opts.version,
		}
		var buf bytes.Buffer
		if err := audio.Encode(&buf, pcm, sr, opts.audioFormat, audio.Options{Tags: tags}); err != nil {
			return nil, false, fmt.Errorf("encoding %s: %w", opts.audioFormat, err)
		}
		return buf.Bytes(), true, nil
	default:
		return nil, false, nil
	}
}
