package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockChangeECMInformation struct {
	count          int
	lastCtx        sirius.Context
	err            error
	DeputyDetails  sirius.DeputyDetails
	EcmTeamDetails []sirius.TeamMember
	Error          string
	Errors         sirius.ValidationErrors
	Success        bool
	DefaultPaTeam  int
}

func (m *mockChangeECMInformation) GetDeputyTeamMembers(ctx sirius.Context, deputyId int, deputyDetails sirius.DeputyDetails) ([]sirius.TeamMember, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.EcmTeamDetails, m.err
}

func (m *mockChangeECMInformation) ChangeECM(ctx sirius.Context, changeEcmForm sirius.ExecutiveCaseManagerOutgoing, deputyDetails sirius.DeputyDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestGetChangeECM(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangeECMInformation{}
	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForChangeECM(client, defaultPATeam, template)
	err := handler(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, client.count)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(changeECMHubVars{
		Path: "/path",
	}, template.lastVars)
}

func TestPostChangeECM(t *testing.T) {
	assert := assert.New(t)
	client := &mockChangeECMInformation{}
	defaultPATeam := 23

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/ecm", strings.NewReader("{ecmId:26}"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/ecm", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForChangeECM(client, defaultPATeam, template)(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Nil(returnedError)
}

func TestPostChangeECMReturnsErrorWithNoECM(t *testing.T) {
	assert := assert.New(t)
	client := &mockChangeECMInformation{}
	defaultPATeam := 23

	validationErrors := sirius.ValidationErrors{
		"Change ECM": {"": "Select an executive case manager"},
	}

	client.err = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/ecm", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/ecm", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForChangeECM(client, defaultPATeam, template)(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationError := sirius.ValidationError{
		Errors: sirius.ValidationErrors{
			"Change ECM": {
				"": "Select an executive case manager",
			},
		},
	}

	assert.Equal(expectedValidationError, returnedError)
}
