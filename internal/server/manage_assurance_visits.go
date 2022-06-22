package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ManageAssuranceVisit interface {
	GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error)
	UpdateAssuranceVisit(ctx sirius.Context, requestedDate string, userId, deputyId int) error
	GetVisitors(ctx sirius.Context) (sirius.Visitors, error)
}

type ManageAssuranceVisitVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	ErrorMessage  string
	Visitors      sirius.Visitors
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

		switch r.Method {
		case http.MethodGet:
			visitors, err := client.GetVisitors(ctx)
			vars.Visitors = visitors
			if err != nil {
				return err
			}

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
