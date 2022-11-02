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

	client.assuranceVisits = append(client.assuranceVisits, sirius.AssuranceVisits{AssuranceType: sirius.AssuranceTypes{
		Handle: "VISIT",
		Label:  "Visit",
	}})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForAssuranceVisits(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(template.lastVars.(AssuranceVisitsVars).ErrorMessage, "You cannot add anything until the current assurance process has a review date and RAG status or is cancelled")
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
		AssuranceType: sirius.AssuranceTypes{
			Handle: "VISIT",
			Label:  "Visit",
		},
	})

	client.assuranceVisits = append(client.assuranceVisits, sirius.AssuranceVisits{
		ReportReviewDate: "01/01/2021",
		AssuranceType: sirius.AssuranceTypes{
			Handle: "VISIT",
			Label:  "Visit",
		},
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForAssuranceVisits(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(template.lastVars.(AssuranceVisitsVars).ErrorMessage, "")
	assert.False(template.lastVars.(AssuranceVisitsVars).AddVisitDisabled)
}

func TestIsCurrentVisitReviewedOrCancelled(t *testing.T) {
	tests := []struct {
		name               string
		visits             []sirius.AssuranceVisits
		want               bool
		wantedErrorMessage string
	}{
		{
			"No visits",
			[]sirius.AssuranceVisits{},
			false,
			"",
		},
		{
			name: "Latest visit is reviewed",
			visits: []sirius.AssuranceVisits{
				{
					ReportReviewDate: "01/01/2022",
					VisitReportMarkedAs: sirius.VisitRagRatingTypes{
						Label:  "RED",
						Handle: "RED",
					},
					AssuranceType: sirius.AssuranceTypes{
						Handle: "VISIT",
						Label:  "Visit",
					},
				},
				{},
			},
			want:               false,
			wantedErrorMessage: "",
		},
		{
			name: "Latest PDR visit is reviewed",
			visits: []sirius.AssuranceVisits{
				{
					ReportReviewDate: "01/01/2022",
					VisitReportMarkedAs: sirius.VisitRagRatingTypes{
						Label:  "RED",
						Handle: "RED",
					},
					AssuranceType: sirius.AssuranceTypes{
						Handle: "PDR",
						Label:  "PDR",
					},
				},
				{},
			},
			want:               false,
			wantedErrorMessage: "",
		},
		{
			"Latest visit has no review date",
			[]sirius.AssuranceVisits{
				{
					VisitReportMarkedAs: sirius.VisitRagRatingTypes{
						Label:  "RED",
						Handle: "RED",
					},
					AssuranceType: sirius.AssuranceTypes{
						Handle: "VISIT",
						Label:  "Visit",
					},
				},
				{
					AssuranceType: sirius.AssuranceTypes{
						Handle: "VISIT",
						Label:  "Visit",
					},
				},
			},
			true,
			"You cannot add anything until the current assurance process has a review date and RAG status or is cancelled",
		},
		{
			"Latest PDR visit has no review date",
			[]sirius.AssuranceVisits{
				{
					VisitReportMarkedAs: sirius.VisitRagRatingTypes{
						Label:  "RED",
						Handle: "RED",
					},
					AssuranceType: sirius.AssuranceTypes{
						Handle: "PDR",
						Label:  "PDR",
					},
				},
				{},
			},
			true,
			"You cannot add anything until the current assurance process has a review date or is marked as 'Not received'",
		},
		{
			"Latest visit has no RAG",
			[]sirius.AssuranceVisits{
				{
					ReportReviewDate: "01/01/2022",
					AssuranceType: sirius.AssuranceTypes{
						Handle: "VISIT",
						Label:  "Visit",
					},
				},
				{},
			},
			true,
			"You cannot add anything until the current assurance process has a review date and RAG status or is cancelled",
		},
		{
			"Latest PDR visit has no RAG",
			[]sirius.AssuranceVisits{
				{
					AssuranceType: sirius.AssuranceTypes{
						Handle: "PDR",
						Label:  "PDR",
					},
				},
				{},
			},
			false,
			"",
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
					AssuranceType: sirius.AssuranceTypes{
						Handle: "VISIT",
						Label:  "Visit",
					},
				},
			},
			true,
			"You cannot add anything until the current assurance process has a review date and RAG status or is cancelled",
		},
		{
			"Latest PDR visit not reviewed but previous one is",
			[]sirius.AssuranceVisits{
				{
					AssuranceType: sirius.AssuranceTypes{
						Handle: "PDR",
						Label:  "PDR",
					},
					VisitReportMarkedAs: sirius.VisitRagRatingTypes{
						Label:  "RED",
						Handle: "RED",
					},
				},
				{
					ReportReviewDate: "01/01/2022",
				},
			},
			true,
			"You cannot add anything until the current assurance process has a review date or is marked as 'Not received'",
		},
		{
			"Latest visit is cancelled",
			[]sirius.AssuranceVisits{
				{
					VisitOutcome: sirius.VisitOutcomeTypes{
						Label:  "Cancelled",
						Handle: "CANCELLED",
					},
					AssuranceType: sirius.AssuranceTypes{
						Handle: "VISIT",
						Label:  "Visit",
					},
				},
				{},
			},
			false,
			"",
		},
		{
			"Latest PDR visit is not received",
			[]sirius.AssuranceVisits{
				{
					PdrOutcome: sirius.PdrOutcomeTypes{
						Label:  "Not received",
						Handle: "NOT_RECEIVED",
					},
					AssuranceType: sirius.AssuranceTypes{
						Handle: "PDR",
						Label:  "PDR",
					},
					VisitReportMarkedAs: sirius.VisitRagRatingTypes{
						Label:  "RED",
						Handle: "RED",
					},
				},
				{},
			},
			false,
			"",
		},
		{
			"Latest visit is not cancelled",
			[]sirius.AssuranceVisits{
				{
					VisitOutcome: sirius.VisitOutcomeTypes{
						Label:  "Successful",
						Handle: "SUCCESSFUL",
					},
					AssuranceType: sirius.AssuranceTypes{
						Handle: "VISIT",
						Label:  "Visit",
					},
				},
				{},
			},
			true,
			"You cannot add anything until the current assurance process has a review date and RAG status or is cancelled",
		},
		{
			"Latest PDR visit is not cancelled",
			[]sirius.AssuranceVisits{
				{
					PdrOutcome: sirius.PdrOutcomeTypes{
						Label:  "Successful",
						Handle: "SUCCESSFUL",
					},
					AssuranceType: sirius.AssuranceTypes{
						Handle: "PDR",
						Label:  "PDR",
					},
					VisitReportMarkedAs: sirius.VisitRagRatingTypes{
						Label:  "RED",
						Handle: "RED",
					},
				},
				{},
			},
			true,
			"You cannot add anything until the current assurance process has a review date or is marked as 'Not received'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErrorMessage := isAddVisitDisabled(tt.visits)
			if got != tt.want {
				t.Errorf("isAddVisitDisabled() = %v, want %v", got, tt.want)
			}
			if gotErrorMessage != tt.wantedErrorMessage {
				t.Errorf("isAddVisitDisabled() = %v, want %v", gotErrorMessage, tt.wantedErrorMessage)
			}
		})
	}
}
