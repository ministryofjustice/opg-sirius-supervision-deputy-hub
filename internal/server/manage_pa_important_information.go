package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ManagePaDeputyImportantInformation interface {
	GetDeputyDetails(sirius.Context, int, int) (sirius.DeputyDetails, error)
	UpdateImportantInformation(sirius.Context, int, sirius.ImportantPaInformationDetails) error
	GetDeputyAnnualInvoiceBillingTypes(ctx sirius.Context) ([]sirius.DeputyAnnualBillingInvoiceTypes, error)
}

type managePaDeputyImportantInformationVars struct {
	Path                      string
	XSRFToken                 string
	DeputyDetails          	sirius.DeputyDetails
	Error                     string
	Errors                    sirius.ValidationErrors
	DeputyId                  int
	AnnualBillingInvoiceTypes []sirius.DeputyAnnualBillingInvoiceTypes
}

func renderTemplateForPaImportantInformation(client ManagePaDeputyImportantInformation, defaultPATeam int, tmpl Template) Handler {
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

		switch r.Method {
		case http.MethodGet:

			vars := managePaDeputyImportantInformationVars{
				Path:                      r.URL.Path,
				XSRFToken:                 ctx.XSRFToken,
				DeputyId:                  deputyId,
				DeputyDetails:          	deputyDetails,
				AnnualBillingInvoiceTypes: annualBillingInvoiceTypes,
			}
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:

			importantInfoForm := sirius.ImportantPaInformationDetails{
				Complaints:                r.PostFormValue("complaints"),
				PanelDeputy:               true,
				AnnualBillingInvoice:      r.PostFormValue("annual-billing"),
				OtherImportantInformation: r.PostFormValue("other-info-note"),
			}

			err = client.UpdateImportantInformation(ctx, deputyId, importantInfoForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				//verr.Errors = renameUpdateAdditionalInformationValidationErrorMessages(verr.Errors)

				vars := managePaDeputyImportantInformationVars{
					Path:                      r.URL.Path,
					XSRFToken:                 ctx.XSRFToken,
					DeputyId:                  deputyId,
					DeputyDetails:          deputyDetails,
					Errors:                    verr.Errors,
					AnnualBillingInvoiceTypes: annualBillingInvoiceTypes,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/deputy/%d?success=importantInformation", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

//func renameUpdateAdditionalInformationValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
//	errorCollection := sirius.ValidationErrors{}
//
//	for fieldName, value := range siriusError {
//		for errorType, errorMessage := range value {
//			err := make(map[string]string)
//			if fieldName == "otherImportantInformation" && errorType == "stringLengthTooLong" {
//				err[errorType] = "The other important information must be 1000 characters or fewer"
//				errorCollection["otherImportantInformation"] = err
//			} else {
//				err[errorType] = errorMessage
//				errorCollection[fieldName] = err
//			}
//		}
//	}
//	return errorCollection
//}
