package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type GetAssurancesClient interface {
	GetAssurances(ctx sirius.Context, deputyId int) ([]model.Assurance, error)
}

type AssurancesVars struct {
	AddVisitDisabled bool
	SuccessMessage   string
	Assurances       []model.Assurance
	ErrorMessage     string
	AppVars
}

func renderTemplateForAssurances(client GetAssurancesClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "addAssurance":
			successMessage = "Assurance process updated"
		case "manageVisit":
			successMessage = "Assurance visit updated"
		case "managePDR":
			successMessage = "PDR updated"
		default:
			successMessage = ""
		}

		assurances, err := client.GetAssurances(ctx, app.DeputyId())
		if err != nil {
			return err
		}

		vars := AssurancesVars{
			SuccessMessage: successMessage,
			Assurances:     assurances,
			AppVars:        app,
		}

		vars.AddVisitDisabled, vars.ErrorMessage = isAddVisitDisabled(assurances)

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func isAddVisitDisabled(assurances []model.Assurance) (bool, string) {
	if len(assurances) > 0 {
		if assurances[0].Type.Label == "PDR" {
			if assurances[0].PdrOutcome.Handle == "NOT_RECEIVED" || (assurances[0].ReportReviewDate != "") {
				return false, ""
			}
			return true, "You cannot add anything until the current assurance process has a review date or is marked as 'Not received'"
		}
		if (assurances[0].ReportReviewDate != "" && assurances[0].ReportMarkedAs.Label != "") || assurances[0].VisitOutcome.Label == "Cancelled" {
			return false, ""
		}
		return true, "You cannot add anything until the current assurance process has a review date and RAG status or is cancelled"
	}
	return false, ""
}
