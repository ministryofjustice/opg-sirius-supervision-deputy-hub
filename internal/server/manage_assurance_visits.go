package server

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ManageAssuranceVisit interface {
	UpdateAssuranceVisit(ctx sirius.Context, manageAssuranceVisitForm sirius.AssuranceVisitDetails, deputyId, visitId int) error
	GetVisitors(ctx sirius.Context) (sirius.Visitors, error)
	GetVisitRagRatingTypes(ctx sirius.Context) ([]sirius.VisitRagRatingTypes, error)
	GetVisitOutcomeTypes(ctx sirius.Context) ([]sirius.VisitOutcomeTypes, error)
	GetPdrOutcomeTypes(ctx sirius.Context) ([]sirius.PdrOutcomeTypes, error)
	GetAssuranceVisitById(ctx sirius.Context, deputyId int, visitId int) (sirius.AssuranceVisit, error)
}

type ManageAssuranceVisitVars struct {
	Visitors            sirius.Visitors
	VisitRagRatingTypes []sirius.VisitRagRatingTypes
	VisitOutcomeTypes   []sirius.VisitOutcomeTypes
	PdrOutcomeTypes     []sirius.PdrOutcomeTypes
	Visit               sirius.AssuranceVisit
	ErrorNote           string
	AppVars
}

func renderTemplateForManageAssuranceVisit(client ManageAssuranceVisit, visitTmpl Template, pdrTmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		visitId, _ := strconv.Atoi(routeVars["visitId"])

		vars := ManageAssuranceVisitVars{AppVars: app}
		tmpl := visitTmpl

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			visit, err := client.GetAssuranceVisitById(ctx, app.DeputyId(), visitId)
			if err != nil {
				return err
			}
			vars.Visit = visit
			if visit.AssuranceType.Handle == "PDR" {
				tmpl = pdrTmpl
			}
			return nil
		})

		group.Go(func() error {
			visitors, err := client.GetVisitors(ctx)
			if err != nil {
				return err
			}
			vars.Visitors = visitors

			return nil
		})

		group.Go(func() error {
			visitRagRatingTypes, err := client.GetVisitRagRatingTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			vars.VisitRagRatingTypes = visitRagRatingTypes
			return nil
		})

		group.Go(func() error {
			visitOutcomeTypes, err := client.GetVisitOutcomeTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			vars.VisitOutcomeTypes = visitOutcomeTypes
			return nil
		})

		group.Go(func() error {
			pdrOutcomeTypes, err := client.GetPdrOutcomeTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			vars.PdrOutcomeTypes = pdrOutcomeTypes
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			reportReviewDate := r.PostFormValue("report-review-date")
			reviewedBy := 0
			if reportReviewDate != "" {
				reviewedBy = app.UserDetails.ID
			}

			pdrOutcome := ""
			if r.PostFormValue("pdr-outcome") == "Not received" {
				pdrOutcome = "NOT_RECEIVED"
			} else if r.PostFormValue("pdr-outcome") == "Received" {
				pdrOutcome = "RECEIVED"
			}

			manageAssuranceVisitForm := sirius.AssuranceVisitDetails{
				CommissionedDate:    r.PostFormValue("commissioned-date"),
				VisitorAllocated:    r.PostFormValue("visitor-allocated"),
				ReportDueDate:       r.PostFormValue("report-due-date"),
				ReportReceivedDate:  r.PostFormValue("report-received-date"),
				VisitOutcome:        r.PostFormValue("visit-outcome"),
				PdrOutcome:          pdrOutcome,
				ReportReviewDate:    reportReviewDate,
				VisitReportMarkedAs: r.PostFormValue("visit-report-marked-as"),
				ReviewedBy:          reviewedBy,
				Note:                strings.TrimSpace(r.PostFormValue("note")),
			}

			err := client.UpdateAssuranceVisit(ctx, manageAssuranceVisitForm, app.DeputyId(), visitId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				vars.ErrorNote = r.PostFormValue("note")

				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			success := "manageAssuranceVisit"
			if vars.Visit.AssuranceType.Handle == "PDR" {
				success = "managePDR"
			}

			return Redirect(fmt.Sprintf("/%d/assurance-visits?success=%s", app.DeputyId(), success))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
