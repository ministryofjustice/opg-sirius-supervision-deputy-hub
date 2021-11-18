package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type EditDeputyHubInformation interface {
	GetDeputyDetails(sirius.Context, int, int) (sirius.DeputyDetails, error)
	EditDeputyDetails(sirius.Context, sirius.DeputyDetails) error
}

type editDeputyHubVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	Success       bool
}

func renderTemplateForEditDeputyHub(client EditDeputyHubInformation, defaultPATeam int, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		deputyDetails, err := client.GetDeputyDetails(ctx, defaultPATeam, deputyId)

		switch r.Method {
		case http.MethodGet:

			if err != nil {
				return err
			}

			vars := editDeputyHubVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			editDeputyDetailForm := sirius.DeputyDetails{
				ID:                               deputyId,
				OrganisationName:                 r.PostFormValue("deputy-name"),
				OrganisationTeamOrDepartmentName: r.PostFormValue("organisationTeamOrDepartmentName"),
				Email:                            r.PostFormValue("email"),
				PhoneNumber:                      r.PostFormValue("telephone"),
				AddressLine1:                     r.PostFormValue("address-line-1"),
				AddressLine2:                     r.PostFormValue("address-line-2"),
				AddressLine3:                     r.PostFormValue("address-line-3"),
				Town:                             r.PostFormValue("town"),
				County:                           r.PostFormValue("county"),
				Postcode:                         r.PostFormValue("postcode"),
			}

			err := client.EditDeputyDetails(ctx, editDeputyDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameEditDeputyValidationErrorMessages(verr.Errors)

				vars := editDeputyHubVars{
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

func renameEditDeputyValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	errorCollection := sirius.ValidationErrors{}

	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)

			if fieldName == "organisationName" && errorType == "stringLengthTooLong" {
				err[errorType] = "The deputy name must be 255 characters or fewer"
				errorCollection["organisationName"] = err
			} else if fieldName == "organisationName" && errorType == "isEmpty" {
				err[errorType] = "Enter a deputy name"
				errorCollection["organisationName"] = err
			} else if fieldName == "workPhoneNumber" && errorType == "stringLengthTooLong" {
				err[errorType] = "The telephone number must be 255 characters or fewer"
				errorCollection["workPhoneNumber"] = err
			} else if fieldName == "email" && errorType == "stringLengthTooLong" {
				err[errorType] = "The email number must be 255 characters or fewer"
				errorCollection["email"] = err
			} else if fieldName == "organisationTeamOrDepartmentName" && errorType == "stringLengthTooLong" {
				err[errorType] = "The team or department must be 255 characters or fewer"
				errorCollection["organisationTeamOrDepartmentName"] = err
			} else if fieldName == "addressLine1" && errorType == "stringLengthTooLong" {
				err[errorType] = "The building or street must be 255 characters or fewer"
				errorCollection["addressLine1"] = err
			} else if fieldName == "addressLine2" && errorType == "stringLengthTooLong" {
				err[errorType] = "Address line 2 must be 255 characters or fewer"
				errorCollection["addressLine2"] = err
			} else if fieldName == "addressLine3" && errorType == "stringLengthTooLong" {
				err[errorType] = "AddressLine 3 must be 255 characters or fewer"
				errorCollection["addressLine3"] = err
			} else if fieldName == "town" && errorType == "stringLengthTooLong" {
				err[errorType] = "The town or city must be 255 characters or fewer"
				errorCollection["town"] = err
			} else if fieldName == "county" && errorType == "stringLengthTooLong" {
				err[errorType] = "The county must be 255 characters or fewer"
				errorCollection["county"] = err
			} else if fieldName == "postcode" && errorType == "stringLengthTooLong" {
				err[errorType] = "The postcode must be 255 characters or fewer"
				errorCollection["postcode"] = err
			} else {
				err[errorType] = errorMessage
				errorCollection[fieldName] = err
			}
		}
	}
	return errorCollection
}
