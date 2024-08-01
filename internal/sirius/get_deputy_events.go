package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"net/http"
	"sort"
	"strings"
	"time"
)

type DeputyEvents []model.DeputyEvent

func (c *ApiClient) GetDeputyEvents(ctx Context, deputyId int) (DeputyEvents, error) {
	var de DeputyEvents

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/timeline/%d", deputyId), nil)

	if err != nil {
		return de, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return de, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return de, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return de, newStatusError(resp)
	}
	err = json.NewDecoder(resp.Body).Decode(&de)

	var (
		taskTypes TaskTypeMap
		terr      error
	)
	if includesTaskEvent(de) {
		taskTypes, terr = c.getTaskTypesMap(ctx)
		if terr != nil {
			return nil, terr
		}
	}

	DeputyEvents := editDeputyEvents(de, taskTypes)

	return DeputyEvents, err

}

func editDeputyEvents(events DeputyEvents, taskTypes TaskTypeMap) DeputyEvents {
	var list DeputyEvents
	for _, e := range events {
		event := model.DeputyEvent{
			Timestamp:  FormatDateTime(IsoDateTime, e.Timestamp, SiriusDateTime),
			EventType:  reformatEventType(e.EventType),
			ID:         e.ID,
			User:       e.User,
			Event:      updateTaskInfo(e.Event, taskTypes),
			IsNewEvent: isNewEvent(e.Event.Changes),
		}

		list = append(list, event)
	}
	sortByTimelineAsc(list)
	return list
}

func isNewEvent(changes []model.Changes) bool {
	var isNew bool
	for _, c := range changes {
		isNew = c.OldValue == ""
	}
	return isNew
}

func reformatEventType(s string) string {
	eventTypeArray := strings.Split(s, "\\")
	return eventTypeArray[len(eventTypeArray)-1]
}

func sortByTimelineAsc(events DeputyEvents) DeputyEvents {
	sort.Slice(events, func(i, j int) bool {
		iTime, _ := time.Parse(SiriusDateTime, events[i].Timestamp)
		jTime, _ := time.Parse(SiriusDateTime, events[j].Timestamp)
		return jTime.Before(iTime)
	})
	return events
}

func includesTaskEvent(events DeputyEvents) bool {
	for _, e := range events {
		if e.Event.TaskType > "" {
			return true
		}
	}
	return false
}

func updateTaskInfo(event model.Event, taskTypes TaskTypeMap) model.Event {
	if event.TaskType > "" {
		event.TaskType = taskTypes[event.TaskType].Description
		event.DueDate = FormatDateTime(IsoDateTime, event.DueDate, SiriusDate)
	}
	return event
}
