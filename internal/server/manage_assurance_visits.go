package server

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ManageAssuranceVisit interface {
	GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error)
	UpdateAssuranceVisit(ctx sirius.Context, manageAssuranceVisitForm sirius.AssuranceVisitDetails, deputyId, visitId int) error
	GetVisitors(ctx sirius.Context) (sirius.Visitors, error)
	GetVisitRagRatingTypes(ctx sirius.Context) ([]sirius.VisitRagRatingTypes, error)
	GetVisitOutcomeTypes(ctx sirius.Context) ([]sirius.VisitOutcomeTypes, error)
	GetAssuranceVisitById(ctx sirius.Context, deputyId int, visitId int) (sirius.AssuranceVisit, error)
}

type ManageAssuranceVisitVars struct {
	Path                string
	XSRFToken           string
	DeputyDetails       sirius.DeputyDetails
	Error               string
	Errors              sirius.ValidationErrors
	Visitors            sirius.Visitors
	VisitRagRatingTypes []sirius.VisitRagRatingTypes
	VisitOutcomeTypes   []sirius.VisitOutcomeTypes
	Visit               sirius.AssuranceVisit
}

func renderTemplateForManageAssuranceVisit(client ManageAssuranceVisit, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		visitId, _ := strconv.Atoi(routeVars["visitId"])

		vars := ManageAssuranceVisitVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
		}

		visit, err := client.GetAssuranceVisitById(ctx, deputyId, visitId)
		vars.Visit = visit
		if err != nil {
			return err
		}

		visitors, err := client.GetVisitors(ctx)
		vars.Visitors = visitors
		if err != nil {
			return err
		}

		group, groupCtx := errgroup.WithContext(ctx.Context)

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

		if err := group.Wait(); err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var err error
			user, err := client.GetUserDetails(ctx)
			if err != nil {
				return err
			}

			reportReviewDate := r.PostFormValue("report-review-date")
			reviewedBy := 0
			if reportReviewDate != "" {
				reviewedBy = user.ID
			}

			manageAssuranceVisitForm := sirius.AssuranceVisitDetails{
				CommissionedDate:    r.PostFormValue("commissioned-date"),
				VisitorAllocated:    r.PostFormValue("visitor-allocated"),
				ReportDueDate:       r.PostFormValue("report-due-date"),
				ReportReceivedDate:  r.PostFormValue("report-received-date"),
				VisitOutcome:        r.PostFormValue("visit-outcome"),
				ReportReviewDate:    reportReviewDate,
				VisitReportMarkedAs: r.PostFormValue("visit-report-marked-as"),
				ReviewedBy:          reviewedBy,
			}

			err = client.UpdateAssuranceVisit(ctx, manageAssuranceVisitForm, deputyId, visitId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := ManageAssuranceVisitVars{
					Path:                r.URL.Path,
					XSRFToken:           ctx.XSRFToken,
					Errors:              verr.Errors,
					VisitRagRatingTypes: vars.VisitRagRatingTypes,
					VisitOutcomeTypes:   vars.VisitOutcomeTypes,
					Visitors:            visitors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/%d/assurance-visits?success=manageAssuranceVisit", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
