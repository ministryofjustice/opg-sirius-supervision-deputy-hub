package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type FirmId struct {
	FirmId int `json:"firmId"`
}

func (c *ApiClient) AssignDeputyToFirm(ctx Context, deputyId int, firmId int) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(FirmId{
		FirmId: firmId,
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/deputies/%d/firm", deputyId)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

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
