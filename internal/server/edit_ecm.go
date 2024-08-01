package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
	"strconv"
)

type EditDeputyEcm interface {
	GetDeputyTeamMembers(sirius.Context, int, sirius.DeputyDetails) ([]model.TeamMember, error)
	ChangeECM(sirius.Context, sirius.ExecutiveCaseManagerOutgoing, sirius.DeputyDetails) error
}

type editDeputyEcm struct {
	EcmTeamDetails []model.TeamMember
	SuccessMessage string
	AppVars
}

type EditDeputyEcmHandler struct {
	router
}

func (h *EditDeputyEcmHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	deputyId, _ := strconv.Atoi(r.PathValue("deputyId"))

	ecmTeamDetails, err := h.Client().GetDeputyTeamMembers(ctx, v.DefaultPaTeam, v.DeputyDetails)
	if err != nil {
		return err
	}

	v.PageName = "Change Executive Case Manager"

	vars := editDeputyEcm{
		EcmTeamDetails: ecmTeamDetails,
		AppVars:        v,
	}

	switch r.Method {
	case http.MethodGet:
		var successMessage string
		if r.URL.Query().Get("success") == "true" {
			successMessage = "new ecm is" + v.DeputyDetails.ExecutiveCaseManager.EcmName
		}

		vars.SuccessMessage = successMessage
		return h.execute(w, r, vars, vars.AppVars)

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
			return h.execute(w, r, vars, vars.AppVars)
		}

		EcmIdValue, err := strconv.Atoi(EcmIdStringValue)
		if err != nil {
			return err
		}

		changeECMForm := sirius.ExecutiveCaseManagerOutgoing{EcmId: EcmIdValue}

		err = h.Client().ChangeECM(ctx, changeECMForm, v.DeputyDetails)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)

			return h.execute(w, r, vars, vars.AppVars)
		}

		if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d?success=ecm", deputyId))

	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
