package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type DocumentsClient interface {
	GetDeputyDocuments(ctx sirius.Context, deputyId int) (sirius.DocumentList, error)
	GetDocument(ctx sirius.Context, documentId int) (model.Document, error)
}

type DocumentsVars struct {
	DocumentList   sirius.DocumentList
	SuccessMessage string
	AppVars
}

func renderTemplateForDocuments(client DocumentsClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		app.PageName = "Documents"

		var successMessage string
		switch r.URL.Query().Get("success") {
		case "addDocument":
			filename := r.URL.Query().Get("filename")
			successMessage = fmt.Sprintf("Document %s added", filename)
		}

		ctx := getContext(r)
		//routeVars := mux.Vars(r)
		//documentId, _ := strconv.Atoi(routeVars["documentId"])

		//if documentId != 0 {
		// Base64 decode the document, put it in a file and put that in the header
		//	document, err := client.GetDocument(ctx, documentId)
		//	if err != nil {
		//		return err
		//	}
		//}

		documentList, err := client.GetDeputyDocuments(ctx, app.DeputyId())
		if err != nil {
			return err
		}

		vars := DocumentsVars{
			AppVars:        app,
			DocumentList:   documentList,
			SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)

	}

}
