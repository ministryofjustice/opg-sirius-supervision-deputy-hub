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

func TestDeputyDetailsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `    {
      	"id": 1,
      	"deputyCasrecId": 10000000,
      	"organisationName": "Test Organisation",
		"email": "deputyship@essexcounty.gov.uk",
		"phoneNumber": "0115 876 5574",
		"addressLine1": "Deputyship Team",
		"addressLine2": "Seax House",
		"addressLine3": "19 Market Rd",
		"town": "Chelmsford",
		"county": "Essex",
		"postcode": "CM1 1GG"
    }`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := DeputyDetails{
		ID:               1,
		DeputyCasrecId:   10000000,
		OrganisationName: "Test Organisation",
		Email:            "deputyship@essexcounty.gov.uk",
		PhoneNumber:      "0115 876 5574",
		AddressLine1:     "Deputyship Team",
		AddressLine2:     "Seax House",
		AddressLine3:     "19 Market Rd",
		Town:             "Chelmsford",
		County:           "Essex",
		Postcode:         "CM1 1GG",
	}

	deputyDetails, err := client.GetDeputyDetails(getContext(nil), "", 1)

	assert.Equal(t, expectedResponse, deputyDetails)
	assert.Equal(t, nil, err)
}

func TestGetDeputyDetailsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyDetails, err := client.GetDeputyDetails(getContext(nil), "", 1)

	expectedResponse := DeputyDetails{
		ID:               0,
		DeputyCasrecId:   0,
		OrganisationName: "",
	}

	assert.Equal(t, expectedResponse, deputyDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/1",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyDetailsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyDetails, err := client.GetDeputyDetails(getContext(nil), "", 1)

	expectedResponse := DeputyDetails{
		ID:               0,
		DeputyCasrecId:   0,
		OrganisationName: "",
	}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyDetails)
}
