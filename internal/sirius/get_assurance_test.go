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

func TestAssuranceReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
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
				"note" : "This is just to see the notes and it is below 1000 characters"
			}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := model.Assurance{
		Id:                 3,
		Type:               model.RefData{Handle: "VISIT", Label: "Assurance"},
		RequestedDate:      "2022-06-25",
		RequestedBy:        model.User{ID: 53, Name: "case manager"},
		CommissionedDate:   "2022-01-01",
		ReportDueDate:      "2022-01-07",
		ReportReceivedDate: "2022-01-07",
		VisitOutcome:       model.RefData{Label: "Cancelled", Handle: "CANCELLED"},
		PdrOutcome:         model.RefData{Label: "Received", Handle: "RECEIVED"},
		ReportReviewDate:   "2022-02-02",
		ReportMarkedAs:     model.RAGRating{Label: "Red", Handle: "RED"},
		Note:               "This is just to see the notes and it is below 1000 characters",
		VisitorAllocated:   "Jane Janeson",
		ReviewedBy:         model.User{ID: 53, Name: "case manager"},
	}

	assurance, err := client.GetAssuranceById(getContext(nil), 76, 3)

	assert.Equal(t, expectedResponse, assurance)
	assert.Equal(t, nil, err)
}

func TestGetAssuranceReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	assurance, err := client.GetAssuranceById(getContext(nil), 76, 1)

	expectedResponse := model.Assurance{}

	assert.Equal(t, expectedResponse, assurance)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/deputies/76/assurances/1",
		Method: http.MethodGet,
	}, err)
}

func TestGetAssuranceReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	assurance, err := client.GetAssuranceById(getContext(nil), 76, 1)

	expectedResponse := model.Assurance{}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, assurance)
}
