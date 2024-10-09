package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestDeputyClientsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := ` {
  "clients": [
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
			"appliesFrom": "01/12/2020",
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
			 "appliesFrom": "01/12/2017",
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
      "riskScore": 5,
		"hasActiveREMWarning": true
    }
  ],
  "pages": {
    "current": 1,
    "total": 1
  },
  "metadata": {
    "totalActiveClients": 1
  },
  "total": 1
} `

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	clients := []DeputyClient{
		DeputyClient{
			ClientId:            67,
			Firstname:           "John",
			Surname:             "Fearless",
			CourtRef:            "67422477",
			RiskScore:           5,
			ClientAccommodation: model.RefData{Handle: "FAMILY MEMBER/FRIEND'S HOME", Label: "Family Member/Friend's Home (including spouse/civil partner)"},
			OrderStatus:         "Active",
			OldestReport: Report{
				DueDate:        "01/01/2016",
				RevisedDueDate: "01/05/2016",
				Status:         model.RefData{Label: "Pending"},
			},
			SupervisionLevel:    "General",
			HasActiveREMWarning: true,
		},
	}

	expectedResponse := ClientList{
		Clients: clients,
		Pages: Page{
			PageCurrent: 1,
			PageTotal:   1,
		},
		Metadata:     Metadata{TotalActiveClients: 1},
		TotalClients: 1,
	}

	deputyClientDetails, err := client.GetDeputyClients(getContext(nil), ClientListParams{
		1,
		25,
		1,
		"PA",
		"",
		[]string{},
		[]string{},
		[]string{},
	})

	assert.Equal(t, 1, deputyClientDetails.Metadata.TotalActiveClients)
	assert.Equal(t, expectedResponse, deputyClientDetails)
	assert.Equal(t, nil, err)
}

func TestGetDeputyClientsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	clientList, err := client.GetDeputyClients(getContext(nil), ClientListParams{
		1,
		25,
		1,
		"PA",
		"",
		[]string{"ACTIVE"},
		[]string{"COUNCIL RENTED", "NO ACCOMMODATION TYPE"},
		[]string{"GENERAL", "MINIMAL"},
	})

	expectedResponse := ClientList{}
	assert.Equal(t, expectedResponse, clientList)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/pa/1/clients?&limit=25&page=1&sort=&filter=order-status:ACTIVE,accommodation:COUNCIL%20RENTED,accommodation:NO%20ACCOMMODATION%20TYPE,supervision-level:GENERAL,supervision-level:MINIMAL",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyClientsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	clientList, err := client.GetDeputyClients(getContext(nil), ClientListParams{
		1,
		25,
		1,
		"PA",
		"",
		[]string{},
		[]string{},
		[]string{},
	})

	expectedResponse := ClientList{}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, clientList)
}

func Test_GetOrderStatus(t *testing.T) {
	tests := []struct {
		Scenario       string
		Order1Date     string
		Order1Status   string
		Order2Date     string
		Order2Status   string
		ExpectedOutput string
	}{
		{
			Scenario:       "Returns oldest active order",
			Order1Status:   "Active",
			Order1Date:     "12/01/2014",
			Order2Status:   "Open",
			Order2Date:     "12/01/2017",
			ExpectedOutput: "Active",
		},
		{
			Scenario:       "Returns oldest non active order",
			Order1Status:   "Closed",
			Order1Date:     "12/01/2014",
			Order2Status:   "Open",
			Order2Date:     "12/01/2017",
			ExpectedOutput: "Closed",
		},
		{
			Scenario:       "Ignores nil order date",
			Order1Status:   "Open",
			Order1Date:     "",
			Order2Status:   "Active",
			Order2Date:     "12/01/2014",
			ExpectedOutput: "Active",
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1)+" given:"+tc.Scenario, func(t *testing.T) {
			orderData := []Order{
				{
					OrderStatus:            model.RefData{Label: tc.Order1Status},
					LatestSupervisionLevel: latestSupervisionLevel{},
					OrderDate:              tc.Order1Date,
				},
				{
					OrderStatus:            model.RefData{Label: tc.Order2Status},
					LatestSupervisionLevel: latestSupervisionLevel{},
					OrderDate:              tc.Order2Date,
				},
			}

			assert.Equal(t, getOrderStatus(orderData), tc.ExpectedOutput)
		})
	}
}

func Test_GetMostRecentSupervisionLevel(t *testing.T) {
	tests := []struct {
		Scenario               string
		Order1AppliesFrom      string
		Order1SupervisionLevel string
		Order1Date             string
		Order2AppliesFrom      string
		Order2SupervisionLevel string
		Order2Date             string
		ExpectedOutput         string
	}{
		{
			Scenario:               "Returns most recent supervision level",
			Order1AppliesFrom:      "01/02/2020",
			Order1SupervisionLevel: "General",
			Order1Date:             "01/02/2020",
			Order2AppliesFrom:      "03/02/2020",
			Order2SupervisionLevel: "Minimal",
			Order2Date:             "12/01/2020",
			ExpectedOutput:         "Minimal",
		},
		{
			Scenario:               "Returns most recent supervision level for a nil",
			Order1AppliesFrom:      "01/02/2020",
			Order1SupervisionLevel: "General",
			Order1Date:             "01/02/2020",
			Order2AppliesFrom:      "",
			Order2SupervisionLevel: "Minimal",
			Order2Date:             "12/01/2020",
			ExpectedOutput:         "General",
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1)+" given:"+tc.Scenario, func(t *testing.T) {
			test := Order{
				LatestSupervisionLevel: latestSupervisionLevel{
					AppliesFrom:      tc.Order1AppliesFrom,
					SupervisionLevel: model.RefData{Label: tc.Order1SupervisionLevel},
				},
				OrderDate: tc.Order1Date,
			}

			test2 := Order{
				LatestSupervisionLevel: latestSupervisionLevel{
					AppliesFrom:      tc.Order2AppliesFrom,
					SupervisionLevel: model.RefData{Label: tc.Order2SupervisionLevel},
				},
				OrderDate: tc.Order2Date,
			}

			orderData := []Order{
				test,
				test2,
			}

			assert.Equal(t, getMostRecentSupervisionLevel(orderData), tc.ExpectedOutput)
		})
	}
}

func Test_GetActivePfaOrderMadeDate(t *testing.T) {
	tests := []struct {
		Scenario          string
		Order1CaseSubType string
		Order2CaseSubType string
		ExpectedOutput    string
	}{
		{
			Scenario:          "Returns null if only hw orders",
			Order1CaseSubType: "hw",
			Order2CaseSubType: "hw",
			ExpectedOutput:    "",
		},
		{
			Scenario:          "Returns pfa over hw orders",
			Order1CaseSubType: "pfa",
			Order2CaseSubType: "hw",
			ExpectedOutput:    "01/02/2020",
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1)+" given:"+tc.Scenario, func(t *testing.T) {
			orderData := []Order{
				Order{
					OrderDate:   "01/02/2020",
					CaseSubType: tc.Order1CaseSubType,
					OrderStatus: model.RefData{
						Label: "Active",
					},
				},
				Order{
					OrderDate:   "02/02/2020",
					CaseSubType: tc.Order2CaseSubType,
					OrderStatus: model.RefData{
						Label: "Active",
					},
				},
			}

			assert.Equal(t, getActivePfaOrderMadeDate(orderData), tc.ExpectedOutput)
		})
	}
}

func Test_GetHasHwOrder(t *testing.T) {
	tests := []struct {
		Scenario          string
		Order1CaseSubType string
		Order1OrderStatus string
		Order2CaseSubType string
		Order2OrderStatus string
		ExpectedOutput    bool
	}{
		{
			Scenario:          "Returns false if no hw orders",
			Order1CaseSubType: "pfa",
			Order1OrderStatus: "Active",
			Order2CaseSubType: "pfa",
			Order2OrderStatus: "Closed",
			ExpectedOutput:    false,
		},
		{
			Scenario:          "Returns false if no active hw orders",
			Order1CaseSubType: "hw",
			Order1OrderStatus: "Closed",
			Order2CaseSubType: "hw",
			Order2OrderStatus: "Closed",
			ExpectedOutput:    false,
		},
		{
			Scenario:          "Returns true if active hw order",
			Order1CaseSubType: "pfa",
			Order1OrderStatus: "Active",
			Order2CaseSubType: "hw",
			Order2OrderStatus: "Active",
			ExpectedOutput:    true,
		},
		{
			Scenario:          "Returns true if only active hw order",
			Order1CaseSubType: "pfa",
			Order1OrderStatus: "Closed",
			Order2CaseSubType: "hw",
			Order2OrderStatus: "Active",
			ExpectedOutput:    true,
		},
		{
			Scenario:          "Returns true if two active hw order",
			Order1CaseSubType: "hw",
			Order1OrderStatus: "Active",
			Order2CaseSubType: "hw",
			Order2OrderStatus: "Active",
			ExpectedOutput:    true,
		},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1)+" given:"+tc.Scenario, func(t *testing.T) {
			orderData := []Order{
				Order{
					OrderDate:   "01/02/2020",
					CaseSubType: tc.Order1CaseSubType,
					OrderStatus: model.RefData{
						Label: tc.Order1OrderStatus,
					},
				},
				Order{
					OrderDate:   "02/02/2020",
					CaseSubType: tc.Order2CaseSubType,
					OrderStatus: model.RefData{
						Label: tc.Order2OrderStatus,
					},
				},
			}

			assert.Equal(t, hasHwOrder(orderData), tc.ExpectedOutput)
		})
	}
}

func Test_GetActivePfaOrderMadeDateReturnsOnlyActivePfaOrders(t *testing.T) {
	tests := []struct {
		Input          string
		ExpectedOutput string
	}{
		{Input: "Closed", ExpectedOutput: ""},
		{Input: "Open", ExpectedOutput: ""},
		{Input: "Duplicate", ExpectedOutput: ""},
		{Input: "Active", ExpectedOutput: "01/02/2020"},
	}
	for k, tc := range tests {
		t.Run("scenario "+strconv.Itoa(k+1)+" given:"+tc.Input, func(t *testing.T) {
			orderData := []Order{
				Order{
					OrderDate:   "01/02/2020",
					CaseSubType: "pfa",
					OrderStatus: model.RefData{
						Label: tc.Input,
					},
				},
			}

			assert.Equal(t, getActivePfaOrderMadeDate(orderData), tc.ExpectedOutput)
		})
	}
}
