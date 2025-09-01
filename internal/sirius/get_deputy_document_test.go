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

func TestDocumentReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `
	{
		"id":5,
		"uuid":"9c287e0b-fd4d-44cd-8cf0-4069b02765a8",
		"type":"Assurance visit",
		"friendlyDescription":"nice-file.png",
		"title":"Correspondence",
		"createdDate":"23\/07\/2024 17:02:58",
		"direction":"Incoming",
		"filename":"ba8d223457a94ec9bf2d503e26ee1298_nice-file.png",
		"mimeType":"image\/png",
		"note":{
			"id":16,
			"type":"Catch-up call",
			"description":"test",
			"name":"Document nice-file.png added",
			"createdTime":"23\/07\/2024 16:02:58","direction":"Incoming"
		},
		"caseItems":[],
		"persons":[
			{
				"id":67,
				"uId":"7000-0000-2472"
			}
		],
		"replacedDate":"25\/07\/2024 12:24:48",
		"createdBy":{
			"id":51,
			"displayName":"system admin",
			"email":"system.admin@opgtest.com"
		},
		"replacedBy":{
			"id":101,
			"displayName":"PROTeam1 User1",
			"email":"pro1@opgtest.com"
		},
		"receivedDateTime":"04\/07\/2024 01:00:00",
		"documentSource":"UPLOAD",
		"metadata":[],
		"childCount":0,
		"subtype":"Catch-up call"
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := model.Document{
		Id:                  5,
		Type:                "Assurance visit",
		FriendlyDescription: "nice-file.png",
		CreatedDate:         "23/07/2024 17:02:58",
		Direction:           "Incoming",
		Filename:            "ba8d223457a94ec9bf2d503e26ee1298_nice-file.png",
		CreatedBy: model.User{
			ID:          51,
			Name:        "system admin",
			PhoneNumber: "",
			Email:       "system.admin@opgtest.com",
		},
		ReceivedDateTime: "04/07/2024 01:00:00",
		Note: model.DocumentNote{
			Description: "test",
			Name:        "Document nice-file.png added",
		},
	}

	document, err := client.GetDocumentById(getContext(nil), 1, 1)

	assert.Equal(t, expectedResponse, document)
	assert.Equal(t, nil, err)
}

func TestGetDocumentReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	document, err := client.GetDocumentById(getContext(nil), 76, 1)

	expectedResponse := model.Document{}

	assert.Equal(t, expectedResponse, document)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/deputies/76/documents/1",
		Method: http.MethodGet,
	}, err)
}

func TestGetDocumentReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	document, err := client.GetDocumentById(getContext(nil), 76, 1)

	expectedResponse := model.Document{}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, document)
}
