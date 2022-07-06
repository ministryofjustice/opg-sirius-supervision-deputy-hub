package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHub(client, template)
	err := handler(sirius.PermissionSet{}, sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnEcmSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=ecm")
	SuccessMessage := getSuccessFromUrl(u, "Jon Snow", "defaultPATeam")
	assert.Equal(t, SuccessMessage, "Ecm changed to Jon Snow")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnTeamDetailsSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=teamDetails")
	SuccessMessage := getSuccessFromUrl(u, "Jon Snow", "defaultPATeam")
	assert.Equal(t, SuccessMessage, "Team details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnDeputyContactDetailsSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=deputyDetails")
	SuccessMessage := getSuccessFromUrl(u, "Jon Snow", "defaultPATeam")
	assert.Equal(t, SuccessMessage, "Deputy details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilForAnyOtherText(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=otherMessage")
	SuccessMessage := getSuccessFromUrl(u, "Jon Snow", "defaultPATeam")
	assert.Equal(t, SuccessMessage, "")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilIfNoSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/")
	SuccessMessage := getSuccessFromUrl(u, "Jon Snow", "defaultPATeam")
	assert.Equal(t, SuccessMessage, "")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnDeputyDetailsSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=deputyDetails")
	SuccessMessage := getSuccessFromUrl(u, "Jon Snow", "defaultPATeam")
	assert.Equal(t, SuccessMessage, "Deputy details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageUseExistingFirmSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=firm")
	SuccessMessage := getSuccessFromUrl(u, "Jon Snow", "defaultPATeam")
	assert.Equal(t, SuccessMessage, "Firm changed to defaultPATeam")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageAddFirmSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/deputy/76/?success=newFirm")
	SuccessMessage := getSuccessFromUrl(u, "Jon Snow", "defaultPATeam")
	assert.Equal(t, SuccessMessage, "Firm added")
}
