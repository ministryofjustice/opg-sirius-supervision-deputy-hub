package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type FirmInformation interface {
	AddFirmDetails(sirius.Context, sirius.FirmDetails) (int, error)
	AssignDeputyToFirm(sirius.Context, int, int) error
	GetDeputyDetails(sirius.Context, int, int) (sirius.DeputyDetails, error)
}

type addFirmVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	Success       bool
	DeputyId      int
}

func renderTemplateForAddFirm(client FirmInformation, defaultPATeam int, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		deputyDetails, err := client.GetDeputyDetails(ctx, defaultPATeam, deputyId)

		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:
			vars := addFirmVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyId:      deputyId,
				DeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:

			addFirmDetailForm := sirius.FirmDetails{
				FirmName:     r.PostFormValue("name"),
				AddressLine1: r.PostFormValue("address-line-1"),
				AddressLine2: r.PostFormValue("address-line-2"),
				AddressLine3: r.PostFormValue("address-line-3"),
				Town:         r.PostFormValue("town"),
				County:       r.PostFormValue("county"),
				Postcode:     r.PostFormValue("postcode"),
				PhoneNumber:  r.PostFormValue("telephone"),
				Email:        r.PostFormValue("email"),
			}

			firmId, err := client.AddFirmDetails(ctx, addFirmDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameAddFirmValidationErrorMessages(verr.Errors)
				vars := addFirmVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			assignDeputyToFirmErr := client.AssignDeputyToFirm(ctx, deputyId, firmId)
			if assignDeputyToFirmErr != nil {
				return assignDeputyToFirmErr
			}

			return Redirect(fmt.Sprintf("/%d", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func renameAddFirmValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	errorCollection := sirius.ValidationErrors{}
	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)
			if fieldName == "firmName" && errorType == "stringLengthTooLong" {
				err[errorType] = "The firm name must be 255 characters or fewer"
				errorCollection["firmName"] = err
			} else if fieldName == "firmName" && errorType == "isEmpty" {
				err[errorType] = "The firm name is required and can't be empty"
				errorCollection["firmName"] = err
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
			} else if fieldName == "phoneNumber" && errorType == "stringLengthTooLong" {
				err[errorType] = "The telephone number must be 255 characters or fewer"
				errorCollection["phoneNumber"] = err
			} else if fieldName == "email" && errorType == "stringLengthTooLong" {
				err[errorType] = "The email must be 255 characters or fewer"
				errorCollection["email"] = err
			} else {
				err[errorType] = errorMessage
				errorCollection[fieldName] = err
			}
		}
	}
	return errorCollection
}
