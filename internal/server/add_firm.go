package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
)

type addFirmVars struct {
	AppVars
}

type AddFirmHandler struct {
	router
}

func (h *AddFirmHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	v.PageName = "Create new firm"

	vars := addFirmVars{
		AppVars: v,
	}

	switch r.Method {
	case http.MethodGet:
		return h.execute(w, r, vars)

	case http.MethodPost:

		addFirmDetailForm := sirius.FirmDetails{
			FirmName:     r.PostFormValue("name"),
			AddressLine1: r.PostFormValue("address-line-1"),
			AddressLine2: r.PostFormValue("address-line-2"),
			AddressLine3: r.PostFormValue("address-line-3"),
			Town:         r.PostFormValue("town"),
			County:       r.PostFormValue("county"),
			Postcode:     r.PostFormValue("postcode"),
			PhoneNumber:  r.PostFormValue("telephone"),
			Email:        r.PostFormValue("email"),
		}

		firmId, err := h.Client().AddFirmDetails(ctx, addFirmDetailForm)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)
			return h.execute(w, r, vars)
		}
		if err != nil {
			return err
		}

		assignDeputyToFirmErr := h.Client().AssignDeputyToFirm(ctx, v.DeputyId(), firmId)
		if assignDeputyToFirmErr != nil {
			return assignDeputyToFirmErr
		}

		return Redirect(fmt.Sprintf("/%d?success=newFirm", v.DeputyId()))
	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
