package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAssurancesReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"assurances":
			[
				{
					"id":3,
					"assuranceType": {
					  "handle": "VISIT",
					  "label": "Assurance",
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
					"pdrOutcome": null,
					"reportReviewDate": "2022-02-02T00:00:00+00:00",
					"reportMarkedAs": {
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
					"requestedDate":"2022-06-25T12:16:34+00:00",
					"requestedBy": {
							"id":53,
							"displayName":"case manager"
					},
					"commissionedDate": "2022-01-01T00:00:00+00:00",
					"reportDueDate": "2022-01-07T00:00:00+00:00",
					"reportReceivedDate": "2022-01-07T00:00:00+00:00",
					"assuranceVisitOutcome": null,
					"pdrOutcome": {
					  "handle": "RECEIVED",
					  "label": "Received",
					  "deprecated": null
					},
					"reportReviewDate": "2022-02-02T00:00:00+00:00",
					"reportMarkedAs": {
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

	expectedResponse := []model.Assurance{
		{
			Id:                 3,
			Type:               model.AssuranceType{Handle: "VISIT", Label: "Assurance"},
			RequestedDate:      "25/06/2022",
			RequestedBy:        model.User{ID: 53, Name: "case manager"},
			DeputyId:           1,
			CommissionedDate:   "01/01/2022",
			ReportDueDate:      "07/01/2022",
			ReportReceivedDate: "07/01/2022",
			VisitOutcome:       model.VisitOutcomeType{Label: "Cancelled", Handle: "CANCELLED"},
			ReportReviewDate:   "02/02/2022",
			ReportMarkedAs:     model.RagRatingType{Label: "Red", Handle: "RED"},
			Note:               "This is just notes for something to show",
			VisitorAllocated:   "Jane Janeson",
			ReviewedBy:         model.User{ID: 53, Name: "case manager"},
		},
		{
			Id:                 4,
			Type:               model.AssuranceType{Handle: "PDR", Label: "PDR"},
			RequestedDate:      "25/06/2022",
			RequestedBy:        model.User{ID: 53, Name: "case manager"},
			DeputyId:           1,
			CommissionedDate:   "01/01/2022",
			ReportDueDate:      "07/01/2022",
			ReportReceivedDate: "07/01/2022",
			PdrOutcome:         model.PdrOutcomeType{Label: "Received", Handle: "RECEIVED"},
			ReportReviewDate:   "02/02/2022",
			ReportMarkedAs:     model.RagRatingType{Label: "Red", Handle: "RED"},
			Note:               "",
			VisitorAllocated:   "Jane Janeson",
			ReviewedBy:         model.User{ID: 53, Name: "case manager"},
		},
	}

	assurances, err := client.GetAssurances(getContext(nil), 1)

	assert.Equal(t, expectedResponse, assurances)
	assert.Equal(t, nil, err)
}

func TestGetAssurancesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	assurances, err := client.GetAssurances(getContext(nil), 76)

	expectedResponse := []model.Assurance(nil)

	assert.Equal(t, expectedResponse, assurances)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/76/assurances",
		Method: http.MethodGet,
	}, err)
}

func TestGetAssurancesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	assurances, err := client.GetAssurances(getContext(nil), 76)

	expectedResponse := []model.Assurance(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, assurances)
}
