package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockGetAssurancesClient struct {
	count      int
	lastCtx    sirius.Context
	assurances []model.Assurance
	err        error
}

func (m *mockGetAssurancesClient) GetAssurances(ctx sirius.Context, deputyId int) ([]model.Assurance, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.assurances, m.err
}

func TestGetAssurances_LatestNotReviewed(t *testing.T) {
	assert := assert.New(t)

	client := &mockGetAssurancesClient{}
	template := &mockTemplates{}

	client.assurances = append(client.assurances, model.Assurance{Type: model.AssuranceType{
		Handle: "VISIT",
		Label:  "Assurance",
	}})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForAssurances(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(template.lastVars.(AssurancesVars).ErrorMessage, "You cannot add anything until the current assurance process has a review date and RAG status or is cancelled")
	assert.True(template.lastVars.(AssurancesVars).AddVisitDisabled)
}

func TestGetAssurances_LatestReviewed(t *testing.T) {
	assert := assert.New(t)

	client := &mockGetAssurancesClient{
		assurances: []model.Assurance{
			{
				ReportReviewDate: "01/01/2022",
				ReportMarkedAs: model.RAGRating{
					Label:  "RED",
					Handle: "RED",
				},
				Type: model.AssuranceType{
					Handle: "VISIT",
					Label:  "Assurance",
				},
			},
			{
				ReportReviewDate: "01/01/2021",
				Type: model.AssuranceType{
					Handle: "VISIT",
					Label:  "Assurance",
				},
			},
		},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForAssurances(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(template.lastVars.(AssurancesVars).ErrorMessage, "")
	assert.False(template.lastVars.(AssurancesVars).AddVisitDisabled)
}

func TestIsCurrentVisitReviewedOrCancelled(t *testing.T) {
	tests := []struct {
		name               string
		assurances         []model.Assurance
		want               bool
		wantedErrorMessage string
	}{
		{
			"No assurances",
			[]model.Assurance{},
			false,
			"",
		},
		{
			name: "Latest visit is reviewed",
			assurances: []model.Assurance{
				{
					ReportReviewDate: "01/01/2022",
					ReportMarkedAs: model.RAGRating{
						Label:  "RED",
						Handle: "RED",
					},
					Type: model.AssuranceType{
						Handle: "VISIT",
						Label:  "Assurance",
					},
				},
				{},
			},
			want:               false,
			wantedErrorMessage: "",
		},
		{
			name: "Latest PDR visit is reviewed",
			assurances: []model.Assurance{
				{
					ReportReviewDate: "01/01/2022",
					Type: model.AssuranceType{
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
			[]model.Assurance{
				{
					ReportMarkedAs: model.RAGRating{
						Label:  "RED",
						Handle: "RED",
					},
					Type: model.AssuranceType{
						Handle: "VISIT",
						Label:  "Assurance",
					},
				},
				{
					Type: model.AssuranceType{
						Handle: "VISIT",
						Label:  "Assurance",
					},
				},
			},
			true,
			"You cannot add anything until the current assurance process has a review date and RAG status or is cancelled",
		},
		{
			"Latest PDR visit has no review date",
			[]model.Assurance{
				{
					Type: model.AssuranceType{
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
			[]model.Assurance{
				{
					ReportReviewDate: "01/01/2022",
					Type: model.AssuranceType{
						Handle: "VISIT",
						Label:  "Assurance",
					},
				},
				{},
			},
			true,
			"You cannot add anything until the current assurance process has a review date and RAG status or is cancelled",
		},
		{
			"Latest visit not reviewed but previous one is",
			[]model.Assurance{
				{},
				{
					ReportReviewDate: "01/01/2022",
					ReportMarkedAs: model.RAGRating{
						Label:  "RED",
						Handle: "RED",
					},
					Type: model.AssuranceType{
						Handle: "VISIT",
						Label:  "Assurance",
					},
				},
			},
			true,
			"You cannot add anything until the current assurance process has a review date and RAG status or is cancelled",
		},
		{
			"Latest PDR visit not reviewed but previous one is",
			[]model.Assurance{
				{
					Type: model.AssuranceType{
						Handle: "PDR",
						Label:  "PDR",
					},
					ReportMarkedAs: model.RAGRating{
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
			[]model.Assurance{
				{
					VisitOutcome: model.VisitOutcomeType{
						Label:  "Cancelled",
						Handle: "CANCELLED",
					},
					Type: model.AssuranceType{
						Handle: "VISIT",
						Label:  "Assurance",
					},
				},
				{},
			},
			false,
			"",
		},
		{
			"Latest PDR visit is not received",
			[]model.Assurance{
				{
					PdrOutcome: model.PdrOutcomeType{
						Label:  "Not received",
						Handle: "NOT_RECEIVED",
					},
					Type: model.AssuranceType{
						Handle: "PDR",
						Label:  "PDR",
					},
					ReportMarkedAs: model.RAGRating{
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
			[]model.Assurance{
				{
					VisitOutcome: model.VisitOutcomeType{
						Label:  "Successful",
						Handle: "SUCCESSFUL",
					},
					Type: model.AssuranceType{
						Handle: "VISIT",
						Label:  "Assurance",
					},
				},
				{},
			},
			true,
			"You cannot add anything until the current assurance process has a review date and RAG status or is cancelled",
		},
		{
			"Latest PDR visit is not cancelled",
			[]model.Assurance{
				{
					PdrOutcome: model.PdrOutcomeType{
						Label:  "Successful",
						Handle: "SUCCESSFUL",
					},
					Type: model.AssuranceType{
						Handle: "PDR",
						Label:  "PDR",
					},
					ReportMarkedAs: model.RAGRating{
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
			got, gotErrorMessage := isAddVisitDisabled(tt.assurances)
			if got != tt.want {
				t.Errorf("isAddVisitDisabled() = %v, want %v", got, tt.want)
			}
			if gotErrorMessage != tt.wantedErrorMessage {
				t.Errorf("isAddVisitDisabled() = %v, want %v", gotErrorMessage, tt.wantedErrorMessage)
			}
		})
	}
}

func TestGetAssurancesReturnsNonValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockGetAssurancesClient{
		err: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))

	returnedError := renderTemplateForAssurances(client, template)(AppVars{}, w, r)

	assert.Equal(client.err, returnedError)
}
