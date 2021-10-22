package sirius

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
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

	expectedResponse := DeputyEventCollection{
		DeputyEvent{
			TimelineEventId: 300,
			Timestamp:       "2021-09-09 14:01:59",
			EventType:       "DeputyLinkedToOrder",
			User:            User{UserId: 41, UserDisplayName: "system admin", UserPhoneNumber: "12345678"},
			Event: Event{
				DeputyID:    "76",
				DeputyName:  "Mx Bob Builder",
				OrderType:   "pfa",
				SiriusId:    "7000-0000-1995",
				OrderNumber: "03305972",
				Client:      []ClientPerson{{ClientName: "Test Name", ClientId: "63", ClientUid: "7000-0000-1961", ClientCourtRef: "40124126"}},
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

	expectedResponse := DeputyEventCollection(nil)

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

	expectedResponse := DeputyEventCollection(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyEvents)
}

func TestReformatEventType(t *testing.T) {
	expectedResponse := "DeputyOrderDetailsChanged"
	testDeputyEvent := "Opg\\Core\\Model\\Event\\Order\\DeputyOrderDetailsChanged"
	assert.Equal(t, expectedResponse, ReformatEventType(testDeputyEvent))
}

func TestSortTimeLineNewestOneFirst(t *testing.T) {
	unsortedData := []DeputyEventCollection{
		DeputyEvent{
			TimelineEventId: 388,
			Timestamp: "2020-01-20 14:39:34",
			EventType: "PaDetailsChanged",
			User: User{
				UserId: 51,
				UserDisplayName: "case manager",
				UserPhoneNumber: ,
			}
			Event: Event{
				OrderType: ,
			},
		}
	}
}

pre sort {388 2020-01-20 14:39:34 PaDetailsChanged {51 case manager 12345678} {   76   [{Deputy name  pa deputy} {Telephone  } {Email  } {Team or department name  PA Team 1 - (Supervision)}] []}} {389 2020-02-20 14:39:34 PersonContactDetailsChanged {51 case manager 12345678} {   76   [{mobileNumber  } {homePhoneNumber  }] []}} {390 2020-02-21 14:39:34 PADeputyCreated {51 case manager 12345678} {   76   [] []}} {391 2020-03-19 14:39:34 DeputyLinkedToOrder {51 case manager 12345678} {pfa 7000-0000-1995 22036651 76   [] [{Duke John Fearless 63 7000-0000-1961 66323745}]}} {392 2020-04-18 14:39:34 PersonStatusChanged {51 case manager 12345678} {   76   [] []}} {393 2020-05-16 14:39:34 DeputyOrderDetailsChanged {51 case manager 12345678} {pfa 7000-0000-1995 22036651 76   [{deputyType  PA}] [{Duke John Fearless 63 7000-0000-1961 66323745}]}} {394 2020-06-12 14:39:34 PaDetailsChanged {51 case manager 12345678} {   76   [{AddressLine1  pa deputy}] []}} {395 2020-01-10 14:39:34 DeputyLinkedToOrder {51 case manager 12345678} {pfa 7000-0000-2290 66345286 76   [] [{Duchess Agnes Burgundy 68 7000-0000-2175 7077934T}]}} {396 2020-01-08 14:39:34 PersonStatusChanged {51 case manager 12345678} {   76   [] []}} {397 2020-01-11 14:39:34 DeputyOrderDetailsChanged {51 case manager 12345678} {pfa 7000-0000-2290 66345286 76   [{deputyType  PA}] [{Duchess Agnes Burgundy 68 7000-0000-2175 7077934T}]}} {398 2020-06-06 08:39:34 DeputyLinkedToOrder {51 case manager 12345678} {pfa 7000-0000-2472 65784899 76   [] [{King Louis Dauphin 74 7000-0000-2449 58574224}]}} {399 2020-06-04 08:39:34 DeputyOrderDetailsChanged {51 case manager 12345678} {pfa 7000-0000-2472 65784899 76   [{deputyType  PA}] [{King Louis Dauphin 74 7000-0000-2449 58574224}]}} {400 2020-06-05 12:39:34 PaDetailsChanged {51 case manager 12345678} {   76   [{AddressLine2  1}] []}} {401 2020-06-06 10:39:34 PaDetailsChanged {51 case manager 12345678} {   76   [{AddressLine3  2}] []}}]

sorted

[{394 2020-06-12 14:39:34 PaDetailsChanged {51 case manager 12345678} {   76   [{AddressLine1  pa deputy}] []}} {401 2020-06-06 10:39:34 PaDetailsChanged {51 case manager 12345678} {   76   [{AddressLine3  2}] []}} {398 2020-06-06 08:39:34 DeputyLinkedToOrder {51 case manager 12345678} {pfa 7000-0000-2472 65784899 76   [] [{King Louis Dauphin 74 7000-0000-2449 58574224}]}} {400 2020-06-05 12:39:34 PaDetailsChanged {51 case manager 12345678} {   76   [{AddressLine2  1}] []}} {399 2020-06-04 08:39:34 DeputyOrderDetailsChanged {51 case manager 12345678} {pfa 7000-0000-2472 65784899 76   [{deputyType  PA}] [{King Louis Dauphin 74 7000-0000-2449 58574224}]}} {393 2020-05-16 14:39:34 DeputyOrderDetailsChanged {51 case manager 12345678} {pfa 7000-0000-1995 22036651 76   [{deputyType  PA}] [{Duke John Fearless 63 7000-0000-1961 66323745}]}} {392 2020-04-18 14:39:34 PersonStatusChanged {51 case manager 12345678} {   76   [] []}} {391 2020-03-19 14:39:34 DeputyLinkedToOrder {51 case manager 12345678} {pfa 7000-0000-1995 22036651 76   [] [{Duke John Fearless 63 7000-0000-1961 66323745}]}} {390 2020-02-21 14:39:34 PADeputyCreated {51 case manager 12345678} {   76   [] []}} {389 2020-02-20 14:39:34 PersonContactDetailsChanged {51 case manager 12345678} {   76   [{mobileNumber  } {homePhoneNumber  }] []}} {388 2020-01-20 14:39:34 PaDetailsChanged {51 case manager 12345678} {   76   [{Deputy name  pa deputy} {Telephone  } {Email  } {Team or department name  PA Team 1 - (Supervision)}] []}} {397 2020-01-11 14:39:34 DeputyOrderDetailsChanged {51 case manager 12345678} {pfa 7000-0000-2290 66345286 76   [{deputyType  PA}] [{Duchess Agnes Burgundy 68 7000-0000-2175 7077934T}]}} {395 2020-01-10 14:39:34 DeputyLinkedToOrder {51 case manager 12345678} {pfa 7000-0000-2290 66345286 76   [] [{Duchess Agnes Burgundy 68 7000-0000-2175 7077934T}]}} {396 2020-01-08 14:39:34 PersonStatusChanged {51 case manager 12345678} {   76   [] []}}]