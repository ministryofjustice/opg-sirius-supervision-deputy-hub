package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockChangeECMInformation struct {
	count      int
	lastCtx    sirius.Context
	err        error
	deputyData sirius.DeputyDetails
}

func (m *mockChangeECMInformation) GetDeputyDetails(ctx sirius.Context, deputyId int) (sirius.DeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func (m *mockChangeECMInformation) ChangeECM(ctx sirius.Context, deputyDetails sirius.DeputyDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestNavigateToChangeECM(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangeECMInformation{}
	template := &mockTemplates{}
	defaultPATeam := "PA"

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHub(client, defaultPATeam, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
