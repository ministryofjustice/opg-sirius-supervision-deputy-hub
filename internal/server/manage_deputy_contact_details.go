package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyContactDetailsInformation interface {
	UpdateDeputyContactDetails(sirius.Context, int, sirius.DeputyContactDetails) error
}

type manageDeputyContactDetailsVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	DeputyId      int
}

func renderTemplateForManageDeputyContactDetails(client DeputyContactDetailsInformation, defaultPATeam int, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, dd sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		deputyDetails := dd

		switch r.Method {
		case http.MethodGet:

			vars := manageDeputyContactDetailsVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyId:      deputyId,
				DeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			form := sirius.DeputyContactDetails{
				DeputySubType:    deputyDetails.DeputySubType.SubType,
				DeputyFirstName:  r.PostFormValue("deputy-first-name"),
				DeputySurname:    r.PostFormValue("deputy-last-name"),
				OrganisationName: r.PostFormValue("organisation-name"),
				AddressLine1:     r.PostFormValue("address-line-1"),
				AddressLine2:     r.PostFormValue("address-line-2"),
				AddressLine3:     r.PostFormValue("address-line-3"),
				Town:             r.PostFormValue("town"),
				County:           r.PostFormValue("county"),
				Postcode:         r.PostFormValue("postcode"),
				PhoneNumber:      r.PostFormValue("telephone"),
				Email:            r.PostFormValue("email"),
			}

			err := client.UpdateDeputyContactDetails(ctx, deputyId, form)

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameManageDeputyContactDetailsValidationErrorMessages(verr.Errors)
				vars := manageDeputyContactDetailsVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					DeputyId:      deputyId,
					DeputyDetails: deputyDetails,
					Errors:        verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=deputyDetails", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func renameManageDeputyContactDetailsValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	errorCollection := sirius.ValidationErrors{}

	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)

			if fieldName == "firstname" && errorType == "stringLengthTooLong" {
				err[errorType] = "The deputy first name must be 255 characters or fewer"
				errorCollection["firstname"] = err
			} else if fieldName == "firstname" && errorType == "isEmpty" {
				err[errorType] = "The deputy first name is required and can't be empty"
				errorCollection["firstname"] = err
			} else if fieldName == "surname" && errorType == "stringLengthTooLong" {
				err[errorType] = "The deputy surname must be 255 characters or fewer"
				errorCollection["surname"] = err
			} else if fieldName == "surname" && errorType == "isEmpty" {
				err[errorType] = "The deputy surname is required and can't be empty"
				errorCollection["surname"] = err
			} else if fieldName == "organisationName" && errorType == "stringLengthTooLong" {
				err[errorType] = "The organisation name must be 255 characters or fewer"
				errorCollection["organisationName"] = err
			} else if fieldName == "organisationName" && errorType == "isEmpty" {
				err[errorType] = "The organisation name is required and can't be empty"
				errorCollection["organisationName"] = err
			} else if fieldName == "workPhoneNumber" && errorType == "stringLengthTooLong" {
				err[errorType] = "The telephone number must be 255 characters or fewer"
				errorCollection["workPhoneNumber"] = err
			} else if fieldName == "email" && errorType == "stringLengthTooLong" {
				err[errorType] = "The email must be 255 characters or fewer"
				errorCollection["email"] = err
			} else if fieldName == "addressLine1" && errorType == "stringLengthTooLong" {
				err[errorType] = "The building or street must be 255 characters or fewer"
				errorCollection["addressLine1"] = err
			} else if fieldName == "addressLine2" && errorType == "stringLengthTooLong" {
				err[errorType] = "Address line 2 must be 255 characters or fewer"
				errorCollection["addressLine2"] = err
			} else if fieldName == "addressLine3" && errorType == "stringLengthTooLong" {
				err[errorType] = "Address line 3 must be 255 characters or fewer"
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
