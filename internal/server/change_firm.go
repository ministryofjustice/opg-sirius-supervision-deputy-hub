package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
	"strconv"
)

type changeFirmVars struct {
	Firms          []sirius.FirmForList
	Success        bool
	SuccessMessage string
	AppVars
}

type EditFirmHandler struct {
	router
}

func (h *EditFirmHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	firms, err := h.Client().GetFirms(ctx)
	if err != nil {
		return err
	}

	v.PageName = "Change firm"

	vars := changeFirmVars{
		Firms:   firms,
		AppVars: v,
	}

	switch r.Method {
	case http.MethodGet:
		return h.execute(w, r, vars)

	case http.MethodPost:
		var vars changeFirmVars
		newFirm := r.PostFormValue("select-firm")
		AssignToExistingFirmStringIdValue := r.PostFormValue("select-existing-firm")

		if newFirm == "new-firm" {
			return Redirect(fmt.Sprintf("/%d/firm/add", v.DeputyId()))
		}

		AssignToFirmId := 0
		if AssignToExistingFirmStringIdValue != "" {
			AssignToFirmId, err = strconv.Atoi(AssignToExistingFirmStringIdValue)
			if err != nil {
				return err
			}
		}

		assignDeputyToFirmErr := h.Client().AssignDeputyToFirm(ctx, v.DeputyId(), AssignToFirmId)

		if verr, ok := assignDeputyToFirmErr.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)
			return h.execute(w, r, vars)
		}

		if assignDeputyToFirmErr != nil {
			return assignDeputyToFirmErr
		}

		return Redirect(fmt.Sprintf("/%d?success=firm", v.DeputyId()))

	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
