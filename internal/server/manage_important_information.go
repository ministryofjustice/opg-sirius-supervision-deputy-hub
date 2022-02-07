package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ManageProDeputyImportantInformation interface {
	GetDeputyDetails(sirius.Context, int, int) (sirius.DeputyDetails, error)
	UpdateProImportantInformation(sirius.Context, int, sirius.ImportantProInformationDetails) error
	UpdatePaImportantInformation(sirius.Context, int, sirius.ImportantPaInformationDetails) error
	GetDeputyAnnualInvoiceBillingTypes(ctx sirius.Context) ([]sirius.DeputyAnnualBillingInvoiceTypes, error)
	GetDeputyBooleanTypes(ctx sirius.Context) ([]sirius.DeputyBooleanTypes, error)
	GetDeputyReportSystemTypes(ctx sirius.Context) ([]sirius.DeputyReportSystemTypes, error)
}

type manageDeputyImportantInformationVars struct {
	Path                      string
	XSRFToken                 string
	DeputyDetails             sirius.DeputyDetails
	Error                     string
	Errors                    sirius.ValidationErrors
	DeputyId                  int
	AnnualBillingInvoiceTypes []sirius.DeputyAnnualBillingInvoiceTypes
	DeputyBooleanTypes        []sirius.DeputyBooleanTypes
	DeputyReportSystemTypes   []sirius.DeputyReportSystemTypes
}

func renderTemplateForImportantInformation(client ManageProDeputyImportantInformation, defaultPATeam int, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		deputyDetails, err := client.GetDeputyDetails(ctx, defaultPATeam, deputyId)
		if err != nil {
			return err
		}

		annualBillingInvoiceTypes, err := client.GetDeputyAnnualInvoiceBillingTypes(ctx)
		if err != nil {
			return err
		}

		deputyBooleanTypes, err := client.GetDeputyBooleanTypes(ctx)
		if err != nil {
			return err
		}

		deputyReportSystemTypes, err := client.GetDeputyReportSystemTypes(ctx)
		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:

			vars := manageDeputyImportantInformationVars{
				Path:                      r.URL.Path,
				XSRFToken:                 ctx.XSRFToken,
				DeputyId:                  deputyId,
				DeputyDetails:             deputyDetails,
				AnnualBillingInvoiceTypes: annualBillingInvoiceTypes,
				DeputyBooleanTypes:        deputyBooleanTypes,
				DeputyReportSystemTypes:   deputyReportSystemTypes,
			}
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			//alter based on deputy type calling different sirius files based on that

			if r.PostFormValue("complaints") != "" {
				var panelDeputyBool bool

				if r.PostFormValue("panel-deputy") != "" {
					panelDeputyBool, err = strconv.ParseBool(r.PostFormValue("panel-deputy"))
					if err != nil {
						return err
					}
				}

				importantInfoForm := sirius.ImportantProInformationDetails{
					DeputyType:                "Pro",
					Complaints:                r.PostFormValue("complaints"),
					PanelDeputy:               panelDeputyBool,
					AnnualBillingInvoice:      r.PostFormValue("annual-billing"),
					OtherImportantInformation: r.PostFormValue("other-info-note"),
				}

				err = client.UpdateProImportantInformation(ctx, deputyId, importantInfoForm)

			} else if r.PostFormValue("monthly-spreadsheet") != "" {

				reportSystemType := ""
				if r.PostFormValue("report-system") == "OPG Digital" {
					reportSystemType = "OPGDigital"
				} else if r.PostFormValue("report-system") == "OPG Paper" {
					reportSystemType = "OPGPaper"
				} else {
					reportSystemType = r.PostFormValue("report-system")
				}

				importantInfoForm := sirius.ImportantPaInformationDetails{
					DeputyType:                "Pa",
					MonthlySpreadsheet:        r.PostFormValue("monthly-spreadsheet"),
					IndependentVisitorCharges: r.PostFormValue("independent-visitor-charges"),
					BankCharges:               r.PostFormValue("bank-charges"),
					APAD:                      r.PostFormValue("apad"),
					ReportSystem:              reportSystemType,
					AnnualBillingInvoice:      r.PostFormValue("annual-billing"),
					OtherImportantInformation: r.PostFormValue("other-info-note"),
				}

				err = client.UpdatePaImportantInformation(ctx, deputyId, importantInfoForm)
			}

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameUpdateAdditionalInformationValidationErrorMessages(verr.Errors)

				vars := manageDeputyImportantInformationVars{
					Path:                      r.URL.Path,
					XSRFToken:                 ctx.XSRFToken,
					DeputyId:                  deputyId,
					DeputyDetails:             deputyDetails,
					Errors:                    verr.Errors,
					AnnualBillingInvoiceTypes: annualBillingInvoiceTypes,
					DeputyBooleanTypes:        deputyBooleanTypes,
					DeputyReportSystemTypes:   deputyReportSystemTypes,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=importantInformation", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func renameUpdateAdditionalInformationValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	errorCollection := sirius.ValidationErrors{}

	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)
			if fieldName == "otherImportantInformation" && errorType == "stringLengthTooLong" {
				err[errorType] = "The other important information must be 1000 characters or fewer"
				errorCollection["otherImportantInformation"] = err
			} else {
				err[errorType] = errorMessage
				errorCollection[fieldName] = err
			}
		}
	}
	return errorCollection
}
