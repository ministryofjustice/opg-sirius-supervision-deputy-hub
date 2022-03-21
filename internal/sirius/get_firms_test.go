package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFirmsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
	{
		"id":1,
		"firmName":"new firm 1",
		"firmNumber":1000000
	},
	{
		"id":2,
		"firmName":"firm 2",
		"firmNumber":1000001
	},
	{
		"id":3,
		"firmName":"firm 3",
		"firmNumber":1000002
	}
]`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []Firm{
		{
			Id:         1,
			FirmName:   "new firm 1",
			FirmNumber: 1000000,
		},
		{
			Id:         2,
			FirmName:   "firm 2",
			FirmNumber: 1000001,
		},
		{
			Id:         3,
			FirmName:   "firm 3",
			FirmNumber: 1000002,
		},
	}

	firms, err := client.GetFirms(getContext(nil))

	assert.Equal(t, expectedResponse, firms)
	assert.Equal(t, nil, err)
}

func TestGetFirmsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	firms, err := client.GetFirms(getContext(nil))

	assert.Equal(t, []Firm(nil), firms)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/firms",
		Method: http.MethodGet,
	}, err)
}

func TestGetFirmsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	firms, err := client.GetFirms(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []Firm(nil), firms)
}
