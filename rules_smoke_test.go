package samfile

import (
	"testing"
)

// TestPhase3RegistryGrowth pins the expected rule count once Phase 3
// is fully implemented. It will fail in Task 1 (only 1 rule registered)
// and pass once Tasks 2-6 land the remaining 19 rules. This is a
// regression gate: if any rule is accidentally removed or never
// registered, this test fails.
func TestPhase3RegistryGrowth(t *testing.T) {
	if got := len(Rules()); got != 20 {
		t.Errorf("len(Rules()) = %d; want 20 (1 smoke + 19 phase-3 rules)", got)
	}
}

func TestDiskNotEmptyRulePositive(t *testing.T) {
	di := NewDiskImage()
	if err := di.AddCodeFile("A", []byte("hello"), 0x8000, 0); err != nil {
		t.Fatalf("AddCodeFile: %v", err)
	}
	findings := checkDiskNotEmpty(&CheckContext{
		Disk:    di,
		Journal: di.DiskJournal(),
		Dialect: DialectUnknown,
	})
	if len(findings) != 0 {
		t.Errorf("checkDiskNotEmpty on populated disk returned %d findings; want 0", len(findings))
	}
}

func TestDiskNotEmptyRuleNegative(t *testing.T) {
	di := NewDiskImage()
	findings := checkDiskNotEmpty(&CheckContext{
		Disk:    di,
		Journal: di.DiskJournal(),
		Dialect: DialectUnknown,
	})
	if len(findings) != 1 {
		t.Fatalf("checkDiskNotEmpty on empty disk returned %d findings; want 1", len(findings))
	}
	f := findings[0]
	if f.RuleID != "DISK-NOT-EMPTY" {
		t.Errorf("RuleID = %q; want DISK-NOT-EMPTY", f.RuleID)
	}
	if f.Severity != SeverityInconsistency {
		t.Errorf("Severity = %v; want inconsistency", f.Severity)
	}
	if !f.Location.IsDiskWide() {
		t.Errorf("Location.IsDiskWide() = false; want true")
	}
	if f.Message == "" {
		t.Error("Message empty")
	}
}

func TestDiskNotEmptyRegistered(t *testing.T) {
	found := false
	for _, r := range Rules() {
		if r.ID == "DISK-NOT-EMPTY" {
			found = true
			break
		}
	}
	if !found {
		t.Error("DISK-NOT-EMPTY rule not in registry")
	}
}
