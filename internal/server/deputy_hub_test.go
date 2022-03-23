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
	firms      []sirius.Firm
}

func (m *mockDeputyHubInformation) GetFirms(ctx sirius.Context) ([]sirius.Firm, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.firms, m.err
}

func (m *mockDeputyHubInformation) AssignDeputyToFirm(ctx sirius.Context, i int, i2 int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func (m *mockDeputyHubInformation) GetDeputyDetails(ctx sirius.Context, defaultPATeam int, deputyId int) (sirius.DeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func TestNavigateToDeputyHub(t *testing.T) {
	client := &mockDeputyHubInformation{}
	template := &mockTemplates{}
	defaultPATeam := 23

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHub(client, defaultPATeam, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(t, err)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageAddFirmSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/?success=newFirm", "Firm Name", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Firm added")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnChangeFirmSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/?success=firm", "firm Name", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Firm changed to firm Name")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnDeputyContactDetailsSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/?success=deputyDetails", "Firm Name", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Deputy details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnEcmSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/?success=ecm", "Firm Name", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Ecm changed to Jon Snow")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnTeamDetailsSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/?success=teamDetails", "Firm Name", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Team details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilForAnyOtherText(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/?success=otherMessage", "Firm Name", "Jon Snow")
	assert.Equal(t, false, Success)
	assert.Equal(t, SuccessMessage, "")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilIfNoSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/", "Firm Name", "Jon Snow")
	assert.Equal(t, false, Success)
	assert.Equal(t, SuccessMessage, "")
}

func TestCheckForDefaultEcmIdReturnsMessageIfTrue(t *testing.T) {
	assert.Equal(t, "An executive case manager has not been assigned. ", checkForDefaultEcmId(23, 23))
}

func TestCheckForDefaultEcmIdReturnsNullIfFalse(t *testing.T) {
	assert.Equal(t, "Firm Name", checkForDefaultEcmId(25, 23))
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnDeputyDetailsSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/76/?success=deputyDetails", "Firm", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Deputy details updated")
}
