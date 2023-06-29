package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

var (
	resp = `{
		"task_types":{
			"AAA":{"handle":"AAA","incomplete":"Pro only","complete":"Pro only","user":true,"category":"supervision","proDeputyTask":true,"paDeputyTask":false},
			"BBB":{"handle":"BBB","incomplete":"PA only","complete":"PA only","user":true,"category":"supervision","proDeputyTask":false,"paDeputyTask":true},
			"CCC":{"handle":"CCC","incomplete":"Both","complete":"Both","user":true,"category":"supervision","proDeputyTask":true,"paDeputyTask":true},
			"DDD":{"handle":"DDD","incomplete":"Neither","complete":"Neither","user":true,"category":"supervision","proDeputyTask":false,"paDeputyTask":false}
		}
    }`
)

func TestGetTaskTypes_PA(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	r := io.NopCloser(bytes.NewReader([]byte(resp)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	deputy := DeputyDetails{
		DeputyType: DeputyType{Handle: "PRO"},
	}

	expectedResponse := []TaskType{
		{
			"AAA",
			"Pro only",
			true,
			false,
		}, {
			"CCC",
			"Both",
			true,
			true,
		},
	}

	taskTypes, err := client.GetTaskTypes(getContext(nil), deputy)

	sort.Slice(taskTypes, func(i, j int) bool {
		return taskTypes[i].Handle < taskTypes[j].Handle
	})

	assert.Equal(t, expectedResponse, taskTypes)
	assert.Equal(t, nil, err)
}

func TestGetTaskTypes_PRO(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	r := io.NopCloser(bytes.NewReader([]byte(resp)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	deputy := DeputyDetails{
		DeputyType: DeputyType{Handle: "PA"},
	}

	expectedResponse := []TaskType{
		{
			"BBB",
			"PA only",
			false,
			true,
		}, {
			"CCC",
			"Both",
			true,
			true,
		},
	}

	taskTypes, err := client.GetTaskTypes(getContext(nil), deputy)

	sort.Slice(taskTypes, func(i, j int) bool {
		return taskTypes[i].Handle < taskTypes[j].Handle
	})

	assert.Equal(t, expectedResponse, taskTypes)
	assert.Equal(t, nil, err)
}

func TestGetTaskTypes_statusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	taskTypes, err := client.GetTaskTypes(getContext(nil), DeputyDetails{})

	assert.Equal(t, []TaskType(nil), taskTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/tasktypes/deputy",
		Method: http.MethodGet,
	}, err)
}

func TestGetTaskTypes_unauthorised(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	taskTypes, err := client.GetTaskTypes(getContext(nil), DeputyDetails{})

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []TaskType(nil), taskTypes)
}
