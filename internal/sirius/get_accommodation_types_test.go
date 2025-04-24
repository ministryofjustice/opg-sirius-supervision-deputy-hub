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

func TestGetAccommodationTypes(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `
            {
				"clientAccommodation": [
					{"handle": "NO ACCOMMODATION TYPE", "label": "No Accommodation Type", "deprecated": false},
					{"handle": "COUNCIL RENTED", "label": "Council Rented", "deprecated": true}
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
			Handle: "HIGH RISK LIVING",
			Label:  "High Risk Living",
		},
		{
			Handle: "NO ACCOMMODATION TYPE",
			Label:  "No Accommodation Type",
		},
		{
			Handle:     "COUNCIL RENTED",
			Label:      "Council Rented",
			Deprecated: true,
		},
	}

	accommodationTypes, err := client.GetAccommodationTypes(getContext(nil))

	assert.Equal(t, expectedResponse, accommodationTypes)
	assert.Equal(t, nil, err)
}

func TestGetAccommodationTypesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	accommodationTypes, err := client.GetAccommodationTypes(getContext(nil))

	assert.Equal(t, []model.RefData(nil), accommodationTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/reference-data?filter=clientAccommodation",
		Method: http.MethodGet,
	}, err)
}

func TestGetAccommodationTypesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	accommodationTypes, err := client.GetAccommodationTypes(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []model.RefData(nil), accommodationTypes)
}
