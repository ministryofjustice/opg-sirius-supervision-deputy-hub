package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
)

type AddAssuranceVisit interface {
	AddAssuranceVisit(ctx sirius.Context, assuranceType string, requestedDate string, userId, deputyId int) error
}

type AddAssuranceVisitVars struct {
	AppVars
}

func renderTemplateForAddAssuranceVisit(client AddAssuranceVisit, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		app.PageName = "Add assurance visit"

		vars := AddAssuranceVisitVars{
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

			vars.Errors = util.RenameErrors(vars.Errors)

			if len(vars.Errors) > 0 {
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			err := client.AddAssuranceVisit(ctx, assuranceType, requestedDate, app.UserDetails.ID, app.DeputyId())

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
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
