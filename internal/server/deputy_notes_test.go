package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubNotesInformation struct {
	count            int
	lastCtx          sirius.Context
	err              error
	deputyData       sirius.DeputyDetails
	deputyNotesData sirius.DeputyNoteList
	addNote error
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

	return m.err
}

func (m *mockDeputyHubNotesInformation) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userDetailsData, m.err
}

func TestRenameValidationErrorMessages(t *testing.T) {
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
	}, template.lastVars)
}

//func TestPostAddNote(t *testing.T) {
//	assert := assert.New(t)
//
//	client := &mockDeputyHubNotesInformation{
//		deputyData: sirius.DeputyDetails{
//			ID: 123,
//		},
//	}
//	client.deputyData.ID = 123
//	client.userDetailsData.ID = 50
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/api/v1/deputy/123/create-note", strings.NewReader(""))
//	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//
//	handler := renderTemplateForDeputyHubNotes(client, template)
//	err := handler(sirius.PermissionSet{}, w, r)
//	assert.Equal(RedirectError("/deputy/123/notes"), err)
//
//	assert.Equal(2, client.count)
//
//	assert.Equal(getContext(r), client.lastCtx)
//
//	assert.Equal(1, template.count)
//	assert.Equal("page", template.lastName)
//	assert.Equal(addNoteVars{
//		Path:   "/path",
//		Title:	"Title",
//		Note:	"Note",
//	}, template.lastVars)
//}
