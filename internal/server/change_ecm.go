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
	ChangeECM(sirius.Context, sirius.ExecutiveCaseManagerOutgoing, sirius.DeputyDetails) (sirius.ExecutiveCaseManager, error)
}

type changeECMHubVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	EcmTeamDetails []sirius.TeamMember
	Error          string
	Errors         sirius.ValidationErrors
	Success        bool
	SuccessMessage string
	DefaultPaTeam  int
	Ecm sirius.ExecutiveCaseManager
}

func renderTemplateForChangeECM(client ChangeECMInformation, defaultPATeam int, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		ecmTeamDetails, err := client.GetPaDeputyTeamMembers(ctx, defaultPATeam)
		if err != nil {
			return err
		}


		switch r.Method {
		case http.MethodGet:
			deputyDetails, err := client.GetDeputyDetails(ctx, defaultPATeam, deputyId)
			if err != nil {
				return err
			}

			hasSuccess := hasSuccessInUrl(r.URL.String(), "/deputy/"+strconv.Itoa(deputyId))
			SuccessMessage := "new ecm is" + deputyDetails.ExecutiveCaseManager.EcmName

			vars := changeECMHubVars{
				Path:           r.URL.Path,
				XSRFToken:      ctx.XSRFToken,
				DeputyDetails:  deputyDetails,
				DefaultPaTeam:  defaultPATeam,
				EcmTeamDetails: ecmTeamDetails,
				Success:        hasSuccess,
				SuccessMessage: SuccessMessage,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			deputyDetails, err := client.GetDeputyDetails(ctx, defaultPATeam, deputyId)
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

			EcmIdStringValue := r.PostFormValue("select-ecm")

			fmt.Println("EcmIdStringValue")
			fmt.Println(EcmIdStringValue)

			if EcmIdStringValue == "" {
				vars.Errors = sirius.ValidationErrors{
					"Change ECM": {"": "Select an executive case manager"},
				}
			}

			fmt.Println("vars.Errors")
			fmt.Println(vars.Errors)

			var newValue string


			if EcmIdStringValue == "" {
				newValue = "0"
			} else {
				newValue = EcmIdStringValue
			}

			fmt.Println("newValue")
			fmt.Println(newValue)

			EcmIdValue, err := strconv.Atoi(newValue);
			if err != nil {
				return err
			}

			changeECMForm := sirius.ExecutiveCaseManagerOutgoing{EcmId: EcmIdValue}

			Ecm, err := client.ChangeECM(ctx, changeECMForm, deputyDetails)

			vars.Ecm = Ecm

			fmt.Println("length of vars errors")
			fmt.Println(len(vars.Errors))

			if len(vars.Errors) >= 1 {
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameEditDeputyValidationErrorMessages(verr.Errors)
				vars.Errors = verr.Errors

				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			return Redirect(fmt.Sprintf("/deputy/%d/?success=ecm", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

