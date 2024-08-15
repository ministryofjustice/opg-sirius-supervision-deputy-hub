package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGcmIssue(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
	  "caseRecNumber": "123456",
	  "gcmIssueType": {
		"handle": "MISSING_INFORMATION",
		"label": "Missing information"
		"deprecated": "false",
	  },
      "notes": "test note"
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	err := client.AddGcmIssue(getContext(nil),
		"123456",
		"notes",
		model.RefData{
			Handle:     "MISSING_INFORMATION",
			Label:      "Missing information",
			Deprecated: false,
		},
		76,
	)
	assert.Nil(t, err)
}

func TestAddGcmIssueReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.AddGcmIssue(getContext(nil),
		"123456",
		"notes",
		model.RefData{
			Handle:     "MISSING_INFORMATION",
			Label:      "Missing information",
			Deprecated: false,
		},
		76,
	)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/76/case-manager-issues",
		Method: http.MethodPost,
	}, err)
}

func TestAddGcmIssueReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.AddGcmIssue(getContext(nil),
		"123456",
		"notes",
		model.RefData{
			Handle:     "MISSING_INFORMATION",
			Label:      "Missing information",
			Deprecated: false,
		},
		76,
	)
	assert.Equal(t, ErrUnauthorized, err)
}
