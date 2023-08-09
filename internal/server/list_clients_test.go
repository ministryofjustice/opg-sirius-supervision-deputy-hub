package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func (m *mockDeputyHubClientInformation) GetDeputyClients(ctx sirius.Context, deputyId, displayClientLimit, search int, deputyType, columnBeingSorted, sortOrder string) (sirius.ClientList, sirius.AriaSorting, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyClientData, m.ariaSorting, m.err
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
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestParseUrlReturnsColumnAndSortOrder(t *testing.T) {
	urlParams := url.Values{}
	urlParams.Set("sort", "crec:desc")
	expectedResponseColumnBeingSorted, sortOrder := "crec", "desc"
	resultColumnBeingSorted, resultSortOrder := parseUrl(urlParams)

	assert.Equal(t, expectedResponseColumnBeingSorted, resultColumnBeingSorted)
	assert.Equal(t, resultSortOrder, sortOrder)
}

func TestParseUrlReturnsEmptyStrings(t *testing.T) {
	urlParams := url.Values{}
	expectedResponseColumnBeingSorted, sortOrder := "", ""
	resultColumnBeingSorted, resultSortOrder := parseUrl(urlParams)

	assert.Equal(t, expectedResponseColumnBeingSorted, resultColumnBeingSorted)
	assert.Equal(t, resultSortOrder, sortOrder)
}

func TestListClientsHandlesErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyHubClientInformation{
		err: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	returnedError := renderTemplateForClientTab(client, template)(sirius.DeputyDetails{}, w, r)

	assert.Equal(client.err, returnedError)

}
