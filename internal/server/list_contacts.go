package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type DeputyHubContactInformation interface {
	GetDeputyContacts(sirius.Context, int) (sirius.ContactList, error)
}

type ListContactsVars struct {
	SuccessMessage string
	ContactList    sirius.ContactList
	AppVars
}

func renderTemplateForContactTab(client DeputyHubContactInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		contactList, err := client.GetDeputyContacts(ctx, app.DeputyId())
		if err != nil {
			return err
		}

		contactName := r.URL.Query().Get("contactName")

		var successMessage string

		switch r.URL.Query().Get("success") {
		case "newContact":
			successMessage = "Contact added"
		case "updatedContact":
			if contactName != "" {
				successMessage = contactName + "'s details updated"
			}
		case "deletedContact":
			if contactName != "" {
				successMessage = contactName + "'s details removed"
			}
		default:
			successMessage = ""
		}

		app.PageName = "Contacts"

		vars := ListContactsVars{
			AppVars:        app,
			ContactList:    contactList,
			SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
