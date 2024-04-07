package srv

import (
	"fmt"
	"github.com/heather92115/verdure-admin/internal/db/mock"
	"github.com/heather92115/verdure-admin/internal/mdl"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestValidateFixit(t *testing.T) {
	tests := []struct {
		name      string
		fixit     mdl.Fixit
		wantError bool
		errorMsg  string
	}{
		{
			name: "Valid Fixit",
			fixit: mdl.Fixit{
				FieldName: "Valid Field Name",
				Comments:  "Valid comment within allowed length.",
			},
			wantError: false,
		},
		{
			name: "FieldName Exceeds Max Length",
			fixit: mdl.Fixit{
				FieldName: strings.Repeat("a", maxFixitFieldNameLen+1), // Exceeds max length
				Comments:  "Valid comment.",
			},
			wantError: true,
			errorMsg:  fmt.Sprintf("Field Name must be shorter than %d characters", maxFixitFieldNameLen),
		},
		{
			name: "Comments Exceed Max Length",
			fixit: mdl.Fixit{
				FieldName: "Valid Field Name",
				Comments:  strings.Repeat("b", maxFixitCommitLen+1), // Exceeds max length
			},
			wantError: true,
			errorMsg:  fmt.Sprintf("Commits must be shorter than %d characters", maxFixitCommitLen),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFixit(&tt.fixit)
			if (err != nil) != tt.wantError {
				t.Errorf("validateFixit() error = %v, wantError %v", err, tt.wantError)
			}
			if tt.wantError && err != nil && !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("validateFixit() error = %v, expected to contain errorMsg %v", err, tt.errorMsg)
			}
		})
	}
}

func TestFixitService_FindFixitByID(t *testing.T) {
	// Create an instance of FixitService with mocks
	fixitService := createMockFixitService()

	testFixit := &mdl.Fixit{
		ID:        1,
		VocabID:   100,
		Status:    "pending",
		FieldName: "Definition",
		Comments:  "Initial comment",
		CreatedBy: "tester",
		Created:   time.Now(),
	}
	_ = fixitService.CreateFixit(testFixit)

	fixit, err := fixitService.FindFixitByID(1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if fixit.ID != testFixit.ID {
		t.Errorf("Expected fixit ID %d, got %d", testFixit.ID, fixit.ID)
	}

	_, err = fixitService.FindFixitByID(999)
	if err == nil {
		t.Error("Expected an error for non-existing fixit, but got nil")
	}
}

// TestFixitService_FindFixits tests the functionality of the FindFixits method.
func TestFixitService_FindFixits(t *testing.T) {
	// Create an instance of FixitService with mocks
	fixitService := createMockFixitService()

	// Seed the mock repository with test fixits
	testFixit1 := &mdl.Fixit{
		ID:        1,
		VocabID:   101,
		Status:    "pending",
		FieldName: "Definition",
		Comments:  "Initial comment",
		CreatedBy: "tester",
		Created:   time.Now(),
	}
	testFixit2 := &mdl.Fixit{
		ID:        2,
		VocabID:   102,
		Status:    "completed",
		FieldName: "",
		Comments:  "Resolved issue",
		CreatedBy: "tester",
		Created:   time.Now(),
	}
	_ = fixitService.CreateFixit(testFixit1)
	_ = fixitService.CreateFixit(testFixit2)

	// Define test cases
	tests := []struct {
		name           string
		status         mdl.StatusType
		vocabID        int
		duration       *mdl.Duration
		limit          int
		expectedFixits []mdl.Fixit
	}{
		{
			name:           "Find pending Fixits",
			status:         "pending",
			vocabID:        0,   // Any vocab ID
			duration:       nil, // Any time duration
			limit:          10,
			expectedFixits: []mdl.Fixit{*testFixit1},
		},
		{
			name:           "Find completed Fixits",
			status:         "completed",
			vocabID:        0,   // Any vocab ID
			duration:       nil, // Any time duration
			limit:          10,
			expectedFixits: []mdl.Fixit{*testFixit2},
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixits, err := fixitService.FindFixits(tt.status, tt.vocabID, tt.duration, tt.limit)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(*fixits, tt.expectedFixits) {
				t.Errorf("Expected fixits %+v, got %+v", tt.expectedFixits, *fixits)
			}
		})
	}
}

func TestFixitService_CreateFixit(t *testing.T) {
	// Create an instance of FixitService with mocks
	fixitService := createMockFixitService()

	// Define test cases
	tests := []struct {
		name    string
		fixit   *mdl.Fixit
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful fixit creation",
			fixit: &mdl.Fixit{
				VocabID:   101,
				Status:    mdl.StatusType("pending"),
				FieldName: "New field name",
				Comments:  "Initial comment",
				CreatedBy: "tester",
			},
			wantErr: false,
		},
		{
			name: "fixit field name too long",
			fixit: &mdl.Fixit{
				VocabID:   102,
				Status:    mdl.StatusType("completed"),
				FieldName: string(make([]rune, 41)), // 41 characters,
				Comments:  "Should fail due to field name",
				CreatedBy: "tester",
			},
			wantErr: true,
			errMsg:  "Field Name must be shorter than 40 characters",
		},
		{
			name: "Valid fixit - missing field name",
			fixit: &mdl.Fixit{
				VocabID:   103,
				Status:    mdl.StatusType("in_progress"),
				Comments:  "Missing field name",
				CreatedBy: "tester",
			},
			wantErr: false,
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fixitService.CreateFixit(tt.fixit)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: CreateFixit() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			} else if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("%s: CreateFixit() error = %v, wantErrMsg to contain %v", tt.name, err, tt.errMsg)
			}
		})
	}
}

func TestFixitService_UpdateFixit(t *testing.T) {

	// Create an instance of FixitService with mocks
	fixitService := createMockFixitService()

	// Seed the mock repository with an existing fixit for update
	existingFixit := &mdl.Fixit{
		ID:        1,
		VocabID:   101,
		Status:    mdl.StatusType("pending"),
		FieldName: "Existing field name",
		Comments:  "Existing comment",
		CreatedBy: "tester",
	}
	_ = fixitService.CreateFixit(existingFixit)

	// Define test cases
	tests := []struct {
		name    string
		fixit   *mdl.Fixit
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful fixit update",
			fixit: &mdl.Fixit{
				ID:        1,
				VocabID:   101,
				Status:    mdl.StatusType("completed"),
				FieldName: "Updated field name",
				Comments:  "Updated comment",
				CreatedBy: "tester",
			},
			wantErr: false,
		},
		{
			name: "FieldName exceeds max length",
			fixit: &mdl.Fixit{
				ID:        1,
				VocabID:   101,
				Status:    mdl.StatusType("completed"),
				FieldName: string(make([]rune, maxFixitFieldNameLen+1)), // Exceeds max length
				Comments:  "Should fail due to field name length",
				CreatedBy: "tester",
			},
			wantErr: true,
			errMsg:  fmt.Sprintf(errFmtStrLen, "Field Name", maxFixitFieldNameLen),
		},
		{
			name: "Fixit does not exist",
			fixit: &mdl.Fixit{
				ID:        999, // Non-existing ID
				VocabID:   102,
				Status:    mdl.StatusType("in_progress"),
				FieldName: "Non-existing field name",
				Comments:  "Non-existing comment",
				CreatedBy: "tester",
			},
			wantErr: true,
			errMsg:  "fixit not found",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedFixit, err := fixitService.UpdateFixit(tt.fixit)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: UpdateFixit() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			} else if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("%s: UpdateFixit() error = %v, wantErrMsg to contain %v", tt.name, err, tt.errMsg)
			}
			if !tt.wantErr && updatedFixit != nil && updatedFixit.FieldName != tt.fixit.FieldName {
				t.Errorf("%s: FieldName not updated correctly, got = %v, want = %v", tt.name, updatedFixit.FieldName, tt.fixit.FieldName)
			}
		})
	}
}

func createMockFixitService() FixitService {
	// Initialize the mock repositories
	mockFixitRepo := mock.NewMockFixitRepository()
	mockAuditRepo := mock.NewMockAuditRepository()
	mockAuditService := &AuditService{repo: mockAuditRepo}

	// Create an instance of FixitService with mocks
	fixitService := FixitService{
		repo:         mockFixitRepo,
		auditService: *mockAuditService,
	}

	return fixitService
}
