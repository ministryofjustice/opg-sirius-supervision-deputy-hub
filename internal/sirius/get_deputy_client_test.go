package sirius

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func TestDeputyClientReturned(t *testing.T) {
//	mockClient := &mocks.MockClient{}
//	client, _ := NewClient(mockClient, "http://localhost:3000")
//
//	json := `{
//               "id": 767,
//               "firstname": "Test",
//               "surname": "Client",
//               "caseRecNumber": "999555111",
//			}`
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
//	expectedResponse := ClientWithOrderDeputy{
//		ClientId:  767,
//		Firstname: "Test",
//		Surname:   "Client",
//		CourtRef:  "999555111",
//	}
//
//	expectedClient, err := client.GetDeputyClient(getContext(nil), "999555111", 76)
//
//	assert.Equal(t, expectedResponse, expectedClient)
//	assert.Equal(t, nil, err)
//}

func TestGetDeputyClientReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	contact, err := client.GetDeputyClient(getContext(nil), "123456", 76)

	expectedResponse := ClientWithOrderDeputy{}

	assert.Equal(t, expectedResponse, contact)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/clients/caserec/123456",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyClientReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	contact, err := client.GetDeputyClient(getContext(nil), "123456", 76)

	expectedResponse := ClientWithOrderDeputy{}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, contact)
}

func TestValidationErrorReturnedIfClientNotLinkedToDeputy(t *testing.T) {
	testClient := ClientWithOrderDeputy{
		ClientId:  123,
		Firstname: "test",
		Surname:   "Client",
		CourtRef:  "12345",
		Cases: []Case{
			{
				[]Deputies{
					{
						Deputy: OrderDeputy{
							77,
						},
					},
				},
			},
		},
	}

	assert.Equal(t, checkIfClientLinkedToDeputy(testClient, 77), true)
	assert.Equal(t, checkIfClientLinkedToDeputy(testClient, 76), false)
}

func TestCheckIfClientLinkedToDeputy(t *testing.T) {
	testClient := ClientWithOrderDeputy{
		ClientId:  123,
		Firstname: "test",
		Surname:   "Client",
		CourtRef:  "12345",
		Cases: []Case{
			{
				[]Deputies{
					{
						Deputy: OrderDeputy{
							77,
						},
					},
				},
			},
		},
	}

	assert.Equal(t, checkIfClientLinkedToDeputy(testClient, 77), true)
	assert.Equal(t, checkIfClientLinkedToDeputy(testClient, 76), false)
}
