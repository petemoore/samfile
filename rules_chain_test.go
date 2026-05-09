package samfile

import "testing"

// walkChain unit tests come first — exercise the helper directly with
// fabricated chains so the rule tests can trust the walker behaves.

func TestWalkChainClean(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 600) // 2 sectors worth of payload
	result := walkChain(di, di.DiskJournal()[0].FirstSector)
	if !result.Terminated {
		t.Errorf("clean chain: Terminated = false; want true")
	}
	if result.Cycle != nil {
		t.Errorf("clean chain: Cycle = %v; want nil", result.Cycle)
	}
	if result.Bailed {
		t.Errorf("clean chain: Bailed = true; want false")
	}
	if len(result.Steps) < 2 {
		t.Errorf("clean chain: %d steps; want >= 2", len(result.Steps))
	}
}

func TestWalkChainCycleDetection(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 600)
	fe := dj[0]
	first := fe.FirstSector
	// Force sector 1's NextSector to point back at itself.
	sd, err := di.SectorData(first)
	if err != nil {
		t.Fatalf("SectorData: %v", err)
	}
	raw := sd[:]
	raw[510] = first.Track
	raw[511] = first.Sector
	di.WriteSector(first, sd)

	result := walkChain(di, first)
	if result.Cycle == nil {
		t.Errorf("Cycle = nil; want non-nil (chain points at itself)")
	}
	if result.Terminated {
		t.Errorf("Terminated = true; want false")
	}
}

// Now the three rule tests.

func TestChainTerminatorZeroZeroPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkChainTerminatorZeroZero(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean chain: %d findings; want 0", len(findings))
	}
}

func TestChainTerminatorZeroZeroNegative(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	first := di.DiskJournal()[0].FirstSector
	// Overwrite the terminator with a fake link (it would loop forever
	// if walkChain weren't bounded).
	sd, _ := di.SectorData(first)
	raw := sd[:]
	raw[510] = first.Track
	raw[511] = first.Sector
	di.WriteSector(first, sd)

	findings := checkChainTerminatorZeroZero(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "CHAIN-TERMINATOR-ZERO-ZERO" {
		t.Fatalf("got %d findings, first=%+v; want 1 CHAIN-TERMINATOR-ZERO-ZERO", len(findings), findings)
	}
}

func TestChainNoCyclePositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	findings := checkChainNoCycle(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean chain: %d findings; want 0", len(findings))
	}
}

func TestChainNoCycleNegative(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 100)
	first := di.DiskJournal()[0].FirstSector
	sd, _ := di.SectorData(first)
	raw := sd[:]
	raw[510] = first.Track
	raw[511] = first.Sector
	di.WriteSector(first, sd)

	findings := checkChainNoCycle(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "CHAIN-NO-CYCLE" {
		t.Fatalf("got %d findings, first=%+v; want 1 CHAIN-NO-CYCLE", len(findings), findings)
	}
}

func TestChainMatchesSAMPositive(t *testing.T) {
	di, _ := cleanSingleFileDisk(t, "TEST", 1500) // ~3 sectors worth
	findings := checkChainMatchesSAM(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 0 {
		t.Errorf("clean disk: %d findings; want 0", len(findings))
	}
}

func TestChainMatchesSAMNegative(t *testing.T) {
	di, dj := cleanSingleFileDisk(t, "TEST", 1500)
	// Clear a bit that IS set in the map so walked > mapSet.
	for i, b := range dj[0].SectorAddressMap {
		if b != 0 {
			dj[0].SectorAddressMap[i] &^= (b & -b) // clear the lowest set bit
			break
		}
	}
	di.WriteFileEntry(dj, 0)
	findings := checkChainMatchesSAM(&CheckContext{
		Disk: di, Journal: di.DiskJournal(),
	})
	if len(findings) != 1 || findings[0].RuleID != "CHAIN-MATCHES-SAM" {
		t.Fatalf("got %d findings, first=%+v; want 1 CHAIN-MATCHES-SAM", len(findings), findings)
	}
}
