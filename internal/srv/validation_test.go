package srv

import (
	"fmt"
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
