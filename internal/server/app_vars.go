package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type AppVars struct {
	Path           string
	XSRFToken      string
	UserDetails    sirius.UserDetails
	DeputyDetails  sirius.DeputyDetails
	SuccessMessage string
	PageName       string
	Error          string
	Errors         sirius.ValidationErrors
	EnvironmentVars
}

func (a AppVars) DeputyId() int {
	return a.DeputyDetails.ID
}

func (a AppVars) DeputyType() string {
	return a.DeputyDetails.DeputyType.Handle
}

type AppVarsClient interface {
	GetUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetDeputyDetails(sirius.Context, int, int, int) (sirius.DeputyDetails, error)
}

func NewAppVars(client AppVarsClient, r *http.Request, envVars EnvironmentVars) (*AppVars, error) {
	ctx := getContext(r)
	group, groupCtx := errgroup.WithContext(ctx.Context)
	deputyId, _ := strconv.Atoi(r.PathValue("id"))

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
