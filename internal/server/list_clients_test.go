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
	deputyClientData sirius.ClientList
	pageDetails      sirius.PageDetails
	ariaSorting      sirius.AriaSorting
}

func (m *mockDeputyHubClientInformation) GetDeputyClients(ctx sirius.Context, deputyId, displayClientLimit, search int, deputyType, columnBeingSorted, sortOrder string) (sirius.ClientList, sirius.AriaSorting, int, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyClientData, m.ariaSorting, 0, m.err
}

func (m *mockDeputyHubClientInformation) GetPageDetails(sirius.Context, sirius.ClientList, int, int) sirius.PageDetails {
	m.count += 1

	return m.pageDetails
}

func TestNavigateToClientTab(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubClientInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForClientTab(client, template)
	err := handler(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestParseUrlReturnsColumnAndSortOrder(t *testing.T) {
	urlPassedin := "http://localhost:8888/supervision/deputies/78/clients?sort=crec:desc"
	expectedResponseColumnBeingSorted, sortOrder := "sort=crec", "desc"
	resultColumnBeingSorted, resultSortOrder := parseUrl(urlPassedin)

	assert.Equal(t, expectedResponseColumnBeingSorted, resultColumnBeingSorted)
	assert.Equal(t, resultSortOrder, sortOrder)

}

func TestParseUrlReturnsEmptyStrings(t *testing.T) {
	urlPassedin := "http://localhost:8888/supervision/deputies/78/clients"
	expectedResponseColumnBeingSorted, sortOrder := "", ""
	resultColumnBeingSorted, resultSortOrder := parseUrl(urlPassedin)

	assert.Equal(t, expectedResponseColumnBeingSorted, resultColumnBeingSorted)
	assert.Equal(t, resultSortOrder, sortOrder)

}
