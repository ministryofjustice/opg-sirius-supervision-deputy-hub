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

func TestGetPdrOutcomeTypes(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
		{
			"handle": "RECEIVED",
			"label": "Received"
		},
		{
			"handle": "NOT_RECEIVED",
			"label": "Not received"
		}
]`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []model.RefData{
		{
			Handle: "RECEIVED",
			Label:  "Received",
		},
		{
			Handle: "NOT_RECEIVED",
			Label:  "Not received",
		},
	}

	pdrOutcomeTypes, err := client.GetPdrOutcomeTypes(getContext(nil))

	assert.Equal(t, expectedResponse, pdrOutcomeTypes)
	assert.Equal(t, nil, err)
}

func TestGetPdrOutcomeTypesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	pdrOutcomeTypes, err := client.GetPdrOutcomeTypes(getContext(nil))

	assert.Equal(t, []model.RefData(nil), pdrOutcomeTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reference-data/pdrOutcome",
		Method: http.MethodGet,
	}, err)
}

func TestGetPdrOutcomeTypesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	pdrOutcomeTypes, err := client.GetPdrOutcomeTypes(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []model.RefData(nil), pdrOutcomeTypes)
}
