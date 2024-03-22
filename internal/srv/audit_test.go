package srv

import (
	"fmt"
	"github.com/heather92115/translator/internal/db/mock"
	"github.com/heather92115/translator/internal/mdl"
	"testing"
	"time"
)

func TestJsonDiff(t *testing.T) {

	jsonStr1 := `{"name":"Alice", "age":30, "car":null}`
	jsonStr2 := `{"name":"Alice", "age":31, "car":"Tesla"}`

	diffs := CompareJSON(jsonStr1, jsonStr2)

	fmt.Println("Differences:", diffs)
}

func TestVocabJsonDiff(t *testing.T) {
	tests := []struct {
		name     string
		before   mdl.Vocab
		after    mdl.Vocab
		expected string
	}{
		{
			name: "ser test",
			before: mdl.Vocab{
				LearningLang:     "ser",
				FirstLang:        "",
				Alternatives:     "Alternative",
				Skill:            "Beginner",
				Infinitive:       "ser",
				Pos:              "verb",
				Hint:             "A hint",
				KnownLangCode:    "en",
				LearningLangCode: "es",
			},
			after: mdl.Vocab{
				LearningLang:     "ser",
				FirstLang:        "to be",
				Alternatives:     "",
				Skill:            "Beginner",
				Infinitive:       "",
				Pos:              "verb",
				Hint:             "",
				KnownLangCode:    "en",
				LearningLangCode: "es",
			},
			expected: `[{"key":"'alternatives'","before":"Alternative","after":""},{"key":"'first_lang'","before":"","after":"to be"},{"key":"'hint'","before":"A hint","after":""},{"key":"'infinitive'","before":"ser","after":""}]`,
		},
		{
			name: "perro test",
			before: mdl.Vocab{
				LearningLang:     "perro",
				FirstLang:        "",
				Alternatives:     "",
				Skill:            "Pets",
				Infinitive:       "",
				Pos:              "noun",
				Hint:             "not gato",
				KnownLangCode:    "en",
				LearningLangCode: "es",
			},
			after: mdl.Vocab{
				LearningLang:     "perro",
				FirstLang:        "dog",
				Alternatives:     "perra, perros, perras",
				Skill:            "Pets",
				Infinitive:       "",
				Pos:              "noun",
				Hint:             "starts with pe",
				KnownLangCode:    "en",
				LearningLangCode: "es",
			},
			expected: `[{"key":"'alternatives'","before":"","after":"perra, perros, perras"},{"key":"'first_lang'","before":"","after":"dog"},{"key":"'hint'","before":"not gato","after":"starts with pe"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fmt.Printf("before %s\n", tt.before.JSON())
			fmt.Printf("after %s\n", tt.after.JSON())

			diffs := CompareJSON(tt.before.JSON(), tt.after.JSON())

			if diffs != tt.expected {
				t.Errorf("CompareJSON() mismatch, \nexpected %s, \nactual   %s\n", tt.expected, diffs)
				fmt.Println("Differences:", diffs)
			}
		})
	}
}

// TestAuditService_FindAudits tests the FindAudits method of AuditService.
func TestAuditService_FindAudits(t *testing.T) {
	// Setup
	mockRepo := mock.NewMockAuditRepository()
	service := &AuditService{repo: mockRepo}

	// Seed some audit data into the mock repository
	_ = mockRepo.CreateAudit(&mdl.Audit{
		ID:        1,
		ObjectID:  123,
		TableName: "users",
		Created:   time.Now(),
	})
	_ = mockRepo.CreateAudit(&mdl.Audit{
		ID:        2,
		ObjectID:  456,
		TableName: "products",
		Created:   time.Now(),
	})

	// Define test cases
	tests := []struct {
		name        string
		tableName   string
		expectCount int
	}{
		{
			name:        "Find audits for users table",
			tableName:   "users",
			expectCount: 1,
		},
		{
			name:        "Find audits for non-existing table",
			tableName:   "non_existing",
			expectCount: 0,
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := &mdl.Duration{Start: time.Now().Add(-24 * time.Hour), End: time.Now()}
			audits, err := service.FindAudits(tt.tableName, duration, 10)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(*audits) != tt.expectCount {
				t.Errorf("Expected %d audits, got %d", tt.expectCount, len(*audits))
			}
		})
	}
}

// TestAuditService_FindAuditByID tests the functionality of FindAuditByID method.
func TestAuditService_FindAuditByID(t *testing.T) {
	// Initialize the mock repository and service
	mockRepo := mock.NewMockAuditRepository()
	service := &AuditService{repo: mockRepo}

	// Seed the mock repository with a test audit
	testAudit := &mdl.Audit{
		ID:        1,
		ObjectID:  101,
		TableName: "test_table",
		Created:   time.Now(),
	}
	_ = mockRepo.CreateAudit(testAudit)

	// Test finding an existing audit
	t.Run("Find existing audit", func(t *testing.T) {
		audit, err := service.FindAuditByID(1)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if audit.ID != testAudit.ID {
			t.Errorf("Expected audit ID %d, got %d", testAudit.ID, audit.ID)
		}
	})

	// Test finding a non-existing audit
	t.Run("Find non-existing audit", func(t *testing.T) {
		_, err := service.FindAuditByID(999)
		if err == nil {
			t.Error("Expected an error for non-existing audit, but got nil")
		}
	})
}

// TestAuditService_CreateAudit tests the functionality of the CreateAudit method.
func TestAuditService_CreateAudit(t *testing.T) {
	mockRepo := mock.NewMockAuditRepository()
	service := &AuditService{repo: mockRepo}

	// Define test cases
	tests := []struct {
		name       string
		tableName  string
		objectId   int
		comments   string
		createdBy  string
		beforeJson string
		afterJson  string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "Valid audit creation",
			tableName:  "vocab",
			objectId:   1,
			comments:   "Updating vocab entry",
			createdBy:  "tester",
			beforeJson: `{"LearningLang":"Hello", "FirstLang":"Hola"}`,
			afterJson:  `{"LearningLang":"Hello Updated", "FirstLang":"Hola Updated"}`,
			wantErr:    false,
		},
		{
			name:       "Invalid comments length",
			tableName:  "vocab",
			objectId:   1,
			comments:   string(make([]rune, 1001)), // 1001 characters
			createdBy:  "tester",
			beforeJson: `{"LearningLang":"Hello", "FirstLang":"Hola"}`,
			afterJson:  `{"LearningLang":"Hello Updated", "FirstLang":"Hola Updated"}`,
			wantErr:    true,
			errMsg:     "comments must be shorter than 1000 characters",
		},
		{
			name:       "Empty afterJson",
			tableName:  "vocab",
			objectId:   1,
			comments:   "This should still proceed",
			createdBy:  "tester",
			beforeJson: `{"LearningLang":"Hello", "FirstLang":"Hola"}`,
			afterJson:  "",
			wantErr:    false, // Assuming empty afterJson is considered valid for creation
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateAudit(tt.tableName, tt.objectId, tt.comments, tt.createdBy, tt.beforeJson, tt.afterJson)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: CreateAudit() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			} else if err != nil && err.Error() != tt.errMsg {
				t.Errorf("%s: CreateAudit() error = %v, wantErrMsg %v", tt.name, err, tt.errMsg)
			}
		})
	}
}

// TestAuditService_CreateVocabAudit tests the functionality of CreateVocabAudit method.
func TestAuditService_CreateVocabAudit(t *testing.T) {
	// Setup
	mockRepo := mock.NewMockAuditRepository()
	service := &AuditService{repo: mockRepo}

	beforeVocab := &mdl.Vocab{
		ID:           1,
		LearningLang: "English",
		FirstLang:    "Español",
		Created:      time.Now(),
	}

	afterVocab := &mdl.Vocab{
		ID:           1,
		LearningLang: "English Updated",
		FirstLang:    "Español Updated",
		Created:      time.Now(),
	}

	// Define test cases
	tests := []struct {
		name      string
		comments  string
		createdBy string
		before    *mdl.Vocab
		after     *mdl.Vocab
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Valid audit creation",
			comments:  "Updating vocab entry",
			createdBy: "tester",
			before:    beforeVocab,
			after:     afterVocab,
			wantErr:   false,
		},
		{
			name:      "Invalid comments length",
			comments:  string(make([]rune, 1001)), // 1001 characters
			createdBy: "tester",
			before:    beforeVocab,
			after:     afterVocab,
			wantErr:   true,
			errMsg:    "comments must be shorter than 1000 characters",
		},
		{
			name:      "After value is nil",
			comments:  "This should fail",
			createdBy: "tester",
			before:    beforeVocab,
			after:     nil,
			wantErr:   true,
			errMsg:    "after value for vocab is required",
		},
		{
			name:      "Before and after ID mismatch",
			comments:  "Mismatch IDs",
			createdBy: "tester",
			before:    &mdl.Vocab{ID: 2}, // Different ID
			after:     afterVocab,
			wantErr:   true,
			errMsg:    "audit before id 2 and after id 1 mismatch",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateVocabAudit(tt.comments, tt.createdBy, tt.before, tt.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateVocabAudit() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && err.Error() != tt.errMsg {
				t.Errorf("CreateVocabAudit() error = %v, wantErrMsg %v", err, tt.errMsg)
			}
		})
	}
}

// TestAuditService_CreateFixitAudit tests the functionality of the CreateFixitAudit method.
func TestAuditService_CreateFixitAudit(t *testing.T) {
	mockRepo := mock.NewMockAuditRepository()
	service := &AuditService{repo: mockRepo}

	beforeFixit := &mdl.Fixit{
		ID:        1,
		VocabID:   100,
		Status:    mdl.StatusType("pending"),
		FieldName: "Definition",
		Comments:  "Initial comment",
		CreatedBy: "tester",
		Created:   time.Now(),
	}

	afterFixit := &mdl.Fixit{
		ID:        1,
		VocabID:   100,
		Status:    mdl.StatusType("completed"),
		FieldName: "Definition",
		Comments:  "Updated comment",
		CreatedBy: "tester",
		Created:   time.Now(),
	}

	// Define test cases
	tests := []struct {
		name      string
		comments  string
		createdBy string
		before    *mdl.Fixit
		after     *mdl.Fixit
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Valid fixit audit creation",
			comments:  "Updating fixit entry",
			createdBy: "tester",
			before:    beforeFixit,
			after:     afterFixit,
			wantErr:   false,
		},
		{
			name:      "After value is nil",
			comments:  "This should fail",
			createdBy: "tester",
			before:    beforeFixit,
			after:     nil,
			wantErr:   true,
			errMsg:    "after value for fixit is required",
		},
		{
			name:      "Before and after ID mismatch",
			comments:  "Mismatch IDs",
			createdBy: "tester",
			before:    &mdl.Fixit{ID: 2}, // Different ID than afterFixit
			after:     afterFixit,
			wantErr:   true,
			errMsg:    "fixit before id 2 and after id 1 mismatch",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateFixitAudit(tt.comments, tt.createdBy, tt.before, tt.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: CreateFixitAudit() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			} else if err != nil && err.Error() != tt.errMsg {
				t.Errorf("%s: CreateFixitAudit() error = %v, wantErrMsg %v", tt.name, err, tt.errMsg)
			}
		})
	}
}
