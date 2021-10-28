package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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

type apiReport struct {
	DueDate        string `json:"dueDate"`
	RevisedDueDate string `json:"revisedDueDate"`
	Status         struct {
		Label string `json:"label"`
	} `json:"status"`
}

type reportReturned struct {
	DueDate        string
	RevisedDueDate string
	StatusLabel    string
}

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
		Orders       apiOrders `json:"orders"`
		OldestReport apiReport `json:"oldestNonLodgedAnnualReport"`
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
	OldestReport      reportReturned
}

type DeputyClientDetails []DeputyClient

type AriaSorting struct {
	SurnameAriaSort   string
	ReportDueAriaSort string
	CRECAriaSort      string
}

func (c *Client) GetDeputyClients(ctx Context, deputyId int, columnBeingSorted string, sortOrder string) (DeputyClientDetails, AriaSorting, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/pa/%d/clients", deputyId), nil)
	if err != nil {
		return nil, AriaSorting{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, AriaSorting{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, AriaSorting{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, AriaSorting{}, newStatusError(resp)
	}

	var v apiClients
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, AriaSorting{}, err
	}

	var clients DeputyClientDetails
	for _, t := range v.Clients {
		orders := restructureOrders(t.Orders)
		if len(orders) > 0 {
			var client = DeputyClient{
				ClientId:          t.ClientId,
				Firstname:         t.Firstname,
				Surname:           t.Surname,
				CourtRef:          t.CourtRef,
				RiskScore:         t.RiskScore,
				AccommodationType: t.ClientAccommodation.Label,
				OrderStatus:       getOrderStatus(orders),
				SupervisionLevel:  getMostRecentSupervisionLevel(orders),
				OldestReport:      reportReturned{t.OldestReport.DueDate, t.OldestReport.RevisedDueDate, t.OldestReport.Status.Label},
			}

			clients = append(clients, client)
		}
	}

	var aria AriaSorting
	aria.SurnameAriaSort = changeSortButtonDirection(sortOrder, columnBeingSorted, "sort=surname")
	aria.ReportDueAriaSort = changeSortButtonDirection(sortOrder, columnBeingSorted, "sort=reportdue")
	aria.CRECAriaSort = changeSortButtonDirection(sortOrder, columnBeingSorted, "sort=crec")

	switch columnBeingSorted {
	case "sort=reportdue":
		reportDueScoreSort(clients, sortOrder)
	case "sort=crec":
		crecScoreSort(clients, sortOrder)
	default:
		alphabeticalSort(clients, sortOrder)
	}

	return clients, aria, err
}

/*
	GetOrderStatus returns the status of the oldest active order for a client.
  If there isn’t one, the status of the oldest order is returned.
*/
func getOrderStatus(orders Orders) string {
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

func getMostRecentSupervisionLevel(orders Orders) string {
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].OrderDate.After(orders[j].OrderDate)
	})
	return orders[0].SupervisionLevel
}

func restructureOrders(apiOrders apiOrders) Orders {
	orders := make(Orders, len(apiOrders))

	for i, t := range apiOrders {
		// reformatting order date to yyyy-dd-mm
		reformattedDate := formatDate(t.OrderDate)

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

func formatDate(dateString string) time.Time {
	dateTime, _ := time.Parse("02/01/2006", dateString)
	return dateTime
}

func removeOpenStatusOrders(orders Orders) Orders {
	/* An order is open when it's with the Allocations team,
	and so not yet supervised by the PA team */

	var updatedOrders Orders
	for _, o := range orders {
		if o.OrderStatus != "Open" {
			updatedOrders = append(updatedOrders, o)
		}
	}
	return updatedOrders
}

func alphabeticalSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
	if len(clients) > 1 {
		sort.Slice(clients, func(i, j int) bool {
			if sortOrder == "asc" {
				return clients[i].Surname < clients[j].Surname
			} else {
				return clients[i].Surname > clients[j].Surname
			}
		})
	}
	return clients
}

func crecScoreSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
	sort.Slice(clients, func(i, j int) bool {
		if sortOrder == "asc" {
			return clients[i].RiskScore < clients[j].RiskScore
		} else {
			return clients[i].RiskScore > clients[j].RiskScore
		}
	})
	return clients
}

func setDueDateForSort(dueDate, revisedDueDate string) string {
	if revisedDueDate != "" {
		return revisedDueDate
	} else if dueDate != "" {
		return dueDate
	} else {
		return "12/12/9999"
	}
}

func reportDueScoreSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
	sort.Slice(clients, func(i, j int) bool {
		x := setDueDateForSort(clients[i].OldestReport.DueDate, clients[i].OldestReport.RevisedDueDate)
		y := setDueDateForSort(clients[j].OldestReport.DueDate, clients[j].OldestReport.RevisedDueDate)
		dateTimeI := formatDate(x)
		dateTimeJ := formatDate(y)

		if sortOrder == "asc" {
			return dateTimeI.Before(dateTimeJ)
		} else {
			return dateTimeJ.Before(dateTimeI)
		}
	})
	return clients
}

func changeSortButtonDirection(sortOrder string, columnBeingSorted string, functionCalling string) string {
	if functionCalling == columnBeingSorted {
		if sortOrder == "asc" {
			return "ascending"
		} else if sortOrder == "desc" {
			return "descending"
		}
		return "none"
	} else {
		return "none"
	}

}
