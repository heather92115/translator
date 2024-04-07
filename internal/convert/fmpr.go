package convert

import (
	"fmt"
	"github.com/heather92115/verdure-admin/graph/model"
	"github.com/heather92115/verdure-admin/internal/mdl"
	"strconv"
)

// FixitToGql maps a mdl.Fixit struct to a graph model's Fixit struct.
func FixitToGql(from *mdl.Fixit) (*model.Fixit, error) {
	if from == nil {
		return nil, fmt.Errorf("expected a fixit record but found nothing")
	}

	status, err := fixitStatusToGql(from.Status)
	if err != nil {
		return nil, err
	}

	return &model.Fixit{
		ID:        strconv.Itoa(from.ID),
		VocabID:   strconv.Itoa(from.VocabID),
		Status:    status,
		FieldName: from.FieldName,
		Comments:  from.Comments,
		CreatedBy: from.CreatedBy,
		Created:   timeToGQLDateTime(from.Created),
	}, nil
}

// FixitsQueryMapper converts GraphQL query parameters into their corresponding internal representations.
// It takes a fixit status as a string, a vocabID as a string, and start and end times as ISO 8601 formatted strings.
// It returns the internal status type, vocabID as an integer, a duration struct representing the time range, and an error if any conversions fail.
// This function is primarily used to prepare parameters for querying the database with filters provided by a GraphQL request.
//
// Parameters:
//   - status: Fixit status in GraphQL enum format (e.g., "PENDING").
//   - vocabID: Vocab ID as a string, which should be convertible to an integer.
//   - startTime: Start of the duration as an ISO 8601 formatted string.
//   - endTime: End of the duration as an ISO 8601 formatted string.
//
// Returns:
//   - StatusType: The internal representation of the fixit status.
//   - int: The vocabID converted to an integer.
//   - *Duration: A pointer to a Duration struct representing the time range.
//   - error: An error if any of the conversions fail (invalid status, non-integer vocabID, or invalid date formats).
func FixitsQueryMapper(status model.Status, vocabID string, startTime string, endTime string) (mdl.StatusType, int, *mdl.Duration, error) {
	fStatus, err := FixitStatusFromGql(status)
	if err != nil {
		return "", 0, nil, err
	}

	fVocabID, err := strconv.Atoi(vocabID)
	if err != nil {
		return "", 0, nil, err
	}

	duration, err := GqlDateTimeToDuration(startTime, endTime)
	if err != nil {
		return "", 0, nil, err
	}

	return fStatus, fVocabID, duration, nil
}

// FixitsToGql maps a slice of mdl.Fixit structs to a slice of graph model's Fixit struct.
// This converts internal fixit structs to the graphql schema form.
func FixitsToGql(from *[]mdl.Fixit) ([]*model.Fixit, error) {
	if from == nil {
		return nil, fmt.Errorf("expected a list of fixit records but found nothing")
	}

	result := make([]*model.Fixit, len(*from))
	for i, v := range *from {
		gqlFixit, err := FixitToGql(&v)
		if err != nil {
			return nil, err // Propagate errors from FixitToGql
		}
		result[i] = gqlFixit
	}

	return result, nil
}

// NewFixitFromGql maps a model.NewFixit struct to a mdl.Fixit struct.
func NewFixitFromGql(from *model.NewFixit) (*mdl.Fixit, error) {
	if from == nil {
		return nil, fmt.Errorf("expected a NewFixit from gql, but found nothing")
	}

	vocabID, err := strconv.Atoi(from.VocabID)
	if err != nil {
		return nil, fmt.Errorf("invalid VocabID %v", from.VocabID)
	}

	status, err := FixitStatusFromGql(from.Status)
	if err != nil {
		return nil, err
	}

	return &mdl.Fixit{
		VocabID:   vocabID,
		Status:    status,
		FieldName: from.FieldName,
		Comments:  from.Comments,
	}, nil
}

// UpdateFixitFromGql maps a model.UpdateFixit struct to a mdl.Fixit struct.
func UpdateFixitFromGql(from *model.UpdateFixit) (*mdl.Fixit, error) {
	if from == nil {
		return nil, fmt.Errorf("expected an UpdateFixit from gql, but found nothing")
	}

	id, err := strconv.Atoi(from.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID %v", from.ID)
	}

	status, err := FixitStatusFromGql(from.Status)
	if err != nil {
		return nil, err
	}

	return &mdl.Fixit{
		ID:        id,
		Status:    status,
		FieldName: from.FieldName,
		Comments:  from.Comments,
	}, nil
}

// FixitStatusFromGql converts the status enum from GraphQL to internal model
func FixitStatusFromGql(gqlStatus model.Status) (mdl.StatusType, error) {
	switch gqlStatus {
	case "PENDING":
		return mdl.Pending, nil
	case "IN_PROGRESS":
		return mdl.InProgress, nil
	case "COMPLETED":
		return mdl.Completed, nil
	default:
		return "", fmt.Errorf("invalid status: %s", gqlStatus)
	}
}

// Convert the status enum from internal model to GraphQL
func fixitStatusToGql(status mdl.StatusType) (model.Status, error) {
	switch status {
	case mdl.Pending:
		return "PENDING", nil
	case mdl.InProgress:
		return "IN_PROGRESS", nil
	case mdl.Completed:
		return "COMPLETED", nil
	default:
		return "", fmt.Errorf("unknown status: %s", status)
	}
}
