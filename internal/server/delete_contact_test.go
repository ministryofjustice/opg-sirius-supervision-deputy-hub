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

type mockDeleteContact struct {
	count               int
	lastCtx             sirius.Context
	getContactByIdError error
	deleteContactError  error
}

func (m *mockDeleteContact) GetContactById(ctx sirius.Context, deputyId int, contactId int) (sirius.Contact, error) {
	m.count += 1
	m.lastCtx = ctx

	return sirius.Contact{}, m.getContactByIdError
}

func (m *mockDeleteContact) DeleteContact(ctx sirius.Context, deputyId int, contactId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.deleteContactError
}

func TestGetDeleteContact(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteContact{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeleteContact(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostDeleteContact(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeleteContact{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123/1", strings.NewReader(""))
	r.SetPathValue("id", "123")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/{contactId}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForDeleteContact(client, nil)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(Redirect("/123/contacts?success=deletedContact&contactName="), returnedError)
}
