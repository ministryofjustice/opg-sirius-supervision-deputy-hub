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

func TestCloseGCMIssues(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{}`
	r := io.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 204,
			Body:       r,
		}, nil
	}

	err := client.CloseGCMIssues(getContext(nil), []string{"23"})
	assert.Equal(t, nil, err)
}

func TestCloseGCMIssuesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	err := client.CloseGCMIssues(getContext(nil), []string{"23"})

	assert.Equal(t, ErrUnauthorized, err)
}
