package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockEditDeputyHubInformation struct {
	count   int
	lastCtx sirius.Context
	err     error
}

func (m *mockEditDeputyHubInformation) EditDeputyDetails(ctx sirius.Context, deputyDetails sirius.DeputyDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestNavigateToEditDeputyHub(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditDeputyHubInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForEditDeputyHub(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestPostEditDeputyHub(t *testing.T) {
	assert := assert.New(t)
	client := &mockEditDeputyHubInformation{}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForEditDeputyHub(client, template)(AppVars{DeputyDetails: testDeputy}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(returnedError, Redirect("/123?success=teamDetails"))
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

	app := AppVars{
		PageName: "Manage team details",
		Errors:   util.RenameErrors(validationErrors),
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
	returnedError := renderTemplateForEditDeputyHub(client, template)(app, w, r)

	assert.Equal(1, client.count)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editDeputyHubVars{AppVars: app}, template.lastVars)

	assert.Nil(returnedError)
}

func TestEditDeputyHubHandlesErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockEditDeputyHubInformation{
		err: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	returnedError := renderTemplateForEditDeputyHub(client, template)(AppVars{}, w, r)

	assert.Equal(client.err, returnedError)
}
