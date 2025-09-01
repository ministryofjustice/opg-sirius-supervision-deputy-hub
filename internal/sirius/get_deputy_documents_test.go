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

func TestDeputyDocumentsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `
	{
		"limit":25,
		"metadata":{"doctype":[],"direction":[]},
		"pages":{"current":1,"total":0},
		"total":0,
		"documents":[
			{
			  "displayDate":"30\/05\/2024 01:00:00",
			  "id":5,
			  "uuid":"f64b6b38-7f99-4f4e-82d3-bd908c6589f2",
			  "type":"Catch-up call",
			  "friendlyDescription":"Screenshot_2024_06_21_at_14_12_30.png",
			  "title":"Correspondence",
			  "createdDate":"24\/06\/2024 15:17:32",
			  "direction":"Outgoing",
			  "filename":"3311d50c3d744d3bab02e0ad5e8e5eeb_Screenshot_2024_06_21_at_14_12_30.png",
			  "mimeType":"image\/png",
			  "caseItems":[],
			  "persons":[{"uId":"7000-0000-1276"}],
			  "createdBy":{
				 "id":51,
				 "name":"system",
				 "displayName":"system admin",
				 "email":"system.admin@opgtest.com",
				 "surname":"admin"
			  },
			  "receivedDateTime":"30\/05\/2024 01:00:00",
			  "documentSource":"UPLOAD",
			  "childCount":0,
			  "subtype":"Catch-up call"
       		},
			 {
				  "displayDate":"01\/06\/2024 01:00:00",
				  "id":6,
				  "uuid":"1a382bbb-f14f-451f-81f4-6ff1d0c4ce64",
				  "type":"General",
				  "friendlyDescription":"Screenshot_2024_06_21_at_15_23_12.png",
				  "title":"Correspondence",
				  "createdDate":"24\/06\/2024 15:20:05",
				  "direction":"Incoming",
				  "filename":"1245b837fa40441e986a1b576db37592_Screenshot_2024_06_21_at_15_23_12.png",
				  "mimeType":"image\/png",
				  "caseItems":[],
				  "persons":[{"uId":"7000-0000-1276"}],
				  "createdBy":{
					 "id":51,
					 "name":"system",
					 "displayName":"system admin",
					 "email":"system.admin@opgtest.com",
					 "surname":"admin"
				  },
				  "receivedDateTime":"01\/06\/2024 01:00:00",
				  "documentSource":"UPLOAD",
				  "childCount":0,
				  "subtype":"General"
			}
		]
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := DocumentList{
		Documents: []model.Document{
			{
				Id:                  5,
				Type:                "Catch-up call",
				FriendlyDescription: "Screenshot_2024_06_21_at_14_12_30.png",
				CreatedDate:         "24/06/2024",
				Direction:           "Outgoing",
				Filename:            "3311d50c3d744d3bab02e0ad5e8e5eeb_Screenshot_2024_06_21_at_14_12_30.png",
				CreatedBy: model.User{
					ID:          51,
					Name:        "system admin",
					PhoneNumber: "",
					Email:       "system.admin@opgtest.com",
				},
				ReceivedDateTime: "30/05/2024",
			},
			{
				Id:                  6,
				Type:                "General",
				FriendlyDescription: "Screenshot_2024_06_21_at_15_23_12.png",
				CreatedDate:         "24/06/2024",
				Direction:           "Incoming",
				Filename:            "1245b837fa40441e986a1b576db37592_Screenshot_2024_06_21_at_15_23_12.png",
				CreatedBy: model.User{
					ID:          51,
					Name:        "system admin",
					PhoneNumber: "",
					Email:       "system.admin@opgtest.com",
				},
				ReceivedDateTime: "01/06/2024",
			},
		},
		Metadata:       Metadata{},
		TotalDocuments: 0,
		Pages:          Page{PageCurrent: 1},
	}

	deputyDocuments, err := client.GetDeputyDocuments(getContext(nil), 1, "receiveddatetime:desc")

	assert.Equal(t, expectedResponse, deputyDocuments)
	assert.Equal(t, nil, err)
}

func TestGetDeputyDocumentsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyDocuments, err := client.GetDeputyDocuments(getContext(nil), 76, "receiveddatetime:desc")

	expectedResponse := DocumentList{}

	assert.Equal(t, expectedResponse, deputyDocuments)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/persons/76/documents?&sort=receiveddatetime:desc",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyDocumentsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyDocuments, err := client.GetDeputyDocuments(getContext(nil), 76, "receiveddatetime:desc")

	expectedResponse := DocumentList{}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyDocuments)
}
