package sirius

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePaImportantInformation(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"id": 3,
		"monthlySpreadsheet": {
			"handle": "NO",
			"label": "No"
		},
		"independentVisitorCharges": {
			"handle": "NO",
			"label": "No"
		},
		"bankCharges": {
			"handle": "NO",
			"label": "No"
		},
		"apad": {
			"handle": "NO",
			"label": "No"
		},
		"reportSystemType": {
			"handle": "CONTROCC",
			"label": "Controcc"
		},
		"annualBillingInvoice": {
			"handle": "SCHEDULE AND INVOICE",
			"label": "Schedule and Invoice"
		},
		"otherImportantInformation": "important info"
	}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	formData := ImportantPaInformationDetails{
		DeputyType:                "PA",
		MonthlySpreadsheet:        "No",
		IndependentVisitorCharges: "No",
		BankCharges:               "No",
		APAD:                      "No",
		ReportSystem:              "Controcc",
		AnnualBillingInvoice:      "Schedule and Invoice",
		OtherImportantInformation: "important info",
	}

	err := client.UpdatePaImportantInformation(getContext(nil), ID, formData)
	assert.Nil(t, err)
}

func TestUpdatePaImportantInformationReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := ImportantPaInformationDetails{
		DeputyType:                "PA",
		MonthlySpreadsheet:        "No",
		IndependentVisitorCharges: "No",
		BankCharges:               "No",
		APAD:                      "No",
		ReportSystem:              "Controcc",
		AnnualBillingInvoice:      "Schedule and Invoice",
		OtherImportantInformation: "important info",
	}

	err := client.UpdatePaImportantInformation(getContext(nil), ID, formData)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    fmt.Sprintf("%v/api/v1/deputies/%d/important-information", svr.URL, ID),
		Method: http.MethodPut,
	}, err)
}

func TestUpdatePaImportantInformationReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := ImportantPaInformationDetails{
		DeputyType:                "PA",
		MonthlySpreadsheet:        "No",
		IndependentVisitorCharges: "No",
		BankCharges:               "No",
		APAD:                      "No",
		ReportSystem:              "Controcc",
		AnnualBillingInvoice:      "Schedule and Invoice",
		OtherImportantInformation: "important info",
	}

	err := client.UpdatePaImportantInformation(getContext(nil), ID, formData)

	assert.Equal(t, ErrUnauthorized, err)
}
