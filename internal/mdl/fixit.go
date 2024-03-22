package mdl

import (
	"encoding/json"
	"fmt"
	"time"
)

type StatusType string

const (
	Pending    StatusType = "pending"
	InProgress StatusType = "in_progress"
	Completed  StatusType = "completed"
)

// Fixit represents a correction or modification suggestion for a Vocab entry.
// It is used to track proposed changes or enhancements to vocabulary records,
// including status tracking, the field targeted for correction, and commentary
// on the suggestion. Each Fixit is associated with a specific Vocab record.
//
// Fields:
//   - ID: The unique identifier for the Fixit record, automatically incremented.
//   - VocabID: The ID of the associated Vocab record that this Fixit suggestion pertains to.
//   - Status: The current status of the Fixit suggestion, represented as a StatusType
//     (e.g., Pending, Approved, Rejected). The specific status types are defined by the
//     StatusType type and are stored in the database as a 'status_type' enum.
//   - FieldName: The name of the field in the Vocab record that the Fixit suggestion
//     aims to correct or modify. This could refer to any textual field within a Vocab
//     record that is subject to correction, such as 'LearningLang', 'FirstLang', etc.
//   - Comments: Optional commentary or rationale provided by the creator of the Fixit
//     suggestion, offering context or justification for the proposed change.
//   - CreatedBy: The identifier (e.g., username or user ID) of the user who created
//     the Fixit suggestion. This field is used to track who is responsible for the
//     suggestion and to enable follow-up or attribution.
//   - Created: The timestamp when the Fixit record was created, automatically set to
//     the current date and time when the record is created in the database.
//
// This struct is typically used within an application that allows users to suggest
// edits or improvements to vocabulary entries, facilitating collaborative refinement
// and accuracy in a vocabulary management system.
type Fixit struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	VocabID   int        `json:"vocab_id" gorm:"foreignKey:Vocab"`
	Status    StatusType `gorm:"type:status_type"`
	FieldName string     `json:"field_name" gorm:"default:''"`
	Comments  string     `gorm:"default:''"`
	CreatedBy string     `json:"created_by" gorm:"not null"`
	Created   time.Time  `json:"created" gorm:"index:idx_fixit_created,not null;default:now()"`
}

// JSON Creates a JSON string from a Fixit object.
func (o *Fixit) JSON() string {
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}
	return string(b)
}

// Clone creates a deep copy of the Fixit object and returns a new instance containing the same data.
// This method is useful for creating distinct instances of Fixit objects when copying is needed
// to prevent modifications to the original object from affecting the copied one.
//
// Returns:
// - A pointer to a new Fixit instance that is a clone of the original.
func (f *Fixit) Clone() *Fixit {
	return &Fixit{
		ID:        f.ID,
		VocabID:   f.VocabID,
		Status:    f.Status,
		FieldName: f.FieldName,
		Comments:  f.Comments,
		CreatedBy: f.CreatedBy,
		Created:   f.Created,
	}
}
