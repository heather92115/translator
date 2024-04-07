package srv

import (
	"fmt"
	"github.com/heather92115/verdure-admin/internal/db"
	"github.com/heather92115/verdure-admin/internal/mdl"
)

// FixitService handles business logic for Fixit entities.
type FixitService struct {
	repo         db.FixitRepository
	auditService AuditService
}

// NewFixitService creates a new instance of FixitService.
func NewFixitService() (*FixitService, error) {

	repo, err := db.NewSqlFixitRepository()
	if err != nil {
		return nil, err
	}

	auditService, err := NewAuditService()
	if err != nil {
		return nil, err
	}

	return &FixitService{repo: repo, auditService: *auditService}, nil
}

// FindFixitByID retrieves a single Fixit record by its primary ID.
//
// This method searches the database for a Fixit record corresponding to the specified ID.
// It logs the search attempt and delegates the actual database query to the repository layer.
// If the record is found, it is returned along with a nil error. If the record is not found
// or if any database errors occur, the function returns nil and the error respectively.
//
// Parameters:
// - id: The primary ID of the Fixit record to retrieve.
//
// Returns:
// - A pointer to the found mdl.Fixit record, or nil if no record is found or an error occurs.
// - An error if the retrieval fails due to a database error or the record does not exist.
//
// Usage example:
// fixit, err := fixitService.FindFixitByID(123)
//
//	if err != nil {
//	    log.Printf("Failed to find fixit with ID 123: %v", err)
//	} else {
//
//	    fmt.Printf("Found fixit: %+v\n", fixit)
//	}
func (s *FixitService) FindFixitByID(id int) (*mdl.Fixit, error) {
	return s.repo.FindFixitByID(id)
}

func (s *FixitService) FindFixits(
	status mdl.StatusType,
	vocabID int,
	duration *mdl.Duration,
	limit int) (fixits *[]mdl.Fixit, err error) {
	return s.repo.FindFixits(status, vocabID, duration, limit)
}

// CreateFixit attempts to create a new Fixit record in the database.
// Before creation, it validates the Fixit struct fields to ensure they meet defined criteria.
//
// Parameters:
// - fixit: A pointer to the mdl.Fixit struct to be created.
//
// Returns:
//   - An error if validation fails or if there's an error during the creation process. Returns nil if the record is successfully created.
//
// Usage example:
// err := fixitService.CreateFixit(&fixit)
//
//	if err != nil {
//	    log.Printf("Failed to create fixit: %v", err)
//	}
func (s *FixitService) CreateFixit(fixit *mdl.Fixit) (err error) {

	if err = validateFixit(fixit); err != nil {
		return
	}

	err = s.repo.CreateFixit(fixit)

	err = s.auditService.CreateFixitAudit("created fixit", "sys", nil, fixit)

	return
}

func (s *FixitService) UpdateFixit(updating *mdl.Fixit) (fixit *mdl.Fixit, err error) {

	if err = validateFixit(updating); err != nil {
		return
	}

	before, err := s.repo.FindFixitByID(updating.ID)
	if err != nil {
		return
	} else if before == nil {
		err = fmt.Errorf("expected to find existing fixit with id %d", updating.ID)
		return
	}

	fixit = before.Clone()

	// Update allowed to change fields
	if fixit.Status != updating.Status || fixit.FieldName != updating.FieldName || fixit.Comments != updating.Comments {
		fixit.Status = updating.Status
		fixit.FieldName = updating.FieldName
		fixit.Comments = updating.Comments
	} else {
		return nil, fmt.Errorf("update for fixit %d has no changes", fixit.ID)
	}

	err = s.repo.UpdateFixit(fixit)
	if err != nil {
		return
	}

	err = s.auditService.CreateFixitAudit("updated fixit", "sys", before, fixit)

	return
}

const (
	maxFixitFieldNameLen = 40
	maxFixitCommitLen    = 2000
)

// validateFixit checks the validity of a Fixit entity's fields against specified constraints.
// This function ensures that the 'FieldName' and 'Comments' of a Fixit do not exceed their
// maximum allowed lengths, adhering to the defined maximum length constants.
//
// Parameters:
// - fixit: A pointer to the Fixit struct whose fields are to be validated.
//
// Returns:
//   - An error if any field fails the validation checks, specifying the nature of the failure.
//     If all fields pass the validation, nil is returned.
//
// Usage:
// The function is typically used before creating or updating a Fixit entity in the database
// to ensure that the data conforms to expected standards. This prevents data integrity issues
// and ensures consistency across the application.
//
// Example:
// err := validateFixit(fixit)
//
//	if err != nil {
//	    // Handle the validation error
//	    log.Printf("Validation failed: %v", err)
//	}
//
// This function relies on validateFieldContent to perform the actual validation of each field,
// using 'maxFixitFieldNameLen' and 'maxFixitCommitLen' as the maximum length constraints for
// the 'FieldName' and 'Comments' fields, respectively.
func validateFixit(fixit *mdl.Fixit) error {

	if err := validateFieldContent(fixit.FieldName, "Field Name", maxFixitFieldNameLen); err != nil {
		return err
	}
	if err := validateFieldContent(fixit.Comments, "Commits", maxFixitCommitLen); err != nil {
		return err
	}

	return nil
}
