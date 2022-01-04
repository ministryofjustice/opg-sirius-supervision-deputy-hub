package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubClientInformation struct {
	count            int
	lastCtx          sirius.Context
	err              error
	deputyData       sirius.DeputyDetails
	deputyClientData sirius.DeputyClientDetails
	ariaSorting      sirius.AriaSorting
}

func (m *mockDeputyHubClientInformation) GetDeputyDetails(ctx sirius.Context, defaultPATeam int, deputyId int) (sirius.DeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func (m *mockDeputyHubClientInformation) GetDeputyClients(ctx sirius.Context, deputyId int, deputyType string, columnBeingSorted string, sortOrder string) (sirius.DeputyClientDetails, sirius.AriaSorting, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyClientData, m.ariaSorting, m.err
}

func TestNavigateToClientTab(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubClientInformation{}
	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForClientTab(client, defaultPATeam, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestParseUrlReturnsColumnAndSortOrder(t *testing.T) {
	urlPassedin := "http://localhost:8888/supervision/deputies/public-authority/deputy/78/clients?sort=crec:desc"
	expectedResponseColumnBeingSorted, sortOrder := "sort=crec", "desc"
	resultColumnBeingSorted, resultSortOrder := parseUrl(urlPassedin)

	assert.Equal(t, expectedResponseColumnBeingSorted, resultColumnBeingSorted)
	assert.Equal(t, resultSortOrder, sortOrder)

}

func TestParseUrlReturnsEmptyStrings(t *testing.T) {
	urlPassedin := "http://localhost:8888/supervision/deputies/public-authority/deputy/78/clients"
	expectedResponseColumnBeingSorted, sortOrder := "", ""
	resultColumnBeingSorted, resultSortOrder := parseUrl(urlPassedin)

	assert.Equal(t, expectedResponseColumnBeingSorted, resultColumnBeingSorted)
	assert.Equal(t, resultSortOrder, sortOrder)

}
