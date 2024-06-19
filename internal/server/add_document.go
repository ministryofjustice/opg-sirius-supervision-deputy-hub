package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"io"
	"net/http"
	"os"
)

type AddDocumentClient interface {
	GetDeputyDocuments(ctx sirius.Context, deputyId int) (sirius.DocumentList, error)
	GetDocument(ctx sirius.Context, documentId int) (model.Document, error)
}

type AddDocumentVars struct {
	SuccessMessage string
	AppVars
}

func renderTemplateForAddDocument(client AddDocumentClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		app.PageName = "Add a document"
		if r.Method != http.MethodGet {
			//routeVars := mux.Vars(r)
			//deputyId, _ := strconv.Atoi(routeVars["id"])

			// Specify max file size to 100mb
			err := r.ParseMultipartForm(100 << 20)
			if err != nil {
				fmt.Println("Error Parsing the Form")
				fmt.Println(err)
			}

			file, handler, err := r.FormFile("document-upload")
			if err != nil {
				fmt.Println("Error Retrieving the File")
				fmt.Println(err)
			}

			defer file.Close()

			//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
			//fmt.Printf("File Size: %+v\n", handler.Size)
			//fmt.Printf("MIME Header: %+v\n", handler.Header)

			tempFile, err := os.OpenFile("./temp-files/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
			}
			defer tempFile.Close()
			io.Copy(tempFile, file)

			return StatusError(http.StatusMethodNotAllowed)
		}

		vars := AddDocumentVars{
			AppVars: app,
			//SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)

	}

}
