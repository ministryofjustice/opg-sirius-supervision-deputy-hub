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

//type T struct {
//	Limit    int           `json:"limit"`
//	Metadata []interface{} `json:"metadata"`
//	Pages    struct {
//		Current int `json:"current"`
//		Total   int `json:"total"`
//	} `json:"pages"`
//	Total          int `json:"total"`
//	TimelineEvents []struct {
//		Id        int    `json:"id"`
//		Hash      string `json:"hash"`
//		Timestamp string `json:"timestamp"`
//		EventType string `json:"eventType"`
//		User      struct {
//			Id          int    `json:"id"`
//			PhoneNumber string `json:"phoneNumber"`
//			DisplayName string `json:"displayName"`
//			Email       string `json:"email"`
//		} `json:"user"`
//		Event struct {
//			IsCaseEvent          bool    `json:"isCaseEvent,omitempty"`
//			IsPersonEvent        bool    `json:"isPersonEvent,omitempty"`
//			TaskId               int     `json:"taskId,omitempty"`
//			TaskType             string  `json:"taskType,omitempty"`
//			DueDate              string  `json:"dueDate,omitempty"`
//			Description          string  `json:"description,omitempty"`
//			Name                 string  `json:"name,omitempty"`
//			AssigneeId           string  `json:"assigneeId,omitempty"`
//			Assignee             string  `json:"assignee,omitempty"`
//			IsCaseOwnerTask      bool    `json:"isCaseOwnerTask,omitempty"`
//			PersonType           string  `json:"personType"`
//			PersonId             string  `json:"personId"`
//			PersonUid            string  `json:"personUid"`
//			PersonName           string  `json:"personName"`
//			PersonCourtRef       *string `json:"personCourtRef"`
//			OrderType            string  `json:"orderType,omitempty"`
//			OrderUid             string  `json:"orderUid,omitempty"`
//			OrderId              string  `json:"orderId,omitempty"`
//			OrderCourtRef        string  `json:"orderCourtRef,omitempty"`
//			CourtReferenceNumber string  `json:"courtReferenceNumber,omitempty"`
//			CourtReference       string  `json:"courtReference,omitempty"`
//			Changes              []struct {
//				FieldName string      `json:"fieldName"`
//				OldValue  interface{} `json:"oldValue"`
//				NewValue  string      `json:"newValue"`
//				Type      string      `json:"type"`
//			} `json:"changes,omitempty"`
//			DeputyName           string `json:"deputyName,omitempty"`
//			ExecutiveCaseManager string `json:"executiveCaseManager,omitempty"`
//		} `json:"event"`
//	} `json:"timelineEvents"`
//}

type DeputyEvents []model.DeputyEvent
type TimelineList struct {
	Limit    int           `json:"limit"`
	Metadata []interface{} `json:"metadata"`
	Pages    struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"pages"`
	Total        int                 `json:"total"`
	DeputyEvents []model.DeputyEvent `json:"timelineEvents"`
}

func (c *Client) GetDeputyEvents(ctx Context, deputyId int, pageNumber int, timelineEventsPerPage int) (TimelineList, error) {
	var de TimelineList

	endpoint := fmt.Sprintf(
		"/api/v1/timeline/%d/deputy?limit=%d&page=%d",
		deputyId,
		timelineEventsPerPage,
		pageNumber,
	)

	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)

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
	if includesTaskEvent(de.DeputyEvents) {
		taskTypes, terr = c.getTaskTypesMap(ctx)
		if terr != nil {
			return TimelineList{}, terr
		}
	}

	de.DeputyEvents = editDeputyEvents(de.DeputyEvents, taskTypes)
	return de, err

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
