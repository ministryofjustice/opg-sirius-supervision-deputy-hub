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

	returnedTypes := []model.TaskType{
		{
			Handle:        "CCC",
			Description:   "Both",
			ProDeputyTask: true,
			PaDeputyTask:  true,
		},
		{
			Handle:        "AAA",
			Description:   "Pro only",
			ProDeputyTask: true,
			PaDeputyTask:  false,
		},
	}

	taskTypes, err := client.GetTaskTypesForDeputyType(getContext(nil), "PRO")

	assert.Equal(t, returnedTypes, taskTypes)
	assert.Equal(t, nil, err)
}

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

	returnedTypes := []model.TaskType{
		{
			Handle:        "CCC",
			Description:   "Both",
			ProDeputyTask: true,
			PaDeputyTask:  true,
		},
		{
			Handle:        "BBB",
			Description:   "PA only",
			ProDeputyTask: false,
			PaDeputyTask:  true,
		},
	}

	taskTypes, err := client.GetTaskTypesForDeputyType(getContext(nil), "PA")

	assert.Equal(t, returnedTypes, taskTypes)
	assert.Equal(t, nil, err)
}

func TestGetTaskTypes_statusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	taskTypes, err := client.GetTaskTypesForDeputyType(getContext(nil), "PRO")

	assert.Equal(t, []model.TaskType(nil), taskTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/tasktypes/deputy",
		Method: http.MethodGet,
	}, err)
}

func TestGetTaskTypes_unauthorised(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	taskTypes, err := client.GetTaskTypesForDeputyType(getContext(nil), "PRO")

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []model.TaskType(nil), taskTypes)
}
