package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockCheckDocumentDownloadClient struct {
	count          int
	lastCtx        sirius.Context
	lastDocumentId int
	err            error
}

func (m *mockCheckDocumentDownloadClient) CheckDocumentDownload(ctx sirius.Context, documentId int) error {
	m.count += 1
	m.lastCtx = ctx
	m.lastDocumentId = documentId
	return m.err
}

func TestCheckDocumentDownloadWhenSuccessful(t *testing.T) {
	assert := assert.New(t)

	client := &mockCheckDocumentDownloadClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/456/documents/123/check", nil)
	r.SetPathValue("id", "456")
	r.SetPathValue("documentId", "123")

	handler := checkDocument(client)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)
	assert.Equal(1, client.count)
	assert.Equal(123, client.lastDocumentId)
	assert.Equal(http.StatusOK, w.Code)
}

func TestCheckDocumentDownloadWhenMethodNotAllowed(t *testing.T) {
	assert := assert.New(t)

	client := &mockCheckDocumentDownloadClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/456/documents/123/check", nil)
	r.SetPathValue("id", "456")
	r.SetPathValue("documentId", "123")

	handler := checkDocument(client)
	err := handler(AppVars{}, w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
	assert.Equal(0, client.count)
}

func TestCheckDocumentDownloadReturnsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := sirius.StatusError{Code: http.StatusInternalServerError}
	client := &mockCheckDocumentDownloadClient{err: expectedError}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/456/documents/123/check", nil)
	r.SetPathValue("id", "456")
	r.SetPathValue("documentId", "123")

	handler := checkDocument(client)
	err := handler(AppVars{}, w, r)

	assert.Equal(expectedError, err)
	assert.Equal(1, client.count)
	assert.Equal(123, client.lastDocumentId)
}
