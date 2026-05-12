package samfile

import (
	"testing"

	"github.com/petemoore/samfile/v3/sambasic"
)

// buildBasicDisk returns a samfile-built disk containing one BASIC
// program with one line (10 REM "hi") and auto-RUN at line 10. The
// returned dj is the journal at construction time; callers can
// mutate slot 0 and call di.WriteFileEntry(dj, 0) to test
// negative cases.
//
// The defaults produce a SAMDOS-2-canonical disk: NumericVars=92
// bytes + Gap=512 bytes = SAVARS-NVARS=604 (sambasic/file.go:3-6).
func buildBasicDisk(t *testing.T) (*DiskImage, *DiskJournal) {
	t.Helper()
	bf := &sambasic.File{
		StartLine: 10,
		Lines: []sambasic.Line{
			{Number: 10, Tokens: []sambasic.Token{sambasic.REM, sambasic.String("hi")}},
		},
	}
	di := NewDiskImage()
	if err := di.AddBasicFile("DEMO", bf); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}
	return di, di.DiskJournal()
}

// buildBasicDiskNoAutoRun builds a BASIC disk where StartLine is
// 0xFFFF (no auto-RUN). Used by BASIC-STARTLINE-* rule tests.
func buildBasicDiskNoAutoRun(t *testing.T) (*DiskImage, *DiskJournal) {
	t.Helper()
	bf := &sambasic.File{
		StartLine: 0xFFFF,
		Lines: []sambasic.Line{
			{Number: 10, Tokens: []sambasic.Token{sambasic.REM, sambasic.String("hi")}},
		},
	}
	di := NewDiskImage()
	if err := di.AddBasicFile("DEMO", bf); err != nil {
		t.Fatalf("AddBasicFile (no auto-RUN): %v", err)
	}
	return di, di.DiskJournal()
}

// ----- BASIC-FILETYPEINFO-TRIPLETS -----

func TestBasicFileTypeInfoTripletsPositive(t *testing.T) {
	di, _ := buildBasicDisk(t)
	findings := checkBasicFileTypeInfoTriplets(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("clean BASIC disk: %d findings; want 0", len(findings))
	}
}

func TestBasicFileTypeInfoTripletsNegative(t *testing.T) {
	di, dj := buildBasicDisk(t)
	// Zero out FileTypeInfo so the triplets decode to all zero.
	for i := range dj[0].FileTypeInfo {
		dj[0].FileTypeInfo[i] = 0
	}
	di.WriteFileEntry(dj, 0)
	findings := checkBasicFileTypeInfoTriplets(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BASIC-FILETYPEINFO-TRIPLETS" {
		t.Fatalf("got %d findings, first=%+v; want 1 BASIC-FILETYPEINFO-TRIPLETS", len(findings), findings)
	}
}

// ----- BASIC-VARS-GAP-INVARIANT -----

func TestBasicVarsGapInvariantSAMDOS2Clean(t *testing.T) {
	di, _ := buildBasicDisk(t)
	// AddBasicFile defaults: SAMDOS-2 gap is 604. ctx.Dialect = SAMDOS2 → no finding.
	findings := checkBasicVarsGapInvariant(&CheckContext{
		Disk: di, Journal: di.DiskJournal(), Dialect: DialectSAMDOS2,
	})
	if len(findings) != 0 {
		t.Errorf("SAMDOS-2 clean BASIC: %d findings; want 0", len(findings))
	}
}

func TestBasicVarsGapInvariantSAMDOS2BadGap(t *testing.T) {
	// Pass a non-canonical NumericVars length so the gap (NumericVars
	// + Gap) becomes 92+1 + 512 = 605 (≠ 604, ≠ 2156). Under SAMDOS-2
	// dialect, the rule fires.
	bf := &sambasic.File{
		StartLine: 10,
		Lines: []sambasic.Line{
			{Number: 10, Tokens: []sambasic.Token{sambasic.REM, sambasic.String("hi")}},
		},
		NumericVars: make([]byte, 93), // default+1 → gap = 605
	}
	di := NewDiskImage()
	if err := di.AddBasicFile("DEMO", bf); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}
	findings := checkBasicVarsGapInvariant(&CheckContext{
		Disk: di, Journal: di.DiskJournal(), Dialect: DialectSAMDOS2,
	})
	if len(findings) != 1 || findings[0].RuleID != "BASIC-VARS-GAP-INVARIANT" {
		t.Fatalf("got %d findings, first=%+v; want 1 BASIC-VARS-GAP-INVARIANT", len(findings), findings)
	}
}

func TestBasicVarsGapInvariantMasterDOSClean(t *testing.T) {
	// Pass Gap=2064 so NumericVars+Gap = 92 + 2064 = 2156 (MasterDOS
	// canonical). Under MasterDOS dialect, no finding.
	bf := &sambasic.File{
		StartLine: 10,
		Lines: []sambasic.Line{
			{Number: 10, Tokens: []sambasic.Token{sambasic.REM, sambasic.String("hi")}},
		},
		Gap: make([]byte, 2064), // 92 default NumericVars + 2064 Gap = 2156
	}
	di := NewDiskImage()
	if err := di.AddBasicFile("DEMO", bf); err != nil {
		t.Fatalf("AddBasicFile: %v", err)
	}
	findings := checkBasicVarsGapInvariant(&CheckContext{
		Disk: di, Journal: di.DiskJournal(), Dialect: DialectMasterDOS,
	})
	if len(findings) != 0 {
		t.Errorf("MasterDOS clean (gap=2156): %d findings; want 0", len(findings))
	}
}

func TestBasicVarsGapInvariantUnknownDialect(t *testing.T) {
	di, _ := buildBasicDisk(t)
	// Default gap is 604 (SAMDOS-2 canonical). Under Unknown dialect
	// the rule accepts both 604 and 2156 silently.
	findings := checkBasicVarsGapInvariant(&CheckContext{
		Disk: di, Journal: di.DiskJournal(), Dialect: DialectUnknown,
	})
	if len(findings) != 0 {
		t.Errorf("Unknown dialect + canonical gap: %d findings; want 0", len(findings))
	}
}

// ----- BASIC-PROG-END-SENTINEL -----

func TestBasicProgEndSentinelPositive(t *testing.T) {
	di, _ := buildBasicDisk(t)
	findings := checkBasicProgEndSentinel(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("clean BASIC disk: %d findings; want 0", len(findings))
	}
}

func TestBasicProgEndSentinelNegative(t *testing.T) {
	// buildBasicDisk produces: 10 REM "hi" = 9 bytes program area.
	// sentinel is at body[progLen-1] = body[8] = raw[17] = sd[17].
	di, _ := buildBasicDisk(t)
	dj := di.DiskJournal()
	progLen := dj[0].ProgramLength() // should be 9
	// sentinel lives in first sector at sd[9 + progLen - 1] = sd[9+progLen-1]
	sentinelOffset := 9 + int(progLen) - 1
	mutateFirstSectorByte(t, di, sentinelOffset, 0xAA)
	findings := checkBasicProgEndSentinel(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BASIC-PROG-END-SENTINEL" {
		t.Fatalf("got %d findings, first=%+v; want 1 BASIC-PROG-END-SENTINEL", len(findings), findings)
	}
}

// ----- BASIC-LINE-NUMBER-BE -----

func TestBasicLineNumberBEPositive(t *testing.T) {
	di, _ := buildBasicDisk(t)
	findings := checkBasicLineNumberBE(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("clean BASIC disk: %d findings; want 0", len(findings))
	}
}

func TestBasicLineNumberBENegative(t *testing.T) {
	// Corrupt body bytes 9..12 (the first line's 4-byte header).
	// Setting length to 0xFF/0xFF makes the parser see a line body
	// extending past the program area → parse error.
	di, _ := buildBasicDisk(t)
	mutateFirstSectorByte(t, di, 11, 0xFF) // lineLen high byte → huge length
	mutateFirstSectorByte(t, di, 10, 0xFF) // lineLen low byte
	findings := checkBasicLineNumberBE(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BASIC-LINE-NUMBER-BE" {
		t.Fatalf("got %d findings, first=%+v; want 1 BASIC-LINE-NUMBER-BE", len(findings), findings)
	}
}

// ----- BASIC-STARTLINE-FF-DISABLES -----

func TestBasicStartLineFFDisablesPositive(t *testing.T) {
	di, _ := buildBasicDisk(t)
	findings := checkBasicStartLineFFDisables(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("clean BASIC disk (auto-RUN): %d findings; want 0", len(findings))
	}
}

func TestBasicStartLineFFDisablesNoAutoRunPositive(t *testing.T) {
	di, _ := buildBasicDiskNoAutoRun(t)
	findings := checkBasicStartLineFFDisables(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("BASIC disk (no auto-RUN): %d findings; want 0", len(findings))
	}
}

func TestBasicStartLineFFDisablesNegative(t *testing.T) {
	di, dj := buildBasicDisk(t)
	dj[0].ExecutionAddressDiv16K = 0x42 // neither 0x00 nor 0xFF
	di.WriteFileEntry(dj, 0)
	findings := checkBasicStartLineFFDisables(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BASIC-STARTLINE-FF-DISABLES" {
		t.Fatalf("got %d findings, first=%+v; want 1 BASIC-STARTLINE-FF-DISABLES", len(findings), findings)
	}
}

// ----- BASIC-STARTLINE-WITHIN-PROG -----

func TestBasicStartLineWithinProgPositive(t *testing.T) {
	di, _ := buildBasicDisk(t) // line 10 exists; StartLine=10
	findings := checkBasicStartLineWithinProg(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("clean BASIC disk: %d findings; want 0", len(findings))
	}
}

func TestBasicStartLineWithinProgNegative(t *testing.T) {
	di, dj := buildBasicDisk(t)
	// SAMBASICStartLine and ExecutionAddressMod16K share the same on-disk
	// bytes 0xF3-0xF4. Raw() serialises ExecutionAddressMod16K, so we must
	// set both fields to make WriteFileEntry persist the change correctly.
	dj[0].SAMBASICStartLine = 99  // line 99 not in program (only line 10 exists)
	dj[0].ExecutionAddressMod16K = 99
	di.WriteFileEntry(dj, 0)
	findings := checkBasicStartLineWithinProg(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BASIC-STARTLINE-WITHIN-PROG" {
		t.Fatalf("got %d findings, first=%+v; want 1 BASIC-STARTLINE-WITHIN-PROG", len(findings), findings)
	}
}

// ----- BASIC-MGTFLAGS-20 -----

func TestBasicMGTFlags20Positive(t *testing.T) {
	di, _ := buildBasicDisk(t) // AddBasicFile sets MGTFlags=0x20
	findings := checkBasicMGTFlags20(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("clean BASIC disk (MGTFlags=0x20): %d findings; want 0", len(findings))
	}
}

func TestBasicMGTFlags20Negative(t *testing.T) {
	di, dj := buildBasicDisk(t)
	dj[0].MGTFlags = 0x80
	di.WriteFileEntry(dj, 0)
	findings := checkBasicMGTFlags20(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "BASIC-MGTFLAGS-20" {
		t.Fatalf("got %d findings, first=%+v; want 1 BASIC-MGTFLAGS-20", len(findings), findings)
	}
}
