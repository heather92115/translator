// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Audit struct {
	ID        string `json:"id"`
	ObjectID  string `json:"object_id"`
	TableName string `json:"table_name"`
	Diff      string `json:"diff"`
	Before    string `json:"before"`
	After     string `json:"after"`
	Comments  string `json:"comments"`
	CreatedBy string `json:"created_by"`
	Created   string `json:"created"`
}

type Fixit struct {
	ID        string `json:"id"`
	VocabID   string `json:"vocab_id"`
	Status    Status `json:"status"`
	FieldName string `json:"field_name"`
	Comments  string `json:"comments"`
	CreatedBy string `json:"created_by"`
	Created   string `json:"created"`
}

type Mutation struct {
}

type NewFixit struct {
	VocabID   string `json:"vocab_id"`
	Status    Status `json:"status"`
	FieldName string `json:"field_name"`
	Comments  string `json:"comments"`
}

type NewVocab struct {
	LearningLang     string `json:"learning_lang"`
	FirstLang        string `json:"first_lang"`
	Alternatives     string `json:"alternatives"`
	Skill            string `json:"skill"`
	Infinitive       string `json:"infinitive"`
	Pos              string `json:"pos"`
	Hint             string `json:"hint"`
	NumLearningWords int    `json:"num_learning_words"`
	KnownLangCode    string `json:"known_lang_code"`
	LearningLangCode string `json:"learning_lang_code"`
}

type Query struct {
}

type UpdateFixit struct {
	ID        string `json:"id"`
	Status    Status `json:"status"`
	FieldName string `json:"field_name"`
	Comments  string `json:"comments"`
}

type UpdateVocab struct {
	ID               string `json:"id"`
	FirstLang        string `json:"first_lang"`
	Alternatives     string `json:"alternatives"`
	Skill            string `json:"skill"`
	Infinitive       string `json:"infinitive"`
	Pos              string `json:"pos"`
	Hint             string `json:"hint"`
	NumLearningWords int    `json:"num_learning_words"`
}

type Vocab struct {
	ID               string `json:"id"`
	LearningLang     string `json:"learning_lang"`
	FirstLang        string `json:"first_lang"`
	Alternatives     string `json:"alternatives"`
	Skill            string `json:"skill"`
	Infinitive       string `json:"infinitive"`
	Pos              string `json:"pos"`
	Hint             string `json:"hint"`
	NumLearningWords int    `json:"num_learning_words"`
	KnownLangCode    string `json:"known_lang_code"`
	LearningLangCode string `json:"learning_lang_code"`
}

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
)

var AllStatus = []Status{
	StatusPending,
	StatusInProgress,
	StatusCompleted,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusPending, StatusInProgress, StatusCompleted:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
