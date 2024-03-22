package srv

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

const (
	errFmtStrLen = "%s must be shorter than %d characters"
)

// validateFieldContent checks a string field's content for compliance with specified length and character restrictions.
// It ensures the field does not exceed a maximum length and does not contain characters that could be used for XSS or injection attacks.
// This function is intended for basic validation and sanitization of input fields to prevent common security vulnerabilities.
//
// Parameters:
// - fieldValue: The content of the field to validate.
// - fieldName: The name of the field, used in the error message to identify the field with invalid content.
// - maxLength: The maximum allowed length of the field content in Unicode code points.
//
// Returns:
//   - An error if the field content exceeds the maxLength or contains restricted characters, specifying the nature of the validation failure.
//     Returns nil if the field content passes all validation checks.
//
// Usage example:
// err := validateFieldContent(userInput, "username", 50)
//
//	if err != nil {
//	    log.Printf("Validation error: %v", err)
//	}
func validateFieldContent(fieldValue, fieldName string, maxLength int) error {
	if utf8.RuneCountInString(fieldValue) > maxLength {
		return fmt.Errorf(errFmtStrLen, fieldName, maxLength)
	}
	// Example basic check against common XSS/injection patterns. Expand as necessary.
	if strings.ContainsAny(fieldValue, "<>") && strings.ContainsAny(fieldValue, "\"/") {
		log.Printf("Validation error on fieldName %s, fieldValue %s ", fieldName, fieldValue)

		return fmt.Errorf("%s contains invalid characters", fieldName)
	}
	return nil
}
