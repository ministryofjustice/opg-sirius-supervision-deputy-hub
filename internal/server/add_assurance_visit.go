package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type AddAssuranceVisit interface {
	AddAssuranceVisit(ctx sirius.Context, assuranceType string, requestedDate string, userId, deputyId int) error
}

type AddAssuranceVisitVars struct {
	AppVars
}

func renderTemplateForAddAssuranceVisit(client AddAssuranceVisit, tmpl Template) Handler {
	return func(appVars AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		vars := AddAssuranceVisitVars{
			AppVars: appVars,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var assuranceType = r.PostFormValue("assurance-type")
			var requestedDate = r.PostFormValue("requested-date")

			vars.Errors = sirius.ValidationErrors{}

			if assuranceType == "" {
				vars.Errors["assurance-type"] = map[string]string{"": "Select an assurance type"}
			}

			if requestedDate == "" {
				vars.Errors["requested-date"] = map[string]string{"": "Enter a requested date"}
			}

			if len(vars.Errors) > 0 {
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			err := client.AddAssuranceVisit(ctx, assuranceType, requestedDate, appVars.UserDetails.ID, appVars.DeputyId())

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/assurance-visits?success=addAssuranceVisit", appVars.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
