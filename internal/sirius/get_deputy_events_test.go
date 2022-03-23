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
			Timestamp:       "09/09/2021 14:01:59",
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
	assert.Equal(t, expectedResponse, reformatEventType(testDeputyEvent))
}

func TestSortTimeLineNewestOneFirst(t *testing.T) {
	unsortedData := DeputyEventCollection{
		DeputyEvent{
			TimelineEventId: 388,
			Timestamp:       "19/10/2020 10:12:08",
			EventType:       "PersonContactDetailsChanged",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []Changes{
					{
						FieldName: "mobileNumber",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "homePhoneNumber",
						OldValue:  "null",
						NewValue:  "null",
					},
				},
				Client: []ClientPerson{},
			},
		},
		DeputyEvent{
			TimelineEventId: 387,
			Timestamp:       "18/10/2020 10:12:08",
			EventType:       "PaDetailsChanged",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []Changes{
					{
						FieldName: "Deputy name",
						OldValue:  "null",
						NewValue:  "PaDeputy",
					},
					{
						FieldName: "Telephone",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "Email",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "Teamordepartmentname",
						OldValue:  "null",
						NewValue:  "PA Team 1 - (Supervision)",
					},
				},
				Client: []ClientPerson{},
			},
		},
		DeputyEvent{
			TimelineEventId: 390,
			Timestamp:       "20/09/2020 10:11:08",
			EventType:       "DeputyLinkedToOrder",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "pfa",
				SiriusId:         "7000-0000-2381",
				OrderNumber:      "18372470",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []Changes{},
				Client: []ClientPerson{
					{
						ClientName:     "Duke John Fearless",
						ClientId:       "72",
						ClientUid:      "7000-0000-2357",
						ClientCourtRef: "2001022T",
					},
				},
			},
		},
		DeputyEvent{
			TimelineEventId: 389,
			Timestamp:       "16/10/2020 10:11:08",
			EventType:       "PADeputyCreated",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []Changes{},
				Client:           []ClientPerson{},
			},
		},
	}
	expectedResponse := DeputyEventCollection{
		DeputyEvent{
			TimelineEventId: 388,
			Timestamp:       "19/10/2020 10:12:08",
			EventType:       "PersonContactDetailsChanged",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []Changes{
					{
						FieldName: "mobileNumber",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "homePhoneNumber",
						OldValue:  "null",
						NewValue:  "null",
					},
				},
				Client: []ClientPerson{},
			},
		},
		DeputyEvent{
			TimelineEventId: 387,
			Timestamp:       "18/10/2020 10:12:08",
			EventType:       "PaDetailsChanged",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []Changes{
					{
						FieldName: "Deputy name",
						OldValue:  "null",
						NewValue:  "PaDeputy",
					},
					{
						FieldName: "Telephone",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "Email",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "Teamordepartmentname",
						OldValue:  "null",
						NewValue:  "PA Team 1 - (Supervision)",
					},
				},
				Client: []ClientPerson{},
			},
		},
		DeputyEvent{
			TimelineEventId: 389,
			Timestamp:       "16/10/2020 10:11:08",
			EventType:       "PADeputyCreated",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []Changes{},
				Client:           []ClientPerson{},
			},
		},
		DeputyEvent{
			TimelineEventId: 390,
			Timestamp:       "20/09/2020 10:11:08",
			EventType:       "DeputyLinkedToOrder",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "pfa",
				SiriusId:         "7000-0000-2381",
				OrderNumber:      "18372470",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []Changes{},
				Client: []ClientPerson{
					{
						ClientName:     "Duke John Fearless",
						ClientId:       "72",
						ClientUid:      "7000-0000-2357",
						ClientCourtRef: "2001022T",
					},
				},
			},
		},
	}
	assert.Equal(t, expectedResponse, sortTimeLineNewestOneFirst(unsortedData))
}

func TestEditDeputyEvents(t *testing.T) {
	unsortedData := DeputyEventCollection{
		DeputyEvent{
			TimelineEventId: 388,
			Timestamp:       "2020-10-18 10:11:08",
			EventType:       "Opg\\Core\\Model\\Event\\Order\\PersonContactDetailsChanged",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []Changes{
					{
						FieldName: "mobileNumber",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "homePhoneNumber",
						OldValue:  "null",
						NewValue:  "null",
					},
				},
				Client: []ClientPerson{},
			},
		},
		DeputyEvent{
			TimelineEventId: 387,
			Timestamp:       "2020-10-18 11:12:08",
			EventType:       "Opg\\Core\\Model\\Event\\Order\\PaDetailsChanged",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []Changes{
					{
						FieldName: "Deputy name",
						OldValue:  "null",
						NewValue:  "PaDeputy",
					},
					{
						FieldName: "Telephone",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "Email",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "Teamordepartmentname",
						OldValue:  "null",
						NewValue:  "PA Team 1 - (Supervision)",
					},
				},
				Client: []ClientPerson{},
			},
		},
		DeputyEvent{
			TimelineEventId: 390,
			Timestamp:       "2020-09-20 10:11:08",
			EventType:       "Opg\\Core\\Model\\Event\\Order\\DeputyLinkedToOrder",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "pfa",
				SiriusId:         "7000-0000-2381",
				OrderNumber:      "18372470",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []Changes{},
				Client: []ClientPerson{
					{
						ClientName:     "Duke John Fearless",
						ClientId:       "72",
						ClientUid:      "7000-0000-2357",
						ClientCourtRef: "2001022T",
					},
				},
			},
		},
		DeputyEvent{
			TimelineEventId: 389,
			Timestamp:       "2020-10-16 10:11:08",
			EventType:       "Opg\\Core\\Model\\Event\\Order\\PADeputyCreated",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []Changes{},
				Client:           []ClientPerson{},
			},
		},
	}
	expectedResponse := DeputyEventCollection{
		DeputyEvent{
			TimelineEventId: 387,
			Timestamp:       "18/10/2020 11:12:08",
			EventType:       "PaDetailsChanged",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []Changes{
					{
						FieldName: "Deputy name",
						OldValue:  "null",
						NewValue:  "PaDeputy",
					},
					{
						FieldName: "Telephone",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "Email",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "Teamordepartmentname",
						OldValue:  "null",
						NewValue:  "PA Team 1 - (Supervision)",
					},
				},
				Client: []ClientPerson{},
			},
		},
		DeputyEvent{
			TimelineEventId: 388,
			Timestamp:       "18/10/2020 10:11:08",
			EventType:       "PersonContactDetailsChanged",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []Changes{
					{
						FieldName: "mobileNumber",
						OldValue:  "null",
						NewValue:  "null",
					},
					{
						FieldName: "homePhoneNumber",
						OldValue:  "null",
						NewValue:  "null",
					},
				},
				Client: []ClientPerson{},
			},
		},
		DeputyEvent{
			TimelineEventId: 389,
			Timestamp:       "16/10/2020 10:11:08",
			EventType:       "PADeputyCreated",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []Changes{},
				Client:           []ClientPerson{},
			},
		},
		DeputyEvent{
			TimelineEventId: 390,
			Timestamp:       "20/09/2020 10:11:08",
			EventType:       "DeputyLinkedToOrder",
			User: User{
				UserId:          51,
				UserDisplayName: "case manager",
				UserPhoneNumber: "12345678",
			},
			Event: Event{
				OrderType:        "pfa",
				SiriusId:         "7000-0000-2381",
				OrderNumber:      "18372470",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []Changes{},
				Client: []ClientPerson{
					{
						ClientName:     "Duke John Fearless",
						ClientId:       "72",
						ClientUid:      "7000-0000-2357",
						ClientCourtRef: "2001022T",
					},
				},
			},
		},
	}
	assert.Equal(t, expectedResponse, editDeputyEvents(unsortedData))
}

func TestFormatDateAndTime(t *testing.T) {
	unsortedData := "2020-10-18 10:11:08"
	expectedResponse := "18/10/2020 10:11:08"
	assert.Equal(t, expectedResponse, formatDateAndTime(unsortedData))
}

func TestCalculateIfNewEvent(t *testing.T) {
	assert.Equal(t, true, calculateIfNewEvent(
		[]Changes{
			{
				FieldName: "firm",
				NewValue:  "new firm name",
			},
			{
				FieldName: "firmNumber",
				NewValue:  "1000028",
			},
		}))
	assert.Equal(t, false, calculateIfNewEvent(
		[]Changes{
			{
				FieldName: "firm",
				NewValue:  "a new firm name",
				OldValue:  "old firm name",
			},
			{
				FieldName: "firmNumber",
				NewValue:  "1000028",
				OldValue:  "1000021",
			},
		}))
}
