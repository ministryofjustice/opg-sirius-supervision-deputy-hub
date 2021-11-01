package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ChangeDeputyECMInformation interface {
	GetDeputyDetails(sirius.Context, int) (sirius.DeputyDetails, error)
	ChangeECM(sirius.Context, sirius.DeputyDetails) error
}

type changeDeputyECMHubVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	Success       bool
}

func renderTemplateForChangeDeputyECMHub(client ChangeDeputyECMInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		deputyDetails, err := client.GetDeputyDetails(ctx, deputyId)

		switch r.Method {
		case http.MethodGet:

			if err != nil {
				return err
			}

			vars := changeDeputyECMHubVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			fmt.Println("in server post method")

			fmt.Println("new ecm post value is ")
			fmt.Println(r.PostFormValue("new-ecm"))

			changeDeputyECMForm := sirius.DeputyDetails{
				ID:                               deputyId,
				OrganisationName:                 deputyDetails.OrganisationName,
				OrganisationTeamOrDepartmentName: r.PostFormValue("new-ecm"),
				Email:                            deputyDetails.Email,
				PhoneNumber:                      deputyDetails.PhoneNumber,
				AddressLine1:                     deputyDetails.AddressLine1,
				AddressLine2:                     deputyDetails.AddressLine2,
				AddressLine3:                     deputyDetails.AddressLine3,
				Town:                             deputyDetails.Town,
				County:                           deputyDetails.County,
				Postcode:                         deputyDetails.Postcode,
			}

			err := client.ChangeECM(ctx, changeDeputyECMForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameEditDeputyValidationErrorMessages(verr.Errors)

				vars := changeDeputyECMHubVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					DeputyDetails: deputyDetails,
					Errors:        verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/deputy/%d/?success=true", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

// func renameEditDeputyValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
// 	errorCollection := sirius.ValidationErrors{}

// 	for fieldName, value := range siriusError {
// 		for errorType, errorMessage := range value {
// 			err := make(map[string]string)

// 			if fieldName == "organisationTeamOrDepartmentName" && errorType == "stringLengthTooLong" {
// 				err[errorType] = "The team or department must be 255 characters or fewer"
// 				errorCollection["organisationTeamOrDepartmentName"] = err
// 			}
// 		}
// 	}
// 	return errorCollection
// }
