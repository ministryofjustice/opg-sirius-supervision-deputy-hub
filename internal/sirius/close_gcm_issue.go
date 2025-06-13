package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type GcmIssueIds struct {
	Ids []string `json:"gcmIssueIds"`
}

func (c *Client) CloseGCMIssues(ctx Context, gcmIssueIds []string) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(GcmIssueIds{
		Ids: gcmIssueIds,
	})

	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPut, SupervisionAPIPath + "/v1/gcm-issues/close", &body)

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

	return nil
}
