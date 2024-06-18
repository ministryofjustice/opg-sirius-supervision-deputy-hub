package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestRAGRating_GetRiskMarker(t *testing.T) {
	tests := []struct {
		rag  RAGRating
		want VisitReportMarkedAs
	}{
		{
			rag:  RAGRating{},
			want: VisitReportMarkedAs{},
		},
		{
			rag:  RAGRating{Handle: "red"},
			want: VisitReportMarkedAs{Name: "High risk", Colour: "red"},
		},
		{
			rag:  RAGRating{Handle: "amber"},
			want: VisitReportMarkedAs{Name: "Medium risk", Colour: "orange"},
		},
		{
			rag:  RAGRating{Handle: "green"},
			want: VisitReportMarkedAs{Name: "Low risk", Colour: "green"},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.want, test.rag.GetRiskMarker())
		})
	}
}
