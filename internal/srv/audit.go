package srv

import (
	"encoding/json"
	"fmt"
	"github.com/heather92115/translator/internal/db"
	"github.com/heather92115/translator/internal/mdl"
	"sort"
)

// AuditService handles business logic for Audit entities.
type AuditService struct {
	repo db.AuditRepository
}

// NewAuditService creates a new instance of AuditService.
func NewAuditService() (*AuditService, error) {

	repo, err := db.NewSqlAuditRepository()
	if err != nil {
		return nil, err
	}

	return &AuditService{repo: repo}, nil
}

// FindAuditByID retrieves a single Audit record by its primary ID.
//
// This method searches the database for an Audit record corresponding to the specified ID.
// It logs the search attempt and delegates the actual database query to the repository layer.
// If the record is found, it is returned along with a nil error. If the record is not found
// or if any database errors occur, the function returns nil and the error respectively.
//
// Parameters:
// - id: The primary ID of the Audit record to retrieve.
//
// Returns:
// - A pointer to the found mdl.Audit record, or nil if no record is found or an error occurs.
// - An error if the retrieval fails due to a database error or the record does not exist.
//
// Usage example:
// Audit, err := AuditService.FindAuditByID(123)
//
//	if err != nil {
//	    log.Printf("Failed to find Audit with ID 123: %v", err)
//	} else {
//
//	    fmt.Printf("Found Audit: %+v\n", Audit)
//	}
func (s *AuditService) FindAuditByID(id int) (*mdl.Audit, error) {

	return s.repo.FindAuditByID(id)
}

// FindAudits retrieves a slice of Audit records filtered based on the provided criteria.
// It filters audits by the name of the table, a duration within which the audits were created,
// and limits the number of returned audits to the specified limit. This method is useful
// for fetching audit records for specific database tables and time ranges, supporting
// both broad and narrow queries.
//
// Parameters:
//   - tableName: The name of the database table for which to retrieve audit records.
//     If an empty string is provided, audit records for all tables are considered.
//   - duration: A pointer to an mdl.Duration struct specifying the start and end time
//     for the time range filter. If nil, no time-based filtering is applied.
//   - limit: The maximum number of audit records to retrieve. A limit of 0 indicates
//     no limit, thus all matching records are returned.
//
// Returns:
//   - A pointer to a slice of mdl.Audit structs containing the retrieved audit records.
//     Returns nil if an error occurs during query execution.
//   - An error if there is an issue establishing a connection to the database, executing
//     the query, or applying the specified filters. Returns nil if the query is
//     successful without errors.
//
// Example usage:
// audits, err := auditService.FindAudits("users", &mdl.Duration{Start: startTime, End: endTime}, 10)
//
//	if err != nil {
//	    log.Printf("Error retrieving audits: %v", err)
//	} else {
//
//	    for _, audit := range *audits {
//	        fmt.Printf("Audit ID: %d, Table: %s\n", audit.ID, audit.TableName)
//	    }
//	}
func (s *AuditService) FindAudits(tableName string, duration *mdl.Duration, limit int) (Audits *[]mdl.Audit, err error) {
	return s.repo.FindAudits(tableName, duration, limit)
}

// CreateVocabAudit records an audit trail for vocabulary modifications. This function
// is called after creating or updating a vocabulary entry to log the changes made.
// It validates the length of the comments, checks the integrity of the before and after
// states of the vocab entry, and then creates an audit record with the provided information.
//
// Parameters:
//   - comments: A string containing comments about the changes made. This field is validated
//     to ensure it does not exceed 1000 characters.
//   - createdBy: The identifier of the user or system that made the changes. This could be a user
//     ID or a system name.
//   - before: A pointer to a Vocab struct representing the state of the vocabulary entry before
//     the changes. This parameter can be nil if the audit is for a newly created entry.
//   - after: A pointer to a Vocab struct representing the state of the vocabulary entry after
//     the changes. This parameter must not be nil.
//
// Returns:
//   - An error if validation fails, if the 'after' parameter is nil, if there is a mismatch between
//     the IDs of the 'before' and 'after' states, or if there is an error creating the audit record
//     in the repository. Returns nil if the audit record is successfully created.
//
// This function handles the serialization of the 'before' and 'after' Vocab states to JSON for
// logging purposes and calculates the diff between these states if both are provided. The diff,
// along with the comments, the ID of the modified entry, and the creator's identifier, are stored
// in an Audit struct and persisted to the repository.
//
// Example usage:
// err := auditService.CreateVocabAudit("Updated definition", "admin_user", beforeVocab, afterVocab)
//
//	if err != nil {
//	    log.Printf("Failed to create vocab audit: %v", err)
//	}
func (s *AuditService) CreateVocabAudit(comments string, createdBy string, before *mdl.Vocab, after *mdl.Vocab) (err error) {

	// validate the comments
	if err = validateFieldContent(comments, "comments", 1000); err != nil {
		return err
	}

	if after == nil {
		err = fmt.Errorf("after value is required")
		return
	}

	if before != nil && before.ID != after.ID {
		err = fmt.Errorf("audit before id %d and after id %d mismatch", before.ID, after.ID)
		return
	}
	afterJson := after.JSON()
	diff := ""

	beforeJson := ""
	if before != nil {
		beforeJson = before.JSON()
		diff = CompareJSON(beforeJson, afterJson)
	}

	audit := mdl.Audit{
		TableName: "vocab",
		ObjectID:  after.ID,
		Comments:  comments,
		Before:    beforeJson,
		After:     afterJson,
		Diff:      diff,
		CreatedBy: createdBy,
	}

	err = s.repo.CreateAudit(&audit)

	return
}

type DiffResult struct {
	Key    string      `json:"key"`
	Before interface{} `json:"before"`
	After  interface{} `json:"after"`
}

// CompareJSON takes two JSON strings as input and compares them to find any differences.
// It identifies keys that are present in one JSON object but not the other, and keys
// with differing values between the two JSON objects. The comparison is recursive,
// so nested objects are fully explored for differences as well.
//
// This function leverages findDiffs internally to perform the actual comparison and
// generate a slice of DiffResult structs representing the detected differences. Each
// DiffResult includes the key (or full key path for nested structures) along with the
// values before and after the change. For keys that are added or removed, the
// corresponding before or after value is provided, if applicable.
//
// Parameters:
// - jsonStr1: The first JSON string to be compared.
// - jsonStr2: The second JSON string to be compared.
//
// Returns:
//   - A JSON string representing a slice of DiffResult structs. Each DiffResult struct
//     includes the key, and, when applicable, the values before and after the change.
//     The returned JSON string is ready to be logged, displayed, or processed further
//     to analyze the differences between the two input JSON strings.
//
// Example usage:
// jsonStr1 := `{"name": "John", "age": 30}`
// jsonStr2 := `{"name": "Jane", "age": 31}`
// diffsJSON := CompareJSON(jsonStr1, jsonStr2)
// fmt.Println(diffsJSON)
//
// This function is particularly useful for debugging, logging changes, or comparing
// JSON representations of data structures to understand how they differ.
func CompareJSON(jsonStr1, jsonStr2 string) string {
	var obj1, obj2 map[string]interface{}

	_ = json.Unmarshal([]byte(jsonStr1), &obj1)
	_ = json.Unmarshal([]byte(jsonStr2), &obj2)

	diffs := findDiffs(obj1, obj2, "")

	diffJSON, _ := json.Marshal(diffs)
	return string(diffJSON)
}

// findDiffs compares two maps of string keys to interface{} values and returns a slice of DiffResult
// indicating the differences between them. Differences include keys that are present in one map
// but not the other (indicating addition or removal) and keys with differing values between the
// two maps. The function also recursively compares nested maps to identify deep differences.
//
// The path parameter is used to keep track of the nested level during recursive comparisons,
// allowing the function to accurately report the full key path of any differences found.
//
// Parameters:
//   - a: The first map to be compared.
//   - b: The second map to be compared.
//   - path: A string representing the current path in the nested structure, used for tracking
//     differences in nested maps. It should be an empty string when called for the top-level comparison.
//
// Returns:
//   - A sorted slice of DiffResult structs, each representing a detected difference. Differences are
//     sorted alphabetically by the full key path for easier readability and analysis.
//
// DiffResult structs include the key (or full key path for nested structures), and, when applicable,
// the values before and after the change. For keys that are added or removed, the corresponding
// before or after value is included, if applicable.
//
// Example usage:
// a := map[string]interface{}{"name": "John", "age": 30, "details": map[string]interface{}{"city": "New York"}}
// b := map[string]interface{}{"name": "Jane", "age": 30, "details": map[string]interface{}{"city": "Boston"}}
// diffs := findDiffs(a, b, "")
//
//	for _, diff := range diffs {
//	    fmt.Println(diff)
//	}
//
// This function is useful for debugging, logging, or otherwise needing to understand the
// differences between two map representations, perhaps of JSON objects or similar data structures.
func findDiffs(a, b map[string]interface{}, path string) (diffs []DiffResult) {
	for key, aValue := range a {
		bValue, exists := b[key]
		fullKey := fmt.Sprintf("%s%s", path, key)
		if !exists {
			// Key removed or added
			description := fmt.Sprintf("'%s' removed", fullKey)
			if path == "" { // Direct comparison implies key was in 'a' but not 'b', indicating removal
				diffs = append(diffs, DiffResult{Key: description})
			} else { // When called with 'b' as 'a', this indicates addition
				diffs = append(diffs, DiffResult{Key: fmt.Sprintf("'%s' added", fullKey)})
			}
			continue
		}

		if aValueTyped, ok := aValue.(map[string]interface{}); ok {
			if bValueTyped, ok := bValue.(map[string]interface{}); ok {
				subDiffs := findDiffs(aValueTyped, bValueTyped, fullKey+".")
				diffs = append(diffs, subDiffs...)
			}
		} else {
			if aValue != bValue {
				diffs = append(diffs, DiffResult{
					Key:    fmt.Sprintf("'%s'", fullKey),
					Before: aValue,
					After:  bValue,
				})
			}
		}
	}

	// Sort the diffs slice alphabetically by the Key field.
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Key < diffs[j].Key
	})

	return diffs
}
