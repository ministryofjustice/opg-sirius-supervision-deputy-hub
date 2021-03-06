package server

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockManageAssuranceVisitInformation struct {
	mock.Mock
}

func (m *mockManageAssuranceVisitInformation) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	args := m.Called(ctx)

	return args.Get(0).(sirius.UserDetails), args.Error(1)
}

func (m *mockManageAssuranceVisitInformation) GetVisitOutcomeTypes(ctx sirius.Context) ([]sirius.VisitOutcomeTypes, error) {
	args := m.Called(ctx)

	return args.Get(0).([]sirius.VisitOutcomeTypes), args.Error(1)
}

func (m *mockManageAssuranceVisitInformation) GetVisitRagRatingTypes(ctx sirius.Context) ([]sirius.VisitRagRatingTypes, error) {
	args := m.Called(ctx)

	return args.Get(0).([]sirius.VisitRagRatingTypes), args.Error(1)
}

func (m *mockManageAssuranceVisitInformation) GetVisitors(ctx sirius.Context) (sirius.Visitors, error) {
	args := m.Called(ctx)

	return args.Get(0).(sirius.Visitors), args.Error(1)
}

func (m *mockManageAssuranceVisitInformation) UpdateAssuranceVisit(ctx sirius.Context, form sirius.AssuranceVisitDetails, deputyId, visitId int) error {
	args := m.Called(ctx, form, deputyId, visitId)

	return args.Error(0)
}

func (m *mockManageAssuranceVisitInformation) GetAssuranceVisitById(ctx sirius.Context, deputyId, visitId int) (sirius.AssuranceVisit, error) {
	args := m.Called(ctx, deputyId, visitId)

	return args.Get(0).(sirius.AssuranceVisit), args.Error(1)
}

func TestGetManageAssurance(t *testing.T) {
	assert := assert.New(t)

	deputyDetails := sirius.DeputyDetails{ID: 123}
	visitors := sirius.Visitors{sirius.Visitor{ID: 1, Name: "Mr Visitor"}}
	visitRagRatingTypes := []sirius.VisitRagRatingTypes{{Handle: "x", Label: "y"}}
	visitOutcomeTypes := []sirius.VisitOutcomeTypes{{Handle: "x", Label: "w"}}
	visit := sirius.AssuranceVisit{Id: 1, RequestedDate: "2022-01-02", RequestedBy: sirius.User{UserId: 2}}

	client := &mockManageAssuranceVisitInformation{}
	client.On("GetAssuranceVisitById", mock.Anything, 0, 0).Return(visit, nil)
	client.On("GetUserDetails", mock.Anything).Return(sirius.UserDetails{Roles: []string{"Finance Manager"}}, nil)
	client.On("GetVisitors", mock.Anything).Return(visitors, nil)
	client.On("GetVisitRagRatingTypes", mock.Anything).Return(visitRagRatingTypes, nil)
	client.On("GetVisitOutcomeTypes", mock.Anything).Return(visitOutcomeTypes, nil)

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForManageAssuranceVisit(client, template)
	err := handler(sirius.PermissionSet{}, sirius.DeputyDetails{ID: 123}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(ManageAssuranceVisitVars{
		DeputyDetails:       deputyDetails,
		Visitors:            visitors,
		VisitRagRatingTypes: visitRagRatingTypes,
		VisitOutcomeTypes:   visitOutcomeTypes,
		Visit:               visit,
	}, template.lastVars)
}

func TestPostManageAssuranceVisit(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageAssuranceVisitInformation{}
	client.On("GetAssuranceVisitById", mock.Anything, 123, 1).Return(sirius.AssuranceVisit{}, nil)
	client.On("GetUserDetails", mock.Anything).Return(sirius.UserDetails{Roles: []string{"Finance Manager"}}, nil)
	client.On("GetVisitors", mock.Anything).Return(sirius.Visitors{}, nil)
	client.On("GetVisitRagRatingTypes", mock.Anything).Return([]sirius.VisitRagRatingTypes{}, nil)
	client.On("GetVisitOutcomeTypes", mock.Anything).Return([]sirius.VisitOutcomeTypes{}, nil)
	client.On("UpdateAssuranceVisit", mock.Anything, sirius.AssuranceVisitDetails{}, 123, 1).Return(nil)
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/assurance-visits/1", strings.NewReader("{commissionedDate:'2200/10/20'}"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var redirect error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/assurance-visits/{visitId}", func(w http.ResponseWriter, r *http.Request) {
		redirect = renderTemplateForManageAssuranceVisit(client, template)(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(Redirect("/123/assurance-visits?success=manageAssuranceVisit"), redirect)
}
