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

func TestDeputyNotesReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
    {
      "personId": 1,
      "userId": 68,
      "userDisplayName": "Finance User Testing",
      "userEmail": "finance.user.testing@opgtest.com",
      "userPhoneNumber": "12345678",
      "id": 65,
      "type": "NEW DEPUTY",
      "noteType": "ORDER_CREATED",
      "description": "notes",
      "name": "This is a HW order...",
      "createdTime": "20/09/2021 08:50:13",
      "direction": null
    },
    {
      "personId": 1,
      "userId": 68,
      "userDisplayName": "Finance User Testing",
      "userEmail": "finance.user.testing@opgtest.com",
      "userPhoneNumber": "12345678",
      "id": 64,
      "type": "OPEN",
      "noteType": "ORDER_STATUS_UPDATED",
      "description": "...and here are the order status notes",
      "name": null,
      "createdTime": "20/09/2021 08:50:12",
      "direction": null
    }
  ]`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := DeputyNoteCollection{
			DeputyNote{
				ID:              65,
				DeputyId:        1,
				UserId:          68,
				UserDisplayName: "Finance User Testing",
				UserEmail:       "finance.user.testing@opgtest.com",
				UserPhoneNumber: "12345678",
				Type:            "NEW DEPUTY",
				NoteType:        "ORDER_CREATED",
				NoteText:     "notes",
				Name:            "This is a HW order...",
				Timestamp:       "20/09/2021 08:50:13",
				Direction:       "",
			},
			DeputyNote{
				ID:              64,
				DeputyId:        1,
				UserId:          68,
				UserDisplayName: "Finance User Testing",
				UserEmail:       "finance.user.testing@opgtest.com",
				UserPhoneNumber: "12345678",
				Type:            "OPEN",
				NoteType:        "ORDER_STATUS_UPDATED",
				NoteText:     "...and here are the order status notes",
				Name:            "",
				Timestamp:       "20/09/2021 08:50:12",
				Direction:       "",
			},
	}

	deputyNotes, err := client.GetDeputyNotes(getContext(nil), 1)

	assert.Equal(t, expectedResponse, deputyNotes)
	assert.Equal(t, nil, err)
}

func TestGetDeputyNotesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyNotes, err := client.GetDeputyNotes(getContext(nil), 76)

	expectedResponse := DeputyNoteCollection(nil)

	assert.Equal(t, expectedResponse, deputyNotes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputy/76/notes",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyNotesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyNotes, err := client.GetDeputyNotes(getContext(nil), 76)

	expectedResponse := DeputyNoteCollection(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyNotes)
}