package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type SupervisionLevelList struct {
	SupervisionLevel []struct {
		Handle     string `json:"handle"`
		Label      string `json:"label"`
		Deprecated bool   `json:"deprecated"`
	} `json:"supervisionLevel"`
}

func (c *ApiClient) GetSupervisionLevels(ctx Context) ([]model.RefData, error) {
	endpoint := "/api/v1/reference-data?filter=supervisionLevel"

	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v SupervisionLevelList
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	var supervisionLevels []model.RefData
	for _, u := range v.SupervisionLevel {
		supervisionLevels = append(supervisionLevels, u)
	}

	return supervisionLevels, nil
}
