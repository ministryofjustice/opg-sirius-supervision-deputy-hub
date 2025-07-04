package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

type PageInformation struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type TaskList struct {
	Tasks      []model.Task    `json:"tasks"`
	TotalTasks int             `json:"total"`
	Pages      PageInformation `json:"pages"`
}

func (c *Client) GetTasks(ctx Context, deputyId int) (TaskList, error) {
	var t TaskList

	requestURL := fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/tasks?filter=status:Not+started&sort=dueDate:asc", deputyId)
	req, err := c.newRequest(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return t, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return t, err
	}

	defer unchecked(resp.Body.Close)

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
