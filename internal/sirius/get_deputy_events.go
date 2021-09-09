package sirius

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DeputyEvents struct {
	DeputyID               int    `json:"personId"`
	DeputyName   string `json:"personName"`
	//OrganisationName string `json:"organisationName"`
	//User struct {
	//	UserId int `json:"id"`
	//	UserDisplayName int `json:"displayName"`
	//} `json:"user"`
	//Event struct {
	//	OrderNumber int `json:"orderId"`
	//	SiriusId int `json:"orderUid"`
	//	OrderType int `json:"orderType"`
	//	Client struct {
	//		ClientName string `json:"personName"`
	//	} `json:"additionalPersons"`
	//} `json:"event"`
}

func (c *Client) GetDeputyEvents(ctx Context, deputyId int) (DeputyEvents, error) {
	var v DeputyEvents

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/timeline/76"), nil)

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
	b, err := io.ReadAll(resp.Body)
	fmt.Println(string(b))
	return v, err

}