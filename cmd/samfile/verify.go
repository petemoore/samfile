package main

import (
	"fmt"
	"strings"

	"github.com/petemoore/samfile/v3"
)

// runVerify is the entry point for the `samfile verify` subcommand.
// Phase 1 implements the minimum useful behaviour: load the disk,
// call Verify, print findings grouped by severity, return nil
// unless a fatal finding is present (in which case return a non-nil
// error so main exits non-zero).
//
// CLI flags (--severity, --json, --dialect, --rule, --quiet, --all)
// are deferred to a later phase. Phase 1 always shows every finding
// regardless of severity, so the smoke test can see DISK-NOT-EMPTY.
//
// Signature note: existing subcommands (ls, cat, etc.) take a
// docopt.Opts map and call log.Fatal on error. runVerify deviates
// deliberately — Task 9 defines it as an isolated, testable unit
// (image-path in, error out); Task 10 wires it into the dispatcher
// by extracting arguments["-i"].(string) and calling log.Fatal on
// a non-nil return.
func runVerify(imagePath string) error {
	di, err := samfile.Load(imagePath)
	if err != nil {
		return fmt.Errorf("verify: %w", err)
	}

	report := di.Verify()

	fmt.Printf("samfile verify: results for %s\n", imagePath)
	fmt.Printf("detected dialect: %s\n", report.Dialect)
	fmt.Println()

	if len(report.Findings) == 0 {
		fmt.Println("no findings.")
		return nil
	}

	// Group by severity, highest first.
	severities := []samfile.Severity{
		samfile.SeverityFatal,
		samfile.SeverityStructural,
		samfile.SeverityInconsistency,
		samfile.SeverityCosmetic,
	}
	for _, s := range severities {
		findings := report.BySeverity(s)
		if len(findings) == 0 {
			continue
		}
		fmt.Printf("%s (%d):\n", strings.ToUpper(s.String()), len(findings))
		for _, f := range findings {
			fmt.Printf("  [%s]", f.RuleID)
			// TODO(phase 3+): when a Rule uses SectorLocation, render
			// Sector (track/sector) and ByteOffset here too. Phase 1 has
			// no SectorLocation users so this formatter only handles
			// disk-wide and slot-level locations.
			if !f.Location.IsDiskWide() {
				fmt.Printf(" slot %d", f.Location.Slot)
				if f.Location.Filename != "" {
					fmt.Printf(" (%s)", f.Location.Filename)
				}
			}
			fmt.Println()
			fmt.Printf("    %s\n", f.Message)
			fmt.Printf("    citation: %s\n", f.Citation)
		}
		fmt.Println()
	}

	fmt.Printf("%d finding(s).\n", len(report.Findings))
	if report.HasFatal() {
		return fmt.Errorf("verify: %d fatal finding(s)", len(report.BySeverity(samfile.SeverityFatal)))
	}
	return nil
}
