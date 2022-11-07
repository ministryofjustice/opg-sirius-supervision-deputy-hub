package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateAssuranceVisit(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"commissionedDate": "2022-06-17"
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	formData := AssuranceVisitDetails{
		CommissionedDate: "2022-06-17",
	}

	err := client.UpdateAssuranceVisit(getContext(nil), formData, 53, 76)
	assert.Nil(t, err)
}

func TestUpdateAssuranceVisitReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.UpdateAssuranceVisit(getContext(nil), AssuranceVisitDetails{}, 53, 76)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/53/assurance-visits/76",
		Method: http.MethodPut,
	}, err)
}

func TestUpdateAssuranceVisitReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.UpdateAssuranceVisit(getContext(nil), AssuranceVisitDetails{}, 53, 76)

	assert.Equal(t, ErrUnauthorized, err)
}
