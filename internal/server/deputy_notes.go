package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
	"strings"
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
	Error         string
	Errors        sirius.ValidationErrors
	Success		  bool
	SuccessMessage string
}

type addNoteVars struct {
	Path      	string
	XSRFToken 	string
	Title      	string
	Note   		string
	Success   	bool
	Error    	sirius.ValidationError
	Errors    	sirius.ValidationErrors
	DeputyDetails sirius.DeputyDetails
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

				//check if there is a success in url
				urlTrim := strings.TrimPrefix(r.URL.String(), "/deputy/" + strconv.Itoa(deputyId) + "/notes" )
				var hasSuccess bool
				if urlTrim == "?success=true" {
					hasSuccess = true
				} else {
					hasSuccess = false
				}

				vars := deputyHubNotesVars{
					Path:          r.URL.Path,
					XSRFToken:     ctx.XSRFToken,
					DeputyDetails: deputyDetails,
					DeputyNotes:   deputyNotes,
					Success: hasSuccess,
					SuccessMessage: "Note added",
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

				deputyDetails, err := client.GetDeputyDetails(ctx, deputyId)
				if err != nil {
					return err
				}

				err = client.AddNote(ctx, title, note, deputyId, userId.ID)


				if verr, ok := err.(sirius.ValidationError); ok {

					verr.Errors = renameValidationErrorMessages(verr.Errors)

					vars = addNoteVars{
						Path:      r.URL.Path,
						XSRFToken: ctx.XSRFToken,
						Title:      title,
						Note:   note,
						Errors:    verr.Errors,
						DeputyDetails: deputyDetails,
					}

					w.WriteHeader(http.StatusBadRequest)
					return tmpl.ExecuteTemplate(w, "page", vars)
				} else if err != nil {
					return err
				}

				return RedirectError(fmt.Sprintf("/deputy/%d/notes?success=true", deputyId))

		default:
				return StatusError(http.StatusMethodNotAllowed)
			}
	}
}

func renameValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	errorCollection := sirius.ValidationErrors{}

	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)
			err[errorType] = errorMessage

			if fieldName == "name" && errorType == "stringLengthTooLong" {
				err[errorType] = "The title must be 255 characters or fewer"
			}

			if fieldName == "name" && errorType == "isEmpty" {
				err[errorType] = "Enter a title for the note"
			}

			if fieldName == "description" && errorType == "stringLengthTooLong" {
				err[errorType] = "The note must be 1000 characters or fewer"
			}

			if fieldName == "description" && errorType == "isEmpty" {
				err[errorType] = "Enter a note"
			}

			errorCollection[fieldName] = err
		}
	}
	return errorCollection
}
