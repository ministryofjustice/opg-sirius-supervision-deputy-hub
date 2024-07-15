package sirius

import "github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"

func (c *Client) GetDeputyBooleanTypes(ctx Context) ([]model.RefData, error) {
	deputyBooleanTypes, err := c.getRefData(ctx, "/deputyBooleanType")
	return deputyBooleanTypes, err
}
