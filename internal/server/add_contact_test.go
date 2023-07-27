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

func TestGetContact(t *testing.T) {
	assert := assert.New(t)

	client := &mockContactInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddContact(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

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
		returnedError = renderTemplateForAddContact(client, nil)(sirius.DeputyDetails{}, w, r)
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
		"email": {
			"isEmpty": "Enter an email address",
		},
		"phoneNumber": {
			"isEmpty": "Enter a telephone number",
		},
		"isMainContact": {
			"isEmpty": "Select whether this contact is a main contact",
		},
		"isNamedDeputy": {
			"isEmpty": "Select whether this contact is the named deputy",
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
		returnedError = renderTemplateForAddContact(client, template)(sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"contactName": {
			"isEmpty": "Enter a name",
		},
		"email": {
			"isEmpty": "Enter an email address",
		},
		"phoneNumber": {
			"isEmpty": "Enter a telephone number",
		},
		"isMainContact": {
			"isEmpty": "Select whether this contact is a main contact",
		},
		"isNamedDeputy": {
			"isEmpty": "Select whether this contact is the named deputy",
		},
	}

	assert.Equal(addContactVars{
		Path:   "/133",
		Errors: expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestAddContactFormatValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockContactInformation{}

	validationErrors := sirius.ValidationErrors{
		"contactName": {
			"stringLengthTooLong": "The name must be 255 characters or fewer",
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
		returnedError = renderTemplateForAddContact(client, template)(sirius.DeputyDetails{}, w, r)
	})

	expectedValidationErrors := sirius.ValidationErrors{
		"contactName": {
			"stringLengthTooLong": "The name must be 255 characters or fewer",
		},
	}

	testHandler.ServeHTTP(w, r)

	assert.Equal(addContactVars{
		Path:   "/133",
		Errors: expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}
