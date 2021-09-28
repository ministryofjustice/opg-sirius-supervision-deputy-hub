package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type DeputyHubNotesInformation interface {
	GetDeputyDetails(sirius.Context, int) (sirius.DeputyDetails, error)
	GetDeputyNotes(sirius.Context, int) (sirius.DeputyNoteList, error)
	AddNote(ctx sirius.Context, title, note string, deputyId, userId int) error
	GetUserDetails(sirius.Context) (sirius.UserDetails, error)
}

type deputyHubNotesVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	DeputyNotes   sirius.DeputyNoteList
	ErrorMessage bool
	Error         string
	Errors        sirius.ValidationErrors
}

type addNoteVars struct {
	Path      	string
	XSRFToken 	string
	Title      	string
	Note   		string
	Success   	bool
	Error    	sirius.ValidationError
	Errors    	sirius.ValidationErrors
}

func renderTemplateForDeputyHubNotes(client DeputyHubNotesInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
			case http.MethodGet:

				deputyDetails, err := client.GetDeputyDetails(ctx, deputyId)
				if err != nil {
					return err
				}
				deputyNotes, err := client.GetDeputyNotes(ctx, deputyId)
				if err != nil {
					return err
				}

				vars := deputyHubNotesVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					DeputyDetails: deputyDetails,
					DeputyNotes:   deputyNotes,
					ErrorMessage: false,
				}

				return tmpl.ExecuteTemplate(w, "page", vars)

			case http.MethodPost:
				var vars addNoteVars
				var (
					title   = r.PostFormValue("title")
					note  	= r.PostFormValue("note")
				)

				userId, err := client.GetUserDetails(ctx)
				if err != nil {
					return err
				}
				err = client.AddNote(ctx, title, note, deputyId, userId.ID)

				if verr, ok := err.(sirius.ValidationError); ok {

					verr.Errors = replaceErrorStrings(verr.Errors)

					vars = addNoteVars{
						Path:      r.URL.Path,
						XSRFToken: ctx.XSRFToken,
						Title:      title,
						Note:   note,
						Errors:    verr.Errors,
					}

					w.WriteHeader(http.StatusBadRequest)
					return tmpl.ExecuteTemplate(w, "page", vars)
				} else if err != nil {
					return err
				}

				return RedirectError(fmt.Sprintf("/deputy/%d/notes", deputyId))

		default:
				return StatusError(http.StatusMethodNotAllowed)
			}
	}
}

func replaceErrorStrings(siriusError sirius.ValidationErrors) sirius.ValidationErrors {

	errorCollection := sirius.ValidationErrors{}

	titleLengthTooShort := make (map[string]string)
	titleLengthTooShort["titleLengthTooShort"] = "Enter a title for the note"

	titleLengthTooLong := make (map[string]string)
	titleLengthTooLong["titleLengthTooLong"] = "The title must be 255 characters or fewer"

	descriptionLengthTooShort := make (map[string]string)
	descriptionLengthTooShort["titleLengthTooShort"] = "Enter a note"

	descriptionLengthTooLong := make (map[string]string)
	descriptionLengthTooLong["titleLengthTooLong"] = "The note must be 1000 characters or fewer"


	for fieldName, value := range siriusError {
		if fieldName == "name" {
			for key, _ := range value {
				if key == "stringLengthTooLong" {
					errorCollection["name"] = titleLengthTooLong
				}
				if key == "isEmpty" {
					errorCollection["name"] = titleLengthTooShort
				}
			}
		}
		if fieldName == "description" {
			for key, _ := range value {
				if key == "stringLengthTooLong" {
					errorCollection["description"] = descriptionLengthTooLong
				}
				if key == "isEmpty" {
					errorCollection["description"] = descriptionLengthTooShort
				}
			}
		}

	}
	return errorCollection

}
