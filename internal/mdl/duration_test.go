package mdl

import (
	"testing"
	"time"
)

func TestDuration_Contains(t *testing.T) {
	start := time.Now()
	end := start.Add(24 * time.Hour) // 24 hours later
	duration := NewDuration(start, end)

	// Test time within the range
	testTime := start.Add(12 * time.Hour)
	if !duration.Contains(testTime) {
		t.Errorf("Duration should contain the time %v", testTime)
	}

	// Test time before the range
	testTimeBefore := start.Add(-1 * time.Hour)
	if duration.Contains(testTimeBefore) {
		t.Errorf("Duration should not contain the time %v", testTimeBefore)
	}

	// Test time after the range
	testTimeAfter := end.Add(1 * time.Hour)
	if duration.Contains(testTimeAfter) {
		t.Errorf("Duration should not contain the time %v", testTimeAfter)
	}
}
