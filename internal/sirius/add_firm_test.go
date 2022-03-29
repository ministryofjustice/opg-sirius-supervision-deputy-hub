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

func TestAddFirm(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"firmName":"The Firm Name",
		"addressLine1":"Address 1",
		"addressLine2":"Address 2",
		"addressLine3":"Address 3",
		"town":"City",
		"county":"Country",
		"postcode":"ff11bc",
		"email":"Email_address@address.com",
		"phoneNumber":"11111111"
		}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	firmDetails := FirmDetails{
		FirmName:     "The Firm Name",
		AddressLine1: "Address 1",
		AddressLine2: "Address 2",
		AddressLine3: "Address 3",
		Town:         "City",
		County:       "Country",
		Postcode:     "ff11bc",
		Email:        "Email_address@address.com",
		PhoneNumber:  "11111111",
	}

	int, err := client.AddFirmDetails(getContext(nil), firmDetails)
	assert.Nil(t, err)
	assert.Equal(t, 0, int)
}

func TestAddFirmReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	int, err := client.AddFirmDetails(getContext(nil), FirmDetails{})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/firms",
		Method: http.MethodPost,
	}, err)

	assert.Equal(t, 0, int)
}

func TestAddFirmReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	int, err := client.AddFirmDetails(getContext(nil), FirmDetails{})

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, 0, int)

}
