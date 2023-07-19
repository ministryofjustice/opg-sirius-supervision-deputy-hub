package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type ManageContact interface {
	GetContactById(ctx sirius.Context, deputyId int, contactId int) (sirius.Contact, error)
}

type ManageContactVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	Contact       sirius.Contact
	ErrorNote     string
}

func renderTemplateForManageContact(client ManageContact, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		contactId, _ := strconv.Atoi(routeVars["contactId"])

		vars := ManageContactVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
		}

		switch r.Method {
		case http.MethodGet:
			contact, err := client.GetContactById(ctx, deputyId, contactId)
			if err != nil {
				return err
			}

			vars.Contact = contact
			fmt.Println(contact)
			return tmpl.ExecuteTemplate(w, "page", vars)
		//
		//case http.MethodPost:
		//	var err error
		//	user, err := client.GetUserDetails(ctx)
		//	if err != nil {
		//		return err
		//	}
		//
		//	reportReviewDate := r.PostFormValue("report-review-date")
		//	reviewedBy := 0
		//	if reportReviewDate != "" {
		//		reviewedBy = user.ID
		//	}
		//
		//	pdrOutcome := ""
		//	if r.PostFormValue("pdr-outcome") == "Not received" {
		//		pdrOutcome = "NOT_RECEIVED"
		//	} else if r.PostFormValue("pdr-outcome") == "Received" {
		//		pdrOutcome = "RECEIVED"
		//	}
		//
		//	manageAssuranceVisitForm := sirius.AssuranceVisitDetails{
		//		CommissionedDate:    r.PostFormValue("commissioned-date"),
		//		VisitorAllocated:    r.PostFormValue("visitor-allocated"),
		//		ReportDueDate:       r.PostFormValue("report-due-date"),
		//		ReportReceivedDate:  r.PostFormValue("report-received-date"),
		//		VisitOutcome:        r.PostFormValue("visit-outcome"),
		//		PdrOutcome:          pdrOutcome,
		//		ReportReviewDate:    reportReviewDate,
		//		VisitReportMarkedAs: r.PostFormValue("visit-report-marked-as"),
		//		ReviewedBy:          reviewedBy,
		//		Note:                strings.TrimSpace(r.PostFormValue("note")),
		//	}
		//
		//	err = client.UpdateAssuranceVisit(ctx, manageAssuranceVisitForm, deputyId, visitId)
		//
		//	if verr, ok := err.(sirius.ValidationError); ok {
		//		vars := ManageAssuranceVisitVars{
		//			Path:                r.URL.Path,
		//			XSRFToken:           ctx.XSRFToken,
		//			Errors:              verr.Errors,
		//			VisitRagRatingTypes: vars.VisitRagRatingTypes,
		//			VisitOutcomeTypes:   vars.VisitOutcomeTypes,
		//			Visitors:            visitors,
		//			ErrorNote:           r.PostFormValue("note"),
		//			DeputyDetails:       deputyDetails,
		//		}
		//		return tmpl.ExecuteTemplate(w, "page", vars)
		//	}
		//
		//	success := "manageAssuranceVisit"
		//	if visit.AssuranceType.Handle == "PDR" {
		//		success = "managePDR"
		//	}
		//
		//	return Redirect(fmt.Sprintf("/%d/assurance-visits?success=%s", deputyId, success))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
