package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

func (c *ApiClient) GetVisitOutcomeTypes(ctx Context) ([]model.VisitOutcomeType, error) {
	var v []model.VisitOutcomeType

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/reference-data/visitOutcome", nil)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	return v, err
}
