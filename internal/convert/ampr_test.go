package convert

import (
	"github.com/heather92115/translator/graph/model"
	"github.com/heather92115/translator/internal/mdl"
	"reflect"
	"testing"
	"time"
)

func TestAuditsToGql(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		from    *[]mdl.Audit
		want    []*model.Audit
		wantErr bool
	}{
		{
			name: "Convert non-empty slice of Audits",
			from: &[]mdl.Audit{
				{
					ID:        1,
					ObjectID:  101,
					TableName: "TestTable1",
					Diff:      "TestDiff1",
					Before:    "TestBefore1",
					After:     "TestAfter1",
					Comments:  "TestComment1",
					CreatedBy: "TestUser1",
					Created:   now,
				},
				{
					ID:        2,
					ObjectID:  102,
					TableName: "TestTable2",
					Diff:      "TestDiff2",
					Before:    "TestBefore2",
					After:     "TestAfter2",
					Comments:  "TestComment2",
					CreatedBy: "TestUser2",
					Created:   now.Add(24 * time.Hour),
				},
			},
			want: []*model.Audit{
				{
					ID:        "1",
					ObjectID:  "101",
					TableName: "TestTable1",
					Diff:      "TestDiff1",
					Before:    "TestBefore1",
					After:     "TestAfter1",
					Comments:  "TestComment1",
					CreatedBy: "TestUser1",
					Created:   timeToGQLDateTime(now),
				},
				{
					ID:        "2",
					ObjectID:  "102",
					TableName: "TestTable2",
					Diff:      "TestDiff2",
					Before:    "TestBefore2",
					After:     "TestAfter2",
					Comments:  "TestComment2",
					CreatedBy: "TestUser2",
					Created:   timeToGQLDateTime(now.Add(24 * time.Hour)),
				},
			},
			wantErr: false,
		},
		{
			name:    "Convert nil slice of Audits",
			from:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Convert empty slice of Audits",
			from:    &[]mdl.Audit{},
			want:    []*model.Audit{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AuditsToGql(tt.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuditsToGql() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuditsToGql() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuditQueryMapper(t *testing.T) {
	// Prepare valid start and end time for duration testing
	validStartTime := time.Now().Truncate(time.Second).Format(time.RFC3339)
	validEndTime := time.Now().Add(24 * time.Hour).Truncate(time.Second).Format(time.RFC3339)

	// Expected duration for valid start and end times
	expectedStart := time.Now().Truncate(time.Second) // This truncates to second precision
	expectedEnd := expectedStart.Add(24 * time.Hour)
	expectedDuration := &mdl.Duration{Start: expectedStart, End: expectedEnd}

	tests := []struct {
		name         string
		objectID     string
		startTime    string
		endTime      string
		wantObjectID int
		wantDuration *mdl.Duration
		wantErr      bool
	}{
		{
			name:         "valid inputs",
			objectID:     "123",
			startTime:    validStartTime,
			endTime:      validEndTime,
			wantObjectID: 123,
			wantDuration: expectedDuration,
			wantErr:      false,
		},
		{
			name:      "invalid objectID",
			objectID:  "not_an_int",
			startTime: validStartTime,
			endTime:   validEndTime,
			wantErr:   true,
		},
		{
			name:      "invalid duration",
			objectID:  "123",
			startTime: "invalid_start_time",
			endTime:   "invalid_end_time",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotObjectID, gotDuration, err := AuditQueryMapper(tt.objectID, tt.startTime, tt.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuditQueryMapper() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			if gotObjectID != tt.wantObjectID {
				t.Errorf("AuditQueryMapper() gotObjectID = %v, want %v", gotObjectID, tt.wantObjectID)
				return
			}
			if !gotDuration.Start.Truncate(time.Second).Equal(tt.wantDuration.Start) || !gotDuration.End.Truncate(time.Second).Equal(tt.wantDuration.End) {
				t.Errorf("AuditQueryMapper() gotDuration = &{%v %v}, want &{%v %v}",
					gotDuration.Start, gotDuration.End, tt.wantDuration.Start, tt.wantDuration.End)
			}
		})
	}
}
