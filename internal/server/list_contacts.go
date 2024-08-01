package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type ListContactsVars struct {
	SuccessMessage string
	ContactList    sirius.ContactList
	AppVars
}

type ListContactsHandler struct {
	router
}

func (h *ListContactsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	contactList, err := h.Client().GetDeputyContacts(ctx, v.DeputyId())
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

	v.PageName = "Contacts"

	vars := ListContactsVars{
		AppVars:        v,
		ContactList:    contactList,
		SuccessMessage: successMessage,
	}

	return h.execute(w, r, vars)
}
