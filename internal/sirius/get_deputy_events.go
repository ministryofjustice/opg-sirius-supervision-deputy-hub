package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeputyEvents []DeputyEvent

type DeputyEvent struct {
	TaskId int `json:"id"`
	Timestamp string `json:"timestamp"`
	EventType string `json:"eventType"`
	DeputyID               int    `json:"personId"`
	DeputyName   string `json:"personName"`
	OrganisationName string `json:"organisationName"`
	User struct {
		UserId int `json:"id"`
		UserDisplayName string `json:"displayName"`
	} `json:"user"`
	Event struct {
		OrderType string `json:"orderType"`
		SiriusId string `json:"orderUid"`
		OrderNumber string `json:"orderId"`
		Changes []struct {
			FieldName string `json:"fieldName"`
			OldValue string `json:"oldValue"`
			NewValue string `json:"newValue"`
		}`json:"changes"`
		Client []struct {
			ClientName string `json:"personName"`
			ClientId string `json:"personId"`
			ClientUid string `json:"personUid"`
			ClientCourtRef string `json:"personCourtRef"`
		}`json:"additionalPersons"`
	} `json:"event"`
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
			TaskId:           s.TaskId,
			DeputyID:         s.DeputyID,
			DeputyName:       s.DeputyName,
			OrganisationName: s.OrganisationName,
			User:             s.User,
			Event: s.Event,
		}

		list = append(list, event)
	}
		return list
}

func ReformatTimestamp(DeputyEvent) string {
	return "2020-01-01"
}

func ReformatEventType(DeputyEvent) string {
	return "good event"
}