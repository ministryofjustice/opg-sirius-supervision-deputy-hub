package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type DeputyEvents []DeputyEvent

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"displayName"`
	PhoneNumber string `json:"phoneNumber"`
}

type Event struct {
	OrderType            string         `json:"orderType"`
	SiriusId             string         `json:"orderUid"`
	OrderNumber          string         `json:"orderCourtRef"`
	DeputyID             string         `json:"personId"`
	DeputyName           string         `json:"personName"`
	OrganisationName     string         `json:"organisationName"`
	ExecutiveCaseManager string         `json:"executiveCaseManager"`
	Changes              []Changes      `json:"changes"`
	Client               []ClientPerson `json:"additionalPersons"`
	RequestedBy          string         `json:"requestedBy"`
	RequestedDate        string         `json:"requestedDate"`
	CommissionedDate     string         `json:"commissionedDate"`
	ReportDueDate        string         `json:"reportDueDate"`
	ReportReceivedDate   string         `json:"reportReceivedDate"`
	VisitOutcome         string         `json:"assuranceVisitOutcome"`
	ReportReviewDate     string         `json:"reportReviewDate"`
	VisitReportMarkedAs  string         `json:"assuranceVisitReportMarkedAs"`
	VisitorAllocated     string         `json:"visitorAllocated"`
	ReviewedBy           string         `json:"reviewedBy"`
	Contact              ContactPerson  `json:"deputyContact"`
	TaskType             string         `json:"taskType"`
	Assignee             string         `json:"assignee"`
	DueDate              string         `json:"dueDate"`
	Notes                string         `json:"description"`
	OldAssigneeName      string         `json:"oldAssigneeName"`
}

type Changes struct {
	FieldName string `json:"fieldName"`
	OldValue  string `json:"oldValue"`
	NewValue  string `json:"newValue"`
}

type ClientPerson struct {
	ID       string `json:"personId"`
	Uid      string `json:"personUid"`
	Name     string `json:"personName"`
	CourtRef string `json:"personCourtRef"`
}

type ContactPerson struct {
	Name             string `json:"name"`
	JobTitle         string `json:"jobTitle"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phoneNumber"`
	OtherPhoneNumber string `json:"otherPhoneNumber"`
	Notes            string `json:"notes"`
}

type DeputyEvent struct {
	ID         int    `json:"id"`
	Timestamp  string `json:"timestamp"`
	EventType  string `json:"eventType"`
	User       User   `json:"user"`
	Event      Event  `json:"event"`
	IsNewEvent bool
}

func (c *Client) GetDeputyEvents(ctx Context, deputyId int) (DeputyEvents, error) {
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
		event := DeputyEvent{
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

func isNewEvent(changes []Changes) bool {
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

func updateTaskInfo(event Event, taskTypes TaskTypeMap) Event {
	if event.TaskType > "" {
		event.TaskType = taskTypes[event.TaskType].Description
		event.DueDate = FormatDateTime(IsoDateTime, event.DueDate, SiriusDate)
	}
	return event
}
