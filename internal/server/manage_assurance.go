package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type ManageAssuranceClient interface {
	UpdateAssurance(ctx sirius.Context, manageAssuranceForm sirius.UpdateAssuranceDetails, deputyId, visitId int) error
	GetVisitors(ctx sirius.Context) ([]model.Visitor, error)
	GetRagRatingTypes(ctx sirius.Context) ([]model.RAGRating, error)
	GetVisitOutcomeTypes(ctx sirius.Context) ([]model.RefData, error)
	GetPdrOutcomeTypes(ctx sirius.Context) ([]model.RefData, error)
	GetAssuranceById(ctx sirius.Context, deputyId int, visitId int) (model.Assurance, error)
}

type ManageAssuranceVars struct {
	Visitors          []model.Visitor
	RagRatingTypes    []model.RAGRating
	VisitOutcomeTypes []model.RefData
	PdrOutcomeTypes   []model.RefData
	Assurance         model.Assurance
	ErrorNote         string
	AppVars
}

func parseAssuranceForm(assuranceForm sirius.UpdateAssuranceDetails) model.Assurance {
	return model.Assurance{
		CommissionedDate:   assuranceForm.CommissionedDate,
		VisitorAllocated:   assuranceForm.VisitorAllocated,
		ReportDueDate:      assuranceForm.ReportDueDate,
		ReportReceivedDate: assuranceForm.ReportReceivedDate,
		VisitOutcome:       model.RefData{Label: assuranceForm.VisitOutcome},
		PdrOutcome:         model.RefData{Label: assuranceForm.PdrOutcome},
		ReportReviewDate:   assuranceForm.ReportReviewDate,
		ReportMarkedAs:     model.RAGRating{Label: assuranceForm.ReportMarkedAs},
		ReviewedBy:         model.User{ID: assuranceForm.ReviewedBy},
		Note:               assuranceForm.Note,
	}
}

//func renderTemplateForManageAssurance(client ManageAssuranceClient, visitTmpl Template, pdrTmpl Template) Handler {
//	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
//		ctx := getContext(r)
//		routeVars := mux.Vars(r)
//		visitId, _ := strconv.Atoi(routeVars["visitId"])
//
//		app.PageName = "Manage assurance visit"
//
//		vars := ManageAssuranceVars{AppVars: app}
//		tmpl := visitTmpl
//
//		group, groupCtx := errgroup.WithContext(ctx.Context)
//
//		group.Go(func() error {
//			visit, err := client.GetAssuranceById(ctx, app.DeputyId(), visitId)
//			if err != nil {
//				return err
//			}
//			vars.Assurance = visit
//			if visit.Type.Handle == "PDR" {
//				tmpl = pdrTmpl
//				vars.AppVars.PageName = "Manage PDR"
//			}
//			return nil
//		})
//
//		group.Go(func() error {
//			visitors, err := client.GetVisitors(ctx)
//			if err != nil {
//				return err
//			}
//			vars.Visitors = visitors
//
//			return nil
//		})
//
//		group.Go(func() error {
//			ragRatingTypes, err := client.GetRagRatingTypes(ctx.With(groupCtx))
//			if err != nil {
//				return err
//			}
//
//			vars.RagRatingTypes = ragRatingTypes
//			return nil
//		})
//
//		group.Go(func() error {
//			visitOutcomeTypes, err := client.GetVisitOutcomeTypes(ctx.With(groupCtx))
//			if err != nil {
//				return err
//			}
//
//			vars.VisitOutcomeTypes = visitOutcomeTypes
//			return nil
//		})
//
//		group.Go(func() error {
//			pdrOutcomeTypes, err := client.GetPdrOutcomeTypes(ctx.With(groupCtx))
//			if err != nil {
//				return err
//			}
//
//			vars.PdrOutcomeTypes = pdrOutcomeTypes
//			return nil
//		})
//
//		if err := group.Wait(); err != nil {
//			return err
//		}
//
//		switch r.Method {
//		case http.MethodGet:
//			return tmpl.ExecuteTemplate(w, "page", vars)
//
//		case http.MethodPost:
//			reportReviewDate := r.PostFormValue("report-review-date")
//			reviewedBy := 0
//			if reportReviewDate != "" {
//				reviewedBy = app.UserDetails.ID
//			}
//
//			pdrOutcome := ""
//			if r.PostFormValue("pdr-outcome") == "Not received" {
//				pdrOutcome = "NOT_RECEIVED"
//			} else if r.PostFormValue("pdr-outcome") == "Received" {
//				pdrOutcome = "RECEIVED"
//			}
//
//			manageAssuranceForm := sirius.UpdateAssuranceDetails{
//				CommissionedDate:   r.PostFormValue("commissioned-date"),
//				VisitorAllocated:   r.PostFormValue("visitor-allocated"),
//				ReportDueDate:      r.PostFormValue("report-due-date"),
//				ReportReceivedDate: r.PostFormValue("report-received-date"),
//				VisitOutcome:       r.PostFormValue("visit-outcome"),
//				PdrOutcome:         pdrOutcome,
//				ReportReviewDate:   reportReviewDate,
//				ReportMarkedAs:     r.PostFormValue("visit-report-marked-as"),
//				ReviewedBy:         reviewedBy,
//				Note:               strings.TrimSpace(r.PostFormValue("note")),
//			}
//
//			err := client.UpdateAssurance(ctx, manageAssuranceForm, app.DeputyId(), visitId)
//
//			if verr, ok := err.(sirius.ValidationError); ok {
//				vars.Errors = util.RenameErrors(verr.Errors)
//				vars.ErrorNote = r.PostFormValue("note")
//				vars.Assurance = parseAssuranceForm(manageAssuranceForm)
//
//				return tmpl.ExecuteTemplate(w, "page", vars)
//			}
//
//			success := "manageVisit"
//			if vars.Assurance.Type.Handle == "PDR" {
//				success = "managePDR"
//			}
//
//			return Redirect(fmt.Sprintf("/%d/assurances?success=%s", app.DeputyId(), success))
//		default:
//			return StatusError(http.StatusMethodNotAllowed)
//		}
//	}
//}
