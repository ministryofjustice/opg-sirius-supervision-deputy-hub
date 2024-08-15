package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetDeputyClient(ctx Context, caseRecNumber string, deputyId int) (DeputyClient, error) {
	var k []DeputyClient

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/client/%s", deputyId, caseRecNumber), nil)

	if err != nil {
		return DeputyClient{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return DeputyClient{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return DeputyClient{}, ErrUnauthorized
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {

		if err := json.NewDecoder(resp.Body).Decode(&k); err != nil {
			if len(k) == 0 {
				validationErrors := ValidationErrors{
					"caseRecNumber": {
						"": "Case number not recognised",
					},
				}
				err = ValidationError{
					Errors: validationErrors,
				}
			}
			return DeputyClient{}, err
		}

		return DeputyClient{}, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&k); err != nil {
		return k[0], err
	}

	return k[0], err
}
