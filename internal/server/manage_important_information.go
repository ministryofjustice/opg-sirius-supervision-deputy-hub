package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

type editProImportantInformation struct {
	AnnualBillingInvoiceTypes []model.RefData
	DeputyBooleanTypes        []model.RefData
	DeputyReportSystemTypes   []model.RefData
	AppVars
}

type EditProImportantInformationHandler struct {
	router
}

func (h *EditProImportantInformationHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	v.PageName = "Manage important information"

	vars := editProImportantInformation{AppVars: v}

	group, groupCtx := errgroup.WithContext(ctx.Context)

	group.Go(func() error {
		annualBillingInvoiceTypes, err := h.Client().GetDeputyAnnualInvoiceBillingTypes(ctx.With(groupCtx))
		if err != nil {
			return err
		}

		vars.AnnualBillingInvoiceTypes = annualBillingInvoiceTypes
		return nil
	})

	group.Go(func() error {
		deputyBooleanTypes, err := h.Client().GetDeputyBooleanTypes(ctx.With(groupCtx))
		if err != nil {
			return err
		}

		vars.DeputyBooleanTypes = deputyBooleanTypes
		return nil
	})

	group.Go(func() error {
		deputyReportSystemTypes, err := h.Client().GetDeputyReportSystemTypes(ctx.With(groupCtx))
		if err != nil {
			return err
		}

		vars.DeputyReportSystemTypes = deputyReportSystemTypes
		return nil
	})

	if err := group.Wait(); err != nil {
		return err
	}

	switch r.Method {
	case http.MethodGet:
		return h.execute(w, r, vars)

	case http.MethodPost:
		var panelDeputyBool bool
		var err error

		if r.PostFormValue("panel-deputy") != "" {
			panelDeputyBool, err = strconv.ParseBool(r.PostFormValue("panel-deputy"))
			if err != nil {
				return err
			}
		}

		annualBillingInvoice := vars.AppVars.DeputyDetails.DeputyImportantInformation.AnnualBillingInvoice.Handle
		if r.PostFormValue("annual-billing") != "" {
			annualBillingInvoice = r.PostFormValue("annual-billing")
		} else if annualBillingInvoice == "" {
			annualBillingInvoice = "UNKNOWN"
		}

		importantInfoForm := sirius.ImportantInformationDetails{
			DeputyType:                vars.AppVars.DeputyType(),
			Complaints:                r.PostFormValue("complaints"),
			PanelDeputy:               panelDeputyBool,
			AnnualBillingInvoice:      annualBillingInvoice,
			OtherImportantInformation: r.PostFormValue("other-info-note"),
			MonthlySpreadsheet:        r.PostFormValue("monthly-spreadsheet"),
			IndependentVisitorCharges: r.PostFormValue("independent-visitor-charges"),
			BankCharges:               r.PostFormValue("bank-charges"),
			APAD:                      r.PostFormValue("apad"),
			ReportSystem:              r.PostFormValue("report-system"),
		}

		err = h.Client().UpdateImportantInformation(ctx, v.DeputyId(), importantInfoForm)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)
			return h.execute(w, r, vars)
		} else if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d?success=importantInformation", v.DeputyId()))
	default:
		return StatusError(http.StatusMethodNotAllowed)
	}
}
