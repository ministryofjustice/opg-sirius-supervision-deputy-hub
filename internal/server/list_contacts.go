package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubContactInformation interface {
	GetDeputyContacts(sirius.Context, int, int, int, string, string, string) (sirius.ContactList, sirius.AriaSorting, error)
	// GetPageDetails(sirius.Context, sirius.ContactList, int, int) sirius.PageDetails
}

type ListContactsVars struct {
	Path                 string
	XSRFToken            string
	AriaSorting          sirius.AriaSorting
	DeputyDetails        sirius.DeputyDetails
	ContactList			 sirius.ContactList
	PageDetails          sirius.PageDetails
	SuccessMessage       string
	Error                string
}

func renderTemplateForContactTab(client DeputyHubContactInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		urlParams := r.URL.Query()

		deputyId, _ := strconv.Atoi(routeVars["id"])
		search, _ := strconv.Atoi(r.FormValue("page"))
		displayContactLimit, _ := strconv.Atoi(r.FormValue("limit"))
		if displayContactLimit == 0 {
			displayContactLimit = 25
		}

		columnBeingSorted, sortOrder := parseUrl(urlParams)

		contactList, ariaSorting, err := client.GetDeputyContacts(ctx, deputyId, displayContactLimit, search, deputyDetails.DeputyType.Handle, columnBeingSorted, sortOrder)
		if err != nil {
			return err
		}

		// pageDetails := client.GetPageDetails(ctx, contactList, search, displayContactLimit)

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "newContact":
			successMessage = "Contact added"
		default:
			successMessage = ""
		}

		vars := ListContactsVars{
			Path:                 r.URL.Path,
			XSRFToken:            ctx.XSRFToken,
			AriaSorting:          ariaSorting,
			ContactList:          contactList,
			// PageDetails:          pageDetails,
			DeputyDetails:        deputyDetails,
			SuccessMessage:       successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}