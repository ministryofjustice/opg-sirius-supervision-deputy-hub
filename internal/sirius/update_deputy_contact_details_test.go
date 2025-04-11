package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

const ID = 32

func TestUpdateDeputyContactDetails(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
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

	formData := DeputyContactDetails{
		Email:        "email",
		PhoneNumber:  "telephone",
		AddressLine1: "address-line-1",
		AddressLine2: "address-line-2",
		AddressLine3: "address-line-3",
		Town:         "town",
		County:       "county",
		Postcode:     "postcode",
	}

	err := client.UpdateDeputyContactDetails(getContext(nil), ID, formData)
	assert.Nil(t, err)
}

func TestUpdateDeputyContactDetailsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := DeputyContactDetails{
		Email:        "email",
		PhoneNumber:  "telephone",
		AddressLine1: "address-line-1",
		AddressLine2: "address-line-2",
		AddressLine3: "address-line-3",
		Town:         "town",
		County:       "county",
		Postcode:     "postcode",
	}

	err := client.UpdateDeputyContactDetails(getContext(nil), ID, formData)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    fmt.Sprintf("%v/v1/deputies/%d/contact-details", svr.URL + SupervisionAPIPath, ID),
		Method: http.MethodPut,
	}, err)
}

func TestUpdateDeputyContactDetailsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := DeputyContactDetails{
		Email:        "email",
		PhoneNumber:  "telephone",
		AddressLine1: "address-line-1",
		AddressLine2: "address-line-2",
		AddressLine3: "address-line-3",
		Town:         "town",
		County:       "county",
		Postcode:     "postcode",
	}

	err := client.UpdateDeputyContactDetails(getContext(nil), ID, formData)

	assert.Equal(t, ErrUnauthorized, err)
}
