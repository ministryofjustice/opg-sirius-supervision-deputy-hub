package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ChangeECMInformation interface {
	GetDeputyTeamMembers(sirius.Context, int, sirius.DeputyDetails) ([]model.TeamMember, error)
	ChangeECM(sirius.Context, sirius.ExecutiveCaseManagerOutgoing, sirius.DeputyDetails) error
}

type changeECMHubVars struct {
	EcmTeamDetails []model.TeamMember
	SuccessMessage string
	AppVars
}

func renderTemplateForChangeECM(client ChangeECMInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		deputyId, _ := strconv.Atoi(r.PathValue("id"))

		ecmTeamDetails, err := client.GetDeputyTeamMembers(ctx, app.DefaultPaTeam, app.DeputyDetails)
		if err != nil {
			return err
		}

		app.PageName = "Change Executive Case Manager"

		vars := changeECMHubVars{
			EcmTeamDetails: ecmTeamDetails,
			AppVars:        app,
		}

		switch r.Method {
		case http.MethodGet:
			var successMessage string
			if r.URL.Query().Get("success") == "true" {
				successMessage = "new ecm is" + app.DeputyDetails.ExecutiveCaseManager.EcmName
			}

			vars.SuccessMessage = successMessage
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			if err != nil {
				return err
			}

			EcmIdStringValue := r.PostFormValue("select-ecm")

			if EcmIdStringValue == "" {
				selectECMError := sirius.ValidationErrors{
					"select-ecm": {"": "Select an executive case manager"},
				}

				vars.Errors = util.RenameErrors(selectECMError)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			EcmIdValue, err := strconv.Atoi(EcmIdStringValue)
			if err != nil {
				return err
			}

			changeECMForm := sirius.ExecutiveCaseManagerOutgoing{EcmId: EcmIdValue}

			err = client.ChangeECM(ctx, changeECMForm, app.DeputyDetails)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)

				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=ecm", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
