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
	count             int
	lastCtx           sirius.Context
	AddContactErr     error
	GetContactByIdErr error
	UpdateContactErr  error
}

func (m *mockManageContact) AddContact(ctx sirius.Context, deputyId int, contact sirius.ContactForm) error {
	m.count += 1
	m.lastCtx = ctx

	return m.AddContactErr
}

func (m *mockManageContact) GetContactById(ctx sirius.Context, deputyId int, contactId int) (sirius.Contact, error) {
	m.count += 1
	m.lastCtx = ctx

	return sirius.Contact{}, m.GetContactByIdErr
}

func (m *mockManageContact) UpdateContact(ctx sirius.Context, deputyId int, contactId int, contact sirius.ContactForm) error {
	m.count += 1
	m.lastCtx = ctx

	return m.UpdateContactErr
}

func TestGetCreateContact(t *testing.T) {
	assert := assert.New(t)

	app := AppVars{PageName: "Add new contact"}
	client := &mockManageContact{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForManageContact(client, template)
	err := handler(app, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(ManageContactVars{
		AppVars:      app,
		IsNewContact: true,
	}, template.lastVars)
}

func TestGetManageContact(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageContact{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	r.SetPathValue("id", "133")
	r.SetPathValue("contactId", "1")

	vars := map[string]string{
		"id":        "133",
		"contactId": "1",
	}

	r = mux.SetURLVars(r, vars)

	app := AppVars{PageName: "Manage contact"}

	handler := renderTemplateForManageContact(client, template)
	err := handler(app, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(ManageContactVars{
		AppVars:                       app,
		IsNamedDeputy:                 "false",
		IsMainContact:                 "false",
		IsMonthlySpreadsheetRecipient: "false",
		IsNewContact:                  false,
		ContactId:                     1,
	}, template.lastVars)
}

func TestPostAddContact(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageContact{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))

	r.SetPathValue("id", "123")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManageContact(client, nil)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(Redirect("/123/contacts?success=newContact"), returnedError)
}

func TestPostManageContact(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageContact{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/1", strings.NewReader(""))
	r.SetPathValue("id", "123")
	r.SetPathValue("contactId", "1")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/{contactId}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManageContact(client, nil)(AppVars{}, w, r)
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

	client.AddContactErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
	r.SetPathValue("id", "133")

	returnedError := renderTemplateForManageContact(client, template)(AppVars{}, w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"contactName": {
			"isEmpty": "Enter a name",
		},
	}

	assert.Equal(ManageContactVars{
		AppVars: AppVars{
			Errors:   expectedValidationErrors,
			PageName: "Add new contact",
		},
		IsNewContact: true,
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestPostAddContactReturnsNonValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageContact{
		AddContactErr: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()

	r, _ := http.NewRequest("POST", "/123/contacts", strings.NewReader(""))
	r.SetPathValue("id", "123")

	manageContactReturnedError := renderTemplateForManageContact(client, template)(AppVars{}, w, r)
	assert.Equal(client.AddContactErr, manageContactReturnedError)
}

func TestPostManageContactReturnsNonValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageContact{
		UpdateContactErr: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	vars := map[string]string{
		"id":        "123",
		"contactId": "1",
	}

	w := httptest.NewRecorder()

	r, _ := http.NewRequest("POST", "/123/contacts/1", strings.NewReader(""))
	r = mux.SetURLVars(r, vars)
	r.SetPathValue("id", "123")
	r.SetPathValue("contactId", "1")

	manageContactReturnedError := renderTemplateForManageContact(client, template)(AppVars{}, w, r)
	assert.Equal(client.UpdateContactErr, manageContactReturnedError)
}

func TestGetManageContactReturnsNonValidationErrors(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageContact{
		GetContactByIdErr: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	vars := map[string]string{
		"id":        "123",
		"contactId": "1",
	}

	w := httptest.NewRecorder()

	r, _ := http.NewRequest("GET", "/123/contacts/1", strings.NewReader(""))
	r = mux.SetURLVars(r, vars)
	r.SetPathValue("id", "123")
	r.SetPathValue("contactId", "1")

	app := AppVars{PageName: "Manage contact"}

	getContactByIdReturnedErr := renderTemplateForManageContact(client, template)(app, w, r)
	assert.Equal(client.GetContactByIdErr, getContactByIdReturnedErr)
}
