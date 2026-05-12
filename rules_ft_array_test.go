package samfile

import "testing"

func TestArrayFileTypeInfoTLBYTENamePositive(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Type = FT_NUM_ARRAY
	dj[0].FileTypeInfo[0] = 0x42 // TLBYTE
	copy(dj[0].FileTypeInfo[1:], []byte("ARR       "))
	di.WriteFileEntry(dj, 0)
	findings := checkArrayFileTypeInfoTLBYTEName(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("array file with populated FileTypeInfo: %d findings; want 0", len(findings))
	}
}

func TestArrayFileTypeInfoTLBYTENameNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].Type = FT_STR_ARRAY
	// FileTypeInfo is zero by default (AddCodeFile leaves it that way).
	di.WriteFileEntry(dj, 0)
	findings := checkArrayFileTypeInfoTLBYTEName(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "ARRAY-FILETYPEINFO-TLBYTE-NAME" {
		t.Fatalf("got %d findings, first=%+v; want 1 ARRAY-FILETYPEINFO-TLBYTE-NAME", len(findings), findings)
	}
}
