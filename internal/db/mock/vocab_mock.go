package mock

import (
	"errors"
	"fmt"
	"github.com/heather92115/translator/internal/mdl"
)

type MockVocabRepository struct {
	vocabs map[int]*mdl.Vocab
}

// NewMockVocabRepository initializes and returns a new instance of MockVocabRepository.
func NewMockVocabRepository() *MockVocabRepository {
	return &MockVocabRepository{
		vocabs: make(map[int]*mdl.Vocab),
	}
}

func (m *MockVocabRepository) FindVocabByID(id int) (*mdl.Vocab, error) {
	if vocab, exists := m.vocabs[id]; exists {
		return vocab, nil
	}
	return nil, fmt.Errorf("error finding vocab with id %d", id)
}

func (m *MockVocabRepository) FindVocabByLearningLang(learningLang string) (vocab *mdl.Vocab, err error) {
	for _, v := range m.vocabs {
		if v.LearningLang == learningLang {
			return v, nil
		}
	}
	return nil, fmt.Errorf("error finding vocab with learning lang %s", learningLang)
}

func (m *MockVocabRepository) FindVocabs(learningCode string, hasFirst bool, limit int) (*[]mdl.Vocab, error) {
	result := make([]mdl.Vocab, 0)
	count := 0
	for _, v := range m.vocabs {
		if v.LearningLangCode == learningCode && (!hasFirst && v.FirstLang == "" || hasFirst && v.FirstLang != "") {
			result = append(result, *v)
			count++
			if limit > 0 && count >= limit {
				break
			}
		}
	}
	return &result, nil
}

func (m *MockVocabRepository) CreateVocab(vocab *mdl.Vocab) error {
	if _, exists := m.vocabs[vocab.ID]; exists {
		return errors.New("vocab already exists")
	}
	m.vocabs[vocab.ID] = vocab
	return nil
}

func (m *MockVocabRepository) UpdateVocab(vocab *mdl.Vocab) error {
	if _, exists := m.vocabs[vocab.ID]; !exists {
		return fmt.Errorf("error finding vocab with id %d", vocab.ID)
	}
	m.vocabs[vocab.ID] = vocab
	return nil
}
