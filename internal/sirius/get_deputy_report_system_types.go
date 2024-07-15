package sirius

import "github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"

func (c *Client) GetDeputyReportSystemTypes(ctx Context) ([]model.RefData, error) {
	deputyReportSystemTypes, err := c.getRefData(ctx, "/deputyReportSystem")
	return deputyReportSystemTypes, err
}
