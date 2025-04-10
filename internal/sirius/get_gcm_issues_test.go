package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGcmIssuesReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `
	[
		{
			"id":1,
			"client":{"id":66,"caseRecNumber":"48217682","firstname":"Hamster","surname":"Person"},
			"createdDate":"13\/08\/2024",
			"createdByUser":{"id":101,"name":"PROTeam1","phoneNumber":"12345678","displayName":"PROTeam1 User1","deleted":false,"email":"pro1@opgtest.com","firstname":"PROTeam1","surname":"User1","roles":["OPG User","Case Manager"],"locked":false,"suspended":false},
			"notes":"Not happy we are missing info here",
			"GCMIssueType":{"handle":"MISSING_INFORMATION","label":"Missing information"},
			"closedByUser":{"id":101,"name":"PROTeam1","phoneNumber":"12345678","displayName":"PROTeam1 User1","deleted":false,"email":"pro1@opgtest.com","firstname":"PROTeam1","surname":"User1","roles":["OPG User","Case Manager"],"locked":false,"suspended":false},
			"closedOnDate":"15\/09\/2025"
		},
		{
			"id":2,
			"client":{"id":77,"caseRecNumber":"48215555","firstname":"Spongebob","surname":"Squarepants"},
			"createdDate":"01\/09\/2024",
			"createdByUser":{"id":102,"name":"OtherUser","phoneNumber":"12345678","displayName":"OtherUser Person2","deleted":false,"email":"pro1@opgtest.com","firstname":"OtherUser","surname":"Person2","roles":["OPG User","Case Manager"],"locked":false,"suspended":false},
			"notes":"Why can they not calculate the fees correctly",
			"GCMIssueType":{"handle":"DEPUTY_FEES_INCORRECT","label":"Deputy fees incorrect"}
		}
	]`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []GcmIssue{
		{
			Id: 1,
			Client: GcmClient{
				Id:            66,
				CaseRecNumber: "48217682",
				Firstname:     "Hamster",
				Surname:       "Person",
			},
			CreatedDate: "13/08/2024",
			CreatedByUser: UserInformation{
				Id:          101,
				Name:        "PROTeam1",
				DisplayName: "PROTeam1 User1",
			},
			Notes: "Not happy we are missing info here",
			GcmIssueType: model.RefData{
				Handle:     "MISSING_INFORMATION",
				Label:      "Missing information",
				Deprecated: false,
			},
			ClosedOnDate: "15/09/2025",
			ClosedByUser: UserInformation{
				Id:          101,
				Name:        "PROTeam1",
				DisplayName: "PROTeam1 User1",
			},
		},
		{
			Id: 2,
			Client: GcmClient{
				Id:            77,
				CaseRecNumber: "48215555",
				Firstname:     "Spongebob",
				Surname:       "Squarepants",
			},
			CreatedDate: "01/09/2024",
			CreatedByUser: UserInformation{
				Id:          102,
				Name:        "OtherUser",
				DisplayName: "OtherUser Person2",
			},
			Notes: "Why can they not calculate the fees correctly",
			GcmIssueType: model.RefData{
				Handle:     "DEPUTY_FEES_INCORRECT",
				Label:      "Deputy fees incorrect",
				Deprecated: false,
			},
		},
	}

	expectedClient, err := client.GetGCMIssues(getContext(nil), 76, GcmIssuesParams{})

	assert.Equal(t, expectedResponse, expectedClient)
	assert.Equal(t, nil, err)
}

func TestGetGcmIssuesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	contact, err := client.GetGCMIssues(getContext(nil), 76, GcmIssuesParams{})

	expectedResponse := []GcmIssue(nil)

	assert.Equal(t, expectedResponse, contact)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/deputies/76/gcm-issues?&filter=&sort=",
		Method: http.MethodGet,
	}, err)
}

func TestGetGcmIssuesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	contact, err := client.GetGCMIssues(getContext(nil), 76, GcmIssuesParams{})

	expectedResponse := []GcmIssue(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, contact)
}
