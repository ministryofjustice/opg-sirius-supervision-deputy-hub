package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockManageAssuranceClient struct {
	mock.Mock
}

func (m *mockManageAssuranceClient) GetRefData(ctx sirius.Context, refDataTypeUrl string) ([]model.RefData, error) {
	args := m.Called(ctx, refDataTypeUrl)

	return args.Get(0).([]model.RefData), args.Error(1)
}

func (m *mockManageAssuranceClient) GetVisitors(ctx sirius.Context) ([]model.Visitor, error) {
	args := m.Called(ctx)

	return args.Get(0).([]model.Visitor), args.Error(1)
}

func (m *mockManageAssuranceClient) UpdateAssurance(ctx sirius.Context, form sirius.UpdateAssuranceDetails, deputyId, visitId int) error {
	args := m.Called(ctx, form, deputyId, visitId)

	return args.Error(0)
}

func (m *mockManageAssuranceClient) GetAssuranceById(ctx sirius.Context, deputyId, visitId int) (model.Assurance, error) {
	args := m.Called(ctx, deputyId, visitId)

	return args.Get(0).(model.Assurance), args.Error(1)
}

var manageAssuranceAppVars = AppVars{
	DeputyDetails: sirius.DeputyDetails{
		ID: 123,
	},
}

func TestGetManageAssurance(t *testing.T) {
	assert := assert.New(t)

	manageAssuranceAppVars.PageName = "Manage assurance visit"

	visitors := []model.Visitor{{ID: 1, Name: "Mr Visitor"}}
	ragRatingTypes := []model.RefData{{Handle: "x", Label: "y"}}
	visitOutcomeTypes := []model.RefData{{Handle: "x", Label: "w"}}
	pdrOutcomeTypes := []model.RefData{{Handle: "x", Label: "z"}}
	assurance := model.Assurance{Id: 1, RequestedDate: "2022-01-02", RequestedBy: model.User{ID: 2}}

	client := &mockManageAssuranceClient{}
	client.On("GetAssuranceById", mock.Anything, manageAssuranceAppVars.DeputyId(), 0).Return(assurance, nil)
	client.On("GetUserDetails", mock.Anything).Return(sirius.UserDetails{Roles: []string{"Finance Manager"}}, nil)
	client.On("GetVisitors", mock.Anything).Return(visitors, nil)
	client.On("GetRefData", mock.Anything, "/ragRating").Return(ragRatingTypes, nil)
	client.On("GetRefData", mock.Anything, "/visitOutcome").Return(visitOutcomeTypes, nil)
	client.On("GetRefData", mock.Anything, "/pdrOutcome").Return(pdrOutcomeTypes, nil)

	visitTemplate := &mockTemplates{}
	pdrTemplate := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForManageAssurance(client, visitTemplate, pdrTemplate)
	err := handler(manageAssuranceAppVars, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, visitTemplate.count)
	assert.Equal("page", visitTemplate.lastName)
	assert.Equal(ManageAssuranceVars{
		AppVars:           manageAssuranceAppVars,
		Visitors:          visitors,
		RagRatingTypes:    ragRatingTypes,
		VisitOutcomeTypes: visitOutcomeTypes,
		PdrOutcomeTypes:   pdrOutcomeTypes,
		Assurance:         assurance,
	}, visitTemplate.lastVars)
}

func TestPostManageAssurance(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageAssuranceClient{}
	client.On("GetAssuranceById", mock.Anything, 123, 1).Return(model.Assurance{}, nil)
	client.On("GetUserDetails", mock.Anything).Return(sirius.UserDetails{Roles: []string{"Finance Manager"}}, nil)
	client.On("GetVisitors", mock.Anything).Return([]model.Visitor{}, nil)
	client.On("GetRefData", mock.Anything, "/ragRating").Return([]model.RefData{}, nil)
	client.On("GetRefData", mock.Anything, "/visitOutcome").Return([]model.RefData{}, nil)
	client.On("GetRefData", mock.Anything, "/pdrOutcome").Return([]model.RefData{}, nil)
	client.On("UpdateAssurance", mock.Anything, sirius.UpdateAssuranceDetails{}, 123, 1).Return(nil)
	visitTemplate := &mockTemplates{}
	pdrTemplate := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/assurances/1", strings.NewReader("{commissionedDate:'2200/10/20'}"))

	var redirect error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/assurances/{visitId}", func(w http.ResponseWriter, r *http.Request) {
		redirect = renderTemplateForManageAssurance(client, visitTemplate, pdrTemplate)(manageAssuranceAppVars, w, r)
	})

	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(Redirect("/123/assurances?success=manageVisit"), redirect)
}

func TestGetManagePDR(t *testing.T) {
	assert := assert.New(t)

	manageAssuranceAppVars.PageName = "Manage PDR"

	visitors := []model.Visitor{{ID: 1, Name: "Mr Visitor"}}
	ragRatingTypes := []model.RefData{{Handle: "x", Label: "y"}}
	visitOutcomeTypes := []model.RefData{{Handle: "x", Label: "w"}}
	pdrOutcomeTypes := []model.RefData{{Handle: "x", Label: "z"}}
	assurance := model.Assurance{Id: 1, Type: model.AssuranceType{Handle: "PDR", Label: "PDR"}, RequestedDate: "2022-01-02", RequestedBy: model.User{ID: 2}}

	client := &mockManageAssuranceClient{}
	client.On("GetAssuranceById", mock.Anything, manageAssuranceAppVars.DeputyId(), 0).Return(assurance, nil)
	client.On("GetUserDetails", mock.Anything).Return(sirius.UserDetails{Roles: []string{"Finance Manager"}}, nil)
	client.On("GetVisitors", mock.Anything).Return(visitors, nil)
	client.On("GetRefData", mock.Anything, "/ragRating").Return(ragRatingTypes, nil)
	client.On("GetRefData", mock.Anything, "/visitOutcome").Return(visitOutcomeTypes, nil)
	client.On("GetRefData", mock.Anything, "/pdrOutcome").Return(pdrOutcomeTypes, nil)

	visitTemplate := &mockTemplates{}
	pdrTemplate := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForManageAssurance(client, visitTemplate, pdrTemplate)
	err := handler(manageAssuranceAppVars, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, pdrTemplate.count)
	assert.Equal("page", pdrTemplate.lastName)
	assert.Equal(ManageAssuranceVars{
		AppVars:           manageAssuranceAppVars,
		Visitors:          visitors,
		RagRatingTypes:    ragRatingTypes,
		VisitOutcomeTypes: visitOutcomeTypes,
		PdrOutcomeTypes:   pdrOutcomeTypes,
		Assurance:         assurance,
	}, pdrTemplate.lastVars)
}

func TestPostManagePDR(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageAssuranceClient{}
	client.On("GetAssuranceById", mock.Anything, 123, 1).Return(model.Assurance{Type: model.AssuranceType{Handle: "PDR", Label: "PDR"}}, nil)
	client.On("GetUserDetails", mock.Anything).Return(sirius.UserDetails{Roles: []string{"Finance Manager"}}, nil)
	client.On("GetVisitors", mock.Anything).Return([]model.Visitor{}, nil)
	client.On("GetRefData", mock.Anything, "/ragRating").Return([]model.RefData{}, nil)
	client.On("GetRefData", mock.Anything, "/visitOutcome").Return([]model.RefData{}, nil)
	client.On("GetRefData", mock.Anything, "/pdrOutcome").Return([]model.RefData{}, nil)
	client.On("UpdateAssurance", mock.Anything, sirius.UpdateAssuranceDetails{}, 123, 1).Return(nil)
	visitTemplate := &mockTemplates{}
	pdrTemplate := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/assurances/1", strings.NewReader("{commissionedDate:'2200/10/20'}"))

	var redirect error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/assurances/{visitId}", func(w http.ResponseWriter, r *http.Request) {
		redirect = renderTemplateForManageAssurance(client, visitTemplate, pdrTemplate)(manageAssuranceAppVars, w, r)
	})

	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(Redirect("/123/assurances?success=managePDR"), redirect)
}

func TestParseVisitForm(t *testing.T) {
	assert := assert.New(t)

	form := sirius.UpdateAssuranceDetails{
		CommissionedDate:   "2020-01-01",
		ReportDueDate:      "2020-01-02",
		ReportReceivedDate: "2020-01-03",
		VisitOutcome:       "Cancelled",
		PdrOutcome:         "Cancelled",
		ReportReviewDate:   "2020-01-04",
		ReportMarkedAs:     "Successful",
		VisitorAllocated:   "John Johnson",
		ReviewedBy:         1,
		Note:               "Test notes",
	}

	expectedAssurance := model.Assurance{
		CommissionedDate:   "2020-01-01",
		ReportDueDate:      "2020-01-02",
		ReportReceivedDate: "2020-01-03",
		VisitOutcome:       model.VisitOutcomeType{Label: "Cancelled"},
		PdrOutcome:         model.PdrOutcomeType{Label: "Cancelled"},
		ReportReviewDate:   "2020-01-04",
		ReportMarkedAs:     model.RAGRating{Label: "Successful"},
		Note:               "Test notes",
		VisitorAllocated:   "John Johnson",
		ReviewedBy:         model.User{ID: 1},
	}

	assert.Equal(expectedAssurance, parseAssuranceForm(form))
}
