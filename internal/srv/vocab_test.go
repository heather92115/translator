package srv

import (
	"fmt"
	"github.com/heather92115/translator/internal/db/mock"
	"github.com/heather92115/translator/internal/mdl"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestValidateVocab(t *testing.T) {
	validLangCode := "en"
	invalidLangCode := "123"

	tests := []struct {
		name    string
		vocab   mdl.Vocab
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid input",
			vocab: mdl.Vocab{
				LearningLang:     "English",
				FirstLang:        "Español",
				Alternatives:     "Alternative text",
				Skill:            "Beginner",
				Infinitive:       "to be",
				Pos:              "verb",
				Hint:             "A hint",
				KnownLangCode:    validLangCode,
				LearningLangCode: validLangCode,
			},
			wantErr: false,
		},
		{
			name: "Invalid LearningLang length",
			vocab: mdl.Vocab{
				LearningLang:     strings.Repeat("a", maxLearningLangLen+1),
				FirstLang:        "Español",
				Alternatives:     "Alternative text",
				Skill:            "Beginner",
				Infinitive:       "to be",
				Pos:              "verb",
				Hint:             "A hint",
				KnownLangCode:    validLangCode,
				LearningLangCode: validLangCode,
			},
			wantErr: true,
			errMsg:  fmt.Sprintf(errFmtStrLen, "Learning language", maxLearningLangLen),
		},
		{
			name: "Invalid KnownLangCode format",
			vocab: mdl.Vocab{
				LearningLang:     "English",
				FirstLang:        "Español",
				Alternatives:     "Alternative text",
				Skill:            "Beginner",
				Infinitive:       "to be",
				Pos:              "verb",
				Hint:             "A hint",
				KnownLangCode:    invalidLangCode,
				LearningLangCode: validLangCode,
			},
			wantErr: true,
			errMsg:  fmt.Sprintf(errFmtStrLangCode, "Language codes"),
		},
		{
			name: "Missing LearningLang format",
			vocab: mdl.Vocab{
				LearningLang:     "",
				FirstLang:        "Español",
				Alternatives:     "Alternative text",
				Skill:            "Beginner",
				Infinitive:       "to be",
				Pos:              "verb",
				Hint:             "A hint",
				KnownLangCode:    invalidLangCode,
				LearningLangCode: validLangCode,
			},
			wantErr: true,
			errMsg:  "learning lang field is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVocab(&tt.vocab)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateVocab() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && err.Error() != tt.errMsg {
				t.Errorf("validateVocab() error = %v, wantErrMsg %v", err, tt.errMsg)
			}
		})
	}
}

// TestVocabService_FindVocabByID tests the functionality of FindVocabByID method.
func TestVocabService_FindVocabByID(t *testing.T) {
	// Initialize the mock repositories
	mockVocabRepo := mock.NewMockVocabRepository()
	mockAuditRepo := mock.NewMockAuditRepository()
	mockAuditService := &AuditService{repo: mockAuditRepo}

	// Create an instance of VocabService with mocks
	vocabService := VocabService{
		repo:         mockVocabRepo,
		auditService: *mockAuditService,
	}

	// Seed the mock repository with a test vocab
	testVocab := &mdl.Vocab{
		ID:           123,
		LearningLang: "están",
		FirstLang:    "they are",
		Created:      time.Now(),
	}
	_ = mockVocabRepo.CreateVocab(testVocab)

	// Execute the test
	vocab, err := vocabService.FindVocabByID(123)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if vocab.ID != testVocab.ID {
		t.Errorf("Expected vocab ID %d, got %d", testVocab.ID, vocab.ID)
	}

	// Test for a non-existing vocab
	_, err = vocabService.FindVocabByID(999) // Assuming 999 is a non-existing ID
	if err == nil {
		t.Error("Expected an error for non-existing vocab, but got nil")
	}
}

// TestVocabService_FindVocabs tests the functionality of FindVocabs method.
func TestVocabService_FindVocabs(t *testing.T) {
	// Initialize the mock repositories
	mockVocabRepo := mock.NewMockVocabRepository()

	// Create an instance of VocabService with mocks
	vocabService := VocabService{
		repo: mockVocabRepo,
	}

	// Seed the mock repository with test vocabs
	testVocab1 := &mdl.Vocab{
		ID:               1,
		LearningLang:     "hola",
		FirstLang:        "hello",
		LearningLangCode: "es",
	}
	testVocab2 := &mdl.Vocab{
		ID:               2,
		LearningLang:     "desafortunadamente",
		FirstLang:        "",
		LearningLangCode: "es",
	}
	_ = mockVocabRepo.CreateVocab(testVocab1)
	_ = mockVocabRepo.CreateVocab(testVocab2)

	// Define test cases
	tests := []struct {
		name           string
		learningCode   string
		hasFirst       bool
		limit          int
		expectedVocabs []mdl.Vocab
	}{
		{
			name:           "Find es Vocabs",
			learningCode:   "es",
			hasFirst:       true,
			limit:          10,
			expectedVocabs: []mdl.Vocab{*testVocab1},
		},
		{
			name:           "Find es with no FirstLang",
			learningCode:   "es",
			hasFirst:       false,
			limit:          10,
			expectedVocabs: []mdl.Vocab{*testVocab2},
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vocabs, err := vocabService.FindVocabs(tt.learningCode, tt.hasFirst, tt.limit)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(*vocabs, tt.expectedVocabs) {
				t.Errorf("Expected vocabs %+v, got %+v", tt.expectedVocabs, *vocabs)
			}
		})
	}
}

// TestVocabService_CreateVocab tests the functionality of CreateVocab method.
func TestVocabService_CreateVocab(t *testing.T) {
	// Setup
	mockVocabRepo := mock.NewMockVocabRepository()
	mockAuditRepo := mock.NewMockAuditRepository()
	mockAuditService := &AuditService{repo: mockAuditRepo}

	vocabService := VocabService{
		repo:         mockVocabRepo,
		auditService: *mockAuditService,
	}

	// Test cases
	tests := []struct {
		name    string
		vocab   *mdl.Vocab
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful vocab creation",
			vocab: &mdl.Vocab{
				ID:               1,
				LearningLang:     "hola",
				FirstLang:        "hello",
				Created:          time.Now(),
				LearningLangCode: "es",
				KnownLangCode:    "en",
			},
			wantErr: false,
		},
		{
			name: "Duplicate learning language",
			vocab: &mdl.Vocab{
				ID:               2,
				LearningLang:     "desafortunadamente",
				FirstLang:        "unfortunately",
				Created:          time.Now(),
				LearningLangCode: "es",
				KnownLangCode:    "en",
			},
			wantErr: true,
			errMsg:  "vocab with learning lang desafortunadamente and id 2 already exists",
		},
		{
			name: "Invalid vocab - missing learning language",
			vocab: &mdl.Vocab{
				ID:               3,
				FirstLang:        "hello",
				Created:          time.Now(),
				LearningLangCode: "es",
				KnownLangCode:    "en",
			},
			wantErr: true,
			errMsg:  "learning lang field is required",
		},
	}

	// Seed initial vocab for testing duplicate scenario
	_ = mockVocabRepo.CreateVocab(&mdl.Vocab{
		ID:               2,
		LearningLang:     "desafortunadamente",
		FirstLang:        "unfortunately",
		Created:          time.Now(),
		LearningLangCode: "es",
		KnownLangCode:    "en",
	})

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vocabService.CreateVocab(tt.vocab)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateVocab() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && err.Error() != tt.errMsg {
				t.Errorf("CreateVocab() error = %v, wantErrMsg %v", err, tt.errMsg)
			}
		})
	}
}

// TestVocabService_UpdateVocab tests the functionality of UpdateVocab method.
func TestVocabService_UpdateVocab(t *testing.T) {
	// Setup
	mockVocabRepo := mock.NewMockVocabRepository()
	mockAuditRepo := mock.NewMockAuditRepository()
	mockAuditService := &AuditService{repo: mockAuditRepo}

	vocabService := VocabService{
		repo:         mockVocabRepo,
		auditService: *mockAuditService,
	}

	// Seed the mock repository with a vocab for update tests
	existingVocab := &mdl.Vocab{
		ID:               1,
		LearningLang:     "hola",
		FirstLang:        "hello",
		Created:          time.Now(),
		LearningLangCode: "es",
		KnownLangCode:    "en",
	}
	_ = vocabService.CreateVocab(existingVocab)

	// Define test cases
	tests := []struct {
		name    string
		vocab   *mdl.Vocab // Vocab to update
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful vocab update",
			vocab: &mdl.Vocab{
				ID:               1, // Assumes ID 1 exists
				LearningLang:     "hola",
				FirstLang:        "hello updated",
				LearningLangCode: "es",
				KnownLangCode:    "en",
			},
			wantErr: false,
		},
		{
			name: "Update non-existing vocab",
			vocab: &mdl.Vocab{
				ID:               999, // Assumes ID 999 does not exist
				LearningLang:     "adios",
				FirstLang:        "goodbye",
				LearningLangCode: "es",
				KnownLangCode:    "en",
			},
			wantErr: true,
			errMsg:  "error finding vocab with id 999",
		},
		{
			name: "Invalid vocab - missing learning language",
			vocab: &mdl.Vocab{
				ID:               1,
				FirstLang:        "hello again",
				LearningLangCode: "es",
				KnownLangCode:    "en",
			},
			wantErr: true,
			errMsg:  "learning lang field is required",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedVocab, err := vocabService.UpdateVocab(tt.vocab)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVocab() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && !tt.wantErr && updatedVocab.FirstLang != tt.vocab.FirstLang {
				t.Errorf("UpdateVocab() failed to update fields properly. Expected firstLang %v, got %v", tt.vocab.FirstLang, updatedVocab.FirstLang)
			} else if err != nil && err.Error() != tt.errMsg {
				t.Errorf("UpdateVocab() error = %v, wantErrMsg %v", err, tt.errMsg)
			}
		})
	}
}
