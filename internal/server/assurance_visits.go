package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type AssuranceVisit interface {
	GetAssuranceVisits(ctx sirius.Context, deputyId int) ([]sirius.AssuranceVisits, error)
}

type AssuranceVisitsVars struct {
	AddVisitDisabled bool
	SuccessMessage   string
	AssuranceVisits  []sirius.AssuranceVisits
	ErrorMessage     string
	AppVars
}

func renderTemplateForAssuranceVisits(client AssuranceVisit, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "addAssuranceVisit":
			successMessage = "Assurance process updated"
		case "manageAssuranceVisit":
			successMessage = "Assurance visit updated"
		case "managePDR":
			successMessage = "PDR updated"
		default:
			successMessage = ""
		}

		visits, err := client.GetAssuranceVisits(ctx, app.DeputyId())
		if err != nil {
			return err
		}

		vars := AssuranceVisitsVars{
			SuccessMessage:  successMessage,
			AssuranceVisits: visits,
			AppVars:         app,
		}

		vars.AddVisitDisabled, vars.ErrorMessage = isAddVisitDisabled(visits)

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func isAddVisitDisabled(visits []sirius.AssuranceVisits) (bool, string) {
	if len(visits) > 0 {
		if visits[0].AssuranceType.Label == "PDR" {
			if visits[0].PdrOutcome.Handle == "NOT_RECEIVED" || (visits[0].ReportReviewDate != "") {
				return false, ""
			}
			return true, "You cannot add anything until the current assurance process has a review date or is marked as 'Not received'"
		}
		if (visits[0].ReportReviewDate != "" && visits[0].VisitReportMarkedAs.Label != "") || visits[0].VisitOutcome.Label == "Cancelled" {
			return false, ""
		}
		return true, "You cannot add anything until the current assurance process has a review date and RAG status or is cancelled"
	}
	return false, ""
}
