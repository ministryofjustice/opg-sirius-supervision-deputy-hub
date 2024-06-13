package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/util"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type ManageProDeputyImportantInformation interface {
	UpdateImportantInformation(sirius.Context, int, sirius.ImportantInformationDetails) error
	GetDeputyAnnualInvoiceBillingTypes(ctx sirius.Context) ([]sirius.DeputyAnnualBillingInvoiceTypes, error)
	GetDeputyBooleanTypes(ctx sirius.Context) ([]sirius.DeputyBooleanTypes, error)
	GetDeputyReportSystemTypes(ctx sirius.Context) ([]sirius.DeputyReportSystemTypes, error)
}

type manageDeputyImportantInformationVars struct {
	AnnualBillingInvoiceTypes []sirius.DeputyAnnualBillingInvoiceTypes
	DeputyBooleanTypes        []sirius.DeputyBooleanTypes
	DeputyReportSystemTypes   []sirius.DeputyReportSystemTypes
	AppVars
}

func renderTemplateForImportantInformation(client ManageProDeputyImportantInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		app.PageName = "Manage important information"

		vars := manageDeputyImportantInformationVars{AppVars: app}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			annualBillingInvoiceTypes, err := client.GetDeputyAnnualInvoiceBillingTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			vars.AnnualBillingInvoiceTypes = annualBillingInvoiceTypes
			return nil
		})

		group.Go(func() error {
			deputyBooleanTypes, err := client.GetDeputyBooleanTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			vars.DeputyBooleanTypes = deputyBooleanTypes
			return nil
		})

		group.Go(func() error {
			deputyReportSystemTypes, err := client.GetDeputyReportSystemTypes(ctx.With(groupCtx))
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
			return tmpl.ExecuteTemplate(w, "page", vars)

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

			err = client.UpdateImportantInformation(ctx, app.DeputyId(), importantInfoForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=importantInformation", app.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
