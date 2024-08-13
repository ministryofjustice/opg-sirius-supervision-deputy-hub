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
	AddGcmIssue(ctx sirius.Context, caseRecNumber, receivedDate, notes string, gcmIssueType model.RefData, deputyId int) error
}

type AddGcmIssueVars struct {
	AppVars
	GcmIssueTypes []model.RefData
	CaseRecNumber string
	Client        sirius.ClientWithOrderDeputy
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
			var receivedDate = r.PostFormValue("received_date")
			var gcmIssueType = r.PostFormValue("issue-type")
			var notes = r.PostFormValue("notes")

			vars.CaseRecNumber = caseNumber

			vars.Errors = sirius.ValidationErrors{}
			if caseNumber == "" {
				vars.Errors["caseRecNumber"] = map[string]string{"": "Enter a case number"}
				vars.Errors = util.RenameErrors(vars.Errors)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if gcmIssueType != "" && receivedDate != "" && notes != "" {
				//replace this with setting hidden input and then checking if thats null
				fmt.Println("final submit")

				fmt.Println("string")
				fmt.Println(gcmIssueType)

				gcmIssue := getRefDataForGcmIssueType(gcmIssueType, vars.GcmIssueTypes)
				fmt.Println("returned ref data")
				fmt.Println(gcmIssue)

				err := client.AddGcmIssue(ctx, caseNumber, receivedDate, notes, gcmIssue, app.DeputyId())
				if err != nil {
					return err
				}

				return Redirect(fmt.Sprintf("/%d/gcm-issues?success=addGcmIssue&%s", app.DeputyId(), caseNumber))
			} else {
				fmt.Println("first submit")
				client, err := client.GetDeputyClient(ctx, caseNumber)

				if verr, ok := err.(sirius.ValidationError); ok {
					fmt.Println("vars errors")
					vars.Errors = util.RenameErrors(verr.Errors)
					//return tmpl.ExecuteTemplate(w, "page", vars)
				}

				if err != nil {
					return err
				}

				linked := checkIfClientLinkedToDeputy(client, app.DeputyId())
				if linked == true {
					vars.Client = client
					return tmpl.ExecuteTemplate(w, "page", vars)

				} else {
					vars.Errors["caseRecNumber"] = map[string]string{"": "Case client not linked to deputy"}
					vars.Errors = util.RenameErrors(vars.Errors)
					return tmpl.ExecuteTemplate(w, "page", vars)
				}
			}

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
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

func getRefDataForGcmIssueType(issueLabelGiven string, refData []model.RefData) model.RefData {
	for i := 0; i < len(refData); {
		if refData[i].Label == issueLabelGiven {
			return refData[i]
		}
		i++
	}
	return model.RefData{}
}
