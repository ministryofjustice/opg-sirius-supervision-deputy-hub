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

type mockManageContact struct {
	count   int
	lastCtx sirius.Context
	err     error
}

func (m *mockManageContact) AddContact(ctx sirius.Context, deputyId int, contact sirius.ContactForm) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func (m *mockManageContact) GetContactById(ctx sirius.Context, deputyId int, contactId int) (sirius.Contact, error) {
	m.count += 1
	m.lastCtx = ctx

	return sirius.Contact{}, m.err
}

func (m *mockManageContact) UpdateContact(ctx sirius.Context, deputyId int, contactId int, contact sirius.ContactForm) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestGetCreateContact(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageContact{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForManageContact(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(ManageContactVars{
		Path:         "/path",
		IsNewContact: true,
	}, template.lastVars)
}

func TestGetManageContact(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageContact{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	vars := map[string]string{
		"id":        "133",
		"contactId": "1",
	}

	r = mux.SetURLVars(r, vars)

	handler := renderTemplateForManageContact(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(ManageContactVars{
		Path:          "/path",
		IsNamedDeputy: "false",
		IsMainContact: "false",
		IsNewContact:  false,
	}, template.lastVars)
}

func TestPostAddContact(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageContact{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManageContact(client, nil)(sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(Redirect("/123/contacts?success=newContact"), returnedError)
}

func TestPostManageContact(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageContact{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/1", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/{contactId}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManageContact(client, nil)(sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(Redirect("/123/contacts?success=updatedContact&contactName="), returnedError)
}

func TestAddContactEmptyValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageContact{}

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
		returnedError = renderTemplateForManageContact(client, template)(sirius.DeputyDetails{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"contactName": {
			"isEmpty": "Enter a name",
		},
	}

	assert.Equal(ManageContactVars{
		Path:   "/133",
		Errors: expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}
