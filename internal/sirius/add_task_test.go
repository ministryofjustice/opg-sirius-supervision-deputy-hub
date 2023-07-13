package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTask(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
	  "id": 1,
	  "type": "Task Type",
	  "createdByUser": {
		"id": 1
	  },
      "assignee": {
		"id": 1
	  },
	  "description": "<p>Note text<\/p>",
	  "createdTime": "26\/06\/2023 00:00:00"
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	err := client.AddTask(getContext(nil), 1, "AAAA", "", "2022-04-02", "test note", 1)
	assert.Nil(t, err)
}

func TestAddTaskReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.AddTask(getContext(nil), 1, "AAAA", "", "2022-04-02", "test note", 1)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/tasks",
		Method: http.MethodPost,
	}, err)
}

func TestAddTaskReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.AddTask(getContext(nil), 1, "AAAA", "", "2022-04-02", "test note", 1)

	assert.Equal(t, ErrUnauthorized, err)
}
