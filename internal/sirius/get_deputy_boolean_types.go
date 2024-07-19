package sirius

import "github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"

type DeputyBooleanTypes struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (c *ApiClient) GetDeputyBooleanTypes(ctx Context) ([]model.RefData, error) {
	deputyBooleanTypes, err := c.getRefData(ctx, "/deputyBooleanType")
	return deputyBooleanTypes, err
}
