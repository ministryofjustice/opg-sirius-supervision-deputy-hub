package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type DocumentsClient interface {
	GetDeputyDocuments(ctx sirius.Context, deputyId int) (*[]model.Document, error)
}

type DocumentsVars struct {
	Documents      []model.Document
	SuccessMessage string
	AppVars
}

func renderTemplateForDocuments(client DocumentsClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)

		documents, err := client.GetDeputyDocuments(ctx, app.DeputyId())
		if err != nil {
			return err
		}

		//taskType := r.URL.Query().Get("taskType")

		//var successMessage string
		//switch r.URL.Query().Get("success") {
		//case "add":
		//	successMessage = fmt.Sprintf("%s task added", taskType)
		//case "manage":
		//	successMessage = fmt.Sprintf("%s task updated", taskType)
		//case "complete":
		//	successMessage = fmt.Sprintf("%s task completed", taskType)
		//default:
		//	successMessage = ""
		//}
		//
		//taskList, err := client.GetTasks(ctx, app.DeputyId())
		//if err != nil {
		//	return err
		//}
		//
		//app.PageName = "Deputy tasks"
		//
		vars := DocumentsVars{
			AppVars:   app,
			Documents: *documents,
			//SuccessMessage: successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)

	}

}
