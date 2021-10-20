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

func TestDeputyClientReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := ` {
    "persons": [
      {
        "id": 67,
        "caseRecNumber": "67422477",
        "email": "john.fearless@example.com",
        "firstname": "John",
        "surname": "Fearless",
        "addressLine1": "94 Duckpit Lane",
        "addressLine2": "Upper Oddington",
        "addressLine3": "Canvey Island",
        "town": "",
        "county": "",
        "postcode": "GL566WQ",
        "country": "",
        "phoneNumber": "07960209814",
        "clientAccommodation": {
          "handle": "FAMILY MEMBER/FRIEND'S HOME",
          "label": "Family Member/Friend's Home (including spouse/civil partner)",
          "deprecated": false
        },
        "orders": [
          {
            "id": 59,
            "latestSupervisionLevel": {
              "supervisionLevel": {
                "handle": "GENERAL",
                "label": "General",
                "deprecated": null
              }
            },
            "orderDate": "01/12/2020",
            "orderStatus": {
              "handle": "ACTIVE",
              "label": "Active",
              "deprecated": false
            }
          },
          {
            "id": 60,
            "latestSupervisionLevel": {
              "supervisionLevel": {
                "handle": "GENERAL",
                "label": "General",
                "deprecated": null
              }
            },
            "orderDate": "01/12/2017",
            "orderStatus": {
              "handle": "ACTIVE",
              "label": "Active",
              "deprecated": false
            }
          }
        ],
        "oldestNonLodgedAnnualReport": {
          "dueDate": "01/01/2016",
          "revisedDueDate": "01/05/2016",
          "status": {
            "label": "Pending"
          }
        },
        "riskScore": 5
      }
    ]
  } `

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := DeputyClientDetails{
		DeputyClient{
			ClientId:          67,
			Firstname:         "John",
			Surname:           "Fearless",
			CourtRef:          "67422477",
			RiskScore:         5,
			AccommodationType: "Family Member/Friend's Home (including spouse/civil partner)",
			OrderStatus:       "Active",
			OldestReport: reportReturned{
				DueDate:        "01/01/2016",
				RevisedDueDate: "01/05/2016",
				StatusLabel:    "Pending",
			},
			SupervisionLevel: "General",
		},
	}

	deputyClientDetails, ariaTags, err := client.GetDeputyClients(getContext(nil), 1, "", "")

	assert.Equal(t, expectedResponse, deputyClientDetails)
	assert.Equal(t, ariaTags, AriaSorting{SurnameAriaSort: "none", ReportDueAriaSort: "none", CRECAriaSort: "none"})
	assert.Equal(t, nil, err)
}

func TestGetDeputyClientReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	deputyClientDetails, ariaTags, err := client.GetDeputyClients(getContext(nil), 1, "", "")

	expectedResponse := DeputyClientDetails(nil)
	assert.Equal(t, ariaTags, AriaSorting{SurnameAriaSort: "", ReportDueAriaSort: "", CRECAriaSort: ""})
	assert.Equal(t, expectedResponse, deputyClientDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/pa/1/clients",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyClientsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	deputyClientDetails, ariaTags, err := client.GetDeputyClients(getContext(nil), 1, "", "")
	assert.Equal(t, ariaTags, AriaSorting{SurnameAriaSort: "", ReportDueAriaSort: "", CRECAriaSort: ""})
	expectedResponse := DeputyClientDetails(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyClientDetails)
}

func SetUpTestData() DeputyClientDetails {
	clients := DeputyClientDetails{
		DeputyClient{
			ClientId:    99,
			Firstname:   "Go",
			Surname:     "Taskforce",
			RiskScore:   1,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "05/01/2017",
				RevisedDueDate: "05/05/2017",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    92,
			Firstname:   "Louis",
			Surname:     "Dauphin",
			RiskScore:   1,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "01/01/2000",
				RevisedDueDate: "null",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    87,
			Firstname:   "Margaret",
			Surname:     "Bavaria-Straubing",
			RiskScore:   2,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "22/01/2018",
				RevisedDueDate: "22/06/2018",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    76,
			Firstname:   "Agnes",
			Surname:     "Burgundy",
			RiskScore:   5,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "23/01/2017",
				RevisedDueDate: "null",
				StatusLabel:    "Non-compliant",
			},
		},
	}
	return clients
}

func TestAlphabeticalSort(t *testing.T) {
	testData := SetUpTestData()
	expectedAscendingResponse := DeputyClientDetails{
		DeputyClient{
			ClientId:    87,
			Firstname:   "Margaret",
			Surname:     "Bavaria-Straubing",
			RiskScore:   2,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "22/01/2018",
				RevisedDueDate: "22/06/2018",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    76,
			Firstname:   "Agnes",
			Surname:     "Burgundy",
			RiskScore:   5,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "23/01/2017",
				RevisedDueDate: "null",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    92,
			Firstname:   "Louis",
			Surname:     "Dauphin",
			RiskScore:   1,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "01/01/2000",
				RevisedDueDate: "null",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    99,
			Firstname:   "Go",
			Surname:     "Taskforce",
			RiskScore:   1,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "05/01/2017",
				RevisedDueDate: "05/05/2017",
				StatusLabel:    "Non-compliant",
			},
		},
	}

	expectedDescendingResponse := DeputyClientDetails{
		DeputyClient{
			ClientId:    99,
			Firstname:   "Go",
			Surname:     "Taskforce",
			RiskScore:   1,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "05/01/2017",
				RevisedDueDate: "05/05/2017",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    92,
			Firstname:   "Louis",
			Surname:     "Dauphin",
			RiskScore:   1,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "01/01/2000",
				RevisedDueDate: "null",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    76,
			Firstname:   "Agnes",
			Surname:     "Burgundy",
			RiskScore:   5,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "23/01/2017",
				RevisedDueDate: "null",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    87,
			Firstname:   "Margaret",
			Surname:     "Bavaria-Straubing",
			RiskScore:   2,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "22/01/2018",
				RevisedDueDate: "22/06/2018",
				StatusLabel:    "Non-compliant",
			},
		},
	}
	assert.Equal(t, AlphabeticalSort(testData, "asc"), expectedAscendingResponse)
	assert.Equal(t, AlphabeticalSort(testData, "desc"), expectedDescendingResponse)
}

// func TestCrecScoreSort(t *testing.T) {
// 	testData := SetUpTestData()
// 	expectedAscendingResponse := DeputyClientDetails{
// 		DeputyClient{
// 			ClientId:    92,
// 			Firstname:   "Louis",
// 			Surname:     "Dauphin",
// 			RiskScore:   1,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "01/01/2000",
// 				RevisedDueDate: "null",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    99,
// 			Firstname:   "Go",
// 			Surname:     "Taskforce",
// 			RiskScore:   1,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "05/01/2017",
// 				RevisedDueDate: "05/05/2017",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    87,
// 			Firstname:   "Margaret",
// 			Surname:     "Bavaria-Straubing",
// 			RiskScore:   2,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "22/01/2018",
// 				RevisedDueDate: "22/06/2018",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    76,
// 			Firstname:   "Agnes",
// 			Surname:     "Burgundy",
// 			RiskScore:   5,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "23/01/2017",
// 				RevisedDueDate: "null",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 	}

// 	expectedDescendingResponse := DeputyClientDetails{
// 		DeputyClient{
// 			ClientId:    76,
// 			Firstname:   "Agnes",
// 			Surname:     "Burgundy",
// 			RiskScore:   5,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "23/01/2017",
// 				RevisedDueDate: "null",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    87,
// 			Firstname:   "Margaret",
// 			Surname:     "Bavaria-Straubing",
// 			RiskScore:   2,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "22/01/2018",
// 				RevisedDueDate: "22/06/2018",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    99,
// 			Firstname:   "Go",
// 			Surname:     "Taskforce",
// 			RiskScore:   1,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "05/01/2017",
// 				RevisedDueDate: "05/05/2017",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    92,
// 			Firstname:   "Louis",
// 			Surname:     "Dauphin",
// 			RiskScore:   1,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "01/01/2000",
// 				RevisedDueDate: "null",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 	}
// 	assert.Equal(t, crecScoreSort(testData, "asc"), expectedAscendingResponse)
// 	assert.Equal(t, crecScoreSort(testData, "desc"), expectedDescendingResponse)
// }

// func TestReportDueScoreSortkate(t *testing.T) {
// 	testData := SetUpTestData()
// 	expectedAscendingResponse := DeputyClientDetails{
// 		DeputyClient{
// 			ClientId:    76,
// 			Firstname:   "Agnes",
// 			Surname:     "Burgundy",
// 			RiskScore:   5,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "23/01/2017",
// 				RevisedDueDate: "null",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    92,
// 			Firstname:   "Louis",
// 			Surname:     "Dauphin",
// 			RiskScore:   1,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "01/01/2000",
// 				RevisedDueDate: "null",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},

// 		DeputyClient{
// 			ClientId:    99,
// 			Firstname:   "Go",
// 			Surname:     "Taskforce",
// 			RiskScore:   1,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "05/01/2017",
// 				RevisedDueDate: "05/05/2017",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    87,
// 			Firstname:   "Margaret",
// 			Surname:     "Bavaria-Straubing",
// 			RiskScore:   2,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "22/01/2018",
// 				RevisedDueDate: "22/06/2018",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 	}

// 	expectedDescendingResponse := DeputyClientDetails{

// 		DeputyClient{
// 			ClientId:    87,
// 			Firstname:   "Margaret",
// 			Surname:     "Bavaria-Straubing",
// 			RiskScore:   2,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "22/01/2018",
// 				RevisedDueDate: "22/06/2018",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    99,
// 			Firstname:   "Go",
// 			Surname:     "Taskforce",
// 			RiskScore:   1,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "05/01/2017",
// 				RevisedDueDate: "05/05/2017",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    76,
// 			Firstname:   "Agnes",
// 			Surname:     "Burgundy",
// 			RiskScore:   5,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "23/01/2017",
// 				RevisedDueDate: "null",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 		DeputyClient{
// 			ClientId:    92,
// 			Firstname:   "Louis",
// 			Surname:     "Dauphin",
// 			RiskScore:   1,
// 			OrderStatus: "Active",
// 			OldestReport: reportReturned{
// 				DueDate:        "01/01/2000",
// 				RevisedDueDate: "null",
// 				StatusLabel:    "Non-compliant",
// 			},
// 		},
// 	}
// 	assert.Equal(t, ReportDueScoreSort(testData, "asc"), expectedAscendingResponse)
// 	assert.Equal(t, ReportDueScoreSort(testData, "desc"), expectedDescendingResponse)
// }

// func TestChangeSortButtonDirection(t *testing.T) {
// 	expectedResponse := "none"
// 	assert.Equal(t, "expectedResponse", changeSortButtonDirection("asc", "sort=report_due", "sort=surname"))
// }
