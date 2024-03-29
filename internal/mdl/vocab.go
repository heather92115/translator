package mdl

import (
	"encoding/json"
	"fmt"
	"time"
)

// Vocab represents a vocabulary item in a language learning application.
//
// This struct is used to store information about words or phrases that users are learning,
// including translations, alternatives, and metadata to assist in the learning process.
//
// Fields:
// - ID: Primary key used to uniquely identify a vocabulary item in the data layer.
// - LearningLang: The word or phrase in the language being learned.
// - FirstLang: The translation of the word or phrase into the user's first language, used as a prompt.
// - Created: Timestamp when the vocabulary item was created. It is typically set automatically to the current time.
// - Alternatives: Optional. Additional correct answers or variations in the learning language.
// - Skill: Optional. The skill or category associated with the vocabulary item, used for organizing content.
// - Infinitive: Optional. For verbs, the infinitive form of the word. Empty for non-verb vocabulary items.
// - Pos: Optional. The part of speech of the vocabulary item, aiding in the application of grammatical rules.
// - Hint: Optional. A hint provided to assist users in translating the word or phrase.
// - NumLearningWords: The number of words contained in the `learning_lang` field, calculated for analytical purposes.
// - KnownLangCode: Language code for the known language.
// - LearningLangCode: Language code for the learning language.
//
// Usage:
// This struct is primarily used with GORM for querying and manipulating vocabulary data in a PostgreSQL db.
// It is annotated with JSON and GORM tags to map it to the `vocab` table and ensure compatibility with the PostgreSQL backend.
type Vocab struct {
	ID               int       `json:"id" gorm:"primaryKey;autoIncrement"`
	LearningLang     string    `json:"learning_lang" gorm:"not null;unique"`
	FirstLang        string    `json:"first_lang" gorm:"not null"`
	Created          time.Time `json:"created" gorm:"not null;default:now()"`
	Alternatives     string    `json:"alternatives" gorm:"default:''"`
	Skill            string    `json:"skill" gorm:"default:''"`
	Infinitive       string    `json:"infinitive" gorm:"default:''"`
	Pos              string    `json:"pos" gorm:"default:''"`
	Hint             string    `json:"hint" gorm:"default:''"`
	NumLearningWords int       `json:"num_learning_words" gorm:"not null;default:1;check:num_learning_words >= 1"`
	KnownLangCode    string    `json:"known_lang_code" gorm:"default:'en'"`
	LearningLangCode string    `json:"learning_lang_code" gorm:"default:'es'"`
}

// JSON Creates a JSON string from a Vocab object.
func (o *Vocab) JSON() string {
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}
	return string(b)
}

// Clone a Vocab object into a new instance
func (v *Vocab) Clone() *Vocab {
	return &Vocab{
		ID:               v.ID,
		LearningLang:     v.LearningLang,
		FirstLang:        v.FirstLang,
		Created:          v.Created,
		Alternatives:     v.Alternatives,
		Skill:            v.Skill,
		Infinitive:       v.Infinitive,
		Pos:              v.Pos,
		Hint:             v.Hint,
		NumLearningWords: v.NumLearningWords,
		KnownLangCode:    v.KnownLangCode,
		LearningLangCode: v.LearningLangCode,
	}
}

// Compare two Vocab instances for equivalence
func (v *Vocab) Compare(other *Vocab) bool {
	return v.ID == other.ID &&
		v.LearningLang == other.LearningLang &&
		v.FirstLang == other.FirstLang &&
		v.Created.Equal(other.Created) &&
		v.Alternatives == other.Alternatives &&
		v.Skill == other.Skill &&
		v.Infinitive == other.Infinitive &&
		v.Pos == other.Pos &&
		v.Hint == other.Hint &&
		v.NumLearningWords == other.NumLearningWords &&
		v.KnownLangCode == other.KnownLangCode &&
		v.LearningLangCode == other.LearningLangCode
}
