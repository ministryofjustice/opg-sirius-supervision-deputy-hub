package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Client interface {
	ErrorHandlerClient
	DeputyHubInformation
	DeputyHubClientInformation
	DeputyHubEventInformation
	DeputyHubNotesInformation
	EditDeputyHubInformation
	ChangeECMInformation
	FirmInformation
	DeputyContactDetailsInformation
	ManageProDeputyImportantInformation
	DeputyChangeFirmInformation
	AddAssuranceVisit
	AssuranceVisit
	ManageAssuranceVisit
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *logging.Logger, client Client, templates map[string]*template.Template, prefix, siriusPublicURL, webDir string, defaultPATeam int) http.Handler {
	wrap := errorHandler(logger, client, templates["error.gotmpl"], prefix, siriusPublicURL, defaultPATeam)

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/health-check", healthCheck())

	pageRouter := router.PathPrefix("/{id}").Subrouter()
	pageRouter.Use(logging.Use(logger))

	pageRouter.Handle("",
		wrap(
			renderTemplateForDeputyHub(client, templates["deputy-details.gotmpl"])))

	pageRouter.Handle("/clients",
		wrap(
			renderTemplateForClientTab(client, templates["clients.gotmpl"])))

	pageRouter.Handle("/timeline",
		wrap(
			renderTemplateForDeputyHubEvents(client, templates["timeline.gotmpl"])))

	pageRouter.Handle("/notes",
		wrap(
			renderTemplateForDeputyHubNotes(client, templates["notes.gotmpl"])))

	pageRouter.Handle("/notes/add-note",
		wrap(
			renderTemplateForDeputyHubNotes(client, templates["add-notes.gotmpl"])))

	pageRouter.Handle("/manage-team-details",
		wrap(
			renderTemplateForEditDeputyHub(client, templates["manage-team-details.gotmpl"])))

	pageRouter.Handle("/change-ecm",
		wrap(
			renderTemplateForChangeECM(client, defaultPATeam, templates["change-ecm.gotmpl"])))

	pageRouter.Handle("/change-firm",
		wrap(
			renderTemplateForChangeFirm(client, templates["change-firm.gotmpl"])))

	pageRouter.Handle("/add-firm",
		wrap(
			renderTemplateForAddFirm(client, templates["add-firm.gotmpl"])))

	pageRouter.Handle("/manage-deputy-contact-details",
		wrap(
			renderTemplateForManageDeputyContactDetails(client, templates["manage-deputy-contact-details.gotmpl"])))

	pageRouter.Handle("/manage-important-information",
		wrap(
			renderTemplateForImportantInformation(client, templates["manage-important-information.gotmpl"])))

	pageRouter.Handle("/assurance-visits",
		wrap(
			renderTemplateForAssuranceVisits(client, templates["assurance-visit.gotmpl"])))

	pageRouter.Handle("/add-assurance-visit",
		wrap(
			renderTemplateForAddAssuranceVisit(client, templates["add-assurance-visit.gotmpl"])))

	pageRouter.Handle("/manage-assurance-visit/{visitId}",
		wrap(
			renderTemplateForManageAssuranceVisit(client, templates["manage-assurance-visit.gotmpl"])))

	static := staticFileHandler(webDir)
	router.PathPrefix("/assets/").Handler(static)
	router.PathPrefix("/javascript/").Handler(static)
	router.PathPrefix("/stylesheets/").Handler(static)

	router.NotFoundHandler = notFoundHandler(templates["error.gotmpl"], siriusPublicURL)

	return http.StripPrefix(prefix, securityheaders.Use(router))
}

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

type Handler func(d sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	Firstname string
	Surname   string
	SiriusURL string
	Path      string
	Code      int
	Error     string
	Errors    string
}

type ErrorHandlerClient interface {
	GetDeputyDetails(sirius.Context, int, int) (sirius.DeputyDetails, error)
}

func errorHandler(logger *logging.Logger, client ErrorHandlerClient, tmplError Template, prefix, siriusURL string, defaultPATeam int) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			deputyId, _ := strconv.Atoi(mux.Vars(r)["id"])
			deputyDetails, err := client.GetDeputyDetails(getContext(r), defaultPATeam, deputyId)

			if err == nil {
				err = next(deputyDetails, w, r)
			}

			if err != nil {
				if err == sirius.ErrUnauthorized {
					http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
					return
				}

				if redirect, ok := err.(Redirect); ok {
					http.Redirect(w, r, prefix+redirect.To(), http.StatusFound)
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
					Firstname: "",
					Surname:   "",
					SiriusURL: siriusURL,
					Path:      "",
					Code:      code,
					Error:     err.Error(),
				})

				if err != nil {
					logger.Request(r, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}

func notFoundHandler(tmplError Template, siriusURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = tmplError.ExecuteTemplate(w, "page", errorVars{
			SiriusURL: siriusURL,
			Code:      http.StatusNotFound,
			Error:     "Not Found",
		})
	}
}

func getContext(r *http.Request) sirius.Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	return sirius.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}

func staticFileHandler(webDir string) http.Handler {
	h := http.FileServer(http.Dir(webDir + "/static"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "must-revalidate")
		h.ServeHTTP(w, r)
	})
}
