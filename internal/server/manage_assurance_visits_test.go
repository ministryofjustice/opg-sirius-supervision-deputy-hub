package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockManageAssuranceVisitInformation struct {
	count               int
	lastCtx             sirius.Context
	err                 error
	userDetails         sirius.UserDetails
	visitors            sirius.Visitors
	VisitRagRatingTypes []sirius.VisitRagRatingTypes
	VisitOutcomeTypes   []sirius.VisitOutcomeTypes
	visit               sirius.AssuranceVisit
}

func (m *mockManageAssuranceVisitInformation) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userDetails, m.err
}

func (m *mockManageAssuranceVisitInformation) UpdateAssuranceVisit(ctx sirius.Context, requestedDate string, userId, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func (m *mockManageAssuranceVisitInformation) GetVisitors(ctx sirius.Context) (sirius.Visitors, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.visitors, m.err
}

func (m *mockManageAssuranceVisitInformation) GetVisitRagRatingTypes(ctx sirius.Context) ([]sirius.VisitRagRatingTypes, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.VisitRagRatingTypes, m.err
}

func (m *mockManageAssuranceVisitInformation) GetVisitOutcomeTypes(ctx sirius.Context) ([]sirius.VisitOutcomeTypes, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.VisitOutcomeTypes, m.err
}

func (m *mockManageAssuranceVisitInformation) GetAssuranceVisitById(ctx sirius.Context, deputyId, visitId int) (sirius.AssuranceVisit, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.visit, m.err
}

func TestPostManageAssuranceVisit(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageAssuranceVisitInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/manage-assurance-visit", strings.NewReader("{requestedDate:'2022/10/20', requestedBy:22}"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/manage-assurance-visit", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManageAssuranceVisit(client, template)(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Nil(returnedError)
}
