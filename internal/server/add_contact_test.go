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

type mockContactInformation struct {
	count   int
	lastCtx sirius.Context
	err     error
}

func (m *mockContactInformation) AddContact(ctx sirius.Context, deputyId int, contact sirius.Contact) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

var addContactAppVars = AppVars{
	DeputyDetails: sirius.DeputyDetails{
		ID: 123,
	},
}

func TestGetContact(t *testing.T) {
	assert := assert.New(t)

	client := &mockContactInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddContact(client, template)
	err := handler(addContactAppVars, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostAddContact(t *testing.T) {
	assert := assert.New(t)
	client := &mockContactInformation{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForAddContact(client, nil)(addContactAppVars, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/123/contacts?success=newContact"))
}

func TestAddContactEmptyValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockContactInformation{}

	validationErrors := sirius.ValidationErrors{
		"contactName": {
			"isEmpty": "Enter a name",
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
		returnedError = renderTemplateForAddContact(client, template)(addContactAppVars, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"contactName": {
			"isEmpty": "Enter a name",
		},
	}

	assert.Equal(addContactVars{
		AppVars: AppVars{
			DeputyDetails: addContactAppVars.DeputyDetails,
			Errors:        expectedValidationErrors,
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}
