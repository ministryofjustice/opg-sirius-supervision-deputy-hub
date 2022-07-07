package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyContactDetailsInformation interface {
	UpdateDeputyContactDetails(sirius.Context, int, sirius.DeputyContactDetails) error
}

type manageDeputyContactDetailsVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	DeputyId      int
}

func renderTemplateForManageDeputyContactDetails(client DeputyContactDetailsInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:

			vars := manageDeputyContactDetailsVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyId:      deputyId,
				DeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			form := sirius.DeputyContactDetails{
				DeputySubType:    deputyDetails.DeputySubType.SubType,
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
				vars := manageDeputyContactDetailsVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					DeputyId:      deputyId,
					DeputyDetails: deputyDetails,
					Errors:        verr.Errors,
				}
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
