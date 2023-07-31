package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubContactInformation struct {
	count             int
	lastCtx           sirius.Context
	err               error
	deputyContactData sirius.ContactList
}

func (m *mockDeputyHubContactInformation) GetDeputyContacts(ctx sirius.Context, deputyId int) (sirius.ContactList, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyContactData, m.err
}

func TestNavigateToContactTab(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubContactInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForContactTab(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
