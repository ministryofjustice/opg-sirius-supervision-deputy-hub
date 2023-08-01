package sirius

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAddContact(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"contactName":"Contact Name",
		"email":"Email_address@address.com",
		"phoneNumber":"11111111"
		}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}
  
	contact := ContactForm{
		ContactName: "Contact Name",
		Email:       "Email_address@address.com",
		PhoneNumber: "11111111",
	}

	deputyId := 76

	err := client.AddContact(getContext(nil), deputyId, contact)
	assert.Nil(t, err)
}

func TestAddContactReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyId := 76

	err := client.AddContact(getContext(nil), deputyId, ContactForm{})

	url := fmt.Sprintf("/api/v1/deputies/%d/contacts", deputyId)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + url,
		Method: http.MethodPost,
	}, err)
}

func TestAddContactReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyId := 76

	err := client.AddContact(getContext(nil), deputyId, ContactForm{})

	assert.Equal(t, ErrUnauthorized, err)

}