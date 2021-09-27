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
	AddNote(ctx sirius.Context, title, note string, deputyId, userId int) (int, error)
	UserDetails(sirius.Context) (sirius.UserDetails, error)
}

type deputyHubNotesVars struct {
	Path          string
	XSRFToken     string
	DeputyDetails sirius.DeputyDetails
	DeputyNotes   sirius.DeputyNoteList
	Error         string
	Errors        sirius.ValidationErrors
}

type addNoteVars struct {
	Path      	string
	XSRFToken 	string
	Title      	string
	Note   		string
	Success   	bool
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
				}

				return tmpl.ExecuteTemplate(w, "page", vars)

			case http.MethodPost:

				var (
					title   = r.PostFormValue("title")
					note  	= r.PostFormValue("note")
				)

				userId, err := client.UserDetails(ctx)
				if err != nil {
					return err
				}

				id, err := client.AddNote(ctx, title, note, deputyId, userId.ID)

				if verr, ok := err.(sirius.ValidationError); ok {

					vars := addNoteVars{
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

				//why error?
				return RedirectError(fmt.Sprintf("/deputy/%d/notes", id))

		default:
				return StatusError(http.StatusMethodNotAllowed)
			}
	}
}
