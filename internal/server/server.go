package server

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
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
	AddAssuranceClient
	GetAssurancesClient
	ManageAssuranceClient
	ManageContact
	DeleteContact
	DeleteDeputy
	AddTasksClient
	TasksClient
	DocumentsClient
	ReplaceDocumentClient
	AddDocumentClient
	ManageTasks
	CompleteTask
	AddGcmIssue
	GetGcmIssues
	CheckDocumentDownload
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *slog.Logger, client Client, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {
	mux := http.NewServeMux()

	wrap := wrapHandler(logger, client, templates["error.gotmpl"], envVars)

	// Static file routes
	static := staticFileHandler(envVars.WebDir)
	mux.Handle("/static/assets/", static)
	mux.Handle("/static/javascript/", static)
	mux.Handle("/static/stylesheets/", static)

	mux.Handle("/{id}", wrap(renderTemplateForDeputyHub(client, templates["deputy-details.gotmpl"])))
	mux.Handle("/{id}/clients", wrap(renderTemplateForClientTab(client, templates["clients.gotmpl"])))
	mux.Handle("/{id}/contacts", wrap(renderTemplateForContactTab(client, templates["contacts.gotmpl"])))
	mux.Handle("/{id}/contacts/add-contact", wrap(renderTemplateForManageContact(client, templates["manage-contact.gotmpl"])))
	mux.Handle("/{id}/contacts/{contactId}", wrap(renderTemplateForManageContact(client, templates["manage-contact.gotmpl"])))
	mux.Handle("/{id}/contacts/{contactId}/delete", wrap(renderTemplateForDeleteContact(client, templates["delete-contact.gotmpl"])))
	mux.Handle("/{id}/timeline", wrap(renderTemplateForDeputyHubEvents(client, templates["timeline.gotmpl"], envVars)))
	mux.Handle("/{id}/notes", wrap(renderTemplateForDeputyHubNotes(client, templates["notes.gotmpl"])))
	mux.Handle("/{id}/notes/add-note", wrap(renderTemplateForDeputyHubNotes(client, templates["add-notes.gotmpl"])))
	mux.Handle("/{id}/tasks", wrap(renderTemplateForTasks(client, templates["tasks.gotmpl"])))
	mux.Handle("/{id}/tasks/add-task", wrap(renderTemplateForAddTask(client, templates["add-task.gotmpl"])))
	mux.Handle("/{id}/tasks/{taskId}", wrap(renderTemplateForManageTasks(client, templates["manage-task.gotmpl"])))
	mux.Handle("/{id}/tasks/complete/{taskId}", wrap(renderTemplateForCompleteTask(client, templates["complete-task.gotmpl"])))
	mux.Handle("/{id}/documents", wrap(renderTemplateForDocuments(client, templates["documents.gotmpl"])))
	mux.Handle("/{id}/documents/add", wrap(renderTemplateForAddDocument(client, templates["add-document.gotmpl"])))
	mux.Handle("/{id}/documents/{documentId}/replace", wrap(renderTemplateForReplaceDocument(client, templates["replace-document.gotmpl"])))
	mux.Handle("/{id}/manage-team-details", wrap(renderTemplateForEditDeputyHub(client, templates["manage-team-details.gotmpl"])))
	mux.Handle("/{id}/change-ecm", wrap(renderTemplateForChangeECM(client, templates["change-ecm.gotmpl"])))
	mux.Handle("/{id}/change-firm", wrap(renderTemplateForChangeFirm(client, templates["change-firm.gotmpl"])))
	mux.Handle("/{id}/delete-deputy", wrap(renderTemplateForDeleteDeputy(client, templates["delete-deputy.gotmpl"])))
	mux.Handle("/{id}/add-firm", wrap(renderTemplateForAddFirm(client, templates["add-firm.gotmpl"])))
	mux.Handle("/{id}/manage-deputy-contact-details", wrap(renderTemplateForManageDeputyContactDetails(client, templates["manage-deputy-contact-details.gotmpl"])))
	mux.Handle("/{id}/manage-important-information", wrap(renderTemplateForImportantInformation(client, templates["manage-important-information.gotmpl"])))
	mux.Handle("/{id}/assurances", wrap(renderTemplateForAssurances(client, templates["assurances.gotmpl"])))
	mux.Handle("/{id}/add-assurance", wrap(renderTemplateForAddAssurance(client, templates["add-assurance.gotmpl"])))
	mux.Handle("/{id}/manage-assurance/{visitId}", wrap(renderTemplateForManageAssurance(client, templates["manage-visit.gotmpl"], templates["manage-pdr.gotmpl"])))
	mux.Handle("/{id}/gcm-issues/open-issues", wrap(renderTemplateForGcmIssues(client, templates["gcm-issues-list.gotmpl"])))
	mux.Handle("/{id}/gcm-issues/closed-issues", wrap(renderTemplateForGcmIssues(client, templates["gcm-issues-list.gotmpl"])))
	mux.Handle("/{id}/gcm-issues/add", wrap(renderTemplateForAddGcmIssue(client, templates["add-gcm-issue.gotmpl"])))
	mux.Handle("/{id}/documents/{documentId}/check", wrap(renderTemplateForCheckDocument(client)))

	// Health check
	mux.Handle("/health-check", healthCheck())

	// Fallback
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = templates["error.gotmpl"].ExecuteTemplate(w, "page", ErrorVars{
			Code:            http.StatusNotFound,
			Error:           "Page not found",
			EnvironmentVars: envVars,
		})
	})

	// Wrap all with security headers and telemetry
	return http.StripPrefix(envVars.Prefix, securityheaders.Use(telemetry.Middleware(logger)(mux)))
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
	h := http.FileServer(http.Dir(webDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "must-revalidate")
		h.ServeHTTP(w, r)
	})
}
