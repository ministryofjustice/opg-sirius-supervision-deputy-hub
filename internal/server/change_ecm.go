package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ChangeECMInformation interface {
	GetDeputyDetails(sirius.Context, int, int) (sirius.DeputyDetails, error)
	GetPaDeputyTeamMembers(sirius.Context, int) ([]sirius.TeamMember, error)
	ChangeECM(sirius.Context, sirius.ExecutiveCaseManagerOutgoing, sirius.DeputyDetails) error
}

type changeECMHubVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	EcmTeamDetails []sirius.TeamMember
	Error         string
	Errors        sirius.ValidationErrors
	Success       bool
	DefaultPaTeam int
}

func renderTemplateForChangeECM(client ChangeECMInformation, defaultPATeam int, tmpl Template) Handler {
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

			ecmTeamDetails, err := client.GetPaDeputyTeamMembers(ctx, defaultPATeam)

			if err != nil {
				return err
			}

			vars := changeECMHubVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
				DefaultPaTeam: defaultPATeam,
				EcmTeamDetails: ecmTeamDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:

			EcmIdValue, err := strconv.Atoi(r.PostFormValue("new-ecm"))

			if err != nil {
				return err
			}

			changeECMForm := sirius.ExecutiveCaseManagerOutgoing{EcmId: EcmIdValue}

			vars := changeECMHubVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
				DefaultPaTeam: defaultPATeam,
			}

			//error := validateInput(changeECMForm.EcmId)
			//if error != nil {
			//	vars.Errors = error
			//	return tmpl.ExecuteTemplate(w, "page", vars)
			//}

			err = client.ChangeECM(ctx, changeECMForm, deputyDetails)

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

//func validateInput(newECM int) sirius.ValidationErrors {
//	if newECM < 1 {
//		newError := sirius.ValidationErrors{
//			"Change ECM": {"": "Select an executive case manager"},
//		}
//		return newError
//	}
//	return nil
//}
