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

func TestGetManageAssuranceVisits(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageAssuranceVisit{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForAssuranceVisits(client, template)
	err := handler(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
