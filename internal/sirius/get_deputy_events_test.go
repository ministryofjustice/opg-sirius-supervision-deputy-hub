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

	json := `    {
		"personId": 555,
		"personName": "kate",
   }`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := DeputyEvents{
		DeputyEvent{
			DeputyID: 555,
			DeputyName: "kate",
		},
	}

	deputyEvents, err := client.GetDeputyEvents(getContext(nil), 76)

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

	expectedResponse := DeputyEvents{
		DeputyEvent{
			DeputyID: 555,
			DeputyName: "kate",
		},
	}

	assert.Equal(t, expectedResponse, deputyEvents)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/76",
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

	//expectedResponse := DeputyEvents {
	//	DeputyEvent{
	//		DeputyID: 555,
	//		DeputyName: "kate",
	//	},
	//}

	expectedResponse := DeputyEvents{
			DeputyEvent{
				TaskId:    392,
				Timestamp: "2021-09-09 14:01:59",
				EventType: "Opg\\Core\\Model\\Event\\Order\\DeputyLinkedToOrder",
				Event{
					OrderType:   "pfa",
					SiriusId:    "7000-0000-1995",
					OrderNumber: "58",
					DeputyID:    76,
				}
			}
	}


	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyEvents)
}