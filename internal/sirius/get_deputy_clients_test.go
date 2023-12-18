package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeputyClientReturned(t *testing.T) {
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
			ClientAccommodation: label{Label: "Family Member/Friend's Home (including spouse/civil partner)"},
			OrderStatus:         "Active",
			OldestReport: Report{
				DueDate:        "01/01/2016",
				RevisedDueDate: "01/05/2016",
				Status:         label{"Pending"},
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

func TestGetDeputyClientReturnsNewStatusError(t *testing.T) {
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

func TestGetOrderStatusReturnsOldestActiveOrder(t *testing.T) {
	orderData := []Order{
		{
			OrderStatus:            label{"Active"},
			LatestSupervisionLevel: latestSupervisionLevel{},
			OrderDate:              "12/01/2014",
		},
		{
			OrderStatus:            label{"Open"},
			LatestSupervisionLevel: latestSupervisionLevel{},
			OrderDate:              "12/01/2017",
		},
	}
	expectedResponse := "Active"
	result := getOrderStatus(orderData)

	assert.Equal(t, expectedResponse, result)
}

func TestGetOrderStatusNil(t *testing.T) {
	orderData := []Order{
		{
			OrderStatus:            label{"Open"},
			LatestSupervisionLevel: latestSupervisionLevel{},
			OrderDate:              "",
		},
		{
			OrderStatus:            label{"Active"},
			LatestSupervisionLevel: latestSupervisionLevel{},
			OrderDate:              "12/01/2014",
		},
		{
			OrderStatus:            label{"Open"},
			LatestSupervisionLevel: latestSupervisionLevel{},
			OrderDate:              "12/01/2017",
		},
	}
	expectedResponse := "Active"
	result := getOrderStatus(orderData)

	assert.Equal(t, expectedResponse, result)
}

func TestGetOrderStatusReturnsOldestNonActiveOrder(t *testing.T) {
	orderData := []Order{
		{
			OrderStatus:            label{"Closed"},
			LatestSupervisionLevel: latestSupervisionLevel{},
			OrderDate:              "12/01/2014",
		},
		{
			OrderStatus:            label{"Open"},
			LatestSupervisionLevel: latestSupervisionLevel{},
			OrderDate:              "12/01/2017",
		},
	}
	expectedResponse := "Closed"
	result := getOrderStatus(orderData)

	assert.Equal(t, expectedResponse, result)
}

func TestGetMostRecentSupervisionLevel(t *testing.T) {

	test := Order{
		LatestSupervisionLevel: latestSupervisionLevel{
			AppliesFrom:      "01/02/2020",
			SupervisionLevel: label{Label: "General"},
		},
		OrderDate: "01/02/2020",
	}

	test2 := Order{
		LatestSupervisionLevel: latestSupervisionLevel{
			AppliesFrom:      "03/02/2020",
			SupervisionLevel: label{Label: "Minimal"},
		},
		OrderDate: "12/01/2020",
	}

	orderData := []Order{
		test,
		test2,
	}
	expectedResponse := "Minimal"
	result := getMostRecentSupervisionLevel(orderData)

	assert.Equal(t, expectedResponse, result)
}

func TestGetMostRecentSupervisionLevelForANil(t *testing.T) {

	test := Order{
		LatestSupervisionLevel: latestSupervisionLevel{
			AppliesFrom:      "01/02/2020",
			SupervisionLevel: label{Label: "General"},
		},
		OrderDate: "01/02/2020",
	}

	test2 := Order{
		LatestSupervisionLevel: latestSupervisionLevel{
			AppliesFrom:      "03/02/2020",
			SupervisionLevel: label{Label: "Minimal"},
		},
		OrderDate: "12/01/2020",
	}

	test3 := Order{
		LatestSupervisionLevel: latestSupervisionLevel{
			AppliesFrom:      "",
			SupervisionLevel: label{Label: "Minimal"},
		},
		OrderDate: "12/01/2020",
	}

	orderData := []Order{
		test3,
		test,
		test2,
	}
	expectedResponse := "Minimal"
	result := getMostRecentSupervisionLevel(orderData)

	assert.Equal(t, expectedResponse, result)
}
