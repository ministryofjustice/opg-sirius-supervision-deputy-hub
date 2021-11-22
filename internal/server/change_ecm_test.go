package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockChangeECMInformation struct {
	count          int
	lastCtx        sirius.Context
	err            error
	DeputyDetails  sirius.DeputyDetails
	EcmTeamDetails []sirius.TeamMember
	Error          string
	Errors         sirius.ValidationErrors
	Success        bool
	DefaultPaTeam  int
}

func (m *mockChangeECMInformation) GetDeputyDetails(ctx sirius.Context, defaultPATeam int, deputyId int) (sirius.DeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.DeputyDetails, m.err
}

func (m *mockChangeECMInformation) GetPaDeputyTeamMembers(ctx sirius.Context, deputyId int) ([]sirius.TeamMember, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.EcmTeamDetails, m.err
}

func (m *mockChangeECMInformation) ChangeECM(ctx sirius.Context, changeEcmForm sirius.ExecutiveCaseManagerOutgoing, deputyDetails sirius.DeputyDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestGetChangeECM(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangeECMInformation{}
	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForChangeECM(client, defaultPATeam, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(2, client.count)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(changeECMHubVars{
		Path:          "/path",
		DefaultPaTeam: 23,
		SuccessMessage: "new ecm is",
	}, template.lastVars)
}

//func TestPostChangeECM(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockChangeECMInformation{}
//	defaultPATeam := 23
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/76/ecm", strings.NewReader({"Ecm":23}))
//	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}/ecm", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = renderTemplateForChangeECM(client, defaultPATeam, nil)(sirius.PermissionSet{}, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//	assert.Equal(returnedError, Redirect("/deputies/76"))
//}
//
//func TestErrorMessageWhenStringLengthTooLong(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockDeputyHubNotesInformation{}
//
//	validationErrors := sirius.ValidationErrors{
//		"name": {
//			"stringLengthTooLong": "This team type is already in use",
//		},
//		"description": {
//			"stringLengthTooLong": "This team type is already in use",
//		},
//	}
//	client.addNote = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//	defaultPATeam := 23
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
//	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = renderTemplateForDeputyHubNotes(client, defaultPATeam, template)(sirius.PermissionSet{}, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//
//	expectedValidationErrors := sirius.ValidationErrors{
//		"1-title": {
//			"stringLengthTooLong": "The title must be 255 characters or fewer",
//		},
//		"2-note": {
//			"stringLengthTooLong": "The note must be 1000 characters or fewer",
//		},
//	}
//
//	assert.Equal(3, client.count)
//
//	assert.Equal(1, template.count)
//	assert.Equal("page", template.lastName)
//	assert.Equal(addNoteVars{
//		Path:   "/123",
//		Errors: expectedValidationErrors,
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//}
//
//func TestErrorMessageWhenIsEmpty(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockDeputyHubNotesInformation{}
//
//	validationErrors := sirius.ValidationErrors{
//		"name": {
//			"isEmpty": "This team type is already in use",
//		},
//		"description": {
//			"isEmpty": "This team type is already in use",
//		},
//	}
//	client.addNote = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//	defaultPATeam := 23
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
//	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = renderTemplateForDeputyHubNotes(client, defaultPATeam, template)(sirius.PermissionSet{}, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//
//	expectedValidationErrors := sirius.ValidationErrors{
//		"1-title": {
//			"isEmpty": "Enter a title for the note",
//		},
//		"2-note": {
//			"isEmpty": "Enter a note",
//		},
//	}
//
//	assert.Equal(3, client.count)
//
//	assert.Equal(1, template.count)
//	assert.Equal("page", template.lastName)
//	assert.Equal(addNoteVars{
//		Path:   "/123",
//		Errors: expectedValidationErrors,
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//}