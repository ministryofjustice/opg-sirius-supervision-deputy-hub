package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
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

type ErrorVars struct {
	Code  int
	Error string
	EnvironmentVars
}

type DeputyHubClient interface {
	GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error)
	GetDeputyDetails(sirius.Context, int, int, int) (sirius.DeputyDetails, error)
}

type ExpandedError interface {
	Title() string
	Data() interface{}
}

func LoggerRequest(l *slog.Logger, r *http.Request, err error) {
	if ee, ok := err.(ExpandedError); ok {
		l.Info(ee.Title(),
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()),
			slog.Any("data", ee.Data()))
	} else if err != nil {
		l.Info(err.Error(),
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()))
	} else {
		l.Info("",
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()))
	}
}

func wrapHandler(logger *slog.Logger, client DeputyHubClient, tmplError Template, envVars EnvironmentVars) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars, err := NewAppVars(client, r, envVars)

			if err == nil {
				err = next(*vars, w, r)
			}

			if err != nil {
				if err == sirius.ErrUnauthorized {
					redirect := ""

					if r.RequestURI != "" {
						redirect = "?redirect=" + r.RequestURI
					}

					http.Redirect(w, r, envVars.SiriusPublicURL+"/auth"+redirect, http.StatusFound)
					return
				}

				if redirect, ok := err.(Redirect); ok {
					http.Redirect(w, r, envVars.Prefix+redirect.To(), http.StatusFound)
					return
				}

				LoggerRequest(logger, r, err)

				code := http.StatusInternalServerError
				if serverStatusError, ok := err.(StatusError); ok {
					code = serverStatusError.Code()
				}
				if siriusStatusError, ok := err.(sirius.StatusError); ok {
					code = siriusStatusError.Code
				}

				w.WriteHeader(code)
				errVars := ErrorVars{
					Code:            code,
					Error:           err.Error(),
					EnvironmentVars: envVars,
				}
				err = tmplError.ExecuteTemplate(w, "page", errVars)

				if err != nil {
					LoggerRequest(logger, r, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}
