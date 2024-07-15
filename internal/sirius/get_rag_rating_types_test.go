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

func TestGetRagRatingTypes(t *testing.T) {
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

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []model.RAGRating{
		{
			Handle: "RED",
			Label:  "Red",
		},
		{
			Handle: "AMBER",
			Label:  "Amber",
		},
		{
			Handle: "GREEN",
			Label:  "Green",
		},
	}

	ragRatingTypes, err := client.GetRagRatingTypes(getContext(nil))

	assert.Equal(t, expectedResponse, ragRatingTypes)
	assert.Equal(t, nil, err)
}

func TestGetRagRatingTypesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	ragRatingTypes, err := client.GetRagRatingTypes(getContext(nil))

	assert.Equal(t, []model.RAGRating(nil), ragRatingTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reference-data/ragRating",
		Method: http.MethodGet,
	}, err)
}

func TestGetRagRatingTypesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	ragRatingTypes, err := client.GetRagRatingTypes(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []model.RAGRating(nil), ragRatingTypes)
}
