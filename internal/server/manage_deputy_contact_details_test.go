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

type mockManageDeputyContactDetailsInformation struct {
	count     int
	lastCtx   sirius.Context
	updateErr error
}

func (m *mockManageDeputyContactDetailsInformation) UpdateDeputyContactDetails(ctx sirius.Context, _ int, _ sirius.DeputyContactDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.updateErr
}

var manageContactAppVars = AppVars{
	DeputyDetails: sirius.DeputyDetails{
		ID: 123,
	},
}

func TestGetManageDeputyDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageDeputyContactDetailsInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForManageDeputyContactDetails(client, template)
	err := handler(manageContactAppVars, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestPostManageDeputyDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageDeputyContactDetailsInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))

	var redirect error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		redirect = renderTemplateForManageDeputyContactDetails(client, template)(manageContactAppVars, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(redirect, Redirect("/123?success=deputyDetails"))
}

func TestManageDeputyDetailsValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageDeputyContactDetailsInformation{}

	validationErrors := sirius.ValidationErrors{
		"firstname": {
			"stringLengthTooLong": "What sirius gives us",
		},
	}

	client.updateErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManageDeputyContactDetails(client, template)(manageContactAppVars, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(1, client.count)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(manageDeputyContactDetailsVars{
		AppVars: AppVars{
			DeputyDetails: manageContactAppVars.DeputyDetails,
			Errors:        util.RenameErrors(validationErrors),
			PageName:      "Manage deputy contact details",
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestDeputyContactDetailsHandlesErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageDeputyContactDetailsInformation{
		updateErr: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))

	returnedError := renderTemplateForManageDeputyContactDetails(client, template)(AppVars{}, w, r)

	assert.Equal(client.updateErr, returnedError)

}
