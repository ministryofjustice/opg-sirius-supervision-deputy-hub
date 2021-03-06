package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type Logger interface {
	Request(*http.Request, error)
}

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

func New(logger Logger, client Client, templates map[string]*template.Template, prefix, siriusPublicURL, webDir string, defaultPATeam int) http.Handler {
	wrap := errorHandler(logger, client, templates["error.gotmpl"], prefix, siriusPublicURL, defaultPATeam)

	router := mux.NewRouter()
	router.Handle("/health-check", healthCheck())

	router.Handle("/{id}",
		wrap(
			renderTemplateForDeputyHub(client, templates["deputy-details.gotmpl"])))

	router.Handle("/{id}/clients",
		wrap(
			renderTemplateForClientTab(client, templates["clients.gotmpl"])))

	router.Handle("/{id}/timeline",
		wrap(
			renderTemplateForDeputyHubEvents(client, templates["timeline.gotmpl"])))

	router.Handle("/{id}/notes",
		wrap(
			renderTemplateForDeputyHubNotes(client, templates["notes.gotmpl"])))

	router.Handle("/{id}/notes/add-note",
		wrap(
			renderTemplateForDeputyHubNotes(client, templates["add-notes.gotmpl"])))

	router.Handle("/{id}/manage-team-details",
		wrap(
			renderTemplateForEditDeputyHub(client, templates["manage-team-details.gotmpl"])))

	router.Handle("/{id}/change-ecm",
		wrap(
			renderTemplateForChangeECM(client, defaultPATeam, templates["change-ecm.gotmpl"])))

	router.Handle("/{id}/change-firm",
		wrap(
			renderTemplateForChangeFirm(client, templates["change-firm.gotmpl"])))

	router.Handle("/{id}/add-firm",
		wrap(
			renderTemplateForAddFirm(client, templates["add-firm.gotmpl"])))

	router.Handle("/{id}/manage-deputy-contact-details",
		wrap(
			renderTemplateForManageDeputyContactDetails(client, templates["manage-deputy-contact-details.gotmpl"])))

	router.Handle("/{id}/manage-important-information",
		wrap(
			renderTemplateForImportantInformation(client, templates["manage-important-information.gotmpl"])))

	router.Handle("/{id}/assurance-visits",
		wrap(
			renderTemplateForAssuranceVisits(client, templates["assurance-visit.gotmpl"])))

	router.Handle("/{id}/add-assurance-visit",
		wrap(
			renderTemplateForAddAssuranceVisit(client, templates["add-assurance-visit.gotmpl"])))

	router.Handle("/{id}/manage-assurance-visit/{visitId}",
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

type Handler func(perm sirius.PermissionSet, d sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error

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
	MyPermissions(sirius.Context) (sirius.PermissionSet, error)
	GetDeputyDetails(sirius.Context, int, int) (sirius.DeputyDetails, error)
}

func errorHandler(logger Logger, client ErrorHandlerClient, tmplError Template, prefix, siriusURL string, defaultPATeam int) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			myPermissions, err := client.MyPermissions(getContext(r))

			if err == nil {
				routeVars := mux.Vars(r)
				deputyId, _ := strconv.Atoi(routeVars["id"])
				var deputyDetails sirius.DeputyDetails
				deputyDetails, err = client.GetDeputyDetails(getContext(r), defaultPATeam, deputyId)

				if err == nil {
					err = next(myPermissions, deputyDetails, w, r)
				}
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
