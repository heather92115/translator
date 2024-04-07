package convert

import (
	"fmt"
	"github.com/heather92115/verdure-admin/graph/model"
	"github.com/heather92115/verdure-admin/internal/mdl"
	"strconv"
)

// VocabToGql maps a mdl.Vocab struct to a model.Vocab struct.
func VocabToGql(from *mdl.Vocab) (*model.Vocab, error) {
	if from == nil {
		return nil, fmt.Errorf("expected a vocab record but found nothing")
	}

	return &model.Vocab{
		ID:               strconv.Itoa(from.ID), // Convert int ID to string
		LearningLang:     from.LearningLang,
		FirstLang:        from.FirstLang,
		Alternatives:     from.Alternatives,
		Skill:            from.Skill,
		Infinitive:       from.Infinitive,
		Pos:              from.Pos,
		Hint:             from.Hint,
		NumLearningWords: from.NumLearningWords,
		KnownLangCode:    from.KnownLangCode,
		LearningLangCode: from.LearningLangCode,
	}, nil
}

// VocabsToGql maps a slice of mdl.Vocab structs to a slice of model.Vocabs struct.
// This converts internal vocab structs to the graphql schema form.
func VocabsToGql(from *[]mdl.Vocab) ([]*model.Vocab, error) {
	if from == nil {
		return nil, fmt.Errorf("expected a list of vocab records but found nothing")
	}

	result := make([]*model.Vocab, len(*from))
	for i, v := range *from {
		gqlVocab, err := VocabToGql(&v)
		if err != nil {
			return nil, err // Propagate errors from VocabToGql
		}
		result[i] = gqlVocab
	}

	return result, nil
}

// VocabFromGql maps a model.Vocab struct to a mdl.Vocab struct.
func VocabFromGql(from *model.UpdateVocab) (*mdl.Vocab, error) {
	if from == nil {
		return nil, fmt.Errorf("expected an vocab from gql, but found nothing")
	}

	// Convert string ID to int. Handle potential conversion error by defaulting to 0 or logging.
	id, err := strconv.Atoi(from.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id %v", from.ID)
	}

	return &mdl.Vocab{
		ID:               id,
		FirstLang:        from.FirstLang,
		Alternatives:     from.Alternatives,
		Skill:            from.Skill,
		Infinitive:       from.Infinitive,
		Pos:              from.Pos,
		Hint:             from.Hint,
		NumLearningWords: from.NumLearningWords,
	}, nil
}

// VocabFromNewGql maps a model.NewVocab struct to a mdl.Vocab struct.
func VocabFromNewGql(from *model.NewVocab) (*mdl.Vocab, error) {
	if from == nil {
		return nil, fmt.Errorf("expected an vocab from gql, but found nothing")
	}

	return &mdl.Vocab{
		LearningLang:     from.LearningLang,
		FirstLang:        from.FirstLang,
		Alternatives:     from.Alternatives,
		Skill:            from.Skill,
		Infinitive:       from.Infinitive,
		Pos:              from.Pos,
		Hint:             from.Hint,
		NumLearningWords: from.NumLearningWords,
		KnownLangCode:    from.KnownLangCode,
		LearningLangCode: from.LearningLangCode,
	}, nil
}
