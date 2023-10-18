package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockManageDeputyImportantInformation struct {
	mock.Mock
}

func (m *mockManageDeputyImportantInformation) GetDeputyAnnualInvoiceBillingTypes(ctx sirius.Context) ([]sirius.DeputyAnnualBillingInvoiceTypes, error) {
	args := m.Called(ctx)

	return args.Get(0).([]sirius.DeputyAnnualBillingInvoiceTypes), args.Error(1)
}

func (m *mockManageDeputyImportantInformation) GetDeputyBooleanTypes(ctx sirius.Context) ([]sirius.DeputyBooleanTypes, error) {
	args := m.Called(ctx)

	return args.Get(0).([]sirius.DeputyBooleanTypes), args.Error(1)
}

func (m *mockManageDeputyImportantInformation) GetDeputyReportSystemTypes(ctx sirius.Context) ([]sirius.DeputyReportSystemTypes, error) {
	args := m.Called(ctx)

	return args.Get(0).([]sirius.DeputyReportSystemTypes), args.Error(1)
}

func (m *mockManageDeputyImportantInformation) UpdateImportantInformation(ctx sirius.Context, deputyID int, form sirius.ImportantInformationDetails) error {
	args := m.Called(ctx, deputyID, form)

	return args.Error(0)
}

func TestGetManageImportantInformation(t *testing.T) {
	assert := assert.New(t)

	deputyDetails := sirius.DeputyDetails{ID: 123}
	app := AppVars{
		DeputyDetails: deputyDetails,
		PageName:      "Manage important information",
	}
	invoiceTypes := []sirius.DeputyAnnualBillingInvoiceTypes{{Handle: "x", Label: "y"}}
	booleanTypes := []sirius.DeputyBooleanTypes{{Handle: "x", Label: "w"}}
	reportTypes := []sirius.DeputyReportSystemTypes{{Handle: "x", Label: "z"}}

	client := &mockManageDeputyImportantInformation{}
	client.On("GetUserDetails", mock.Anything).Return(sirius.UserDetails{Roles: []string{"Finance Manager"}}, nil)
	client.On("GetDeputyAnnualInvoiceBillingTypes", mock.Anything).Return(invoiceTypes, nil)
	client.On("GetDeputyBooleanTypes", mock.Anything).Return(booleanTypes, nil)
	client.On("GetDeputyReportSystemTypes", mock.Anything).Return(reportTypes, nil)

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForImportantInformation(client, template)
	err := handler(app, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(manageDeputyImportantInformationVars{
		AppVars:                   app,
		AnnualBillingInvoiceTypes: invoiceTypes,
		DeputyBooleanTypes:        booleanTypes,
		DeputyReportSystemTypes:   reportTypes,
	}, template.lastVars)
}

func TestPostManageImportantInformation(t *testing.T) {
	testCases := map[string]struct {
		app                         AppVars
		form                        url.Values
		importantInformationDetails sirius.ImportantInformationDetails
	}{
		"default": {
			app: AppVars{DeputyDetails: sirius.DeputyDetails{
				ID:         123,
				DeputyType: sirius.DeputyType{Handle: "x"},
			}},
			form: url.Values{},
			importantInformationDetails: sirius.ImportantInformationDetails{
				DeputyType:           "x",
				AnnualBillingInvoice: "Unknown",
			},
		},
		"previous value": {
			app: AppVars{DeputyDetails: sirius.DeputyDetails{
				ID:         123,
				DeputyType: sirius.DeputyType{Handle: "x"},
				DeputyImportantInformation: sirius.DeputyImportantInformation{
					AnnualBillingInvoice: sirius.HandleLabel{Label: "last-value"},
				},
			}},
			form: url.Values{},
			importantInformationDetails: sirius.ImportantInformationDetails{
				DeputyType:           "x",
				AnnualBillingInvoice: "last-value",
			},
		},
		"form value": {
			app: AppVars{DeputyDetails: sirius.DeputyDetails{
				ID:         123,
				DeputyType: sirius.DeputyType{Handle: "x"},
				DeputyImportantInformation: sirius.DeputyImportantInformation{
					AnnualBillingInvoice: sirius.HandleLabel{Label: "last-value"},
				},
			}},
			form: url.Values{"annual-billing": {"new-value"}},
			importantInformationDetails: sirius.ImportantInformationDetails{
				DeputyType:           "x",
				AnnualBillingInvoice: "new-value",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockManageDeputyImportantInformation{}
			client.On("GetUserDetails", mock.Anything).Return(sirius.UserDetails{}, nil)
			client.On("GetDeputyAnnualInvoiceBillingTypes", mock.Anything).Return([]sirius.DeputyAnnualBillingInvoiceTypes{}, nil)
			client.On("GetDeputyBooleanTypes", mock.Anything).Return([]sirius.DeputyBooleanTypes{}, nil)
			client.On("GetDeputyReportSystemTypes", mock.Anything).Return([]sirius.DeputyReportSystemTypes{}, nil)
			client.On("UpdateImportantInformation", mock.Anything, 123, tc.importantInformationDetails).Return(nil)

			template := &mockTemplates{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/123", strings.NewReader(tc.form.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			var redirect error

			testHandler := mux.NewRouter()
			testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
				redirect = renderTemplateForImportantInformation(client, template)(tc.app, w, r)
			})

			testHandler.ServeHTTP(w, r)

			assert.Equal(Redirect("/123?success=importantInformation"), redirect)
		})
	}
}

func TestCheckForReportSystemType(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(checkForReportSystemType("OPG Digital"), "OPGDigital")
	assert.Equal(checkForReportSystemType("OPG Paper"), "OPGPaper")
	assert.Equal(checkForReportSystemType("Other type"), "Other type")
}
