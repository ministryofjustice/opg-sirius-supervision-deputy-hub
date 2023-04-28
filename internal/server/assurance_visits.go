package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type AssuranceVisit interface {
	GetAssuranceVisits(ctx sirius.Context, deputyId int) ([]sirius.AssuranceVisits, error)
}

type AssuranceVisitsVars struct {
	Path             string
	XSRFToken        string
	DeputyDetails    sirius.DeputyDetails
	Error            string
	AddVisitDisabled bool
	SuccessMessage   string
	AssuranceVisits  []sirius.AssuranceVisits
	ErrorMessage     string
}

func renderTemplateForAssuranceVisits(client AssuranceVisit, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

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

		visits, err := client.GetAssuranceVisits(ctx, deputyId)
		if err != nil {
			return err
		}

		vars := AssuranceVisitsVars{
			Path:            r.URL.Path,
			XSRFToken:       ctx.XSRFToken,
			DeputyDetails:   deputyDetails,
			SuccessMessage:  successMessage,
			AssuranceVisits: visits,
		}
		vars.AddVisitDisabled, vars.ErrorMessage = isAddVisitDisabled(visits)

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func isAddVisitDisabled(visits []sirius.AssuranceVisits) (bool, string) {
	nullDate := sirius.GetNullDate()

	if len(visits) > 0 {
		if visits[0].AssuranceType.Label == "PDR" {
			if visits[0].PdrOutcome.Handle == "NOT_RECEIVED" || (visits[0].ReportReviewDate != nullDate) {
				return false, ""
			}
			return true, "You cannot add anything until the current assurance process has a review date or is marked as 'Not received'"
		}
		if (visits[0].ReportReviewDate != nullDate && visits[0].VisitReportMarkedAs.Label != "") || visits[0].VisitOutcome.Label == "Cancelled" {
			return false, ""
		}
		return true, "You cannot add anything until the current assurance process has a review date and RAG status or is cancelled"
	}
	return false, ""
}
