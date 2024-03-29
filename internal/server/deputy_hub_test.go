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
	deputyClientData    sirius.ClientList
}

func (m *mockDeputyHubInformation) GetDeputyClients(ctx sirius.Context, params sirius.ClientListParams) (sirius.ClientList, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyClientData, m.GetDeputyClientsErr
}

var testDeputy = sirius.DeputyDetails{
	ID: 123,
	ExecutiveCaseManager: sirius.ExecutiveCaseManager{
		EcmId:   1,
		EcmName: "Jon Snow",
	},
	Firm: sirius.Firm{
		FirmName: "defaultPATeam",
	},
}

func TestNavigateToDeputyHub(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHub(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnEcmSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=ecm")
	SuccessMessage := getSuccessFromUrl(u, testDeputy)
	assert.Equal(t, SuccessMessage, "<abbr title='Executive Case Manager'>ECM</abbr> changed to Jon Snow")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnTeamDetailsSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=teamDetails")
	SuccessMessage := getSuccessFromUrl(u, testDeputy)
	assert.Equal(t, SuccessMessage, "Team details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnDeputyContactDetailsSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=deputyDetails")
	SuccessMessage := getSuccessFromUrl(u, testDeputy)
	assert.Equal(t, SuccessMessage, "Deputy details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilForAnyOtherText(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=otherMessage")
	SuccessMessage := getSuccessFromUrl(u, testDeputy)
	assert.Equal(t, SuccessMessage, "")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilIfNoSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/")
	SuccessMessage := getSuccessFromUrl(u, testDeputy)
	assert.Equal(t, SuccessMessage, "")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnDeputyDetailsSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=deputyDetails")
	SuccessMessage := getSuccessFromUrl(u, testDeputy)
	assert.Equal(t, SuccessMessage, "Deputy details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageUseExistingFirmSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/76/?success=firm")
	SuccessMessage := getSuccessFromUrl(u, testDeputy)
	assert.Equal(t, SuccessMessage, "Firm changed to defaultPATeam")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageAddFirmSuccess(t *testing.T) {
	u, _ := url.Parse("http::deputyhub/deputy/76/?success=newFirm")
	SuccessMessage := getSuccessFromUrl(u, testDeputy)
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
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {

			client := tc.Client
			template := &mockTemplates{}
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/123", strings.NewReader(""))
			deputyHubReturnedError := renderTemplateForDeputyHub(client, template)(AppVars{}, w, r)
			assert.Equal(t, returnedError, deputyHubReturnedError)
		})
	}
}
