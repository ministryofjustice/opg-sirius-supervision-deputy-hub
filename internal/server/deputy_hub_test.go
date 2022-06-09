package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubInformation struct {
	count            int
	lastCtx          sirius.Context
	err              error
	deputyClientData sirius.DeputyClientDetails
	ariaSorting      sirius.AriaSorting
	userDetails      sirius.UserDetails
}

func (m *mockDeputyHubInformation) GetDeputyClients(ctx sirius.Context, deputyId int, deputyType string, columnBeingSorted string, sortOrder string) (sirius.DeputyClientDetails, sirius.AriaSorting, int, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyClientData, m.ariaSorting, 0, m.err
}

func (m *mockDeputyHubInformation) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userDetails, m.err
}

func TestNavigateToDeputyHub(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubInformation{}
	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHub(client, defaultPATeam, template)
	err := handler(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnEcmSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/?success=ecm", "Jon Snow", "defaultPATeam")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Ecm changed to Jon Snow")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnTeamDetailsSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/?success=teamDetails", "Jon Snow", "defaultPATeam")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Team details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnDeputyContactDetailsSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/?success=deputyDetails", "Jon Snow", "defaultPATeam")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Deputy details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilForAnyOtherText(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/?success=otherMessage", "Jon Snow", "defaultPATeam")
	assert.Equal(t, false, Success)
	assert.Equal(t, SuccessMessage, "")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilIfNoSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/", "Jon Snow", "defaultPATeam")
	assert.Equal(t, false, Success)
	assert.Equal(t, SuccessMessage, "")
}

func TestCheckForDefaultEcmIdReturnsMessageIfTrue(t *testing.T) {
	assert.Equal(t, "An executive case manager has not been assigned. ", checkForDefaultEcmId(23, 23))
}

func TestCheckForDefaultEcmIdReturnsNullIfFalse(t *testing.T) {
	assert.Equal(t, "", checkForDefaultEcmId(25, 23))
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnDeputyDetailsSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/?success=deputyDetails", "Jon Snow", "defaultPATeam")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Deputy details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageUseExistingFirmSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/?success=firm", "Jon Snow", "defaultPATeam")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Firm changed to defaultPATeam")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageAddFirmSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/?success=newFirm", "Jon Snow", "defaultPATeam")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Firm added")
}
