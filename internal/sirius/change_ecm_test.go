package sirius

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestChangeECM(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"ID":                               32,
    "OrganisationName":                 "deputyDetails.OrganisationName",
    "OrganisationTeamOrDepartmentName": "r.PostFormValue("new-ecm")",
    "Email":                            "deputyDetails.Email",
    "PhoneNumber":                      "deputyDetails.PhoneNumber",
    "AddressLine1":                     "deputyDetails.AddressLine1",
    "AddressLine2":                     "deputyDetails.AddressLine2",
    "AddressLine3":                     "deputyDetails.AddressLine3",
    "Town":                             "deputyDetails.Town",
    "County":                           "deputyDetails.County",
    "Postcode":                         "deputyDetails.Postcode",
	}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	changeECMForm := DeputyDetails{
		ID:                               32,
		OrganisationName:                 "deputy-name",
		OrganisationTeamOrDepartmentName: "organisationTeamOrDepartmentName",
		Email:                            "email",
		PhoneNumber:                      "telephone",
		AddressLine1:                     "address-line-1",
		AddressLine2:                     "address-line-2",
		AddressLine3:                     "address-line-3",
		Town:                             "town",
		County:                           "county",
		Postcode:                         "postcode",
	}

	err := client.EditDeputyDetails(getContext(nil), changeECMForm)
	assert.Nil(t, err)
}

func TestChangeECMReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	changeECMForm := DeputyDetails{
		ID:                               32,
		OrganisationName:                 "deputy-name",
		OrganisationTeamOrDepartmentName: "organisationTeamOrDepartmentName",
		Email:                            "email",
		PhoneNumber:                      "telephone",
		AddressLine1:                     "address-line-1",
		AddressLine2:                     "address-line-2",
		AddressLine3:                     "address-line-3",
		Town:                             "town",
		County:                           "county",
		Postcode:                         "postcode",
	}

	err := client.EditDeputyDetails(getContext(nil), changeECMForm)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/32",
		Method: http.MethodPut,
	}, err)
}

func TestChangeECMReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	changeECMForm := DeputyDetails{
		ID:                               32,
		OrganisationName:                 "deputy-name",
		OrganisationTeamOrDepartmentName: "organisationTeamOrDepartmentName",
		Email:                            "email",
		PhoneNumber:                      "telephone",
		AddressLine1:                     "address-line-1",
		AddressLine2:                     "address-line-2",
		AddressLine3:                     "address-line-3",
		Town:                             "town",
		County:                           "county",
		Postcode:                         "postcode",
	}

	err := client.EditDeputyDetails(getContext(nil), changeECMForm)

	assert.Equal(t, ErrUnauthorized, err)
}
