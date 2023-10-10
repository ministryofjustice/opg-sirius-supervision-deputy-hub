package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type AddAssuranceClient interface {
	AddAssurance(ctx sirius.Context, assuranceType string, requestedDate string, userId, deputyId int) error
}

type AddAssuranceVars struct {
	AppVars
}

func renderTemplateForAddAssurance(client AddAssuranceClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		vars := AddAssuranceVars{
			AppVars: app,
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

			err := client.AddAssurance(ctx, assuranceType, requestedDate, app.UserDetails.ID, app.DeputyId())

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d/assurance-visits?success=addAssuranceVisit", app.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
