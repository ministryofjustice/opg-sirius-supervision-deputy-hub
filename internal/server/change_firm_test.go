package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyChangeFirmInformation struct {
	count    int
	lastCtx  sirius.Context
	err      error
	firmData []sirius.FirmForList
}

func (m *mockDeputyChangeFirmInformation) GetFirms(ctx sirius.Context) ([]sirius.FirmForList, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.firmData, m.err
}

func (m *mockDeputyChangeFirmInformation) AssignDeputyToFirm(ctx sirius.Context, defaultPATeam int, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestNavigateToChangeFirm(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyChangeFirmInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForChangeFirm(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
