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

func TestGetDeputyTeamUsersReturnedPa(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
    "id": 23,
    "name": "PA Team 1 - (Supervision)",
    "phoneNumber": "0123456789",
    "displayName": "PA Team 1 - (Supervision)",
    "deleted": false,
    "email": "PATeam1.team@opgtest.com",
    "members": [
        {
            "id": 92,
            "name": "PATeam1",
            "phoneNumber": "12345678",
            "displayName": "PATeam1 User1",
            "deleted": false,
            "email": "pa1@opgtest.com"
        }
    ],
    "children": [],
    "teamType": {
        "handle": "PA",
        "label": "PA"
    }
	}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []TeamMember{
		{
			ID:          92,
			DisplayName: "PATeam1 User1",
			CurrentEcm:  1,
		},
	}

	deputyDetails := DeputyDetails{ID: 76, DeputyType: DeputyType{Handle: "PA", Label: "Public Authority"}, ExecutiveCaseManager: ExecutiveCaseManager{EcmId: 1}}
	paDeputyTeam, err := client.GetDeputyTeamMembers(getContext(nil), 23, deputyDetails)

	assert.Equal(t, expectedResponse, paDeputyTeam)
	assert.Equal(t, nil, err)
}

func TestGetDeputyTeamUsersReturnedPro(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[{
		"id": 25,
		"name": "Pro Team 1 - (Supervision)",
		"phoneNumber": "0123456789",
		"displayName": "Pro Team 1 - (Supervision)",
		"deleted": false,
		"email": "ProTeam1.team@opgtest.com",
		"members": [
			{
				"id": 94,
				"name": "PROTeam1",
				"phoneNumber": "12345678",
				"displayName": "ProTeam1 User1",
				"deleted": false,
				"email": "pro1@opgtest.com"
			}
		],
		"children": [],
		"teamType": {
			"handle": "PRO",
			"label": "Pro"
    	}
	},
	{
		"id": 26,
		"name": "Pro Team 2 - (Supervision)",
		"phoneNumber": "0123456789",
		"displayName": "Pro Team 2 - (Supervision)",
		"deleted": false,
		"email": "ProTeam2.team@opgtest.com",
		"members": [
			{
				"id": 95,
				"name": "PROTeam2",
				"phoneNumber": "12345678",
				"displayName": "ProTeam2 User1",
				"deleted": false,
				"email": "pro2@opgtest.com"
			}
		],
		"children": [],
		"teamType": {
			"handle": "PRO",
			"label": "Pro"
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

	expectedResponse := []TeamMember{
		{
			ID:          94,
			DisplayName: "ProTeam1 User1",
			CurrentEcm:  1,
		},
		{
			ID:          95,
			DisplayName: "ProTeam2 User1",
			CurrentEcm:  1,
		},
	}

	deputyDetails := DeputyDetails{ID: 76, DeputyType: DeputyType{Handle: "PRO", Label: "Professional"}, ExecutiveCaseManager: ExecutiveCaseManager{EcmId: 1}}
	proDeputyTeam, err := client.GetDeputyTeamMembers(getContext(nil), 23, deputyDetails)

	assert.Equal(t, expectedResponse, proDeputyTeam)
	assert.Equal(t, nil, err)
}

func TestGetDeputyTeamUsersReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	paDeputyTeam, err := client.GetDeputyTeamMembers(getContext(nil), 23, DeputyDetails{})

	expectedResponse := []TeamMember([]TeamMember{})

	assert.Equal(t, expectedResponse, paDeputyTeam)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/teams/23",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyTeamUsersReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	paDeputyTeam, err := client.GetDeputyTeamMembers(getContext(nil), 23, DeputyDetails{})

	expectedResponse := []TeamMember([]TeamMember{})

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, paDeputyTeam)
}
