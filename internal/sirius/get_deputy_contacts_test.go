package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeputyContactReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := ` [
      {
		"name": "Test Contact",
        "jobTitle": "Software Tester",
        "email": "test.contact@email.com",
        "phoneNumber": "0123456789",
        "otherPhoneNumber": "9876543210",
        "isMainContact": true,
		"isNamedDeputy": false
      }
  	] `

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedContacts := ContactList{
		{
			Name:             "Test Contact",
			JobTitle:         "Software Tester",
			Email:            "test.contact@email.com",
			PhoneNumber:      "0123456789",
			OtherPhoneNumber: "9876543210",
			IsMainContact:    true,
			IsNamedDeputy:    false,
		},
	}

	contacts, err := client.GetDeputyContacts(getContext(nil), 1)

	assert.Equal(t, expectedContacts, contacts)
	assert.Equal(t, nil, err)
}

func TestGetDeputyContactReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	contacts, err := client.GetDeputyContacts(getContext(nil), 1)

	expectedResponse := ContactList(nil)

	assert.Equal(t, expectedResponse, contacts)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/1/contacts",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyContactsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	contacts, err := client.GetDeputyContacts(getContext(nil), 1)

	expectedResponse := ContactList(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, contacts)
}
