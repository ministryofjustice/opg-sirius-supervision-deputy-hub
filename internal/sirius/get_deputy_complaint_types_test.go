package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDeputyComplaintTypes(t *testing.T) {
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

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []DeputyComplaintTypes{
		{
			"YES",
			"Yes",
		},
		{
			"NO",
			"No",
		},
		{
			"UNKNOWN",
			"Unknown",
		},
	}

	complaints, err := client.GetDeputyComplaintTypes(getContext(nil))

	assert.Equal(t, expectedResponse, complaints)
	assert.Equal(t, nil, err)
}

func TestGetDeputyComplaintTypesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	complaints, err := client.GetDeputyComplaintTypes(getContext(nil))

	assert.Equal(t, []DeputyComplaintTypes(nil), complaints)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reference-data/complaints",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyComplaintTypesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	complaints, err := client.GetDeputyComplaintTypes(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []DeputyComplaintTypes(nil), complaints)
}
