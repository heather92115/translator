package mdl

import (
	"time"
)

// Audit represents a record of changes made to a database entity. It is designed
// to track modifications by storing the differences, along with metadata about
// the change, such as the entity affected, who made the change, and when.
//
// Fields:
//   - ID: The unique identifier for the audit record.
//   - ObjectID: The identifier of the entity that was changed.
//   - TableName: The name of the table where the entity resides.
//   - Diff: A representation of the changes made to the entity. This could be in any format
//     that best represents the difference, such as JSON.
//   - Before: The state of the entity before the changes were made, possibly serialized as a string.
//   - After: The state of the entity after the changes were made, possibly serialized as a string.
//   - Comments: Optional comments or notes about the changes made.
//   - CreatedBy: The identifier of the user or process that made the changes.
//   - Created: The timestamp when the audit record was created.
//
// This struct is typically used to populate an audit log, allowing for a historical
// review of changes for accountability and possibly restoration of previous states.
type Audit struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ObjectID  int       `json:"object_id" gorm:"not null"`
	TableName string    `json:"table_name" gorm:"not null"`
	Diff      string    `json:"diff"`   // Serialized representation of the differences
	Before    string    `json:"before"` // State before the changes
	After     string    `json:"after"`  // State after the changes
	Comments  string    `gorm:"default:''"`
	CreatedBy string    `json:"created_by" gorm:"not null"`
	Created   time.Time `json:"created" gorm:"not null;default:now()"`
}
