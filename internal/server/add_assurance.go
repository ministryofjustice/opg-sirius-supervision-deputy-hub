package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
)

type AddAssuranceVars struct {
	AppVars
}

type AddAssuranceHandler struct {
	router
}

func (h *AddAssuranceHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	v.PageName = "Add assurance visit"

	vars := AddAssuranceVars{
		AppVars: v,
	}

	switch r.Method {
	case http.MethodGet:
		return h.execute(w, r, vars)

	case http.MethodPost:
		var assuranceType = r.PostFormValue("assurance-type")
		var requestedDate = r.PostFormValue("requested-date")

		vars.Errors = sirius.ValidationErrors{}

		if assuranceType == "" {
			vars.Errors["assurance-type"] = map[string]string{"": "Select an assurance type"}
		}

		if requestedDate == "" {
			vars.Errors["requested-date"] = map[string]string{"": "Enter a requested date"}
		}

		vars.Errors = util.RenameErrors(vars.Errors)

		if len(vars.Errors) > 0 {
			return h.execute(w, r, vars)
		}

		err := h.Client().AddAssurance(ctx, assuranceType, requestedDate, v.UserDetails.ID, v.DeputyId())

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)
			return h.execute(w, r, vars)
		}
		if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d/assurances?success=addAssurance", v.DeputyId()))
	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
