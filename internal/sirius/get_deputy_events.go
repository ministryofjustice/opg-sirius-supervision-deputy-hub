package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeputyEvents []DeputyEvent

type DeputyEvent struct {
	TaskId int `json:"id"`
	//DeputyID               int    `json:"personId"`
	//DeputyName   string `json:"personName"`
	//OrganisationName string `json:"organisationName"`
	//User struct {
	//	UserId int `json:"id"`
	//	UserDisplayName int `json:"displayName"`
	//} `json:"user"`
	Event struct {
		OrderNumber string `json:"orderId"`
		SiriusId string `json:"orderUid"`
		OrderType string `json:"orderType"`
		//Client struct {
		//	ClientName string `json:"personName"`
		//} `json:"additionalPersons"`
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
	return v, err

}