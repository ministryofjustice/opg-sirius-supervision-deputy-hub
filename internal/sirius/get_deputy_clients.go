package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Order struct {
	OrderStatus struct {
		Label string `json:"label"`
	}
	LatestSupervisionLevel struct {
		SupervisionLevel struct {
			Label string `json:"label"`
		}
	}
	OrderDate string `json:"orderDate"`
}

type Orders []Order

type apiClients struct {
	Clients []struct {
		ID                  int    `json:"id"`
		Firstname           string `json:"firstname"`
		Surname             string `json:"surname"`
		CourtRef            string `json:"caseRecNumber"`
		RiskScore           int    `json:"riskScore"`
		ClientAccommodation struct {
			Label string `json:"label"`
		}
		Orders Orders
	} `json:"persons"`
}

type DeputyClient struct {
	ID                int
	Firstname         string
	Surname           string
	CourtRef          string
	RiskScore         int
	AccommodationType string
	OrderStatus       string
	SupervisionLevel  string
}

type DeputyClientDetails []DeputyClient


func (c *Client) GetDeputyClients(ctx Context, deputyId int) (DeputyClientDetails, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/clients", deputyId), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v apiClients
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	clients := make(DeputyClientDetails, len(v.Clients))

	for i, t := range v.Clients {
		clients[i] = DeputyClient{
			ID:                t.ID,
			Firstname:         t.Firstname,
			Surname:           t.Surname,
			CourtRef:          t.CourtRef,
			RiskScore:         t.RiskScore,
			AccommodationType: t.ClientAccommodation.Label,
			OrderStatus:       CalculateOrderStatus(t.Orders),
			SupervisionLevel:  GetLatestSupervisionLevel(t.Orders),
		}
	}

	return clients, err

}

func CalculateOrderStatus(o Orders) string {
	// return the status of the oldest active order for a client
	// if there isn’t one, the status of the oldest order

	return o[0].OrderStatus.Label
}

func GetLatestSupervisionLevel(o Orders) string {

	return o[0].LatestSupervisionLevel.SupervisionLevel.Label
}