package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
)

type FirmInformation interface {
	AddFirmDetails(sirius.Context, sirius.FirmDetails) (int, error)
	AssignDeputyToFirm(sirius.Context, int, int) error
}

type addFirmVars struct {
	AppVars
}

func renderTemplateForAddFirm(client FirmInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		app.PageName = "Create new firm"

		vars := addFirmVars{
			AppVars: app,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

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

			firmId, err := client.AddFirmDetails(ctx, addFirmDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			assignDeputyToFirmErr := client.AssignDeputyToFirm(ctx, app.DeputyId(), firmId)
			if assignDeputyToFirmErr != nil {
				return assignDeputyToFirmErr
			}

			return Redirect(fmt.Sprintf("/%d?success=newFirm", app.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
