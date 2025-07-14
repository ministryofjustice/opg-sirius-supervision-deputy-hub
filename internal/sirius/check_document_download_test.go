package sirius

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckDocumentDownloadReturnsNilOnSuccess(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/supervision-api/v1/documents/123/download", r.URL.Path)
		assert.Equal(t, http.MethodHead, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.CheckDocumentDownload(getContext(nil), 123)
	assert.Nil(t, err)
}

func TestCheckDocumentDownloadReturnsUnauthorizedError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.CheckDocumentDownload(getContext(nil), 123)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestCheckDocumentDownloadReturnsBadRequestError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.CheckDocumentDownload(getContext(nil), 123)
	assert.Equal(t, StatusError{
		Code:   http.StatusBadRequest,
		URL:    svr.URL + SupervisionAPIPath + "/v1/documents/123/download",
		Method: http.MethodHead,
	}, err)
}

func TestCheckDocumentDownloadReturnsStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.CheckDocumentDownload(getContext(nil), 123)
	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + SupervisionAPIPath + "/v1/documents/123/download",
		Method: http.MethodHead,
	}, err)
}
