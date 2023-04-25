package sirius

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func TestDeputyEventsReturned(t *testing.T) {
//	mockClient := &mocks.MockClient{}
//	client, _ := NewClient(mockClient, "http://localhost:3000")
//
//	json := `[
//  {
//    "id": 300,
//    "hash": "AW",
//    "timestamp": "2021-09-09T14:01:59+00:00",
//    "eventType": "Opg\\Core\\Model\\Event\\Order\\DeputyLinkedToOrder",
//    "user": {
//      "id": 41,
//      "phoneNumber": "12345678",
//      "displayName": "system admin",
//      "email": "system.admin@opgtest.com"
//    },
//    "event": {
//      "orderType": "pfa",
//      "orderUid": "7000-0000-1995",
//      "orderId": "58",
//      "orderCourtRef": "03305972",
//      "courtReferenceNumber": "03305972",
//      "courtReference": "03305972",
//      "personType": "Deputy",
//      "personId": "76",
//      "personUid": "7000-0000-2530",
//      "personName": "Mx Bob Builder",
//      "personCourtRef": null,
//      "additionalPersons": [
//        {
//          "personType": "Client",
//          "personId": "63",
//          "personUid": "7000-0000-1961",
//          "personName": "Test Name",
//          "personCourtRef": "40124126"
//        }
//      ]
//    }
//  }
//]`
//
//	r := io.NopCloser(bytes.NewReader([]byte(json)))
//
//	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
//		return &http.Response{
//			StatusCode: 200,
//			Body:       r,
//		}, nil
//	}
//
//	expectedResponse := DeputyEventCollection{
//		DeputyEvent{
//			TimelineEventId: 300,
//			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2021-09-09T14:01:59+00:00"),
//			EventType:       "DeputyLinkedToOrder",
//			User:            User{UserId: 41, UserDisplayName: "system admin", UserPhoneNumber: "12345678"},
//			Event: Event{
//				DeputyID:    "76",
//				DeputyName:  "Mx Bob Builder",
//				OrderType:   "pfa",
//				SiriusId:    "7000-0000-1995",
//				OrderNumber: "03305972",
//				Client:      []ClientPerson{{ClientName: "Test Name", ClientId: "63", ClientUid: "7000-0000-1961", ClientCourtRef: "40124126"}},
//			},
//		},
//	}
//
//	deputyEvents, err := client.GetDeputyEvents(getContext(nil), 1)
//
//	assert.Equal(t, expectedResponse, deputyEvents)
//	assert.Equal(t, nil, err)
//}

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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2021-09-09T10:12:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-10-18T10:12:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-09-20T10:11:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-10-16T10:11:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2021-09-09T10:12:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-10-18T10:12:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-10-16T10:11:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-09-20T10:11:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-10-18T10:11:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-10-18T11:12:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-09-20T10:11:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-10-16T10:11:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-10-18T11:12:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-10-18T10:11:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-10-16T10:11:08+00:00"),
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
			Timestamp:       FormatDateTimeStringIntoDateTime("2006-01-02T15:04:05+00:00", "2020-09-20T10:11:08+00:00"),
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
