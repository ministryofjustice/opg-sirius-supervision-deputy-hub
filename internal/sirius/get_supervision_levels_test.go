package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSupervisionLevels(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `
            {
				"supervisionLevel": [
					{"handle": "GENERAL", "label": "General"},
					{"handle": "MINIMAL", "label": "Minimal"}
				]
			}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []model.RefData{
		{
			Handle: "GENERAL",
			Label:  "General",
		},
		{
			Handle: "MINIMAL",
			Label:  "Minimal",
		},
	}

	supervisionLevel, err := client.GetSupervisionLevels(getContext(nil))

	assert.Equal(t, expectedResponse, supervisionLevel)
	assert.Equal(t, nil, err)
}

func TestGetSupervisionLevelsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	supervisionLevel, err := client.GetSupervisionLevels(getContext(nil))

	assert.Equal(t, []model.RefData(nil), supervisionLevel)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reference-data?filter=supervisionLevel",
		Method: http.MethodGet,
	}, err)
}

func TestGetSupervisionLevelsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	supervisionLevel, err := client.GetSupervisionLevels(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []model.RefData(nil), supervisionLevel)
}
