package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeputyDetailsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `    {
      "id": 76,
      "deputyCasrecId": 10000000,
      "organisationName": "Test Organisation"
    }`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := DeputyDetails{
		ID: 76,
		DeputyCasrecId: 10000000,
		OrganisationName: "Test Organisation",
	}

	deputyDetails, err := client.GetDeputyDetails(getContext(nil), 76)

	assert.Equal(t, expectedResponse, deputyDetails)
	assert.Equal(t, nil, err)
}

func TestGetDeputyDetailsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyDetails, err := client.GetDeputyDetails(getContext(nil), 76)

	expectedResponse := DeputyDetails{
		ID: 0,
		DeputyCasrecId: 0,
		OrganisationName: "",
	}

	assert.Equal(t, expectedResponse, deputyDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/76",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyDetailsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyDetails, err := client.GetDeputyDetails(getContext(nil), 76)

	expectedResponse := DeputyDetails{
		ID: 0,
		DeputyCasrecId: 0,
		OrganisationName: "",
	}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyDetails)
}