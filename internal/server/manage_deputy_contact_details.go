package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type DeputyContactDetailsInformation interface {
	UpdateDeputyContactDetails(sirius.Context, int, sirius.DeputyContactDetails) error
}

type manageDeputyContactDetailsVars struct {
	AppVars
}

func renderTemplateForManageDeputyContactDetails(client DeputyContactDetailsInformation, tmpl Template) Handler {
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		vars := manageDeputyContactDetailsVars{AppVars: appVars}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			form := sirius.DeputyContactDetails{
				DeputySubType:    appVars.DeputyDetails.DeputySubType.SubType,
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

			err := client.UpdateDeputyContactDetails(ctx, deputyId, form)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=deputyDetails", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
