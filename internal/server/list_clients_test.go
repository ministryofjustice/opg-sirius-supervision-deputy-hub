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
	supervisionLevels  []model.RefData
	SuccessMessage     string
}

func (m *mockDeputyHubClientInformation) AssignAssuranceVisitToClients(ctx sirius.Context, params sirius.AssignAssuranceVisitToClientsParams, deputyId int) (string, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.SuccessMessage, m.err
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

func (m *mockDeputyHubClientInformation) GetSupervisionLevels(ctx sirius.Context) ([]model.RefData, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.supervisionLevels, m.err
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

func TestGetFiltersFromParamsWithOrderStatus(t *testing.T) {
	params := url.Values{}
	params.Set("order-status", "ACTIVE")

	var expectedResponseForAccommodation, expectedResponseForSupervisionLevel []string
	expectedResponseForOrderStatus := []string{"ACTIVE"}
	resultOrderStatus, resultAccommodation, resultSupervisionLevel := getFiltersFromParams(params)

	assert.Equal(t, resultOrderStatus, expectedResponseForOrderStatus)
	assert.Equal(t, expectedResponseForAccommodation, resultAccommodation)
	assert.Equal(t, expectedResponseForSupervisionLevel, resultSupervisionLevel)
}

func TestGetFiltersFromParamsWithAccommodation(t *testing.T) {
	params := url.Values{}
	params.Set("accommodation", "COUNCIL RENTED")

	expectedResponseForAccommodation := []string{"COUNCIL RENTED"}
	var expectedResponseForOrderStatus, expectedResponseForSupervisionLevel []string
	resultOrderStatus, resultAccommodation, resultSupervisionLevel := getFiltersFromParams(params)

	assert.Equal(t, resultOrderStatus, expectedResponseForOrderStatus)
	assert.Equal(t, expectedResponseForAccommodation, resultAccommodation)
	assert.Equal(t, expectedResponseForSupervisionLevel, resultSupervisionLevel)
}

func TestGetFiltersFromParamsWithSupervisionLevel(t *testing.T) {
	params := url.Values{}
	params.Set("supervision-level", "GENERAL")

	expectedResponseForSupervisionLevel := []string{"GENERAL"}
	var expectedResponseForOrderStatus, expectedResponseForAccommodation []string
	resultOrderStatus, resultAccommodation, resultSupervisionLevel := getFiltersFromParams(params)

	assert.Equal(t, resultOrderStatus, expectedResponseForOrderStatus)
	assert.Equal(t, expectedResponseForAccommodation, resultAccommodation)
	assert.Equal(t, expectedResponseForSupervisionLevel, resultSupervisionLevel)
}

func TestGetFiltersFromParamsWithAllFilters(t *testing.T) {
	params := url.Values{}
	params.Set("order-status", "ACTIVE")
	params.Set("accommodation", "COUNCIL RENTED")
	params.Set("supervision-level", "GENERAL")

	expectedResponseForAccommodation := []string{"COUNCIL RENTED"}
	expectedResponseForOrderStatus := []string{"ACTIVE"}
	expectedResponseForSupervisionLevel := []string{"GENERAL"}
	resultOrderStatus, resultAccommodation, resultSupervisionLevel := getFiltersFromParams(params)

	assert.Equal(t, resultOrderStatus, expectedResponseForOrderStatus)
	assert.Equal(t, expectedResponseForAccommodation, resultAccommodation)
	assert.Equal(t, expectedResponseForSupervisionLevel, resultSupervisionLevel)
}
