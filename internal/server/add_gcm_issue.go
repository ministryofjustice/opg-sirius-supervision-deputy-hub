package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type GetGcmIssue interface {
	GetGCMIssueTypes(ctx sirius.Context) ([]model.RefData, error)
}

type AddGcmIssueVars struct {
	AppVars
	GcmIssueTypes []model.RefData
}

func renderTemplateForAddGcmIssue(client GetGcmIssue, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		app.PageName = "Add GCM Issue"
		vars := AddGcmIssueVars{
			AppVars: app,
		}

		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			gcmIssueTypes, err := client.GetGCMIssueTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}
			vars.GcmIssueTypes = gcmIssueTypes
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			//var assuranceType = r.PostFormValue("assurance-type")
			//var requestedDate = r.PostFormValue("requested-date")
			//
			//vars.Errors = sirius.ValidationErrors{}
			//
			//if assuranceType == "" {
			//	vars.Errors["assurance-type"] = map[string]string{"": "Select an assurance type"}
			//}
			//
			//if requestedDate == "" {
			//	vars.Errors["requested-date"] = map[string]string{"": "Enter a requested date"}
			//}
			//
			//vars.Errors = util.RenameErrors(vars.Errors)
			//
			//if len(vars.Errors) > 0 {
			//	return tmpl.ExecuteTemplate(w, "page", vars)
			//}
			//
			//err := client.AddAssurance(ctx, assuranceType, requestedDate, app.UserDetails.ID, app.DeputyId())
			//
			//if verr, ok := err.(sirius.ValidationError); ok {
			//	vars.Errors = util.RenameErrors(verr.Errors)
			//	return tmpl.ExecuteTemplate(w, "page", vars)
			//}
			//if err != nil {
			//	return err
			//}

			return Redirect(fmt.Sprintf("/%d/assurances?success=addAssurance", app.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
