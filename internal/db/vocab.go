// Package db defines interfaces and implementations for interacting with
// entities in the database. It includes the VocabRepository interface, which outlines
// operations for querying and mutating Vocab records, and the SQLVocabRepository
// struct, which provides a concrete implementation of the VocabRepository using GORM.
//
// The SQLVocabRepository implementation leverages GORM to perform CRUD operations
// and complex queries on the database, abstracting the data layer away from the
// service layer to facilitate easier testing and maintenance.
package db

import (
	"fmt"
	"github.com/heather92115/verdure-admin/internal/mdl"
	"gorm.io/gorm"
	"log"
)

// VocabRepository defines the operations available for a Vocab entity.
type VocabRepository interface {
	FindVocabByID(id int) (*mdl.Vocab, error)
	FindVocabByLearningLang(learningLang string) (vocab *mdl.Vocab, err error)
	FindVocabs(learningCode string, hasFirst bool, limit int) (*[]mdl.Vocab, error)
	CreateVocab(vocab *mdl.Vocab) error
	UpdateVocab(vocab *mdl.Vocab) error
}

// SQLVocabRepository provides a GORM-based implementation of the VocabRepository interface.
type SQLVocabRepository struct {
	db *gorm.DB
}

// NewSqlVocabRepository initializes a new SQLVocabRepository with a database connection.
func NewSqlVocabRepository() (repo *SQLVocabRepository, err error) {
	db, err := GetConnection()
	if err != nil {
		return
	}

	repo = &SQLVocabRepository{db: db}

	return
}

// FindVocabByID retrieves a single Vocab record from the database using its primary ID.
//
// The function attempts to establish a database connection and then queries the Vocab table
// for a record matching the specified ID. It is designed to fetch exactly one record or return
// an error if the record does not exist or in case of a connection or query execution error.
//
// Parameters:
// - id: An integer representing the primary ID of the Vocab record to retrieve.
//
// Returns:
//   - *mdl.Vocab: A pointer to a Vocab struct representing the found record. If no record is found
//     or in case of an error, nil is returned.
//   - error: An error object detailing any issues encountered during the database connection
//     attempt or query execution. Errors could include connection failures, issues executing
//     the query, or the situation where no record is found matching the provided ID.
//     In cases where the operation succeeds and a record is found, nil is returned for the error.
//
// Usage example:
// vocab, err := FindVocabByID(123)
//
//	if err != nil {
//	    log.Printf("An error occurred: %v", err)
//	} else {
//		log.Printf("Retrieved vocab: %+v\n", vocab)
//	}
func (repo *SQLVocabRepository) FindVocabByID(id int) (vocab *mdl.Vocab, err error) {

	db, err := GetConnection()
	if err != nil {
		return
	}

	result := db.First(&vocab, id) // `First` method adds `WHERE id = ?` to the query
	if result.Error != nil {
		err = fmt.Errorf("error finding vocab with id %d: %v", id, result.Error)
	}

	return
}

// FindVocabByLearningLang retrieves a Vocab record from the database based on the learning language.
//
// This function searches the database for a Vocab record that matches the specified learning language string.
// The learning language is expected to be unique for each record, hence only one record should match the criteria.
// If the connection to the database cannot be established, or if no record is found matching the given learning language,
// the function returns an error detailing the issue encountered.
//
// Parameters:
// - learningLang: A string representing the learning language of the Vocab record to retrieve.
//
// Returns:
//   - *mdl.Vocab: A pointer to the retrieved Vocab record. If no record is found or in case of an error, nil is returned.
//   - error: An error object that details any issues encountered during the database connection attempt or query execution.
//     Possible errors include connection failures, issues executing the query, or the case where no record is found
//     matching the provided learning language. In cases where the operation succeeds, nil is returned for the error.
//
// Usage example:
// vocab, err := FindVocabByLearningLang("English")
//
//	if err != nil {
//	    log.Printf("An error occurred: %v", err)
//	} else {
//
//	    fmt.Printf("Retrieved vocab: %+v\n", vocab)
//	}
func (repo *SQLVocabRepository) FindVocabByLearningLang(learningLang string) (vocab *mdl.Vocab, err error) {
	db, err := GetConnection()
	if err != nil {
		return
	}

	// Use the `Where` method to specify the search condition
	result := db.Where("learning_lang = ?", learningLang).First(&vocab)
	if result.Error != nil {
		err = fmt.Errorf("error finding vocab with learning lang %s: %v", learningLang, result.Error)
	}

	return
}

// FindVocabs retrieves a list of Vocab records filtered by learning language code and
// the presence or absence of a first language translation. It limits the number of records
// returned based on the specified limit. The function filters records based on the learningCode
// provided and whether each record has a non-empty first language translation as indicated by
// the hasFirst parameter. If hasFirst is true, only records with a first language translation
// are included. If hasFirst is false, it returns records without a first language translation.
//
// Parameters:
//   - learningCode: The code of the learning language to filter records by.
//   - hasFirst: A boolean flag indicating whether to filter for records with (true) or
//     without (false) a first language translation.
//   - limit: The maximum number of records to return.
//
// Returns:
// - vocabs: A pointer to a slice of Vocab records matching the criteria, or nil if an error occurs.
// - err: An error object if an error occurs during the query execution, otherwise nil.
//
// Example of usage:
// vocabs, err := FindVocabs("es", true, 10)
//
//	if err != nil {
//	    log.Println("Error fetching vocabs:", err)
//	} else {
//	    for _, vocab := range *vocabs {
//	        fmt.Println(vocab)
//	    }
//	}
func (repo *SQLVocabRepository) FindVocabs(learningCode string, hasFirst bool, limit int) (vocabs *[]mdl.Vocab, err error) {
	db, err := GetConnection()
	if err != nil {
		return
	}

	vocabs = &[]mdl.Vocab{}

	query := db.Limit(limit)

	// Filter by LearningLangCode
	query = query.Where("learning_lang_code = ?", learningCode)

	// Conditionally filter based on the presence/absence of FirstLang
	if hasFirst {
		// Find records WITH a FirstLang value
		query = query.Where("first_lang != '' AND first_lang IS NOT NULL")
	} else {
		// Look for records where FirstLang is NOT present (assumed to be empty string or NULL)
		query = query.Where("first_lang = '' OR first_lang IS NULL")
	}

	// Execute the query
	err = query.Find(vocabs).Error
	if err != nil {
		log.Printf("Error finding %d vocab records with learning code '%s': %v", limit, learningCode, err)
	}

	return
}

// CreateVocab inserts a new Vocab record into the database.
// It establishes a database connection, then attempts to insert the provided Vocab instance.
// Returns an error if the database connection fails or if the insert operation encounters an error.
func (repo *SQLVocabRepository) CreateVocab(vocab *mdl.Vocab) error {
	db, err := GetConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to the db, error: %v", err)
	}

	result := db.Create(vocab)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateVocab updates an existing Vocab record in the database.
// It establishes a database connection, then attempts to update the Vocab instance based on its ID.
// Returns an error if the database connection fails or if the update operation encounters an error.
func (repo *SQLVocabRepository) UpdateVocab(vocab *mdl.Vocab) error {
	db, err := GetConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to the db, error: %v", err)
	}

	result := db.Save(vocab)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
