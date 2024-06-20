package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"mime/multipart"
	"net/http"
)

type AddDocumentClient interface {
	AddDocument(ctx sirius.Context, file multipart.File, filename string, documentType string, direction string, date string, notes string, deputyId int) error
	GetRefData(ctx sirius.Context, refDataUrlType string) ([]model.RefData, error)
}

type AddDocumentVars struct {
	SuccessMessage string
	AppVars
	DocumentDirectionRefData []model.RefData
	DocumentTypes            []model.RefData
}

func renderTemplateForAddDocument(client AddDocumentClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		app.PageName = "Add a document"

		vars := AddDocumentVars{
			AppVars: app,
		}

		documentDirectionRefData, err := client.GetRefData(getContext(r), "/documentDirection")
		if err != nil {
			return err
		}
		vars.DocumentDirectionRefData = documentDirectionRefData

		documentTypes, err := client.GetRefData(getContext(r), "?filter=noteType:deputy")
		if err != nil {
			return err
		}
		vars.DocumentTypes = documentTypes

		if r.Method == http.MethodPost {
			vars.Errors = sirius.ValidationErrors{}

			// Specify max file size to 100mb
			err := r.ParseMultipartForm(100 << 20)
			if err != nil {
				fmt.Println("Error Parsing the Form")
				fmt.Println(err)
			}

			file, handler, err := r.FormFile("document-upload")
			if err != nil {
				fmt.Println(err)
				vars.Errors["document-upload"] = map[string]string{"": "Error uploading the file"}
			}

			documentType := r.PostFormValue("type")
			direction := r.PostFormValue("direction")

			date := r.PostFormValue("date")
			notes := r.PostFormValue("notes")

			if direction == "" {
				vars.Errors["direction"] = map[string]string{"": "Select a direction"}
			}

			if date == "" {
				vars.Errors["date"] = map[string]string{"": "Select a date"}
			}

			if len(vars.Errors) > 0 {
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			ctx := getContext(r)
			err = client.AddDocument(ctx, file, handler.Filename, documentType, direction, date, notes, vars.DeputyDetails.ID)

			if err != nil {
				fmt.Println(err)
			}

			return Redirect(fmt.Sprintf("/%d/documents?success=addDocument", app.DeputyId()))
		}

		return tmpl.ExecuteTemplate(w, "page", vars)

	}

}
