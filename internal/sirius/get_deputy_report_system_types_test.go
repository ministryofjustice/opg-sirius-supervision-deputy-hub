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

func TestGetDeputyReportSystemTypes(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
		{
			"handle": "CASPAR",
			"label": "Caspar"
		},
		{
			"handle": "SOFTBOX",
			"label": "Softbox"
		},
		{
			"handle": "CONTROCC",
			"label": "Controcc"
		},
		{
			"handle": "CASHFAC",
			"label": "CASHFAC"
		},
		{
			"handle": "OPGDIGITAL",
			"label": "OPG Digital"
		},
		{
			"handle": "OPGPAPER",
			"label": "OPG Paper"
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
			Handle: "CASPAR",
			Label:  "Caspar",
		},
		{
			Handle: "SOFTBOX",
			Label:  "Softbox",
		},
		{
			Handle: "CONTROCC",
			Label:  "Controcc",
		},
		{
			Handle: "CASHFAC",
			Label:  "CASHFAC",
		},
		{
			Handle: "OPGDIGITAL",
			Label:  "OPG Digital",
		},
		{
			Handle: "OPGPAPER",
			Label:  "OPG Paper",
		},
		{
			Handle: "UNKNOWN",
			Label:  "Unknown",
		},
	}

	deputyReportSystemTypes, err := client.GetDeputyReportSystemTypes(getContext(nil))

	assert.Equal(t, expectedResponse, deputyReportSystemTypes)
	assert.Equal(t, nil, err)
}

func TestGetDeputyReportSystemTypesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyReportSystemTypes, err := client.GetDeputyReportSystemTypes(getContext(nil))

	assert.Equal(t, []model.RefData(nil), deputyReportSystemTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reference-data/deputyReportSystem",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyReportSystemTypesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyReportSystemTypes, err := client.GetDeputyReportSystemTypes(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []model.RefData(nil), deputyReportSystemTypes)
}
