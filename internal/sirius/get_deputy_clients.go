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
				OldestReport:      reportReturned{t.OldestReport.DueDate, t.OldestReport.RevisedDueDate, t.OldestReport.Status.Label},
			}

			clients = append(clients, client)
		}
	}

	var aria AriaSorting
	aria.SurnameAriaSort = ChangeSortButtonDirection(sortOrder, columnBeingSorted, "sort=surname")
	aria.ReportDueAriaSort = ChangeSortButtonDirection(sortOrder, columnBeingSorted, "sort=reportdue")
	aria.CRECAriaSort = ChangeSortButtonDirection(sortOrder, columnBeingSorted, "sort=crec")

	switch columnBeingSorted {
	case "sort=reportdue":
		ReportDueScoreSort(clients, sortOrder)
	case "sort=crec":
		CrecScoreSort(clients, sortOrder)
	default:
		AlphabeticalSort(clients, sortOrder)
	}

	return clients, aria, err
}

/*
	GetOrderStatus returns the status of the oldest active order for a client.
  If there isn’t one, the status of the oldest order is returned.
*/
func GetOrderStatus(orders Orders) string {
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
		reformattedDate := FormattingDate(t.OrderDate)

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

	updatedOrders := RemoveOpenStatusOrders(orders)
	return updatedOrders
}

func FormattingDate(dateString string) time.Time {
	reformattedDate := strings.ReplaceAll(dateString, "/", "-")
	dateTime, _ := time.Parse("02-01-2006", reformattedDate)
	return dateTime
}

func RemoveOpenStatusOrders(orders Orders) Orders {
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

func AlphabeticalSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
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

func CrecScoreSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
	sort.Slice(clients, func(i, j int) bool {
		if sortOrder == "asc" {
			return clients[i].RiskScore < clients[j].RiskScore
		} else {
			return clients[i].RiskScore > clients[j].RiskScore
		}
	})
	return clients
}

func SetDueDateForSort(dueDate, revisedDueDate string) string {

	if revisedDueDate != "" {
		return revisedDueDate
	} else if dueDate != "" {
		return dueDate
	} else {
		return "12/12/9999"
	}
}

func ReportDueScoreSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
	sort.Slice(clients, func(i, j int) bool {
		x := SetDueDateForSort(clients[i].OldestReport.DueDate, clients[i].OldestReport.RevisedDueDate)
		y := SetDueDateForSort(clients[j].OldestReport.DueDate, clients[j].OldestReport.RevisedDueDate)
		dateTimeI := FormattingDate(x)
		dateTimeJ := FormattingDate(y)

		if sortOrder == "asc" {
			return dateTimeI.Before(dateTimeJ)
		} else {
			return dateTimeJ.Before(dateTimeI)
		}
	})
	return clients
}

func ChangeSortButtonDirection(sortOrder string, columnBeingSorted string, functionCalling string) string {
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
