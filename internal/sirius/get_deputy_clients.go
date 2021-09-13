package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)


type DeputyClientDetails struct {
	Clients     []struct {
		ID          int    `json:"id"`
		Firstname string `json:"firstname"`
		Surname	string `json:"surname"`
	} `json:"persons"`
}

func (c *Client) GetDeputyClients(ctx Context, deputyId int) (DeputyClientDetails, error) {
	var v DeputyClientDetails

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/clients", deputyId), nil)

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

	clients := DeputyClientDetails{
		Clients: v.Clients,
	}

	return clients, err

}