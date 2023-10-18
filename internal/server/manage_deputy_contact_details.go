package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
)

type DeputyContactDetailsInformation interface {
	UpdateDeputyContactDetails(sirius.Context, int, sirius.DeputyContactDetails) error
}

type manageDeputyContactDetailsVars struct {
	AppVars
}

func renderTemplateForManageDeputyContactDetails(client DeputyContactDetailsInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		app.PageName = "Manage deputy contact details"

		vars := manageDeputyContactDetailsVars{AppVars: app}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			form := sirius.DeputyContactDetails{
				DeputySubType:    app.DeputyDetails.DeputySubType.SubType,
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

			err := client.UpdateDeputyContactDetails(ctx, app.DeputyId(), form)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=deputyDetails", app.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
