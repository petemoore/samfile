package samfile

import "fmt"

// §14 Cosmetic / canonical-output rules (catalog docs/disk-validity-rules.md §14).
// Rules in this file warn when dir-entry bytes diverge from
// the conventions real ROM SAVE produces, without affecting
// runtime behaviour. They apply to all dialects.

// Ensure fmt is used (will be used by rules added in Task 3).
var _ = fmt.Sprintf
