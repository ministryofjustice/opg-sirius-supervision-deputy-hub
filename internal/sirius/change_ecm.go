package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ChangeECM(ctx Context, changeECMForm ExecutiveCaseManagerOutgoing, deputyDetails DeputyDetails) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(ExecutiveCaseManagerOutgoing{EcmId: changeECMForm.EcmId})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/ecm", deputyDetails.ID)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return err
	}
	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
