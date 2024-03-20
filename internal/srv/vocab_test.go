package srv

import (
	"fmt"
	"github.com/heather92115/translator/internal/mdl"
	"strings"
	"testing"
)

func TestValidateFieldContent(t *testing.T) {

	tests := []struct {
		name       string
		fieldValue string
		fieldName  string
		maxLength  int
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "Valid input",
			fieldValue: "test",
			fieldName:  "username",
			maxLength:  50,
			wantErr:    false,
		},
		{
			name:       "Input exceeds max length",
			fieldValue: "thisisareallylonginputthatexceedsthemaximumallowedlengthforthisfield",
			fieldName:  "username",
			maxLength:  50,
			wantErr:    true,
			errMsg:     fmt.Sprintf(errFmtStrLen, "username", 50),
		},
		{
			name:       "Input contains invalid characters",
			fieldValue: "test<script>",
			fieldName:  "username",
			maxLength:  50,
			wantErr:    true,
			errMsg:     "username contains invalid characters",
		},
		{
			name:       "Empty input",
			fieldValue: "",
			fieldName:  "username",
			maxLength:  50,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldContent(tt.fieldValue, tt.fieldName, tt.maxLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFieldContent() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && err.Error() != tt.errMsg {
				t.Errorf("validateFieldContent() error = %v, wantErrMsg %v", err, tt.errMsg)
			}
		})
	}
}

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
