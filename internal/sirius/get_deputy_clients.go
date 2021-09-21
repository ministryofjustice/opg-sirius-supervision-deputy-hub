package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type apiOrder struct {
	OrderStatus struct {
		Label string `json:"label"`
	}
	LatestSupervisionLevel *struct {
		SupervisionLevel struct {
			Label string `json:"label"`
		}
	}
	OrderDate string `json:"orderDate"`
}

type apiOrders []apiOrder

type apiClients struct {
	Clients []struct {
		ClientId            int    `json:"id"`
		Firstname           string `json:"firstname"`
		Surname             string `json:"surname"`
		CourtRef            string `json:"caseRecNumber"`
		RiskScore           int    `json:"riskScore"`
		ClientAccommodation struct {
			Label string `json:"label"`
		}
		Orders apiOrders
	} `json:"persons"`
}

type Order struct {
	OrderStatus      string
	SupervisionLevel string
	OrderDate        time.Time
}

type Orders []Order

type DeputyClient struct {
	ClientId          int
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

	var clients DeputyClientDetails
	for _, t := range v.Clients {
		orders := RestructureOrders(t.Orders)
		if len(orders) > 0 {
			var client = DeputyClient{
				ClientId:          t.ClientId,
				Firstname:         t.Firstname,
				Surname:           t.Surname,
				CourtRef:          t.CourtRef,
				RiskScore:         t.RiskScore,
				AccommodationType: t.ClientAccommodation.Label,
				OrderStatus:       GetOrderStatus(orders),
				SupervisionLevel:  GetMostRecentSupervisionLevel(orders),
			}
			clients = append(clients, client)
		}
	}
	return clients, err
}

func GetOrderStatus(orders Orders) string {
	// return the status of the oldest active order for a client
	// if there isn’t one, the status of the oldest order

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].OrderDate.Before(orders[j].OrderDate)
	})

	for _, o := range orders {
		if o.OrderStatus == "Active" {
			return o.OrderStatus
		}
	}
	return orders[0].OrderStatus
}

func GetMostRecentSupervisionLevel(orders Orders) string {
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].OrderDate.After(orders[j].OrderDate)
	})
	return orders[0].SupervisionLevel
}

func RestructureOrders(apiOrders apiOrders) Orders {
	orders := make(Orders, len(apiOrders))

	for i, t := range apiOrders {
		// reformatting order date to yyyy-dd-mm
		reformattedDate := reformatOrderDate(t.OrderDate)

		var supervisionLevel string
		if t.LatestSupervisionLevel != nil {
			supervisionLevel = t.LatestSupervisionLevel.SupervisionLevel.Label
		} else {
			supervisionLevel = ""
		}

		orders[i] = Order{
			OrderStatus:      t.OrderStatus.Label,
			SupervisionLevel: supervisionLevel,
			OrderDate:        reformattedDate,
		}
	}

	updatedOrders := removeOpenStatusOrders(orders)
	return updatedOrders
}

func reformatOrderDate(orderDate string) time.Time {
	dashDateString := strings.Replace(orderDate, "/", "-", 2)
	reformattedDate := fmt.Sprintf("%s%s%s%s%s", dashDateString[6:], "-", dashDateString[3:5], "-", dashDateString[:2])
	date, _ := time.Parse("2006-01-02", reformattedDate)
	return date
}

func removeOpenStatusOrders(orders Orders) Orders {
	// an order is open when it's with the Allocations team,
	// and so not yet supervised by the PA team

	var updatedOrders Orders
	for _, o := range orders {
		if o.OrderStatus != "Open" {
			updatedOrders = append(updatedOrders, o)
		}
	}
	return updatedOrders
}
