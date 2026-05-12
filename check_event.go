package samfile

// CheckEvent is the per-subject record emitted by Verify when a rule
// declares Scope + CheckSubject. One event per (rule, applicable
// subject) pair. Legacy rules (Check-only) emit one event per finding
// they produce, with Outcome=fail and a synthetic Subject ref.
type CheckEvent struct {
	Version int            `json:"v"`
	Disk    string         `json:"disk"`
	RuleID  string         `json:"rule_id"`
	Scope   string         `json:"scope"`
	Ref     string         `json:"ref"`
	Outcome string         `json:"outcome"` // pass | fail | not_applicable
	Attrs   map[string]any `json:"attrs,omitempty"`
	Finding *Finding       `json:"finding,omitempty"`
}

// EventRecorder collects CheckEvents during a Verify run. The verify
// CLI installs a recorder when --format jsonl is set; for normal text
// output the recorder is nil and no events are kept.
type EventRecorder struct {
	Events []CheckEvent
	Disk   string
}

func (r *EventRecorder) Record(e CheckEvent) {
	if r == nil {
		return
	}
	e.Version = 1
	if e.Disk == "" {
		e.Disk = r.Disk
	}
	r.Events = append(r.Events, e)
}
