package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
)

func (c *Client) GetTask(ctx Context, taskId int) (model.Task, error) {
	var t model.Task
	requestURL := fmt.Sprintf(SupervisionAPIPath+"/v1/tasks/%d", taskId)
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
