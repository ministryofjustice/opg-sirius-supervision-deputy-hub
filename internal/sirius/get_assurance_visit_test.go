package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
				"requestedDate":"2023-04-01T15:04:05Z",
				"requestedBy": {
						"id":53,
						"displayName":"case manager"
				},
				"commissionedDate": "2023-05-01T15:04:05Z",
				"reportDueDate": "2023-05-11T15:04:05Z",
				"reportReceivedDate": "2023-04-22T15:04:05Z",
				"assuranceVisitOutcome": {
				  "handle": "CANCELLED",
				  "label": "Cancelled",
				  "deprecated": null
				},
				"pdrOutcome": null,
				"reportReviewDate": "2023-10-01T15:04:05Z",
				"assuranceVisitReportMarkedAs": {
				  "handle": "RED",
				  "label": "Red",
				  "deprecated": null
				},
				"visitorAllocated": "Jane Janeson",
				"reviewedBy": {
				  "id": 53,
				  "displayName": "case manager"
				},
				"note": "This is just notes for something to show"
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
		RequestedDate:       GenerateTimeForTest(2023, time.April, 01, 15, 04, 5),
		RequestedBy:         User{UserId: 53, UserDisplayName: "case manager"},
		CommissionedDate:    GenerateTimeForTest(2023, time.May, 01, 15, 04, 5),
		ReportDueDate:       GenerateTimeForTest(2023, time.May, 11, 15, 04, 5),
		ReportReceivedDate:  GenerateTimeForTest(2023, time.April, 22, 15, 04, 5),
		VisitOutcome:        VisitOutcomeTypes{Label: "Cancelled", Handle: "CANCELLED"},
		ReportReviewDate:    GenerateTimeForTest(2023, time.October, 01, 15, 04, 5),
		VisitReportMarkedAs: VisitRagRatingTypes{Label: "Red", Handle: "RED"},
		Note:                "This is just notes for something to show",
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
