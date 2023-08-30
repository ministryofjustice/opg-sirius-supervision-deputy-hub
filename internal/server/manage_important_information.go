package server

import (
	"fmt"
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
			reportSystemType := checkForReportSystemType(r.PostFormValue("report-system"))

			annualBillingInvoice := r.PostFormValue("annual-billing")
			if annualBillingInvoice == "" {
				annualBillingInvoice = vars.AppVars.DeputyDetails.DeputyImportantInformation.AnnualBillingInvoice.Label
			}
			if annualBillingInvoice == "" {
				annualBillingInvoice = "Unknown"
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
				ReportSystem:              reportSystemType,
			}

			err = client.UpdateImportantInformation(ctx, app.DeputyId(), importantInfoForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors

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

func checkForReportSystemType(reportType string) string {
	if reportType == "OPG Digital" {
		return "OPGDigital"
	} else if reportType == "OPG Paper" {
		return "OPGPaper"
	} else {
		return reportType
	}
}
