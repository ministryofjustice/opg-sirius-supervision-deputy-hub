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

	return taskTypes.TaskTypes, err
}
