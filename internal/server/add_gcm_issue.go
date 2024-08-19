package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type AddGcmIssue interface {
	GetGCMIssueTypes(ctx sirius.Context) ([]model.RefData, error)
	GetDeputyClient(ctx sirius.Context, caseRecNumber string, deputyId int) (sirius.DeputyClient, error)
	AddGcmIssue(ctx sirius.Context, caseRecNumber, notes string, gcmIssueType model.RefData, deputyId int) error
}

type AddGcmIssueVars struct {
	AppVars
	GcmIssueTypes  []model.RefData
	CaseRecNumber  string
	Client         sirius.DeputyClient
	HasFoundClient string
	GcmIssueType   model.RefData
	Notes          string
}

func renderTemplateForAddGcmIssue(client AddGcmIssue, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		app.PageName = "Add a GCM issue"
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
			var caseRecNumber = r.PostFormValue("case-number")
			var gcmIssueType = r.PostFormValue("issue-type")
			var notes = r.PostFormValue("notes")
			var searchForClient = r.PostFormValue("search-for-client")
			var submitForm = r.PostFormValue("submit-form")

			vars.CaseRecNumber = caseRecNumber
			vars.Notes = notes

			if caseRecNumber == "" {
				vars.Errors = sirius.ValidationErrors{}
				vars.Errors["client-case-number"] = map[string]string{"": "Enter a case number"}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			siriusClient, err := client.GetDeputyClient(ctx, caseRecNumber, app.DeputyId())

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			vars.Client = siriusClient
			if searchForClient == "search-for-client" {
				//first submit to get the client name from caserec
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if vars.Client.ClientId != 0 && submitForm == "submit-form" {
				//second submit to add the issue
				gcmIssue, _ := getRefDataForGcmIssueType(gcmIssueType, vars.GcmIssueTypes)

				err := client.AddGcmIssue(ctx, caseRecNumber, notes, gcmIssue, app.DeputyId())

				if verr, ok := err.(sirius.ValidationError); ok {
					vars.Client = siriusClient
					vars.CaseRecNumber = caseRecNumber
					vars.GcmIssueType = gcmIssue
					vars.Notes = notes
					vars.Errors = util.RenameErrors(verr.Errors)
					return tmpl.ExecuteTemplate(w, "page", vars)
				}

				if err != nil {
					return err
				}

				return Redirect(fmt.Sprintf("/%d/gcm-issues?success=addGcmIssue&%s", app.DeputyId(), caseRecNumber))
			}

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
		return StatusError(http.StatusMethodNotAllowed)
	}
}

func getRefDataForGcmIssueType(issueHandleGiven string, refData []model.RefData) (model.RefData, sirius.ValidationErrors) {
	for i := 0; i < len(refData); {
		if refData[i].Handle == issueHandleGiven {
			return refData[i], nil
		}
		i++
	}
	return model.RefData{},
		sirius.ValidationErrors{
			"caseRecNumber": map[string]string{"invalid": "Select a valid issue type"},
		}
}
