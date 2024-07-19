package sirius

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
)

type SupervisionLevelList struct {
	SupervisionLevel []struct {
		Handle     string `json:"handle"`
		Label      string `json:"label"`
		Deprecated bool   `json:"deprecated"`
	} `json:"supervisionLevel"`
}

func (c *ApiClient) GetSupervisionLevels(ctx Context) ([]model.RefData, error) {
	supervisionLevels, err := c.getRefData(ctx, "?filter=supervisionLevel")
	return supervisionLevels, err
}
