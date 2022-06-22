package sirius

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAssuranceVisitsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"assuranceVisits":
			[
				{
					"id":3,
					"requestedDate":"2022-06-25T12:16:34+00:00",
					"requestedBy":
						{
							"id":53,
							"displayName":"case manager"
						}
				}
			]
		}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []AssuranceVisits{
		{
			VisitId:       3,
			RequestedDate: "25/06/2022",
			RequestedBy:   User{UserId: 53, UserDisplayName: "case manager"},
			DeputyId:      1,
		},
	}

	assuranceVisits, err := client.GetAssuranceVisits(getContext(nil), 1)

	assert.Equal(t, expectedResponse, assuranceVisits)
	assert.Equal(t, nil, err)
}

func TestGetAssuranceVisitsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	assuranceVisits, err := client.GetAssuranceVisits(getContext(nil), 76)

	expectedResponse := []AssuranceVisits(nil)

	assert.Equal(t, expectedResponse, assuranceVisits)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/76/assurance-visits",
		Method: http.MethodGet,
	}, err)
}

func TestGetAssuranceVisitsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	assuranceVisits, err := client.GetAssuranceVisits(getContext(nil), 76)

	expectedResponse := []AssuranceVisits(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, assuranceVisits)
}
