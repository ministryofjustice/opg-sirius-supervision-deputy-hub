package server

import (
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type ApiClient interface {
	DeputyHubClient
	//DeputyHubInformation
	//DeputyHubClientInformation
	Contacts
	DeleteContact
	ManageContact
	AddDocument
	Documents
	Notes
	Tasks
	AddTask
	ManageTasks
	CompleteTask
	Timeline
	//EditDeputyHubInformation
	//ChangeECMInformation
	//FirmInformation
	//DeputyContactDetailsInformation
	//ManageProDeputyImportantInformation
	//DeputyChangeFirmInformation
	//AddAssuranceClient
	//GetAssurancesClient
	//ManageAssuranceClient
	//DeleteDeputy

}

type Template interface {
	Execute(wr io.Writer, data any) error
	ExecuteTemplate(io.Writer, string, interface{}) error
}

type router interface {
	Client() ApiClient
	execute(http.ResponseWriter, *http.Request, any, AppVars) error
}

func New(logger *slog.Logger, client ApiClient, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {
	wrap := wrapHandler(client, logger, templates["error.gotmpl"], envVars)
	mux := http.NewServeMux()

	mux.Handle("GET /{deputyId}/contacts", wrap(&ListContactsHandler{&route{client: client, tmpl: templates["contacts.gotmpl"], partial: "contacts"}}))
	mux.Handle("GET /{deputyId}/contacts/add-contact", wrap(&ManageContactsHandler{&route{client: client, tmpl: templates["manage-contact.gotmpl"], partial: "manage-contact"}}))
	mux.Handle("POST /{deputyId}/contacts/add-contact", wrap(&ManageContactsHandler{&route{client: client, tmpl: templates["manage-contact.gotmpl"], partial: "manage-contact"}}))
	mux.Handle("GET /{deputyId}/contacts/{contactId}", wrap(&ManageContactsHandler{&route{client: client, tmpl: templates["manage-contact.gotmpl"], partial: "manage-contact"}}))
	mux.Handle("POST /{deputyId}/contacts/{contactId}", wrap(&ManageContactsHandler{&route{client: client, tmpl: templates["manage-contact.gotmpl"], partial: "manage-contact"}}))
	mux.Handle("GET /{deputyId}/contacts/{contactId}/delete", wrap(&DeleteContactHandler{&route{client: client, tmpl: templates["delete-contact.gotmpl"], partial: "delete-contact"}}))
	mux.Handle("POST /{deputyId}/contacts/{contactId}/delete", wrap(&DeleteContactHandler{&route{client: client, tmpl: templates["delete-contact.gotmpl"], partial: "delete-contact"}}))

	mux.Handle("GET /{deputyId}/documents", wrap(&ListDocumentsHandler{&route{client: client, tmpl: templates["documents.gotmpl"], partial: "documents"}}))
	mux.Handle("GET /{deputyId}/documents/add", wrap(&AddDocumentHandler{&route{client: client, tmpl: templates["add-document.gotmpl"], partial: "add-document"}}))
	mux.Handle("POST /{deputyId}/documents/add", wrap(&AddDocumentHandler{&route{client: client, tmpl: templates["add-document.gotmpl"], partial: "add-document"}}))

	mux.Handle("GET /{deputyId}/timeline", wrap(&TimelineHandler{&route{client: client, tmpl: templates["timeline.gotmpl"], partial: "timeline"}}))

	mux.Handle("GET /{deputyId}/notes", wrap(&NotesHandler{&route{client: client, tmpl: templates["notes.gotmpl"], partial: "notes"}}))
	mux.Handle("GET /{deputyId}/notes/add-note", wrap(&NotesHandler{&route{client: client, tmpl: templates["add-notes.gotmpl"], partial: "add-notes"}}))
	mux.Handle("POST /{deputyId}/notes/add-note", wrap(&NotesHandler{&route{client: client, tmpl: templates["add-notes.gotmpl"], partial: "add-notes"}}))

	mux.Handle("GET /{deputyId}/tasks", wrap(&TasksHandler{&route{client: client, tmpl: templates["tasks.gotmpl"], partial: "tasks"}}))
	mux.Handle("GET /{deputyId}/tasks/add-task", wrap(&AddTaskHandler{&route{client: client, tmpl: templates["add-task.gotmpl"], partial: "add-task"}}))
	mux.Handle("POST /{deputyId}/tasks/add-task", wrap(&AddTaskHandler{&route{client: client, tmpl: templates["add-task.gotmpl"], partial: "add-task"}}))
	mux.Handle("GET /{deputyId}/tasks/{taskId}", wrap(&ManageTaskHandler{&route{client: client, tmpl: templates["manage-task.gotmpl"], partial: "manage-task"}}))
	mux.Handle("POST /{deputyId}/tasks/{taskId}", wrap(&ManageTaskHandler{&route{client: client, tmpl: templates["manage-task.gotmpl"], partial: "manage-task"}}))
	mux.Handle("GET /{deputyId}/tasks/complete/{taskId}", wrap(&CompleteTaskHandler{&route{client: client, tmpl: templates["complete-task.gotmpl"], partial: "complete-task"}}))
	mux.Handle("POST /{deputyId}/tasks/complete/{taskId}", wrap(&CompleteTaskHandler{&route{client: client, tmpl: templates["complete-task.gotmpl"], partial: "complete-task"}}))

	//pageRouter.Handle("",
	//	wrap(
	//		renderTemplateForDeputyHub(client, templates["deputy-details.gotmpl"])))
	//

	//pageRouter.Handle("/clients",
	//	wrap(
	//		renderTemplateForClientTab(client, templates["clients.gotmpl"])))
	//

	//pageRouter.Handle("/manage-team-details",
	//	wrap(
	//		renderTemplateForEditDeputyHub(client, templates["manage-team-details.gotmpl"])))
	//
	//pageRouter.Handle("/change-ecm",
	//	wrap(
	//		renderTemplateForChangeECM(client, templates["change-ecm.gotmpl"])))
	//
	//pageRouter.Handle("/change-firm",
	//	wrap(
	//		renderTemplateForChangeFirm(client, templates["change-firm.gotmpl"])))
	//
	//pageRouter.Handle("/delete-deputy",
	//	wrap(
	//		renderTemplateForDeleteDeputy(client, templates["delete-deputy.gotmpl"])))
	//
	//pageRouter.Handle("/add-firm",
	//	wrap(
	//		renderTemplateForAddFirm(client, templates["add-firm.gotmpl"])))
	//
	//pageRouter.Handle("/manage-deputy-contact-details",
	//	wrap(
	//		renderTemplateForManageDeputyContactDetails(client, templates["manage-deputy-contact-details.gotmpl"])))
	//
	//pageRouter.Handle("/manage-important-information",
	//	wrap(
	//		renderTemplateForImportantInformation(client, templates["manage-important-information.gotmpl"])))
	//
	//pageRouter.Handle("/assurances",
	//	wrap(
	//		renderTemplateForAssurances(client, templates["assurances.gotmpl"])))
	//
	//pageRouter.Handle("/add-assurance",
	//	wrap(
	//		renderTemplateForAddAssurance(client, templates["add-assurance.gotmpl"])))
	//
	//pageRouter.Handle("/manage-assurance/{visitId}",
	//	wrap(
	//		renderTemplateForManageAssurance(client, templates["manage-visit.gotmpl"], templates["manage-pdr.gotmpl"])))

	//static := staticFileHandler(envVars.WebDir)
	//router.PathPrefix("/assets/").Handler(static)
	//router.PathPrefix("/javascript/").Handler(static)
	//router.PathPrefix("/stylesheets/").Handler(static)
	//
	mux.Handle("/health-check", healthCheck())

	//static := http.FileServer(http.Dir(envVars.WebDir + "/static"))
	//mux.Handle("/assets/", static)
	//mux.Handle("/javascript/", static)
	//mux.Handle("/stylesheets/", static)

	return otelhttp.NewHandler(http.StripPrefix(envVars.Prefix, securityheaders.Use(mux)), "supervision-deputy-hub")
}

//func notFoundHandler(tmplError Template, envVars EnvironmentVars) Handler {
//	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
//		_ = tmplError.ExecuteTemplate(w, "page", ErrorVars{
//			Code:            http.StatusNotFound,
//			Error:           "Page not found",
//			EnvironmentVars: envVars,
//		})
//		return nil
//	}
//}

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
