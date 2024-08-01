package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
	"strconv"
)

type DeleteDeputyVars struct {
	SuccessMessage string
	AppVars
}

type DeleteDeputyHandler struct {
	router
}

func (h *DeleteDeputyHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	deputyId, _ := strconv.Atoi(r.PathValue("deputyId"))
	v.PageName = "Delete deputy"
	vars := DeleteDeputyVars{
		AppVars: v,
	}

	switch r.Method {
	case http.MethodGet:
		return h.execute(w, r, vars)
	case http.MethodPost:
		err := h.Client().DeleteDeputy(ctx, deputyId)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)

			w.WriteHeader(http.StatusBadRequest)
			return h.execute(w, r, vars)
		} else if err != nil {
			return err
		}

		successVar := fmt.Sprintf("%s %d has been deleted.", vars.DeputyDetails.DisplayName, vars.DeputyDetails.DeputyNumber)
		return Redirect(fmt.Sprintf("/%d/delete?success=%s", deputyId, successVar))

	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
