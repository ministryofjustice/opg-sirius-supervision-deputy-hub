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

func TestBulkAssignAssuranceVisitTasksToClients(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{"dueDate": "2256-02-03", "clientId": ["1", "2"]}`
	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	bulkAssignVisitTasksForm := BulkAssignAssuranceVisitTasksToClientsParams{DueDate: "2256-02-03", ClientIds: []string{"1", "2"}}

	_, err := client.BulkAssignAssuranceVisitTasksToClients(getContext(nil), bulkAssignVisitTasksForm, 76)
	assert.Equal(t, nil, err)
}

func TestBulkAssignAssuranceVisitTasksToClientsErrorIfNoId(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	bulkAssignVisitTasksForm := BulkAssignAssuranceVisitTasksToClientsParams{DueDate: "2256-02-03", ClientIds: []string{"1", "2"}}

	_, err := client.BulkAssignAssuranceVisitTasksToClients(getContext(nil), bulkAssignVisitTasksForm, 0)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/deputies/0/bulk-assurance-visit-tasks",
		Method: http.MethodPost,
	}, err)
}

func TestBulkAssignAssuranceVisitTasksToClientsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	bulkAssignVisitTasksForm := BulkAssignAssuranceVisitTasksToClientsParams{DueDate: "2256-02-03", ClientIds: []string{"1", "2"}}

	_, err := client.BulkAssignAssuranceVisitTasksToClients(getContext(nil), bulkAssignVisitTasksForm, 76)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/deputies/76/bulk-assurance-visit-tasks",
		Method: http.MethodPost,
	}, err)
}
