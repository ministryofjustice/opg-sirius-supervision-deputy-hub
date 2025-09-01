package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceDocument(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"date": "14/06/2024",
		"description": "<p>Note content</p>",
		"direction": {
			"handle": "INCOMING",
			"label": "Incoming"
		},
		"name": "Test",
		"type": {
			"handle": "CASE_FORUM",
			"label": "Case forum"
		},
		"personId": 68,
		"fileName": "testfile.png",
		"file": {
			"name": "testfile.png",
			"type": "image/png",
			"source": "VBORw0KGgoAAAANSUhEUgAABg0AAAMOCA",
		},
		"fileSource" : "VBORw0KGgoAAAANSUhEUgAABg0AAAMOCA"
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	tempFile, _ := os.Create("testfile.txt")

	err := client.ReplaceDocument(getContext(nil), tempFile, "file_title.pdf", "Call", "INCOMING", "2020-01-01", "Some notes about my file", 68, 5)
	assert.Nil(t, err)
}

func TestReplaceDocumentReturnsNewStatusError(t *testing.T) {
	tempFile, _ := os.Create("testfile.txt")
	_, _ = tempFile.Write([]byte("test string"))

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.ReplaceDocument(
		getContext(nil),
		tempFile,
		"file_title.pdf",
		"Call",
		"INCOMING",
		"2020-01-01",
		"Some notes about my file",
		68,
		5,
	)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/deputies/68/documents/5",
		Method: http.MethodPut,
	}, err)
}

func TestReplaceDocumentReturnsUnauthorisedClientError(t *testing.T) {
	tempFile, _ := os.Create("testfile.txt")
	_, _ = tempFile.Write([]byte("test string"))

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()
	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.ReplaceDocument(
		getContext(nil),
		tempFile,
		"file_title.pdf",
		"Call",
		"INCOMING",
		"2020-01-01",
		"Some notes about my file",
		68,
		5,
	)
	assert.Equal(t, ErrUnauthorized, err)
}
