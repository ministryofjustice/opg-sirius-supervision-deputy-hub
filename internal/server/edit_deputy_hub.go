package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type EditDeputyHubInformation interface {
	EditDeputyDetails(sirius.Context, sirius.DeputyDetails) error
}

type editDeputyHubVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
}

func renderTemplateForEditDeputyHub(client EditDeputyHubInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:

			vars := editDeputyHubVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			editDeputyDetailForm := sirius.DeputyDetails{
				ID:                               deputyId,
				OrganisationName:                 r.PostFormValue("deputy-name"),
				OrganisationTeamOrDepartmentName: r.PostFormValue("organisationTeamOrDepartmentName"),
				Email:                            r.PostFormValue("email"),
				PhoneNumber:                      r.PostFormValue("telephone"),
				AddressLine1:                     r.PostFormValue("address-line-1"),
				AddressLine2:                     r.PostFormValue("address-line-2"),
				AddressLine3:                     r.PostFormValue("address-line-3"),
				Town:                             r.PostFormValue("town"),
				County:                           r.PostFormValue("county"),
				Postcode:                         r.PostFormValue("postcode"),
			}
			fmt.Println("deputy id")
			fmt.Println(deputyId)

			err := client.EditDeputyDetails(ctx, editDeputyDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := editDeputyHubVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					DeputyDetails: deputyDetails,
					Errors:        verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=teamDetails", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
