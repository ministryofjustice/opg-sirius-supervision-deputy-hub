package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockGcmIssues struct {
	mock.Mock
}

func (m *mockGcmIssues) GetGCMIssues(ctx sirius.Context, deputyId int) ([]sirius.GcmIssue, error) {
	args := m.Called(ctx, deputyId)

	return args.Get(0).([]sirius.GcmIssue), args.Error(1)
}

func TestNavigateToGcmIssuesTab(t *testing.T) {
	assert := assert.New(t)

	client := &mockGcmIssues{}

	deputyDetails := sirius.DeputyDetails{ID: 123}
	app := AppVars{
		DeputyDetails: deputyDetails,
		PageName:      "General Case Manager issues",
	}
	gcmIssues := []sirius.GcmIssue{
		{
			Id:            1,
			Client:        sirius.GcmClient{},
			CreatedDate:   "2024-01-01",
			CreatedByUser: sirius.CreatedByUser{},
			Notes:         "Problem here",
			GcmIssueType: model.RefData{
				Handle:     "MISSING_INFORMATION",
				Label:      "Missing information",
				Deprecated: false,
			},
		},
	}

	client.On("GetGCMIssues", mock.Anything, 123).Return(gcmIssues, nil)

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForGcmIssues(client, template)
	err := handler(app, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(GcmIssuesVars{
		AppVars:   app,
		GcmIssues: gcmIssues,
	}, template.lastVars)
}
