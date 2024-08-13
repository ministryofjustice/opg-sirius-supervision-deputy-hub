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
	GetDeputyClient(ctx sirius.Context, caseRecNumber string) (sirius.ClientWithOrderDeputy, error)
	AddGcmIssue(ctx sirius.Context, caseRecNumber, notes string, gcmIssueType model.RefData, deputyId int) error
}

type AddGcmIssueVars struct {
	AppVars
	GcmIssueTypes  []model.RefData
	CaseRecNumber  string
	Client         sirius.ClientWithOrderDeputy
	HasFoundClient string
	GcmIssueType   model.RefData
	GcmIssueLabel  string
	Notes          string
}

func renderTemplateForAddGcmIssue(client AddGcmIssue, tmpl Template) Handler {
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
			var caseNumber = r.PostFormValue("case_number")
			var gcmIssueType = r.PostFormValue("issue-type")
			var notes = r.PostFormValue("notes")
			var searchForClient = r.PostFormValue("search-for-client")
			var submitForm = r.PostFormValue("submit-form")

			vars.CaseRecNumber = caseNumber
			vars.Notes = notes

			//	they are looking for client
			if caseNumber == "" {
				vars.Errors = sirius.ValidationErrors{}
				vars.Errors["caseRecNumber"] = map[string]string{"": "Enter a case number"}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			siriusClient, err := client.GetDeputyClient(ctx, caseNumber)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if err != nil {
				return err
			}

			linked := checkIfClientLinkedToDeputy(siriusClient, app.DeputyId())
			if linked == true {
				vars.Client = siriusClient
				if searchForClient == "search-for-client" {
					return tmpl.ExecuteTemplate(w, "page", vars)
				}
			} else {
				vars.Errors["caseRecNumber"] = map[string]string{"": "Case number does not belong to this deputy"}
				vars.Errors = util.RenameErrors(vars.Errors)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if vars.Client.ClientId != 0 && submitForm == "submit-form" {
				//second submit to add the issue
				gcmIssue, _ := getRefDataForGcmIssueType(gcmIssueType, vars.GcmIssueTypes)
				err := client.AddGcmIssue(ctx, caseNumber, notes, gcmIssue, app.DeputyId())

				if verr, ok := err.(sirius.ValidationError); ok {
					vars.Client = siriusClient
					vars.CaseRecNumber = caseNumber
					vars.GcmIssueType = gcmIssue
					vars.Notes = notes
					vars.Errors = util.RenameErrors(verr.Errors)
					return tmpl.ExecuteTemplate(w, "page", vars)
				}

				if err != nil {
					fmt.Println("normal error")
					return err
				}

				return Redirect(fmt.Sprintf("/%d/gcm-issues?success=addGcmIssue&%s", app.DeputyId(), caseNumber))
			}

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
		return StatusError(http.StatusMethodNotAllowed)
	}
}

func checkIfClientLinkedToDeputy(client sirius.ClientWithOrderDeputy, deputyId int) bool {
	for i := 0; i < len(client.Cases); {
		deputiesForOrder := client.Cases[i].Deputies
		for j := 0; j < len(deputiesForOrder); {
			if deputiesForOrder[j].Deputy.Id == deputyId {
				return true
			}
			j++
		}
		i++
	}
	return false
}

func getRefDataForGcmIssueType(issueLabelGiven string, refData []model.RefData) (model.RefData, sirius.ValidationErrors) {
	for i := 0; i < len(refData); {
		if refData[i].Label == issueLabelGiven {
			return refData[i], nil
		}
		i++
	}
	return model.RefData{},
		sirius.ValidationErrors{
			"caseRecNumber": map[string]string{"invalid": "Select a valid issue type"},
		}
}
