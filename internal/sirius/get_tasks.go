package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type TaskList struct {
	Tasks      []model.Task `json:"tasks"`
	TotalTasks int          `json:"total"`
	Pages      struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"pages"`
}

func (c *Client) GetTasks(ctx Context, deputyId string) (TaskList, error) {
	var t TaskList

	requestURL := fmt.Sprintf("/api/v1/deputies/%s/tasks", deputyId)
	req, err := c.newRequest(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return t, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return t, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return t, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return t, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return t, err
	}

	return t, err
}
