package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

func (c *Client) GetDeputyAnnualInvoiceBillingTypes(ctx Context) ([]model.RefData, error) {
	deputyAnnualBillingInvoiceTypes, err := c.getRefData(ctx, "/annualBillingInvoice")
	return deputyAnnualBillingInvoiceTypes, err
}
