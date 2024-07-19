package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type Documents interface {
	GetDeputyDocuments(ctx sirius.Context, deputyId int) (sirius.DocumentList, error)
}

type DocumentsVars struct {
	DocumentList   sirius.DocumentList
	SuccessMessage string
	AppVars
}

type ListDocumentsHandler struct {
	router
}

func (h *ListDocumentsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	v.PageName = "Documents"

	var successMessage string
	switch r.URL.Query().Get("success") {
	case "addDocument":
		filename := r.URL.Query().Get("filename")
		successMessage = fmt.Sprintf("Document %s added", filename)
	}

	documentList, err := h.Client().GetDeputyDocuments(ctx, v.DeputyId())
	if err != nil {
		return err
	}

	vars := DocumentsVars{
		AppVars:        v,
		DocumentList:   documentList,
		SuccessMessage: successMessage,
	}

	return h.execute(w, r, vars, v)
}
