package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type AddAssuranceVisit interface {
	GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error)
	AddAssuranceVisit(ctx sirius.Context, requestedDate string, userId, deputyId int) error
}

type AddAssuranceVisitVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
}

func renderTemplateForAddAssuranceVisit(client AddAssuranceVisit, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		vars := AddAssuranceVisitVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var requestedDate = r.PostFormValue("requested-date")

			if requestedDate == "" {
				vars.Errors = sirius.ValidationErrors{
					"requested-date": {"": "Enter a real date"},
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			user, err := client.GetUserDetails(ctx)
			if err != nil {
				return err
			}

			err = client.AddAssuranceVisit(ctx, requestedDate, user.ID, deputyId)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := AddAssuranceVisitVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/%d/assurance-visits?success=addAssuranceVisit", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
