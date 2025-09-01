package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteContact(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	r := io.NopCloser(bytes.NewReader([]byte("")))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 204,
			Body:       r,
		}, nil
	}

	err := client.DeleteContact(getContext(nil), 76, 1)
	assert.Nil(t, err)
}

func TestDeleteContactReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.DeleteContact(getContext(nil), 76, 1)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/deputies/76/contacts/1",
		Method: http.MethodDelete,
	}, err)
}

func TestDeleteContactReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.DeleteContact(getContext(nil), 76, 1)

	assert.Equal(t, ErrUnauthorized, err)
}
