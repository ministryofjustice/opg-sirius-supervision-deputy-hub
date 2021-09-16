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

func TestDeputyEventsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
    {
      "id": 300,
      "hash": "AW",
      "timestamp": "2021-09-09 14:01:59",
      "eventType": "Opg\\Core\\Model\\Event\\Order\\DeputyLinkedToOrder",
      "user": {
        "id": 41,
        "phoneNumber": "12345678",
        "displayName": "system admin",
        "email": "system.admin@opgtest.com"
      },
      "event": {
        "orderType": "pfa",
        "orderUid": "7000-0000-1995",
        "orderId": "58",
        "orderCourtRef": "03305972",
        "courtReferenceNumber": "03305972",
        "courtReference": "03305972",
        "personType": "Deputy",
        "personId": "76",
        "personUid": "7000-0000-2530",
        "personName": "Mx Bob Builder",
        "personCourtRef": null,
        "additionalPersons": [
          {
            "personType": "Client",
            "personId": "63",
            "personUid": "7000-0000-1961",
            "personName": "Test Name",
            "personCourtRef": "40124126"
          }
        ]
      }
    }
  ]`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := DeputyEvents{
		DeputyEvent{
			TimelineEventId:    300,
			Timestamp: "2021-09-09 14:01:59",
			EventType: "DeputyLinkedToOrder",
			User:  User{UserId: 41, UserDisplayName: "system admin", UserPhoneNumber: "12345678",},
			Event: Event{
				DeputyID: "76",
				DeputyName: "Mx Bob Builder",
				OrderType:   "pfa",
				SiriusId:    "7000-0000-1995",
				OrderNumber: "58",
				Client: []ClientPerson{{ClientName: "Test Name", ClientId: "63", ClientUid: "7000-0000-1961", ClientCourtRef: "40124126" }},
			},
		},
	}

	deputyEvents, err := client.GetDeputyEvents(getContext(nil), 1)

	assert.Equal(t, expectedResponse, deputyEvents)
	assert.Equal(t, nil, err)
}

func TestGetDeputyEventsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyEvents, err := client.GetDeputyEvents(getContext(nil), 76)

	expectedResponse := DeputyEvents(nil)

	assert.Equal(t, expectedResponse, deputyEvents)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/timeline/76",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyEventsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyEvents, err := client.GetDeputyEvents(getContext(nil), 76)

	expectedResponse := DeputyEvents(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyEvents)
}