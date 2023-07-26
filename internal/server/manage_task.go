package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type ManageTasks interface {
	GetTask(sirius.Context, int) (model.Task, error)
	//UpdateTask
	//ReassignTask
}

type manageTaskVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	TaskDetails    model.Task
	Error          string
	Errors         sirius.ValidationErrors
	Success        bool
	SuccessMessage string
}

func renderTemplateForManageTasks(client ManageTasks, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		//deputyId, _ := strconv.Atoi(routeVars["id"])
		taskId, _ := strconv.Atoi(routeVars["taskId"])

		taskDetails, err := client.GetTask(ctx, taskId)
		if err != nil {
			return err
		}

		taskDetails.DueDate = sirius.FormatDateTime(sirius.SiriusDate, taskDetails.DueDate, sirius.IsoDate)

		switch r.Method {
		case http.MethodGet:

			vars := manageTaskVars{
				Path:          r.URL.Path,
				XSRFToken:     ctx.XSRFToken,
				DeputyDetails: deputyDetails,
				TaskDetails:   taskDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		//case http.MethodPost:
		//	var vars changeFirmVars
		//	newFirm := r.PostFormValue("select-firm")
		//	AssignToExistingFirmStringIdValue := r.PostFormValue("select-existing-firm")
		//
		//	if newFirm == "new-firm" {
		//		return Redirect(fmt.Sprintf("/%d/add-firm", deputyId))
		//	}
		//
		//	AssignToFirmId := 0
		//	if AssignToExistingFirmStringIdValue != "" {
		//		AssignToFirmId, err = strconv.Atoi(AssignToExistingFirmStringIdValue)
		//		if err != nil {
		//			return err
		//		}
		//	}
		//
		//	assignDeputyToFirmErr := client.AssignDeputyToFirm(ctx, deputyId, AssignToFirmId)
		//
		//	if verr, ok := assignDeputyToFirmErr.(sirius.ValidationError); ok {
		//		vars = changeFirmVars{
		//			Path:      r.URL.Path,
		//			XSRFToken: ctx.XSRFToken,
		//			Errors:    verr.Errors,
		//		}
		//
		//		return tmpl.ExecuteTemplate(w, "page", vars)
		//	} else if err != nil {
		//		return err
		//	}
		//	return Redirect(fmt.Sprintf("/%d?success=firm", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
