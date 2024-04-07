// Package db defines interfaces and implementations for interacting with
// entities in the database. It includes the FixitRepository interface, which outlines
// operations for querying and mutating Fixit records, and the SQLFixitRepository
// struct, which provides a concrete implementation of the FixitRepository using GORM.
//
// The SQLFixitRepository implementation leverages GORM to perform CRUD operations
// and complex queries on the database, abstracting the data layer away from the
// service layer to facilitate easier testing and maintenance.
package db

import (
	"fmt"
	"github.com/heather92115/verdure-admin/internal/mdl"
	"gorm.io/gorm"
	"log"
)

// FixitRepository defines the operations available for a Fixit entity.
type FixitRepository interface {
	FindFixitByID(id int) (*mdl.Fixit, error)
	FindFixits(
		status mdl.StatusType,
		vocabID int,
		duration *mdl.Duration,
		limit int) (fixits *[]mdl.Fixit, err error)

	CreateFixit(Fixit *mdl.Fixit) error
	UpdateFixit(fixit *mdl.Fixit) error
}

// SQLFixitRepository provides a GORM-based implementation of the FixitRepository interface.
type SQLFixitRepository struct {
	db *gorm.DB
}

// NewSqlFixitRepository initializes a new SQLFixitRepository with a database connection.
func NewSqlFixitRepository() (repo *SQLFixitRepository, err error) {
	db, err := GetConnection()
	if err != nil {
		return
	}

	repo = &SQLFixitRepository{db: db}

	return
}

// FindFixitByID retrieves a single Fixit record from the database using its primary ID.
//
// The function attempts to establish a database connection and then queries the Fixit table
// for a record matching the specified ID. It is designed to fetch exactly one record or return
// an error if the record does not exist or in case of a connection or query execution error.
//
// Parameters:
// - id: An integer representing the primary ID of the Fixit record to retrieve.
//
// Returns:
//   - *mdl.Fixit: A pointer to a Fixit struct representing the found record. If no record is found
//     or in case of an error, nil is returned.
//   - error: An error object detailing any issues encountered during the database connection
//     attempt or query execution. Errors could include connection failures, issues executing
//     the query, or the situation where no record is found matching the provided ID.
//     In cases where the operation succeeds and a record is found, nil is returned for the error.
//
// Usage example:
// Fixit, err := FindFixitByID(123)
//
//	if err != nil {
//	    log.Printf("An error occurred: %v", err)
//	} else {
//		log.Printf("Retrieved Fixit: %+v\n", Fixit)
//	}
func (repo *SQLFixitRepository) FindFixitByID(id int) (fixit *mdl.Fixit, err error) {

	db, err := GetConnection()
	if err != nil {
		return
	}

	result := db.First(&fixit, id) // `First` method adds `WHERE id = ?` to the query
	if result.Error != nil {
		err = fmt.Errorf("error finding Fixit with id %d: %v", id, result.Error)
	}

	return
}

// FindFixits retrieves a slice of Fixit entities from the database that match the given criteria.
// It filters Fixits based on their status, associated vocab ID, and creation date within a specified duration.
// This method supports pagination through a 'limit' parameter, allowing clients to specify the maximum
// number of Fixits to retrieve.
//
// Parameters:
//   - status: A StatusType value to filter Fixits by their current status.
//   - vocabID: An integer representing the vocab ID. If greater than 0, the method filters Fixits associated with this vocab ID.
//   - duration: A pointer to a Duration struct specifying the start and end time for filtering Fixits based on their creation date.
//     Fixits created within this duration are included in the result. If nil, no time-based filtering is applied.
//   - limit: An integer defining the maximum number of Fixits to return. Useful for pagination.
//
// Returns:
// - A pointer to a slice of Fixit entities matching the provided criteria.
// - An error if there's a problem executing the database query.
//
// Example usage:
// fixits, err := fixitService.FindFixits(mdl.StatusType("pending"), 101, &mdl.Duration{Start: time.Now().Add(-7*24*time.Hour), End: time.Now()}, 10)
//
//	if err != nil {
//	    log.Printf("Error retrieving Fixits: %v", err)
//	} else {
//
//	    for _, fixit := range *fixits {
//	        fmt.Printf("Found Fixit: %+v\n", fixit)
//	    }
//	}
func (repo *SQLFixitRepository) FindFixits(
	status mdl.StatusType,
	vocabID int,
	duration *mdl.Duration,
	limit int) (fixits *[]mdl.Fixit, err error) {

	db, err := GetConnection()
	if err != nil {
		return
	}

	fixits = &[]mdl.Fixit{}

	query := db.Limit(limit)
	query = query.Where("status = ?", status)

	if vocabID > 0 {
		query = query.Where("vocab_id = ?", vocabID)
	}

	if duration != nil {
		query = query.Where("created >= ? and created <= ?", duration.Start, duration.End)
	}

	// Execute the query
	err = query.Find(fixits).Error
	if err != nil {
		log.Printf("Error finding %d Fixit records with: status %v, vocab id '%d', : %v", limit, status, vocabID, err)
	}

	return
}

// CreateFixit inserts a new Fixit record into the database.
// It establishes a database connection, then attempts to insert the provided Fixit instance.
// Returns an error if the database connection fails or if the insert operation encounters an error.
func (repo *SQLFixitRepository) CreateFixit(fixit *mdl.Fixit) error {
	db, err := GetConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to the db, error: %v", err)
	}

	result := db.Create(fixit)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateFixit updates an existing Fixit record into the database.
// It establishes a database connection, then attempts to find and update the provided Fixit instance.
// Returns an error if the database connection fails or if the update operation encounters an error.
func (repo *SQLFixitRepository) UpdateFixit(fixit *mdl.Fixit) error {
	db, err := GetConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to the db, error: %v", err)
	}

	result := db.Save(fixit)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
