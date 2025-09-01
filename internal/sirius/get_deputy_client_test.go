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

func TestDeputyClientReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := ` {
		"id":66,
		"caseRecNumber":"43787324",
		"firstname":"Hamster",
		"surname":"Person"
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := DeputyClient{
		ClientId:  66,
		Firstname: "Hamster",
		Surname:   "Person",
		CourtRef:  "43787324",
	}

	deputyClient, err := client.GetDeputyClient(getContext(nil), "43787324", 67)

	assert.Equal(t, expectedResponse, deputyClient)
	assert.Equal(t, nil, err)
}

func TestDeputyClientReturnsValidationError(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := ` {"validation_errors":{"client-case-number":{"error":"Case number does not belong to this deputy"}}}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	expectedResponse := ValidationError(
		ValidationError{
			Message: "",
			Errors: ValidationErrors{
				"client-case-number": map[string]string{"error": "Case number does not belong to this deputy"},
			},
		},
	)

	actualClient, actualError := client.GetDeputyClient(getContext(nil), "43787324", 999)
	assert.Equal(t, expectedResponse, actualError)
	assert.Equal(t, DeputyClient{}, actualClient)

}

func TestGetDeputyClientReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	contact, err := client.GetDeputyClient(getContext(nil), "123456", 76)

	expectedResponse := DeputyClient{}

	assert.Equal(t, expectedResponse, contact)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/deputies/76/client/123456",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyClientReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	contact, err := client.GetDeputyClient(getContext(nil), "123456", 76)

	expectedResponse := DeputyClient{}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, contact)
}
