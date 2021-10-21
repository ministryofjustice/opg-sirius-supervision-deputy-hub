package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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
		reformattedDate := ReformatOrderDate(t.OrderDate)

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

func ReformatOrderDate(orderDate string) time.Time {
	dashDateString := strings.Replace(orderDate, "/", "-", 2)
	reformattedDate := fmt.Sprintf("%s%s%s%s%s", dashDateString[6:], "-", dashDateString[3:5], "-", dashDateString[:2])
	date, _ := time.Parse("2006-01-02", reformattedDate)
	return date
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
	} else {
		return dueDate
	}
}

func ReportDueScoreSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
	fmt.Println("before")
	fmt.Println(clients)
	sort.Slice(clients, func(i, j int) bool {
		x := SetDueDateForSort(clients[i].OldestReport.DueDate, clients[i].OldestReport.RevisedDueDate)
		y := SetDueDateForSort(clients[j].OldestReport.DueDate, clients[j].OldestReport.RevisedDueDate)

		fmt.Println("x")
		fmt.Println(x)
		fmt.Println("y")
		fmt.Println(y)

		clientiReportDueDateArray := strings.Split(x, "/")
		fmt.Println("client i report revised array")
		fmt.Println(clientiReportDueDateArray)
		fmt.Println("check array value")
		fmt.Println(reflect.ValueOf(clientiReportDueDateArray).Kind())
		fmt.Println("element 0")
		fmt.Println(clientiReportDueDateArray[0])
		fmt.Println("element 1")
		fmt.Println(clientiReportDueDateArray[1])
		fmt.Println("element 2")
		fmt.Println(clientiReportDueDateArray[2])
		fmt.Println("split the array")
		fmt.Println(clientiReportDueDateArray)
		clientiRestructuredDueDate := clientiReportDueDateArray[2] + "-" + clientiReportDueDateArray[1] + "-" + clientiReportDueDateArray[0]
		fmt.Println("did i get past the using split bit")
		clientjReportDueDate := clients[j].OldestReport.DueDate
		fmt.Println("client j report revised due date")
		fmt.Println(clientjReportDueDate)
		clientjReportDueDateArray := strings.Split(y, "/")
		clientjRestructuredDueDate := clientjReportDueDateArray[2] + "-" + clientjReportDueDateArray[1] + "-" + clientjReportDueDateArray[0]
		iDueDateTime, _ := time.Parse("2006-01-02", clientiRestructuredDueDate)
		jDueDateTime, _ := time.Parse("2006-01-02", clientjRestructuredDueDate)
		fmt.Println("iDueDateTime")
		fmt.Println(iDueDateTime)
		fmt.Println("jDueDateTime")
		fmt.Println(jDueDateTime)
		if sortOrder == "asc" {
			return iDueDateTime.Before(jDueDateTime)
		} else {
			return jDueDateTime.Before(iDueDateTime)
		}
	})
	fmt.Println("after")
	fmt.Println(clients)
	return clients
}

// func formatDate(unformattedDate string) string {
// 	// log date to see if null
// 	fmt.Println("unformatted date at start of format date")
// 	fmt.Println(unformattedDate)
// 	unformattedDateArray := strings.Split(unformattedDate, "/")
// 	fmt.Println("Unformatted date array is")
// 	fmt.Println(unformattedDateArray)
// 	var formattedDate string
// 	fmt.Println("Unformatted date array length is")
// 	fmt.Println(len(unformattedDateArray))
// 	if len(unformattedDateArray) > 2 {
// 		formattedDate = unformattedDateArray[2] + "-" + unformattedDateArray[1] + "-" + unformattedDateArray[0]
// 		// formattedDate, _ = time.Parse("2006-01-02", restructuredDate)
// 		fmt.Println("formattedDate is")
// 		fmt.Println(formattedDate)
// 		return formattedDate
// 	} else {
// 		return formattedDate
// 	}

// 	return formattedDate

// }

// func compareDates(sortOrder string, iDueDate, jDueDate time.Time) bool {
// 	if sortOrder == "asc" {
// 		return iDueDate.Before(jDueDate)
// 	} else {
// 		return jDueDate.Before(iDueDate)
// 	}
// }

// func ReportDueScoreSort(clients DeputyClientDetails, sortOrder string) DeputyClientDetails {
// 	sort.Slice(clients, func(i, j int) bool {
// 		fmt.Println("client is")
// 		fmt.Println(clients[i].Firstname)
// 		fmt.Println("due date is ")
// 		fmt.Println(clients[i].OldestReport.DueDate)
// 		fmt.Println("revised due date is ")
// 		fmt.Println(clients[i].OldestReport.RevisedDueDate)
// 		fmt.Println("type of due date")
// 		fmt.Printf("%T", clients[i].OldestReport.DueDate)
// 		if clients[i].OldestReport.RevisedDueDate != "null" {
// 			clientiDueDate := formatDate(clients[i].OldestReport.RevisedDueDate)
// 			fmt.Println("client i due date after format date")
// 			fmt.Println(clientiDueDate)
// 			clientiDueDateTime, _ := time.Parse("2006-01-02", clientiDueDate)
// 			fmt.Println("client i due date time within report due score sort")
// 			fmt.Println(clientiDueDateTime)
// 			if clients[j].OldestReport.RevisedDueDate != "null" {
// 				clientjDueDate := formatDate(clients[j].OldestReport.RevisedDueDate)
// 				fmt.Println("client j due date after format date")
// 				fmt.Println(clientjDueDate)
// 				clientjDueDateTime, _ := time.Parse("2006-01-02", clientjDueDate)
// 				fmt.Println("client i due date time within report due score sort")
// 				fmt.Println(clientjDueDateTime)
// 				return compareDates(sortOrder, clientiDueDateTime, clientjDueDateTime)
// 			} else {
// 				clientjDueDate := formatDate(clients[j].OldestReport.DueDate)
// 				fmt.Println("client j due date after format date")
// 				fmt.Println(clientjDueDate)
// 				clientjDueDateTime, _ := time.Parse("2006-01-02", clientjDueDate)
// 				fmt.Println("client j due date time within report due score sort")
// 				fmt.Println(clientjDueDateTime)
// 				return compareDates(sortOrder, clientiDueDateTime, clientjDueDateTime)
// 			}
// 		} else {
// 			clientiDueDate := formatDate(clients[i].OldestReport.DueDate)
// 			fmt.Println("client i due date after format date")
// 			fmt.Println(clientiDueDate)
// 			clientiDueDateTime, _ := time.Parse("2006-01-02", clientiDueDate)
// 			fmt.Println("client i due date time within report due score sort")
// 			fmt.Println(clientiDueDateTime)
// 			if clients[j].OldestReport.RevisedDueDate != "null" {

// 				clientjDueDate := formatDate(clients[j].OldestReport.RevisedDueDate)
// 				fmt.Println("client j due date after format date")
// 				fmt.Println(clientjDueDate)
// 				clientjDueDateTime, _ := time.Parse("2006-01-02", clientjDueDate)
// 				fmt.Println("client j due date within report due score sort")
// 				fmt.Println(clientjDueDate)
// 				return compareDates(sortOrder, clientiDueDateTime, clientjDueDateTime)
// 			} else {
// 				clientjDueDate := formatDate(clients[j].OldestReport.DueDate)
// 				fmt.Println("client j due date after format date")
// 				fmt.Println(clientjDueDate)
// 				clientjDueDateTime, _ := time.Parse("2006-01-02", clientjDueDate)
// 				fmt.Println("client j due date time within report due score sort")
// 				fmt.Println(clientjDueDateTime)
// 				return compareDates(sortOrder, clientiDueDateTime, clientjDueDateTime)
// 			}
// 		}
// 	})
// 	return clients
// }

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
