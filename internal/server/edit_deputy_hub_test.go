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

type mockEditDeputyHubInformation struct {
	count            int
	lastCtx          sirius.Context
	err              error
	deputyClientData sirius.DeputyClientDetails
	ariaSorting      sirius.AriaSorting
	userDetails      sirius.UserDetails
}

func (m *mockEditDeputyHubInformation) EditDeputyDetails(ctx sirius.Context, deputyDetails sirius.DeputyDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func (m *mockEditDeputyHubInformation) GetDeputyClients(ctx sirius.Context, deputyId int, deputyType string, columnBeingSorted string, sortOrder string) (sirius.DeputyClientDetails, sirius.AriaSorting, int, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyClientData, m.ariaSorting, 0, m.err
}

func (m *mockEditDeputyHubInformation) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userDetails, m.err
}

func TestNavigateToEditDeputyHub(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditDeputyHubInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHub(client, template)
	err := handler(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestEditDeputyValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockEditDeputyHubInformation{}

	validationErrors := sirius.ValidationErrors{
		"organisationName": {
			"stringLengthTooLong": "What sirius gives us",
		},
	}

	client.err = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForEditDeputyHub(client, template)(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editDeputyHubVars{
		Path:   "/133",
		Errors: validationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}
