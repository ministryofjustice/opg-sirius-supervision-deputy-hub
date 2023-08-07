package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"html/template"
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	DeputyHubClient
	DeputyHubInformation
	DeputyHubClientInformation
	DeputyHubContactInformation
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
	ContactInformation
	AddTasksClient
	TasksClient
	ManageTasks
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *logging.Logger, client Client, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {
	wrap := wrapHandler(logger, client, templates["error.gotmpl"], envVars)

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/health-check", healthCheck())

	pageRouter := router.PathPrefix("/{id}").Subrouter()
	pageRouter.Use(logging.Use(logger))

	pageRouter.Handle("",
		wrap(
			renderTemplateForDeputyHub(client, templates["deputy-details.gotmpl"])))

	pageRouter.Handle("/contacts",
		wrap(
			renderTemplateForContactTab(client, templates["contacts.gotmpl"])))

	pageRouter.Handle("/contacts/add-contact",
		wrap(
			renderTemplateForAddContact(client, templates["add-contact.gotmpl"])))

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

	pageRouter.Handle("/tasks",
		wrap(
			renderTemplateForTasks(client, templates["tasks.gotmpl"])))

	pageRouter.Handle("/tasks/add-task",
		wrap(
			renderTemplateForAddTask(client, templates["add-task.gotmpl"])))

	pageRouter.Handle("/tasks/{taskId}",
		wrap(
			renderTemplateForManageTasks(client, templates["manage-task.gotmpl"])))

	pageRouter.Handle("/manage-team-details",
		wrap(
			renderTemplateForEditDeputyHub(client, templates["manage-team-details.gotmpl"])))

	pageRouter.Handle("/change-ecm",
		wrap(
			renderTemplateForChangeECM(client, templates["change-ecm.gotmpl"])))

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
			renderTemplateForAssuranceVisits(client, templates["assurance-visits.gotmpl"])))

	pageRouter.Handle("/add-assurance-visit",
		wrap(
			renderTemplateForAddAssuranceVisit(client, templates["add-assurance-visit.gotmpl"])))

	pageRouter.Handle("/manage-assurance-visit/{visitId}",
		wrap(
			renderTemplateForManageAssuranceVisit(client, templates["manage-assurance-visit.gotmpl"], templates["manage-pdr.gotmpl"])))

	static := staticFileHandler(envVars.WebDir)
	router.PathPrefix("/assets/").Handler(static)
	router.PathPrefix("/javascript/").Handler(static)
	router.PathPrefix("/stylesheets/").Handler(static)

	router.NotFoundHandler = notFoundHandler(templates["error.gotmpl"], envVars.SiriusPublicURL)

	return http.StripPrefix(envVars.Prefix, securityheaders.Use(router))
}

func notFoundHandler(tmplError Template, siriusURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = tmplError.ExecuteTemplate(w, "page", ErrorVars{
			Code:  http.StatusNotFound,
			Error: "Not Found",
			AppVars: AppVars{
				EnvironmentVars: EnvironmentVars{
					SiriusURL: siriusURL,
				},
			},
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
