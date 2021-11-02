package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ChangeECMInformation interface {
	GetDeputyDetails(sirius.Context, int) (sirius.DeputyDetails, error)
	ChangeECM(sirius.Context, sirius.DeputyDetails, sirius.DeputyDetails) error
}

type changeECMHubVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	Success       bool
	DefaultPaTeam string
}

func renderTemplateForChangeECM(client ChangeECMInformation, defaultPATeam string, tmpl Template) Handler {
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

			vars := changeECMHubVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
				DefaultPaTeam: defaultPATeam,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			changeECMForm := sirius.DeputyDetails{
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

			vars := changeECMHubVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
				DefaultPaTeam: defaultPATeam,
			}

			error := validateInput(changeECMForm.OrganisationTeamOrDepartmentName)
			if error != nil {
				vars.Errors = error
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			err := client.ChangeECM(ctx, changeECMForm, deputyDetails)

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameEditDeputyValidationErrorMessages(verr.Errors)
				vars.Errors = verr.Errors

				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			return Redirect(fmt.Sprintf("/deputy/%d/", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func validateInput(newECM string) sirius.ValidationErrors {
	if len(newECM) < 1 {
		newError := sirius.ValidationErrors{
			"Change ECM": {"": "Select an executive case manager"},
		}
		return newError
	}
	return nil
}
