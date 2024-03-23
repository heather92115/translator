package srv

import (
	"fmt"
	"github.com/heather92115/translator/internal/db"
	"github.com/heather92115/translator/internal/mdl"
	"regexp"
)

// VocabService handles business logic for Vocab entities.
type VocabService struct {
	repo         db.VocabRepository
	auditService AuditService
}

// NewVocabService creates a new instance of VocabService.
func NewVocabService() (*VocabService, error) {

	repo, err := db.NewSqlVocabRepository()
	if err != nil {
		return nil, err
	}

	auditService, err := NewAuditService()
	if err != nil {
		return nil, err
	}

	return &VocabService{repo: repo, auditService: *auditService}, nil
}

// FindVocabByID retrieves a single Vocab record by its primary ID.
//
// This method searches the database for a Vocab record corresponding to the specified ID.
// It logs the search attempt and delegates the actual database query to the repository layer.
// If the record is found, it is returned along with a nil error. If the record is not found
// or if any database errors occur, the function returns nil and the error respectively.
//
// Parameters:
// - id: The primary ID of the Vocab record to retrieve.
//
// Returns:
// - A pointer to the found mdl.Vocab record, or nil if no record is found or an error occurs.
// - An error if the retrieval fails due to a database error or the record does not exist.
//
// Usage example:
// vocab, err := vocabService.FindVocabByID(123)
//
//	if err != nil {
//	    log.Printf("Failed to find vocab with ID 123: %v", err)
//	} else {
//
//	    fmt.Printf("Found vocab: %+v\n", vocab)
//	}
func (s *VocabService) FindVocabByID(id int) (*mdl.Vocab, error) {
	return s.repo.FindVocabByID(id)
}

// FindVocabs retrieves a list of Vocab records from the database based on the specified criteria.
// It filters records by the learning language code and the presence of a first language translation,
// returning up to a specified limit of records.
//
// Parameters:
// - learningCode: The code of the learning language used to filter the Vocab records.
// - hasFirst: A boolean indicating whether to filter records that have (true) or lack (false) a first language translation.
// - limit: The maximum number of Vocab records to retrieve.
//
// Returns:
// - A pointer to a slice of mdl.Vocab records matching the search criteria.
// - An error if there's an issue retrieving the records from the database.
//
// Usage example:
// vocabs, err := vocabService.FindVocabs("es", true, 10)
//
//	if err != nil {
//	    log.Printf("Error finding vocabs: %v", err)
//	} else {
//
//	    for _, vocab := range *vocabs {
//	        fmt.Println(vocab)
//	    }
//	}
func (s *VocabService) FindVocabs(learningCode string, hasFirst bool, limit int) (vocabs *[]mdl.Vocab, err error) {
	return s.repo.FindVocabs(learningCode, hasFirst, limit)
}

// CreateVocab attempts to create a new Vocab record in the database.
// Before creation, it validates the Vocab struct's fields to ensure they meet defined criteria
// and checks if a Vocab record with the same learning language already exists in the database.
// If the record exists, or if validation fails, it returns an error.
//
// Parameters:
// - vocab: A pointer to the mdl.Vocab struct to be created.
//
// Returns:
//   - An error if validation fails, if a record with the same learning language already exists,
//     or if there's an error during the creation process. Returns nil if the record is successfully created.
//
// Usage example:
// err := vocabService.CreateVocab(&vocab)
//
//	if err != nil {
//	    log.Printf("Failed to create vocab: %v", err)
//	}
func (s *VocabService) CreateVocab(vocab *mdl.Vocab) (err error) {

	if err = validateVocab(vocab); err != nil {
		return
	}

	existing, err := s.repo.FindVocabByLearningLang(vocab.LearningLang)
	if err == nil && existing != nil {
		return fmt.Errorf("vocab with learning lang %s and id %d already exists", vocab.LearningLang, existing.ID)
	}

	err = s.repo.CreateVocab(vocab)

	err = s.auditService.CreateVocabAudit("created vocab", "sys", nil, vocab)

	return
}

func (s *VocabService) UpdateVocab(updating *mdl.Vocab) (vocab *mdl.Vocab, err error) {

	if err = validateVocabUpdate(updating); err != nil {
		return
	}

	before, err := s.repo.FindVocabByID(updating.ID)
	if err != nil {
		return
	} else if before == nil {
		err = fmt.Errorf("expected to find existing vocab with id %d", updating.ID)
		return
	}

	vocab = before.Clone()

	if vocab.Hint != updating.Hint ||
		vocab.Pos != updating.Pos ||
		vocab.Skill != updating.Skill ||
		vocab.FirstLang != updating.FirstLang ||
		vocab.Infinitive != updating.Infinitive ||
		vocab.Alternatives != updating.Alternatives ||
		vocab.NumLearningWords != updating.NumLearningWords {
		// Update allowed to change fields
		vocab.Hint = updating.Hint
		vocab.Pos = updating.Pos
		vocab.Skill = updating.Skill
		vocab.FirstLang = updating.FirstLang
		vocab.Infinitive = updating.Infinitive
		vocab.Alternatives = updating.Alternatives
		vocab.NumLearningWords = updating.NumLearningWords
	} else {
		return nil, fmt.Errorf("update for vocab %d has no changes", vocab.ID)
	}

	err = s.repo.UpdateVocab(vocab)
	if err != nil {
		return
	}

	err = s.auditService.CreateVocabAudit("updated vocab", "sys", before, vocab)

	return
}

const (
	maxLearningLangLen = 40
	maxFirstLangLen    = 40
	maxAlternativesLen = 255
	maxSkillLen        = 100
	maxInfinitiveLen   = 40
	maxPosLen          = 40
	maxHintLen         = 255
	errFmtStrLangCode  = "%s must consist of two lowercase letters"
)

// validateVocab checks the validity of a Vocab struct's fields against defined constraints.
// It ensures that string fields do not exceed their maximum lengths and do not contain characters
// potentially harmful in the context of HTML or SQL. Specifically, it validates the 'LearningLang',
// 'FirstLang', 'Alternatives', 'Skill', 'Infinitive', 'Pos', and 'Hint' fields for length and restricted
// characters, and ensures 'KnownLangCode' and 'LearningLangCode' match the expected pattern for language codes.
//
// Parameters:
// - vocab: A pointer to the Vocab struct to validate.
//
// Returns:
//   - An error if any validation fails, detailing the issue with the corresponding field.
//     If all fields pass validation, nil is returned.
//
// Usage example:
// err := validateVocab(&vocab)
//
//	if err != nil {
//	    log.Printf("Validation failed: %v", err)
//	}
func validateVocab(vocab *mdl.Vocab) error {

	if err := validateFieldContent(vocab.LearningLang, "Learning language", maxLearningLangLen); err != nil {

		return err
	} else if len(vocab.LearningLang) == 0 {
		return fmt.Errorf("learning lang field is required")
	}

	if err := validateVocabUpdate(vocab); err != nil {
		return err
	}

	// Validate language codes with a more specific pattern
	langCodePattern := regexp.MustCompile(`^[a-z]{2}$`)
	if !langCodePattern.MatchString(vocab.KnownLangCode) || !langCodePattern.MatchString(vocab.LearningLangCode) {
		return fmt.Errorf(errFmtStrLangCode, "Language codes")
	}

	return nil
}

// validateVocabUpdate checks the validity of a Vocab struct's fields in the context of an update against defined
// constraints. It ensures that string fields do not exceed their maximum lengths and do not contain characters
// potentially harmful in the context of HTML or SQL. Specifically, it validates the 'LearningLang',
// 'FirstLang', 'Alternatives', 'Skill', 'Infinitive', 'Pos', and 'Hint' fields for length and restricted
// characters, and ensures 'KnownLangCode' and 'LearningLangCode' match the expected pattern for language codes.
//
// Parameters:
// - vocab: A pointer to the Vocab struct to validate.
//
// Returns:
//   - An error if any validation fails, detailing the issue with the corresponding field.
//     If all fields pass validation, nil is returned.
//
// Usage example:
// err := validateVocabUpdate(&vocab)
//
//	if err != nil {
//	    log.Printf("Validation failed: %v", err)
//	}
func validateVocabUpdate(vocab *mdl.Vocab) error {

	if err := validateFieldContent(vocab.FirstLang, "First language", maxFirstLangLen); err != nil {
		return err
	}
	if err := validateFieldContent(vocab.Alternatives, "Alternatives", maxAlternativesLen); err != nil {
		return err
	}
	if err := validateFieldContent(vocab.Skill, "Skill", maxSkillLen); err != nil {
		return err
	}
	if err := validateFieldContent(vocab.Infinitive, "Infinitive", maxInfinitiveLen); err != nil {
		return err
	}
	if err := validateFieldContent(vocab.Pos, "Part of speech", maxPosLen); err != nil {
		return err
	}
	if err := validateFieldContent(vocab.Hint, "Hint", maxHintLen); err != nil {
		return err
	}

	return nil
}
