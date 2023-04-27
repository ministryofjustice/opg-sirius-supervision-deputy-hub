package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
				},
				{
					"id":4,
					"assuranceType": {
					  "handle": "PDR",
					  "label": "PDR",
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
					"assuranceVisitOutcome": null,
					"pdrOutcome": {
					  "handle": "RECEIVED",
					  "label": "Received",
					  "deprecated": null
					},
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
					"note": ""
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

	expectedResponse := []AssuranceVisits{
		{
			VisitId:             3,
			AssuranceType:       AssuranceTypes{Handle: "VISIT", Label: "Visit"},
			RequestedDate:       GenerateTimeForTest(2023, time.April, 01, 15, 04, 5),
			RequestedBy:         User{UserId: 53, UserDisplayName: "case manager"},
			DeputyId:            1,
			CommissionedDate:    GenerateTimeForTest(2023, time.May, 01, 15, 04, 5),
			ReportDueDate:       GenerateTimeForTest(2023, time.May, 11, 15, 04, 5),
			ReportReceivedDate:  GenerateTimeForTest(2023, time.April, 22, 15, 04, 5),
			VisitOutcome:        VisitOutcomeTypes{Label: "Cancelled", Handle: "CANCELLED"},
			ReportReviewDate:    GenerateTimeForTest(2023, time.October, 01, 15, 04, 5),
			VisitReportMarkedAs: VisitRagRatingTypes{Label: "Red", Handle: "RED"},
			Note:                "This is just notes for something to show",
			VisitorAllocated:    "Jane Janeson",
			ReviewedBy:          User{UserId: 53, UserDisplayName: "case manager"},
		},
		{
			VisitId:             4,
			AssuranceType:       AssuranceTypes{Handle: "PDR", Label: "PDR"},
			RequestedDate:       GenerateTimeForTest(2023, time.April, 01, 15, 04, 5),
			RequestedBy:         User{UserId: 53, UserDisplayName: "case manager"},
			DeputyId:            1,
			CommissionedDate:    GenerateTimeForTest(2023, time.May, 01, 15, 04, 5),
			ReportDueDate:       GenerateTimeForTest(2023, time.May, 11, 15, 04, 5),
			ReportReceivedDate:  GenerateTimeForTest(2023, time.April, 22, 15, 04, 5),
			PdrOutcome:          PdrOutcomeTypes{Label: "Received", Handle: "RECEIVED"},
			ReportReviewDate:    GenerateTimeForTest(2023, time.October, 01, 15, 04, 5),
			VisitReportMarkedAs: VisitRagRatingTypes{Label: "Red", Handle: "RED"},
			Note:                "",
			VisitorAllocated:    "Jane Janeson",
			ReviewedBy:          User{UserId: 53, UserDisplayName: "case manager"},
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
