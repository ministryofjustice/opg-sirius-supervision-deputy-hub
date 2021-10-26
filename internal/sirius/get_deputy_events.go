package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type DeputyEventCollection []DeputyEvent

type User struct {
	UserId          int    `json:"id"`
	UserDisplayName string `json:"displayName"`
	UserPhoneNumber string `json:"phoneNumber"`
}

type Event struct {
	OrderType        string         `json:"orderType"`
	SiriusId         string         `json:"orderUid"`
	OrderNumber      string         `json:"orderCourtRef"`
	DeputyID         string         `json:"personId"`
	DeputyName       string         `json:"personName"`
	OrganisationName string         `json:"organisationName"`
	Changes          []Changes      `json:"changes"`
	Client           []ClientPerson `json:"additionalPersons"`
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

	DeputyEvents := EditDeputyEvents(v)

	return DeputyEvents, err

}

func EditDeputyEvents(v DeputyEventCollection) DeputyEventCollection {
	var list DeputyEventCollection
	for _, s := range v {
		event := DeputyEvent{
			Timestamp:       s.Timestamp,
			EventType:       ReformatEventType(s.EventType),
			TimelineEventId: s.TimelineEventId,
			User:            s.User,
			Event:           s.Event,
		}

		list = append(list, event)
	}
	SortTimeLineNewestOneFirst(list)
	return list
}

func ReformatEventType(s string) string {
	eventTypeArray := strings.Split(s, "\\")
	eventTypeArrayLength := len(eventTypeArray)
	eventTypeName := eventTypeArray[eventTypeArrayLength-1]
	return eventTypeName
}

func SortTimeLineNewestOneFirst(v DeputyEventCollection) DeputyEventCollection {
	sort.Slice(v, func(i, j int) bool {
		return v[i].Timestamp > v[j].Timestamp
	})
	return v
}
