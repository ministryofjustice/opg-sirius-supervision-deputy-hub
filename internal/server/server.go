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
	Assurances
	AddAssurance
	ManageAssurance
	Clients
	Contacts
	EditContact
	DeleteContact
	ManageContact
	Deputy
	DeleteDeputy
	DeputyHubClient
	EditDeputyTeam
	EditDeputyEcm
	Documents
	AddDocument
	ReplaceDocument
	AddFirm
	EditFirm
	Notes
	EditProImportantInformation
	Tasks
	AddTask
	CompleteTask
	ManageTasks
	Timeline
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

	http.Handle("/assets/", http.StripPrefix("/static/", http.FileServer(http.Dir(envVars.WebDir+"/assets"))))
	http.Handle("/javascript/", http.StripPrefix("/static/", http.FileServer(http.Dir(envVars.WebDir+"/javascript"))))
	http.Handle("/stylesheets/", http.StripPrefix("/static/", http.FileServer(http.Dir(envVars.WebDir+"/stylesheets"))))

	mux.Handle("GET /{deputyId}/clients", wrap(&ClientHandler{&route{client: client, tmpl: templates["clients.gotmpl"], partial: "clients"}}))

	mux.Handle("GET /{deputyId}/assurances", wrap(&AssurancesHandler{&route{client: client, tmpl: templates["assurances.gotmpl"], partial: "assurances"}}))
	mux.Handle("GET /{deputyId}/assurance/add", wrap(&AddAssuranceHandler{&route{client: client, tmpl: templates["add-assurance.gotmpl"], partial: "add-assurance"}}))
	mux.Handle("POST /{deputyId}/assurance/add", wrap(&AddAssuranceHandler{&route{client: client, tmpl: templates["add-assurance.gotmpl"], partial: "add-assurance"}}))
	mux.Handle("GET /{deputyId}/manage-visit/{visitId}", wrap(&ManageAssuranceHandler{&route{client: client, tmpl: templates["manage-visit.gotmpl"], partial: "manage-visit"}}))
	mux.Handle("GET /{deputyId}/manage-assurance/{visitId}", wrap(&ManageAssuranceHandler{&route{client: client, tmpl: templates["manage-pdr.gotmpl"], partial: "manage-pdr"}}))
	mux.Handle("POST /{deputyId}/manage-visit/{visitId}", wrap(&ManageAssuranceHandler{&route{client: client, tmpl: templates["manage-visit.gotmpl"], partial: "manage-visit"}}))
	mux.Handle("POST /{deputyId}/manage-assurance/{visitId}", wrap(&ManageAssuranceHandler{&route{client: client, tmpl: templates["manage-pdr.gotmpl"], partial: "manage-pdr"}}))

	mux.Handle("GET /{deputyId}/contacts", wrap(&ListContactsHandler{&route{client: client, tmpl: templates["contacts.gotmpl"], partial: "contacts"}}))
	mux.Handle("GET /{deputyId}/contacts/add", wrap(&ManageContactsHandler{&route{client: client, tmpl: templates["manage-contact.gotmpl"], partial: "manage-contact"}}))
	mux.Handle("POST /{deputyId}/contacts/add", wrap(&ManageContactsHandler{&route{client: client, tmpl: templates["manage-contact.gotmpl"], partial: "manage-contact"}}))
	mux.Handle("GET /{deputyId}/contacts/{contactId}", wrap(&ManageContactsHandler{&route{client: client, tmpl: templates["manage-contact.gotmpl"], partial: "manage-contact"}}))
	mux.Handle("POST /{deputyId}/contacts/{contactId}", wrap(&ManageContactsHandler{&route{client: client, tmpl: templates["manage-contact.gotmpl"], partial: "manage-contact"}}))
	mux.Handle("GET /{deputyId}/contacts/{contactId}/delete", wrap(&DeleteContactHandler{&route{client: client, tmpl: templates["delete-contact.gotmpl"], partial: "delete-contact"}}))
	mux.Handle("POST /{deputyId}/contacts/{contactId}/delete", wrap(&DeleteContactHandler{&route{client: client, tmpl: templates["delete-contact.gotmpl"], partial: "delete-contact"}}))

	mux.Handle("GET /{deputyId}", wrap(&DeputyHandler{&route{client: client, tmpl: templates["deputy-details.gotmpl"], partial: "deputy-details"}}))
	mux.Handle("GET /{deputyId}/delete", wrap(&DeleteDeputyHandler{&route{client: client, tmpl: templates["delete-deputy.gotmpl"], partial: "delete-deputy"}}))
	mux.Handle("POST /{deputyId}/delete", wrap(&DeleteDeputyHandler{&route{client: client, tmpl: templates["delete-deputy.gotmpl"], partial: "delete-deputy"}}))

	mux.Handle("GET /{deputyId}/edit-deputy-team", wrap(&EditDeputyTeamHandler{&route{client: client, tmpl: templates["manage-deputy-team.gotmpl"], partial: "manage-team-details"}}))
	mux.Handle("POST /{deputyId}/edit-deputy-team", wrap(&EditDeputyTeamHandler{&route{client: client, tmpl: templates["manage-deputy-team.gotmpl"], partial: "manage-team-details"}}))
	mux.Handle("GET /{deputyId}/edit-important-information", wrap(&EditDeputyEcmHandler{&route{client: client, tmpl: templates["manage-important-information.gotmpl"], partial: "manage-important-information"}}))
	mux.Handle("POST /{deputyId}/edit-important-information", wrap(&EditDeputyEcmHandler{&route{client: client, tmpl: templates["manage-important-information.gotmpl"], partial: "manage-important-information"}}))
	mux.Handle("GET /{deputyId}/change-ecm", wrap(&EditDeputyEcmHandler{&route{client: client, tmpl: templates["change-ecm.gotmpl"], partial: "change-ecm"}}))
	mux.Handle("POST /{deputyId}/change-ecm", wrap(&EditDeputyEcmHandler{&route{client: client, tmpl: templates["change-ecm.gotmpl"], partial: "change-ecm"}}))

	mux.Handle("GET /{deputyId}/firm/add", wrap(&AddFirmHandler{&route{client: client, tmpl: templates["add-firm.gotmpl"], partial: "add-firm"}}))
	mux.Handle("POST /{deputyId}/firm/add", wrap(&AddFirmHandler{&route{client: client, tmpl: templates["add-firm.gotmpl"], partial: "add-firm"}}))
	mux.Handle("GET /{deputyId}/firm/change", wrap(&EditFirmHandler{&route{client: client, tmpl: templates["change-firm.gotmpl"], partial: "change-firm"}}))
	mux.Handle("POST /{deputyId}/firm/change", wrap(&EditFirmHandler{&route{client: client, tmpl: templates["change-firm.gotmpl"], partial: "change-firm"}}))

	mux.Handle("GET /{deputyId}/documents", wrap(&ListDocumentsHandler{&route{client: client, tmpl: templates["documents.gotmpl"], partial: "documents"}}))
	mux.Handle("GET /{deputyId}/documents/add", wrap(&AddDocumentHandler{&route{client: client, tmpl: templates["add-document.gotmpl"], partial: "add-document"}}))
	mux.Handle("POST /{deputyId}/documents/add", wrap(&AddDocumentHandler{&route{client: client, tmpl: templates["add-document.gotmpl"], partial: "add-document"}}))
	mux.Handle("GET /{deputyId}/documents/{documentId}/replace", wrap(&ReplaceDocumentHandler{&route{client: client, tmpl: templates["replace-document.gotmpl"], partial: "replace-document"}}))
	mux.Handle("POST /{deputyId}/documents/{documentId}/replace", wrap(&ReplaceDocumentHandler{&route{client: client, tmpl: templates["replace-document.gotmpl"], partial: "replace-document"}}))

	mux.Handle("GET /{deputyId}/notes", wrap(&NotesHandler{&route{client: client, tmpl: templates["notes.gotmpl"], partial: "notes"}}))
	mux.Handle("GET /{deputyId}/notes/add", wrap(&NotesHandler{&route{client: client, tmpl: templates["add-notes.gotmpl"], partial: "add-notes"}}))
	mux.Handle("POST /{deputyId}/notes/add", wrap(&NotesHandler{&route{client: client, tmpl: templates["add-notes.gotmpl"], partial: "add-notes"}}))

	mux.Handle("GET /{deputyId}/timeline", wrap(&TimelineHandler{&route{client: client, tmpl: templates["timeline.gotmpl"], partial: "timeline"}}))

	mux.Handle("GET /{deputyId}/tasks", wrap(&TasksHandler{&route{client: client, tmpl: templates["tasks.gotmpl"], partial: "tasks"}}))
	mux.Handle("GET /{deputyId}/tasks/add", wrap(&AddTaskHandler{&route{client: client, tmpl: templates["add-task.gotmpl"], partial: "add-task"}}))
	mux.Handle("POST /{deputyId}/tasks/add", wrap(&AddTaskHandler{&route{client: client, tmpl: templates["add-task.gotmpl"], partial: "add-task"}}))
	mux.Handle("GET /{deputyId}/tasks/{taskId}", wrap(&ManageTaskHandler{&route{client: client, tmpl: templates["manage-task.gotmpl"], partial: "manage-task"}}))
	mux.Handle("POST /{deputyId}/tasks/{taskId}", wrap(&ManageTaskHandler{&route{client: client, tmpl: templates["manage-task.gotmpl"], partial: "manage-task"}}))
	mux.Handle("GET /{deputyId}/tasks/{taskId}/complete", wrap(&CompleteTaskHandler{&route{client: client, tmpl: templates["complete-task.gotmpl"], partial: "complete-task"}}))
	mux.Handle("POST /{deputyId}/tasks/{taskId}/complete", wrap(&CompleteTaskHandler{&route{client: client, tmpl: templates["complete-task.gotmpl"], partial: "complete-task"}}))

	//static := staticFileHandler(envVars.WebDir)
	//router.PathPrefix("/assets/").Handler(static)
	//router.PathPrefix("/javascript/").Handler(static)
	//router.PathPrefix("/stylesheets/").Handler(static)
	////
	mux.Handle("GET /health-check", healthCheck())

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
