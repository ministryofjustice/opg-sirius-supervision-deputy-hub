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
	UpdateAssuranceVisit(ctx sirius.Context, requestedDate string, userId, deputyId int) error
	GetVisitors(ctx sirius.Context) (sirius.Visitors, error)
	GetVisitRagRatingTypes(ctx sirius.Context) ([]sirius.VisitRagRatingTypes, error)
	GetVisitOutcomeTypes(ctx sirius.Context) ([]sirius.VisitOutcomeTypes, error)
}

type ManageAssuranceVisitVars struct {
	Path                string
	XSRFToken           string
	DeputyDetails       sirius.DeputyDetails
	Error               string
	Errors              sirius.ValidationErrors
	ErrorMessage        string
	Visitors            sirius.Visitors
	VisitRagRatingTypes []sirius.VisitRagRatingTypes
	VisitOutcomeTypes   []sirius.VisitOutcomeTypes
}

func renderTemplateForManageAssuranceVisit(client ManageAssuranceVisit, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		vars := ManageAssuranceVisitVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
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
			var requestedDate = r.PostFormValue("requested-date")

			if requestedDate == "" {
				vars.Errors = sirius.ValidationErrors{
					"commissioned-date":    {"": "Enter a real date"},
					"report-due-date":      {"": "Enter a real date"},
					"report-received-date": {"": "Enter a real date"},
					"report-review-date":   {"": "Enter a real date"},
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			user, err := client.GetUserDetails(ctx)
			if err != nil {
				return err
			}

			err = client.UpdateAssuranceVisit(ctx, requestedDate, user.ID, deputyId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := ManageAssuranceVisitVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/%d/assurance-visits?success=manageAssuranceVisit", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
