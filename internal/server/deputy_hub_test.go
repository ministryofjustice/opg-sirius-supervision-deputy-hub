package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubInformation struct {
	count      int
	lastCtx    sirius.Context
	err        error
	deputyData sirius.DeputyDetails
}

func (m *mockDeputyHubInformation) GetDeputyDetails(ctx sirius.Context, defaultPATeam int, deputyId int) (sirius.DeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func TestNavigateToDeputyHub(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubInformation{}
	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHub(client, defaultPATeam, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/?success=ecm", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Ecm changed to Jon Snow")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilIfNoSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/", "Jon Snow")
	assert.Equal(t, false, Success)
	assert.Equal(t, SuccessMessage, "")
}

func TestCheckForDefaultEcmIdReturnsMessageIfTrue(t *testing.T) {
	assert.Equal(t, "An executive case manager has not been assigned. ", checkForDefaultEcmId(23, 23))
}

func TestCheckForDefaultEcmIdReturnsNullIfFalse(t *testing.T) {
	assert.Equal(t, "", checkForDefaultEcmId(25, 23))
}