package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type Redirect string

func (e Redirect) Error() string {
	return "redirect to " + string(e)
}

func (e Redirect) To() string {
	return string(e)
}

type StatusError int

func (e StatusError) Error() string {
	code := e.Code()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func (e StatusError) Code() int {
	return int(e)
}

type Handler func(v AppVars, w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	Code  int
	Error string
	AppVars
}

type DeputyHubClient interface {
	GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error)
	GetDeputyDetails(sirius.Context, int, int, int) (sirius.DeputyDetails, error)
}

func wrapHandler(logger *logging.Logger, client DeputyHubClient, tmplError Template, envVars EnvironmentVars) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars, err := NewAppVars(client, r, envVars)

			if err == nil {
				err = next(*vars, w, r)
			}

			if err != nil {
				if err == sirius.ErrUnauthorized {
					http.Redirect(w, r, vars.SiriusURL+"/auth", http.StatusFound)
					return
				}

				if redirect, ok := err.(Redirect); ok {
					http.Redirect(w, r, vars.Prefix+redirect.To(), http.StatusFound)
					return
				}

				logger.Request(r, err)

				code := http.StatusInternalServerError
				if status, ok := err.(StatusError); ok {
					if status.Code() == http.StatusForbidden || status.Code() == http.StatusNotFound {
						code = status.Code()
					}
				}

				w.WriteHeader(code)
				err = tmplError.ExecuteTemplate(w, "page", errorVars{
					Code:    code,
					Error:   err.Error(),
					AppVars: *vars,
				})

				if err != nil {
					logger.Request(r, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}
