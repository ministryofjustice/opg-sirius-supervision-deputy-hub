package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
)

type editDeputyTeamVars struct {
	AppVars
}

type EditDeputyTeamHandler struct {
	router
}

func (h *EditDeputyTeamHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	v.PageName = "Manage team details"

	vars := editDeputyTeamVars{
		AppVars: v,
	}

	switch r.Method {
	case http.MethodGet:
		return h.execute(w, r, vars)

	case http.MethodPost:
		editDeputyDetailForm := sirius.DeputyDetails{
			ID:                               v.DeputyId(),
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

		err := h.Client().EditDeputyTeamDetails(ctx, editDeputyDetailForm)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)
			return h.execute(w, r, vars)
		}
		if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d?success=teamDetails", v.DeputyId()))

	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
