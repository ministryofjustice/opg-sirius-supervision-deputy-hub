package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ChangeECMInformation interface {
	GetDeputyTeamMembers(sirius.Context, int, sirius.DeputyDetails) ([]sirius.TeamMember, error)
	ChangeECM(sirius.Context, sirius.ExecutiveCaseManagerOutgoing, sirius.DeputyDetails) error
}

type changeECMHubVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	EcmTeamDetails []sirius.TeamMember
	Error          string
	Errors         sirius.ValidationErrors
	SuccessMessage string
}

func renderTemplateForChangeECM(client ChangeECMInformation, defaultPATeam int, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		ecmTeamDetails, err := client.GetDeputyTeamMembers(ctx, defaultPATeam, deputyDetails)
		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:
			var successMessage string
			if r.URL.Query().Get("success") == "true" {
				successMessage = "new ecm is" + deputyDetails.ExecutiveCaseManager.EcmName
			}

			vars := changeECMHubVars{
				Path:           r.URL.Path,
				XSRFToken:      ctx.XSRFToken,
				DeputyDetails:  deputyDetails,
				EcmTeamDetails: ecmTeamDetails,
				SuccessMessage: successMessage,
			}
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			if err != nil {
				return err
			}

			vars := changeECMHubVars{
				Path:           r.URL.Path,
				XSRFToken:      ctx.XSRFToken,
				DeputyDetails:  deputyDetails,
				EcmTeamDetails: ecmTeamDetails,
			}

			EcmIdStringValue := r.PostFormValue("select-ecm")

			if EcmIdStringValue == "" {
				vars.Errors = sirius.ValidationErrors{
					"Change ECM": {"": "Select an executive case manager"},
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			EcmIdValue, err := strconv.Atoi(EcmIdStringValue)
			if err != nil {
				return err
			}

			changeECMForm := sirius.ExecutiveCaseManagerOutgoing{EcmId: EcmIdValue}

			err = client.ChangeECM(ctx, changeECMForm, deputyDetails)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors

				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			return Redirect(fmt.Sprintf("/%d?success=ecm", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
