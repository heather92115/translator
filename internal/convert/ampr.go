package convert

import (
	"fmt"
	"github.com/heather92115/translator/graph/model"
	"github.com/heather92115/translator/internal/mdl"
	"strconv"
)

// AuditToGql maps a mdl.Audit struct to a corresponding GraphQL model.Audit struct.
func AuditToGql(from *mdl.Audit) (*model.Audit, error) {
	if from == nil {
		return nil, fmt.Errorf("expected an audit record but found nothing")
	}

	return &model.Audit{
		ID:        strconv.Itoa(from.ID),
		ObjectID:  strconv.Itoa(from.ObjectID),
		TableName: from.TableName,
		Diff:      from.Diff,
		Before:    from.Before,
		After:     from.After,
		Comments:  from.Comments,
		CreatedBy: from.CreatedBy,
		Created:   timeToGQLDateTime(from.Created),
	}, nil
}

// AuditsToGql maps a slice of mdl.Audit structs to a slice of GraphQL model.Audit structs.
// This converter facilitates converting a collection of internal audit structs to the GraphQL schema form.
func AuditsToGql(from *[]mdl.Audit) ([]*model.Audit, error) {
	if from == nil {
		return nil, fmt.Errorf("expected a list of audit records but found nothing")
	}

	result := make([]*model.Audit, len(*from))
	for i, audit := range *from {
		gqlAudit, err := AuditToGql(&audit)
		if err != nil {
			return nil, err // Propagate errors from AuditToGql
		}
		result[i] = gqlAudit
	}

	return result, nil
}

// AuditQueryMapper converts GraphQL query parameters for an audit search into internal representations.
// It parses the objectID from a string to an integer and converts startTime and endTime from
// GraphQL DateTime strings to a mdl.Duration struct representing the time range of interest.
//
// Parameters:
// - objectID: The ID of the object associated with the audits as a string.
// - startTime: The start of the search period as a GraphQL DateTime string (ISO 8601 format).
// - endTime: The end of the search period as a GraphQL DateTime string (ISO 8601 format).
//
// Returns:
// - The object ID as an integer.
// - A pointer to a mdl.Duration struct representing the time range.
// - An error if the conversion of objectID or the parsing of start and end times fails.
func AuditQueryMapper(objectID string, startTime string, endTime string) (int, *mdl.Duration, error) {

	aObjectID, err := strconv.Atoi(objectID)
	if err != nil {
		return 0, nil, err
	}

	duration, err := GqlDateTimeToDuration(startTime, endTime)
	if err != nil {
		return 0, nil, err
	}

	return aObjectID, duration, nil
}
