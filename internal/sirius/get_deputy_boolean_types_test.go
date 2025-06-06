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

func TestGetDeputyBooleanTypes(t *testing.T) {
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

	deputyBooleanTypes, err := client.GetDeputyBooleanTypes(getContext(nil))

	assert.Equal(t, expectedResponse, deputyBooleanTypes)
	assert.Equal(t, nil, err)
}

func TestGetDeputyBooleanTypesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyBooleanTypes, err := client.GetDeputyBooleanTypes(getContext(nil))

	assert.Equal(t, []model.RefData(nil), deputyBooleanTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/reference-data/deputyBooleanType",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyBooleanTypesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyBooleanTypes, err := client.GetDeputyBooleanTypes(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []model.RefData(nil), deputyBooleanTypes)
}
