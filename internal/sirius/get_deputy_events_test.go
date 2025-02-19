package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func AmendDateForDST(date string) string {
	return FormatDateTime(SiriusDateTime, date, SiriusDateTime)
}

type mockEventClient struct {
	responses []io.ReadCloser
	count     int
}

func (m *mockEventClient) Do(*http.Request) (*http.Response, error) {
	m.count++
	return &http.Response{
		StatusCode: 200,
		Body:       m.responses[m.count-1],
	}, nil
}

func TestDeputyEventsReturned(t *testing.T) {
	mockClient := mockEventClient{}
	client, _ := NewClient(&mockClient, "http://localhost:3000")

	eventJson := `
	{
		"limit":1,
		"metadata":[],
		"pages":{"current":1,"total":62},
		"total":62,
		"timelineEvents":[
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
		  },
		  {
			"id": 397,
			"hash": "AY",
			"timestamp": "2021-01-10 15:01:59",
			"eventType": "Opg\\Core\\Model\\Event\\Common\\TaskCreated",
			"user": {
			  "id": 21,
			  "phoneNumber": "0123456789",
			  "displayName": "Lay Team 1 - (Supervision)",
			  "email": "LayTeam1.team@opgtest.com"
			},
			"event": {
				"isCaseEvent": false,
				"isPersonEvent": true,
				"taskId": 249,
				"taskType": "AVFU",
				"dueDate": "2023-07-13 00:00:00",
				"notes": "This is a note",
				"name": "",
				"assignee": "PA Team Workflow",
				"isCaseOwnerTask": false,
				"personType": "Deputy",
				"personId": "78",
				"personUid": "7000-0000-2530",
				"personName": "Bobby Deputiser"
			  }
		  },
			{
				"id": 369,
				"hash": "A9",
				"timestamp": "2023-07-31 08:45:22",
				"eventType": "Opg\\Core\\Model\\Event\\Task\\TaskEdited",
				"user": {
				  "id": 21,
				  "phoneNumber": "0123456789",
				  "displayName": "Lay Team 1 - (Supervision)",
				  "email": "LayTeam1.team@opgtest.com"
				},
				"event": {
					"isCaseEvent": false,
					"isPersonEvent": true,
					"taskId": 184,
					"taskType": "AVFU",
					"dueDate": "2023-03-01 00:00:00",
					"notes": "Edited notes for edited task",
					"name": "",
					"assigneeId": 21,
					"assignee": "Lay Team 1 - (Supervision)",
					"isCaseOwnerTask": false,
					"oldAssigneeId": 60,
					"oldAssigneeName": "case manager",
					"wasCaseOwnerTask": false,
					"personType": "Deputy",
					"personId": "78",
					"personUid": "7000-0000-2530",
					"personName": "Bobby Deputiser",
					"changes": [
						{
							"fieldName": "dueDate",
							"oldValue": "01/03/2015",
							"newValue": "01/03/2023",
							"type": "string"
						},
						{
							"fieldName": "notes",
							"oldValue": "OG notes for edited task",
							"newValue": "Edited notes for edited task",
							"type": "string"
						}
					]
				}
			}
		]
	}
	`
	taskTypesJson := `
	{
           "task_types": {
               "SAP": {
                   "handle": "SAP",
                   "incomplete": "Start Assurance process",
                   "complete": "Start Assurance process",
                   "user": true,
                   "category": "deputy",
                   "ecmTask": false,
                   "proDeputyTask": true,
                   "paDeputyTask": false
               },
               "AVFU": {
                   "handle": "AVFU",
                   "incomplete": "Assurance visit follow up",
                   "complete": "Assurance visit follow up",
                   "user": true,
                   "category": "deputy",
                   "ecmTask": false,
                   "proDeputyTask": true,
                   "paDeputyTask": true
               }
           }
       }
	`

	mockClient.responses = append(
		mockClient.responses,
		io.NopCloser(bytes.NewReader([]byte(eventJson))),
		io.NopCloser(bytes.NewReader([]byte(taskTypesJson))))

	expectedResponse := TimelineList{
		Limit:    1,
		Metadata: []interface{}{},
		Pages: struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		}{
			Current: 1,
			Total:   62,
		},
		Total: 62,
		DeputyEvents: []model.DeputyEvent{
			model.DeputyEvent{
				ID:        369,
				Timestamp: AmendDateForDST("31/07/2023 08:45:22"),
				EventType: "TaskEdited",
				User:      model.User{ID: 21, Name: "Lay Team 1 - (Supervision)", PhoneNumber: "0123456789", Email: "LayTeam1.team@opgtest.com"},
				Event: model.Event{
					DeputyID:        "78",
					DeputyName:      "Bobby Deputiser",
					TaskType:        "Assurance visit follow up",
					Assignee:        "Lay Team 1 - (Supervision)",
					OldAssigneeName: "case manager",
					DueDate:         "01/03/2023",
					Notes:           "Edited notes for edited task",
					Changes: []model.Changes{
						{
							FieldName: "dueDate",
							OldValue:  "01/03/2015",
							NewValue:  "01/03/2023",
						},
						{
							FieldName: "notes",
							OldValue:  "OG notes for edited task",
							NewValue:  "Edited notes for edited task",
						},
					},
				},
			},
			model.DeputyEvent{
				ID:        300,
				Timestamp: AmendDateForDST("09/09/2021 14:01:59"),
				EventType: "DeputyLinkedToOrder",
				User:      model.User{ID: 41, Name: "system admin", PhoneNumber: "12345678", Email: "system.admin@opgtest.com"},
				Event: model.Event{
					DeputyID:    "76",
					DeputyName:  "Mx Bob Builder",
					OrderType:   "pfa",
					SiriusId:    "7000-0000-1995",
					OrderNumber: "03305972",
					Client:      []model.Client{{Name: "Test Name", ID: "63", Uid: "7000-0000-1961", CourtRef: "40124126"}},
				},
			},
			model.DeputyEvent{
				ID:        397,
				Timestamp: AmendDateForDST("10/01/2021 15:01:59"),
				EventType: "TaskCreated",
				User:      model.User{ID: 21, Name: "Lay Team 1 - (Supervision)", PhoneNumber: "0123456789", Email: "LayTeam1.team@opgtest.com"},
				Event: model.Event{
					DeputyID:   "78",
					DeputyName: "Bobby Deputiser",
					TaskType:   "Assurance visit follow up",
					Assignee:   "PA Team Workflow",
					DueDate:    "13/07/2023",
					Notes:      "This is a note",
				},
			},
		},
	}

	deputyEvents, err := client.GetDeputyEvents(getContext(nil), 1, 1, 25)

	assert.Equal(t, expectedResponse, deputyEvents)
	assert.Equal(t, nil, err)
}

func TestGetDeputyEventsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyEvents, err := client.GetDeputyEvents(getContext(nil), 76, 1, 25)

	expectedResponse := TimelineList{}

	assert.Equal(t, expectedResponse, deputyEvents)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/timeline/76/deputy?limit=25&page=1",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyEventsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyEvents, err := client.GetDeputyEvents(getContext(nil), 76, 1, 25)

	expectedResponse := TimelineList{}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyEvents)
}

func TestReformatEventType(t *testing.T) {
	expectedResponse := "DeputyOrderDetailsChanged"
	testDeputyEvent := "Opg\\Core\\Model\\Event\\Order\\DeputyOrderDetailsChanged"
	assert.Equal(t, expectedResponse, reformatEventType(testDeputyEvent))
}

func TestSortByTimelineAsc(t *testing.T) {
	unsortedData := DeputyEvents{
		model.DeputyEvent{
			ID:        388,
			Timestamp: "19/10/2020 10:12:08",
			EventType: "PersonContactDetailsChanged",
		},
		model.DeputyEvent{
			ID:        387,
			Timestamp: "18/10/2020 10:12:08",
			EventType: "PaDetailsChanged",
		},
		model.DeputyEvent{
			ID:        390,
			Timestamp: "20/09/2020 10:11:08",
			EventType: "DeputyLinkedToOrder",
		},
		model.DeputyEvent{
			ID:        389,
			Timestamp: "16/10/2020 10:11:08",
			EventType: "PADeputyCreated",
		},
	}
	expectedResponse := DeputyEvents{
		model.DeputyEvent{
			ID:        388,
			Timestamp: "19/10/2020 10:12:08",
			EventType: "PersonContactDetailsChanged",
		},
		model.DeputyEvent{
			ID:        387,
			Timestamp: "18/10/2020 10:12:08",
			EventType: "PaDetailsChanged",
		},
		model.DeputyEvent{
			ID:        389,
			Timestamp: "16/10/2020 10:11:08",
			EventType: "PADeputyCreated",
		},
		model.DeputyEvent{
			ID:        390,
			Timestamp: "20/09/2020 10:11:08",
			EventType: "DeputyLinkedToOrder",
		},
	}
	assert.Equal(t, expectedResponse, sortByTimelineAsc(unsortedData))
}

func TestEditDeputyEvents(t *testing.T) {
	unsortedData := DeputyEvents{
		model.DeputyEvent{
			ID:        388,
			Timestamp: "2020-10-18 10:11:08",
			EventType: "Opg\\Core\\Model\\Event\\Order\\PersonContactDetailsChanged",
			User: model.User{
				ID:          51,
				Name:        "case manager",
				PhoneNumber: "12345678",
			},
			Event: model.Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []model.Changes{
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
				Client: []model.Client{},
			},
		},
		model.DeputyEvent{
			ID:        387,
			Timestamp: "2020-10-18 11:12:08",
			EventType: "Opg\\Core\\Model\\Event\\Order\\PaDetailsChanged",
			User: model.User{
				ID:          51,
				Name:        "case manager",
				PhoneNumber: "12345678",
			},
			Event: model.Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []model.Changes{
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
				Client: []model.Client{},
			},
		},
		model.DeputyEvent{
			ID:        390,
			Timestamp: "2020-09-20 10:11:08",
			EventType: "Opg\\Core\\Model\\Event\\Order\\DeputyLinkedToOrder",
			User: model.User{
				ID:          51,
				Name:        "case manager",
				PhoneNumber: "12345678",
			},
			Event: model.Event{
				OrderType:        "pfa",
				SiriusId:         "7000-0000-2381",
				OrderNumber:      "18372470",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []model.Changes{},
				Client: []model.Client{
					{
						Name:     "Duke John Fearless",
						ID:       "72",
						Uid:      "7000-0000-2357",
						CourtRef: "2001022T",
					},
				},
			},
		},
		model.DeputyEvent{
			ID:        389,
			Timestamp: "2020-10-16 10:11:08",
			EventType: "Opg\\Core\\Model\\Event\\Order\\PADeputyCreated",
			User: model.User{
				ID:          51,
				Name:        "case manager",
				PhoneNumber: "12345678",
			},
			Event: model.Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []model.Changes{},
				Client:           []model.Client{},
			},
		},
	}
	expectedResponse := DeputyEvents{
		model.DeputyEvent{
			ID:        387,
			Timestamp: AmendDateForDST("18/10/2020 11:12:08"),
			EventType: "PaDetailsChanged",
			User: model.User{
				ID:          51,
				Name:        "case manager",
				PhoneNumber: "12345678",
			},
			Event: model.Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []model.Changes{
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
				Client: []model.Client{},
			},
		},
		model.DeputyEvent{
			ID:        388,
			Timestamp: AmendDateForDST("18/10/2020 10:11:08"),
			EventType: "PersonContactDetailsChanged",
			User: model.User{
				ID:          51,
				Name:        "case manager",
				PhoneNumber: "12345678",
			},
			Event: model.Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes: []model.Changes{
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
				Client: []model.Client{},
			},
		},
		model.DeputyEvent{
			ID:        389,
			Timestamp: AmendDateForDST("16/10/2020 10:11:08"),
			EventType: "PADeputyCreated",
			User: model.User{
				ID:          51,
				Name:        "case manager",
				PhoneNumber: "12345678",
			},
			Event: model.Event{
				OrderType:        "null",
				SiriusId:         "null",
				OrderNumber:      "null",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []model.Changes{},
				Client:           []model.Client{},
			},
		},
		model.DeputyEvent{
			ID:        390,
			Timestamp: AmendDateForDST("20/09/2020 10:11:08"),
			EventType: "DeputyLinkedToOrder",
			User: model.User{
				ID:          51,
				Name:        "case manager",
				PhoneNumber: "12345678",
			},
			Event: model.Event{
				OrderType:        "pfa",
				SiriusId:         "7000-0000-2381",
				OrderNumber:      "18372470",
				DeputyID:         "76",
				DeputyName:       "null",
				OrganisationName: "null",
				Changes:          []model.Changes{},
				Client: []model.Client{
					{
						Name:     "Duke John Fearless",
						ID:       "72",
						Uid:      "7000-0000-2357",
						CourtRef: "2001022T",
					},
				},
			},
		},
	}
	assert.Equal(t, expectedResponse, editDeputyEvents(unsortedData, TaskTypeMap{}))
}
func TestIsNewEvent(t *testing.T) {
	assert.Equal(t, true, isNewEvent(
		[]model.Changes{
			{
				FieldName: "firm",
				NewValue:  "new firm name",
			},
			{
				FieldName: "firmNumber",
				NewValue:  "1000028",
			},
		}))
	assert.Equal(t, false, isNewEvent(
		[]model.Changes{
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
