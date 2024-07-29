package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"golang.org/x/sync/errgroup"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type ReplaceDocumentClient interface {
	ReplaceDocument(ctx sirius.Context, file multipart.File, filename, documentType, direction, date, notes string, deputyId, documentId int) error
	GetDocumentDirections(ctx sirius.Context) ([]model.RefData, error)
	GetDocumentTypes(ctx sirius.Context) ([]model.RefData, error)
	GetDocumentById(ctx sirius.Context, deputyId, documentId int) (model.Document, error)
}

type ReplaceDocumentVars struct {
	SuccessMessage   string
	OriginalDocument model.Document
	AppVars
	DocumentDirectionRefData []model.RefData
	DocumentTypes            []model.RefData
	DocumentType             string
	Direction                string
	Date                     string
	Notes                    string
}

func renderTemplateForReplaceDocument(client ReplaceDocumentClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		app.PageName = "Replace a document"

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		documentId, _ := strconv.Atoi(routeVars["documentId"])

		vars := ReplaceDocumentVars{
			AppVars: app,
			Date:    time.Now().Format("2006-01-02"),
		}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			documentDirectionRefData, err := client.GetDocumentDirections(ctx.With(groupCtx))
			if err != nil {
				return err
			}
			vars.DocumentDirectionRefData = documentDirectionRefData
			return nil
		})

		fmt.Println("after get doc directions")

		group.Go(func() error {
			documentTypes, err := client.GetDocumentTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}
			vars.DocumentTypes = documentTypes
			return nil
		})

		fmt.Println("after get doc types")

		group.Go(func() error {
			originalDocument, err := client.GetDocumentById(ctx.With(groupCtx), vars.DeputyDetails.ID, documentId)
			if err != nil {
				return err
			}
			if (originalDocument != model.Document{}) {
				vars.OriginalDocument = originalDocument
				newTime, err := time.Parse("02/01/2006 15:04:05", originalDocument.ReceivedDateTime)
				if err != nil {
					return err
				}
				vars.OriginalDocument.ReformattedTime = newTime.Format("02/01/2006")
			}
			return nil
		})

		fmt.Println("after get doc by id")

		if err := group.Wait(); err != nil {
			return err
		}
		fmt.Println("before post")

		if r.Method == http.MethodPost {
			fmt.Println("in post")

			vars.Errors = sirius.ValidationErrors{}

			// Specify max file size to 100mb
			err := r.ParseMultipartForm(100 << 20)
			if err != nil {
				return err
			}

			file, handler, err := r.FormFile("document-upload")
			if err != nil {
				vars.Errors["document-upload"] = map[string]string{"": "Select a file to attach"}
			}
			fmt.Println("after select file to attach")

			documentType := r.PostFormValue("documentType")
			direction := r.PostFormValue("documentDirection")
			date := r.PostFormValue("documentDate")
			notes := r.PostFormValue("notes")

			vars.DocumentType = documentType
			vars.Direction = direction
			vars.Date = date
			vars.Notes = notes

			if len(vars.Errors) > 0 {
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			fmt.Println("after if errors more than 0")

			ctx := getContext(r)
			err = client.ReplaceDocument(ctx, file, handler.Filename, documentType, direction, date, notes, vars.DeputyDetails.ID, vars.OriginalDocument.Id)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			fmt.Println("after validation errors")

			if err != nil {
				return err
			}

			fmt.Println("after errors")

			return Redirect(fmt.Sprintf("/%d/documents?success=replaceDocument&previousFilename=%s&filename=%s", app.DeputyId(), vars.OriginalDocument.FriendlyDescription, handler.Filename))
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}

}
