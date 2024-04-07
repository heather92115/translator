package convert

import (
	"github.com/heather92115/verdure-admin/graph/model"
	"github.com/heather92115/verdure-admin/internal/mdl"
	"reflect"
	"testing"
)

func TestVocabsToGql(t *testing.T) {
	// Define test cases
	tests := []struct {
		name    string
		from    *[]mdl.Vocab
		want    []*model.Vocab
		wantErr bool
	}{
		{
			name: "Convert non-empty slice of Vocabs",
			from: &[]mdl.Vocab{
				{ID: 1, LearningLang: "Hola", FirstLang: "Hello"},
				{ID: 2, LearningLang: "Bonjour", FirstLang: "Hello"},
			},
			want: []*model.Vocab{
				{ID: "1", LearningLang: "Hola", FirstLang: "Hello"},
				{ID: "2", LearningLang: "Bonjour", FirstLang: "Hello"},
			},
			wantErr: false,
		},
		{
			name:    "Convert nil slice of Vocabs",
			from:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Convert empty slice of Vocabs",
			from:    &[]mdl.Vocab{},
			want:    []*model.Vocab{},
			wantErr: false,
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VocabsToGql(tt.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("VocabsToGql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VocabsToGql() got = %v, want %v", got, tt.want)
			}
		})
	}
}
