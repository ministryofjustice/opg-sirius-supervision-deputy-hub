package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type DeputyEventCollection []DeputyEvent

type User struct {
	UserId          int    `json:"id"`
	UserDisplayName string `json:"displayName"`
	UserPhoneNumber string `json:"phoneNumber"`
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
}

type Changes struct {
	FieldName string `json:"fieldName"`
	OldValue  string `json:"oldValue"`
	NewValue  string `json:"newValue"`
}

type ClientPerson struct {
	ClientName     string `json:"personName"`
	ClientId       string `json:"personId"`
	ClientUid      string `json:"personUid"`
	ClientCourtRef string `json:"personCourtRef"`
}

type DeputyEvent struct {
	TimelineEventId int    `json:"id"`
	Timestamp       string `json:"timestamp"`
	EventType       string `json:"eventType"`
	User            User   `json:"user"`
	Event           Event  `json:"event"`
	IsNewEvent      bool
}

func (c *Client) GetDeputyEvents(ctx Context, deputyId int) (DeputyEventCollection, error) {
	var v DeputyEventCollection

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/timeline/%d", deputyId), nil)

	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}
	err = json.NewDecoder(resp.Body).Decode(&v)

	DeputyEvents := editDeputyEvents(v)

	return DeputyEvents, err

}

func editDeputyEvents(v DeputyEventCollection) DeputyEventCollection {
	var list DeputyEventCollection
	for _, s := range v {
		event := DeputyEvent{
			Timestamp:       FormatDateAndTime("2006-01-02 15:04:05", s.Timestamp, "02/01/2006 15:04:05"),
			EventType:       reformatEventType(s.EventType),
			TimelineEventId: s.TimelineEventId,
			User:            s.User,
			Event:           s.Event,
			IsNewEvent:      calculateIfNewEvent(s.Event.Changes),
		}

		list = append(list, event)
	}
	sortTimeLineNewestOneFirst(list)
	return list
}

func calculateIfNewEvent(changes []Changes) bool {
	var answer bool
	for _, s := range changes {
		if s.OldValue != "" {
			answer = false
		} else {
			answer = true
		}
	}
	return answer
}

func reformatEventType(s string) string {
	eventTypeArray := strings.Split(s, "\\")
	eventTypeArrayLength := len(eventTypeArray)
	eventTypeName := eventTypeArray[eventTypeArrayLength-1]
	return eventTypeName
}

func sortTimeLineNewestOneFirst(v DeputyEventCollection) DeputyEventCollection {
	sort.Slice(v, func(i, j int) bool {
		changeToTimeTypeI, _ := time.Parse("02/01/2006 15:04:05", v[i].Timestamp)
		changeToTimeTypeJ, _ := time.Parse("02/01/2006 15:04:05", v[j].Timestamp)
		return changeToTimeTypeJ.Before(changeToTimeTypeI)
	})
	return v
}
