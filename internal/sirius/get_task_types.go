package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"sort"
)

type TaskTypesMap struct {
	TaskTypes map[string]model.TaskType `json:"task_types"`
}

func (c *Client) GetTaskTypes(ctx Context, deputy DeputyDetails) ([]model.TaskType, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/tasktypes/deputy", nil)

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

	var taskTypes TaskTypesMap
	if err = json.NewDecoder(resp.Body).Decode(&taskTypes); err != nil {
		return nil, err
	}

	isPro := deputy.DeputyType.Handle == "PRO"

	var deputyTaskTypes []model.TaskType
	for _, t := range taskTypes.TaskTypes {
		if t.ProDeputyTask && isPro {
			deputyTaskTypes = append(deputyTaskTypes, t)
		} else if t.PaDeputyTask && !isPro {
			deputyTaskTypes = append(deputyTaskTypes, t)
		}
	}

	sort.Slice(deputyTaskTypes, func(i, j int) bool {
		return deputyTaskTypes[i].Handle < deputyTaskTypes[j].Handle
	})

	return deputyTaskTypes, err
}
