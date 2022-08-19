package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockManageAssuranceVisit struct {
	count                int
	lastCtx              sirius.Context
	assuranceVisits      []sirius.AssuranceVisits
	assuranceVisitsError error
}

func (m *mockManageAssuranceVisit) GetAssuranceVisits(ctx sirius.Context, deputyId int) ([]sirius.AssuranceVisits, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.assuranceVisits, m.assuranceVisitsError
}

func TestGetManageAssuranceVisits_latestNotReviewed(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageAssuranceVisit{}
	template := &mockTemplates{}

	client.assuranceVisits = append(client.assuranceVisits, sirius.AssuranceVisits{})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForAssuranceVisits(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.True(template.lastVars.(AssuranceVisitsVars).AddVisitDisabled)
}

func TestGetManageAssuranceVisits_latestReviewed(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageAssuranceVisit{}
	template := &mockTemplates{}

	client.assuranceVisits = append(client.assuranceVisits, sirius.AssuranceVisits{
		ReportReviewDate: "01/01/2022",
		VisitReportMarkedAs: sirius.VisitRagRatingTypes{
			Label:  "RED",
			Handle: "RED",
		},
	})

	client.assuranceVisits = append(client.assuranceVisits, sirius.AssuranceVisits{
		ReportReviewDate: "01/01/2021",
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForAssuranceVisits(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.False(template.lastVars.(AssuranceVisitsVars).AddVisitDisabled)
}

func TestIsCurrentVisitReviewedOrCancelled(t *testing.T) {
	tests := []struct {
		name   string
		visits []sirius.AssuranceVisits
		want   bool
	}{
		{
			"No visits",
			[]sirius.AssuranceVisits{},
			true,
		},
		{
			"Latest visit is reviewed",
			[]sirius.AssuranceVisits{
				{
					ReportReviewDate: "01/01/2022",
					VisitReportMarkedAs: sirius.VisitRagRatingTypes{
						Label:  "RED",
						Handle: "RED",
					},
				},
				{},
			},
			true,
		},
		{
			"Latest visit has no review date",
			[]sirius.AssuranceVisits{
				{
					VisitReportMarkedAs: sirius.VisitRagRatingTypes{
						Label:  "RED",
						Handle: "RED",
					},
				},
				{},
			},
			false,
		},
		{
			"Latest visit has no RAG",
			[]sirius.AssuranceVisits{
				{
					ReportReviewDate: "01/01/2022",
				},
				{},
			},
			false,
		},
		{
			"Latest visit not reviewed but previous one is",
			[]sirius.AssuranceVisits{
				{},
				{
					ReportReviewDate: "01/01/2022",
					VisitReportMarkedAs: sirius.VisitRagRatingTypes{
						Label:  "RED",
						Handle: "RED",
					},
				},
			},
			false,
		},
		{
			"Latest visit is cancelled",
			[]sirius.AssuranceVisits{
				{
					VisitOutcome: sirius.VisitOutcomeTypes{
						Label:  "Cancelled",
						Handle: "CANCELLED",
					},
				},
				{},
			},
			true,
		},
		{
			"Latest visit is not cancelled",
			[]sirius.AssuranceVisits{
				{
					VisitOutcome: sirius.VisitOutcomeTypes{
						Label:  "Successful",
						Handle: "SUCCESSFUL",
					},
				},
				{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCurrentVisitReviewedOrCancelled(tt.visits); got != tt.want {
				t.Errorf("isCurrentVisitReviewed() = %v, want %v", got, tt.want)
			}
		})
	}
}
