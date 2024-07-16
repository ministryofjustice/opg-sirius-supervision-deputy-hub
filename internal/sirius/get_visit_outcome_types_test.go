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

func TestGetVisitOutcomeTypes(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
		{
			"handle": "SUCCESSFUL",
			"label": "Successful"
		},
		{
			"handle": "ABORTED",
			"label": "Aborted"
		},
		{
			"handle": "CANCELLED",
			"label": "Cancelled"
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
			Handle: "SUCCESSFUL",
			Label:  "Successful",
		},
		{
			Handle: "ABORTED",
			Label:  "Aborted",
		},
		{
			Handle: "CANCELLED",
			Label:  "Cancelled",
		},
	}

	visitOutcomeTypes, err := client.GetVisitOutcomeTypes(getContext(nil))

	assert.Equal(t, expectedResponse, visitOutcomeTypes)
	assert.Equal(t, nil, err)
}

func TestGetVisitOutcomeTypesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	visitOutcomeTypes, err := client.GetVisitOutcomeTypes(getContext(nil))

	assert.Equal(t, []model.RefData(nil), visitOutcomeTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reference-data/visitOutcome",
		Method: http.MethodGet,
	}, err)
}

func TestGetVisitOutcomeTypesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	visitOutcomeTypes, err := client.GetVisitOutcomeTypes(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []model.RefData(nil), visitOutcomeTypes)
}
