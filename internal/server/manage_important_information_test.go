package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
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

func (m *mockManageDeputyImportantInformation) GetDeputyAnnualInvoiceBillingTypes(ctx sirius.Context) ([]model.RefData, error) {
	args := m.Called(ctx)

	return args.Get(0).([]model.RefData), args.Error(1)
}

func (m *mockManageDeputyImportantInformation) GetDeputyBooleanTypes(ctx sirius.Context) ([]model.RefData, error) {
	args := m.Called(ctx)

	return args.Get(0).([]model.RefData), args.Error(1)
}

func (m *mockManageDeputyImportantInformation) GetDeputyReportSystemTypes(ctx sirius.Context) ([]model.RefData, error) {
	args := m.Called(ctx)

	return args.Get(0).([]model.RefData), args.Error(1)
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
	invoiceTypes := []model.RefData{{Handle: "TYPEA", Label: "TypeA"}, {Handle: "TYPEB", Label: "TypeB"}}
	booleanTypes := []model.RefData{{Handle: "x", Label: "w"}}
	reportTypes := []model.RefData{{Handle: "x", Label: "z"}}

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
				AnnualBillingInvoice: "UNKNOWN",
			},
		},
		"previous value": {
			app: AppVars{DeputyDetails: sirius.DeputyDetails{
				ID:         123,
				DeputyType: sirius.DeputyType{Handle: "x"},
				DeputyImportantInformation: sirius.DeputyImportantInformation{
					AnnualBillingInvoice: sirius.HandleLabel{Label: "TypeA", Handle: "TYPEA"},
				},
			}},
			form: url.Values{},
			importantInformationDetails: sirius.ImportantInformationDetails{
				DeputyType:           "x",
				AnnualBillingInvoice: "TYPEA",
			},
		},
		"form value": {
			app: AppVars{DeputyDetails: sirius.DeputyDetails{
				ID:         123,
				DeputyType: sirius.DeputyType{Handle: "x"},
				DeputyImportantInformation: sirius.DeputyImportantInformation{
					AnnualBillingInvoice: sirius.HandleLabel{Handle: "TYPEA", Label: "TypeA"},
				},
			}},
			form: url.Values{"annual-billing": {"TYPEB"}},
			importantInformationDetails: sirius.ImportantInformationDetails{
				DeputyType:           "x",
				AnnualBillingInvoice: "TYPEB",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockManageDeputyImportantInformation{}
			client.On("GetUserDetails", mock.Anything).Return(sirius.UserDetails{}, nil)
			client.On("GetDeputyAnnualInvoiceBillingTypes", mock.Anything).Return([]model.RefData{{Handle: "TYPEA", Label: "TypeA"}, {Handle: "TYPEB", Label: "TypeB"}}, nil)
			client.On("GetDeputyBooleanTypes", mock.Anything).Return([]model.RefData{}, nil)
			client.On("GetDeputyReportSystemTypes", mock.Anything).Return([]model.RefData{}, nil)
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
