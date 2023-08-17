package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

type AppVars struct {
	Path          string
	XSRFToken     string
	UserDetails   sirius.UserDetails
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	EnvironmentVars
}

func (a AppVars) DeputyId() int {
	return a.DeputyDetails.ID
}

type AppVarsClient interface {
	GetUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetDeputyDetails(sirius.Context, int, int, int) (sirius.DeputyDetails, error)
}

func NewAppVars(client AppVarsClient, r *http.Request, envVars EnvironmentVars) (*AppVars, error) {
	ctx := getContext(r)
	group, groupCtx := errgroup.WithContext(ctx.Context)

	vars := AppVars{
		Path:            r.URL.Path,
		XSRFToken:       ctx.XSRFToken,
		EnvironmentVars: envVars,
	}

	group.Go(func() error {
		user, err := client.GetUserDetails(ctx.With(groupCtx))
		if err != nil {
			return err
		}
		vars.UserDetails = user
		return nil
	})
	group.Go(func() error {
		deputyId, _ := strconv.Atoi(mux.Vars(r)["id"])
		deputy, err := client.GetDeputyDetails(ctx.With(groupCtx), vars.DefaultPaTeam, vars.DefaultProTeam, deputyId)
		if err != nil {
			return err
		}
		vars.DeputyDetails = deputy
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return &vars, nil
}
