package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type ManageAssuranceVisits interface {
	GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error)
	CreateAssuranceVisit(ctx sirius.Context, deputyId int, assuranceVisitForm sirius.AssuranceVisit) error
}

type ManageAssuranceVisitsVars struct {
	Path                      string
	XSRFToken                 string
	DeputyDetails             sirius.DeputyDetails
	Error                     string
	Errors                    sirius.ValidationErrors
	IsFinanceManager          bool
}

func renderTemplateForAssuranceVisits(client ManageAssuranceVisits, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		vars := ManageAssuranceVisitsVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
		}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			userDetails, err := client.GetUserDetails(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			vars.IsFinanceManager = userDetails.IsFinanceManager()
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

			assuranceVisitForm := sirius.AssuranceVisit{
			}

			err = client.CreateAssuranceVisit(ctx, deputyId, assuranceVisitForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors

				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=assuranceVisit", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

