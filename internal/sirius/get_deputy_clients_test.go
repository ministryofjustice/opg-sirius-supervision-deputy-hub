package sirius

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	deputyClientDetails, ariaTags, err := client.GetDeputyClients(getContext(nil), 1, "PA", "", "")

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
	deputyClientDetails, ariaTags, err := client.GetDeputyClients(getContext(nil), 1, "PA", "", "")

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
	deputyClientDetails, ariaTags, err := client.GetDeputyClients(getContext(nil), 1, "PA", "", "")
	assert.Equal(t, ariaTags, AriaSorting{SurnameAriaSort: "", ReportDueAriaSort: "", CRECAriaSort: ""})
	expectedResponse := DeputyClientDetails(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyClientDetails)
}

func SetUpTestData() DeputyClientDetails {
	clients := DeputyClientDetails{
		DeputyClient{
			ClientId:    92,
			Firstname:   "Louis",
			Surname:     "Dauphin",
			RiskScore:   1,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "01/01/2000",
				RevisedDueDate: "05/05/3000",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    87,
			Firstname:   "Margaret",
			Surname:     "Bavaria-Straubing",
			RiskScore:   3,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "09/11/3018",
				RevisedDueDate: "",
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
				DueDate:        "03/01/2017",
				RevisedDueDate: "05/05/2017",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    99,
			Firstname:   "Go",
			Surname:     "Taskforce",
			RiskScore:   2,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "05/01/2017",
				RevisedDueDate: "",
				StatusLabel:    "Non-compliant",
			},
		},
	}

	return clients
}

func TestAlphabeticalSortAsc(t *testing.T) {
	testData := SetUpTestData()
	expectedAscendingResponse := DeputyClientDetails{
		DeputyClient{
			ClientId:    87,
			Firstname:   "Margaret",
			Surname:     "Bavaria-Straubing",
			RiskScore:   3,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "09/11/3018",
				RevisedDueDate: "",
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
				DueDate:        "03/01/2017",
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
				RevisedDueDate: "05/05/3000",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    99,
			Firstname:   "Go",
			Surname:     "Taskforce",
			RiskScore:   2,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "05/01/2017",
				RevisedDueDate: "",
				StatusLabel:    "Non-compliant",
			},
		},
	}

	assert.Equal(t, alphabeticalSort(testData, "asc"), expectedAscendingResponse)
}

func TestAlphabeticalSortDesc(t *testing.T) {
	testData := SetUpTestData()

	expectedDescendingResponse := DeputyClientDetails{
		DeputyClient{
			ClientId:    99,
			Firstname:   "Go",
			Surname:     "Taskforce",
			RiskScore:   2,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "05/01/2017",
				RevisedDueDate: "",
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
				RevisedDueDate: "05/05/3000",
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
				DueDate:        "03/01/2017",
				RevisedDueDate: "05/05/2017",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    87,
			Firstname:   "Margaret",
			Surname:     "Bavaria-Straubing",
			RiskScore:   3,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "09/11/3018",
				RevisedDueDate: "",
				StatusLabel:    "Non-compliant",
			},
		},
	}
	assert.Equal(t, alphabeticalSort(testData, "desc"), expectedDescendingResponse)
}

func TestCrecScoreSortAsc(t *testing.T) {
	testData := SetUpTestData()
	expectedAscendingResponse := DeputyClientDetails{
		DeputyClient{
			ClientId:    92,
			Firstname:   "Louis",
			Surname:     "Dauphin",
			RiskScore:   1,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "01/01/2000",
				RevisedDueDate: "05/05/3000",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    99,
			Firstname:   "Go",
			Surname:     "Taskforce",
			RiskScore:   2,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "05/01/2017",
				RevisedDueDate: "",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    87,
			Firstname:   "Margaret",
			Surname:     "Bavaria-Straubing",
			RiskScore:   3,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "09/11/3018",
				RevisedDueDate: "",
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
				DueDate:        "03/01/2017",
				RevisedDueDate: "05/05/2017",
				StatusLabel:    "Non-compliant",
			},
		},
	}

	assert.Equal(t, expectedAscendingResponse, crecScoreSort(testData, "asc"))
}

func TestCrecScoreSortDesc(t *testing.T) {
	testData := SetUpTestData()

	expectedDescendingResponse := DeputyClientDetails{
		DeputyClient{
			ClientId:    76,
			Firstname:   "Agnes",
			Surname:     "Burgundy",
			RiskScore:   5,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "03/01/2017",
				RevisedDueDate: "05/05/2017",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    87,
			Firstname:   "Margaret",
			Surname:     "Bavaria-Straubing",
			RiskScore:   3,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "09/11/3018",
				RevisedDueDate: "",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    99,
			Firstname:   "Go",
			Surname:     "Taskforce",
			RiskScore:   2,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "05/01/2017",
				RevisedDueDate: "",
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
				RevisedDueDate: "05/05/3000",
				StatusLabel:    "Non-compliant",
			},
		},
	}

	assert.Equal(t, expectedDescendingResponse, crecScoreSort(testData, "desc"))
}

func TestReportDueScoreSortAsc(t *testing.T) {
	testData := SetUpTestData()

	expectedAscendingResponse := DeputyClientDetails{
		DeputyClient{
			ClientId:    99,
			Firstname:   "Go",
			Surname:     "Taskforce",
			RiskScore:   2,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "05/01/2017",
				RevisedDueDate: "",
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
				DueDate:        "03/01/2017",
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
				RevisedDueDate: "05/05/3000",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    87,
			Firstname:   "Margaret",
			Surname:     "Bavaria-Straubing",
			RiskScore:   3,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "09/11/3018",
				RevisedDueDate: "",
				StatusLabel:    "Non-compliant",
			},
		},
	}

	assert.Equal(t, expectedAscendingResponse, reportDueScoreSort(testData, "asc"))
}

func TestReportDueScoreSortDesc(t *testing.T) {
	testData := SetUpTestData()

	expectedDescendingResponse := DeputyClientDetails{
		DeputyClient{
			ClientId:    87,
			Firstname:   "Margaret",
			Surname:     "Bavaria-Straubing",
			RiskScore:   3,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "09/11/3018",
				RevisedDueDate: "",
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
				RevisedDueDate: "05/05/3000",
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
				DueDate:        "03/01/2017",
				RevisedDueDate: "05/05/2017",
				StatusLabel:    "Non-compliant",
			},
		},
		DeputyClient{
			ClientId:    99,
			Firstname:   "Go",
			Surname:     "Taskforce",
			RiskScore:   2,
			OrderStatus: "Active",
			OldestReport: reportReturned{
				DueDate:        "05/01/2017",
				RevisedDueDate: "",
				StatusLabel:    "Non-compliant",
			},
		},
	}

	assert.Equal(t, expectedDescendingResponse, reportDueScoreSort(testData, "desc"))
}

func TestChangeSortButtonDirection(t *testing.T) {
	tests := []struct {
		sortOrder         string
		columnBeingSorted string
		functionCalling   string
		expectedResponse  string
	}{
		{sortOrder: "asc", columnBeingSorted: "sort=report_due", functionCalling: "sort=surname", expectedResponse: "none"},
		{sortOrder: "other", columnBeingSorted: "sort=report_due", functionCalling: "sort=report_due", expectedResponse: "none"},
		{sortOrder: "asc", columnBeingSorted: "sort=report_due", functionCalling: "sort=report_due", expectedResponse: "ascending"},
		{sortOrder: "desc", columnBeingSorted: "sort=report_due", functionCalling: "sort=report_due", expectedResponse: "descending"},
	}

	for _, tc := range tests {
		result := changeSortButtonDirection(tc.sortOrder, tc.columnBeingSorted, tc.functionCalling)
		assert.Equal(t, tc.expectedResponse, result)
	}
}

func TestSetDueDateForSortReturnDueDate(t *testing.T) {
	expectedResponse := "01/01/2021"
	result := setDueDateForSort("01/01/2021", "")
	assert.Equal(t, expectedResponse, result)
}

func TestSetDueDateForSortReturnRevisedDueDate(t *testing.T) {
	expectedResponse := "20/12/2021"
	result := setDueDateForSort("", "20/12/2021")
	assert.Equal(t, expectedResponse, result)
}

func TestSetDueDateForSortReturnZeroDateForNoDueOrRevisedDueDate(t *testing.T) {
	expectedResponse := "12/12/9999"
	result := setDueDateForSort("", "")
	assert.Equal(t, expectedResponse, result)
}

func TestFormatDate(t *testing.T) {
	expectedResponse, _ := time.Parse("2006-01-02", "2021-01-01")
	result := formatDate("01/01/2021")
	assert.Equal(t, expectedResponse, result)
}

func TestReformatCompletedDateReturnsNullIfNoDate(t *testing.T) {
	assert.Equal(t, "", reformatCompletedDate(""))
}

func TestReformatCompletedDateReturnsDateInCorrectFormat(t *testing.T) {
	assert.Equal(t, "25/02/2020", reformatCompletedDate("2020-02-25T00:00:00+00:00"))
}

func TestGetOrderStatusReturnsOldestActiveOrder(t *testing.T) {
	dateOne, _ := time.Parse("2006-02-01 00:00:00 +0000 UTC", "2014-12-01 00:00:00 +0000 UTC")
	dateTwo, _ := time.Parse("2006-02-01 00:00:00 +0000 UTC", "2017-12-01 00:00:00 +0000 UTC")

	orderData := Orders{
		Order{OrderStatus: "Active", SupervisionLevel: "General", OrderDate: dateOne},
		Order{OrderStatus: "Open", SupervisionLevel: "General", OrderDate: dateTwo},
	}
	expectedResponse := "Active"
	result := getOrderStatus(orderData)

	assert.Equal(t, expectedResponse, result)
}
func TestGetOrderStatusReturnsOldestNonActiveOrder(t *testing.T) {
	dateOne, _ := time.Parse("2006-02-01 00:00:00 +0000 UTC", "2014-12-01 00:00:00 +0000 UTC")
	dateTwo, _ := time.Parse("2006-02-01 00:00:00 +0000 UTC", "2017-12-01 00:00:00 +0000 UTC")

	orderData := Orders{
		Order{OrderStatus: "Close", SupervisionLevel: "General", OrderDate: dateOne},
		Order{OrderStatus: "Open", SupervisionLevel: "General", OrderDate: dateTwo},
	}
	expectedResponse := "Close"
	result := getOrderStatus(orderData)

	assert.Equal(t, expectedResponse, result)
}

func TestGetMostRecentSupervisionLevel(t *testing.T) {
	dateOne, _ := time.Parse("2006-02-01 00:00:00 +0000 UTC", "2014-12-01 00:00:00 +0000 UTC")
	dateTwo, _ := time.Parse("2006-02-01 00:00:00 +0000 UTC", "2017-12-01 00:00:00 +0000 UTC")

	orderData := Orders{
		Order{OrderStatus: "Close", SupervisionLevel: "Minimal", OrderDate: dateOne},
		Order{OrderStatus: "Open", SupervisionLevel: "General", OrderDate: dateTwo},
	}
	expectedResponse := "General"
	result := getMostRecentSupervisionLevel(orderData)

	assert.Equal(t, expectedResponse, result)
}

func TestRestructureOrders(t *testing.T) {

	unformattedDataOrder1 := apiOrder{}
	unformattedDataOrder1.OrderStatus.Label = "Active"
	unformattedDataOrder1.LatestSupervisionLevel.SupervisionLevel.Label = "Minimal"
	unformattedDataOrder1.OrderDate = "01/12/2014"

	unformattedData := apiOrders{
		unformattedDataOrder1,
	}

	dateOne, _ := time.Parse("02/01/2006", unformattedDataOrder1.OrderDate)

	expectedResponse := Orders{
		Order{OrderStatus: "Active", SupervisionLevel: "Minimal", OrderDate: dateOne},
	}
	assert.Equal(t, expectedResponse, restructureOrders(unformattedData))
}

func TestRestructureOrdersReturnsEmptySupervisionLevel(t *testing.T) {

	unformattedDataOrder1 := apiOrder{}
	unformattedDataOrder1.OrderStatus.Label = "Active"
	unformattedDataOrder1.LatestSupervisionLevel.SupervisionLevel.Label = ""
	unformattedDataOrder1.OrderDate = "01/12/2014"

	unformattedData := apiOrders{
		unformattedDataOrder1,
	}

	dateOne, _ := time.Parse("02/01/2006", unformattedDataOrder1.OrderDate)

	expectedResponse := Orders{
		Order{OrderStatus: "Active", SupervisionLevel: "", OrderDate: dateOne},
	}
	assert.Equal(t, expectedResponse, restructureOrders(unformattedData))
}

func TestRestructureOrdersReturnsNilForAnOpenOrder(t *testing.T) {

	unformattedDataOrder1 := apiOrder{}
	unformattedDataOrder1.OrderStatus.Label = "Open"
	unformattedDataOrder1.LatestSupervisionLevel.SupervisionLevel.Label = "General"
	unformattedDataOrder1.OrderDate = "01/12/2014"

	unformattedData := apiOrders{
		unformattedDataOrder1,
	}

	assert.Nil(t, restructureOrders(unformattedData))
}

func TestRestructureOrdersReturnsEmptyStringForNilSupervisionLevel(t *testing.T) {

	unformattedDataOrder1 := apiOrder{}
	unformattedDataOrder1.OrderStatus.Label = "Active"
	unformattedDataOrder1.OrderDate = "01/12/2014"

	unformattedData := apiOrders{
		unformattedDataOrder1,
	}

	dateOne, _ := time.Parse("02/01/2006", unformattedDataOrder1.OrderDate)

	expectedResponse := Orders{
		Order{OrderStatus: "Active", SupervisionLevel: "", OrderDate: dateOne},
	}
	assert.Equal(t, expectedResponse, restructureOrders(unformattedData))
}
