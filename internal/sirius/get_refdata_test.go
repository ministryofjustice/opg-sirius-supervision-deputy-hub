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

func TestGetRefDataWithFilter(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"clientAccommodation": [
			{
				"handle": "CARE/NURSING/RESIDENTIAL HOME (PRIVATE/LA/REGISTERED)",
				"label": "Care/Nursing/Residential Home (Private/LA/Registered)",
				"deprecated": false
			},
			{
				"handle": "COUNCIL RENTED",
				"label": "Council Rented",
				"deprecated": false
			},
			{
				"handle": "FAMILY MEMBER/FRIEND'S HOME",
				"label": "Family Member/Friend's Home (including spouse/civil partner)",
				"deprecated": false
			}
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
			Handle: "CARE/NURSING/RESIDENTIAL HOME (PRIVATE/LA/REGISTERED)",
			Label:  "Care/Nursing/Residential Home (Private/LA/Registered)",
		},
		{
			Handle: "COUNCIL RENTED",
			Label:  "Council Rented",
		},
		{
			Handle: "FAMILY MEMBER/FRIEND'S HOME",
			Label:  "Family Member/Friend's Home (including spouse/civil partner)",
		},
	}

	accommodationTypes, err := client.GetRefData(getContext(nil), "?filter=clientAccommodation")

	assert.Equal(t, expectedResponse, accommodationTypes)
	assert.Equal(t, nil, err)
}

func TestGetRefDataWithoutFilter(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
       {
            "handle": "YES",
            "label": "Yes"
        },
        {
            "handle": "NO",
            "label": "No"
        },
        {
            "handle": "UNKNOWN",
            "label": "Unknown"
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
			Handle: "YES",
			Label:  "Yes",
		},
		{
			Handle: "NO",
			Label:  "No",
		},
		{
			Handle: "UNKNOWN",
			Label:  "Unknown",
		},
	}

	booleanTypes, err := client.GetRefData(getContext(nil), "deputyBooleanType")

	assert.Equal(t, expectedResponse, booleanTypes)
	assert.Equal(t, nil, err)
}

func TestGetRefDataReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	accommodationTypes, err := client.GetRefData(getContext(nil), "?filter=clientAccommodation")

	assert.Equal(t, []model.RefData(nil), accommodationTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reference-data?filter=clientAccommodation",
		Method: http.MethodGet,
	}, err)
}

func TestTestGetRefDataReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	booleanTypes, err := client.GetRefData(getContext(nil), "deputyBooleanType")

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []model.RefData(nil), booleanTypes)
}
