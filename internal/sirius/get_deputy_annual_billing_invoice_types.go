package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

type DeputyAnnualBillingInvoiceTypes struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (c *ApiClient) GetDeputyAnnualInvoiceBillingTypes(ctx Context) ([]model.RefData, error) {
	deputyAnnualBillingInvoiceTypes, err := c.getRefData(ctx, "/annualBillingInvoice")
	return deputyAnnualBillingInvoiceTypes, err
}
