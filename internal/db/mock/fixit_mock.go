package mock

import (
	"errors"
	"github.com/heather92115/verdure-admin/internal/mdl"
)

type MockFixitRepository struct {
	fixits map[int]*mdl.Fixit
	seq    int
}

// NewMockFixitRepository initializes and returns a new instance of MockFixitRepository.
func NewMockFixitRepository() *MockFixitRepository {
	return &MockFixitRepository{
		fixits: make(map[int]*mdl.Fixit),
	}
}

func (m *MockFixitRepository) FindFixitByID(id int) (*mdl.Fixit, error) {
	if fixit, exists := m.fixits[id]; exists {
		return fixit, nil
	}
	return nil, errors.New("fixit not found")
}

func (m *MockFixitRepository) FindFixits(status mdl.StatusType, vocabID int, duration *mdl.Duration, limit int) (*[]mdl.Fixit, error) {
	result := make([]mdl.Fixit, 0)
	count := 0
	for _, f := range m.fixits {
		if (status == "" || f.Status == status) &&
			(vocabID == 0 || f.VocabID == vocabID) &&
			(duration == nil || (f.Created.After(duration.Start) && f.Created.Before(duration.End))) {
			result = append(result, *f)
			count++
			if limit > 0 && count >= limit {
				break
			}
		}
	}
	return &result, nil
}

func (m *MockFixitRepository) CreateFixit(fixit *mdl.Fixit) error {
	m.seq += 1
	fixit.ID = m.seq
	m.fixits[fixit.ID] = fixit
	return nil
}

func (m *MockFixitRepository) UpdateFixit(fixit *mdl.Fixit) error {
	if _, exists := m.fixits[fixit.ID]; !exists {
		return errors.New("fixit does not exist")
	}
	m.fixits[fixit.ID] = fixit
	return nil
}
