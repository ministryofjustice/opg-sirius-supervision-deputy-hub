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

type mockDeputyHubNotesInformation struct {
	count           int
	lastCtx         sirius.Context
	err             error
	addNote         error
	deputyNotesData sirius.DeputyNoteCollection
}

func (m *mockDeputyHubNotesInformation) GetDeputyNotes(ctx sirius.Context, deputyId int) (sirius.DeputyNoteCollection, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyNotesData, m.err
}

func (m *mockDeputyHubNotesInformation) AddNote(ctx sirius.Context, title, note string, deputyId, usedId int, deputyType string) error {
	m.count += 1
	m.lastCtx = ctx

	return m.addNote
}

func TestGetNotes(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubNotesInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?success=true", nil)

	handler := renderTemplateForDeputyHubNotes(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, client.count)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(deputyHubNotesVars{
		SuccessMessage: "Note added",
	}, template.lastVars)
}

func TestPostAddNote(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyHubNotesInformation{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForDeputyHubNotes(client, nil)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/123/notes?success=true"))
}

func TestErrorMessageWhenStringLengthTooLong(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyHubNotesInformation{}

	validationErrors := sirius.ValidationErrors{
		"name": {
			"stringLengthTooLong": "This team type is already in use",
		},
		"description": {
			"stringLengthTooLong": "This team type is already in use",
		},
	}
	client.addNote = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForDeputyHubNotes(client, template)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addNoteVars{
		AppVars: AppVars{
			Errors: validationErrors,
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestErrorMessageWhenIsEmpty(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeputyHubNotesInformation{}

	validationErrors := sirius.ValidationErrors{
		"name": {
			"isEmpty": "This team type is already in use",
		},
		"description": {
			"isEmpty": "This team type is already in use",
		},
	}
	client.addNote = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForDeputyHubNotes(client, template)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addNoteVars{
		AppVars: AppVars{
			Errors: validationErrors,
		},
	}, template.lastVars)

	assert.Nil(returnedError)
}
