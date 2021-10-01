package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubNotesInformation struct {
	count            int
	lastCtx          sirius.Context
	err              error
	addNote			 error
	deputyData       sirius.DeputyDetails
	deputyNotesData sirius.DeputyNoteList
	userDetailsData sirius.UserDetails
}

func (m *mockDeputyHubNotesInformation) GetDeputyDetails(ctx sirius.Context, deputyId int) (sirius.DeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func (m *mockDeputyHubNotesInformation) GetDeputyNotes(ctx sirius.Context, deputyId int) (sirius.DeputyNoteList, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyNotesData, m.err
}

func (m *mockDeputyHubNotesInformation) AddNote(ctx sirius.Context, title, note string, deputyId, usedId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.addNote
}

func (m *mockDeputyHubNotesInformation) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userDetailsData, m.err
}

func TestGetNotes(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubNotesInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHubNotes(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(2, client.count)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(deputyHubNotesVars{
		Path:      "/path",
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

	testHandler := mux.NewRouter();
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForDeputyHubNotes(client, nil)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, RedirectError("/deputy/123/notes?success=true"))
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

	testHandler := mux.NewRouter();
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForDeputyHubNotes(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"1-title": {
			"stringLengthTooLong": "The title must be 255 characters or fewer",
		},
		"2-note": {
			"stringLengthTooLong": "The note must be 1000 characters or fewer",
		},
	}

	assert.Equal(3, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addNoteVars{
		Path:    "/123",
		Errors:  expectedValidationErrors,
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

	testHandler := mux.NewRouter();
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForDeputyHubNotes(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"1-title": {
			"isEmpty": "Enter a title for the note",
		},
		"2-note": {
			"isEmpty": "Enter a note",
		},
	}

	assert.Equal(3, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addNoteVars{
		Path:    "/123",
		Errors:  expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}
