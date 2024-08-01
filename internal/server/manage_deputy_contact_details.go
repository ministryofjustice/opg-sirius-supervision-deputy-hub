package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
)

type EditContact interface {
	UpdateDeputyContactDetails(sirius.Context, int, sirius.DeputyContactDetails) error
}

type editContactVars struct {
	AppVars
}

type EditContactHandler struct {
	router
}

func (h *EditContactHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	v.PageName = "Manage deputy contact details"

	vars := editContactVars{AppVars: v}

	switch r.Method {
	case http.MethodGet:
		return h.execute(w, r, vars, vars.AppVars)

	case http.MethodPost:
		form := sirius.DeputyContactDetails{
			DeputySubType:    v.DeputyDetails.DeputySubType.SubType,
			DeputyFirstName:  r.PostFormValue("deputy-first-name"),
			DeputySurname:    r.PostFormValue("deputy-last-name"),
			OrganisationName: r.PostFormValue("organisation-name"),
			AddressLine1:     r.PostFormValue("address-line-1"),
			AddressLine2:     r.PostFormValue("address-line-2"),
			AddressLine3:     r.PostFormValue("address-line-3"),
			Town:             r.PostFormValue("town"),
			County:           r.PostFormValue("county"),
			Postcode:         r.PostFormValue("postcode"),
			PhoneNumber:      r.PostFormValue("telephone"),
			Email:            r.PostFormValue("email"),
		}

		err := h.Client().UpdateDeputyContactDetails(ctx, v.DeputyId(), form)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)
			return h.execute(w, r, vars, vars.AppVars)
		} else if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d?success=deputyDetails", v.DeputyId()))
	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
