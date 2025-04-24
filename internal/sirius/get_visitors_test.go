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

func TestVisitorsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
    {
      "id": 1,
      "name": "John Johnson"
    },
    {
      "id": 2,
      "name": "Richard Richardson"
    },
    {
      "id": 3,
      "name": "Jack Jackson"
    }
  ]`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []model.Visitor{
		{
			ID:   1,
			Name: "John Johnson",
		},
		{
			ID:   2,
			Name: "Richard Richardson",
		},
		{
			ID:   3,
			Name: "Jack Jackson",
		},
	}

	visitors, err := client.GetVisitors(getContext(nil))

	assert.Equal(t, expectedResponse, visitors)
	assert.Equal(t, nil, err)
}

func TestGetVisitorsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	visitors, err := client.GetVisitors(getContext(nil))

	var expectedResponse []model.Visitor

	assert.Equal(t, expectedResponse, visitors)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/visitors",
		Method: http.MethodGet,
	}, err)
}

func TestGetVisitorsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	visitors, err := client.GetVisitors(getContext(nil))

	var expectedResponse []model.Visitor

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, visitors)
}
