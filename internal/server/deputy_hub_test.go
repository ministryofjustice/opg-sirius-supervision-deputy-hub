package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubInformation struct {
	count               int
	lastCtx             sirius.Context
	GetDeputyClientsErr error
	GetUserDetailsErr   error
	deputyClientData    sirius.ClientList
	ariaSorting         sirius.AriaSorting
	userDetails         sirius.UserDetails
}

func (m *mockDeputyHubInformation) GetDeputyClients(ctx sirius.Context, deputyId, displayClientLimit, search int, deputyType, columnBeingSorted, sortOrder string) (sirius.ClientList, sirius.AriaSorting, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyClientData, m.ariaSorting, m.GetDeputyClientsErr
}

func (m *mockDeputyHubInformation) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userDetails, m.GetUserDetailsErr
}

func TestNavigateToDeputyHub(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHub(client, template)
	err := handler(sirius.DeputyDetails{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnEcmSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=ecm")
	SuccessMessage := getSuccessFromUrl(u, "Jon Snow", "defaultPATeam")
	assert.Equal(t, SuccessMessage, "<abbr title='Executive Case Manager'>ECM</abbr> changed to Jon Snow")
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

func TestDeputyHubHandlesErrorsInOtherClientFiles(t *testing.T) {
	returnedError := sirius.StatusError{Code: 500}
	tests := []struct {
		Client *mockDeputyHubInformation
	}{
		{
			Client: &mockDeputyHubInformation{
				GetDeputyClientsErr: returnedError,
			},
		},
		{
			Client: &mockDeputyHubInformation{
				GetUserDetailsErr: returnedError,
			},
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k), func(t *testing.T) {

			client := tc.Client
			template := &mockTemplates{}
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/123", strings.NewReader(""))
			deputyHubReturnedError := renderTemplateForDeputyHub(client, template)(sirius.DeputyDetails{}, w, r)
			assert.Equal(t, returnedError, deputyHubReturnedError)
		})
	}
}
