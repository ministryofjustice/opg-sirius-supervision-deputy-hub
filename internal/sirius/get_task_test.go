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

func TestGetTask(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `
		{
			"id":119,
			"type":"CWGN",
			"status":"Not started",
			"dueDate":"29\/11\/2022",
			"name":"",
			"description":"Case work general notes",
			"ragRating":1,
			"assignee":{"id":60,"displayName":"Case manager"},
			"createdTime":"14\/11\/2022 12:02:01",
			"caseItems":[],
			"persons":[{"id":61,"uId":"7000-0000-1870","caseRecNumber":"92902877","salutation":"Maquis","firstname":"Antoine","middlenames":"","surname":"Burgundy","supervisionCaseOwner":{"id":22,"teams":[],"displayName":"Allocations - (Supervision)"}}],
			"clients":[{"id":61,"uId":"7000-0000-1870","caseRecNumber":"92902877","salutation":"Maquis","firstname":"Antoine","middlenames":"","surname":"Burgundy","supervisionCaseOwner":{"id":22,"teams":[],"displayName":"Allocations - (Supervision)"}}],
			"caseOwnerTask":false
    	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := model.Task{
		Id:      119,
		Type:    "CWGN",
		DueDate: "29/11/2022",
		Name:    "",
		Assignee: model.Assignee{
			Id:          60,
			Teams:       nil,
			DisplayName: "Case manager",
		},
		CreatedTime:   "14/11/2022 12:02:01",
		CaseOwnerTask: false,
		Notes:         "Case work general notes",
	}

	task, err := client.GetTask(getContext(nil), 119)

	assert.Equal(t, expectedResponse, task)
	assert.Equal(t, nil, err)
}

func TestGetTaskReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	task, err := client.GetTask(getContext(nil), 119)

	assert.Equal(t, model.Task{}, task)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/tasks/119",
		Method: http.MethodGet,
	}, err)
}

func TestGetTaskReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	task, err := client.GetTask(getContext(nil), 119)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, model.Task{}, task)
}
