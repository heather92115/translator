package mock

import (
	"errors"
	"github.com/heather92115/translator/internal/mdl"
)

type MockAuditRepository struct {
	audits map[int]*mdl.Audit
	seq    int
}

// NewMockAuditRepository initializes and returns a new instance of MockAuditRepository.
func NewMockAuditRepository() *MockAuditRepository {
	return &MockAuditRepository{
		audits: make(map[int]*mdl.Audit),
	}
}

func (m *MockAuditRepository) FindAuditByID(id int) (*mdl.Audit, error) {
	if audit, exists := m.audits[id]; exists {
		return audit, nil
	}
	return nil, errors.New("audit not found")
}

func (m *MockAuditRepository) FindAudits(tableName string, objectId int, duration *mdl.Duration, limit int) (*[]mdl.Audit, error) {
	result := make([]mdl.Audit, 0)
	count := 0
	for _, a := range m.audits {
		if (tableName == "" || a.TableName == tableName) && (objectId == 0 || a.ObjectID == objectId) &&
			(duration == nil || (a.Created.After(duration.Start) && a.Created.Before(duration.End))) {
			result = append(result, *a)
			count++
			if limit > 0 && count >= limit {
				break
			}
		}
	}
	return &result, nil
}

func (m *MockAuditRepository) CreateAudit(audit *mdl.Audit) error {

	m.seq += 1
	audit.ID = m.seq
	m.audits[audit.ID] = audit
	return nil
}
