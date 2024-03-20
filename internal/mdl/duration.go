package mdl

import "time"

// Duration represents a range of time from a start to an end.
// Both the start and end are inclusive.
type Duration struct {
	Start time.Time
	End   time.Time
}

// NewDuration creates and returns a new TimeRange given a start and end time.
func NewDuration(start, end time.Time) *Duration {
	return &Duration{Start: start, End: end}
}

// Contains checks if a given time is within the time range, inclusive of start and end.
func (tr *Duration) Contains(t time.Time) bool {
	return !t.Before(tr.Start) && !t.After(tr.End)
}
