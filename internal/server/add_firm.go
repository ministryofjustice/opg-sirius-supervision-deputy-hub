package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
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
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		vars := addFirmVars{
			AppVars: appVars,
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
				vars.Errors = verr.Errors
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			assignDeputyToFirmErr := client.AssignDeputyToFirm(ctx, appVars.DeputyId(), firmId)
			if assignDeputyToFirmErr != nil {
				return assignDeputyToFirmErr
			}

			return Redirect(fmt.Sprintf("/%d?success=newFirm", appVars.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
