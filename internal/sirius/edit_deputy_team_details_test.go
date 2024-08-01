package sirius

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestEditDeputyTeamDetails(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
	"ID":                               32,
    "OrganisationName":                 "r.PostFormValue("deputy-name")",
    "OrganisationTeamOrDepartmentName": "r.PostFormValue("organisationTeamOrDepartmentName")",
    "Email":                            "r.PostFormValue("email")",
    "PhoneNumber":                      "r.PostFormValue("telephone")",
    "AddressLine1":                     "r.PostFormValue("address-line-1")",
    "AddressLine2":                     "r.PostFormValue("address-line-2")",
    "AddressLine3":                     "r.PostFormValue("address-line-3")",
    "Town":                             "r.PostFormValue("town")",
    "County":                           "r.PostFormValue("county")",
    "Postcode":                         "r.PostFormValue("postcode")",
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	editDeputyDetailForm := DeputyDetails{
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

	err := client.EditDeputyTeamDetails(getContext(nil), editDeputyDetailForm)
	assert.Nil(t, err)
}

func TestEditDeputyTeamDetailsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	editDeputyDetailForm := DeputyDetails{
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

	err := client.EditDeputyTeamDetails(getContext(nil), editDeputyDetailForm)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/32",
		Method: http.MethodPut,
	}, err)
}

func TestEditDeputyTeamDetailsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	editDeputyDetailForm := DeputyDetails{
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

	err := client.EditDeputyTeamDetails(getContext(nil), editDeputyDetailForm)

	assert.Equal(t, ErrUnauthorized, err)
}
