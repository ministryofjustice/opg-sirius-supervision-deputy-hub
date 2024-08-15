package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddGCMIssueClient struct {
	count          int
	lastCtx        sirius.Context
	AddGCMIssueErr error
}

func (m *mockAddGCMIssueClient) GetGCMIssueTypes(ctx sirius.Context) ([]model.RefData, error) {
	return []model.RefData{}, nil
}

func (m *mockAddGCMIssueClient) GetDeputyClient(ctx sirius.Context, caseRecNumber string, deputyId int) (sirius.DeputyClient, error) {
	return sirius.DeputyClient{}, nil
}

func (m *mockAddGCMIssueClient) AddGcmIssue(ctx sirius.Context, caseRecNumber, notes string, gcmIssueType model.RefData, deputyId int) error {
	return nil
}

var addGCMIssueAppVars = AppVars{
	DeputyDetails: testDeputy,
}

func TestGetAddGCMIssue(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddGCMIssueClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddGcmIssue(client, template)
	err := handler(addGCMIssueAppVars, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}
