package convert

import (
	"github.com/heather92115/verdure-admin/graph/model"
	"github.com/heather92115/verdure-admin/internal/mdl"
	"reflect"
	"testing"
	"time"
)

func TestFixitToGql(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		from    *mdl.Fixit
		want    *model.Fixit
		wantErr bool
	}{
		{
			name: "Convert valid Fixit",
			from: &mdl.Fixit{
				ID:        1,
				VocabID:   101,
				Status:    mdl.Pending,
				FieldName: "TestField",
				Comments:  "TestComment",
				CreatedBy: "TestUser",
				Created:   now,
			},
			want: &model.Fixit{
				ID:        "1",
				VocabID:   "101",
				Status:    "PENDING",
				FieldName: "TestField",
				Comments:  "TestComment",
				CreatedBy: "TestUser",
				Created:   timeToGQLDateTime(now),
			},
			wantErr: false,
		},
		{
			name:    "Nil input",
			from:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid status",
			from: &mdl.Fixit{
				ID:        2,
				VocabID:   102,
				Status:    "unknown",
				FieldName: "TestField",
				Comments:  "TestComment",
				CreatedBy: "TestUser",
				Created:   now,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FixitToGql(tt.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("FixitToGql() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FixitToGql() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixitsToGql(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		from    *[]mdl.Fixit
		want    []*model.Fixit
		wantErr bool
	}{
		{
			name: "Convert non-empty slice of Fixits",
			from: &[]mdl.Fixit{
				{
					ID:        1,
					VocabID:   101,
					Status:    mdl.Pending,
					FieldName: "TestField1",
					Comments:  "TestComment1",
					CreatedBy: "TestUser1",
					Created:   now,
				},
				{
					ID:        2,
					VocabID:   102,
					Status:    mdl.Completed,
					FieldName: "TestField2",
					Comments:  "TestComment2",
					CreatedBy: "TestUser2",
					Created:   now.Add(24 * time.Hour),
				},
			},
			want: []*model.Fixit{
				{
					ID:        "1",
					VocabID:   "101",
					Status:    "PENDING",
					FieldName: "TestField1",
					Comments:  "TestComment1",
					CreatedBy: "TestUser1",
					Created:   timeToGQLDateTime(now),
				},
				{
					ID:        "2",
					VocabID:   "102",
					Status:    "COMPLETED",
					FieldName: "TestField2",
					Comments:  "TestComment2",
					CreatedBy: "TestUser2",
					Created:   timeToGQLDateTime(now.Add(24 * time.Hour)),
				},
			},
			wantErr: false,
		},
		{
			name:    "Convert nil slice of Fixits",
			from:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Convert empty slice of Fixits",
			from:    &[]mdl.Fixit{},
			want:    []*model.Fixit{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FixitsToGql(tt.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("FixitsToGql() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FixitsToGql() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFixitFromGql(t *testing.T) {
	tests := []struct {
		name    string
		from    *model.NewFixit
		want    *mdl.Fixit
		wantErr bool
	}{
		{
			name: "Valid NewFixit conversion",
			from: &model.NewFixit{
				VocabID:   "101",
				Status:    "PENDING",
				FieldName: "Test Field",
				Comments:  "Test Comment",
			},
			want: &mdl.Fixit{
				VocabID:   101,
				Status:    mdl.Pending,
				FieldName: "Test Field",
				Comments:  "Test Comment",
			},
			wantErr: false,
		},
		{
			name:    "Nil input",
			from:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid VocabID format",
			from: &model.NewFixit{
				VocabID:   "invalid",
				Status:    "PENDING",
				FieldName: "Test Field",
				Comments:  "Test Comment",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid status value",
			from: &model.NewFixit{
				VocabID:   "101",
				Status:    "UNKNOWN",
				FieldName: "Test Field",
				Comments:  "Test Comment",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFixitFromGql(tt.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFixitFromGql() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFixitFromGql() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateFixitFromGql(t *testing.T) {
	tests := []struct {
		name    string
		from    *model.UpdateFixit
		want    *mdl.Fixit
		wantErr bool
	}{
		{
			name: "Valid UpdateFixit conversion",
			from: &model.UpdateFixit{
				ID:        "1",
				Status:    "COMPLETED",
				FieldName: "Updated Field",
				Comments:  "Updated Comment",
			},
			want: &mdl.Fixit{
				ID:        1,
				Status:    mdl.Completed,
				FieldName: "Updated Field",
				Comments:  "Updated Comment",
			},
			wantErr: false,
		},
		{
			name:    "Nil input",
			from:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid ID format",
			from: &model.UpdateFixit{
				ID:        "invalid",
				Status:    "COMPLETED",
				FieldName: "Updated Field",
				Comments:  "Updated Comment",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid status value",
			from: &model.UpdateFixit{
				ID:        "1",
				Status:    "UNKNOWN",
				FieldName: "Updated Field",
				Comments:  "Updated Comment",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateFixitFromGql(tt.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateFixitFromGql() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateFixitFromGql() = %v, want %v", got, tt.want)
			}
		})
	}
}

// FixitsQueryMapper converts GraphQL query parameters into their corresponding internal representations.
// It takes a fixit status as a string, a vocabID as a string, and start and end times as ISO 8601 formatted strings.
// It returns the internal status type, vocabID as an integer, a duration struct representing the time range, and an error if any conversions fail.
// This function is primarily used to prepare parameters for querying the database with filters provided by a GraphQL request.
//
// Parameters:
//   - status: Fixit status in GraphQL enum format (e.g., "PENDING").
//   - vocabID: Vocab ID as a string, which should be convertible to an integer.
//   - startTime: Start of the duration as an ISO 8601 formatted string.
//   - endTime: End of the duration as an ISO 8601 formatted string.
//
// Returns:
//   - StatusType: The internal representation of the fixit status.
//   - int: The vocabID converted to an integer.
//   - *Duration: A pointer to a Duration struct representing the time range.
//   - error: An error if any of the conversions fail (invalid status, non-integer vocabID, or invalid date formats).
func TestFixitsQueryMapper(t *testing.T) {
	// Prepare a valid start and end time for duration testing
	validStartTime := time.Now().Truncate(time.Second).Format(time.RFC3339)
	validEndTime := time.Now().Add(24 * time.Hour).Truncate(time.Second).Format(time.RFC3339)

	// Expected duration for the above valid values
	expectedStart := time.Now().Truncate(time.Second) // This truncates to second precision
	expectedEnd := expectedStart.Add(24 * time.Hour)
	expectedDuration := &mdl.Duration{Start: expectedStart, End: expectedEnd}

	tests := []struct {
		name         string
		status       model.Status
		vocabID      string
		startTime    string
		endTime      string
		wantStatus   mdl.StatusType
		wantVocabID  int
		wantDuration *mdl.Duration
		wantErr      bool
	}{
		{
			name:         "valid inputs",
			status:       "PENDING",
			vocabID:      "123",
			startTime:    validStartTime,
			endTime:      validEndTime,
			wantStatus:   mdl.Pending,
			wantVocabID:  123,
			wantDuration: expectedDuration,
			wantErr:      false,
		},
		{
			name:      "invalid status",
			status:    "INVALID_STATUS",
			vocabID:   "123",
			startTime: validStartTime,
			endTime:   validEndTime,
			wantErr:   true,
		},
		{
			name:      "invalid vocabID",
			status:    "PENDING",
			vocabID:   "not_an_int",
			startTime: validStartTime,
			endTime:   validEndTime,
			wantErr:   true,
		},
		{
			name:      "invalid duration",
			status:    "PENDING",
			vocabID:   "123",
			startTime: "invalid_start_time",
			endTime:   "invalid_end_time",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotVocabID, gotDuration, err := FixitsQueryMapper(tt.status, tt.vocabID, tt.startTime, tt.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("FixitsQueryMapper() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			if gotStatus != tt.wantStatus {
				t.Errorf("FixitsQueryMapper() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
				return
			}
			if gotVocabID != tt.wantVocabID {
				t.Errorf("FixitsQueryMapper() gotVocabID = %v, want %v", gotVocabID, tt.wantVocabID)
				return
			}
			if !gotDuration.Start.Truncate(time.Second).Equal(tt.wantDuration.Start) || !gotDuration.End.Truncate(time.Second).Equal(tt.wantDuration.End) {
				t.Errorf("FixitsQueryMapper() gotDuration = &{%v %v}, want &{%v %v}",
					gotDuration.Start, gotDuration.End, tt.wantDuration.Start, tt.wantDuration.End)
			}
		})
	}
}
