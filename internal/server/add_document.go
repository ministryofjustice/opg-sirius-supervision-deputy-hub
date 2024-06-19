package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"mime/multipart"
	"net/http"
)

type AddDocumentClient interface {
	AddDocument(ctx sirius.Context, file multipart.File, documentType string, direction string, date string, notes string) error
}

type AddDocumentVars struct {
	SuccessMessage string
	AppVars
}

func renderTemplateForAddDocument(client AddDocumentClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		app.PageName = "Add a document"

		vars := AddDocumentVars{
			AppVars: app,
			//SuccessMessage: successMessage,
		}

		if r.Method == http.MethodPost {
			vars.Errors = sirius.ValidationErrors{}

			// Specify max file size to 100mb
			err := r.ParseMultipartForm(100 << 20)
			if err != nil {
				fmt.Println("Error Parsing the Form")
				fmt.Println(err)
			}

			file, _, err := r.FormFile("document-upload")
			if err != nil {
				fmt.Println(err)
				vars.Errors["document-upload"] = map[string]string{"": "Error uploading the file"}
			}

			//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
			//fmt.Printf("File Size: %+v\n", handler.Size)
			//fmt.Printf("MIME Header: %+v\n", handler.Header)

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

			//defer file.Close()

			// Temporarily upload it to temp-files
			// Then, pass the filename to the sirius side, who can upload it as a request
			// It'd be better if we can pass this formFile directly to the sirius side to add to a new request, but we'll see
			//tempFile, err := os.OpenFile("./temp-files/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			//if err != nil {
			//	fmt.Println(err)
			//}
			//
			//defer tempFile.Close()
			//io.Copy(tempFile, file)

			fmt.Println(documentType)
			fmt.Println(notes)

			ctx := getContext(r)
			err = client.AddDocument(ctx, file, documentType, direction, date, notes)

			if err != nil {
				panic(err)
			}

			return StatusError(http.StatusMethodNotAllowed)
		}

		return tmpl.ExecuteTemplate(w, "page", vars)

	}

}
