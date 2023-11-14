package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubClientInformation struct {
	count              int
	lastCtx            sirius.Context
	err                error
	deputyClientData   sirius.ClientList
	accommodationTypes []model.RefData
}

func (m *mockDeputyHubClientInformation) GetDeputyClients(ctx sirius.Context, params sirius.ClientListParams) (sirius.ClientList, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyClientData, m.err
}

func (m *mockDeputyHubClientInformation) GetAccommodationTypes(ctx sirius.Context) ([]model.RefData, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.accommodationTypes, m.err
}

func TestNavigateToClientTab(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubClientInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForClientTab(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestParseUrlReturnsColumnAndSortOrder(t *testing.T) {
	urlParams := url.Values{}
	urlParams.Set("sort", "crec:desc")
	expectedResponseColumnBeingSorted, sortOrder, expectedsortBool := "crec", "desc", false
	resultColumnBeingSorted, resultSortOrder, sortBool := parseUrl(urlParams)

	assert.Equal(t, expectedResponseColumnBeingSorted, resultColumnBeingSorted)
	assert.Equal(t, resultSortOrder, sortOrder)
	assert.Equal(t, expectedsortBool, sortBool)

}

func TestParseUrlReturnsEmptyStrings(t *testing.T) {
	urlParams := url.Values{}
	expectedResponseColumnBeingSorted, sortOrder, expectedSortBool := "", "", false
	resultColumnBeingSorted, resultSortOrder, sortBool := parseUrl(urlParams)

	assert.Equal(t, expectedResponseColumnBeingSorted, resultColumnBeingSorted)
	assert.Equal(t, resultSortOrder, sortOrder)
	assert.Equal(t, expectedSortBool, sortBool)

}

func TestListClientsHandlesErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyHubClientInformation{
		err: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/123", strings.NewReader(""))

	returnedError := renderTemplateForClientTab(client, template)(AppVars{}, w, r)

	assert.Equal(client.err, returnedError)

}
