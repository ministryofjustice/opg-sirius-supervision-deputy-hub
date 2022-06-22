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

func TestGetVisitRagRatingTypes(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
		{
			"handle": "RED",
			"label": "Red"
		},
		{
			"handle": "AMBER",
			"label": "Amber"
		},
		{
			"handle": "GREEN",
			"label": "Green"
		}
]`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []VisitRagRatingTypes{
		{
			"RED",
			"Red",
		},
		{
			"AMBER",
			"Amber",
		},
		{
			"GREEN",
			"Green",
		},
	}

	visitRagRatingTypes, err := client.GetVisitRagRatingTypes(getContext(nil))

	assert.Equal(t, expectedResponse, visitRagRatingTypes)
	assert.Equal(t, nil, err)
}

func TestGetVisitRagRatingTypesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	visitRagRatingTypes, err := client.GetVisitRagRatingTypes(getContext(nil))

	assert.Equal(t, []VisitRagRatingTypes(nil), visitRagRatingTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reference-data/ragRating",
		Method: http.MethodGet,
	}, err)
}

func TestGetVisitRagRatingTypesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	visitRagRatingTypes, err := client.GetVisitRagRatingTypes(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []VisitRagRatingTypes(nil), visitRagRatingTypes)
}
