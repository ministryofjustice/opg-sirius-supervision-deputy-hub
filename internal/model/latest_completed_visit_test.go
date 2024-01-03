package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestAssurance_GetRAGRating(t *testing.T) {
	tests := []struct {
		visit LatestCompletedVisit
		want  RAGRating
	}{
		{
			visit: LatestCompletedVisit{},
			want:  RAGRating{},
		},
		{
			visit: LatestCompletedVisit{RagRatingLowerCase: "red"},
			want:  RAGRating{Name: "High risk", Colour: "red"},
		},
		{
			visit: LatestCompletedVisit{RagRatingLowerCase: "amber"},
			want:  RAGRating{Name: "Medium risk", Colour: "orange"},
		},
		{
			visit: LatestCompletedVisit{RagRatingLowerCase: "green"},
			want:  RAGRating{Name: "Low risk", Colour: "green"},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.want, test.visit.GetRAGRating())
		})
	}
}
