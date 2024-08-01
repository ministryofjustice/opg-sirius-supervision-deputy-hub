package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockChangeECMInformation struct {
	count                   int
	lastCtx                 sirius.Context
	GetDeputyTeamMembersErr error
	ChangeECMErr            error
	DeputyDetails           sirius.DeputyDetails
	EcmTeamDetails          []model.TeamMember
	DefaultPaTeam           int
}

func (m *mockChangeECMInformation) GetDeputyTeamMembers(ctx sirius.Context, deputyId int, deputyDetails sirius.DeputyDetails) ([]model.TeamMember, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.EcmTeamDetails, m.GetDeputyTeamMembersErr
}

func (m *mockChangeECMInformation) ChangeECM(ctx sirius.Context, changeEcmForm sirius.ExecutiveCaseManagerOutgoing, deputyDetails sirius.DeputyDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.ChangeECMErr
}

func TestGetChangeECM(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangeECMInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	app := AppVars{PageName: "Change Executive Case Manager"}

	handler := renderTemplateForChangeECM(client, template)
	err := handler(app, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, client.count)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(changeECMHubVars{AppVars: app}, template.lastVars)
}

func TestPostChangeECM(t *testing.T) {
	assert := assert.New(t)
	client := &mockChangeECMInformation{}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/ecm", strings.NewReader("{ecmId:26}"))

	form := url.Values{}
	form.Add("select-ecm", "26")
	r.PostForm = form

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/ecm", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForChangeECM(client, template)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(returnedError, Redirect("/76?success=ecm"))
}

func TestPostChangeECMReturnsErrorWithNoECM(t *testing.T) {
	assert := assert.New(t)
	client := &mockChangeECMInformation{}

	form := url.Values{}
	form.Add("select-ecm", "")

	template := &mockTemplates{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/ecm", strings.NewReader(form.Encode()))

	returnedError := renderTemplateForChangeECM(client, template)(AppVars{}, w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"select-ecm": {"": "Select an executive case manager"},
	}

	assert.Equal(changeECMHubVars{
		AppVars: AppVars{
			PageName: "Change Executive Case Manager",
			Errors:   expectedValidationErrors,
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestPutChangeECMReturnsStatusMethodError(t *testing.T) {
	assert := assert.New(t)
	client := &mockChangeECMInformation{}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/76/ecm", strings.NewReader(""))

	returnedError := renderTemplateForChangeECM(client, template)(AppVars{}, w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), returnedError)
}

func TestPostChangeECMReturnsTimeoutError(t *testing.T) {
	assert := assert.New(t)
	client := &mockChangeECMInformation{}

	template := &mockTemplates{}

	client.ChangeECMErr = StatusError(http.StatusGatewayTimeout)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/ecm", strings.NewReader("{ecmId:26}"))

	form := url.Values{}
	form.Add("select-ecm", "26")
	r.PostForm = form

	returnedError := renderTemplateForChangeECM(client, template)(AppVars{}, w, r)

	assert.Equal(StatusError(http.StatusGatewayTimeout), returnedError)
}

func TestChangeECMHandlesErrorsInOtherClientFiles(t *testing.T) {
	returnedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockChangeECMInformation
	}{
		{
			Client: &mockChangeECMInformation{
				GetDeputyTeamMembersErr: returnedError,
			},
		},
		{
			Client: &mockChangeECMInformation{
				ChangeECMErr: returnedError,
			},
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {

			client := tc.Client
			template := &mockTemplates{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/76/ecm", strings.NewReader(""))

			form := url.Values{}
			form.Add("select-ecm", "26")
			r.PostForm = form

			changeEcmReturnedError := renderTemplateForChangeECM(client, template)(AppVars{}, w, r)
			assert.Equal(t, returnedError, changeEcmReturnedError)
		})
	}
}
