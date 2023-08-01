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

func TestContactReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
                "id": 2,
                "name": "Test Contact",
                "phoneNumber": "0123456789",
                "email": "test@email.com",
                "isNamedDeputy": false,
                "isMainContact": false
			}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := Contact{
		ContactName:   "Test Contact",
		PhoneNumber:   "0123456789",
		Email:         "test@email.com",
		IsNamedDeputy: false,
		IsMainContact: false,
	}

	contact, err := client.GetContactById(getContext(nil), 76, 2)

	assert.Equal(t, expectedResponse, contact)
	assert.Equal(t, nil, err)
}

func TestGetContactReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	contact, err := client.GetContactById(getContext(nil), 76, 1)

	expectedResponse := Contact{}

	assert.Equal(t, expectedResponse, contact)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/76/contacts/1",
		Method: http.MethodGet,
	}, err)
}

func TestGetContactReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	contact, err := client.GetContactById(getContext(nil), 76, 1)

	expectedResponse := Contact{}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, contact)
}
