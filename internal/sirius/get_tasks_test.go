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

func TestGetTasksReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `
	{
		"limit":25,
		"pages":{"current":1,"total":1},
		"total":2,
		"tasks":[
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
    	},
		{
			"id":183,
			"type":"ORAL",
			"status":"Not started",
			"dueDate":"29\/11\/2022",
			"name":"",
			"description":"A client has been created",
			"ragRating":1,
			"assignee":{"id":61,"displayName":"Spongebob Squarepants"},
			"createdTime":"14\/11\/2022 12:02:01",
			"caseItems":[],
			"persons":[{"id":61,"uId":"7000-0000-1870","caseRecNumber":"92902877","salutation":"Maquis","firstname":"Antoine","middlenames":"","surname":"Burgundy","supervisionCaseOwner":{"id":22,"teams":[],"displayName":"Allocations - (Supervision)"}}],
			"clients":[{"id":61,"uId":"7000-0000-1870","caseRecNumber":"92902877","salutation":"Maquis","firstname":"Antoine","middlenames":"","surname":"Burgundy","supervisionCaseOwner":{"id":22,"teams":[],"displayName":"Allocations - (Supervision)"}}],
			"caseOwnerTask":false
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

	expectedResponse := TaskList{
		Tasks: []model.Task{
			{
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
			},
			{
				Id:      183,
				Type:    "ORAL",
				DueDate: "29/11/2022",
				Name:    "",
				Assignee: model.Assignee{
					Id:          61,
					Teams:       nil,
					DisplayName: "Spongebob Squarepants",
				},
				CreatedTime:   "14/11/2022 12:02:01",
				CaseOwnerTask: false,
				Notes:         "A client has been created",
			},
		},
		TotalTasks: 2,
		Pages: PageInformation{
			Current: 1,
			Total:   1,
		},
	}

	tasks, err := client.GetTasks(getContext(nil), 76)

	assert.Equal(t, expectedResponse, tasks)
	assert.Equal(t, nil, err)
}

func TestGetTasksReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	tasks, err := client.GetTasks(getContext(nil), 76)

	assert.Equal(t, TaskList{Tasks: []model.Task(nil), TotalTasks: 0, Pages: PageInformation{Current: 0, Total: 0}}, tasks)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/deputies/76/tasks?filter=status:Not+started&sort=dueDate:asc",
		Method: http.MethodGet,
	}, err)
}

func TestGetTasksReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	tasks, err := client.GetTasks(getContext(nil), 76)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, TaskList{Tasks: []model.Task(nil), TotalTasks: 0, Pages: PageInformation{Current: 0, Total: 0}}, tasks)
}
