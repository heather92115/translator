package convert

import (
	"errors"
	"fmt"
	"github.com/heather92115/translator/internal/mdl"
	"time"
)

// GqlDateTimeToDuration converts GraphQL DateTime strings into a mdl.Duration struct,
// representing a time range. It accepts start and end time as strings in ISO 8601 format.
// If the start time is not provided, it defaults to one hour ago from the current time.
// If the end time is not provided, it defaults to the current time. The function ensures
// that the start time is chronologically before the end time. If the start time is after
// the end time, an error is returned. This function facilitates converting user-provided
// date ranges from GraphQL inputs into internal representations of time ranges.
//
// Parameters:
// - startTime: The start time as a GraphQL DateTime string. If empty, defaults to one hour ago.
// - endTime: The end time as a GraphQL DateTime string. If empty, defaults to the current time.
//
// Returns:
// - A pointer to a mdl.Duration struct containing the parsed or defaulted start and end times.
// - An error if the start or end time strings are invalid, or if the start time is after the end time.
//
// Example:
// - Given valid start and end times, it returns a Duration with those times.
// - Given an empty start time, it defaults to one hour ago from now.
// - Given an empty end time, it defaults to the current time.
// - If the start time is provided as after the end time, it returns an error.
func GqlDateTimeToDuration(startTime string, endTime string) (*mdl.Duration, error) {
	// Default to start time as one hour ago and end time as current time
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now()

	var err error
	// Only attempt to parse startTime if it's provided
	if startTime != "" {
		start, err = gqlDateTimeToTime(startTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start time: %w", err)
		}
	}

	// Only attempt to parse endTime if it's provided
	if endTime != "" {
		end, err = gqlDateTimeToTime(endTime)
		if err != nil {
			return nil, fmt.Errorf("invalid end time: %w", err)
		}
	}

	// Ensure the start time is before the end time
	if start.After(end) {
		return nil, errors.New("start time must be before end time")
	}

	return &mdl.Duration{
		Start: start,
		End:   end,
	}, nil
}

// Convert from Go time.Time to GraphQL DateTime (ISO 8601 string)
func timeToGQLDateTime(t time.Time) string {
	// Ensure the time is in UTC before formatting
	utcTime := t.UTC()
	return utcTime.Format(time.RFC3339)
}

// Convert from GraphQL DateTime (ISO 8601 string) to Go time.Time
func gqlDateTimeToTime(gqlDateTime string) (time.Time, error) {
	return time.Parse(time.RFC3339, gqlDateTime)
}
