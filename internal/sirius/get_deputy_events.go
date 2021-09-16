package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DeputyEvents []DeputyEvent


type User struct {
	UserId int `json:"id"`
	UserDisplayName string `json:"displayName"`
	UserPhoneNumber string `json:"phoneNumber"`
}

type Event struct {
	OrderType        string         `json:"orderType"`
	SiriusId         string         `json:"orderUid"`
	OrderNumber      string         `json:"orderId"`
	DeputyID         string            `json:"personId"`
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
	TimelineEventId  int    `json:"id"`
	Timestamp        string `json:"timestamp"`
	EventType        string `json:"eventType"`
	User             User   `json:"user"`
	Event            Event  `json:"event"`
}

func (c *Client) GetDeputyEvents(ctx Context, deputyId int) (DeputyEvents, error) {
	var v DeputyEvents

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

	DeputyEvents := v
	DeputyEvents = EditDeputyEvents(v)

	return DeputyEvents, err

}

func EditDeputyEvents(v DeputyEvents) DeputyEvents {
	var list DeputyEvents
	for _, s := range v {
		event := DeputyEvent{
			Timestamp:        ReformatTimestamp(s),
			EventType:        ReformatEventType(s),
			TimelineEventId:  s.TimelineEventId,
			User:             s.User,
			Event: s.Event,
		}

		list = append(list, event)
	}
		return list
}

func ReformatTimestamp(s DeputyEvent) string {
	//edit time to format required by timeline
	return s.Timestamp
}

func ReformatEventType(s DeputyEvent) string {
	stringsArray := strings.Split(s.EventType, "\\")
	string := stringsArray[5]
	//add spaces into the name
	return string
}