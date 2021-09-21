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
		DeputyClient{ClientId: 67, Firstname: "John", Surname: "Fearless", CourtRef: "67422477", RiskScore: 5, AccommodationType: "Family Member/Friend's Home (including spouse/civil partner)", OrderStatus: "Active", SupervisionLevel: "General"},
	}

	deputyClientDetails, err := client.GetDeputyClients(getContext(nil), 1)

	assert.Equal(t, expectedResponse, deputyClientDetails)
	assert.Equal(t, nil, err)
}

func TestGetDeputyClientReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	deputyClientDetails, err := client.GetDeputyClients(getContext(nil), 1)

	expectedResponse := DeputyClientDetails(nil)

	assert.Equal(t, expectedResponse, deputyClientDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/1/clients",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyClientsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	deputyClientDetails, err := client.GetDeputyClients(getContext(nil), 1)

	expectedResponse := DeputyClientDetails(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyClientDetails)
}
