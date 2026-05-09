package samfile

import (
	"testing"
)

func TestSeverityOrdering(t *testing.T) {
	if !(SeverityCosmetic < SeverityInconsistency &&
		SeverityInconsistency < SeverityStructural &&
		SeverityStructural < SeverityFatal) {
		t.Fatalf("Severity constants out of order: cosmetic=%d inconsistency=%d structural=%d fatal=%d",
			SeverityCosmetic, SeverityInconsistency, SeverityStructural, SeverityFatal)
	}
}

func TestSeverityString(t *testing.T) {
	cases := []struct {
		sev  Severity
		want string
	}{
		{SeverityCosmetic, "cosmetic"},
		{SeverityInconsistency, "inconsistency"},
		{SeverityStructural, "structural"},
		{SeverityFatal, "fatal"},
	}
	for _, c := range cases {
		if got := c.sev.String(); got != c.want {
			t.Errorf("Severity(%d).String() = %q; want %q", c.sev, got, c.want)
		}
	}
}

func TestDialectString(t *testing.T) {
	cases := []struct {
		d    Dialect
		want string
	}{
		{DialectUnknown, "unknown"},
		{DialectSAMDOS1, "samdos1"},
		{DialectSAMDOS2, "samdos2"},
		{DialectMasterDOS, "masterdos"},
	}
	for _, c := range cases {
		if got := c.d.String(); got != c.want {
			t.Errorf("Dialect(%d).String() = %q; want %q", c.d, got, c.want)
		}
	}
}
