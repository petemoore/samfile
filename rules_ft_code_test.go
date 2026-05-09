package samfile

import "testing"

func TestCodeLoadAboveROMPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100) // load 0x8000
	findings := checkCodeLoadAboveROM(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("clean CODE file at 0x8000: %d findings; want 0", len(findings))
	}
}

func TestCodeLoadAboveROMNegative(t *testing.T) {
	// AddCodeFile rejects load < 0x4000 (samfile.go:799-801), so we
	// can't build a violating file via the public API. Patch the
	// dir entry directly to point below ROM.
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].StartAddressPage = 0
	dj[0].StartAddressPageOffset = 0
	// Subtract 1 to land at 0x3FFF (below ROM boundary).
	// Decoded Start() = ((StartPage & 0x1F)+1)<<14 | (PageOffset & 0x3FFF).
	// We need < 0x4000, so the (+1)<<14 path with StartPage=0 always
	// gives 0x4000. We need an off-by-one: set StartPage to a value
	// that, after +1 shift, produces 0x3FFF or below. The only way is
	// for the formula's & 0x3FFF mask of PageOffset to interact with
	// (page+1)<<14 — which it can't, since the mask isolates bits.
	//
	// Conclusion: samfile.Start()'s +1 shift makes load < 0x4000
	// unreachable via legal field values. This rule will never fire
	// on a samfile-parsed disk. Skip the negative test and document.
	t.Skip("samfile.Start()'s +1 shift makes Start()<0x4000 unreachable via FileEntry fields; rule is documentation-only and exists for parity with the catalog")
	_ = di
}

func TestCodeLoadFitsInMemoryPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkCodeLoadFitsInMemory(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("100-byte file at 0x8000: %d findings; want 0", len(findings))
	}
}

func TestCodeLoadFitsInMemoryNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	// Patch Pages to 31 (the off-disk pseudo-page marker via samfile's +1)
	// so length decodes huge AND the load address is near the top of RAM.
	dj[0].Pages = 31           // length = 31 * 16384 + (LengthMod16K & 0x3FFF)
	dj[0].LengthMod16K = 0x3FFF // max bits in low 14 bits
	di.WriteFileEntry(dj, 0)
	findings := checkCodeLoadFitsInMemory(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "CODE-LOAD-FITS-IN-MEMORY" {
		t.Fatalf("got %d findings, first=%+v; want 1 CODE-LOAD-FITS-IN-MEMORY", len(findings), findings)
	}
}

func TestCodeExecWithinLoadedRangePositive(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("TEST", make([]byte, 100), 0x8000, 0x8010); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	findings := checkCodeExecWithinLoadedRange(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("exec 0x8010 inside [0x8000, 0x8064): %d findings; want 0", len(findings))
	}
}

func TestCodeExecWithinLoadedRangeNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	// Set a real exec address (clear the 0xFF marker), but place it
	// far outside the loaded region [0x8000, 0x8064).
	dj[0].ExecutionAddressDiv16K = 0x05      // page 5 = 0x14000
	dj[0].ExecutionAddressMod16K = 0x8000    // offset 0 (PageOffset form)
	di.WriteFileEntry(dj, 0)
	findings := checkCodeExecWithinLoadedRange(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "CODE-EXEC-WITHIN-LOADED-RANGE" {
		t.Fatalf("got %d findings, first=%+v; want 1 CODE-EXEC-WITHIN-LOADED-RANGE", len(findings), findings)
	}
}

func TestCodeFileTypeInfoEmptyPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkCodeFileTypeInfoEmpty(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestCodeFileTypeInfoEmptyNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	dj[0].FileTypeInfo[5] = 0xAA // neither 0x00 (samfile) nor 0xFF (ROM SAVE)
	di.WriteFileEntry(dj, 0)
	findings := checkCodeFileTypeInfoEmpty(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 1 || findings[0].RuleID != "CODE-FILETYPEINFO-EMPTY" {
		t.Fatalf("got %d findings, first=%+v; want 1 CODE-FILETYPEINFO-EMPTY", len(findings), findings)
	}
}

func TestCodeFileTypeInfoEmptyAcceptsAllFF(t *testing.T) {
	// Real ROM SAMDOS-2 SAVE 0xFF-fills FileTypeInfo via HDCLP2
	// (rom-disasm:22076-22080); observed on the M0 boot disk's slot-4
	// OUT entry. This is a legitimate "unused" marker — the rule must
	// not fire when every byte is 0xFF.
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	for i := range dj[0].FileTypeInfo {
		dj[0].FileTypeInfo[i] = 0xFF
	}
	di.WriteFileEntry(dj, 0)
	findings := checkCodeFileTypeInfoEmpty(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("FileTypeInfo = 0xFF × 11 (ROM SAMDOS-2 convention): %d findings; want 0", len(findings))
	}
}

func TestCodeFileTypeInfoEmptyAcceptsAll0x20(t *testing.T) {
	// Iteration 1 SCOPE: 0x20 added as a third legitimate "unused"
	// marker (HDR space-fill leakage from ROM HDCLP at
	// rom-disasm:22070-22074). 99% of corpus FileTypeInfo-mismatch
	// fires are byte 0x20 — the rule must not fire on this value.
	di, dj := cleanSingleFileDisk(t, "TEST", 100)
	for i := range dj[0].FileTypeInfo {
		dj[0].FileTypeInfo[i] = 0x20
	}
	di.WriteFileEntry(dj, 0)
	findings := checkCodeFileTypeInfoEmpty(&CheckContext{Disk: di, Journal: di.DiskJournal()})
	if len(findings) != 0 {
		t.Errorf("FileTypeInfo = 0x20 × 11 (HDR space-fill leakage): %d findings; want 0", len(findings))
	}
}
