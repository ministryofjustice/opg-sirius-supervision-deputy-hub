package sirius

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAssuranceVisitReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
				"id":3,
				"assuranceType": {
				  "handle": "VISIT",
				  "label": "Visit",
				  "deprecated": null
				},
				"requestedDate":"2022-06-25T12:16:34+00:00",
				"requestedBy": {
						"id":53,
						"displayName":"case manager"
				},
				"commissionedDate": "2022-01-01T00:00:00+00:00",
				"reportDueDate": "2022-01-07T00:00:00+00:00",
				"reportReceivedDate": "2022-01-07T00:00:00+00:00",
				"assuranceVisitOutcome": {
				  "handle": "CANCELLED",
				  "label": "Cancelled",
				  "deprecated": null
				},
				"pdrOutcome": {
				  "handle": "RECEIVED",
				  "label": "Received",
				  "deprecated": null
				},
				"reportReviewDate": "2022-02-02T00:00:00+00:00",
				"assuranceVisitReportMarkedAs": {
				  "handle": "RED",
				  "label": "Red",
				  "deprecated": null
				},
				"visitorAllocated": "Jane Janeson",
				"reviewedBy": {
				  "id": 53,
				  "displayName": "case manager"
				}
			}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := AssuranceVisit{
		Id:                  3,
		AssuranceType:       AssuranceTypes{Handle: "VISIT", Label: "Visit"},
		RequestedDate:       "2022-06-25",
		RequestedBy:         User{UserId: 53, UserDisplayName: "case manager"},
		CommissionedDate:    "2022-01-01",
		ReportDueDate:       "2022-01-07",
		ReportReceivedDate:  "2022-01-07",
		VisitOutcome:        VisitOutcomeTypes{Label: "Cancelled", Handle: "CANCELLED"},
		PdrOutcome:          PdrOutcomeTypes{Label: "Received", Handle: "RECEIVED"},
		ReportReviewDate:    "2022-02-02",
		VisitReportMarkedAs: VisitRagRatingTypes{Label: "Red", Handle: "RED"},
		VisitorAllocated:    "Jane Janeson",
		ReviewedBy:          User{UserId: 53, UserDisplayName: "case manager"},
	}

	assuranceVisit, err := client.GetAssuranceVisitById(getContext(nil), 76, 3)

	assert.Equal(t, expectedResponse, assuranceVisit)
	assert.Equal(t, nil, err)
}

func TestGetAssuranceVisitReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	assuranceVisit, err := client.GetAssuranceVisitById(getContext(nil), 76, 1)

	expectedResponse := AssuranceVisit{}

	assert.Equal(t, expectedResponse, assuranceVisit)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/76/assurance-visits/1",
		Method: http.MethodGet,
	}, err)
}

func TestGetAssuranceVisitReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	assuranceVisit, err := client.GetAssuranceVisitById(getContext(nil), 76, 1)

	expectedResponse := AssuranceVisit{}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, assuranceVisit)
}
