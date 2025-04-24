package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUpdateImportantInformationForPaInfo(t *testing.T) {
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

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	formData := ImportantInformationDetails{
		DeputyType:                "PA",
		MonthlySpreadsheet:        "NO",
		IndependentVisitorCharges: "NO",
		BankCharges:               "NO",
		APAD:                      "NO",
		ReportSystem:              "CONTROCC",
		AnnualBillingInvoice:      "SCHEDULE AND INVOICE",
		OtherImportantInformation: "important info",
	}

	err := client.UpdateImportantInformation(getContext(nil), ID, formData)
	assert.Nil(t, err)
}

func TestUpdateImportantInformationForProData(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"id": 3,
		"complaints": {
			"handle": "NO",
			"label": "No"
		},
		"panelDeputy": true,
		"annualBillingInvoice": {
			"handle": "SCHEDULE AND INVOICE",
			"label": "Schedule and Invoice"
		},
		"otherImportantInformation": "important info"
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	formData := ImportantInformationDetails{
		PanelDeputy:               true,
		Complaints:                "NO",
		AnnualBillingInvoice:      "SCHEDULE AND INVOICE",
		OtherImportantInformation: "important info",
	}

	err := client.UpdateImportantInformation(getContext(nil), ID, formData)
	assert.Nil(t, err)
}

func TestUpdateImportantInformationReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := ImportantInformationDetails{
		DeputyType:                "PA",
		MonthlySpreadsheet:        "NO",
		IndependentVisitorCharges: "NO",
		BankCharges:               "NO",
		APAD:                      "NO",
		ReportSystem:              "CONTROCC",
		AnnualBillingInvoice:      "SCHEDULE AND INVOICE",
		OtherImportantInformation: "important info",
	}

	err := client.UpdateImportantInformation(getContext(nil), ID, formData)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    fmt.Sprintf("%v/v1/deputies/%d/important-information", svr.URL + SupervisionAPIPath, ID),
		Method: http.MethodPut,
	}, err)
}

func TestUpdateImportantInformationReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := ImportantInformationDetails{
		DeputyType:                "PA",
		MonthlySpreadsheet:        "NO",
		IndependentVisitorCharges: "NO",
		BankCharges:               "NO",
		APAD:                      "NO",
		ReportSystem:              "CONTROCC",
		AnnualBillingInvoice:      "SCHEDULE AND INVOICE",
		OtherImportantInformation: "important info",
	}

	err := client.UpdateImportantInformation(getContext(nil), ID, formData)

	assert.Equal(t, ErrUnauthorized, err)
}
