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
		Clients: sirius.ClientList{
			Clients: []sirius.DeputyClient(nil),
			Pages: sirius.Page{
				PageCurrent: 0,
				PageTotal:   0,
			},
			TotalClients: 0,
			Metadata: sirius.Metadata{
				TotalActiveClients: 0,
			},
		},
		SuccessMessage: "",
		ListPage: ListPage{
			AppVars: AppVars{
				Path:      "",
				XSRFToken: "",
				UserDetails: sirius.UserDetails{
					ID:       0,
					Roles:    []string(nil),
					Username: "",
				},
				DeputyDetails: sirius.DeputyDetails{
					ID:              0,
					DeputyFirstName: "",
					DeputySurname:   "",
					DeputyCasrecId:  0,
					DisplayName:     "",
					CanDelete:       false,
					DeputyNumber:    0,
					DeputySubType: sirius.DeputySubType{
						SubType: "",
					},
					DeputyStatus: "",
					DeputyImportantInformation: sirius.DeputyImportantInformation{
						Id: 0,
						AnnualBillingInvoice: sirius.HandleLabel{
							Handle: "",
							Label:  "",
						},
						APAD: sirius.HandleLabel{
							Handle: "",
							Label:  "",
						},
						BankCharges: sirius.HandleLabel{
							Handle: "",
							Label:  "",
						},
						Complaints: sirius.HandleLabel{
							Handle: "",
							Label:  "",
						},
						IndependentVisitorCharges: sirius.HandleLabel{
							Handle: "",
							Label:  "",
						},
						MonthlySpreadsheet: sirius.HandleLabel{
							Handle: "",
							Label:  "",
						},
						PanelDeputy: false,
					},
				},
				SuccessMessage: "",
				PageName:       "Clients",
				Error:          "",
				Errors:         expectedValidationErrors,
				EnvironmentVars: EnvironmentVars{
					Port:            "",
					WebDir:          "",
					SiriusURL:       "",
					SiriusPublicURL: "",
					FirmHubURL:      "",
					Prefix:          "",
					DefaultPaTeam:   0,
					DefaultProTeam:  0,
					Features:        []string(nil),
				},
			},
			AppliedFilters: []string(nil),
			Sort: urlbuilder.Sort{
				OrderBy:    "surname",
				Descending: false,
			},
			Error: "",
			Pagination: paginate.Pagination{
				CurrentPage:     0,
				TotalPages:      0,
				TotalElements:   0,
				ElementsPerPage: 25,
				ElementName:     "clients",
				PerPageOptions:  []int{25, 50, 100},
				UrlBuilder: urlbuilder.UrlBuilder{
					OriginalPath:    "clients",
					SelectedPerPage: 25,
					SelectedFilters: []urlbuilder.Filter{{
						Name:                  "order-status",
						SelectedValues:        []string(nil),
						ClearBetweenTeamViews: false},
						{
							Name:                  "accommodation",
							SelectedValues:        []string(nil),
							ClearBetweenTeamViews: false},
						{
							Name:                  "supervision-level",
							SelectedValues:        []string(nil),
							ClearBetweenTeamViews: false},
					},
					SelectedSort: urlbuilder.Sort{
						OrderBy:    "surname",
						Descending: false},
				},
			},
			PerPage: 25,
			UrlBuilder: urlbuilder.UrlBuilder{
				OriginalPath:    "clients",
				SelectedPerPage: 25,
				SelectedFilters: []urlbuilder.Filter{{
					Name:                  "order-status",
					SelectedValues:        []string(nil),
					ClearBetweenTeamViews: false,
				},
					{
						Name:                  "accommodation",
						SelectedValues:        []string(nil),
						ClearBetweenTeamViews: false,
					},
					{
						Name:                  "supervision-level",
						SelectedValues:        []string(nil),
						ClearBetweenTeamViews: false},
				},
				SelectedSort: urlbuilder.Sort{
					OrderBy:    "surname",
					Descending: false,
				},
			},
		},
		FilterByOrderStatus: FilterByOrderStatus{
			ListPage: ListPage{
				AppVars: AppVars{
					Path:      "",
					XSRFToken: "",
					UserDetails: sirius.UserDetails{
						ID:       0,
						Roles:    []string(nil),
						Username: "",
					},
					DeputyDetails: sirius.DeputyDetails{
						ID:              0,
						DeputyFirstName: "",
						DeputySurname:   "",
						DeputyCasrecId:  0,
						DisplayName:     "",
						CanDelete:       false,
						DeputyNumber:    0,
						DeputySubType: sirius.DeputySubType{
							SubType: "",
						},
						DeputyStatus: "",
						DeputyImportantInformation: sirius.DeputyImportantInformation{
							Id: 0,
							AnnualBillingInvoice: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							APAD: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							BankCharges: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							Complaints: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							IndependentVisitorCharges: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							MonthlySpreadsheet: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							PanelDeputy: false,
							ReportSystem: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							OtherImportantInformation: "",
						},
						OrganisationName:                 "",
						OrganisationTeamOrDepartmentName: "",
						Email:                            "",
						PhoneNumber:                      "",
						AddressLine1:                     "",
						AddressLine2:                     "",
						AddressLine3:                     "",
						Town:                             "",
						County:                           "",
						Postcode:                         "",
						ExecutiveCaseManager: sirius.ExecutiveCaseManager{
							EcmId:     0,
							EcmName:   "",
							IsDefault: false,
						},
						DeputyType: sirius.DeputyType{
							Handle: "",
							Label:  "",
						},
						Firm: sirius.Firm{
							FirmName:   "",
							FirmId:     0,
							FirmNumber: 0,
						},
					},
					SuccessMessage: "",
					PageName:       "",
					Error:          "",
					Errors:         sirius.ValidationErrors(nil),
					EnvironmentVars: EnvironmentVars{
						Port:            "",
						WebDir:          "",
						SiriusURL:       "",
						SiriusPublicURL: "",
						FirmHubURL:      "",
						Prefix:          "",
						DefaultPaTeam:   0,
						DefaultProTeam:  0,
						Features:        []string(nil)}},
				AppliedFilters: []string(nil),
				Sort: urlbuilder.Sort{
					OrderBy:    "",
					Descending: false,
				},
				Error: "",
				Pagination: paginate.Pagination{
					CurrentPage:     0,
					TotalPages:      0,
					TotalElements:   0,
					ElementsPerPage: 0,
					ElementName:     "",
					PerPageOptions:  []int(nil),
					UrlBuilder:      paginate.UrlBuilder(nil)},
				PerPage: 0,
				UrlBuilder: urlbuilder.UrlBuilder{
					OriginalPath:    "",
					SelectedPerPage: 0,
					SelectedFilters: []urlbuilder.Filter(nil),
					SelectedSort: urlbuilder.Sort{
						OrderBy:    "",
						Descending: false},
				}},
			OrderStatusOptions: []model.RefData{
				{
					Handle:     "ACTIVE",
					Label:      "Active",
					Deprecated: false,
				},
				{
					Handle:     "CLOSED",
					Label:      "Closed",
					Deprecated: false,
				},
			},
			SelectedOrderStatuses: []string(nil),
			OrderStatuses: []model.OrderStatus{
				{
					Handle:      "ACTIVE",
					Incomplete:  "Active",
					Category:    "Active",
					Complete:    "Active",
					StatusCount: 0,
				},
				{
					Handle:      "CLOSED",
					Incomplete:  "Closed",
					Category:    "Closed",
					Complete:    "Closed",
					StatusCount: 0,
				},
			},
		},
		FilterByAccommodation: FilterByAccommodation{
			ListPage: ListPage{
				AppVars: AppVars{
					Path:      "",
					XSRFToken: "",
					UserDetails: sirius.UserDetails{
						ID:       0,
						Roles:    []string(nil),
						Username: "",
					},
					DeputyDetails: sirius.DeputyDetails{
						ID:              0,
						DeputyFirstName: "",
						DeputySurname:   "",
						DeputyCasrecId:  0,
						DisplayName:     "",
						CanDelete:       false,
						DeputyNumber:    0,
						DeputySubType: sirius.DeputySubType{
							SubType: "",
						},
						DeputyStatus: "",
						DeputyImportantInformation: sirius.DeputyImportantInformation{
							Id: 0,
							AnnualBillingInvoice: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							APAD: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							BankCharges: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							Complaints: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							IndependentVisitorCharges: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							MonthlySpreadsheet: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							PanelDeputy: false,
							ReportSystem: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							OtherImportantInformation: "",
						},
						OrganisationName:                 "",
						OrganisationTeamOrDepartmentName: "",
						Email:                            "",
						PhoneNumber:                      "",
						AddressLine1:                     "",
						AddressLine2:                     "",
						AddressLine3:                     "",
						Town:                             "",
						County:                           "",
						Postcode:                         "",
						ExecutiveCaseManager: sirius.ExecutiveCaseManager{
							EcmId:     0,
							EcmName:   "",
							IsDefault: false,
						},
						DeputyType: sirius.DeputyType{
							Handle: "",
							Label:  "",
						},
						Firm: sirius.Firm{
							FirmName:   "",
							FirmId:     0,
							FirmNumber: 0,
						}},
					SuccessMessage: "",
					PageName:       "",
					Error:          "",
					Errors:         sirius.ValidationErrors(nil),
					EnvironmentVars: EnvironmentVars{
						Port:            "",
						WebDir:          "",
						SiriusURL:       "",
						SiriusPublicURL: "",
						FirmHubURL:      "",
						Prefix:          "",
						DefaultPaTeam:   0,
						DefaultProTeam:  0,
						Features:        []string(nil)}},
				AppliedFilters: []string(nil),
				Sort: urlbuilder.Sort{
					OrderBy:    "",
					Descending: false,
				},
				Error: "",
				Pagination: paginate.Pagination{
					CurrentPage:     0,
					TotalPages:      0,
					TotalElements:   0,
					ElementsPerPage: 0,
					ElementName:     "",
					PerPageOptions:  []int(nil),
					UrlBuilder:      paginate.UrlBuilder(nil)},
				PerPage: 0,
				UrlBuilder: urlbuilder.UrlBuilder{
					OriginalPath:    "",
					SelectedPerPage: 0,
					SelectedFilters: []urlbuilder.Filter(nil),
					SelectedSort:    urlbuilder.Sort{OrderBy: "", Descending: false}}},
			AccommodationTypes:         []model.RefData(nil),
			SelectedAccommodationTypes: []string(nil),
		},
		FilterBySupervisionLevel: FilterBySupervisionLevel{
			ListPage: ListPage{
				AppVars: AppVars{
					Path:      "",
					XSRFToken: "",
					UserDetails: sirius.UserDetails{
						ID:       0,
						Roles:    []string(nil),
						Username: ""},
					DeputyDetails: sirius.DeputyDetails{
						ID:              0,
						DeputyFirstName: "",
						DeputySurname:   "",
						DeputyCasrecId:  0,
						DisplayName:     "",
						CanDelete:       false,
						DeputyNumber:    0,
						DeputySubType: sirius.DeputySubType{
							SubType: "",
						},
						DeputyStatus: "",
						DeputyImportantInformation: sirius.DeputyImportantInformation{
							Id: 0,
							AnnualBillingInvoice: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							APAD: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							BankCharges: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							Complaints: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							IndependentVisitorCharges: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							MonthlySpreadsheet: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							PanelDeputy: false,
							ReportSystem: sirius.HandleLabel{
								Handle: "",
								Label:  "",
							},
							OtherImportantInformation: "",
						},
						OrganisationName:                 "",
						OrganisationTeamOrDepartmentName: "",
						Email:                            "",
						PhoneNumber:                      "",
						AddressLine1:                     "",
						AddressLine2:                     "",
						AddressLine3:                     "",
						Town:                             "",
						County:                           "",
						Postcode:                         "",
						ExecutiveCaseManager: sirius.ExecutiveCaseManager{
							EcmId:     0,
							EcmName:   "",
							IsDefault: false,
						},
						DeputyType: sirius.DeputyType{
							Handle: "",
							Label:  "",
						},
						Firm: sirius.Firm{
							FirmName:   "",
							FirmId:     0,
							FirmNumber: 0}},
					SuccessMessage: "",
					PageName:       "",
					Error:          "",
					Errors:         sirius.ValidationErrors(nil),
					EnvironmentVars: EnvironmentVars{
						Port:            "",
						WebDir:          "",
						SiriusURL:       "",
						SiriusPublicURL: "",
						FirmHubURL:      "",
						Prefix:          "",
						DefaultPaTeam:   0,
						DefaultProTeam:  0,
						Features:        []string(nil)}},
				AppliedFilters: []string(nil),
				Sort: urlbuilder.Sort{
					OrderBy:    "",
					Descending: false,
				},
				Error: "",
				Pagination: paginate.Pagination{
					CurrentPage:     0,
					TotalPages:      0,
					TotalElements:   0,
					ElementsPerPage: 0,
					ElementName:     "",
					PerPageOptions:  []int(nil),
					UrlBuilder:      paginate.UrlBuilder(nil)},
				PerPage: 0,
				UrlBuilder: urlbuilder.UrlBuilder{
					OriginalPath:    "",
					SelectedPerPage: 0,
					SelectedFilters: []urlbuilder.Filter(nil),
					SelectedSort: urlbuilder.Sort{
						OrderBy:    "",
						Descending: false}}},
			SupervisionLevels:         []model.RefData(nil),
			SelectedSupervisionLevels: []string(nil),
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
