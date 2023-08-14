package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestDeleteContactReturnsError(t *testing.T) {
	tests := []struct {
		statusCode    int
		isClientError bool
		expectedErr   error
	}{
		{
			http.StatusMethodNotAllowed,
			false,
			nil,
		},
		{
			http.StatusUnauthorized,
			true,
			ErrUnauthorized,
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			}))

			defer svr.Close()
			client, _ := NewClient(http.DefaultClient, svr.URL)

			err := client.DeleteContact(getContext(nil), 76, 1)
			if tc.isClientError {
				assert.Equal(t, err, tc.expectedErr)
			} else {
				assert.Equal(t, err, StatusError{
					Code:   tc.statusCode,
					URL:    svr.URL + "/api/v1/deputies/76/contacts/1",
					Method: "DELETE"})
			}
		})
	}
}
