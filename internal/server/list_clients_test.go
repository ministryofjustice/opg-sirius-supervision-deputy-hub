package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
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

func (m *mockDeputyHubClientInformation) BulkAssignAssuranceVisitTasksToClients(ctx sirius.Context, params sirius.BulkAssignAssuranceVisitTasksToClientsParams, deputyId int) (string, error) {
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

func TestPostBulkAssuranceVisitTasksReturnsErrorWithDueDate(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyHubClientInformation{}

	form := url.Values{}
	form.Add("due-date", "")

	template := &mockTemplates{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/1/bulk-assurance-visit-tasks", strings.NewReader(form.Encode()))

	returnedError := renderTemplateForClientTab(client, template)(AppVars{}, w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"due-date": {"": "Enter a due date"},
	}

	assert.Equal(ListClientsVars{
		ListPage: ListPage{
			AppVars: AppVars{
				PageName: "Clients",
				Errors:   expectedValidationErrors,
			},
			Sort: urlbuilder.Sort{
				OrderBy: "surname",
			},
			Pagination: paginate.Pagination{
				ElementsPerPage: 25,
				ElementName:     "clients",
				PerPageOptions:  []int{25, 50, 100},
				UrlBuilder: urlbuilder.UrlBuilder{
					OriginalPath:    "clients",
					SelectedPerPage: 25,
					SelectedFilters: []urlbuilder.Filter{
						{Name: "order-status"},
						{Name: "accommodation"},
						{Name: "supervision-level"},
					},
					SelectedSort: urlbuilder.Sort{
						OrderBy: "surname",
					},
				},
			},
			PerPage: 25,
			UrlBuilder: urlbuilder.UrlBuilder{
				OriginalPath:    "clients",
				SelectedPerPage: 25,
				SelectedFilters: []urlbuilder.Filter{
					{Name: "order-status"},
					{Name: "accommodation"},
					{Name: "supervision-level"},
				},
				SelectedSort: urlbuilder.Sort{
					OrderBy: "surname",
				},
			},
		},
		FilterByOrderStatus: FilterByOrderStatus{
			OrderStatusOptions: []model.RefData{
				{
					Handle: "ACTIVE",
					Label:  "Active",
				},
				{
					Handle: "CLOSED",
					Label:  "Closed",
				},
			},
			OrderStatuses: []model.OrderStatus{
				{
					Handle:     "ACTIVE",
					Incomplete: "Active",
					Category:   "Active",
					Complete:   "Active",
				},
				{
					Handle:     "CLOSED",
					Incomplete: "Closed",
					Category:   "Closed",
					Complete:   "Closed",
				},
			},
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestPostBulkAssuranceVisitTask(t *testing.T) {
	assert := assert.New(t)
	template := &mockTemplates{}

	client := &mockDeputyHubClientInformation{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/bulk-assurance-visit-tasks", strings.NewReader("{dueDate:2626-02-02, clientIds: [1,2]}"))

	form := url.Values{}
	form.Add("due-date", "2626-02-02")
	form.Add("selected-clients", "1")
	form.Add("selected-clients", "2")
	r.PostForm = form
	r.SetPathValue("id", "76")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/bulk-assurance-visit-tasks", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForClientTab(client, template)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	fmt.Print(returnedError)
	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Nil(returnedError)
}
