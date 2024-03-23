package convert

import (
	"testing"
	"time"
)

func TestTimeToGQLDateTime(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "Convert UTC time to GQL DateTime",
			time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "2023-01-01T00:00:00Z",
		},
		{
			name: "Convert non-UTC time to GQL DateTime",
			time: time.Date(2023, 1, 1, 15, 30, 45, 0, time.FixedZone("EST", -5*3600)),
			want: "2023-01-01T20:30:45Z", // EST is UTC-5
		},
		{
			name: "Leap year and second",
			time: time.Date(2020, 2, 29, 23, 59, 59, 0, time.UTC),
			want: "2020-02-29T23:59:59Z",
		},
		{
			name: "With milliseconds",
			time: time.Date(2023, 8, 15, 12, 0, 0, 123456789, time.UTC),
			want: "2023-08-15T12:00:00Z", // Milliseconds are not included in ISO 8601 string format
		},
		{
			name: "Different timezone (Asia/Tokyo)",
			time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.FixedZone("JST", 9*3600)),
			want: "2022-12-31T15:00:00Z", // JST is UTC+9, so it goes back to the previous day
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeToGQLDateTime(tt.time); got != tt.want {
				t.Errorf("timeToGQLDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGqlDateTimeToTime(t *testing.T) {
	// Define test cases
	tests := []struct {
		name    string
		gqlDate string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Valid GQL DateTime",
			gqlDate: "2023-01-01T00:00:00Z",
			want:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Valid GQL DateTime with offset",
			gqlDate: "2023-01-01T05:00:00+05:00",
			want:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Valid GQL DateTime with milliseconds",
			gqlDate: "2023-01-01T00:00:00.123Z",
			want:    time.Date(2023, 1, 1, 0, 0, 0, 123000000, time.UTC),
			wantErr: false,
		},
		{
			name:    "Invalid GQL DateTime format",
			gqlDate: "invalid",
			wantErr: true,
		},
		{
			name:    "Invalid GQL DateTime with wrong timezone format",
			gqlDate: "2023-01-01T00:00:00Z+01:00",
			wantErr: true,
		},
		{
			name:    "Empty GQL DateTime string",
			gqlDate: "",
			wantErr: true,
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gqlDateTimeToTime(tt.gqlDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("gqlDateTimeToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("gqlDateTimeToTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGqlDateTimeToDuration(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		startTime string
		endTime   string
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		{
			name:      "Valid datetime range",
			startTime: "2023-01-01T00:00:00Z",
			endTime:   "2023-01-02T00:00:00Z",
			wantStart: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr:   false,
		},
		{
			name:      "Invalid start time format",
			startTime: "invalid",
			endTime:   "2023-01-02T00:00:00Z",
			wantErr:   true,
		},
		{
			name:      "Invalid end time format",
			startTime: "2023-01-01T00:00:00Z",
			endTime:   "invalid",
			wantErr:   true,
		},
		{
			name:      "Start time after end time",
			startTime: "2023-01-03T00:00:00Z",
			endTime:   "2023-01-02T00:00:00Z",
			wantErr:   true,
		},
		{
			name:      "Empty start and end times",
			startTime: "",
			endTime:   "",
			wantStart: time.Now().Add(-1 * time.Hour), // This is expected to be approximately 1 hour ago from the current time.
			wantEnd:   time.Now(),
			wantErr:   false,
		},
		{
			name:      "Empty start time, valid end time",
			startTime: "",
			endTime:   "2023-01-02T00:00:00Z",         // in the past
			wantStart: time.Now().Add(-1 * time.Hour), // This is expected to be approximately 1 hour ago from the current time.
			wantEnd:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr:   true,
		},
		{
			name:      "Valid start time, empty end time",
			startTime: "2023-01-01T00:00:00Z",
			endTime:   "",
			wantStart: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Now(),
			wantErr:   false,
		},
	}

	// Define a reasonable variance allowance, e.g., 2 seconds
	variance := 2 * time.Second

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GqlDateTimeToDuration(tt.startTime, tt.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GqlDateTimeToDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				startDiff := got.Start.Sub(tt.wantStart)
				endDiff := got.End.Sub(tt.wantEnd)
				if startDiff < -variance || startDiff > variance || endDiff < -variance || endDiff > variance {
					t.Errorf("GqlDateTimeToDuration() got = %v to %v, want %v to %v with allowance %v", got.Start, got.End, tt.wantStart, tt.wantEnd, variance)
				}
			}
		})
	}
}
