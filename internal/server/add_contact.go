package server

import (
	"net/http"
	"strconv"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ContactInformation interface {
	AddContactDetails(sirius.Context, int, sirius.ContactDetails) (error)
}

type addContactVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	Error         string
	Errors        sirius.ValidationErrors
	DeputyId      int
	Form          sirius.ContactDetails
}

func renderTemplateForAddContact(client ContactInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:
			vars := addContactVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyId:      deputyId,
				DeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)
		case http.MethodPost:
			addContactDetailForm := sirius.ContactDetails{
				ContactName:      r.PostFormValue("name"),
				JobTitle:         r.PostFormValue("job-title"),
				Email:            r.PostFormValue("email"),
				PhoneNumber:      r.PostFormValue("phone"),
				OtherPhoneNumber: r.PostFormValue("other-phone"),
				Notes:            r.PostFormValue("notes"),
				IsNamedDeputy:    r.PostFormValue("is-named-deputy"),
				IsMainContact:    r.PostFormValue("is-main-contact"),
			}

			err := client.AddContactDetails(ctx, deputyId, addContactDetailForm)

			//fmt.Println(err)

			if verr, ok := err.(sirius.ValidationError); ok {
				fmt.Println(verr.Errors)
				vars := addContactVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
					DeputyDetails: deputyDetails,
					Form: addContactDetailForm,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if(err != nil) {
				return err	
			}

			return Redirect(fmt.Sprintf("/%d/contacts?success=newContact", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}

	}
}