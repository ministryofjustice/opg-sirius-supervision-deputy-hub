package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

type DeputyReportSystemTypes struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (c *ApiClient) GetDeputyReportSystemTypes(ctx Context) ([]model.RefData, error) {
	deputyReportSystemTypes, err := c.getRefData(ctx, "/deputyReportSystem")
	return deputyReportSystemTypes, err
}
