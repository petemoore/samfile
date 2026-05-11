package sambasic

import (
	"bytes"
	"testing"
)

func TestParseRoundTrip(t *testing.T) {
	original := &File{
		Lines: []Line{
			{
				Number: 10,
				Tokens: []Token{
					CLEAR,
					Number(32767),
					Literal(':'),
					LOAD,
					Literal('"'),
					String("stub"),
					Literal('"'),
					CODE,
					Number(32768),
				},
			},
		},
		StartLine: 10,
	}

	body := original.Bytes()

	parsed, err := Parse(body)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	parsed.NumericVars = original.NumericVars
	parsed.Gap = original.Gap
	parsed.StringArrayVars = original.StringArrayVars
	parsed.StartLine = original.StartLine

	got := parsed.Bytes()
	if !bytes.Equal(got, body) {
		t.Errorf("roundtrip mismatch: len(got)=%d, len(want)=%d", len(got), len(body))
	}
}
