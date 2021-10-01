package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddNote(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
	"personId":76,
	"userId":51,
	"userDisplayName":"case manager",
	"userEmail":"case.manager@opgtest.com",
	"userPhoneNumber":"12345678",
	"id":127,
	"type":null,
	"noteType":"PA_DEPUTY_NOTE_CREATED",
	"description":"fake note text",
	"name":"fake note title",
	"createdTime":"28\/09\/2021 09:30:27",
	"direction":null
	}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	err := client.AddNote(getContext(nil), "fake note title", "fake note text", 76, 51)
	assert.Nil(t, err)
}

func TestAddDeputyNoteReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.AddNote(getContext(nil), "test title", "test note", 76, 51)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputy/76/create-note",
		Method: http.MethodPost,
	}, err)
}

func TestAddDeputyNotesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.AddNote(getContext(nil), "test title", "test note", 76, 51)

	assert.Equal(t, ErrUnauthorized, err)
}