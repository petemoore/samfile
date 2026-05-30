package main

import (
	"bytes"
	"fmt"

	"github.com/petemoore/samfile/v3"
)

// convertedName returns the output filename for a file body that is
// being converted (the -c/--convert option). SAM BASIC files gain a
// ".txt" suffix (they hold detokenised text), SCREEN$ files gain a
// ".png" suffix, and every other type keeps its raw name. The bool
// reports whether conversion applies to fe at all.
func convertedName(base string, fe *samfile.FileEntry) (string, bool) {
	switch fe.Type {
	case samfile.FT_SAM_BASIC:
		return base + ".txt", true
	case samfile.FT_SCREEN:
		return base + ".png", true
	default:
		return base, false
	}
}

// convertBody renders a file body according to its directory type: SAM
// BASIC bodies are detokenised to a plain-text listing and SCREEN$
// bodies are rendered to a PNG. Any other type (or a type that cannot
// be decoded) returns ok=false so the caller falls back to the raw
// body. fe supplies the type and, for SCREEN$, the stored MODE byte.
func convertBody(body []byte, fe *samfile.FileEntry) (out []byte, ok bool, err error) {
	switch fe.Type {
	case samfile.FT_SAM_BASIC:
		var buf bytes.Buffer
		sb := samfile.NewSAMBasic(body)
		if err := sb.WriteText(&buf); err != nil {
			return nil, false, err
		}
		return buf.Bytes(), true, nil
	case samfile.FT_SCREEN:
		var buf bytes.Buffer
		scr := samfile.NewScreen(body, fe.FileTypeInfo[0])
		if err := scr.WritePNG(&buf); err != nil {
			return nil, false, fmt.Errorf("decoding SCREEN$: %w", err)
		}
		return buf.Bytes(), true, nil
	default:
		return nil, false, nil
	}
}
