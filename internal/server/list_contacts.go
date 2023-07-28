package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubContactInformation interface {
	GetDeputyContacts(sirius.Context, int) (sirius.ContactList, error)
}

type ListContactsVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	ContactList    sirius.ContactList
	SuccessMessage string
	Error          string
}

func renderTemplateForContactTab(client DeputyHubContactInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)

		deputyId, _ := strconv.Atoi(routeVars["id"])

		contactList, err := client.GetDeputyContacts(ctx, deputyId)
		if err != nil {
			return err
		}

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "newContact":
			successMessage = "Contact added"
		default:
			successMessage = ""
		}

		vars := ListContactsVars{
			Path:           r.URL.Path,
			XSRFToken:      ctx.XSRFToken,
			ContactList:    contactList,
			DeputyDetails:  deputyDetails,
			SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}