package sirius

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (c *Client) GetDeputyClient(ctx Context, caseRecNumber string, deputyId int) (DeputyClient, error) {
	var k DeputyClient

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/client/%s", deputyId, caseRecNumber), nil)

	if err != nil {
		return DeputyClient{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return DeputyClient{}, err
	}

	io.Copy(os.Stdout, resp.Body)

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return DeputyClient{}, ErrUnauthorized
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {

		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			fmt.Println("sirius")
			fmt.Println(err)
			fmt.Println(v.ValidationErrors)
			fmt.Println(v)
			return DeputyClient{}, ValidationError{Errors: v.ValidationErrors}
		}

		return DeputyClient{}, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&k); err != nil {
		return k, err
	}

	return k, err
}
