package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type EditDeputyHubInformation interface {
	EditDeputyDetails(sirius.Context, sirius.DeputyDetails) error
}

type editDeputyHubVars struct {
	AppVars
}

func renderTemplateForEditDeputyHub(client EditDeputyHubInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		vars := editDeputyHubVars{
			AppVars: app,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			editDeputyDetailForm := sirius.DeputyDetails{
				ID:                               app.DeputyId(),
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

			err := client.EditDeputyDetails(ctx, editDeputyDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=teamDetails", app.DeputyId()))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
