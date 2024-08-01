package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"sort"
)

type TaskTypeMap map[string]model.TaskType

type TaskTypes struct {
	TaskTypes TaskTypeMap `json:"task_types"`
}

func (c *ApiClient) GetTaskTypesForDeputyType(ctx Context, deputyType string) ([]model.TaskType, error) {
	taskTypes, err := c.getTaskTypesMap(ctx)
	if err != nil {
		return nil, err
	}

	var deputyTaskTypes []model.TaskType
	for _, t := range taskTypes {
		if t.ProDeputyTask && deputyType == "PRO" {
			deputyTaskTypes = append(deputyTaskTypes, t)
		} else if t.PaDeputyTask && deputyType == "PA" {
			deputyTaskTypes = append(deputyTaskTypes, t)
		}
	}

	sort.Slice(deputyTaskTypes, func(i, j int) bool {
		return deputyTaskTypes[i].Description < deputyTaskTypes[j].Description
	})

	return deputyTaskTypes, err
}

func (c *ApiClient) getTaskTypesMap(ctx Context) (TaskTypeMap, error) {
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

	var taskTypes TaskTypes
	err = json.NewDecoder(resp.Body).Decode(&taskTypes)

	return taskTypes.TaskTypes, err
}
